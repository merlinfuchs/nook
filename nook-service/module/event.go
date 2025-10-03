package module

import (
	"encoding/json"

	"github.com/disgoorg/disgo/bot"
	"github.com/merlinfuchs/nook/nook-service/common"
)

type EventType string

const (
	EventTypeConfigure EventType = "configure"
	EventTypeDiscord   EventType = "discord"
)

type Event interface {
	Type() EventType
}

type ConfigureEvent struct {
	GuildID common.ID
	Enabled bool
	Config  json.RawMessage
}

func (e *ConfigureEvent) Type() EventType {
	return EventTypeConfigure
}

type DiscordEvent struct {
	bot.Event
}

func (e *DiscordEvent) Type() EventType {
	return EventTypeDiscord
}
