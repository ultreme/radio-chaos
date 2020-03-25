package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
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
		commands := []string{
			"!history",
			"!manfred",
			"!il-est-pas-quelle-heure",
			"!discord",
			"!radio",
			"!zoom",
			"!coucou",
			"!pouet",
			"!podcast",
			"!calendrier",
			"!ultreme",
			"!soundcloud",
		}
		sort.Strings(commands)
		out := strings.Join(commands, ", ")
		s.ChannelMessageSend(m.ChannelID, out)
	case m.Content == "!soundcloud":
		s.ChannelMessageSend(m.ChannelID, "https://soundcloud.com/ultreme-reporters")
	case m.Content == "!ultreme":
		s.ChannelMessageSend(m.ChannelID, "https://ultre.me")
	case m.Content == "!calendrier":
		s.ChannelMessageSend(m.ChannelID, "https://calendrier.ultre.me")
	case m.Content == "!pouet":
		s.ChannelMessageSend(m.ChannelID, "https://calendrier.ultre.me/2019/pouet/")
	case m.Content == "!manfred":
		msgs := []string{
			"c'est ce qu'elles disent toutes",
			"plus tu cours moins vite, moins t'avances plus vite",
		}
		msg := msgs[rand.Intn(len(msgs))]
		s.ChannelMessageSend(m.ChannelID, msg)
	case m.Content == "!il-est-pas-quelle-heure":
		out := fmt.Sprintf("%d%d:%d%d",
			rand.Intn(3),
			rand.Intn(10),
			rand.Intn(6),
			rand.Intn(10),
		)
		s.ChannelMessageSend(m.ChannelID, out)
	case m.Content == "!discord":
		s.ChannelMessageSend(m.ChannelID, "https://ultre.me/discord")
	case m.Content == "!radio":
		s.ChannelMessageSend(m.ChannelID, "http://salutcestcool.com/radio")
	case m.Content == "!zoom":
		s.ChannelMessageSend(m.ChannelID, `
Sur internet/via une appli: https://zoom.us/j/129255108
Depuis un téléphone: 01.70.37.22.46, puis taper 129 255 108#
`)
	case m.Content == "!coucou":
		msgs := []string{
			"SALUT ÇA VA !?",
			"bonjour à toutes et à tous",
			"bonjour",
			"enchanté de vous rencontrer",
			"coucou",
			"salut",
			"yo",
		}
		msg := msgs[rand.Intn(len(msgs))]
		s.ChannelMessageSend(m.ChannelID, msg)
	}
}
