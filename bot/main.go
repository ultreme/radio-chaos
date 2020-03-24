package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"moul.io/godev"
)

func main() {
	var token string
	flag.StringVar(&token, "t", "", "Bot Token")
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

	dg.AddHandler(messageCreate)

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

func sendError(s *discordgo.Session, m *discordgo.MessageCreate, err error) {
	log.Printf("ERROR: %v", err)
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("ERROR!"))
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	switch {
	case m.Content == "!history":
		resp, err := http.Get("http://radio-admin.casse-tete.solutions/?action=infos&format=json")
		if err != nil {
			sendError(s, m, err)
			break
		}
		var parsed GetInfosResponse

		err = json.NewDecoder(resp.Body).Decode(&parsed)
		if err != nil {
			sendError(s, m, err)
			break
		}

		fmt.Println(godev.PrettyJSON(parsed))

		msg := fmt.Sprintf("current: %s\n", parsed.Current.Pretty())
		msg += "\nhistory:\n"
		for _, history := range parsed.History {
			msg += fmt.Sprintf("  - %s\n", history.Pretty())
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	case m.Content == "!coucou":
		s.ChannelMessageSend(m.ChannelID, "SALUT ÇA VA !?")
	}
}
