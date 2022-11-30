package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/builder"
	"strings"
	"time"
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

	TickerStats       *time.Ticker
	TickerFail        *time.Ticker
	Id                []string
	DockerSys         []*builder.DockerSystem
	ErrorCh           chan error
	ImageBuildOptions types.ImageBuildOptions
}

func (el *Manager) New(errorCh chan error) {
	el.ErrorCh = errorCh

	var err error
	el.Id = make([]string, 0)
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

func (el *Manager) ContainerFromImage(imageName string) (containerFromImage *ContainerFromImage) {
	if !strings.Contains(imageName, "delete") { //todo: function?
		imageName = "delete_" + imageName
	}

	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	containerFromImage.imageName = imageName
	containerFromImage.command = "fromImage" //fixme: contante
	return
}

func (el *Manager) ContainerFromFolder(imageName, buildPath string) (containerFromImage *ContainerFromImage) {
	if !strings.Contains(imageName, "delete") { //todo: function?
		imageName = "delete_" + imageName
	}

	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	containerFromImage.buildPath = buildPath
	containerFromImage.imageName = imageName
	containerFromImage.command = "fromFolder" //fixme: contante
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
