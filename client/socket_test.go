package client

import (
	"strings"
	"testing"
)

func TestSocketClient(t *testing.T) {
	t.Run("CheckServerVersion - checks major and minor versions", func(tt *testing.T) {
		connection := &AeroSpaceSocketConnection{
			MinMajorVersion: 1,
			MinMinorVersion: 0,
			Conn:            nil, // Not used in this test
			SocketPath:      "/tmp/aerospace.sock",
		}

		err := connection.CheckServerVersion("2.0.0-beta xxxxx")
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(err.Error(), "server major version 2.0.0") ||
			!strings.Contains(err.Error(), "minimum required 1.0.x") {
			t.Fatalf("expected error about minimum version, got %v", err)
		}

		err = connection.CheckServerVersion("1.2.0-beta xxxxx")
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(err.Error(), "server minor version 1.2.0") ||
			!strings.Contains(err.Error(), "minimum required 1.0.x") {
			t.Fatalf("expected error about minimum version, got %v", err)
		}
	})
}
