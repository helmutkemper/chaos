package docker

import (
	"github.com/docker/docker/api/types/container"
)

// SetImageBuildOptionsIsolationProcess
//
// English:
//
//	Set process isolation mode
//
// Português:
//
//	Determina o método de isolamento do processo
func (e *ContainerBuilder) SetImageBuildOptionsIsolationProcess() {
	e.buildOptions.Isolation = container.IsolationProcess
}
