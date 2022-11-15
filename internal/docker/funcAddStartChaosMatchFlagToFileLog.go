package docker

import (
	"github.com/helmutkemper/util"
	"path/filepath"
	"strings"
)

// AddStartChaosMatchFlagToFileLog
//
// Similar:
//
//	AddStartChaosMatchFlag(), AddStartChaosMatchFlagToFileLog(), AddFilterToStartChaos()
//
// English:
//
//	Adds a filter to the container's standard output to look for a textual value releasing the start
//	of the chaos test.
//
//	 Input:
//	   value: Error text
//	   logDirectoryPath: File path where the container's standard output filed in a `log.N.log` file
//	     will be saved, where N is an automatically incremented number.
//	     e.g.: "./bug/critical/"
//
//	 Output:
//	   err: Default error object
//
// Português:
//
//	Adiciona um filtro na saída padrão do container para procurar um valor textual liberando o início
//	do teste de caos.
//
//	 Entrada:
//	   value: Texto indicativo de erro
//	   logDirectoryPath: Caminho do arquivo onde será salva a saída padrão do container arquivada em
//	     um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
//	     Ex.: "./bug/critical/"
//
//	 Output:
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) AddStartChaosMatchFlagToFileLog(value, logDirectoryPath string) (err error) {
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
