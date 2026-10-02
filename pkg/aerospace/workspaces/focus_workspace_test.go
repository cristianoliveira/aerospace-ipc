package workspaces

import (
	"errors"
	"strings"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestFocusWorkspace(t *testing.T) {
	t.Run("focuses the named workspace", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("workspace", []string{"--", "dev"}).
			Return(&client.Response{ExitCode: 0}, nil)

		if err := service.FocusWorkspace("dev"); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("treats informational already-focused stderr as success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("workspace", []string{"--", "dev"}).
			Return(&client.Response{
				ExitCode: 0,
				StdErr:   "Workspace 'dev' is already focused.",
			}, nil)

		if err := service.FocusWorkspace("dev"); err != nil {
			t.Fatalf("informational stderr should not be an error: %v", err)
		}
	})

	t.Run("preserves invalid-name server error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("Workspace names are not allowed to contain comma")
		mockConn.EXPECT().
			SendCommand("workspace", []string{"--", "bad,name"}).
			Return(nil, wantErr)

		err := service.FocusWorkspace("bad,name")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected server error %v, got %v", wantErr, err)
		}
	})

	t.Run("preserves connection error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("socket connection lost")
		mockConn.EXPECT().
			SendCommand("workspace", []string{"--", "dev"}).
			Return(nil, wantErr)

		err := service.FocusWorkspace("dev")
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected connection error %v, got %v", wantErr, err)
		}
	})

	t.Run("returns nonzero response stderr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		const serverMessage = "workspace is already focused"
		mockConn.EXPECT().
			SendCommand("workspace", []string{"--", "dev"}).
			Return(&client.Response{
				ExitCode: 1,
				StdErr:   serverMessage,
			}, nil)

		err := service.FocusWorkspace("dev")
		if err == nil || !strings.Contains(err.Error(), serverMessage) {
			t.Fatalf("expected error containing %q, got %v", serverMessage, err)
		}
	})
}
