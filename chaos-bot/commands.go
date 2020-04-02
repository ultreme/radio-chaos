package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"

	"github.com/bwmarrin/discordgo"
	yaml "gopkg.in/yaml.v2"
	"ultre.me/recettator"
)

type commandFunc func(s *discordgo.Session, m *discordgo.MessageCreate) error

var commands map[string]commandFunc

func init() {
	//
	// register custom commands
	//
	commands = map[string]commandFunc{
		"!history":                 doHistory,
		"!help":                    doHelp,
		"!il-est-pas-quelle-heure": doIlEstPasQuelleHeure,
		"!bite":                    doBite,
		"!recettator":              doRecettator,
	}
	// FIXME: !pause 5min
	// FIXME: !pipotron
	// FIXME: !blague

	//
	// generate commands based on `replies.yml`
	//
	var repliesYaml map[string][]string
	yamlFile, err := ioutil.ReadFile("replies.yml")
	if err != nil {
		log.Fatalf("yamlFile.Get: %v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &repliesYaml)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	for key, msgs := range repliesYaml {
		commands["!"+key] = genericRepliesYaml(msgs)
	}
}

func doRecettator(s *discordgo.Session, m *discordgo.MessageCreate) error {
	rctt := recettator.New(int64(rand.Intn(1000))) // FIXME: make it overridable by arg
	rctt.SetSettings(recettator.Settings{
		MainIngredients:      uint64(rand.Intn(2) + 1),
		SecondaryIngredients: uint64(rand.Intn(2) + 1),
		Steps:                uint64(rand.Intn(4) + 3),
	})
	markdown, err := rctt.Markdown()
	if err != nil {
		return err
	}
	s.ChannelMessageSend(m.ChannelID, "```markdown\n"+markdown+"\n```")
	return nil
}

func doHistory(s *discordgo.Session, m *discordgo.MessageCreate) error {
	resp, err := http.Get("http://radio-admin.casse-tete.solutions/?action=infos&format=json")
	if err != nil {
		return err
	}
	var parsed GetInfosResponse

	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return err
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
	return nil
}

func doHelp(s *discordgo.Session, m *discordgo.MessageCreate) error {
	keys := make([]string, len(commands))
	i := 0
	for key := range commands {
		keys[i] = key
		i++
	}

	sort.Strings(keys)
	out := strings.Join(keys, ", ")
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

func doBite(s *discordgo.Session, m *discordgo.MessageCreate) error {
	switch m.Author.Username {
	case "manfred":
		s.ChannelMessageSend(m.ChannelID, "B"+strings.Repeat("=", rand.Intn(10)+42)+"D")
	case "sassou":
		s.ChannelMessageSend(m.ChannelID, "{(.)}")
	default:
		s.ChannelMessageSend(m.ChannelID, "B"+strings.Repeat("=", rand.Intn(42)+1)+"D")
	}
	return nil
}

func doIlEstPasQuelleHeure(s *discordgo.Session, m *discordgo.MessageCreate) error {
	out := fmt.Sprintf("%02d:%02d",
		rand.Intn(24),
		rand.Intn(60),
	)
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

// see replies.yml
func genericRepliesYaml(msgs []string) commandFunc {
	return func(s *discordgo.Session, m *discordgo.MessageCreate) error {
		msg := msgs[rand.Intn(len(msgs))]
		s.ChannelMessageSend(m.ChannelID, msg)
		return nil
	}
}
