package docker

// SetCacheEnable
//
// English:
//
//	When true, looks for an image named `chache:latest` as a basis for creating new images when the
//	MakeDefaultDockerfileForMe() function is used.
//
//	 Input:
//	   value: true to enable the use of image named cache:latest as the basis for new images if it
//	     exists
//
// Note:
//
//   - This function is extremely useful when developing new applications, reducing the time to create
//     images with each new test.
//
// Example:
//
//	Folder: cache
//	File: Dockerfile-iotmaker
//	Need: Image with nats.io drive installed
//	Content:
//
//	FROM golang:1.16-alpine as builder
//	RUN mkdir -p /root/.ssh/ && \
//	    apk update && \
//	    apk add --no-cache build-base && \
//	    apk add --no-cache alpine-sdk && \
//	    rm -rf /var/cache/apk/*
//	ARG CGO_ENABLED=0
//
//	RUN go get -u github.com/nats-io/nats.go
//
//	Code Golang:
//
//	var imageCacheName = "cache:latest"
//	var imageId string
//	var container = &dockerBuilder.ContainerBuilder{}
//
//	imageId, err = container.ImageFindIdByName(imageCacheName)
//	if err != nil && err.Error() != "image name not found" {
//	  return
//	}
//
//	if imageId != "" {
//	  return
//	}
//
//	container.SetImageName(imageCacheName)
//	container.SetPrintBuildOnStrOut()
//	container.SetContainerName(imageCacheName)
//	container.SetBuildFolderPath("./cache")
//	err = container.Init()
//	if err != nil {
//	  return
//	}
//
//	err = container.ImageBuildFromFolder()
//	if err != nil {
//	  return
//	}
//
// Português:
//
//	Quando true, procura por uma imagem de nome `chache:latest` como base para a criação de novas
//	imagens quando a função MakeDefaultDockerfileForMe() é usada.
//
//	 Entrada:
//	   value: true para habilitar o uso da imagem de nome cache:latest como base para novas imagens,
//	     caso a mesma exista
//
// Nota:
//
//   - Esta função é extremamente útil no desenvolvimento de novas aplicações, reduzindo o tempo de
//     criação de imagens a cada novo teste.
//
// Exemplo:
//
//	Pasta: cache
//	Arquivo: Dockerfile-iotmaker
//	Necessidade: Imagem com o drive do nats.io instalada
//	Conteúdo:
//
//	FROM golang:1.16-alpine as builder
//	RUN mkdir -p /root/.ssh/ && \
//	    apk update && \
//	    apk add --no-cache build-base && \
//	    apk add --no-cache alpine-sdk && \
//	    rm -rf /var/cache/apk/*
//	ARG CGO_ENABLED=0
//
//	RUN go get -u github.com/nats-io/nats.go
//
//	Código Golang:
//
//	var imageCacheName = "cache:latest"
//	var imageId string
//	var container = &dockerBuilder.ContainerBuilder{}
//
//	imageId, err = container.ImageFindIdByName(imageCacheName)
//	if err != nil && err.Error() != "image name not found" {
//	  return
//	}
//
//	if imageId != "" {
//	  return
//	}
//
//	container.SetImageName(imageCacheName)
//	container.SetPrintBuildOnStrOut()
//	container.SetContainerName(imageCacheName)
//	container.SetBuildFolderPath("./cache")
//	err = container.Init()
//	if err != nil {
//	  return
//	}
//
//	err = container.ImageBuildFromFolder()
//	if err != nil {
//	  return
//	}
func (e *ContainerBuilder) SetCacheEnable(value bool) {
	e.enableCache = value
}
