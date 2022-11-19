package manager

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/standalone"
	"testing"
)

func TestNetwork_NetworkCreate(t *testing.T) {
	var err error
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

	var id string
	var data types.NetworkResource

	mng.Primordial().NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1")

	if id, err = mng.DockerSys[0].NetworkFindIdByName("delete_before_test"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data, err = mng.DockerSys[0].NetworkInspect(id); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data.IPAM.Config[0].Subnet != "10.0.0.0/16" || data.IPAM.Config[0].Gateway != "10.0.0.1" {
		t.Log("subnet or gateway do not match with 10.0.0.0/16 and 10.0.0.1")
		t.FailNow()
	}

	mng.Primordial().NetworkCreate("delete_before_test", "192.168.63.0/16", "192.168.63.1")

	if id, err = mng.DockerSys[0].NetworkFindIdByName("delete_before_test"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data, err = mng.DockerSys[0].NetworkInspect(id); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data.IPAM.Config[0].Subnet != "192.168.63.0/16" || data.IPAM.Config[0].Gateway != "192.168.63.1" {
		t.Log("subnet or gateway do not match with 192.168.63.0/16 and 192.168.63.1")
		t.FailNow()
	}
}
