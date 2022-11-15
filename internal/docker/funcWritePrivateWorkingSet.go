package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writePrivateWorkingSet(file *os.File, stats *types.Stats) (tab bool, err error) {
	// private working set
	if e.rowsToPrint&KLogColumnPrivateWorkingSet == KLogColumnPrivateWorkingSet {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.PrivateWorkingSet)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPrivateWorkingSetComa != 0
	}

	return
}
