package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
)

// FindTextInsideContainerLog
//
// English:
//
//	Search for text in standard container output.
//
//	 Input:
//	   value: searched text
//
//	 Output:
//	   contains: true if text was found
//	   err: standard error object
//
// Português:
//
//	Procurar por um texto na saída padrão do container.
//
//	 Entrada:
//	   value: texto procurado
//
//	 Saída:
//	   contains: true se o texto foi encontrado
//	   err: objeto de erro padrão
func (e *ContainerBuilder) FindTextInsideContainerLog(value string) (contains bool, err error) {
	var logs []byte
	logs, err = e.GetContainerLog()
	if err != nil {
		util.TraceToLog()
		return
	}

	contains = bytes.Contains(logs, []byte(value))
	return
}
