package main

import (
	"context"
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
	"github.com/peterbourgon/ff/v3/ffcli"
	"moul.io/godev"
)

const manfredChannel = "691928992495960124"

func main() {
	var (
		rootFlagSet    = flag.NewFlagSet("chaos-bot", flag.ExitOnError)
		devMode        = rootFlagSet.Bool("dev", false, "only reply in PV")
		debug          = rootFlagSet.Bool("debug", false, "verbose")
		discordFlagSet = flag.NewFlagSet("discord", flag.ExitOnError)
		discordToken   = discordFlagSet.String("discord-token", "", "Discord Bot Token")
	)

	discordBot := &ffcli.Command{
		Name:    "discord-bot",
		FlagSet: discordFlagSet,
		Exec: func(_ context.Context, _ []string) error {
			return discordBot(*discordToken, *devMode, *debug)
		},
	}

	root := &ffcli.Command{
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{discordBot},
	}

	err := root.ParseAndRun(context.Background(), os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func discordBot(token string, devMode bool, debug bool) error {
	log.Printf("starting bot, devMode=%v, debug=%v", devMode, debug)
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return err
	}

	hostname, _ := os.Hostname()
	dg.ChannelMessageSend(manfredChannel, fmt.Sprintf("COUCOU JE VIENS DE BOOT (%s)", hostname))

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
