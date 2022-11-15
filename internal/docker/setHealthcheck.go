package docker

import (
	"github.com/docker/docker/api/types/container"
)

// SetHealthcheck (English): The HEALTHCHECK instruction has two forms:
//
//	HEALTHCHECK [OPTIONS] CMD command (check container health by running a command inside
//	  the container)
//	HEALTHCHECK NONE (disable any healthcheck inherited from the base image)
//
// The HEALTHCHECK instruction tells Docker how to test a container to check that it is
// still working. This can detect cases such as a web server that is stuck in an infinite
// loop and unable to handle new connections, even though the server process is still
// running.
//
// When a container has a healthcheck specified, it has a health status in addition to its
// normal status. This status is initially starting. Whenever a health check passes, it
// becomes healthy (whatever state it was previously in). After a certain number of
// consecutive failures, it becomes unhealthy.
//
// The options that can appear before CMD are:
//
//	--interval=DURATION (default: 30s)
//	--timeout=DURATION (default: 30s)
//	--start-period=DURATION (default: 0s)
//	--retries=N (default: 3)
//
// The health check will first run interval seconds after the container is started, and
// then again interval seconds after each previous check completes.
//
// If a single run of the check takes longer than timeout seconds then the check is
// considered to have failed.
//
// It takes retries consecutive failures of the health check for the container to be
// considered unhealthy.
//
// start period provides initialization time for containers that need time to bootstrap.
// Probe failure during that period will not be counted towards the maximum number of
// retries. However, if a health check succeeds during the start period, the container is
// considered started and all consecutive failures will be counted towards the maximum
// number of retries.
//
// There can only be one HEALTHCHECK instruction in a Dockerfile. If you list more than
// one then only the last HEALTHCHECK will take effect.
//
// The command after the CMD keyword can be either a shell command (e.g. HEALTHCHECK CMD
// /bin/check-running) or an exec array (as with other Dockerfile commands; see e.g.
// ENTRYPOINT for details).
//
// The command’s exit status indicates the health status of the container. The possible
// values are:
//
//	0: success - the container is healthy and ready for use
//	1: unhealthy - the container is not working correctly
//	2: reserved - do not use this exit code
//
// For example, to check every five minutes or so that a web-server is able to serve the
// site’s main page within three seconds:
//
//	HEALTHCHECK --interval=5m --timeout=3s \
//	  CMD curl -f http://localhost/ || exit 1
//
// To help debug failing probes, any output text (UTF-8 encoded) that the command writes
// on stdout or stderr will be stored in the health status and can be queried with docker
// inspect. Such output should be kept short (only the first 4096 bytes are stored
// currently).
//
// When the health status of a container changes, a health_status event is generated with
// the new status.
//
// https://docs.docker.com/engine/reference/builder/#healthcheck
func (el *DockerSystem) SetHealthcheck(healthcheck HealthConfig) {
	var check = container.HealthConfig(healthcheck)
	el.healthcheck = &check
}
