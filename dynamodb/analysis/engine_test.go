package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineSplit(t *testing.T) {
	script := "SELECT * FROM Music;\n-- c\n;\n\tSELECT * FROM Books"
	segments := Engine.Split(script)
	if len(segments) != 3 || segments[0].Text != "SELECT * FROM Music;" || !segments[1].Empty || segments[2].Text != "\n\tSELECT * FROM Books" {
		t.Fatalf("Split = %+v", segments)
	}
	for _, s := range segments {
		if script[s.Start:s.End] != s.Text {
			t.Fatalf("segment %+v is not the script's text", s)
		}
	}
}

func TestEngineValidateQuery(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("SELECT * FROM Music"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery(SELECT) = %v %v %v", readOnly, returnsData, err)
	}
	if readOnly, returnsData, err := Engine.ValidateQuery("DELETE FROM Music WHERE Artist = 'x'"); err != nil || readOnly || returnsData {
		t.Fatalf("ValidateQuery(DELETE) = %v %v %v", readOnly, returnsData, err)
	}
	_, _, err := Engine.ValidateQuery("SELECT * FROM Music;\nSELECT FROM WHERE")
	var syntaxErr *analysis.SyntaxError
	if !errors.As(err, &syntaxErr) {
		t.Fatalf("err = %v", err)
	}
	if syntaxErr.Position != (analysis.Position{Line: 2, Column: 8}) || syntaxErr.Message != "syntax error at position 28: "+syntaxErr.RawMessage {
		t.Fatalf("err = %+v", syntaxErr)
	}
}
