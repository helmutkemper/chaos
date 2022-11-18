package manager

import (
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/chaos/internal/standalone"
	"testing"
)

func ExamplePrimordial_NetworkCreate() {
	var err error
	mng := &Manager{}
	if err = mng.New(); err != nil {
		fmt.Print(err.Error())
		return
	}

	if err = mng.Primordial().NetworkCreate("delete", "10.0.0.0/16", "10.0.0.1"); err != nil {
		fmt.Print(err.Error())
		return
	}

	standalone.GarbageCollector()

	// Output:
	//
}

func TestNetwork_NetworkCreate(t *testing.T) {
	var err error

	standalone.GarbageCollector()
	t.Cleanup(func() {
		standalone.GarbageCollector()
	})

	mng := &Manager{}
	if err = mng.New(); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	net := Network{}
	net.New(mng)

	var id string
	var data types.NetworkResource

	if err = net.NetworkCreate("delete_before_test", "10.0.0.0/16", "10.0.0.1"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if id, err = mng.DockerSys.NetworkFindIdByName("delete_before_test"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data, err = mng.DockerSys.NetworkInspect(id); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data.IPAM.Config[0].Subnet != "10.0.0.0/16" || data.IPAM.Config[0].Gateway != "10.0.0.1" {
		t.Log("subnet or gateway do not match with 10.0.0.0/16 and 10.0.0.1")
		t.FailNow()
	}

	if err = net.NetworkCreate("delete_before_test", "192.168.63.0/16", "192.168.63.1"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if id, err = mng.DockerSys.NetworkFindIdByName("delete_before_test"); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data, err = mng.DockerSys.NetworkInspect(id); err != nil {
		t.Error(err.Error())
		t.FailNow()
	}

	if data.IPAM.Config[0].Subnet != "192.168.63.0/16" || data.IPAM.Config[0].Gateway != "192.168.63.1" {
		t.Log("subnet or gateway do not match with 192.168.63.0/16 and 192.168.63.1")
		t.FailNow()
	}
}
