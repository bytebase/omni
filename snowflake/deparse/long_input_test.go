package deparse_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/enginetest/cputime"
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
		cputime.Within(t, 2*time.Second, fmt.Sprintf("%s: deparse of %d bytes", tc.name, len(tc.sql)), func() {
			_, err = deparse.DeparseFile(file)
		})
		if err != nil {
			t.Fatalf("%s: deparse error = %v", tc.name, err)
		}
	}
}
