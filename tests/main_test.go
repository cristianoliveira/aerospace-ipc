//go:build integration

package tests

import (
	"testing"

	"github.com/cristianoliveira/aerospace-ipc/pkg/aerospace"
)

func TestIntegration(t *testing.T) {
	t.Run("aerospace.Windows().GetAllWindows: retrieves all windows", func(t *testing.T) {
		client, err := aerospace.NewClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer client.CloseConnection()

		windows, err := client.Windows().GetAllWindows()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(windows) == 0 {
			t.Fatal("expected at least one window, got none")
		}

		t.Logf("Retrieved %d windows", len(windows))
	})

	t.Run("aerospace.Client().Connection().GetSocketPath: retrieves the socket path", func(t *testing.T) {
		client, err := aerospace.NewClient()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		defer client.CloseConnection()

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
