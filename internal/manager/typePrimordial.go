package manager

import (
	"fmt"
	"github.com/helmutkemper/chaos/internal/monitor"
	"github.com/helmutkemper/chaos/internal/standalone"
	"strings"
	"testing"
	"time"
)

var networkManagerGlobal *dockerNetwork

type Primordial struct {
	manager *Manager
}

func (el *Primordial) Test(t *testing.T, names ...string) (ref *Primordial) {
	t.Cleanup(func() {
		el.GarbageCollector(names...)
	})

	return el
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

	if _, err = el.manager.networkCreate(name, subnet, gateway); err != nil {
		el.manager.ErrorCh <- fmt.Errorf("primordial.NetworkCreate().error: %v", err)
		return el
	}

	return el
}

// Monitor
//
// Monitors the test for errors while waiting for the test to end.
//
//	Notes:
//	  * When the test timer ends, Monitor() waits for all test pipelines to finish, hooking up all containers at the
//	    end of the test
func (el *Primordial) Monitor(duration time.Duration) (pass bool) {
	var timer = time.NewTimer(duration)
	go func() {
		<-timer.C
		monitor.EndAll()
	}()

	return monitor.Monitor()
}

// GarbageCollector
//
// Deletes all Docker elements with `delete` in the name.
//
//	Input:
//	  names: additional list of terms to be deleted
//
//	Example:
//	  GarbageCollector("mongo") will delete any docker element with the term `mongo` contained in the name, which
//	  includes the image named `mongo:latest`, container named `mongodb` and network name `mongodb_network`
func (el *Primordial) GarbageCollector(names ...string) (ref *Primordial) {
	standalone.GarbageCollector(names...)
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
