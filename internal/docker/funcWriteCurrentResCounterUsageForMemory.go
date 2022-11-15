package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeCurrentResCounterUsageForMemory(file *os.File, stats *types.Stats) (tab bool, err error) {
	// current res_counter usage for memory
	if e.rowsToPrint&KLogColumnCurrentResCounterUsageForMemory == KLogColumnCurrentResCounterUsageForMemory {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.Usage)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KCurrentResCounterUsageForMemoryComa != 0
	}

	return
}
