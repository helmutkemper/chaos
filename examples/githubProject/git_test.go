package githubproject

import (
	"github.com/helmutkemper/chaos/factory"
	"testing"
	"time"
)

func TestLinear(t *testing.T) {

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1")
	//Test(t)

	factory.NewContainerFromGit(
		"barco:latest",
		"https://github.com/polarstreams/polar.git",
	).
		ReplaceBeforeBuild("Dockerfile", "Dockerfile").
		EnvironmentVar([]string{"BARCO_DEV_MODE=true"}).
		//PrivateRepositoryAutoConfig().
		VulnerabilityScanner().
		Ports("tcp", 3000, 3000).
		Create("barco", 1).
		Start()

	if !primordial.Monitor(3 * time.Minute) {
		t.Fail()
	}
}
