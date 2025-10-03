package nook

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/merlinfuchs/nook/nook-service/module"
)

type EventListener[T any] interface {
	OnEvent(c module.Context[T], e bot.Event)
}
