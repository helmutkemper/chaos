package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelCurrentNumberOfOidsInTheCGroup(file *os.File) (tab bool, err error) {
	// Linux specific stats, not populated on Windows.
	// Current is the number of pids in the cgroup
	if e.rowsToPrint&KLogColumnCurrentNumberOfOidsInTheCGroup == KLogColumnCurrentNumberOfOidsInTheCGroup {
		_, err = file.Write([]byte("Linux specific stats - not populated on Windows. Current is the number of pids in the cgroup"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KCurrentNumberOfOidsInTheCGroupComa != 0
	}

	return
}
