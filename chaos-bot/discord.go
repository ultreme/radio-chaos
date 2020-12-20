package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/etherlabsio/errors"
	"github.com/etherlabsio/pkg/httputil"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httplog"
	gosocketio "github.com/graarh/golang-socketio"
	"github.com/hako/durafmt"
	"github.com/oklog/run"
	"github.com/rs/cors"
	"moul.io/godev"
)

type discordbotOpts struct {
	devMode        bool
	debug          bool
	manfredChannel string
	discordToken   string
	serverBind     string
}

func discordBotCmd(opts discordbotOpts) error {
	log.Printf("starting bot, args=%s", godev.PrettyJSON(opts))

	db := &discordBot{
		startedAt: time.Now(),
		opts:      opts,
	}

	var g run.Group

	{ // DISCORD
		dg, err := discordgo.New("Bot " + opts.discordToken)
		if err != nil {
			return err
		}

		hostname, _ := os.Hostname()
		dg.ChannelMessageSend(opts.manfredChannel, fmt.Sprintf("**COUCOU JE VIENS DE BOOT (%s)!**", hostname))

		db.discordSession = dg
		commands["!info"] = db.doInfo
		commands["!radiosay"] = db.doRadioSay

		dg.AddHandler(db.onMessageCreate)

		if err = dg.Open(); err != nil {
			return err
		}
		defer dg.Close()
		log.Print("Discord bot started")
	}

	{ // SOCKET IO
		sio, err := newSIOClient("mixer")
		if err != nil {
			return err
		}
		defer sio.Close()

		sio.c.On("event:join", func(h *gosocketio.Channel, args Message) {
			db.discordSession.ChannelMessageSend(opts.manfredChannel, fmt.Sprintf("sio event:join: %s", godev.JSON(args)))
		})
		sio.c.On("event:disconnect", func(h *gosocketio.Channel, args Message) {
			db.discordSession.ChannelMessageSend(opts.manfredChannel, fmt.Sprintf("sio event:disconnect: %s", godev.JSON(args)))
		})
		sio.c.On("event:broadcast", func(h *gosocketio.Channel, args Message) {
			db.discordSession.ChannelMessageSend(opts.manfredChannel, fmt.Sprintf("sio event:broadcast: %s", godev.JSON(args)))
		})
		db.sio = sio
		log.Print("Socket.IO client started")
	}

	{ // HTTP API
		r := chi.NewRouter()
		//r.Use(middleware.DefaultCompress)
		r.Use(middleware.StripSlashes)
		r.Use(middleware.Recoverer)
		r.Use(cors.New(cors.Options{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		}).Handler)
		logger := httplog.NewLogger("http", httplog.Options{JSON: false})
		r.Use(httplog.RequestLogger(logger))
		//r.Use(middleware.Heartbeat("/ping"))

		httpStatusCodeFrom := func(err error) int {
			switch errors.KindOf(err) {
			case errors.Invalid:
				return http.StatusBadRequest
			case errors.Internal:
				return http.StatusInternalServerError
			default:
				return http.StatusOK
			}
		}
		httpErrorEncoder := httputil.JSONErrorEncoder(httpStatusCodeFrom)
		httpJSONResponseEncoder := httputil.EncodeJSONResponse(httpErrorEncoder)

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			resp := struct {
				Pong string `json:"pong"`
			}{
				Pong: "pong",
			}
			httpJSONResponseEncoder(ctx, w, resp)
		})

		httpListener, err := net.Listen("tcp", db.opts.serverBind)
		if err != nil {
			return err
		}
		g.Add(func() error {
			return http.Serve(httpListener, r)
		}, func(error) {
			httpListener.Close()
		})

		log.Printf("HTTP API started (%q)", db.opts.serverBind)
	}

	g.Add(run.SignalHandler(context.TODO(), os.Interrupt))
	log.Print("Press ctrl-C to exit")
	return g.Run()
}

type discordBot struct {
	startedAt      time.Time
	seenMessages   int
	seenCommands   int
	seenErrors     int
	discordSession *discordgo.Session
	lock           sync.Mutex
	sio            *sioClient
	opts           discordbotOpts
}

func (db *discordBot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// avoid loop
	if m.Author.ID == s.State.User.ID {
		return
	}
	if db.opts.debug {
		log.Println(godev.JSON(m))
	}
	if db.opts.devMode && m.GuildID != "" {
		return
	}

	args := strings.Split(m.Content, " ")

	db.seenMessages++
	command, found := commands[args[0]] // FIXME: split the content and support args
	if !found || command == nil {
		return
	}
	db.seenCommands++

	db.lock.Lock()
	defer db.lock.Unlock()

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
		db.seenErrors++
		sendError(s, m, err)
	}
}

func (db *discordBot) doRadioSay(s *discordgo.Session, m *discordgo.MessageCreate) error {
	content := strings.Split(m.Content, " ")
	say := strings.Join(content[1:], " ")
	say = strings.TrimSpace(say)
	if say == "" {
		return nil
	}

	input := Message{
		Room: "radio-chaos",
		Msg: &Msg{
			Text: ProtocolMsg{
				Action: "radiosay",
				Args:   []string{say},
			}.ToJSON(),
		},
	}
	return db.sio.c.Emit("broadcast", input)
}

func (db *discordBot) doInfo(s *discordgo.Session, m *discordgo.MessageCreate) error {
	uptime := time.Since(db.startedAt)
	msg := fmt.Sprintf(
		"uptime: %v, messages: %d, commands: %d, errors: %d\n",
		durafmt.ParseShort(uptime).String(), db.seenMessages, db.seenCommands, db.seenErrors,
	)
	msg += fmt.Sprintf("source: https://github.com/ultreme/radio-chaos/tree/master/chaos-bot")
	s.ChannelMessageSend(m.ChannelID, msg)
	return nil
}
