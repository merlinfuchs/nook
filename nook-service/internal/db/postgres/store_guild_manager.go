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
)

var _ store.GuildManagerStore = &Client{}

func (c *Client) UpsertGuildManager(ctx context.Context, manager model.GuildManager) (*model.GuildManager, error) {
	row, err := c.Q.UpsertGuildManager(ctx, pgmodel.UpsertGuildManagerParams{
		GuildID:   int64(manager.GuildID),
		UserID:    int64(manager.UserID),
		Role:      string(manager.Role),
		CreatedAt: pgtype.Timestamp{Time: manager.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: manager.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return rowToGuildManager(row), nil
}

func (c *Client) GetGuildManager(ctx context.Context, guildID common.ID, userID common.ID) (*model.GuildManager, error) {
	row, err := c.Q.GetGuildManager(ctx, pgmodel.GetGuildManagerParams{
		GuildID: int64(guildID),
		UserID:  int64(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToGuildManager(row), nil
}

func (c *Client) DeleteGuildManager(ctx context.Context, guildID common.ID, userID common.ID) error {
	return c.Q.DeleteGuildManager(ctx, pgmodel.DeleteGuildManagerParams{
		GuildID: int64(guildID),
		UserID:  int64(userID),
	})
}

func (c *Client) GetGuildManagers(ctx context.Context, guildID common.ID) ([]model.GuildManager, error) {
	rows, err := c.Q.GetGuildManagers(ctx, int64(guildID))
	if err != nil {
		return nil, err
	}

	res := make([]model.GuildManager, 0, len(rows))
	for _, row := range rows {
		res = append(res, *rowToGuildManager(row))
	}

	return res, nil
}

func (c *Client) GetGuildManagersWithUsers(ctx context.Context, guildID common.ID) ([]model.GuildManagerWithUser, error) {
	rows, err := c.Q.GetGuildManagersWithUsers(ctx, int64(guildID))
	if err != nil {
		return nil, err
	}

	res := make([]model.GuildManagerWithUser, 0, len(rows))
	for _, row := range rows {
		res = append(res, model.GuildManagerWithUser{
			GuildManager: *rowToGuildManager(row.GuildManager),
			User:         *rowToUser(row.User),
		})
	}

	return res, nil
}

func rowToGuildManager(row pgmodel.GuildManager) *model.GuildManager {
	return &model.GuildManager{
		GuildID:   common.ID(row.GuildID),
		UserID:    common.ID(row.UserID),
		Role:      model.GuildManagerRole(row.Role),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}
}
