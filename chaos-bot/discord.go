package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/hako/durafmt"
	"github.com/oklog/run"
	"moul.io/godev"
)

func discordBotCmd(token string, devMode bool, debug bool) error {
	log.Printf("starting bot, devMode=%v, debug=%v", devMode, debug)

	var db *discordBot
	var dg *discordgo.Session
	{ // DISCORD
		var err error
		dg, err = discordgo.New("Bot " + token)
		if err != nil {
			return err
		}

		hostname, _ := os.Hostname()
		dg.ChannelMessageSend(manfredChannel, fmt.Sprintf("**COUCOU JE VIENS DE BOOT (%s)!**", hostname))

		db = &discordBot{
			startedAt: time.Now(),
			devMode:   devMode,
			debug:     debug,
		}
		commands["!info"] = db.doInfo
		commands["!radiosay"] = db.doRadioSay

		dg.AddHandler(db.onMessageCreate)

		if err = dg.Open(); err != nil {
			return err
		}
		defer dg.Close()
	}

	{ // SOCKET IO
		sio, err := newSIOClient("mixer")
		if err != nil {
			return err
		}
		defer sio.Close()

		sio.c.On("event:join", func(h *gosocketio.Channel, args Message) {
			dg.ChannelMessageSend(manfredChannel, fmt.Sprintf("sio event:join: %s", godev.JSON(args)))
		})
		sio.c.On("event:disconnect", func(h *gosocketio.Channel, args Message) {
			dg.ChannelMessageSend(manfredChannel, fmt.Sprintf("sio event:disconnect: %s", godev.JSON(args)))
		})
		sio.c.On("event:broadcast", func(h *gosocketio.Channel, args Message) {
			dg.ChannelMessageSend(manfredChannel, fmt.Sprintf("sio event:broadcast: %s", godev.JSON(args)))
		})
		db.sio = sio
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	var g run.Group
	g.Add(run.SignalHandler(context.TODO(), os.Interrupt))
	return g.Run()
}

type discordBot struct {
	startedAt    time.Time
	seenMessages int
	seenCommands int
	seenErrors   int
	devMode      bool
	debug        bool
	lock         sync.Mutex
	sio          *sioClient
}

func (b *discordBot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// avoid loop
	if m.Author.ID == s.State.User.ID {
		return
	}
	if b.debug {
		log.Println(godev.JSON(m))
	}
	if b.devMode && m.GuildID != "" {
		return
	}

	args := strings.Split(m.Content, " ")

	b.seenMessages++
	command, found := commands[args[0]] // FIXME: split the content and support args
	if !found || command == nil {
		return
	}
	b.seenCommands++

	b.lock.Lock()
	defer b.lock.Unlock()

	channel := m.ChannelID
	if m.GuildID == "" {
		channel = "PV"
	}
	log.Printf("channel: %q, author: %s#%s, input: %q", channel, m.Author.Username, m.Author.Discriminator, m.Content)

	defer func() {
		if r := recover(); r != nil {
			sendError(s, m, fmt.Errorf("panic: %v", r))
		}
	}()
	err := command(s, m)
	if err != nil {
		b.seenErrors++
		sendError(s, m, err)
	}
}

func (b *discordBot) doRadioSay(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Split(m.Content, " ")
	say := strings.Join(content[1:], " ")
	say = strings.TrimSpace(say)
	if say == "" {
		return nil
	}

	input := Message{
		Room: "radio-chaos",
		Msg: &Msg{
			Text: ProtocolMsg{
				Action: "radiosay",
				Args:   []string{say},
			}.ToJSON(),
		},
	}
	return b.sio.c.Emit("broadcast", input)
}

func (b *discordBot) doInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	uptime := time.Since(b.startedAt)
	msg := fmt.Sprintf(
		"uptime: %v, messages: %d, commands: %d, errors: %d\n",
		durafmt.ParseShort(uptime).String(), b.seenMessages, b.seenCommands, b.seenErrors,
	)
	msg += fmt.Sprintf("source: https://github.com/ultreme/radio-chaos/tree/master/chaos-bot")
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}
