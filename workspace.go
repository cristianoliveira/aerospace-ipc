package aerospace

import (
	"encoding/json"
	"fmt"
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
//	workspace, err := aerospace.GetFocusedWorkspace()
//	fmworkspace.Prworkspacentln("Workspace:", workspace)
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (a *AeroSpaceWM) GetFocusedWorkspace() (*Workspace, error) {
	response, err := a.Conn.SendCommand(
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
//	err := aerospace.MoveWindowToWorkspace(12345, "my-workspace")
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (a *AeroSpaceWM) MoveWindowToWorkspace(windowID int, workspaceName string) error {
	response, err := a.Conn.SendCommand(
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
