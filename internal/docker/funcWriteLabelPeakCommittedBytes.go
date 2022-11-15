package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelPeakCommittedBytes(file *os.File) (tab bool, err error) {
	// peak committed bytes
	if e.rowsToPrint&KLogColumnPeakCommittedBytes == KLogColumnPeakCommittedBytes {
		_, err = file.Write([]byte("Peak committed bytes"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KPeakCommittedBytesComa != 0
	}

	return
}
