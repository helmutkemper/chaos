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
var IpAddress = make(map[string]string)

var Err bool

var counterEndFunc = 0

func AddIpAddress(container, ip string) {
	IpAddress[container] = ip
}

func GetIpAddress(container string) (ip string) {
	return IpAddress[container]
}

func AddChaosFunc(f ...func()) {
	ChaosFunc = append(ChaosFunc, f...)
}

func AddEndFunc(f func()) {
	counterEndFunc += 1
	EndFunc = append(EndFunc, f)
}

func EndAll() {
	for k := range EndFunc {
		EndFunc[k]()
	}
}

func Monitor() (pass bool) {
	if !Err {
		for k := range ChaosFunc {
			if ChaosFunc[k] != nil {
				ChaosFunc[k]()
			} else {
				log.Printf("bug: chaos func is nil")
			}
		}
	}

	eventError := mergeErrorChannels(ErrorChList...)
	eventFail := mergeFailChannels(FailChList...)
	eventDone := mergeChannels(DoneChList...)

	for {
		select {
		case err := <-eventError:
			log.Printf("test error: %v", err)
			return
		case fail := <-eventFail:
			log.Printf("test fail: %v", fail)
			return
		case <-eventDone:
			counterEndFunc -= 1
			if counterEndFunc <= 0 {
				pass = true
				return
			}
		}
	}
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
