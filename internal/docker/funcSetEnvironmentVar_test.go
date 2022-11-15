package docker

import (
	dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
	"time"
)

func ExampleContainerBuilder_SetEnvironmentVar() {
	var err error

	SaGarbageCollector()

	var netDocker = dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		panic(err)
	}

	// create a network named cache_delete_after_test, subnet 10.0.0.0/16 e gatway 10.0.0.1
	err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		panic(err)
	}

	// At this point in the code, the network has been created and is ready for use

	var mongoDocker = &ContainerBuilder{}
	// set a docker network
	//mongoDocker.SetNetworkDocker(&netDocker)
	// set an image name
	mongoDocker.SetImageName("mongo:latest")
	// set a container name
	mongoDocker.SetContainerName("container_delete_mongo_after_test")
	// set a port to expose
	mongoDocker.AddPortToExpose("27017")
	// se a environment var list
	mongoDocker.SetEnvironmentVar(
		[]string{
			"--host 0.0.0.0",
		},
	)
	// set a MongoDB data dir to ./test/data
	err = mongoDocker.AddFileOrFolderToLinkBetweenComputerHostAndContainer("./test/data", "/data")
	if err != nil {
		panic(err)
	}
	// set a text indicating for container ready for use
	mongoDocker.SetWaitStringWithTimeout(`"msg":"Waiting for connections","attr":{"port":27017`, 20*time.Second)

	// inicialize the object before sets
	err = mongoDocker.Init()
	if err != nil {
		panic(err)
	}

	// build a container
	err = mongoDocker.ContainerBuildAndStartFromImage()
	if err != nil {
		panic(err)
	}

	// Output:
	//

	// At this point, the MongoDB is ready for use on port 27017

	// Stop and delete the container
	// SaGarbageCollector()

	// Output:
	//
}
