package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres/pgmodel"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
)

var _ store.ModuleSettingsStore = &Client{}

func (c *Client) UpdateModuleSettings(ctx context.Context, settings model.ModuleSettings) error {
	commandOverwrites, err := json.Marshal(settings.CommandOverwrites)
	if err != nil {
		return err
	}

	return c.Q.UpsertModuleSettings(ctx, pgmodel.UpsertModuleSettingsParams{
		GuildID:           int64(settings.GuildID),
		ModuleID:          settings.ModuleID,
		Enabled:           settings.Enabled,
		CommandOverwrites: commandOverwrites,
		Config:            settings.Config,
		UpdatedAt:         pgtype.Timestamp{Time: settings.UpdatedAt, Valid: true},
	})
}

func (c *Client) ModuleSettings(ctx context.Context, guildID common.ID, moduleID string) (model.ModuleSettings, error) {
	row, err := c.Q.GetModuleSettings(ctx, pgmodel.GetModuleSettingsParams{
		GuildID:  int64(guildID),
		ModuleID: moduleID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ModuleSettings{}, store.ErrNotFound
		}
		return model.ModuleSettings{}, err
	}

	var commandOverwrites map[string]model.ModuleCommandOverwrite
	if err := json.Unmarshal(row.CommandOverwrites, &commandOverwrites); err != nil {
		return model.ModuleSettings{}, err
	}

	return model.ModuleSettings{
		GuildID:           common.ID(row.GuildID),
		ModuleID:          row.ModuleID,
		Enabled:           row.Enabled,
		CommandOverwrites: commandOverwrites,
		Config:            row.Config,
		UpdatedAt:         row.UpdatedAt.Time,
	}, nil
}

func (c *Client) DeleteModuleSettings(ctx context.Context, guildID common.ID, moduleID string) error {
	return c.Q.DeleteModuleSettings(ctx, pgmodel.DeleteModuleSettingsParams{
		GuildID:  int64(guildID),
		ModuleID: moduleID,
	})
}

func (c *Client) GuildModuleSettings(ctx context.Context, guildID common.ID) ([]model.ModuleSettings, error) {
	rows, err := c.Q.GetModuleSettingsByGuildID(ctx, int64(guildID))
	if err != nil {
		return nil, err
	}

	res := make([]model.ModuleSettings, 0, len(rows))
	for _, row := range rows {
		var commandOverwrites map[string]model.ModuleCommandOverwrite
		if err := json.Unmarshal(row.CommandOverwrites, &commandOverwrites); err != nil {
			return nil, err
		}

		res = append(res, model.ModuleSettings{
			GuildID:           common.ID(row.GuildID),
			ModuleID:          row.ModuleID,
			Enabled:           row.Enabled,
			CommandOverwrites: commandOverwrites,
			Config:            row.Config,
			UpdatedAt:         row.UpdatedAt.Time,
		})
	}

	return res, nil
}
