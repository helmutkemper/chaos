package docker

// SetImageBuildOptionsCPUQuota
//
// English:
//
//	Defines the host machine’s CPU cycles.
//
//	 Input:
//	   value: machine’s CPU cycles. (Default: 1024)
//
//	Set this flag to a value greater or less than the default of 1024 to increase or reduce the
//	container’s weight, and give it access to a greater or lesser proportion of the host machine’s
//	CPU cycles.
//
//	This is only enforced when CPU cycles are constrained. When plenty of CPU cycles are available,
//	all containers use as much CPU as they need. In that way, this is a soft limit. --cpu-shares does
//	not prevent containers from being scheduled in swarm mode. It prioritizes container CPU resources
//	for the available CPU cycles.
//
//	It does not guarantee or reserve any specific CPU access.
//
// Português:
//
//	Define os ciclos de CPU da máquina hospedeira.
//
//	 Entrada:
//	   value: ciclos de CPU da máquina hospedeira. (Default: 1024)
//
//	Defina este flag para um valor maior ou menor que o padrão de 1024 para aumentar ou reduzir o peso
//	do container e dar a ele acesso a uma proporção maior ou menor dos ciclos de CPU da máquina
//	hospedeira.
//
//	Isso só é aplicado quando os ciclos da CPU são restritos. Quando muitos ciclos de CPU estão
//	disponíveis, todos os containeres usam a quantidade de CPU de que precisam. Dessa forma, é um
//	limite flexível. --cpu-shares não impede que os containers sejam agendados no modo swarm. Ele
//	prioriza os recursos da CPU do container para os ciclos de CPU disponíveis.
//
//	Não garante ou reserva nenhum acesso específico à CPU.
func (e *ContainerBuilder) SetImageBuildOptionsCPUQuota(value int64) {
	e.buildOptions.CPUQuota = value

	e.addProblem("The SetImageBuildOptionsCPUQuota() function can generate an error when building the image.")
}
