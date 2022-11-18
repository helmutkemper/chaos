package iotmakerdocker

import (
	"errors"
)

func ExampleDockerSystem_NetworkVerifyName() {

	var err error
	var networkId string
	var dockerSys *DockerSystem

	// English: create a new default client. Please, use: err, dockerSys = factoryDocker.NewClient()
	// Português: cria um novo cliente com configurações padrão. Por favor, usr: err, dockerSys = factoryDocker.NewClient()
	dockerSys = &DockerSystem{}
	dockerSys.ContextCreate()
	err = dockerSys.ClientCreate()
	if err != nil {
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// English: create a network named 'network_delete_before_test'
	// Português: cria uma nova rede de nome 'network_delete_before_test'
	networkId, _, err = dockerSys.NetworkCreate(
		"network_delete_before_test",
		KNetworkDriveBridge,
		"local",
		"10.0.0.0/16",
		"10.0.0.1",
	)
	if err != nil {
		panic(err)
	}

	if networkId == "" {
		err = errors.New("network id was not generated")
		panic(err)
	}

	var networkFound bool
	networkFound, err = dockerSys.NetworkVerifyName("network_delete_before_test")
	if err != nil {
		panic(err)
	}

	if networkFound == false {
		err = errors.New("network network_delete_before_test not found")
		panic(err)
	}

	err = dockerSys.NetworkRemove(networkId)
	if err != nil {
		panic(err)
	}

	networkFound, err = dockerSys.NetworkVerifyName("network_delete_before_test")
	if err != nil {
		panic(err)
	}

	if networkFound == true {
		err = errors.New("network removal erro")
		panic(err)
	}

	// English: garbage collector and deletes networks and images whose name contains the term 'delete'
	// Português: coletor de lixo e apaga redes e imagens cujo o nome contém o temo 'delete'
	err = dockerSys.RemoveAllByNameContains("delete")
	if err != nil {
		panic(err)
	}

	// Output:
	//
}
