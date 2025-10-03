package model

import (
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type Session struct {
	KeyHash string
	UserID  common.ID

	TokenAccess    string
	TokenRefresh   string
	TokenScopes    []string
	TokenExpiresAt time.Time

	Guilds []SessionGuild

	CreatedAt time.Time
	ExpiresAt time.Time
}

type SessionGuild struct {
	ID          common.ID           `json:"id"`
	Name        string              `json:"n,omitzero"`
	Icon        null.String         `json:"i,omitzero"`
	Owner       bool                `json:"o,omitzero"`
	Permissions discord.Permissions `json:"p,omitzero"`
}
