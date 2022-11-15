package docker

// AddSuccessMatchFlag
//
// English:
//
//	Adds a text to be searched for in the container's standard output, indicating test success
//
//	 Input:
//	   value: Text searched for in the container's standard output
//
// Português:
//
//	Adiciona um texto a ser procurado na saída padrão do conteiner, indicando sucesso do teste
//
//	 Entrada:
//	   value: Texto procurado na saída padrão do container
func (e *ContainerBuilder) AddSuccessMatchFlag(value string) {
	if e.chaos.filterSuccess == nil {
		e.chaos.filterSuccess = make([]LogFilter, 0)
	}

	e.chaos.filterSuccess = append(e.chaos.filterSuccess, LogFilter{Match: value})
}
