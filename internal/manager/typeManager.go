package manager

import (
	"fmt"
	"github.com/helmutkemper/chaos/internal/builder"
)

type Manager struct {
	DockerSys *builder.DockerSystem
}

func (el *Manager) New() (err error) {
	el.DockerSys = new(builder.DockerSystem)
	err = el.DockerSys.Init()
	if err != nil {
		err = fmt.Errorf("chaos.Manager.New().error: %v. Usually this error occurs when docker is not running", err)
		return
	}

	return
}

func (el *Manager) Primordial() (primordial *Primordial) {
	primordial = new(Primordial)
	primordial.Manager = el
	return
}
