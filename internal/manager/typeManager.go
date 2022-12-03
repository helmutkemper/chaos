package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/builder"
	"github.com/helmutkemper/chaos/internal/monitor"
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
	ImageBuildOptions types.ImageBuildOptions

	DoneCh  chan struct{}
	ErrorCh chan error
	FailCh  chan string
}

func (el *Manager) New() {
	var err error
	el.Id = make([]string, 0)
	el.DockerSys = make([]*builder.DockerSystem, 1)
	el.DockerSys[0] = new(builder.DockerSystem)

	el.DoneCh = make(chan struct{})
	el.ErrorCh = make(chan error)
	el.FailCh = make(chan string)

	el.addMonitor()

	err = el.DockerSys[0].Init()
	if err != nil {
		el.ErrorCh <- fmt.Errorf("chaos.Manager.New().error: %v. Usually this error occurs when docker is not running", err)
		return
	}

	return
}

func (el *Manager) addMonitor() {
	monitor.DoneChList = append(monitor.DoneChList, el.DoneCh)
	monitor.ErrorChList = append(monitor.ErrorChList, el.ErrorCh)
	monitor.FailChList = append(monitor.FailChList, el.FailCh)
}

func (el *Manager) Primordial() (primordial *Primordial) {
	primordial = new(Primordial)
	primordial.manager = el
	return
}

func (el *Manager) ContainerFromImage(imageName string) (containerFromImage *ContainerFromImage) {
	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	containerFromImage.imageName = imageName
	containerFromImage.command = "fromImage" //fixme: contante
	return
}

func (el *Manager) ContainerFromFolder(imageName, buildPath string) (containerFromImage *ContainerFromImage) {
	if !strings.Contains(imageName, "delete") {
		imageName = "delete_" + imageName
	}

	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	containerFromImage.buildPath = buildPath
	containerFromImage.imageName = imageName
	containerFromImage.command = "fromFolder" //fixme: contante
	return
}

func (el *Manager) ContainerFromGit(imageName, serverPath string) (containerFromImage *ContainerFromImage) {
	if !strings.Contains(imageName, "delete") {
		imageName = "delete_" + imageName
	}

	containerFromImage = new(ContainerFromImage)
	containerFromImage.manager = el
	containerFromImage.gitUrl = serverPath
	containerFromImage.imageName = imageName
	containerFromImage.command = "fromServer" //fixme: contante
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
