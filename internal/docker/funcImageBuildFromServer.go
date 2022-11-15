package docker

import (
	"errors"
	"github.com/docker/docker/api/types"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/helmutkemper/util"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ImageBuildFromServer
//
// English:
//
//	Build a docker image from a project contained in a git repository.
//
//	 Output:
//	   inspect: Contém informações sobre a imagem criada
//	   err: standard object error
//
// Note:
//
//   - The repository must be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh();
//   - SetPrivateRepositoryAutoConfig() copies the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig;
//   - The SetGitConfigFile(), SetSshIdRsaFile() and SetSshKnownHostsFile() functions can be used to
//     set git security and configuration files manually.
//
// Português: Monta uma imagem docker a partir de um projeto contido em um repositório git.
//
//	Saída:
//	  inspect: contém informações sobre a imagem criada
//	  err: objeto de erro padrão
//
// Nota:
//
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh();
//   - SetPrivateRepositoryAutoConfig() copia as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig;
//   - As funções SetGitConfigFile(), SetSshIdRsaFile() e SetSshKnownHostsFile() podem ser usadas para
//     definir os arquivos de configurações se segurança do git manualmente.
func (e *ContainerBuilder) ImageBuildFromServer() (inspect types.ImageInspect, err error) {
	err = e.verifyImageName()
	if err != nil {
		util.TraceToLog()
		return
	}

	e.imageID, _ = e.ImageFindIdByName(e.imageName)
	if e.imageID != "" && e.imageExpirationTimeIsValid() == true {
		return
	}

	var tmpDirPath string
	var publicKeys *ssh.PublicKeys
	var gitCloneConfig *git.CloneOptions

	publicKeys, err = e.gitMakePublicSshKey()
	if err != nil {
		util.TraceToLog()
		return
	}

	tmpDirPath, err = ioutil.TempDir(os.TempDir(), "iotmaker.docker.builder.git.*")
	if err != nil {
		util.TraceToLog()
		return
	}

	defer os.RemoveAll(tmpDirPath)
	if e.gitData.sshPrivateKeyPath != "" || e.contentIdRsaFile != "" {
		gitCloneConfig = &git.CloneOptions{
			URL:      e.gitData.url,
			Auth:     publicKeys,
			Progress: nil,
		}
	} else if e.gitData.privateToke != "" {
		gitCloneConfig = &git.CloneOptions{
			// The intended use of a GitHub personal access token is in replace of your password
			// because access tokens can easily be revoked.
			// https://help.github.com/articles/creating-a-personal-access-token-for-the-command-line/
			Auth: &http.BasicAuth{
				Username: "abc123", // yes, this can be anything except an empty string
				Password: e.gitData.privateToke,
			},
			URL:      e.gitData.url,
			Progress: nil,
		}
	} else if e.gitData.user != "" && e.gitData.password != "" {
		gitCloneConfig = &git.CloneOptions{
			Auth: &http.BasicAuth{
				Username: e.gitData.user,
				Password: e.gitData.password,
			},
			URL:      e.gitData.url,
			Progress: nil,
		}
	} else {
		gitCloneConfig = &git.CloneOptions{
			URL:      e.gitData.url,
			Progress: nil,
		}
	}

	_, err = git.PlainClone(tmpDirPath, false, gitCloneConfig)
	if err != nil {
		util.TraceToLog()
		return
	}

	var data []byte
	for _, file := range e.addFileToServerBeforeBuild {
		data, err = ioutil.ReadFile(file.Src)
		if err != nil {
			return
		}

		err = ioutil.WriteFile(path.Join(tmpDirPath, file.Dst), data, fs.ModePerm)
		if err != nil {
			return
		}
	}

	if e.makeDefaultDockerfile == true {
		var dockerfile string
		var fileList []fs.FileInfo

		fileList, err = ioutil.ReadDir(tmpDirPath)
		if err != nil {
			util.TraceToLog()
			return
		}

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

		var cacheID string
		if e.enableCache == true {
			cacheID, err = e.dockerSys.ImageFindIdByName(e.imageCacheName)
			if err != nil && err.Error() != "image name not found" {
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
			cacheID != "",
			e.imageCacheName,
		)
		if err != nil {
			util.TraceToLog()
			return
		}

		var dockerfilePath = filepath.Join(tmpDirPath, "Dockerfile-iotmaker")
		err = ioutil.WriteFile(dockerfilePath, []byte(dockerfile), os.ModePerm)
		if err != nil {
			util.TraceToLog()
			return
		}
	} else if e.replaceDockerfile != "" {
		var dockerfilePath = filepath.Join(tmpDirPath, "Dockerfile-iotmaker")
		err = ioutil.WriteFile(dockerfilePath, []byte(e.replaceDockerfile), os.ModePerm)
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
		tmpDirPath,
		e.imageName,
		[]string{},
		e.buildOptions,
		&e.changePointer,
	)

	if err != nil {
		util.TraceToLog()
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

func (e *ContainerBuilder) ReplaceDockerfileFromServer(filePath string) (err error) {
	var data []byte
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.replaceDockerfile = string(data)
	return
}

func (e *ContainerBuilder) AddFileToServerBeforeBuild(dst, src string) {
	if e.addFileToServerBeforeBuild == nil {
		e.addFileToServerBeforeBuild = make([]CopyFile, 0)
	}

	e.addFileToServerBeforeBuild = append(e.addFileToServerBeforeBuild, CopyFile{Dst: dst, Src: src})
}
