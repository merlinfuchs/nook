package manager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/store"
)

func moduleSettingsCacheKey(guildID common.ID, moduleID string) string {
	return fmt.Sprintf("%s:%s", guildID, moduleID)
}

type ModuleUpdateListener func(guildID common.ID, module ModuleWithSettings)

type ModuleManager struct {
	modules        []module.GenericModule
	eventListeners []ModuleUpdateListener

	moduleSettingsStore store.ModuleSettingsStore

	cache *common.QueuedCache[string, model.ModuleSettings]
}

func NewModuleManager(inner store.ModuleSettingsStore) *ModuleManager {
	return &ModuleManager{
		moduleSettingsStore: inner,
		cache:               common.NewQueuedCache[string, model.ModuleSettings](time.Minute * 1),
	}
}

// AddModule adds a module to the manager
// this shouldn't be called after the client has started receiving events
func (s *ModuleManager) AddModule(module module.GenericModule) {
	s.modules = append(s.modules, module)
}

func (s *ModuleManager) Modules() []module.GenericModule {
	return s.modules
}

func (s *ModuleManager) Module(moduleID string) module.GenericModule {
	for _, module := range s.modules {
		if module.ModuleID() == moduleID {
			return module
		}
	}
	return nil
}

func (s *ModuleManager) EnabledModuleIDs(ctx context.Context, guildID common.ID) ([]string, error) {
	// TODO: Should we have a separate cache for this? It's called on every event
	modules, err := s.GuildModules(ctx, guildID)
	if err != nil {
		return nil, err
	}

	res := make([]string, len(modules))
	for i, mod := range modules {
		if mod.Settings.Enabled {
			res[i] = mod.Module.ModuleID()
		}
	}

	return res, nil
}

func (s *ModuleManager) IsModuleEnabled(ctx context.Context, guildID common.ID, moduleID string) (bool, error) {
	settings, err := s.ModuleSettings(ctx, guildID, moduleID)
	if err != nil {
		return false, err
	}
	return settings.Enabled, nil
}

func (s *ModuleManager) GuildModules(ctx context.Context, guildID common.ID) ([]ModuleWithSettings, error) {
	modules := s.Modules()

	res := make([]ModuleWithSettings, 0, len(modules))
	for _, module := range modules {
		settings, err := s.ModuleSettings(ctx, guildID, module.ModuleID())
		if err != nil {
			return nil, fmt.Errorf("error on getting module settings: %w", err)
		}

		res = append(res, ModuleWithSettings{
			Module:   module,
			Settings: settings,
		})
	}

	return res, nil
}

func (s *ModuleManager) GuildModule(ctx context.Context, guildID common.ID, moduleID string) (ModuleWithSettings, error) {
	module := s.Module(moduleID)
	if module == nil {
		return ModuleWithSettings{}, store.ErrNotFound
	}

	settings, err := s.ModuleSettings(ctx, guildID, moduleID)
	if err != nil {
		return ModuleWithSettings{}, err
	}

	return ModuleWithSettings{
		Module:   module,
		Settings: settings,
	}, nil
}

func (s *ModuleManager) AddModuleUpdateListener(listener ModuleUpdateListener) {
	s.eventListeners = append(s.eventListeners, listener)
}

func (s *ModuleManager) ModuleSettings(ctx context.Context, guildID common.ID, moduleID string) (model.ModuleSettings, error) {
	cacheKey := moduleSettingsCacheKey(guildID, moduleID)
	settings, err := s.cache.GetOrSet(cacheKey, func() (model.ModuleSettings, error) {
		module := s.Module(moduleID)
		if module == nil {
			return model.ModuleSettings{}, store.ErrNotFound
		}

		settings, err := s.moduleSettingsStore.ModuleSettings(ctx, guildID, moduleID)
		if errors.Is(err, store.ErrNotFound) {
			defaultConfig, err := json.Marshal(module.Metadata().DefaultConfig)
			if err != nil {
				return model.ModuleSettings{}, fmt.Errorf("failed to marshal default config: %w", err)
			}

			settings = model.ModuleSettings{
				GuildID:   guildID,
				ModuleID:  moduleID,
				Enabled:   module.Metadata().DefaultEnabled,
				UpdatedAt: time.Now().UTC(),
				Config:    defaultConfig,
			}
		} else if err != nil {
			return model.ModuleSettings{}, err
		}

		return settings, nil
	})
	if err != nil {
		return model.ModuleSettings{}, err
	}

	return settings, nil
}

func (s *ModuleManager) UpdateModuleSettings(ctx context.Context, settings model.ModuleSettings) error {
	mod := s.Module(settings.ModuleID)
	if mod == nil {
		return store.ErrNotFound
	}

	err := s.moduleSettingsStore.UpdateModuleSettings(ctx, settings)
	if err != nil {
		return err
	}

	s.cache.Set(moduleSettingsCacheKey(settings.GuildID, settings.ModuleID), settings)
	for _, listener := range s.eventListeners {
		go listener(settings.GuildID, ModuleWithSettings{
			Module:   mod,
			Settings: settings,
		})
	}
	return nil
}

type ModuleWithSettings struct {
	Module   module.GenericModule
	Settings model.ModuleSettings
}
