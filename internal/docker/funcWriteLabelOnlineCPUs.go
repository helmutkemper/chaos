package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelOnlineCPUs(file *os.File) (tab bool, err error) {
	// Online CPUs. Linux only.
	if e.rowsToPrint&KLogColumnOnlineCPUs == KLogColumnOnlineCPUs {
		_, err = file.Write([]byte("Online CPUs. Linux only."))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KOnlineCPUsComa != 0
	}

	return
}
