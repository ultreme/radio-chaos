package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
)

func main() {
	var token string
	flag.StringVar(&token, "t", "", "Bot Token")
	// FIXME: --dev (only in pv)
	// FIXME: --debug
	flag.Parse()

	err := run(token)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(token string) error {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	b := bot{
		startedAt: time.Now(),
	}
	commands["!info"] = b.doInfo

	dg.AddHandler(b.onMessageCreate)

	if err = dg.Open(); err != nil {
		return err
	}
	defer dg.Close()

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	return nil
}

type bot struct {
	startedAt    time.Time
	seenMessages int
	seenCommands int
	seenErrors   int
}

func (b *bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// avoid loop
	if m.Author.ID == s.State.User.ID {
		return
	}

	b.seenMessages++
	command, found := commands[m.Content] // FIXME: split the content and support args
	if !found {
		return
	}
	b.seenCommands++

	log.Printf("channel: %q, author: %s#%s, input: %q", m.ChannelID, m.Author.Username, m.Author.Discriminator, m.Content)

	err := command(s, m)
	if err != nil {
		b.seenErrors++
		sendError(s, m, err)
	}
}

func (b *bot) doInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	uptime := time.Since(b.startedAt)
	msg := fmt.Sprintf(
		"uptime: %v, messages: %d, commands: %d, errors: %d\n",
		durafmt.ParseShort(uptime).String(), b.seenMessages, b.seenCommands, b.seenErrors,
	)
	msg += fmt.Sprintf("source: https://github.com/ultreme/radio-chaos/tree/master/chaos-bot")
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}
