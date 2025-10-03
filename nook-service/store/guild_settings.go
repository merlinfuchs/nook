package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type GuildSettingsStore interface {
	UpdateGuildSettings(ctx context.Context, settings model.GuildSettings) error
	GuildSettings(ctx context.Context, guildID common.ID) (model.GuildSettings, error)
	DeleteGuildSettings(ctx context.Context, guildID common.ID) error
	AllGuildSettings(ctx context.Context) ([]model.GuildSettings, error)
}
