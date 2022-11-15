package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
	"time"
)

// example:padrão

func ExampleContainerBuilder_AddFilterToCvsLog() {
	var err error
	var imageInspect types.ImageInspect

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

	// English: Mount a default dockerfile for golang where the `main.go` file and the `go.mod` file should be in the root folder
	//
	// Português: Monta um dockerfile padrão para o golang onde o arquivo `main.go` e o arquivo `go.mod` devem está na pasta raiz
	container.MakeDefaultDockerfileForMeWithInstallExtras()

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("delete:latest")

	// English: Defines the path where the golang code to be transformed into a docker image is located.
	//
	// Português: Define o caminho onde está o código golang a ser transformado em imagem docker.
	container.SetBuildFolderPath("./test/counter")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_counter_delete_after_test")

	// English: Defines the maximum amount of memory to be used by the docker container.
	//
	// Português: Define a quantidade máxima de memória a ser usada pelo container docker.
	container.SetImageBuildOptionsMemory(100 * KMegaByte)

	// English: Prints the name of the constant used for each column of the container on the first line
	// of the report.
	// This function must be used in conjunction with the container.SetCsvFileRowsToPrint() function.
	//
	// Português: Imprime na primeira linha do relatório o nome da constante usada para cada coluna
	// do container.
	// Esta função deve ser usada em conjunto com a função container.SetCsvFileRowsToPrint()
	container.SetCsvFileReader(true)

	// English: Defines the log file path with container statistical data
	//
	// Português: Define o caminho do arquivo de log com dados estatísticos do container
	container.SetCsvLogPath("./test.counter.log.csv", true)

	// English: Adds a search filter to the standard output of the container, to save the information in the log file
	//
	// Português: Adiciona um filtro de busca na saída padrão do container, para salvar a informação no arquivo de log
	container.AddFilterToCvsLogWithReplace(
		// English: Defines the column name
		//
		// Português: Define o nome da coluna
		"contador",

		// English: Defines the text to be searched
		//
		// Português: Define o texto a ser procurado
		"counter",

		// English: Defines the regular expression to be applied on the found text
		//
		// Português: Define a expressão regular a ser aplicada no texto encontrado
		"^.*?counter: (?P<valueToGet>[\\d\\.]+)",

		// English: Defines the text to be replaced on the found text
		//
		// Português: Define o texto a ser substituído no texto encontrado
		"\\.",

		// English: Defines the text to be written on replaced text
		//
		// Português: Define o texto a ser escrito no texto substituído
		":",
	)

	// English: Adds a filter to look for a value in the container's standard output indicating the success of the test.
	//
	// Português: Adiciona um filtro para procurar um valor na saída padrão do container indicando o sucesso do teste.
	container.AddFilterToSuccess(
		// English: Simple text searched in the container's standard output to activate the filter
		//
		// Português: Texto simples procurado na saída padrão do container para ativar o filtro
		"done!",

		// English: Regular expression used to filter what goes into the log using the `valueToGet` parameter.
		//
		// Português: Expressão regular usada para filtrar o que vai para o log usando o parâmetro `valueToGet`.
		"^.*?(?P<valueToGet>\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ done!).*",

		// English: Regular expression used for search and replacement in the text found in the previous step [optional].
		//          Use this feature to clear the message printed on test success event.
		//
		// Português: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
		//            Use esta funcionalidade para limpar a mensagem impressa no evento de sucesso do teste.
		"(?P<date>\\d+/\\d+/\\d+)\\s+(?P<hour>\\d+:\\d+:\\d+)\\s+(?P<value>done!).*",
		"${value}",
	)

	// English: Adds a filter to look for a value in the container's standard output indicating the fail of the test.
	//
	// Português: Adiciona um filtro para procurar um valor na saída padrão do container indicando a falha do teste.
	container.AddFilterToFail(
		// English: Simple text searched in the container's standard output to activate the filter
		//
		// Português: Texto simples procurado na saída padrão do container para ativar o filtro
		"counter: 40000000",

		// English: Regular expression used to filter what goes into the log using the `valueToGet` parameter.
		//
		// Português: Expressão regular usada para filtrar o que vai para o log usando o parâmetro `valueToGet`.
		"^.*?(?P<valueToGet>\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ counter: [\\d\\.]+).*",

		// English: Regular expression used for search and replacement in the text found in the previous step [optional].
		//
		// Português: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
		"(?P<date>\\d+/\\d+/\\d+)\\s+(?P<hour>\\d+:\\d+:\\d+)\\s+counter:\\s+(?P<value>[\\d\\.]+).*",
		"Test Fail! Counter Value: ${value} - Hour: ${hour} - Date: ${date}",
	)

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

	// English: Starts container monitoring at two second intervals. This functionality generates the log and monitors the standard output of the container.
	//
	// Português: Inicializa o monitoramento do container com intervalos de dois segundos. Esta funcionalidade gera o log e monitora a saída padrão do container.
	container.StartMonitor()

	// English: Gets the event channel pointer inside the container.
	//
	// Português: Pega o ponteiro do canal de eventos dentro do container.
	event := container.GetChaosEvent()

	for {
		var pass = false
		select {
		case e := <-event:
			if e.Done == true || e.Error == true || e.Fail == true {
				pass = true

				fmt.Printf("container name: %v\n", e.ContainerName)
				fmt.Printf("done: %v\n", e.Done)
				fmt.Printf("fail: %v\n", e.Fail)
				fmt.Printf("error: %v\n", e.Error)
				fmt.Printf("message: %v\n", e.Message)

				break
			}
		}

		if pass == true {
			break
		}
	}

	// English: Stop container monitoring.
	//
	// Português: Para o monitoramento do container.
	err = container.StopMonitor()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	SaGarbageCollector()

	// Output:
	// image size: 1.4 MB
	// image os: linux
	// container name: container_counter_delete_after_test
	// done: true
	// fail: false
	// error: false
	// message: done!
}
