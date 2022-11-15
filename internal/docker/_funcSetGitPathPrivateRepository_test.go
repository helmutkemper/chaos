package docker

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// this test only work on my acount (sorry)
func ExampleContainerBuilder_SetGitPathPrivateRepository() {
	var err error

	SaGarbageCollector()

	var container = ContainerBuilder{}
	container.SetPrintBuildOnStrOut()
	container.SetGitPathPrivateRepository("github.com/helmutkemper")
	// new image name delete:latest
	container.SetImageName("delete:latest")
	// container name container_delete_server_after_test
	container.SetContainerName("container_delete_server_after_test")
	// git project to clone git@github.com:helmutkemper/iotmaker.docker.builder.private.example.git
	container.SetGitCloneToBuild("git@github.com:helmutkemper/iotmaker.docker.builder.private.example.git")
	container.MakeDefaultDockerfileForMeWithInstallExtras()

	err = container.SetPrivateRepositoryAutoConfig()
	if err != nil {
		panic(err)
	}
	// set a waits for the text to appear in the standard container output to proceed [optional]
	container.SetWaitStringWithTimeout("Stating server on port 3000", 10*time.Second)
	// change and open port 3000 to 3030
	container.AddPortToChange("3000", "3030")
	// replace container folder /static to host folder ./test/static
	err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./test/static", "/static")
	if err != nil {
		panic(err)
	}

	// inicialize container object
	err = container.Init()
	if err != nil {
		panic(err)
	}

	// builder new image from git project
	_, err = container.ImageBuildFromServer()
	if err != nil {
		panic(err)
	}

	// build a new container from image
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		panic(err)
	}

	// At this point, the container is ready for use on port 3030

	// read server inside a container on address http://localhost:3030/
	var resp *http.Response
	resp, err = http.Get("http://localhost:3030/")
	if err != nil {
		log.Printf("http.Get().error: %v", err.Error())
		panic(err)
	}

	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("http.Get().error: %v", err.Error())
		panic(err)
	}

	// print output
	fmt.Printf("%s", body)

	SaGarbageCollector()

	// Output:
	// <html><body><p>C is life! Golang is a evolution of C</p></body></html>
}
