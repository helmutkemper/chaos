package builder

// NewClient (English): Prepare docker system
//
//	Example:
//
//	  err, dockerSys := factoryDocker.NewClient()
//	  if err != nil {
//	    panic(err)
//	  }
//	  dockerSys.ContainerCreateChangeExposedPortAndStart(...)
//
// NewClient (Português): Prepara o docker
//
//	Exemplo:
//
//	  err, dockerSys := factoryDocker.NewClient()
//	  if err != nil {
//	    panic(err)
//	  }
//	  dockerSys.ContainerCreateChangeExposedPortAndStart(...)
func NewClient() (
	dockerSystem *DockerSystem,
	err error,
) {
	dockerSystem = &DockerSystem{}
	dockerSystem.ContextCreate()
	err = dockerSystem.ClientCreate()

	return
}
