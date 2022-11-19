package manager

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/chaos/internal/builder"
	"path/filepath"
	"strconv"
)

type container struct {
	IPV4Address []string

	//port inside container and host computer port
	portsContainer []nat.Port
	portsHost      [][]int64

	volumeContainer []string
	volumeHost      [][]string
}

type ContainerFromImage struct {
	container

	manager *Manager

	imageId       string
	imageName     string
	containerName string
	copies        int
}

func (el *ContainerFromImage) New(manager *Manager) {
	el.manager = manager
}

func (el *ContainerFromImage) Volumes(containerPath string, hostPath ...string) (ref *ContainerFromImage) {
	var err error

	if el.volumeContainer == nil {
		el.volumeContainer = make([]string, 0)
		el.volumeHost = make([][]string, 0)
	}

	var path string
	var absolutePath []string
	for k := range hostPath {
		if hostPath[k] != "" {
			path, err = filepath.Abs(hostPath[k])
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("containerFromImage.Volumes().error: %v", err)
				return el
			}
		} else {
			path = ""
		}

		absolutePath = append(absolutePath, path)
	}

	el.volumeContainer = append(el.volumeContainer, containerPath)
	el.volumeHost = append(el.volumeHost, absolutePath)
	return el
}

// Ports
//
// Defines which port of the container will be exposed to the world
//
//	Input:
//	  containerProtocol: network protocol `tcp` or `utc`
//	  containerPort: port number on the container. eg: 27017 for MongoDB
//	  localPort: port number on the host computer. eg: 27017 for MongoDB
//
//	Notes:
//	  * When `localPort` receives one more value, each container created will receive a different value.
//	    - Imagine creating 3 containers and passing the values 27016 and 27015. The first container created will receive
//	    27016, the second, 27015 and the third will not receive value;
//	    - Imagine creating 3 containers and passing the values 27016, 0 and 27015. The first container created will
//	    receive 27016, the second will not receive value, and the third receive 27015.
func (el *ContainerFromImage) Ports(containerProtocol string, containerPort int64, localPort ...int64) (ref *ContainerFromImage) {
	if el.portsContainer == nil {
		el.portsContainer = make([]nat.Port, 0)
		el.portsHost = make([][]int64, 0)
	}

	port, err := nat.NewPort(containerProtocol, strconv.FormatInt(containerPort, 10))
	if err != nil {
		el.manager.ErrorCh <- fmt.Errorf("containerFromImage.ExposePorts().error: %v", err)
		return
	}

	el.portsContainer = append(el.portsContainer, port)

	el.portsHost = append(el.portsHost, localPort)
	return el
}

func (el *ContainerFromImage) Create(imageName, containerName string, copies int) (ref *ContainerFromImage) {
	var err error

	if copies == 0 {
		return el
	}

	// adjust image name to have version tag
	el.imageName = el.manager.DockerSys[0].AdjustImageName(imageName)
	el.containerName = containerName
	el.copies = copies

	// if the image does not exist, download the image
	if err = el.imagePull(); err != nil {
		el.manager.ErrorCh <- err
		return el
	}

	var ipAddress string
	var netConfig *networkTypes.NetworkingConfig
	el.IPV4Address = make([]string, 0)
	for i := 0; i != copies; i += 1 {

		// index zero is created when the manager object is created, the other indexes are created here, in case there is
		// more than one container to be created
		if i != 0 {
			var dockerSys = new(builder.DockerSystem)
			_ = dockerSys.Init()
			el.manager.DockerSys = append(el.manager.DockerSys, dockerSys)
		}

		// get the next ip address from network
		if el.manager.network != nil {
			ipAddress, netConfig, err = el.manager.network.generator.GetNext()
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("container.network().GetNext().error: %v", err)
				return
			}
			el.IPV4Address = append(el.IPV4Address, ipAddress)
		}

		// map the port container:host[copiesKey]
		var portConfig = nat.PortMap{}
		for kContainer := range el.portsContainer {
			portBind := make([]nat.PortBinding, 0)
			if len(el.portsHost[kContainer]) > i && el.portsHost[kContainer][i] > 0 {
				portBind = append(portBind, nat.PortBinding{HostPort: strconv.FormatInt(el.portsHost[kContainer][i], 10)})
			}

			portConfig[el.portsContainer[kContainer]] = portBind
		}

		var volumes = make([]mount.Mount, 0)
		for k := range el.volumeContainer {
			volume := mount.Mount{}
			if len(el.volumeContainer[k]) > i && el.volumeHost[k][i] != "" {
				volume.Type = builder.KVolumeMountTypeBindString
				volume.Source = el.volumeHost[k][i]
				volume.Target = el.volumeContainer[k]

				volumes = append(volumes, volume)
			}
		}

		// create the container, link container and network, but, don't start the container
		_, err = el.manager.DockerSys[i].ContainerCreate(
			imageName,
			containerName+"_"+strconv.FormatInt(int64(i), 10),
			builder.KRestartPolicyNo,
			portConfig,
			volumes,
			netConfig,
		)
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerCreate().error: %v", i, err)
			return
		}
	}

	return el
}

// imagePull
//
// If the image exists on the local computer, it does nothing, otherwise it tries to download the image
func (el *ContainerFromImage) imagePull() (err error) {
	el.imageId, _ = el.manager.DockerSys[0].ImageFindIdByName(el.imageName)
	if el.imageId != "" {
		return
	}

	// English: make a channel to end goroutine
	// Português: monta um canal para terminar a goroutine
	var chProcessEnd = make(chan bool, 1)

	// English: make a channel [optional] to print build output
	// Português: monta o canal [opcional] para imprimir a saída do build
	var chStatus = make(chan builder.ContainerPullStatusSendToChannel, 1)

	// English: make a thread to monitoring and print channel data
	// Português: monta uma thread para imprimir os dados do canal
	go func(chStatus chan builder.ContainerPullStatusSendToChannel, chProcessEnd chan bool) {

		for {
			select {
			case <-chProcessEnd:
				// English: Eliminate this goroutine after process end
				// Português: Elimina a goroutine após o fim do processo
				return

			case status := <-chStatus:
				// English: remove this comment to see all build status
				// Português: remova este comentário para vê _todo o status da criação da imagem
				fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					fmt.Println("image pull complete!")
				}
			}
		}

	}(chStatus, chProcessEnd)

	defer func() {
		// English: ends a goroutine
		// Português: termina a goroutine
		chProcessEnd <- true
	}()

	// docker pull
	el.imageId, el.imageName, err = el.manager.DockerSys[0].ImagePull(el.imageName, &chStatus)
	if err != nil {
		err = fmt.Errorf("containerFromImage.Primordial().imagePull().error: %v", err)
		return
	}

	return
}
