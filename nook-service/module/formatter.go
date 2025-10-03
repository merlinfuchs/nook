package module

import (
	"encoding/json"
	"log"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
)

type MessageIcon string

const (
	MessageIconError    MessageIcon = "error"
	MessageIconSuccess  MessageIcon = "success"
	MessageIconInfo     MessageIcon = "info"
	MessageIconWarning  MessageIcon = "warning"
	MessageIconQuestion MessageIcon = "question"
)

type MessageBuilder[C any] struct {
	c       Context[C]
	guildID common.ID

	title       string
	description string
	url         string
	icon        MessageIcon
	fields      []discord.EmbedField
	timestamp   time.Time
	thumbnail   string
}

func FormatMessage[C any](c Context[C], guildID common.ID) *MessageBuilder[C] {
	return &MessageBuilder[C]{
		c:       c,
		guildID: guildID,
	}
}

func (b *MessageBuilder[C]) Title(title string) *MessageBuilder[C] {
	b.title = title
	return b
}

func (b *MessageBuilder[C]) Description(description string) *MessageBuilder[C] {
	b.description = description
	return b
}

func (b *MessageBuilder[C]) URL(url string) *MessageBuilder[C] {
	b.url = url
	return b
}

func (b *MessageBuilder[C]) Icon(icon MessageIcon) *MessageBuilder[C] {
	b.icon = icon
	return b
}

func (b *MessageBuilder[C]) Field(name, value string, inline ...bool) *MessageBuilder[C] {
	inlineValue := len(inline) > 0 && inline[0]

	b.fields = append(b.fields, discord.EmbedField{
		Name:   name,
		Value:  value,
		Inline: &inlineValue,
	})
	return b
}

func (b *MessageBuilder[C]) Timestamp(timestamp time.Time) *MessageBuilder[C] {
	b.timestamp = timestamp
	return b
}

func (b *MessageBuilder[C]) Thumbnail(thumbnail string) *MessageBuilder[C] {
	b.thumbnail = thumbnail
	return b
}

func (b *MessageBuilder[C]) BuildMessageCreate() discord.MessageCreate {
	guildSettings, err := b.c.GuildSettings(b.guildID)
	if err != nil {
		slog.Error("Failed to get guild settings", "error", err)
	}

	embed := discord.NewEmbedBuilder().
		SetAuthorName(b.title).
		SetAuthorURL(b.url).
		SetDescription(b.description).
		SetAuthorIcon(string(b.icon)).
		SetColor(guildSettings.Color()).
		SetFields(b.fields...)

	if b.thumbnail != "" {
		embed.SetThumbnail(b.thumbnail)
	}

	if !b.timestamp.IsZero() {
		embed.SetTimestamp(b.timestamp)
	}

	return discord.NewMessageCreateBuilder().
		AddEmbeds(embed.Build()).
		Build()
}

func (b *MessageBuilder[C]) Serialize() SerializedMessageFormat {
	fields := make([]SerializedMessageFormatField, len(b.fields))
	for i, field := range b.fields {
		fields[i] = SerializedMessageFormatField{
			Name:   field.Name,
			Value:  field.Value,
			Inline: *field.Inline,
		}
	}
	return SerializedMessageFormat{
		Title:       b.title,
		Description: b.description,
		URL:         b.url,
		Icon:        b.icon,
		Fields:      fields,
		Timestamp:   b.timestamp,
		Thumbnail:   b.thumbnail,
	}
}

type SerializedMessageFormat struct {
	Title       string                         `json:"title,omitzero"`
	Description string                         `json:"description,omitzero"`
	URL         string                         `json:"url,omitzero"`
	Icon        MessageIcon                    `json:"icon,omitzero"`
	Fields      []SerializedMessageFormatField `json:"fields,omitzero"`
	Timestamp   time.Time                      `json:"timestamp,omitzero"`
	Thumbnail   string                         `json:"thumbnail,omitzero"`
}

func (f SerializedMessageFormat) MustMarshalString() string {
	json, err := json.Marshal(f)
	if err != nil {
		log.Fatal("Failed to marshal JSON", "err", err)
	}
	return string(json)
}

type SerializedMessageFormatField struct {
	Name   string `json:"name,omitzero"`
	Value  string `json:"value,omitzero"`
	Inline bool   `json:"inline,omitzero"`
}
