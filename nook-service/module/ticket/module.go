package ticket

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/disgoorg/disgo/rest"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/manager"
	"github.com/merlinfuchs/nook/nook-service/module"
)

const ModuleID = "ticket"

type TicketModule struct {
	ctx          context.Context
	rest         rest.Rest
	moduleValues *manager.ScopedModuleValueManager
}

func NewTicketModule(
	ctx context.Context,
	rest rest.Rest,
	moduleValues *manager.ScopedModuleValueManager,
) *TicketModule {
	return &TicketModule{
		ctx:          ctx,
		rest:         rest,
		moduleValues: moduleValues,
	}
}

func (m *TicketModule) Metadata() module.ModuleMetadata {
	return module.ModuleMetadata{
		Name:           "Ticket",
		Description:    "Add a ticket system to your server to allow users to submit requests.",
		Icon:           "ticket",
		Internal:       false,
		DefaultEnabled: false,
		ConfigSchema:   configSchema,
		ConfigUISchema: configUISchema,
	}
}

func (m *TicketModule) Configure(guildID common.ID, enabled bool, config json.RawMessage) error {
	var cfg TicketConfig
	err := json.Unmarshal(config, &cfg)
	if err != nil {
		return err
	}

	// TODO: Create or update panel message

	fmt.Println("Configure", guildID, enabled, cfg)
	return nil
}
