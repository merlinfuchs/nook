package moderation

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/module"
)

func (m *ModerationModule) handleBanCommand(c module.Context[ModerationConfig], e *events.ApplicationCommandInteractionCreate) error {
	config, err := c.Config(*e.GuildID())
	if err != nil {
		return fmt.Errorf("failed to get moderation config: %w", err)
	}

	data := e.SlashCommandInteractionData()
	guild, ok := e.Guild()
	if !ok {
		return fmt.Errorf("guild not found")
	}

	user := data.User("user")
	reason := data.String("reason")
	purgeMessageDays := data.Int("purge_message_days")
	purgeDuration := config.Ban.DefaultPurgeDuration
	if purgeMessageDays > 0 {
		purgeDuration = time.Duration(purgeMessageDays) * 24 * time.Hour
	}

	err = c.Rest().AddBan(guild.ID, user.ID, purgeDuration, rest.WithReason(reason))
	if err != nil {
		return fmt.Errorf("failed to add ban: %w", err)
	}

	if config.Ban.NotifyUser {
		dmChannel, err := c.Rest().CreateDMChannel(user.ID, rest.WithCtx(c.Context()))
		if err != nil {
			return fmt.Errorf("failed to create DM: %w", err)
		}

		_, err = c.Rest().CreateMessage(dmChannel.ID(), discord.NewMessageCreateBuilder().
			SetContent(fmt.Sprintf(
				"You have been banned from %s (%d) for the following reason: %s",
				guild.Name,
				guild.ID,
				reason,
			)).
			Build(),
			rest.WithCtx(c.Context()),
		)
		if err != nil {
			return fmt.Errorf("failed to create message: %w", err)
		}
	}

	return nil
}
