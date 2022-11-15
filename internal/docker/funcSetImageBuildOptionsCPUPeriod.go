package docker

// SetImageBuildOptionsCPUPeriod
//
// English:
//
//	Specify the CPU CFS scheduler period, which is used alongside --cpu-quota.
//
//	 Input:
//	   value: CPU CFS scheduler period
//
//	Defaults to 100000 microseconds (100 milliseconds). Most users do not change this from the
//	default.
//
//	For most use-cases, --cpus is a more convenient alternative.
//
// Português:
//
//	Especifique o período do agendador CFS da CPU, que é usado junto com --cpu-quota.
//
//	 Entrada:
//	   value: período do agendador CFS da CPU
//
//	O padrão é 100.000 microssegundos (100 milissegundos). A maioria dos usuários não altera o padrão.
//
//	Para a maioria dos casos de uso, --cpus é uma alternativa mais conveniente.
func (e *ContainerBuilder) SetImageBuildOptionsCPUPeriod(value int64) {
	e.buildOptions.CPUPeriod = value

	e.addProblem("The SetImageBuildOptionsCPUPeriod() function can generate an error when building the image.")
}
