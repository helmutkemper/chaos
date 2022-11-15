package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"io/ioutil"
	"log"
	"os"
	"time"
)

// example:padrão

func ExampleContainerBuilder_AddFailMatchFlagToFileLog() {
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
	// [optional/opcional]
	SaGarbageCollector()

	var container = ContainerBuilder{}

	// English: print the standard output of the container
	//
	// Português: imprime a saída padrão do container
	// [optional/opcional]
	container.SetPrintBuildOnStrOut()

	// English: If there is an image named `cache:latest`, it will be used as a base to create the container.
	//
	// Português: Caso exista uma imagem de nome `cache:latest`, ela será usada como base para criar o container.
	// [optional/opcional]
	container.SetCacheEnable(true)

	// English: Mount a default dockerfile for golang where the `main.go` file and the `go.mod` file should be in the root folder
	//
	// Português: Monta um dockerfile padrão para o golang onde o arquivo `main.go` e o arquivo `go.mod` devem está na pasta raiz
	// [optional/opcional]
	container.MakeDefaultDockerfileForMe()

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

	// English: Golang project path to be turned into docker image
	//
	// Português: Caminho do projeto em Golang a ser transformado em imagem docker
	container.SetBuildFolderPath("./test/bug")

	// English: Defines the name of the docker container to be created.
	//
	// Português: Define o nome do container docker a ser criado.
	container.SetContainerName("container_counter_delete_after_test")

	// English: Defines the maximum amount of memory to be used by the docker container.
	//
	// Português: Define a quantidade máxima de memória a ser usada pelo container docker.
	// [optional/opcional]
	container.SetImageBuildOptionsMemory(100 * KMegaByte)

	// English: Defines a log, in the form of a CSV file, of the container's performance, with indicators of memory consumption and access times. Note: The log format varies by platform, macos, windows, linux.
	//
	// Português: Define um log, na forma de arquivo CSV, de desempenho do container, com indicadores de consumo de memória e tempos de acesso. Nota: O formato do log varia de acordo com a plataforma, macos, windows, linux.
	// [optional/opcional]
	container.SetCsvLogPath("./test.counter.log.csv", true)

	// English: Swaps the comma by tab, making the file compatible with floating-point numbers
	//
	// Português: Troca a virgula por tabulação, compatibilizando o arquivo com números de ponto flutuante
	container.SetCsvFileValueSeparator("\t")

	// English: Prints in the header of the file the name of the constant responsible for printing the column in the log.
	//
	// Português: Imprime no cabeçalho do arquivo o nome da constante responsável por imprimir a coluna no log.
	// [optional/opcional]
	container.SetCsvFileReader(true)

	// English: Defines which columns to print in the log. To see all columns, set SetCsvFileRowsToPrint(KLogColumnAll) and SetCsvFileReader(true).
	// Open the log file, define the columns to be printed in the log, and then use SetCsvFileRowsToPrint(KReadingTime | KCurrentNumberOfOidsInTheCGroup | KLimitOnTheNumberOfPidsInTheCGroup | ...)
	//
	// Português: Define quais colunas imprimir no log. Para vê todas as colunas, defina SetCsvFileRowsToPrint(KLogColumnAll) e SetCsvFileReader(true).
	// Abra o arquivo de log, defina as colunas a serem impressas no log e em seguida, use SetCsvFileRowsToPrint(KReadingTime | KCurrentNumberOfOidsInTheCGroup | ...)
	// [optional/opcional]
	container.SetCsvFileRowsToPrint(KLogColumnAll)

	// English: Sets a text search filter on the container's standard output and writes the text to the log defined by SetCsvLogPath()
	// The container example prints a counter to standard output `log.Printf("counter: %.2f", counter)`. `label` adds the column name; `match` searches for text; `filter` applies a regular expression; `search` and `replace` do a replacement on top of the found value before writing to the log.
	//
	// Português: Define um filtro de busca por texto na saída padrão do container e escreve o texto no log definido por SetCsvLogPath()
	// O container de exemplo imprime um contador na saída padrão `log.Printf("counter: %.2f", counter)`. `label` adiciona o nome da coluna; `match` procura pelo texto; `filter` aplica uma expressão regular; `search` e `replace` fazem uma substuição em cima do valor encontrado antes de escrever no log.
	// [optional/opcional]
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
		",",
	)

	// English: Adds a failure indicator to the project. Failure indicator is a text searched for in the container's standard output and indicates something that should not have happened during the test.
	//
	// Português: Adiciona um indicador de falha ao projeto. Indicador de falha é um texto procurado na saída padrão do container e indica algo que não deveria ter acontecido durante o teste.
	// [optional/opcional]
	container.AddFailMatchFlag(
		"counter: 40",
	)

	// English: Adds a log file write failure indicator to the project. Failure indicator is a text searched for in the container's standard output and indicates something that should not have happened during the test.
	// Some critical failures can be monitored and when they happen, the container's standard output is filed in a `log.N.log` file, where N is an automatically incremented number.
	//
	// Português: Adiciona um indicador de falha com gravação de arquivo em log ao projeto. Indicador de falha é um texto procurado na saída padrão do container e indica algo que não deveria ter acontecido durante o teste.
	// Algumas falhas críticas podem ser monitoradas e quando elas acontecem, a saída padrão do container é arquivada em um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
	// [optional/opcional]
	err = container.AddFailMatchFlagToFileLog(
		"bug:",
		"./log1/log2/log3",
	)
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

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

	// English: Starts container monitoring at two second intervals. This functionality monitors the container's standard output and generates the log defined by the SetCsvLogPath() function.
	//
	// Português: Inicializa o monitoramento do container com intervalos de dois segundos. Esta funcionalidade monitora a saída padrão do container e gera o log definido pela função SetCsvLogPath().
	// StartMonitor() é usado durante o teste de caos e na geração do log de desempenho do container.
	// [optional/opcional]
	container.StartMonitor()

	// English: Gets the event channel inside the container.
	//
	// Português: Pega o canal de eventos dentro do container.
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

				break
			}
		}

		if pass == true {
			break
		}
	}

	// English: For container monitoring. Note: This function should be used to avoid trying to read a container that no longer exists, erased by the SaGarbageCollector() function.
	//
	// Português: Para o monitoramento do container. Nota: Esta função deve ser usada para evitar tentativa de leitura em um container que não existe mais, apagado pela função SaGarbageCollector().
	// [optional/opcional]
	_ = container.StopMonitor()

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	var data []byte
	data, err = ioutil.ReadFile("./log1/log2/log3/log.0.log")
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	if len(data) == 0 {
		fmt.Println("log file error")
	}

	_ = os.Remove("./log1/log2/log3/log.0.log")
	_ = os.Remove("./log1/log2/log3/")
	_ = os.Remove("./log1/log2/")
	_ = os.Remove("./log1/")

	// Output:
	// image size: 1.4 MB
	// image os: linux
	// container name: container_counter_delete_after_test
	// done: false
	// fail: true
	// error: false
}
