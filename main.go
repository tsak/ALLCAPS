package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"math/rand"
	"os"
	"strings"
	"time"
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

// ContainsLowercase signals if a given string endWith lowercase, ignoring Slack @mentions, #channels, URLs and :emoji:
func ContainsLowercase(m string) bool {
	// Split string into chars
	chars := strings.Split(m, "")
	l := len(chars)
	out := ""
	before := ""
	next := ""

	// State
	channel := false
	mention := false
	url := false
	emoji := false

	for i, char := range chars {
		if i < l-1 {
			next = chars[i+1]
		}
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
			if emoji && next != ":" {
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

var (
	responses = []string{
		":male-police-officer: ALLCAPS POLICE :male-police-officer: IS YOUR CAPS LOCK BROKEN?\n\n> %s",
		":female-police-officer: ALLCAPS POLICE :female-police-officer: CAPS AND REGISTRATION PLEASE!\n\n> %s",
		":male-police-officer: ALLCAPS POLICE :female-police-officer: PLEASE KEEP YOUR CAPS ON THE LOCK!\n\n> %s",
		":female-police-officer: ALLCAPS POLICE :male-police-officer: HAVE YOU SEEN THESE CAPS BEFORE?\n\n> %s",
	}
	rl = len(responses)
)

// Enforcement wraps a given message into a helpful ALLCAPS POLICE response
func Enforcement(m string) string {
	// Ignore empty string
	if m == "" {
		return ":male-police-officer: ALLCAPS POLICE :female-police-officer: NOTHING TO SEE HERE, MOVE ALONG CAPS"
	}
	// Avoid double quoting
	if strings.HasPrefix(m, "> ") {
		m = m[2:]
	}
	return fmt.Sprintf(responses[rand.Intn(rl)], strings.ToUpper(m))
}

func main() {
	token := getenv("SLACKTOKEN")
	api := slack.New(token)
	rtm := api.NewRTM()
	rand.Seed(time.Now().UnixNano())
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
					rtm.SendMessage(rtm.NewOutgoingMessage(Enforcement(ev.Text), ev.Channel))
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
