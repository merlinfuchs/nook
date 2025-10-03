package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

type UserStore interface {
	User(ctx context.Context, id common.ID) (*model.User, error)
	UpsertUser(ctx context.Context, user *model.User) (*model.User, error)
}
