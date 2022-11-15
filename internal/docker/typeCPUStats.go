package docker

// CPUStats aggregates and wraps all CPU related info of container

// CPUStats
//
// English:
//
//	Aggregates and wraps all CPU related info of container
//
// Português:
//
//	Agrega e embrulha todas as informações de CPU do container
type CPUStats struct {
	// CPU Usage. Linux and Windows.
	CPUUsage CPUUsage `json:"cpu_usage"`

	// System Usage. Linux only.
	SystemUsage uint64 `json:"system_cpu_usage,omitempty"`

	// Online CPUs. Linux only.
	OnlineCPUs uint32 `json:"online_cpus,omitempty"`

	// Throttling Data. Linux only.
	ThrottlingData ThrottlingData `json:"throttling_data,omitempty"`
}
