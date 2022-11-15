package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelNumberOfTimesMemoryUsageHitsLimits(file *os.File) (tab bool, err error) {
	// number of times memory usage hits limits.
	if e.rowsToPrint&KLogColumnNumberOfTimesMemoryUsageHitsLimits == KLogColumnNumberOfTimesMemoryUsageHitsLimits {
		_, err = file.Write([]byte("Number of times memory usage hits limits."))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfTimesMemoryUsageHitsLimitsComa != 0
	}

	return
}
