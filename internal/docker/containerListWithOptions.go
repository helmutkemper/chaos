package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// ContainerListWithOptions (English): List containers
//
//	quiet: bool
//	size: bool populate types.Container.SiseRw and types.Container.SizeRootFs
//	all: bool false for running containers only
//	latest: bool
//	since: string example: "2020-09-08T00:39:53.613203298Z"
//	before: string example: "2020-09-08T00:39:53.613203298Z"
//	limit: int
//
// ContainerListWithOptions (Português): Lista containers
//
//	quiet: bool
//	size: bool popula types.Container.SiseRw e types.Container.SizeRootFs
//	all: bool false retorna apenas containers rodando
//	latest: bool
//	since: string exemplo: "2020-09-08T00:39:53.613203298Z"
//	before: string exemplo: "2020-09-08T00:39:53.613203298Z"
//	limit: int
func (el *DockerSystem) ContainerListWithOptions(
	quiet bool,
	// English: populate types.Container.SiseRw and SizeRootFs
	size bool,
	all bool,
	latest bool,

	// English: example: "2020-09-08T00:39:53.613203298Z"
	since string,

	// English: example: "2020-09-08T00:39:53.613203298Z"
	before string,
	limit int,
	filters filters.Args,
) (
	list []types.Container,
	err error,
) {

	list, err = el.cli.ContainerList(
		el.ctx,
		types.ContainerListOptions{
			Quiet:   quiet,
			Size:    size,
			All:     all,
			Latest:  latest,
			Since:   since,
			Before:  before,
			Limit:   limit,
			Filters: filters,
		},
	)

	return
}
