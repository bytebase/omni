package compat

import (
	"errors"
	"testing"
)

// serverError is a stand-in for a driver's server error.
type serverError uint16

func (e serverError) Error() string { return "server error" }

func code(err error) (uint16, bool) {
	var s serverError
	if !errors.As(err, &s) {
		return 0, false
	}
	return uint16(s), true
}

// TestPrepareVerdictFailsWithoutVerdict pins that a server that never
// answered is no verdict, never silent agreement.
func TestPrepareVerdictFailsWithoutVerdict(t *testing.T) {
	unprobed := func() error {
		t.Errorf("the verdict asked the probe")
		return nil
	}
	if _, err := PrepareVerdict(errors.New("connection lost"), code, unprobed); err == nil {
		t.Errorf("a dropped connection came back as a verdict")
	}
	if _, err := PrepareVerdict(serverError(1146), code, func() error { return errors.New("connection lost") }); err == nil {
		t.Errorf("a probe on a dropped connection came back as a verdict")
	}
	if _, err := PrepareVerdict(serverError(1146), code, func() error { return nil }); err == nil {
		t.Errorf("a probe the server took came back as a verdict")
	}
}

// TestPrepareVerdictReadsTheProbe pins how a verdict reads the code: a
// grammar action's code comes back with the tail after the statement, and
// anything raised after the parse gives way to the tail's 1064.
func TestPrepareVerdictReadsTheProbe(t *testing.T) {
	for _, c := range []struct {
		name  string
		err   error
		probe uint16 // 0: the verdict must not ask
		want  bool
	}{
		{"prepared", nil, 0, true},
		{"parse error", serverError(1064), 0, false},
		{"a name it looks up", serverError(1193), 0, true},
		{"table missing", serverError(1146), 1064, true},
		{"grammar action", serverError(1149), 1149, false},
		{"a code raised both ways, after the parse", serverError(1105), 1064, true},
		{"a code raised both ways, while parsing", serverError(1105), 1105, false},
	} {
		got, err := PrepareVerdict(c.err, code, func() error {
			if c.probe == 0 {
				t.Errorf("%s: the verdict asked the probe", c.name)
				return nil
			}
			return serverError(c.probe)
		})
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		} else if got != c.want {
			t.Errorf("%s: accepts=%v, want %v", c.name, got, c.want)
		}
	}
}
