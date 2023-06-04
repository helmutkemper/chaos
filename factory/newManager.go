package factory

import (
	"github.com/helmutkemper/chaos/internal/manager"
	"github.com/helmutkemper/chaos/internal/standalone"
	"strconv"
)

type ProxyConfig struct {
	LocalPort   int64
	Destination string
	MinDelay    int64
	MaxDelay    int64
}

// NewContainerNetworkProxy
//
// Create a container with a proxy simulating a slow network.
//
//	Input:
//	  containerName: name of container
//	  localPort: connection port. eg.: 27016 for MongoDB
//	  destination: container destination. eg. delete_mongo_0:27017 for MongoDB
//	  minDelay: min delay in milliseconds for block of 32k bytes. Use 0 for default value
//	  maxDelay: max delay in milliseconds for block of 32k bytes. Use 0 for default value
func NewContainerNetworkProxy(containerName string, config []ProxyConfig) (reference *manager.ContainerFromImage) {

	envFinal := make([][]string, 0)
	for _, conf := range config {
		localPortString := ":" + strconv.FormatInt(conf.LocalPort, 10)
		environmentVars := make([]string, 0)
		if conf.MinDelay != 0 {
			environmentVars = append(environmentVars, "CHAOS_NETWORK_MIN_DELAY="+strconv.FormatInt(conf.MinDelay, 10))
		}
		if conf.MaxDelay != 0 {
			environmentVars = append(environmentVars, "CHAOS_NETWORK_MAX_DELAY="+strconv.FormatInt(conf.MaxDelay, 10))
		}
		environmentVars = append(environmentVars, "CHAOS_NETWORK_LOCAL_PORT="+localPortString)
		environmentVars = append(environmentVars, "CHAOS_NETWORK_REMOTE_CONTAINER="+conf.Destination)

		envFinal = append(envFinal, environmentVars)
	}

	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromGit(
		"delay",
		"https://github.com/helmutkemper/chaos.network.git",
	).
		//EnvironmentVar(environmentVars).
		EnvironmentVar(envFinal...).
		//Ports("tcp", localPort, localPort).
		MakeDockerfile().
		Create(containerName, len(config)).
		Start()
}

func NewContainerFromGit(imageName, serverPath string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromGit(imageName, serverPath).
		Reports()
}

func NewContainerFromFolder(imageName, buildPath string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromFolder(imageName, buildPath).
		Reports()
}

func NewContainerFromImage(imageName string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromImage(imageName).
		Reports()
}

func NewPrimordial() (reference *manager.Primordial) {
	standalone.GarbageCollector()

	ref := new(manager.Manager)
	ref.New()
	reference = ref.Primordial()
	return
}
