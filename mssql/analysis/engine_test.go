package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/mssql/ast"
	"github.com/bytebase/omni/mssql/parser"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT 1;\nGO\nSELECT 2;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT 1;" || !segments[1].Empty || segments[2].Text != "\nSELECT 2;" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v does not tile the script", s)
		}
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "\nSELECT 2;" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(""); err != nil || parsed != nil {
		t.Fatalf("Parse of an empty script = %#v, %v", parsed, err)
	}
}

func TestEngineParseErrorKeepsTheSegmentsColumn(t *testing.T) {
	_, err := Engine.Parse("SELECT 1;\nSELECT 2; SELECT * FROM;")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 15}) || syntaxErr.Message != `Syntax error at line 2: syntax error at or near ";"` || syntaxErr.RawMessage != `syntax error at or near ";"` {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineDiagnose(t *testing.T) {
	if diagnostics := Engine.Diagnose("SELECT 1;"); diagnostics == nil || len(diagnostics) != 0 {
		t.Fatalf("Diagnose of valid SQL = %#v", diagnostics)
	}
	diagnostics := Engine.Diagnose("SELECT 1; SELECT * FROM;")
	if len(diagnostics) != 1 || diagnostics[0].Start != (analysis.Position{Line: 1, Column: 24}) || diagnostics[0].Message != `syntax error at or near ";"` {
		t.Fatalf("Diagnose = %+v", diagnostics)
	}
}

func TestEngineValidateQueryAndLimit(t *testing.T) {
	if valid, allQuery, err := Engine.ValidateQuery("SET NOCOUNT ON; SELECT 1;"); err != nil || !valid || allQuery {
		t.Fatalf("ValidateQuery = %v %v %v", valid, allQuery, err)
	}
	var parseErr *parser.ParseError
	if _, _, err := Engine.ValidateQuery("SELECT FROM ;"); !errors.As(err, &parseErr) {
		t.Fatalf("ValidateQuery err = %#v, want the parser's own error", err)
	}
	if got := Engine.ResultLimit("SELECT * FROM dbo.t;", 10, ""); got != "SELECT TOP 10 * FROM dbo.t;" {
		t.Fatalf("ResultLimit = %q", got)
	}
	if got := Engine.ResultLimit("SELECT FROM ;", 10, ""); got != "WITH result AS (SELECT FROM) SELECT TOP 10 * FROM result;" {
		t.Fatalf("ResultLimit fallback = %q", got)
	}
}
