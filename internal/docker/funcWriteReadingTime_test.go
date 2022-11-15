package docker

import (
	"fmt"
	pb "github.com/helmutkemper/iotmaker.docker.problem"
	"log"
	"runtime"
	"testing"
)

func TestContainerBuilder_writeReadingTime(t *testing.T) {
	var err error

	SaGarbageCollector()

	var logFile = "./test.counter.log.1.csv"

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

	container.SetCsvLogPath(logFile, true)
	container.AddFilterToCvsLog(
		"contador",
		"counter",
		"^.*?counter: (?P<valueToGet>[\\d\\.]+)",
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

	container.SetCsvFileRowsToPrint(KLogColumnReadingTime)

	err = container.Init()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		t.Fail()
		return
	}

	_, err = container.ImageBuildFromFolder()
	if err != nil {
		fmt.Printf("error: %v", err.Error())
		SaGarbageCollector()
		t.Fail()
		return
	}

	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		t.Fail()
		return
	}

	container.StartMonitor()

	event := container.GetChaosEvent()

	select {
	case e := <-event:
		fmt.Printf("container name: %v\n", e.ContainerName)
		fmt.Printf("done: %v\n", e.Done)
		fmt.Printf("fail: %v\n", e.Fail)
		fmt.Printf("error: %v\n", e.Error)
		fmt.Printf("message: %v\n", e.Message)
	}

	err = container.StopMonitor()
	if err != nil {
		log.Printf("error: %v", err.Error())
		SaGarbageCollector()
		t.Fail()
		return
	}

	var problem pb.Problem
	var logTest = TestContainerLog{}
	var listUnderTest []parserLog
	var osName = runtime.GOOS
	switch osName {
	case "darwin":
		listUnderTest = []parserLog{
			{
				KLogReadingTimeLabel,
				KLogReadingTimeValue,
				KLogReadingTimeRegexp,
			},
		}
	}

	problem = logTest.makeTest(logFile, &listUnderTest, t)
	if problem != nil {
		var file, funcName string
		var line int
		file, line, funcName, _ = problem.Trace()
		log.Printf("Error: %v", problem.Error())
		log.Printf("Cause: %v", problem.Cause())
		log.Printf("File: %v", file)
		log.Printf("Function: [%v]: %v", line, funcName)
		return
	}

	//_ = os.Remove(logFile)

	SaGarbageCollector()
}
