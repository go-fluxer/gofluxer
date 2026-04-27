# gofluxer

gofluxer is an API wrapper that allows you to create fluxer.app bots as well as use fluxer.app webhooks

If you need help with using the gofluxer package, join the [gofluxer Fluxer server](https://fluxer.gg/KTXuuy7k). 

**NOTICE**: This library is still a heavy work in progress. Expect unfinished and buggy features.

This Go package is not officially endorsed by or affiliated with fluxer.app

## Getting Started

This assumes you already have a Go environment on your system. If not, download Go [from here](https://go.dev/dl). Make sure you have Go 1.21 or newer.

Install gostoat by using the following command.

```sh
go get github.com/go-fluxer/gofluxer
```

After installing the package, import the gostoat package using this within your code.

```go
import (
	"fmt"
	"strings"
	"github.com/go-fluxer/gofluxer"
)
```

Here is a basic example of a ping pong bot.

```go
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
```

You can refer to the examples directory for more examples to see how to use this Go package

# gofluxer Update Notes:

### Version 0.1.1 - April 26th, 2026

- gofluxer will now try to reconnect if connection to fluxer has dropped.
- Updated heartbeat to fix a potiental issue where data could still try to be sent to closed connections.

### Version 0.1.0 - April 15th, 2026

- The first early release version of gofluxer.
- Supports message and command handlers.
- Supports NSFW channel and bot owner checks.
- Supports sending messages to webhooks.
- Uses Apache License 2.0.
- Added a example_pingpong.go example and a example_commands.go example.