package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
)

func main() {
	client, err := aerospace.NewClient()
	if err != nil {
		log.Fatalf("failed to create AeroSpace client: %v", err)
	}
	defer client.CloseConnection()

	// Example usage - get all windows
	windows, err := client.Windows().GetAllWindows()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting windows: %v\n", err)
		os.Exit(1)
	}

	// Print windows
	for _, window := range windows {
		fmt.Println(window)
	}
}
