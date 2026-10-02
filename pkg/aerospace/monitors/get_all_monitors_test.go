package monitors

import (
	"errors"
	"strings"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestGetAllMonitors(t *testing.T) {
	t.Run("returns all monitors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{
				"--json",
				"--format",
				"%{monitor-id} %{monitor-name}",
			}).
			Return(&client.Response{StdOut: `[{"monitor-id":1,"monitor-name":"built-in"},{"monitor-id":2,"monitor-name":"external"}]`}, nil)

		got, err := service.GetAllMonitors()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []Monitor{
			{MonitorID: 1, MonitorName: "built-in"},
			{MonitorID: 2, MonitorName: "external"},
		}
		if len(got) != len(want) {
			t.Fatalf("expected %d monitors, got %d", len(want), len(got))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("monitor[%d]: expected %+v, got %+v", i, want[i], got[i])
			}
		}
	})

	t.Run("returns an empty array", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: "[]"}, nil)

		got, err := service.GetAllMonitors()
		if err != nil {
			t.Fatalf("unexpected error for empty inventory: %v", err)
		}
		if got == nil || len(got) != 0 {
			t.Fatalf("expected empty non-nil result, got %#v", got)
		}
	})

	t.Run("rejects malformed and non-array responses", func(t *testing.T) {
		invalidResponses := []struct {
			name string
			json string
		}{
			{name: "malformed JSON", json: `[{`},
			{name: "null", json: "null"},
			{name: "object instead of array", json: `{"monitor-id":1,"monitor-name":"built-in"}`},
		}

		for _, tt := range invalidResponses {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)
				mockConn.EXPECT().
					SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
					Return(&client.Response{StdOut: tt.json}, nil)

				if got, err := service.GetAllMonitors(); err == nil {
					t.Fatalf("expected parse error, got monitors %+v", got)
				}
			})
		}
	})

	t.Run("rejects incomplete or invalid monitor fields", func(t *testing.T) {
		invalidResults := []struct {
			name string
			json string
		}{
			{name: "missing monitor ID", json: `[{"monitor-name":"built-in"}]`},
			{name: "null monitor ID", json: `[{"monitor-id":null,"monitor-name":"built-in"}]`},
			{name: "zero monitor ID", json: `[{"monitor-id":0,"monitor-name":"built-in"}]`},
			{name: "negative monitor ID", json: `[{"monitor-id":-1,"monitor-name":"built-in"}]`},
			{name: "missing monitor name", json: `[{"monitor-id":1}]`},
			{name: "null monitor name", json: `[{"monitor-id":1,"monitor-name":null}]`},
			{name: "blank monitor name", json: `[{"monitor-id":1,"monitor-name":"  "}]`},
		}

		for _, tt := range invalidResults {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)
				mockConn.EXPECT().
					SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
					Return(&client.Response{StdOut: tt.json}, nil)

				if got, err := service.GetAllMonitors(); err == nil {
					t.Fatalf("expected invalid monitor error, got %+v", got)
				}
			})
		}
	})

	t.Run("rejects duplicate monitor IDs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{StdOut: `[{"monitor-id":1,"monitor-name":"one"},{"monitor-id":1,"monitor-name":"two"}]`}, nil)

		if got, err := service.GetAllMonitors(); err == nil {
			t.Fatalf("expected duplicate ID error, got %+v", got)
		}
	})

	t.Run("preserves server and connection errors", func(t *testing.T) {
		for _, tt := range []struct {
			name string
			err  error
		}{
			{name: "server error", err: errors.New("command failed with exit code 1")},
			{name: "connection error", err: errors.New("socket connection lost")},
		} {
			t.Run(tt.name, func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
				service := NewService(mockConn)
				mockConn.EXPECT().
					SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
					Return(nil, tt.err)

				if got, err := service.GetAllMonitors(); !errors.Is(err, tt.err) {
					t.Fatalf("expected error %v, got monitors %+v and error %v", tt.err, got, err)
				}
			})
		}
	})

	t.Run("returns nonzero response stderr", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		const serverMessage = "failed to enumerate monitors"
		mockConn.EXPECT().
			SendCommand("list-monitors", []string{"--json", "--format", "%{monitor-id} %{monitor-name}"}).
			Return(&client.Response{ExitCode: 1, StdErr: serverMessage}, nil)

		got, err := service.GetAllMonitors()
		if err == nil || !strings.Contains(err.Error(), serverMessage) {
			t.Fatalf("expected error containing %q, got monitors %+v and error %v", serverMessage, got, err)
		}
	})
}
