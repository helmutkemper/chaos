package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelPreCPUSystemUsage(file *os.File) (tab bool, err error) {
	// System Usage. Linux only.
	if e.rowsToPrint&KLogColumnPreCPUSystemUsage == KLogColumnPreCPUSystemUsage {
		_, err = file.Write([]byte("System Usage. (Linux only)"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPreCPUSystemUsageComa != 0
	}

	return
}
