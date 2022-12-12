package factory

import (
	"github.com/helmutkemper/chaos/internal/manager"
	"github.com/helmutkemper/chaos/internal/standalone"
)

func NewManager() (reference *manager.Manager) {
	reference = new(manager.Manager)
	reference.New()
	return
}

func NewPrimordial() (reference *manager.Primordial) {
	standalone.GarbageCollector()

	ref := new(manager.Manager)
	ref.New()
	reference = ref.Primordial()
	return
}
