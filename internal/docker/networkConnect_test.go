package docker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/iotmaker.docker/util"
	"os"
	"reflect"
)

func ExampleDockerSystem_NetworkConnect() {

	var err error
	var containerId string
	var imageId string
	var networkId string
	var dockerSys *DockerSystem

	// English: make a channel to end goroutine
	// Português: monta um canal para terminar a goroutine
	var chProcessEnd = make(chan bool, 1)

	// English: make a channel [optional] to print build output
	// Português: monta o canal [opcional] para imprimir a saída do build
	var chStatus = make(chan ContainerPullStatusSendToChannel, 1)

	// English: make a thread to monitoring and print channel data
	// Português: monta uma thread para imprimir os dados do canal
	go func(chStatus chan ContainerPullStatusSendToChannel, chProcessEnd chan bool) {

		for {
			select {
			case <-chProcessEnd:
				return

			case status := <-chStatus:
				// English: remove this comment to see all build status
				// Português: remova este comentário para vê todo o status da criação da imagem
				//fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					// fmt.Println("image pull complete!")

					// English: Eliminate this goroutine after process end
					// Português: Elimina a goroutine após o fim do processo
					// return
				}
			}
		}

	}(chStatus, chProcessEnd)

	// English: searches for the folder containing the test server
	// Português: procura pela pasta contendo o servidor de teste
	var smallServerPath string
	smallServerPath, err = util.FileFindRecursivelyFullPath("small_test_server_port_3000")
	if err != nil {
		panic(err)
	}

	// English: 'static' folder path
	// Português: caminho da pasta 'static'
	var smallServerPathStatic string
	smallServerPathStatic = smallServerPath + string(os.PathSeparator) + "static"

	// English: create a new default client. Please, use: err, dockerSys = factoryDocker.NewClient()
	// Português: cria um novo cliente com configurações padrão. Por favor, usr: err, dockerSys = factoryDocker.NewClient()
	dockerSys = &DockerSystem{}
	dockerSys.ContextCreate()
	err = dockerSys.ClientCreate()
	if err != nil {
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// English: create a network named 'network_delete_before_test'
	// Português: cria uma nova rede de nome 'network_delete_before_test'
	networkId, _, err = dockerSys.NetworkCreate(
		"network_delete_before_test",
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)
	if err != nil {
		panic(err)
	}

	if networkId == "" {
		err = errors.New("network id was not generated")
		panic(err)
	}

	// English: build a new image from folder 'small_test_server_port_3000'
	// Português: monta uma imagem a partir da pasta 'small_test_server_port_3000'
	imageId, err = dockerSys.ImageBuildFromFolder(
		smallServerPath,
		"image_server_delete_before_test:latest",
		[]string{},
		types.ImageBuildOptions{},
		&chStatus, // [channel|nil]
	)
	if err != nil {
		panic(err)
	}

	if imageId == "" {
		err = errors.New("image ID was not generated")
		panic(err)
	}

	// English: building a multi-step image leaves large and useless images, taking up space on the HD.
	// Português: construir uma imagem de múltiplas etapas deixa imagens grandes e sem serventia, ocupando espaço no HD.
	err = dockerSys.ImageGarbageCollector()
	if err != nil {
		panic(err)
	}

	// English: mount and start a container
	// Português: monta i inicializa o container
	containerId, err = dockerSys.ContainerCreate(
		// image name
		"image_server_delete_before_test:latest",
		// container name
		"container_delete_before_test",
		// restart policy
		KRestartPolicyUnlessStopped,
		// portMap
		nat.PortMap{
			// container port number/protocol [tpc/udp]
			"3000/tcp": []nat.PortBinding{ // server original port
				{
					// server output port number
					HostPort: "9002",
				},
			},
		},
		// mount volumes
		[]mount.Mount{
			{
				// bind - is the type for mounting host dir (real folder inside computer where
				// this code work)
				Type: KVolumeMountTypeBindString,
				// path inside host machine
				Source: smallServerPathStatic,
				// path inside image
				Target: "/static",
			},
		},
		// nil or container network configuration
		nil,
	)
	if err != nil {
		panic(err)
	}

	if containerId == "" {
		err = errors.New("container id was not generated")
		panic(err)
	}

	err = dockerSys.NetworkConnect(networkId, containerId, nil)
	if err != nil {
		panic(err)
	}

	// English: container start
	// Português: inicia o container
	err = dockerSys.ContainerStart(containerId)
	if err != nil {
		panic(err)
	}

	// English: ends a goroutine
	// Português: termina a goroutine
	chProcessEnd <- true

	// English: inspect a container by id
	// Português: inspeciona um container por id
	var inspect types.ContainerJSON
	inspect, err = dockerSys.ContainerInspect(containerId)
	if err != nil {
		panic(err)
	}

	if inspect.Name != "/container_delete_before_test" {
		err = errors.New("wrong container name")
		panic(err)
	}

	if inspect.State == nil {
		err = errors.New("container not running")
		panic(err)
	}

	if inspect.State.Running == false {
		err = errors.New("container not running")
		panic(err)
	}

	if len(inspect.Config.ExposedPorts) == 0 {
		err = errors.New("container not running")
		panic(err)
	}

	if reflect.ValueOf(inspect.Config.ExposedPorts["3000/tcp"]).IsZero() == false {
		err = errors.New("exposed ports fail")
		panic(err)
	}

	if inspect.NetworkSettings.Networks["network_delete_before_test"].IPAddress != "10.0.0.2" {
		err = errors.New("IPv4 address error")
		panic(err)
	}

	err = dockerSys.ContainerStop(containerId)
	if err != nil {
		panic(err)
	}

	err = dockerSys.NetworkDisconnect(networkId, containerId, false)
	if err != nil {
		panic(err)
	}

	err = dockerSys.NetworkRemove(networkId)
	if err != nil {
		panic(err)
	}

	networkId, err = dockerSys.NetworkFindIdByName("bridge")
	if err != nil {
		panic(err)
	}

	err = dockerSys.NetworkConnect(networkId, containerId, nil)
	if err != nil {
		panic(err)
	}

	// English: container start
	// Português: inicia o container
	err = dockerSys.ContainerStart(containerId)
	if err != nil {
		panic(err)
	}

	// English: inspect a container by id
	// Português: inspeciona um container por id
	inspect, err = dockerSys.ContainerInspect(containerId)
	if err != nil {
		panic(err)
	}

	if inspect.Name != "/container_delete_before_test" {
		err = errors.New("wrong container name")
		panic(err)
	}

	if inspect.State == nil {
		err = errors.New("container not running")
		panic(err)
	}

	if inspect.State.Running == false {
		err = errors.New("container not running")
		panic(err)
	}

	if reflect.ValueOf(inspect.Config.ExposedPorts["3000/tcp"]).IsZero() == false {
		err = errors.New("exposed ports fail")
		panic(err)
	}

	if inspect.NetworkSettings.Networks["bridge"].IPAddress == "10.0.0.2" {
		err = errors.New("IPv4 address error")
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// Output:
	//
}
