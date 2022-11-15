package docker

const (
	// KKiloByte
	//
	// English: 1024 Bytes multiplier
	//
	// Example:
	//   5 * KKiloByte = 5 KBytes
	//
	// Português: multiplicador de 1024 Bytes
	//
	// Exemplo:
	//   5 * KKiloByte = 5 KBytes
	KKiloByte = 1024

	// KMegaByte
	//
	// English: 1024 KBytes multiplier
	//
	// Example:
	//   5 * KMegaByte = 5 MBytes
	//
	// Português: multiplicador de 1024 KBytes
	//
	// Exemplo:
	//   5 * KMegaByte = 5 MBytes
	KMegaByte = 1024 * 1024

	// KGigaByte
	//
	// English: 1024 MBytes multiplier
	//
	// Example:
	//   5 * KGigaByte = 5 GBytes
	//
	// Português: multiplicador de 1024 MBytes
	//
	// Exemplo:
	//   5 * KGigaByte = 5 GBytes
	KGigaByte = 1024 * 1024 * 1024

	// KTeraByte (
	//
	// English: 1024 GBytes multiplier
	//
	// Example:
	//   5 * KTeraByte = 5 TBytes
	//
	// Português: multiplicador de 1024 GBytes
	//
	// Exemplo:
	//   5 * KTeraByte = 5 TBytes
	KTeraByte = 1024 * 1024 * 1024 * 1024

	// KLogColumnAll
	//
	// English: Enable all values to log
	KLogColumnAll = 0x7FFFFFFFFFFFFFF

	// KLogColumnReadingTime
	//
	// English: Reading time
	KLogColumnReadingTime = 0b0000000000000000000000000000000000000000000000000000000000000001
	KReadingTimeComa      = 0b0111111111111111111111111111111111111111111111111111111111111110

	KFilterLog              = 0b0000000000000000000000000000000000000000000000000000000000000010
	KLogColumnFilterLogComa = 0b0111111111111111111111111111111111111111111111111111111111111100

	// KLogColumnCurrentNumberOfOidsInTheCGroup
	//
	// English: Linux specific stats, not populated on Windows. Current is the number of pids in the cgroup
	KLogColumnCurrentNumberOfOidsInTheCGroup = 0b0000000000000000000000000000000000000000000000000000000000000100
	KCurrentNumberOfOidsInTheCGroupComa      = 0b0111111111111111111111111111111111111111111111111111111111111000

	// KLogColumnLimitOnTheNumberOfPidsInTheCGroup
	//
	// English: Linux specific stats, not populated on Windows. Limit is the hard limit on the number of pids in the cgroup. A "Limit" of 0 means that there is no limit.
	KLogColumnLimitOnTheNumberOfPidsInTheCGroup = 0b0000000000000000000000000000000000000000000000000000000000001000
	KLimitOnTheNumberOfPidsInTheCGroupComa      = 0b0111111111111111111111111111111111111111111111111111111111110000

	// KLogColumnTotalCPUTimeConsumed
	//
	// English: Total CPU time consumed. (Units: nanoseconds on Linux, Units: 100's of nanoseconds on Windows)
	KLogColumnTotalCPUTimeConsumed = 0b0000000000000000000000000000000000000000000000000000000000010000
	KTotalCPUTimeConsumedComa      = 0b0111111111111111111111111111111111111111111111111111111111100000

	// KLogColumnTotalCPUTimeConsumedPerCore
	//
	// English: Total CPU time consumed. (Units: nanoseconds on Linux, Units: 100's of nanoseconds on Windows)
	KLogColumnTotalCPUTimeConsumedPerCore = 0b0000000000000000000000000000000000000000000000000000000000100000
	KTotalCPUTimeConsumedPerCoreComa      = 0b0111111111111111111111111111111111111111111111111111111111000000

	// KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode
	//
	// English: Time spent by tasks of the cgroup in kernel mode (Units: nanoseconds on Linux). Time spent by all container processes in kernel mode (Units: 100's of nanoseconds on Windows.Not populated for Hyper-V Containers.)
	KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode = 0b0000000000000000000000000000000000000000000000000000000001000000
	KTimeSpentByTasksOfTheCGroupInKernelModeComa      = 0b0111111111111111111111111111111111111111111111111111111110000000

	// KLogColumnTimeSpentByTasksOfTheCGroupInUserMode
	//
	// English: Time spent by tasks of the cgroup in user mode (Units: nanoseconds on Linux). Time spent by all container processes in user mode (Units: 100's of nanoseconds on Windows. Not populated for Hyper-V Containers)
	KLogColumnTimeSpentByTasksOfTheCGroupInUserMode = 0b0000000000000000000000000000000000000000000000000000000010000000
	KTimeSpentByTasksOfTheCGroupInUserModeComa      = 0b0111111111111111111111111111111111111111111111111111111100000000

	// KLogColumnSystemUsage
	//
	// English: System Usage. Linux only.
	KLogColumnSystemUsage = 0b0000000000000000000000000000000000000000000000000000000100000000
	KSystemUsageComa      = 0b0111111111111111111111111111111111111111111111111111111000000000

	// KOnlineCPUs
	//
	// English: Online CPUs. Linux only.
	KLogColumnOnlineCPUs = 0b0000000000000000000000000000000000000000000000000000001000000000
	KOnlineCPUsComa      = 0b0111111111111111111111111111111111111111111111111111110000000000

	// KLogColumnNumberOfPeriodsWithThrottlingActive
	//
	// English: Throttling Data. Linux only. Number of periods with throttling active.
	KLogColumnNumberOfPeriodsWithThrottlingActive = 0b0000000000000000000000000000000000000000000000000000010000000000
	KNumberOfPeriodsWithThrottlingActiveComa      = 0b0111111111111111111111111111111111111111111111111111100000000000

	// KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit
	//
	// English: Throttling Data. Linux only. Number of periods when the container hits its throttling limit.
	KLogColumnNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit = 0b0000000000000000000000000000000000000000000000000000100000000000
	KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitComa      = 0b0111111111111111111111111111111111111111111111111111000000000000

	// KAggregateTimeTheContainerWasThrottledForInNanoseconds
	//
	// English: Throttling Data. Linux only. Aggregate time the container was throttled for in nanoseconds.
	KLogColumnAggregateTimeTheContainerWasThrottledForInNanoseconds = 0b0000000000000000000000000000000000000000000000000001000000000000
	KAggregateTimeTheContainerWasThrottledForInNanosecondsComa      = 0b0111111111111111111111111111111111111111111111111110000000000000

	// KLogColumnTotalPreCPUTimeConsumed
	//
	// English: Total CPU time consumed per core (Units: nanoseconds on Linux). Not used on Windows.
	KLogColumnTotalPreCPUTimeConsumed = 0b0000000000000000000000000000000000000000000000000010000000000000
	KTotalPreCPUTimeConsumedComa      = 0b0111111111111111111111111111111111111111111111111100000000000000

	// KLogColumnTotalPreCPUTimeConsumedPerCore
	//
	// English: Total CPU time consumed per core (Units: nanoseconds on Linux). Not used on Windows.
	KLogColumnTotalPreCPUTimeConsumedPerCore = 0b0000000000000000000000000000000000000000000000000100000000000000
	KTotalPreCPUTimeConsumedPerCoreComa      = 0b0111111111111111111111111111111111111111111111111000000000000000

	// KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInKernelMode
	//
	// English: Time spent by tasks of the cgroup in kernel mode (Units: nanoseconds on Linux) - Time spent by all container processes in kernel mode (Units: 100's of nanoseconds on Windows - Not populated for Hyper-V Containers.)
	KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInKernelMode = 0b0000000000000000000000000000000000000000000000001000000000000000
	KTimeSpentByPreCPUTasksOfTheCGroupInKernelModeComa      = 0b0111111111111111111111111111111111111111111111110000000000000000

	// KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInUserMode
	//
	// English: Time spent by tasks of the cgroup in user mode (Units: nanoseconds on Linux) - Time spent by all container processes in user mode (Units: 100's of nanoseconds on Windows. Not populated for Hyper-V Containers)
	KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInUserMode = 0b0000000000000000000000000000000000000000000000010000000000000000
	KTimeSpentByPreCPUTasksOfTheCGroupInUserModeComa      = 0b0111111111111111111111111111111111111111111111100000000000000000

	// KLogColumnPreCPUSystemUsage
	//
	// English: System Usage. (Linux only)
	KLogColumnPreCPUSystemUsage = 0b0000000000000000000000000000000000000000000000100000000000000000
	KPreCPUSystemUsageComa      = 0b0111111111111111111111111111111111111111111111000000000000000000

	// KLogColumnOnlinePreCPUs
	//
	// English: Online CPUs. (Linux only)
	KLogColumnOnlinePreCPUs = 0b0000000000000000000000000000000000000000000001000000000000000000
	KOnlinePreCPUsComa      = 0b0111111111111111111111111111111111111111111110000000000000000000

	// KLogColumnAggregatePreCPUTimeTheContainerWasThrottled
	//
	// English: Throttling Data. (Linux only) - Aggregate time the container was throttled for in nanoseconds
	KLogColumnAggregatePreCPUTimeTheContainerWasThrottled = 0b0000000000000000000000000000000000000000000010000000000000000000
	KAggregatePreCPUTimeTheContainerWasThrottledComa      = 0b0111111111111111111111111111111111111111111100000000000000000000

	// KLogColumnNumberOfPeriodsWithPreCPUThrottlingActive
	//
	// English: Throttling Data. (Linux only) - Number of periods with throttling active
	KLogColumnNumberOfPeriodsWithPreCPUThrottlingActive = 0b0000000000000000000000000000000000000000000100000000000000000000
	KNumberOfPeriodsWithPreCPUThrottlingActiveComa      = 0b0111111111111111111111111111111111111111111000000000000000000000

	// KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit
	//
	// English: Throttling Data. (Linux only) - Number of periods when the container hits its throttling limit.
	KLogColumnNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit = 0b0000000000000000000000000000000000000000001000000000000000000000
	KNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitComa      = 0b0111111111111111111111111111111111111111110000000000000000000000

	// KLogColumnCurrentResCounterUsageForMemory
	//
	// English: Current res_counter usage for memory
	KLogColumnCurrentResCounterUsageForMemory = 0b0000000000000000000000000000000000000000010000000000000000000000
	KCurrentResCounterUsageForMemoryComa      = 0b0111111111111111111111111111111111111111100000000000000000000000

	// KLogColumnMaximumUsageEverRecorded
	//
	// English: Maximum usage ever recorded
	KLogColumnMaximumUsageEverRecorded = 0b0000000000000000000000000000000000000000100000000000000000000000
	KMaximumUsageEverRecordedComa      = 0b0111111111111111111111111111111111111111000000000000000000000000

	// KLogColumnNumberOfTimesMemoryUsageHitsLimits
	//
	// English: Number of times memory usage hits limits
	KLogColumnNumberOfTimesMemoryUsageHitsLimits = 0b0000000000000000000000000000000000000001000000000000000000000000
	KNumberOfTimesMemoryUsageHitsLimitsComa      = 0b0111111111111111111111111111111111111110000000000000000000000000

	// KLogColumnMemoryLimit
	//
	// English: Memory limit
	KLogColumnMemoryLimit = 0b0000000000000000000000000000000000000010000000000000000000000000
	KMemoryLimitComa      = 0b0111111111111111111111111111111111111100000000000000000000000000

	// KLogColumnCommittedBytes
	//
	// English: Committed bytes
	KLogColumnCommittedBytes = 0b0000000000000000000000000000000000000100000000000000000000000000
	KCommittedBytesComa      = 0b0111111111111111111111111111111111111000000000000000000000000000

	// KLogColumnPeakCommittedBytes
	//
	// English: Peak committed bytes
	KLogColumnPeakCommittedBytes = 0b0000000000000000000000000000000000001000000000000000000000000000
	KPeakCommittedBytesComa      = 0b0111111111111111111111111111111111110000000000000000000000000000

	// KLogColumnPrivateWorkingSet
	//
	// English: Private working set
	KLogColumnPrivateWorkingSet = 0b0000000000000000000000000000000000010000000000000000000000000000
	KPrivateWorkingSetComa      = 0b0111111111111111111111111111111111100000000000000000000000000000

	KLogColumnBlkioIoServiceBytesRecursive = 0b0000000000000000000000000000000000100000000000000000000000000000
	KBlkioIoServiceBytesRecursiveComa      = 0b0111111111111111111111111111111111000000000000000000000000000000

	KLogColumnBlkioIoServicedRecursive = 0b0000000000000000000000000000000001000000000000000000000000000000
	KBlkioIoServicedRecursiveComa      = 0b0111111111111111111111111111111110000000000000000000000000000000

	KLogColumnBlkioIoQueuedRecursive = 0b0000000000000000000000000000000010000000000000000000000000000000
	KBlkioIoQueuedRecursiveComa      = 0b0111111111111111111111111111111100000000000000000000000000000000

	KLogColumnBlkioIoServiceTimeRecursive = 0b0000000000000000000000000000000100000000000000000000000000000000
	KBlkioIoServiceTimeRecursiveComa      = 0b0111111111111111111111111111111000000000000000000000000000000000

	KLogColumnBlkioIoWaitTimeRecursive = 0b0000000000000000000000000000001000000000000000000000000000000000
	KBlkioIoWaitTimeRecursiveComa      = 0b0111111111111111111111111111110000000000000000000000000000000000

	KLogColumnBlkioIoMergedRecursive = 0b0000000000000000000000000000010000000000000000000000000000000000
	KBlkioIoMergedRecursiveComa      = 0b0111111111111111111111111111100000000000000000000000000000000000

	KLogColumnBlkioIoTimeRecursive = 0b0000000000000000000000000000100000000000000000000000000000000000
	KBlkioIoTimeRecursiveComa      = 0b0111111111111111111111111111000000000000000000000000000000000000

	KLogColumnBlkioSectorsRecursive = 0b0000000000000000000000000001000000000000000000000000000000000000
	KBlkioSectorsRecursiveComa      = 0b0111111111111111111111111110000000000000000000000000000000000000

	// KLogColumnMacOsLogWithAllCores
	//
	// English: Mac OS Log
	KLogColumnMacOsLogWithAllCores = KLogColumnReadingTime |
		KLogColumnCurrentNumberOfOidsInTheCGroup |
		KLogColumnTotalCPUTimeConsumed |
		KLogColumnTotalCPUTimeConsumedPerCore |
		KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode |
		KLogColumnSystemUsage |
		KLogColumnOnlineCPUs |
		KLogColumnTotalPreCPUTimeConsumed |
		KLogColumnTotalPreCPUTimeConsumedPerCore |
		KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInKernelMode |
		KLogColumnPreCPUSystemUsage |
		KLogColumnOnlinePreCPUs |
		KLogColumnCurrentResCounterUsageForMemory |
		KLogColumnMaximumUsageEverRecorded |
		KLogColumnMemoryLimit |
		KLogColumnBlkioIoServiceBytesRecursive | // não aparece no mac
		KLogColumnBlkioIoServicedRecursive | // não aparece no mac
		KLogColumnBlkioIoQueuedRecursive | // não aparece no mac
		KLogColumnBlkioIoServiceTimeRecursive | // não aparece no mac
		KLogColumnBlkioIoWaitTimeRecursive | // não aparece no mac
		KLogColumnBlkioIoMergedRecursive | // não aparece no mac
		KLogColumnBlkioIoTimeRecursive | // não aparece no mac
		KLogColumnBlkioSectorsRecursive // não aparece no mac

	// KLogColumnMacOs
	//
	// English: Mac OS Log
	KLogColumnMacOs = KLogColumnReadingTime |
		KLogColumnCurrentNumberOfOidsInTheCGroup |
		KLogColumnTotalCPUTimeConsumed |
		KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode |
		KLogColumnSystemUsage |
		KLogColumnOnlineCPUs |
		KLogColumnTotalPreCPUTimeConsumed |
		KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInKernelMode |
		KLogColumnPreCPUSystemUsage |
		KLogColumnOnlinePreCPUs |
		KLogColumnCurrentResCounterUsageForMemory |
		KLogColumnMaximumUsageEverRecorded |
		KLogColumnMemoryLimit |
		KLogColumnBlkioIoServiceBytesRecursive | // não aparece no mac
		KLogColumnBlkioIoServicedRecursive | // não aparece no mac
		KLogColumnBlkioIoQueuedRecursive | // não aparece no mac
		KLogColumnBlkioIoServiceTimeRecursive | // não aparece no mac
		KLogColumnBlkioIoWaitTimeRecursive | // não aparece no mac
		KLogColumnBlkioIoMergedRecursive | // não aparece no mac
		KLogColumnBlkioIoTimeRecursive | // não aparece no mac
		KLogColumnBlkioSectorsRecursive // não aparece no mac

	KLogColumnWindows = KLogColumnReadingTime |
		KLogColumnCurrentNumberOfOidsInTheCGroup |
		KLogColumnTotalCPUTimeConsumed |
		KLogColumnTotalCPUTimeConsumedPerCore |
		KLogColumnTimeSpentByTasksOfTheCGroupInKernelMode |
		KLogColumnTimeSpentByTasksOfTheCGroupInUserMode |
		KLogColumnSystemUsage |
		KLogColumnOnlineCPUs |
		KLogColumnTotalPreCPUTimeConsumed |
		KLogColumnTotalPreCPUTimeConsumedPerCore |
		KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInKernelMode |
		KLogColumnTimeSpentByPreCPUTasksOfTheCGroupInUserMode |
		KLogColumnPreCPUSystemUsage |
		KLogColumnOnlinePreCPUs |
		KLogColumnCurrentResCounterUsageForMemory |
		KLogColumnMaximumUsageEverRecorded |
		KLogColumnMemoryLimit
)
