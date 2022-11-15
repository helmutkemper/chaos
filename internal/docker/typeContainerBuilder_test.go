package docker

import (
	dockerNetwork "github.com/helmutkemper/iotmaker.docker.builder.network"
	"github.com/helmutkemper/util"
	"testing"
	"time"
)

func TestContainer_1(t *testing.T) {
	var err error

	var dockerSys DockerSystem
	err = dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	var netDocker = dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	var container = ContainerBuilder{}
	container.SetNetworkDocker(&netDocker)
	container.SetImageName("nats:latest")
	container.SetContainerName("container_delete_nats_after_test")
	container.AddPortToExpose("4222")
	container.SetWaitString("Listening for route connections on 0.0.0.0:6222")

	err = container.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = container.imagePull()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}
}

func TestContainer_2(t *testing.T) {
	var err error

	var dockerSys DockerSystem
	err = dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	var netDocker = dockerNetwork.ContainerBuilderNetwork{}
	err = netDocker.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	var container = ContainerBuilder{}
	container.SetNetworkDocker(&netDocker)
	container.SetImageName("nats:latest")
	container.SetContainerName("container_delete_nats_after_test")
	container.AddPortToChange("4222", "4200")
	container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)

	err = container.Init()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = container.imagePull()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = container.ContainerBuildAndStartFromImage()
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}

	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		util.TraceToLog()
		t.FailNow()
	}
}
