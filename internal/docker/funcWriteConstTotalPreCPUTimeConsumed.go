package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstTotalPreCPUTimeConsumed(file *os.File) (tab bool, err error) {
	// CPU Usage. Linux and Windows.
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	if e.rowsToPrint&KLogColumnTotalPreCPUTimeConsumed == KLogColumnTotalPreCPUTimeConsumed {
		_, err = file.Write([]byte("KTotalPreCPUTimeConsumed"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KTotalPreCPUTimeConsumedComa != 0
	}

	return
}
