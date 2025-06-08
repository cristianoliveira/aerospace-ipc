package client

import (
	"encoding/json"
	"errors"
	"io"
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	"github.com/cristianoliveira/aerospace-ipc/internal/mocks/net"
	"go.uber.org/mock/gomock"
)

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
			title:           "CheckServerVersion - fails when minor different than minimum version",
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
			minMajorVersion := tc.minMajorVersion
			minMinorVersion := tc.minMinorVersion
			serverVersion := tc.serverVersion

			mockedResponse := Response{
				ServerVersion: serverVersion,
				StdErr:        "",
				StdOut:        "",
				ExitCode:      0,
			}
			cmdBytes, err := json.Marshal(mockedResponse)
			if err != nil {
				t.Fatalf("failed to marshal mocked response: %v", err)
			}

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockConn := net_mock.NewMockConn(ctrl)

			readCount := 0
			gomock.InOrder(
				mockConn.EXPECT().
					Write(gomock.Any()).
					Return(0, nil),

				mockConn.EXPECT().
					Read(gomock.Any()).
					DoAndReturn(
						func(p []byte) (int, error) {
							n := copy(p, cmdBytes)
							readCount++
							if readCount == 1 {
								return n, nil
							}
							return n, io.EOF
						},
					).
					Times(1),
			)

			connection := &AeroSpaceSocketConnection{
				MinMajorVersion: minMajorVersion,
				MinMinorVersion: minMinorVersion,
				Conn:            mockConn, // Not used in this test
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
			Conn:            nil, // Not used in this test
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
