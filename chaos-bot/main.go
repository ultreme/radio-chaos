package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

func main() {
	var (
		discordbotOpts = discordbotOpts{}
		rootFlagSet    = flag.NewFlagSet("chaos-bot", flag.ExitOnError)
		devMode        = rootFlagSet.Bool("dev", false, "only reply in PV")
		debug          = rootFlagSet.Bool("debug", false, "verbose")
		discordFlagSet = flag.NewFlagSet("discord", flag.ExitOnError)
		mixerFlagSet   = flag.NewFlagSet("mixer-bot", flag.ExitOnError)
	)

	discordFlagSet.StringVar(&discordbotOpts.discordToken, "discord-token", "", "Discord Bot Token")
	discordFlagSet.StringVar(&discordbotOpts.serverBind, "server-bind", ":4747", "API listening address")
	discordFlagSet.StringVar(&discordbotOpts.manfredChannel, "manfred-channel", "691928992495960124", "manfred private channel (for debugging)")

	discordBot := &ffcli.Command{
		Name:    "discord-bot",
		FlagSet: discordFlagSet,
		Exec: func(_ context.Context, _ []string) error {
			discordbotOpts.debug = *debug
			discordbotOpts.devMode = *devMode
			return discordBotCmd(discordbotOpts)
		},
	}

	mixerBot := &ffcli.Command{
		Name:    "mixer-bot",
		FlagSet: mixerFlagSet,
		Exec: func(_ context.Context, _ []string) error {
			return mixerBotCmd()
		},
	}

	root := &ffcli.Command{
		FlagSet:     rootFlagSet,
		Subcommands: []*ffcli.Command{discordBot, mixerBot},
		Exec: func(context.Context, []string) error {
			return flag.ErrHelp
		},
	}

	err := root.ParseAndRun(context.Background(), os.Args[1:])
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}
