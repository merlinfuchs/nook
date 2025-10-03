package store

import (
	"context"

	"github.com/merlinfuchs/nook/nook-service/model"
)

type SessionStore interface {
	CreateSession(ctx context.Context, session *model.Session) error
	DeleteSession(ctx context.Context, keyHash string) error
	Session(ctx context.Context, keyHash string) (*model.Session, error)
}
