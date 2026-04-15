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

func NewWebhookClient(id, token string) *WebhookClient {
	return &WebhookClient{
		ID:    id,
		Token: token,
		URL:   fmt.Sprintf("https://api.fluxer.app/v1/webhooks/%s/%s", id, token),
	}
}

type WebhookPayload struct {
	Content   string `json:"content"`
	Username  string `json:"username,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
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
		return fmt.Errorf("failed to execute webhook: %s", resp.Status)
	}

	return nil
}