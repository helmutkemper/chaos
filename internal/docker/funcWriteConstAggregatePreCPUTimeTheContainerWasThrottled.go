package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstAggregatePreCPUTimeTheContainerWasThrottled(file *os.File) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Aggregate time the container was throttled for in nanoseconds.
	if e.rowsToPrint&KLogColumnAggregatePreCPUTimeTheContainerWasThrottled == KLogColumnAggregatePreCPUTimeTheContainerWasThrottled {
		_, err = file.Write([]byte("KAggregatePreCPUTimeTheContainerWasThrottled"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KAggregatePreCPUTimeTheContainerWasThrottledComa != 0
	}

	return
}
