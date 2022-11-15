package docker

// NewRestartPolicyRestartPolicyNoRestart (English): Container restart policy, no restart
//
// NewRestartPolicyRestartPolicyNoRestart (Português): Política de reinício do container,
// não reiniciar
func NewRestartPolicyRestartPolicyNoRestart() RestartPolicy {
	return KRestartPolicyNo
}
