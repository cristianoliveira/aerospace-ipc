package client

import (
	"errors"
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/internal/exceptions"
)

func TestSocketClient(t *testing.T) {
	t.Run("CheckServerVersion - checks major and minor versions", func(tt *testing.T) {
		connection := &AeroSpaceSocketConnection{
			MinMajorVersion: 2,
			MinMinorVersion: 10,
			Conn:            nil, // Not used in this test
			socketPath:      "/tmp/aerospace.sock",
		}

		err := connection.CheckServerVersion("3.10.0-beta xxxxx")
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !errors.Is(err, exceptions.ErrVersion) {
			t.Fatalf("expected error about minimum version, got %v", err)
		}
	})

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
