package docker

import (
	"github.com/helmutkemper/util"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

// SetPrivateRepositoryAutoConfig
//
// English:
//
//	Copies the ssh ~/.ssh/id_rsa file and the ~/.gitconfig file to the SSH_ID_RSA_FILE and
//	GITCONFIG_FILE variables.
//
//	 Output:
//	   err: Standard error object
//
//	 Notes:
//	   * For change ssh key file name, use SetSshKeyFileName() function.
//
// Português:
//
//	Copia o arquivo ssh ~/.ssh/id_rsa e o arquivo ~/.gitconfig para as variáveis SSH_ID_RSA_FILE e
//	GITCONFIG_FILE.
//
//	 Saída:
//	   err: Objeto de erro padrão
//
//	 Notas:
//	   * Para mudar o nome do arquivo ssh usado como chave, use a função SetSshKeyFileName().
func (e *ContainerBuilder) SetPrivateRepositoryAutoConfig() (err error) {
	var userData *user.User
	var fileData []byte
	var filePathToRead string

	userData, err = user.Current()
	if err != nil {
		return
	}

	if e.sshDefaultFileName == "" {
		e.sshDefaultFileName, err = e.GetSshKeyFileName(userData.HomeDir)
		if err != nil {
			return
		}
	}

	filePathToRead = filepath.Join(userData.HomeDir, ".ssh", e.sshDefaultFileName)
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		return
	}

	e.contentIdRsaFile = string(fileData)
	e.contentIdRsaFileWithScape = strings.ReplaceAll(e.contentIdRsaFile, `"`, `\"`)

	filePathToRead = filepath.Join(userData.HomeDir, ".ssh", "known_hosts")
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.contentKnownHostsFile = string(fileData)
	e.contentKnownHostsFileWithScape = strings.ReplaceAll(e.contentKnownHostsFile, `"`, `\"`)

	filePathToRead = filepath.Join(userData.HomeDir, ".gitconfig")
	fileData, err = ioutil.ReadFile(filePathToRead)
	if err != nil {
		util.TraceToLog()
		return
	}

	e.contentGitConfigFile = string(fileData)
	e.contentGitConfigFileWithScape = strings.ReplaceAll(e.contentGitConfigFile, `"`, `\"`)

	e.addImageBuildOptionsGitCredentials()
	return
}
