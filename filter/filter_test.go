package filter

import (
	"testing"
)

func TestMatchHost(t *testing.T) {
	tests := []struct {
		host    string
		pattern string
		want    bool
	}{
		// Exact match
		{"api.ipify.org", "api.ipify.org", true},
		{"api.ipify.org", "api.ipify.com", false},

		// Subdomain wildcard *.domain.com
		{"www.youtube.com", "*.youtube.com", true},
		{"m.youtube.com", "*.youtube.com", true},
		{"youtube.com", "*.youtube.com", false}, // root domain — no match

		// Multi-segment wildcard
		{"video.googlevideo.com", "*.googlevideo.com", true},

		// Contains wildcard *discord*.*
		{"discord.com", "*discord*.*", true},
		{"cdn.discordapp.com", "*discord*.*", true},
		{"notdiscord.net", "*discord*.*", true}, // "discord" is in "notdiscord" → matches
		{"example.com", "*discord*.*", false},   // no "discord" → no match

		// Short host with wildcard
		{"youtu.be", "youtu.be", true},

		// Port stripping
		{"www.youtube.com:443", "*.youtube.com", true},
		{"api.ipify.org:80", "api.ipify.org", true},

		// Case insensitivity
		{"WWW.YouTube.COM", "*.youtube.com", true},

		// No match
		{"example.com", "*.youtube.com", false},
	}

	for _, tt := range tests {
		name := tt.host + " ~ " + tt.pattern
		t.Run(name, func(t *testing.T) {
			got := MatchHost(tt.host, tt.pattern)
			if got != tt.want {
				t.Errorf("MatchHost(%q, %q) = %v, want %v",
					tt.host, tt.pattern, got, tt.want)
			}
		})
	}
}
