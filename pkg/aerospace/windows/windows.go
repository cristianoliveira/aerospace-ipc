package windows

import (
	"encoding/json"
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
)

// Window represents a window managed by the AeroSpace window manager.
//
// See: aerospace list-windows --all --json
//
// Example JSON response:
//
//	[
//	  {
//	    "window-id" : 6231,
//	    "workspace" : "8",
//	    "app-bundle-id" : "com.brave.Browser",
//	    "app-name" : "Brave Browser",
//	  },
//	  {
//	    "window-id" : 10772,
//	    "workspace" : ".scratchpad",
//	    "app-name" : "WhatsApp",
//	    "app-bundle-id" : "net.whatsapp.WhatsApp"
//	  }
//	]
type Window struct {
	WindowID    int    `json:"window-id"`
	WindowTitle string `json:"window-title"`
	AppName     string `json:"app-name"`
	AppBundleID string `json:"app-bundle-id"`
	Workspace   string `json:"workspace"`
}

// String returns a string representation of the Window struct.
//
// It includes the window ID, application name, window title (if available),
//
// Example:
//
//	window := Window{
//	  WindowID:    6231,
//	  AppName:     "Brave Browser",
//	  WindowTitle: "Github Page",
//	  Workspace:   "8",
//	  AppBundleID: "com.brave.Browser",
//	}
//	fmt.Println(window)
//
//	// Output: 6231  | Brave Browser | Github Page | 8 | com.brave.Browser
func (w Window) String() string {
	builder := fmt.Sprintf("%d | %s ", w.WindowID, w.AppName)
	if w.WindowTitle != "" {
		builder += fmt.Sprintf("| %s", w.WindowTitle)
	}
	if w.Workspace != "" {
		builder += fmt.Sprintf(" | %s", w.Workspace)
	}
	if w.AppBundleID != "" {
		builder += fmt.Sprintf(" | %s", w.AppBundleID)
	}

	return builder
}

// Service provides methods to interact with windows in AeroSpaceWM.
type Service struct {
	client client.AeroSpaceConnection
}

// SetFocusArgs contains required arguments for SetFocus.
// Exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set.
type SetFocusArgs struct {
	// WindowID specifies the window ID to focus.
	// Use this to focus a window by its ID.
	WindowID *int

	// Direction specifies the direction to focus: left, down, up, or right.
	// Use this to focus the nearest window in the given direction.
	Direction *string

	// DFSDirection specifies the DFS direction: "dfs-next" or "dfs-prev".
	// Use this to focus the window before or after the current window in depth-first order.
	DFSDirection *string

	// DFSIndex specifies the DFS index (0-based) to focus.
	// Use this to focus a window by its DFS index.
	DFSIndex *int
}

// SetFocusOpts contains optional parameters for SetFocus.
type SetFocusOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	// It may be useful for more reliable scripting.
	IgnoreFloating bool

	// Boundaries defines focus boundaries.
	// Possible values: "workspace" (default), "all-monitors-outer-frame"
	// Only applicable when using Direction or DFSDirection.
	Boundaries *string

	// BoundariesAction defines the behavior when requested to cross the boundary.
	// Possible values: "stop" (default), "fail", "wrap-around-the-workspace", "wrap-around-all-monitors"
	// Only applicable when using Direction or DFSDirection.
	BoundariesAction *string
}

// SetLayoutArgs contains required arguments for SetLayout.
type SetLayoutArgs struct {
	// Layouts can be one or more of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
	// If multiple layouts are provided, finds the first that doesn't describe the currently active layout and applies it.
	// This is useful for toggling between layouts.
	Layouts []string
}

// SetLayoutOpts contains optional parameters for SetLayout.
type SetLayoutOpts struct {
	// WindowID specifies the window ID to set layout for. If not set, the focused window is used.
	WindowID *int
}

// WindowsService defines the interface for window operations in AeroSpaceWM.
type WindowsService interface {
	// GetAllWindows returns all windows currently managed by the window manager.
	GetAllWindows() ([]Window, error)

	// GetAllWindowsByWorkspace returns all windows in a specified workspace.
	GetAllWindowsByWorkspace(workspaceName string) ([]Window, error)

	// GetFocusedWindow returns the currently focused window.
	GetFocusedWindow() (*Window, error)

	// SetFocus sets the focus using the specified method.
	// Exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set in args.
	SetFocus(args SetFocusArgs) error

	// SetFocusWithOpts sets the focus using the specified method with options.
	// Exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set in args.
	// opts must be provided and contains optional parameters.
	SetFocusWithOpts(args SetFocusArgs, opts SetFocusOpts) error

	// SetLayout sets the layout for the focused window.
	SetLayout(args SetLayoutArgs) error

	// SetLayoutWithOpts sets the layout for a window with options.
	// opts must be provided and contains optional parameters.
	SetLayoutWithOpts(args SetLayoutArgs, opts SetLayoutOpts) error
}

// NewService creates a new window service with the given AeroSpace client connection.
func NewService(client client.AeroSpaceConnection) *Service {
	return &Service{client: client}
}

// GetAllWindows returns all windows currently managed by the window manager.
//
// It is equivalent to running the command:
//
//	aerospace list-windows --all --json
//
// The result is returned a list of Window structs.
//
// Usage:
//
//	windows, err := windowService.GetAllWindows()
//	fmt.Println("Windows:", windows)
//	fmt.Println("Error:", err)
func (s *Service) GetAllWindows() ([]Window, error) {
	response, err := s.client.SendCommand(
		"list-windows",
		[]string{
			"--all",
			"--json",
			"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
		},
	)
	if err != nil {
		return nil, err
	}
	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal windows: %w\nOut:%s\nErr:%s",
			err,
			response.StdOut,
			response.StdErr,
		)
	}
	return windows, nil
}

// GetAllWindowsByWorkspace returns all windows in a specified workspace.
//
// It is equivalent to running the command:
//
//	aerospace list-windows --workspace <workspace> --json
//
// The result is returned as a list of Window structs.
//
// Usage:
//
//	windows, err := windowService.GetAllWindowsByWorkspace("my-workspace")
//	fmt.Println("Windows:", windows)
//	fmt.Println("Error:", err)
func (s *Service) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	response, err := s.client.SendCommand(
		"list-windows",
		[]string{
			"--workspace", workspaceName,
			"--json",
			"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
		},
	)
	if err != nil {
		return nil, err
	}

	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal windows: %w\nOut:%s\nErr:%s",
			err,
			response.StdOut,
			response.StdErr,
		)
	}
	return windows, nil
}

// GetFocusedWindow returns the currently focused window.
//
// It is equivalent to running the command:
//
//	aerospace list-windows --focused --json
//
// The result is returned as a Window struct.
//
// Usage:
//
//	window, err := windowService.GetFocusedWindow()
//	fmt.Println("Window:", window)
//	fmt.Println("Error:", err)
func (s *Service) GetFocusedWindow() (*Window, error) {
	response, err := s.client.SendCommand(
		"list-windows",
		[]string{
			"--focused",
			"--json",
			"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
		},
	)
	if err != nil {
		return nil, err
	}

	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to unmarshal windows: %w\nOut:%s\nErr:%s",
			err,
			response.StdOut,
			response.StdErr,
		)
	}
	if len(windows) == 0 {
		return nil, fmt.Errorf("no windows focused found")
	}

	return &windows[0], nil
}

// SetFocus sets the focus using the specified method.
//
// Exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set in args.
//
// It is equivalent to running one of the following commands:
//
//	aerospace focus --window-id <window-id> [--ignore-floating]
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (left|down|up|right)
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (dfs-next|dfs-prev)
//	aerospace focus --dfs-index <dfs-index>
//
// Returns an error if the operation fails or if more than one focus method is specified.
//
// Usage:
//
//	// Focus by window ID
//	windowID := 12345
//	err := windowService.SetFocus(windows.SetFocusArgs{
//	    WindowID: &windowID,
//	})
//
//	// Focus by direction
//	direction := "left"
//	err := windowService.SetFocus(windows.SetFocusArgs{
//	    Direction: &direction,
//	})
//
//	// Focus by DFS direction
//	dfsDir := "dfs-next"
//	err := windowService.SetFocus(windows.SetFocusArgs{
//	    DFSDirection: &dfsDir,
//	})
//
//	// Focus by DFS index
//	dfsIndex := 0
//	err := windowService.SetFocus(windows.SetFocusArgs{
//	    DFSIndex: &dfsIndex,
//	})
func (s *Service) SetFocus(args SetFocusArgs) error {
	return s.SetFocusWithOpts(args, SetFocusOpts{})
}

// SetFocusWithOpts sets the focus using the specified method with options.
//
// Exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set in args.
// opts must be provided and contains optional parameters.
//
// It is equivalent to running one of the following commands:
//
//	aerospace focus --window-id <window-id> [--ignore-floating]
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (left|down|up|right)
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (dfs-next|dfs-prev)
//	aerospace focus --dfs-index <dfs-index>
//
// Returns an error if the operation fails or if more than one focus method is specified.
//
// Usage:
//
//	// Focus window by ID ignoring floating windows
//	windowID := 12345
//	err := windowService.SetFocusWithOpts(windows.SetFocusArgs{
//	    WindowID: &windowID,
//	}, windows.SetFocusOpts{
//	    IgnoreFloating: true,
//	})
//
//	// Focus by direction with all options
//	direction := "left"
//	boundaries := "workspace"
//	action := "wrap-around-the-workspace"
//	err := windowService.SetFocusWithOpts(windows.SetFocusArgs{
//	    Direction: &direction,
//	}, windows.SetFocusOpts{
//	    IgnoreFloating:  true,
//	    Boundaries:      &boundaries,
//	    BoundariesAction: &action,
//	})
//
//	// Focus by DFS direction with options
//	dfsDir := "dfs-next"
//	err := windowService.SetFocusWithOpts(windows.SetFocusArgs{
//	    DFSDirection: &dfsDir,
//	}, windows.SetFocusOpts{
//	    IgnoreFloating:  true,
//	    BoundariesAction: &action,
//	})
func (s *Service) SetFocusWithOpts(args SetFocusArgs, opts SetFocusOpts) error {
	// Validate that exactly one focus method is specified
	count := 0
	if args.WindowID != nil {
		count++
	}
	if args.Direction != nil {
		count++
	}
	if args.DFSDirection != nil {
		count++
	}
	if args.DFSIndex != nil {
		count++
	}

	if count == 0 {
		return fmt.Errorf("exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set")
	}
	if count > 1 {
		return fmt.Errorf("only one of WindowID, Direction, DFSDirection, or DFSIndex can be set")
	}

	var cmdArgs []string
	var errorMsg string

	// Build command arguments based on which focus method is specified
	if args.WindowID != nil {
		cmdArgs = []string{
			"--window-id", fmt.Sprintf("%d", *args.WindowID),
		}
		errorMsg = fmt.Sprintf("failed to focus window with ID %d", *args.WindowID)
	} else if args.Direction != nil {
		cmdArgs = []string{*args.Direction}
		errorMsg = fmt.Sprintf("failed to focus window in direction %s", *args.Direction)
	} else if args.DFSDirection != nil {
		cmdArgs = []string{*args.DFSDirection}
		errorMsg = fmt.Sprintf("failed to focus window using DFS direction %s", *args.DFSDirection)
	} else if args.DFSIndex != nil {
		cmdArgs = []string{
			"--dfs-index", fmt.Sprintf("%d", *args.DFSIndex),
		}
		errorMsg = fmt.Sprintf("failed to focus window with DFS index %d", *args.DFSIndex)
	}

	// Add optional flags
	if opts.IgnoreFloating {
		cmdArgs = append(cmdArgs, "--ignore-floating")
	}
	if opts.Boundaries != nil {
		// Boundaries only apply to Direction and DFSDirection
		if args.Direction != nil || args.DFSDirection != nil {
			cmdArgs = append(cmdArgs, "--boundaries", *opts.Boundaries)
		}
	}
	if opts.BoundariesAction != nil {
		// BoundariesAction only applies to Direction and DFSDirection
		if args.Direction != nil || args.DFSDirection != nil {
			cmdArgs = append(cmdArgs, "--boundaries-action", *opts.BoundariesAction)
		}
	}

	response, err := s.client.SendCommand("focus", cmdArgs)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("%s\n%s", errorMsg, response.StdErr)
	}

	return nil
}

// SetLayout sets the layout for the focused window.
//
// args.Layouts can be one or more of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
// If multiple layouts are provided, finds the first that doesn't describe the currently active layout and applies it.
// This is useful for toggling between layouts.
//
// It is equivalent to running the command:
//
//	aerospace layout <layout>...
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Set a single layout
//	err := windowService.SetLayout(windows.SetLayoutArgs{
//	    Layouts: []string{"floating"},
//	})
//
//	// Toggle between layouts (order doesn't matter)
//	err := windowService.SetLayout(windows.SetLayoutArgs{
//	    Layouts: []string{"floating", "tiling"},
//	})
//	err := windowService.SetLayout(windows.SetLayoutArgs{
//	    Layouts: []string{"horizontal", "vertical"},
//	})
func (s *Service) SetLayout(args SetLayoutArgs) error {
	return s.SetLayoutWithOpts(args, SetLayoutOpts{})
}

// SetLayoutWithOpts sets the layout for a window with options.
//
// args.Layouts can be one or more of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
// If multiple layouts are provided, finds the first that doesn't describe the currently active layout and applies it.
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace layout <layout>... [--window-id <window-id>]
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Set layout for specific window
//	windowID := 12345
//	err := windowService.SetLayoutWithOpts(windows.SetLayoutArgs{
//	    Layouts: []string{"floating"},
//	}, windows.SetLayoutOpts{
//	    WindowID: &windowID,
//	})
//
//	// Toggle layout for specific window
//	err := windowService.SetLayoutWithOpts(windows.SetLayoutArgs{
//	    Layouts: []string{"floating", "tiling"},
//	}, windows.SetLayoutOpts{
//	    WindowID: &windowID,
//	})
func (s *Service) SetLayoutWithOpts(args SetLayoutArgs, opts SetLayoutOpts) error {
	if len(args.Layouts) == 0 {
		return fmt.Errorf("at least one layout must be provided")
	}

	cmdArgs := make([]string, 0, len(args.Layouts)+2)
	cmdArgs = append(cmdArgs, args.Layouts...)

	if opts.WindowID != nil {
		cmdArgs = append(cmdArgs, "--window-id", fmt.Sprintf("%d", *opts.WindowID))
	}

	if _, err := s.client.SendCommand("layout", cmdArgs); err != nil {
		return fmt.Errorf(
			"failed to set layout(s) %v\nReason:%w",
			args.Layouts,
			err,
		)
	}

	return nil
}
