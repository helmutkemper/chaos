package docker

// AddFilterToRestartContainer
//
// Similar:
//
//	AddFilterToRestartContainer(), AddRestartMatchFlag(), AddRestartMatchFlagToFileLog()
//
// Português:
//
//	Adiciona um filtro na saída padrão do container para procurar um valor textual liberando a
//	possibilidade do container ser reinicado durante o teste de caos.
//
//	 Entrada:
//	   match: Texto simples procurado na saída padrão do container para ativar o filtro
//	   filter: Expressão regular usada para filtrar o que vai para o log usando o parâmetro
//	     `valueToGet`.
//	   search: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior
//	     [opcional].
//	   replace: Elemento da troca da expressão regular [opcional].
//
// Nota:
//
//   - Teste de caos é um teste feito quando há a necessidade de simular falhas dos microsserviços
//     envolvidos no projeto.
//     Durante o teste de caos, o container pode ser pausado, para simular um container não
//     respondendo devido a sobrecarga, ou parado e reiniciado, simulando uma queda crítica, onde um
//     microsserviço foi reinicializado depois de um tempo sem resposta.
//
// English:
//
//	Adds a filter to the standard output of the container to look for a textual value releasing the
//	possibility of the container being restarted during the chaos test.
//
//	 Input:
//	   match: Simple text searched in the container's standard output to activate the filter
//	   filter: Regular expression used to filter what goes into the log using the `valueToGet`
//	     parameter.
//	   search: Regular expression used for search and replacement in the text found in the previous
//	     step [optional].
//	   replace: Regular expression replace element [optional].
//
// Note:
//
//   - Chaos testing is a test performed when there is a need to simulate failures of the
//     microservices involved in the project.
//     During chaos testing, the container can be paused, to simulate a container not responding due
//     to overload, or stopped and restarted, simulating a critical crash, where a microservice was
//     restarted after an unresponsive time.
func (e *ContainerBuilder) AddFilterToRestartContainer(match, filter, search, replace string) {
	if e.chaos.filterRestart == nil {
		e.chaos.filterRestart = make([]LogFilter, 0)
	}

	e.chaos.filterRestart = append(
		e.chaos.filterRestart,
		LogFilter{
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
