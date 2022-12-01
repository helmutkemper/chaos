package manager

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	sshGit "github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/helmutkemper/chaos/internal/builder"
	"github.com/helmutkemper/chaos/internal/dockerfileGolang"
	"github.com/helmutkemper/chaos/internal/util/utilCopy"
	"io/fs"
	"io/ioutil"
	"log"
	"math"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	kDst = 0
	kSrc = 1
)

type containerCommon struct {
	IPV4Address []string

	//port inside container and host computer port
	portsContainer []nat.Port
	portsHost      [][]int64

	volumeContainer []string
	volumeHost      [][]string

	manager *Manager

	imageExpirationTime time.Duration
	buildPath           string
	serverPath          string
	replaceBeforeBuild  [][]string
	command             string
	imageId             string
	imageName           string
	containerName       string
	copies              int
	csvPath             string
	failPath            string
	failFlag            []string
	failLogsLastSize    []int

	environment [][]string
	testTimeout time.Duration

	makeDefaultDockerfile       bool
	makeDefaultDockerfileExtras bool

	enableCache    bool
	imageCacheName string
	autoDockerfile DockerfileAuto

	contentGitConfigFile           string
	contentKnownHostsFile          string
	contentIdRsaFile               string
	contentIdEcdsaFile             string
	gitPathPrivateRepository       string
	sshDefaultFileName             string
	contentIdRsaFileWithScape      string
	contentKnownHostsFileWithScape string
	contentGitConfigFileWithScape  string

	gitUrl               string
	gitPrivateToke       string
	gitUser              string
	gitPassword          string
	gitSshPrivateKeyPath string
}

type ContainerFromImage struct {
	containerCommon
}

//func (el *ContainerFromImage) New(manager *Manager) {
//	el.manager = manager
//}

// MakeDockerfile
//
// Mounts a standard Dockerfile automatically
func (el *ContainerFromImage) MakeDockerfile() (ref *ContainerFromImage) {
	el.makeDefaultDockerfile = true
	el.makeDefaultDockerfileExtras = true

	return el
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
func (el *ContainerFromImage) EnvironmentVar(env ...[]string) (ref *ContainerFromImage) {
	if len(env) == 0 {
		env = nil
		return el
	}

	el.environment = env
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
// Timeout to stop the container after command `container.Stop()`
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

// ReplaceBeforeBuild
//
// Replaces or adds files to the project, in the temporary folder, before the image is created.
func (el *ContainerFromImage) ReplaceBeforeBuild(dst, src string) (ref *ContainerFromImage) {
	var err error
	if el.replaceBeforeBuild == nil {
		el.replaceBeforeBuild = make([][]string, 0)
	}

	src, err = filepath.Abs(src)
	if err != nil {
		el.manager.ErrorCh <- fmt.Errorf("container.ReplaceBeforeBuild().error: %v", err)
		return el
	}

	el.replaceBeforeBuild = append(el.replaceBeforeBuild, []string{dst, src})

	return el
}

// TestDuration
//
// Defines the duration of the test
func (el *ContainerFromImage) TestDuration(timeout time.Duration) (ref *ContainerFromImage) {
	el.testTimeout = timeout //todo: fazer
	return el
}

// FailFlag
//
// Define um texto, que quando encontrado na saída padrão do container, define o teste como falho.
//
//	Input:
//	  path: path to save the container standard output
//	  flags: texts to be searched for in the container standard output
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

// Start
//
// Start the container after build
func (el *ContainerFromImage) Start() (ref *ContainerFromImage) {
	var err error

	for i := 0; i != el.copies; i += 1 {
		err = el.manager.DockerSys[i].ContainerStart(el.manager.Id[i])
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Start().ContainerStart().error: %v", i, err)
			return el
		}
	}

	var inspect types.ContainerJSON
	for i := 0; i != el.copies; i += 1 {
		inspect, err = el.manager.DockerSys[i].ContainerInspect(el.manager.Id[i])
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Start().ContainerInspect().error: %v", i, err)
			return el
		}

		if inspect.State == nil || inspect.State.Running == false {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Start().error: %v", i, "container is't running")
			return el
		}
	}

	el.failFlagThread()

	if el.csvPath == "" {
		return
	}

	el.statsThread()

	//if el.testTimeout != 0 {
	//  var timeout = time.NewTimer(el.testTimeout)
	//  go func() {
	//    select {
	//    case <-timeout.C:
	//      el.manager.Done <- struct{}{}
	//    }
	//  }()
	//}
	//
	//<-el.manager.Done

	return el
}

// failFlagThread
//
// ticker that monitors the standard output of the container looking for test failure flags
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
						el.manager.ErrorCh <- fmt.Errorf("container.logsSearchAndReplaceIntoText().ioutil.ReadDir(%v).error: %v", pathLog, err)
						return
					}
					var totalOfFiles = strconv.Itoa(len(dirList))
					var path = filepath.Join(pathLog, el.containerName+"_"+strconv.FormatInt(int64(key), 10)+"."+totalOfFiles+".fail.log")
					err = ioutil.WriteFile(path, *logs, fs.ModePerm)
					if err != nil {
						el.manager.ErrorCh <- fmt.Errorf("container.logsSearchAndReplaceIntoText().ioutil.WriteFile(%v).error: %v", path, err)
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

// statsThread
//
// Inspects the container and saves container statistics information to a CSV file every 10 seconds
func (el *ContainerFromImage) statsThread() {
	var err error
	el.manager.TickerStats = time.NewTicker(10 * time.Second)
	go func() {
		var stats types.Stats
		var line [][]string
		var writer *csv.Writer

		var file = make([]*os.File, el.copies)
		for i := 0; i != el.copies; i += 1 {
			var filePath = filepath.Join(el.csvPath, fmt.Sprintf("stats.%v.%v.csv", el.containerName, i))
			_ = os.Remove(filePath)
			file[i], err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, fs.ModePerm)
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
			"time",

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
						time.Now().Format("2006-01-02 15:04:05"),

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

// Create
//
// Cria o container.
//
// Before this function is called, all settings functions must have been defined.
//
//	Input:
//	  containerName: name from container
//	  copies: number total of containers
func (el *ContainerFromImage) Create(containerName string, copies int) (ref *ContainerFromImage) {
	var err error

	if copies == 0 {
		return el
	}

	if el.autoDockerfile == nil {
		el.autoDockerfile = new(dockerfileGolang.DockerfileGolang)
	}

	if !strings.Contains(containerName, "delete") {
		containerName = "delete_" + containerName
	}

	el.failLogsLastSize = make([]int, copies)
	// adjust image name to have version tag
	el.imageName = el.manager.DockerSys[0].AdjustImageName(el.imageName)
	el.containerName = containerName
	el.copies = copies

	var config = el.manager.DockerSys[0].GetConfig()

	err = el.imageBuild()
	if err != nil {
		el.manager.ErrorCh <- fmt.Errorf("container.Create().imageBuild().error: %v", err)
		return el
	}

	var ipAddress string
	var netConfig *networkTypes.NetworkingConfig
	el.IPV4Address = make([]string, 0)
	for iCopy := 0; iCopy != copies; iCopy += 1 {

		// index zero is created when the manager object is created, the other indexes are created here, in case there is
		// more than one container to be created
		if iCopy != 0 {
			var dockerSys = new(builder.DockerSystem)
			_ = dockerSys.Init()
			el.manager.DockerSys = append(el.manager.DockerSys, dockerSys)
		}

		// get the next ip address from network
		if el.manager.network != nil {
			ipAddress, netConfig, err = el.manager.network.generator.GetNext()
			if err != nil {
				el.manager.ErrorCh <- fmt.Errorf("container.Create().network.GetNext().error: %v", err)
				return el
			}
			el.IPV4Address = append(el.IPV4Address, ipAddress)
		}

		// map the port container:host[copiesKey]
		var portConfig = el.mapContainerPorts(iCopy)
		var volumes = el.mapVolumes(iCopy)

		config.Image = el.imageName

		// todo: documentar isto
		if len(el.environment) > iCopy {
			config.Env = el.environment[iCopy]
		} else if len(el.environment) == 1 {
			config.Env = el.environment[0]
		} else {
			config.Env = nil
		}

		// create the container, link container and network, but, don't start the container
		var warnings []string
		var id string
		id, warnings, err = el.manager.DockerSys[iCopy].ContainerCreateWithConfig(
			config,
			containerName+"_"+strconv.FormatInt(int64(iCopy), 10),
			builder.KRestartPolicyNo,
			portConfig,
			volumes,
			netConfig,
		)
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Create().ContainerCreateWithConfig().error: %v", iCopy, err)
			return el
		}

		// id de todos os containers criados para a função start()
		el.manager.Id = append(el.manager.Id, id)

		//todo: fazer warnings - não deve ser erro
		if len(warnings) != 0 {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].Create().ContainerCreateWithConfig().warnings: %v", iCopy, strings.Join(warnings, "; "))
			return el
		}
	}

	return el
}

// mapVolumes
//
// Mount the container volumes
func (el *ContainerFromImage) mapVolumes(iCopy int) (volumes []mount.Mount) {
	volumes = make([]mount.Mount, 0)

	for k := range el.volumeContainer {
		volume := mount.Mount{}
		if len(el.volumeContainer[k]) > iCopy && el.volumeHost[k][iCopy] != "" {
			volume.Type = builder.KVolumeMountTypeBindString
			volume.Source = el.volumeHost[k][iCopy]
			volume.Target = el.volumeContainer[k]

			volumes = append(volumes, volume)
		}
	}

	return
}

// mapContainerPorts
//
// Maps container ports
func (el *ContainerFromImage) mapContainerPorts(iCopy int) (portConfig nat.PortMap) {
	portConfig = make(map[nat.Port][]nat.PortBinding)

	// map the port container:host[copiesKey]
	for kContainer := range el.portsContainer {
		portBind := make([]nat.PortBinding, 0)
		if len(el.portsHost[kContainer]) > iCopy && el.portsHost[kContainer][iCopy] > 0 {
			portBind = append(portBind, nat.PortBinding{HostPort: strconv.FormatInt(el.portsHost[kContainer][iCopy], 10)})
		}

		portConfig[el.portsContainer[kContainer]] = portBind
	}

	return
}

// imageBuild
//
// project image build
func (el *ContainerFromImage) imageBuild() (err error) {
	switch el.command {
	case "fromServer":
		if el.serverPath == "" {
			err = fmt.Errorf("set server path first")
			return
		}

		if el.checkImageExpirationTimeIsValid() {
			return
		}

		var tmpDir string
		var publicKeys *sshGit.PublicKeys
		var gitCloneConfig *git.CloneOptions
		publicKeys, err = el.gitMakePublicSshKey()
		if err != nil {
			err = fmt.Errorf("container.imageBuild().gitMakePublicSshKey().error: %v", err)
			return
		}

		tmpDir, err = el.makeTmpDir()
		if err != nil {
			err = fmt.Errorf("container.imageBuild().makeTmpDir().error: %v", err)
			return
		}
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		if el.gitSshPrivateKeyPath != "" || el.contentIdRsaFile != "" {
			gitCloneConfig = &git.CloneOptions{
				URL:      el.gitUrl,
				Auth:     publicKeys,
				Progress: nil,
			}
		} else if el.gitPrivateToke != "" {
			gitCloneConfig = &git.CloneOptions{
				// The intended use of a GitHub personal access token is in replace of your password
				// because access tokens can easily be revoked.
				// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
				Auth: &http.BasicAuth{
					Username: "abc123", // yes, this can be anything except an empty string
					Password: el.gitPrivateToke,
				},
				URL:      el.gitUrl,
				Progress: nil,
			}
		} else if el.gitUser != "" && el.gitPassword != "" {
			gitCloneConfig = &git.CloneOptions{
				Auth: &http.BasicAuth{
					Username: el.gitUser,
					Password: el.gitPassword,
				},
				URL:      el.gitUrl,
				Progress: nil,
			}
		} else {
			gitCloneConfig = &git.CloneOptions{
				URL:      el.gitUrl,
				Progress: nil,
			}
		}

		_, err = git.PlainClone(tmpDir, false, gitCloneConfig)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().PlainClone().error: %v", err)
			return
		}

		err = el.replaceFilesBeforeBuild(tmpDir)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().replaceFilesBeforeBuild().error: %v", err)
			return
		}

		el.buildPath = tmpDir

		var volumes = make([]mount.Mount, 0)
		err = el.makeDefaultDockerfileForMe(volumes)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().makeDefaultDockerfileForMe().error: %v", err)
			return
		}

		el.autoDockerfile.Prayer()

		var changePointer = make(chan builder.ContainerPullStatusSendToChannel)
		go el.imageBuildStdOutputToLogOutput(changePointer)

		el.manager.DockerSys[0].ImageBuildFromRemoteServer()

	case "fromFolder":
		if el.buildPath == "" {
			err = fmt.Errorf("set build folder path first")
			return
		}

		if el.checkImageExpirationTimeIsValid() {
			return
		}

		var tmpDir string
		tmpDir, err = el.copyBuildPathToTmpDir()
		if err != nil {
			err = fmt.Errorf("container.imageBuild().copyBuildPathToTmpDir().error: %v", err)
			return
		}
		defer func() {
			_ = os.RemoveAll(tmpDir)
		}()

		err = el.replaceFilesBeforeBuild(tmpDir)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().replaceFilesBeforeBuild().error: %v", err)
			return
		}

		el.buildPath = tmpDir

		var volumes = make([]mount.Mount, 0)
		err = el.makeDefaultDockerfileForMe(volumes)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().makeDefaultDockerfileForMe().error: %v", err)
			return
		}

		el.autoDockerfile.Prayer()

		var changePointer = make(chan builder.ContainerPullStatusSendToChannel)
		go el.imageBuildStdOutputToLogOutput(changePointer)

		el.imageId, err = el.manager.DockerSys[0].ImageBuildFromFolder(
			el.buildPath,
			el.imageName,
			[]string{},
			el.manager.ImageBuildOptions,
			changePointer,
		)
		if err != nil {
			err = fmt.Errorf("container.imageBuild().ImageBuildFromFolder().error: %v", err)
			return
		}

		if el.imageId == "" {
			err = fmt.Errorf("container.imageBuild().ImageBuildFromFolder().error: %v", "image ID was not generated")
			return
		}

		// Construir uma imagem de múltiplas etapas deixa imagens grandes e sem serventia, ocupando espaço no HD.
		_ = el.manager.DockerSys[0].ImageGarbageCollector()

	case "fromImage":
		// if the image does not exist, download the image
		if err = el.imagePull(); err != nil {
			err = fmt.Errorf("container.imageBuild().imagePull().error: %v", err)
			return
		}
	}

	return
}

// imageBuildStdOutputToLogOutput
//
// Turns the container's standard output into a log during the image creation or download process
func (el *ContainerFromImage) imageBuildStdOutputToLogOutput(ch chan builder.ContainerPullStatusSendToChannel) {

	for {
		select {
		case event := <-ch:
			var stream = event.Stream
			stream = strings.ReplaceAll(stream, "\n", "")
			stream = strings.ReplaceAll(stream, "\r", "")
			stream = strings.Trim(stream, " ")

			if stream == "" {
				continue
			}

			log.Printf("%v", stream)

			if event.Closed == true {
				return
			}
		}
	}
}

// makeDefaultDockerfileForMe
//
// When enabled by the user, it mounts the dockerfile automatically. Requires a go.mod file in the project root
func (el *ContainerFromImage) makeDefaultDockerfileForMe(volumes []mount.Mount) (err error) {
	if !el.makeDefaultDockerfile {
		return
	}

	var dockerfile string
	var fileList []fs.FileInfo

	fileList, err = ioutil.ReadDir(el.buildPath)
	if err != nil {
		err = fmt.Errorf("container.makeDefaultDockerfileForMe().ioutil.ReadDir().error: %v", err)
		return
	}

	// fixme: modificar isto
	// deve ir para a interface{} fazer a verificação
	var pass = false
	for _, file := range fileList {
		if file.Name() == "go.mod" {
			pass = true
			break
		}
	}
	if pass == false {
		err = errors.New("go.mod file not found")
		return
	}

	if el.enableCache == true && el.manager.ImageBuildOptions.NoCache != true {
		_, err = el.manager.DockerSys[0].ImageFindIdByName(el.imageCacheName)
		if err != nil && err.Error() == "image name not found" { //todo: isto deveria ser um var inf = errors.New("image name not found")
			err = nil
			el.enableCache = false
		}
		if err != nil {
			err = fmt.Errorf("container.makeDefaultDockerfileForMe().ImageFindIdByName().error: %v", err)
			return
		}
	}

	dockerfile, err = el.autoDockerfile.MountDefaultDockerfile(
		el.manager.ImageBuildOptions.BuildArgs,
		el.portsContainer,
		volumes,
		el.makeDefaultDockerfileExtras,
		el.enableCache,
		el.imageCacheName,
	)
	if err != nil {
		err = fmt.Errorf("container.makeDefaultDockerfileForMe().autoDockerfile.MountDefaultDockerfile().error: %v", err)
		return
	}

	var dockerfilePath = filepath.Join(el.buildPath, "Dockerfile-iotmaker")
	err = ioutil.WriteFile(dockerfilePath, []byte(dockerfile), os.ModePerm)
	if err != nil {
		err = fmt.Errorf("container.makeDefaultDockerfileForMe().ioutil.WriteFile().error: %v", err)
		return
	}

	return
}

// replaceFilesBeforeBuild
//
// Copies files and folders defined before testing into the project folder
func (el *ContainerFromImage) replaceFilesBeforeBuild(tmpDir string) (err error) {
	var fileInfo os.FileInfo
	for k := range el.replaceBeforeBuild {

		fileInfo, err = os.Stat(el.replaceBeforeBuild[k][kSrc])
		if err != nil {
			err = fmt.Errorf("container.replaceFilesBeforeBuild().Stat().error: %v", err)
			return
		}

		if fileInfo.IsDir() {
			err = utilCopy.Dir(filepath.Join(tmpDir, el.replaceBeforeBuild[k][kDst]), el.replaceBeforeBuild[k][kSrc])
			if err != nil {
				err = fmt.Errorf("container.replaceFilesBeforeBuild().utilCopy.Dir(1).error: %v", err)
				return
			}
		} else {
			err = utilCopy.File(filepath.Join(tmpDir, el.replaceBeforeBuild[k][kDst]), el.replaceBeforeBuild[k][kSrc])
			if err != nil {
				err = fmt.Errorf("container.replaceFilesBeforeBuild().utilCopy.File(0).error: %v", err)
				return
			}
		}
	}

	return
}

// checkImageExpirationTimeIsValid
//
// Checks the image validity time to not recreate the same image in the same test
func (el *ContainerFromImage) checkImageExpirationTimeIsValid() (isValid bool) {
	el.imageId, _ = el.manager.DockerSys[0].ImageFindIdByName(el.imageName)
	return el.imageId != "" && el.imageExpirationTimeIsValid() == true
}

// makeTmpDir
//
// make a tmp dir
func (el *ContainerFromImage) makeTmpDir() (tmpDir string, err error) {
	tmpDir, err = os.MkdirTemp("", "chaos__")
	if err != nil {
		err = fmt.Errorf("container.makeTmpDir().error: %v", err)
		return
	}

	return
}

// copyBuildPathToTmpDir
//
// Create a temporary directory and copy the project to it, before making the image.
// This allows changing project files without damaging the original project.
func (el *ContainerFromImage) copyBuildPathToTmpDir() (tmpDir string, err error) {
	tmpDir, err = el.makeTmpDir()
	if err != nil {
		err = fmt.Errorf("container.copyBuildPathToTmpDir().makeTmpDir().error: %v", err)
		return
	}

	el.buildPath, err = filepath.Abs(el.buildPath)
	if err != nil {
		err = fmt.Errorf("container.copyBuildPathToTmpDir().Abs().error: %v", err)
		return
	}

	err = utilCopy.Dir(tmpDir, el.buildPath)
	if err != nil {
		err = fmt.Errorf("container.copyBuildPathToTmpDir().Dir().error: %v", err)
		return
	}

	return
}

// imageExpirationTimeIsValid
//
// English:
//
//	Detects if the image is within the expiration date.
//
//	 Output:
//	   valid: true, if the image is within the expiry date.
//
// Português:
//
//	Detecta se a imagem está dentro do prazo de validade.
//
//	 Saída:
//	   valid: true, se a imagem está dentro do prazo de validade.
func (el *ContainerFromImage) imageExpirationTimeIsValid() (valid bool) {
	if el.imageExpirationTime == 0 {
		return
	}

	var err error
	var inspect types.ImageInspect
	inspect, err = el.manager.DockerSys[0].ImageInspect(el.imageId)
	if err != nil {
		return
	}

	var imageCreated time.Time
	imageCreated, err = time.Parse(time.RFC3339Nano, inspect.Created)
	if err != nil {
		el.manager.ErrorCh <- fmt.Errorf("container.imageExpirationTimeIsValid().error: %v", err)
		return
	}
	return imageCreated.Add(el.imageExpirationTime).After(time.Now())
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
	el.imageId, el.imageName, err = el.manager.DockerSys[0].ImagePull(el.imageName, chStatus)
	if err != nil {
		err = fmt.Errorf("containerFromImage.Primordial().imagePull().error: %v", err)
		return
	}

	return
}

// DockerfileAuto
//
// English: Interface from automatic Dockerfile generator.
//
//	Note: To be able to access private repositories from inside the container, build the image in two or more
//	steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh folder and the .gitconfig
//	file to the /root folder.
//
//	One way to do this automatically is to use the Dockerfile example below, where the arguments SSH_ID_RSA_FILE
//	contains the file ~/.ssh/id_rsa, KNOWN_HOSTS_FILE contains the file ~/.ssh/known_hosts and GITCONFIG_FILE
//	contains the file ~/.gitconfig.
//
//	If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the password.
//
//	If you want to copy the files into the image automatically, use SetPrivateRepositoryAutoConfig() and the
//	function will copy the files ~/.ssh/id_rsa, ~/.ssh/known_hosts and ~/.gitconfig to the viable arguments
//	located above.
//
//	If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and SetSshIdRsaFile()
//	to define the files manually.
//
//	The Dockerfile below can be used as a base
//
//	  # (en) first stage of the process
//	  # (pt) primeira etapa do processo
//	  FROM golang:1.16-alpine as builder
//
//	  # (en) enable the argument variables
//	  # (pt) habilita as variáveis de argumento
//	  ARG SSH_ID_RSA_FILE
//	  ARG KNOWN_HOSTS_FILE
//	  ARG GITCONFIG_FILE
//	  ARG GIT_PRIVATE_REPO
//
//	  # (en) creates the .ssh directory within the root directory
//	  # (pt) cria o diretório .ssh dentro do diretório root
//	  RUN mkdir -p /root/.ssh/ && \
//	      # (en) creates the id_esa file inside the .ssh directory
//	      # (pt) cria o arquivo id_esa dentro do diretório .ssh
//	      echo "$SSH_ID_RSA_FILE" > /root/.ssh/id_rsa && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/ && \
//	      # (en) creates the known_hosts file inside the .ssh directory
//	      # (pt) cria o arquivo known_hosts dentro do diretório .ssh
//	      echo "$KNOWN_HOSTS_FILE" > /root/.ssh/known_hosts && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/known_hosts && \
//	      # (en) creates the .gitconfig file at the root of the root directory
//	      # (pt) cria o arquivo .gitconfig na raiz do diretório /root
//	      echo "$GITCONFIG_FILE" > /root/.gitconfig && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.gitconfig && \
//	      # (en) prepares the OS for installation
//	      # (pt) prepara o OS para instalação
//	      apk update && \
//	      # (en) install git and openssh
//	      # (pt) instala o git e o opnssh
//	      apk add --no-cache build-base git openssh && \
//	      # (en) clear the cache
//	      # (pt) limpa a cache
//	      rm -rf /var/cache/apk/*
//
//	  # (en) creates the /app directory, where your code will be installed
//	  # (pt) cria o diretório /app, onde seu código vai ser instalado
//	  WORKDIR /app
//	  # (en) copy your project into the /app folder
//	  # (pt) copia seu projeto para dentro da pasta /app
//	  COPY . .
//	  # (en) enables the golang compiler to run on an extremely simple OS, scratch
//	  # (pt) habilita o compilador do golang para rodar em um OS extremamente simples, o scratch
//	  ARG CGO_ENABLED=0
//	  # (en) adjust git to work with shh
//	  # (pt) ajusta o git para funcionar com shh
//	  RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
//	  # (en) defines the path of the private repository
//	  # (pt) define o caminho do repositório privado
//	  RUN echo "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//	  # (en) install the dependencies in the go.mod file
//	  # (pt) instala as dependências no arquivo go.mod
//	  RUN go mod tidy
//	  # (en) compiles the main.go file
//	  # (pt) compila o arquivo main.go
//	  RUN go build -ldflags="-w -s" -o /app/main /app/main.go
//	  # (en) creates a new scratch-based image
//	  # (pt) cria uma nova imagem baseada no scratch
//	  # (en) scratch is an extremely simple OS capable of generating very small images
//	  # (pt) o scratch é um OS extremamente simples capaz de gerar imagens muito reduzidas
//	  # (en) discarding the previous image erases git access credentials for your security and reduces the size of the
//	  #      image to save server space
//	  # (pt) descartar a imagem anterior apaga as credenciais de acesso ao git para a sua segurança e reduz o tamanho
//	  #      da imagem para poupar espaço no servidor
//	  FROM scratch
//	  # (en) copy your project to the new image
//	  # (pt) copia o seu projeto para a nova imagem
//	  COPY --from=builder /app/main .
//	  # (en) execute your project
//	  # (pt) executa o seu projeto
//	  CMD ["/main"]
//
// Português: Interface do gerador de dockerfile automático.
//
//	Nota: Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas ou mais
//	etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta /root/.ssh e o arquivo
//	.gitconfig para a pasta /root/.
//
//	Uma maneira de fazer isto de forma automática é usar o exemplo de Dockerfile abaixo, onde os argumentos
//	SSH_ID_RSA_FILE contém o arquivo ~/.ssh/id_rsa, KNOWN_HOSTS_FILE contém o arquivo ~/.ssh/known_hosts e
//	GITCONFIG_FILE contém o arquivo ~/.gitconfig.
//
//	Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para definir a senha da
//	mesma.
//
//	Caso você queira copiar os arquivos para dentro da imagem de forma automática, use
//	SetPrivateRepositoryAutoConfig() e a função copiará os arquivos ~/.ssh/id_rsa, ~/.ssh/known_hosts e
//	~/.gitconfig para as viáveis de argumentos sitada anteriormente.
//
//	Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile() e
//	SetSshIdRsaFile() para definir os arquivos de forma manual.
//
//	O arquivo Dockerfile abaixo pode ser usado como base
//
//	  # (en) first stage of the process
//	  # (pt) primeira etapa do processo
//	  FROM golang:1.16-alpine as builder
//
//	  # (en) enable the argument variables
//	  # (pt) habilita as variáveis de argumento
//	  ARG SSH_ID_RSA_FILE
//	  ARG KNOWN_HOSTS_FILE
//	  ARG GITCONFIG_FILE
//	  ARG GIT_PRIVATE_REPO
//
//	  # (en) creates the .ssh directory within the root directory
//	  # (pt) cria o diretório .ssh dentro do diretório root
//	  RUN mkdir -p /root/.ssh/ && \
//	      # (en) creates the id_esa file inside the .ssh directory
//	      # (pt) cria o arquivo id_esa dentro do diretório .ssh
//	      echo "$SSH_ID_RSA_FILE" > /root/.ssh/id_rsa && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/ && \
//	      # (en) creates the known_hosts file inside the .ssh directory
//	      # (pt) cria o arquivo known_hosts dentro do diretório .ssh
//	      echo "$KNOWN_HOSTS_FILE" > /root/.ssh/known_hosts && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/known_hosts && \
//	      # (en) creates the .gitconfig file at the root of the root directory
//	      # (pt) cria o arquivo .gitconfig na raiz do diretório /root
//	      echo "$GITCONFIG_FILE" > /root/.gitconfig && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.gitconfig && \
//	      # (en) prepares the OS for installation
//	      # (pt) prepara o OS para instalação
//	      apk update && \
//	      # (en) install git and openssh
//	      # (pt) instala o git e o opnssh
//	      apk add --no-cache build-base git openssh && \
//	      # (en) clear the cache
//	      # (pt) limpa a cache
//	      rm -rf /var/cache/apk/*
//
//	  # (en) creates the /app directory, where your code will be installed
//	  # (pt) cria o diretório /app, onde seu código vai ser instalado
//	  WORKDIR /app
//	  # (en) copy your project into the /app folder
//	  # (pt) copia seu projeto para dentro da pasta /app
//	  COPY . .
//	  # (en) enables the golang compiler to run on an extremely simple OS, scratch
//	  # (pt) habilita o compilador do golang para rodar em um OS extremamente simples, o scratch
//	  ARG CGO_ENABLED=0
//	  # (en) adjust git to work with shh
//	  # (pt) ajusta o git para funcionar com shh
//	  RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
//	  # (en) defines the path of the private repository
//	  # (pt) define o caminho do repositório privado
//	  RUN echo "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//	  # (en) install the dependencies in the go.mod file
//	  # (pt) instala as dependências no arquivo go.mod
//	  RUN go mod tidy
//	  # (en) compiles the main.go file
//	  # (pt) compila o arquivo main.go
//	  RUN go build -ldflags="-w -s" -o /app/main /app/main.go
//	  # (en) creates a new scratch-based image
//	  # (pt) cria uma nova imagem baseada no scratch
//	  # (en) scratch is an extremely simple OS capable of generating very small images
//	  # (pt) o scratch é um OS extremamente simples capaz de gerar imagens muito reduzidas
//	  # (en) discarding the previous image erases git access credentials for your security and reduces the size of the
//	  #      image to save server space
//	  # (pt) descartar a imagem anterior apaga as credenciais de acesso ao git para a sua segurança e reduz o tamanho
//	  #      da imagem para poupar espaço no servidor
//	  FROM scratch
//	  # (en) copy your project to the new image
//	  # (pt) copia o seu projeto para a nova imagem
//	  COPY --from=builder /app/main .
//	  # (en) execute your project
//	  # (pt) executa o seu projeto
//	  CMD ["/main"]
type DockerfileAuto interface {
	MountDefaultDockerfile(args map[string]*string, ports []nat.Port, volumes []mount.Mount, installExtraPackages bool, useCache bool, imageCacheName string) (dockerfile string, err error)
	Prayer()
	SetFinalImageName(name string)
	AddCopyToFinalImage(src, dst string)
	SetDefaultSshFileName(name string)
}

// SetImageBuildOptionsSecurityOpt
//
// English:
//
//	Set the container security options
//
//	 Input:
//	   values: container security options
//
// Examples:
//
//	label=user:USER        — Set the label user for the container
//	label=role:ROLE        — Set the label role for the container
//	label=type:TYPE        — Set the label type for the container
//	label=level:LEVEL      — Set the label level for the container
//	label=disable          — Turn off label confinement for the container
//	apparmor=PROFILE       — Set the apparmor profile to be applied to the container
//	no-new-privileges:true — Disable container processes from gaining new privileges
//	seccomp=unconfined     — Turn off seccomp confinement for the container
//	seccomp=profile.json   — White-listed syscalls seccomp Json file to be used as a seccomp filter
//
// Português:
//
//	Modifica as opções de segurança do container
//
//	 Entrada:
//	   values: opções de segurança do container
//
// Exemplos:
//
//	label=user:USER        — Determina o rótulo user para o container
//	label=role:ROLE        — Determina o rótulo role para o container
//	label=type:TYPE        — Determina o rótulo type para o container
//	label=level:LEVEL      — Determina o rótulo level para o container
//	label=disable          — Desliga o confinamento do rótulo para o container
//	apparmor=PROFILE       — Habilita o perfil definido pelo apparmor do linux para ser definido ao container
//	no-new-privileges:true — Impede o processo do container a ganhar novos privilégios
//	seccomp=unconfined     — Desliga o confinamento causado pelo seccomp do linux ao container
//	seccomp=profile.json   — White-listed syscalls seccomp Json file to be used as a seccomp filter
func (el *ContainerFromImage) SetImageBuildOptionsSecurityOpt(value []string) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.SecurityOpt = value
	return el
}

// AddImageBuildOptionsBuildArgs
//
// English:
//
//	Set build-time variables (--build-arg)
//
//	 Input:
//	   key: Argument name
//	   value: Argument value
//
// Example:
//
//	key:   argument key (e.g. Dockerfile: ARG key)
//	value: argument value
//
//	https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg
//	docker build --build-arg HTTP_PROXY=http://10.20.30.2:1234
//
//	  code:
//	    var key = "GIT_PRIVATE_REPO"
//	    var value = "github.com/yourgit"
//
//	    var container = ContainerBuilder{}
//	    container.AddImageBuildOptionsBuildArgs(key, &value)
//
//	  Dockerfile:
//	    FROM golang:1.16-alpine as builder
//	    ARG GIT_PRIVATE_REPO
//	    RUN go env -w GOPRIVATE=$GIT_PRIVATE_REPO
//
// Português:
//
//	Adiciona uma variável durante a construção (--build-arg)
//
//	 Input:
//	   key: Nome do argumento.
//	   value: Valor do argumento.
//
// Exemplo:
//
//	key:   chave do argumento (ex. Dockerfile: ARG key)
//	value: valor do argumento
//
//	https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg
//	docker build --build-arg HTTP_PROXY=http://10.20.30.2:1234
//
//	  code:
//	    var key = "GIT_PRIVATE_REPO"
//	    var value = "github.com/yourgit"
//
//	    var container = ContainerBuilder{}
//	    container.AddImageBuildOptionsBuildArgs(key, &value)
//
//	  Dockerfile:
//	    FROM golang:1.16-alpine as builder
//	    ARG GIT_PRIVATE_REPO
//	    RUN go env -w GOPRIVATE=$GIT_PRIVATE_REPO
func (el *ContainerFromImage) AddImageBuildOptionsBuildArgs(key string, value *string) (ref *ContainerFromImage) {
	if el.manager.ImageBuildOptions.BuildArgs == nil {
		el.manager.ImageBuildOptions.BuildArgs = make(map[string]*string)
	}

	el.manager.ImageBuildOptions.BuildArgs[key] = value
	return el
}

// addImageBuildOptionsGitCredentials
//
// English:
//
//	Prepare the git credentials.
//
//	Called from SetPrivateRepositoryAutoConfig()
//
// Português:
//
//	Prepara as credenciais do git.
//
//	Chamada por SetPrivateRepositoryAutoConfig()
func (el *ContainerFromImage) addImageBuildOptionsGitCredentials() {

	if el.manager.ImageBuildOptions.BuildArgs == nil {
		el.manager.ImageBuildOptions.BuildArgs = make(map[string]*string)
	}

	if el.contentGitConfigFile != "" {
		el.manager.ImageBuildOptions.BuildArgs["GITCONFIG_FILE"] = &el.contentGitConfigFile
	}

	if el.contentKnownHostsFile != "" {
		el.manager.ImageBuildOptions.BuildArgs["KNOWN_HOSTS_FILE"] = &el.contentKnownHostsFile
	}

	if el.contentIdRsaFile != "" {
		el.manager.ImageBuildOptions.BuildArgs["SSH_ID_RSA_FILE"] = &el.contentIdRsaFile
	}

	if el.contentIdEcdsaFile != "" {
		el.manager.ImageBuildOptions.BuildArgs["SSH_ID_ECDSA_FILE"] = &el.contentIdEcdsaFile
	}

	if el.gitPathPrivateRepository != "" {
		el.manager.ImageBuildOptions.BuildArgs["GIT_PRIVATE_REPO"] = &el.gitPathPrivateRepository
	}

	return
}

// Target
//
// English:
//
//	Build the specified stage as defined inside the Dockerfile.
//
//	 Input:
//	   value: stage name
//
// Note:
//
//   - See the multi-stage build docs for details.
//     See https://docs.docker.com/develop/develop-images/multistage-build/
//
// Português:
//
//	Monta o container a partir do estágio definido no arquivo Dockerfile.
//
//	 Entrada:
//	   value: nome do estágio
//
// Nota:
//
//   - Veja a documentação de múltiplos estágios para mais detalhes.
//     See https://docs.docker.com/develop/develop-images/multistage-build/
func (el *ContainerFromImage) Target(value string) {
	el.manager.ImageBuildOptions.Target = value
}

// Squash
//
// English:
//
//	Squash the resulting image's layers to the parent preserves the original image and creates a new
//	one from the parent with all the changes applied to a single layer
//
//	 Input:
//	   value: true preserve the original image and creates a new one from the parent
//
// Português:
//
//	Usa o conteúdo dos layers da imagem pai para criar uma imagem nova, preservando a imagem pai, e
//	aplica todas as mudanças a um novo layer
//
//	 Entrada:
//	   value: true preserva a imagem original e cria uma nova imagem a partir da imagem pai
func (el *ContainerFromImage) Squash(value bool) {
	el.manager.ImageBuildOptions.Squash = value
}

// Platform
//
// English:
//
//	Target platform containers for this service will run on, using the os[/arch[/variant]] syntax.
//
//	 Input:
//	   value: target platform
//
// Examples:
//
//	osx
//	windows/amd64
//	linux/arm64/v8
//
// Português:
//
//	Especifica a plataforma de container onde o serviço vai rodar, usando a sintaxe
//	os[/arch[/variant]]
//
//	 Entrada:
//	   value: plataforma de destino
//
// Exemplos:
//
//	osx
//	windows/amd64
//	linux/arm64/v8
func (el *ContainerFromImage) Platform(value string) {
	el.manager.ImageBuildOptions.Platform = value
}

// NoCache
//
// English:
//
//	Set image build no cache
//
// Português:
//
//	Define a opção `sem cache` para a construção da imagem
func (el *ContainerFromImage) NoCache() {
	el.enableCache = false
	el.manager.ImageBuildOptions.NoCache = true
}

//User memory constraints🔗
//We have four ways to set user memory usage:
//
//Option	Result
//memory=inf, memory-swap=inf (default)	There is no memory limit for the container. The container can use as much memory as needed.
//memory=L<inf, memory-swap=inf	(specify memory and set memory-swap as -1) The container is not allowed to use more than L bytes of memory, but can use as much swap as is needed (if the host supports swap memory).
//memory=L<inf, memory-swap=2*L	(specify memory without memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is double of that.
//memory=L<inf, memory-swap=S<inf, L<=S	(specify both memory and memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is limited by S.
//Examples:
//
//$ docker run -it ubuntu:14.04 /bin/bash
//We set nothing about memory, this means the processes in the container can use as much memory and swap memory as they need.
//
//$ docker run -it -m 300M --memory-swap -1 ubuntu:14.04 /bin/bash
//We set memory limit and disabled swap memory limit, this means the processes in the container can use 300M memory and as much swap memory as they need (if the host supports swap memory).
//
//$ docker run -it -m 300M ubuntu:14.04 /bin/bash
//We set memory limit only, this means the processes in the container can use 300M memory and 300M swap memory, by default, the total virtual memory size (--memory-swap) will be set as double of memory, in this case, memory + swap would be 2*300M, so processes can use 300M swap memory as well.
//
//$ docker run -it -m 300M --memory-swap 1G ubuntu:14.04 /bin/bash
//We set both memory and swap memory, so the processes in the container can use 300M memory and 700M swap memory.
//
//Memory reservation is a kind of memory soft limit that allows for greater sharing of memory. Under normal circumstances, containers can use as much of the memory as needed and are constrained only by the hard limits set with the -m/--memory option. When memory reservation is set, Docker detects memory contention or low memory and forces containers to restrict their consumption to a reservation limit.
//
//Always set the memory reservation value below the hard limit, otherwise the hard limit takes precedence. A reservation of 0 is the same as setting no reservation. By default (without reservation set), memory reservation is the same as the hard memory limit.
//
//Memory reservation is a soft-limit feature and does not guarantee the limit won’t be exceeded. Instead, the feature attempts to ensure that, when memory is heavily contended for, memory is allocated based on the reservation hints/setup.
//
//The following example limits the memory (-m) to 500M and sets the memory reservation to 200M.
//
//$ docker run -it -m 500M --memory-reservation 200M ubuntu:14.04 /bin/bash
//Under this configuration, when the container consumes memory more than 200M and less than 500M, the next system memory reclaim attempts to shrink container memory below 200M.
//
//The following example set memory reservation to 1G without a hard memory limit.
//
//$ docker run -it --memory-reservation 1G ubuntu:14.04 /bin/bash
//The container can use as much memory as it needs. The memory reservation setting ensures the container doesn’t consume too much memory for long time, because every memory reclaim shrinks the container’s consumption to the reservation.
//
//By default, kernel kills processes in a container if an out-of-memory (OOM) error occurs. To change this behaviour, use the --oom-kill-disable option. Only disable the OOM killer on containers where you have also set the -m/--memory option. If the -m flag is not set, this can result in the host running out of memory and require killing the host’s system processes to free memory.
//
//The following example limits the memory to 100M and disables the OOM killer for this container:
//
//$ docker run -it -m 100M --oom-kill-disable ubuntu:14.04 /bin/bash
//The following example, illustrates a dangerous way to use the flag:
//
//$ docker run -it --oom-kill-disable ubuntu:14.04 /bin/bash
//The container has unlimited memory which can cause the host to run out memory and require killing system processes to free memory. The --oom-score-adj parameter can be changed to select the priority of which containers will be killed when the system is out of memory, with negative scores making them less likely to be killed, and positive scores more likely.

// MemorySwap
//
// English:
//
//	Set memory swap (--memory-swap)
//
// Note:
//
//   - Use value * KKiloByte, value * KMegaByte and value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
//
// Português:
//
//	habilita a opção memory swp
//
// Note:
//
//   - Use value * KKiloByte, value * KMegaByte e value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
func (el *ContainerFromImage) MemorySwap(value int64) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.MemorySwap = value

	//e.addProblem("The SetImageBuildOptionsMemorySwap() function can generate an error when building the image.")
	return el
}

//User memory constraints🔗
//We have four ways to set user memory usage:
//
//Option	Result
//memory=inf, memory-swap=inf (default)	There is no memory limit for the container. The container can use as much memory as needed.
//memory=L<inf, memory-swap=inf	(specify memory and set memory-swap as -1) The container is not allowed to use more than L bytes of memory, but can use as much swap as is needed (if the host supports swap memory).
//memory=L<inf, memory-swap=2*L	(specify memory without memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is double of that.
//memory=L<inf, memory-swap=S<inf, L<=S	(specify both memory and memory-swap) The container is not allowed to use more than L bytes of memory, swap plus memory usage is limited by S.
//Examples:
//
//$ docker run -it ubuntu:14.04 /bin/bash
//We set nothing about memory, this means the processes in the container can use as much memory and swap memory as they need.
//
//$ docker run -it -m 300M --memory-swap -1 ubuntu:14.04 /bin/bash
//We set memory limit and disabled swap memory limit, this means the processes in the container can use 300M memory and as much swap memory as they need (if the host supports swap memory).
//
//$ docker run -it -m 300M ubuntu:14.04 /bin/bash
//We set memory limit only, this means the processes in the container can use 300M memory and 300M swap memory, by default, the total virtual memory size (--memory-swap) will be set as double of memory, in this case, memory + swap would be 2*300M, so processes can use 300M swap memory as well.
//
//$ docker run -it -m 300M --memory-swap 1G ubuntu:14.04 /bin/bash
//We set both memory and swap memory, so the processes in the container can use 300M memory and 700M swap memory.
//
//Memory reservation is a kind of memory soft limit that allows for greater sharing of memory. Under normal circumstances, containers can use as much of the memory as needed and are constrained only by the hard limits set with the -m/--memory option. When memory reservation is set, Docker detects memory contention or low memory and forces containers to restrict their consumption to a reservation limit.
//
//Always set the memory reservation value below the hard limit, otherwise the hard limit takes precedence. A reservation of 0 is the same as setting no reservation. By default (without reservation set), memory reservation is the same as the hard memory limit.
//
//Memory reservation is a soft-limit feature and does not guarantee the limit won’t be exceeded. Instead, the feature attempts to ensure that, when memory is heavily contended for, memory is allocated based on the reservation hints/setup.
//
//The following example limits the memory (-m) to 500M and sets the memory reservation to 200M.
//
//$ docker run -it -m 500M --memory-reservation 200M ubuntu:14.04 /bin/bash
//Under this configuration, when the container consumes memory more than 200M and less than 500M, the next system memory reclaim attempts to shrink container memory below 200M.
//
//The following example set memory reservation to 1G without a hard memory limit.
//
//$ docker run -it --memory-reservation 1G ubuntu:14.04 /bin/bash
//The container can use as much memory as it needs. The memory reservation setting ensures the container doesn’t consume too much memory for long time, because every memory reclaim shrinks the container’s consumption to the reservation.
//
//By default, kernel kills processes in a container if an out-of-memory (OOM) error occurs. To change this behaviour, use the --oom-kill-disable option. Only disable the OOM killer on containers where you have also set the -m/--memory option. If the -m flag is not set, this can result in the host running out of memory and require killing the host’s system processes to free memory.
//
//The following example limits the memory to 100M and disables the OOM killer for this container:
//
//$ docker run -it -m 100M --oom-kill-disable ubuntu:14.04 /bin/bash
//The following example, illustrates a dangerous way to use the flag:
//
//$ docker run -it --oom-kill-disable ubuntu:14.04 /bin/bash
//The container has unlimited memory which can cause the host to run out memory and require killing system processes to free memory. The --oom-score-adj parameter can be changed to select the priority of which containers will be killed when the system is out of memory, with negative scores making them less likely to be killed, and positive scores more likely.

// Memory
//
// English:
//
//	The maximum amount of memory the container can use.
//
//	 Input:
//	   value: amount of memory in bytes
//
// Note:
//
//   - If you set this option, the minimum allowed value is 4 * 1024 * 1024 (4 megabyte);
//   - Use value * KKiloByte, value * KMegaByte and value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
//
// Português:
//
//	Memória máxima total que o container pode usar.
//
//	 Entrada:
//	   value: Quantidade de memória em bytes
//
// Nota:
//
//   - Se você vai usar esta opção, o máximo permitido é 4 * 1024 * 1024 (4 megabyte)
//   - Use value * KKiloByte, value * KMegaByte e value * KGigaByte
//     See https://docs.docker.com/engine/reference/run/#user-memory-constraints
func (el *ContainerFromImage) Memory(value int64) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.Memory = value

	//e.addProblem("The SetImageBuildOptionsMemory() function can generate an error when building the image.")
	return el
}

// IsolationProcess
//
// English:
//
//	Set process isolation mode
//
// Português:
//
//	Determina o método de isolamento do processo
func (el *ContainerFromImage) IsolationProcess() (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.Isolation = dockerContainer.IsolationProcess
	return el
}

// IsolationHyperV
//
// English:
//
//	Set HyperV isolation mode
//
// Português:
//
//	Define o método de isolamento como sendo HyperV
func (el *ContainerFromImage) IsolationHyperV() (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.Isolation = dockerContainer.IsolationHyperV
	return el
}

// IsolationDefault
//
// English:
//
//	Set default isolation mode on current daemon
//
// Português:
//
//	Define o método de isolamento do processo como sendo o mesmo do deamon
func (el *ContainerFromImage) IsolationDefault() (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.Isolation = dockerContainer.IsolationDefault
	return el
}

// ExtraHosts
//
// English:
//
//	Add hostname mappings at build-time. Use the same values as the docker client --add-host
//	parameter.
//
//	 Input:
//	   values: hosts to mapping
//
// Example:
//
//	values = []string{
//	  "somehost:162.242.195.82",
//	  "otherhost:50.31.209.229",
//	}
//
//	An entry with the ip address and hostname is created in /etc/hosts inside containers for this
//	build, e.g:
//
//	  162.242.195.82 somehost
//	  50.31.209.229 otherhost
//
// Português:
//
//	Adiciona itens ao mapa de hostname durante o processo de construção da imagem. Use os mesmos
//	valores que em docker client --add-host parameter.
//
//	 Entrada:
//	   values: hosts para mapeamento
//
// Exemplo:
//
//	values = []string{
//	  "somehost:162.242.195.82",
//	  "otherhost:50.31.209.229",
//	}
//
//	Uma nova entrada com o endereço ip e hostname será criada dentro de /etc/hosts do container.
//	Exemplo:
//
//	  162.242.195.82 somehost
//	  50.31.209.229 otherhost
func (el *ContainerFromImage) ExtraHosts(values []string) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.ExtraHosts = values
	return el
}

// CacheFrom
//
// English:
//
//	Specifies images that are used for matching cache.
//
//	 Entrada:
//	   values: images that are used for matching cache.
//
// Note:
//
//	Images specified here do not need to have a valid parent chain to match cache.
//
// Português:
//
//	Especifica imagens que são usadas para correspondência de cache.
//
//	 Entrada:
//	   values: imagens que são usadas para correspondência de cache.
//
// Note:
//
//	As imagens especificadas aqui não precisam ter uma cadeia pai válida para corresponder a cache.
func (el *ContainerFromImage) CacheFrom(values []string) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.CacheFrom = values
	return el
}

// Shares
//
// English:
//
//	Set the CPU shares of the image build options.
//
//	 Input:
//	   value: CPU shares (Default: 1024)
//
//	Set this flag to a value greater or less than the default of 1024 to increase or reduce the
//	container’s weight, and give it access to a greater or lesser proportion of the host machine’s
//	CPU cycles.
//
//	This is only enforced when CPU cycles are constrained.
//
//	When plenty of CPU cycles are available, all containers use as much CPU as they need.
//
//	In that way, this is a soft limit. --cpu-shares does not prevent containers from being scheduled
//	in swarm mode.
//
//	It prioritizes container CPU resources for the available CPU cycles.
//
//	It does not guarantee or reserve any specific CPU access.
//
// Português:
//
//	Define o compartilhamento de CPU na construção da imagem.
//
//	 Entrada:
//	   value: Compartilhamento de CPU (Default: 1024)
//
//	Defina este sinalizador para um valor maior ou menor que o padrão de 1024 para aumentar ou reduzir
//	o peso do container e dar a ele acesso a uma proporção maior ou menor dos ciclos de CPU da máquina
//	host.
//
//	Isso só é aplicado quando os ciclos da CPU são restritos. Quando muitos ciclos de CPU estão
//	disponíveis, todos os container usam a quantidade de CPU de que precisam. Dessa forma, este é um
//	limite flexível. --cpu-shares não impede que os containers sejam agendados no modo swarm.
//
//	Ele prioriza os recursos da CPU do container para os ciclos de CPU disponíveis.
//
//	Não garante ou reserva nenhum acesso específico à CPU.
func (el *ContainerFromImage) Shares(value int64) (ref *ContainerFromImage) { //cpu
	el.manager.ImageBuildOptions.CPUShares = value
	return el
}

// Mems
//
// English:
//
//	Define a memory nodes (MEMs) (--cpuset-mems)
//
//	 Input:
//	   value: string with the format "0-3,5-7"
//
//	--cpuset-mems="" Memory nodes (MEMs) in which to allow execution (0-3, 0,1). Only effective on
//	NUMA systems.
//
//	If you have four memory nodes on your system (0-3), use --cpuset-mems=0,1 then processes in your
//	Docker container will only use memory from the first two memory nodes.
//
// Português:
//
//	Define memory node (MEMs) (--cpuset-mems)
//
//	 Entrada:
//	   value: string com o formato "0-3,5-7"
//
//	--cpuset-mems="" Memory nodes (MEMs) no qual permitir a execução (0-3, 0,1). Só funciona em
//	sistemas NUMA.
//
//	Se você tiver quatro nodes de memória em seu sistema (0-3), use --cpuset-mems=0,1 então, os
//	processos em seu container do Docker usarão apenas a memória dos dois primeiros nodes.
func (el *ContainerFromImage) Mems(value string) (ref *ContainerFromImage) { //cpu
	el.manager.ImageBuildOptions.CPUSetMems = value
	return el
}

// CPUs
//
// English:
//
//	Limit the specific CPUs or cores a container can use.
//
//	 Input:
//	   value: string with the format "1,2,3"
//
//	A comma-separated list or hyphen-separated range of CPUs a container can use, if you have more
//	than one CPU.
//
// The first CPU is numbered 0.
//
//	A valid value might be 0-3 (to use the first, second, third, and fourth CPU) or 1,3 (to use the
//	second and fourth CPU).
//
// Português:
//
//	Limite a quantidade de CPUs ou núcleos específicos que um container pode usar.
//
//	 Entrada:
//	   value: string com o formato "1,2,3"
//
//	Uma lista separada por vírgulas ou intervalo separado por hífen de CPUs que um container pode
//	usar, se você tiver mais de uma CPU.
//
//	A primeira CPU é numerada como 0.
//
//	Um valor válido pode ser 0-3 (para usar a primeira, segunda, terceira e quarta CPU) ou 1,3 (para
//	usar a segunda e a quarta CPU).
func (el *ContainerFromImage) CPUs(value string) (ref *ContainerFromImage) { //cpu
	el.manager.ImageBuildOptions.CPUSetCPUs = value

	//e.addProblem("The SetImageBuildOptionsCPUSetCPUs() function can generate an error when building the image.")
	return el
}

// Quota
//
// English:
//
//	Defines the host machine’s CPU cycles.
//
//	 Input:
//	   value: machine’s CPU cycles. (Default: 1024)
//
//	Set this flag to a value greater or less than the default of 1024 to increase or reduce the
//	container’s weight, and give it access to a greater or lesser proportion of the host machine’s
//	CPU cycles.
//
//	This is only enforced when CPU cycles are constrained. When plenty of CPU cycles are available,
//	all containers use as much CPU as they need. In that way, this is a soft limit. --cpu-shares does
//	not prevent containers from being scheduled in swarm mode. It prioritizes container CPU resources
//	for the available CPU cycles.
//
//	It does not guarantee or reserve any specific CPU access.
//
// Português:
//
//	Define os ciclos de CPU da máquina hospedeira.
//
//	 Entrada:
//	   value: ciclos de CPU da máquina hospedeira. (Default: 1024)
//
//	Defina este flag para um valor maior ou menor que o padrão de 1024 para aumentar ou reduzir o peso
//	do container e dar a ele acesso a uma proporção maior ou menor dos ciclos de CPU da máquina
//	hospedeira.
//
//	Isso só é aplicado quando os ciclos da CPU são restritos. Quando muitos ciclos de CPU estão
//	disponíveis, todos os containeres usam a quantidade de CPU de que precisam. Dessa forma, é um
//	limite flexível. --cpu-shares não impede que os containers sejam agendados no modo swarm. Ele
//	prioriza os recursos da CPU do container para os ciclos de CPU disponíveis.
//
//	Não garante ou reserva nenhum acesso específico à CPU.
func (el *ContainerFromImage) Quota(value int64) (ref *ContainerFromImage) { //cpu
	el.manager.ImageBuildOptions.CPUQuota = value

	//e.addProblem("The SetImageBuildOptionsCPUQuota() function can generate an error when building the image.")
	return el
}

// Period
//
// English:
//
//	Specify the CPU CFS scheduler period, which is used alongside --cpu-quota.
//
//	 Input:
//	   value: CPU CFS scheduler period
//
//	Defaults to 100000 microseconds (100 milliseconds). Most users do not change this from the
//	default.
//
//	For most use-cases, --cpus is a more convenient alternative.
//
// Português:
//
//	Especifique o período do agendador CFS da CPU, que é usado junto com --cpu-quota.
//
//	 Entrada:
//	   value: período do agendador CFS da CPU
//
//	O padrão é 100.000 microssegundos (100 milissegundos). A maioria dos usuários não altera o padrão.
//
//	Para a maioria dos casos de uso, --cpus é uma alternativa mais conveniente.
func (el *ContainerFromImage) Period(value int64) (ref *ContainerFromImage) { //cpu
	el.manager.ImageBuildOptions.CPUPeriod = value

	//e.addProblem("The SetImageBuildOptionsCPUPeriod() function can generate an error when building the image.")
	return el
}

// DockerfilePath
//
// English:
//
// Defines a Dockerfile to build the image.
//
// Português:
//
// Define um arquivo Dockerfile para construir a imagem.
func (el *ContainerFromImage) DockerfilePath(path string) (ref *ContainerFromImage) {
	el.manager.ImageBuildOptions.Dockerfile = path
	return el
}

// AutoDockerfileGenerator
//
// Defines the dockerfile generator object
func (el *ContainerFromImage) AutoDockerfileGenerator(autoDockerfile DockerfileAuto) (ref *ContainerFromImage) {
	el.autoDockerfile = autoDockerfile
	return el
}

// PrivateRepositoryAutoConfig
//
// English:
//
//	Copies the ssh ~/.ssh/id_rsa file and the ~/.gitconfig file to the SSH_ID_RSA_FILE and
//	GITCONFIG_FILE variables.
//
//	 Output:
//	   err: Standard error object
//
//	 Notes:
//	   * For change ssh key file name, use SetSshKeyFileName() function.
//
// Português:
//
//	Copia o arquivo ssh ~/.ssh/id_rsa e o arquivo ~/.gitconfig para as variáveis SSH_ID_RSA_FILE e
//	GITCONFIG_FILE.
//
//	 Saída:
//	   err: Objeto de erro padrão
//
//	 Notas:
//	   * Para mudar o nome do arquivo ssh usado como chave, use a função SetSshKeyFileName().
func (el *ContainerFromImage) PrivateRepositoryAutoConfig() (err error) {
	var userData *user.User
	var fileData []byte
	var filePathToRead string

	userData, err = user.Current()
	if err != nil {
		return
	}

	if el.sshDefaultFileName == "" {
		el.sshDefaultFileName, err = el.GetSshKeyFileName(userData.HomeDir)
		if err != nil {
			return
		}
	}

	filePathToRead = filepath.Join(userData.HomeDir, ".ssh", el.sshDefaultFileName)
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		return
	}

	el.contentIdRsaFile = string(fileData)
	el.contentIdRsaFileWithScape = strings.ReplaceAll(el.contentIdRsaFile, `"`, `\"`)

	filePathToRead = filepath.Join(userData.HomeDir, ".ssh", "known_hosts")
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		err = fmt.Errorf("container.PrivateRepositoryAutoConfig().ReadFile().error: %v", err)
		return
	}

	el.contentKnownHostsFile = string(fileData)
	el.contentKnownHostsFileWithScape = strings.ReplaceAll(el.contentKnownHostsFile, `"`, `\"`)

	filePathToRead = filepath.Join(userData.HomeDir, ".gitconfig")
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		err = fmt.Errorf("container.PrivateRepositoryAutoConfig().ReadFile().error: %v", err)
		return
	}

	el.contentGitConfigFile = string(fileData)
	el.contentGitConfigFileWithScape = strings.ReplaceAll(el.contentGitConfigFile, `"`, `\"`)

	el.addImageBuildOptionsGitCredentials()
	return
}

// GetSshKeyFileName
//
// English:
//
//	Returns the name of the last generated ssh key.
//
// Português:
//
//	Retorna o nome da chave ssh gerada por último.
func (el *ContainerFromImage) GetSshKeyFileName(dir string) (fileName string, err error) {
	var folderPath = path.Join(dir, ".ssh")

	var minDate = int64(math.MaxInt64)

	var files []fs.FileInfo
	files, err = ioutil.ReadDir(folderPath)

	for _, file := range files {
		var name = file.Name()
		var date = file.ModTime().UnixNano()

		if file.IsDir() == true {
			continue
		}

		if strings.HasPrefix(name, "id_") == true && strings.HasSuffix(name, ".pub") == false {
			if minDate >= date {
				minDate = date
				fileName = name
			}
		}
	}

	return
}

// GitCloneToBuildWithPrivateToken
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//	   privateToken: token defined on your git tool portal
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//	   privateToken: token definido no portal da sua ferramenta git
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
func (el *ContainerFromImage) GitCloneToBuildWithPrivateToken(url, privateToken string) {
	el.gitUrl = url
	el.gitPrivateToke = privateToken
}

// GitCloneToBuild
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e SetGitCloneToBuildWithUserPassworh().
func (el *ContainerFromImage) GitCloneToBuild(url string) {
	el.gitUrl = url
}

// GetGitCloneToBuild
//
// English:
//
//	Returns the URL of the repository to clone for image transformation
//
// Note:
//
//   - See the SetGitCloneToBuild() function for more details.
//
// Português:
//
//	Retorna a URL do repositório a ser clonado para a transformação em imagem
//
// Nota:
//
//   - Veja a função SetGitCloneToBuild() para mais detalhes.
func (el *ContainerFromImage) GetGitCloneToBuild() (url string) {
	return el.gitUrl
}

// GitSshPassword
//
// English:
//
//	Sets the password for the ssh key for private git repositories.
//
//	 Input:
//	   password: git ssh certificate password
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// Português:
//
//	Define a senha da chave ssh para repositórios git privados.
//
//	 Entrada:
//	   password: senha da chave ssh
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
func (el *ContainerFromImage) GitSshPassword(password string) {
	el.gitPassword = password
}

// GitCloneToBuildWithPrivateSSHKey
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//	   privateSSHKeyPath: this is the path of the private ssh key compatible with the public key
//	     registered in git
//	   password: password used when the ssh key was generated or empty string
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var privateSSHKeyPath string
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	privateSSHKeyPath = filepath.Join(usr.HomeDir, ".shh", "id_ecdsa")
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateSSHKey(url, privateSSHKeyPath, password)
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//	   privateSSHKeyPath: este é o caminho da chave ssh privada compatível com a chave pública
//	     cadastrada no git
//	   password: senha usada no momento que a chave ssh foi gerada ou string em branco
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var privateSSHKeyPath string
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	privateSSHKeyPath = filepath.Join(usr.HomeDir, ".shh", "id_ecdsa")
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateSSHKey(url, privateSSHKeyPath, password)
//	container.SetGitConfigFile(string(file))
func (el *ContainerFromImage) GitCloneToBuildWithPrivateSSHKey(url, privateSSHKeyPath, password string) {
	el.gitUrl = url
	el.gitSshPrivateKeyPath = privateSSHKeyPath
	el.gitPassword = password
}

// gitMakePublicSshKey
//
// English:
//
//	Mount the ssl certificate for the git clone function
//
//	 Output:
//	   publicKeys: Ponteiro de objeto compatível com o objeto ssh.PublicKeys
//	   err: standard error object
//
// Português:
//
//	 Monta o certificado ssl para a função de git clone
//
//		 Saída:
//	    publicKeys: Ponteiro de objeto compatível com o objeto ssh.PublicKeys
//	    err: objeto de erro padrão
func (el *ContainerFromImage) gitMakePublicSshKey() (publicKeys *sshGit.PublicKeys, err error) {
	if el.gitSshPrivateKeyPath != "" {
		_, err = os.Stat(el.gitSshPrivateKeyPath)
		if err != nil {
			err = fmt.Errorf("container.gitMakePublicSshKey().Stat().error: %v", err)
			return
		}
		publicKeys, err = sshGit.NewPublicKeysFromFile("git", el.gitSshPrivateKeyPath, el.gitPassword)
		if err != nil {
			err = fmt.Errorf("container.gitMakePublicSshKey().NewPublicKeysFromFile().error: %v", err)
			return
		}
	} else if el.contentIdEcdsaFile != "" {
		publicKeys, err = sshGit.NewPublicKeys("git", []byte(el.contentIdEcdsaFile), el.gitPassword)
		if err != nil {
			err = fmt.Errorf("container.gitMakePublicSshKey().NewPublicKeys().error: %v", err)
			return
		}
	} else if el.contentIdRsaFile != "" {
		publicKeys, err = sshGit.NewPublicKeys("git", []byte(el.contentIdRsaFile), el.gitPassword)
		if err != nil {
			err = fmt.Errorf("container.gitMakePublicSshKey().NewPublicKeys().error: %v", err)
			return
		}
	}

	return
}
