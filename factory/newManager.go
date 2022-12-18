package factory

import (
	"github.com/helmutkemper/chaos/internal/manager"
	"github.com/helmutkemper/chaos/internal/standalone"
)

func NewContainerFromGit(imageName, serverPath string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromGit(imageName, serverPath).
		Reports()
}

func NewContainerFromFolder(imageName, buildPath string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromFolder(imageName, buildPath).
		Reports()
}

func NewContainerFromImage(imageName string) (reference *manager.ContainerFromImage) {
	ref := new(manager.Manager)
	ref.New()
	return ref.ContainerFromImage(imageName).
		Reports()
}

func NewPrimordial() (reference *manager.Primordial) {
	standalone.GarbageCollector()

	ref := new(manager.Manager)
	ref.New()
	reference = ref.Primordial()
	return
}
