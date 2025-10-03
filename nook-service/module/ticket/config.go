package ticket

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(TicketConfig{})

var configUISchema = module.ConfigUISchema{
	Properties: map[string]module.ConfigUISchema{
		"admin_role_id": {
			Widget: module.ConfigUIWidgetRoleSelect,
		},
		"log_channel_id": {
			Widget: module.ConfigUIWidgetChannelSelect,
			ChannelTypes: []discord.ChannelType{
				discord.ChannelTypeGuildText,
				discord.ChannelTypeGuildNews,
			},
		},
		"panel_channel_id": {
			Widget: module.ConfigUIWidgetChannelSelect,
			ChannelTypes: []discord.ChannelType{
				discord.ChannelTypeGuildText,
				discord.ChannelTypeGuildNews,
			},
		},
		"panel_message_data": {
			Widget: module.ConfigUIWidgetMessage,
		},
	},
}

// TicketConfig is the settings for the ticket module
type TicketConfig struct {
	AdminRoleID      common.ID `json:"admin_role_id" title:"Admin Role" description:"The ID role that can manage tickets" required:"true"`
	LogChannelID     common.ID `json:"log_channel_id" title:"Log Channel" description:"The channel in which the ticket logs are located" required:"true"`
	PanelChannelID   common.ID `json:"panel_channel_id" title:"Panel Channel" description:"The channel in which the ticket panel is located" required:"true"`
	PanelMessageData string    `json:"panel_message_data" title:"Panel Message" description:"The message that is sent when a ticket is created"`
}
