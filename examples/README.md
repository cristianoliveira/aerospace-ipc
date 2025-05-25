# Examples for AeroSpace IPC Library

This directory contains example applications demonstrating how to use the socket client library to interact with the AeroSpace window manager.

## Using Examples as Templates

These examples serve as templates for your own extensions. You can modify them to suit your specific needs or use them as a starting point for developing new features.

## Running the Examples

To run an example, navigate to the example's directory and execute the following command:

```bash
cd examples/aerospace-ls
go run main.go
```

Ensure that the AeroSpace window manager is running and accessible via the Unix socket before executing the examples.

## Available Examples

- **aerospace-ls**: Demonstrates listing all windows and their details.
- **aerospace-custom**: Shows how to create a configured client with custom settings and use focus commands.

Feel free to explore and modify the examples to better understand how the AeroSpace IPC library can be integrated into your projects.
