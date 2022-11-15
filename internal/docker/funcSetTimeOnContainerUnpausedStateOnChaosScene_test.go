package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
	"time"
)

func ExampleContainerBuilder_SetTimeOnContainerUnpausedStateOnChaosScene() {
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
	container.MakeDefaultDockerfileForMe()

	// English: Name of the new image to be created.
	//
	// Português: Nome da nova imagem a ser criada.
	container.SetImageName("delete:latest")

	// English: Defines the path where the golang code to be transformed into a docker image is located.
	//
	// Português: Define o caminho onde está o código golang a ser transformado em imagem docker.
	container.SetBuildFolderPath("./test/chaos")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_counter_delete_after_test")

	// English: Defines the maximum amount of memory to be used by the docker container.
	//
	// Português: Define a quantidade máxima de memória a ser usada pelo container docker.
	container.SetImageBuildOptionsMemory(100 * KMegaByte)

	// English: Defines the log file path with container statistical data
	//
	// Português: Define o caminho do arquivo de log com dados estatísticos do container
	container.SetCsvLogPath("./test.counter.log.csv", true)

	// English: Defines the separator used in the CSV file
	//
	// Português: Define o separador usado no arquivo CSV
	container.SetCsvFileValueSeparator("\t")

	// English: Adds a search filter to the standard output of the container, to save the information in the log file
	//
	// Português: Adiciona um filtro de busca na saída padrão do container, para salvar a informação no arquivo de log
	container.AddFilterToCvsLogWithReplace(
		// English: Label to be written to log file
		//
		// Português: Rótulo a ser escrito no arquivo de log
		"contador",

		// English: Simple text searched in the container's standard output to activate the filter
		//
		// Português: Texto simples procurado na saída padrão do container para ativar o filtro
		"counter",

		// English: Regular expression used to filter what goes into the log using the `valueToGet` parameter.
		//
		// Português: Expressão regular usada para filtrar o que vai para o log usando o parâmetro `valueToGet`.
		"^.*?counter: (?P<valueToGet>[\\d\\.]+)",

		// English: Regular expression used for search and replacement in the text found in the previous step [optional].
		//
		// Português: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
		"\\.",
		",",
	)

	// English: Adds a filter to look for a value in the container's standard output indicating the possibility of restarting the container.
	//
	// Português: Adiciona um filtro para procurar um valor na saída padrão do container indicando a possibilidade de reiniciar o container.
	container.AddFilterToRestartContainer(
		// English: Simple text searched in the container's standard output to activate the filter
		//
		// Português: Texto simples procurado na saída padrão do container para ativar o filtro
		"restart-me!",

		// English: Regular expression used to filter what goes into the log using the `valueToGet` parameter.
		//
		// Português: Expressão regular usada para filtrar o que vai para o log usando o parâmetro `valueToGet`.
		"^.*?(?P<valueToGet>restart-me!)",

		// English: Regular expression used for search and replacement in the text found in the previous step [optional].
		//
		// Português: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
		"",
		"",
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
		//
		// Português: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
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
		"counter: 340",

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

	// English: Adds a filter to look for a value in the container's standard output releasing the chaos test to be started
	//
	// Português: Adiciona um filtro para procurar um valor na saída padrão do container liberando o início do teste de caos
	container.AddFilterToStartChaos(
		"chaos enable",
		"chaos enable",
		"",
		"",
	)

	// English: Defines the probability of the container restarting and changing the IP address in the process.
	//
	// Português: Define a probalidade do container reiniciar e mudar o endereço IP no processo.
	container.SetRestartProbability(0.9, 1.0, 1)

	// English: Defines a time window used to start chaos testing after container initialized
	//
	// Português: Define uma janela de tempo usada para começar o teste de caos depois do container inicializado
	container.SetTimeToStartChaosOnChaosScene(2*time.Second, 5*time.Second)

	// English: Sets a time window used to release container restart after the container has been initialized
	//
	// Português: Define uma janela de tempo usada para liberar o reinício do container depois do container ter sido inicializado
	container.SetTimeBeforeStartChaosInThisContainerOnChaosScene(2*time.Second, 5*time.Second)

	// English: Defines a time window used to pause the container
	//
	// Português: Define uma janela de tempo usada para pausar o container
	container.SetTimeOnContainerPausedStateOnChaosScene(2*time.Second, 5*time.Second)

	// English: Defines a time window used to unpause the container
	//
	// Português: Define uma janela de tempo usada para remover a pausa do container
	container.SetTimeOnContainerUnpausedStateOnChaosScene(2*time.Second, 5*time.Second)

	// English: Sets a time window used to restart the container after stopping
	//
	// Português: Define uma janela de tempo usada para reiniciar o container depois de parado
	container.SetTimeToRestartThisContainerAfterStopEventOnChaosScene(2*time.Second, 5*time.Second)

	// English: Enable chaos test
	//
	// Português: Habilita o teste de caos
	container.EnableChaosScene(true)

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

	// English: Let the example run until a failure happens to terminate the test
	//
	// Português: Deixa o exemplo rodar até que uma falha aconteça para terminar o teste
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
