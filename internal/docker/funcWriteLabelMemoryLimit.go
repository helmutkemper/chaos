package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelMemoryLimit(file *os.File) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnMemoryLimit == KLogColumnMemoryLimit {
		_, err = file.Write([]byte("Memory limit"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KMemoryLimitComa != 0
	}

	return
}
