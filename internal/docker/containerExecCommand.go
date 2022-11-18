package iotmakerdocker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"io/ioutil"
)

type ExecResult struct {
	StdOut   string
	StdErr   string
	ExitCode int
}

func (el *DockerSystem) ContainerExecCommand(
	id string,
	commands []string,
) (
	exitCode int,
	runing bool,
	stdOutput []byte,
	stdError []byte,
	err error,
) {

	var idResponse types.IDResponse
	idResponse, err = el.cli.ContainerExecCreate(
		el.ctx,
		id,
		types.ExecConfig{
			Cmd:          commands,
			Privileged:   true,
			AttachStderr: true,
			AttachStdin:  true,
			AttachStdout: true,
		},
	)
	if err != nil {
		return
	}

	var resp types.HijackedResponse
	resp, err = el.cli.ContainerExecAttach(el.ctx, idResponse.ID, types.ExecStartCheck{})
	if err != nil {
		return
	}
	defer resp.Close()

	stderr := new(bytes.Buffer)

	var i types.ContainerExecInspect
	i, err = el.cli.ContainerExecInspect(el.ctx, idResponse.ID)
	if err != nil {
		return
	}

	stdOutput, err = ioutil.ReadAll(resp.Reader)

	stdError, err = ioutil.ReadAll(stderr)

	exitCode = i.ExitCode
	runing = i.Running

	return
}
