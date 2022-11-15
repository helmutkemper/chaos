package docker

import "runtime"

// Must be first function call
func (el *DockerSystem) Init() (err error) {

	el.ContextCreate()
	return el.ClientCreate()
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}
