package main

import (
	"github.com/helmutkemper/chaos/factory"
	"log"
)

func main() {
	var errorCh = make(chan error)
	go func() {
		select {
		case err := <-errorCh:
			log.Fatal(err)
		}
	}()

	factory.NewManager(errorCh).
		Primordial().
		NetworkCreate("delete_before_test", "192.168.63.0/16", "192.168.63.1")
}
