package manager

import (
	"fmt"
	"github.com/docker/docker/api/types/mount"
	networkTypes "github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/chaos/internal/builder"
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

	imageId       string
	imageName     string
	containerName string
	copies        int
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

	return el
}

func (el *ContainerFromImage) Create(imageName, containerName string, copies int) (ref *ContainerFromImage) {
	var err error

	if copies == 0 {
		return el
	}

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
				return
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
		var containerID string
		var warnings []string
		containerID, warnings, err = el.manager.DockerSys[i].ContainerCreateWithConfig(
			config,
			containerName+"_"+strconv.FormatInt(int64(i), 10),
			builder.KRestartPolicyNo,
			portConfig,
			volumes,
			netConfig,
		)
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerCreate().error: %v", i, err)
			return
		}

		//todo: fazer warnings - não deve ser erro
		if len(warnings) != 0 {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerCreate().warnings: %v", i, strings.Join(warnings, "; "))
			return
		}

		// fixme: apagar start
		err = el.manager.DockerSys[i].ContainerStart(containerID)
		if err != nil {
			el.manager.ErrorCh <- fmt.Errorf("container[%v].ContainerStart().error: %v", i, err)
			return
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
