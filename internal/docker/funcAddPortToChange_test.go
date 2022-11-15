package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func ExampleContainerBuilder_AddPortToChange() {

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

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("delete:latest")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_delete_server_after_test")

	// English: Project to be cloned from github
	//
	// Português: Projeto para ser clonado do github
	container.SetGitCloneToBuild("https://github.com/helmutkemper/iotmaker.docker.builder.public.example.git")

	// English: See SetGitCloneToBuildWithUserPassworh(), SetGitCloneToBuildWithPrivateSshKey() and
	// SetGitCloneToBuildWithPrivateToken()
	//
	// Português: SetGitCloneToBuildWithUserPassworh(), SetGitCloneToBuildWithPrivateSshKey() and
	// SetGitCloneToBuildWithPrivateToken()

	// English: set a waits for the text to appear in the standard container output to proceed [optional]
	//
	// Português: Define a espera pelo texto na saída padrão do container para prosseguir [opcional]
	container.SetWaitStringWithTimeout("Stating server on port 3000", 10*time.Second)

	// English: change and open port 3000 to 3030
	//
	// English: troca a porta 3000 pela porta 3030
	container.AddPortToChange("3000", "3030")

	// English: replace container folder /static to host folder ./test/static
	//
	// Português: substitui a pasta do container /static pela pasta do host ./test/static
	err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./test/static", "/static")
	if err != nil {
		log.Printf("container.AddFileOrFolderToLinkBetweenComputerHostAndContainer().error: %v", err.Error())
		util.TraceToLog()
		panic(err)
	}

	// English: inicialize container object
	//
	// Português: inicializa o objeto container
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

	// English: container build from image delete:latest
	//
	// Português: monta o container a partir da imagem delete:latest
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		log.Printf("container.ContainerBuildAndStartFromImage().error: %v", err.Error())
		panic(err)
	}

	// English: container "container_delete_server_after_test" running and ready for use on this code point on port 3030
	//
	// Português: container "container_delete_server_after_test" executando e pronto para uso neste ponto de código na porta 3030

	// English: read server inside a container on address http://localhost:3030/
	//
	// Português: lê o servidor dentro do container na porta http://localhost:3030/
	var resp *http.Response
	resp, err = http.Get("http://localhost:3030/")
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

	fmt.Printf("%s", body)

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	// Output:
	// <html><body><p>C is life! Golang is a evolution of C</p></body></html>
}
