package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"regexp"
	"strings"
)

var evilLowerCase = regexp.MustCompile("[a-z]")

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	token := getenv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

Loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			//fmt.Printf("Event Received: %s %+v\n", msg.Type, msg.Data)
			switch ev := msg.Data.(type) {

			case *slack.MessageEvent:
				info := rtm.GetInfo()

				matched := evilLowerCase.MatchString(ev.Text)

				if ev.Msg.User != info.User.ID && ev.SubType == "" && matched {
					rtm.SendMessage(rtm.NewOutgoingMessage(fmt.Sprintf("*ALL CAPS PLEASE!* %s", strings.ToUpper(ev.Text)), ev.Channel))
				}

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break Loop

			default:
				// Take no action
			}
		}
	}
}
