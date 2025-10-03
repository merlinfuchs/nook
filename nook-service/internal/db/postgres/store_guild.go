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

var _ store.GuildStore = &Client{}

func (c *Client) UpsertGuild(ctx context.Context, guild model.Guild) (*model.Guild, error) {
	row, err := c.Q.UpsertGuild(ctx, pgmodel.UpsertGuildParams{
		ID:          int64(guild.ID),
		Name:        guild.Name,
		Description: pgtype.Text{String: guild.Description.String, Valid: guild.Description.Valid},
		Icon:        pgtype.Text{String: guild.Icon.String, Valid: guild.Icon.Valid},
		Unavailable: guild.Unavailable,
		OwnerUserID: int64(guild.OwnerUserID),
		CreatedAt:   pgtype.Timestamp{Time: guild.CreatedAt, Valid: true},
		UpdatedAt:   pgtype.Timestamp{Time: guild.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return rowToGuild(row), nil
}

func (c *Client) UpdateGuildDeleted(ctx context.Context, guildID common.ID, deleted bool) error {
	return c.Q.UpdateGuildDeleted(ctx, pgmodel.UpdateGuildDeletedParams{
		ID:      int64(guildID),
		Deleted: deleted,
	})
}

func (c *Client) UpdateGuildUnavailable(ctx context.Context, guildID common.ID, unavailable bool) error {
	return c.Q.UpdateGuildUnavailable(ctx, pgmodel.UpdateGuildUnavailableParams{
		ID:          int64(guildID),
		Unavailable: unavailable,
	})
}

func (c *Client) Guild(ctx context.Context, guildID common.ID) (*model.Guild, error) {
	row, err := c.Q.GetGuild(ctx, int64(guildID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToGuild(row), nil
}

func (c *Client) GuildsByOwnerUserID(ctx context.Context, ownerUserID common.ID) ([]model.Guild, error) {
	rows, err := c.Q.GetGuildsByOwnerUserID(ctx, int64(ownerUserID))
	if err != nil {
		return nil, err
	}

	res := make([]model.Guild, 0, len(rows))
	for _, row := range rows {
		res = append(res, *rowToGuild(row))
	}

	return res, nil
}
func rowToGuild(row pgmodel.Guild) *model.Guild {
	return &model.Guild{
		ID:          common.ID(row.ID),
		Name:        row.Name,
		Description: null.NewString(row.Description.String, row.Description.Valid),
		Icon:        null.NewString(row.Icon.String, row.Icon.Valid),
		Unavailable: row.Unavailable,
		Deleted:     row.Deleted,
		OwnerUserID: common.ID(row.OwnerUserID),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
}
