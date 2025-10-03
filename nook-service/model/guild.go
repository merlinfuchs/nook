package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type Guild struct {
	ID          common.ID
	Name        string
	Description null.String
	Icon        null.String
	Unavailable bool
	Deleted     bool
	OwnerUserID common.ID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
