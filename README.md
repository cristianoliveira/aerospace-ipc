# aerospace-ipc

A Go library for interacting with the AeroSpace window manager via Unix socket.

## Description

This package allows interaction with the AeroSpace window manager using a Unix socket. 
The socket is typically located at `/tmp/\(aeroSpaceAppId)-\(unixUserName).sock` ([see](https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/AppBundle/server.swift#L9)).

## Features

- Connect to the AeroSpace window manager via Unix socket.
- Send commands and receive responses in mapped types.
- Send raw commands and receive responses in pure JSON format.
- Manage windows and workspaces programmatically.

## Installation

To use this library in your Go project, add it as a dependency:

```bash
go get -u github.com/cristianoliveira/aerospace-ipc
```

## Usage

### Example Usage

To use the library, import it into your Go project and create a new AeroSpace connection:

```go
import aerospacecli "github.com/cristianoliveira/aerospace-ipc"

func main() {
    client, err := aerospacecli.NewAeroSpaceConnection()
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
}
```

See also in [examples](examples) for more detailed usage examples.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the AeroSpace WM maintainers for their support and documentation.
- Inspired by the need for efficient window management solutions.
- This library is heavily inspired by [go-i3](https://github.com/i3/go-i3)
