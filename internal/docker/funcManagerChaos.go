package docker

import (
	"log"
	"strconv"
	"time"
)

// managerChaos
//
// English: manages the log and state of the container
//
// Português: gerencia o log e o estado do container
func (e *ContainerBuilder) managerChaos() {
	var err error
	var logs []byte
	var lineList [][]byte
	var line []byte
	var found bool
	var timeToNextEvent time.Duration
	var probability float64
	var lineNumber int
	var event Event

	var inspect ContainerInspect

	// English: Probability of container restart
	// Português: Probabilidade do container reiniciar
	probability = e.getProbalityNumber()

	inspect, err = e.ContainerInspect()
	if err != nil {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.Metadata = e.metadata
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
		event.Error = true
		if len(e.chaos.event) == 0 {
			e.chaos.event <- event
		}
		return
	}

	if e.verifyStatusError(inspect) == true {
		return
	}

	logs, err = e.GetContainerLog()
	if err != nil {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.Metadata = e.metadata
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
		event.Error = true
		if len(e.chaos.event) == 0 {
			e.chaos.event <- event
		}
		return
	}

	lineList = e.logsCleaner(logs)

	err = e.writeContainerLogToFile(e.chaos.logPath, lineList)
	if err != nil {
		_, lineNumber = e.traceCodeLine()
		event.clear()
		event.Metadata = e.metadata
		event.ContainerName = e.GetContainerName()
		event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
		event.Error = true
		if len(e.chaos.event) == 0 {
			e.chaos.event <- event
		}
		return
	}

	event.clear()
	line, e.chaos.foundSuccess = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterSuccess)
	if e.chaos.foundSuccess == true {
		event.Message = string(line)
	} else {
		e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterMonitor)

		line, e.chaos.foundFail = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterFail)
		if e.chaos.foundFail == true {
			event.Message = string(line)
		}
	}
	_, lineNumber = e.traceCodeLine()

	event.Metadata = e.metadata
	event.ContainerName = e.GetContainerName()
	event.Fail = e.chaos.foundFail
	event.Done = e.chaos.foundSuccess
	if len(e.chaos.event) == 0 {
		e.chaos.event <- event
	}

	if e.chaos.enableChaos == false {
		return
	}

	if e.chaos.chaosStarted == false {

		timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToStartChaos, e.chaos.minimumTimeToStartChaos)

		if e.chaos.filterToStart != nil && e.chaos.minimumTimeToStartChaos > 0 {

			_, found = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterToStart)
			if found == true {
				if e.chaos.serviceStartedAt.Add(timeToNextEvent).Before(time.Now()) == true {
					e.chaos.chaosStarted = true
				}
			}

		} else if e.chaos.filterToStart != nil {

			_, found = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterToStart)
			if found == true {
				e.chaos.chaosStarted = true
			}

		} else if e.chaos.serviceStartedAt.Add(timeToNextEvent).Before(time.Now()) == true {
			e.chaos.chaosStarted = true
		}

		if e.chaos.chaosStarted == true {
			timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToStartChaos, e.chaos.minimumTimeToStartChaos)
		} else {
			return
		}

	}

	if e.chaos.chaosCanRestartContainer == false {

		timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeBeforeRestart, e.chaos.minimumTimeBeforeRestart)

		if e.chaos.filterRestart != nil && e.chaos.minimumTimeBeforeRestart > 0 {

			_, found = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterRestart)
			if found == true {
				if e.chaos.serviceStartedAt.Add(timeToNextEvent).Before(time.Now()) == true {
					e.chaos.chaosCanRestartContainer = true
				}
			}

		} else if e.chaos.filterRestart != nil {

			_, found = e.logsSearchAndReplaceIntoText(&logs, lineList, e.chaos.filterRestart)
			if found == true {
				e.chaos.chaosCanRestartContainer = true
			}

		} else if e.chaos.serviceStartedAt.Add(timeToNextEvent).Before(time.Now()) == true {
			e.chaos.chaosCanRestartContainer = true
		}

	}

	//if e.chaos.containerStarted == false {
	//	return
	//}

	//var restartEnable = time.Now().After(e.chaos.minimumTimeBeforeRestart) == true || time.Now().Equal(e.chaos.minimumTimeBeforeRestart) == true

	if time.Now().After(e.chaos.eventNext) == true || time.Now().Equal(e.chaos.eventNext) == true {

		if e.chaos.containerPaused == true {

			theater.SetContainerUnPaused(e.chaos.sceneName)

			log.Printf("%v: unpause()", e.containerName)
			e.chaos.containerPaused = false
			err = e.ContainerUnpause()
			if err != nil {
				_, lineNumber = e.traceCodeLine()
				event.clear()
				event.Metadata = e.metadata
				event.ContainerName = e.GetContainerName()
				event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
				event.Error = true
				if len(e.chaos.event) == 0 {
					e.chaos.event <- event
				}
				return
			}
			timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToPause, e.chaos.minimumTimeToPause)
			e.chaos.eventNext = time.Now().Add(timeToNextEvent)

		} else if e.chaos.containerStopped == true {

			theater.SetContainerUnStopped(e.chaos.sceneName)

			log.Printf("%v: start()", e.containerName)
			e.chaos.containerStopped = false

			probability = e.getProbalityNumber()
			if e.network != nil && e.chaos.restartChangeIpProbability > 0.0 && e.chaos.restartChangeIpProbability >= probability {
				err = e.NetworkChangeIp()
				if err != nil {
					_, lineNumber = e.traceCodeLine()
					event.clear()
					event.Metadata = e.metadata
					event.ContainerName = e.GetContainerName()
					event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
					event.Error = true
					if len(e.chaos.event) == 0 {
						e.chaos.event <- event
					}
					return
				}
			}

			err = e.ContainerStart()
			if err != nil {
				_, lineNumber = e.traceCodeLine()
				event.clear()
				event.Metadata = e.metadata
				event.ContainerName = e.GetContainerName()
				event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
				event.Error = true
				if len(e.chaos.event) == 0 {
					e.chaos.event <- event
				}
				return
			}
			e.chaos.chaosCanRestartContainer = false
			e.chaos.serviceStartedAt = time.Now().Add(e.selectBetweenMaxAndMin(e.chaos.maximumTimeToRestart, e.chaos.minimumTimeToRestart))
			timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToPause, e.chaos.minimumTimeToPause)
			e.chaos.eventNext = time.Now().Add(timeToNextEvent)

		} else if e.chaos.chaosCanRestartContainer == true && e.chaos.restartProbability != 0.0 && e.chaos.restartProbability >= probability && e.chaos.restartLimit > 0 {

			if e.chaos.disableStopContainer == true || theater.SetContainerStopped(e.chaos.sceneName) == true {
				return
			}

			log.Printf("%v: stop()", e.containerName)
			e.chaos.containerStopped = true
			err = e.ContainerStop()
			if err != nil {
				_, lineNumber = e.traceCodeLine()
				event.clear()
				event.Metadata = e.metadata
				event.ContainerName = e.GetContainerName()
				event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
				event.Error = true
				if len(e.chaos.event) == 0 {
					e.chaos.event <- event
				}
				return
			}
			e.chaos.restartLimit -= 1
			timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToRestart, e.chaos.minimumTimeToRestart)
			e.chaos.eventNext = time.Now().Add(timeToNextEvent)

		} else {

			if e.chaos.disablePauseContainer == true || theater.SetContainerPaused(e.chaos.sceneName) == true {
				return
			}

			log.Printf("%v: pause()", e.containerName)
			e.chaos.containerPaused = true
			err = e.ContainerPause()
			if err != nil {
				_, lineNumber = e.traceCodeLine()
				event.clear()
				event.Metadata = e.metadata
				event.ContainerName = e.GetContainerName()
				event.Message = "[" + strconv.Itoa(lineNumber) + "]: " + err.Error()
				event.Error = true
				if len(e.chaos.event) == 0 {
					e.chaos.event <- event
				}
				return
			}
			timeToNextEvent = e.selectBetweenMaxAndMin(e.chaos.maximumTimeToUnpause, e.chaos.minimumTimeToUnpause)
			e.chaos.eventNext = time.Now().Add(timeToNextEvent)

		}
	}
}
