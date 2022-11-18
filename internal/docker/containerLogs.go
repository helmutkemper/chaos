package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"io"
	"io/ioutil"
)

// ContainerLogs (English): Returns container std out
//
// ContainerLogs (Português): Retorna a saída padrão do container
func (el *DockerSystem) ContainerLogs(
	id string,
) (
	log []byte,
	err error,
) {

	var reader io.ReadCloser

	reader, err = el.cli.ContainerLogs(el.ctx, id, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Timestamps: true,
		Follow:     false,
		Details:    false,
	})
	if err != nil {
		return
	}

	log, err = ioutil.ReadAll(reader)

	return
}
