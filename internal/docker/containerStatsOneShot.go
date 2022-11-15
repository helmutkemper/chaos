package docker

import (
	"encoding/json"
	"github.com/docker/docker/api/types"
	"io/ioutil"
)

// ContainerStatisticsOneShot (English): Returns the performance information of the
// container in a timely manner
//
//	id: string container id
//
// ContainerStatisticsOneShot (Português): Retorna as informações de desempenho do
// container de forma pontual
//
//	id: string container id
func (el *DockerSystem) ContainerStatisticsOneShot(
	id string,
) (
	statsRet types.Stats,
	err error,
) {

	var stats types.ContainerStats
	var body []byte

	stats, err = el.cli.ContainerStats(el.ctx, id, false)
	if err != nil {
		return
	}

	body, err = ioutil.ReadAll(stats.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &statsRet)
	return
}
