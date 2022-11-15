package docker

import (
	"bytes"
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/util"
	"io"
	"io/ioutil"
	"path/filepath"
)

// FindDockerFile (English): Find dockerfile in folder tree.
//
//	Priority order: './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//	'.*Dockerfile.*', '.*dockerfile.*'
//
// FindDockerFile (Português): Procura pelo arquivo dockerfile na árvore de diretórios.
//
//	Ordem de prioridade: './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//	'.*Dockerfile.*', '.*dockerfile.*'
func (el *DockerSystem) FindDockerFile(folderPath string) (fullPathInsideTarFile string, err error) {
	var fileExists bool

	folderPath, err = filepath.Abs(folderPath)
	if err != nil {
		return
	}

	_, err = ioutil.ReadDir(folderPath)
	if err != nil {
		return
	}

	fileExists = util.VerifyFileExists(folderPath + "/Dockerfile-iotmaker")
	if fileExists == true {
		fullPathInsideTarFile = "/Dockerfile-iotmaker"
		return
	}

	fileExists = util.VerifyFileExists(folderPath + "/Dockerfile")
	if fileExists == true {
		fullPathInsideTarFile = "/Dockerfile"
		return
	}

	fileExists = util.VerifyFileExists(folderPath + "/dockerfile")
	if fileExists == true {
		fullPathInsideTarFile = "/dockerfile"
		return
	}

	fullPathInsideTarFile, err = util.FileFindHasPrefixRecursively("Dockerfile-iotmaker")
	if err == nil {
		return
	}

	fullPathInsideTarFile, err = util.FileFindHasPrefixRecursively("dockerfile-iotmaker")
	if err == nil {
		return
	}

	fullPathInsideTarFile, err = util.FileFindHasPrefixRecursively("Dockerfile")
	if err == nil {
		return
	}

	fullPathInsideTarFile, err = util.FileFindHasPrefixRecursively("dockerfile")
	if err == nil {
		return
	}

	fullPathInsideTarFile, err = util.FileFindContainsRecursively("Dockerfile")
	if err == nil {
		return
	}

	fullPathInsideTarFile, err = util.FileFindContainsRecursively("dockerfile")

	return
}

// ImageBuildFromFolder (English): Make a image from folder path content
//
//	folderPath: string absolute folder path
//	tags: []string image tags
//	channel: *chan channel of pull/build data
//
//	  Note: dockerfile priority order: './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*',
//	  'dockerfile.*', '.*Dockerfile.*', '.*dockerfile.*'
//
// ImageBuildFromFolder (Português): Monta uma imagem a partir de um diretório
//
//	folderPath: string caminho absoluto do diretório
//	tags: []string tags da imagem
//	channel: *chan channel com dados do pull/build da imagem
//
//	  Nota: ordem de prioridade do dockerfile: './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*',
//	  'dockerfile.*', '.*Dockerfile.*', '.*dockerfile.*'
func (el *DockerSystem) ImageBuildFromFolder(
	folderPath string,
	imageName string,
	tags []string,
	imageBuildOptions types.ImageBuildOptions,
	channel *chan ContainerPullStatusSendToChannel,
) (
	imageID string,
	err error,
) {

	var tarFileReader *bytes.Reader
	var reader io.Reader
	var dockerFilePath string
	var dockerFileName string

	if len(tags) == 0 {
		tags = []string{
			imageName,
		}
	} else {
		tags = append(tags, imageName)
	}

	tarFileReader, err = el.ImageBuildPrepareFolderContext(folderPath)
	if err != err {
		return
	}

	if len(imageBuildOptions.Tags) == 0 {
		imageBuildOptions.Tags = tags
	} else {
		imageBuildOptions.Tags = append(imageBuildOptions.Tags, tags...)
	}

	imageBuildOptions.Remove = true
	if imageBuildOptions.Dockerfile == "" {
		dockerFilePath, err = el.FindDockerFile(folderPath)
		if err != nil {
			return
		}

		dockerFileName = filepath.Base(dockerFilePath)
		imageBuildOptions.Dockerfile = dockerFileName
	}

	reader, err = el.ImageBuild(tarFileReader, imageBuildOptions)
	if err != nil {
		return
	}

	var successfully bool
	successfully, err = el.processBuildAndPullReaders(&reader, channel)
	if successfully == false || err != nil {
		if err != nil {
			return
		}

		err = errors.New("image build error")
		return
	}

	imageID, err = el.ImageFindIdByName(imageBuildOptions.Tags[0])
	if err != nil {
		return
	}

	return
}
