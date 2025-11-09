package main

import (
	"fmt"
	"log"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
)

func main() {
	client, err := aerospace.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer func() {
		if err := client.CloseConnection(); err != nil {
			log.Fatalf("Failed to close connection: %v", err)
		}
	}()

	windows, err := client.Windows().GetAllWindows()
	if err != nil {
		log.Fatalf("Failed to get windows: %v", err)
	}

	for _, window := range windows {
		fmt.Println(window)
	}

	fmt.Println("Listed all windows successfully.")

	socketPath, err := client.Connection().GetSocketPath()
	if err != nil {
		log.Fatalf("Failed to get socket path: %v", err)
	}
	fmt.Printf("Socket path: %s\n", socketPath)
}
