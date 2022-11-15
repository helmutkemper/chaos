package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func ImageBuildViewer(ch *chan ContainerPullStatusSendToChannel) {
	go func(ch *chan ContainerPullStatusSendToChannel) {
		for {

			select {
			case event := <-*ch:
				var stream = event.Stream
				stream = strings.ReplaceAll(stream, "\n", "")
				stream = strings.ReplaceAll(stream, "\r", "")
				stream = strings.Trim(stream, " ")

				if stream == "" {
					continue
				}

				log.Printf("%v", stream)

				if event.Closed == true {
					return
				}
			}
		}
	}(ch)
}

func ExampleContainerBuilder_ImageBuildFromFolder() {
	var err error

	SaGarbageCollector()

	var container = ContainerBuilder{}
	// new image name delete:latest
	container.SetImageName("delete:latest")
	// set a folder path to make a new image
	container.SetBuildFolderPath("./test/server")
	container.MakeDefaultDockerfileForMeWithInstallExtras()
	// container name container_delete_server_after_test
	container.SetContainerName("container_delete_server_after_test")
	// set a waits for the text to appear in the standard container output to proceed [optional]
	container.SetWaitStringWithTimeout("starting server at port 3000", 10*time.Second)
	// change and open port 3000 to 3030
	container.AddPortToExpose("3000")
	// replace container folder /static to host folder ./test/static
	err = container.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./test/static", "/static")
	if err != nil {
		panic(err)
	}

	// show image build stram on std out
	ImageBuildViewer(container.GetChannelEvent())

	// inicialize container object
	err = container.Init()
	if err != nil {
		panic(err)
	}

	// todo: fazer o teste do inspect

	// builder new image from folder
	_, err = container.ImageBuildFromFolder()
	if err != nil {
		panic(err)
	}

	// build a new container from image
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		panic(err)
	}

	// Server is ready for use o port 3000

	// read server inside a container on address http://localhost:3000/
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

	SaGarbageCollector()

	// Output:
	// <html><body><p>C is life! Golang is a evolution of C</p></body></html>
}
