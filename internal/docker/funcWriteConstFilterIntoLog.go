package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstFilterIntoLog(file *os.File, filter []LogFilter) (tab bool, err error) {
	var lineToFile = make([]byte, 0)

	for filterLine := 0; filterLine != len(filter); filterLine += 1 {
		if filter[filterLine].Label == "" {
			continue
		}

		lineToFile = append(lineToFile, []byte("")...)
		lineToFile = append(lineToFile, []byte(e.csvValueSeparator)...)

		tab = e.rowsToPrint&KLogColumnFilterLogComa != 0
	}

	lineToFile = bytes.TrimSuffix(lineToFile, []byte(e.csvValueSeparator))
	_, err = file.Write(lineToFile)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	return
}
