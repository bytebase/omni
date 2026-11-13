package testcontainer

import (
	"errors"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"

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

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*mysql.MySQLError{
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: syntaxErrorPreamble + "'x'\nat line 2"},
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: "a syntax error worded otherwise"},
		{Number: 1105, Message: "no SQLSTATE"},
	} {
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("%v is not encoded", in)
		}
		if strings.Contains(text, syntaxErrorPreamble) {
			t.Errorf("%q holds the preamble every syntax error opens with", text)
		}
		var out *mysql.MySQLError
		if err := (Codec{}).DecodeError(text); !errors.As(err, &out) || *out != *in {
			t.Errorf("%q decodes to %#v, want %#v", text, err, in)
		}
	}
	if _, ok := (Codec{}).EncodeError(mysql.ErrInvalidConn); ok {
		t.Error("a dropped connection was encoded as the engine's answer")
	}
}
