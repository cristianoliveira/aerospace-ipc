package focus

import (
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// SetFocusOpts contains optional parameters for focus operations.
type SetFocusOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	// It may be useful for more reliable scripting.
	IgnoreFloating bool

	// Boundaries defines focus boundaries.
	// Used with SetFocusByDirection and SetFocusByDFS.
	// Possible values: "workspace" (default), "all-monitors-outer-frame"
	Boundaries *string

	// BoundariesAction defines the behavior when requested to cross the boundary.
	// Used with SetFocusByDirection and SetFocusByDFS.
	// Possible values: "stop" (default), "fail", "wrap-around-the-workspace", "wrap-around-all-monitors"
	BoundariesAction *string
}

// Service provides methods to interact with focus in AeroSpaceWM.
type Service struct {
	client client.AeroSpaceConnection
}

// FocusService defines the interface for focus operations in AeroSpaceWM.
type FocusService interface {
	// SetFocusByWindowID sets focus to a window specified by its ID.
	SetFocusByWindowID(windowID int, opts ...SetFocusOpts) error

	// SetFocusByDirection sets focus to the nearest window in the given direction.
	SetFocusByDirection(direction string, opts ...SetFocusOpts) error

	// SetFocusByDFS sets focus to the window before or after the current window in depth-first order.
	SetFocusByDFS(direction string, opts ...SetFocusOpts) error

	// SetFocusByDFSIndex sets focus to a window by its DFS index.
	SetFocusByDFSIndex(dfsIndex int) error
}

// NewService creates a new focus service with the given AeroSpace client connection.
func NewService(client client.AeroSpaceConnection) *Service {
	return &Service{client: client}
}

// SetFocusByWindowID sets focus to a window specified by its ID.
//
// It is equivalent to running the command:
//
//	aerospace focus --window-id <window-id> [--ignore-floating]
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus by window ID
//	err := focusService.SetFocusByWindowID(12345)
//
//	// Focus by window ID with options
//	err := focusService.SetFocusByWindowID(12345, focus.SetFocusOpts{
//	    IgnoreFloating: true,
//	})
func (s *Service) SetFocusByWindowID(windowID int, opts ...SetFocusOpts) error {
	cmdArgs := []string{
		"--window-id", fmt.Sprintf("%d", windowID),
	}

	var opt SetFocusOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.IgnoreFloating {
		cmdArgs = append(cmdArgs, "--ignore-floating")
	}

	response, err := s.client.SendCommand("focus", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with ID %d\n%s", windowID, response.StdErr)
	}

	return nil
}

// SetFocusByDirection sets focus to the nearest window in the given direction.
//
// Direction must be one of: "left", "down", "up", "right"
//
// It is equivalent to running the command:
//
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (left|down|up|right)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus by direction
//	err := focusService.SetFocusByDirection("left")
//
//	// Focus by direction with all options
//	boundaries := "workspace"
//	action := "wrap-around-the-workspace"
//	err := focusService.SetFocusByDirection("left", focus.SetFocusOpts{
//	    IgnoreFloating:  true,
//	    Boundaries:      &boundaries,
//	    BoundariesAction: &action,
//	})
func (s *Service) SetFocusByDirection(direction string, opts ...SetFocusOpts) error {
	// Validate direction value
	validDirections := map[string]bool{"left": true, "down": true, "up": true, "right": true}
	if !validDirections[direction] {
		return fmt.Errorf("invalid direction %q, must be one of: left, down, up, right", direction)
	}

	cmdArgs := []string{direction}

	var opt SetFocusOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.IgnoreFloating {
		cmdArgs = append(cmdArgs, "--ignore-floating")
	}
	if opt.Boundaries != nil {
		cmdArgs = append(cmdArgs, "--boundaries", *opt.Boundaries)
	}
	if opt.BoundariesAction != nil {
		cmdArgs = append(cmdArgs, "--boundaries-action", *opt.BoundariesAction)
	}

	response, err := s.client.SendCommand("focus", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window in direction %s\n%s", direction, response.StdErr)
	}

	return nil
}

// SetFocusByDFS sets focus to the window before or after the current window in depth-first order.
//
// Direction must be one of: "dfs-next", "dfs-prev"
//
// It is equivalent to running the command:
//
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (dfs-next|dfs-prev)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus by DFS
//	err := focusService.SetFocusByDFS("dfs-next")
//
//	// Focus by DFS with options
//	action := "wrap-around-the-workspace"
//	err := focusService.SetFocusByDFS("dfs-prev", focus.SetFocusOpts{
//	    IgnoreFloating:  true,
//	    BoundariesAction: &action,
//	})
func (s *Service) SetFocusByDFS(direction string, opts ...SetFocusOpts) error {
	// Validate DFS direction value
	if direction != "dfs-next" && direction != "dfs-prev" {
		return fmt.Errorf("invalid DFS direction %q, must be one of: dfs-next, dfs-prev", direction)
	}

	cmdArgs := []string{direction}

	var opt SetFocusOpts
	if len(opts) > 0 {
		opt = opts[0]
	}

	if opt.IgnoreFloating {
		cmdArgs = append(cmdArgs, "--ignore-floating")
	}
	if opt.Boundaries != nil {
		cmdArgs = append(cmdArgs, "--boundaries", *opt.Boundaries)
	}
	if opt.BoundariesAction != nil {
		cmdArgs = append(cmdArgs, "--boundaries-action", *opt.BoundariesAction)
	}

	response, err := s.client.SendCommand("focus", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window using DFS direction %s\n%s", direction, response.StdErr)
	}

	return nil
}

// SetFocusByDFSIndex sets focus to a window by its DFS index.
//
// It is equivalent to running the command:
//
//	aerospace focus --dfs-index <dfs-index>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus by DFS index
//	err := focusService.SetFocusByDFSIndex(0)
func (s *Service) SetFocusByDFSIndex(dfsIndex int) error {
	cmdArgs := []string{
		"--dfs-index", fmt.Sprintf("%d", dfsIndex),
	}

	response, err := s.client.SendCommand("focus", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with DFS index %d\n%s", dfsIndex, response.StdErr)
	}

	return nil
}

// Helper functions for creating pointers (useful for API usage)

// IntPtr returns a pointer to the given int value.
func IntPtr(v int) *int {
	return &v
}

// StringPtr returns a pointer to the given string value.
func StringPtr(v string) *string {
	return &v
}
