package docker

import (
	"github.com/docker/docker/api/types/mount"
	dockerfileGolang "github.com/helmutkemper/iotmaker.docker.builder.golang.dockerfile"
)

// DockerfileAuto
//
// English: Interface from automatic Dockerfile generator.
//
//	Note: To be able to access private repositories from inside the container, build the image in two or more
//	steps and in the first step, copy the id_rsa and known_hosts files to the /root/.ssh folder and the .gitconfig
//	file to the /root folder.
//
//	One way to do this automatically is to use the Dockerfile example below, where the arguments SSH_ID_RSA_FILE
//	contains the file ~/.ssh/id_rsa, KNOWN_HOSTS_FILE contains the file ~/.ssh/known_hosts and GITCONFIG_FILE
//	contains the file ~/.gitconfig.
//
//	If the ~/.ssh/id_rsa key is password protected, use the SetGitSshPassword() function to set the password.
//
//	If you want to copy the files into the image automatically, use SetPrivateRepositoryAutoConfig() and the
//	function will copy the files ~/.ssh/id_rsa, ~/.ssh/known_hosts and ~/.gitconfig to the viable arguments
//	located above.
//
//	If you want to define the files manually, use SetGitConfigFile(), SetSshKnownHostsFile() and SetSshIdRsaFile()
//	to define the files manually.
//
//	The Dockerfile below can be used as a base
//
//	  # (en) first stage of the process
//	  # (pt) primeira etapa do processo
//	  FROM golang:1.16-alpine as builder
//
//	  # (en) enable the argument variables
//	  # (pt) habilita as variáveis de argumento
//	  ARG SSH_ID_RSA_FILE
//	  ARG KNOWN_HOSTS_FILE
//	  ARG GITCONFIG_FILE
//	  ARG GIT_PRIVATE_REPO
//
//	  # (en) creates the .ssh directory within the root directory
//	  # (pt) cria o diretório .ssh dentro do diretório root
//	  RUN mkdir -p /root/.ssh/ && \
//	      # (en) creates the id_esa file inside the .ssh directory
//	      # (pt) cria o arquivo id_esa dentro do diretório .ssh
//	      echo "$SSH_ID_RSA_FILE" > /root/.ssh/id_rsa && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/ && \
//	      # (en) creates the known_hosts file inside the .ssh directory
//	      # (pt) cria o arquivo known_hosts dentro do diretório .ssh
//	      echo "$KNOWN_HOSTS_FILE" > /root/.ssh/known_hosts && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/known_hosts && \
//	      # (en) creates the .gitconfig file at the root of the root directory
//	      # (pt) cria o arquivo .gitconfig na raiz do diretório /root
//	      echo "$GITCONFIG_FILE" > /root/.gitconfig && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.gitconfig && \
//	      # (en) prepares the OS for installation
//	      # (pt) prepara o OS para instalação
//	      apk update && \
//	      # (en) install git and openssh
//	      # (pt) instala o git e o opnssh
//	      apk add --no-cache build-base git openssh && \
//	      # (en) clear the cache
//	      # (pt) limpa a cache
//	      rm -rf /var/cache/apk/*
//
//	  # (en) creates the /app directory, where your code will be installed
//	  # (pt) cria o diretório /app, onde seu código vai ser instalado
//	  WORKDIR /app
//	  # (en) copy your project into the /app folder
//	  # (pt) copia seu projeto para dentro da pasta /app
//	  COPY . .
//	  # (en) enables the golang compiler to run on an extremely simple OS, scratch
//	  # (pt) habilita o compilador do golang para rodar em um OS extremamente simples, o scratch
//	  ARG CGO_ENABLED=0
//	  # (en) adjust git to work with shh
//	  # (pt) ajusta o git para funcionar com shh
//	  RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
//	  # (en) defines the path of the private repository
//	  # (pt) define o caminho do repositório privado
//	  RUN echo "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//	  # (en) install the dependencies in the go.mod file
//	  # (pt) instala as dependências no arquivo go.mod
//	  RUN go mod tidy
//	  # (en) compiles the main.go file
//	  # (pt) compila o arquivo main.go
//	  RUN go build -ldflags="-w -s" -o /app/main /app/main.go
//	  # (en) creates a new scratch-based image
//	  # (pt) cria uma nova imagem baseada no scratch
//	  # (en) scratch is an extremely simple OS capable of generating very small images
//	  # (pt) o scratch é um OS extremamente simples capaz de gerar imagens muito reduzidas
//	  # (en) discarding the previous image erases git access credentials for your security and reduces the size of the
//	  #      image to save server space
//	  # (pt) descartar a imagem anterior apaga as credenciais de acesso ao git para a sua segurança e reduz o tamanho
//	  #      da imagem para poupar espaço no servidor
//	  FROM scratch
//	  # (en) copy your project to the new image
//	  # (pt) copia o seu projeto para a nova imagem
//	  COPY --from=builder /app/main .
//	  # (en) execute your project
//	  # (pt) executa o seu projeto
//	  CMD ["/main"]
//
// Português: Interface do gerador de dockerfile automático.
//
//	Nota: Para conseguir acessar repositórios privados de dentro do container, construa a imagem em duas ou mais
//	etapas e na primeira etapa, copie os arquivos id_rsa e known_hosts para a pasta /root/.ssh e o arquivo
//	.gitconfig para a pasta /root/.
//
//	Uma maneira de fazer isto de forma automática é usar o exemplo de Dockerfile abaixo, onde os argumentos
//	SSH_ID_RSA_FILE contém o arquivo ~/.ssh/id_rsa, KNOWN_HOSTS_FILE contém o arquivo ~/.ssh/known_hosts e
//	GITCONFIG_FILE contém o arquivo ~/.gitconfig.
//
//	Caso a chave ~/.ssh/id_rsa seja protegida com senha, use a função SetGitSshPassword() para definir a senha da
//	mesma.
//
//	Caso você queira copiar os arquivos para dentro da imagem de forma automática, use
//	SetPrivateRepositoryAutoConfig() e a função copiará os arquivos ~/.ssh/id_rsa, ~/.ssh/known_hosts e
//	~/.gitconfig para as viáveis de argumentos sitada anteriormente.
//
//	Caso queira definir os arquivos de forma manual, use SetGitConfigFile(), SetSshKnownHostsFile() e
//	SetSshIdRsaFile() para definir os arquivos de forma manual.
//
//	O arquivo Dockerfile abaixo pode ser usado como base
//
//	  # (en) first stage of the process
//	  # (pt) primeira etapa do processo
//	  FROM golang:1.16-alpine as builder
//
//	  # (en) enable the argument variables
//	  # (pt) habilita as variáveis de argumento
//	  ARG SSH_ID_RSA_FILE
//	  ARG KNOWN_HOSTS_FILE
//	  ARG GITCONFIG_FILE
//	  ARG GIT_PRIVATE_REPO
//
//	  # (en) creates the .ssh directory within the root directory
//	  # (pt) cria o diretório .ssh dentro do diretório root
//	  RUN mkdir -p /root/.ssh/ && \
//	      # (en) creates the id_esa file inside the .ssh directory
//	      # (pt) cria o arquivo id_esa dentro do diretório .ssh
//	      echo "$SSH_ID_RSA_FILE" > /root/.ssh/id_rsa && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/ && \
//	      # (en) creates the known_hosts file inside the .ssh directory
//	      # (pt) cria o arquivo known_hosts dentro do diretório .ssh
//	      echo "$KNOWN_HOSTS_FILE" > /root/.ssh/known_hosts && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.ssh/known_hosts && \
//	      # (en) creates the .gitconfig file at the root of the root directory
//	      # (pt) cria o arquivo .gitconfig na raiz do diretório /root
//	      echo "$GITCONFIG_FILE" > /root/.gitconfig && \
//	      # (en) adjust file access security
//	      # (pt) ajusta a segurança de acesso do arquivo
//	      chmod -R 600 /root/.gitconfig && \
//	      # (en) prepares the OS for installation
//	      # (pt) prepara o OS para instalação
//	      apk update && \
//	      # (en) install git and openssh
//	      # (pt) instala o git e o opnssh
//	      apk add --no-cache build-base git openssh && \
//	      # (en) clear the cache
//	      # (pt) limpa a cache
//	      rm -rf /var/cache/apk/*
//
//	  # (en) creates the /app directory, where your code will be installed
//	  # (pt) cria o diretório /app, onde seu código vai ser instalado
//	  WORKDIR /app
//	  # (en) copy your project into the /app folder
//	  # (pt) copia seu projeto para dentro da pasta /app
//	  COPY . .
//	  # (en) enables the golang compiler to run on an extremely simple OS, scratch
//	  # (pt) habilita o compilador do golang para rodar em um OS extremamente simples, o scratch
//	  ARG CGO_ENABLED=0
//	  # (en) adjust git to work with shh
//	  # (pt) ajusta o git para funcionar com shh
//	  RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/
//	  # (en) defines the path of the private repository
//	  # (pt) define o caminho do repositório privado
//	  RUN echo "go env -w GOPRIVATE=$GIT_PRIVATE_REPO"
//	  # (en) install the dependencies in the go.mod file
//	  # (pt) instala as dependências no arquivo go.mod
//	  RUN go mod tidy
//	  # (en) compiles the main.go file
//	  # (pt) compila o arquivo main.go
//	  RUN go build -ldflags="-w -s" -o /app/main /app/main.go
//	  # (en) creates a new scratch-based image
//	  # (pt) cria uma nova imagem baseada no scratch
//	  # (en) scratch is an extremely simple OS capable of generating very small images
//	  # (pt) o scratch é um OS extremamente simples capaz de gerar imagens muito reduzidas
//	  # (en) discarding the previous image erases git access credentials for your security and reduces the size of the
//	  #      image to save server space
//	  # (pt) descartar a imagem anterior apaga as credenciais de acesso ao git para a sua segurança e reduz o tamanho
//	  #      da imagem para poupar espaço no servidor
//	  FROM scratch
//	  # (en) copy your project to the new image
//	  # (pt) copia o seu projeto para a nova imagem
//	  COPY --from=builder /app/main .
//	  # (en) execute your project
//	  # (pt) executa o seu projeto
//	  CMD ["/main"]
type DockerfileAuto interface {
	MountDefaultDockerfile(args map[string]*string, changePorts []dockerfileGolang.ChangePort, openPorts []string, exposePorts []string, volumes []mount.Mount, installExtraPackages bool, useCache bool, imageCacheName string) (dockerfile string, err error)
	Prayer()
	SetFinalImageName(name string)
	AddCopyToFinalImage(src, dst string)
	SetDefaultSshFileName(name string)
}
