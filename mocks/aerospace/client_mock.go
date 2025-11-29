package aerospace_mock

import (
	windows_mock "github.com/cristianoliveira/aerospace-ipc/mocks/aerospace/windows"
	workspaces_mock "github.com/cristianoliveira/aerospace-ipc/mocks/aerospace/workspaces"
	client_mock "github.com/cristianoliveira/aerospace-ipc/mocks/client"
	"github.com/cristianoliveira/aerospace-ipc/pkg/client"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	Conn client.AeroSpaceConnection

	// Services
	windows    *windows_mock.MockWindowsService
	workspaces *workspaces_mock.MockWorkspacesService
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	conn := client_mock.NewMockAeroSpaceConnection(ctrl)
	windows := windows_mock.NewMockWindowsService(ctrl)
	workspaces := workspaces_mock.NewMockWorkspacesService(ctrl)

	mock := &MockClient{
		Conn: conn,
		windows: windows,
		workspaces: workspaces,
	}

	return mock
}

func (m *MockClient) Windows() *windows_mock.MockWindowsService {
	return m.windows
}

func (m *MockClient) Workspaces() *workspaces_mock.MockWorkspacesService {
	return m.workspaces
}

func (m *MockClient) Connection() client.AeroSpaceConnection { 
	return m.Conn
}

// CloseConnection mocks base method.
func (m *MockClient) CloseConnection() error {
	return nil
}
