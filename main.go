package aerospace

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/client"
	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
)

// ErrVersionMismatch indicates that the server version does not match the minimum required version.
var ErrVersionMismatch = exceptions.ErrVersion

// AeroSpaceClient defines the interface for interacting with AeroSpaceWM.
type AeroSpaceClient interface {
	// Windows Methods

	// GetAllWindows returns all windows currently managed by the window manager.
	//
	// It is equivalent to running the command:
	//   aerospace list-windows --all --json
	//
	// The result is returned a list of Window structs.
	GetAllWindows() ([]Window, error)

	// GetAllWindowsByWorkspace returns all windows in a specified workspace.
	//
	// It is equivalent to running the command:
	//   aerospace list-windows --workspace <workspace> --json
	//
	// The result is returned as a list of Window structs.
	GetAllWindowsByWorkspace(workspaceName string) ([]Window, error)

	// GetFocusedWindow returns the currently focused window.
	//
	// It is equivalent to running the command:
	//   aerospace list-windows --focused --json
	//
	// The result is returned as a Window struct.
	GetFocusedWindow() (*Window, error)

	// SetFocusByWindowID sets the focus to a window specified by its ID.
	//
	// It is equivalent to running the command:
	//   aerospace focus --window-id <window-id>
	//
	// Returns an error if the operation fails.
	SetFocusByWindowID(windowID int) error

	// GetFocusedWorkspace returns the currently focused workspace.
	//
	// It is equivalent to running the command:
	//   aerospace list-workspaces --focused --json
	//
	// The result is returned as a Workspace struct.
	GetFocusedWorkspace() (*Workspace, error)

	// MoveWindowToWorkspace moves a window to a specified workspace.
	//
	// It is equivalent to running the command:
	//   aerospace move-node-to-workspace <workspace> --window-id <window-id>
	//
	// Returns an error if the operation fails.
	MoveWindowToWorkspace(windowID int, workspaceName string) error

	// Layout Methods
	// SetLayout sets the layout for a specified window.
	//
	// It is equivalent to running the command:
	//   aerospace layout <floating|tiled> --window-id <window-id>
	//
	// Returns an error if the operation fails.
	SetLayout(windowID int, layout string) error

	// Connection Methods

	// Client returns the AeroSpaceWM client.
	//
	// Returns the AeroSpaceSocketConn interface for further operations.
	Client() client.AeroSpaceSocketConn

	// CloseConnection closes the AeroSpaceWM connection and releases resources.
	//
	// Returns an error if the operation fails.
	CloseConnection() error
}

// AeroSpaceWM implements the AeroSpaceClient interface.
type AeroSpaceWM struct {
	Conn client.AeroSpaceSocketConn
}

func (a *AeroSpaceWM) Client() client.AeroSpaceSocketConn {
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
func NewAeroSpaceConnection() (*AeroSpaceWM, error) {
	conn, err := GetDefaultConnector().Connect()
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
	// ValidateVersion is the version to validate against the AeroSpace server.
	ValidateVersion bool
}

// NewAeroSpaceCustomConnection creates a new AeroSpaceClient with a custom socket path.
func NewAeroSpaceCustomConnection(opts AeroSpaceCustomConnectionOpts) (*AeroSpaceWM, error) {
	if opts.SocketPath == "" {
		return nil, fmt.Errorf("socket path cannot be empty")
	}

	connector := &AeroSpaceCustomConnector{
		SocketPath:      opts.SocketPath,
		ValidateVersion: opts.ValidateVersion,
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
