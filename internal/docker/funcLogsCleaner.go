package docker

import (
	"bytes"
)

// logsCleaner
//
// English:
//
//	Clear blank lines of the container's standard output
//
//	Input:
//	  logs: container's standard output
//
//	Output:
//	  logsLine: List of lines of the container's standard output
//
// Português:
//
//	Limpa as linhas em branco da saída padrão do container
//
//	 Entrada:
//	   logs: saída padrão do container
//
//	 Saída:
//	   logsLine: Lista de linhas da saída padrão do container
func (e *ContainerBuilder) logsCleaner(logs []byte) (logsLine [][]byte) {

	size := len(logs)

	// faz o log só lê a parte mais recente do mesmo
	logs = logs[e.logsLastSize:]
	e.logsLastSize = size

	logs = bytes.ReplaceAll(logs, []byte("\r"), []byte(""))
	return bytes.Split(logs, []byte("\n"))
}
