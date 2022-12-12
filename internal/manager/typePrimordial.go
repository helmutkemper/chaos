package manager

import (
	"fmt"
	"github.com/helmutkemper/chaos/internal/monitor"
	"github.com/helmutkemper/chaos/internal/standalone"
	"strings"
	"time"
)

var networkManagerGlobal *dockerNetwork

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
	var err error

	if !strings.Contains(name, "delete") {
		name = "delete_" + name
	}

	if err = el.manager.networkCreate(name, subnet, gateway); err != nil {
		el.manager.ErrorCh <- fmt.Errorf("primordial.NetworkCreate().error: %v", err)
		return el
	}

	return el
}

func (el *Primordial) Monitor(duration time.Duration) (pass bool) {
	var timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		monitor.EndAll()
	}()

	return monitor.Monitor()
}

func (el *Primordial) GarbageCollector() (ref *Primordial) {
	standalone.GarbageCollector()
	return el
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
