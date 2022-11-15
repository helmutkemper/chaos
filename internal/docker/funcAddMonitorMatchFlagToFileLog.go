package docker

import (
	"github.com/helmutkemper/util"
	"path/filepath"
	"strings"
)

// AddMonitorMatchFlagToFileLog
//
// English:
//
//	Looks for text in the container's standard output and saves it to a log file on the host
//	computer
//
//	 Input:
//	   value: text
//	   logDirectoryPath: File path where the container's standard output filed in a `log.N.log` file
//	     will be saved, where N is an automatically incremented number. e.g.: "./bug/critical/"
//
//	 Output:
//	   err: Default error object
//
// Português:
//
//	Procura por um texto indicativo na saída padrão do container e o salva em um arquivo de log no computador hospedeiro
//
//	 Entrada:
//	   value: Texto indicativo
//	   logDirectoryPath: Caminho do arquivo onde será salva a saída padrão do container arquivada em
//	     um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
//	     Ex.: "./bug/critical/"
//
//	 Output:
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) AddMonitorMatchFlagToFileLog(value, logDirectoryPath string) (err error) {
	if e.chaos.filterMonitor == nil {
		e.chaos.filterMonitor = make([]LogFilter, 0)
	}

	if strings.HasPrefix(logDirectoryPath, string(filepath.Separator)) == false {
		logDirectoryPath += string(filepath.Separator)
	}

	err = util.DirMake(logDirectoryPath)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.chaos.filterMonitor = append(e.chaos.filterMonitor, LogFilter{Match: value, LogPath: logDirectoryPath})

	return
}
