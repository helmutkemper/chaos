package githubproject

import (
	"bytes"
	"github.com/helmutkemper/chaos/factory"
	"io/fs"
	"os"
	"testing"
	"time"
)

func TestLinear(t *testing.T) {
	var err error

	t.Cleanup(func() {
		factory.NewPrimordial().GarbageCollector()
		_ = os.Remove("./data/ignore.dataSent.txt")
		_ = os.Remove("./data/ignore.dataReceived.txt")
	})

	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1") //.
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

	consumer := factory.NewContainerFromFolder(
		"consumer",
		"./consumer",
	).
		MakeDockerfile().
		Volumes("/data", "./data").
		Create("consumer", 1).
		Start()

	factory.NewContainerFromFolder(
		"producer",
		"./producer",
	).
		MakeDockerfile().
		Volumes("/data", "./data").
		Create("producer", 1).
		Start()

	if !primordial.Monitor(2 * time.Minute) {
		t.Fail()
	}

	_ = os.WriteFile("./data/ignore.end.empty", nil, fs.ModePerm)
	consumer.WaitStatusNotRunning()

	var data []byte
	data, err = os.ReadFile("./data/ignore.dataSent.txt")
	if err != nil {
		t.Logf("read sent data log error: %v", err)
		t.FailNow()
	}
	linesSent := bytes.Split(data, []byte("\n"))

	data, err = os.ReadFile("./data/ignore.dataReceived.txt")
	if err != nil {
		t.Logf("read receiceved data log error: %v", err)
		t.FailNow()
	}

	linesReceived := bytes.Split(data, []byte("\n"))
	for kSent := range linesSent {
		var pass = false

		for kReceived := range linesReceived {
			if bytes.Equal(linesSent[kSent], linesReceived[kReceived]) {
				pass = true
				break
			}
		}

		if !pass {
			t.Log("found an inconsistency in the shared data")
			t.FailNow()
		}
	}
}
