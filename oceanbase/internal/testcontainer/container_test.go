package testcontainer

import (
	"context"
	"errors"
	"os"
	"slices"
	"strings"
	"testing"
	"time"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/enginetest/tape"
)

func TestMain(m *testing.M) {
	code := m.Run()
	OceanBase.Terminate()
	os.Exit(code)
}

// TestRetryOceanBase pins what the retrying connector promises a verdict helper:
// a transient answer is asked again, a real answer comes back as the server's
// error, and a transient answer that never clears is no answer at all.
func TestRetryOceanBase(t *testing.T) {
	window, pause := retryWindow, retryPause
	t.Cleanup(func() { retryWindow, retryPause = window, pause })
	retryWindow, retryPause = 50*time.Millisecond, time.Millisecond
	ctx := context.Background()
	shutdown := &gomysql.MySQLError{Number: 4550, Message: "rootservice is shutdown"}

	calls := 0
	if err := retry(ctx, func() error {
		calls++
		if calls < 3 {
			return shutdown
		}
		return nil
	}); err != nil || calls != 3 {
		t.Errorf("a transient answer that clears: err=%v after %d calls, want nil after 3", err, calls)
	}

	calls = 0
	err := retry(ctx, func() error {
		calls++
		return &gomysql.MySQLError{Number: 1064, Message: "syntax"}
	})
	var me *gomysql.MySQLError
	if !errors.As(err, &me) || me.Number != 1064 || calls != 1 {
		t.Errorf("a parse error: err=%v after %d calls, want 1064 after 1", err, calls)
	}

	err = retry(ctx, func() error { return shutdown })
	if err == nil || errors.As(err, &me) {
		t.Errorf("a transient answer that never clears came back as %#v, want a plain error", err)
	}
}

func TestCodecRoundTrips(t *testing.T) {
	for _, in := range []*gomysql.MySQLError{
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: "You have an error in your SQL syntax; near 'x'\nat line 2"},
		{Number: 4013, Message: "no SQLSTATE"},
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: syntaxErrorPreamble + "'()' at line 1"},
		{Number: 1064, SQLState: [5]byte{'4', '2', '0', '0', '0'}, Message: "Parse error"},
	} {
		text, ok := Codec{}.EncodeError(in)
		if !ok {
			t.Fatalf("%v is not encoded", in)
		}
		if len(text) > 100 && strings.Contains(in.Message, "'()'") {
			t.Errorf("the parser's syntax error is written in full: %q", text)
		}
		var out *gomysql.MySQLError
		if err := (Codec{}).DecodeError(text); !errors.As(err, &out) || *out != *in {
			t.Errorf("%q decodes to %#v, want %#v", text, err, in)
		}
	}
	// A transient answer that outlasts the retry window is a plain error, and
	// must not reach a tape as the engine's answer.
	if _, ok := (Codec{}).EncodeError(errors.New("OceanBase stayed unavailable")); ok {
		t.Error("an error the server did not send was encoded as its answer")
	}
}

// TestCodecStandIns pins the shape of the names OceanBase makes up and of what
// a tape holds in their place.
func TestCodecStandIns(t *testing.T) {
	var _ tape.Renamer = Codec{}
	re := Codec{}.MadeUp()
	const stored = "CONSTRAINT `c_OBFK_1794529609863994` FOREIGN KEY (`a`), CONSTRAINT `t_OBCHECK_1794529978914345` CHECK ((`a` > 0)), CONSTRAINT `fk_1794529609863994` FOREIGN KEY (`b`)"
	got := re.FindAllString(stored, -1)
	if want := []string{"_OBFK_1794529609863994", "_OBCHECK_1794529978914345"}; !slices.Equal(got, want) {
		t.Fatalf("made-up names in %q: %q, want %q", stored, got, want)
	}
	for i, name := range got {
		in := Codec{}.StandIn(name, i+1)
		if in == name || re.FindString(in) != in {
			t.Errorf("the stand-in for %q is %q, which must differ from it and still read as a made-up name", name, in)
		}
	}
	if in := (Codec{}).StandIn("_OBCHECK_1794529978914345", 12); in != "_OBCHECK_0000000000000012" {
		t.Errorf("stand-in 12 is %q", in)
	}
}
