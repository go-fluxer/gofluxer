package main

import (
	"fmt"
	"strings"
	"github.com/go-fluxer/gofluxer"
)

func main() {
	bot := gofluxer.NewBot("FLUXERBOTTOKEN", "!")
	// Replace FLUXERBOTTOKEN with your actual fluxer.app bot token.

	bot.OnMessage(func(m *gofluxer.Message) {
		fmt.Printf("[%s]: %s\n", m.Author.Username, m.Content)
		
		if strings.ToLower(m.Content) == "ping" {
			bot.SendMessage(m.ChannelID, "Pong!")
		}
	})

	fmt.Println("Gofluxer Bot is getting Ready")
	if err := bot.Run(); err != nil {
		fmt.Printf("Gofluxer Bot stopped: %v\n", err)
	}
}