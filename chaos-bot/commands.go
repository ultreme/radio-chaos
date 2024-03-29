package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gohugoio/hugo/common/maps"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	tpyo "github.com/tpyolang/tpyo-cli"
	hpeg "github.com/ultreme/histoire-pour-enfant-generator"
	yaml "gopkg.in/yaml.v2"
	"moul.io/moulsay/moulsay"
	ntw "moul.io/number-to-words"
	"moul.io/pipotron/dict"
	"moul.io/pipotron/pipotron"
	"ultre.me/recettator"
	"ultre.me/smsify/smsify"
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
		"!histoire-pour-enfant":    doHistoirePourEnfant,
		"!moulsay":                 doMoulsay,
		"!roulette":                doRoulette,
		"!nombre":                  doNombre,
		"!smsify":                  doSmsify,
		"!ntw":                     doNtw,
		"!tpyo":                    doTpyo,
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

	//
	// pipotron
	//
	dicts := []string{"proverbe-africain", "marabout", "reve", "whatsapp-message-in-case-of-pandemic", "question-baleze-raw", "prenom-compose", "moijaime", "insulte-mignonne", "horoscope", "fuu", "excuse-a-2-balles", "bingo-winner", "asv", "accords", "project-idea", "compliment", "blague", "complot", "popular-science-bestseller-title"}
	//fmt.Println(godev.PrettyJSON(dict.Box.List()))
	for _, dictName := range dicts {
		commands["!"+dictName] = genericPipotron(dictName)
	}
}

// see https://github.com/moul/pipotron
func genericPipotron(name string) commandFunc {
	dictFile, err := dict.Box.Find(name + ".yml")
	if err != nil {
		log.Printf("warn: %v", err)
		return nil
	}
	var context pipotron.Context
	err = yaml.Unmarshal(dictFile, &context.Dict)
	if err != nil {
		log.Printf("warn: %v", err)
		return nil
	}
	context.Scratch = maps.NewScratch()
	return func(s *discordgo.Session, m *discordgo.MessageCreate) error {
		out, err := pipotron.Generate(&context)
		if err != nil {
			return err
		}
		s.ChannelMessageSend(m.ChannelID, out)
		return nil
	}
}

func doNtw(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Split(m.Content, " ")
	if len(content) < 2 {
		return nil
	}
	input, err := strconv.Atoi(content[1])
	if err != nil {
		return err
	}
	out := ntw.IntegerToFrFr(input)
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

func doHistoirePourEnfant(s *discordgo.Session, m *discordgo.MessageCreate) error {
	story := hpeg.NewStory()

	for i := 0; i < 4; i++ {
		story.AddElement(hpeg.NewAnimal())
	}

	lines := story.Tell()
	lines[0] = "# " + lines[0]
	markdown := strings.Join(lines, "\n")
	msg, err := sendMarkdown(s, m, markdown)
	if err != nil {
		return err
	}
	s.MessageReactionAdd(m.ChannelID, msg.ID, "👶")
	return nil
}

func doRecettator(s *discordgo.Session, m *discordgo.MessageCreate) error {
	rctt := recettator.New(int64(rand.Intn(1000))) // FIXME: make it overridable by arg
	rctt.SetSettings(recettator.Settings{
		MainIngredients:      uint64(rand.Intn(2) + 2),
		SecondaryIngredients: uint64(rand.Intn(2) + 2),
		Steps:                uint64(rand.Intn(4) + 3),
	})
	markdown, err := rctt.Markdown()
	if err != nil {
		return err
	}
	msg, err := sendMarkdown(s, m, markdown)
	if err != nil {
		return err
	}
	s.MessageReactionAdd(m.ChannelID, msg.ID, "😋")
	s.MessageReactionAdd(m.ChannelID, msg.ID, "🤮")
	return nil
}

func doNombre(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Split(m.Content, " ")
	if len(content) < 3 {
		s.ChannelMessageSend(m.ChannelID, "exemple: !nombre 1 100")
		return nil
	}
	a, err := strconv.Atoi(content[1])
	if err != nil {
		return err
	}
	b, err := strconv.Atoi(content[2])
	if err != nil {
		return err
	}
	nbr := rand.Intn(b) + a
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d", nbr))
	return nil
}

func doSmsify(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Join(strings.Split(m.Content, " ")[1:], " ")
	msg := smsify.Smsify(content)
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func doTpyo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Join(strings.Split(m.Content, " ")[1:], " ")
	tpyo := tpyo.NewTpyo()
	msg := tpyo.Enocde(content)
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}

func doRoulette(s *discordgo.Session, m *discordgo.MessageCreate) error {
	os.MkdirAll("data", 0777)
	db, err := gorm.Open("sqlite3", "data/roulette.db")
	if err != nil {
		return err
	}
	defer db.Close()
	type roulette struct {
		gorm.Model
		Count int
	}
	if err := db.AutoMigrate(&roulette{}).Error; err != nil {
		return err
	}
	var entry roulette
	db.First(&entry)
	if entry.Count == 0 {
		entry.Count = rand.Intn(6) + 1
		s.ChannelMessageSend(m.ChannelID, "**BANG!**")
	} else {
		entry.Count--
		s.ChannelMessageSend(m.ChannelID, "* **click** *")
	}
	db.Save(&entry)
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
	keys := []string{}
	for key, fn := range commands {
		if fn == nil {
			continue
		}
		keys = append(keys, key)
	}

	sort.Strings(keys)
	out := strings.Join(keys, ", ")
	s.ChannelMessageSend(m.ChannelID, out)
	return nil
}

func doMoulsay(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Split(m.Content, " ")
	say := strings.Join(content[1:], " ")
	say = strings.TrimSpace(say)
	if say == "" {
		return nil
	}
	out, err := moulsay.Say(say, 60)
	if err != nil {
		return err
	}
	sendBlock(s, m, out)
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
