package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/exec"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/oklog/run"
	"moul.io/godev"
)

func mixerBotCmd(devMode bool) error {
	name := "mixer"
	if devMode {
		name += "-dev"
	}
	sio, err := newSIOClient(name)
	if err != nil {
		return err
	}
	defer sio.Close()

	sio.c.On("event:broadcast", func(h *gosocketio.Channel, args Message) {
		if args.Msg == nil || args.Msg.Text == "" {
			return
		}
		var pm ProtocolMsg
		err := json.Unmarshal([]byte(args.Msg.Text), &pm)
		if err != nil {
			sio.logError(err)
			return
		}
		switch pm.Action {
		case "radiosay":
			cmd := exec.Command("say", pm.Args[0])
			err := cmd.Run()
			if err != nil {
				sio.logError(err)
				return
			}
		default:
			log.Printf("unsupported action: %s", godev.JSON(pm))
		}
	})

	var g run.Group
	g.Add(run.SignalHandler(context.TODO(), os.Interrupt))
	return g.Run()
}
