package aerospace

import "testing"

func TestAeroSpaceWM(t *testing.T) {
	t.Run("Monitors returns a cached monitor service", func(t *testing.T) {
		client := &AeroSpaceWM{}

		service := client.Monitors()
		if service == nil {
			t.Fatal("expected monitor service")
		}
		if got := client.Monitors(); got != service {
			t.Fatal("expected the same monitor service instance")
		}
	})

	t.Run("Implements the Client interface", func(t *testing.T) {
		var client any
		aeroSpaceWM, _ := NewClient()
		client = aeroSpaceWM

		if _, ok := client.(Client); !ok {
			t.Fatal("AeroSpaceWM does not implement Client interface")
		}
		t.Log("AeroSpaceWM implements Client interface")
	})
}
