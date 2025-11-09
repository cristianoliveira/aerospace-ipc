package aerospace

import "testing"

func TestAeroSpaceWM(t *testing.T) {
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
