package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelTotalPreCPUTimeConsumed(file *os.File) (tab bool, err error) {
	// CPU Usage. Linux and Windows.
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	if e.rowsToPrint&KLogColumnTotalPreCPUTimeConsumed == KLogColumnTotalPreCPUTimeConsumed {
		_, err = file.Write([]byte("Total CPU time consumed. (Units: nanoseconds on Linux. Units: 100's of nanoseconds on Windows)"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KTotalPreCPUTimeConsumedComa != 0
	}

	return
}
