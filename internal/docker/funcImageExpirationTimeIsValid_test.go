package docker

import (
	"testing"
	"time"
)

func TestContainerBuilder_imageExpirationTimeIsValid(t *testing.T) {
	var err error

	t.Cleanup(func() {
		SaGarbageCollector()
	})

	// English: Deletes all docker elements with the term `delete` in the name.
	//
	// Português: Apaga todos os elementos docker com o termo `delete` no nome.
	// [optional/opcional]
	SaGarbageCollector()

	var docker = ContainerBuilder{}

	docker.SetImageName("delete:latest")
	docker.SetContainerName("delete_test")
	docker.SetPrintBuildOnStrOut()
	docker.SetBuildFolderPath("./test/done")

	err = docker.Init()
	if err != nil {
		t.Logf("docker.Init().error: %v", err)
		t.FailNow()
		return
	}

	_, err = docker.ImageBuildFromFolder()
	if err != nil {
		t.Logf("docker.ImageBuildFromFolder().error: %v", err)
		t.FailNow()
		return
	}

	_, err = docker.ImageInspect()
	if err != nil {
		t.Logf("docker.ImageInspect().error: %v", err)
		t.FailNow()
		return
	}

	docker.SetImageExpirationTime(time.Minute)
	if !docker.imageExpirationTimeIsValid() {
		t.Logf("docker.ImageInspect().error: %v", err)
		t.FailNow()
	}
}
