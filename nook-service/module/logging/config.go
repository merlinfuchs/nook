package logging

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(LoggingConfig{})

var configUISchema = module.ConfigUISchema{
	Layout: module.ConfigLayout{
		Items: []string{"default_channel_id"},
		Children: []module.ConfigLayout{
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "guild_events",
				Items: []string{
					"guild_events.channel_id",
					"guild_events.channel_create_events",
					"guild_events.channel_update_events",
					"guild_events.channel_delete_events",
					"guild_events.role_create_events",
					"guild_events.role_update_events",
					"guild_events.role_delete_events",
					"guild_events.guild_update_events",
					"guild_events.emoji_update_events",
					"guild_events.sticker_update_events",
				},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "member_events",
				Items: []string{
					"member_events.channel_id",
					"member_events.member_join_events",
					"member_events.member_leave_events",
					"member_events.member_ban_events",
					"member_events.member_kick_events",
					"member_events.member_timeout_events",
				},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "message_events",
				Items: []string{
					"message_events.channel_id",
					"message_events.message_delete_events",
					"message_events.message_edit_events",
					"message_events.message_purge_events",
				},
			},
		},
	},
	Properties: map[string]module.ConfigUISchema{
		"default_channel_id": {
			Widget: module.ConfigUIWidgetChannelSelect,
			ChannelTypes: []discord.ChannelType{
				discord.ChannelTypeGuildText,
				discord.ChannelTypeGuildNews,
			},
		},
		"guild_events": {
			Properties: map[string]module.ConfigUISchema{
				"channel_id": {
					Widget: module.ConfigUIWidgetChannelSelect,
				},
			},
		},
		"member_events": {
			Properties: map[string]module.ConfigUISchema{
				"channel_id": {
					Widget: module.ConfigUIWidgetChannelSelect,
				},
			},
		},
		"message_events": {
			Properties: map[string]module.ConfigUISchema{
				"channel_id": {
					Widget: module.ConfigUIWidgetChannelSelect,
				},
			},
		},
	},
}

var defaultConfig = LoggingConfig{}

type LoggingConfig struct {
	DefaultChannelID common.ID `json:"default_channel_id,omitzero" title:"Default Channel" description:"The channel to log events to by default" required:"true"`

	GuildEvents   GuildEventsConfig   `json:"guild_events,omitzero" title:"Guild Events" description:"The events to log related to the server"`
	MemberEvents  MemberEventsConfig  `json:"member_events,omitzero" title:"Member Events" description:"The events to log related to members"`
	MessageEvents MessageEventsConfig `json:"message_events,omitzero" title:"Message Events" description:"The events to log related to messages"`
}

func (c LoggingConfig) GuildEventsChannelID() common.ID {
	if c.GuildEvents.ChannelID != 0 {
		return c.GuildEvents.ChannelID
	}
	return c.DefaultChannelID
}

func (c LoggingConfig) MemberEventsChannelID() common.ID {
	if c.MemberEvents.ChannelID != 0 {
		return c.MemberEvents.ChannelID
	}
	return c.DefaultChannelID
}

func (c LoggingConfig) MessageChannelID() common.ID {
	if c.MessageEvents.ChannelID != 0 {
		return c.MessageEvents.ChannelID
	}
	return c.DefaultChannelID
}

type GuildEventsConfig struct {
	ChannelID common.ID `json:"channel_id,omitzero" title:"Channel ID" description:"The channel ID to log server events to"`

	ChannelUpdateEvents bool `json:"channel_update_events,omitzero" title:"Channel Update Events" description:"Whether to log channel update events"`
	RoleUpdateEvents    bool `json:"role_update_events,omitzero" title:"Role Update Events" description:"Whether to log role update events"`
	GuildUpdateEvents   bool `json:"guild_update_events,omitzero" title:"Guild Update Events" description:"Whether to log guild update events"`
	EmojiUpdateEvents   bool `json:"emoji_update_events,omitzero" title:"Emoji Update Events" description:"Whether to log emoji update events"`
	StickerUpdateEvents bool `json:"sticker_update_events,omitzero" title:"Sticker Update Events" description:"Whether to log sticker update events"`
}

type MemberEventsConfig struct {
	ChannelID common.ID `json:"channel_id,omitzero" title:"Channel ID" description:"The channel ID to log member events to"`

	MemberJoinEvents  bool `json:"member_join_events,omitzero" title:"Member Join Events" description:"Whether to log member join events"`
	MemberLeaveEvents bool `json:"member_leave_events,omitzero" title:"Member Leave Events" description:"Whether to log member leave events"`
}

type MessageEventsConfig struct {
	ChannelID common.ID `json:"channel_id,omitzero" title:"Channel ID" description:"The channel ID to log message events to"`

	MessageDeleteEvents     bool `json:"message_delete_events,omitzero" title:"Message Delete Events" description:"Whether to log message delete events"`
	MessageEditEvents       bool `json:"message_edit_events,omitzero" title:"Message Edit Events" description:"Whether to log message edit events"`
	MessageBulkDeleteEvents bool `json:"message_bulk_delete_events,omitzero" title:"Message Bulk Delete Events" description:"Whether to log message bulk delete events"`
}
