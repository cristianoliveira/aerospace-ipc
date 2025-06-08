package aerospace

import "testing"

func TestAeroSpaceWM(t *testing.T) {
	// This is a bit annoying but I had a case where I pushed a change that
	// didn't implement the AeroSpaceClient interface, and I only found out
	// attempting use it later.
	t.Run("Implements the AeroSpaceClient interface", func(t *testing.T) {
		var client any
		aeroSpaceWM, _ := NewAeroSpaceClient()
		client = aeroSpaceWM

		if _, ok := client.(AeroSpaceClient); !ok {
			t.Fatal("AeroSpaceWM does not implement AeroSpaceClient interface")
		}
		t.Log("AeroSpaceWM implements AeroSpaceClient interface")
	})
}
