package testcontainer

import (
	"errors"
	"testing"

	"github.com/go-sql-driver/mysql"
)

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*mysql.MySQLError{
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: "You have an error in your SQL syntax; near 'x'\nat line 2"},
		{Number: 1105, Message: "no SQLSTATE"},
	} {
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("%v is not encoded", in)
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
