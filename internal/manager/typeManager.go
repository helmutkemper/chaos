package manager

import (
	"fmt"
	"github.com/helmutkemper/chaos/internal/builder"
)

//
//    +-----------+          +--------------+
//    |  manager  +------+---+  primordial  |
//    +-----------+      |   +--------------+
//                       |
//                       |   +-------------+
//                       +---+  container  |
//                       |   +-------------+
//                       |
//                       |   +---------+
//                       +---+  basic  |
//                           +---------+
//
//    +--------------+       +--------------------+
//    |  primordial  +---+---+  garbage colector  |
//    +--------------+   |   +--------------------+
//                       |
//                       |   +------------------+
//                       +---+  create network  |
//                           +------------------+
//
//    +-------------+       +--------------+
//    |  container  +---+---+  from image  +---+--- ports()
//    +-------------+   |   +--------------+   |
//                      |                      +--- volumes()
//                      |                      |
//                      |                      +---
//                      |
//                      |
//                      |
//                      |
//                      |   +---------------+
//                      +---+  from folder  |
//                      |   +---------------+
//                      |
//                      |   +---------------+
//                      +---+  from server  |
//                          +---------------+
//
//
//
//
//
//

type dockerNetwork struct {
	generator   *builder.NextNetworkAutoConfiguration
	networkID   string
	networkName string
}

type Manager struct {
	network *dockerNetwork

	DockerSys []*builder.DockerSystem
	ErrorCh   chan error
}

func (el *Manager) New(errorCh chan error) {
	el.ErrorCh = errorCh

	var err error
	el.DockerSys = make([]*builder.DockerSystem, 1)
	el.DockerSys[0] = new(builder.DockerSystem)
	err = el.DockerSys[0].Init()
	if err != nil {
		el.ErrorCh <- fmt.Errorf("chaos.Manager.New().error: %v. Usually this error occurs when docker is not running", err)
		return
	}

	return
}

func (el *Manager) Primordial() (primordial *Primordial) {
	primordial = new(Primordial)
	primordial.manager = el
	return
}

func (el *Manager) ContainerFromImage() (containerFromImage *ContainerFromImage) {
	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	return
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
