package docker

// SetContainerRestartPolicyUnlessStopped
//
// English:
//
//	Similar to always, except that when the container is stopped (manually or otherwise), it is not
//	restarted even after Docker daemon restarts.
//
// Português:
//
//	Define a política de reinício do container como sempre reinicia o container, caso ele não tenha
//	sido parado manualmente.
func (e *ContainerBuilder) SetContainerRestartPolicyUnlessStopped() {
	e.restartPolicy = KRestartPolicyUnlessStopped
}
