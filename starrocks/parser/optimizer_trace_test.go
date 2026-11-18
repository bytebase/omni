package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestOptimizerTrace covers optimizerTrace (StarRocks.g4:2584): every mode,
// the optional trace identifier, and the query it wraps.
func TestOptimizerTrace(t *testing.T) {
	for _, tc := range []struct {
		sql        string
		mode       string
		identifier string
	}{
		{"TRACE ALL SELECT 1", "ALL", ""},
		{"TRACE LOGS SELECT 1", "LOGS", ""},
		{"TRACE TIMES SELECT 1", "TIMES", ""},
		{"TRACE VALUES SELECT 1", "VALUES", ""},
		{"TRACE REASON SELECT 1", "REASON", ""},
		{"TRACE ALL trace_name SELECT 1 FROM t", "ALL", "trace_name"},
		{"TRACE ALL WITH c AS (SELECT 1) SELECT * FROM c", "ALL", ""},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := parseForTest(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt, ok := file.Stmts[0].(*ast.OptimizerTraceStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.OptimizerTraceStmt", file.Stmts[0])
			}
			if stmt.Mode != tc.mode || stmt.Identifier != tc.identifier {
				t.Errorf("(Mode, Identifier) = (%q, %q), want (%q, %q)", stmt.Mode, stmt.Identifier, tc.mode, tc.identifier)
			}
			if stmt.Query == nil {
				t.Error("Query is nil")
			}
		})
	}
}

func TestOptimizerTraceRejects(t *testing.T) {
	for _, sql := range []string{
		"TRACE SELECT 1",
		"TRACE ALL zzz",
		"TRACE ALL a zzz SELECT 1",
		"TRACE ALL INSERT INTO t SELECT 1",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}
