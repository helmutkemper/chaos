package docker

// SetImageBuildOptionsCPUSetCPUs
//
// English:
//
//	Limit the specific CPUs or cores a container can use.
//
//	 Input:
//	   value: string with the format "1,2,3"
//
//	A comma-separated list or hyphen-separated range of CPUs a container can use, if you have more
//	than one CPU.
//
// The first CPU is numbered 0.
//
//	A valid value might be 0-3 (to use the first, second, third, and fourth CPU) or 1,3 (to use the
//	second and fourth CPU).
//
// Português:
//
//	Limite a quantidade de CPUs ou núcleos específicos que um container pode usar.
//
//	 Entrada:
//	   value: string com o formato "1,2,3"
//
//	Uma lista separada por vírgulas ou intervalo separado por hífen de CPUs que um container pode
//	usar, se você tiver mais de uma CPU.
//
//	A primeira CPU é numerada como 0.
//
//	Um valor válido pode ser 0-3 (para usar a primeira, segunda, terceira e quarta CPU) ou 1,3 (para
//	usar a segunda e a quarta CPU).
func (e *ContainerBuilder) SetImageBuildOptionsCPUSetCPUs(value string) {
	e.buildOptions.CPUSetCPUs = value

	e.addProblem("The SetImageBuildOptionsCPUSetCPUs() function can generate an error when building the image.")
}
