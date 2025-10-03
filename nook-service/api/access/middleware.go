package access

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/model"
)

func (m *AccessManager) GuildAccess(next api.HandlerFunc) api.HandlerFunc {
	return func(c *api.Context) error {
		guildID, err := common.ParseID(c.Param("guildID"))
		if err != nil {
			return api.ErrBadRequest("invalid_guild_id", "Invalid guild ID")
		}

		var guild *model.SessionGuild
		for _, g := range c.Session.Guilds {
			if g.ID == guildID {
				guild = &g
				break
			}
		}

		if guild == nil || (!guild.Owner && !guild.Permissions.Has(discord.PermissionAdministrator)) {
			return api.ErrForbidden("missing_access", "Access to guild missing")
		}

		c.Guild = guild
		return next(c)
	}
}
