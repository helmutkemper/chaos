package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeCommittedBytes(file *os.File, stats *types.Stats) (tab bool, err error) {
	// committed bytes
	if e.rowsToPrint&KLogColumnCommittedBytes == KLogColumnCommittedBytes {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.Commit)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KCommittedBytesComa != 0
	}

	return
}
