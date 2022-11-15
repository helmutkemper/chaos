package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func ExampleContainerBuilder_AddPortToExpose() {
	var err error

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	var container = ContainerBuilder{}

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	container.SetPrintBuildOnStrOut()

	// English: If there is an image named `cache:latest`, it will be used as a base to create the container.
	//
	// Português: Caso exista uma imagem de nome `cache:latest`, ela será usada como base para criar o container.
	container.SetCacheEnable(true)

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("delete:latest")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_delete_server_after_test")

	// English: git project to clone https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git
	//
	// Português: repositório git a ser clonado https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git
	container.SetGitCloneToBuild("https://github.com/helmutkemper/iotmaker.docker.builder.public.example.git")

	// see SetGitCloneToBuildWithUserPassworh(), SetGitCloneToBuildWithPrivateSshKey() and
	// SetGitCloneToBuildWithPrivateToken()

	// English: Set a waits for the text to appear in the standard container output to proceed [optional]
	//
	// Português: Define a espera pelo texto aguardado aparecer na saída padrão do container para prosseguir [opcional]
	container.SetWaitStringWithTimeout("Stating server on port 3000", 20*time.Second)

	// English: open port 3000
	//
	// Português: abre porta 3000
	container.AddPortToExpose("3000")

	// English: Replace container folder /static to host folder ./test/static
	//
	// Português: Substitua a pasta do container /static para a pasta da máquina ./test/static
	err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./test/static", "/static")
	if err != nil {
		log.Printf("container.AddFileOrFolderToLinkBetweenComputerHostAndContainer().error: %v", err.Error())
		util.TraceToLog()
		panic(err)
	}

	// English: Initializes the container manager object.
	//
	// Português: Inicializa o objeto gerenciador de container.
	err = container.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	// English: builder new image from git project
	//
	// Português: monta a nova imagem a partir do projeto git
	_, err = container.ImageBuildFromServer()
	if err != nil {
		util.TraceToLog()
		log.Printf("container.ImageBuildFromServer().error: %v", err.Error())
		panic(err)
	}

	// English: Creates and initializes the container based on the created image.
	//
	// Português: Cria e inicializa o container baseado na imagem criada.
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		log.Printf("container.ContainerBuildAndStartFromImage().error: %v", err.Error())
		panic(err)
	}

	// English: container "container_delete_server_after_test" running and ready for use on this code point on port 3000
	//
	// Português: container "container_delete_server_after_test" executando e pronto para uso neste ponto de código na porta 3000

	// English: read server inside a container on address http://localhost:3030/
	//
	// Português: lê o servidor dentro do container na porta http://localhost:3030/
	var resp *http.Response
	resp, err = http.Get("http://localhost:3000/")
	if err != nil {
		util.TraceToLog()
		log.Printf("http.Get().error: %v", err.Error())
		panic(err)
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		util.TraceToLog()
		log.Printf("http.Get().error: %v", err.Error())
		panic(err)
	}

	// print output
	fmt.Printf("%s", body)

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	// Output:
	// <html><body><p>C is life! Golang is a evolution of C</p></body></html>
}
