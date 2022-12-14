package githubproject

import (
	"github.com/helmutkemper/chaos/factory"
	"log"
	"testing"
	"time"
)

func TestLinear(t *testing.T) {

	primordial := factory.NewPrimordial()

	primordial.NetworkCreate("mongo", "10.0.0.0/16", "10.0.0.1")

	factory.NewContainerFromGit(
		"public:latest",
		"https://github.com/helmutkemper/chaos.public.example.git",
	).
		Ports("tcp", 3000, 3000).
		Create("public", 1).
		Start()

	if !primordial.Monitor(1 * time.Minute) {
		log.Print("fail!")
	}

	primordial.GarbageCollector()
}
