package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstReadingTime(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnReadingTime == KLogColumnReadingTime {
		_, err = file.Write([]byte("KReadingTime"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KReadingTimeComa != 0
	}

	return
}
