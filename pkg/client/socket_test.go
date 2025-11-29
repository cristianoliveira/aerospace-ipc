package client

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
	"github.com/cristianoliveira/aerospace-ipc/internal/mocks/net"
	"go.uber.org/mock/gomock"
)

// containsSubstring checks if a string contains a substring (case-insensitive)
func containsSubstring(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
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

// TestCheckServerVersion provides comprehensive unit tests for CheckServerVersion function
func TestCheckServerVersion(t *testing.T) {
	t.Run("success cases", func(t *testing.T) {
		successCases := []struct {
			name             string
			minMajorVersion  int
			minMinorVersion  int
			serverVersion    string
			expectNoError    bool
		}{
			{
				name:            "same major and minor version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.0-beta abc123",
				expectNoError:   true,
			},
			{
				name:            "same major, higher minor version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.25.0-beta abc123",
				expectNoError:   true,
			},
			{
				name:            "exact major match with higher minor",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.30.0-beta abc123",
				expectNoError:   true,
			},
			{
				name:            "version without hash suffix",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.0",
				expectNoError:   true,
			},
			{
				name:            "version with patch number",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.20.5-beta abc123",
				expectNoError:   true,
			},
		}

		for _, tc := range successCases {
			t.Run(tc.name, func(tt *testing.T) {
				mockedResponse := Response{
					ServerVersion: tc.serverVersion,
					StdErr:        "",
					StdOut:        "",
					ExitCode:      0,
				}
				cmdBytes, err := json.Marshal(mockedResponse)
				if err != nil {
					tt.Fatalf("failed to marshal mocked response: %v", err)
				}

				ctrl := gomock.NewController(tt)
				defer ctrl.Finish()

				mockConn := net_mock.NewMockConn(ctrl)

				readCount := 0
				gomock.InOrder(
					mockConn.EXPECT().
						Write(gomock.Any()).
						Return(len(cmdBytes), nil),

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
					MinMajorVersion: tc.minMajorVersion,
					MinMinorVersion: tc.minMinorVersion,
					Conn:            mockConn,
					socketPath:      "/tmp/aerospace.sock",
				}

				err = connection.CheckServerVersion()
				if tc.expectNoError && err != nil {
					tt.Fatalf("expected no error, got %v", err)
				}
				if !tc.expectNoError && err == nil {
					tt.Fatalf("expected error, got nil")
				}
			})
		}
	})

	t.Run("error cases", func(t *testing.T) {
		errorCases := []struct {
			name             string
			minMajorVersion  int
			minMinorVersion  int
			serverVersion    string
			setupMock        func(*gomock.Controller, *net_mock.MockConn)
			expectedErrorMsg string
		}{
			{
				name:            "GetServerVersion fails - connection not established",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					// No expectations - connection will be nil
				},
				expectedErrorMsg: "connection is not established",
			},
			{
				name:            "GetServerVersion fails - SendCommand error",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "",
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
				serverVersion:   "",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "server version is empty",
			},
			{
				name:            "invalid version format - only major version (causes panic)",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "panic", // The function panics when versionParts[1] is accessed
			},
			{
				name:            "non-numeric major version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "abc.20.0-beta xyz",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "abc.20.0-beta xyz",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "failed to parse major version",
			},
			{
				name:            "non-numeric minor version",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.abc.0-beta xyz",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.abc.0-beta xyz",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "failed to parse minor version",
			},
			{
				name:            "version mismatch - different major version",
				minMajorVersion: 1,
				minMinorVersion: 0,
				serverVersion:   "0.20.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.20.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - higher major version (not allowed)",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "1.0.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "1.0.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - same major, lower minor",
				minMajorVersion: 0,
				minMinorVersion: 25,
				serverVersion:   "0.20.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.20.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.19.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.19.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.19.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.18.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.18.5-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.18.5-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.15.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.15.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.15.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.1.x below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.1.0-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.1.0-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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
				},
				expectedErrorMsg: "version mismatch",
			},
			{
				name:            "version mismatch - 0.19.9 below minimum 0.20",
				minMajorVersion: 0,
				minMinorVersion: 20,
				serverVersion:   "0.19.9-beta abc123",
				setupMock: func(ctrl *gomock.Controller, mockConn *net_mock.MockConn) {
					mockedResponse := Response{
						ServerVersion: "0.19.9-beta abc123",
						StdErr:        "",
						StdOut:        "",
						ExitCode:      0,
					}
					cmdBytes, _ := json.Marshal(mockedResponse)

					readCount := 0
					gomock.InOrder(
						mockConn.EXPECT().
							Write(gomock.Any()).
							Return(len(cmdBytes), nil),

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

				// For "connection not established" test, set Conn to nil
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

				// Handle panic case for invalid version format
				if tc.expectedErrorMsg == "panic" {
					defer func() {
						if r := recover(); r == nil {
							tt.Fatalf("expected panic for invalid version format, but no panic occurred")
						}
					}()
					_ = connection.CheckServerVersion()
					tt.Fatalf("expected panic, but function completed")
					return
				}

				err := connection.CheckServerVersion()
				if err == nil {
					tt.Fatalf("expected error containing '%s', got nil", tc.expectedErrorMsg)
				}

				// Check for specific error types
				if tc.expectedErrorMsg == "version mismatch" {
					if !errors.Is(err, exceptions.ErrVersion) {
						tt.Fatalf("expected ErrVersion error, got %v", err)
					}
				} else {
					// Check that error message contains expected text
					errMsg := err.Error()
					if errMsg == "" {
						tt.Fatalf("expected error message containing '%s', got empty error", tc.expectedErrorMsg)
					}
					// Verify error message contains the expected substring
					if !containsSubstring(errMsg, tc.expectedErrorMsg) {
						tt.Fatalf("expected error message containing '%s', got '%s'", tc.expectedErrorMsg, errMsg)
					}
				}
			})
		}
	})
}
