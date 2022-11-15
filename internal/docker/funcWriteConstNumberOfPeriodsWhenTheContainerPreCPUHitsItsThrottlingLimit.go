package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit(file *os.File) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Number of periods when the container hits its throttling limit.
	if e.rowsToPrint&KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit == KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit {
		_, err = file.Write([]byte("KNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitComa != 0
	}

	return
}
