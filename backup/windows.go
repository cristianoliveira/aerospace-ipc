package aerospace

import (
	"encoding/json"
	"fmt"
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
//	windows, err := aerospace.GetAllWindows()
//	fmt.Println("Windows:", windows)
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (c *AeroSpaceWM) GetAllWindows() ([]Window, error) {
	response, err := c.Conn.SendCommand(
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
//	windows, err := aerospace.GetAllWindowsByWorkspace("my-workspace")
//	fmt.Println("Windows:", windows)
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (c *AeroSpaceWM) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	response, err := c.Conn.SendCommand(
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
//	window, err := aerospace.GetAllFocusedWindow()
//	fmt.Println("Window:", window)
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (c *AeroSpaceWM) GetFocusedWindow() (*Window, error) {
	response, err := c.Conn.SendCommand(
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
//	err := aerospace.SetFocusByWindowID(12345)
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (c *AeroSpaceWM) SetFocusByWindowID(windowID int) error {
	response, err := c.Conn.SendCommand(
		"focus",
		[]string{
			"--window-id", fmt.Sprintf("%d", windowID),
		},
	)
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with ID %d\n%s", windowID, response.StdErr)
	}

	return nil
}

// Layout Methods

// SetLayout sets the layout for a specified window.
//
// It is equivalent to running the command:
//
//	aerospace layout <layout> --window-id <window-id>
//
// Available layouts:
//
//	accordion|tiles|horizontal|vertical|h_accordion|v_accordion|h_tiles|v_tiles|tiling|floating
//
// Returns an error if the operation fails.
//
// Usage:
//
//	err := aerospace.SetLayout(12345, "floating")
//	fmt.Println("Error:", err)
//
// More:
// https://github.com/cristianoliveira/aerospace-ipc/tree/main/examples
func (a *AeroSpaceWM) SetLayout(windowID int, layout string) error {
	windowStr := fmt.Sprintf("%d", windowID)
	if _, err := a.Conn.SendCommand(
		"layout",
		[]string{
			layout,
			"--window-id", windowStr,
		},
	); err != nil {
		return fmt.Errorf(
			"failed to set layout '%s' for window %d\nReason:%w",
			layout,
			windowID,
			err,
		)
	}

	return nil
}
