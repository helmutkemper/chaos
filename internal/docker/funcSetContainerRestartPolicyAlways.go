package docker

// SetContainerRestartPolicyAlways
//
// English:
//
//	Always restart the container if it stops. If it is manually stopped, it is restarted only when
//	Docker daemon restarts or the container itself is manually restarted.
//
// Português:
//
//	Define a política de reinício do container como sempre reinicia o container quando ele para, mesmo
//	quando ele é parado manualmente.
func (e *ContainerBuilder) SetContainerRestartPolicyAlways() {
	e.restartPolicy = KRestartPolicyAlways
}
