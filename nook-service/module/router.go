package module

import (
	"fmt"
	"log/slog"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
)

type ConfigureHandler[C any] func(c Context[C], e *ConfigureEvent) error

type Router[C any] struct {
	configure        ConfigureHandler[C]
	discordListeners []eventListener[C]
	commands         []commandHandler[C]
	components       []componentHandler[C]
}

func NewRouter[C any]() Router[C] {
	return Router[C]{}
}

func (r Router[C]) Configure(handler ConfigureHandler[C]) Router[C] {
	r.configure = handler
	return r
}

func (r Router[C]) Handle(listener eventListener[C]) Router[C] {
	r.discordListeners = append(r.discordListeners, listener)
	return r
}

func (r Router[C]) Command(name string, fn func(c Context[C], e *events.ApplicationCommandInteractionCreate) error) Router[C] {
	r.commands = append(r.commands, newCommandHandler(name, fn))
	return r
}

func (r Router[C]) Component(name string, fn func(c Context[C], e *events.ComponentInteractionCreate) error) Router[C] {
	r.components = append(r.components, newComponentHandler(name, fn))
	return r
}

func (r Router[C]) OnEvent(c Context[C], e Event) error {
	switch e := e.(type) {
	case *ConfigureEvent:
		if r.configure != nil {
			return r.configure(c, e)
		}
	case *DiscordEvent:
		for _, listener := range r.discordListeners {
			// go func() {
			if err := listener.OnEvent(c, e.Event); err != nil {
				slog.Error(
					"Failed to handle discord event",
					slog.String("module_id", c.ModuleID()),
					slog.String("event_type", fmt.Sprintf("%T", e.Event)),
					slog.Any("err", err),
				)
			}
			// }()
		}

		switch e := e.Event.(type) {
		case *events.ApplicationCommandInteractionCreate:
			for _, command := range r.commands {
				command.OnEvent(c, e)
			}
		case *events.ComponentInteractionCreate:
			for _, component := range r.components {
				component.OnEvent(c, e)
			}
		}
	}
	return nil
}

func (r Router[C]) Generic() GenericRouter {
	return &genericRouterImpl[C]{
		inner: r,
	}
}

type GenericRouter interface {
	OnEvent(c GenericContext, e Event) error
}

type genericRouterImpl[C any] struct {
	inner Router[C]
}

func (r *genericRouterImpl[C]) OnEvent(gc GenericContext, e Event) error {
	c := NewContext[C](gc)
	return r.inner.OnEvent(c, e)
}

type eventListenerImpl[C any, E bot.Event] struct {
	fn func(c Context[C], e E) error
}

func newEventListener[C any, E bot.Event](fn func(c Context[C], e E) error) eventListener[C] {
	return &eventListenerImpl[C, E]{
		fn: fn,
	}
}

func (l *eventListenerImpl[C, E]) OnEvent(c Context[C], e bot.Event) error {
	if e, ok := e.(E); ok {
		return l.fn(c, e)
	}
	return nil
}

type eventListener[C any] interface {
	OnEvent(c Context[C], e bot.Event) error
}

func ListenerFunc[C any, E bot.Event](handler func(c Context[C], e E) error) eventListener[C] {
	return newEventListener(handler)
}

type commandHandler[C any] struct {
	pattern string
	fn      func(c Context[C], e *events.ApplicationCommandInteractionCreate) error
}

func newCommandHandler[C any](pattern string, fn func(c Context[C], e *events.ApplicationCommandInteractionCreate) error) commandHandler[C] {
	return commandHandler[C]{
		pattern: pattern,
		fn:      fn,
	}
}

func (h *commandHandler[C]) OnEvent(c Context[C], e *events.ApplicationCommandInteractionCreate) {
	if e.Data.CommandName() == h.pattern {
		go func() {
			if err := h.fn(c, e); err != nil {
				slog.Error(
					"Failed to handle command event",
					slog.String("module_id", c.ModuleID()),
					slog.String("event_type", fmt.Sprintf("%T", e)),
					slog.Any("err", err),
				)
			}
		}()
	}
}

type componentHandler[C any] struct {
	pattern string
	fn      func(c Context[C], e *events.ComponentInteractionCreate) error
}

func newComponentHandler[C any](pattern string, fn func(c Context[C], e *events.ComponentInteractionCreate) error) componentHandler[C] {
	return componentHandler[C]{
		pattern: pattern,
		fn:      fn,
	}
}

func (h *componentHandler[C]) OnEvent(c Context[C], e *events.ComponentInteractionCreate) {
	if e.Data.CustomID() == h.pattern {
		go func() {
			if err := h.fn(c, e); err != nil {
				slog.Error(
					"Failed to handle component event",
					slog.String("module_id", c.ModuleID()),
					slog.String("event_type", fmt.Sprintf("%T", e)),
					slog.Any("err", err),
				)
			}
		}()
	}
}
