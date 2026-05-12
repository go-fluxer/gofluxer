package gofluxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	DefaultBaseURL    = "https://api.fluxer.app/v1"
	DefaultGatewayURL = "wss://gateway.fluxer.app/?v=1"
)
type Bot struct {
	Token         string
	CommandPrefix string
	BaseURL       string
	GatewayURL    string
	Handlers      []func(m *Message)
	Commands      map[string]func(m *Message, args []string)
	Conn          *websocket.Conn
}

func NewBot(token, prefix string) *Bot {
	return &Bot{
		Token:         token,
		CommandPrefix: prefix,
		BaseURL:       DefaultBaseURL,
		GatewayURL:    DefaultGatewayURL,
		Commands:      make(map[string]func(m *Message, args []string)),
	}
}
func (b *Bot) NewBotInstance(apiBase, gatewayURL string) {
	b.BaseURL = apiBase
	b.GatewayURL = gatewayURL
}



func (b *Bot) OnMessage(handler func(m *Message)) {
	b.Handlers = append(b.Handlers, handler)
}

func (b *Bot) AddCommand(name string, handler func(m *Message, args []string)) {
	b.Commands[name] = handler
}

func (b *Bot) checkRateLimit(statusCode int) {
	if statusCode == http.StatusTooManyRequests {
		fmt.Println("[gofluxer]: API returned a status 429 rate limit. Stopping process.")
		os.Exit(1)
	}
}

func (b *Bot) IsOwner(m *Message) bool {
	if m.GuildID == "" {
		return false
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/guilds/%s", b.BaseURL, m.GuildID), nil)
	req.Header.Set("Authorization", "Bot "+b.Token)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)

	var guild struct {
		OwnerID string `json:"owner_id"`
	}
	json.NewDecoder(resp.Body).Decode(&guild)
	return m.Author.ID == guild.OwnerID
}

func (b *Bot) IsNSFW(channelID string) bool {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/channels/%s", b.BaseURL, channelID), nil)
	req.Header.Set("Authorization", "Bot "+b.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)

	var channel struct {
		NSFW bool `json:"nsfw"`
	}
	json.NewDecoder(resp.Body).Decode(&channel)
	return channel.NSFW
}

func (b *Bot) SendMessage(channelID string, content string) {
	body, _ := json.Marshal(map[string]string{"content": content})
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/channels/%s/messages", b.BaseURL, channelID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		b.checkRateLimit(resp.StatusCode)
		resp.Body.Close()
	}
}

func (b *Bot) SendEmbed(channelID string, embed interface{}) {
	body, _ := json.Marshal(map[string]interface{}{"embeds": []interface{}{embed}})
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/channels/%s/messages", b.BaseURL, channelID), bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		b.checkRateLimit(resp.StatusCode)
		resp.Body.Close()
	}
}

func (b *Bot) Run() error {
	for {
		fmt.Println("[gofluxer]: Attempting to connect to Fluxer Gateway...")
		fmt.Printf("[gofluxer]: Connecting to %s...\n", b.GatewayURL)
		conn, _, err := websocket.DefaultDialer.Dial(b.GatewayURL, nil)
		if err != nil {
			fmt.Printf("[gofluxer]: Connection failed: %v. Retrying in 5 seconds...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		b.Conn = conn
		fmt.Println("[gofluxer]: Connected to Fluxer")
		err = b.listen(conn)
		fmt.Printf("[gofluxer]: Connection lost: %v. Reconnecting...\n", err)
		conn.Close()
		time.Sleep(2 * time.Second)
	}
}
func (b *Bot) listen(conn *websocket.Conn) error {
	for {
		var payload struct {
			Op int             `json:"op"`
			T  string          `json:"t"`
			D  json.RawMessage `json:"d"`
		}

		if err := conn.ReadJSON(&payload); err != nil {
			return err
		}

		switch payload.Op {
		case 10:
			var hello struct {
				HeartbeatInterval int `json:"heartbeat_interval"`
			}
			json.Unmarshal(payload.D, &hello)
			go b.heartbeat(time.Duration(hello.HeartbeatInterval) * time.Millisecond)
			b.identify()
		case 0:
			if payload.T == "MESSAGE_CREATE" {
				var m Message
				json.Unmarshal(payload.D, &m)
				if m.Author.Bot {
					continue
				}
				for _, h := range b.Handlers {
					h(&m)
				}
				if strings.HasPrefix(m.Content, b.CommandPrefix) {
					cleanContent := m.Content[len(b.CommandPrefix):]
					parts := strings.Fields(cleanContent)
					if len(parts) > 0 {
						cmdName := parts[0]
						args := parts[1:]
						if cmd, ok := b.Commands[cmdName]; ok {
							cmd(&m, args)
						}
					}
				}
			}
		}
	}
}

func (b *Bot) identify() {
	payload := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token":   b.Token,
			"intents": 512,
			"properties": map[string]string{
				"os":      "linux",
				"browser": "gofluxer",
				"device":  "gofluxer",
			},
		},
	}
	b.Conn.WriteJSON(payload)
}

func (b *Bot) heartbeat(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	currentConn := b.Conn

	for range ticker.C {
		if b.Conn != currentConn {
			return
		}
		err := b.Conn.WriteJSON(map[string]interface{}{"op": 1, "d": nil})
		if err != nil {
			return
		}
	}
}