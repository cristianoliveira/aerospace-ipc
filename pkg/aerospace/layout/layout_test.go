package layout

import (
	"fmt"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestLayoutService(t *testing.T) {
	t.Run("Happy path", func(tt *testing.T) {
		t.Run("SetLayout single layout", func(ttt *testing.T) {
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

			err := service.SetLayout([]string{"floating"})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetLayout toggle between layouts", func(ttt *testing.T) {
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

			err := service.SetLayout([]string{"floating", "tiling"})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetLayout toggle orientation", func(ttt *testing.T) {
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

			err := service.SetLayout([]string{"horizontal", "vertical"})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetLayout with window ID", func(ttt *testing.T) {
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

			err := service.SetLayout([]string{"floating"}, SetLayoutOpts{
				WindowID: &windowID,
			})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})

		t.Run("SetLayout toggle for specific window", func(ttt *testing.T) {
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

			err := service.SetLayout([]string{"floating", "tiling"}, SetLayoutOpts{
				WindowID: &windowID,
			})
			if err != nil {
				ttt.Fatalf("unexpected error: %v", err)
			}
		})
	})

	t.Run("Error cases", func(tt *testing.T) {
		t.Run("Empty layouts", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			err := service.SetLayout([]string{})
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
			if err.Error() != "at least one layout must be provided" {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})

		t.Run("Command execution error", func(ttt *testing.T) {
			ctrl := gomock.NewController(ttt)
			defer ctrl.Finish()

			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("layout", []string{"floating"}).
				Return(nil, fmt.Errorf("connection error"))

			err := service.SetLayout([]string{"floating"})
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
				SendCommand("layout", []string{"floating"}).
				Return(
					&client.Response{
						ServerVersion: "1.0",
						StdOut:        "",
						StdErr:        "Invalid layout",
						ExitCode:      1,
					},
					nil,
				)

			err := service.SetLayout([]string{"floating"})
			if err == nil {
				ttt.Fatal("expected error, got nil")
			}
			if err.Error() != "failed to set layout(s) [floating]\nInvalid layout" {
				ttt.Errorf("unexpected error message: %v", err)
			}
		})
	})
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
}
