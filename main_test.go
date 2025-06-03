package aerospace

import "testing"

func TestIntegration(t *testing.T) {
	t.Run("TestGetAllWindows", func(t *testing.T) {
		client, err := NewAeroSpaceConnection()
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
}
