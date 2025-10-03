package module

import (
	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/api"
	"github.com/swaggest/jsonschema-go"
)

type GenericModule interface {
	ModuleID() string
	Metadata() ModuleMetadata
}

type GenericModuleWithCommands interface {
	GenericModule
	Commands() []discord.ApplicationCommandCreate
}

type GenericModuleWithEndpoints interface {
	GenericModule
	Endpoints(mx api.HandlerGroup)
}

type Module[C any] interface {
	GenericModule

	Router() Router[C]
}

type ModuleMetadata struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Icon           string `json:"icon"`
	Internal       bool   `json:"internal"`
	DefaultEnabled bool   `json:"default_enabled"`

	DefaultConfig  any               `json:"default_config,omitzero"`
	ConfigSchema   jsonschema.Schema `json:"config_schema,omitzero"`
	ConfigUISchema ConfigUISchema    `json:"config_ui_schema,omitzero"`
}
