package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type Entitlement struct {
	ID             common.ID
	GuildID        common.ID
	SubscriptionID common.NullID
	Type           string
	PlanIDs        []string
	StartsAt       null.Time
	EndsAt         null.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
