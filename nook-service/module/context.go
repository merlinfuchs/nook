package module

import (
	"context"
	"encoding/json"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

type GenericContext interface {
	Context() context.Context
	ModuleID() string

	Rest() rest.Rest
	Cache() cache.Caches
	Client() bot.Client
	KV() KV
	Config(guildID common.ID) (json.RawMessage, error)
	GuildSettings(guildID common.ID) (model.ResolvedGuildSettings, error)
}

type Context[C any] interface {
	Context() context.Context
	ModuleID() string

	Rest() rest.Rest
	Cache() cache.Caches
	Client() bot.Client
	KV() KV
	Config(guildID common.ID) (C, error)
	GuildSettings(guildID common.ID) (model.ResolvedGuildSettings, error)
}

type KV interface {
	Get(ctx context.Context, guildID common.ID, key string) (thing.Thing, error)
	Set(ctx context.Context, guildID common.ID, key string, value thing.Thing) error
	Update(ctx context.Context, guildID common.ID, op thing.Operation, key string, value thing.Thing) (thing.Thing, error)
	Delete(ctx context.Context, guildID common.ID, key string) error
}

type contextImpl[C any] struct {
	GenericContext
}

func (c *contextImpl[C]) Config(guildID common.ID) (C, error) {
	var res C

	rawConfig, err := c.GenericContext.Config(guildID)
	if err != nil {
		return res, err
	}

	err = json.Unmarshal(rawConfig, &res)
	if err != nil {
		return res, err
	}

	return res, nil
}

func NewContext[C any](c GenericContext) Context[C] {
	return &contextImpl[C]{
		GenericContext: c,
	}
}
