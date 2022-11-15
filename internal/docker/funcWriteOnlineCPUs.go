package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeOnlineCPUs(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Online CPUs. Linux only.
	if e.rowsToPrint&KLogColumnOnlineCPUs == KLogColumnOnlineCPUs {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.CPUStats.OnlineCPUs)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KOnlineCPUsComa != 0
	}

	return
}
