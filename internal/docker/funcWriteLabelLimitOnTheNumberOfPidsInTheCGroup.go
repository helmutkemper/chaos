package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelLimitOnTheNumberOfPidsInTheCGroup(file *os.File) (tab bool, err error) {
	// Linux specific stats, not populated on Windows.
	// Limit is the hard limit on the number of pids in the cgroup.
	// A "Limit" of 0 means that there is no limit.
	if e.rowsToPrint&KLogColumnLimitOnTheNumberOfPidsInTheCGroup == KLogColumnLimitOnTheNumberOfPidsInTheCGroup {
		_, err = file.Write([]byte("Linux specific stats. Not populated on Windows. Limit is the hard limit on the number of pids in the cgroup. A \"Limit\" of 0 means that there is no limit."))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KLimitOnTheNumberOfPidsInTheCGroupComa != 0
	}

	return
}
