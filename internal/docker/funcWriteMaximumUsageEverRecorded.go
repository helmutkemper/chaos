package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeMaximumUsageEverRecorded(file *os.File, stats *types.Stats) (tab bool, err error) {
	// maximum usage ever recorded.
	if e.rowsToPrint&KLogColumnMaximumUsageEverRecorded == KLogColumnMaximumUsageEverRecorded {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.MaxUsage)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KMaximumUsageEverRecordedComa != 0
	}

	return
}
