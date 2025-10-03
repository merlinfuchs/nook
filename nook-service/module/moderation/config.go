package moderation

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(ModerationConfig{})

var configUISchema = module.ConfigUISchema{
	Layout: module.ConfigLayout{
		Items: []string{"log_channel_id"},
		Children: []module.ConfigLayout{
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "ban",
				Items:  []string{"ban.default_purge_duration", "ban.notify_user"},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "tempban",
				Items:  []string{"tempban.default_duration", "tempban.default_purge_duration", "tempban.notify_user"},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "kick",
				Items:  []string{"kick.notify_user"},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "warn",
				Items:  []string{"warn.notify_user"},
			},
			{
				Type:   module.ConfigLayoutTypeContainer,
				Header: "timeout",
				Items:  []string{"timeout.default_duration", "timeout.notify_user"},
			},
		},
	},
	Properties: map[string]module.ConfigUISchema{
		"log_channel_id": {
			Widget: module.ConfigUIWidgetChannelSelect,
			ChannelTypes: []discord.ChannelType{
				discord.ChannelTypeGuildText,
				discord.ChannelTypeGuildNews,
			},
		},
		"ban": {
			Properties: map[string]module.ConfigUISchema{
				"default_purge_duration": {
					Widget: module.ConfigUIWidgetDurationDays,
				},
			},
		},
		"tempban": {
			Properties: map[string]module.ConfigUISchema{
				"default_duration": {
					Widget: module.ConfigUIWidgetDurationSeconds,
				},
				"default_purge_duration": {
					Widget: module.ConfigUIWidgetDurationDays,
				},
			},
		},
		"timeout": {
			Properties: map[string]module.ConfigUISchema{
				"default_duration": {
					Widget: module.ConfigUIWidgetDurationSeconds,
				},
			},
		},
	},
}

var defaultConfig = ModerationConfig{
	Ban: ModerationBanConfig{
		NotifyUser:           true,
		DefaultPurgeDuration: 2 * 24 * time.Hour,
	},
	Tempban: ModerationTempbanConfig{
		NotifyUser:           true,
		DefaultDuration:      7 * 24 * time.Hour,
		DefaultPurgeDuration: 2 * 24 * time.Hour,
	},
	Kick: ModerationKickConfig{
		NotifyUser: true,
	},
	Warn: ModerationWarnConfig{
		NotifyUser: true,
	},
	Timeout: ModerationTimeoutConfig{
		NotifyUser:      true,
		DefaultDuration: 1 * time.Hour,
	},
}

// ModerationConfig is the settings for the moderation module
type ModerationConfig struct {
	LogChannelID common.ID `json:"log_channel_id,omitzero" title:"Modlog Channel" description:"The channel in which moderation actions are logged"`

	Ban     ModerationBanConfig     `json:"ban,omitzero" title:"Ban" description:"Settings related to the ban command" required:"true"`
	Tempban ModerationTempbanConfig `json:"tempban,omitzero" title:"Tempban" description:"Settings related to the tempban command" required:"true"`
	Kick    ModerationKickConfig    `json:"kick,omitzero" title:"Kick" description:"Settings related to the kick command" required:"true"`
	Warn    ModerationWarnConfig    `json:"warn,omitzero" title:"Warn" description:"Settings related to the warn command" required:"true"`
	Timeout ModerationTimeoutConfig `json:"timeout,omitzero" title:"Timeout" description:"Settings related to the timeout command" required:"true"`
}

type ModerationBanConfig struct {
	NotifyUser           bool          `json:"notify_user" title:"Notify User" description:"Whether to DM the user when they are banned" required:"true"`
	DefaultPurgeDuration time.Duration `json:"default_purge_duration" title:"Default Purge Duration" description:"The number of days to purge the message for" required:"true"`
}

type ModerationTempbanConfig struct {
	NotifyUser           bool          `json:"notify_user" title:"Notify User" description:"Whether to DM the user when they are tempbanned" required:"true"`
	DefaultDuration      time.Duration `json:"default_duration" title:"Default Duration" description:"The duration of the tempban" required:"true"`
	DefaultPurgeDuration time.Duration `json:"default_purge_duration" title:"Default Purge Duration" description:"The number of days to purge the message for" required:"true"`
}

type ModerationKickConfig struct {
	NotifyUser bool `json:"notify_user" title:"Notify User" description:"Whether to DM the user when they are kicked" required:"true"`
}

type ModerationWarnConfig struct {
	NotifyUser bool `json:"notify_user" title:"Notify User" description:"Whether to DM the user when they are warned" required:"true"`
}

type ModerationTimeoutConfig struct {
	NotifyUser      bool          `json:"notify_user" title:"Notify User" description:"Whether to DM the user when they are timed out" required:"true"`
	DefaultDuration time.Duration `json:"default_duration" title:"Default Duration" description:"The duration of the timeout" required:"true"`
}
