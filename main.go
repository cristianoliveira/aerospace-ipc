package aerospace

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// ErrVersionMismatch indicates that the server version does not match the minimum required version.
var ErrVersionMismatch = exceptions.ErrVersion

// AeroSpaceClient defines the interface for interacting with AeroSpaceWM.
type AeroSpaceClient interface {
	// Windows Methods

	// GetAllWindows returns all windows currently managed by the window manager.
	GetAllWindows() ([]Window, error)

	// GetAllWindowsByWorkspace returns all windows in a specified workspace.
	GetAllWindowsByWorkspace(workspaceName string) ([]Window, error)

	// GetFocusedWindow returns the currently focused window.
	GetFocusedWindow() (*Window, error)

	// SetFocusByWindowID sets the focus to a window specified by its ID.
	SetFocusByWindowID(windowID int) error

	// GetFocusedWorkspace returns the currently focused workspace.
	GetFocusedWorkspace() (*Workspace, error)

	// MoveWindowToWorkspace moves a window to a specified workspace.
	MoveWindowToWorkspace(windowID int, workspaceName string) error

	// Layout Methods

	// SetLayout sets the layout for a specified window.
	SetLayout(windowID int, layout string) error

	// Low-Level Client Methods

	// Connection returns the AeroSpaceWM client.
	//
	// Returns the AeroSpaceSocketConn interface for further operations.
	Connection() client.AeroSpaceConnection

	// CloseConnection closes the AeroSpaceWM connection and releases resources.
	//
	// Returns an error if the operation fails.
	CloseConnection() error
}

// AeroSpaceWM implements the AeroSpaceClient interface.
type AeroSpaceWM struct {
	Conn client.AeroSpaceConnection
}

// Connection returns the AeroSpaceConnection
// which allows low-level interaction with the AeroSpace socket.
func (a *AeroSpaceWM) Connection() client.AeroSpaceConnection {
	if a.Conn == nil {
		panic("ASSERTION: AeroSpaceWM client is not initialized")
	}

	return a.Conn
}

func (a *AeroSpaceWM) CloseConnection() error {
	if a.Conn == nil {
		panic("ASSERTION: AeroSpaceWM client is not initialized")
	}

	return a.Conn.CloseConnection()
}

// NewAeroSpaceClient creates a new AeroSpaceClient with the default socket path.
//
// It checks for environment variable AEROSPACESOCK or uses the default socket path.
//
//	Default: /tmp/bobko.aerospace-<username>.sock
//
// Returns an AeroSpaceWM client or an error if the connection fails.
//
// Usage:
//
//	client, err := aerospace.NewAeroSpaceClient()
//	if err != nil {
//	    log.Fatalf("failed to create AeroSpace client: %v", err)
//	}
//	defer client.CloseConnection()
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func NewAeroSpaceClient() (*AeroSpaceWM, error) {
	conn, err := client.GetDefaultConnector().Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		Conn: conn,
	}

	return client, nil
}

type AeroSpaceCustomConnectionOpts struct {
	// SocketPath is the custom socket path for the AeroSpace connection.
	SocketPath string
}

// NewAeroSpaceCustomClient creates a new AeroSpaceClient with a custom socket path.
//
// It allows specifying a custom socket path and whether to validate the version.
// Returns an AeroSpaceWM client or an error if the connection fails.
// Usage:
//
//	client, err := aerospace.NewAeroSpaceCustomClient(aerospace.AeroSpaceCustomConnectionOpts{
//	    SocketPath:      "/path/to/custom/socket",
//	})
//	if err != nil {
//	    log.Fatalf("failed to create AeroSpace client: %v", err)
//	}
//	defer client.CloseConnection()
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func NewAeroSpaceCustomClient(opts AeroSpaceCustomConnectionOpts) (*AeroSpaceWM, error) {
	if opts.SocketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	connector := &client.AeroSpaceCustomConnector{
		SocketPath:      opts.SocketPath,
	}

	conn, err := connector.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		Conn: conn,
	}

	return client, nil
}
