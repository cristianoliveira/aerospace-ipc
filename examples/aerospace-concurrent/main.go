package main

import (
	"log"
	"sync"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
)

func main() {
	conn, err := aerospace.NewClient()
	if err != nil {
		log.Fatalf("Error creating connection: %v", err)
	}
	defer func() {
		if err := conn.CloseConnection(); err != nil {
			log.Fatalf("Error closing connection: %v", err)
		}
	}()

	var wg sync.WaitGroup
	numGoroutines := 5

	for i := range numGoroutines {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			response, err := conn.Connection().SendCommand("list-windows", []string{"--all", "--json"})
			if err != nil {
				log.Printf("Goroutine %d: Error sending command: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: Received response: %v", id, response)
		}(i)
	}

	wg.Wait()
}
