package postgres

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres/pgmodel"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
)

var _ store.SessionStore = &Client{}

func (c *Client) CreateSession(ctx context.Context, session *model.Session) error {
	rawGuilds, err := json.Marshal(session.Guilds)
	if err != nil {
		return fmt.Errorf("failed to marshal guilds: %w", err)
	}

	err = c.Q.CreateSession(ctx, pgmodel.CreateSessionParams{
		KeyHash:        session.KeyHash,
		UserID:         int64(session.UserID),
		TokenAccess:    session.TokenAccess,
		TokenRefresh:   session.TokenRefresh,
		TokenScopes:    session.TokenScopes,
		Guilds:         rawGuilds,
		TokenExpiresAt: pgtype.Timestamp{Time: session.TokenExpiresAt.UTC(), Valid: true},
		CreatedAt:      pgtype.Timestamp{Time: session.CreatedAt.UTC(), Valid: true},
		ExpiresAt:      pgtype.Timestamp{Time: session.ExpiresAt.UTC(), Valid: true},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteSession(ctx context.Context, keyHash string) error {
	return c.Q.DeleteSession(ctx, keyHash)
}

func (c *Client) Session(ctx context.Context, keyHash string) (*model.Session, error) {
	row, err := c.Q.GetSession(ctx, keyHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToSession(row)
}

func rowToSession(row pgmodel.Session) (*model.Session, error) {
	var guilds []model.SessionGuild
	if err := json.Unmarshal(row.Guilds, &guilds); err != nil {
		return nil, fmt.Errorf("failed to unmarshal guilds: %w", err)
	}

	return &model.Session{
		KeyHash:        row.KeyHash,
		UserID:         common.ID(row.UserID),
		TokenAccess:    row.TokenAccess,
		TokenRefresh:   row.TokenRefresh,
		TokenScopes:    row.TokenScopes,
		TokenExpiresAt: row.TokenExpiresAt.Time,
		Guilds:         guilds,
		CreatedAt:      row.CreatedAt.Time,
		ExpiresAt:      row.ExpiresAt.Time,
	}, nil
}
