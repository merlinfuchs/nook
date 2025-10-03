package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type Subscription struct {
	ID                    common.ID
	UserID                common.ID
	Status                string
	PaddleSubscriptionID  string
	PaddleCustomerID      string
	PaddleProductIds      []string
	PaddlePriceIds        []string
	CreatedAt             time.Time
	StartedAt             null.Time
	PausedAt              null.Time
	CanceledAt            null.Time
	CurrentPeriodEndsAt   null.Time
	CurrentPeriodStartsAt null.Time
	UpdatedAt             time.Time
}

func (s Subscription) IsActive() bool {
	return s.Status == "active"
}
