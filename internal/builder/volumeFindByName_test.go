package builder

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/iotmaker.docker/util"
)

func ExampleDockerSystem_VolumeFindByName() {

	var err error
	var containerId string
	var imageId string
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

	// English: build a new image from folder 'small_test_server_port_3000'
	// Português: monta uma imagem a partir da pasta 'small_test_server_port_3000'
	imageId, err = dockerSys.ImageBuildFromFolder(
		smallServerPath,
		"image_server_delete_before_test:latest",
		[]string{},
		types.ImageBuildOptions{},
		chStatus, // [channel|nil]
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
				Type: KVolumeMountTypeVolumeString,
				// path inside host machine
				Source: "static_smallServer",
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

	// English: container start
	// Português: inicia o container
	err = dockerSys.ContainerStart(containerId)
	if err != nil {
		panic(err)
	}

	// English: ends a goroutine
	// Português: termina a goroutine
	chProcessEnd <- true

	var volume types.Volume
	volume, err = dockerSys.VolumeFindByName("static_smallServer")
	if err != nil {
		panic(err)
	}

	if volume.Name != "static_smallServer" {
		err = errors.New("volume name not found")
		panic(err)
	}

	if volume.Driver != "local" {
		err = errors.New("wrong volume drive")
		panic(err)
	}

	if volume.Scope != "local" {
		err = errors.New("wrong volume scope")
		panic(err)
	}

	err = dockerSys.ContainerStopAndRemove(containerId)
	if err != nil {
		panic(err)
	}

	err = dockerSys.VolumeRemove("static_smallServer")
	if err != nil {
		panic(err)
	}

	volume, err = dockerSys.VolumeFindByName("static_smallServer")
	if err != nil && err.Error() != "volume name not found" {
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
