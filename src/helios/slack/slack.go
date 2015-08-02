package slack

import (
	"helios/helios"
	"os"
	"time"

	"github.com/nlopes/slack"
)

var MessageChann chan interface{}

func Service() helios.ServiceHandler {
	return func(h *helios.Engine) error {

		// Setup broadcaster and save channel to a global
		MessageChann = h.NewBroadcastChannel("slack")

		err := initSlackRTM()
		if err != nil {
			return err
		}

		return nil
	}
}

func initSlackRTM() error {
	token := os.Getenv("SLACK_API")

	chReceiver := make(chan slack.SlackEvent)

	api := slack.New(token)

	rtm, err := api.StartRTM("", "http://localhost/")
	if err != nil {
		return err
	}

	go rtm.HandleIncomingEvents(chReceiver)
	go rtm.Keepalive(20 * time.Second)
	go messageHandler(chReceiver)

	return nil
}

func messageHandler(c chan slack.SlackEvent) {
	for {
		select {
		case msg := <-c:
			switch msg.Data.(type) {
			case *slack.MessageEvent:
				MessageChann <- "message"
			}
		}
	}
	panic("Shouldn't be here")
}
