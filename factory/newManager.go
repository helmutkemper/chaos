package factory

import "github.com/helmutkemper/chaos/internal/manager"

func NewManager(errorCh chan error) (reference *manager.Manager) {
	reference = new(manager.Manager)
	reference.New(errorCh)
	return
}
