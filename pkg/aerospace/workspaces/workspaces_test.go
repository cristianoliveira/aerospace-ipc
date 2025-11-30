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
			tt.Run("with window ID", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				windowID := 12345
				workspace := "42"

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							workspace,
							"--window-id", "12345",
						},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspace(workspace, &MoveWindowToWorkspaceOpts{
					WindowID: &windowID,
				})
				if err != nil {
					ttt.Fatalf("unexpected error: %v", err)
				}
			})

			tt.Run("without window ID (focused window)", func(ttt *testing.T) {
				ctrl := gomock.NewController(ttt)
				defer ctrl.Finish()

				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)

				workspace := "42"

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{workspace},
					).
					Return(
						&client.Response{
							StdOut: "",
						},
						nil,
					)

				err := service.MoveWindowToWorkspace(workspace, nil)
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
				workspace := "next"

				mockConn.EXPECT().
					SendCommand(
						"move-node-to-workspace",
						[]string{
							workspace,
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

				err := service.MoveWindowToWorkspace(workspace, &MoveWindowToWorkspaceOpts{
					WindowID:           &windowID,
					FocusFollowsWindow: true,
					FailIfNoop:         true,
					WrapAround:         true,
				})
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
	})
}
