package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
)

func TestEngineValidateQuery(t *testing.T) {
	if readOnly, returnsData, err := Engine.ValidateQuery("GET a\nHGETALL h"); err != nil || !readOnly || !returnsData {
		t.Fatalf("ValidateQuery = %v %v %v", readOnly, returnsData, err)
	}
	if readOnly, _, err := Engine.ValidateQuery("SET a 1"); err != nil || readOnly {
		t.Fatalf("ValidateQuery(SET) = %v %v", readOnly, err)
	}
	var syntaxErr *analysis.SyntaxError
	if _, _, err := Engine.ValidateQuery(`GET "a`); err == nil || errors.As(err, &syntaxErr) {
		t.Fatalf("ValidateQuery returns the splitter's own error, got %v", err)
	}
}
