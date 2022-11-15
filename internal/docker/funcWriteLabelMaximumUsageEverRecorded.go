package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelMaximumUsageEverRecorded(file *os.File) (tab bool, err error) {
	// maximum usage ever recorded.
	if e.rowsToPrint&KLogColumnMaximumUsageEverRecorded == KLogColumnMaximumUsageEverRecorded {
		_, err = file.Write([]byte("Maximum usage ever recorded."))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KMaximumUsageEverRecordedComa != 0
	}

	return
}
