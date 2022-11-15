package main

import (
	"github.com/helmutkemper/chaos/internal/docker"
	"log"
)

func main() {
	var err error

	manager := docker.Builder{}
	err = manager.New()
	if err != nil {
		log.Fatal(err)
	}

	primordial := manager.Primordial()
	err = primordial.SetProject("https://github.com/helmutkemper/iotmaker.docker.builder.public.example.git")
	if err != nil {
		log.Fatal(err)
	}

	err = primordial.SetImageName("delete:latest")
	if err != nil {
		log.Fatal(err)
	}

	primordial.SetContainerName("delete")

	err = manager.Init()
	if err != nil {
		log.Fatal(err)
	}
}
