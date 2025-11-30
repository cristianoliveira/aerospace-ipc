# Migration Guide

This guide helps you migrate from the old API to the new service-based API. If you are 

Minimum version of the AeroSpace socket client required for compatibility

 - AeroSpace 0.15.0 till 0.19.x use <=v0.2.1
 - AeroSpace 0.20.0 onwards use >=v0.3.0

This guide is for migrating from v2.x to v3.x of the Go client library. 

TIP: when migrating using an AI agent, ask them to follow this guide, in my experience it works very well :)  

## Overview

The library has been refactored to use a service-based architecture, providing better organization and clearer API boundaries. This is a breaking change, so you'll need to update your code.

## Key Changes

1. **Import Path**: Changed from root package to `pkg/aerospace`
2. **Client Creation**: Function names changed
3. **Method Access**: Methods are now organized into services
4. **Type Access**: Types are now in their respective service packages

## Migration Steps

### 1. Update Import Path

**Before:**
```go
import "github.com/cristianoliveira/aerospace-ipc"
```

**After:**
```go
import "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
```

### 2. Update Client Creation

**Before:**
```go
client, err := aerospace.NewAeroSpaceClient()
// or
client, err := aerospace.NewAeroSpaceCustomClient(aerospace.AeroSpaceCustomConnectionOpts{
    SocketPath: "/path/to/socket",
})
```

**After:**
```go
client, err := aerospace.NewClient()
// or
client, err := aerospace.NewCustomClient(aerospace.CustomConnectionOpts{
    SocketPath: "/path/to/socket",
})
```

### 3. Update Window Operations

**Before:**
```go
// Get all windows
windows, err := client.GetAllWindows()

// Get focused window
focusedWindow, err := client.GetFocusedWindow()

// Get windows by workspace
windows, err := client.GetAllWindowsByWorkspace("my-workspace")

// Set focus
err := client.SetFocusByWindowID(windowID)

// Set layout
err := client.SetLayout(windowID, "floating")
```

**After:**
```go
// Get all windows
windows, err := client.Windows().GetAllWindows()

// Get focused window
focusedWindow, err := client.Windows().GetFocusedWindow()

// Get windows by workspace
windows, err := client.Windows().GetAllWindowsByWorkspace("my-workspace")

// Set focus (standard)
err := client.Windows().SetFocusByWindowID(windows.SetFocusArgs{
    WindowID: windowID,
})

// Set focus (ignoring floating windows)
err := client.Windows().SetFocusByWindowIDWithOpts(windows.SetFocusArgs{
    WindowID: windowID,
}, windows.SetFocusOpts{
    IgnoreFloating: true,
})

// Set focus by direction (left, down, up, right)
err := client.Windows().SetFocusByDirection(windows.SetFocusByDirectionArgs{
    Direction: "left",
})

// Set focus by direction with options
boundaries := "workspace"
action := "wrap-around-the-workspace"
err := client.Windows().SetFocusByDirectionWithOpts(windows.SetFocusByDirectionArgs{
    Direction: "left",
}, windows.SetFocusByDirectionOpts{
    IgnoreFloating:  true,
    Boundaries:      &boundaries,
    BoundariesAction: &action,
})

// Set focus by DFS (dfs-next, dfs-prev)
err := client.Windows().SetFocusByDFS(windows.SetFocusByDFSArgs{
    Direction: "dfs-next",
})

// Set focus by DFS with options
err := client.Windows().SetFocusByDFSWithOpts(windows.SetFocusByDFSArgs{
    Direction: "dfs-prev",
}, windows.SetFocusByDFSOpts{
    IgnoreFloating: true,
    BoundariesAction: &action,
})

// Set focus by DFS index
err := client.Windows().SetFocusByDFSIndex(windows.SetFocusByDFSIndexArgs{
    DFSIndex: 0,
})

// Set layout for focused window (standard)
err := client.Windows().SetLayout(windows.SetLayoutArgs{
    Layouts: []string{"floating"},
})

// Toggle between layouts (order doesn't matter)
err := client.Windows().SetLayout(windows.SetLayoutArgs{
    Layouts: []string{"floating", "tiling"},
})
err := client.Windows().SetLayout(windows.SetLayoutArgs{
    Layouts: []string{"horizontal", "vertical"},
})

// Set layout for specific window
err := client.Windows().SetLayoutWithOpts(windows.SetLayoutArgs{
    Layouts: []string{"floating"},
}, windows.SetLayoutOpts{
    WindowID: &windowID,
})

// Toggle layout for specific window
err := client.Windows().SetLayoutWithOpts(windows.SetLayoutArgs{
    Layouts: []string{"floating", "tiling"},
}, windows.SetLayoutOpts{
    WindowID: &windowID,
})
```

### 4. Update Workspace Operations

**Before:**
```go
// Get focused workspace
workspace, err := client.GetFocusedWorkspace()

// Move window to workspace
err := client.MoveWindowToWorkspace(windowID, "workspace-name")
```

**After:**
```go
// Get focused workspace
workspace, err := client.Workspaces().GetFocusedWorkspace()

// Move window to workspace (standard - moves focused window)
err := client.Workspaces().MoveWindowToWorkspace(workspaces.MoveWindowToWorkspaceArgs{
    WorkspaceName: "workspace-name",
})

// Move specific window to workspace
windowID := 12345
err := client.Workspaces().MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
    WorkspaceName: "workspace-name",
}, workspaces.MoveWindowToWorkspaceOpts{
    WindowID: &windowID,
})

// Move window with focus follows window
err := client.Workspaces().MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
    WorkspaceName: "workspace-name",
}, workspaces.MoveWindowToWorkspaceOpts{
    WindowID:          &windowID,
    FocusFollowsWindow: true,
})

// Move to next workspace with wrap around
err := client.Workspaces().MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
    WorkspaceName: "next",
}, workspaces.MoveWindowToWorkspaceOpts{
    WrapAround: true,
})
```

### 5. Update Type References

**Before:**
```go
import "github.com/cristianoliveira/aerospace-ipc"

var window aerospace.Window
var workspace aerospace.Workspace
```

**After:**
```go
import (
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
)

var window windows.Window
var workspace workspaces.Workspace
```

### 6. Low-Level Connection Access

The low-level connection access remains the same:

**Before & After:**
```go
// Access low-level connection
conn := client.Connection()

// Send custom command
response, err := conn.SendCommand("list-windows", []string{"--all", "--json"})

// Get socket path
path, err := conn.GetSocketPath()

// Check server version
err := conn.CheckServerVersion()
```

### 7. Error Handling

Error handling remains the same:

**Before & After:**
```go
import (
    "errors"
    "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
)

err := client.Connection().CheckServerVersion()
if errors.Is(err, aerospace.ErrVersionMismatch) {
    // Handle version mismatch
}
```

## Complete Example

**Before:**
```go
package main

import (
    "fmt"
    "log"

    "github.com/cristianoliveira/aerospace-ipc"
)

func main() {
    client, err := aerospace.NewAeroSpaceClient()
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

    workspace, err := client.GetFocusedWorkspace()
    if err != nil {
        log.Fatalf("Failed to get workspace: %v", err)
    }
    fmt.Printf("Focused workspace: %s\n", workspace.Workspace)
}
```

**After:**
```go
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
    defer client.CloseConnection()

    windows, err := client.Windows().GetAllWindows()
    if err != nil {
        log.Fatalf("Failed to get windows: %v", err)
    }

    for _, window := range windows {
        fmt.Println(window)
    }

    workspace, err := client.Workspaces().GetFocusedWorkspace()
    if err != nil {
        log.Fatalf("Failed to get workspace: %v", err)
    }
    fmt.Printf("Focused workspace: %s\n", workspace.Workspace)
}
```

## Benefits of the New API

- **Better Organization**: Methods are grouped by domain (windows, workspaces)
- **Clearer API Boundaries**: Each service has a focused responsibility
- **Easier Testing**: Services can be tested and mocked independently
- **More Intuitive**: Method names clearly indicate which service they belong to

## Need Help?

If you encounter any issues during migration, please:
1. Check the [examples](examples) directory for working code samples
2. Review the [documentation](https://pkg.go.dev/github.com/cristianoliveira/aerospace-ipc/pkg/aerospace)
3. Open an issue on GitHub
