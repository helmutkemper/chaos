package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"strings"
)

// imagePull
//
// English:
//
//	Downloads the image to be mounted. (equivalent to the docker pull image command)
//
//	 Output:
//	   err: standart error object
//
// Português:
//
//	Baixa a imagem a ser montada. (equivale ao comando docker pull image)
//
//	 Saída:
//	   err: objeto de erro padrão
func (e *ContainerBuilder) imagePull() (err error) {
	if e.printBuildOutput == true {
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
		}(&e.changePointer)
	}

	e.imageID, e.imageName, err = e.dockerSys.ImagePull(e.imageName, &e.changePointer)
	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
