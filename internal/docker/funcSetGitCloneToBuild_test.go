package docker

import (
	"time"
)

func ExampleContainerBuilder_SetGitCloneToBuild() {
	var err error

	SaGarbageCollector()

	var container = ContainerBuilder{}
	// new image name delete:latest
	container.SetImageName("delete:latest")
	// container name container_delete_server_after_test
	container.SetContainerName("container_delete_server_after_test")
	// git project to clone https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git
	container.SetGitCloneToBuild("https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git")
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

	// Stop and delete the container
	SaGarbageCollector()

	// Output:
	//
}
