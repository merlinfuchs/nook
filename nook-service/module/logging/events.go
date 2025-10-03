package logging

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/module"
)

func (m *LoggingModule) handleChannelCreate(c module.Context[LoggingConfig], e *events.GuildChannelCreate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.ChannelUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventChannelCreate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.ChannelID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Channel Created").
		Timestamp(time.Now()).
		Field("Channel", fmt.Sprintf("<#%d>", e.ChannelID), true)

	if auditLogEntry != nil {
		resp.Field("Created By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleChannelUpdate(c module.Context[LoggingConfig], e *events.GuildChannelUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.ChannelUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventChannelUpdate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.ChannelID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Channel Updated").
		Timestamp(time.Now()).
		Field("Channel", fmt.Sprintf("<#%d>", e.ChannelID), true)

	if auditLogEntry != nil {
		resp.Field("Updated By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	if e.Channel.Name() != e.OldChannel.Name() {
		resp.Field("Name", fmt.Sprintf("`%s` -> `%s`", e.OldChannel.Name(), e.Channel.Name()), true)
	}

	if e.Channel.Position() != e.OldChannel.Position() {
		resp.Field("Position", fmt.Sprintf("`%d` -> `%d`", e.OldChannel.Position(), e.Channel.Position()), true)
	}

	if e.Channel.Type() != e.OldChannel.Type() {
		resp.Field("Type", fmt.Sprintf("`%d` -> `%d`", e.OldChannel.Type(), e.Channel.Type()), true)
	}

	if newChannel, ok := e.Channel.(discord.GuildMessageChannel); ok {
		oldChannel, ok := e.OldChannel.(discord.GuildMessageChannel)
		if !ok {
			return nil
		}

		if derefOr(newChannel.Topic(), "") != derefOr(oldChannel.Topic(), "") {
			resp.Field("Topic", fmt.Sprintf("`%s` -> `%s`", derefOr(oldChannel.Topic(), "''"), derefOr(newChannel.Topic(), "''")), true)
		}
	}

	if newChannel, ok := e.Channel.(discord.GuildAudioChannel); ok {
		oldChannel, ok := e.OldChannel.(discord.GuildAudioChannel)
		if !ok {
			return nil
		}

		if newChannel.Bitrate() != oldChannel.Bitrate() {
			resp.Field("Bitrate", fmt.Sprintf("`%d` -> `%d`", oldChannel.Bitrate(), newChannel.Bitrate()), true)
		}

		if newChannel.RTCRegion() != oldChannel.RTCRegion() {
			resp.Field("RTC Region", fmt.Sprintf("`%s` -> `%s`", def(oldChannel.RTCRegion(), "''"), def(newChannel.RTCRegion(), "''")), true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleChannelDelete(c module.Context[LoggingConfig], e *events.GuildChannelDelete) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.ChannelUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventChannelDelete &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.ChannelID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Channel Deleted").
		Timestamp(time.Now()).
		Field("Channel", fmt.Sprintf("%s (%d)", e.Channel.Name(), e.ChannelID), true)

	if auditLogEntry != nil {
		resp.Field("Deleted By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleRoleCreate(c module.Context[LoggingConfig], e *events.RoleCreate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.RoleUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventRoleCreate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.RoleID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Role Created").
		Timestamp(time.Now()).
		Field("Role", fmt.Sprintf("<@&%d>", e.RoleID), true)

	if auditLogEntry != nil {
		resp.Field("Created By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleRoleUpdate(c module.Context[LoggingConfig], e *events.RoleUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.RoleUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventRoleUpdate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.RoleID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Role Updated").
		Timestamp(time.Now()).
		Field("Role", fmt.Sprintf("<@&%d>", e.RoleID), true)

	if auditLogEntry != nil {
		resp.Field("Updated By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	if e.Role.Name != e.OldRole.Name {
		resp.Field("Name", fmt.Sprintf("`%s` -> `%s`", e.OldRole.Name, e.Role.Name), true)
	}

	if e.Role.Permissions != e.OldRole.Permissions {
		resp.Field("Permissions", fmt.Sprintf("`%s` -> `%s`", e.OldRole.Permissions, e.Role.Permissions), true)
	}

	if e.Role.Position != e.OldRole.Position {
		resp.Field("Position", fmt.Sprintf("`%d` -> `%d`", e.OldRole.Position, e.Role.Position), true)
	}

	if e.Role.Color != e.OldRole.Color {
		resp.Field("Color", fmt.Sprintf("`#%06X` -> `#%06X`", e.OldRole.Color, e.Role.Color), true)
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleRoleDelete(c module.Context[LoggingConfig], e *events.RoleDelete) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.RoleUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventRoleDelete &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.RoleID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Role Deleted").
		Timestamp(time.Now()).
		Field("Role", fmt.Sprintf("%s (%d)", e.Role.Name, e.RoleID), true)

	if auditLogEntry != nil {
		resp.Field("Deleted By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleGuildUpdate(c module.Context[LoggingConfig], e *events.GuildUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.GuildUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventGuildUpdate
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Guild Updated").
		Timestamp(time.Now()).
		Field("Guild", fmt.Sprintf("%s (%d)", e.Guild.Name, e.GuildID), true)

	if auditLogEntry != nil {
		resp.Field("Updated By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	if e.Guild.Name != e.OldGuild.Name {
		resp.Field("Name", fmt.Sprintf("`%s` -> `%s`", e.OldGuild.Name, e.Guild.Name), true)
	}

	if e.Guild.Description != e.OldGuild.Description {
		resp.Field("Description", fmt.Sprintf("`%s` -> `%s`", derefOr(e.OldGuild.Description, "''"), derefOr(e.Guild.Description, "''")), true)
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleEmojiCreate(c module.Context[LoggingConfig], e *events.EmojiCreate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.EmojiUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventEmojiCreate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Emoji.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Emoji Created").
		Timestamp(time.Now()).
		Field("Emoji", e.Emoji.Mention(), true).
		Field("Name", e.Emoji.Name, true).
		Thumbnail(e.Emoji.URL())

	if auditLogEntry != nil {
		resp.Field("Created By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleEmojiUpdate(c module.Context[LoggingConfig], e *events.EmojiUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.EmojiUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventEmojiUpdate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Emoji.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Emoji Updated").
		Timestamp(time.Now()).
		Field("Emoji", e.Emoji.Mention(), true).
		Field("Name", e.Emoji.Name, true).
		Thumbnail(e.Emoji.URL())

	if auditLogEntry != nil {
		resp.Field("Updated By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	if e.Emoji.Name != e.OldEmoji.Name {
		resp.Field("Name", fmt.Sprintf("`%s` -> `%s`", e.OldEmoji.Name, e.Emoji.Name), true)
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleEmojiDelete(c module.Context[LoggingConfig], e *events.EmojiDelete) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.StickerUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventEmojiDelete &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Emoji.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Emoji Deleted").
		Timestamp(time.Now()).
		Field("Emoji", fmt.Sprintf("%s (%d)", e.Emoji.Name, e.Emoji.ID), true).
		Thumbnail(e.Emoji.URL())

	if auditLogEntry != nil {
		resp.Field("Deleted By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleStickerCreate(c module.Context[LoggingConfig], e *events.StickerCreate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.StickerUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventStickerCreate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Sticker.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Sticker Created").
		Timestamp(time.Now()).
		Field("Sticker", fmt.Sprintf("%s (%d)", e.Sticker.Name, e.Sticker.ID), true).
		Thumbnail(e.Sticker.URL())

	if auditLogEntry != nil {
		resp.Field("Created By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleStickerUpdate(c module.Context[LoggingConfig], e *events.StickerUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.StickerUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventStickerUpdate &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Sticker.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Sticker Updated").
		Timestamp(time.Now()).
		Field("Sticker", fmt.Sprintf("%s (%d)", e.Sticker.Name, e.Sticker.ID), true).
		Thumbnail(e.Sticker.URL())

	if auditLogEntry != nil {
		resp.Field("Updated By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	if e.Sticker.Name != e.OldSticker.Name {
		resp.Field("Name", fmt.Sprintf("`%s` -> `%s`", e.OldSticker.Name, e.Sticker.Name), true)
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleStickerDelete(c module.Context[LoggingConfig], e *events.StickerDelete) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.GuildEventsChannelID()
	if !config.GuildEvents.StickerUpdateEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventStickerDelete &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.Sticker.ID
		},
	)
	if err != nil {
		return err
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title("Sticker Deleted").
		Timestamp(time.Now()).
		Field("Sticker", fmt.Sprintf("%s (%d)", e.Sticker.Name, e.Sticker.ID), true).
		Thumbnail(e.Sticker.URL())

	if auditLogEntry != nil {
		resp.Field("Deleted By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleMessageCreate(c module.Context[LoggingConfig], e *events.GuildMessageCreate) error {
	m.messageCache.Set(messageCacheKey{
		channelID: e.ChannelID,
		messageID: e.MessageID,
	}, e.Message, 0)
	return nil
}

func (m *LoggingModule) handleMessageDelete(c module.Context[LoggingConfig], e *events.GuildMessageDelete) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.MessageChannelID()
	if !config.MessageEvents.MessageDeleteEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				ae.AuditLogEntry.ActionType == discord.AuditLogEventMessageDelete &&
				ae.AuditLogEntry.Options != nil &&
				ae.AuditLogEntry.Options.ChannelID != nil &&
				*ae.AuditLogEntry.Options.ChannelID == e.ChannelID
		},
	)
	if err != nil {
		return err
	}

	cachedMsg, msgInCache := m.messageCache.GetAndDelete(messageCacheKey{
		channelID: e.ChannelID,
		messageID: e.MessageID,
	})

	resp := module.FormatMessage(c, e.GuildID).
		Title("Message Deleted").
		Timestamp(time.Now()).
		Field("Channel", fmt.Sprintf("<#%d>", e.ChannelID), true)

	if msgInCache {
		msg := cachedMsg.Value()
		resp.Field("Author", msg.Author.Mention(), true)
		resp.Field("Created At", discord.FormattedTimestampMention(
			msg.CreatedAt.Unix(),
			discord.TimestampStyleLongDateTime,
		), true)
		if msg.Content != "" {
			resp = resp.Field("Content", "```"+truncate(msg.Content, 950)+"```")
		}
	} else {
		resp.Field("Created At", discord.FormattedTimestampMention(
			e.MessageID.Time().Unix(),
			discord.TimestampStyleLongDateTime,
		), true)
		resp = resp.Field("Message ID", e.MessageID.String(), true)
	}

	if auditLogEntry != nil {
		resp.Field("Deleted By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleMessageEdit(c module.Context[LoggingConfig], e *events.GuildMessageUpdate) error {
	defer func() {
		m.messageCache.Set(messageCacheKey{
			channelID: e.ChannelID,
			messageID: e.MessageID,
		}, e.Message, 0)
	}()

	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.MessageChannelID()
	if !config.MessageEvents.MessageEditEvents || channelID == 0 {
		return nil
	}

	cachedMsg, msgInCache := m.messageCache.GetAndDelete(messageCacheKey{
		channelID: e.ChannelID,
		messageID: e.MessageID,
	})

	resp := module.FormatMessage(c, e.GuildID).
		Title("Message Edited").
		Timestamp(time.Now()).
		Field("Channel", fmt.Sprintf("<#%d>", e.ChannelID), true)

	if msgInCache {
		msg := cachedMsg.Value()
		resp.Field("Author", msg.Author.Mention(), true)
		resp.Field("Created At", discord.FormattedTimestampMention(
			msg.CreatedAt.Unix(),
			discord.TimestampStyleLongDateTime,
		), true)
		if msg.Content != "" {
			resp = resp.Field("Old Content", "```"+truncate(msg.Content, 950)+"```")
		}
	} else {
		resp.Field("Created At", discord.FormattedTimestampMention(
			e.MessageID.Time().Unix(),
			discord.TimestampStyleLongDateTime,
		), true)
		resp = resp.Field("Message ID", e.MessageID.String(), true)
	}

	resp.Field("New Content", "```"+truncate(e.Message.Content, 950)+"```")

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleMemberJoin(c module.Context[LoggingConfig], e *events.GuildMemberJoin) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.MemberEventsChannelID()
	if !config.MemberEvents.MemberJoinEvents || channelID == 0 {
		return nil
	}

	displayName := e.Member.User.Username
	if e.Member.User.GlobalName != nil {
		displayName = *e.Member.User.GlobalName
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title(displayName+" joined").
		Timestamp(time.Now()).
		Field("User", e.Member.Mention(), true).
		Field("Created At", discord.FormattedTimestampMention(
			e.Member.CreatedAt().Unix(),
			discord.TimestampStyleLongDateTime,
		), true).
		Thumbnail(e.Member.User.EffectiveAvatarURL())

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func (m *LoggingModule) handleMemberLeave(c module.Context[LoggingConfig], e *events.GuildMemberLeave) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	channelID := config.MemberEventsChannelID()
	if !config.MemberEvents.MemberJoinEvents || channelID == 0 {
		return nil
	}

	auditLogEntry, err := m.auditLogCollector.WaitForAuditLogEntry(
		c.Context(),
		func(ae *events.GuildAuditLogEntryCreate) bool {
			return ae.GuildID == e.GuildID &&
				(ae.AuditLogEntry.ActionType == discord.AuditLogEventMemberBanAdd ||
					ae.AuditLogEntry.ActionType == discord.AuditLogEventMemberKick) &&
				ae.AuditLogEntry.TargetID != nil &&
				*ae.AuditLogEntry.TargetID == e.User.ID
		},
	)
	if err != nil {
		return err
	}

	displayName := e.User.Username
	if e.User.GlobalName != nil {
		displayName = *e.User.GlobalName
	}

	action := "left"
	if auditLogEntry != nil {
		if auditLogEntry.ActionType == discord.AuditLogEventMemberBanAdd {
			action = "banned"
		} else {
			action = "kicked"
		}
	}

	resp := module.FormatMessage(c, e.GuildID).
		Title(displayName+" "+action).
		Timestamp(time.Now()).
		Field("User", e.User.Mention(), true).
		Field("Created At", discord.FormattedTimestampMention(
			e.User.CreatedAt().Unix(),
			discord.TimestampStyleLongDateTime,
		), true).
		Thumbnail(e.User.EffectiveAvatarURL())

	if !e.Member.JoinedAt.IsZero() {
		resp.Field("Joined At", discord.FormattedTimestampMention(
			e.Member.JoinedAt.Unix(),
			discord.TimestampStyleLongDateTime,
		), true)
	}

	if auditLogEntry != nil {
		if auditLogEntry.ActionType == discord.AuditLogEventMemberBanAdd {
			resp.Field("Banned By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		} else {
			resp.Field("Kicked By", fmt.Sprintf("<@%d>", auditLogEntry.UserID), true)
		}
		if auditLogEntry.Reason != nil {
			resp.Field("Reason", *auditLogEntry.Reason, true)
		}
	}

	_, err = c.Rest().CreateMessage(channelID, resp.BuildMessageCreate(), rest.WithCtx(c.Context()))
	if err != nil {
		return err
	}

	return nil
}

func truncate(s string, maxLength int) string {
	if len(s) > maxLength {
		return s[:maxLength] + "..."
	}
	return s
}

func derefOr[T comparable](v *T, def T) T {
	var zero T
	if v == nil || *v == zero {
		return def
	}
	return *v
}

func def[T comparable](v T, def T) T {
	var zero T
	if v == zero {
		return def
	}
	return v
}
