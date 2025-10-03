package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type GuildSettings struct {
	GuildID       common.ID
	CommandPrefix null.String
	ColorScheme   null.String
	UpdatedAt     time.Time
}

type ResolvedGuildSettings struct {
	CommandPrefix string
	ColorScheme   string
}

func (g ResolvedGuildSettings) Color() int {
	return 0xc3a9fe
}
