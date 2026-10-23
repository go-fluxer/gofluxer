package gofluxer

type Message struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Author    User   `json:"author"`
	Pinned    bool   `json:"pinned"`
	Mention   bool   `json:"mention_everyone"`
	Flags     int    `json:"flags"`
}

type User struct {
	ID             string `json:"id"`
	Username       string `json:"username"`
	Discriminator  string `json:"discriminator"`
	Avatar         string `json:"avatar,omitempty"`
	Bot            bool   `json:"bot"`
	Flags          int    `json:"flags"`
	System         bool   `json:"system"`
}

type Channel struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    int    `json:"type"`
	GuildID string `json:"guild_id"`
	Topic   string `json:"topic"`
	Nsfw    bool   `json:"nsfw"`
}

type GuildInfo struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	OwnerID            string   `json:"owner_id"`
	MemberCount        int      `json:"member_count"`
	Icon               string   `json:"icon"`
	Banner             string   `json:"banner"`
	Permissions        string   `json:"permissions"`
	SysChannel         string   `json:"system_channel_id,omitempty"`
	AfkChannel         string   `json:"afk_channel_id,omitempty"`
	AfkTimeout         string   `json:"afk_timeout,omitempty"`
	Vanity             string   `json:"vanity_url_code,omitempty"`
	VerificationLevel  int      `json:"verification_level"`
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
type MessageUpdateRequest struct {
	Content string        `json:"content,omitempty"`
	Embeds  []interface{} `json:"embeds,omitempty"`
	Flags   int           `json:"flags,omitempty"`
}
type Sticker struct {
	ID          string `json:"id"`
	GuildID     string `json:"guild_id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Tags        string `json:"tags,omitempty"`
	Type        int    `json:"type,omitempty"`
	FormatType  int    `json:"format_type,omitempty"`
	Available   bool   `json:"available,omitempty"`
	User        *User  `json:"user,omitempty"`
}
type MessageReference struct {
	MessageID string `json:"message_id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id,omitempty"`
	Type      int    `json:"type,omitempty"`
}
type ForwardMessageRequest struct {
	MessageReference MessageReference `json:"message_reference"`
	Content          string           `json:"content,omitempty"`
}

type EditChannelRequest struct {
	Name                 string  `json:"name,omitempty"`
	Topic                *string `json:"topic,omitempty"`
	NSFW                 *bool   `json:"nsfw,omitempty"`
	RateLimitPerUser     *int    `json:"rate_limit_per_user,omitempty"`
	Position             *int    `json:"position,omitempty"`
	ParentID             *string `json:"parent_id,omitempty"`
	DefaultAutoArchive   *int    `json:"default_auto_archive_duration,omitempty"`
}

type StartThreadRequest struct {
	Name                string `json:"name"`
	AutoArchiveDuration int    `json:"auto_archive_duration,omitempty"`
	RateLimitPerUser    int    `json:"rate_limit_per_user,omitempty"`
	Type                int    `json:"type,omitempty"`
}
type ForumTag struct {
	ID        string `json:"id,omitempty"`
	Name      string `json:"name"`
	Moderated bool   `json:"moderated,omitempty"`
	EmojiID   string `json:"emoji_id,omitempty"`
	EmojiName string `json:"emoji_name,omitempty"`
}
type ForumPostMessage struct {
	Content string        `json:"content,omitempty"`
	Embeds  []interface{} `json:"embeds,omitempty"`
}
type StartForumThreadRequest struct {
	Name                string           `json:"name"`
	Message             ForumPostMessage `json:"message"`
	AppliedTags         []string         `json:"applied_tags,omitempty"`
	AutoArchiveDuration int              `json:"auto_archive_duration,omitempty"`
}

type PresenceUpdatePayload struct {
	Op int `json:"op"`
	D  struct {
		Since      *int       `json:"since"`
		Status     string     `json:"status"`
		AFK        bool       `json:"afk"`
	} `json:"d"`
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