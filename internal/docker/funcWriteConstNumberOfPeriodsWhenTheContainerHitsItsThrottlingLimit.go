package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit(file *os.File) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Number of periods when the container hits its throttling limit.
	if e.rowsToPrint&KLogColumnNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit == KLogColumnNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit {
		_, err = file.Write([]byte("KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitComa != 0
	}

	return
}
