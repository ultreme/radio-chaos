package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
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

		//fmt.Println(godev.PrettyJSON(parsed))
		msg := fmt.Sprintf("current: %s\n", parsed.Current.Pretty())
		msg += "\nhistory:\n"
		for _, history := range parsed.History {
			pretty := history.Pretty()
			if strings.TrimSpace(pretty) == "-" {
				continue
			}
			msg += fmt.Sprintf("  - %s\n", pretty)
		}
		s.ChannelMessageSend(m.ChannelID, msg)
	case m.Content == "!help":
		commands := []string{"!history", "!manfred", "!il-est-pas-quelle-heure", "!discord", "!radio", "!zoom", "!coucou", "!podcast", "!calendrier", "!ultreme", "!soundcloud"}
		sort.Strings(commands)
		out := ""
		for _, command := range commands {
			out += fmt.Sprintf("- %s\n", command)
		}
		s.ChannelMessageSend(m.ChannelID, out)
	case m.Content == "!soundcloud":
		s.ChannelMessageSend(m.ChannelID, "https://soundcloud.com/ultreme-reporters")
	case m.Content == "!ultreme":
		s.ChannelMessageSend(m.ChannelID, "https://ultre.me")
	case m.Content == "!calendrier":
		s.ChannelMessageSend(m.ChannelID, "https://calendrier.ultre.me")
	case m.Content == "!manfred":
		s.ChannelMessageSend(m.ChannelID, "c'est ce qu'elles disent toutes")
	case m.Content == "!il-est-pas-quelle-heure":
		s.ChannelMessageSend(m.ChannelID, "23:42")
	case m.Content == "!discord":
		s.ChannelMessageSend(m.ChannelID, "https://ultre.me/disord")
	case m.Content == "!radio":
		s.ChannelMessageSend(m.ChannelID, "http://salutcestcool.com/radio")
	case m.Content == "!zoom":
		s.ChannelMessageSend(m.ChannelID, `
Sur internet/via une appli: https://zoom.us/j/129255108
Depuis un téléphone: 01.70.37.22.46, puis taper 129 255 108#
`)
	case m.Content == "!coucou":
		s.ChannelMessageSend(m.ChannelID, "SALUT ÇA VA !?")
	}
}
