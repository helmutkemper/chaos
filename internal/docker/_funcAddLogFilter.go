package docker

// AddLogFilter
//
// English: Adds a filter to the container's standard output to search for a textual value and populate the container's log file.
//   Input:
//     label: Column label in log file
//     match: Simple text searched in the container's standard output to activate the filter
//     filter: Regular expression used to filter what goes into the log using the `valueToGet` parameter.
//     search: Regular expression used for search and replacement in the text found in the previous step [optional].
//     replace: Regular expression replace element [optional].
//
// Português: Adiciona um filtro na saída padrão do container para procurar um valor textual e preencher o arquivo de log do container.
//   Entrada:
//     label: Rótulo da coluna no arquivo de log
//     match: Texto simples procurado na saída padrão do container para ativar o filtro
//     filter: Expressão regular usada para filtrar o que vai para o log usando o parâmetro `valueToGet`.
//     search: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior [opcional].
//     replace: Elemento da troca da expressão regular [opcional].
func (e *ContainerBuilder) AddLogFilter(label, match, filter, search, replace string) {
	if e.chaos.filterFail == nil {
		e.chaos.filterFail = make([]LogFilter, 0)
	}

	e.chaos.filterFail = append(
		e.chaos.filterFail,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
