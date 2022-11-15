package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelReadingTime(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnReadingTime == KLogColumnReadingTime {
		_, err = file.Write([]byte("Reading time"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
		tab = e.rowsToPrint&KReadingTimeComa != 0
	}

	return
}
