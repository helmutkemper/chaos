package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeTotalCPUTimeConsumedPerCore(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnTotalCPUTimeConsumedPerCore == KLogColumnTotalCPUTimeConsumedPerCore {
		// Total CPU time consumed per core (Linux). Not used on Windows.
		// Units: nanoseconds.
		if e.logCpus != 0 && len(stats.CPUStats.CPUUsage.PercpuUsage) == 0 {
			for i := 0; i != e.logCpus; i += 1 {
				_, err = file.Write([]byte{0x30})
				if err != nil {
					log.Printf("writeContainerLogToFile().error: %v", err.Error())
					util.TraceToLog()
					return
				}

				if i != e.logCpus-1 {
					_, err = file.Write([]byte(e.csvValueSeparator))
					if err != nil {
						log.Printf("writeContainerLogToFile().error: %v", err.Error())
						util.TraceToLog()
						return
					}
				}
			}
		} else if e.logCpus != 0 && len(stats.CPUStats.CPUUsage.PercpuUsage) == e.logCpus {

			for i, cpuTime := range stats.CPUStats.CPUUsage.PercpuUsage {
				_, err = file.Write([]byte(fmt.Sprintf("%v", cpuTime)))
				if err != nil {
					log.Printf("writeContainerLogToFile().error: %v", err.Error())
					util.TraceToLog()
					return
				}

				if i != e.logCpus-1 {
					_, err = file.Write([]byte(e.csvValueSeparator))
					if err != nil {
						log.Printf("writeContainerLogToFile().error: %v", err.Error())
						util.TraceToLog()
						return
					}
				}
			}
		}

		tab = e.rowsToPrint&KTotalCPUTimeConsumedPerCoreComa != 0
	}

	return
}
