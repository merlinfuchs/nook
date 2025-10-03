package logging

import (
	"context"
	"sync"
	"time"

	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/merlinfuchs/nook/nook-service/common"
	"github.com/merlinfuchs/nook/nook-service/module"
)

type messageCacheKey struct {
	channelID common.ID
	messageID common.ID
}

type auditLogCollector struct {
	sync.RWMutex
	waiters map[common.ID]chan *events.GuildAuditLogEntryCreate
}

func newAuditLogCollector() *auditLogCollector {
	return &auditLogCollector{
		waiters: make(map[common.ID]chan *events.GuildAuditLogEntryCreate),
	}
}

func (a *auditLogCollector) OnEvent(c module.Context[LoggingConfig], e bot.Event) error {
	if e, ok := e.(*events.GuildAuditLogEntryCreate); ok {
		time.Sleep(250 * time.Millisecond)

		a.RLock()
		for _, waiter := range a.waiters {
			waiter <- e
		}
		a.RUnlock()
	}
	return nil
}

func (a *auditLogCollector) WaitForAuditLogEntry(
	ctx context.Context,
	filter func(e *events.GuildAuditLogEntryCreate) bool,
) (*discord.AuditLogEntry, error) {
	a.Lock()

	waiterID := common.UniqueID()
	ch := make(chan *events.GuildAuditLogEntryCreate)
	a.waiters[waiterID] = ch

	a.Unlock()

	defer func() {
		a.Lock()
		delete(a.waiters, waiterID)
		a.Unlock()
	}()

	for {
		select {
		case e := <-ch:
			if filter(e) {
				return &e.AuditLogEntry, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(1 * time.Second):
			return nil, nil
		}
	}
}
