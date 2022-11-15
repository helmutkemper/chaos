package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeOnlinePreCPUs(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Online CPUs. Linux only.
	if e.rowsToPrint&KLogColumnOnlinePreCPUs == KLogColumnOnlinePreCPUs {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.PreCPUStats.OnlineCPUs)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KOnlinePreCPUsComa != 0
	}

	return
}
