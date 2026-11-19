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
		{"DESC FILES('path' = 's3://x')", true, true},
		{"DESC PIPE db.p", true, true},
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

// CACHE SELECT returns the cache warm-up metrics, so the driver reads rows.
func TestValidateQueryDataCacheSelect(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("CACHE SELECT a FROM t")
	if err != nil {
		t.Fatalf("ValidateQuery: %v", err)
	}
	if readOnly || !returnsData {
		t.Errorf("CACHE SELECT = (readOnly %v, returnsData %v), want (false, true)", readOnly, returnsData)
	}
}

// SUBMIT TASK returns the task name and its status, so the driver must read
// rows even though the statement is not read-only.
func TestValidateQuerySubmitTask(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("SUBMIT TASK AS INSERT INTO t SELECT 1")
	if err != nil {
		t.Fatalf("ValidateQuery: %v", err)
	}
	if readOnly || !returnsData {
		t.Errorf("SUBMIT TASK = (readOnly %v, returnsData %v), want (false, true)", readOnly, returnsData)
	}
}

// A plan-only DESC accepts a SELECT that carries INTO OUTFILE: the plan does
// not perform the export.
func TestValidateQueryPlanOnlyOutfile(t *testing.T) {
	readOnly, returnsData, err := ValidateQuery("DESC SELECT 1 INTO OUTFILE 'x'")
	if err != nil {
		t.Fatalf("ValidateQuery: %v", err)
	}
	if !readOnly || !returnsData {
		t.Errorf("DESC SELECT INTO OUTFILE = (readOnly %v, returnsData %v), want (true, true)", readOnly, returnsData)
	}
}

// EXPLAIN ANALYZE over DML is not read-only, but its profile is still a
// result set the driver must read rather than Exec.
func TestValidateQueryAnalyzeInsert(t *testing.T) {
	for _, sql := range []string{
		"EXPLAIN ANALYZE INSERT INTO t SELECT 1",
		"DESC ANALYZE INSERT INTO t SELECT 1",
	} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (false, true)", sql, readOnly, returnsData)
		}
	}
}

// TRACE plans the query it wraps and never runs it, so it is read-only like
// EXPLAIN, including over a query that writes out a file.
func TestValidateQueryOptimizerTrace(t *testing.T) {
	for _, sql := range []string{
		"TRACE ALL SELECT 1",
		"TRACE LOGS SELECT * FROM t",
		"TRACE ALL SELECT 1 INTO OUTFILE 's3://bucket/x'",
	} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if !readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (true, true)", sql, readOnly, returnsData)
		}
	}
}

// DESC STORAGE VOLUME answers with a storage volume's settings, so the
// editor's read-only mode admits it the way it admits DESC over a table.
func TestValidateQueryDescStorageVolume(t *testing.T) {
	for _, sql := range []string{"DESC STORAGE VOLUME v", "DESCRIBE STORAGE VOLUME v"} {
		readOnly, returnsData, err := ValidateQuery(sql)
		if err != nil {
			t.Errorf("%q: %v", sql, err)
			continue
		}
		if !readOnly || !returnsData {
			t.Errorf("%q = (readOnly %v, returnsData %v), want (true, true)", sql, readOnly, returnsData)
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
	format, executed, ok = DescribeExplain("TRACE ALL SELECT 1")
	if !ok || format != "TRACE" || executed {
		t.Errorf("TRACE = (%q, %v, %v), want (TRACE, false, true)", format, executed, ok)
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
