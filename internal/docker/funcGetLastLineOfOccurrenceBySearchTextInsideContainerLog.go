package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
)

// GetLastLineOfOccurrenceBySearchTextInsideContainerLog
//
// English:
//
//	Returns the last line of output standard output container that contains the text searched
//
//	Input:
//	  value: text to be searched in the standard output of the container
//
//	Output:
//	  text: string with the last line that contains the text searched
//	  contains: true if the text was found
//	  err: default error object
//
// Português:
//
//	Retorna a ultima linha sa saída padrão do container que contém o texto procurado
//
//	 Input:
//	   value: texto procurado na saída padrão do container
//
//	 Saída:
//	   text: string com a última linha que contém o texto procurado
//	   contains: true se o texto foi encontrado
//	   err: objeto de erro padrão
func (e *ContainerBuilder) GetLastLineOfOccurrenceBySearchTextInsideContainerLog(value string) (text string, contains bool, err error) {
	var logs []byte
	var lineList [][]byte
	logs, err = e.GetContainerLog()
	if err != nil {
		util.TraceToLog()
		return
	}

	logs = bytes.ReplaceAll(logs, []byte("\r"), []byte(""))
	lineList = bytes.Split(logs, []byte("\n"))

	for i := len(lineList) - 1; i >= 0; i -= 1 {
		if bytes.Contains(lineList[i], []byte(value)) == true {
			text = string(lineList[i])
			contains = true
			return
		}
	}

	return
}
