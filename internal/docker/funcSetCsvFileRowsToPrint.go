package docker

// SetCsvFileRowsToPrint
//
// English:
//
//	Defines which columns will be printed in the log, in the form of a CSV file, with container
//	performance information, memory consumption indicators and access times.
//
//	 Input:
//	   value: List of columns printed in CSV file. Eg.: KLogColumnMacOs, KLogColumnWindows,
//	     KLogColumnAll or any combination of KLogColumn... concatenated with pipe.
//	     Eg.: KLogColumnReadingTime | KLogColumnCurrentNumberOfOidsInTheCGroup | ...
//
// Nota:
//
//   - To see the complete list of columns, use SetCsvFileRowsToPrint(KLogColumnAll) and
//     SetCsvFileReader(true).
//     This will print the constant names on top of each column in the log.
//
// Português:
//
//	Define quais colunas vão ser impressas no log, na forma de arquivo CSV, com informações de
//	desempenho do container, indicadores de consumo de memória e tempos de acesso.
//
//	 Entrada:
//
//	   value: Lista das colunas impressas no arquivo CSV. Ex.: KLogColumnMacOs, KLogColumnWindows,
//	     KLogColumnAll ou qualquer combinação de KLogColumn... concatenado com pipe.
//	     Ex.: KLogColumnReadingTime | KLogColumnCurrentNumberOfOidsInTheCGroup | ...
//
// Nota:
//
//   - Para vê a lista completa de colunas, use SetCsvFileRowsToPrint(KLogColumnAll) e
//     SetCsvFileReader(true).
//     Isto irá imprimir os nomes das constantes em cima de cada coluna do log.
func (e *ContainerBuilder) SetCsvFileRowsToPrint(value int64) {
	e.rowsToPrint = value
}
