package manager

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
)

type PlanManager struct {
	entitlementStore  store.EntitlementStore
	subscriptionStore store.SubscriptionStore
	userStore         store.UserStore
	plans             []model.Plan
}

func NewPlanManager(
	entitlementStore store.EntitlementStore,
	subscriptionStore store.SubscriptionStore,
	userStore store.UserStore,
	plans []model.Plan,
) *PlanManager {
	return &PlanManager{
		entitlementStore:  entitlementStore,
		subscriptionStore: subscriptionStore,
		userStore:         userStore,
		plans:             plans,
	}
}

func (m *PlanManager) Run(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := m.cleanupFeaturesWithoutAccess(ctx)
			if err != nil {
				slog.Error("Failed to cleanup features without access", slog.String("error", err.Error()))
			}
		}
	}
}

func (m *PlanManager) Plans() []model.Plan {
	return m.plans
}

func (m *PlanManager) PlanByPaddlePriceID(priceID string) *model.Plan {
	for _, plan := range m.plans {
		if plan.PaddleMonthlyPriceID == priceID || plan.PaddleYearlyPriceID == priceID {
			return &plan
		}
	}
	return nil
}

func (m *PlanManager) GuildFeatures(ctx context.Context, guildID common.ID) (model.Features, string) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	entitlements, err := m.entitlementStore.ActiveEntitlementsByGuildID(ctx, guildID, time.Now().UTC())
	if err != nil {
		slog.Error(
			"Failed to get active entitlements",
			slog.Int64("guild_id", int64(guildID)),
			slog.String("error", err.Error()),
		)
	}

	var features model.Features
	var planID string
	for _, plan := range m.plans {
		if plan.Default {
			features = features.Merge(plan.Features)
			planID = plan.ID
			continue
		}

		for _, entitlement := range entitlements {
			if slices.Contains(entitlement.PlanIDs, plan.ID) {
				features = features.Merge(plan.Features)
				planID = plan.ID
				break
			}
		}
	}

	return features, planID
}

func (m *PlanManager) cleanupFeaturesWithoutAccess(ctx context.Context) error {
	return nil
}
