package aerospace

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestWindows(t *testing.T) {
	t.Run("formatting window as string", func(tt *testing.T) {
		testCases := []struct {
			title    string
			window   Window
			expected string
		}{
			{
				title:    "Basic Window Formatting",
				window:   Window{WindowID: 123, WindowTitle: "Test Window", AppName: "TestApp"},
				expected: "123 | TestApp | Test Window",
			},
			{
				title:    "Basic 2 Window Formatting",
				window:   Window{WindowID: 456, WindowTitle: "Another Window", AppName: "AnotherApp"},
				expected: "456 | AnotherApp | Another Window",
			},
			{
				title:    "Window with Empty App Name",
				window:   Window{WindowID: 789, WindowTitle: "Sample Window", AppName: ""},
				expected: "789 |  | Sample Window",
			},
			{
				title:    "Window with Empty Title",
				window:   Window{WindowID: 101, WindowTitle: "", AppName: "EmptyTitleApp"},
				expected: "101 | EmptyTitleApp ",
			},
			{
				title: "Window with more fields",
				window: Window{
					WindowID:    101,
					WindowTitle: "Another Window",
					AppName:     "EmptyTitleApp",
					AppBundleID: "com.example.app",
					Workspace:   "Workspace1",
				},
				expected: "101 | EmptyTitleApp | Another Window | Workspace1 | com.example.app",
			},
		}
		for _, tc := range testCases {
			t.Run(tc.title, func(t *testing.T) {
				result := tc.window.String()
				if result != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, result)
				}
			})
		}
	})
}

func TestAeroSpaceWindowsHappyPaths(t *testing.T) {
	t.Run("GetAllWindows", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		windowsResponse := []Window{
			{WindowID: 123456, WindowTitle: "Terminal - MyApp", AppName: "MyApp"},
			{WindowID: 789012, WindowTitle: "Web Browser - Example", AppName: "Web Browser"},
		}
		windowsJSON, err := json.Marshal(windowsResponse)
		if err != nil {
			t.Fatalf("failed to marshal windows response: %v", err)
		}
		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(
				&client.Response{
					StdOut: string(windowsJSON),
				},
				nil,
			)

		windows, err := aeroSpaceWM.GetAllWindows()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(windows) != 2 {
			t.Errorf("expected 2 windows, got %d", len(windows))
		}

		if windows[0].String() != "123456 | MyApp | Terminal - MyApp" {
			t.Errorf("wrong window 1, got '%s'", windows[0].String())
		}
		if windows[1].String() != "789012 | Web Browser | Web Browser - Example" {
			t.Errorf("wrong window 2, got '%s'", windows[1].String())
		}
	})

	t.Run("GetAllWindowsByWorkspace", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		windowsResponse := []Window{
			{WindowID: 123456, WindowTitle: "Terminal - MyApp", AppName: "MyApp"},
			{WindowID: 789012, WindowTitle: "Web Browser - Example", AppName: "Web Browser"},
		}
		windowsJSON, err := json.Marshal(windowsResponse)
		if err != nil {
			t.Fatalf("failed to marshal windows response: %v", err)
		}
		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--workspace", "my-workspace",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(
				&client.Response{
					StdOut: string(windowsJSON),
				},
				nil,
			)

		windows, err := aeroSpaceWM.GetAllWindowsByWorkspace("my-workspace")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(windows) != 2 {
			t.Errorf("expected 2 windows, got %d", len(windows))
		}

		if windows[0].String() != "123456 | MyApp | Terminal - MyApp" {
			t.Errorf("wrong window 1, got '%s'", windows[0].String())
		}
		if windows[1].String() != "789012 | Web Browser | Web Browser - Example" {
			t.Errorf("wrong window 2, got '%s'", windows[1].String())
		}
	})

	t.Run("GetFocusedWindow", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		focusedWindowResponse := []Window{
			{WindowID: 123456, WindowTitle: "Focused Window", AppName: "FocusedApp"},
		}
		focusedWindowJSON, err := json.Marshal(focusedWindowResponse)
		if err != nil {
			t.Fatalf("failed to marshal focused window response: %v", err)
		}
		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(
				&client.Response{
					StdOut: string(focusedWindowJSON),
				},
				nil,
			)

		window, err := aeroSpaceWM.GetFocusedWindow()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if window.String() != "123456 | FocusedApp | Focused Window" {
			t.Errorf("wrong focused window, got '%s'", window.String())
		}
	})

	t.Run("SetFocusByWindowID", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		mockConn.EXPECT().
			SendCommand("focus", []string{"--window-id", "123456"}).
			Return(
				&client.Response{
					ServerVersion: "1.0",
					StdOut:        "",
					StdErr:        "",
					ExitCode:      0,
				},
				nil,
			)

		err := aeroSpaceWM.SetFocusByWindowID(123456)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
	// End Happy Paths
}

func TestAeroSpaceWindowsErrorPaths(t *testing.T) {
	t.Run("GetAllWindowsError", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--all",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(nil, fmt.Errorf("connection error"))

		_, err := aeroSpaceWM.GetAllWindows()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("GetAllWindowsByWorkspaceError", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		var workspaceName = "nonexistent-workspace"
		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--workspace", workspaceName,
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(nil, fmt.Errorf("workspace not found"))

		_, err := aeroSpaceWM.GetAllWindowsByWorkspace("nonexistent-workspace")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("GetFocusedWindowError", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(nil, fmt.Errorf("no focused window found"))

		_, err := aeroSpaceWM.GetFocusedWindow()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("GetFocusedWindow no focus", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		focusedWindowResponse := []Window{}
		focusedWindowJSON, err := json.Marshal(focusedWindowResponse)
		if err != nil {
			t.Fatalf("failed to marshal focused window response: %v", err)
		}
		mockConn.EXPECT().
			SendCommand(
				"list-windows",
				[]string{
					"--focused",
					"--json",
					"--format", "%{window-id} %{app-name} %{app-bundle-id} %{workspace}",
				},
			).
			Return(
				&client.Response{
					StdOut: string(focusedWindowJSON),
				},
				nil,
			)

		_, err = aeroSpaceWM.GetFocusedWindow()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("SetFocusByWindowIDError", func(tt *testing.T) {
		ctrl := gomock.NewController(tt)
		defer ctrl.Finish()

		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		aeroSpaceWM := AeroSpaceWM{Conn: mockConn}

		mockConn.EXPECT().
			SendCommand("focus", []string{"--window-id", "123456"}).
			Return(nil, fmt.Errorf("failed to focus window"))

		err := aeroSpaceWM.SetFocusByWindowID(123456)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}
