package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/store"
)

const SessionCookieName = "nook-session"

type CreateSessionOpts struct {
	UserID         common.ID
	TokenAccess    string
	TokenRefresh   string
	TokenScopes    []string
	TokenExpiresAt time.Time
	Guilds         []model.SessionGuild
}

type SessionManagerConfig struct {
	StrictCookies bool
	SecureCookies bool
}

type SessionManager struct {
	config       SessionManagerConfig
	sessionStore store.SessionStore
}

func NewSessionManager(config SessionManagerConfig, sessionStore store.SessionStore) *SessionManager {
	return &SessionManager{
		config:       config,
		sessionStore: sessionStore,
	}
}

func (s *SessionManager) CreateSessionCookie(c *api.Context, opts CreateSessionOpts) (string, *model.Session, error) {
	key, session, err := s.CreateSession(c.Context(), opts)
	if err != nil {
		return "", nil, err
	}

	sameSite := http.SameSiteNoneMode
	if s.config.StrictCookies {
		sameSite = http.SameSiteStrictMode
	}

	c.SetCookie(&http.Cookie{
		Name:     SessionCookieName,
		Value:    key,
		Secure:   s.config.SecureCookies,
		HttpOnly: true,
		SameSite: sameSite,
		MaxAge:   int(session.ExpiresAt.Sub(session.CreatedAt).Seconds()),
		Path:     "/",
	})

	return key, session, nil
}

func (s *SessionManager) CreateSession(ctx context.Context, opts CreateSessionOpts) (string, *model.Session, error) {
	key := common.SecureKey()
	keyHash := common.HashKey(key)

	session := &model.Session{
		KeyHash:        keyHash,
		UserID:         opts.UserID,
		CreatedAt:      time.Now().UTC(),
		TokenAccess:    opts.TokenAccess,
		TokenRefresh:   opts.TokenRefresh,
		TokenScopes:    opts.TokenScopes,
		TokenExpiresAt: opts.TokenExpiresAt,
		ExpiresAt:      opts.TokenExpiresAt,
		Guilds:         opts.Guilds,
	}

	err := s.sessionStore.CreateSession(ctx, session)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	return key, session, nil
}

func (s *SessionManager) DeleteSession(c *api.Context) error {
	defer c.DeleteCookie(SessionCookieName)

	key := c.Cookie(SessionCookieName)
	if key == "" {
		return nil
	}

	keyHash := common.HashKey(key)
	if err := s.sessionStore.DeleteSession(c.Context(), keyHash); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil
		}
		return fmt.Errorf("failed to delete session: %w", err)
	}

	return nil
}

func (s *SessionManager) Session(c *api.Context) (*model.Session, error) {
	key := c.Cookie(SessionCookieName)
	if key == "" {
		return nil, nil
	}

	keyHash := common.HashKey(key)

	session, err := s.sessionStore.Session(c.Context(), keyHash)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return session, nil
}
