package docker

// SetImageBuildOptionsCPUShares
//
// English:
//
//	Set the CPU shares of the image build options.
//
//	 Input:
//	   value: CPU shares (Default: 1024)
//
//	Set this flag to a value greater or less than the default of 1024 to increase or reduce the
//	container’s weight, and give it access to a greater or lesser proportion of the host machine’s
//	CPU cycles.
//
//	This is only enforced when CPU cycles are constrained.
//
//	When plenty of CPU cycles are available, all containers use as much CPU as they need.
//
//	In that way, this is a soft limit. --cpu-shares does not prevent containers from being scheduled
//	in swarm mode.
//
//	It prioritizes container CPU resources for the available CPU cycles.
//
//	It does not guarantee or reserve any specific CPU access.
//
// Português:
//
//	Define o compartilhamento de CPU na construção da imagem.
//
//	 Entrada:
//	   value: Compartilhamento de CPU (Default: 1024)
//
//	Defina este sinalizador para um valor maior ou menor que o padrão de 1024 para aumentar ou reduzir
//	o peso do container e dar a ele acesso a uma proporção maior ou menor dos ciclos de CPU da máquina
//	host.
//
//	Isso só é aplicado quando os ciclos da CPU são restritos. Quando muitos ciclos de CPU estão
//	disponíveis, todos os container usam a quantidade de CPU de que precisam. Dessa forma, este é um
//	limite flexível. --cpu-shares não impede que os containers sejam agendados no modo swarm.
//
//	Ele prioriza os recursos da CPU do container para os ciclos de CPU disponíveis.
//
//	Não garante ou reserva nenhum acesso específico à CPU.
func (e *ContainerBuilder) SetImageBuildOptionsCPUShares(value int64) {
	e.buildOptions.CPUShares = value
}
