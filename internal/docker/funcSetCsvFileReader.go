package docker

// SetCsvFileReader
//
// English:
//
//	Prints in the header of the file the name of the constant responsible for printing the column in
//	the log.
//
//	 Input:
//	   value: true to print the name of the constant responsible for printing the column in the log
//	     in the header of the file.
//
// Nota:
//
//   - The constants printed in the first line of the file are used in the SetCsvFileRowsToPrint()
//     function. Simply separate the constants by pipe (|).
//     Example: container.SetCsvFileRowsToPrint( KLogColumnReadingTime
//     | KLogColumnCurrentNumberOfOidsInTheCGroup | ... )
//
// Português:
//
//	Imprime no cabeçalho do arquivo o nome da constante responsável por imprimir a coluna no log.
//
//	 Entrada:
//	   value: true para imprimir no cabeçalho do arquivo o nome da constante responsável por imprimir
//	     a coluna no log.
//
// Nota:
//
//   - As constantes impressas na primeira linha do arquivo são usadas na função
//     SetCsvFileRowsToPrint(). Basta separar as contantes por pipe (|).
//     Exemplo: container.SetCsvFileRowsToPrint( KLogColumnReadingTime
//     | KLogColumnCurrentNumberOfOidsInTheCGroup | ... )
func (e *ContainerBuilder) SetCsvFileReader(value bool) {
	e.csvConstHeader = value
}
