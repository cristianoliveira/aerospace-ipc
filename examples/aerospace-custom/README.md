# aerospace-custom example

A simple example demonstrating how to customize the socket client.

This example demonstrates:
- Creating a custom client with a specific socket path using `aerospace.NewCustomClient()`
- Using the Windows service to list windows via `client.Windows().GetAllWindows()`
- Using the Windows service to set focus via `client.Windows().SetFocusByWindowID()`
- Interactive window selection and focusing

## Usage

```bash
go run main.go
```
