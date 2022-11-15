package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstTotalPreCPUTimeConsumedPerCore(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnTotalPreCPUTimeConsumedPerCore == KLogColumnTotalPreCPUTimeConsumedPerCore {
		for cpuNumber := 0; cpuNumber != e.logCpus; cpuNumber += 1 {
			_, err = file.Write([]byte("KTotalPreCPUTimeConsumedPerCore"))
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
