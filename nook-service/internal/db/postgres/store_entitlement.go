package postgres

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres/pgmodel"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
	"gopkg.in/guregu/null.v4"
)

var _ store.EntitlementStore = &Client{}

func (c *Client) UpsertSubscriptionEntitlement(ctx context.Context, entitlement model.Entitlement) (*model.Entitlement, error) {
	ent, err := c.Q.UpsertSubscriptionEntitlement(ctx, pgmodel.UpsertSubscriptionEntitlementParams{
		ID:             int64(entitlement.ID),
		GuildID:        int64(entitlement.GuildID),
		SubscriptionID: pgtype.Int8{Int64: int64(entitlement.SubscriptionID.ID), Valid: entitlement.SubscriptionID.Valid},
		Type:           entitlement.Type,
		PlanIds:        entitlement.PlanIDs,
		StartsAt:       pgtype.Timestamp{Time: entitlement.StartsAt.Time, Valid: entitlement.StartsAt.Valid},
		EndsAt:         pgtype.Timestamp{Time: entitlement.EndsAt.Time, Valid: entitlement.EndsAt.Valid},
		CreatedAt:      pgtype.Timestamp{Time: entitlement.CreatedAt, Valid: true},
		UpdatedAt:      pgtype.Timestamp{Time: entitlement.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return rowToEntitlement(&ent), nil
}

func (c *Client) EntitlementsByGuildID(ctx context.Context, guildID common.ID) ([]*model.Entitlement, error) {
	rows, err := c.Q.GetEntitlementsByGuildID(ctx, int64(guildID))
	if err != nil {
		return nil, err
	}

	ents := make([]*model.Entitlement, len(rows))
	for i, row := range rows {
		ents[i] = rowToEntitlement(&row)
	}
	return ents, nil
}

func (c *Client) ActiveEntitlementsByGuildID(ctx context.Context, guildID common.ID, now time.Time) ([]*model.Entitlement, error) {
	rows, err := c.Q.GetActiveEntitlementsByGuildID(ctx, pgmodel.GetActiveEntitlementsByGuildIDParams{
		GuildID:  int64(guildID),
		StartsAt: pgtype.Timestamp{Time: now, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	ents := make([]*model.Entitlement, len(rows))
	for i, row := range rows {
		ents[i] = rowToEntitlement(&row)
	}
	return ents, nil
}

func rowToEntitlement(row *pgmodel.Entitlement) *model.Entitlement {
	return &model.Entitlement{
		ID:             common.ID(row.ID),
		SubscriptionID: common.NullID{ID: common.ID(row.SubscriptionID.Int64), Valid: row.SubscriptionID.Valid},
		Type:           row.Type,
		PlanIDs:        row.PlanIds,
		StartsAt:       null.NewTime(row.StartsAt.Time, row.StartsAt.Valid),
		EndsAt:         null.NewTime(row.EndsAt.Time, row.EndsAt.Valid),
		CreatedAt:      row.CreatedAt.Time,
		UpdatedAt:      row.UpdatedAt.Time,
	}
}
