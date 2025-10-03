package user

import (
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/store"
)

type UserModule struct {
	userStore      store.UserStore
	sessionManager *session.SessionManager
}

func NewUserModule(
	userStore store.UserStore,
	sessionManager *session.SessionManager,
) *UserModule {
	return &UserModule{
		userStore:      userStore,
		sessionManager: sessionManager,
	}
}

func (m *UserModule) ModuleID() string {
	return "user"
}

func (m *UserModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "User",
		Description:    "Exposes information about users to the API.",
		Icon:           "user-round",
		Internal:       true,
		DefaultEnabled: true,
	}
}

func (m *UserModule) Endpoints(mx api.HandlerGroup) {
	usersGroup := mx.Group("/users", m.sessionManager.RequireSession)

	usersGroup.Get("/{userID}", api.Typed(func(c *api.Context) (*UserGetResponseWire, error) {
		rawUserID := c.Param("userID")
		var userID common.ID
		if rawUserID == "@me" {
			userID = c.Session.UserID
		} else {
			var err error
			userID, err = common.ParseID(rawUserID)
			if err != nil {
				return nil, api.ErrBadRequest("invalid_user_id", "Invalid user ID")
			}
		}

		user, err := m.userStore.User(c.Context(), userID)
		if err != nil {
			return nil, err
		}

		return UserToWire(user), nil
	}))
}

func (m *UserModule) Router() module.Router[UserConfig] {
	return module.NewRouter[UserConfig]()
}
