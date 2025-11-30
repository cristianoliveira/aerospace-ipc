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

// SetFocusOpts contains optional parameters for SetFocusByWindowID.
type SetFocusOpts struct {
	// IgnoreFloating don't perceive floating windows as part of the tree.
	// It may be useful for more reliable scripting.
	IgnoreFloating bool
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

	// SetFocusByWindowID sets the focus to a window specified by its ID.
	SetFocusByWindowID(windowID int) error

	// SetFocusByWindowIDWithOpts sets the focus to a window specified by its ID with options.
	// opts must be provided and contains optional parameters.
	SetFocusByWindowIDWithOpts(windowID int, opts SetFocusOpts) error

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
//	err := windowService.SetFocusByWindowID(12345)
func (s *Service) SetFocusByWindowID(windowID int) error {
	return s.SetFocusByWindowIDWithOpts(windowID, SetFocusOpts{})
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
//	err := windowService.SetFocusByWindowIDWithOpts(12345, windows.SetFocusOpts{
//	    IgnoreFloating: true,
//	})
func (s *Service) SetFocusByWindowIDWithOpts(windowID int, opts SetFocusOpts) error {
	args := []string{
		"--window-id", fmt.Sprintf("%d", windowID),
	}

	if opts.IgnoreFloating {
		args = append(args, "--ignore-floating")
	}

	response, err := s.client.SendCommand("focus", args)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with ID %d\n%s", windowID, response.StdErr)
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
