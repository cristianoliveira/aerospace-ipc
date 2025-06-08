//go:build integration

package tests

import (
	"testing"

	ipc "github.com/cristianoliveira/aerospace-ipc"
)

func TestIntegration(t *testing.T) {
	t.Run("ipc.GetAllWindows: retrieves all windows", func(t *testing.T) {
		client, err := ipc.NewAeroSpaceClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		windows, err := client.GetAllWindows()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(windows) == 0 {
			t.Fatal("expected at least one window, got none")
		}

		t.Logf("Retrieved %d windows", len(windows))
	})

	t.Run("ipc.Client().GetSocketPath: retrieves the socket path", func(t *testing.T) {
		client, err := ipc.NewAeroSpaceClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		path, err := client.Connection().GetSocketPath()
		if err != nil {
			t.Fatalf("expected no error, got '%v'", err)
		}
		if len(path) == 0 {
			t.Fatal("expected non-empty socket path, got empty")
		}

		t.Logf("Retrieved %s path", path)
	})
}
