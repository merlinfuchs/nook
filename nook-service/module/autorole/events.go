package autorole

import (
	"fmt"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
	"github.com/merlinfuchs/nook/nook-service/thing"
)

func (m *AutoroleModule) handleMemberJoin(c module.Context[AutoroleConfig], e *events.GuildMemberJoin) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	if !config.StickyRolesEnabled {
		return nil
	}

	stickyRoles, err := c.KV().Get(c.Context(), e.GuildID, stickyRolesKey(e.Member.User.ID))
	if err != nil {
		return fmt.Errorf("failed to get sticky roles: %w", err)
	}

	newRolesMap := map[common.ID]bool{}
	// Add all the sticky roles the member had before
	for _, rawRoleID := range stickyRoles.Array() {
		roleID := rawRoleID.ID()
		// Check if the role still exists
		if _, roleExists := c.Cache().Role(e.GuildID, roleID); roleExists {
			newRolesMap[roleID] = true
		}
	}

	// Remove all the roles that are in the sticky roles blacklist
	for _, roleID := range config.StickyRolesBlacklist {
		newRolesMap[roleID] = false
	}

	// Add all the roles the member already has (should be empty)
	for _, roleID := range e.Member.RoleIDs {
		newRolesMap[roleID] = true
	}

	// Add all the roles that are in the auto roles
	for _, roleID := range config.AutoRoles {
		// Check if the role still exists
		if _, roleExists := c.Cache().Role(e.GuildID, roleID); roleExists {
			newRolesMap[roleID] = true
		}
	}

	newRoles := []common.ID{}
	for roleID, shouldAssign := range newRolesMap {
		if shouldAssign {
			newRoles = append(newRoles, roleID)
		}
	}

	_, err = c.Rest().UpdateMember(e.GuildID, e.Member.User.ID, discord.MemberUpdate{
		Roles: &newRoles,
	})
	if err != nil {
		return fmt.Errorf("failed to modify guild member: %w", err)
	}

	return nil
}

func (m *AutoroleModule) handleMemberUpdate(c module.Context[AutoroleConfig], e *events.GuildMemberUpdate) error {
	config, err := c.Config(e.GuildID)
	if err != nil {
		return err
	}

	if !config.StickyRolesEnabled {
		return nil
	}

	var roles []thing.Thing
	for _, role := range e.Member.RoleIDs {
		roles = append(roles, thing.NewString(role.String()))
	}

	err = c.KV().Set(c.Context(), e.GuildID, stickyRolesKey(e.Member.User.ID), thing.NewArray(roles))
	if err != nil {
		return fmt.Errorf("failed to set sticky roles: %w", err)
	}

	return nil
}

func stickyRolesKey(userID common.ID) string {
	return fmt.Sprintf("sticky_roles:%d", userID)
}
