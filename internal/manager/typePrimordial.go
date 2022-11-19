package manager

import "fmt"

type Primordial struct {
	manager *Manager
}

// NetworkCreate
//
// Create a docker network to be used in the chaos test
//
//	Input:
//	  name: network name
//	  subnet: subnet value. eg. 10.0.0.0/16
//	  gateway: gateway value. eg. "10.0.0.1
//
//	Notes:
//	  * If there is already a network with the same name and the same configuration, nothing will be done;
//	  * If a network with the same name and different configuration already exists, the network will be deleted and a new network created.
func (el *Primordial) NetworkCreate(name, subnet, gateway string) (ref *Primordial) {
	network := new(Network)
	network.New(el.manager)

	var err error
	if err = network.NetworkCreate(name, subnet, gateway); err != nil {
		el.manager.ErrorCh <- fmt.Errorf("primordial.NetworkCreate().error: %v", err)
		return el
	}

	return el
}
