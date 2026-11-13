package compat

import (
	"errors"
	"strings"
	"testing"
)

// trinoError is a stand-in for the Trino client's error: what the server
// said about a failed statement.
type trinoError struct{ TrinoError }

func (e trinoError) Error() string { return "trino: " + e.Name }

func readTrino(err error) (TrinoError, bool) {
	var e trinoError
	if !errors.As(err, &e) {
		return TrinoError{}, false
	}
	return e.TrinoError, true
}

func parsing(line, column int) error {
	return trinoError{TrinoError{Name: "SYNTAX_ERROR", Type: TrinoParsingException, Line: line, Column: column}}
}

// TestTrinoVerdictReadsTheParsersException pins how a verdict reads the
// answer: the parser's exception is a rejection at the byte of its
// location, and anything else, a SYNTAX_ERROR raised after the parse
// included, is an acceptance. The answers are those of trinodb/trino:482.
func TestTrinoVerdictReadsTheParsersException(t *testing.T) {
	for _, c := range []struct {
		name      string
		statement string
		err       error
		want      string
	}{
		{"ran", "CREATE TABLE a (a TIMESTAMP)", nil, Accept},
		{"a syntax error", "SELECT 1 FROM FROM", parsing(1, 15), "reject@14"},
		{"at the end of input", "SELECT (1", parsing(1, 10), "reject@9"},
		{"on a later line", "SELECT 1\n  FROM FROM", parsing(2, 8), "reject@16"},
		{"code points, not bytes", "SELECT 'é' FROM FROM", parsing(1, 17), "reject@17"},
		{"a character outside the BMP is one", "SELECT '😀' FROM FROM", parsing(1, 17), "reject@19"},
		{"AstBuilder's refusal", "SELECT 1a", parsing(1, 8), "reject@7"},
		{"no location", "SELECT", parsing(0, 0), Reject},
		{"SYNTAX_ERROR after the parse", "ALTER TABLE a.b.c.d RENAME TO e",
			trinoError{TrinoError{Name: "SYNTAX_ERROR", Type: "io.trino.spi.TrinoException"}}, Accept},
		{"a missing table", "SELECT * FROM a",
			trinoError{TrinoError{Name: "TABLE_NOT_FOUND", Type: "io.trino.spi.TrinoException", Line: 1, Column: 15}}, Accept},
	} {
		got, err := TrinoVerdict(c.statement, c.err, readTrino)
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

// TestTrinoVerdictRefusesWhatIsNoAnswer pins the answers that are no
// verdict, which fail a recording rather than read as one.
func TestTrinoVerdictRefusesWhatIsNoAnswer(t *testing.T) {
	for _, c := range []struct {
		name      string
		statement string
		err       error
		want      string
	}{
		{"not the server's", "SELECT 1", errors.New("connection refused"), "connection refused"},
		{"no exception class", "SELECT 1", trinoError{TrinoError{Name: "SYNTAX_ERROR"}}, "no exception class"},
		{"the parser under another name", "SELECT 1",
			trinoError{TrinoError{Name: "GENERIC_USER_ERROR", Type: TrinoParsingException, Line: 1, Column: 1}}, "not SYNTAX_ERROR"},
		{"past the end", "SELECT (1", parsing(1, 11), "outside the statement"},
		{"past the line", "SELECT 1\nFROM", parsing(1, 10), "outside the statement"},
		{"no such line", "SELECT 1", parsing(2, 1), "outside the statement"},
	} {
		got, err := TrinoVerdict(c.statement, c.err, readTrino)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: %q, %v; want an error holding %q", c.name, got, err, c.want)
		}
	}
}
