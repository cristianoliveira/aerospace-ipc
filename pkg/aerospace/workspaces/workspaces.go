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

// MoveWindowToWorkspaceArgs contains required arguments for MoveWindowToWorkspace.
type MoveWindowToWorkspaceArgs struct {
	// WorkspaceName specifies the workspace name where to move the window.
	// Can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
	// to move to the next or previous workspace.
	WorkspaceName string
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

	// Stdin reads the list of workspaces from stdin.
	// Incompatible with NoStdin.
	Stdin bool

	// NoStdin ignores the list of workspaces from stdin, even if provided.
	// Incompatible with Stdin.
	NoStdin bool
}

// WorkspacesService defines the interface for workspace operations in AeroSpaceWM.
type WorkspacesService interface {
	// GetFocusedWorkspace returns the currently focused workspace.
	GetFocusedWorkspace() (*Workspace, error)

	// MoveWindowToWorkspace moves the focused window to a specified workspace.
	MoveWindowToWorkspace(args MoveWindowToWorkspaceArgs) error

	// MoveWindowToWorkspaceWithOpts moves a window to a specified workspace with options.
	// opts must be provided and contains optional parameters.
	MoveWindowToWorkspaceWithOpts(args MoveWindowToWorkspaceArgs, opts MoveWindowToWorkspaceOpts) error

	// MoveBackAndForth switches between the focused workspace and previously focused workspace.
	MoveBackAndForth() error
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
// args.WorkspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
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
//	err := workspaceService.MoveWindowToWorkspace(workspaces.MoveWindowToWorkspaceArgs{
//	    WorkspaceName: "my-workspace",
//	})
func (s *Service) MoveWindowToWorkspace(args MoveWindowToWorkspaceArgs) error {
	return s.MoveWindowToWorkspaceWithOpts(args, MoveWindowToWorkspaceOpts{})
}

// MoveWindowToWorkspaceWithOpts moves a window to a specified workspace with options.
//
// args.WorkspaceName can be a workspace name (e.g., "42", "terminal") or "next"/"prev"
// to move to the next or previous workspace.
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace move-node-to-workspace [--window-id <window-id>] [--focus-follows-window] [--fail-if-noop] [--wrap-around] [--stdin|--no-stdin] <workspace-name>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Move specific window to workspace
//	err := workspaceService.MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
//	    WorkspaceName: "my-workspace",
//	}, workspaces.MoveWindowToWorkspaceOpts{
//	    WindowID: &windowID,
//	})
//
//	// Move window with focus follows window
//	err := workspaceService.MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
//	    WorkspaceName: "my-workspace",
//	}, workspaces.MoveWindowToWorkspaceOpts{
//	    WindowID:          &windowID,
//	    FocusFollowsWindow: true,
//	})
//
//	// Move to next workspace with wrap around
//	err := workspaceService.MoveWindowToWorkspaceWithOpts(workspaces.MoveWindowToWorkspaceArgs{
//	    WorkspaceName: "next",
//	}, workspaces.MoveWindowToWorkspaceOpts{
//	    WrapAround: true,
//	})
func (s *Service) MoveWindowToWorkspaceWithOpts(args MoveWindowToWorkspaceArgs, opts MoveWindowToWorkspaceOpts) error {
	// Validate incompatible options
	if opts.Stdin && opts.NoStdin {
		return fmt.Errorf("cannot specify both --stdin and --no-stdin options")
	}

	cmdArgs := []string{args.WorkspaceName}

	if opts.WindowID != nil {
		cmdArgs = append(cmdArgs, "--window-id", fmt.Sprintf("%d", *opts.WindowID))
	}
	if opts.FocusFollowsWindow {
		cmdArgs = append(cmdArgs, "--focus-follows-window")
	}
	if opts.FailIfNoop {
		cmdArgs = append(cmdArgs, "--fail-if-noop")
	}
	if opts.WrapAround {
		cmdArgs = append(cmdArgs, "--wrap-around")
	}
	if opts.Stdin {
		cmdArgs = append(cmdArgs, "--stdin")
	}
	if opts.NoStdin {
		cmdArgs = append(cmdArgs, "--no-stdin")
	}

	response, err := s.client.SendCommand("move-node-to-workspace", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to move window to workspace: %s", response.StdErr)
	}

	return nil
}

// MoveBackAndForth switches between the focused workspace and previously focused workspace.
//
// It is equivalent to running the command:
//
//	aerospace workspace-back-and-forth
//
// Unlike focus-back-and-forth, workspace-back-and-forth always succeeds.
// Because unlike windows, workspaces can not be "closed".
// Workspaces are name-addressable objects that are created and destroyed on the fly.
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := workspaceService.MoveBackAndForth()
func (s *Service) MoveBackAndForth() error {
	response, err := s.client.SendCommand("workspace-back-and-forth", []string{})
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to switch workspace back and forth: %s", response.StdErr)
	}

	return nil
}

// MoveWorkspaceToMonitor moves a workspace to a monitor.
//
// Supports three modes:
//  1. Direction-based: Move workspace to monitor in direction relative to the focused monitor (left|down|up|right)
//  2. Order-based: Move workspace to next or previous monitor (next|prev)
//  3. Pattern-based: Move workspace to monitor matching pattern(s)
//
// Exactly one of args.Direction, args.Order, or args.Patterns must be specified.
//
// It is equivalent to running the command:
//
//	aerospace move-workspace-to-monitor [--workspace <workspace>] [--wrap-around] (left|down|up|right)
//	aerospace move-workspace-to-monitor [--workspace <workspace>] [--wrap-around] (next|prev)
//	aerospace move-workspace-to-monitor [--workspace <workspace>] <monitor-pattern>...
//
// Focus follows the focused workspace, so the workspace stays focused.
// The command fails for workspaces that have monitor force assignment.
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Move focused workspace to monitor on the left
//	err := workspaceService.MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
//	    Direction: "left",
//	}, workspaces.MoveWorkspaceToMonitorOpts{})
//
//	// Move specific workspace to next monitor with wrap around
//	workspace := "my-workspace"
//	err := workspaceService.MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
//	    Order: "next",
//	}, workspaces.MoveWorkspaceToMonitorOpts{
//	    Workspace:  &workspace,
//	    WrapAround: true,
//	})
//
//	// Move workspace to monitor matching pattern
//	err := workspaceService.MoveWorkspaceToMonitor(workspaces.MoveWorkspaceToMonitorArgs{
//	    Patterns: []string{"HDMI-1", "DP-1"},
//	}, workspaces.MoveWorkspaceToMonitorOpts{})
func (s *Service) MoveWorkspaceToMonitor(args MoveWorkspaceToMonitorArgs, opts MoveWorkspaceToMonitorOpts) error {
	// Validate that exactly one mode is specified
	modesSet := 0
	if args.Direction != "" {
		modesSet++
	}
	if args.Order != "" {
		modesSet++
	}
	if len(args.Patterns) > 0 {
		modesSet++
	}

	if modesSet == 0 {
		return fmt.Errorf("must specify exactly one of: Direction, Order, or Patterns")
	}
	if modesSet > 1 {
		return fmt.Errorf("cannot specify multiple modes; must specify exactly one of: Direction, Order, or Patterns")
	}

	// Validate direction if specified
	if args.Direction != "" {
		validDirections := map[string]bool{"left": true, "down": true, "up": true, "right": true}
		if !validDirections[args.Direction] {
			return fmt.Errorf("invalid direction %q, must be one of: left, down, up, right", args.Direction)
		}
	}

	// Validate order if specified
	if args.Order != "" {
		if args.Order != "next" && args.Order != "prev" {
			return fmt.Errorf("invalid order %q, must be one of: next, prev", args.Order)
		}
	}

	// Build command arguments
	cmdArgs := []string{}

	// Add optional flags
	if opts.Workspace != nil {
		cmdArgs = append(cmdArgs, "--workspace", *opts.Workspace)
	}
	if opts.WrapAround {
		cmdArgs = append(cmdArgs, "--wrap-around")
	}

	// Add the mode-specific argument(s)
	if args.Direction != "" {
		cmdArgs = append(cmdArgs, args.Direction)
	} else if args.Order != "" {
		cmdArgs = append(cmdArgs, args.Order)
	} else if len(args.Patterns) > 0 {
		cmdArgs = append(cmdArgs, args.Patterns...)
	}

	response, err := s.client.SendCommand("move-workspace-to-monitor", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to move workspace to monitor: %s", response.StdErr)
	}

	return nil
}
