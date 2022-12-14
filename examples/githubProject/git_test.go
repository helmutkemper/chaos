package githubproject

import (
	"github.com/helmutkemper/chaos/factory"
	"testing"
	"time"
)

func TestLinear(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1").
		Test(t)

	factory.NewContainerFromGit(
		"server:latest",
		"https://github.com/helmutkemper/chaos.public.example.git",
	).
		PrivateRepositoryAutoConfig().
		Ports("tcp", 3000, 3000).
		Create("server", 1).
		Start()

	if !primordial.Monitor(3 * time.Minute) {
		t.Fail()
	}
}
