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
		"polar:latest",
		"https://github.com/polarstreams/polar.git",
	).
		ReplaceBeforeBuild("Dockerfile", "Dockerfile").
		EnvironmentVar(
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=0",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=1",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"POLAR_SHUTDOWN_DELAY_SECS=0",
				"POLAR_CONSUMER_ADD_DELAY_MS=5000",
				"POLAR_SEGMENT_FLUSH_INTERVAL_MS=500",
				"POLAR_BROKER_NAMES=delete_polar_0,delete_polar_1,delete_polar_2",
				"POLAR_ORDINAL=2",
			},
		).
		Ports("tcp", 9250, 9250).
		Ports("tcp", 9251, 9251).
		Ports("tcp", 9252, 9252).
		Create("polar", 3).
		Start()

	if !primordial.Monitor(3 * time.Minute) {
		t.Fail()
	}
}
