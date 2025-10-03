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
	"github.com/merlinfuchs/nook/nook-service/thing"
)

var _ store.ModuleValueStore = &Client{}

func (c *Client) ModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) (*model.ModuleValue, error) {
	row, err := c.Q.GetModuleValue(ctx, pgmodel.GetModuleValueParams{
		GuildID:  int64(guildID),
		ModuleID: moduleID,
		Key:      key,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToModuleValue(row)
}

func (c *Client) DeleteModuleValue(ctx context.Context, guildID common.ID, moduleID, key string) error {
	err := c.Q.DeleteModuleValue(ctx, pgmodel.DeleteModuleValueParams{
		GuildID:  int64(guildID),
		ModuleID: moduleID,
		Key:      key,
	})

	return err
}

func (c *Client) SetModuleValue(ctx context.Context, value model.ModuleValue) error {
	_, err := c.setModuleValueWithTx(ctx, nil, value)
	return err
}

func (c *Client) UpdateModuleValue(ctx context.Context, operation thing.Operation, value model.ModuleValue) (*model.ModuleValue, error) {
	if operation == thing.OperationOverwrite {
		return c.setModuleValueWithTx(ctx, nil, value)
	}

	tx, err := c.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	currentValue, err := c.moduleValueWithTx(ctx, tx, value.GuildID, value.ModuleID, value.Key)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Current trasaction is rolled back, we set the value outside of the transaction
			return c.setModuleValueWithTx(ctx, nil, value)
		}
		return nil, fmt.Errorf("failed to get current module value: %w", err)
	}

	value.Value = currentValue.Value.Perform(value.Value, operation)

	newValue, err := c.setModuleValueWithTx(ctx, tx, value)
	if err != nil {
		return nil, fmt.Errorf("failed to set module value: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return newValue, nil
}

func (c *Client) moduleValueWithTx(ctx context.Context, tx pgx.Tx, guildID common.ID, moduleID string, key string) (*model.ModuleValue, error) {
	q := c.Q
	if tx != nil {
		q = c.Q.WithTx(tx)
	}

	row, err := q.GetModuleValueForUpdate(ctx, pgmodel.GetModuleValueForUpdateParams{
		GuildID:  int64(guildID),
		ModuleID: moduleID,
		Key:      key,
	})

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rowToModuleValue(row)
}

func (c *Client) setModuleValueWithTx(ctx context.Context, tx pgx.Tx, value model.ModuleValue) (*model.ModuleValue, error) {
	q := c.Q
	if tx != nil {
		q = c.Q.WithTx(tx)
	}

	data, err := json.Marshal(value.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal module value: %w", err)
	}

	row, err := q.SetModuleValue(ctx, pgmodel.SetModuleValueParams{
		GuildID:   int64(value.GuildID),
		ModuleID:  value.ModuleID,
		Key:       value.Key,
		Value:     data,
		CreatedAt: pgtype.Timestamp{Time: value.CreatedAt, Valid: true},
		UpdatedAt: pgtype.Timestamp{Time: value.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return rowToModuleValue(row)
}

func rowToModuleValue(row pgmodel.ModuleValue) (*model.ModuleValue, error) {
	var data thing.Thing
	err := json.Unmarshal(row.Value, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal module value: %w", err)
	}

	return &model.ModuleValue{
		ID:        uint64(row.ID),
		GuildID:   common.ID(row.GuildID),
		ModuleID:  row.ModuleID,
		Key:       row.Key,
		Value:     data,
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}, nil
}
