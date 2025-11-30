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
	// opts can be nil to use default options.
	SetFocusByWindowID(windowID int, opts *SetFocusOpts) error

	// SetLayout sets the layout for a window.
	// layout is required and can be one of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
	// opts can be nil to use default options (set layout for focused window).
	SetLayout(layout string, opts *SetLayoutOpts) error
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
// opts can be nil to use default options.
//
// It is equivalent to running the command:
//
//	aerospace focus --window-id <window-id> [--ignore-floating]
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Focus window with default options
//	err := windowService.SetFocusByWindowID(12345, nil)
//
//	// Focus window ignoring floating windows
//	err := windowService.SetFocusByWindowID(12345, &windows.SetFocusOpts{
//	    IgnoreFloating: true,
//	})
func (s *Service) SetFocusByWindowID(windowID int, opts *SetFocusOpts) error {
	args := []string{
		"--window-id", fmt.Sprintf("%d", windowID),
	}

	if opts != nil && opts.IgnoreFloating {
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

// SetLayout sets the layout for a window.
//
// layout is required and can be one of: accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
// opts can be nil to use default options (set layout for focused window).
//
// It is equivalent to running the command:
//
//	aerospace layout <layout> [--window-id <window-id>]
//
// Available layouts:
//
//	accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
//
// Returns an error if the operation fails.
//
// Usage:
//
//	// Set layout for focused window
//	err := windowService.SetLayout("floating", nil)
//
//	// Set layout for specific window
//	windowID := 12345
//	err := windowService.SetLayout("floating", &windows.SetLayoutOpts{
//	    WindowID: &windowID,
//	})
func (s *Service) SetLayout(layout string, opts *SetLayoutOpts) error {
	args := []string{layout}

	if opts != nil && opts.WindowID != nil {
		args = append(args, "--window-id", fmt.Sprintf("%d", *opts.WindowID))
	}

	if _, err := s.client.SendCommand("layout", args); err != nil {
		return fmt.Errorf(
			"failed to set layout '%s'\nReason:%w",
			layout,
			err,
		)
	}

	return nil
}
