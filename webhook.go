package gofluxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebhookClient struct {
	ID    string
	Token string
	URL   string
}
type WebhookPayload struct {
	Content   string `json:"content"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

func NewWebhookClient(id, token string) *WebhookClient {
	return &WebhookClient{
		ID:    id,
		Token: token,
		URL:   fmt.Sprintf("https://api.fluxer.app/v1/webhooks/%s/%s", id, token),
	}
}
type Webhook struct {
	ID        string `json:"id"`
	GuildID   string `json:"guild_id,omitempty"`
	ChannelID string `json:"channel_id"`
	Name      string `json:"name"`
	Avatar    string `json:"avatar,omitempty"`
	Token     string `json:"token,omitempty"`
}
type ModifyWebhookRequest struct {
	Name      string `json:"name,omitempty"`
	Avatar    string `json:"avatar,omitempty"`
	ChannelID string `json:"channel_id,omitempty"`
}
type EditWebhookMessageRequest struct {
	Content string  `json:"content,omitempty"`
	Embeds  []interface{} `json:"embeds,omitempty"`
}

func (w *WebhookClient) Execute(content string) error {
	payload := WebhookPayload{Content: content}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(w.URL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("[gofluxer]: Error executing webhook: %s", resp.Status)
	}
	return nil
}
func (b *Bot) GetWebhook(webhookID string) (*Webhook, error) {
	url := fmt.Sprintf("%s/webhooks/%s", b.BaseURL, webhookID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bot "+b.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	var webhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&webhook); err != nil {
		return nil, err
	}
	return &webhook, nil
}
func (b *Bot) DeleteWebhook(webhookID string) error {
	url := fmt.Sprintf("%s/webhooks/%s", b.BaseURL, webhookID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
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
func (w *WebhookClient) AltDeleteWebhook() error {
	req, err := http.NewRequest("DELETE", w.URL, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}
	return nil
}

func (b *Bot) EditWebhook(webhookID string, name, avatar, channelID string) (*Webhook, error) {
	url := fmt.Sprintf("%s/webhooks/%s", b.BaseURL, webhookID)
	payload := ModifyWebhookRequest{
		Name:      name,
		Avatar:    avatar,
		ChannelID: channelID,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bot "+b.Token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned with status code: %d", resp.StatusCode)
	}

	var updatedWebhook Webhook
	if err := json.NewDecoder(resp.Body).Decode(&updatedWebhook); err != nil {
		return nil, err
	}
	return &updatedWebhook, nil
}

func (b *Bot) WhEditMessage(webhookID, webhookToken, messageID, newContent string, embeds []interface{}) (*Message, error) {
	url := fmt.Sprintf("%s/webhooks/%s/%s/messages/%s", b.BaseURL, webhookID, webhookToken, messageID)
	payload := EditWebhookMessageRequest{
		Content: newContent,
		Embeds:  embeds,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("edit webhook message failed with status: %d", resp.StatusCode)
	}

	var editedMsg Message
	if err := json.NewDecoder(resp.Body).Decode(&editedMsg); err != nil {
		return nil, err
	}
	return &editedMsg, nil
}