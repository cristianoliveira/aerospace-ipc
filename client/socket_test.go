package client

import (
	"strings"
	"testing"
)

func TestSocketClient(t *testing.T) {
	t.Run("CheckServerVersion - checks major and minor versions", func(tt *testing.T) {
		connection := &AeroSpaceSocketConnection{
			MinMajorVersion: 2,
			MinMinorVersion: 10,
			Conn:            nil, // Not used in this test
			SocketPath:      "/tmp/aerospace.sock",
		}

		err := connection.CheckServerVersion("3.10.0-beta xxxxx")
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if !strings.Contains(err.Error(), "[VERSION-MISMATCH]") ||
			!strings.Contains(err.Error(), "server major version 3.10.0") ||
			!strings.Contains(err.Error(), "minimum required 2.10.x") {
			t.Fatalf("expected error about minimum version, got %v", err)
		}

		err = connection.CheckServerVersion("1.2.0-beta xxxxx")
		if err == nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !strings.Contains(err.Error(), "[VERSION-MISMATCH]") ||
			!strings.Contains(err.Error(), "server major version 1.2.0") ||
			!strings.Contains(err.Error(), "minimum required 2.10.x") {
			t.Fatalf("major ok but min no-ok, got %v", err)
		}
	})
}
