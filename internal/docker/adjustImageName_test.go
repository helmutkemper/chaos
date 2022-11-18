package iotmakerdocker

import "fmt"

func ExampleDockerSystem_AdjustImageName() {

	var err error
	var dockerSys *DockerSystem
	var correctImageName string

	dockerSys = &DockerSystem{}
	dockerSys.ContextCreate()
	err = dockerSys.ClientCreate()
	if err != nil {
		panic(err)
	}

	correctImageName = dockerSys.AdjustImageName("alpine")
	fmt.Printf("%v\n", correctImageName)

	correctImageName = dockerSys.AdjustImageName("alpine:")
	fmt.Printf("%v\n", correctImageName)

	// Output:
	// alpine:latest
	// alpine:latest
}
