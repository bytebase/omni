package validate

import (
	"testing"

	parser "github.com/bytebase/omni/mysql/parser"
)

func TestValidateEmpty(t *testing.T) {
	diags := Validate(nil, Options{})
	if diags != nil {
		t.Fatalf("expected nil diagnostics, got %v", diags)
	}
}

func TestValidateCleanProcedure(t *testing.T) {
	list, err := parser.Parse(`CREATE PROCEDURE p() BEGIN DECLARE x INT; SET x = 1; END`)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	diags := Validate(list, Options{})
	if len(diags) != 0 {
		t.Fatalf("expected no diagnostics, got %v", diags)
	}
}
