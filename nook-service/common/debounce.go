package common

import (
	"sync"
	"time"

	"github.com/bep/debounce"
)

type DebounceGroup struct {
	sync.Mutex

	interval    time.Duration
	debeouncers map[string]func(func())
}

func NewDebounceGroup(interval time.Duration) *DebounceGroup {
	return &DebounceGroup{
		interval:    interval,
		debeouncers: make(map[string]func(func())),
	}
}

func (d *DebounceGroup) Debounce(key string, fn func()) {
	d.Lock()
	defer d.Unlock()

	debouncer, ok := d.debeouncers[key]
	if !ok {
		debouncer = debounce.New(d.interval)
		d.debeouncers[key] = debouncer
	}

	debouncer(fn)
}
