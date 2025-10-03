package session

import (
	"github.com/merlinfuchs/nook/nook-service/api"
)

func (m *SessionManager) RequireSession(next api.HandlerFunc) api.HandlerFunc {
	return func(c *api.Context) error {
		session, err := m.Session(c)
		if err != nil {
			return err
		}

		if session == nil {
			return api.ErrUnauthorized("unauthorized", "Session required")
		}

		c.Session = session
		return next(c)
	}
}

func (m *SessionManager) OptionalSession(next api.HandlerFunc) api.HandlerFunc {
	return func(c *api.Context) error {
		session, err := m.Session(c)
		if err != nil {
			return err
		}

		c.Session = session
		return next(c)
	}
}
