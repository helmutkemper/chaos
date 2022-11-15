package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstCommittedBytes(file *os.File) (tab bool, err error) {
	// committed bytes
	if e.rowsToPrint&KLogColumnCommittedBytes == KLogColumnCommittedBytes {
		_, err = file.Write([]byte("KCommittedBytes"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KCommittedBytesComa != 0
	}

	return
}
