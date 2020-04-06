package main

import (
	"encoding/json"
	"log"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/protocol"
	"github.com/graarh/golang-socketio/transport"
	"moul.io/godev"
)

type sioClient struct {
	c *gosocketio.Client
}

func (c *sioClient) Close() {
	c.c.Close()
}

func (c *sioClient) logError(err error) {
	log.Printf("error: %v", err)
	// FIXME: send error on socketio so it can be logged on discord
}

func newSIOClient(name string) (*sioClient, error) {
	ws := transport.GetDefaultWebsocketTransport()
	ws.PingTimeout = 15 * time.Second
	c, err := gosocketio.Dial(
		gosocketio.GetUrl("sockethub.moul.io", 443, true),
		ws,
	)
	if err != nil {
		return nil, err
	}

	c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("--- disconnected")
	})
	c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Print("--- connected")
	})
	c.On("event:join", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got event:join: %v", godev.JSON(args))
	})
	c.On("event:disconnect", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got event:disconnect: %v", godev.JSON(args))
	})
	c.On("event:broadcast", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got event:broadcast: %v", godev.JSON(args))
	})

	maxLogEntries := 0
	input := Message{Room: "radio-chaos", MaxLogEntries: &maxLogEntries}
	input.Peer.Name = name
	c.Emit("join", input)

	client := sioClient{c: c}
	// ping loop
	go func() {
		// FIXME: close on Close()
		for {
			c.Emit("ping", protocol.MessageTypePing)
			time.Sleep(5 * time.Second)
		}
	}()
	return &client, nil
}

type Peer struct {
	Name string `json:"name,omitempty"`
}

type Msg struct {
	Text string `json:"text,omitempty"`
}

type Message struct {
	Peer          Peer    `json:"peer,omitempty"`
	Room          string  `json:"room,omitempty"`
	IsLive        *bool   `json:"is_live,omitempty"`
	Msg           *Msg    `json:"msg,omitempty"`
	Peers         *[]Peer `json:"peers,omitempty"`
	MaxLogEntries *int    `json:"max_log_entries,omitempty"`
}

type ProtocolMsg struct {
	Action string
	Args   []string
}

func (pm ProtocolMsg) ToJSON() string {
	out, _ := json.Marshal(pm)
	return string(out)
}
