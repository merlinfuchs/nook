package manager

import (
	"context"
	"errors"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

type ModuleValueManager struct {
	store store.ModuleValueStore
	// TODO?: Caching
}

func NewModuleValueManager(
	store store.ModuleValueStore,
) *ModuleValueManager {
	return &ModuleValueManager{
		store: store,
	}
}

func (m *ModuleValueManager) SetModuleValue(
	ctx context.Context,
	guildID common.ID,
	moduleID string,
	key string,
	value thing.Thing,
) error {
	return m.store.SetModuleValue(ctx, model.ModuleValue{
		GuildID:   guildID,
		ModuleID:  moduleID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
}

func (m *ModuleValueManager) UpdateModuleValue(
	ctx context.Context,
	op thing.Operation,
	guildID common.ID,
	moduleID string,
	key string,
	value thing.Thing,
) (thing.Thing, error) {
	newValue, err := m.store.UpdateModuleValue(ctx, op, model.ModuleValue{
		GuildID:   guildID,
		ModuleID:  moduleID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return thing.Null, err
	}
	return newValue.Value, nil
}

func (m *ModuleValueManager) ModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) (thing.Thing, error) {
	value, err := m.store.ModuleValue(ctx, guildID, moduleID, key)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return thing.Null, err
		}
		return thing.Null, err
	}
	return value.Value, nil
}

func (m *ModuleValueManager) DeleteModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) error {
	return m.store.DeleteModuleValue(ctx, guildID, moduleID, key)
}

func (m *ModuleValueManager) Scoped(moduleID string) *ScopedModuleValueManager {
	return &ScopedModuleValueManager{
		inner:    m,
		moduleID: moduleID,
	}
}

type ScopedModuleValueManager struct {
	inner    *ModuleValueManager
	moduleID string
}

func (s *ScopedModuleValueManager) ModuleValue(ctx context.Context, guildID common.ID, key string) (thing.Thing, error) {
	return s.inner.ModuleValue(ctx, guildID, s.moduleID, key)
}

func (s *ScopedModuleValueManager) UpdateModuleValue(
	ctx context.Context,
	op thing.Operation,
	guildID common.ID,
	key string,
	value thing.Thing,
) (thing.Thing, error) {
	return s.inner.UpdateModuleValue(ctx, op, guildID, s.moduleID, key, value)
}

func (s *ScopedModuleValueManager) DeleteModuleValue(ctx context.Context, guildID common.ID, key string) error {
	return s.inner.DeleteModuleValue(ctx, guildID, s.moduleID, key)
}

func (s *ScopedModuleValueManager) SetModuleValue(ctx context.Context, guildID common.ID, key string, value thing.Thing) error {
	return s.inner.SetModuleValue(ctx, guildID, s.moduleID, key, value)
}
