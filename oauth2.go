package gofluxer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type OAuthClient struct {
	ClientID     string
	ClientSecret string
	RedirectURI  string
	BaseURL      string
}

func NewOAuth(clientID, clientSecret, redirectURI string) *OAuthClient {
	return &OAuthClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		BaseURL:      "https://api.fluxer.app/v1",
	}
}

func (o *OAuthClient) GetLoginURL(scopes []string) string {
	return fmt.Sprintf("https://fluxer.app/oauth2/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=%s", o.ClientID, url.QueryEscape(o.RedirectURI), url.QueryEscape(bytes.NewBufferString("").String()))
}

func (o *OAuthClient) AuthCode(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", o.ClientID)
	data.Set("client_secret", o.ClientSecret)
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", o.RedirectURI)

	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/oauth2/token", o.BaseURL), bytes.NewBufferString(data.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("[gofluxer]: Token exchange failed with status: %d", resp.StatusCode)
	}

	var tr TokenResponse
	json.NewDecoder(resp.Body).Decode(&tr)
	return &tr, nil
}

func (o *OAuthClient) GetUser(accessToken string) (*OAuthUser, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/@me", o.BaseURL), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user OAuthUser
	json.NewDecoder(resp.Body).Decode(&user)
	return &user, nil
}

func (o *OAuthClient) GetGuilds(accessToken string) ([]OAuthGuild, error) {
	req, _ := http.NewRequest("GET", fmt.Sprintf("%s/users/@me/guilds", o.BaseURL), nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	// fmt.Printf("[gofluxer]: Raw API JSON Response string: %s\n", buf.String())
	// Debug logging
	var guilds []OAuthGuild
	json.NewDecoder(&buf).Decode(&guilds)
	return guilds, nil
}