package counting

import (
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

var configSchema = module.MustReflectConfigSchema(CountingConfig{})

var configUISchema = module.ConfigUISchema{
	Properties: map[string]module.ConfigUISchema{
		"channels": {
			Items: &module.ConfigUISchema{
				Properties: map[string]module.ConfigUISchema{
					"id": {
						Widget: module.ConfigUIWidgetChannelSelect,
					},
				},
			},
		},
	},
}

// CountingConfig is the settings for the counting module
type CountingConfig struct {
	// The channels in which the counting module is configured
	Channels []CountingChannelConfig `json:"channels" title:"Channels" description:"The channels in which the counting module is configured"`
}

func (s CountingConfig) ChannelSettings(channelID common.ID) CountingChannelConfig {
	for _, channel := range s.Channels {
		if channel.ID == channelID {
			return channel
		}
	}

	return CountingChannelConfig{
		ID:      channelID,
		Enabled: false,
	}
}

// CountingChannelConfig is the settings for a specific channel in the counting module
type CountingChannelConfig struct {
	// The ID of the channel
	ID common.ID `json:"id" title:"Channel" description:"The channel in which to configure the counting game" required:"true"`
	// Whether the counting module is enabled in the channel
	Enabled bool `json:"enabled" title:"Enabled" description:"Whether the counting module is enabled in the channel" default:"true"`
}
