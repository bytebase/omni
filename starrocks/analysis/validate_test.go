package analysis

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/starrocks/parser"
)

// A statement cut out of a script reports its position in the script, not in
// the segment: a caller highlighting the error must land on the right text.
func TestValidateQueryReportsScriptPositions(t *testing.T) {
	const script = "SELECT 1;\nSELECT FROM"
	_, _, err := ValidateQuery(script)
	if err == nil {
		t.Fatal("ValidateQuery accepted a syntax error")
	}
	errs := parser.AllErrors(err)
	if len(errs) == 0 {
		t.Fatalf("ValidateQuery error = %v, want a ParseError", err)
	}
	if want := strings.Index(script, "FROM"); errs[0].Position != want {
		t.Errorf("error position = %d, want %d", errs[0].Position, want)
	}
}
