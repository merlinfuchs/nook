package manage

import (
	"encoding/json"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/module"
)

type GlobalModuleWire struct {
	ID       string                `json:"id"`
	Metadata module.ModuleMetadata `json:"metadata"`
}

type GlobalModuleListResponseWire = []GlobalModuleWire

type ModuleWire struct {
	ID       string                             `json:"id"`
	Metadata module.ModuleMetadata              `json:"metadata"`
	Commands []discord.ApplicationCommandCreate `json:"commands"`
	Enabled  bool                               `json:"enabled"`
}

type ModuleWithConfigWire struct {
	ModuleWire        `tstype:",extends,required"`
	CommandOverwrites map[string]ModuleCommandOverwriteWire `json:"command_overwrites"`
	Config            json.RawMessage                       `json:"config"`
}

type ModuleGetResponseWire = ModuleWithConfigWire

type ModuleListResponseWire = []ModuleWire

type ModuleConfigureRequestWire struct {
	Enabled           bool                                  `json:"enabled"`
	CommandOverwrites map[string]ModuleCommandOverwriteWire `json:"command_overwrites"`
	Config            json.RawMessage                       `json:"config"`
}

type ModuleConfigureResponseWire = ModuleWithConfigWire

type ModuleCommandOverwriteWire struct {
	Disabled bool `json:"disabled,omitempty"`
}
