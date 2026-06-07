package gofluxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
	"sync"

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
	EnableCache   bool
	MaxCacheSize  int
	Handlers      []func(m *Message)
	Commands      map[string]func(m *Message, args []string)
	Conn          *websocket.Conn
	ReadyHandlers         []func()
	UserJoinHandlers      []func(u *UserJoinPayload)
	UserLeaveHandlers     []func(l *UserLeavePayload)
	MessageDeleteHandlers []func(d *MessageDeletePayload)
	MessageUpdateHandlers []func(e *MessageUpdatePayload)
	ReactionAddHandlers   []func(e *MessageReactionPayload)
	messageCache  map[string]Message
	cacheOrder    []string
	cacheMutex    sync.RWMutex
}

func NewBot(token, prefix string) *Bot {
	return &Bot{
		Token:         token,
		CommandPrefix: prefix,
		BaseURL:       DefaultBaseURL,
		GatewayURL:    DefaultGatewayURL,
		Commands:      make(map[string]func(m *Message, args []string)),
		messageCache:  make(map[string]Message),
		cacheOrder:    make([]string, 0),
		EnableCache:   false,
	}
}
func (b *Bot) NewBotInstance(apiBase, gatewayURL string) {
	b.BaseURL = apiBase
	b.GatewayURL = gatewayURL
}
func (b *Bot) NewBotConfig(enableCache bool, maxCacheSize int) {
	b.EnableCache = enableCache
	if maxCacheSize <= 0 {
		b.MaxCacheSize = 500
	} else {
		b.MaxCacheSize = maxCacheSize
	}
}



func (b *Bot) OnMessage(handler func(m *Message)) {
	b.Handlers = append(b.Handlers, handler)
}
func (b *Bot) AddCommand(name string, handler func(m *Message, args []string)) {
	b.Commands[name] = handler
}

func (b *Bot) OnReady(handler func()) {
	b.ReadyHandlers = append(b.ReadyHandlers, handler)
}

func (b *Bot) OnUserJoin(handler func(u *UserJoinPayload)) {
	b.UserJoinHandlers = append(b.UserJoinHandlers, handler)
}
func (b *Bot) OnUserLeave(handler func(l *UserLeavePayload)) {
	b.UserLeaveHandlers = append(b.UserLeaveHandlers, handler)
}
func (b *Bot) OnMessageDelete(handler func(d *MessageDeletePayload)) {
	b.MessageDeleteHandlers = append(b.MessageDeleteHandlers, handler)
}
func (b *Bot) OnMessageEdit(handler func(e *MessageUpdatePayload)) {
	b.MessageUpdateHandlers = append(b.MessageUpdateHandlers, handler)
}
func (b *Bot) OnMessageReact(handler func(e *MessageReactionPayload)) {
	b.ReactionAddHandlers = append(b.ReactionAddHandlers, handler)
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
func (b *Bot) ReplyMessage(m *Message, content string) {
	payload := map[string]interface{}{
		"content": content,
		"message_reference": map[string]string{
			"message_id": m.ID,
			"channel_id": m.ChannelID,
		},
	}

	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/channels/%s/messages", b.BaseURL, m.ChannelID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
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

func (b *Bot) AddReaction(channelID, messageID, emoji string) {
	url := fmt.Sprintf("%s/channels/%s/messages/%s/reactions/%s", b.BaseURL, channelID, messageID, emoji)
	req, _ := http.NewRequest("PUT", url, nil)
	req.Header.Set("Authorization", "Bot "+b.Token)
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		b.checkRateLimit(resp.StatusCode)
		resp.Body.Close()
	}
}

func (b *Bot) ForwardMessage(targetChannelID string, m *Message) {
	payload := map[string]interface{}{
		"message_reference": map[string]string{
			"channel_id": m.ChannelID,
			"message_id": m.ID,
			"guild_id":   m.GuildID,
			"type": "1",
		},
		"message_snapshots": []map[string]interface{}{
			{
				"content":   m.Content,
				"author_id": m.Author.ID,
			},
		},
	}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/channels/%s/messages", b.BaseURL, targetChannelID)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err == nil {
		b.checkRateLimit(resp.StatusCode)
		resp.Body.Close()
	}
}

func (b *Bot) GetGuild(guildID string) (*GuildInfo, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/guilds/%s", b.BaseURL, guildID), nil)
	req.Header.Set("Authorization", "Bot "+b.Token)
	
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)

	var info GuildInfo
	json.NewDecoder(resp.Body).Decode(&info)
	return &info, nil
}

func (b *Bot) AddRole(guildID, userID, roleID string) error {
	url := fmt.Sprintf("%s/guilds/%s/members/%s/roles/%s", b.BaseURL, guildID, userID, roleID)
	req, _ := http.NewRequest("PUT", url, nil)
	req.Header.Set("Authorization", "Bot "+b.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *Bot) RemoveRole(guildID, userID, roleID string) error {
	url := fmt.Sprintf("%s/guilds/%s/members/%s/roles/%s", b.BaseURL, guildID, userID, roleID)
	req, _ := http.NewRequest("DELETE", url, nil)
	req.Header.Set("Authorization", "Bot "+b.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	return nil
}

func (b *Bot) TimeoutMember(guildID, userID string, durationSeconds int, reason string) error {
	expiration := time.Now().Add(time.Duration(durationSeconds) * time.Second).Format(time.RFC3339)
	payload := map[string]interface{}{
		"communication_disabled_until": expiration,
	}
	if reason != "" {
		payload["timeout_reason"] = reason
	}
	body, _ := json.Marshal(payload)
	url := fmt.Sprintf("%s/guilds/%s/members/%s", b.BaseURL, guildID, userID)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	return nil
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
			switch payload.T {
			case "READY":
				for _, h := range b.ReadyHandlers {
					go h()
				}
			case "MESSAGE_CREATE":
				var m Message
				json.Unmarshal(payload.D, &m)

				if m.Author.Bot {
					continue
				}
				
				if b.EnableCache {
					b.cacheMutex.Lock()
					if len(b.cacheOrder) >= b.MaxCacheSize {
						oldestID := b.cacheOrder[0]
						b.cacheOrder = b.cacheOrder[1:]
						delete(b.messageCache, oldestID)
					}
					b.messageCache[m.ID] = m
					b.cacheOrder = append(b.cacheOrder, m.ID)
					b.cacheMutex.Unlock()
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
			case "GUILD_MEMBER_ADD":
				var uj UserJoinPayload
				json.Unmarshal(payload.D, &uj)
				for _, h := range b.UserJoinHandlers {
					go h(&uj)
				}
			case "GUILD_MEMBER_REMOVE":
				var ul UserLeavePayload
				json.Unmarshal(payload.D, &ul)
				ul.UserID = ul.User.ID

				req, _ := http.NewRequest("GET", fmt.Sprintf("%s/guilds/%s", b.BaseURL, ul.GuildID), nil)
				req.Header.Set("Authorization", "Bot "+b.Token)
				if resp, err := http.DefaultClient.Do(req); err == nil {
					var guild struct{ Name string `json:"name"` }
					json.NewDecoder(resp.Body).Decode(&guild)
					ul.GuildName = guild.Name
					resp.Body.Close()
				}
				for _, h := range b.UserLeaveHandlers {
					go h(&ul)
				}
			case "MESSAGE_DELETE":
				var md MessageDeletePayload
				json.Unmarshal(payload.D, &md)

				b.cacheMutex.RLock()
				cachedMsg, exists := b.messageCache[md.MessageID]
				b.cacheMutex.RUnlock()
				if exists {
					md.CachedContent = cachedMsg.Content
					md.Author = cachedMsg.Author
					b.cacheMutex.Lock()
					delete(b.messageCache, md.MessageID)
					b.cacheMutex.Unlock()
				}

				for _, h := range b.MessageDeleteHandlers {
					go h(&md)
				}
			case "MESSAGE_UPDATE":
				var mu MessageUpdatePayload
				json.Unmarshal(payload.D, &mu)
				if b.EnableCache {
					b.cacheMutex.RLock()
					cachedMsg, exists := b.messageCache[mu.MessageID]
					b.cacheMutex.RUnlock()
					if exists {
						mu.OldContent = cachedMsg.Content
						b.cacheMutex.Lock()
						cachedMsg.Content = mu.NewContent
						b.messageCache[mu.MessageID] = cachedMsg
						b.cacheMutex.Unlock()
					}
				}

				for _, h := range b.MessageUpdateHandlers {
					go h(&mu)
				}
			case "MESSAGE_REACTION_ADD":
				var mr MessageReactionPayload
				json.Unmarshal(payload.D, &mr)
				for _, h := range b.ReactionAddHandlers {
					go h(&mr)
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