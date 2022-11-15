```
	var chanList = make([]<-chan dockerBuilder.Event, 2)
	
	for ... {
	  // English: Gets the event channel inside the container.
		//
		// Português: Pega o canal de eventos dentro do container.
		chanList[i] = simulation[i].GetChaosEvent()
	}
	
	event := mergeChannels(chanList...)
	// English: Let the example run until a failure happens to terminate the test
	//
	// Português: Deixa o exemplo rodar até que uma falha aconteça para terminar o teste
	for {
		var pass = false
		select {
		case e := <-event:
			if e.Error == true || e.Fail == true {
				util.TraceToLog()
				log.Printf("Error: %v", e.Message)
				return
			}
			if e.Done == true || e.Error == true || e.Fail == true {
				pass = true
				
				fmt.Printf("container name: %v\n", e.ContainerName)
				fmt.Printf("done: %v\n", e.Done)
				fmt.Printf("fail: %v\n", e.Fail)
				fmt.Printf("error: %v\n", e.Error)
				
				break
			}
		}
		
		if pass == true {
			break
		}
	}
```

```go
func mergeChannels(cs ...<-chan dockerBuilder.Event) <-chan dockerBuilder.Event {
	out := make(chan dockerBuilder.Event)
	var wg sync.WaitGroup
	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan dockerBuilder.Event) {
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

```