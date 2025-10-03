package welcome

import (
	"github.com/merlinfuchs/nook/nook-service/module"
)

const ModuleID = "welcome"

type WelcomeModule struct {
}

func NewWelcomeModule() *WelcomeModule {
	return &WelcomeModule{}
}

func (m *WelcomeModule) ModuleID() string {
	return ModuleID
}

func (m *WelcomeModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Welcome",
		Description:    "Welcome users to the server.",
		Icon:           "hand",
		Internal:       false,
		DefaultEnabled: false,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
		DefaultConfig:  defaultConfig,
	}
}

func (m *WelcomeModule) Router() module.Router[WelcomeConfig] {
	return module.NewRouter[WelcomeConfig]()
}
