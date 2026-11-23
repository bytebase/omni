package analysis

import (
	"strings"
	"testing"
)

func TestDiagnose(t *testing.T) {
	for _, statement := range []string{
		"SELECT * FROM person",
		"SELECT * FROM person LATERAL VIEW EXPLODE(ARRAY(30, 60)) tableName AS c_age;",
		"SELECT * FROM schema1.t1 JOIN schema2.t2 ON t1.c1 = t2.c1",
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

// Diagnose parses the whole script, so a chunk that holds no statement still
// reports its lex error. SplitStatements and Classify have no error channel:
// the first keeps the chunk marked Empty, and the second answers Unknown for
// text whose first token never arrives, which a permission check denies.
func TestDiagnoseReportsDroppedChunkLexError(t *testing.T) {
	const script = "/*"
	diagnostics := Diagnose(script)
	if len(diagnostics) == 0 {
		t.Fatalf("Diagnose(%q) returned no diagnostics", script)
	}
	if !strings.Contains(diagnostics[0].Message, "unterminated block comment") {
		t.Errorf("message = %q, want the unterminated comment", diagnostics[0].Message)
	}
	if diagnostics[0].Offset != 0 {
		t.Errorf("offset = %d, want 0", diagnostics[0].Offset)
	}
	if got := ClassifyStatement(script); got != QueryTypeUnknown {
		t.Errorf("ClassifyStatement(%q) = %s, want %s", script, got, QueryTypeUnknown)
	}
	if statements := SplitStatements(script); len(statements) != 1 || !statements[0].Empty {
		t.Errorf("SplitStatements(%q) = %+v, want one Empty statement", script, statements)
	}
}
