package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// example:padrão

func ExampleContainerBuilder_AddFileOrFolderToLinkBetweenComputerHostAndContainer() {

	var err error

	// English: Mounts an image cache and makes imaging up to 5x faster
	//
	// Português: Monta uma imagem cache e deixa a criação de imagens até 5x mais rápida
	// [optional/opcional]
	err = SaImageMakeCacheWithDefaultName("./example/cache/", 365*24*60*60*time.Second)
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	var container = ContainerBuilder{}

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	// [optional/opcional]
	container.SetPrintBuildOnStrOut()

	// English: Sets a validity time for the image, preventing the same image from being remade for a period of time.
	// In some tests, the same image is created inside a loop, and adding an expiration date causes the same image to be used without having to redo the same image at each loop iteration.
	//
	// Português: Define uma tempo de validade para a imagem, evitando que a mesma imagem seja refeita durante um período de tempo.
	// Em alguns testes, a mesma imagem é criada dentro de um laço, e adicionar uma data de validade faz a mesma imagem ser usada sem a necessidade de refazer a mesma imagem a cada interação do loop
	// [optional/opcional]
	container.SetImageExpirationTime(5 * time.Minute)

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

	// English: See the functions: SetGitCloneToBuildWithUserPassword(), SetGitCloneToBuildWithPrivateSshKey() and SetGitCloneToBuildWithPrivateToken()
	//
	// Português: Veja as funções: SetGitCloneToBuildWithUserPassword(), SetGitCloneToBuildWithPrivateSshKey() e SetGitCloneToBuildWithPrivateToken()

	// English: Set a waits for the text to appear in the standard container output to proceed [optional]
	//
	// Português: Define a espera pelo texto aguardado aparecer na saída padrão do container para prosseguir [opcional]
	container.SetWaitStringWithTimeout(
		"Stating server on port 3000",
		10*time.Second,
	)

	// English: Change and open port 3000 to 3030
	//
	// Português: Mude e abra a porta 3000 para 3030
	container.AddPortToChange(
		"3000",
		"3030",
	)

	// English: Replace container folder /static to host folder ./test/static
	//
	// Português: Substitua a pasta do container /static para a pasta da máquina ./test/static
	err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer(
		"./test/static",
		"/static",
	)
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

	// English: Creates an image from a project server.
	//
	// Português: Cria uma imagem a partir do servidor com o projeto.
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

	// English: container "container_delete_server_after_test" running and ready for use on this code point on port 3030
	//
	// Português: container "container_delete_server_after_test" executando e pronto para uso nesse ponto do código na porta 3030

	// English: Read server inside a container on address http://localhost:3030/
	//
	// Português: Lê o servidor dentro do container em http://localhost:3030/
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

	// English: Print the output from get() function
	//
	// Português: Imprime a saída da função get()
	fmt.Printf("%s", body)

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	//Output:
	//<html><body><p>C is life! Golang is a evolution of C</p></body></html>
}
