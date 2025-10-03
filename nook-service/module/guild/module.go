package guild

import (
	"fmt"
	"time"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/api/access"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/store"
)

type GuildModule struct {
	cache                cache.Caches
	rest                 rest.Rest
	moduleStore          store.ModuleStore
	guildStore           store.GuildStore
	guildManagerStore    store.GuildManagerStore
	userStore            store.UserStore
	sessionManager       *session.SessionManager
	accessManager        *access.AccessManager
	guildSettingsManager *manager.GuildSettingsManager
}

func NewGuildModule(
	cache cache.Caches,
	rest rest.Rest,
	moduleStore store.ModuleStore,
	guildStore store.GuildStore,
	guildManagerStore store.GuildManagerStore,
	userStore store.UserStore,
	sessionManager *session.SessionManager,
	accessManager *access.AccessManager,
	guildSettingsManager *manager.GuildSettingsManager,
) *GuildModule {
	return &GuildModule{
		cache:                cache,
		rest:                 rest,
		moduleStore:          moduleStore,
		guildStore:           guildStore,
		guildManagerStore:    guildManagerStore,
		userStore:            userStore,
		sessionManager:       sessionManager,
		accessManager:        accessManager,
		guildSettingsManager: guildSettingsManager,
	}
}

func (m *GuildModule) ModuleID() string {
	return "guild"
}

func (m *GuildModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Guild",
		Description:    "Exposes information about guilds to the API.",
		Icon:           "building-2",
		Internal:       true,
		DefaultEnabled: true,
	}
}

func (m *GuildModule) Router() module.Router[GuildConfig] {
	return module.NewRouter[GuildConfig]().
		Handle(module.ListenerFunc(m.handleGuildReady)).
		Handle(module.ListenerFunc(m.handleGuildUnavailable)).
		Handle(module.ListenerFunc(m.handleGuildLeave))
}

func (m *GuildModule) handleGuildReady(c module.Context[GuildConfig], e *events.GuildReady) error {
	_, err := m.guildStore.UpsertGuild(c.Context(), model.Guild{
		ID:          e.Guild.ID,
		Name:        e.Guild.Name,
		Description: common.PtrToNullString(e.Guild.Description),
		Icon:        common.PtrToNullString(e.Guild.Icon),
		Unavailable: e.Guild.Unavailable,
		OwnerUserID: e.Guild.OwnerID,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		return fmt.Errorf("failed to upsert guild: %w", err)
	}
	return nil
}

func (m *GuildModule) handleGuildUnavailable(c module.Context[GuildConfig], e *events.GuildUnavailable) error {
	err := m.guildStore.UpdateGuildUnavailable(c.Context(), e.Guild.ID, true)
	if err != nil {
		return fmt.Errorf("failed to update guild unavailable: %w", err)
	}
	return nil
}

func (m *GuildModule) handleGuildLeave(c module.Context[GuildConfig], e *events.GuildLeave) error {
	err := m.guildStore.UpdateGuildUnavailable(c.Context(), e.Guild.ID, true)
	if err != nil {
		return fmt.Errorf("failed to update guild unavailable: %w", err)
	}
	return nil
}
