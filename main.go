package main

import (
	"fmt"
	"github.com/nlopes/slack"
	"html"
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
	channel := -1
	mention := -1
	url := -1
	emoji := -1

	for i, char := range chars {
		if i < l-1 {
			next = chars[i+1]
		}
		switch char {
		// #channel
		case "#":
			if url != -1 {
				break
			}
			if emoji != -1 {
				// Reset emoji mode and append characters since emoji begin was encountered
				out += m[emoji:i]
				emoji = -1
				break
			}
			if before == "" || before == " " {
				channel = i
			} else {
				out += "#"
			}
		// @mention
		case "@":
			if url != -1 {
				break
			}
			if emoji != -1 {
				// Reset emoji mode and append characters since emoji begin was encountered
				out += m[emoji:i]
				emoji = -1
				break
			}
			if before == "" || before == " " {
				mention = i
			} else {
				out += "@"
			}
		// URL
		case "h":
			if channel != -1 || mention != -1 || url != -1 || emoji != -1 {
				break
			}
			if (before == "" || before == " ") && (l-i > 7 && strings.Join(chars[i:i+7], "") == "http://" || l-i > 8 && strings.Join(chars[i:i+8], "") == "https://") {
				url = i
			} else {
				out += "h"
			}
		// :emoji:
		case ":":
			if url != -1 {
				break
			}
			if channel != -1 {
				// Reset channel mode and append characters since channel begin was encountered
				out += m[channel:i]
				channel = -1
				break
			}
			if mention != -1 {
				// Reset mention mode and append characters since mention begin was encountered
				out += m[mention:i]
				mention = -1
				break
			}
			if emoji != -1 && next != ":" && before != ":" {
				emoji = -1
				break
			}
			if before == "" || before == " " {
				emoji = i
			} else {
				out += ":"
			}
		// Terminate when seeing a space
		case " ":
			if channel != -1 || mention != -1 || url != -1 || emoji != -1 {
				channel = -1
				mention = -1
				url = -1
				emoji = -1
			}
			out += " "
		default:
			if !(channel != -1 || mention != -1 || url != -1 || emoji != -1) {
				out += char
			}
		}

		if debug {
			fmt.Printf("'%s' '%s' '%s' (%d %d %d %d) %d %d %d\n", char, before, out, channel, mention, url, emoji, i, l, l-i)
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

	up := strings.ToUpper(m)

	// Replace skin-tone modifier with lowercase
	up = strings.Replace(up, ":SKIN-TONE-", ":skin-tone-", -1)

	return fmt.Sprintf(responses[rand.Intn(rl)], up)
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

				text := html.UnescapeString(ev.Text)
				if ev.Msg.User != info.User.ID && ev.SubType == "" && ContainsLowercase(text) {
					rtm.SendMessage(rtm.NewOutgoingMessage(Enforcement(text), ev.Channel))
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
