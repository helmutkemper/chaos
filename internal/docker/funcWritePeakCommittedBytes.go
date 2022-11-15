package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writePeakCommittedBytes(file *os.File, stats *types.Stats) (tab bool, err error) {
	// peak committed bytes
	if e.rowsToPrint&KLogColumnPeakCommittedBytes == KLogColumnPeakCommittedBytes {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.CommitPeak)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPeakCommittedBytesComa != 0
	}

	return
}
