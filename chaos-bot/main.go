package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/hako/durafmt"
	"moul.io/godev"
)

func main() {
	var (
		token   string
		devMode bool
		debug   bool
	)

	flag.StringVar(&token, "t", "", "Bot Token")
	flag.BoolVar(&devMode, "dev", false, "Only reply in PV")
	flag.BoolVar(&debug, "debug", false, "Verbose")
	flag.Parse()

	err := run(token, devMode, debug)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func run(token string, devMode bool, debug bool) error {
	log.Printf("starting bot, devMode=%v, debug=%v", devMode, debug)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	b := bot{
		startedAt: time.Now(),
		devMode:   devMode,
		debug:     debug,
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
	devMode      bool
	debug        bool
	lock         sync.Mutex
}

func (b *bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	b.seenMessages++
	command, found := commands[m.Content] // FIXME: split the content and support args
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
