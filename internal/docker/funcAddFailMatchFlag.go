package docker

// AddFailMatchFlag
//
// Similar:
//
//	AddFailMatchFlag(), AddFailMatchFlagToFileLog(), AddFilterToFail()
//
// English:
//
//	Error text searched for in the container's standard output.
//
//	 Input:
//	   value: Error text
//
// Português:
//
//	Texto indicativo de erro procurado na saída padrão do container.
//
//	 Entrada:
//	   value: Texto indicativo de erro
func (e *ContainerBuilder) AddFailMatchFlag(value string) {
	if e.chaos.filterFail == nil {
		e.chaos.filterFail = make([]LogFilter, 0)
	}

	e.chaos.filterFail = append(e.chaos.filterFail, LogFilter{Match: value})
}
