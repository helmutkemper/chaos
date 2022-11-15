package docker

import (
	"fmt"
	"github.com/helmutkemper/util"
)

func ExampleContainerBuilder_ImageListExposedVolumes() {
	var err error
	var volumes []string
	
	SaGarbageCollector()
	
	var container = ContainerBuilder{}
	container.SetPrintBuildOnStrOut()
	
	// new image name delete:latest
	container.SetImageName("delete:latest")
	// git project to clone https://github.com/helmutkemper/iotmaker.docker.util.whaleAquarium.sample.git
	container.SetGitCloneToBuild("https://github.com/helmutkemper/iotmaker.docker.builder.public.example.git")
	err = container.Init()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	// todo: fazer o teste do inspect
	_, err = container.ImageBuildFromServer()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	volumes, err = container.ImageListExposedVolumes()
	if err != nil {
		util.TraceToLog()
		panic(err)
	}
	
	fmt.Printf("%v", volumes[0])
	
	SaGarbageCollector()
	
	// Output:
	// /static
}
