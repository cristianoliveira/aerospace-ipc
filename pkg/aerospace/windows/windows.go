package windows

import (
	"encoding/json"
	"fmt"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace/focus"
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
//	    "window-layout" : "floating",
//	    "window-parent-container-layout" : "floating",
//	    "app-bundle-id" : "com.brave.Browser",
//	    "app-name" : "Brave Browser",
//	  },
//	  {
//	    "window-id" : 10772,
//	    "workspace" : ".scratchpad",
//	    "window-layout" : "h_tiles",
//	    "window-parent-container-layout" : "h_tiles",
//	    "app-name" : "WhatsApp",
//	    "app-bundle-id" : "net.whatsapp.WhatsApp"
//	  }
//	]
type Window struct {
	WindowID                    int    `json:"window-id"`
	WindowTitle                 string `json:"window-title"`
	WindowLayout                string `json:"window-layout"`
	WindowParentContainerLayout string `json:"window-parent-container-layout"`
	AppName                     string `json:"app-name"`
	AppBundleID                 string `json:"app-bundle-id"`
	Workspace                   string `json:"workspace"`
}

const formatArguments = "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace} %{window-layout} %{window-parent-container-layout}"

// String returns a string representation of the Window struct.
//
// It includes the window ID, application name, window title (if available),
// window layout, window parent container layout, workspace, and app bundle ID.
//
// Example:
//
//	window := Window{
//	  WindowID:    6231,
//	  AppName:     "Brave Browser",
//	  WindowTitle: "Github Page",
//    WindowLayout: "floating",
//    WindowParentContainerLayout: "floating",
//	  Workspace:   "8",
//	  AppBundleID: "com.brave.Browser",
//	}
//	fmt.Println(window)
//
//	// Output: 6231 | Brave Browser | Github Page | floating | floating | 8 | com.brave.Browser
func (w Window) String() string {
	builder := fmt.Sprintf("%d | %s ", w.WindowID, w.AppName)
	if w.WindowTitle != "" {
		builder += fmt.Sprintf("| %s", w.WindowTitle)
	}
	if w.WindowLayout != "" {
		builder += fmt.Sprintf(" | %s", w.WindowLayout)
	}
	if w.WindowParentContainerLayout != "" {
		builder += fmt.Sprintf(" | %s", w.WindowParentContainerLayout)
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

// SetFocusArgs contains required arguments for SetFocusByWindowID.
type SetFocusArgs struct {
	// WindowID specifies the window ID to focus.
	WindowID int
}

// SetFocusOpts contains optional parameters for SetFocusByWindowID.
type SetFocusOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	// It may be useful for more reliable scripting.
	IgnoreFloating bool
}

// SetFocusByDirectionArgs contains required arguments for SetFocusByDirection.
type SetFocusByDirectionArgs struct {
	// Direction specifies the direction to focus: left, down, up, or right.
	Direction string
}

// SetFocusByDirectionOpts contains optional parameters for SetFocusByDirection.
type SetFocusByDirectionOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	IgnoreFloating bool

	// Boundaries defines focus boundaries.
	// Possible values: "workspace" (default), "all-monitors-outer-frame"
	Boundaries *string

	// BoundariesAction defines the behavior when requested to cross the boundary.
	// Possible values: "stop" (default), "fail", "wrap-around-the-workspace", "wrap-around-all-monitors"
	BoundariesAction *string
}

// SetFocusByDFSArgs contains required arguments for SetFocusByDFS.
type SetFocusByDFSArgs struct {
	// Direction specifies the DFS direction: "dfs-next" or "dfs-prev".
	Direction string
}

// SetFocusByDFSOpts contains optional parameters for SetFocusByDFS.
type SetFocusByDFSOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	IgnoreFloating bool

	// Boundaries defines focus boundaries. Must be "workspace" (the default) for DFS mode.
	// Possible values: "workspace" (default)
	Boundaries *string

	// BoundariesAction defines the behavior when requested to cross the boundary.
	// Possible values: "stop" (default), "fail", "wrap-around-the-workspace"
	BoundariesAction *string
}

// SetFocusByDFSIndexArgs contains required arguments for SetFocusByDFSIndex.
type SetFocusByDFSIndexArgs struct {
	// DFSIndex specifies the DFS index (0-based) to focus.
	DFSIndex int
}

// WindowsService defines the interface for window operations in AeroSpaceWM.
type WindowsService interface {
	// GetAllWindows returns all windows currently managed by the window manager.
	GetAllWindows() ([]Window, error)

	// GetAllWindowsByWorkspace returns all windows in a specified workspace.
	GetAllWindowsByWorkspace(workspaceName string) ([]Window, error)

	// GetFocusedWindow returns the currently focused window.
	GetFocusedWindow() (*Window, error)

	// SetFocusByWindowID sets the focus to a window specified by its ID.
	SetFocusByWindowID(args SetFocusArgs) error

	// SetFocusByWindowIDWithOpts sets the focus to a window specified by its ID with options.
	// opts must be provided and contains optional parameters.
	SetFocusByWindowIDWithOpts(args SetFocusArgs, opts SetFocusOpts) error

	// SetFocusByDirection sets focus to the nearest window in the given direction.
	SetFocusByDirection(args SetFocusByDirectionArgs) error

	// SetFocusByDirectionWithOpts sets focus to the nearest window in the given direction with options.
	// opts must be provided and contains optional parameters.
	SetFocusByDirectionWithOpts(args SetFocusByDirectionArgs, opts SetFocusByDirectionOpts) error

	// SetFocusByDFS sets focus to the window before or after the current window in depth-first order.
	SetFocusByDFS(args SetFocusByDFSArgs) error

	// SetFocusByDFSWithOpts sets focus to the window before or after the current window in depth-first order with options.
	// opts must be provided and contains optional parameters.
	SetFocusByDFSWithOpts(args SetFocusByDFSArgs, opts SetFocusByDFSOpts) error

	// SetFocusByDFSIndex sets focus to a window by its DFS index.
	SetFocusByDFSIndex(args SetFocusByDFSIndexArgs) error
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
			"--format", formatArguments,
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
			"--format", formatArguments,
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
			"--format", formatArguments,
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

// SetFocusByWindowID sets the focus to a window specified by its ID.
//
// It is equivalent to running the command:
//
//	aerospace focus --window-id <window-id>
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := windowService.SetFocusByWindowID(windows.SetFocusArgs{
//	    WindowID: 12345,
//	})
//
// Deprecated: Use client.Focus().SetFocusByWindowID() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByWindowID(args SetFocusArgs) error {
	return s.SetFocusByWindowIDWithOpts(args, SetFocusOpts{})
}

// SetFocusByWindowIDWithOpts sets the focus to a window specified by its ID with options.
//
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace focus --window-id <window-id> [--ignore-floating]
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus window ignoring floating windows
//	err := windowService.SetFocusByWindowIDWithOpts(windows.SetFocusArgs{
//	    WindowID: 12345,
//	}, windows.SetFocusOpts{
//	    IgnoreFloating: true,
//	})
//
// Deprecated: Use client.Focus().SetFocusByWindowID() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByWindowIDWithOpts(args SetFocusArgs, opts SetFocusOpts) error {
	focusService := focus.NewService(s.client)
	focusOpts := focus.SetFocusOpts{
		IgnoreFloating: opts.IgnoreFloating,
	}
	return focusService.SetFocusByWindowID(args.WindowID, focusOpts)
}

// SetFocusByDirection sets focus to the nearest window in the given direction.
//
// It is equivalent to running the command:
//
//	aerospace focus (left|down|up|right)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := windowService.SetFocusByDirection(windows.SetFocusByDirectionArgs{
//	    Direction: "left",
//	})
//
// Deprecated: Use client.Focus().SetFocusByWindowID() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByDirection(args SetFocusByDirectionArgs) error {
	return s.SetFocusByDirectionWithOpts(args, SetFocusByDirectionOpts{})
}

// SetFocusByDirectionWithOpts sets focus to the nearest window in the given direction with options.
//
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (left|down|up|right)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	boundaries := "workspace"
//	action := "wrap-around-the-workspace"
//	err := windowService.SetFocusByDirectionWithOpts(windows.SetFocusByDirectionArgs{
//	    Direction: "left",
//	}, windows.SetFocusByDirectionOpts{
//	    IgnoreFloating:  true,
//	    Boundaries:      &boundaries,
//	    BoundariesAction: &action,
//	})
//
// Deprecated: Use client.Focus().SetFocusByDirection() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByDirectionWithOpts(args SetFocusByDirectionArgs, opts SetFocusByDirectionOpts) error {
	focusService := focus.NewService(s.client)
	focusOpts := focus.SetFocusOpts{
		IgnoreFloating:  opts.IgnoreFloating,
		Boundaries:      opts.Boundaries,
		BoundariesAction: opts.BoundariesAction,
	}
	return focusService.SetFocusByDirection(args.Direction, focusOpts)
}

// SetFocusByDFS sets focus to the window before or after the current window in depth-first order.
//
// It is equivalent to running the command:
//
//	aerospace focus (dfs-next|dfs-prev)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := windowService.SetFocusByDFS(windows.SetFocusByDFSArgs{
//	    Direction: "dfs-next",
//	})
//
// Deprecated: Use client.Focus().SetFocusByWindowID() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByDFS(args SetFocusByDFSArgs) error {
	return s.SetFocusByDFSWithOpts(args, SetFocusByDFSOpts{})
}

// SetFocusByDFSWithOpts sets focus to the window before or after the current window in depth-first order with options.
//
// opts must be provided and contains optional parameters.
//
// It is equivalent to running the command:
//
//	aerospace focus [--ignore-floating] [--boundaries <boundary>] [--boundaries-action <action>] (dfs-next|dfs-prev)
//
// Returns an error if the operation fails.
//
// Usage:
//
//	action := "wrap-around-the-workspace"
//	err := windowService.SetFocusByDFSWithOpts(windows.SetFocusByDFSArgs{
//	    Direction: "dfs-next",
//	}, windows.SetFocusByDFSOpts{
//	    IgnoreFloating:  true,
//	    BoundariesAction: &action,
//	})
//
// Deprecated: Use client.Focus().SetFocusByDFS() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByDFSWithOpts(args SetFocusByDFSArgs, opts SetFocusByDFSOpts) error {
	focusService := focus.NewService(s.client)
	focusOpts := focus.SetFocusOpts{
		IgnoreFloating:  opts.IgnoreFloating,
		Boundaries:      opts.Boundaries,
		BoundariesAction: opts.BoundariesAction,
	}
	return focusService.SetFocusByDFS(args.Direction, focusOpts)
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
//	err := windowService.SetFocusByDFSIndex(windows.SetFocusByDFSIndexArgs{
//	    DFSIndex: 0,
//	})
//
// Deprecated: Use client.Focus().SetFocusByDFSIndex() instead. This method is kept for backward compatibility.
func (s *Service) SetFocusByDFSIndex(args SetFocusByDFSIndexArgs) error {
	focusService := focus.NewService(s.client)
	return focusService.SetFocusByDFSIndex(args.DFSIndex)
}

