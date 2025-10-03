package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres/pgmodel"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
	"gopkg.in/guregu/null.v4"
)

var _ store.GuildSettingsStore = &Client{}

func (c *Client) UpdateGuildSettings(ctx context.Context, settings model.GuildSettings) error {
	return c.Q.UpsertGuildSettings(ctx, pgmodel.UpsertGuildSettingsParams{
		GuildID:       int64(settings.GuildID),
		CommandPrefix: pgtype.Text{String: settings.CommandPrefix.String, Valid: settings.CommandPrefix.Valid},
		ColorScheme:   pgtype.Text{String: settings.ColorScheme.String, Valid: settings.ColorScheme.Valid},
		UpdatedAt:     pgtype.Timestamp{Time: settings.UpdatedAt, Valid: true},
	})
}

func (c *Client) GuildSettings(ctx context.Context, guildID common.ID) (model.GuildSettings, error) {
	row, err := c.Q.GetGuildSettings(ctx, int64(guildID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.GuildSettings{}, store.ErrNotFound
		}
		return model.GuildSettings{}, err
	}

	return model.GuildSettings{
		GuildID:       common.ID(row.GuildID),
		CommandPrefix: null.NewString(row.CommandPrefix.String, row.CommandPrefix.Valid),
		UpdatedAt:     row.UpdatedAt.Time,
	}, nil
}

func (c *Client) DeleteGuildSettings(ctx context.Context, guildID common.ID) error {
	return c.Q.DeleteGuildSettings(ctx, int64(guildID))
}

func (c *Client) AllGuildSettings(ctx context.Context) ([]model.GuildSettings, error) {
	rows, err := c.Q.GetAllGuildSettings(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]model.GuildSettings, 0, len(rows))
	for _, row := range rows {
		res = append(res, model.GuildSettings{
			GuildID:       common.ID(row.GuildID),
			CommandPrefix: null.NewString(row.CommandPrefix.String, row.CommandPrefix.Valid),
			UpdatedAt:     row.UpdatedAt.Time,
		})
	}

	return res, nil
}
