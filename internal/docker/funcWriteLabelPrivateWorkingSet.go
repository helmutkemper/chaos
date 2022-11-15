package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelPrivateWorkingSet(file *os.File) (tab bool, err error) {
	// private working set
	if e.rowsToPrint&KLogColumnPrivateWorkingSet == KLogColumnPrivateWorkingSet {
		_, err = file.Write([]byte("Private working set"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPrivateWorkingSetComa != 0
	}

	return
}
