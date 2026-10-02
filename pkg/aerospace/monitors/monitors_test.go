package monitors

import (
	"errors"
	"strings"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestServiceInterface(t *testing.T) {
	var _ MonitorsService = (*Service)(nil)
}

func TestGetFocusedMonitor(t *testing.T) {
	t.Run("returns focused monitor", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{
				"--focused",
				"--json",
				"--format",
				"%{monitor-id} %{monitor-name}",
			}).
			Return(&client.Response{StdOut: `[{"monitor-id":2,"monitor-name":"external"}]`}, nil)

		got, err := service.GetFocusedMonitor()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := &Monitor{MonitorID: 2, MonitorName: "external"}
		if *got != *want {
			t.Fatalf("expected %+v, got %+v", want, got)
		}
	})

	t.Run("returns no-focused-monitor error for empty array", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: "[]"}, nil)

		got, err := service.GetFocusedMonitor()
		if err == nil || !strings.Contains(err.Error(), "no focused monitor found") {
			t.Fatalf("expected no-focused-monitor error, got monitor %+v and error %v", got, err)
		}
	})

	t.Run("returns malformed JSON error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: `[{`}, nil)

		if got, err := service.GetFocusedMonitor(); err == nil {
			t.Fatalf("expected malformed JSON error, got %+v", got)
		}
	})

	t.Run("rejects null instead of an array", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: "null"}, nil)

		if got, err := service.GetFocusedMonitor(); err == nil {
			t.Fatalf("expected JSON array error, got %+v", got)
		}
	})

	t.Run("rejects incomplete or invalid monitor fields", func(t *testing.T) {
		invalidResults := []struct {
			name string
			json string
		}{
			{name: "missing monitor ID", json: `[{"monitor-name":"external"}]`},
			{name: "null monitor ID", json: `[{"monitor-id":null,"monitor-name":"external"}]`},
			{name: "zero monitor ID", json: `[{"monitor-id":0,"monitor-name":"external"}]`},
			{name: "negative monitor ID", json: `[{"monitor-id":-1,"monitor-name":"external"}]`},
			{name: "missing monitor name", json: `[{"monitor-id":2}]`},
			{name: "blank monitor name", json: `[{"monitor-id":2,"monitor-name":"  "}]`},
		}

		for _, tt := range invalidResults {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)
				mockConn.EXPECT().
					SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
					Return(&client.Response{StdOut: tt.json}, nil)

				if got, err := service.GetFocusedMonitor(); err == nil {
					t.Fatalf("expected invalid monitor error, got %+v", got)
				}
			})
		}
	})

	t.Run("rejects multiple focused monitors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: `[{"monitor-id":1,"monitor-name":"one"},{"monitor-id":2,"monitor-name":"two"}]`}, nil)

		if got, err := service.GetFocusedMonitor(); err == nil {
			t.Fatalf("expected multiple-monitor error, got %+v", got)
		}
	})

	t.Run("preserves server error", func(t *testing.T) {
		wantErr := errors.New("command failed with exit code 1")
		assertCommandError(t, wantErr)
	})

	t.Run("preserves connection error", func(t *testing.T) {
		wantErr := errors.New("socket connection lost")
		assertCommandError(t, wantErr)
	})
}

func assertCommandError(t *testing.T, wantErr error) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
	service := NewService(mockConn)
	mockConn.EXPECT().
		SendCommand("list-monitors", []string{"--focused", "--json", "--format", "%{monitor-id} %{monitor-name}"}).
		Return(nil, wantErr)

	got, err := service.GetFocusedMonitor()
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected error %v, got monitor %+v and error %v", wantErr, got, err)
	}
}
