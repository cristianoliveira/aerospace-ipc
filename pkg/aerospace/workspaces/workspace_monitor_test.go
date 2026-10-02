package workspaces

import (
	"errors"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestGetAllWorkspacesWithMonitors(t *testing.T) {
	t.Run("returns workspace to monitor mappings", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)

		mockConn.EXPECT().
			SendCommand("list-workspaces", []string{
				"--all",
				"--json",
				"--format",
				"%{workspace} %{monitor-id}",
			}).
			Return(&client.Response{StdOut: `[{"workspace":"dev","monitor-id":1},{"workspace":".scratchpad","monitor-id":2}]`}, nil)

		got, err := service.GetAllWorkspacesWithMonitors()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		want := []WorkspaceMonitor{
			{Workspace: "dev", MonitorID: 1},
			{Workspace: ".scratchpad", MonitorID: 2},
		}
		if len(got) != len(want) {
			t.Fatalf("expected %d workspace mappings, got %d", len(want), len(got))
		}
		for i := range want {
			if got[i] != want[i] {
				t.Errorf("mapping[%d]: expected %+v, got %+v", i, want[i], got[i])
			}
		}
	})

	t.Run("returns an empty result", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-workspaces", []string{"--all", "--json", "--format", "%{workspace} %{monitor-id}"}).
			Return(&client.Response{StdOut: "[]"}, nil)

		got, err := service.GetAllWorkspacesWithMonitors()
		if err != nil {
			t.Fatalf("unexpected error for empty result: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("expected no workspace mappings, got %+v", got)
		}
	})

	t.Run("returns malformed JSON error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		mockConn.EXPECT().
			SendCommand("list-workspaces", []string{"--all", "--json", "--format", "%{workspace} %{monitor-id}"}).
			Return(&client.Response{StdOut: `[{"workspace":"dev","monitor-id":"not-a-number"}]`}, nil)

		got, err := service.GetAllWorkspacesWithMonitors()
		if err == nil {
			t.Fatalf("expected malformed JSON error, got mappings %+v", got)
		}
	})

	t.Run("propagates server error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("command failed with exit code 1\nlist workspaces failed")
		mockConn.EXPECT().
			SendCommand("list-workspaces", []string{"--all", "--json", "--format", "%{workspace} %{monitor-id}"}).
			Return(nil, wantErr)

		got, err := service.GetAllWorkspacesWithMonitors()
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected server error %v, got mappings %+v and error %v", wantErr, got, err)
		}
	})
}
