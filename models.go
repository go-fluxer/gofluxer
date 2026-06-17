package gofluxer

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Author    User   `json:"author"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
	Bot      bool   `json:"bot"`
}

type Channel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	GuildID string `json:"guild_id"`
}

type GuildInfo struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	OwnerID      string   `json:"owner_id"`
	MemberCount  int      `json:"member_count"`
	Icon         string   `json:"icon"`
	Banner       string   `json:"banner"`
	Permissions  string   `json:"permissions"`
	SysChannel   string   `json:"system_channel_id"`
	AfkChannel   string   `json:"afk_channel_id"`
	Vanity       string   `json:"vanity_url_code"`
}

type UserJoinPayload struct {
	GuildName  string `json:"guild_name"`
	GuildID    string `json:"guild_id"`
	UserID     string `json:"user_id"`
	User       User   `json:"user"`
	UserAvatar *string `json:"avatar,omitempty"`
}
type UserLeavePayload struct {
	GuildName string `json:"guild_name"`
	GuildID   string `json:"guild_id"`
	UserID    string `json:"user_id"`
	User      User   `json:"user"`
}
type MessageDeletePayload struct {
	MessageID     string `json:"id"`
	ChannelID     string `json:"channel_id"`
	GuildID       string `json:"guild_id"`
	Author        User   `json:"author_id,omitempty"`
	CachedContent string `json:"content,omitempty"`
}
type MessageUpdatePayload struct {
	MessageID      string `json:"id"`
	ChannelID      string `json:"channel_id"`
	GuildID        string `json:"guild_id"`
	OldContent     string `json:"old_content"`
	NewContent     string `json:"content"`
}
type MessageReactionPayload struct {
	UserID    string `json:"user_id"`
	ChannelID string `json:"channel_id"`
	MessageID string `json:"message_id"`
	GuildID   string `json:"guild_id"`
	Emoji     string `json:"emoji"`
}



type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type OAuthUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email,omitempty"`
	Avatar        string `json:"avatar"`
	MFAEnabled    bool   `json:"mfa_enabled"`
}

type OAuthGuild struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	OwnerID     string `json:"owner_id"`
	Permissions string `json:"permissions"`
}