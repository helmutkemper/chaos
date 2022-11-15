package docker

// SetCsvFileRowSeparator
//
// English:
//
// Defines the log file line separator, in CSV format, containing container usage statistics.
//
//	Input:
//	  value: separador de linha do arquivo CSV (valor padrão: "\n")
//
// Nota:
//
//   - Esta função é usada em conjunto com as funções SetCsvLogPath(), StartMonitor(), StopMonitor(),
//     SetCsvFileRowSeparator(), SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog() e
//     AddFilterToCvsLogWithReplace();
//   - As colunas de dados preenchidos varia de acordo com o sistema operacional.
//
// Português:
//
//	Define o separador de linha do arquivo de log, em formato CSV, contendo estatísticas de uso do container.
//
//	 Entrada:
//	   value: separador de linha do arquivo CSV (valor padrão: "\n")
//
// Nota:
//
//   - Esta função é usada em conjunto com as funções SetCsvLogPath(), StartMonitor(), StopMonitor(),
//     SetCsvFileRowSeparator(), SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog() e
//     AddFilterToCvsLogWithReplace();
//   - As colunas de dados preenchidos varia de acordo com o sistema operacional.
func (e *ContainerBuilder) SetCsvFileRowSeparator(value string) {
	e.csvRowSeparator = value
}
