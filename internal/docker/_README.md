# iotmaker.docker.builder

![./image/docker.png](./image/docker.png)

### English

This project creates a simple way to manipulate containers from Golang code

> Status: Starting to document in English

### Português

> Status: Documentando em ingles

Este projeto cria uma API Golang simples para criar e manipular o docker a partir de um código Golang

## Examples / Exemplos

# Create a docker network / Cria uma rede docker

English: Creates a docker network with subnet 10.0.0.0/16 and gateway 10.0.0.1

Português: Cria uma rede docker com subnet 10.0.0.0/16 e gateway 10.0.0.1

```golang
  var err error
  var netDocker = dockerNetwork.ContainerBuilderNetwork{}
  err = netDocker.Init()
  if err != nil { panic(err) }

  // create a network named cache_delete_after_test, subnet 10.0.0.0/16 e gatway 10.0.0.1
  err = netDocker.NetworkCreate("cache_delete_after_test", "10.0.0.0/16", "10.0.0.1")
  if err != nil { panic(err) }
```

English: use the `container.SetNetworkDocker(&netDocker)` command to link the container to the network

Português: use o comando `container.SetNetworkDocker(&netDocker)` para ligar um container com o docker

# Container nats

English: Creates a nats container, from the https://nats.io/ project and expects it to be ready, monitoring standard 
output and looking for the text "Listening for route connections on 0.0.0.0:6222"

Português: Cria um container nats, do projeto https://nats.io/ e espera o mesmo ficar pronto, monitorando a saída 
padrão e procurando pelo texto "Listening for route connections on 0.0.0.0:6222"

```golang
  var err error

  // create a container
  var container = ContainerBuilder{}
  // set image name for docker pull
  container.SetImageName("nats:latest")
  // link container and network [optional] (next ip address is 10.0.0.2)
  container.SetNetworkDocker(&netDocker)
  // set a container name
  container.SetContainerName("container_delete_nats_after_test")
  // set a waits for the text to appear in the standard container output to proceed [optional]
  container.SetWaitStringWithTimeout("Listening for route connections on 0.0.0.0:6222", 10*time.Second)

  // inialize the container object
  err = container.Init()
  if err != nil { panic(err) }

  // image nats:latest pull command [optional]
  err = container.ImagePull()
  if err != nil { panic(err) }

  // container build and start from image nats:latest
  // waits for the text "Listening for route connections on 0.0.0.0:6222" to appear  in the standard container 
  // output to proceed
  err = container.ContainerBuildFromImage()
  if err != nil { panic(err) }
```

# Container from github project

English: Creates a container based on a golang project contained in a remote git repository.

If you don't want to make the Dockerfile, use the `container.MakeDefaultDockerfileForMe()` command if the `go.mod` 
file is present and the `main.go` file is in the root directory.

If the repository is private, use the `container.SetPrivateRepositoryAutoConfig()` command to automatically copy 
the credentials contained in `~/.ssh` and the `~/.gitconfig` file into the image.

If the repository is private, use the `container.SetPrivateRepositoryAutoConfig()` command to automatically copy 
the credentials contained in `~/.ssh` and the `~/.gitconfig` file into the image.

Note that the image is built in two steps and credentials will be lost.

If the image needs to access a private repository, use the `container.SetGitPathPrivateRepository()` function to 
enter the repository.

Português: Cria um container baseado em um projeto golang contido em um repositório git remoto.

Caso você não queira fazer o Dockerfile, use o comando `container.MakeDefaultDockerfileForMe()` se o arquivo 
`go.mod` estiver presente e o arquivo `main.go` estiver na raiz do repositório.

Se o repositório for privado, use o comando `container.SetPrivateRepositoryAutoConfig()` para copiar as credenciais
contidas em `~/.ssh/` e o arquivo `~/.gitconfig` de forma automática para dentro da imagem.

Perceba que a imagem é construída em duas etapas e as credenciais serão perdidas.

Se a imagem necessitar acessar um repositório privado, use a função `container.SetGitPathPrivateRepository()` para
informar o repositório.

```golang
  var err error
  var container = ContainerBuilder{}
  // new image name delete:latest
  container.SetImageName("delete:latest")
  // container name container_delete_server_after_test
  container.SetContainerName("container_delete_server_after_test")
  // git project to clone https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git
  container.SetGitCloneToBuild("https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git")
    
  // see SetGitCloneToBuildWithUserPassworh(), SetGitCloneToBuildWithPrivateSshKey() and
  // SetGitCloneToBuildWithPrivateToken()
    
  // set a waits for the text to appear in the standard container output to proceed [optional]
  container.SetWaitStringWithTimeout("Stating server on port 3000", 10*time.Second)
  // change and open port 3000 to 3030
  container.AddPortToChange("3000", "3030")
  // replace container folder /static to host folder ./test/static
  err = container.AddFileOrFolderToLinkBetweenConputerHostAndContainer("./test/static", "/static")
  if err != nil { panic(err) }
    
  // inicialize container object
  err = container.Init()
  if err != nil { panic(err) }
    
  // builder new image from git project
  err = container.ImageBuildFromServer()
  if err != nil { panic(err) }

  err = container.ContainerBuildFromImage()
  if err != nil { panic(err) }
```

# MongoDB

English: Create a MongoDB container.

To archive data non-ephemerally, use the `mongoDocker.AddFileOrFolderToLinkBetweenConputerHostAndContainer()` 
command to define where to archive the data on the host computer.

Português: Cria um container MongoDB. 

Para arquivar dados de forma não efêmera, use o comando 
`mongoDocker.AddFileOrFolderToLinkBetweenConputerHostAndContainer()` para definir onde arquivar os dados no 
computador hospedeiro.

```golang
  var err error
  var mongoDocker = &ContainerBuilder{}
  mongoDocker.SetImageName("mongo:latest")
  mongoDocker.SetContainerName("container_delete_mongo_after_test")
  mongoDocker.AddPortToExpose("27017")
  mongoDocker.SetEnvironmentVar(
    []string{
      "--host 0.0.0.0",
    },
  )
  err = mongoDocker.AddFileOrFolderToLinkBetweenConputerHostAndContainer("./test/data", "/data")
  mongoDocker.SetWaitStringWithTimeout(`"msg":"Waiting for connections","attr":{"port":27017`, 20*time.Second)
  err = mongoDocker.Init()
  err = mongoDocker.ContainerBuildFromImage()
```

# Container from folder

English: Mount a container from a folder on the host computer.

If you don't want to make the Dockerfile, use the `container.MakeDefaultDockerfileForMe()` command if the `go.mod`
file is present and the `main.go` file is in the root directory.

If the repository is private, use the `container.SetPrivateRepositoryAutoConfig()` command to automatically copy
the credentials contained in `~/.ssh` and the `~/.gitconfig` file into the image.

If the repository is private, use the `container.SetPrivateRepositoryAutoConfig()` command to automatically copy
the credentials contained in `~/.ssh` and the `~/.gitconfig` file into the image.

Note that the image is built in two steps and credentials will be lost.

If the image needs to access a private repository, use the `container.SetGitPathPrivateRepository()` function to
enter the repository.

Português: Monta um container a partir de uma pasta no computador hospedeiro.

Caso você não queira fazer o Dockerfile, use o comando `container.MakeDefaultDockerfileForMe()` se o arquivo
`go.mod` estiver presente e o arquivo `main.go` estiver na raiz do repositório.

Se o repositório for privado, use o comando `container.SetPrivateRepositoryAutoConfig()` para copiar as credenciais
contidas em `~/.ssh/` e o arquivo `~/.gitconfig` de forma automática para dentro da imagem.

Perceba que a imagem é construída em duas etapas e as credenciais serão perdidas.

Se a imagem necessitar acessar um repositório privado, use a função `container.SetGitPathPrivateRepository()` para
informar o repositório.

```golang
  var err error

  GarbageCollector()
  var container = ContainerBuilder{}
  // new image name delete:latest
  container.SetImageName("delete:latest")
  // set a folder path to make a new image
  container.SetBuildFolderPath("./test/server")
  // container name container_delete_server_after_test
  container.SetContainerName("container_delete_server_after_test")
  // set a waits for the text to appear in the standard container output to proceed [optional]
  container.SetWaitStringWithTimeout("starting server at port 3000", 10*time.Second)
  // change and open port 3000 to 3030
  container.AddPortToExpose("3000")
  // replace container folder /static to host folder ./test/static
  err = container.AddFileOrFolderToLinkBetweenConputerHostAndContainer("./test/static", "/static")
  if err != nil { panic(err) }

  // inicialize container object
  err = container.Init()
  if err != nil { panic(err) }

  // builder new image from folder
  err = container.ImageBuildFromFolder()
  if err != nil { panic(err) }

  // build a new container from image
  err = container.ContainerBuildFromImage()
  if err != nil { panic(err) }
```