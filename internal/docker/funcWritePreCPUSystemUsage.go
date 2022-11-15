package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writePreCPUSystemUsage(file *os.File, stats *types.Stats) (tab bool, err error) {
	// System Usage. Linux only.
	if e.rowsToPrint&KLogColumnPreCPUSystemUsage == KLogColumnPreCPUSystemUsage {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.PreCPUStats.SystemUsage)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPreCPUSystemUsageComa != 0
	}

	return
}
