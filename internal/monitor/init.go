package monitor

import (
	"log"
	"sync"
)

var ErrorChList = make([]<-chan error, 0)
var FailChList = make([]<-chan string, 0)
var DoneChList = make([]<-chan struct{}, 0)
var EndFunc = make([]func(), 0)
var ChaosFunc = make([]func(), 0)

var counterEndFunc = 0

func AddChaosFunc(f ...func()) {
	counterEndFunc += 1
	ChaosFunc = append(ChaosFunc, f...)
}

func AddEndFunc(f func()) {
	EndFunc = append(EndFunc, f)
}

func EndAll() {
	for k := range EndFunc {
		EndFunc[k]()
	}
}

func Monitor() (pass bool) {
	for k := range ChaosFunc {
		if ChaosFunc[k] != nil {
			ChaosFunc[k]()
		} else {
			log.Printf("bug: chaos func is nil")
		}
	}

	eventError := mergeErrorChannels(ErrorChList...)
	eventFail := mergeFailChannels(FailChList...)
	eventDone := mergeChannels(DoneChList...)

	var end bool
	for {
		select {
		case err := <-eventError:
			end = true
			log.Printf("test error: %v", err)
		case fail := <-eventFail:
			end = true
			log.Printf("test fail: %v", fail)
		case <-eventDone:
			counterEndFunc -= 1
			if counterEndFunc <= 0 {
				end = true
				pass = true
				EndAll()
			}
		}

		if end {
			break
		}
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
