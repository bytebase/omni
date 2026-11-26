package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/cosmosdb/parser"
)

func TestEngineValidateQuery(t *testing.T) {
	for _, script := range []string{"", "SELECT * FROM c", "SELECT c.id FROM c WHERE c.a = 1"} {
		if readOnly, returnsData, err := Engine.ValidateQuery(script); err != nil || !readOnly || !returnsData {
			t.Fatalf("ValidateQuery(%q) = %v %v %v", script, readOnly, returnsData, err)
		}
	}
	readOnly, returnsData, err := Engine.ValidateQuery("SELECT FROM")
	var parseErr *parser.ParseError
	if !errors.As(err, &parseErr) || readOnly || returnsData {
		t.Fatalf("ValidateQuery(invalid) = %v %v %v", readOnly, returnsData, err)
	}
}
