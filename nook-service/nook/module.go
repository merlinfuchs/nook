package nook

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

func (o *Nook) handleModuleUpdate(guildID common.ID, mod manager.ModuleWithSettings) {
	modules, err := o.moduleManager.GuildModules(o.ctx, guildID)
	if err != nil {
		slog.Error("error on getting guild modules", slog.Any("err", err))
		return
	}

	o.moduleDispatcher.Dispatch(mod.Module.ModuleID(), &module.ConfigureEvent{
		GuildID: guildID,
		Enabled: mod.Settings.Enabled,
		Config:  mod.Settings.Config,
	})

	o.moduleDebouncers.Debounce(guildID.String(), func() {
		var commands []discord.ApplicationCommandCreate

		for _, mod := range modules {
			if mod.Settings.Enabled {
				if mod, ok := mod.Module.(module.GenericModuleWithCommands); ok {
					settings, err := o.moduleManager.ModuleSettings(o.ctx, guildID, mod.ModuleID())
					if err != nil {
						slog.Error(
							"Failed to get module settings",
							slog.String("guild_id", guildID.String()),
							slog.String("module_id", mod.ModuleID()),
							slog.Any("err", err),
						)
						continue
					}

					for _, command := range mod.Commands() {
						if overwrite := settings.CommandOverwrites[command.CommandName()]; !overwrite.Disabled {
							commands = append(commands, command)
						}
					}
				}
			}
		}

		slog.Info("Deploying commands for guild", slog.String("guild_id", guildID.String()), slog.Int("commands", len(commands)))
		ctx, cancel := context.WithTimeout(o.ctx, 10*time.Second)
		defer cancel()

		_, err = o.client.Rest().SetGuildCommands(o.client.ApplicationID(), guildID, commands, rest.WithCtx(ctx))
		if err != nil {
			slog.Error("error on setting guild commands", slog.Any("err", err))
		}
	})
}

type ModuleDispatcher struct {
	ctx           context.Context
	handlers      map[string]moduleHandler
	moduleManager *manager.ModuleManager
}

func NewModuleDispatcher(ctx context.Context, moduleManager *manager.ModuleManager) *ModuleDispatcher {
	return &ModuleDispatcher{
		ctx:           ctx,
		handlers:      make(map[string]moduleHandler),
		moduleManager: moduleManager,
	}
}

func (o *ModuleDispatcher) OnEvent(e bot.Event) {
	guildID, _ := common.EventGuildID(e)
	if guildID == nil {
		return
	}

	enabledModules, err := o.moduleManager.EnabledModuleIDs(o.ctx, *guildID)
	if err != nil {
		slog.Error("error on getting enabled module IDs", slog.Any("err", err))
		return
	}

	for _, moduleID := range enabledModules {
		handler, ok := o.handlers[moduleID]
		if ok {
			handler.OnEvent(&module.DiscordEvent{
				Event: e,
			})
		}
	}
}

func (o *ModuleDispatcher) Dispatch(moduleID string, e module.Event) {
	handler, ok := o.handlers[moduleID]
	if ok {
		handler.OnEvent(e)
	}
}

func (o *ModuleDispatcher) Broadcast(e module.Event) {
	for _, handler := range o.handlers {
		handler.OnEvent(e)
	}
}

func (o *ModuleDispatcher) addHandler(handl moduleHandler) {
	o.handlers[handl.gc.ModuleID()] = handl
}

type ModuleContext[T any] struct {
	context.Context

	moduleID string

	rest          rest.Rest
	cache         cache.Caches
	moduleManager *manager.ModuleManager
}

func NewModuleContext[T any](ctx context.Context, moduleID string, rest rest.Rest, moduleManager *manager.ModuleManager) *ModuleContext[T] {
	return &ModuleContext[T]{
		Context:       ctx,
		moduleID:      moduleID,
		rest:          rest,
		moduleManager: moduleManager,
	}
}

func (c *ModuleContext[T]) Config(guildID common.ID) (T, error) {
	var res T

	rawConfig, err := c.moduleManager.ModuleSettings(c, guildID, c.moduleID)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(rawConfig.Config, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func (c *ModuleContext[T]) Rest() rest.Rest {
	return c.rest
}

func (c *ModuleContext[T]) Cache() cache.Caches {
	return c.cache
}

func AddModule[C any](o *Nook, mod module.Module[C]) {
	disp := mod.Router().Generic()

	gc := newGenericContextImpl(
		o.ctx,
		mod.ModuleID(),
		o.client,
		o.guildSettingsManager,
		o.moduleManager,
		o.moduleValueManager,
	)

	handl := moduleHandler{
		gc:     gc,
		router: disp,
	}

	o.moduleDispatcher.addHandler(handl)
	o.moduleManager.AddModule(mod)

	if mod, ok := mod.(module.GenericModuleWithEndpoints); ok {
		mod.Endpoints(o.apiGroup)
	}

	slog.Info("Module added", slog.String("id", mod.ModuleID()), slog.String("name", mod.Metadata().Name))
}

type genericContextImpl struct {
	ctx      context.Context
	moduleID string

	client               bot.Client
	kv                   *kvImpl
	guildSettingsManager *manager.GuildSettingsManager
	moduleManager        *manager.ModuleManager
}

func newGenericContextImpl(
	ctx context.Context,
	moduleID string,
	client bot.Client,
	guildSettingsManager *manager.GuildSettingsManager,
	moduleManager *manager.ModuleManager,
	moduleValueManager *manager.ModuleValueManager,
) *genericContextImpl {
	return &genericContextImpl{
		ctx:                  ctx,
		moduleID:             moduleID,
		client:               client,
		kv:                   newKVImpl(moduleID, moduleValueManager),
		guildSettingsManager: guildSettingsManager,
		moduleManager:        moduleManager,
	}
}

func (c *genericContextImpl) Context() context.Context {
	return c.ctx
}

func (c *genericContextImpl) ModuleID() string {
	return c.moduleID
}

func (c *genericContextImpl) Client() bot.Client {
	return c.client
}

func (c *genericContextImpl) Rest() rest.Rest {
	return c.client.Rest()
}

func (c *genericContextImpl) Cache() cache.Caches {
	return c.client.Caches()
}

func (c *genericContextImpl) KV() module.KV {
	return c.kv
}

func (c *genericContextImpl) ModuleManager() *manager.ModuleManager {
	return c.moduleManager
}

func (c *genericContextImpl) Config(guildID common.ID) (json.RawMessage, error) {
	settings, err := c.moduleManager.ModuleSettings(c.Context(), guildID, c.moduleID)
	if err != nil {
		return nil, err
	}

	return settings.Config, nil
}

func (c *genericContextImpl) GuildSettings(guildID common.ID) (model.ResolvedGuildSettings, error) {
	return c.guildSettingsManager.ResolvedGuildSettings(c.Context(), guildID)
}

type kvImpl struct {
	moduleID           string
	moduleValueManager *manager.ModuleValueManager
}

func newKVImpl(moduleID string, moduleValueManager *manager.ModuleValueManager) *kvImpl {
	return &kvImpl{
		moduleID:           moduleID,
		moduleValueManager: moduleValueManager,
	}
}

func (k *kvImpl) Get(ctx context.Context, guildID common.ID, key string) (thing.Thing, error) {
	return k.moduleValueManager.ModuleValue(ctx, guildID, k.moduleID, key)
}

func (k *kvImpl) Set(ctx context.Context, guildID common.ID, key string, value thing.Thing) error {
	return k.moduleValueManager.SetModuleValue(ctx, guildID, k.moduleID, key, value)
}

func (k *kvImpl) Update(ctx context.Context, guildID common.ID, op thing.Operation, key string, value thing.Thing) (thing.Thing, error) {
	return k.moduleValueManager.UpdateModuleValue(ctx, op, guildID, k.moduleID, key, value)
}

func (k *kvImpl) Delete(ctx context.Context, guildID common.ID, key string) error {
	return k.moduleValueManager.DeleteModuleValue(ctx, guildID, k.moduleID, key)
}

type moduleHandler struct {
	gc     module.GenericContext
	router module.GenericRouter
}

func (h *moduleHandler) OnEvent(e module.Event) {
	err := h.router.OnEvent(h.gc, e)
	if err != nil {
		slog.Error(
			"Failed to handle module event",
			slog.String("module_id", h.gc.ModuleID()),
			slog.String("event_type", fmt.Sprintf("%T", e)),
			slog.Any("err", err),
		)
	}
}
