package user

import (
	"time"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"gopkg.in/guregu/null.v4"
)

type UserWire struct {
	ID            common.ID   `json:"id"`
	Username      string      `json:"username"`
	Discriminator string      `json:"discriminator"`
	DisplayName   string      `json:"display_name"`
	Avatar        null.String `json:"avatar"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
}

type UserGetResponseWire = UserWire

func UserToWire(user *model.User) *UserWire {
	if user == nil {
		return nil
	}

	return &UserWire{
		ID:            user.ID,
		Username:      user.Username,
		Discriminator: user.Discriminator,
		DisplayName:   user.DisplayName,
		Avatar:        user.Avatar,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
