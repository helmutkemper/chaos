package docker

import (
	"testing"
)

func TestIncIpV4Address(t *testing.T) {
	var err error
	var c ContainerBuilder
	var next string
	next, err = c.incIpV4Address("10.0.0.1", 1)
	if err != nil {
		t.FailNow()
	}

	if next != "10.0.0.2" {
		t.FailNow()
	}

	next, err = c.incIpV4Address("10.0.0.1", 10)
	if err != nil {
		t.FailNow()
	}

	if next != "10.0.0.11" {
		t.FailNow()
	}
}
