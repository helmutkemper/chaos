package docker

// SetContainerRestartPolicyNo
//
// English:
//
//	Do not automatically restart the container. (the default)
//
// Português:
//
//	Define a política de reinício do container como não reiniciar o container (padrão).
func (e *ContainerBuilder) SetContainerRestartPolicyNo() {
	e.restartPolicy = KRestartPolicyNo
}
