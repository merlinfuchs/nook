package moderation

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/module"
)

const ModuleID = "moderation"

type ModerationModule struct {
}

func NewModerationModule() *ModerationModule {
	return &ModerationModule{}
}

func (m *ModerationModule) ModuleID() string {
	return ModuleID
}

func (m *ModerationModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Moderation",
		Description:    "Moderate your server with ease like banning, kicking, muting, etc.",
		Icon:           "shield",
		Internal:       false,
		DefaultEnabled: false,
		DefaultConfig:  defaultConfig,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
	}
}

func (m *ModerationModule) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ban",
			Description: "Ban a user from the server",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionUser{
					Name:        "user",
					Description: "The user to ban",
					Required:    true,
				},
				discord.ApplicationCommandOptionString{
					Name:        "reason",
					Description: "The reason for the ban",
					Required:    true,
				},
				discord.ApplicationCommandOptionInt{
					Name:        "purge_message_days",
					Description: "The number of days to purge messages by the user for",
					Required:    false,
				},
			},
		},
	}
}

func (m *ModerationModule) Router() module.Router[ModerationConfig] {
	return module.NewRouter[ModerationConfig]().
		Command("ban", m.handleBanCommand)
}
