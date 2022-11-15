package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstAggregateTimeTheContainerWasThrottledForInNanoseconds(file *os.File) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Aggregate time the container was throttled for in nanoseconds.
	if e.rowsToPrint&KLogColumnAggregateTimeTheContainerWasThrottledForInNanoseconds == KLogColumnAggregateTimeTheContainerWasThrottledForInNanoseconds {
		_, err = file.Write([]byte("KAggregateTimeTheContainerWasThrottledForInNanoseconds"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KAggregateTimeTheContainerWasThrottledForInNanosecondsComa != 0
	}

	return
}
