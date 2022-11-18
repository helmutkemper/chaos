package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"reflect"
	"testing"
)

func TestDockerSystem_NetworkFindNetworkTypeBridgePublic(t *testing.T) {
	var err error
	var inspect types.NetworkResource

	dockerSys := DockerSystem{}
	err = dockerSys.Init()
	if err != nil {
		t.Fail()
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	_, _, err = dockerSys.NetworkCreate(
		"delete_before_test",
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)
	if err != nil {
		t.Fail()
		panic(err)
	}

	inspect, err = dockerSys.NetworkFindNetworkTypeBridgePublic()
	if err != nil {
		t.Fail()
		panic(err)
	}

	if reflect.DeepEqual(inspect, types.NetworkResource{}) == true {
		t.Fail()
		panic(errors.New("pubic network not found"))
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}
}
