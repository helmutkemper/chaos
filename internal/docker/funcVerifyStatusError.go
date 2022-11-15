package docker

import (
	"strconv"
)

func (e *ContainerBuilder) verifyStatusError(inspect ContainerInspect) (hasError bool) {
	var event Event
	var lineNumber int

	if inspect.State.OOMKilled == true {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "OOMKilled"
		event.Error = true
		e.chaos.event <- event
		hasError = true
		return
	}

	if inspect.State.Dead == true {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "dead"
		event.Error = true
		e.chaos.event <- event
		hasError = true
		return
	}

	if e.chaos.containerStopped == false && inspect.State.ExitCode != 0 {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "exit code: " + strconv.Itoa(inspect.State.ExitCode)
		event.Error = true
		e.chaos.event <- event
		hasError = true
		return
	}

	if (e.chaos.containerStopped == true || e.chaos.containerPaused == true) != true {

		if inspect.State.Running == false {
			_, lineNumber = e.traceCodeLine()
			event.clear()
			event.ContainerName = e.GetContainerName()
			event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "not running"
			event.Error = true
			e.chaos.event <- event
			hasError = true
			return
		}

		if inspect.State.Paused == true {
			_, lineNumber = e.traceCodeLine()
			event.clear()
			event.ContainerName = e.GetContainerName()
			event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "paused"
			event.Error = true
			e.chaos.event <- event
			hasError = true
			return
		}

		if inspect.State.Restarting == true {
			_, lineNumber = e.traceCodeLine()
			event.clear()
			event.ContainerName = e.GetContainerName()
			event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + "restarting"
			event.Error = true
			e.chaos.event <- event
			hasError = true
			return
		}
	}

	return
}
