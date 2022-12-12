package manager

import (
	"github.com/helmutkemper/chaos/internal/standalone"
	"log"
	"testing"
	"time"
)

func TestContainerFromImage_Primordial(t *testing.T) {
	standalone.GarbageCollector()
	t.Cleanup(func() {
		standalone.GarbageCollector()
	})

	manager := &Manager{}
	manager.New()

	manager.Primordial().
		NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1")

	private := &Manager{}
	private.New()
	private.ContainerFromGit("private", "git@github.com:helmutkemper/iotmaker.docker.builder.private.example.git").
		SaveStatistics("../../").
		PrivateRepositoryAutoConfig().
		GitPathPrivateRepository("github.com/helmutkemper").
		Healthcheck(30*time.Second, 30*time.Second, 30*time.Second, 1, "CMD", "ash", "curl --fail http://localhost:5000/").
		Ports("tcp", 3000, 0, 3000).
		MakeDockerfile().
		EnableChaos(1, 1, 2, 1.0).
		Create("private", 3).
		Start()

	barco := &Manager{}
	barco.New()
	barco.ContainerFromFolder("barco:latest", "/Users/kemper/go/projetos/barcocopy").
		SaveStatistics("../../").
		EnvironmentVar(
			[]string{
				//"BARCO_DEV_MODE=true",
				"BARCO_SHUTDOWN_DELAY_SECS=0",
				"BARCO_CONSUMER_ADD_DELAY_MS=5000",
				"BARCO_SEGMENT_FLUSH_INTERVAL_MS=500",
				"BARCO_BROKER_NAMES=delete_barco_0,delete_barco_1,delete_barco_2",
				"BARCO_ORDINAL=0",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"BARCO_SHUTDOWN_DELAY_SECS=0",
				"BARCO_CONSUMER_ADD_DELAY_MS=5000",
				"BARCO_SEGMENT_FLUSH_INTERVAL_MS=500",
				"BARCO_BROKER_NAMES=delete_barco_0,delete_barco_1,delete_barco_2",
				"BARCO_ORDINAL=1",
			},
			[]string{
				//"BARCO_DEV_MODE=true",
				"BARCO_SHUTDOWN_DELAY_SECS=0",
				"BARCO_CONSUMER_ADD_DELAY_MS=5000",
				"BARCO_SEGMENT_FLUSH_INTERVAL_MS=500",
				"BARCO_BROKER_NAMES=delete_barco_0,delete_barco_1,delete_barco_2",
				"BARCO_ORDINAL=2",
			},
		).
		FailFlag("../../bugs", "\"fatal\"", "panic:").
		ReplaceBeforeBuild("/Dockerfile", "/Users/kemper/go/projetos/barcocopy/internal/test/chaos/simpleRestart/Dockerfile").
		Create("barco", 3).
		Start()

	mongodb := &Manager{}
	mongodb.New()
	mongodb.ContainerFromImage("mongo:latest").
		SaveStatistics("../../").
		EnableChaos(1, 1, 2, 1.0).
		Ports("tcp", 27017, 27016, 27015, 27014).
		Volumes("/data/db", "../../internal/builder/test/data0", "../../internal/builder/test/data1", "../../internal/builder/test/data2").
		EnvironmentVar([]string{"--host 0.0.0.0"}).
		Create("mongo", 1).
		Start()

	if !manager.Primordial().Monitor(5 * time.Minute) {
		t.Fail()
	}

	log.Printf("done!")
}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
