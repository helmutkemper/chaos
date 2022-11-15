package docker

import (
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/helmutkemper/util"
	"sort"
)

func ExampleContainerBuilder_ImageListExposedPorts() {
	var err error
	var portList []nat.Port

	// create a container
	var container = ContainerBuilder{}
	// set image name for docker pull
	container.SetImageName("nats:latest")
	err = container.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	err = container.ImagePull()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	portList, err = container.ImageListExposedPorts()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	var portsToPrint = make([]string, 0)

	for _, p := range portList {
		portsToPrint = append(portsToPrint, fmt.Sprintf("port: %v/%v\n", p.Port(), p.Proto()))
	}

	sort.Strings(portsToPrint)

	for _, print := range portsToPrint {
		fmt.Printf("%v", print)
	}

	err = container.ImageRemove()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}

	// Output:
	// port: 4222/tcp
	// port: 6222/tcp
	// port: 8222/tcp
}
