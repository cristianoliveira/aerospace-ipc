package focus

import (
	"fmt"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestFocusService(t *testing.T) {
	t.Run("Happy path", func(tt *testing.T) {
		t.Run("SetFocusByWindowID", func(ttt *testing.T) {
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

			err := service.SetFocusByWindowID(123456)
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByWindowID with IgnoreFloating", func(ttt *testing.T) {
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

			err := service.SetFocusByWindowID(123456, SetFocusOpts{
				IgnoreFloating: true,
			})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDirection", func(ttt *testing.T) {
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

			err := service.SetFocusByDirection("left")
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDirection with all options", func(ttt *testing.T) {
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

			err := service.SetFocusByDirection("left", SetFocusOpts{
				IgnoreFloating:  true,
				Boundaries:      &boundaries,
				BoundariesAction: &action,
			})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDFS", func(ttt *testing.T) {
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

			err := service.SetFocusByDFS("dfs-next")
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDFS with options", func(ttt *testing.T) {
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

			err := service.SetFocusByDFS("dfs-prev", SetFocusOpts{
				IgnoreFloating:  true,
				BoundariesAction: &action,
			})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDFSIndex", func(ttt *testing.T) {
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

			err := service.SetFocusByDFSIndex(0)
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("FocusBackAndForth", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("focus-back-and-forth", []string{}).
				Return(
					&client.Response{
						ServerVersion: "1.0",
						StdOut:        "",
						StdErr:        "",
						ExitCode:      0,
					},
					nil,
				)

			err := service.FocusBackAndForth()
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetFocusByDirection - all directions", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			directions := []string{"left", "down", "up", "right"}
			for _, dir := range directions {
				ttt.Run(dir, func(tttt *testing.T) {
					mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
					service := NewService(mockConn)

					mockConn.EXPECT().
						SendCommand("focus", []string{dir}).
						Return(
							&client.Response{
								ServerVersion: "1.0",
								StdOut:        "",
								StdErr:        "",
								ExitCode:      0,
							},
							nil,
						)

					err := service.SetFocusByDirection(dir)
					if err != nil {
						tttt.Fatalf("unexpected error: %v", err)
					}
				})
			}
		})

		t.Run("SetFocusByDFS - both directions", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			dfsDirections := []string{"dfs-next", "dfs-prev"}
			for _, dir := range dfsDirections {
				ttt.Run(dir, func(tttt *testing.T) {
					mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
					service := NewService(mockConn)

					mockConn.EXPECT().
						SendCommand("focus", []string{dir}).
						Return(
							&client.Response{
								ServerVersion: "1.0",
								StdOut:        "",
								StdErr:        "",
								ExitCode:      0,
							},
							nil,
						)

					err := service.SetFocusByDFS(dir)
					if err != nil {
						tttt.Fatalf("unexpected error: %v", err)
					}
				})
			}
		})
	})

	t.Run("Error cases", func(tt *testing.T) {
		t.Run("SetFocusByDirection with invalid direction", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			err := service.SetFocusByDirection("invalid")
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
			if err.Error() != `invalid direction "invalid", must be one of: left, down, up, right` {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})

		t.Run("SetFocusByDFS with invalid direction", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			err := service.SetFocusByDFS("invalid")
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
			if err.Error() != `invalid DFS direction "invalid", must be one of: dfs-next, dfs-prev` {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})

		t.Run("Command execution error", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("focus", []string{"--window-id", "123456"}).
				Return(nil, fmt.Errorf("connection error"))

			err := service.SetFocusByWindowID(123456)
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
		})

		t.Run("Command exit code error", func(ttt *testing.T) {
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
						StdErr:        "Window not found",
						ExitCode:      1,
					},
					nil,
				)

			err := service.SetFocusByWindowID(123456)
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
			if err.Error() != "failed to focus window with ID 123456\nWindow not found" {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})

		t.Run("FocusBackAndForth non-zero exit code", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("focus-back-and-forth", []string{}).
				Return(
					&client.Response{
						ServerVersion: "1.0",
						StdOut:        "",
						StdErr:        "no previous window",
						ExitCode:      1,
					},
					nil,
				)

			err := service.FocusBackAndForth()
			if err == nil {
				ttt.Fatal("expected error for non-zero exit code, got nil")
			}
			if err.Error() != "failed to switch focus back and forth: no previous window" {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})

		t.Run("FocusBackAndForth connection error", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("focus-back-and-forth", []string{}).
				Return(nil, fmt.Errorf("connection failed")).
				Times(1)

			err := service.FocusBackAndForth()
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
		})
	})
}

// TestFocusServiceInterface ensures that Service implements FocusService interface.
// This is a compile-time check - if Service doesn't implement all methods, this will fail to compile.
func TestFocusServiceInterface(t *testing.T) {
	var _ FocusService = (*Service)(nil)
}

func TestHelperFunctions(t *testing.T) {
	t.Run("IntPtr", func(tt *testing.T) {
		val := 42
		ptr := IntPtr(val)
		if ptr == nil {
			tt.Fatal("IntPtr returned nil")
		}
		if *ptr != val {
			tt.Errorf("expected %d, got %d", val, *ptr)
		}
	})

	t.Run("StringPtr", func(tt *testing.T) {
		val := "test"
		ptr := StringPtr(val)
		if ptr == nil {
			tt.Fatal("StringPtr returned nil")
		}
		if *ptr != val {
			tt.Errorf("expected %s, got %s", val, *ptr)
		}
	})
}
