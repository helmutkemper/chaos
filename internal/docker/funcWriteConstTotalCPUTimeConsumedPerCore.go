package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstTotalCPUTimeConsumedPerCore(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnTotalCPUTimeConsumedPerCore == KLogColumnTotalCPUTimeConsumedPerCore {
		// Total CPU time consumed per core (Linux). Not used on Windows.
		// Units: nanoseconds.
		for cpuNumber := 0; cpuNumber != e.logCpus; cpuNumber += 1 {
			_, err = file.Write([]byte("KTotalCPUTimeConsumedPerCore"))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			if cpuNumber != e.logCpus-1 {
				_, err = file.Write([]byte(e.csvValueSeparator))
				if err != nil {
					log.Printf("writeContainerLogToFile().error: %v", err.Error())
					util.TraceToLog()
					return
				}
			}
		}

		tab = e.rowsToPrint&KTotalCPUTimeConsumedPerCoreComa != 0
	}

	return
}
