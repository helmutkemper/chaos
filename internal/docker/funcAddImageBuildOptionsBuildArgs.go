package docker

// AddImageBuildOptionsBuildArgs
//
// English:
//
//	Set build-time variables (--build-arg)
//
//	 Input:
//	   key: Argument name
//	   value: Argument value
//
// Example:
//
//	key:   argument key (e.g. Dockerfile: ARG key)
//	value: argument value
//
//	https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg
//	docker build --build-arg HTTP_PROXY=http://10.20.30.2:1234
//
//	  code:
//	    var key = "GIT_PRIVATE_REPO"
//	    var value = "github.com/yourgit"
//
//	    var container = ContainerBuilder{}
//	    container.AddImageBuildOptionsBuildArgs(key, &value)
//
//	  Dockerfile:
//	    FROM golang:1.16-alpine as builder
//	    ARG GIT_PRIVATE_REPO
//	    RUN go env -w GOPRIVATE=$GIT_PRIVATE_REPO
//
// Português:
//
//	Adiciona uma variável durante a construção (--build-arg)
//
//	 Input:
//	   key: Nome do argumento.
//	   value: Valor do argumento.
//
// Exemplo:
//
//	key:   chave do argumento (ex. Dockerfile: ARG key)
//	value: valor do argumento
//
//	https://docs.docker.com/engine/reference/commandline/build/#set-build-time-variables---build-arg
//	docker build --build-arg HTTP_PROXY=http://10.20.30.2:1234
//
//	  code:
//	    var key = "GIT_PRIVATE_REPO"
//	    var value = "github.com/yourgit"
//
//	    var container = ContainerBuilder{}
//	    container.AddImageBuildOptionsBuildArgs(key, &value)
//
//	  Dockerfile:
//	    FROM golang:1.16-alpine as builder
//	    ARG GIT_PRIVATE_REPO
//	    RUN go env -w GOPRIVATE=$GIT_PRIVATE_REPO
func (e *ContainerBuilder) AddImageBuildOptionsBuildArgs(key string, value *string) {
	if e.buildOptions.BuildArgs == nil {
		e.buildOptions.BuildArgs = make(map[string]*string)
	}

	e.buildOptions.BuildArgs[key] = value
}
