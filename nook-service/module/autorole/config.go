package autorole

import (
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(AutoroleConfig{})

var configUISchema = module.ConfigUISchema{
	Properties: map[string]module.ConfigUISchema{
		"auto_roles": {
			Widget:        module.ConfigUIWidgetRoleSelect,
			AllowMultiple: true,
		},
		"sticky_roles_blacklist": {
			Widget:        module.ConfigUIWidgetRoleSelect,
			AllowMultiple: true,
		},
	},
}

var defaultConfig = AutoroleConfig{}

type AutoroleConfig struct {
	AutoRoles            []common.ID `json:"auto_roles,omitzero" title:"Auto Roles" description:"The roles to assign to users when they join the server"`
	StickyRolesEnabled   bool        `json:"sticky_roles_enabled,omitzero" title:"Sticky Roles Enabled" description:"Whether to assign sticky roles to users when they join the server"`
	StickyRolesBlacklist []common.ID `json:"sticky_roles_blacklist,omitzero" title:"Sticky Roles Blacklist" description:"The roles to not assign to users when they join the server"`
}
