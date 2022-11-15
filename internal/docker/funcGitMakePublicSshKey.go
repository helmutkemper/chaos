package docker

import (
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/helmutkemper/util"
	"os"
)

// gitMakePublicSshKey
//
// English:
//
//	Mount the ssl certificate for the git clone function
//
//	 Output:
//	   publicKeys: Ponteiro de objeto compatível com o objeto ssh.PublicKeys
//	   err: standard error object
//
// Português:
//
//	 Monta o certificado ssl para a função de git clone
//
//		 Saída:
//	    publicKeys: Ponteiro de objeto compatível com o objeto ssh.PublicKeys
//	    err: objeto de erro padrão
func (e *ContainerBuilder) gitMakePublicSshKey() (publicKeys *ssh.PublicKeys, err error) {
	if e.gitData.sshPrivateKeyPath != "" {
		_, err = os.Stat(e.gitData.sshPrivateKeyPath)
		if err != nil {
			util.TraceToLog()
			return
		}
		publicKeys, err = ssh.NewPublicKeysFromFile("git", e.gitData.sshPrivateKeyPath, e.gitData.password)
	} else if e.contentIdEcdsaFile != "" {
		publicKeys, err = ssh.NewPublicKeys("git", []byte(e.contentIdEcdsaFile), e.gitData.password)
	} else if e.contentIdRsaFile != "" {
		publicKeys, err = ssh.NewPublicKeys("git", []byte(e.contentIdRsaFile), e.gitData.password)
	}

	if err != nil {
		util.TraceToLog()
		return
	}

	return
}
