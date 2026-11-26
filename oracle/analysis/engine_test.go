package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/oracle/ast"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1 FROM DUAL;\n;\nSELECT 2 FROM DUAL;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1 FROM DUAL" || !segments[1].Empty || segments[2].Text != "\nSELECT 2 FROM DUAL" {
		t.Fatalf("Split = %+v", segments)
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELECT 2 FROM DUAL" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T, want the statement inside the RawStmt", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(""); err != nil || parsed != nil {
		t.Fatalf("Parse of an empty script = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorIsPositionedInTheScript(t *testing.T) {
	for _, tc := range []struct {
		script string
		want   analysis.Position
	}{
		// A column on the segment's first line moves by the segment's column.
		{"SELECT 1 FROM DUAL; SELECT * FROM", analysis.Position{Line: 1, Column: 34}},
		{"SELECT 1 FROM DUAL;\n  SELECT * FROM", analysis.Position{Line: 2, Column: 16}},
	} {
		_, err := Engine.Parse(tc.script)
		var syntaxErr *analysis.SyntaxError
		if !errors.As(err, &syntaxErr) {
			t.Fatalf("Parse(%q) err = %v", tc.script, err)
		}
		if syntaxErr.Position != tc.want || syntaxErr.RawMessage != "syntax error at end of input" {
			t.Fatalf("Parse(%q) err = %+v", tc.script, syntaxErr)
		}
	}
	_, err := Engine.Parse("SELECT 1 FROM DUAL;\n  SELECT * FROM")
	if err.Error() != "Syntax error at line 2:16: syntax error at end of input" {
		t.Fatalf("message = %q", err.Error())
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1 FROM DUAL;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1 FROM DUAL;\nSELECT * FROM")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 2, Column: 14}) || diagnostics[0].Message != "syntax error at end of input" {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateQuery(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT 1 FROM DUAL"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	_, _, err := Engine.ValidateQuery("SELECT 1 FROM DUAL;\nSELECT * FROM")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) || syntaxErr.Message != "Syntax error at line 2:14: syntax error at end of input" {
		t.Fatalf("ValidateQuery err = %#v", err)
	}
}

func TestEngineResultLimit(t *testing.T) {
	for _, tc := range []struct {
		statement, version, want string
	}{
		{"SELECT * FROM t;", "19.3.0.0.0", "SELECT * FROM t FETCH NEXT 10 ROWS ONLY"},
		{"SELECT * FROM t;", "11.2.0.4.0", "SELECT * FROM (SELECT * FROM t) WHERE ROWNUM <= 10"},
		{"SELECT * FROM t", "", "SELECT * FROM (SELECT * FROM t) WHERE ROWNUM <= 10"},
		// A FETCH NEXT rewrite that fails falls back to ROWNUM.
		{"SELECT * FROM t FETCH FIRST 50 PERCENT ROWS ONLY", "19.3.0.0.0", "SELECT * FROM (SELECT * FROM t FETCH FIRST 50 PERCENT ROWS ONLY) WHERE ROWNUM <= 10"},
	} {
		if got := Engine.ResultLimit(tc.statement, 10, tc.version); got != tc.want {
			t.Fatalf("ResultLimit(%q, %q) = %q", tc.statement, tc.version, got)
		}
	}
}
