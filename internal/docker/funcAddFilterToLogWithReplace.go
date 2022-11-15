package docker

// AddFilterToCvsLogWithReplace
//
// Similar:
//
//	AddFilterToCvsLogWithReplace(), AddFilterToCvsLog()
//
// English:
//
//	Adds a filter to search and convert a textual value to a column in the CSV log file.
//
//	 Input:
//	   label: Value to be placed in the log file column.
//	   match: Simple text searched in the container's standard output to activate the filter
//	   filter: Regular expression used to filter what goes into the log using the `valueToGet`
//	     parameter.
//	   search: Regular expression used for search and replacement in the text found in the previous
//	     step [optional].
//	   replace: Regular expression replace element [optional].
//
// Note:
//
//   - This function is used in conjunction with SetCsvLogPath(), StartMonitor(), StopMonitor().
//
// Português:
//
//	Adiciona um filtro para procurar e converter um valor textual em uma coluna no arquivo de log CSV.
//
//	 Entrada:
//	   label: Valor do rótulo a ser colocado na coluna do arquivo de log.
//	   match: Texto simples procurado na saída padrão do container para ativar o filtro
//	   filter: Expressão regular usada para filtrar o que vai para o log usando o parâmetro
//	     `valueToGet`.
//	   search: Expressão regular usada para busca e substituição no texto encontrado na etapa anterior
//	     [opcional].
//	   replace: Elemento da troca da expressão regular [opcional].
//
// Nota:
//
//   - Esta função é usada em conjunto com SetCsvLogPath(), StartMonitor(), StopMonitor()
func (e *ContainerBuilder) AddFilterToCvsLogWithReplace(label, match, filter, search, replace string) {
	if e.chaos.filterLog == nil {
		e.chaos.filterLog = make([]LogFilter, 0)
	}

	e.chaos.filterLog = append(
		e.chaos.filterLog,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  search,
			Replace: replace,
		},
	)
}
