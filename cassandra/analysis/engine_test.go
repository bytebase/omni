package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/cassandra/ast"
)

func TestEngineSplitAndParse(t *testing.T) {
	script := "SELECT * FROM a;\n;\n  SELECT * FROM b;"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT * FROM a" || !segments[1].Empty || segments[2].Text != "SELECT * FROM b" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v is not the script's text", s)
		}
	}

	parsed, err := Engine.Parse(script)
	if err != nil {
		t.Fatal(err)
	}
	if len(parsed) != 2 || parsed[1].Text != "SELECT * FROM b" {
		t.Fatalf("Parse kept %d statements: %+v", len(parsed), parsed)
	}
	if _, ok := parsed[0].Node.(*ast.SelectStmt); !ok {
		t.Fatalf("node = %T", parsed[0].Node)
	}
	if parsed, err := Engine.Parse(";"); err != nil || parsed != nil {
		t.Fatalf("Parse of no statement = %+v, %v", parsed, err)
	}
}

func TestEngineSyntaxErrorShiftsLineAndColumn(t *testing.T) {
	_, err := Engine.Parse("SELECT * FROM t;\n  SELCT * FROM users")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 3}) || syntaxErr.Message != "Syntax error at line 2: "+syntaxErr.RawMessage {
		t.Fatalf("err = %+v", syntaxErr)
	}
}

func TestEngineValidateQueryReturnsTheParsersError(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT * FROM t"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery("SELECT FROM"); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery(SELECT FROM) err = %v", err)
	}
}
