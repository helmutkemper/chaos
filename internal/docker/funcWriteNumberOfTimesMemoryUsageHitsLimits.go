package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeNumberOfTimesMemoryUsageHitsLimits(file *os.File, stats *types.Stats) (tab bool, err error) {
	// number of times memory usage hits limits.
	if e.rowsToPrint&KLogColumnNumberOfTimesMemoryUsageHitsLimits == KLogColumnNumberOfTimesMemoryUsageHitsLimits {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.Failcnt)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfTimesMemoryUsageHitsLimitsComa != 0
	}

	return
}
