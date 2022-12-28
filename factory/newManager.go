package factory

import (
	"github.com/helmutkemper/chaos/internal/manager"
	"github.com/helmutkemper/chaos/internal/standalone"
	"strconv"
)

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
//	  changeRate: percentage value of blocks with damaged values. Values between 0.0 and 1.0 where 0 is disabled and 1.0 is 100%
func NewContainerNetworkProxy(containerName string, localPort int64, destination string, minDelay, maxDelay int64, changeRate float64) (reference *manager.ContainerFromImage) {

	localPortString := ":" + strconv.FormatInt(localPort, 10)
	environmentVars := make([]string, 0)
	if minDelay != 0 {
		environmentVars = append(environmentVars, "CHAOS_NETWORK_MIN_DELAY="+strconv.FormatInt(minDelay, 10))
	}
	if maxDelay != 0 {
		environmentVars = append(environmentVars, "CHAOS_NETWORK_MAX_DELAY="+strconv.FormatInt(maxDelay, 10))
	}
	if changeRate != 0 {
		environmentVars = append(environmentVars, "CHAOS_NETWORK_MAX_DELAY="+strconv.FormatFloat(changeRate, 'g', -1, 64))
	}
	environmentVars = append(environmentVars, "CHAOS_NETWORK_LOCAL_PORT="+localPortString)
	environmentVars = append(environmentVars, "CHAOS_NETWORK_REMOTE_CONTAINER="+destination)

	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromGit(
		"delay",
		"https://github.com/helmutkemper/chaos.network.git",
	).
		EnvironmentVar(environmentVars).
		Ports("tcp", localPort, localPort).
		MakeDockerfile().
		Create(containerName, 1).
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
