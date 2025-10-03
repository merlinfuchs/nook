package autorole

import (
	"github.com/merlinfuchs/nook/nook-service/module"
)

const ModuleID = "autorole"

type AutoroleModule struct{}

func NewAutoroleModule() *AutoroleModule {
	return &AutoroleModule{}

}

func (m *AutoroleModule) ModuleID() string {
	return ModuleID
}

func (m *AutoroleModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Autoroles",
		Description:    "Automatically assign roles to users when they join the server and stick them.",
		Icon:           "tags",
		Internal:       false,
		DefaultEnabled: false,
		DefaultConfig:  defaultConfig,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
	}
}

func (m *AutoroleModule) Router() module.Router[AutoroleConfig] {
	return module.NewRouter[AutoroleConfig]().
		Handle(module.ListenerFunc(m.handleMemberJoin)).
		Handle(module.ListenerFunc(m.handleMemberUpdate))
}
