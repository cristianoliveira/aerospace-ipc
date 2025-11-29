# aerospace-concurrent example

An example demonstrating concurrent usage of the AeroSpace IPC client.

This example demonstrates:
- Creating a client using `aerospace.NewClient()`
- Using the connection interface for low-level commands via `client.Connection().SendCommand()`
- Concurrent access to the client from multiple goroutines
- Thread-safe usage of the AeroSpace client

## Usage

```bash
go run main.go
```

## Notes

The AeroSpace client is designed to be thread-safe and can be safely used from multiple goroutines concurrently. This example shows how to send multiple commands in parallel.
