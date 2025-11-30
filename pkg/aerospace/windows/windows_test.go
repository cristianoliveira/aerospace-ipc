package windows

import (
	"encoding/json"
	"fmt"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestWindow(t *testing.T) {
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

func TestWindowService(t *testing.T) {
	t.Run("Happy path", func(tt *testing.T) {
		t.Run("GetAllWindows", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

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
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(
					&client.Response{
						StdOut: string(windowsJSON),
					},
					nil,
				)

			windows, err := service.GetAllWindows()
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
			service := NewService(mockConn)

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
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(
					&client.Response{
						StdOut: string(windowsJSON),
					},
					nil,
				)

			windows, err := service.GetAllWindowsByWorkspace("my-workspace")
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
			service := NewService(mockConn)

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
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(
					&client.Response{
						StdOut: string(focusedWindowJSON),
					},
					nil,
				)

			window, err := service.GetFocusedWindow()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if window.String() != "123456 | FocusedApp | Focused Window" {
				t.Errorf("wrong focused window, got '%s'", window.String())
			}
		})

		t.Run("SetFocus", func(tt *testing.T) {
			tt.Run("by window ID", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

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

				windowID := 123456
				err := service.SetFocus(SetFocusArgs{
					WindowID: &windowID,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by window ID with options", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("focus", []string{"--window-id", "123456", "--ignore-floating"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				windowID := 123456
				err := service.SetFocusWithOpts(SetFocusArgs{
					WindowID: &windowID,
				}, SetFocusOpts{
					IgnoreFloating: true,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by direction", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("focus", []string{"left"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				direction := "left"
				err := service.SetFocus(SetFocusArgs{
					Direction: &direction,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by direction with all options", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				boundaries := "workspace"
				action := "wrap-around-the-workspace"
				mockConn.EXPECT().
					SendCommand("focus", []string{"left", "--ignore-floating", "--boundaries", "workspace", "--boundaries-action", "wrap-around-the-workspace"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				direction := "left"
				err := service.SetFocusWithOpts(SetFocusArgs{
					Direction: &direction,
				}, SetFocusOpts{
					IgnoreFloating:  true,
					Boundaries:      &boundaries,
					BoundariesAction: &action,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by DFS direction", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("focus", []string{"dfs-next"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				dfsDir := "dfs-next"
				err := service.SetFocus(SetFocusArgs{
					DFSDirection: &dfsDir,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by DFS direction with options", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				action := "wrap-around-the-workspace"
				mockConn.EXPECT().
					SendCommand("focus", []string{"dfs-prev", "--ignore-floating", "--boundaries-action", "wrap-around-the-workspace"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				dfsDir := "dfs-prev"
				err := service.SetFocusWithOpts(SetFocusArgs{
					DFSDirection: &dfsDir,
				}, SetFocusOpts{
					IgnoreFloating:  true,
					BoundariesAction: &action,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("by DFS index", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("focus", []string{"--dfs-index", "0"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				dfsIndex := 0
				err := service.SetFocus(SetFocusArgs{
					DFSIndex: &dfsIndex,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})

		t.Run("SetLayout", func(tt *testing.T) {
			tt.Run("standard single layout (focused window)", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("layout", []string{"floating"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				err := service.SetLayout(SetLayoutArgs{
					Layouts: []string{"floating"},
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("toggle between layouts", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("layout", []string{"floating", "tiling"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				err := service.SetLayout(SetLayoutArgs{
					Layouts: []string{"floating", "tiling"},
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("toggle orientation", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("layout", []string{"horizontal", "vertical"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				err := service.SetLayout(SetLayoutArgs{
					Layouts: []string{"horizontal", "vertical"},
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})

		t.Run("SetLayoutWithOpts", func(tt *testing.T) {
			tt.Run("for specific window", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 123456
				mockConn.EXPECT().
					SendCommand("layout", []string{"floating", "--window-id", "123456"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				err := service.SetLayoutWithOpts(SetLayoutArgs{
					Layouts: []string{"floating"},
				}, SetLayoutOpts{
					WindowID: &windowID,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("toggle layout for specific window", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 123456
				mockConn.EXPECT().
					SendCommand("layout", []string{"floating", "tiling", "--window-id", "123456"}).
					Return(
						&client.Response{
							ServerVersion: "1.0",
							StdOut:        "",
							StdErr:        "",
							ExitCode:      0,
						},
						nil,
					)

				err := service.SetLayoutWithOpts(SetLayoutArgs{
					Layouts: []string{"floating", "tiling"},
				}, SetLayoutOpts{
					WindowID: &windowID,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})

		t.Run("Window titles are properly populated", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			windowsResponse := []Window{
				{WindowID: 12345, WindowTitle: "My Window Title", AppName: "MyApp", AppBundleID: "com.app", Workspace: "1"},
				{WindowID: 67890, WindowTitle: "", AppName: "AnotherApp", AppBundleID: "com.other", Workspace: "2"},
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
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(
					&client.Response{
						StdOut: string(windowsJSON),
					},
					nil,
				)

			windows, err := service.GetAllWindows()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(windows) != 2 {
				t.Errorf("expected 2 windows, got %d", len(windows))
			}

			if windows[0].WindowTitle != "My Window Title" {
				t.Errorf("expected window title 'My Window Title', got '%s'", windows[0].WindowTitle)
			}

			if windows[1].WindowTitle != "" {
				t.Errorf("expected empty window title, got '%s'", windows[1].WindowTitle)
			}
		})
	})

	t.Run("Error cases", func(tt *testing.T) {
		t.Run("GetAllWindowsError", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"list-windows",
					[]string{
						"--all",
						"--json",
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(nil, fmt.Errorf("connection error"))

			_, err := service.GetAllWindows()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("GetAllWindowsByWorkspaceError", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			var workspaceName = "nonexistent-workspace"
			mockConn.EXPECT().
				SendCommand(
					"list-windows",
					[]string{
						"--workspace", workspaceName,
						"--json",
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(nil, fmt.Errorf("workspace not found"))

			_, err := service.GetAllWindowsByWorkspace("nonexistent-workspace")
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("GetFocusedWindowError", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"list-windows",
					[]string{
						"--focused",
						"--json",
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(nil, fmt.Errorf("no focused window found"))

			_, err := service.GetFocusedWindow()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("GetFocusedWindow no focus", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

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
						"--format", "%{window-id} %{window-title} %{app-name} %{app-bundle-id} %{workspace}",
					},
				).
				Return(
					&client.Response{
						StdOut: string(focusedWindowJSON),
					},
					nil,
				)

			_, err = service.GetFocusedWindow()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("SetFocusError", func(tt *testing.T) {
			tt.Run("by window ID", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand("focus", []string{"--window-id", "123456"}).
					Return(nil, fmt.Errorf("failed to focus window"))

				windowID := 123456
				err := service.SetFocus(SetFocusArgs{
					WindowID: &windowID,
				})
				if err == nil {
					ttt.Fatal("expected error, got nil")
				}
			})

			tt.Run("no focus method specified", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				err := service.SetFocus(SetFocusArgs{})
				if err == nil {
					ttt.Fatal("expected error, got nil")
				}
				if err.Error() != "exactly one of WindowID, Direction, DFSDirection, or DFSIndex must be set" {
					ttt.Fatalf("unexpected error message: %v", err)
				}
			})

			tt.Run("multiple focus methods specified", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 123456
				direction := "left"
				err := service.SetFocus(SetFocusArgs{
					WindowID:  &windowID,
					Direction: &direction,
				})
				if err == nil {
					ttt.Fatal("expected error, got nil")
				}
				if err.Error() != "only one of WindowID, Direction, DFSDirection, or DFSIndex can be set" {
					ttt.Fatalf("unexpected error message: %v", err)
				}
			})
		})
	})
}
