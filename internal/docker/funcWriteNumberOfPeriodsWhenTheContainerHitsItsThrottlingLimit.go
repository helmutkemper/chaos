package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Number of periods when the container hits its throttling limit.
	if e.rowsToPrint&KLogColumnNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit == KLogColumnNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.CPUStats.ThrottlingData.ThrottledPeriods)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitComa != 0
	}

	return
}
