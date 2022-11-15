package docker

// AddPortToDockerfileExpose
//
// English:
//
//	Add ports to dockerfile expose tag.
//
//	 Input:
//	   value: port in string form (without a colon, ":")
//
// Português:
//
//	Adiciona portas a tag expose do dockerfile.
//
//	 Entrada:
//	   value: porta na forma de string (sem dois pontos, ":")
func (e *ContainerBuilder) AddPortToDockerfileExpose(value string) {
	if e.exposePortsOnDockerfile == nil {
		e.exposePortsOnDockerfile = make([]string, 0)
	}

	e.exposePortsOnDockerfile = append(e.exposePortsOnDockerfile, value)
}
