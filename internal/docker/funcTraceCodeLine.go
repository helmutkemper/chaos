package docker

import "runtime"

func (e ContainerBuilder) traceCodeLine() (file string, line int) {
	_, file, line, _ = runtime.Caller(1)
	return
}
