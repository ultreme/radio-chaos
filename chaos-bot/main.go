package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

const manfredChannel = "691928992495960124"

func main() {
	var (
		rootFlagSet    = flag.NewFlagSet("chaos-bot", flag.ExitOnError)
		devMode        = rootFlagSet.Bool("dev", false, "only reply in PV")
		debug          = rootFlagSet.Bool("debug", false, "verbose")
		discordFlagSet = flag.NewFlagSet("discord", flag.ExitOnError)
		discordToken   = discordFlagSet.String("discord-token", "", "Discord Bot Token")
		mixerFlagSet   = flag.NewFlagSet("mixer-bot", flag.ExitOnError)
	)

	discordBot := &ffcli.Command{
		Name:    "discord-bot",
		FlagSet: discordFlagSet,
		Exec: func(_ context.Context, _ []string) error {
			return discordBotCmd(*discordToken, *devMode, *debug)
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
