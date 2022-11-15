package docker

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Manager struct {
	docker  *DockerSystem
	builder *ContainerBuilder
}

func (e *Manager) New() (err error) {
	e.docker = new(DockerSystem)
	err = e.docker.Init()
	if err != nil {
		err = fmt.Errorf("manager.New().error: failed to connect with docker: %v", err)
		return
	}

	e.builder = new(ContainerBuilder)
	return
}

func (e *Manager) Builder() (ref *Builder) {
	ref = new(Builder)
	ref.docker = e.docker
	ref.builder = e.builder

	return
}

func (e *Manager) removeContainerBeforeCreate() (err error) {
	if e.builder.containerName == "" {
		return
	}

	// English: Several containers are created with `container name` + NUMBER
	// Português: Vários containers são criados com o `nome do container` + _NUMBER
	var re = regexp.MustCompile("^" + e.builder.containerName + "_\\D+$")
	var list []NameAndId
	list, _ = e.builder.ContainerFindIdByNameContains(e.builder.containerName)
	for k := range list {
		if list[k].Name == e.builder.containerName {
			err = e.builder.ContainerRemove(false)
			if err != nil {
				err = fmt.Errorf("manager.Init().Error: an container named '%v' already exists. The container must be deleted before continuing", e.builder.containerName)
				return
			}
		} else if re.MatchString(list[k].Name) {
			err = e.builder.ContainerRemove(false)
			if err != nil {
				err = fmt.Errorf("manager.Init().Error: an container named '%v' already exists. The container must be deleted before continuing", e.builder.containerName)
				return
			}
		}
	}

	return
}

func (e *Manager) removeImageBeforeCreate() (err error) {
	var id string
	id, err = e.builder.ImageFindIdByName(e.builder.imageName)
	if err == nil {
		err = e.docker.ImageRemove(id, true, true)
		if err != nil {
			err = fmt.Errorf("manager.Init().Error: an image named '%v' already exists. The image must be deleted before continuing", e.builder.imageName)
			return
		}
	}

	return
}

func (e *Manager) Init() (err error) {

	// English: Removes the container, if there is one with the same name, before creating
	// Português: Remove o container, caso haja algum com o mesmo nome, antes de criar
	if err = e.removeContainerBeforeCreate(); err != nil {
		return
	}

	// English: Removes the image, if there is one with the same name, before creating
	// Português: Remove a imagem, caso haja alguma com o mesmo nome, antes de criar
	if err = e.removeImageBeforeCreate(); err != nil {
		return
	}

	if err = e.builder.Init(); err != nil {
		err = fmt.Errorf("manager.Init().Builder.Init().Error: %v", err)
		return
	}

	// English: Download the image
	// Português: Faz o download da imagem
	if !(e.builder.buildFolder || e.builder.buildServer) && e.builder.imageName != "" {
		if err = e.builder.imagePull(); err != nil {
			err = fmt.Errorf("manager.Init().Builder.imagePull().Error: %v", err)
			return
		}
	}

	// English: Set when project path is a folder
	// Português: Definido quando o caminho do projeto é uma pasta
	if e.builder.buildFolder {
		if _, err = e.builder.ImageBuildFromFolder(); err != nil {
			err = fmt.Errorf("manager.Init().Builder.ImageBuildFromFolder().Error: %v", err)
			return
		}

		// English: Set when project path ends in .git
		// Português: Definido quando o caminho do projeto termina em .git
	} else if e.builder.buildServer {
		if _, err = e.builder.ImageBuildFromServer(); err != nil {
			err = fmt.Errorf("manager.Init().Builder.ImageBuildFromServer().Error: %v", err)
			return
		}
	}

	if err = e.builder.ContainerBuildWithoutStartingItFromImage(); err != nil {
		err = fmt.Errorf("manager.Init().Builder.ContainerBuildWithoutStartingItFromImage().Error: %v", err)
		return
	}

	return
}

type Problem struct {
	Manager
}

// FlagToFile
//
// English:
//
//	Looks for error text in the container's standard output and saves it to a log file on the host
//	computer
//
//	 Input:
//	   value: Error text
//	   logDirectoryPath: File path where the container's standard output filed in a `log.N.log` file
//	     will be saved, where N is an automatically incremented number. e.g.: "./bug/critical/"
//
//	 Output:
//	   err: Default error object
//
// Português:
//
//	Procura por um texto indicativo de erro na saída padrão do container e o salva em um arquivo de
//	log no computador hospedeiro
//
//	 Entrada:
//	   value: Texto indicativo de erro
//	   logDirectoryPath: Caminho do arquivo onde será salva a saída padrão do container arquivada em
//	     um arquivo `log.N.log`, onde N é um número incrementado automaticamente.
//	     Ex.: "./bug/critical/"
//
//	 Output:
//	   err: Objeto de erro padrão
func (e *Problem) FlagToFile(value, logDirectoryPath string) (err error) {
	err = e.builder.AddFailMatchFlagToFileLog(
		value,
		logDirectoryPath,
	)
	if err != nil {
		err = fmt.Errorf("problem.FlagToFile().Builder.AddFailMatchFlagToFileLog().Error: %v", err)
		return
	}

	return
}

type Image struct {
	Manager
}

// Pull
//
// English:
//
//	Downloads the image to be mounted. (equivalent to the docker pull image command)
//
//	 Output:
//	   err: standart error object
//
// Português:
//
//	Baixa a imagem a ser montada. (equivale ao comando docker pull image)
//
//	 Saída:
//	   err: objeto de erro padrão
func (e *Image) Pull() (err error) {
	err = e.builder.imagePull()
	return
}

type Security struct {
	Manager
}

// SetPrivateSSHKey
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//	   privateSSHKeyPath: this is the path of the private ssh key compatible with the public key
//	     registered in git
//	   password: password used when the ssh key was generated or empty string
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var privateSSHKeyPath string
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	privateSSHKeyPath = filepath.Join(usr.HomeDir, ".shh", "id_ecdsa")
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateSSHKey(url, privateSSHKeyPath, password)
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//	   privateSSHKeyPath: este é o caminho da chave ssh privada compatível com a chave pública
//	     cadastrada no git
//	   password: senha usada no momento que a chave ssh foi gerada ou string em branco
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var privateSSHKeyPath string
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	privateSSHKeyPath = filepath.Join(usr.HomeDir, ".shh", "id_ecdsa")
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateSSHKey(url, privateSSHKeyPath, password)
//	container.SetGitConfigFile(string(file))
func (e *Security) SetPrivateSSHKey(privateSSHKeyPath, password string) {
	e.builder.gitData.sshPrivateKeyPath = privateSSHKeyPath
	e.builder.gitData.password = password
}

// SetPrivateToken
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//	   privateToken: token defined on your git tool portal
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//	   privateToken: token definido no portal da sua ferramenta git
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
func (e *Security) SetPrivateToken(privateToken string) {
	e.builder.gitData.privateToke = privateToken
}

// SetUserAndPassword
//
// English:
//
//	Defines the path of a repository to be used as the base of the image to be mounted.
//
//	 Input:
//	   url: Address of the repository containing the project
//	   user: git user
//	   password: git password
//
// Note:
//
//   - If the repository is private and the host computer has access to the git server, use
//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
//     settings of ~/.gitconfig automatically;
//   - To be able to access private repositories from inside the container, build the image in two or
//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
//     folder, and the ~/.gitconfig file to the /root folder;
//   - The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
//   - If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
//     password;
//   - If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
//     SetSshIdRsaFile() to define the files manually;
//   - This function must be used with the ImageBuildFromServer() and SetImageName() function to
//     download and generate an image from the contents of a git repository;
//   - The repository must contain a Dockerfile file and it will be searched for in the following
//     order:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' and '.*dockerfile.*';
//   - The repository can be defined by the methods SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
//
// Português:
//
//	Define o caminho de um repositório para ser usado como base da imagem a ser montada.
//
//	 Entrada:
//	   url: Endereço do repositório contendo o projeto
//	   user: git user
//	   password: git password
//
// Nota:
//
//   - Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
//     configurações de ~/.gitconfig de forma automática;
//   - Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
//   - A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
//   - Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
//     definir a senha da mesma;
//   - Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
//   - Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
//     gerar uma imagem a partir do conteúdo de um repositório git;
//   - O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
//     '.*Dockerfile.*' e '.*dockerfile.*';
//   - O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e
//     SetGitCloneToBuildWithUserPassworh().
//
// code:
//
//	var err error
//	var usr *user.User
//	var userGitConfigPath string
//	var file []byte
//	usr, err = user.Current()
//	if err != nil {
//	  panic(err)
//	}
//
//	userGitConfigPath = filepath.Join(usr.HomeDir, ".gitconfig")
//	file, err = ioutil.ReadFile(userGitConfigPath)
//
//	var container = ContainerBuilder{}
//	container.SetGitCloneToBuildWithPrivateToken(url, privateToken)
//	container.SetGitConfigFile(string(file))
func (e *Security) SetUserAndPassword(user, password string) {
	e.builder.gitData.user = user
	e.builder.gitData.password = password
}

type Builder struct {
	Manager
}

func (e *Builder) Get() (ref *Security) {
	ref = new(Security)
	ref.docker = e.docker
	ref.builder = e.builder

	return
}

func (e *Builder) Primordial() (ref *Primordial) {
	ref = new(Primordial)
	ref.docker = e.docker
	ref.builder = e.builder

	return
}

func (e *Builder) Problem() (ref *Problem) {
	ref = new(Problem)
	ref.docker = e.docker
	ref.builder = e.builder

	return
}

type Primordial struct {
	Manager
}

// SetImageName
//
// English:
//
//	Defines the name of the image to be used or created
//
//	 Input:
//	   value: name of the image to be downloaded or created
//
// Português:
//
//	Define o nome da imagem a ser usada ou criada
//
//	 Entrada:
//	   value: noma da imagem a ser baixada ou criada
func (e *Primordial) SetImageName(value string) (err error) {
	if !strings.Contains(value, ":") {
		err = fmt.Errorf("primordial.SetImageName().error: image name must contain version tag. eg: %v:latest", value)
		return
	}

	e.builder.imageName = value
	return
}

// SetContainerName
//
// English:
//
//	Defines the name of the container
//
//	 Input:
//	   value: container name
//
// Português:
//
//	Define o nome do container
//
//	 Entrada:
//	   value: nome do container
func (e *Primordial) SetContainerName(value string) {
	e.builder.containerName = value
}

func (e *Primordial) SetProject(path string) (err error) {
	if strings.HasSuffix(path, ".git") {
		e.builder.buildServer = true

		// SetGitCloneToBuild
		//
		// English:
		//
		//  Defines the path of a repository to be used as the base of the image to be mounted.
		//
		//   Input:
		//     url: Address of the repository containing the project
		//
		// Note:
		//
		//   * If the repository is private and the host computer has access to the git server, use
		//     SetPrivateRepositoryAutoConfig() to copy the git credentials contained in ~/.ssh and the
		//     settings of ~/.gitconfig automatically;
		//   * To be able to access private repositories from inside the container, build the image in two or
		//     more steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh
		//     folder, and the ~/.gitconfig file to the /root folder;
		//   * The MakeDefaultDockerfileForMe() function make a standard dockerfile with the procedures above;
		//   * If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the
		//     password;
		//   * If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and
		//     SetSshIdRsaFile() to define the files manually;
		//   * This function must be used with the ImageBuildFromServer() and SetImageName() function to
		//     download and generate an image from the contents of a git repository;
		//   * The repository must contain a Dockerfile file and it will be searched for in the following
		//     order:
		//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
		//     '.*Dockerfile.*' and '.*dockerfile.*';
		//   * The repository can be defined by the methods SetGitCloneToBuild(),
		//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() and
		//     SetGitCloneToBuildWithUserPassworh().
		//
		// Português:
		//
		//  Define o caminho de um repositório para ser usado como base da imagem a ser montada.
		//
		//   Entrada:
		//     url: Endereço do repositório contendo o projeto
		//
		// Nota:
		//
		//   * Caso o repositório seja privado e o computador hospedeiro tenha acesso ao servidor git, use
		//     SetPrivateRepositoryAutoConfig() para copiar as credências do git contidas em ~/.ssh e as
		//     configurações de ~/.gitconfig de forma automática;
		//   * Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas
		//     ou mais etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta
		//     /root/.ssh e o arquivo .gitconfig para a pasta /root/;
		//   * A função MakeDefaultDockerfileForMe() monta um dockerfile padrão com os procedimentos acima;
		//   * Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para
		//     definir a senha da mesma;
		//   * Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile()
		//     e SetSshIdRsaFile() para definir os arquivos de forma manual;
		//   * Esta função deve ser usada com a função ImageBuildFromServer() e SetImageName() para baixar e
		//     gerar uma imagem a partir do conteúdo de um repositório git;
		//   * O repositório deve contar um arquivo Dockerfile e ele será procurado na seguinte ordem:
		//     './Dockerfile-iotmaker', './Dockerfile', './dockerfile', 'Dockerfile.*', 'dockerfile.*',
		//     '.*Dockerfile.*' e '.*dockerfile.*';
		//   * O repositório pode ser definido pelos métodos SetGitCloneToBuild(),
		//     SetGitCloneToBuildWithPrivateSshKey(), SetGitCloneToBuildWithPrivateToken() e SetGitCloneToBuildWithUserPassworh().
		e.builder.gitData.url = path
		return
	}

	_, err = os.Stat(path)
	if os.IsExist(err) {
		e.builder.buildFolder = true

		// SetBuildFolderPath
		//
		// English:
		//
		//  Defines the path of the folder to be transformed into an image
		//
		//   Input:
		//     value: path of the folder to be transformed into an image
		//
		// Note:
		//
		//   * The folder must contain a dockerfile file, but since different uses can have different
		//     dockerfiles, the following order will be given when searching for the file:
		//     "Dockerfile-iotmaker", "Dockerfile", "dockerfile" in the root folder;
		//   * If not found, a recursive search will be done for "Dockerfile" and "dockerfile";
		//   * If the project is in golang and the main.go file, containing the package main, is contained in
		//     the root folder, with the go.mod file, the MakeDefaultDockerfileForMe() function can be used to
		//     use a standard Dockerfile file
		//
		// Português:
		//
		//  Define o caminho da pasta a ser transformada em imagem
		//
		//   Entrada:
		//     value: caminho da pasta a ser transformada em imagem
		//
		// Nota:
		//
		//   * A pasta deve conter um arquivo dockerfile, mas, como diferentes usos podem ter diferentes
		//     dockerfiles, será dada a seguinte ordem na busca pelo arquivo: "Dockerfile-iotmaker",
		//     "Dockerfile", "dockerfile" na pasta raiz.
		//   * Se não houver encontrado, será feita uma busca recursiva por "Dockerfile" e "dockerfile"
		//   * Caso o projeto seja em golang e o arquivo main.go, contendo o pacote main, esteja contido na
		//     pasta raiz, com o arquivo go.mod, pode ser usada a função MakeDefaultDockerfileForMe() para ser
		//     usado um arquivo Dockerfile padrão
		e.builder.buildPath = path
		return
	}

	err = fmt.Errorf("primordial.SetProject().error: path was not recognized as git server or local folder")
	return
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
