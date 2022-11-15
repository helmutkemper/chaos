package docker

// AddStartChaosMatchFlag
//
// Similar:
//
//	AddStartChaosMatchFlag(), AddStartChaosMatchFlagToFileLog(), AddFilterToStartChaos()
//
// English:
//
//	Adds a filter to the container's standard output to look for a textual value releasing the start of the chaos test.
//
//	 Input:
//	   value: Error text
//
// Português:
//
//	Adiciona um filtro na saída padrão do container para procurar um valor textual liberando o início do teste de caos.
//
//	 Entrada:
//	   value: Texto indicativo de erro
func (e *ContainerBuilder) AddStartChaosMatchFlag(value string) {
	if e.chaos.filterFail == nil {
		e.chaos.filterFail = make([]LogFilter, 0)
	}

	e.chaos.filterFail = append(e.chaos.filterFail, LogFilter{Match: value})
}
