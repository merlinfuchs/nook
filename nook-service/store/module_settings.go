package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type ModuleSettingsStore interface {
	UpdateModuleSettings(ctx context.Context, settings model.ModuleSettings) error
	ModuleSettings(ctx context.Context, guildID common.ID, moduleID string) (model.ModuleSettings, error)
	GuildModuleSettings(ctx context.Context, guildID common.ID) ([]model.ModuleSettings, error)
	DeleteModuleSettings(ctx context.Context, guildID common.ID, moduleID string) error
}
