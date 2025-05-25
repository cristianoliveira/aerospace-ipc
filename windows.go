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
//	    {
//	        "window-id": 123456,
//	        "window-title": "Terminal - MyApp",
//	        "app-name": "MyApp"
//	    },
//	    {
//	        "window-id": 789012,
//	        "window-title": "Web Browser - Example",
//	        "app-name": "Web Browser"
//	    }
//	]
type Window struct {
	WindowID    int    `json:"window-id"`
	WindowTitle string `json:"window-title"`
	AppName     string `json:"app-name"`
}

func (w Window) String() string {
	builder := fmt.Sprintf("%d | %s ", w.WindowID, w.AppName)
	if w.WindowTitle != "" {
		builder += fmt.Sprintf("| %s", w.WindowTitle)
	}

	return builder
}

func (c *AeroSpaceWM) GetAllWindows() ([]Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--all", "--json"})
	if err != nil {
		return nil, err
	}
	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, err
	}
	return windows, nil
}

func (c *AeroSpaceWM) GetAllWindowsByWorkspace(workspaceName string) ([]Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--workspace", workspaceName, "--json"})
	if err != nil {
		return nil, err
	}

	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, err
	}
	return windows, nil
}

func (c *AeroSpaceWM) GetFocusedWindow() (*Window, error) {
	response, err := c.Conn.SendCommand("list-windows", []string{"--focused", "--json"})
	if err != nil {
		return nil, err
	}

	var windows []Window
	err = json.Unmarshal([]byte(response.StdOut), &windows)
	if err != nil {
		return nil, err
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
