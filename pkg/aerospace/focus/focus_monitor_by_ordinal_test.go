package focus

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestFocusMonitorByOrdinal(t *testing.T) {
	t.Run("focuses the monitor by ordinal", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("focus-monitor", []string{"3"}).
			Return(&client.Response{ExitCode: 0}, nil)

		if err := service.FocusMonitorByOrdinal(3); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("rejects nonpositive ordinals without sending a command", func(t *testing.T) {
		for _, ordinal := range []int{0, -1} {
			t.Run(fmt.Sprintf("ordinal %d", ordinal), func(t *testing.T) {
				ctrl := gomock.NewController(t)
				service := NewService(mock_client.NewMockAeroSpaceConnection(ctrl))

				if err := service.FocusMonitorByOrdinal(ordinal); err == nil {
					t.Fatalf("expected validation error for ordinal %d", ordinal)
				}
			})
		}
	})

	t.Run("preserves command error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("monitor is unavailable")
		mockConn.EXPECT().
			SendCommand("focus-monitor", []string{"3"}).
			Return(nil, wantErr)

		if err := service.FocusMonitorByOrdinal(3); !errors.Is(err, wantErr) {
			t.Fatalf("expected command error %v, got %v", wantErr, err)
		}
	})

	t.Run("preserves connection error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("socket connection lost")
		mockConn.EXPECT().
			SendCommand("focus-monitor", []string{"3"}).
			Return(nil, wantErr)

		if err := service.FocusMonitorByOrdinal(3); !errors.Is(err, wantErr) {
			t.Fatalf("expected connection error %v, got %v", wantErr, err)
		}
	})

	t.Run("returns nonzero response stderr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		const serverMessage = "No monitor matches ordinal 3"
		mockConn.EXPECT().
			SendCommand("focus-monitor", []string{"3"}).
			Return(&client.Response{ExitCode: 1, StdErr: serverMessage}, nil)

		err := service.FocusMonitorByOrdinal(3)
		if err == nil || !strings.Contains(err.Error(), serverMessage) {
			t.Fatalf("expected error containing %q, got %v", serverMessage, err)
		}
	})

	t.Run("treats informational no-op stderr as success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("focus-monitor", []string{"3"}).
			Return(&client.Response{
				ExitCode: 0,
				StdErr:   "Monitor is already focused",
			}, nil)

		if err := service.FocusMonitorByOrdinal(3); err != nil {
			t.Fatalf("informational stderr should not be an error: %v", err)
		}
	})
}
