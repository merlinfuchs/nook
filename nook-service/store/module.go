package store

import "github.com/merlinfuchs/nook/nook-service/module"

type ModuleStore interface {
	Modules() []module.GenericModule
	Module(moduleID string) module.GenericModule
}
