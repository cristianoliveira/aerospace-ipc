package workspaces

import (
	"encoding/json"
	"fmt"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestWorkspaceService(t *testing.T) {
	t.Run("Happy path", func(tt *testing.T) {
		t.Run("GetFocusedWorkspace", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			workspaces := []Workspace{
				{Workspace: "42"},
			}

			dataJSON, err := json.Marshal(workspaces)
			if err != nil {
				t.Fatalf("failed to marshal windows response: %v", err)
			}

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(
					&client.Response{
						StdOut: string(dataJSON),
					},
					nil,
				)

			workspace, err := service.GetFocusedWorkspace()
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if workspace.Workspace != "42" {
				t.Fatalf("expected workspace '42', got '%s'", workspace.Workspace)
			}
		})

		t.Run("MoveWindowToWorkspace", func(tt *testing.T) {
			tt.Run("standard (focused window)", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{"42"},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspace(MoveWindowToWorkspaceArgs{
					WorkspaceName: "42",
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})

		t.Run("MoveWindowToWorkspaceWithOpts", func(tt *testing.T) {
			tt.Run("with window ID", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 12345
				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							"42",
							"--window-id", "12345",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
					WorkspaceName: "42",
				}, MoveWindowToWorkspaceOpts{
					WindowID: &windowID,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("with all options", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 12345
				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							"next",
							"--window-id", "12345",
							"--focus-follows-window",
							"--fail-if-noop",
							"--wrap-around",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
					WorkspaceName: "next",
				}, MoveWindowToWorkspaceOpts{
					WindowID:           &windowID,
					FocusFollowsWindow: true,
					FailIfNoop:         true,
					WrapAround:         true,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("with stdin option", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							"42",
							"--stdin",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
					WorkspaceName: "42",
				}, MoveWindowToWorkspaceOpts{
					Stdin: true,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("with no-stdin option", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							"42",
							"--no-stdin",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
					WorkspaceName: "42",
				}, MoveWindowToWorkspaceOpts{
					NoStdin: true,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("with prev workspace", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							"prev",
							"--wrap-around",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
					WorkspaceName: "prev",
				}, MoveWindowToWorkspaceOpts{
					WrapAround: true,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})

		t.Run("MoveBackAndForth", func(tt *testing.T) {
			tt.Run("successful switch", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				mockConn.EXPECT().
					SendCommand(
						"workspace-back-and-forth",
						[]string{},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveBackAndForth()
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})
		})
	})

	t.Run("Error cases", func(tt *testing.T) {
		t.Run("GetFocusedWorkspace return error", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(nil, fmt.Errorf("no focused workspace found")).
				Times(1)

			_, err := service.GetFocusedWorkspace()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("GetFocusedWorkspace return empty", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"list-workspaces",
					[]string{
						"--focused",
						"--json",
					},
				).
				Return(&client.Response{StdOut: "[]"}, nil).
				Times(1)

			_, err := service.GetFocusedWorkspace()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})

		t.Run("MoveWindowToWorkspaceWithOpts incompatible options", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
				WorkspaceName: "42",
			}, MoveWindowToWorkspaceOpts{
				Stdin:   true,
				NoStdin: true,
			})
			if err == nil {
				t.Fatal("expected error for incompatible options, got nil")
			}
			if err.Error() != "cannot specify both --stdin and --no-stdin options" {
				t.Fatalf("expected specific error message, got: %v", err)
			}
		})

		t.Run("MoveWindowToWorkspaceWithOpts non-zero exit code", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"move-node-to-workspace",
					[]string{"42"},
				).
				Return(
					&client.Response{
						ExitCode: 1,
						StdErr:   "window not found",
					},
					nil,
				)

			err := service.MoveWindowToWorkspaceWithOpts(MoveWindowToWorkspaceArgs{
				WorkspaceName: "42",
			}, MoveWindowToWorkspaceOpts{})
			if err == nil {
				t.Fatal("expected error for non-zero exit code, got nil")
			}
			if err.Error() != "failed to move window to workspace: window not found" {
				t.Fatalf("expected specific error message, got: %v", err)
			}
		})

		t.Run("MoveBackAndForth non-zero exit code", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"workspace-back-and-forth",
					[]string{},
				).
				Return(
					&client.Response{
						ExitCode: 1,
						StdErr:   "connection error",
					},
					nil,
				)

			err := service.MoveBackAndForth()
			if err == nil {
				t.Fatal("expected error for non-zero exit code, got nil")
			}
			if err.Error() != "failed to switch workspace back and forth: connection error" {
				t.Fatalf("expected specific error message, got: %v", err)
			}
		})

		t.Run("MoveBackAndForth connection error", func(tt *testing.T) {
			ctrl := gomock.NewController(tt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand(
					"workspace-back-and-forth",
					[]string{},
				).
				Return(nil, fmt.Errorf("connection failed")).
				Times(1)

			err := service.MoveBackAndForth()
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	})
}
