package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type GuildStore interface {
	UpsertGuild(ctx context.Context, guild model.Guild) (*model.Guild, error)
	UpdateGuildDeleted(ctx context.Context, guildID common.ID, deleted bool) error
	UpdateGuildUnavailable(ctx context.Context, guildID common.ID, unavailable bool) error
	Guild(ctx context.Context, guildID common.ID) (*model.Guild, error)
	GuildsByOwnerUserID(ctx context.Context, ownerUserID common.ID) ([]model.Guild, error)
}
