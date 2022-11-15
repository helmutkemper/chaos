package docker

func (e *ContainerBuilder) SetSshKeyFileName(value string) {
	e.sshDefaultFileName = value
}
