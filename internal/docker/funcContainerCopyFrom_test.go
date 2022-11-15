package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
	"os"
)

func ExampleContainerBuilder_ContainerCopyFrom() {
	var err error

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	err = buildGoLintImageCopyFromExample()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	err = buildAlpineImageCopyFromExample()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	_ = os.Remove("./example/golint/golangci-lint")

	// Output:
	// image size: 510.9 MB
	// image os: linux
	// file name: golangci-lint
	// image size: 281.8 MB
	// image os: linux
}

func buildGoLintImageCopyFromExample() (err error) {
	var imageInspect types.ImageInspect
	var container = ContainerBuilder{}

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	// [optional/opcional]
	container.SetPrintBuildOnStrOut()

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("golint_delete:latest")

	// English: Golang project path to be turned into docker image
	//
	// Português: Caminho do projeto em Golang a ser transformado em imagem docker
	container.SetBuildFolderPath("./example/golint/imageGolintBuild")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_golint_delete_after_test")

	// English: Initializes the container manager object.
	//
	// Português: Inicializa o objeto gerenciador de container.
	err = container.Init()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Creates an image from a project folder.
	//
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
	//
	// Português: Cria e inicializa o container baseado na imagem criada.
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	var copyResponse []types.ContainerPathStat
	copyResponse, err = container.ContainerCopyFrom(
		[]string{"/go/pkg/mod/github.com/golangci/golangci-lint@v1.23.6/bin/golangci-lint"},
		[]string{"./example/golint/golangci-lint"},
	)
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	fmt.Printf("file name: %v\n", copyResponse[0].Name)

	return
}

func buildAlpineImageCopyFromExample() (err error) {
	var imageInspect types.ImageInspect
	var container = ContainerBuilder{}

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	// [optional/opcional]
	container.SetPrintBuildOnStrOut()

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("alpine_delete:latest")

	// English: Golang project path to be turned into docker image
	//
	// Português: Caminho do projeto em Golang a ser transformado em imagem docker
	container.SetBuildFolderPath("./example/golint/imageAlpineBuild")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_alpine_delete_after_test")

	// English: Initializes the container manager object.
	//
	// Português: Inicializa o objeto gerenciador de container.
	err = container.Init()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Creates an image from a project folder.
	//
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
	//
	// Português: Cria e inicializa o container baseado na imagem criada.
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	err = container.ContainerCopyTo(
		[]string{"./example/golint/golangci-lint"},
		[]string{"/go"},
	)

	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	var exitCode int
	var running bool
	var stdOutput []byte
	var stdError []byte
	exitCode, running, stdOutput, stdError, err = container.ContainerExecCommand([]string{"ls", "-l"})

	log.Printf("exitCode: %v", exitCode)
	log.Printf("running: %v", running)
	log.Printf("stdOutput: %v", string(stdOutput))
	log.Printf("stdError: %v", string(stdError))

	exitCode, running, stdOutput, stdError, err = container.ContainerExecCommand([]string{"./golangci-lint"})

	log.Printf("exitCode: %v", exitCode)
	log.Printf("running: %v", running)
	log.Printf("stdOutput: %v", string(stdOutput))
	log.Printf("stdError: %v", string(stdError))

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	return
}
