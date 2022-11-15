package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeContainerConstToFile(file *os.File, stats *types.Stats) (err error) {
	var tab bool

	// time ok
	tab, err = e.writeConstReadingTime(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	for _, v := range e.chaos.filterLog {
		if v.Label != "" {
			tab = true
			break
		}
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstFilterIntoLog(file, e.chaos.filterLog)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstCurrentNumberOfOidsInTheCGroup(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstLimitOnTheNumberOfPidsInTheCGroup(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTotalCPUTimeConsumed(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	if len(stats.CPUStats.CPUUsage.PercpuUsage) != 0 {
		e.logCpus = len(stats.CPUStats.CPUUsage.PercpuUsage)
	}

	tab, err = e.writeConstTotalCPUTimeConsumedPerCore(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTimeSpentByTasksOfTheCGroupInKernelMode(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTimeSpentByTasksOfTheCGroupInUserMode(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstSystemUsage(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstOnlineCPUs(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstNumberOfPeriodsWithThrottlingActive(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstAggregateTimeTheContainerWasThrottledForInNanoseconds(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTotalPreCPUTimeConsumed(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTotalPreCPUTimeConsumedPerCore(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTimeSpentByPreCPUTasksOfTheCGroupInKernelMode(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstTimeSpentByPreCPUTasksOfTheCGroupInUserMode(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstPreCPUSystemUsage(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstOnlinePreCPUs(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstAggregatePreCPUTimeTheContainerWasThrottled(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstNumberOfPeriodsWithPreCPUThrottlingActive(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstCurrentResCounterUsageForMemory(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstMaximumUsageEverRecorded(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstNumberOfTimesMemoryUsageHitsLimits(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstMemoryLimit(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstCommittedBytes(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstPeakCommittedBytes(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstPrivateWorkingSet(file)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	flagIoServiceBytesRecursive := len(stats.BlkioStats.IoServiceBytesRecursive) != 0
	flagIoServicedRecursive := len(stats.BlkioStats.IoServicedRecursive) != 0
	flagIoQueuedRecursive := len(stats.BlkioStats.IoQueuedRecursive) != 0
	flagIoServiceTimeRecursive := len(stats.BlkioStats.IoServiceTimeRecursive) != 0
	flagIoWaitTimeRecursive := len(stats.BlkioStats.IoWaitTimeRecursive) != 0
	flagIoMergedRecursive := len(stats.BlkioStats.IoMergedRecursive) != 0
	flagIoTimeRecursive := len(stats.BlkioStats.IoTimeRecursive) != 0
	flagSectorsRecursive := len(stats.BlkioStats.SectorsRecursive) != 0

	condensedFlagIoServiceBytesRecursive := flagIoServiceBytesRecursive || flagIoServicedRecursive || flagIoQueuedRecursive || flagIoServiceTimeRecursive || flagIoWaitTimeRecursive || flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoServicedRecursive := flagIoServicedRecursive || flagIoQueuedRecursive || flagIoServiceTimeRecursive || flagIoWaitTimeRecursive || flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoQueuedRecursive := flagIoQueuedRecursive || flagIoServiceTimeRecursive || flagIoWaitTimeRecursive || flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoServiceTimeRecursive := flagIoServiceTimeRecursive || flagIoWaitTimeRecursive || flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoWaitTimeRecursive := flagIoWaitTimeRecursive || flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoMergedRecursive := flagIoMergedRecursive || flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagIoTimeRecursive := flagIoTimeRecursive || flagSectorsRecursive
	condensedFlagSectorsRecursive := flagSectorsRecursive

	if tab == true && condensedFlagIoServiceBytesRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoServiceBytesRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoServicedRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoServicedRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoQueuedRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoQueuedRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoServiceTimeRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoServiceTimeRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoWaitTimeRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoWaitTimeRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoMergedRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoMergedRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagIoTimeRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	tab, err = e.writeConstBlkioIoTimeRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if tab == true && condensedFlagSectorsRecursive == true {
		_, err = file.Write([]byte(e.csvValueSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	_, err = e.writeConstBlkioSectorsRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	return
}
