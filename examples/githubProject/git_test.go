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

	_ = os.Remove("./data/ignore.dataSent.txt")
	_ = os.Remove("./data/ignore.dataReceived.txt")
	_ = os.Remove("./data/ignore.end.empty")
	factory.NewPrimordial().GarbageCollector()

	// clear data after test
	t.Cleanup(func() {
		factory.NewPrimordial().GarbageCollector()
		_ = os.Remove("./data/ignore.dataSent.txt")
		_ = os.Remove("./data/ignore.dataReceived.txt")
		_ = os.Remove("./data/ignore.end.empty")
	})

	// create a network
	primordial := factory.NewPrimordial().
		NetworkCreate("test_network", "10.0.0.0/16", "10.0.0.1")

	// cloning polar project
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
		EnableChaos(1, 1, 1).
		Create("polar", 3).
		Start()

	// create a polar consuming container
	consumer := factory.NewContainerFromFolder(
		"consumer",
		"./consumer",
	).
		MakeDockerfile().
		DockerfileBuild("/app", "/app/main", "/app/main.go").
		Volumes("/data", "./data").
		Create("consumer", 1).
		Start()

	// create a polar producer container
	factory.NewContainerFromFolder(
		"producer",
		"./producer",
	).
		MakeDockerfile().
		DockerfileBuild("/app", "/app/main", "/app/main.go").
		Volumes("/data", "./data").
		Create("producer", 1).
		Start()

	// define a test timeout
	if !primordial.Monitor(15 * time.Minute) {
		t.FailNow()
	}

	// write this file, indicate test end to producer container

	err = os.WriteFile("./data/ignore.end.empty", nil, fs.ModePerm)
	if err != nil {
		t.Logf("write end simulation error: %v", err)
		t.FailNow()
	}

	t.Logf("end simulation signal sent")
	consumer.WaitStatusNotRunning(2 * time.Minute)

	// test data integrity
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
			t.Logf("%s", linesSent[kSent])
			t.Fail()
		}
	}
}
