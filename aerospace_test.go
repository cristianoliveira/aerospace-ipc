package aerospace

import "testing"

func TestAeroSpaceWM(t *testing.T) {
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
