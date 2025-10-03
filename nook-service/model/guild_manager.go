package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
)

type GuildManagerRole string

const (
	GuildManagerRoleOwner GuildManagerRole = "admin"
)

type GuildManager struct {
	GuildID   common.ID
	UserID    common.ID
	Role      GuildManagerRole
	CreatedAt time.Time
	UpdatedAt time.Time
}

type GuildManagerWithUser struct {
	GuildManager
	User User
}
