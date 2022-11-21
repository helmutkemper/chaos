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

	mng := &Manager{}
	mng.New(errorCh)

	mng.Primordial().
		NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1")
	mng.ContainerFromImage().
		SaveStatistics("../../").
		FailFlag("../../bugs", "Multi threading initialized").
		Ports("tcp", 27017, 27016, 27015, 27014).
		Volumes("/data/db", "../../internal/builder/test/data0", "../../internal/builder/test/data1", "../../internal/builder/test/data2").
		EnvironmentVar("--host 0.0.0.0").
		Create("mongo:latest", "delete_mongo", 3).
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
