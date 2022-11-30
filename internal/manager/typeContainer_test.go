package manager

import (
	"github.com/helmutkemper/chaos/internal/standalone"
	"testing"
)

func TestContainerFromImage_Primordial(t *testing.T) {
	var errorCh = make(chan error)
	go func(t *testing.T) {
		err := <-errorCh
		t.Error(err.Error())
		t.Fail()
	}(t)

	standalone.GarbageCollector()
	t.Cleanup(func() {
		standalone.GarbageCollector()
	})

	mongodb := &Manager{}
	mongodb.New(errorCh)

	mongodb.Primordial().
		NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1")
	//mongodb.ContainerFromImage("mongo:latest").
	//  SaveStatistics("../../").
	//  FailFlag("../../bugs", "Multi threading initialized").
	//  Ports("tcp", 27017, 27016, 27015, 27014).
	//  Volumes("/data/db", "../../internal/builder/test/data0", "../../internal/builder/test/data1", "../../internal/builder/test/data2").
	//  EnvironmentVar("--host 0.0.0.0").
	//  Create("delete_mongo", 3).
	//  Start()

	barco := &Manager{}
	barco.New(errorCh)
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
		ReplaceBeforeBuild("/Dockerfile", "/Users/kemper/go/projetos/barcocopy/internal/test/chaos/simpleRestart/Dockerfile").
		Create("delete_barco", 3).
		Start()

	done := make(chan struct{})
	done <- struct{}{}
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
