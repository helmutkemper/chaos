package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
)

// ContainerStatisticsOneShot
//
// English:
//
//	Returns the container's memory and system consumption data at the time of the query.
//
//	 Output:
//	   stats: Container statistics such as memory, bytes read/written, CPUs, access times, etc.
//	   err: standard error object
//
// Português:
//
//	Retorna os dados de consumo de memória e sistema do container no instante da consulta.
//
//	 Saída:
//	   stats: Estatísticas do conbtainer, como memória, bytes lidos/escritos, CPUs, tempos de acesso,
//	     etc.
//	   err: Objeto de erro padrão
func (e *ContainerBuilder) ContainerStatisticsOneShot() (
	stats types.Stats,
	err error,
) {

	stats, err = e.dockerSys.ContainerStatisticsOneShot(e.containerID)
	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
