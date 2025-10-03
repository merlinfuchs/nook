package model

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"gopkg.in/guregu/null.v4"
)

type User struct {
	ID            common.ID
	Username      string
	Discriminator string
	DisplayName   string
	Avatar        null.String
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
