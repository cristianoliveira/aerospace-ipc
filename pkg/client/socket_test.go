package client

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	net_mock "github.com/cristianoliveira/aerospace-ipc/internal/mocks/net"
	"go.uber.org/mock/gomock"
)

func containsSubstring(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

// framedResponse prepends a 4-byte little-endian length prefix to payload.
func framedResponse(payload []byte) []byte {
	frame := make([]byte, 4, 4+len(payload))
	binary.LittleEndian.PutUint32(frame, uint32(len(payload)))
	return append(frame, payload...)
}

// setupSendCommandMock registers mock expectations for one SendCommand call:
// one framed Write, then two Reads (4-byte length prefix + JSON body).
func setupSendCommandMock(mockConn *net_mock.MockConn, responsePayload []byte) {
	gomock.InOrder(
		mockConn.EXPECT().
			Write(gomock.Any()).
			DoAndReturn(func(p []byte) (int, error) {
				return len(p), nil
			}),
		setupResponseReads(mockConn, responsePayload),
	)
}

func setupResponseReads(mockConn *net_mock.MockConn, responsePayload []byte) *gomock.Call {
	mockConn.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(p []byte) (int, error) {
			prefix := framedResponse(responsePayload)[:4]
			return copy(p, prefix), nil
		}).
		Times(1)

	return mockConn.EXPECT().
		Read(gomock.Any()).
		DoAndReturn(func(p []byte) (int, error) {
			return copy(p, responsePayload), nil
		}).
		Times(1)
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	value, exists := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if exists {
			if err := os.Setenv(key, value); err != nil {
				t.Errorf("failed to restore %s: %v", key, err)
			}
			return
		}
		if err := os.Unsetenv(key); err != nil {
			t.Errorf("failed to clear %s: %v", key, err)
		}
	})
}

func TestSendCommand(t *testing.T) {
	unsetEnv(t, "AEROSPACE_WINDOW_ID")
	unsetEnv(t, "AEROSPACE_WORKSPACE")

	t.Run("sends documented request schema", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		responsePayload, err := json.Marshal(Response{StdOut: "1\n", ExitCode: 0})
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		expectedRequest := []byte(`{"args":["list-workspaces","--focused"],"stdin":"","windowId":null,"workspace":null}`)

		gomock.InOrder(
			mockConn.EXPECT().Write(framedResponse(expectedRequest)).Return(len(framedResponse(expectedRequest)), nil),
			setupResponseReads(mockConn, responsePayload),
		)

		connection := &AeroSpaceSocketConnection{Conn: mockConn}
		response, err := connection.SendCommand("list-workspaces", []string{"--focused"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.StdOut != "1\n" {
			t.Fatalf("expected stdout %q, got %q", "1\n", response.StdOut)
		}
	})

	t.Run("completes a partial frame write", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		responsePayload, err := json.Marshal(Response{ExitCode: 0})
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		requestPayload := []byte(`{"args":["list-workspaces","--focused"],"stdin":"","windowId":null,"workspace":null}`)
		frame := framedResponse(requestPayload)

		gomock.InOrder(
			mockConn.EXPECT().Write(frame).Return(3, nil),
			mockConn.EXPECT().Write(frame[3:]).Return(len(frame)-3, nil),
			setupResponseReads(mockConn, responsePayload),
		)

		connection := &AeroSpaceSocketConnection{Conn: mockConn}
		if _, err := connection.SendCommand("list-workspaces", []string{"--focused"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("accepts stderr when exit code is zero", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		responsePayload, err := json.Marshal(Response{StdErr: "warning", ExitCode: 0})
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		setupSendCommandMock(mockConn, responsePayload)

		connection := &AeroSpaceSocketConnection{Conn: mockConn}
		response, err := connection.SendCommand("list-workspaces", []string{"--focused"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if response.StdErr != "warning" {
			t.Fatalf("expected stderr %q, got %q", "warning", response.StdErr)
		}
	})

	t.Run("rejects subscribe mode", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		connection := &AeroSpaceSocketConnection{Conn: net_mock.NewMockConn(ctrl)}

		_, err := connection.SendCommand("subscribe", []string{"--all"})
		if err == nil || !containsSubstring(err.Error(), "subscribe is not supported") {
			t.Fatalf("expected unsupported subscribe error, got %v", err)
		}
	})

	t.Run("forwards uint32 window ID", func(t *testing.T) {
		t.Setenv("AEROSPACE_WINDOW_ID", "4294967295")
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		responsePayload, err := json.Marshal(Response{ExitCode: 0})
		if err != nil {
			t.Fatalf("failed to marshal response: %v", err)
		}
		expectedRequest := []byte(`{"args":["list-workspaces","--focused"],"stdin":"","windowId":4294967295,"workspace":null}`)
		gomock.InOrder(
			mockConn.EXPECT().Write(framedResponse(expectedRequest)).Return(len(framedResponse(expectedRequest)), nil),
			setupResponseReads(mockConn, responsePayload),
		)

		connection := &AeroSpaceSocketConnection{Conn: mockConn}
		if _, err := connection.SendCommand("list-workspaces", []string{"--focused"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("rejects window ID larger than uint32", func(t *testing.T) {
		t.Setenv("AEROSPACE_WINDOW_ID", "4294967296")
		ctrl := gomock.NewController(t)
		connection := &AeroSpaceSocketConnection{Conn: net_mock.NewMockConn(ctrl)}

		_, err := connection.SendCommand("list-workspaces", []string{"--focused"})
		if err == nil || !containsSubstring(err.Error(), "failed to parse AEROSPACE_WINDOW_ID") {
			t.Fatalf("expected invalid window ID error, got %v", err)
		}
	})
}

func TestCloseAfterConnectionError(t *testing.T) {
	primaryErr := errors.New("handshake failed")

	t.Run("preserves connection error when close succeeds", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		mockConn.EXPECT().Close().Return(nil)

		err := closeAfterConnectionError(mockConn, primaryErr)
		if !errors.Is(err, primaryErr) {
			t.Fatalf("expected primary error, got %v", err)
		}
	})

	t.Run("preserves connection and close errors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockConn := net_mock.NewMockConn(ctrl)
		closeErr := errors.New("close failed")
		mockConn.EXPECT().Close().Return(closeErr)

		err := closeAfterConnectionError(mockConn, primaryErr)
		if !errors.Is(err, primaryErr) {
			t.Fatalf("expected primary error, got %v", err)
		}
		if !errors.Is(err, closeErr) {
			t.Fatalf("expected close error, got %v", err)
		}
	})
}

func TestSocketClient(t *testing.T) {
	testCases := []struct {
		title           string
		minMajorVersion int
		minMinorVersion int
		serverVersion   string
		expectation     func(*testing.T, error)
	}{
		{
			title:           "CheckServerVersion - fails when major different than minimum version",
			minMajorVersion: 2,
			minMinorVersion: 10,
			serverVersion:   "1.10.0-beta xxxxx",
			expectation: func(t *testing.T, err error) {
				if err == nil {
					t.Fatalf("expected error about minimum version, got nil")
				}
				if !errors.Is(err, exceptions.ErrVersion) {
					t.Fatalf("expected error about minimum version, got %v", err)
				}
			},
		},
		{
			title:           "CheckServerVersion - fails when minor isn't greater than to minimum minor version",
			minMajorVersion: 1,
			minMinorVersion: 12,
			serverVersion:   "1.10.0-beta xxxxx",
			expectation: func(t *testing.T, err error) {
				if err == nil {
					t.Fatalf("expected error about minimum version, got nil")
				}
				if !errors.Is(err, exceptions.ErrVersion) {
					t.Fatalf("expected error about minimum version, got %v", err)
				}
			},
		},
		{
			title:           "CheckServerVersion - succeeds with same major and minor version",
			minMajorVersion: 1,
			minMinorVersion: 10,
			serverVersion:   "1.10.0-beta xxxxx",
			expectation: func(t *testing.T, err error) {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.title, func(tt *testing.T) {
			mockedResponse := Response{
				ServerVersion: tc.serverVersion,
				ExitCode:      0,
			}
			cmdBytes, err := json.Marshal(mockedResponse)
			if err != nil {
				t.Fatalf("failed to marshal mocked response: %v", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := net_mock.NewMockConn(ctrl)
			setupSendCommandMock(mockConn, cmdBytes)

			connection := &AeroSpaceSocketConnection{
				MinMajorVersion: tc.minMajorVersion,
				MinMinorVersion: tc.minMinorVersion,
				Conn:            mockConn,
				socketPath:      "/tmp/aerospace.sock",
			}
			err = connection.CheckServerVersion()
			tc.expectation(tt, err)
		})
	}

	t.Run("GetSocketPath - retrieves the socket path", func(tt *testing.T) {
		expectedSocketPath := "/tmp/aerospace.sock"
		connection := &AeroSpaceSocketConnection{
			MinMajorVersion: 2,
			MinMinorVersion: 10,
			Conn:            nil,
			socketPath:      expectedSocketPath,
		}

		socketPath, err := connection.GetSocketPath()
		if err != nil {
			tt.Fatalf("expected no error, got %v", err)
		}

		if socketPath != expectedSocketPath {
			tt.Fatalf("expected socket path %s, got %s", expectedSocketPath, socketPath)
		}
	})
}

// TestCheckServerVersion provides comprehensive unit tests for CheckServerVersion function
func TestCheckServerVersion(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		successCases := []struct {
			name            string
			minMajorVersion int
			minMinorVersion int
			serverVersion   string
		}{
			{
				name:            "same major and minor version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.0-beta abc123",
			},
			{
				name:            "version without hash suffix",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.0",
			},
			{
				name:            "version with patch number",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.5-beta abc123",
			},
		}

		for _, tc := range successCases {
			t.Run(tc.name, func(tt *testing.T) {
				cmdBytes, err := json.Marshal(Response{ServerVersion: tc.serverVersion, ExitCode: 0})
				if err != nil {
					tt.Fatalf("failed to marshal mocked response: %v", err)
				}

				ctrl := gomock.NewController(tt)
				defer ctrl.Finish()

				mockConn := net_mock.NewMockConn(ctrl)
				setupSendCommandMock(mockConn, cmdBytes)

				connection := &AeroSpaceSocketConnection{
					MinMajorVersion: tc.minMajorVersion,
					MinMinorVersion: tc.minMinorVersion,
					Conn:            mockConn,
					socketPath:      "/tmp/aerospace.sock",
				}

				if err := connection.CheckServerVersion(); err != nil {
					tt.Fatalf("expected no error, got %v", err)
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		errorCases := []struct {
			name             string
			minMajorVersion  int
			minMinorVersion  int
			setupMock        func(*gomock.Controller, *net_mock.MockConn)
			expectedErrorMsg string
		}{
			{
				name:             "GetServerVersion fails - connection not established",
				minMajorVersion:  0,
				minMinorVersion:  20,
				setupMock:        func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {},
				expectedErrorMsg: "connection is not established",
			},
			{
				name:            "GetServerVersion fails - SendCommand error",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockConn.EXPECT().
						Write(gomock.Any()).
						Return(0, io.ErrUnexpectedEOF)
				},
				expectedErrorMsg: "failed to send command",
			},
			{
				name:            "empty server version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "server version is empty",
			},
			{
				name:            "invalid version format - only major version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "invalid server version format",
			},
			{
				name:            "non-numeric major version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "abc.20.0-beta xyz", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "failed to parse major version",
			},
			{
				name:            "non-numeric minor version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.abc.0-beta xyz", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "failed to parse minor version",
			},
			{
				name:            "version mismatch - different major version",
				minMajorVersion: 1,
				minMinorVersion: 0,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.20.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - higher major version (not allowed)",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "1.0.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - same major, lower minor",
				minMajorVersion: 0,
				minMinorVersion: 25,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.20.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - same major, higher minor (not allowed)",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.25.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - exact major match with higher minor",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.30.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.19.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.19.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.18.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.18.5-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.15.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.15.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.1.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.1.0-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.19.9 below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					cmdBytes, _ := json.Marshal(Response{ServerVersion: "0.19.9-beta abc123", ExitCode: 0})
					setupSendCommandMock(mockConn, cmdBytes)
				},
				expectedErrorMsg: "version mismatch",
			},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(tt *testing.T) {
				ctrl := gomock.NewController(tt)
				defer ctrl.Finish()

				mockConn := net_mock.NewMockConn(ctrl)
				tc.setupMock(ctrl, mockConn)

				var conn net.Conn = mockConn
				if tc.name == "GetServerVersion fails - connection not established" {
					conn = nil
				}

				connection := &AeroSpaceSocketConnection{
					MinMajorVersion: tc.minMajorVersion,
					MinMinorVersion: tc.minMinorVersion,
					Conn:            conn,
					socketPath:      "/tmp/aerospace.sock",
				}

				err := connection.CheckServerVersion()
				if err == nil {
					tt.Fatalf("expected error containing '%s', got nil", tc.expectedErrorMsg)
				}

				if tc.expectedErrorMsg == "version mismatch" {
					if !errors.Is(err, exceptions.ErrVersion) {
						tt.Fatalf("expected ErrVersion error, got %v", err)
					}
				} else {
					if !containsSubstring(err.Error(), tc.expectedErrorMsg) {
						tt.Fatalf("expected error message containing '%s', got '%s'", tc.expectedErrorMsg, err.Error())
					}
				}
			})
		}
	})
}
