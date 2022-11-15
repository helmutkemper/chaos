package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLimitOnTheNumberOfPidsInTheCGroup(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Linux specific stats, not populated on Windows.
	// Limit is the hard limit on the number of pids in the cgroup.
	// A "Limit" of 0 means that there is no limit.
	if e.rowsToPrint&KLogColumnLimitOnTheNumberOfPidsInTheCGroup == KLogColumnLimitOnTheNumberOfPidsInTheCGroup {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.PidsStats.Limit)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KLimitOnTheNumberOfPidsInTheCGroupComa != 0
	}

	return
}
