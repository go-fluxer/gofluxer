package main

import (
	"fmt"
	"github.com/go-fluxer/gofluxer"
)

func main() {
	bot := gofluxer.NewBot("FLUXERBOTTOKEN", "!")
	// Replace FLUXERBOTTOKEN with your actual fluxer.app bot token.

	bot.OnReady(func() {
		// status must be one of: online, idle, dnd, invisible, offline.
		// custom_status is optional; pass nil to leave the current one alone.
		err := bot.SetPresenceWithCustomStatus("online", &gofluxer.CustomStatus{
			Text:      "🏎️ Hello",
			EmojiName: "🏎️",
		})
		if err != nil {
			fmt.Println("Failed to set presence:", err)
		}
	})

	bot.Run()
}
