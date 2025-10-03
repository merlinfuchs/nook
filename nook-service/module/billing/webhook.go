package billing

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/PaddleHQ/paddle-go-sdk"
	"github.com/PaddleHQ/paddle-go-sdk/pkg/paddlenotification"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
	"gopkg.in/guregu/null.v4"
)

func (m *BillingModule) handleBillingWebhook(c *api.Context) error {
	ok, err := m.webhookVerifier.Verify(c.Request())
	if err != nil {
		return fmt.Errorf("failed to verify webhook: %w", err)
	}

	if !ok {
		return api.ErrUnauthorized("invalid_signature", "Invalid webhook signature")
	}

	var event paddleWebhookEvent
	if err := c.ParseBody(&event); err != nil {
		return fmt.Errorf("failed to parse webhook event: %w", err)
	}

	rawUserID, _ := event.Data.CustomData["user_id"].(string)
	if rawUserID == "" {
		return fmt.Errorf("user id is required")
	}

	userID, err := common.ParseID(rawUserID)
	if err != nil {
		return fmt.Errorf("failed to parse user id: %w", err)
	}

	_, err = m.userStore.User(c.Context(), userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			slog.Info(
				"User for subscription event not found, skipping",
				slog.String("user_id", userID.String()),
			)
			return nil
		}
		return fmt.Errorf("failed to get user: %w", err)
	}

	productIDs := make([]string, len(event.Data.Items))
	priceIDs := make([]string, len(event.Data.Items))
	planIDs := make([]string, 0, len(event.Data.Items))
	for i, item := range event.Data.Items {
		productIDs[i] = item.Product.ID
		priceIDs[i] = item.Price.ID

		plan := m.planManager.PlanByPaddlePriceID(item.Price.ID)
		if plan != nil {
			planIDs = append(planIDs, plan.ID)
		}
	}

	var currentPeriodEndsAt null.Time
	var currentPeriodStartsAt null.Time
	if event.Data.CurrentBillingPeriod != nil {
		currentPeriodEndsAt = parseNullISOTimestamp(&event.Data.CurrentBillingPeriod.EndsAt)
		currentPeriodStartsAt = parseNullISOTimestamp(&event.Data.CurrentBillingPeriod.StartsAt)
	}

	sub, err := m.subscriptionStore.UpsertSubscription(c.Context(), model.Subscription{
		ID:                    common.UniqueID(),
		UserID:                userID,
		Status:                string(event.Data.Status),
		PaddleSubscriptionID:  event.Data.ID,
		PaddleCustomerID:      event.Data.CustomerID,
		PaddleProductIds:      productIDs,
		PaddlePriceIds:        priceIDs,
		CreatedAt:             parseISOTimestamp(event.Data.CreatedAt),
		StartedAt:             parseNullISOTimestamp(event.Data.StartedAt),
		PausedAt:              parseNullISOTimestamp(event.Data.PausedAt),
		CanceledAt:            parseNullISOTimestamp(event.Data.CanceledAt),
		CurrentPeriodEndsAt:   currentPeriodEndsAt,
		CurrentPeriodStartsAt: currentPeriodStartsAt,
		UpdatedAt:             parseISOTimestamp(event.Data.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert subscription: %w", err)
	}

	rawGuildID, _ := event.Data.CustomData["guild_id"].(string)
	if rawGuildID == "" {
		return fmt.Errorf("guild id is required")
	}

	guildID, err := common.ParseID(rawGuildID)
	if err != nil {
		return fmt.Errorf("failed to parse guild id: %w", err)
	}

	endsAt := event.Data.CanceledAt
	if event.Data.CurrentBillingPeriod != nil {
		endsAt = &event.Data.CurrentBillingPeriod.EndsAt
	}

	// What about subscriptions with multiple items?
	// Should we have an entitlement for each item or should an entitlement have a list of plan ids?
	_, err = m.entitlementStore.UpsertSubscriptionEntitlement(c.Context(), model.Entitlement{
		ID:             common.UniqueID(),
		SubscriptionID: common.NullID{ID: sub.ID, Valid: true},
		GuildID:        guildID,
		Type:           "subscription",
		PlanIDs:        planIDs,
		StartsAt:       parseNullISOTimestamp(event.Data.StartedAt),
		EndsAt:         parseNullISOTimestamp(endsAt),
		CreatedAt:      parseISOTimestamp(event.Data.CreatedAt),
		UpdatedAt:      parseISOTimestamp(event.Data.UpdatedAt),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert entitlement: %w", err)
	}

	return c.Send(200, []byte("ok"))
}

type paddleWebhookEvent struct {
	EventID   string                                      `json:"event_id"`
	EventType paddle.EventTypeName                        `json:"event_type"`
	Data      paddlenotification.SubscriptionNotification `json:"data"`
}

func parseISOTimestamp(timestamp string) time.Time {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}
	}
	return t
}

func parseNullISOTimestamp(timestamp *string) null.Time {
	if timestamp == nil {
		return null.Time{}
	}
	t := parseISOTimestamp(*timestamp)
	if t.IsZero() {
		return null.Time{}
	}
	return null.NewTime(t, true)
}
