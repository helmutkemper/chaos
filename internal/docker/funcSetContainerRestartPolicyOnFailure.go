package docker

// SetContainerRestartPolicyOnFailure
//
// English:
//
//	Restart the container if it exits due to an error, which manifests as a non-zero exit code
//
// Português:
//
//	Define a política de reinício do container como reinicia o container se houver um erro (com o
//	manifesto informando um código de erro diferente de zero).
func (e *ContainerBuilder) SetContainerRestartPolicyOnFailure() {
	e.restartPolicy = KRestartPolicyOnFailure
}
