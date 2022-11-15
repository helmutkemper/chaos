package docker

import (
	"github.com/helmutkemper/util"
	"log"
)

func (e *ContainerBuilder) stopMonitorAfterStopped() (err error) {
	if e.containerID == "" {
		err = e.getIdByContainerName()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	e.chaos.linear = true

	if e.chaos.containerPaused == true {

		theater.SetContainerUnPaused(e.chaos.sceneName)
		log.Printf("%v: unpause()", e.containerName)
		e.chaos.containerPaused = false

		err = e.ContainerUnpause()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	if e.chaos.containerStopped == true {

		theater.SetContainerUnStopped(e.chaos.sceneName)
		log.Printf("%v: start()", e.containerName)
		e.chaos.containerStopped = false

		err = e.ContainerStart()
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	return
}
