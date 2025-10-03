package ping

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/module"
)

type PingModule struct {
}

func NewPingModule() *PingModule {
	return &PingModule{}
}

func (m *PingModule) ModuleID() string {
	return "ping"
}

func (m *PingModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Ping",
		Description:    "Provides a ping command that shows the latency of the bot.",
		Icon:           "wifi",
		Internal:       false,
		DefaultEnabled: true,
	}
}

func (m *PingModule) Commands() []discord.ApplicationCommandCreate {
	return []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "ping",
			Description: "Ping? Pong!",
			Options:     []discord.ApplicationCommandOption{},
		},
	}
}

func (m *PingModule) Router() module.Router[PingConfig] {
	return module.NewRouter[PingConfig]().
		Command("ping", m.handlePingCommand).
		Handle(module.ListenerFunc(m.handleMessageCreate))
}

func (m *PingModule) handlePingCommand(c module.Context[PingConfig], e *events.ApplicationCommandInteractionCreate) error {
	latency := time.Since(e.ApplicationCommandInteraction.CreatedAt())

	msg := module.FormatMessage(c, *e.GuildID()).
		Title("Pong!").
		Description(fmt.Sprintf("Latency: %s", latency)).
		BuildMessageCreate()

	err := e.CreateMessage(msg, rest.WithCtx(c.Context()))
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (m *PingModule) handleMessageCreate(c module.Context[PingConfig], e *events.GuildMessageCreate) error {
	latency := time.Since(e.Message.CreatedAt)

	guildSettings, err := c.GuildSettings(e.GuildID)
	if err != nil {
		return fmt.Errorf("failed to get guild settings: %w", err)
	}

	if e.Message.Content == guildSettings.CommandPrefix+"ping" {
		msg := module.FormatMessage(c, e.GuildID).
			Title("Pong!").
			Description(fmt.Sprintf("Latency: %s", latency)).
			BuildMessageCreate()

		_, err := c.Rest().CreateMessage(e.ChannelID, msg, rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	return nil
}
