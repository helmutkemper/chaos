package builder

import "github.com/docker/docker/api/types/container"

// GetConfig
//
// Return the container config pointer
func (el DockerSystem) GetConfig() (config *container.Config) {
	return el.Config
}
