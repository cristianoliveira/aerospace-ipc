package workspaces

import (
	"encoding/json"
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// Workspace represents a workspace in AeroSpaceWM.
//
// See: aerospace list-workspaces --all --json
//
// Example JSON response:
//
//	[
//	  {
//	    "workspace": "42",
//	  },
//	  {
//	    "workspace": "terminal",
//	  }
//	]
type Workspace struct {
	Workspace string `json:"workspace"`
}

// Service provides methods to interact with workspaces in AeroSpaceWM.
type Service struct {
	client client.AeroSpaceConnection
}

// MoveWindowToWorkspaceOpts contains optional parameters for MoveWindowToWorkspace.
type MoveWindowToWorkspaceOpts struct {
	// WindowID specifies the window ID to move. If not set, the focused window is moved.
	WindowID *int

	// FocusFollowsWindow makes the window receive focus after moving.
	// This is a shortcut for manually running aerospace-workspace/aerospace-focus
	// after move-node-to-workspace successful execution.
	FocusFollowsWindow bool

	// FailIfNoop exits with non-zero code if moving the window to a workspace
	// it already belongs to.
	FailIfNoop bool

	// WrapAround makes it possible to jump between first and last workspaces
	// when using "next" or "prev" as workspace name.
	WrapAround bool
}

// WorkspacesService defines the interface for workspace operations in AeroSpaceWM.
type WorkspacesService interface {
	// GetFocusedWorkspace returns the currently focused workspace.
	GetFocusedWorkspace() (*Workspace, error)

	// MoveWindowToWorkspace moves the focused window to a specified workspace.
	// workspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
	// to move to the next or previous workspace.
	MoveWindowToWorkspace(workspaceName string) error

	// MoveWindowToWorkspaceWithOpts moves a window to a specified workspace with options.
	// workspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
	// to move to the next or previous workspace.
	// opts must be provided and contains optional parameters.
	MoveWindowToWorkspaceWithOpts(workspaceName string, opts MoveWindowToWorkspaceOpts) error
}

// NewService creates a new workspace service with the given AeroSpace client connection.
func NewService(client client.AeroSpaceConnection) *Service {
	return &Service{client: client}
}

// GetFocusedWorkspace returns the currently focused workspace.
//
// It is equivalent to running the command:
//
//	aerospace list-workspaces --focused --json
//
// The result differs from the `list-workspaces` command by only returning
// the focused workspace.
//
// Usage:
//
//	workspace, err := workspaceService.GetFocusedWorkspace()
//	fmt.Println("Workspace:", workspace)
//	fmt.Println("Error:", err)
func (s *Service) GetFocusedWorkspace() (*Workspace, error) {
	response, err := s.client.SendCommand(
		"list-workspaces",
		[]string{
			"--focused",
			"--json",
		},
	)
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	err = json.Unmarshal([]byte(response.StdOut), &workspaces)
	if err != nil {
		return nil, err
	}
	if len(workspaces) == 0 {
		return nil, fmt.Errorf("no workspace focused found")
	}

	return &workspaces[0], nil
}

// MoveWindowToWorkspace moves the focused window to a specified workspace.
//
// workspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
// to move to the next or previous workspace.
//
// It is equivalent to running the command:
//
//	aerospace move-node-to-workspace <workspace-name>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := workspaceService.MoveWindowToWorkspace("my-workspace")
func (s *Service) MoveWindowToWorkspace(workspaceName string) error {
	return s.MoveWindowToWorkspaceWithOpts(workspaceName, MoveWindowToWorkspaceOpts{})
}

// MoveWindowToWorkspaceWithOpts moves a window to a specified workspace with options.
//
// workspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
// to move to the next or previous workspace.
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace move-node-to-workspace [--window-id <window-id>] [--focus-follows-window] [--fail-if-noop] [--wrap-around] <workspace-name>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Move specific window to workspace
//	err := workspaceService.MoveWindowToWorkspaceWithOpts("my-workspace", workspaces.MoveWindowToWorkspaceOpts{
//	    WindowID: &windowID,
//	})
//
//	// Move window with focus follows window
//	err := workspaceService.MoveWindowToWorkspaceWithOpts("my-workspace", workspaces.MoveWindowToWorkspaceOpts{
//	    WindowID:          &windowID,
//	    FocusFollowsWindow: true,
//	})
//
//	// Move to next workspace with wrap around
//	err := workspaceService.MoveWindowToWorkspaceWithOpts("next", workspaces.MoveWindowToWorkspaceOpts{
//	    WrapAround: true,
//	})
func (s *Service) MoveWindowToWorkspaceWithOpts(workspaceName string, opts MoveWindowToWorkspaceOpts) error {
	args := []string{workspaceName}

	if opts.WindowID != nil {
		args = append(args, "--window-id", fmt.Sprintf("%d", *opts.WindowID))
	}
	if opts.FocusFollowsWindow {
		args = append(args, "--focus-follows-window")
	}
	if opts.FailIfNoop {
		args = append(args, "--fail-if-noop")
	}
	if opts.WrapAround {
		args = append(args, "--wrap-around")
	}

	response, err := s.client.SendCommand("move-node-to-workspace", args)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to move window to workspace: %s", response.StdErr)
	}

	return nil
}
