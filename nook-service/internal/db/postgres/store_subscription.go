package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres/pgmodel"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
	"gopkg.in/guregu/null.v4"
)

var _ store.SubscriptionStore = &Client{}

func (c *Client) UpsertSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error) {
	sub, err := c.Q.UpsertSubscription(ctx, pgmodel.UpsertSubscriptionParams{
		ID:                    int64(subscription.ID),
		UserID:                int64(subscription.UserID),
		Status:                subscription.Status,
		PaddleSubscriptionID:  subscription.PaddleSubscriptionID,
		PaddleCustomerID:      subscription.PaddleCustomerID,
		PaddleProductIds:      subscription.PaddleProductIds,
		PaddlePriceIds:        subscription.PaddlePriceIds,
		CreatedAt:             pgtype.Timestamp{Time: subscription.CreatedAt, Valid: true},
		StartedAt:             pgtype.Timestamp{Time: subscription.StartedAt.Time, Valid: subscription.StartedAt.Valid},
		PausedAt:              pgtype.Timestamp{Time: subscription.PausedAt.Time, Valid: subscription.PausedAt.Valid},
		CanceledAt:            pgtype.Timestamp{Time: subscription.CanceledAt.Time, Valid: subscription.CanceledAt.Valid},
		CurrentPeriodEndsAt:   pgtype.Timestamp{Time: subscription.CurrentPeriodEndsAt.Time, Valid: subscription.CurrentPeriodEndsAt.Valid},
		CurrentPeriodStartsAt: pgtype.Timestamp{Time: subscription.CurrentPeriodStartsAt.Time, Valid: subscription.CurrentPeriodStartsAt.Valid},
		UpdatedAt:             pgtype.Timestamp{Time: subscription.UpdatedAt, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	return rowToSubscription(&sub), nil
}

func (c *Client) SubscriptionByID(ctx context.Context, id common.ID) (*model.Subscription, error) {
	row, err := c.Q.GetSubscriptionByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return rowToSubscription(&row), nil
}

func (c *Client) SubscriptionsByUserID(ctx context.Context, userID common.ID) ([]*model.Subscription, error) {
	rows, err := c.Q.GetSubscriptionsByUserID(ctx, int64(userID))
	if err != nil {
		return nil, err
	}

	subs := make([]*model.Subscription, len(rows))
	for i, row := range rows {
		subs[i] = rowToSubscription(&row)
	}
	return subs, nil
}

func (c *Client) SubscriptionsByGuildID(ctx context.Context, guildID common.ID) ([]*model.Subscription, error) {
	rows, err := c.Q.GetSubscriptionsByGuildID(ctx, int64(guildID))
	if err != nil {
		return nil, err
	}

	subs := make([]*model.Subscription, len(rows))
	for i, row := range rows {
		subs[i] = rowToSubscription(&row)
	}
	return subs, nil
}

func rowToSubscription(row *pgmodel.Subscription) *model.Subscription {
	return &model.Subscription{
		ID:                    common.ID(row.ID),
		UserID:                common.ID(row.UserID),
		Status:                row.Status,
		PaddleSubscriptionID:  row.PaddleSubscriptionID,
		PaddleCustomerID:      row.PaddleCustomerID,
		PaddleProductIds:      row.PaddleProductIds,
		PaddlePriceIds:        row.PaddlePriceIds,
		CreatedAt:             row.CreatedAt.Time,
		StartedAt:             null.NewTime(row.StartedAt.Time, row.StartedAt.Valid),
		PausedAt:              null.NewTime(row.PausedAt.Time, row.PausedAt.Valid),
		CanceledAt:            null.NewTime(row.CanceledAt.Time, row.CanceledAt.Valid),
		CurrentPeriodEndsAt:   null.NewTime(row.CurrentPeriodEndsAt.Time, row.CurrentPeriodEndsAt.Valid),
		CurrentPeriodStartsAt: null.NewTime(row.CurrentPeriodStartsAt.Time, row.CurrentPeriodStartsAt.Valid),
		UpdatedAt:             row.UpdatedAt.Time,
	}
}
