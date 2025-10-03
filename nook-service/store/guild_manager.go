package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type GuildManagerStore interface {
	UpsertGuildManager(ctx context.Context, manager model.GuildManager) (*model.GuildManager, error)
	GetGuildManager(ctx context.Context, guildID common.ID, userID common.ID) (*model.GuildManager, error)
	DeleteGuildManager(ctx context.Context, guildID common.ID, userID common.ID) error
	GetGuildManagers(ctx context.Context, guildID common.ID) ([]model.GuildManager, error)
	GetGuildManagersWithUsers(ctx context.Context, guildID common.ID) ([]model.GuildManagerWithUser, error)
}
