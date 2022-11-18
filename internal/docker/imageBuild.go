package iotmakerdocker

import (
	"github.com/docker/docker/api/types"
	"io"
)

// ImageBuild (English): Image build from reader. Please, see
// ImageBuildFromFolder(folderPath string, tags []string) and
// ImageBuildFromRemoteServer(server string, tags []string)
//
//	dockerFileTarReader: io.Reader reader from image
//	imageBuildOptions: types.ImageBuildOptions image build options
//
// ImageBuild (Português): Monta uma imagem baseada no header. Por favor, veja,
// ImageBuildFromFolder(folderPath string, tags []string) e
// ImageBuildFromRemoteServer(server string, tags []string)
//
//	dockerFileTarReader: io.Reader reader from image
//	imageBuildOptions: types.ImageBuildOptions configurações da criação da imagem
func (el *DockerSystem) ImageBuild(
	dockerFileTarReader io.Reader,
	imageBuildOptions types.ImageBuildOptions,
) (
	reader io.ReadCloser,
	err error,
) {

	var response types.ImageBuildResponse

	response, err = el.cli.ImageBuild(el.ctx, dockerFileTarReader, imageBuildOptions)
	if err != nil {
		return
	}

	reader = response.Body

	return
}
