package docker

// AddFilterToCvsLog
//
// Similar:
//
//	AddFilterToCvsLogWithReplace(), AddFilterToCvsLog()
//
// English:
//
//	Adds a filter to search and convert a textual value to a column in the log file.
//
//	 Input:
//	   label: Value to be placed in the log file column.
//	   match: Simple text searched in the container's standard output to activate the filter
//	   filter: Regular expression used to filter what goes into the log using the `valueToGet`
//	     parameter.
//
// Note:
//
//   - This function is used in conjunction with SetCsvLogPath(), StartMonitor(), StopMonitor().
//
// Português:
//
//	Adiciona um filtro para procurar e converter um valor textual em uma coluna no arquivo de log.
//
//	 Entrada:
//	   label: Valor do rótulo a ser colocado na coluna do arquivo de log.
//	   match: Texto simples procurado na saída padrão do container para ativar o filtro
//	   filter: Expressão regular usada para filtrar o que vai para o log usando o parâmetro
//	     `valueToGet`.
//
// Nota:
//
//   - Esta função é usada em conjunto com SetCsvLogPath(), StartMonitor(), StopMonitor()
func (e *ContainerBuilder) AddFilterToCvsLog(label, match, filter string) {
	if e.chaos.filterLog == nil {
		e.chaos.filterLog = make([]LogFilter, 0)
	}

	e.chaos.filterLog = append(
		e.chaos.filterLog,
		LogFilter{
			Label:   label,
			Match:   match,
			Filter:  filter,
			Search:  "",
			Replace: "",
		},
	)
}
