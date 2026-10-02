package workspaces

import (
	"errors"
	"testing"

	mock_client "github.com/cristianoliveira/aerospace-ipc/internal/mocks"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	"go.uber.org/mock/gomock"
)

func TestSummonWorkspace(t *testing.T) {
	tests := []struct {
		name       string
		workspace  string
		failIfNoop bool
		wantArgs   []string
		stderr     string
	}{
		{
			name:      "workspace name is separated from command flags",
			workspace: "work",
			wantArgs:  []string{"--", "work"},
		},
		{
			name:       "fail if noop flag precedes separator and workspace",
			workspace:  "work",
			failIfNoop: true,
			wantArgs:   []string{"--fail-if-noop", "--", "work"},
		},
		{
			name:      "workspace names with spaces stay a single argument and stderr tip is not failure",
			workspace: "my workspace",
			wantArgs:  []string{"--", "my workspace"},
			stderr:    "Workspace is already visible; use --fail-if-noop to fail",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("summon-workspace", tt.wantArgs).
				Return(&client.Response{StdErr: tt.stderr}, nil)

			err := service.SummonWorkspace(
				SummonWorkspaceArgs{WorkspaceName: tt.workspace},
				SummonWorkspaceOpts{FailIfNoop: tt.failIfNoop},
			)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestSummonWorkspaceRejectsBlankNameWithoutSendingCommand(t *testing.T) {
	for _, tt := range []struct {
		name  string
		value string
	}{
		{name: "empty", value: ""},
		{name: "spaces", value: " "},
		{name: "whitespace", value: "\t\n"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			err := service.SummonWorkspace(SummonWorkspaceArgs{WorkspaceName: tt.value}, SummonWorkspaceOpts{})
			if err == nil {
				t.Fatal("expected blank workspace name error")
			}
		})
	}
}

func TestSummonWorkspaceDashLeadingNameUsesUpstreamParserValidation(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
	service := NewService(mockConn)

	wantErr := errors.New("command failed with exit code 2\nWorkspace names starting with dash are disallowed")
	mockConn.EXPECT().
		SendCommand("summon-workspace", []string{"--", "-work"}).
		Return(nil, wantErr)

	err := service.SummonWorkspace(SummonWorkspaceArgs{WorkspaceName: "-work"}, SummonWorkspaceOpts{})
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected upstream parser error %v, got %v", wantErr, err)
	}
}

func TestSummonWorkspacePropagatesServerErrors(t *testing.T) {
	for _, tt := range []struct {
		name      string
		workspace string
		options   SummonWorkspaceOpts
		wantError string
	}{
		{
			name:      "unknown workspace",
			workspace: "missing",
			wantError: "command failed with exit code 1\nunknown workspace",
		},
		{
			name:      "fail if noop",
			workspace: "work",
			options:   SummonWorkspaceOpts{FailIfNoop: true},
			wantError: "command failed with exit code 1\nWorkspace is already visible",
		},
		{
			name:      "forced monitor assignment",
			workspace: "work",
			wantError: "command failed with exit code 1\nworkspace-to-monitor-force-assignment doesn't allow it",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)
			wantArgs := []string{"--", tt.workspace}
			if tt.options.FailIfNoop {
				wantArgs = []string{"--fail-if-noop", "--", tt.workspace}
			}
			mockConn.EXPECT().
				SendCommand("summon-workspace", wantArgs).
				Return(nil, errors.New(tt.wantError))

			err := service.SummonWorkspace(
				SummonWorkspaceArgs{WorkspaceName: tt.workspace},
				tt.options,
			)
			if err == nil || err.Error() != tt.wantError {
				t.Fatalf("expected SendCommand error %q, got %v", tt.wantError, err)
			}
		})
	}
}

func TestSummonWorkspaceReturnsConnectionError(t *testing.T) {
	t.Run("connection error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)
		wantErr := errors.New("connection failed")

		mockConn.EXPECT().
			SendCommand("summon-workspace", []string{"--", "work"}).
			Return(nil, wantErr)

		err := service.SummonWorkspace(SummonWorkspaceArgs{WorkspaceName: "work"}, SummonWorkspaceOpts{})
		if !errors.Is(err, wantErr) {
			t.Fatalf("expected connection error %v, got %v", wantErr, err)
		}
	})
}
