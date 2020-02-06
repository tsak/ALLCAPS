package main

import (
	"strings"
	"testing"
)

func TestContainsLowercase(t *testing.T) {
	//Debug(true)
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{name: "nil", input: "", expected: false},
		{name: "simple lowercase", input: "foo", expected: true},
		{name: "SIMPLE UPPERCASE", input: "BAR", expected: false},
		// #channels
		{name: "simple #channel", input: "#channel", expected: false},
		{name: "#channel at beginning with lowercase", input: "#channel bar", expected: true},
		{name: "#channel at beginning with UPPERCASE", input: "#channel BAR", expected: false},
		{name: "#channel in middle, surrounded by lowercase", input: "foo #channel bar", expected: true},
		{name: "#channel in middle, surrounded by UPPERCASE", input: "FOO #channel BAR", expected: false},
		{name: "#channel at end", input: "FOO #channel", expected: false},
		{name: "#channel no space", input: "foo#channel", expected: true},
		{name: "#channel no space", input: "FOO#channel", expected: true},
		{name: "#channel no space", input: "FOO#CHANNEL", expected: false},
		// @mention
		{name: "simple @mention", input: "@mention", expected: false},
		{name: "@mention at beginning with lowercase", input: "@mention bar", expected: true},
		{name: "@mention at beginning with UPPERCASE", input: "@mention BAR", expected: false},
		{name: "@mention in middle, surrounded by lowercase", input: "foo @mention bar", expected: true},
		{name: "@mention in middle, surrounded by UPPERCASE", input: "FOO @mention BAR", expected: false},
		{name: "@mention at end", input: "FOO @mention", expected: false},
		{name: "@mention no space", input: "foo@mention", expected: true},
		{name: "@mention no space", input: "FOO@mention", expected: true},
		{name: "@mention no space", input: "FOO@MENTION", expected: false},
		// HTTP URL
		{name: "simple http://user:pass@domain.com#anchor", input: "http://user:pass@domain.com#anchor", expected: false},
		{name: "http://user:pass@domain.com#anchor at beginning with lowercase", input: "http://user:pass@domain.com#anchor bar", expected: true},
		{name: "http://user:pass@domain.com#anchor at beginning with UPPERCASE", input: "http://user:pass@domain.com#anchor BAR", expected: false},
		{name: "http://user:pass@domain.com#anchor in middle, surrounded by lowercase", input: "foo http://user:pass@domain.com#anchor bar", expected: true},
		{name: "http://user:pass@domain.com#anchor in middle, surrounded by UPPERCASE", input: "FOO http://user:pass@domain.com#anchor BAR", expected: false},
		{name: "http://user:pass@domain.com#anchor at end", input: "FOO http://user:pass@domain.com#anchor", expected: false},
		{name: "http://user:pass@domain.com#anchor no space", input: "foohttp://user:pass@domain.com#anchor", expected: true},
		{name: "http://user:pass@domain.com#anchor no space", input: "FOOhttp://user:pass@domain.com#anchor", expected: true},
		// HTTPS URL
		{name: "simple https://user:pass@domain.com#anchor", input: "https://user:pass@domain.com#anchor", expected: false},
		{name: "https://user:pass@domain.com#anchor at beginning with lowercase", input: "https://user:pass@domain.com#anchor bar", expected: true},
		{name: "https://user:pass@domain.com#anchor at beginning with UPPERCASE", input: "https://user:pass@domain.com#anchor BAR", expected: false},
		{name: "https://user:pass@domain.com#anchor in middle, surrounded by lowercase", input: "foo https://user:pass@domain.com#anchor bar", expected: true},
		{name: "https://user:pass@domain.com#anchor in middle, surrounded by UPPERCASE", input: "FOO https://user:pass@domain.com#anchor BAR", expected: false},
		{name: "https://user:pass@domain.com#anchor at end", input: "FOO https://user:pass@domain.com#anchor", expected: false},
		{name: "https://user:pass@domain.com#anchor no space", input: "foohttps://user:pass@domain.com#anchor", expected: true},
		{name: "https://user:pass@domain.com#anchor no space", input: "FOOhttps://user:pass@domain.com#anchor", expected: true},
		// :emoji:
		{name: "simple :emoji:", input: ":emoji:", expected: false},
		{name: ":emoji: at beginning with lowercase", input: ":emoji: bar", expected: true},
		{name: ":emoji: at beginning with UPPERCASE", input: ":emoji: BAR", expected: false},
		{name: ":emoji: in middle, surrounded by lowercase", input: "foo :emoji: bar", expected: true},
		{name: ":emoji: in middle, surrounded by UPPERCASE", input: "FOO :emoji: BAR", expected: false},
		{name: ":emoji: at end", input: "FOO :emoji:", expected: false},
		{name: ":emoji: no space", input: "foo:emoji:", expected: true},
		{name: ":emoji: no space", input: "FOO:emoji:", expected: true},
		{name: ":emoji: no space", input: "FOO:EMOJI:", expected: false},
		// Mix
		{name: "#channelmoji:", input: "contains a #channelmoji: in the text", expected: true},
		{name: ":emoji#chan", input: ":emoji#chan", expected: true},
		{name: ":emoji@mention", input: ":emoji@mention", expected: true},
		{name: "@mention:moji", input: "@mention:moji", expected: true},
		{name: "@mention:moji:", input: "@mention:moji:", expected: true},
		// Skin tone
		{name: ":skin::tone:", input: ":+1::skin-tone-6:", expected: false},
		{name: ":skin::tone: with text", input: ":+1::skin-tone-6: with lowercase", expected: true},
		{name: ":skin::tone: WITH TEXT", input: ":+1::skin-tone-6: WITH UPPERCASE", expected: false},
		// Unicode
		{name: "sᴍᴀʟʟ ᴄᴀᴘs", input: "sᴍᴀʟʟ ᴄᴀᴘs", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ContainsLowercase(tt.input)

			if got != tt.expected {
				if tt.expected {
					t.Errorf("'%s' was not detected as lowercase", tt.input)
				} else {
					t.Errorf("'%s' WAS NOT DETECTED AS UPPERCASE", tt.input)
				}
			}
		})
	}
}

func TestEnforcement(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		endWith string
	}{
		{name: "nil", input: "", endWith: "NOTHING TO SEE HERE, MOVE ALONG CAPS"},
		{name: "quoted", input: "> foo", endWith: "> FOO"},
		{name: "normal", input: "bar", endWith: "> BAR"},
		{name: "normal", input: "baz", endWith: "> BAZ"},
		{name: "normal", input: "foo", endWith: "> FOO"},
		{name: "normal", input: "bar", endWith: "> BAR"},
		{name: "normal", input: "baz", endWith: "> BAZ"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Enforcement(tt.input)
			if !strings.HasSuffix(got, tt.endWith) {
				t.Errorf("Expected '%s' to end with '%s'", got, tt.endWith)
			}
		})
	}
}
