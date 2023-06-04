package manager

import (
	"bytes"
	"fmt"
	"github.com/helmutkemper/chaos/internal/monitor"
	"github.com/helmutkemper/chaos/internal/standalone"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var networkManagerGlobal *dockerNetwork

type Primordial struct {
	manager *Manager
}

func (el *Primordial) getLogs(id string) (log []byte, err error) {
	log, err = el.manager.DockerSys[0].ContainerLogs(id)
	if err != nil {
		return
	}

	logCopy := make([]byte, len(log))
	copy(logCopy, log)

	// Remove non-printable characters
	for i := 0; i != 32; i += 1 {
		if i == 10 {
			continue
		}

		log = bytes.ReplaceAll(log, []byte{uint8(i)}, []byte(""))
	}

	for {
		log = bytes.ReplaceAll(log, []byte("\n\n"), []byte(""))
		if bytes.Equal(log, logCopy) {
			break
		}

		logCopy = make([]byte, len(log))
		copy(logCopy, log)
	}

	return
}

func (el *Primordial) Test(t *testing.T, pathToSave string, names ...string) (ref *Primordial) {
	var log []byte

	err := os.MkdirAll(pathToSave, fs.ModePerm)
	if err != nil {
		ErrorCh <- fmt.Errorf("primordial.Test().error: %v", err)
		return el
	}

	t.Cleanup(func() {

		// Saves contents of containers before deleting
		containers, err := el.manager.DockerSys[0].ContainerListAll()
		if err != nil {
			ErrorCh <- fmt.Errorf("primordial.NetworkCreate().error: %v", err)
			return
		}

		for _, container := range containers {
			if len(container.Names) == 0 {
				continue
			}

			if strings.Contains(container.Names[0], "delete") {
				log, err = el.getLogs(container.ID)
				if err != nil {
					ErrorCh <- fmt.Errorf("primordial.Test().error: %v", err)
					return
				}

				pathAbs, err := filepath.Abs(pathToSave)
				if err != nil {
					ErrorCh <- fmt.Errorf("primordial.Test().error: %v", err)
					return
				}

				pathData := path.Join(pathAbs, container.Names[0]+".log")
				err = os.WriteFile(pathData, log, fs.ModePerm)
				if err != nil {
					ErrorCh <- fmt.Errorf("primordial.Test().error: %v", err)
					return
				}
			}

			for _, name := range names {
				if strings.Contains(container.Names[0], name) {
					log, err = el.getLogs(container.ID)
					if err != nil {
						ErrorCh <- fmt.Errorf("primordial.Test().error: %v", err)
						return
					}

					pathData := path.Join(pathToSave, container.Names[0]+".log")
					err = os.WriteFile(pathData, log, fs.ModePerm)
				}
			}
		}

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
		ErrorCh <- fmt.Errorf("primordial.NetworkCreate().error: %v", err)
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
		select {
		case <-timer.C:
		case <-el.manager.DoneCh:
		}
		monitor.EndAll()
	}()

	return monitor.Monitor()
}

// Done
//
// End of test before requested time
func (el *Primordial) Done() {
	el.manager.DoneCh <- struct{}{}
}

// GetLastError
//
// Returns the last error from the test
func (el *Primordial) GetLastError() (err error) {
	if monitor.Err == false {
		return nil
	}

	return <-ErrorCh
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
