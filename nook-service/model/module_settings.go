package model

import (
	"encoding/json"
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
)

type ModuleSettings struct {
	GuildID           common.ID
	ModuleID          string
	Enabled           bool
	CommandOverwrites map[string]ModuleCommandOverwrite
	Config            json.RawMessage
	UpdatedAt         time.Time
}

type ModuleCommandOverwrite struct {
	Disabled bool `json:"disabled"`
}
