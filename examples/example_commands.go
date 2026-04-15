package main

import (
	"fmt"
	"strings"
	"github.com/go-fluxer/gofluxer"
)

func main() {
	bot := gofluxer.NewBot("FLUXERBOTTOKEN", "!")
	// Replace FLUXERBOTTOKEN with your actual fluxer.app bot token.

	bot.AddCommand("ping", func(m *gofluxer.Message) {
		bot.SendMessage(m.ChannelID, "Pong!")
	})

	bot.AddCommand("say", func(m *gofluxer.Message, args []string) {
		if len(args) == 0 {
			bot.SendMessage(m.ChannelID, "What do you want me to say?")
			return
		}
		bot.SendMessage(m.ChannelID, strings.Join(args, " "))
	})

	fmt.Println("Gofluxer Bot is getting Ready")
	if err := bot.Run(); err != nil {
		fmt.Printf("Gofluxer Bot stopped: %v\n", err)
	}
}