package builder

type RestartPolicy int

func (el RestartPolicy) String() string {
	return restartPolicies[el]
}

const (
	//Do not automatically restart the container. (the default)
	KRestartPolicyNo RestartPolicy = iota

	//Restart the container if it exits due to an error, which manifests as a non-zero exit
	//code.
	KRestartPolicyOnFailure

	//Always restart the container if it stops. If it is manually stopped, it is restarted
	//only when Docker daemon restarts or the container itself is manually restarted. (See the second bullet listed in restart policy details)
	KRestartPolicyAlways

	//Similar to always, except that when the container is stopped (manually or otherwise),
	//it is not restarted even after Docker daemon restarts.
	KRestartPolicyUnlessStopped
)

var restartPolicies = [...]string{
	"no",
	"on-failure",
	"always",
	"unless-stopped",
}
