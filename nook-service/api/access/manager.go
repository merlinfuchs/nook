package access

import (
	"github.com/merlinfuchs/nook/nook-service/store"
)

type AccessManager struct {
	guildStore store.GuildStore
}

func NewAccessManager(
	guildStore store.GuildStore,
) *AccessManager {
	return &AccessManager{
		guildStore: guildStore,
	}
}
