package docker

import (
	"bytes"
	pb "github.com/helmutkemper/iotmaker.docker.problem"
	"github.com/helmutkemper/util"
	"io"
	"io/fs"
	"log"
	"os"
	"regexp"
	"testing"
)

const (
	KLogReadingTimeLabel = "Reading time"
	KLogReadingTimeValue = "KReadingTime"
	//deixou de funcionar após atualização do docker
	KLogReadingTimeRegexp = "\\d{4}-\\d{2}-\\d{2}\\s\\d{2}:\\d{2}:\\d{2}.\\d{4,}(\\s\\+\\d{4})?\\sUTC"
	//KLogReadingTimeRegexp = "\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.\\d{4,}Z\\s+\\d{4}/\\d{2}/\\d{2}\\s\\d{2}:\\d{2}:\\d{2}"

	KLogCurrentNumberOfOidsInTheCGroupLabel  = "Linux specific stats - not populated on Windows. Current is the number of pids in the cgroup"
	KLogCurrentNumberOfOidsInTheCGroupValue  = "KCurrentNumberOfOidsInTheCGroup"
	KLogCurrentNumberOfOidsInTheCGroupRegexp = "\\d+"

	KLogLimitOnTheNumberOfPidsInTheCGroupLabel  = "Linux specific stats. Not populated on Windows. Limit is the hard limit on the number of pids in the cgroup. A \"Limit\" of 0 means that there is no limit."
	KLogLimitOnTheNumberOfPidsInTheCGroupValue  = "KLimitOnTheNumberOfPidsInTheCGroup"
	KLogLimitOnTheNumberOfPidsInTheCGroupRegexp = "\\d+"

	KLogTotalCPUTimeConsumedLabel  = "Total CPU time consumed. (Units: nanoseconds on Linux - Units: 100's of nanoseconds on Windows)"
	KLogTotalCPUTimeConsumedValue  = "KTotalCPUTimeConsumed"
	KLogTotalCPUTimeConsumedRegexp = "\\d+"

	KLogTotalCPUTimeConsumedPerCoreLabel  = "Total CPU time consumed per core (Units: nanoseconds on Linux). Not used on Windows."
	KLogTotalCPUTimeConsumedPerCoreValue  = "KTotalCPUTimeConsumedPerCore"
	KLogTotalCPUTimeConsumedPerCoreRegexp = "\\d+"

	KLogTimeSpentByTasksOfTheCGroupInKernelModeLabel  = "Time spent by tasks of the cgroup in kernel mode (Units: nanoseconds on Linux). Time spent by all container processes in kernel mode (Units: 100's of nanoseconds on Windows.Not populated for Hyper-V Containers.)."
	KLogTimeSpentByTasksOfTheCGroupInKernelModeValue  = "KTimeSpentByTasksOfTheCGroupInKernelMode"
	KLogTimeSpentByTasksOfTheCGroupInKernelModeRegexp = "\\d+"

	KLogTimeSpentByTasksOfTheCGroupInUserModeLabel  = "Time spent by tasks of the cgroup in user mode (Units: nanoseconds on Linux). Time spent by all container processes in user mode (Units: 100's of nanoseconds on Windows. Not populated for Hyper-V Containers)."
	KLogTimeSpentByTasksOfTheCGroupInUserModeValue  = "KTimeSpentByTasksOfTheCGroupInUserMode"
	KLogTimeSpentByTasksOfTheCGroupInUserModeRegexp = "\\d+"

	KLogSystemUsageLabel  = "System Usage. Linux only."
	KLogSystemUsageValue  = "KSystemUsage"
	KLogSystemUsageRegexp = "\\d+"

	KLogOnlineCPUsLabel  = "Online CPUs. Linux only."
	KLogOnlineCPUsValue  = "KOnlineCPUs"
	KLogOnlineCPUsRegexp = "\\d+"

	KLogNumberOfPeriodsWithThrottlingActiveLabel  = "Throttling Data. Linux only. Number of periods with throttling active."
	KLogNumberOfPeriodsWithThrottlingActiveValue  = "KNumberOfPeriodsWithThrottlingActive"
	KLogNumberOfPeriodsWithThrottlingActiveRegexp = "\\d+"

	KLogNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitLabel  = "Throttling Data. Linux only. Number of periods when the container hits its throttling limit."
	KLogNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitValue  = "KNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimit"
	KLogNumberOfPeriodsWhenTheContainerHitsItsThrottlingLimitRegexp = "\\d+"

	KLogAggregateTimeTheContainerWasThrottledForInNanosecondsLabel  = "Throttling Data. Linux only. Aggregate time the container was throttled for in nanoseconds."
	KLogAggregateTimeTheContainerWasThrottledForInNanosecondsValue  = "KAggregateTimeTheContainerWasThrottledForInNanoseconds"
	KLogAggregateTimeTheContainerWasThrottledForInNanosecondsRegexp = "\\d+"

	KLogTotalPreCPUTimeConsumedLabel  = "Total CPU time consumed. (Units: nanoseconds on Linux. Units: 100's of nanoseconds on Windows)"
	KLogTotalPreCPUTimeConsumedValue  = "KTotalPreCPUTimeConsumed"
	KLogTotalPreCPUTimeConsumedRegexp = "\\d+"

	KLogTotalPreCPUTimeConsumedPerCoreLabel  = "Total CPU time consumed per core (Units: nanoseconds on Linux). Not used on Windows."
	KLogTotalPreCPUTimeConsumedPerCoreValue  = "KTotalPreCPUTimeConsumedPerCore"
	KLogTotalPreCPUTimeConsumedPerCoreRegexp = "\\d+"

	KLogTimeSpentByPreCPUTasksOfTheCGroupInKernelModeLabel  = "Time spent by tasks of the cgroup in kernel mode (Units: nanoseconds on Linux) - Time spent by all container processes in kernel mode (Units: 100's of nanoseconds on Windows - Not populated for Hyper-V Containers.)"
	KLogTimeSpentByPreCPUTasksOfTheCGroupInKernelModeValue  = "KTimeSpentByPreCPUTasksOfTheCGroupInKernelMode"
	KLogTimeSpentByPreCPUTasksOfTheCGroupInKernelModeRegexp = "\\d+"

	KLogTimeSpentByPreCPUTasksOfTheCGroupInUserModeLabel  = "Time spent by tasks of the cgroup in user mode (Units: nanoseconds on Linux) - Time spent by all container processes in user mode (Units: 100's of nanoseconds on Windows. Not populated for Hyper-V Containers)"
	KLogTimeSpentByPreCPUTasksOfTheCGroupInUserModeValue  = "KTimeSpentByPreCPUTasksOfTheCGroupInUserMode"
	KLogTimeSpentByPreCPUTasksOfTheCGroupInUserModeRegexp = "\\d+"

	KLogPreCPUSystemUsageLabel  = "System Usage. (Linux only)"
	KLogPreCPUSystemUsageValue  = "KPreCPUSystemUsage"
	KLogPreCPUSystemUsageRegexp = "\\d+"

	KLogOnlinePreCPUsLabel  = "Online CPUs. (Linux only)"
	KLogOnlinePreCPUsValue  = "KOnlinePreCPUs"
	KLogOnlinePreCPUsRegexp = "\\d+"

	KLogAggregatePreCPUTimeTheContainerWasThrottledLabel  = "Throttling Data. (Linux only) - Aggregate time the container was throttled for in nanoseconds."
	KLogAggregatePreCPUTimeTheContainerWasThrottledValue  = "KAggregatePreCPUTimeTheContainerWasThrottled"
	KLogAggregatePreCPUTimeTheContainerWasThrottledRegexp = "\\d+"

	KLogNumberOfPeriodsWithPreCPUThrottlingActiveLabel  = "Throttling Data. (Linux only) - Number of periods with throttling active."
	KLogNumberOfPeriodsWithPreCPUThrottlingActiveValue  = "KNumberOfPeriodsWithPreCPUThrottlingActive"
	KLogNumberOfPeriodsWithPreCPUThrottlingActiveRegexp = "\\d+"

	KLogNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitLabel  = "Throttling Data. (Linux only) - Number of periods when the container hits its throttling limit."
	KLogNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitValue  = "KNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimit"
	KLogNumberOfPeriodsWhenTheContainerPreCPUHitsItsThrottlingLimitRegexp = "\\d+"

	KLogCurrentResCounterUsageForMemoryLabel  = "Current res_counter usage for memory"
	KLogCurrentResCounterUsageForMemoryValue  = "KCurrentResCounterUsageForMemory"
	KLogCurrentResCounterUsageForMemoryRegexp = "\\d+"

	KLogMaximumUsageEverRecordedLabel  = "Maximum usage ever recorded."
	KLogMaximumUsageEverRecordedValue  = "KMaximumUsageEverRecorded"
	KLogMaximumUsageEverRecordedRegexp = "\\d+"

	KLogNumberOfTimesMemoryUsageHitsLimitsLabel  = "Number of times memory usage hits limits."
	KLogNumberOfTimesMemoryUsageHitsLimitsValue  = "KNumberOfTimesMemoryUsageHitsLimits"
	KLogNumberOfTimesMemoryUsageHitsLimitsRegexp = "\\d+"

	KLogMemoryLimitLabel  = "Memory limit"
	KLogMemoryLimitValue  = "KMemoryLimit"
	KLogMemoryLimitRegexp = "\\d+"

	KLogCommittedBytesLabel  = "Committed bytes"
	KLogCommittedBytesValue  = "KCommittedBytes"
	KLogCommittedBytesRegexp = "\\d+"

	KLogPeakCommittedBytesLabel  = "Peak committed bytes"
	KLogPeakCommittedBytesValue  = "KPeakCommittedBytes"
	KLogPeakCommittedBytesRegexp = "\\d+"

	KLogPrivateWorkingSetLabel  = "Private working set"
	KLogPrivateWorkingSetValue  = "KPrivateWorkingSet"
	KLogPrivateWorkingSetRegexp = "\\d+"

	kLogHeaderLine = 0
	kLogLabelLine  = 1
)

type parserLog struct {
	Label string
	Key   string
	Rule  string
}

type TestContainerLog struct {
	data [][][]byte
}

func (e *TestContainerLog) makeTest(path string, listUnderTest *[]parserLog, t *testing.T) (problem pb.Problem) {
	problem = e.fileToDataFormat(path, t)
	if problem != nil {
		var file, funcName string
		var line int
		file, line, funcName, _ = problem.Trace()
		log.Printf("Error: %v", problem.Error())
		log.Printf("Cause: %v", problem.Cause())
		log.Printf("File: %v", file)
		log.Printf("Function: [%v]: %v", line, funcName)
		t.Fail()
		return
	}

	problem = e.proccessKeyList(listUnderTest, t)
	if problem != nil {
		var file, funcName string
		var line int
		file, line, funcName, _ = problem.Trace()
		log.Printf("Error: %v", problem.Error())
		log.Printf("Cause: %v", problem.Cause())
		log.Printf("File: %v", file)
		log.Printf("Function: [%v]: %v", line, funcName)
		t.Fail()
		return
	}

	return
}

func (e *TestContainerLog) fileToDataFormat(path string, t *testing.T) (problem pb.Problem) {
	var err error
	var f *os.File
	f, err = os.OpenFile(path, os.O_RDONLY, fs.ModePerm)
	if err != nil {
		problem = pb.NewProblem(err.Error(), "")
		t.Fail()
		return
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			util.TraceToLog()
		}
	}(f)

	e.data = make([][][]byte, 0)

	var char = make([]byte, 1)
	var lineOfData = make([][]byte, 0)

	for {
		var lineToSplit = make([]byte, 0)
		for {
			_, err = f.Read(char)
			if err != nil && err != io.EOF {
				break
			}

			if err != nil {
				return
			}

			if bytes.Equal(char, []byte("\n")) == true {
				break
			}

			lineToSplit = append(lineToSplit, char[0])
		}

		lineOfData = bytes.Split(lineToSplit, []byte(","))

		e.data = append(e.data, lineOfData)

		if err != nil && err != io.EOF {
			break
		}
	}

	err = nil
	return
}

func (e TestContainerLog) findKeyIndex(key string, t *testing.T) (index int, problem pb.Problem) {
	if len(e.data) == 0 {
		problem = pb.NewProblem(
			"data log error",
			"fileToDataFormat(path) did not find the file or was not called first",
		)
		t.Fail()
		return
	}

	var value []byte
	var keyToBeFound = []byte(key)
	log.Printf("key: %s", key)
	log.Printf("e.data: %s", e.data)

	for index, value = range e.data[kLogHeaderLine] {
		log.Printf("index: %v", index)
		log.Printf("value: %s", value)
		if bytes.Equal(value, keyToBeFound) == true {
			return
		}
	}

	return
}

func (e TestContainerLog) proccessKeyList(listUnderTest *[]parserLog, t *testing.T) (problem pb.Problem) {
	var index int
	var match bool
	var err error

	for _, test := range *listUnderTest {
		log.Printf("test: %+v", test)
		index, problem = e.findKeyIndex(test.Key, t)
		if problem != nil {
			var file, funcName string
			var line int
			file, line, funcName, _ = problem.Trace()
			log.Printf("Error: %v", problem.Error())
			log.Printf("Cause: %v", problem.Cause())
			log.Printf("File: %v", file)
			log.Printf("Function: [%v]: %v", line, funcName)
			t.Fail()
			return
		}

		log.Printf("1: %s", e.data[kLogHeaderLine][index])
		log.Printf("2: %s", test.Key)
		if bytes.Equal(e.data[kLogHeaderLine][index], []byte(test.Label)) == false {
			log.Printf("log file: %s", e.data[kLogHeaderLine][index])
			log.Printf("key: %v", test.Key)
			log.Printf("index: %v", index)
			log.Printf("log header don't match")
			t.Fail()
			return
		}

		if bytes.HasPrefix(e.data[kLogLabelLine][index], []byte(test.Label)) == false {
			log.Printf("log file: %s", e.data[kLogLabelLine][index])
			log.Printf("label: %v", test.Label)
			log.Printf("index: %v", index)
			log.Printf("log label don't match")
			t.Fail()
			return
		}

		for k := range e.data {
			if k <= kLogLabelLine {
				continue
			}

			match, err = regexp.Match(test.Rule, e.data[k][index])
			if err != nil {
				log.Printf("error: %v", err)
				return
			}

			if match == false {
				log.Printf("value: %s", e.data[k][index])
				log.Printf("rula: %s", test.Rule)
				log.Printf("index: %v", index)
				log.Printf("log rule don't match regexp")
				t.Fail()
				return
			}
		}

	}

	return
}
