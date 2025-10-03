package guild

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module/user"
	"github.com/merlinfuchs/nook/nook-service/store"
	"gopkg.in/guregu/null.v4"
)

func (m *GuildModule) Endpoints(mx api.HandlerGroup) {
	guildsGroup := mx.Group("/guilds", m.sessionManager.RequireSession)
	guildGroup := guildsGroup.Group("/{guildID}", m.accessManager.GuildAccess)

	guildsGroup.Get("/", api.Typed(m.handleListGuilds))
	guildGroup.Get("/", api.Typed(m.handleGetGuild))
	guildGroup.Get("/channels", api.Typed(m.handleListChannels))
	guildGroup.Get("/roles", api.Typed(m.handleListRoles))
	guildGroup.Get("/settings", api.Typed(m.handleGetSettings))
	guildGroup.Put("/settings", api.TypedWithBody(m.handleUpdateSettings))
	guildGroup.Get("/managers", api.Typed(m.handleListManagers))
	guildGroup.Post("/managers", api.TypedWithBody(m.handleAddManager))
	guildGroup.Delete("/managers/{userID}", api.Typed(m.handleRemoveManager))
	guildGroup.Get("/profile", api.Typed(m.handleGetProfile))
	guildGroup.Put("/profile", api.TypedWithBody(m.handleUpdateProfile))
}

func (m *GuildModule) handleListGuilds(c *api.Context) (*GuildListResponseWire, error) {
	guilds := []GuildWire{}
	for _, guild := range c.Session.Guilds {
		_, bot := m.cache.Guild(guild.ID)

		guilds = append(guilds, GuildWire{
			ID:          guild.ID,
			Name:        guild.Name,
			Icon:        guild.Icon,
			Owner:       guild.Owner,
			Permissions: guild.Permissions,
			Bot:         bot,
			Access:      guild.Permissions.Has(discord.PermissionAdministrator),
		})
	}

	return &guilds, nil
}

func (m *GuildModule) handleGetGuild(c *api.Context) (*GuildGetResponseWire, error) {
	for _, g := range c.Session.Guilds {
		if g.ID == c.Guild.ID {
			_, bot := m.cache.Guild(g.ID)

			return &GuildWire{
				ID:          g.ID,
				Name:        g.Name,
				Icon:        g.Icon,
				Owner:       g.Owner,
				Permissions: g.Permissions,
				Bot:         bot,
				Access:      g.Permissions.Has(discord.PermissionAdministrator),
			}, nil
		}
	}

	return nil, api.ErrNotFound("unknown_guild", "Guild not found")
}

func (m *GuildModule) handleListChannels(c *api.Context) (*ChannelListResponseWire, error) {
	channels := []ChannelWire{}
	m.cache.ChannelsForEach(func(channel discord.GuildChannel) {
		if channel.GuildID() == c.Guild.ID {
			channels = append(channels, ChannelWire{
				ID:       channel.ID(),
				Type:     channel.Type(),
				Name:     channel.Name(),
				Position: channel.Position(),
			})
		}
	})

	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})

	return &channels, nil
}

func (m *GuildModule) handleListRoles(c *api.Context) (*RoleListResponseWire, error) {
	channels := []RoleWire{}
	m.cache.RolesForEach(c.Guild.ID, func(role discord.Role) {
		color := "#979c9f"
		if role.Color != 0 {
			color = fmt.Sprintf("#%06X", role.Color)
		}

		channels = append(channels, RoleWire{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: role.Permissions,
			Position:    role.Position,
			Color:       color,
		})
	})

	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Position < channels[j].Position
	})

	return &channels, nil
}

func (m *GuildModule) handleGetSettings(c *api.Context) (*GuildSettingsGetResponseWire, error) {
	settings, err := m.guildSettingsManager.GuildSettings(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	return m.guildSettingsToWire(settings), nil
}

func (m *GuildModule) handleUpdateSettings(c *api.Context, req GuildSettingsUpdateRequestWire) (*GuildSettingsUpdateResponseWire, error) {
	settings, err := m.guildSettingsManager.GuildSettings(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	if req.CommandPrefix != nil {
		commandPrefix := strings.TrimSpace(*req.CommandPrefix)
		settings.CommandPrefix = null.NewString(commandPrefix, commandPrefix != "")
	}

	if req.ColorScheme != nil {
		colorScheme := strings.TrimSpace(*req.ColorScheme)
		settings.ColorScheme = null.NewString(colorScheme, colorScheme != "")
	}

	err = m.guildSettingsManager.UpdateGuildSettings(c.Context(), c.Guild.ID, model.GuildSettings{
		GuildID:       c.Guild.ID,
		CommandPrefix: settings.CommandPrefix,
		ColorScheme:   settings.ColorScheme,
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	return m.guildSettingsToWire(settings), nil
}

func (m *GuildModule) handleListManagers(c *api.Context) (*GuildManagerListResponseWire, error) {
	managers, err := m.guildManagerStore.GetGuildManagersWithUsers(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	res := make([]GuildManagerWire, 0, len(managers))
	for _, manager := range managers {
		res = append(res, GuildManagerWire{
			GuildID:   manager.GuildID,
			User:      user.UserWire(manager.User),
			Role:      manager.Role,
			CreatedAt: manager.CreatedAt,
			UpdatedAt: manager.UpdatedAt,
		})
	}

	return &res, nil
}

func (m *GuildModule) handleAddManager(c *api.Context, req GuildManagerCreateAddWire) (*GuildManagerAddAddResponseWire, error) {
	u, err := m.userStore.User(c.Context(), req.UserID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			u, err = m.tryUpsertUser(c.Context(), req.UserID)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	manager, err := m.guildManagerStore.UpsertGuildManager(c.Context(), model.GuildManager{
		GuildID:   c.Guild.ID,
		UserID:    req.UserID,
		Role:      req.Role,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}

	return &GuildManagerAddAddResponseWire{
		GuildID:   c.Guild.ID,
		User:      user.UserWire(*u),
		Role:      req.Role,
		CreatedAt: manager.CreatedAt,
		UpdatedAt: manager.UpdatedAt,
	}, nil
}

func (m *GuildModule) handleRemoveManager(c *api.Context) (*GuildManagerRemoveResponseWire, error) {
	userID, err := common.ParseID(c.Param("userID"))
	if err != nil {
		return nil, api.ErrBadRequest("invalid_user_id", "Invalid user ID")
	}

	if userID == c.Session.UserID {
		return nil, api.ErrBadRequest("invalid_user_id", "You cannot remove yourself as a manager")
	}

	err = m.guildManagerStore.DeleteGuildManager(c.Context(), c.Guild.ID, userID)
	if err != nil {
		return nil, err
	}
	return &GuildManagerRemoveResponseWire{}, nil
}

func (m *GuildModule) handleGetProfile(c *api.Context) (*GuildProfileGetResponseWire, error) {
	member, err := m.guildCurrentMember(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	return memberToProfileWire(member), nil
}

func (m *GuildModule) handleUpdateProfile(c *api.Context, req GuildProfileUpdateRequestWire) (*GuildProfileUpdateResponseWire, error) {
	member, err := m.updateGuildCurrentMember(c.Context(), c.Guild.ID, req.Name, req.Bio, req.Avatar)
	if err != nil {
		return nil, err
	}

	return memberToProfileWire(member), nil
}
