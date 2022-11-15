package docker

import (
	"github.com/docker/docker/api/types/container"
)

// SetImageBuildOptionsIsolationDefault
//
// English:
//
//	Set default isolation mode on current daemon
//
// Português:
//
//	Define o método de isolamento do processo como sendo o mesmo do deamon
func (e *ContainerBuilder) SetImageBuildOptionsIsolationDefault() {
	e.buildOptions.Isolation = container.IsolationDefault
}
