package deparse_test

import (
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/snowflake/deparse"
	"github.com/bytebase/omni/snowflake/parser"
)

// TestDeparseLongInputsStayLinear: a subscript, the first element of an
// array literal and many other expressions are written with a leading space
// that is then taken back. Taking it back copied the whole output, once per
// such expression, so 100,000 subscripted arguments took 16s and the others
// below took 14s to 36s; they take milliseconds now.
func TestDeparseLongInputsStayLinear(t *testing.T) {
	for _, tc := range []struct{ name, sql string }{
		{"subscripts", "SELECT a" + strings.Repeat("[1]", 90000) + " FROM t"},
		{"arguments", "SELECT f(" + strings.Repeat("a[1], ", 100000) + "1) FROM t"},
		{"array elements", "SELECT [" + strings.Repeat("[1], ", 100000) + "1] FROM t"},
		{"view", "CREATE VIEW v AS SELECT a[1]" + strings.Repeat(" || a[1]", 90000) + " AS c FROM t"},
	} {
		file, err := parser.Parse(tc.sql)
		if err != nil {
			t.Fatalf("%s: parse error = %v", tc.name, err)
		}
		done := make(chan error, 1)
		go func() {
			_, err := deparse.DeparseFile(file)
			done <- err
		}()
		select {
		case err := <-done:
			if err != nil {
				t.Fatalf("%s: deparse error = %v", tc.name, err)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("%s: deparse of %d bytes took over 2s", tc.name, len(tc.sql))
		}
	}
}
