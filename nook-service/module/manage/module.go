package manage

import (
	"context"
	"encoding/json"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/merlinfuchs/nook/nook-service/api/access"
	"github.com/merlinfuchs/nook/nook-service/api/session"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/model"
	"github.com/merlinfuchs/nook/nook-service/module"
)

type ManageModule struct {
	ctx context.Context

	sessionManager *session.SessionManager
	accessManager  *access.AccessManager

	moduleManager *manager.ModuleManager
}

func NewManageModule(
	ctx context.Context,
	sessionManager *session.SessionManager,
	accessManager *access.AccessManager,
	moduleManager *manager.ModuleManager,
) *ManageModule {
	return &ManageModule{
		ctx:            ctx,
		sessionManager: sessionManager,
		accessManager:  accessManager,
		moduleManager:  moduleManager,
	}
}

func (m *ManageModule) ModuleID() string {
	return "manage"
}

func (m *ManageModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Manage",
		Description:    "Manage the module settings",
		Icon:           "gear",
		Internal:       true,
		DefaultEnabled: true,
	}
}

func (m *ManageModule) Commands() []discord.ApplicationCommandCreate {
	return nil
}

func (m *ManageModule) Endpoints(mx api.HandlerGroup) {
	mx.Get("/modules", api.Typed(m.handleGlobalModuleList))

	guildGroup := mx.Group("/guilds/{guildID}", m.sessionManager.RequireSession, m.accessManager.GuildAccess)
	guildGroup.Get("/modules", api.Typed(m.handleModuleList))
	guildGroup.Get("/modules/{moduleId}", api.Typed(m.handleModuleGet))
	guildGroup.Put("/modules/{moduleId}", api.TypedWithBody(m.handleModuleConfigure))
}

func (m *ManageModule) Router() module.Router[ManageConfig] {
	return module.NewRouter[ManageConfig]()
}

func (m *ManageModule) handleGlobalModuleList(c *api.Context) (*GlobalModuleListResponseWire, error) {
	modules := m.moduleManager.Modules()

	res := make([]GlobalModuleWire, 0, len(modules))
	for _, module := range modules {
		metadata := module.Metadata()

		res = append(res, GlobalModuleWire{
			ID:       module.ModuleID(),
			Metadata: metadata,
		})
	}

	return &res, nil
}

func (m *ManageModule) handleModuleList(c *api.Context) (*ModuleListResponseWire, error) {
	modules, err := m.moduleManager.GuildModules(c.Context(), c.Guild.ID)
	if err != nil {
		return nil, err
	}

	res := make([]ModuleWire, 0, len(modules))
	for _, mod := range modules {
		metadata := mod.Module.Metadata()

		var commands []discord.ApplicationCommandCreate
		if mod, ok := mod.Module.(module.GenericModuleWithCommands); ok {
			commands = mod.Commands()
		}

		res = append(res, ModuleWire{
			ID:       mod.Module.ModuleID(),
			Metadata: metadata,
			Commands: commands,
			Enabled:  mod.Settings.Enabled,
		})
	}

	return &res, nil
}

func (m *ManageModule) handleModuleGet(c *api.Context) (*ModuleGetResponseWire, error) {
	mod, err := m.moduleManager.GuildModule(c.Context(), c.Guild.ID, c.Param("moduleId"))
	if err != nil {
		return nil, err
	}

	var commands []discord.ApplicationCommandCreate
	if mod, ok := mod.Module.(module.GenericModuleWithCommands); ok {
		commands = mod.Commands()
	}

	commandOverwrites := make(map[string]ModuleCommandOverwriteWire, len(mod.Settings.CommandOverwrites))
	for commandName, overwrite := range mod.Settings.CommandOverwrites {
		commandOverwrites[commandName] = ModuleCommandOverwriteWire{
			Disabled: overwrite.Disabled,
		}
	}

	metadata := mod.Module.Metadata()
	return &ModuleWithConfigWire{
		ModuleWire: ModuleWire{
			ID:       mod.Module.ModuleID(),
			Metadata: metadata,
			Commands: commands,
			Enabled:  mod.Settings.Enabled,
		},
		CommandOverwrites: commandOverwrites,
		Config:            mod.Settings.Config,
	}, nil
}

func (m *ManageModule) handleModuleConfigure(c *api.Context, req ModuleConfigureRequestWire) (*ModuleConfigureResponseWire, error) {
	mod := m.moduleManager.Module(c.Param("moduleId"))
	if mod == nil {
		return nil, api.ErrNotFound("module_not_found", "Module not found")
	}

	if req.Config == nil {
		req.Config = json.RawMessage("{}")
	}

	err := module.ValidateConfig(req.Config, mod.Metadata().ConfigSchema)
	if err != nil {
		return nil, api.ErrBadRequest("invalid_config", "Invalid config")
	}

	commandOverwrites := make(map[string]model.ModuleCommandOverwrite, len(req.CommandOverwrites))
	for commandName, overwrite := range req.CommandOverwrites {
		commandOverwrites[commandName] = model.ModuleCommandOverwrite{
			Disabled: overwrite.Disabled,
		}
	}

	settings := model.ModuleSettings{
		GuildID:           c.Guild.ID,
		ModuleID:          c.Param("moduleId"),
		Enabled:           req.Enabled,
		CommandOverwrites: commandOverwrites,
		Config:            req.Config,
		UpdatedAt:         time.Now().UTC(),
	}

	err = m.moduleManager.UpdateModuleSettings(c.Context(), settings)
	if err != nil {
		return nil, err
	}

	var commands []discord.ApplicationCommandCreate
	if mod, ok := mod.(module.GenericModuleWithCommands); ok {
		commands = mod.Commands()
	}

	commandOverwritesWire := make(map[string]ModuleCommandOverwriteWire, len(settings.CommandOverwrites))
	for commandName, overwrite := range settings.CommandOverwrites {
		commandOverwritesWire[commandName] = ModuleCommandOverwriteWire{
			Disabled: overwrite.Disabled,
		}
	}

	metadata := mod.Metadata()
	return &ModuleConfigureResponseWire{
		ModuleWire: ModuleWire{
			ID:       mod.ModuleID(),
			Metadata: metadata,
			Commands: commands,
			Enabled:  settings.Enabled,
		},
		CommandOverwrites: commandOverwritesWire,
		Config:            settings.Config,
	}, nil
}
