package main

import (
	"fmt"
	"log"

	"github.com/cristianoliveira/aerospace-ipc"
)

func main() {
	client, err := aerospace.NewAeroSpaceConnection()
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer client.CloseConnection()

	windows, err := client.GetAllWindows()
	if err != nil {
		log.Fatalf("Failed to get windows: %v", err)
	}

	for _, window := range windows {
		fmt.Println(window)
	}

	fmt.Println("Listed all windows successfully.")
}
