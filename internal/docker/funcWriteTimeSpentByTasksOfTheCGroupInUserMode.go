package docker

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeTimeSpentByTasksOfTheCGroupInUserMode(file *os.File, stats *types.Stats) (tab bool, err error) {
	// Time spent by tasks of the cgroup in user mode (Linux).
	// Time spent by all container processes in user mode (Windows).
	// Units: nanoseconds (Linux).
	// Units: 100's of nanoseconds (Windows). Not populated for Hyper-V Containers
	if e.rowsToPrint&KLogColumnTimeSpentByTasksOfTheCGroupInUserMode == KLogColumnTimeSpentByTasksOfTheCGroupInUserMode {
		_, err = file.Write([]byte(fmt.Sprintf("%v", stats.CPUStats.CPUUsage.UsageInUsermode)))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		tab = e.rowsToPrint&KTimeSpentByTasksOfTheCGroupInUserModeComa != 0
	}

	return
}
