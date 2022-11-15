package docker

import (
	"github.com/helmutkemper/util"
	"path/filepath"
	"strings"
)

// AddRestartMatchFlagToFileLog
//
// Similar:
//
//	AddFilterToRestartContainer(), AddRestartMatchFlag(), AddRestartMatchFlagToFileLog()
//
// English:
//
//	Adds a filter to the standard output of the container to look for a textual value releasing the
//	possibility of the container being restarted during the chaos test.
//
//	 Input:
//	   value: Simple text searched in the container's standard output to activate the filter
//	   logDirectoryPath: File path where the container's standard output filed in a `log.N.log` file
//	     will be saved, where N is an automatically incremented number. e.g.: "./bug/critical/"
//
// Note:
//
//   - Chaos testing is a test performed when there is a need to simulate failures of the
//     microservices involved in the project.
//     During chaos testing, the container can be paused, to simulate a container not responding due
//     to overload, or stopped and restarted, simulating a critical crash, where a microservice was
//     restarted after an unresponsive time.
//
// Português:
//
//	Adiciona um filtro na saída padrão do container para procurar um valor textual liberando a
//	possibilidade do container ser reinicado durante o teste de caos.
//
//	 Entrada:
//	   value: Texto simples procurado na saída padrão do container para ativar o filtro
//	   logDirectoryPath: Caminho do arquivo onde será salva a saída padrão do container arquivada em
//	     um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
//	     Ex.: "./bug/critical/"
//
// Nota:
//
//   - Teste de caos é um teste feito quando há a necessidade de simular falhas dos microsserviços
//     envolvidos no projeto.
//     Durante o teste de caos, o container pode ser pausado, para simular um container não
//     respondendo devido a sobrecarga, ou parado e reiniciado, simulando uma queda crítica, onde um
//     microsserviço foi reinicializado depois de um tempo sem resposta.
func (e *ContainerBuilder) AddRestartMatchFlagToFileLog(value, logDirectoryPath string) (err error) {
	if e.chaos.filterRestart == nil {
		e.chaos.filterRestart = make([]LogFilter, 0)
	}

	if strings.HasPrefix(logDirectoryPath, string(filepath.Separator)) == false {
		logDirectoryPath += string(filepath.Separator)
	}

	err = util.DirMake(logDirectoryPath)
	if err != nil {
		return
	}

	e.chaos.filterRestart = append(e.chaos.filterRestart, LogFilter{Match: value, LogPath: logDirectoryPath})

	return
}
