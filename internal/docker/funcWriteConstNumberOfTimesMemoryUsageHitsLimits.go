package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstNumberOfTimesMemoryUsageHitsLimits(file *os.File) (tab bool, err error) {
	// number of times memory usage hits limits.
	if e.rowsToPrint&KLogColumnNumberOfTimesMemoryUsageHitsLimits == KLogColumnNumberOfTimesMemoryUsageHitsLimits {
		_, err = file.Write([]byte("KNumberOfTimesMemoryUsageHitsLimits"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KNumberOfTimesMemoryUsageHitsLimitsComa != 0
	}

	return
}
