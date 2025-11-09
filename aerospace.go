// Package aerospace provides an interface for interacting with the Aerospace Windows Manager.
//
// Deprecated: Import "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
package aerospace

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows"
	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// ErrVersionMismatch indicates that the server version does not match the minimum required version.
var ErrVersionMismatch = exceptions.ErrVersion

// AeroSpaceClient defines the interface for interacting with AeroSpaceWM.
//
// Deprecated: Use aerospace.Client from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
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

// Window represents a window managed by the AeroSpace window manager.
//
// Deprecated: Use windows.Window from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/windows" instead.
type Window = windows.Window

// Workspace represents a workspace in AeroSpaceWM.
//
// Deprecated: Use workspaces.Workspace from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/workspaces" instead.
type Workspace = workspaces.Workspace

// AeroSpaceWM implements the AeroSpaceClient interface.
//
// Deprecated: Use aerospace.AeroSpaceWM from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
type AeroSpaceWM struct {
	Conn client.AeroSpaceConnection
	Impl *aerospace.AeroSpaceWM
}

// GetAllWindows returns all windows currently managed by the window manager.
//
// Deprecated: Use Windows().GetAllWindows() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) GetAllWindows() ([]Window, error) {
	return a.Impl.Windows().GetAllWindows()
}

// GetAllWindowsByWorkspace returns all windows in a specified workspace.
//
// Deprecated: Use Windows().GetAllWindowsByWorkspace() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	return a.Impl.Windows().GetAllWindowsByWorkspace(workspaceName)
}

// GetFocusedWindow returns the currently focused window.
//
// Deprecated: Use Windows().GetFocusedWindow() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) GetFocusedWindow() (*Window, error) {
	return a.Impl.Windows().GetFocusedWindow()
}

// SetFocusByWindowID sets the focus to a window specified by its ID.
//
// Deprecated: Use Windows().SetFocusByWindowID() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) SetFocusByWindowID(windowID int) error {
	return a.Impl.Windows().SetFocusByWindowID(windowID)
}

// GetFocusedWorkspace returns the currently focused workspace.
//
// Deprecated: Use Workspaces().GetFocusedWorkspace() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) GetFocusedWorkspace() (*Workspace, error) {
	return a.Impl.Workspaces().GetFocusedWorkspace()
}

// MoveWindowToWorkspace moves a window to a specified workspace.
//
// Deprecated: Use Workspaces().MoveWindowToWorkspace() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) MoveWindowToWorkspace(windowID int, workspaceName string) error {
	return a.Impl.Workspaces().MoveWindowToWorkspace(windowID, workspaceName)
}

// SetLayout sets the layout for a specified window.
//
// Deprecated: Use Windows().SetLayout() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) SetLayout(windowID int, layout string) error {
	return a.Impl.Windows().SetLayout(windowID, layout)
}

// Connection returns the AeroSpaceConnection
// which allows low-level interaction with the AeroSpace socket.
//
// Deprecated: Use Connection() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) Connection() client.AeroSpaceConnection {
	return a.Conn
}

// CloseConnection closes the AeroSpaceWM connection and releases resources.
//
// Deprecated: Use CloseConnection() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func (a *AeroSpaceWM) CloseConnection() error {
	return a.Impl.CloseConnection()
}

// NewAeroSpaceClient creates a new AeroSpaceClient with the default socket path.
//
// Deprecated: Use aerospace.NewClient() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func NewAeroSpaceClient() (*AeroSpaceWM, error) {
	impl, err := aerospace.NewClient()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		Conn: impl.Connection(),
		Impl: impl,
	}

	return client, nil
}

// AeroSpaceCustomConnectionOpts options for custom connection.
//
// Deprecated: Use aerospace.CustomConnectionOpts from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
type AeroSpaceCustomConnectionOpts = aerospace.CustomConnectionOpts

// NewAeroSpaceCustomClient creates a new AeroSpaceClient with a custom socket path.
//
// Deprecated: Use aerospace.NewCustomClient() from "github.com/cristianoliveira/aerospace-ipc/pkg/aerospace" instead.
func NewAeroSpaceCustomClient(opts AeroSpaceCustomConnectionOpts) (*AeroSpaceWM, error) {
	impl, err := aerospace.NewCustomClient(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to socket\n %w", err)
	}

	client := &AeroSpaceWM{
		Conn: impl.Connection(),
		Impl: impl,
	}

	return client, nil
}
