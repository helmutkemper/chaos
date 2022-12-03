package monitor

import (
	"log"
	"sync"
)

var ErrorChList = make([]<-chan error, 0)
var FailChList = make([]<-chan string, 0)
var DoneChList = make([]<-chan struct{}, 0)

func Monitor() (pass bool) {
	eventError := mergeErrorChannels(ErrorChList...)
	eventFail := mergeFailChannels(FailChList...)
	eventDone := mergeChannels(DoneChList...)

	select {
	case err := <-eventError:
		log.Printf("test error: %v", err)
	case fail := <-eventFail:
		log.Printf("test fail: %v", fail)
	case <-eventDone:
		log.Printf("done!")
		pass = true
	}

	return
}

func mergeErrorChannels(cs ...<-chan error) <-chan error {
	out := make(chan error)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan error) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func mergeFailChannels(cs ...<-chan string) <-chan string {
	out := make(chan string)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan string) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func mergeChannels(cs ...<-chan struct{}) <-chan struct{} {
	out := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan struct{}) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
