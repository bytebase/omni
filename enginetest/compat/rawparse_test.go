package compat

import (
	"errors"
	"testing"
)

// pgError is a stand-in for a driver's server error: a SQLSTATE and the
// position the server reported, 0 for none.
type pgError struct {
	code     string
	position int
}

func (e pgError) Error() string { return "server error " + e.code }

func pgState(err error) (string, int, bool) {
	var e pgError
	if !errors.As(err, &e) {
		return "", 0, false
	}
	return e.code, e.position, true
}

// at is the position the server reports for an error at character chars of
// the statement: characters from 1 into the head and the statement.
func at(chars int) int { return len(RawParseHead) + chars + 1 }

// TestRawParseVerdictReadsTheProbe pins how a verdict reads the answer: the
// head's own failure is an acceptance, and every other server error is the
// raw parser's, at the byte the server named. The answers are those of
// postgres:17.12-alpine.
func TestRawParseVerdictReadsTheProbe(t *testing.T) {
	for _, c := range []struct {
		name      string
		statement string
		err       error
		want      string
	}{
		{"parsed", "SELECT 1", pgError{"22012", 0}, Accept},
		{"a syntax error", "SELECT FROM FROM", pgError{"42601", at(12)}, "reject@12"},
		{"at the first byte", "FROM", pgError{"42601", at(0)}, "reject@0"},
		{"at the end of input", "SELECT (", pgError{"42601", at(8)}, "reject@8"},
		{"characters, not bytes", "SELECT 'é', 'ü' FROM FROM", pgError{"42601", at(21)}, "reject@23"},
		{"the lexer's error", `SELECT U&'\00'`, pgError{"42601", at(7)}, "reject@7"},
		{"a grammar action's error", "SELECT ROW(1) OVERLAPS ROW(2)", pgError{"42601", at(7)}, "reject@7"},
		{"a location left at zero", "ALTER PUBLICATION a ADD a *", pgError{"42601", 1}, "reject@0"},
		{"another code, and no position", "CREATE TABLE a (b int, CONSTRAINT c CHECK (b>0) DEFERRABLE)", pgError{"0A000", 0}, Reject},
	} {
		got, err := RawParseVerdict(c.statement, c.err, pgState)
		if err != nil {
			t.Errorf("%s: %v", c.name, err)
		} else if got != c.want {
			t.Errorf("%s: %s, want %s", c.name, got, c.want)
		}
		if err := validVerdict(got); err != nil {
			t.Errorf("%s: %v", c.name, err)
		}
	}
}

// TestRawParseVerdictFailsWithoutVerdict pins that an answer the probe
// cannot read is no verdict, never silent agreement.
func TestRawParseVerdictFailsWithoutVerdict(t *testing.T) {
	for _, c := range []struct {
		name string
		err  error
	}{
		{"the head ran clean", nil},
		{"a dropped connection", errors.New("connection lost")},
		{"an error inside the head", pgError{"42601", 3}},
		{"an error past the statement", pgError{"42601", at(len("SELECT 1") + 1)}},
	} {
		if v, err := RawParseVerdict("SELECT 1", c.err, pgState); err == nil {
			t.Errorf("%s: came back as the verdict %s", c.name, v)
		}
	}
}
