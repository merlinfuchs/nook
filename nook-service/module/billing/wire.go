package billing

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"gopkg.in/guregu/null.v4"
)

type BillingPlanWire struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Default     bool   `json:"default"`
	Popular     bool   `json:"popular"`
	Hidden      bool   `json:"hidden"`

	PaddleMonthlyPriceID string `json:"paddle_monthly_price_id,omitempty"`
	PaddleYearlyPriceID  string `json:"paddle_yearly_price_id,omitempty"`

	DiscordRoleID string `json:"discord_role_id"`

	Features BillingFeaturesWire `json:"features"`
}

type BillingPlanListResponseWire []BillingPlanWire

type BillingFeaturesWire struct {
	PlanID      string `json:"plan_id"`
	BasicAccess bool   `json:"basic_access"`
}

type BillingFeaturesGetResponseWire = BillingFeaturesWire

func BillingPlanToWire(plan model.Plan) BillingPlanWire {
	return BillingPlanWire{
		ID:                   plan.ID,
		Title:                plan.Title,
		Description:          plan.Description,
		Default:              plan.Default,
		Popular:              plan.Popular,
		Hidden:               plan.Hidden,
		PaddleMonthlyPriceID: plan.PaddleMonthlyPriceID,
		PaddleYearlyPriceID:  plan.PaddleYearlyPriceID,
		DiscordRoleID:        plan.DiscordRoleID,
		Features:             BillingFeaturesToWire(plan.Features, plan.ID),
	}
}

func BillingFeaturesToWire(features model.Features, planID string) BillingFeaturesWire {
	return BillingFeaturesWire{
		PlanID:      planID,
		BasicAccess: features.BasicAccess,
	}
}

type SubscriptionWire struct {
	ID                    string    `json:"id"`
	UserID                string    `json:"user_id"`
	Status                string    `json:"status"`
	PaddleSubscriptionID  string    `json:"paddle_subscription_id"`
	PaddleCustomerID      string    `json:"paddle_customer_id"`
	PaddleProductIds      []string  `json:"paddle_product_ids"`
	PaddlePriceIds        []string  `json:"paddle_price_ids"`
	CreatedAt             time.Time `json:"created_at"`
	StartedAt             null.Time `json:"started_at"`
	PausedAt              null.Time `json:"paused_at"`
	CanceledAt            null.Time `json:"canceled_at"`
	CurrentPeriodEndsAt   null.Time `json:"current_period_ends_at"`
	CurrentPeriodStartsAt null.Time `json:"current_period_starts_at"`
	UpdatedAt             time.Time `json:"updated_at"`
	Manageable            bool      `json:"manageable"`
}

type SubscriptionListResponseWire []SubscriptionWire

type SubscriptionManageResponseWire struct {
	CancelURL              string `json:"cancel_url"`
	UpdatePaymentMethodURL string `json:"update_payment_method_url"`
}

func SubscriptionToWire(subscription *model.Subscription, userID common.ID) SubscriptionWire {
	return SubscriptionWire{
		ID:                    subscription.ID.String(),
		UserID:                subscription.UserID.String(),
		Status:                subscription.Status,
		PaddleSubscriptionID:  subscription.PaddleSubscriptionID,
		PaddleCustomerID:      subscription.PaddleCustomerID,
		PaddleProductIds:      subscription.PaddleProductIds,
		PaddlePriceIds:        subscription.PaddlePriceIds,
		CreatedAt:             subscription.CreatedAt,
		StartedAt:             subscription.StartedAt,
		PausedAt:              subscription.PausedAt,
		CanceledAt:            subscription.CanceledAt,
		CurrentPeriodEndsAt:   subscription.CurrentPeriodEndsAt,
		CurrentPeriodStartsAt: subscription.CurrentPeriodStartsAt,
		UpdatedAt:             subscription.UpdatedAt,
		Manageable:            subscription.UserID == userID,
	}
}
