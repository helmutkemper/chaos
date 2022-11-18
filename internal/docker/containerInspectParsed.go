package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/go-connections/nat"
	"time"
)

type ContainerInspect struct {
	ID      string
	Name    string
	Created time.Time
	Args    []string

	ImageId      string
	RestartCount int
	Volumes      ContainerInspectVolumes
	State        ContainerInspectState
	Host         ContainerInspectHost
	Network      ContainerNetwork
}

type ContainerInspectState struct {
	Status     string
	Running    bool
	Paused     bool
	Restarting bool
	OOMKilled  bool
	Dead       bool
	ExitCode   int
	StartedAt  time.Time
	FinishedAt time.Time
}

type ContainerInspectHost struct {
	NetworkMode                        string
	Binds                              []string
	IsDefault                          bool
	IsHost                             bool
	IsContainer                        bool
	IsPrivate                          bool
	IsNone                             bool
	IsBridge                           bool
	IsUserDefined                      bool
	UserDefined                        string
	PortBindings                       nat.PortMap
	PortRestartPolicyName              string
	PortRestartPolicyMaximumRetryCount int
	ExposedPorts                       nat.PortSet
}

type ContainerInspectVolumes struct {
	VolumeDriver string
	VolumesFrom  []string
	Mounts       []InspectVolumeMount
}

type Type string
type Consistency string

type InspectVolumeMount struct {
	Type Type `json:",omitempty"`
	// Source specifies the name of the mount. Depending on mount type, this
	// may be a volume name or a host path, or even ignored.
	// Source is not supported for tmpfs (must be an empty value)
	Source      string      `json:",omitempty"`
	Target      string      `json:",omitempty"`
	ReadOnly    bool        `json:",omitempty"`
	Consistency Consistency `json:",omitempty"`
}

type EndpointIPAMConfig struct {
	IPv4Address  string   `json:",omitempty"`
	IPv6Address  string   `json:",omitempty"`
	LinkLocalIPs []string `json:",omitempty"`
}

type EndpointSettings struct {
	// Configurations
	IPAMConfig EndpointIPAMConfig
	Links      []string
	Aliases    []string
	// Operational data
	NetworkID           string
	EndpointID          string
	Gateway             string
	IPAddress           string
	IPPrefixLen         int
	IPv6Gateway         string
	GlobalIPv6Address   string
	GlobalIPv6PrefixLen int
	MacAddress          string
	DriverOpts          map[string]string
}

type ContainerNetwork struct {
	Ports       nat.PortMap
	Gateway     string
	IPAddress   string
	IPPrefixLen int
	IPv6Gateway string
	MacAddress  string
	Networks    map[string]EndpointSettings
}

// ContainerInspectParsed testing. do not use
func (el *DockerSystem) ContainerInspectParsed(
	id string,
) (
	parsed ContainerInspect,
	err error,
) {

	var inspect types.ContainerJSON
	inspect, _, err = el.cli.ContainerInspectWithRaw(el.ctx, id, true)
	if err != nil {
		return
	}

	parsed.Network.Networks = make(map[string]EndpointSettings)

	if inspect.NetworkSettings != nil {

		for k, network := range inspect.NetworkSettings.Networks {
			var settings = EndpointSettings{
				Links:               (*network).Links,
				Aliases:             (*network).Aliases,
				NetworkID:           (*network).NetworkID,
				EndpointID:          (*network).EndpointID,
				Gateway:             (*network).Gateway,
				IPAddress:           (*network).IPAddress,
				IPPrefixLen:         (*network).IPPrefixLen,
				IPv6Gateway:         (*network).IPv6Gateway,
				GlobalIPv6Address:   (*network).GlobalIPv6Address,
				GlobalIPv6PrefixLen: (*network).GlobalIPv6PrefixLen,
				MacAddress:          (*network).MacAddress,
				DriverOpts:          (*network).DriverOpts,
			}

			if (*network).IPAMConfig != nil {
				settings.IPAMConfig.IPv4Address = (*network).IPAMConfig.IPv4Address
				settings.IPAMConfig.IPv6Address = (*network).IPAMConfig.IPv6Address
				settings.IPAMConfig.LinkLocalIPs = (*network).IPAMConfig.LinkLocalIPs
			}

			parsed.Network.Networks[k] = settings
		}

		parsed.Network.Ports = inspect.NetworkSettings.Ports
		parsed.Network.Gateway = inspect.NetworkSettings.Gateway
		parsed.Network.IPAddress = inspect.NetworkSettings.IPAddress
		parsed.Network.IPPrefixLen = inspect.NetworkSettings.IPPrefixLen
		parsed.Network.IPv6Gateway = inspect.NetworkSettings.IPv6Gateway
		parsed.Network.MacAddress = inspect.NetworkSettings.MacAddress
	}

	if inspect.HostConfig != nil {
		parsed.Volumes.VolumeDriver = inspect.HostConfig.VolumeDriver
		parsed.Volumes.VolumesFrom = inspect.HostConfig.VolumesFrom

		parsed.Volumes.Mounts = make([]InspectVolumeMount, len(inspect.HostConfig.Mounts))
		for k, mount := range inspect.HostConfig.Mounts {
			parsed.Volumes.Mounts[k].Type = Type(mount.Type)
			parsed.Volumes.Mounts[k].Source = mount.Source
			parsed.Volumes.Mounts[k].Target = mount.Target
			parsed.Volumes.Mounts[k].ReadOnly = mount.ReadOnly
		}
	}

	if inspect.Config != nil {
		parsed.Host.ExposedPorts = inspect.Config.ExposedPorts
	}

	if inspect.HostConfig != nil {
		parsed.Host.Binds = inspect.HostConfig.Binds
		parsed.Host.NetworkMode = string(inspect.HostConfig.NetworkMode)
		parsed.Host.IsDefault = inspect.HostConfig.NetworkMode.IsDefault()
		parsed.Host.IsHost = inspect.HostConfig.NetworkMode.IsHost()
		parsed.Host.IsContainer = inspect.HostConfig.NetworkMode.IsContainer()
		parsed.Host.IsPrivate = inspect.HostConfig.NetworkMode.IsPrivate()
		parsed.Host.IsNone = inspect.HostConfig.NetworkMode.IsNone()
		parsed.Host.IsBridge = inspect.HostConfig.NetworkMode.IsBridge()
		parsed.Host.IsUserDefined = inspect.HostConfig.NetworkMode.IsUserDefined()
		parsed.Host.UserDefined = inspect.HostConfig.NetworkMode.UserDefined()

		parsed.Host.PortBindings = inspect.HostConfig.PortBindings
		parsed.Host.PortRestartPolicyName = inspect.HostConfig.RestartPolicy.Name
		parsed.Host.PortRestartPolicyMaximumRetryCount = inspect.HostConfig.RestartPolicy.MaximumRetryCount
	}

	if inspect.ContainerJSONBase != nil {
		parsed.ID = inspect.ContainerJSONBase.ID
		parsed.Name = inspect.ContainerJSONBase.Name
		parsed.Created, err = time.Parse(time.RFC3339, inspect.ContainerJSONBase.Created)
		if err != nil {
			return
		}
		parsed.Args = inspect.ContainerJSONBase.Args
		parsed.ImageId = inspect.ContainerJSONBase.Image
		parsed.RestartCount = inspect.ContainerJSONBase.RestartCount

		if inspect.ContainerJSONBase.State != nil {
			parsed.State.Status = inspect.ContainerJSONBase.State.Status
			parsed.State.Running = inspect.ContainerJSONBase.State.Running
			parsed.State.Paused = inspect.ContainerJSONBase.State.Paused
			parsed.State.Restarting = inspect.ContainerJSONBase.State.Restarting
			parsed.State.OOMKilled = inspect.ContainerJSONBase.State.OOMKilled
			parsed.State.Dead = inspect.ContainerJSONBase.State.Dead
			parsed.State.ExitCode = inspect.ContainerJSONBase.State.ExitCode
			parsed.State.StartedAt, err = time.Parse(time.RFC3339, inspect.ContainerJSONBase.State.StartedAt)
			if err != nil {
				return
			}
			parsed.State.FinishedAt, err = time.Parse(time.RFC3339, inspect.ContainerJSONBase.State.FinishedAt)
			if err != nil {
				return
			}
		}
	}

	return
}
