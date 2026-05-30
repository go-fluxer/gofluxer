package main

import (
	"fmt"
	"github.com/go-fluxer/gofluxer"
)

func main() {
	// Webhook urls look something like this: https://api.fluxer.app/webhooks/<WEBHOOKID>/<WEBHOOKTOKEN>
	wh := gofluxer.NewWebhookClient("WEBHOOKID", "WEBHOOKTOKEN")
	wh.Execute("Hello World!")
	fmt.Printf("[gofluxer]: Message sent to webhook.")
}