package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

type ModuleValueStore interface {
	SetModuleValue(ctx context.Context, value model.ModuleValue) error
	UpdateModuleValue(ctx context.Context, operation thing.Operation, value model.ModuleValue) (*model.ModuleValue, error)
	ModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) (*model.ModuleValue, error)
	DeleteModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) error
}
