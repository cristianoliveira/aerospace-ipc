package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cristianoliveira/aerospace-ipc"
)

func main() {
	socketPath := fmt.Sprintf("/tmp/bobko.%s-%s.sock", "aerospace", os.Getenv("USER"))
	client, err := aerospace.NewAeroSpaceCustomConnection(
		aerospace.AeroSpaceCustomConnectionOpts{
			SocketPath:      socketPath,
			ValidateVersion: true,
		},
	)
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
