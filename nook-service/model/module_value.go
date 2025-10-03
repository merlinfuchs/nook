package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

type ModuleValue struct {
	ID        uint64      `json:"id"`
	GuildID   common.ID   `json:"guild_id"`
	ModuleID  string      `json:"module_id"`
	Key       string      `json:"key"`
	Value     thing.Thing `json:"value"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
