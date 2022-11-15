package docker

import (
	"github.com/helmutkemper/util"
	"log"
)

// ContainerBuildAndStartFromImage
//
// English:
//
//	Transforms an image downloaded by ImagePull() or created by ImageBuildFromFolder() into a container
//	and start it.
//
//	 Output:
//	   err: Default object error from golang
//
// Português:
//
//	Transforma uma imagem baixada por ImagePull() ou criada por ImageBuildFromFolder() em container e o
//	inicializa.
//
//	 Saída:
//	   err: Objeto padrão de erro golang
func (e *ContainerBuilder) ContainerBuildAndStartFromImage() (err error) {
	err = e.ContainerBuildWithoutStartingItFromImage()
	if err != nil {
		util.TraceToLog()
		return
	}

	err = e.ContainerStartAfterBuild()
	if err != nil {
		util.TraceToLog()
		return
	}

	if e.waitString != "" && e.waitStringTimeout == 0 {
		_, err = e.dockerSys.ContainerLogsWaitText(e.containerID, e.waitString, log.Writer())
		if err != nil {
			util.TraceToLog()
			return
		}

	} else if e.waitString != "" {
		_, err = e.dockerSys.ContainerLogsWaitTextWithTimeout(e.containerID, e.waitString, e.waitStringTimeout, log.Writer())
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	if e.network == nil {
		e.IPV4Address, err = e.FindCurrentIPV4Address()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	if len(*e.onContainerReady) == 0 {
		*e.onContainerReady <- true
	}

	return
}
