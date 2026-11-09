package gofluxer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const (
	PermissionCreateInstantInvite   int64 = 1 << 0
	PermissionKickMembers           int64 = 1 << 1
	PermissionBanMembers            int64 = 1 << 2
	PermissionAdministrator         int64 = 1 << 3
	PermissionManageChannels        int64 = 1 << 4
	PermissionManageGuild           int64 = 1 << 5
	PermissionAddReactions          int64 = 1 << 6
	PermissionViewAuditLog          int64 = 1 << 7
	PermissionPrioritySpeaker       int64 = 1 << 8
	PermissionStream                int64 = 1 << 9
	PermissionViewChannel           int64 = 1 << 10
	PermissionSendMessages          int64 = 1 << 11
	PermissionSendTTSMessages       int64 = 1 << 12
	PermissionManageMessages        int64 = 1 << 13
	PermissionEmbedLinks            int64 = 1 << 14
	PermissionAttachFiles           int64 = 1 << 15
	PermissionReadMessageHistory    int64 = 1 << 16
	PermissionMentionEveryone       int64 = 1 << 17
	PermissionUseExternalEmojis     int64 = 1 << 18
	PermissionViewGuildInsights     int64 = 1 << 19
	PermissionConnect               int64 = 1 << 20
	PermissionSpeak                 int64 = 1 << 21
	PermissionMuteMembers           int64 = 1 << 22
	PermissionDeafenMembers         int64 = 1 << 23
	PermissionMoveMembers           int64 = 1 << 24
	PermissionUseVAD                int64 = 1 << 25
	PermissionChangeNickname        int64 = 1 << 26
	PermissionManageNicknames       int64 = 1 << 27
	PermissionManageRoles           int64 = 1 << 28
	PermissionManageWebhooks        int64 = 1 << 29
	PermissionManageEmojisAndStickers int64 = 1 << 30
	PermissionUseApplicationCommands  int64 = 1 << 31
	PermissionRequestToSpeak        int64 = 1 << 32
	PermissionManageEvents          int64 = 1 << 33
	PermissionManageThreads         int64 = 1 << 34
	PermissionCreatePublicThreads   int64 = 1 << 35
	PermissionCreatePrivateThreads  int64 = 1 << 36
	PermissionUseExternalStickers   int64 = 1 << 37
	PermissionSendMessagesInThreads int64 = 1 << 38
	PermissionUseEmbeddedActivities int64 = 1 << 39
	PermissionModerateMembers       int64 = 1 << 40
)

func (b *Bot) HasPermission(m *Message, requiredPermission int64) bool {
	if m.GuildID == "" {
		return false
	}
	if b.IsOwner(m) {
		return true
	}

	url := fmt.Sprintf("%s/guilds/%s/roles", b.BaseURL, m.GuildID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req.Header.Set("Authorization", "Bot "+b.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	b.checkRateLimit(resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return false
	}

	var roles []Role
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return false
	}
	rolePermissions := make(map[string]int64)
	for _, r := range roles {
		permVal, _ := strconv.ParseInt(r.Permissions, 10, 64)
		rolePermissions[r.ID] = permVal
	}

	memberURL := fmt.Sprintf("%s/guilds/%s/members/%s", b.BaseURL, m.GuildID, m.Author.ID)
	memberReq, err := http.NewRequest("GET", memberURL, nil)
	if err != nil {
		return false
	}
	memberReq.Header.Set("Authorization", "Bot "+b.Token)
	memberResp, err := http.DefaultClient.Do(memberReq)
	if err != nil {
		return false
	}
	defer memberResp.Body.Close()
	b.checkRateLimit(memberResp.StatusCode)
	if memberResp.StatusCode != http.StatusOK {
		return false
	}

	var member struct {
		Roles []string `json:"roles"`
	}
	if err := json.NewDecoder(memberResp.Body).Decode(&member); err != nil {
		return false
	}

	userPermissions := rolePermissions[m.GuildID]
	for _, roleID := range member.Roles {
		if p, exists := rolePermissions[roleID]; exists {
			userPermissions |= p
		}
	}
	if (userPermissions & PermissionAdministrator) == PermissionAdministrator {
		return true
	}

	return (userPermissions & requiredPermission) == requiredPermission
}