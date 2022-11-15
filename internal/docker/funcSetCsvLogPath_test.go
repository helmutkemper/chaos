package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"log"
)

// IDE: Goland
// Test Framework: gotest
// Test Kind: File
// Files: /Users/kemper/go/Libraries/src/github.com/helmutkemper/iotmaker.docker.builder/funcSetLogPath_test.go
// Working Directory: /Users/kemper/go/Libraries/src/github.com/helmutkemper/iotmaker.docker.builder
// Module: iotmaker.docker.builder

func ExampleContainerBuilder_SetCsvLogPath() {
	var err error
	var imageInspect types.ImageInspect

	SaGarbageCollector()

	var container = ContainerBuilder{}
	// imprime a saída padrão do container
	container.SetPrintBuildOnStrOut()
	// caso exista uma imagem de nome cache:latest, ela será usada como base para criar o container
	container.SetCacheEnable(true)
	// monta um dockerfile padrão para o golang onde o arquivo main.go e o arquivo go.mod devem está na pasta raiz
	container.MakeDefaultDockerfileForMe()
	// new image name delete:latest
	container.SetImageName("delete:latest")
	// set a folder path to make a new image
	container.SetBuildFolderPath("./test/counter")
	// container name container_delete_server_after_test
	container.SetContainerName("container_counter_delete_after_test")
	// define o limite de memória
	container.SetImageBuildOptionsMemory(100 * KMegaByte)

	container.SetCsvLogPath("./test.counter.log.csv", true)
	container.SetCsvFileValueSeparator("\t")
	container.AddFilterToCvsLogWithReplace(
		"contador",
		"counter",
		"^.*?counter: (?P<valueToGet>[\\d\\.]+)",
		"\\.",
		",",
	)
	container.AddFilterToSuccess(
		"done!",
		"^.*?(?P<valueToGet>\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ done!).*",
		"(?P<date>\\d+/\\d+/\\d+)\\s+(?P<hour>\\d+:\\d+:\\d+)\\s+(?P<value>done!).*",
		"${value}",
	)
	container.AddFilterToFail(
		"counter: 40",
		"^.*?(?P<valueToGet>\\d+/\\d+/\\d+ \\d+:\\d+:\\d+ counter: [\\d\\.]+).*",
		"(?P<date>\\d+/\\d+/\\d+)\\s+(?P<hour>\\d+:\\d+:\\d+)\\s+counter:\\s+(?P<value>[\\d\\.]+).*",
		"Test Fail! Counter Value: ${value} - Hour: ${hour} - Date: ${date}",
	)

	err = container.Init()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	imageInspect, err = container.ImageBuildFromFolder()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	fmt.Printf("image size: %v\n", container.SizeToString(imageInspect.Size))
	fmt.Printf("image os: %v\n", imageInspect.Os)

	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		return
	}

	container.StartMonitor()

	event := container.GetChaosEvent()

	for {
		e := <-event

		if e.Error || e.Fail {
			fmt.Printf("container name: %v\n", e.ContainerName)
			log.Printf("Error: %v", e.Message)
			return
		}
		if e.Done || e.Error || e.Fail {

			fmt.Printf("container name: %v\n", e.ContainerName)
			fmt.Printf("done: %v\n", e.Done)
			fmt.Printf("fail: %v\n", e.Fail)
			fmt.Printf("error: %v\n", e.Error)
			fmt.Printf("message: %v\n", e.Message)

			break
		}
	}

	container.StopMonitor()

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
