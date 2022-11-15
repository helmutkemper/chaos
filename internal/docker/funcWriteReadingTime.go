package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeReadingTime(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnReadingTime == KLogColumnReadingTime {
		_, err = file.Write([]byte(stats.Read.String()))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
		tab = e.rowsToPrint&KReadingTimeComa != 0
	}

	return
}
