package main

import "github.com/helmutkemper/chaos/factory"

func main() {
	primordial := factory.NewPrimordial()
	primordial.NetworkCreate("test", "10.0.0.0/16", "10.0.0.1")

	factory.NewManager().
		ContainerFromImage("mongo:latest").
		SaveStatistics("./").
		EnableChaos(1, 1, 2, 1, 1.0).
		Ports("tcp", 27017, 27016, 27015, 27014).
		Volumes("/data/db", "./mongoData/data0", "./mongoData/data1", "./mongoData/data2").
		EnvironmentVar([]string{"--host 0.0.0.0"}).
		Create("mongo", 3).
		Start()
}
