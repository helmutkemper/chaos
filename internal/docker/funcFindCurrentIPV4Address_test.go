package docker

import (
	"errors"
	"github.com/helmutkemper/util"
	"time"
)

func ExampleContainerBuilder_FindCurrentIPV4Address() {
	var err error

	SaGarbageCollector()

	var container = ContainerBuilder{}

	// set image name for docker pull
	container.SetImageName("nats:latest")

	// set a container name
	container.SetContainerName("container_delete_nats_after_test")

	// set a waits for the text to appear in the standard container output to proceed [optional]
	container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)

	// inialize the container object
	err = container.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	// image nats:latest pull command [optional]
	err = container.ImagePull()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	// container build and start from image nats:latest
	// waits for the text "Listening for route connections on 0.0.0.0:6222" to appear  in the standard container output
	// to proceed
	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	var IP string
	IP, err = container.FindCurrentIPV4Address()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	if IP != container.GetIPV4Address() {
		err = errors.New("all ip address must be a samer IP")
		util.TraceToLog()
		panic(err)
	}

	// container "container_delete_nats_after_test" running and ready for use on this code point on var IP
	// all nats ports are open
	// you can use AddPortToExpose("4222"), to open only ports defineds inside code;
	// you can use AddPortToChange("4222", "1111") to open only ports defineds inside code and change port 4222 to port
	// 1111;
	// you can use SetDoNotOpenContainersPorts() to not open containers ports

	SaGarbageCollector()

	// use this function to remove image, ONLY before container stoped and deleted
	err = container.ImageRemoveByName("nats:latest")
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	// Output:
	//
}
