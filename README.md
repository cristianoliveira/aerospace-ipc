# aerospace-ipc
[![Go project version](https://badge.fury.io/go/github.com%2Fcristianoliveira%2Faerospace-ipc.svg)](https://badge.fury.io/go/github.com%2Fcristianoliveira%2Faerospace-ipc)
[![Go Reference](https://pkg.go.dev/badge/github.com/cristianoliveira/aerospace-ipc.svg)](https://pkg.go.dev/github.com/cristianoliveira/aerospace-ipc)
[![Quick Checks](https://github.com/cristianoliveira/aerospace-ipc/actions/workflows/on-push.yml/badge.svg)](https://github.com/cristianoliveira/aerospace-ipc/actions/workflows/on-push.yml)

A Go library for interacting with [AeroSpace WM](https://github.com/nikitabobko/AeroSpace) via IPC

*Minimum AeroSpace version:* `v0.20.x` (check [docs/MIGRATION](https://github.com/cristianoliveira/aerospace-ipc/blob/main/docs/MIGRATION-2-to-3.md) for older version)

## Description

This package allows interacting with [AeroSpace WM](https://github.com/nikitabobko/AeroSpace) via IPC.
It uses the available Unix Socket to communicate. The socket is typically located at `/tmp/\(aeroSpaceAppId)-\(unixUserName).sock` ([see](https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/AppBundle/server.swift#L9)).

## Features

As of now, this library only covers the functionality necessary for implementing
[aerospace-marks](https://github.com/cristianoliveira/aerospace-marks) and [aerospace-scratchpad](https://github.com/cristianoliveira/aerospace-scratchpad) which is:

    - Windows Service (`client.Windows()`)
        - Get all windows
        - Get focused window
        - Get windows by workspace
 
    - Workspaces Service (`client.Workspaces()`)
        - Get focused workspace
        - Move window to workspace
        - Move workspace back and forth (switch between focused and previous workspace)
        - Move workspace to monitor (direction-based, order-based, or pattern-based)

    - Focus Service (`client.Focus()`)
        - Set focus by window ID
        - Set focus by direction (left, down, up, right)
        - Set focus by DFS (dfs-next, dfs-prev)
        - Set focus by DFS index

    - Layout Service (`client.Layout()`)
        - Set window layout
        - Toggle between layouts

For the remaining functionality, this library exposes [an AeroSpaceConnection interface](https://github.com/cristianoliveira/aerospace-ipc/blob/main/pkg/client/socket.go#L40), which allows you to send raw commands and receive responses in pure JSON format. Access it via `client.Connection()`.

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

To use the library, import it into your Go project and create a new AeroSpace client:

```go
import (
    "errors"
    "fmt"
    "log"

    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/focus"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/layout"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
)

func main() {
    client, err := aerospace.NewClient()
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer client.CloseConnection()

    // This isn't strictly necessary, but it's a good practice to check the server version
    err = client.Connection().CheckServerVersion()
    if err != nil {
        if errors.Is(err, aerospace.ErrVersionMismatch) {
            fmt.Printf("[WARN] %s\n", err)
        } else {
            log.Fatalf("Failed to connect: %v", err)
        }
    }

    // Use the Windows service to interact with windows
    windows, err := client.Windows().GetAllWindows()
    if err != nil {
        log.Fatalf("Failed to get windows: %v", err)
    }

    for _, window := range windows {
        fmt.Println(window)
    }

    // Use the Workspaces service to interact with workspaces
    workspace, err := client.Workspaces().GetFocusedWorkspace()
    if err != nil {
        log.Fatalf("Failed to get focused workspace: %v", err)
    }
    fmt.Printf("Focused workspace: %s\n", workspace.Workspace)

    // Move window to workspace
    err = client.Workspaces().MoveWindowToWorkspace(workspaces.MoveWindowToWorkspaceArgs{
        WorkspaceName: "my-workspace",
    })
    if err != nil {
        log.Fatalf("Failed to move window: %v", err)
    }

    // Switch between focused and previous workspace
    err = client.Workspaces().MoveBackAndForth()
    if err != nil {
        log.Fatalf("Failed to switch workspace: %v", err)
    }

    // Move workspace to monitor (direction-based)
    err = client.Workspaces().MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
        Direction: "left",
    }, workspaces.MoveWorkspaceToMonitorOpts{})
    if err != nil {
        log.Fatalf("Failed to move workspace to monitor: %v", err)
    }

    // Move workspace to monitor (order-based with wrap-around)
    workspaceName := "my-workspace"
    err = client.Workspaces().MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
        Order: "next",
    }, workspaces.MoveWorkspaceToMonitorOpts{
        Workspace:  &workspaceName,
        WrapAround: true,
    })
    if err != nil {
        log.Fatalf("Failed to move workspace to monitor: %v", err)
    }

    // Move workspace to monitor (pattern-based)
    err = client.Workspaces().MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
        Patterns: []string{"HDMI-1", "DP-1"},
    }, workspaces.MoveWorkspaceToMonitorOpts{})
    if err != nil {
        log.Fatalf("Failed to move workspace to monitor: %v", err)
    }

    // Use the Focus service to set focus
    err = client.Focus().SetFocusByWindowID(12345, focus.SetFocusOpts{
        IgnoreFloating: true,
    })
    if err != nil {
        log.Fatalf("Failed to set focus: %v", err)
    }

    // Use the Layout service to set layout
    err = client.Layout().SetLayout([]string{"floating"})
    if err != nil {
        log.Fatalf("Failed to set layout: %v", err)
    }

    // Toggle between layouts
    err = client.Layout().SetLayout([]string{"floating", "tiling"})
    if err != nil {
        log.Fatalf("Failed to toggle layout: %v", err)
    }
}
```

See also in [examples](examples) for more detailed usage examples.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
