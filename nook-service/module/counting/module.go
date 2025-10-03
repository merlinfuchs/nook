package counting

import (
	"fmt"
	"strconv"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

const ModuleID = "counting"

type CountingModule struct {
}

func NewCountingModule() *CountingModule {
	return &CountingModule{}
}

func (m *CountingModule) ModuleID() string {
	return ModuleID
}

func (m *CountingModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Counting",
		Description:    "Create counting channels where users can try to count up.",
		Icon:           "calculator",
		Internal:       false,
		DefaultEnabled: false,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
	}
}

func (m *CountingModule) Router() module.Router[CountingConfig] {
	return module.NewRouter[CountingConfig]().
		Handle(module.ListenerFunc(m.handleMessageCreate))
}

func (m *CountingModule) handleMessageCreate(c module.Context[CountingConfig], e *events.MessageCreate) error {
	if e.GuildID == nil {
		return nil
	}
	guildID := *e.GuildID

	config, err := c.Config(guildID)
	if err != nil {
		return fmt.Errorf("failed to get module config: %w", err)
	}

	channelSettings := config.ChannelSettings(e.ChannelID)
	if !channelSettings.Enabled {
		return nil
	}

	num, err := strconv.Atoi(e.Message.Content)
	if err != nil {
		return fmt.Errorf("failed to convert message content to number: %w", err)
	}

	newValue, err := c.KV().Update(
		c.Context(),
		guildID,
		thing.OperationIncrement,
		channelCountKey(e.ChannelID),
		thing.NewInt(1),
	)
	if err != nil {
		return fmt.Errorf("failed to set module value: %w", err)
	}

	if newValue.Int() != int64(num) {
		err = c.KV().Delete(c.Context(), guildID, channelCountKey(e.ChannelID))
		if err != nil {
			return fmt.Errorf("failed to delete module value: %w", err)
		}

		err = c.Rest().AddReaction(e.ChannelID, e.Message.ID, "❌", rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to add reaction: %w", err)
		}

		_, err = c.Rest().CreateMessage(e.ChannelID, discord.NewMessageCreateBuilder().
			SetContent(fmt.Sprintf("%s messed it up.", e.Message.Author.Mention())).
			Build(),
			rest.WithCtx(c.Context()),
		)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}
		return nil
	}

	err = c.Rest().AddReaction(e.ChannelID, e.Message.ID, "✅", rest.WithCtx(c.Context()))
	if err != nil {
		return fmt.Errorf("failed to add reaction: %w", err)
	}

	return nil
}

func channelCountKey(channelID common.ID) string {
	return fmt.Sprintf("count:%s", channelID)
}
