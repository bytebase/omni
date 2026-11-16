package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/starrocks/parser"
)

// A statement cut out of a script reports its position in the script, not in
// the segment: a caller highlighting the error must land on the right text.
func TestValidateQueryReportsScriptPositions(t *testing.T) {
	const script = "SELECT 1;\nSELECT FROM"
	_, _, err := ValidateQuery(script)
	if err == nil {
		t.Fatal("ValidateQuery accepted a syntax error")
	}
	errs := parser.AllErrors(err)
	if len(errs) == 0 {
		t.Fatalf("ValidateQuery error = %v, want a ParseError", err)
	}
	if want := strings.Index(script, "FROM"); errs[0].Position != want {
		t.Errorf("error position = %d, want %d", errs[0].Position, want)
	}
}

// A statement that changes state outside the session is not read-only, and the
// session forms stay allowed.
func TestValidateQueryServerState(t *testing.T) {
	tests := []struct {
		statement   string
		readOnly    bool
		returnsData bool
	}{
		{"SET query_timeout = 3600", true, false},
		{"SET @@session.query_timeout = 3600", true, false},
		{"SET GLOBAL query_timeout = 3600", false, false},
		{"SET @@global.query_timeout = 3600", false, false},
		{"SELECT 1", true, true},
		{"ADMIN SHOW FRONTEND CONFIG", false, true},
	}
	for _, tc := range tests {
		readOnly, returnsData, err := ValidateQuery(tc.statement)
		if err != nil {
			t.Errorf("%q: %v", tc.statement, err)
			continue
		}
		if readOnly != tc.readOnly || returnsData != tc.returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (%v, %v)", tc.statement, readOnly, returnsData, tc.readOnly, tc.returnsData)
		}
	}
}

// An explained statement is a plan, and only StarRocks's EXPLAIN ANALYZE runs
// the statement it plans.
func TestDescribeExplain(t *testing.T) {
	for _, statement := range []string{"SELECT 1", "SHOW TABLES"} {
		if _, _, ok := DescribeExplain(statement); ok {
			t.Errorf("%q describes a plan", statement)
		}
	}
	format, executed, ok := DescribeExplain("EXPLAIN SELECT 1")
	if !ok || format != "TEXT" || executed {
		t.Errorf("EXPLAIN = (%q, %v, %v), want (TEXT, false, true)", format, executed, ok)
	}
	format, executed, ok = DescribeExplain("EXPLAIN VERBOSE SELECT 1")
	if !ok || format != "VERBOSE" || executed {
		t.Errorf("EXPLAIN VERBOSE = (%q, %v, %v), want (VERBOSE, false, true)", format, executed, ok)
	}
}

// A buffer that cannot lex is not read-only anything: the chunk it sits in
// holds no statement, so parser.Split drops it and a per-segment parse never
// sees it. The check runs on the whole text before splitting.
func TestValidateQueryFailsClosedOnLexError(t *testing.T) {
	for _, script := range []string{"/*", "SELECT 1; /*"} {
		if readOnly, returnsData, err := ValidateQuery(script); err == nil {
			t.Errorf("ValidateQuery(%q) = %v, %v, nil; want the lex error", script, readOnly, returnsData)
		}
	}
}
