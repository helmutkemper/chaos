package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/builder"
	"github.com/helmutkemper/chaos/internal/monitor"
	"strings"
	"time"
)

var ErrorCh chan error

type dockerNetwork struct {
	generator   *builder.NextNetworkAutoConfiguration
	networkID   string
	networkName string
}

type chaosAction struct {
	time    time.Time
	action  func(string) error
	display string
	id      string
}

type chaosConfig struct {
	minimumTimeDelay         time.Duration
	maximumTimeDelay         time.Duration
	minimumTimeToUnpause     time.Duration
	maximumTimeToUnpause     time.Duration
	minimumTimeBeforeRestart time.Duration
	maximumTimeBeforeRestart time.Duration

	minimumTimeToStartChaos time.Duration
	maximumTimeToStartChaos time.Duration
	minimumTimeToPause      time.Duration
	maximumTimeToPause      time.Duration

	minimumTimeToRestart       time.Duration
	maximumTimeToRestart       time.Duration
	restartProbability         float64
	restartChangeIpProbability float64
	restartLimit               int
}

type Chaos struct {
	Type   string
	Action []chaosAction
}

type Manager struct {
	network *dockerNetwork

	TickerStats       *time.Ticker
	TickerFail        *time.Ticker
	Id                []string
	DockerSys         []*builder.DockerSystem
	Chaos             []Chaos
	ChaosConfig       chaosConfig
	ImageBuildOptions types.ImageBuildOptions

	DoneCh chan struct{}

	FailCh chan string
}

func (el *Manager) New() {
	var err error

	if ErrorCh == nil {
		ErrorCh = make(chan error, 10)
	}

	el.Id = make([]string, 0)
	el.DockerSys = make([]*builder.DockerSystem, 1)
	el.DockerSys[0] = new(builder.DockerSystem)

	el.DoneCh = make(chan struct{})
	el.FailCh = make(chan string)

	el.ChaosConfig.maximumTimeDelay = 90 * time.Second
	el.ChaosConfig.minimumTimeDelay = 30 * time.Second

	el.ChaosConfig.maximumTimeBeforeRestart = 90 * time.Second
	el.ChaosConfig.minimumTimeBeforeRestart = 30 * time.Second

	el.ChaosConfig.maximumTimeToUnpause = 90 * time.Second
	el.ChaosConfig.minimumTimeToUnpause = 30 * time.Second

	el.addMonitor()

	err = el.DockerSys[0].Init()
	if err != nil {
		ErrorCh <- fmt.Errorf("chaos.Manager.New().error: %v. Usually this error occurs when docker is not running", err)
		return
	}

	return
}

func (el *Manager) addMonitor() {
	monitor.DoneChList = append(monitor.DoneChList, el.DoneCh)
	monitor.ErrorChList = append(monitor.ErrorChList, ErrorCh)
	monitor.FailChList = append(monitor.FailChList, el.FailCh)
}

func (el *Manager) Primordial() (primordial *Primordial) {
	primordial = new(Primordial)
	primordial.manager = el
	return
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
//	  * If a network with the same name and different configuration already exists, the network will be deleted, and a new network created.
func (el *Manager) networkCreate(name, subnet, gateway string) (networkId string, err error) {
	el.network = new(dockerNetwork)
	el.network.networkName = name

	var networkList []types.NetworkResource
	networkList, err = el.DockerSys[0].NetworkList()
	if err != nil {
		return
	}

	for _, networkData := range networkList {
		if networkData.Name == name {
			el.network.networkID = networkData.ID

			var data types.NetworkResource
			if data, err = el.DockerSys[0].NetworkInspect(networkData.ID); err != nil {
				err = fmt.Errorf("network.NetworkCreate().NetworkInspect().error: %v", err)
				return
			}
			if data.IPAM.Config[0].Subnet != subnet || data.IPAM.Config[0].Gateway != gateway {
				if err = el.DockerSys[0].NetworkRemove(networkData.ID); err != nil {
					err = fmt.Errorf("network.NetworkCreate().NetworkRemove().error: %v", err)
					return
				}

				break
			}

			el.network.generator = el.DockerSys[0].NetworkGetGenerator(name)
			return
		}
	}

	if el.network.networkID, el.network.generator, err = el.DockerSys[0].NetworkCreate(name, builder.KNetworkDriveBridge, "local", subnet, gateway); err != nil {
		err = fmt.Errorf("network.NetworkCreate().NetworkCreate().error: %v", err)
		return
	}

	networkManagerGlobal = el.network

	networkId = el.network.networkID
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
