package docker

import (
	"github.com/helmutkemper/util"
	"path/filepath"
	"strings"
)

// AddFailMatchFlagToFileLog
//
// Similar:
//
//	AddFailMatchFlag(), AddFailMatchFlagToFileLog(), AddFilterToFail()
//
// English:
//
//	Looks for error text in the container's standard output and saves it to a log file on the host
//	computer
//
//	 Input:
//	   value: Error text
//	   logDirectoryPath: File path where the container's standard output filed in a `log.N.log` file
//	     will be saved, where N is an automatically incremented number. e.g.: "./bug/critical/"
//
//	 Output:
//	   err: Default error object
//
// Português:
//
//	Procura por um texto indicativo de erro na saída padrão do container e o salva em um arquivo de
//	log no computador hospedeiro
//
//	 Entrada:
//	   value: Texto indicativo de erro
//	   logDirectoryPath: Caminho do arquivo onde será salva a saída padrão do container arquivada em
//	     um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
//	     Ex.: "./bug/critical/"
//
//	 Output:
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) AddFailMatchFlagToFileLog(value, logDirectoryPath string) (err error) {
	if e.chaos.filterFail == nil {
		e.chaos.filterFail = make([]LogFilter, 0)
	}

	if strings.HasPrefix(logDirectoryPath, string(filepath.Separator)) == false {
		logDirectoryPath += string(filepath.Separator)
	}

	err = util.DirMake(logDirectoryPath)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.chaos.filterFail = append(e.chaos.filterFail, LogFilter{Match: value, LogPath: logDirectoryPath})

	return
}
