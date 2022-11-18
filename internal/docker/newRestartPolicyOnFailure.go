package iotmakerdocker

// NewRestartPolicyOnFailureRestart (English): Container restart policy, on failure
//
// NewRestartPolicyOnFailureRestart (Português): Política de reinício do container,
// quando falhar
func NewRestartPolicyOnFailureRestart() RestartPolicy {
	return KRestartPolicyOnFailure
}
