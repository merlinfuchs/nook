package nook

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/sharding"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/rs/cors"
)

const APIVersion = "v1"
const APIPath = "/" + APIVersion

type NookConfig struct {
	Host string
	Port int

	APIPublicBaseURL string
	AppPublicBaseURL string

	DiscordToken string
	PublicKey    string
}

type Nook struct {
	ctx context.Context

	cfg    NookConfig
	client bot.Client

	moduleManager        *manager.ModuleManager
	guildSettingsManager *manager.GuildSettingsManager
	moduleValueManager   *manager.ModuleValueManager

	httpMux    *http.ServeMux
	httpServer *http.Server
	apiGroup   api.HandlerGroup

	moduleDispatcher *ModuleDispatcher
	moduleDebouncers *common.DebounceGroup
}

func NewNook(
	ctx context.Context,
	cfg NookConfig,
	moduleManager *manager.ModuleManager,
	guildSettingsManager *manager.GuildSettingsManager,
	moduleValueManager *manager.ModuleValueManager,
) (*Nook, error) {
	httpMux := http.NewServeMux()
	httpHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{cfg.AppPublicBaseURL},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
	}).Handler(httpMux)

	httpAddress := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	httpServer := &http.Server{
		Addr:    httpAddress,
		Handler: httpHandler,
	}

	moduleDispatcher := NewModuleDispatcher(ctx, moduleManager)

	client, err := disgo.New(cfg.DiscordToken,
		bot.WithShardManagerConfigOpts(
			sharding.WithAutoScaling(false),
			sharding.WithGatewayConfigOpts(
				gateway.WithIntents(
					gateway.IntentGuilds,
					gateway.IntentGuildMembers,
					gateway.IntentGuildExpressions,
					gateway.IntentGuildMessages,
					gateway.IntentDirectMessages,
					gateway.IntentMessageContent,
					gateway.IntentGuildModeration,
				),
				gateway.WithPresenceOpts(
					gateway.WithCustomActivity("🏕️ nooks.chat"),
				),
			),
		),
		bot.WithEventManagerConfigOpts(
			bot.WithAsyncEventsEnabled(),
		),
		bot.WithCacheConfigOpts(
			cache.WithCaches(
				cache.FlagGuilds,
				cache.FlagChannels,
				cache.FlagRoles,
				cache.FlagEmojis,
			),
		),
		bot.WithEventListenerFunc(func(e *events.Ready) {
			slog.Info(
				"Shard is ready",
				slog.String("shard_id", strconv.Itoa(e.ShardID())),
				slog.String("user_id", e.User.ID.String()),
				slog.String("username", e.User.Username),
			)
		}),
		bot.WithEventListeners(moduleDispatcher),
	)
	if err != nil {
		return nil, fmt.Errorf("error on creating client: %w", err)
	}

	n := &Nook{
		ctx:                  ctx,
		cfg:                  cfg,
		moduleManager:        moduleManager,
		guildSettingsManager: guildSettingsManager,
		moduleValueManager:   moduleValueManager,
		client:               client,
		httpMux:              httpMux,
		httpServer:           httpServer,
		apiGroup:             api.Group(httpMux, APIPath),

		moduleDispatcher: moduleDispatcher,
		moduleDebouncers: common.NewDebounceGroup(10 * time.Second),
	}

	moduleManager.AddModuleUpdateListener(n.handleModuleUpdate)

	return n, nil
}

func (o *Nook) Start(ctx context.Context) error {
	if err := o.client.OpenShardManager(ctx); err != nil {
		return err
	}

	go func() {
		if err := o.httpServer.ListenAndServe(); err != nil {
			slog.Error("error on listening http server", slog.Any("err", err))
		}
	}()

	return nil
}

func (o *Nook) Close(ctx context.Context) {
	o.client.Close(ctx)
}

func (o *Nook) Client() bot.Client {
	return o.client
}

func (o *Nook) APIURL() string {
	return o.cfg.APIPublicBaseURL + APIPath
}
