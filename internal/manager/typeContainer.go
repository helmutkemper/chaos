package manager

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/mount"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/chaos/internal/builder"
	"github.com/helmutkemper/util"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dockerContainer "github.com/docker/docker/api/types/container"
)

type container struct {
	IPV4Address []string

	//port inside container and host computer port
	portsContainer []nat.Port
	portsHost      [][]int64

	volumeContainer []string
	volumeHost      [][]string
}

type ContainerFromImage struct {
	container

	manager *Manager

	imageId          string
	imageName        string
	containerName    string
	copies           int
	csvPath          string
	failPath         string
	failFlag         []string
	failLogsLastSize []int
}

func (el *ContainerFromImage) New(manager *Manager) {
	el.manager = manager
}

// Volumes
//
// List of volumes (mounts) used for the container
//
//	 Input:
//	   containerPath: folder or file path inside the container
//	   hostPath: list of folders or files within the host computer
//
//		Notes:
//		  * When `hostPath` receives one more value, each container created will receive a different value.
//		    - Imagine creating 3 containers and passing the values `pathA` and `pathB`. The first container created will
//		    receive `pathA`, the second, `pathB` and the third will not receive value;
//		    - Imagine creating 3 containers and passing the values `pathA`, `` and `pathB`. The first container created will
//		    receive `pathA`, the second will not receive value, and the third receive `pathB`.
func (el *ContainerFromImage) Volumes(containerPath string, hostPath ...string) (ref *ContainerFromImage) {
	var err error

	if el.volumeContainer == nil {
		el.volumeContainer = make([]string, 0)
		el.volumeHost = make([][]string, 0)
	}

	var path string
	var absolutePath []string
	for k := range hostPath {
		if hostPath[k] != "" {
			path, err = filepath.Abs(hostPath[k])
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("containerFromImage.Volumes().error: %v", err)
				return el
			}
		} else {
			path = ""
		}

		absolutePath = append(absolutePath, path)
	}

	el.volumeContainer = append(el.volumeContainer, containerPath)
	el.volumeHost = append(el.volumeHost, absolutePath)
	return el
}

// Ports
//
// Defines which port of the container will be exposed to the world
//
//	Input:
//	  containerProtocol: network protocol `tcp` or `utc`
//	  containerPort: port number on the container. eg: 27017 for MongoDB
//	  localPort: port number on the host computer. eg: 27017 for MongoDB
//
//	Notes:
//	  * When `localPort` receives one more value, each container created will receive a different value.
//	    - Imagine creating 3 containers and passing the values 27016 and 27015. The first container created will receive
//	    27016, the second, 27015 and the third will not receive value;
//	    - Imagine creating 3 containers and passing the values 27016, 0 and 27015. The first container created will
//	    receive 27016, the second will not receive value, and the third receive 27015.
func (el *ContainerFromImage) Ports(containerProtocol string, containerPort int64, localPort ...int64) (ref *ContainerFromImage) {
	if el.portsContainer == nil {
		el.portsContainer = make([]nat.Port, 0)
		el.portsHost = make([][]int64, 0)
	}

	port, err := nat.NewPort(containerProtocol, strconv.FormatInt(containerPort, 10))
	if err != nil {
		el.manager.ErrorCh <- fmt.Errorf("containerFromImage.ExposePorts().error: %v", err)
		return
	}

	el.portsContainer = append(el.portsContainer, port)

	el.portsHost = append(el.portsHost, localPort)
	return el
}

// OnBuild
//
// The ONBUILD instruction adds to the image a trigger instruction
// to be executed at a later time, when the image is used as the base for another build.
// The trigger will be executed in the context of the downstream build, as if it had been
// inserted immediately after the FROM instruction in the downstream Dockerfile.
//
// Any build instruction can be registered as a trigger.
//
// This is useful if you are building an image which will be used as a base to build other
// images, for example an application build environment or a daemon which may be
// customized with user-specific configuration.
//
// For example, if your image is a reusable Python application builder, it will require
// application source code to be added in a particular directory, and it might require a
// build script to be called after that. You can’t just call ADD and RUN now, because you
// don’t yet have access to the application source code, and it will be different for each
// application build. You could simply provide application developers with a boilerplate
// Dockerfile to copy-paste into their application, but that is inefficient, error-prone
// and difficult to update because it mixes with application-specific code.
//
// The solution is to use ONBUILD to register advance instructions to run later, during
// the next build stage.
//
// Here’s how it works:
//
// When it encounters an ONBUILD instruction, the builder adds a trigger to the metadata
// of the image being built. The instruction does not otherwise affect the current build.
// At the end of the build, a list of all triggers is stored in the image manifest, under
// the key OnBuild. They can be inspected with the docker inspect command.
// Later the image may be used as a base for a new build, using the FROM instruction. As
// part of processing the FROM instruction, the downstream builder looks for ONBUILD
// triggers, and executes them in the same order they were registered. If any of the
// triggers fail, the FROM instruction is aborted which in turn causes the build to fail.
// If all triggers succeed, the FROM instruction completes and the build continues as
// usual.
// Triggers are cleared from the final image after being executed. In other words they are
// not inherited by “grand-children” builds.
// For example you might add something like this:
//
// ONBUILD ADD . /app/src
// ONBUILD RUN /usr/local/bin/python-build --dir /app/src
//
//	Warning:
//	Chaining ONBUILD instructions using ONBUILD ONBUILD isn’t allowed.
//
//	Warning:
//	The ONBUILD instruction may not trigger FROM or MAINTAINER instructions.
//
// https://docs.docker.com/engine/reference/builder/#onbuild
func (el *ContainerFromImage) OnBuild(onBuild ...string) (ref *ContainerFromImage) {
	if len(onBuild) == 0 {
		onBuild = nil
	}

	el.manager.DockerSys[0].Config.OnBuild = onBuild
	return el
}

// HostName
//
// Defines the hostname of the container
func (el *ContainerFromImage) HostName(name string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.Hostname = name
	return el
}

// DomainName
//
// Defines the domain name of the container
func (el *ContainerFromImage) DomainName(name string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.Domainname = name
	return el
}

// User
//
// User that will run the command(s) inside the container, also support user:group
func (el *ContainerFromImage) User(name string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.User = name
	return el
}

// Tty
//
// Attach standard streams to a tty, including stdin if it is not closed
func (el *ContainerFromImage) Tty(tty bool) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.Tty = tty
	return el
}

// OpenStdin
//
// Open stdin
func (el *ContainerFromImage) OpenStdin(open bool) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.OpenStdin = open
	return el
}

// StdinOnce
//
// If true, close stdin after the 1 attached client disconnects
func (el *ContainerFromImage) StdinOnce(once bool) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.StdinOnce = once
	return el
}

// EnvironmentVar
//
// List of environment variable to set in the container
func (el *ContainerFromImage) EnvironmentVar(env ...string) (ref *ContainerFromImage) {
	if len(env) == 0 {
		env = nil
	}

	el.manager.DockerSys[0].Config.Env = env
	return el
}

// Cmd
//
// Command to run when starting the container
func (el *ContainerFromImage) Cmd(cmd ...string) (ref *ContainerFromImage) {
	if len(cmd) == 0 {
		cmd = nil
	}

	el.manager.DockerSys[0].Config.Cmd = cmd
	return el
}

// ArgsEscaped
//
// True if command is already escaped (meaning treat as a command line) (Windows specific)
func (el *ContainerFromImage) ArgsEscaped(argsEscaped bool) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.ArgsEscaped = argsEscaped
	return el
}

// WorkingDir
//
// Current directory (PWD) in the command will be launched
func (el *ContainerFromImage) WorkingDir(workingDir string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.WorkingDir = workingDir
	return el
}

// Entrypoint
//
// Entrypoint to run when starting the container
func (el *ContainerFromImage) Entrypoint(entrypoint ...string) (ref *ContainerFromImage) {
	if len(entrypoint) == 0 {
		entrypoint = nil
	}

	el.manager.DockerSys[0].Config.Entrypoint = entrypoint
	return el
}

// NetworkDisabled
//
// Is network disabled
func (el *ContainerFromImage) NetworkDisabled(disabled bool) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.NetworkDisabled = disabled
	return el
}

// MacAddress
//
// Mac Address of the container
func (el *ContainerFromImage) MacAddress(macAddress string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.MacAddress = macAddress
	return el
}

// Labels
//
// List of labels set to this container
func (el *ContainerFromImage) Labels(labels map[string]string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.Labels = labels
	return el
}

// StopSignal
//
// Signal to stop a container
func (el *ContainerFromImage) StopSignal(signal string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.StopSignal = signal
	return el
}

// StopTimeout
//
// Timeout to stop a container
func (el *ContainerFromImage) StopTimeout(timeout time.Duration) (ref *ContainerFromImage) {
	timeoutRef := int(timeout)
	el.manager.DockerSys[0].Config.StopTimeout = &timeoutRef
	return el
}

// Shell
//
// Shell for shell-form of RUN, CMD, ENTRYPOINT
func (el *ContainerFromImage) Shell(shell ...string) (ref *ContainerFromImage) {
	if len(shell) == 0 {
		shell = nil
	}

	el.manager.DockerSys[0].Config.Shell = shell
	return el
}

// Healthcheck
//
// Check the health of the container.
//
//	Input:
//	  interval: time to wait between checks (zero means to inherit);
//	  timeout: time to wait before considering the check to have hung (zero means to inherit);
//	  startPeriod: start period for the container to initialize before the retries starts to count down (zero means to inherit);
//	  retries: number of consecutive failures needed to consider a container as unhealthy (zero means to inherit);
//	  test: test to perform to check that the container is healthy;
//	    * An empty slice means to inherit the default.
//	    * {} : inherit healthcheck
//	    * {"NONE"} : disable healthcheck
//	    * {"CMD", args...} : exec arguments directly
//	    * {"CMD-SHELL", command} : run command with system's default shell
func (el *ContainerFromImage) Healthcheck(interval, timeout, startPeriod time.Duration, retries int, test ...string) (ref *ContainerFromImage) {
	el.manager.DockerSys[0].Config.Healthcheck = &dockerContainer.HealthConfig{
		Test:        test,
		Interval:    interval,
		Timeout:     timeout,
		StartPeriod: startPeriod,
		Retries:     retries,
	}

	//Healthcheck(1*time.Second, 1*time.Second, 1*time.Second, 5, "CMD-SHELL", "mongod --shutdown || exit 1").

	return el
}

// SaveStatistics
//
// Salva um arquivo com as estatísticas de consumo de memória e processamento do container durante os testes
//
//	| read      | pre read  | pids - current (linux) | pids - limit (linux) | num of process (windows) | storage - read count (windows) | storage - write count (windows) | cpu - online | cpu - system usage | cpu - usage in user mode | cpu - usage in kernel mode | cpu - total usage | cpu - throttled time | cpu - throttled periods | cpu - throttling periods | pre cpu - online | pre cpu - system usage | pre cpu - usage in user mode | pre cpu - usage in kernel mode | pre cpu - total usage | pre cpu - throttled time | pre cpu - throttled periods | pre cpu - throttling periods | memory - limit | memory - commit peak | memory - commit | memory - fail cnt | memory - usage | memory - max usage | memory - private working set |
//	|-----------|-----------|------------------------|----------------------|--------------------------|--------------------------------|---------------------------------|--------------|--------------------|--------------------------|----------------------------|-------------------|----------------------|-------------------------|--------------------------|------------------|------------------------|------------------------------|--------------------------------|-----------------------|--------------------------|-----------------------------|------------------------------|----------------|----------------------|-----------------|-------------------|----------------|--------------------|------------------------------|
//	| 270355545 | 267925794 | 36                     | -1                   | 0                        | 0                              | 0                               | 8            | 128396690000000    | 1333036000               | 273231000                  | 1606267000        | 0                    | 0                       | 0                        | 8                | 128388860000000        | 1122134000                   | 188896000                      | 1311030000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 67489792       | 0                  | 0                            |
//	| 315625547 | 312487880 | 36                     | -1                   | 0                        | 0                              | 0                               | 8            | 128443910000000    | 2428358000               | 705437000                  | 3133796000        | 0                    | 0                       | 0                        | 8                | 128436100000000        | 2261894000                   | 623029000                      | 2884924000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 74043392       | 0                  | 0                            |
//	| 331017884 | 328716175 | 37                     | -1                   | 0                        | 0                              | 0                               | 8            | 128490870000000    | 3388019000               | 1217971000                 | 4605991000        | 0                    | 0                       | 0                        | 8                | 128483010000000        | 3235788000                   | 1129258000                     | 4365046000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 79872000       | 0                  | 0                            |
//	| 375934470 | 373538303 | 37                     | -1                   | 0                        | 0                              | 0                               | 8            | 128538150000000    | 4373956000               | 1736955000                 | 6110912000        | 0                    | 0                       | 0                        | 8                | 128530300000000        | 4209072000                   | 1648809000                     | 5857882000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 85491712       | 0                  | 0                            |
//	| 392846000 | 389797833 | 37                     | -1                   | 0                        | 0                              | 0                               | 8            | 128585060000000    | 5392002000               | 2341771000                 | 7733774000        | 0                    | 0                       | 0                        | 8                | 128577290000000        | 5213464000                   | 2236247000                     | 7449711000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 91275264       | 0                  | 0                            |
//	| 438223378 | 435128169 | 36                     | -1                   | 0                        | 0                              | 0                               | 8            | 128632160000000    | 6476036000               | 2913993000                 | 9390029000        | 0                    | 0                       | 0                        | 8                | 128624350000000        | 6290689000                   | 2803815000                     | 9094505000            | 0                        | 0                           | 0                            | 12544057344    | 0                    | 0               | 0                 | 97112064       | 0                  | 0                            |
func (el *ContainerFromImage) SaveStatistics(path string) (ref *ContainerFromImage) {
	var err error
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if err = os.MkdirAll(path, fs.ModePerm); err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container.SaveStatistics().MkdirAll().error: %v", "directory not found")
			return el
		}
	} else if !fileInfo.IsDir() {
		el.manager.ErrorCh <- fmt.Errorf("container.SaveStatistics().error: %v", "directory not found")
		return el
	}

	el.csvPath = path
	return el
}

func (el *ContainerFromImage) FailFlag(path string, flags ...string) (ref *ContainerFromImage) {
	var err error
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(path); err != nil {
		if err = os.MkdirAll(path, fs.ModePerm); err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container.FailFlag().MkdirAll().error: %v", "directory not found")
			return el
		}
	} else if !fileInfo.IsDir() {
		el.manager.ErrorCh <- fmt.Errorf("container.FailFlag().error: %v", "directory not found")
		return el
	}

	el.failPath = path
	el.failFlag = flags

	return el
}

func (el *ContainerFromImage) Start() (ref *ContainerFromImage) {
	var err error
	for i := 0; i != el.copies; i += 1 {
		err = el.manager.DockerSys[i].ContainerStart(el.manager.Id[i])
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Start().ContainerStart().error: %v", i, err)
			return el
		}
	}

	if el.csvPath == "" {
		return
	}

	el.statsThread()
	el.failFlagThread()

	return el
}

func (el *ContainerFromImage) failFlagThread() {
	var err error
	var logs []byte
	var lineList [][]byte
	el.manager.TickerFail = time.NewTicker(10 * time.Second)
	go func() {
		for {
			select {
			case <-el.manager.TickerFail.C:
				for i := 0; i != el.copies; i += 1 {

					logs, err = el.manager.DockerSys[i].ContainerLogs(el.manager.Id[i])
					if err != nil {
						el.manager.ErrorCh <- fmt.Errorf("container[%v].failFlagThread().ContainerLogs().error: %v", i, err)
						return
					}

					lineList = el.logsCleaner(logs, i)
					el.logsSearchAndReplaceIntoText(i, &logs, lineList, el.failPath, el.failFlag)

				}
			}
		}
	}()
}

// logsCleaner
//
// English:
//
//	Clear blank lines of the container's standard output
//
//	Input:
//	  logs: container's standard output
//
//	Output:
//	  logsLine: List of lines of the container's standard output
func (el *ContainerFromImage) logsCleaner(logs []byte, key int) (logsLine [][]byte) {

	size := len(logs) - 1

	// faz o log só lê a parte mais recente do mesmo
	logs = logs[el.failLogsLastSize[key]:]
	el.failLogsLastSize[key] = size

	logs = bytes.ReplaceAll(logs, []byte("\r"), []byte(""))
	return bytes.Split(logs, []byte("\n"))
}

func (el *ContainerFromImage) logsSearchAndReplaceIntoText(key int, logs *[]byte, lineList [][]byte, pathLog string, failFlags []string) (line []byte, found bool) {
	var err error
	var dirList []fs.FileInfo

	for logLine := len(lineList) - 1; logLine >= 0; logLine -= 1 {

		for filterLine := 0; filterLine != len(failFlags); filterLine += 1 {
			line = lineList[logLine]
			if bytes.Contains(line, []byte(failFlags[filterLine])) == true {

				if pathLog != "" {
					dirList, err = ioutil.ReadDir(pathLog)
					if err != nil {
						log.Printf("ioutil.ReadDir().error: %v", err.Error())
						util.TraceToLog()
						return
					}
					var totalOfFiles = strconv.Itoa(len(dirList))
					err = ioutil.WriteFile(filepath.Join(pathLog, el.containerName+"_"+strconv.FormatInt(int64(key), 10)+"."+totalOfFiles+".fail.log"), *logs, fs.ModePerm)
					if err != nil {
						log.Printf("ioutil.WriteFile().error: %v", err.Error())
						util.TraceToLog()
						return
					}
				}

				found = true
				return
			}
		}
	}

	return
}

func (el *ContainerFromImage) statsThread() {
	var err error
	el.manager.TickerStats = time.NewTicker(10 * time.Second)
	go func() {
		var stats types.Stats
		var line [][]string
		var writer *csv.Writer

		var file = make([]*os.File, el.copies)
		for i := 0; i != el.copies; i += 1 {
			file[i], err = os.OpenFile(filepath.Join(el.csvPath, fmt.Sprintf("stats.%v.%v.csv", el.containerName, i)), os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.ModePerm)
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("container[%v].statsThread().OpenFile().error: %v", i, err)
				return
			}
		}

		defer func() {
			for i := 0; i != el.copies; i += 1 {
				_ = file[i].Close()
			}
		}()

		line = [][]string{{
			"read",
			"pre read",

			"pids - current (linux)",
			"pids - limit (linux)",

			"num of process (windows)",

			"storage - read count (windows)",
			"storage - write count (windows)",

			"cpu - online",
			"cpu - system usage",
			"cpu - usage in user mode",
			"cpu - usage in kernel mode",
			"cpu - total usage",
			"cpu - throttled time",
			"cpu - throttled periods",
			"cpu - throttling periods",

			"pre cpu - online",
			"pre cpu - system usage",
			"pre cpu - usage in user mode",
			"pre cpu - usage in kernel mode",
			"pre cpu - total usage",
			"pre cpu - throttled time",
			"pre cpu - throttled periods",
			"pre cpu - throttling periods",

			"memory - limit",
			"memory - commit peak",
			"memory - commit",
			"memory - fail cnt",
			"memory - usage",
			"memory - max usage",
			"memory - private working set",
		}}

		for i := 0; i != el.copies; i += 1 {
			writer = csv.NewWriter(file[i])
			err = writer.WriteAll(line)
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("container[%v].statsThread().WriteAll(0).error: %v", i, err)
				return
			}
		}

		for {
			select {
			case <-el.manager.TickerStats.C:
				for i := 0; i != el.copies; i += 1 {

					stats, err = el.manager.DockerSys[i].ContainerStatisticsOneShot(el.manager.Id[i])
					if err != nil {
						el.manager.ErrorCh <- fmt.Errorf("container[%v].statsThread().ContainerInspect().error: %v", i, err)
						continue
					}

					line = [][]string{{
						strconv.FormatInt(int64(stats.Read.Nanosecond()), 10),
						strconv.FormatInt(int64(stats.PreRead.Nanosecond()), 10),

						//linux
						strconv.FormatInt(int64(stats.PidsStats.Current), 10),
						strconv.FormatInt(int64(stats.PidsStats.Limit), 10),

						//windows
						strconv.FormatInt(int64(stats.NumProcs), 10),
						strconv.FormatInt(int64(stats.StorageStats.ReadCountNormalized), 10),
						strconv.FormatInt(int64(stats.StorageStats.WriteCountNormalized), 10),

						// Shared stats
						strconv.FormatUint(uint64(stats.CPUStats.OnlineCPUs), 10),
						strconv.FormatUint(stats.CPUStats.SystemUsage, 10),
						strconv.FormatUint(stats.CPUStats.CPUUsage.UsageInUsermode, 10),
						strconv.FormatUint(stats.CPUStats.CPUUsage.UsageInKernelmode, 10),
						strconv.FormatUint(stats.CPUStats.CPUUsage.TotalUsage, 10),
						strconv.FormatUint(stats.CPUStats.ThrottlingData.ThrottledTime, 10),
						strconv.FormatUint(stats.CPUStats.ThrottlingData.ThrottledPeriods, 10),
						strconv.FormatUint(stats.CPUStats.ThrottlingData.Periods, 10),

						strconv.FormatUint(uint64(stats.PreCPUStats.OnlineCPUs), 10),
						strconv.FormatUint(stats.PreCPUStats.SystemUsage, 10),
						strconv.FormatUint(stats.PreCPUStats.CPUUsage.UsageInUsermode, 10),
						strconv.FormatUint(stats.PreCPUStats.CPUUsage.UsageInKernelmode, 10),
						strconv.FormatUint(stats.PreCPUStats.CPUUsage.TotalUsage, 10),
						strconv.FormatUint(stats.PreCPUStats.ThrottlingData.ThrottledTime, 10),
						strconv.FormatUint(stats.PreCPUStats.ThrottlingData.ThrottledPeriods, 10),
						strconv.FormatUint(stats.PreCPUStats.ThrottlingData.Periods, 10),

						strconv.FormatUint(stats.MemoryStats.Limit, 10),
						strconv.FormatUint(stats.MemoryStats.CommitPeak, 10),
						strconv.FormatUint(stats.MemoryStats.Commit, 10),
						strconv.FormatUint(stats.MemoryStats.Failcnt, 10),
						strconv.FormatUint(stats.MemoryStats.Usage, 10),
						strconv.FormatUint(stats.MemoryStats.MaxUsage, 10),
						strconv.FormatUint(stats.MemoryStats.PrivateWorkingSet, 10),
					}}

					writer = csv.NewWriter(file[i])
					err = writer.WriteAll(line)
					if err != nil {
						el.manager.ErrorCh <- fmt.Errorf("container[%v].statsThread().WriteAll(1).error: %v", i, err)
						return
					}

				}
			}
		}
	}()
}

func (el *ContainerFromImage) Create(imageName, containerName string, copies int) (ref *ContainerFromImage) {
	var err error

	if copies == 0 {
		return el
	}

	el.failLogsLastSize = make([]int, copies)
	// adjust image name to have version tag
	el.imageName = el.manager.DockerSys[0].AdjustImageName(imageName)
	el.containerName = containerName
	el.copies = copies

	// if the image does not exist, download the image
	if err = el.imagePull(); err != nil {
		el.manager.ErrorCh <- err
		return el
	}

	var ipAddress string
	var netConfig *networkTypes.NetworkingConfig
	el.IPV4Address = make([]string, 0)
	for i := 0; i != copies; i += 1 {

		// index zero is created when the manager object is created, the other indexes are created here, in case there is
		// more than one container to be created
		if i != 0 {
			var dockerSys = new(builder.DockerSystem)
			_ = dockerSys.Init()
			el.manager.DockerSys = append(el.manager.DockerSys, dockerSys)
		}

		// get the next ip address from network
		if el.manager.network != nil {
			ipAddress, netConfig, err = el.manager.network.generator.GetNext()
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("container.network().GetNext().error: %v", err)
				return el
			}
			el.IPV4Address = append(el.IPV4Address, ipAddress)
		}

		// map the port container:host[copiesKey]
		var portConfig = nat.PortMap{}
		for kContainer := range el.portsContainer {
			portBind := make([]nat.PortBinding, 0)
			if len(el.portsHost[kContainer]) > i && el.portsHost[kContainer][i] > 0 {
				portBind = append(portBind, nat.PortBinding{HostPort: strconv.FormatInt(el.portsHost[kContainer][i], 10)})
			}

			portConfig[el.portsContainer[kContainer]] = portBind
		}

		var volumes = make([]mount.Mount, 0)
		for k := range el.volumeContainer {
			volume := mount.Mount{}
			if len(el.volumeContainer[k]) > i && el.volumeHost[k][i] != "" {
				volume.Type = builder.KVolumeMountTypeBindString
				volume.Source = el.volumeHost[k][i]
				volume.Target = el.volumeContainer[k]

				volumes = append(volumes, volume)
			}
		}

		var config = el.manager.DockerSys[0].GetConfig()
		config.Image = imageName

		// create the container, link container and network, but, don't start the container
		var warnings []string
		var id string
		id, warnings, err = el.manager.DockerSys[i].ContainerCreateWithConfig(
			config,
			containerName+"_"+strconv.FormatInt(int64(i), 10),
			builder.KRestartPolicyNo,
			portConfig,
			volumes,
			netConfig,
		)
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerCreate().error: %v", i, err)
			return el
		}

		// id de todos os containers criados para a função start()
		el.manager.Id = append(el.manager.Id, id)

		//todo: fazer warnings - não deve ser erro
		if len(warnings) != 0 {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerCreate().warnings: %v", i, strings.Join(warnings, "; "))
			return el
		}
	}

	return el
}

// imagePull
//
// If the image exists on the local computer, it does nothing, otherwise it tries to download the image
func (el *ContainerFromImage) imagePull() (err error) {
	el.imageId, _ = el.manager.DockerSys[0].ImageFindIdByName(el.imageName)
	if el.imageId != "" {
		return
	}

	// English: make a channel to end goroutine
	// Português: monta um canal para terminar a goroutine
	var chProcessEnd = make(chan bool, 1)

	// English: make a channel [optional] to print build output
	// Português: monta o canal [opcional] para imprimir a saída do build
	var chStatus = make(chan builder.ContainerPullStatusSendToChannel, 1)

	// English: make a thread to monitoring and print channel data
	// Português: monta uma thread para imprimir os dados do canal
	go func(chStatus chan builder.ContainerPullStatusSendToChannel, chProcessEnd chan bool) {

		for {
			select {
			case <-chProcessEnd:
				// English: Eliminate this goroutine after process end
				// Português: Elimina a goroutine após o fim do processo
				return

			case status := <-chStatus:
				// English: remove this comment to see all build status
				// Português: remova este comentário para vê _todo o status da criação da imagem
				fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					fmt.Println("image pull complete!")
				}
			}
		}

	}(chStatus, chProcessEnd)

	defer func() {
		// English: ends a goroutine
		// Português: termina a goroutine
		chProcessEnd <- true
	}()

	// docker pull
	el.imageId, el.imageName, err = el.manager.DockerSys[0].ImagePull(el.imageName, &chStatus)
	if err != nil {
		err = fmt.Errorf("containerFromImage.Primordial().imagePull().error: %v", err)
		return
	}

	return
}
