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
//		[
//	   {
//	     "window-id" : 6231,
//	     "workspace" : "8",
//	     "app-bundle-id" : "com.brave.Browser.app.agimnkijcaahngcdmfeangaknmldooml",
//	     "app-name" : "YouTube"
//	   },
//	   {
//	     "window-id" : 10772,
//	     "workspace" : ".scratchpad",
//	     "app-name" : "‎WhatsApp",
//	     "app-bundle-id" : "net.whatsapp.WhatsApp"
//	   }
//		]
type Window struct {
	WindowID    int    `json:"window-id"`
	WindowTitle string `json:"window-title"`
	AppName     string `json:"app-name"`
	AppBundleID string `json:"app-bundle-id"`
	Workspace   string `json:"workspace"`
}

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

func (c *AeroSpaceWM) GetAllWindows() ([]Window, error) {
	response, err := c.Conn.SendCommand(
		"list-windows",
		[]string{
			"--all",
			"--json",
			"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
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

func (c *AeroSpaceWM) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	response, err := c.Conn.SendCommand(
		"list-windows",
		[]string{
			"--workspace", workspaceName,
			"--json",
			"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
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

func (c *AeroSpaceWM) GetFocusedWindow() (*Window, error) {
	response, err := c.Conn.SendCommand(
		"list-windows",
		[]string{
			"--focused",
			"--json",
			"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
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

func (c *AeroSpaceWM) SetFocusByWindowID(windowID int) error {
	response, err := c.Conn.SendCommand("focus", []string{"--window-id", fmt.Sprintf("%d", windowID)})
	if err != nil {
		return err
	}

	if response.ExitCode != 0 {
		return fmt.Errorf("failed to focus window with ID %d\n%s", windowID, response.StdErr)
	}

	return nil
}

// Window layouts related functions

func (a *AeroSpaceWM) SetLayout(windowID int, layout string) error {
	windowStr := fmt.Sprintf("%d", windowID)
	args := []string{layout, "--window-id", windowStr}
	if res, err := a.Conn.SendCommand("layout", args); err != nil {
		fmt.Println("Error setting layout:", res, err)
		return fmt.Errorf("failed to set layout for window %d: %w", windowID, err)
	}

	return nil
}
