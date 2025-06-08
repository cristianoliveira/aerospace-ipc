package main

import (
	"log"
	"sync"

	client "github.com/cristianoliveira/aerospace-ipc"
)

func main() {
	conn, err := client.NewAeroSpaceClient()
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
			response, err := conn.Client().SendCommand("list-windows", []string{"--all", "--json"})
			if err != nil {
				log.Printf("Goroutine %d: Error sending command: %v", id, err)
				return
			}
			log.Printf("Goroutine %d: Received response: %v", id, response)
		}(i)
	}

	wg.Wait()
}
