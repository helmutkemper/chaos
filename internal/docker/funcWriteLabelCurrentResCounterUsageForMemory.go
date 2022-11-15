package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelCurrentResCounterUsageForMemory(file *os.File) (tab bool, err error) {
	// current res_counter usage for memory
	if e.rowsToPrint&KLogColumnCurrentResCounterUsageForMemory == KLogColumnCurrentResCounterUsageForMemory {
		_, err = file.Write([]byte("Current res_counter usage for memory"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KCurrentResCounterUsageForMemoryComa != 0
	}

	return
}
