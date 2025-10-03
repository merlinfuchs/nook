package module

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/merlinfuchs/nook/nook-service/common"
	jsonschemavalidate "github.com/santhosh-tekuri/jsonschema/v6"
	"github.com/swaggest/jsonschema-go"
)

type ConfigSchema = jsonschema.Schema

type ConfigUISchema struct {
	Widget        ConfigUIWidget        `json:"ui:widget,omitzero"`
	ChannelTypes  []discord.ChannelType `json:"ui:channel_types,omitzero"`
	SelectValues  []ConfigUISelectValue `json:"ui:select_values,omitzero"`
	AllowMultiple bool                  `json:"ui:allow_multiple,omitzero"`

	Properties map[string]ConfigUISchema `json:"properties,omitzero"`
	Items      *ConfigUISchema           `json:"items,omitzero"`
	Layout     ConfigLayout              `json:"layout,omitzero"`
}

type ConfigUIWidget string

const (
	ConfigUIWidgetSelect          ConfigUIWidget = "select"
	ConfigUIWidgetChannelSelect   ConfigUIWidget = "channel_select"
	ConfigUIWidgetRoleSelect      ConfigUIWidget = "role_select"
	ConfigUIWidgetMessage         ConfigUIWidget = "message"
	ConfigUIWidgetMessageFormat   ConfigUIWidget = "message_format"
	ConfigUIWidgetDurationSeconds ConfigUIWidget = "duration_seconds"
	ConfigUIWidgetDurationMinutes ConfigUIWidget = "duration_minutes"
	ConfigUIWidgetDurationHours   ConfigUIWidget = "duration_hours"
	ConfigUIWidgetDurationDays    ConfigUIWidget = "duration_days"
)

type ConfigUISelectValue struct {
	Value string `json:"value"`
	Label string `json:"label"`
}

type ConfigLayoutType string

const (
	ConfigLayoutTypeContainer ConfigLayoutType = "container"
	ConfigLayoutTypeFlexWrap  ConfigLayoutType = "flex_wrap"
	ConfigLayoutTypeRow       ConfigLayoutType = "row"
)

var ConfigLayoutTypeValues = []ConfigLayoutType{
	ConfigLayoutTypeRow,
}

type ConfigLayout struct {
	Type     ConfigLayoutType `json:"type"`
	Header   string           `json:"header,omitempty"`
	Items    []string         `json:"items"`
	Children []ConfigLayout   `json:"children"`
}

func ReflectConfigSchema(config interface{}) (ConfigSchema, error) {
	reflector := jsonschema.Reflector{}

	schema, err := reflector.Reflect(config, jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
		if !params.Processed {
			return nil
		}

		switch params.Field.Type {
		case reflect.TypeOf(common.ID(0)):
			params.PropertySchema.Type = &jsonschema.Type{
				SimpleTypes: ptr(jsonschema.String),
			}
			params.PropertySchema.Pattern = ptr("^[0-9]+$")
		case reflect.TypeOf([]common.ID{}):
			params.PropertySchema.Type = &jsonschema.Type{
				SimpleTypes: ptr(jsonschema.Array),
			}
			params.PropertySchema.Items = &jsonschema.Items{
				SchemaOrBool: &jsonschema.SchemaOrBool{
					TypeObject: &jsonschema.Schema{
						Type: &jsonschema.Type{
							SimpleTypes: ptr(jsonschema.String),
						},
					},
				},
			}
		case reflect.TypeOf(json.RawMessage{}):
			params.PropertySchema.Type = &jsonschema.Type{
				SimpleTypes: ptr(jsonschema.Object),
			}
		case reflect.TypeOf(time.Duration(0)):
			params.PropertySchema.Type = &jsonschema.Type{
				SimpleTypes: ptr(jsonschema.Integer),
			}
		}

		if v, ok := params.Field.Tag.Lookup("type"); ok {
			params.PropertySchema.Type = &jsonschema.Type{
				SimpleTypes: ptr(jsonschema.SimpleType(v)),
			}
		}

		return nil
	}))
	if err != nil {
		return ConfigSchema{}, err
	}

	return schema, nil
}

func MustReflectConfigSchema(config interface{}) ConfigSchema {
	schema, err := ReflectConfigSchema(config)
	if err != nil {
		log.Fatal("Failed to reflect config schema", "err", err)
	}

	return schema
}

func ValidateConfig(config json.RawMessage, schema ConfigSchema) error {
	rawSchema, err := json.Marshal(schema)
	if err != nil {
		return fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	inst, err := jsonschemavalidate.UnmarshalJSON(bytes.NewReader(config))
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	schamResource, err := jsonschemavalidate.UnmarshalJSON(bytes.NewReader(rawSchema))
	if err != nil {
		return fmt.Errorf("failed to unmarshal schema: %w", err)
	}

	c := jsonschemavalidate.NewCompiler()
	if err := c.AddResource("./schema", schamResource); err != nil {
		return fmt.Errorf("failed to add resource: %w", err)
	}

	sch, err := c.Compile("./schema")
	if err != nil {
		return fmt.Errorf("failed to compile schema: %w", err)
	}

	err = sch.Validate(inst)
	if err != nil {
		return fmt.Errorf("failed to validate config: %w", err)
	}

	return nil
}

func ptr[T any](s T) *T {
	return &s
}
