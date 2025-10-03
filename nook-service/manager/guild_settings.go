package manager

import (
	"context"
	"errors"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
)

type GuildSettingsManagerConfig struct {
	DefaultPrefix      string
	DefaultColorScheme string
}

type GuildSettingsManager struct {
	guildSettingsStore store.GuildSettingsStore

	config GuildSettingsManagerConfig
	cache  *common.QueuedCache[common.ID, model.GuildSettings]
}

func NewGuildSettingsManager(inner store.GuildSettingsStore, config GuildSettingsManagerConfig) *GuildSettingsManager {
	return &GuildSettingsManager{
		guildSettingsStore: inner,
		config:             config,
		cache:              common.NewQueuedCache[common.ID, model.GuildSettings](time.Minute * 1),
	}
}

func (s *GuildSettingsManager) InvalidateCache(guildID common.ID) {
	s.cache.Delete(guildID)
}

func (s *GuildSettingsManager) GuildSettings(ctx context.Context, guildID common.ID) (model.GuildSettings, error) {
	settings, err := s.cache.GetOrSet(guildID, func() (model.GuildSettings, error) {
		settings, err := s.guildSettingsStore.GuildSettings(ctx, guildID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				return model.GuildSettings{}, nil
			}
			return model.GuildSettings{}, err
		}
		return settings, nil
	})

	return settings, err
}

func (s *GuildSettingsManager) ResolvedGuildSettings(ctx context.Context, guildID common.ID) (model.ResolvedGuildSettings, error) {
	settings, err := s.GuildSettings(ctx, guildID)
	if err != nil {
		return model.ResolvedGuildSettings{}, err
	}
	return s.resolveGuildSettings(settings), nil
}

func (s *GuildSettingsManager) UpdateGuildSettings(ctx context.Context, guildID common.ID, settings model.GuildSettings) error {
	err := s.guildSettingsStore.UpdateGuildSettings(ctx, settings)
	if err != nil {
		return err
	}
	s.cache.Set(guildID, settings)
	return nil
}

func (s *GuildSettingsManager) DefaultSettings() model.ResolvedGuildSettings {
	return s.resolveGuildSettings(model.GuildSettings{})
}

func (s *GuildSettingsManager) resolveGuildSettings(settings model.GuildSettings) model.ResolvedGuildSettings {
	commandPrefix := settings.CommandPrefix.String
	if commandPrefix == "" {
		commandPrefix = s.config.DefaultPrefix
	}

	colorScheme := settings.ColorScheme.String
	if colorScheme == "" {
		colorScheme = s.config.DefaultColorScheme
	}

	return model.ResolvedGuildSettings{
		CommandPrefix: commandPrefix,
		ColorScheme:   colorScheme,
	}
}
