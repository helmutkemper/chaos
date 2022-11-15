package docker

import (
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"log"
	"strings"
)

// ContainerBuildWithoutStartingItFromImage
//
// English:
//
//	Transforms an image downloaded by ImagePull() or created by ImageBuildFromFolder() into a container
//
//	 Output:
//	   err: Default object error from golang
//
// Português:
//
//	Transforma uma imagem baixada por ImagePull() ou criada por ImageBuildFromFolder() em container
//
//	 Saída:
//	   err: Objeto padrão de erro golang
func (e *ContainerBuilder) ContainerBuildWithoutStartingItFromImage() (err error) {
	err = e.verifyImageName()
	if err != nil {
		return
	}

	_, err = e.dockerSys.ImageFindIdByName(e.imageName)
	if err != nil {
		return
	}

	var netConfig *network.NetworkingConfig
	if e.network != nil {
		e.IPV4Address, netConfig, err = e.network.GetConfiguration()
		if err != nil {
			return
		}
	}

	var portMap = nat.PortMap{}
	var originalImagePortlist []nat.Port
	var originalImagePortlistAsString string
	originalImagePortlist, err = e.dockerSys.ImageListExposedPortsByName(e.imageName)

	if err != nil {
		return
	}

	for k, v := range originalImagePortlist {
		if k != 0 {
			originalImagePortlistAsString += ", "
		}
		originalImagePortlistAsString += v.Port()
	}

	if e.openAllPorts == true {
		for _, port := range originalImagePortlist {
			portMap[port] = []nat.PortBinding{{HostPort: port.Port()}}
		}
	} else if e.openPorts != nil {
		var port nat.Port
		for _, portToOpen := range e.openPorts {
			//var pass = false
			//for _, portToVerify := range originalImagePortlist {
			//	if portToVerify.Port() == portToOpen {
			//		//pass = true
			//		break
			//	}
			//}

			//comentado - nem sempre funciona verificar - início
			//if pass == false {
			//	err = errors.New("port " + portToOpen + " not found in image port list. port list: " + originalImagePortlistAsString)
			//				//	return
			//}
			//comentado - nem sempre funciona verificar - fim

			port, err = nat.NewPort("tcp", portToOpen)
			if err != nil {
				return
			}

			portMap[port] = []nat.PortBinding{{HostPort: port.Port()}}
		}
	} else if e.changePorts != nil {
		var imagePort nat.Port
		var newPort nat.Port

		for _, newPortLinkMap := range e.changePorts {
			imagePort, err = nat.NewPort("tcp", newPortLinkMap.OldPort)
			if err != nil {
				return
			}

			//var pass = false
			//for _, portToVerify := range originalImagePortlist {
			//	if portToVerify.Port() == newPortLinkMap.OldPort {
			//		pass = true
			//		break
			//	}
			//}
			//
			//if pass == false {
			//	err = errors.New("port " + newPortLinkMap.OldPort + " not found in image port list. port list: " + originalImagePortlistAsString)
			//				//	return
			//}

			newPort, err = nat.NewPort("tcp", newPortLinkMap.NewPort)
			if err != nil {
				return
			}
			portMap[imagePort] = []nat.PortBinding{{HostPort: newPort.Port()}}
		}
	}

	if e.printBuildOutput == true {
		go func(ch *chan ContainerPullStatusSendToChannel) {
			for {

				select {
				case event := <-*ch:
					var stream = event.Stream
					stream = strings.ReplaceAll(stream, "\n", "")
					stream = strings.ReplaceAll(stream, "\r", "")
					stream = strings.Trim(stream, " ")

					if stream == "" {
						continue
					}

					log.Printf("%v", stream)

					//if event.Closed == true {
					//	return
					//}
				}
			}
		}(&e.changePointer)
	}

	e.containerConfig.OpenStdin = true
	e.containerConfig.AttachStderr = true
	e.containerConfig.AttachStdin = true
	e.containerConfig.AttachStdout = true
	e.containerConfig.Env = e.environmentVar
	e.containerConfig.Image = e.imageName

	e.containerID, err = e.dockerSys.ContainerCreateWithConfig(
		&e.containerConfig,
		e.containerName,
		e.restartPolicy,
		portMap,
		e.volumes,
		netConfig,
	)
	if err != nil {
		return
	}

	return
}
