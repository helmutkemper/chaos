package docker

import "os"

// SetCsvLogPath
//
// English:
//
//	Defines the log file path, in CSV format, containing container usage statistics.
//
//	 Input:
//	   path: Log file path.
//	   removeOldFile: true deletes the file if it exists; false adds more records to the existing
//	     file.
//
// Note:
//
//   - This function must be used in conjunction with the StartMonitor() and StopMonitor() functions;
//   - The data columns populated varies by operating system;
//   - See the SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog(),
//     AddFilterToCvsLogWithReplace(), SetCsvFileValueSeparator() and SetCsvFileRowSeparator() functions
//     to change some log settings.
//
// Português:
//
//	Define o caminho do arquivo de log, em formato CSV, contendo estatísticas de uso do container.
//
//	 Entrada:
//	   path: Caminho do arquivo de log.
//	   removeOldFile: true apaga o arquivo caso o mesmo exista; false adiciona mais registros ao arquivo existente.
//
// Nota:
//
//   - Esta função deve ser usada em conjunto com as funções StartMonitor() e StopMonitor();
//   - As colunas de dados preenchidos varia de acordo com o sistema operacional;
//   - Veja as funções SetCsvFileReader(), SetCsvFileRowsToPrint(), AddFilterToCvsLog(),
//     AddFilterToCvsLogWithReplace(), SetCsvFileValueSeparator() e SetCsvFileRowSeparator() para alterar
//     algumas configurações do log.
func (e *ContainerBuilder) SetCsvLogPath(path string, removeOldFile bool) {

	if removeOldFile == true {
		_ = os.Remove(path)
	}

	e.chaos.logPath = path
}
