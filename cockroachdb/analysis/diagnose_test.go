package analysis

import (
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
		// CockroachDB-only syntax the PostgreSQL grammar rejects.
		{script: "UPSERT INTO t VALUES (1); SELECT * FROM t AS OF SYSTEM TIME '-10s'; SELECT * FROM t@idx"},
		{script: "ALTER TABLE t CONFIGURE ZONE USING num_replicas = 3"},
		{script: "SELEC 1", want: &DiagnoseError{Offset: 0, Message: `at or near "selec": syntax error`}},
		{script: "SELECT 1;\nSELECT * FROM", want: &DiagnoseError{Offset: len("SELECT 1;\nSELECT * FROM"), Message: `at or near "EOF": syntax error`}},
		{script: "SELECT '你好' FROM WHERE", want: &DiagnoseError{Offset: len("SELECT '你好' FROM "), Message: `at or near "where": syntax error`}},
	} {
		t.Run(tc.script, func(t *testing.T) {
			if diff := cmp.Diff(tc.want, Diagnose(tc.script)); diff != "" {
				t.Errorf("Diagnose() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDiagnoseAgreesWithEngine(t *testing.T) {
	for _, row := range corpus(t) {
		accepted := Diagnose(row.Statement) == nil
		if want := row.Expected == compat.Accept; accepted != want {
			t.Errorf("Diagnose(%q) accepts=%v, CockroachDB accepts=%v", row.Statement, accepted, want)
		}
	}
}
