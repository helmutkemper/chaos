package docker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ImageBuildFromFolder
//
// English:
//
//	Transforms the contents of the folder defined in SetBuildFolderPath() into a docker image
//
//	 Output:
//	   inspect: Contém informações sobre a imagem criada
//	   err: standard object error
//
// Note:
//
//   - The folder must contain a dockerfile file, but since different uses can have different
//     dockerfiles, the following order will be given when searching for the file:
//     "Dockerfile-iotmaker", "Dockerfile", "dockerfile" in the root folder;
//   - If not found, a recursive search will be done for "Dockerfile" and "dockerfile";
//   - If the project is in golang and the main.go file, containing the package main, is contained in
//     the root folder, with the go.mod file, the MakeDefaultDockerfileForMe() and
//     MakeDefaultDockerfileForMeWithInstallExtras() functions can be used to make a standard
//     Dockerfile file automatically.
//
// Português:
//
//	Transforma o conteúdo da pasta definida em SetBuildFolderPath() em uma imagem docker
//
//	 Saída:
//	   inspect: contém informações sobre a imagem criada
//	   err: objeto de erro padrão
//
// Nota:
//
//   - A pasta deve conter um arquivo dockerfile, mas, como diferentes usos podem ter diferentes
//     dockerfiles, será dada a seguinte ordem na busca pelo arquivo: "Dockerfile-iotmaker",
//     "Dockerfile", "dockerfile" na pasta raiz;
//   - Se não houver encontrado, será feita uma busca recursiva por "Dockerfile" e "dockerfile";
//   - Caso o projeto seja em golang e o arquivo main.go, contendo o pacote main, esteja contido na
//     pasta raiz, com o arquivo go.mod, podem ser usadas as funções MakeDefaultDockerfileForMe() e
//     MakeDefaultDockerfileForMeWithInstallExtras() para ser gerar um arquivo Dockerfile padrão de
//     forma automática.
func (e *ContainerBuilder) ImageBuildFromFolder() (inspect types.ImageInspect, err error) {
	err = e.verifyImageName()
	if err != nil {
		util.TraceToLog()
		return
	}

	e.imageID, _ = e.ImageFindIdByName(e.imageName)
	if e.imageID != "" && e.imageExpirationTimeIsValid() == true {
		return
	}

	if e.buildPath == "" {
		util.TraceToLog()
		err = errors.New("set build folder path first")
		return
	}

	e.buildPath, err = filepath.Abs(e.buildPath)
	if err != nil {
		util.TraceToLog()
		return
	}

	if e.makeDefaultDockerfile == true {
		var dockerfile string
		var fileList []fs.FileInfo

		fileList, err = ioutil.ReadDir(e.buildPath)
		if err != nil {
			util.TraceToLog()
			return
		}

		// fixme: modificar isto
		// deve ir para a interface{} fazer a verificação
		var pass = false
		for _, file := range fileList {
			if file.Name() == "go.mod" {
				pass = true
				break
			}
		}
		if pass == false {
			util.TraceToLog()
			err = errors.New("go.mod file not found")
			return
		}

		if e.enableCache == true {
			_, err = e.dockerSys.ImageFindIdByName(e.imageCacheName)
			if err != nil && err.Error() == "image name not found" {
				err = nil
				e.enableCache = false
			}
			if err != nil {
				util.TraceToLog()
				return
			}
		}

		dockerfile, err = e.autoDockerfile.MountDefaultDockerfile(
			e.buildOptions.BuildArgs,
			e.changePorts,
			e.openPorts,
			e.exposePortsOnDockerfile,
			e.volumes,
			e.imageInstallExtras,
			e.enableCache,
			e.imageCacheName,
		)
		if err != nil {
			util.TraceToLog()
			return
		}

		var dockerfilePath = filepath.Join(e.buildPath, "Dockerfile-iotmaker")
		err = ioutil.WriteFile(dockerfilePath, []byte(dockerfile), os.ModePerm)
		if err != nil {
			util.TraceToLog()
			return
		}
	}

	if e.printBuildOutput == true {
		e.autoDockerfile.Prayer()

		go func(ch *chan ContainerPullStatusSendToChannel) {
			for {

				select {
				case event := <-*ch:
					var stream = event.Stream
					stream = strings.ReplaceAll(stream, "\n", "")
					stream = strings.ReplaceAll(stream, "\r", "")
					stream = strings.Trim(stream, " ")

					if stream == "" {
						continue
					}

					log.Printf("%v", stream)

					//if event.Closed == true {
					//	return
					//}
				}
			}
		}(&e.changePointer)
	}

	e.imageID, err = e.dockerSys.ImageBuildFromFolder(
		e.buildPath,
		e.imageName,
		[]string{},
		e.buildOptions,
		&e.changePointer,
	)
	if err != nil {
		util.TraceToLog()

		err = errors.New(err.Error() + "\nfolder path: " + e.buildPath)
		return
	}

	if e.imageID == "" {
		util.TraceToLog()
		err = errors.New("image ID was not generated")
		return
	}

	// Construir uma imagem de múltiplas etapas deixa imagens grandes e sem serventia, ocupando espaço no HD.
	_ = e.dockerSys.ImageGarbageCollector()
	//if err != nil {
	//	return
	//}

	inspect, err = e.ImageInspect()
	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
