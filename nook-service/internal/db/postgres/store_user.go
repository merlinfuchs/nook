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

var _ store.UserStore = &Client{}

func (c *Client) User(ctx context.Context, id common.ID) (*model.User, error) {
	row, err := c.Q.GetUser(ctx, int64(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToUser(row), nil
}

func (c *Client) UpsertUser(ctx context.Context, user *model.User) (*model.User, error) {
	row, err := c.Q.UpsertUser(ctx, pgmodel.UpsertUserParams{
		ID:            int64(user.ID),
		Username:      user.Username,
		Discriminator: pgtype.Text{String: user.Discriminator, Valid: true},
		DisplayName:   user.DisplayName,
		Avatar:        pgtype.Text{String: user.Avatar.String, Valid: user.Avatar.Valid},
		CreatedAt: pgtype.Timestamp{
			Time:  user.CreatedAt.UTC(),
			Valid: true,
		},
		UpdatedAt: pgtype.Timestamp{
			Time:  user.UpdatedAt.UTC(),
			Valid: true,
		},
	})
	if err != nil {
		return nil, err
	}

	return rowToUser(row), nil
}

func rowToUser(row pgmodel.User) *model.User {
	return &model.User{
		ID:            common.ID(row.ID),
		Username:      row.Username,
		Discriminator: row.Discriminator.String,
		DisplayName:   row.DisplayName,
		Avatar:        null.NewString(row.Avatar.String, row.Avatar.Valid),
		CreatedAt:     row.CreatedAt.Time,
		UpdatedAt:     row.UpdatedAt.Time,
	}
}
