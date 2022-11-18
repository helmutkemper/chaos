package builder

// NewRestartPolicyRestartPolicyUnlessStopped (English): Container restart policy, unless
// stopped
//
// NewRestartPolicyRestartPolicyUnlessStopped (Português): Política de reinício do
// container, a menos que seja interrompido
func NewRestartPolicyRestartPolicyUnlessStopped() RestartPolicy {
	return KRestartPolicyUnlessStopped
}
