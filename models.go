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