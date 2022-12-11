package factory

import "github.com/helmutkemper/chaos/internal/manager"

func NewManager() (reference *manager.Manager) {
	reference = new(manager.Manager)
	reference.New()
	return
}

func NewPrimordial() (reference *manager.Primordial) {
	ref := new(manager.Manager)
	ref.New()
	reference = ref.Primordial()
	return
}
