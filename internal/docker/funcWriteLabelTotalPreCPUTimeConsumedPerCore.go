package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelTotalPreCPUTimeConsumedPerCore(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnTotalPreCPUTimeConsumedPerCore == KLogColumnTotalPreCPUTimeConsumedPerCore {
		for cpuNumber := 0; cpuNumber != e.logCpus; cpuNumber += 1 {
			_, err = file.Write([]byte(fmt.Sprintf("Total CPU time consumed per core (Units: nanoseconds on Linux). Not used on Windows. CPU: %v", cpuNumber)))
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

		tab = e.rowsToPrint&KTotalPreCPUTimeConsumedPerCoreComa != 0
	}

	return
}
