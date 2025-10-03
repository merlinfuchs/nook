package billing

import (
	"errors"
	"fmt"

	"github.com/PaddleHQ/paddle-go-sdk"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/api/access"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/store"
)

type BillingModuleConfig struct {
	PaddleWebhookSecret string
	PaddleAPIKey        string
	PaddleEnvironment   string
}

type BillingModule struct {
	cfg BillingModuleConfig

	planManager       *manager.PlanManager
	sessionManager    *session.SessionManager
	accessManager     *access.AccessManager
	entitlementStore  store.EntitlementStore
	subscriptionStore store.SubscriptionStore
	userStore         store.UserStore

	client          *paddle.SDK
	webhookVerifier *paddle.WebhookVerifier
}

func NewBillingModule(
	cfg BillingModuleConfig,
	planManager *manager.PlanManager,
	sessionManager *session.SessionManager,
	accessManager *access.AccessManager,
	entitlementStore store.EntitlementStore,
	subscriptionStore store.SubscriptionStore,
	userStore store.UserStore,
) *BillingModule {
	webhookVerifier := paddle.NewWebhookVerifier(cfg.PaddleWebhookSecret)

	options := []paddle.Option{}
	if cfg.PaddleEnvironment == "sandbox" {
		options = append(options, paddle.WithBaseURL(paddle.SandboxBaseURL))
	} else {
		options = append(options, paddle.WithBaseURL(paddle.ProductionBaseURL))
	}

	client, err := paddle.New(cfg.PaddleAPIKey, options...)
	if err != nil {
		panic(fmt.Errorf("failed to create paddle client: %w", err))
	}

	return &BillingModule{
		cfg: cfg,

		planManager:       planManager,
		sessionManager:    sessionManager,
		accessManager:     accessManager,
		entitlementStore:  entitlementStore,
		subscriptionStore: subscriptionStore,
		userStore:         userStore,

		client:          client,
		webhookVerifier: webhookVerifier,
	}
}

func (m *BillingModule) ModuleID() string {
	return "billing"
}

func (m *BillingModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Billing",
		Description:    "Provides billing functionality for the bot.",
		Icon:           "credit-card",
		Internal:       true,
		DefaultEnabled: true,
	}
}

func (m *BillingModule) Endpoints(mx api.HandlerGroup) {
	mx.Get("/billing/plans", api.Typed(m.handleBillingPlanList))
	mx.Post("/billing/webhooks", m.handleBillingWebhook)

	guildGroup := mx.Group("/guilds/{guildID}", m.sessionManager.RequireSession, m.accessManager.GuildAccess)

	guildGroup.Get("/billing/features", api.Typed(m.handleFeaturesGet))
	guildGroup.Get("/billing/subscriptions", api.Typed(m.handleAppSubscriptionList))
	guildGroup.Get("/billing/subscriptions/{subscriptionID}/manage", api.Typed(m.handleAppSubscriptionManage))
}

func (m *BillingModule) handleBillingPlanList(c *api.Context) (*BillingPlanListResponseWire, error) {
	plans := m.planManager.Plans()

	res := make(BillingPlanListResponseWire, len(plans))
	for i, plan := range plans {
		res[i] = BillingPlanToWire(plan)
	}

	return &res, nil
}

func (m *BillingModule) handleFeaturesGet(c *api.Context) (*BillingFeaturesGetResponseWire, error) {
	features, planID := m.planManager.GuildFeatures(c.Context(), c.Guild.ID)

	res := BillingFeaturesToWire(features, planID)
	return &res, nil
}

func (m *BillingModule) handleAppSubscriptionList(c *api.Context) (*SubscriptionListResponseWire, error) {
	subscriptions, err := m.subscriptionStore.SubscriptionsByGuildID(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	res := make(SubscriptionListResponseWire, len(subscriptions))
	for i, subscription := range subscriptions {
		res[i] = SubscriptionToWire(subscription, c.Session.UserID)
	}

	return &res, nil
}

func (m *BillingModule) handleAppSubscriptionManage(c *api.Context) (*SubscriptionManageResponseWire, error) {
	subscriptionID, err := common.ParseID(c.Param("subscriptionID"))
	if err != nil {
		return nil, err
	}

	subscription, err := m.subscriptionStore.SubscriptionByID(c.Context(), subscriptionID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, api.ErrNotFound("unknown_subscription", "The subscription you are trying to manage does not exist.")
		}
		return nil, err
	}

	if subscription.UserID != c.Session.UserID {
		return nil, api.ErrForbidden("missing_access", "You are not authorized to manage this subscription.")
	}

	sub, err := m.client.SubscriptionsClient.GetSubscription(c.Context(), &paddle.GetSubscriptionRequest{
		SubscriptionID: subscription.PaddleSubscriptionID,
	})
	if err != nil {
		return nil, err
	}

	return &SubscriptionManageResponseWire{
		CancelURL:              sub.ManagementURLs.Cancel,
		UpdatePaymentMethodURL: *sub.ManagementURLs.UpdatePaymentMethod,
	}, nil
}

func (m *BillingModule) Router() module.Router[BillingConfig] {
	return module.NewRouter[BillingConfig]()
}
