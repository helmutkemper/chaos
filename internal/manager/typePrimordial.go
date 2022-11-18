package manager

import "fmt"

type Primordial struct {
	*Manager
}

// NetworkCreate
//
// Create a docker network to be used in the chaos test
//
//	 Input:
//	    name: network name
//			subnet: subnet value. eg. 10.0.0.0/16
//			gateway: gateway value. eg. "10.0.0.1
//
//		Notes:
//	   * If there is already a network with the same name and the same configuration, nothing will be done;
//		  * If a network with the same name and different configuration already exists, the network will be deleted and a new network created.
func (el *Primordial) NetworkCreate(name, subnet, gateway string) (err error) {
	network := new(Network)
	network.New(el.Manager)

	if err = network.NetworkCreate(name, subnet, gateway); err != nil {
		err = fmt.Errorf("primordial.NetworkCreate().error: %v", err)
		return
	}

	return
}
