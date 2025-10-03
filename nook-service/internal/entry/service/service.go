package service

import (
	"context"
	"fmt"
	"time"

	"github.com/merlinfuchs/nook/nook-service/api/access"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/config"
	"github.com/merlinfuchs/nook/nook-service/internal/db/postgres"
	"github.com/merlinfuchs/nook/nook-service/internal/logging"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module/auth"
	"github.com/merlinfuchs/nook/nook-service/module/autorole"
	"github.com/merlinfuchs/nook/nook-service/module/billing"
	"github.com/merlinfuchs/nook/nook-service/module/counting"
	"github.com/merlinfuchs/nook/nook-service/module/guild"
	loggingmodule "github.com/merlinfuchs/nook/nook-service/module/logging"
	"github.com/merlinfuchs/nook/nook-service/module/manage"
	"github.com/merlinfuchs/nook/nook-service/module/moderation"
	"github.com/merlinfuchs/nook/nook-service/module/ping"
	"github.com/merlinfuchs/nook/nook-service/module/user"
	"github.com/merlinfuchs/nook/nook-service/module/welcome"
	"github.com/merlinfuchs/nook/nook-service/nook"
)

func RunService(ctx context.Context) error {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	logging.SetupLogger(logging.LoggerConfig(cfg.Logging))

	pg, err := postgres.New(postgres.ClientConfig{
		Host:     cfg.Database.Postgres.Host,
		Port:     cfg.Database.Postgres.Port,
		DBName:   cfg.Database.Postgres.DBName,
		User:     cfg.Database.Postgres.User,
		Password: cfg.Database.Postgres.Password,
	})
	if err != nil {
		return fmt.Errorf("failed to create postgres client: %w", err)
	}

	moduleManager := manager.NewModuleManager(pg)
	moduleValueManager := manager.NewModuleValueManager(pg)
	guildSettingsManager := manager.NewGuildSettingsManager(pg, manager.GuildSettingsManagerConfig{
		DefaultPrefix:      cfg.Defaults.CommandPrefix,
		DefaultColorScheme: cfg.Defaults.ColorScheme,
	})

	orb, err := nook.NewNook(ctx, nook.NookConfig{
		Host:             cfg.API.Host,
		Port:             cfg.API.Port,
		APIPublicBaseURL: cfg.API.PublicBaseURL,
		AppPublicBaseURL: cfg.App.PublicBaseURL,
		DiscordToken:     cfg.Discord.BotToken,
		PublicKey:        cfg.Discord.PublicKey,
	}, moduleManager, guildSettingsManager, moduleValueManager)
	if err != nil {
		return fmt.Errorf("failed to create nook: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sessionManager := session.NewSessionManager(session.SessionManagerConfig{
		StrictCookies: cfg.API.StrictCookies,
		SecureCookies: cfg.API.SecureCookies,
	}, pg)
	accessManager := access.NewAccessManager(pg)
	planManager := manager.NewPlanManager(pg, pg, pg, getConfigPlans(cfg))

	nook.AddModule(orb, guild.NewGuildModule(
		orb.Client().Caches(),
		orb.Client().Rest(),
		moduleManager,
		pg,
		pg,
		pg,
		sessionManager,
		accessManager,
		guildSettingsManager,
	))
	nook.AddModule(orb, user.NewUserModule(pg, sessionManager))
	nook.AddModule(orb, auth.NewAuthModule(
		auth.AuthModuleConfig{
			ClientID:      cfg.Discord.ClientID,
			ClientSecret:  cfg.Discord.ClientSecret,
			APIURL:        orb.APIURL(),
			AppURL:        cfg.App.PublicBaseURL,
			SecureCookies: cfg.API.SecureCookies,
			StrictCookies: cfg.API.StrictCookies,
		},
		pg,
		pg,
		sessionManager,
	))
	nook.AddModule(orb, billing.NewBillingModule(
		billing.BillingModuleConfig{
			PaddleWebhookSecret: cfg.Billing.PaddleWebhookSecret,
			PaddleAPIKey:        cfg.Billing.PaddleAPIKey,
			PaddleEnvironment:   cfg.Billing.PaddleEnvironment,
		},
		planManager,
		sessionManager,
		accessManager,
		pg,
		pg,
		pg,
	))
	nook.AddModule(orb, manage.NewManageModule(
		ctx,
		sessionManager,
		accessManager,
		moduleManager,
	))
	nook.AddModule(orb, ping.NewPingModule())
	nook.AddModule(orb, moderation.NewModerationModule())
	nook.AddModule(orb, counting.NewCountingModule())
	nook.AddModule(orb, loggingmodule.NewLoggingModule())
	nook.AddModule(orb, autorole.NewAutoroleModule())
	nook.AddModule(orb, welcome.NewWelcomeModule())

	if err := orb.Start(ctx); err != nil {
		return fmt.Errorf("failed to open gateway: %w", err)
	}

	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer closeCancel()
		orb.Client().Close(closeCtx)
	}()

	<-ctx.Done()
	return nil
}

func getConfigPlans(cfg *config.Config) []model.Plan {
	plans := make([]model.Plan, len(cfg.Billing.Plans))
	for i, plan := range cfg.Billing.Plans {
		plans[i] = model.Plan{
			ID:                   plan.ID,
			Title:                plan.Title,
			Description:          plan.Description,
			Default:              plan.Default,
			Popular:              plan.Popular,
			Hidden:               plan.Hidden,
			PaddleMonthlyPriceID: plan.PaddleMonthlyPriceID,
			PaddleYearlyPriceID:  plan.PaddleYearlyPriceID,
			DiscordRoleID:        plan.DiscordRoleID,
			Features: model.Features{
				BasicAccess: plan.FeatureBasicAccess,
			},
		}
	}
	return plans
}
