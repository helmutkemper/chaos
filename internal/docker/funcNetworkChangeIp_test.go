package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
	"log"
	"time"
)

func ExampleContainerBuilder_NetworkChangeIp() {
	var err error
	var imageInspect types.ImageInspect

	// English: Deletes all docker elements with the term `delete` in the name.
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	var netDocker = &dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		panic(err)
	}

	// create a network named delete_after_test, subnet 10.0.0.0/16 e gatway 10.0.0.1
	err = netDocker.NetworkCreate("delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		panic(err)
	}

	var container = ContainerBuilder{}

	container.SetNetworkDocker(netDocker)

	// English: print the standard output of the container
	// Português: imprime a saída padrão do container
	container.SetPrintBuildOnStrOut()

	// English: If there is an image named `cache:latest`, it will be used as a base to create the container.
	// Português: Caso exista uma imagem de nome `cache:latest`, ela será usada como base para criar o container.
	container.SetCacheEnable(true)

	// English: Mount a default dockerfile for golang where the `main.go` file and the `go.mod` file should be in the root folder
	// Português: Monta um dockerfile padrão para o golang onde o arquivo `main.go` e o arquivo `go.mod` devem está na pasta raiz
	container.MakeDefaultDockerfileForMe()

	// English: Name of the new image to be created.
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("delete:latest")

	// English: Defines the path where the golang code to be transformed into a docker image is located.
	// Português: Define o caminho onde está o código golang a ser transformado em imagem docker.
	container.SetBuildFolderPath("./test/doNothing")

	// English: Defines the name of the docker container to be created.
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_counter_delete_after_test")

	// English: Defines the maximum amount of memory to be used by the docker container.
	// Português: Define a quantidade máxima de memória a ser usada pelo container docker.
	container.SetImageBuildOptionsMemory(100 * KMegaByte)

	// English: Wait for a textual event on the container's standard output before continuing
	// Português: Espera por um evento textual na saída padrão do container antes de continuar
	container.SetWaitStringWithTimeout("done!", 15*time.Second)

	// English: Initializes the container manager object.
	// Português: Inicializa o objeto gerenciador de container.
	err = container.Init()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Creates an image from a project folder.
	// Português: Cria uma imagem a partir de uma pasta de projeto.
	imageInspect, err = container.ImageBuildFromFolder()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	fmt.Printf("image size: %v\n", container.SizeToString(imageInspect.Size))
	fmt.Printf("image os: %v\n", imageInspect.Os)

	// English: Creates and initializes the container based on the created image.
	// Português: Cria e inicializa o container baseado na imagem criada.
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	var containerInspect ContainerInspect
	containerInspect, err = container.ContainerInspect()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	fmt.Printf("IP: %v\n", containerInspect.Network.Networks["delete_after_test"].IPAddress)

	err = container.ContainerStop()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	err = container.NetworkChangeIp()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	err = container.ContainerStart()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	containerInspect, err = container.ContainerInspect()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	fmt.Printf("IP: %v\n", containerInspect.Network.Networks["delete_after_test"].IPAddress)

	// English: Deletes all docker elements with the term `delete` in the name.
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	// Output:
	// image size: 1.3 MB
	// image os: linux
	// IP: 10.0.0.2
	// IP: 10.0.0.3
}
