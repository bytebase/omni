package testcontainer

import (
	"testing"

	"github.com/bytebase/omni/enginetest"
)

// TestImageIsThePin holds Image, which the packages not yet on tapes start,
// to the server pins.tsv names.
func TestImageIsThePin(t *testing.T) {
	pin, err := enginetest.LoadPin(Pins, "default")
	if err != nil {
		t.Fatal(err)
	}
	if pin.Image != Image {
		t.Errorf("pins.tsv pins %s, and Image is %s", pin.Image, Image)
	}
}
