package docker

// NewKRestartPolicyAlwaysRestart (English): Container restart policy, always restart
//
// NewKRestartPolicyAlwaysRestart (Português): Política de reinício do container, sempre
// reinicia
func NewKRestartPolicyAlwaysRestart() RestartPolicy {
	return KRestartPolicyOnFailure
}
