package docker

import (
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstTimeSpentByTasksOfTheCGroupInKernelMode(file *os.File) (tab bool, err error) {
	// Time spent by tasks of the cgroup in kernel mode (Linux).
	// Time spent by all container processes in kernel mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers.
	if e.rowsToPrint&KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode == KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode {
		_, err = file.Write([]byte("KTimeSpentByTasksOfTheCGroupInKernelMode"))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KTimeSpentByTasksOfTheCGroupInKernelModeComa != 0
	}

	return
}
