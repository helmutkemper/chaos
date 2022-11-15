package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelSystemUsage(file *os.File) (tab bool, err error) {
	// System Usage. Linux only.
	if e.rowsToPrint&KLogColumnSystemUsage == KLogColumnSystemUsage {
		_, err = file.Write([]byte("System Usage. Linux only."))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KSystemUsageComa != 0
	}

	return
}
