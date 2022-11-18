package manager

import "fmt"

type Primordial struct {
	*Manager
}

func (el *Primordial) NetworkCreate(name, subnet, gateway string) (err error) {
	network := new(Network)
	network.New(el.Manager)

	if err = network.NetworkCreate(name, subnet, gateway); err != nil {
		err = fmt.Errorf("primordial.NetworkCreate().error: %v", err)
		return
	}

	return
}
