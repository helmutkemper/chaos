package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeContainerStatsToFile(file *os.File, stats *types.Stats, lineList *[][]byte) (err error) {
	var tab bool

	// time ok
	tab, err = e.writeReadingTime(file, stats)
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

	tab, err = e.writeFilterIntoLog(file, e.chaos.filterLog, lineList)
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

	tab, err = e.writeCurrentNumberOfOidsInTheCGroup(file, stats)
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

	tab, err = e.writeLimitOnTheNumberOfPidsInTheCGroup(file, stats)
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

	tab, err = e.writeTotalCPUTimeConsumed(file, stats)
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

	tab, err = e.writeTotalCPUTimeConsumedPerCore(file, stats)
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

	tab, err = e.writeTimeSpentByTasksOfTheCGroupInKernelMode(file, stats)
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

	tab, err = e.writeTimeSpentByTasksOfTheCGroupInUserMode(file, stats)
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

	tab, err = e.writeSystemUsage(file, stats)
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

	tab, err = e.writeOnlineCPUs(file, stats)
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

	tab, err = e.writeNumberOfPeriodsWithThrottlingActive(file, stats)
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

	tab, err = e.writeNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit(file, stats)
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

	tab, err = e.writeAggregateTimeTheContainerWasThrottledForInNanoseconds(file, stats)
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

	tab, err = e.writeTotalPreCPUTimeConsumed(file, stats)
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

	tab, err = e.writeTotalPreCPUTimeConsumedPerCore(file, stats)
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

	tab, err = e.writeTimeSpentByPreCPUTasksOfTheCGroupInKernelMode(file, stats)
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

	tab, err = e.writeTimeSpentByPreCPUTasksOfTheCGroupInUserMode(file, stats)
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

	tab, err = e.writePreCPUSystemUsage(file, stats)
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

	tab, err = e.writeOnlinePreCPUs(file, stats)
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

	tab, err = e.writeAggregatePreCPUTimeTheContainerWasThrottled(file, stats)
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

	tab, err = e.writeNumberOfPeriodsWithPreCPUThrottlingActive(file, stats)
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

	tab, err = e.writeNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit(file, stats)
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

	tab, err = e.writeCurrentResCounterUsageForMemory(file, stats)
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

	tab, err = e.writeMaximumUsageEverRecorded(file, stats)
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

	tab, err = e.writeNumberOfTimesMemoryUsageHitsLimits(file, stats)
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

	tab, err = e.writeMemoryLimit(file, stats)
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

	tab, err = e.writeCommittedBytes(file, stats)
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

	tab, err = e.writePeakCommittedBytes(file, stats)
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

	tab, err = e.writePrivateWorkingSet(file, stats)
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

	tab, err = e.writeBlkioIoServiceBytesRecursive(file, stats)
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

	tab, err = e.writeBlkioIoServicedRecursive(file, stats)
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

	tab, err = e.writeBlkioIoQueuedRecursive(file, stats)
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

	tab, err = e.writeBlkioIoServiceTimeRecursive(file, stats)
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

	tab, err = e.writeBlkioIoWaitTimeRecursive(file, stats)
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

	tab, err = e.writeBlkioIoMergedRecursive(file, stats)
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

	tab, err = e.writeBlkioIoTimeRecursive(file, stats)
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

	_, err = e.writeBlkioSectorsRecursive(file, stats)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	return
}
