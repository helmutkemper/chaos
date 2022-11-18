package builder

import (
	"github.com/docker/docker/client"
)

// ClientCreate (English): Create a new client from docker
//
//	    Example:
//	      dockerSys = &iotmakerDocker.DockerSystem{}
//		     dockerSys.ContextCreate()
//		     err := dockerSys.ClientCreate()
//	      if err != nil {
//	        panic(err)
//	      }
//	      dockerSys.ContainerCreateChangeExposedPortAndStart(...)
//
// ClientCreate (Português): Cria um novo cliente para o docker
//
//	    Exemplo:
//	      dockerSys = &iotmakerDocker.DockerSystem{}
//		     dockerSys.ContextCreate()
//		     err := dockerSys.ClientCreate()
//	      if err != nil {
//	        panic(err)
//	      }
//	      dockerSys.ContainerCreateChangeExposedPortAndStart(...)
func (el *DockerSystem) ClientCreate() (
	err error,
) {

	el.cli, err = client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)

	return err
}
