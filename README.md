# aerospace-ipc
[![Go project version](https://badge.fury.io/go/github.com%2Fcristianoliveira%2Faerospace-ipc.svg)](https://badge.fury.io/go/github.com%2Fcristianoliveira%2Faerospace-ipc)
[![Go Reference](https://pkg.go.dev/badge/github.com/cristianoliveira/aerospace-ipc.svg)](https://pkg.go.dev/github.com/cristianoliveira/aerospace-ipc)
[![Quick Checks](https://github.com/cristianoliveira/aerospace-ipc/actions/workflows/on-push.yml/badge.svg)](https://github.com/cristianoliveira/aerospace-ipc/actions/workflows/on-push.yml)

A Go library for interacting with [AeroSpace WM](https://github.com/nikitabobko/AeroSpace) via IPC

*Minimum AeroSpace version:* `v0.15.x`

## Description

This package allows interacting with [AeroSpace WM](https://github.com/nikitabobko/AeroSpace) via IPC.
It uses the available Unix Socket to communicate. The socket is typically located at `/tmp/\(aeroSpaceAppId)-\(unixUserName).sock` ([see](https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/AppBundle/server.swift#L9)).

## Features

As of now, this library only covers the functionality necessary for implementing
[aerospace-marks](https://github.com/cristianoliveira/aerospace-marks) and [aerospace-scratchpad](https://github.com/cristianoliveira/aerospace-scratchpad) which is:

    - Windows
        - Get all windows
        - Get focused window
        - Get window by ID
        - Get windows by workspace
        - Move window to workspace
        - Set window layout
 
    - Workspaces
        - Get focused workspace
        - Move workspace to another workspace

For the remaining functionality, this library exposes [a SocketClient interface](https://github.com/cristianoliveira/aerospace-ipc/blob/b02bec38820a70895785880b60002a4cf6d5a09b/pkg/client/socket.go#L34), which allows you to send raw commands and receive responses in pure JSON format.

See [documentation](https://pkg.go.dev/github.com/cristianoliveira/aerospace-ipc) for the full list of available methods.

If you need a specific functionality that is not yet implemented, feel free to open an issue or a pull request.

## When to use this library

If you are creating an extension that relies a lot on querying the AeroSpace WM. Because it uses IPC (Inter-Process Communication) to communicate directly with the AeroSpace Unix socket, just like the built-in AeroSpace CLI. It allows you to avoid repeated process spawning. This approach offers lower latency and better efficiency.

## Installation

To use this library in your Go project, add it as a dependency:

```bash
go get -u github.com/cristianoliveira/aerospace-ipc
```

## Usage

### Example Usage

To use the library, import it into your Go project and create a new AeroSpace connection:

```go
import (
    "errors"
    "fmt"
    "log"

    aerospacecli "github.com/cristianoliveira/aerospace-ipc"
)

func main() {
    client, err := aerospacecli.NewAeroSpaceConnection()
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.CloseConnection()

    // This isn't strictly necessary, but it's a good practice to check the server version
    err = client.Connection().CheckServerVersion()
    if err != nil {
        if error.Is(err, aerospacecli.ErrVersionMismatch) {
            fmt.Printf("[WARN] %s\n", err)
        } else {
            log.Fatalf("Failed to connect: %v", err)
        }
    }

    windows, err := client.GetAllWindows()
    if err != nil {
        log.Fatalf("Failed to get windows: %v", err)
    }

    for _, window := range windows {
        fmt.Println(window)
    }
}
```

See also in [examples](examples) for more detailed usage examples.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
