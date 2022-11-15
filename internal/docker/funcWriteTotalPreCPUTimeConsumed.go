package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeTotalPreCPUTimeConsumed(file *os.File, stats *types.Stats) (tab bool, err error) {
	// CPU Usage. Linux and Windows.
	// Total CPU time consumed.
	// Units: nanoseconds (Linux)
	// Units: 100's of nanoseconds (Windows)
	if e.rowsToPrint&KLogColumnTotalPreCPUTimeConsumed == KLogColumnTotalPreCPUTimeConsumed {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.PreCPUStats.CPUUsage.TotalUsage)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KTotalPreCPUTimeConsumedComa != 0
	}

	return
}
