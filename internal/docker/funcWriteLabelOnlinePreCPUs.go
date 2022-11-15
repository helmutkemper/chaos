package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelOnlinePreCPUs(file *os.File) (tab bool, err error) {
	// Online CPUs. Linux only.
	if e.rowsToPrint&KLogColumnOnlinePreCPUs == KLogColumnOnlinePreCPUs {
		_, err = file.Write([]byte("Online CPUs. (Linux only)"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KOnlinePreCPUsComa != 0
	}

	return
}
