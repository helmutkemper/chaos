package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	ticker := time.NewTicker(15 * time.Second)
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Printf("%v\n", t)
				fmt.Printf("Bye...\n")
				os.Exit(0)
			}
		}
	}()

	fmt.Printf("starting server at port 3000\n")
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
