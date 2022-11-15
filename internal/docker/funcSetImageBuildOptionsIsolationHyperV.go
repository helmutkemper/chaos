package docker

import (
	"github.com/docker/docker/api/types/container"
)

// SetImageBuildOptionsIsolationHyperV
//
// English:
//
//	Set HyperV isolation mode
//
// Português:
//
//	Define o método de isolamento como sendo HyperV
func (e *ContainerBuilder) SetImageBuildOptionsIsolationHyperV() {
	e.buildOptions.Isolation = container.IsolationHyperV
}
