package welcome

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(WelcomeConfig{})

var configUISchema = module.ConfigUISchema{
	Properties: map[string]module.ConfigUISchema{
		"welcome_channel_id": {
			Widget: module.ConfigUIWidgetChannelSelect,
			ChannelTypes: []discord.ChannelType{
				discord.ChannelTypeGuildText,
				discord.ChannelTypeGuildNews,
			},
		},
		"welcome_message_format": {
			Widget: module.ConfigUIWidgetMessageFormat,
		},
	},
}

var defaultConfig = WelcomeConfig{
	WelcomeMessageFormat: module.SerializedMessageFormat{
		Title:       "Welcome to the server!",
		Description: "Welcome to the server! Please read the rules and have a great time!",
		URL:         "https://example.com",
	}.MustMarshalString(),
}

type WelcomeConfig struct {
	WelcomeChannelID     common.ID `json:"welcome_channel_id,omitzero" title:"Welcome Channel" description:"The channel to send the welcome message to"`
	WelcomeMessageFormat string    `json:"welcome_message_format,omitzero" title:"Welcome Message" description:"The message to send when a user joins the server"`
}
