# iotmaker.docker.builder.golang.dockerfile

https://github.com/helmutkemper/iotmaker.docker.builder

This module is a standard **Golang** Dockerfile generator used in the **iotmaker.docker.builder** project.

Este módulo é um gerador de **Dockerfile** padrão para **Golang** usado no projeto **iotmaker.docker.builder**.

## iotmaker.docker.builder

Golang container generator to be integrated into golang codes in a simple and practical way, allowing the creation of 
integration tests in unit test format, or to write your own container manager.

Gerador de container em golang para ser integrado em códigos golang de forma simples e prática, permitindo a criação
de testes de integração em formato de testes unitários, ou escrever seu próprio gerenciador de containers.

## Sample code

```golang
package main

import (
  "fmt"
  iotmakerdocker "github.com/helmutkemper/iotmaker.docker/v1.0.1"
  "github.com/helmutkemper/util"
  "io/ioutil"
  "log"
  "net/http"
  "strings"
  "time"
)

func ImageBuildViewer(ch *chan iotmakerdocker.ContainerPullStatusSendToChannel) {
  go func(ch *chan iotmakerdocker.ContainerPullStatusSendToChannel) {
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

func main() {
  var err error

  var container = ContainerBuilder{}
  // new image name delete:latest
  container.SetImageName("delete:latest")
  // set a folder path to make a new image
  container.SetBuildFolderPath("./test")
  container.MakeDefaultDockerfileForMe()
  // container name container_delete_server_after_test
  container.SetContainerName("container_delete_server_after_test")
  // set a waits for the text to appear in the standard container output to proceed [optional]
  container.SetWaitStringWithTimeout("starting server at port 3000", 10*time.Second)
  // change and open port 3000 to 3030
  container.AddPortToOpen("3000")
  // replace container folder /static to host folder ./test/static
  err = container.AddFiileOrFolderToLinkBetweenConputerHostAndContainer("./test/static", "/static")
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

  // builder new image from folder
  err = container.ImageBuildFromFolder()
  if err != nil {
    panic(err)
  }

  // build a new container from image
  err = container.ContainerBuildFromImage()
  if err != nil {
    panic(err)
  }

  // Server is ready for use o port 3000

}
```

### ./test/static/index.html file
```html
<!DOCTYPE html><html><body>server is running</body></html>
```

### ./test/go.mod file
```golang
module teste
go 1.16
```

### ./test/main.go
```golang
package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	fmt.Printf("starting server at port 3000\n")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
```