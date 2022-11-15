package docker

// SetCsvFileValueSeparator
//
// English:
//
//	Defines the column separator of the log file, in CSV format, containing container usage
//	statistics.
//
//	 Input:
//	   value: CSV file column separator (default value: ",")
//
// Note:
//
//   - This function is used in conjunction with the SetCsvLogPath(), StartMonitor(), StopMonitor(),
//     SetCsvFileRowSeparator(), SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog() and
//     AddFilterToCvsLogWithReplace() functions;
//   - The data columns populated varies by operating system.
//
// Português:
//
//	Define o separador de coluna do arquivo de log, em formato CSV, contendo estatísticas de uso do
//	container.
//
//	 Entrada:
//	   value: separador de coluna do arquivo CSV (valor padrão: ",")
//
// Nota:
//
//   - Esta função é usada em conjunto com as funções SetCsvLogPath(), StartMonitor(), StopMonitor(),
//     SetCsvFileRowSeparator(), SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog() e
//     AddFilterToCvsLogWithReplace();
//   - As colunas de dados preenchidos varia de acordo com o sistema operacional.
func (e *ContainerBuilder) SetCsvFileValueSeparator(value string) {
	e.csvValueSeparator = value
}
