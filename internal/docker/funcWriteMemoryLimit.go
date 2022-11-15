package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeMemoryLimit(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnMemoryLimit == KLogColumnMemoryLimit {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.MemoryStats.Limit)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KMemoryLimitComa != 0
	}

	return
}
