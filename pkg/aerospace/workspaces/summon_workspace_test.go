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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
			service := NewService(mockConn)

			mockConn.EXPECT().
				SendCommand("summon-workspace", tt.wantArgs).
				Return(&client.Response{}, nil)

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

	mockConn.EXPECT().
		SendCommand("summon-workspace", []string{"--", "-work"}).
		Return(&client.Response{ExitCode: 2, StdErr: "Workspace names starting with dash are disallowed"}, nil)

	err := service.SummonWorkspace(SummonWorkspaceArgs{WorkspaceName: "-work"}, SummonWorkspaceOpts{})
	if err == nil || err.Error() != "failed to summon workspace: Workspace names starting with dash are disallowed" {
		t.Fatalf("expected upstream parser error, got %v", err)
	}
}

func TestSummonWorkspaceReturnsCommandAndConnectionErrors(t *testing.T) {
	t.Run("server error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := mock_client.NewMockAeroSpaceConnection(ctrl)
		service := NewService(mockConn)

		mockConn.EXPECT().
			SendCommand("summon-workspace", []string{"--", "missing"}).
			Return(&client.Response{ExitCode: 1, StdErr: "unknown workspace"}, nil)

		err := service.SummonWorkspace(SummonWorkspaceArgs{WorkspaceName: "missing"}, SummonWorkspaceOpts{})
		if err == nil || err.Error() != "failed to summon workspace: unknown workspace" {
			t.Fatalf("expected server error, got %v", err)
		}
	})

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
