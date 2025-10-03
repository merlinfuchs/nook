package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type SubscriptionStore interface {
	UpsertSubscription(ctx context.Context, subscription model.Subscription) (*model.Subscription, error)
	SubscriptionByID(ctx context.Context, id common.ID) (*model.Subscription, error)
	SubscriptionsByUserID(ctx context.Context, userID common.ID) ([]*model.Subscription, error)
	SubscriptionsByGuildID(ctx context.Context, guildID common.ID) ([]*model.Subscription, error)
}
