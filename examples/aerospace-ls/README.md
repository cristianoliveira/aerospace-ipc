# aerospace-ls

A simple example CLI for listing all windows from AeroSpace WM.

This example demonstrates:
- Creating a client using `aerospace.NewClient()`
- Using the Windows service via `client.Windows().GetAllWindows()`
- Accessing the connection via `client.Connection().GetSocketPath()`

## Usage

```bash
go run main.go
```

## Your own extension

Use this example as a starting point to create your own CLI tool that interacts with the AeroSpace window manager. You can extend the functionality by adding more commands or features as needed.
