package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"os"
	"strings"
)

var debug = false

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

// Debug enables or disables debug mode
func Debug(state bool) {
	debug = state
}

// ContainsLowercase signals if a given string contains lowercase, ignoring Slack @mentions, #channels, URLs and :emoji:
func ContainsLowercase(m string) bool {
	// Split string into chars
	chars := strings.Split(m, "")
	l := len(chars)
	out := ""
	before := ""

	// State
	channel := false
	mention := false
	url := false
	emoji := false

	for i, char := range chars {
		switch char {
		// #channel
		case "#":
			if url {
				break
			}
			if before == "" || before == " " {
				channel = true
			} else {
				out += "#"
			}
		// @mention
		case "@":
			if url {
				break
			}
			if before == "" || before == " " {
				mention = true
			} else {
				out += "@"
			}
		// URL
		case "h":
			if channel || mention || url || emoji {
				break
			}
			if (before == "" || before == " ") && (l-i > 7 && strings.Join(chars[i:i+7], "") == "http://" || l-i > 8 && strings.Join(chars[i:i+8], "") == "https://") {
				url = true
			} else {
				out += "h"
			}
		// :emoji:
		case ":":
			if url {
				break
			}
			if emoji {
				emoji = false
				break
			}
			if before == "" || before == " " {
				emoji = true
			} else {
				out += ":"
			}
		// Terminate when seeing a space
		case " ":
			if channel || mention || url || emoji {
				channel = false
				mention = false
				url = false
				emoji = false
			}
			out += " "
		default:
			if !(channel || mention || url || emoji) {
				out += char
			}
		}

		if debug {
			fmt.Printf("'%s' '%s' '%s' (%t %t %t %t) %d %d %d\n", char, before, out, channel, mention, url, emoji, i, l, l-i)
		}

		before = char
	}

	if debug {
		fmt.Printf("'%s' => '%s'\n", m, out)
	}
	return strings.ToUpper(out) != out
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

				if ev.Msg.User != info.User.ID && ev.SubType == "" && ContainsLowercase(ev.Text) {
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
