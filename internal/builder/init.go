package builder

import (
	dockerContainer "github.com/docker/docker/api/types/container"
	"runtime"
)

// Must be first function call
func (el *DockerSystem) Init() (err error) {

	el.ContextCreate()
	el.Config = new(dockerContainer.Config)
	el.Config.AttachStderr = true
	el.Config.AttachStdin = true
	el.Config.AttachStdout = true

	return el.ClientCreate()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
