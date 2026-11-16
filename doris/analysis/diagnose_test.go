package analysis

import (
	"strings"
	"testing"
)

func TestDiagnose(t *testing.T) {
	for _, statement := range []string{
		"SELECT * FROM person",
		"SELECT a FROM t WHERE b = 1",
		"CREATE TABLE t (id INT)",
		"INSERT INTO t VALUES (1)",
		"",
	} {
		if got := Diagnose(statement); len(got) != 0 {
			t.Errorf("Diagnose(%q) = %+v, want no diagnostics", statement, got)
		}
	}
}

func TestDiagnoseReportsSyntaxError(t *testing.T) {
	const statement = "SELECT a > (select max(a) from t1) FROM"
	diagnostics := Diagnose(statement)
	if len(diagnostics) == 0 {
		t.Fatalf("Diagnose(%q) returned no diagnostics", statement)
	}
	if !strings.Contains(diagnostics[0].Message, "syntax error") {
		t.Errorf("message = %q, want a syntax error", diagnostics[0].Message)
	}
	if diagnostics[0].Offset < 0 || diagnostics[0].Offset > len(statement) {
		t.Errorf("offset = %d, want within the statement", diagnostics[0].Offset)
	}
}

func TestDiagnoseFiltersUnsupportedStub(t *testing.T) {
	// The parser reports a statement-dispatch gap as a diagnostic; it is not
	// an error in the SQL, so Diagnose drops it.
	if got := Diagnose("REFRESH DATABASE db"); len(got) != 0 {
		t.Errorf("Diagnose(REFRESH DATABASE db) = %+v, want no diagnostics", got)
	}
}
