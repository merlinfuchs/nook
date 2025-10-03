package store

import (
	"context"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type EntitlementStore interface {
	UpsertSubscriptionEntitlement(ctx context.Context, entitlement model.Entitlement) (*model.Entitlement, error)
	EntitlementsByGuildID(ctx context.Context, guildID common.ID) ([]*model.Entitlement, error)
	ActiveEntitlementsByGuildID(ctx context.Context, guildID common.ID, now time.Time) ([]*model.Entitlement, error)
}
