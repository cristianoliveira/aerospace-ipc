package aerospace

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// ErrVersionMismatch indicates that the server version does not match the minimum required version.
var ErrVersionMismatch = exceptions.ErrVersion

// Client defines the interface for interacting with AeroSpaceWM.
type Client interface {
	// Windows returns the windows service for interacting with windows.
	Windows() *windows.Service

	// Workspaces returns the workspace service for interacting with workspaces.
	Workspaces() *workspaces.Service

	// Connection returns the AeroSpaceWM client.
	//
	// Returns the AeroSpaceConnection interface for further operations.
	Connection() client.AeroSpaceConnection

	// CloseConnection closes the AeroSpaceWM connection and releases resources.
	//
	// Returns an error if the operation fails.
	CloseConnection() error
}

// AeroSpaceWM implements the Client interface.
type AeroSpaceWM struct {
	conn client.AeroSpaceConnection

	// Services
	windowsService    *windows.Service
	workspacesService *workspaces.Service
}

// Windows returns the windows service for interacting with windows.
func (a *AeroSpaceWM) Windows() *windows.Service {
	if a.windowsService == nil {
		a.windowsService = windows.NewService(a.conn)
	}
	return a.windowsService
}

// Workspaces returns the workspace service for interacting with workspaces.
func (a *AeroSpaceWM) Workspaces() *workspaces.Service {
	if a.workspacesService == nil {
		a.workspacesService = workspaces.NewService(a.conn)
	}
	return a.workspacesService
}

// Connection returns the AeroSpaceConnection
// which allows low-level interaction with the AeroSpace socket.
func (a *AeroSpaceWM) Connection() client.AeroSpaceConnection {
	if a.conn == nil {
		panic("ASSERTION: AeroSpaceWM client is not initialized")
	}

	return a.conn
}

func (a *AeroSpaceWM) CloseConnection() error {
	if a.conn == nil {
		panic("ASSERTION: AeroSpaceWM client is not initialized")
	}

	return a.conn.CloseConnection()
}

// NewClient creates a new Client with the default socket path.
//
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
//
//	Default: /tmp/bobko.aerospace-<username>.sock
//
// Returns an AeroSpaceWM client or an error if the connection fails.
//
// Usage:
//
//	client, err := aerospace.NewClient()
//	if err != nil {
//	    log.Fatalf("failed to create AeroSpace client: %v", err)
//	}
//	defer client.CloseConnection()
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func NewClient() (*AeroSpaceWM, error) {
	conn, err := client.GetDefaultConnector().Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		conn: conn,
	}

	return client, nil
}

type CustomConnectionOpts struct {
	// SocketPath is the custom socket path for the AeroSpace connection.
	SocketPath string
	// ValidateVersion is the version to validate against the AeroSpace server.
	ValidateVersion bool
}

// NewCustomClient creates a new Client with a custom socket path.
//
// It allows specifying a custom socket path and whether to validate the version.
// Returns an AeroSpaceWM client or an error if the connection fails.
// Usage:
//
//	client, err := aerospace.NewCustomClient(aerospace.CustomConnectionOpts{
//	    SocketPath:      "/path/to/custom/socket",
//	    ValidateVersion: true, // Set to true to validate the server version
//	})
//	if err != nil {
//	    log.Fatalf("failed to create AeroSpace client: %v", err)
//	}
//	defer client.CloseConnection()
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func NewCustomClient(opts CustomConnectionOpts) (*AeroSpaceWM, error) {
	if opts.SocketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	connector := &client.AeroSpaceCustomConnector{
		SocketPath:      opts.SocketPath,
		ValidateVersion: opts.ValidateVersion,
	}

	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		conn: conn,
	}

	return client, nil
}
