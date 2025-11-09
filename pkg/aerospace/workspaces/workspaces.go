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

// MoveWindowToWorkspace moves a window to a specified workspace.
//
// It is equivalent to running the command:
//
//	aerospace move-node-to-workspace <workspace> --window-id <window-id>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := workspaceService.MoveWindowToWorkspace(12345, "my-workspace")
//	fmt.Println("Error:", err)
func (s *Service) MoveWindowToWorkspace(windowID int, workspaceName string) error {
	response, err := s.client.SendCommand(
		"move-node-to-workspace",
		[]string{workspaceName, "--window-id", fmt.Sprintf("%d", windowID)},
	)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to move window to workspace: %s", response.StdErr)
	}

	return nil
}
