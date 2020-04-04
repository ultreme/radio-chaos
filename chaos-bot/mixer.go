package main

import (
	"context"
	"log"
	"os"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/oklog/run"
)

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func mixerBot() error {
	c, err := gosocketio.Dial(
		gosocketio.GetUrl("sockethub.moul.io", 443, true),
		transport.GetDefaultWebsocketTransport(),
	)
	if err != nil {
		return err
	}

	err = c.On("/message", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got chat message: %v", args)
	})
	if err != nil {
		return err
	}
	err = c.On("event:join", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got event:join: %v", args)
	})
	if err != nil {
		return err
	}
	err = c.On("/event:join", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got /event:join: %v", args)
	})
	if err != nil {
		return err
	}
	err = c.On("event:broadcast", func(h *gosocketio.Channel, args Message) {
		log.Printf("--- got event:broadcast: %v", args)
	})
	if err != nil {
		return err
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("--- disconnected")
	})
	if err != nil {
		return err
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Print("--- connected")
	})
	if err != nil {
		return err
	}

	/*
		go sendJoin(c)
		go sendJoin(c)
		go sendJoin(c)
		go sendJoin(c)
		go sendJoin(c)
	*/

	var g run.Group
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			<-cancel
			return nil
		}, func(err error) {
			c.Close()
			close(cancel)
		})
	}
	g.Add(run.SignalHandler(context.TODO(), os.Interrupt))

	return g.Run()
}

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}
