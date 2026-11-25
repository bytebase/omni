package analysis

import (
	"path/filepath"
	"testing"

	"github.com/bytebase/omni/enginetest/compat"

	"github.com/google/go-cmp/cmp"
)

func TestDiagnose(t *testing.T) {
	for _, tc := range []struct {
		script string
		want   *DiagnoseError
	}{
		{script: "SELECT 1"},
		{script: ""},
		{script: "-- comment only"},
		{script: ";;"},
		// ClickHouse-only syntax a MySQL grammar rejects.
		{script: "SELECT * FROM t FINAL SAMPLE 0.1 SETTINGS max_threads = 1; SELECT arrayJoin([1, 2]) FORMAT JSONEachRow"},
		{script: "CREATE TABLE t (a UInt64) ENGINE = MergeTree ORDER BY a"},
		{script: "INSERT INTO t FORMAT CSV\n1,2;3\n"},
		{script: "SELEC 1", want: &DiagnoseError{Offset: 0, Message: `unexpected token "SELEC" at start of query`}},
		{script: "SELECT 1;\nSELECT * FROM", want: &DiagnoseError{Offset: len("SELECT 1;\nSELECT * FROM"), Message: `expected identifier, got ""`}},
		{script: "SELECT '你好' FROM t WHERE )", want: &DiagnoseError{Offset: len("SELECT '你好' FROM t WHERE "), Message: `unexpected token ")" in expression`}},
		{script: "SELECT 'open", want: &DiagnoseError{Offset: len("SELECT "), Message: "unterminated string literal"}},
	} {
		t.Run(tc.script, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, Diagnose(tc.script)); diff != "" {
				t.Errorf("Diagnose() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestDiagnoseAgreesWithEngine holds Diagnose to ClickHouse's verdict on
// every statement of the parser's compat.tsv that omni judges as ClickHouse
// does.
func TestDiagnoseAgreesWithEngine(t *testing.T) {
	f, err := compat.ReadFile(filepath.Join("..", "parser", "testdata", "compat.tsv"))
	if err != nil {
		t.Fatal(err)
	}
	for _, row := range f.Rows {
		if row.Omni != row.Expected {
			continue
		}
		accepted := Diagnose(row.Statement) == nil
		if want := row.Expected == compat.Accept; accepted != want {
			t.Errorf("Diagnose(%q) accepts=%v, ClickHouse accepts=%v", row.Statement, accepted, want)
		}
	}
}
