package common

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

func EventGuildID(e bot.Event) (*ID, bool) {
	switch e := e.(type) {
	case *events.Ready:
		return nil, true
	case *events.GuildReady:
		return &e.GuildID, true
	case *events.GuildUpdate:
		return &e.GuildID, true
	case *events.GuildLeave:
		return &e.GuildID, true
	case *events.GuildUnavailable:
		return &e.GuildID, true
	case *events.MessageCreate:
		return e.GuildID, false
	case *events.MessageUpdate:
		return e.GuildID, false
	case *events.MessageDelete:
		return e.GuildID, false
	case *events.GuildMessageCreate:
		return &e.GuildID, false
	case *events.GuildMessageDelete:
		return &e.GuildID, false
	case *events.GuildMessageUpdate:
		return &e.GuildID, false
	case *events.MessageReactionAdd:
		return e.GuildID, false
	case *events.MessageReactionRemove:
		return e.GuildID, false
	case *events.ApplicationCommandInteractionCreate:
		return e.GuildID(), false
	case *events.ComponentInteractionCreate:
		return e.GuildID(), false
	case *events.ModalSubmitInteractionCreate:
		return e.GuildID(), false
	case *events.AutocompleteInteractionCreate:
		return e.GuildID(), false
	case *events.GuildAuditLogEntryCreate:
		return &e.GuildID, false
	case *events.GuildMemberJoin:
		return &e.GuildID, false
	case *events.GuildMemberLeave:
		return &e.GuildID, false
	case *events.GuildMemberUpdate:
		return &e.GuildID, false
	case *events.GuildChannelCreate:
		return &e.GuildID, false
	case *events.GuildChannelUpdate:
		return &e.GuildID, false
	case *events.GuildChannelDelete:
		return &e.GuildID, false
	case *events.RoleCreate:
		return &e.GuildID, false
	case *events.RoleUpdate:
		return &e.GuildID, false
	case *events.RoleDelete:
		return &e.GuildID, false
	case *events.EmojiCreate:
		return &e.GuildID, false
	case *events.EmojiUpdate:
		return &e.GuildID, false
	case *events.EmojiDelete:
		return &e.GuildID, false
	case *events.StickerCreate:
		return &e.GuildID, false
	case *events.StickerUpdate:
		return &e.GuildID, false
	case *events.StickerDelete:
		return &e.GuildID, false
	}
	return nil, false
}
