package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Throttling Data. Linux only.
	// Number of periods when the container hits its throttling limit.
	if e.rowsToPrint&KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit == KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.PreCPUStats.ThrottlingData.ThrottledPeriods)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitComa != 0
	}

	return
}
