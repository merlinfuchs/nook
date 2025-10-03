package guild

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module/user"
	"gopkg.in/guregu/null.v4"
)

type GuildWire struct {
	ID          common.ID           `json:"id"`
	Name        string              `json:"name"`
	Icon        null.String         `json:"icon"`
	Owner       bool                `json:"owner"`
	Permissions discord.Permissions `json:"permissions"`
	Bot         bool                `json:"bot"`
	Access      bool                `json:"access"`
}

type GuildGetResponseWire = GuildWire

type GuildListResponseWire = []GuildWire
type ChannelWire struct {
	ID       common.ID           `json:"id"`
	Type     discord.ChannelType `json:"type"`
	Name     string              `json:"name"`
	Position int                 `json:"position"`
}

type ChannelGetResponseWire = ChannelWire

type ChannelListResponseWire = []ChannelWire

type RoleWire struct {
	ID          common.ID           `json:"id"`
	Name        string              `json:"name"`
	Permissions discord.Permissions `json:"permissions"`
	Position    int                 `json:"position"`
	Color       string              `json:"color"`
}

type RoleGetResponseWire = RoleWire

type RoleListResponseWire = []RoleWire

type GuildSettingsWire struct {
	CommandPrefix null.String              `json:"command_prefix"`
	ColorScheme   null.String              `json:"color_scheme"`
	Default       DefaultGuildSettingsWire `json:"default"`
}

type DefaultGuildSettingsWire struct {
	CommandPrefix string `json:"command_prefix"`
	ColorScheme   string `json:"color_scheme"`

	ResolvedColor string `json:"resolved_color"`
}

type GuildSettingsGetResponseWire = GuildSettingsWire

type GuildSettingsUpdateRequestWire struct {
	CommandPrefix *string `json:"command_prefix"`
	ColorScheme   *string `json:"color_scheme"`
}

type GuildSettingsUpdateResponseWire = GuildSettingsWire

type GuildManagerWire struct {
	GuildID   common.ID              `json:"guild_id"`
	User      user.UserWire          `json:"user"`
	Role      model.GuildManagerRole `json:"role"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

type GuildManagerListResponseWire = []GuildManagerWire

type GuildManagerCreateAddWire struct {
	UserID common.ID              `json:"user_id"`
	Role   model.GuildManagerRole `json:"role"`
}

type GuildManagerAddAddResponseWire = GuildManagerWire

type GuildManagerRemoveResponseWire struct{}

type GuildProfileWire struct {
	UserID           common.ID `json:"user_id"`
	CustomName       string    `json:"custom_name,omitempty"`
	DefaultName      string    `json:"default_name"`
	CustomBio        string    `json:"custom_bio,omitempty"`
	DefaultBio       string    `json:"default_bio,omitempty"`
	CustomAvatarURL  string    `json:"custom_avatar_url,omitempty"`
	DefaultAvatarURL string    `json:"default_avatar_url,omitempty"`
}

type GuildProfileGetResponseWire = GuildProfileWire

type GuildProfileUpdateRequestWire struct {
	Name   string  `json:"name"`
	Bio    string  `json:"bio"`
	Avatar *string `json:"avatar"`
}

type GuildProfileUpdateResponseWire = GuildProfileWire
