package builder

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"time"
)

const (
	kWaitTextLoopSleep = 500 * time.Millisecond
)

type DockerSystem struct {
	cli              *client.Client
	ctx              context.Context
	networkId        map[string]string
	imageId          map[string]string
	ContainerName    string
	container        map[string]container.ContainerCreateCreatedBody
	networkGenerator map[string]*NextNetworkAutoConfiguration
	healthcheck      *container.HealthConfig
	Config           *container.Config
}
