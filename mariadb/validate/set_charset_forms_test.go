package validate

import (
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
)

// TestSetCharsetFormsNameNoVariable checks NAMES and the charset forms raise no undeclared
// variable in a routine, wherever they sit in the list.
func TestSetCharsetFormsNameNoVariable(t *testing.T) {
	for _, body := range []string{
		"SET NAMES utf8mb4",
		"SET NAMES utf8mb4 COLLATE utf8mb4_bin",
		"SET NAMES DEFAULT",
		"SET CHARSET utf8mb4",
		"SET CHARACTER SET utf8mb4",
		"SET CHAR SET utf8mb4",
		"SET @a = 1, NAMES utf8mb4",
		"SET CHARSET utf8mb4, @a = 1",
	} {
		t.Run(body, func(t *testing.T) {
			list, err := parser.Parse("CREATE PROCEDURE p() BEGIN " + body + "; END")
			if err != nil {
				t.Fatal(err)
			}
			if diags := Validate(list, Options{}); len(diags) != 0 {
				t.Errorf("findings = %v, want none", diags)
			}
		})
	}
	list, err := parser.Parse("CREATE PROCEDURE p() BEGIN SET NAMES utf8mb4, nope = 1; END")
	if err != nil {
		t.Fatal(err)
	}
	if diags := Validate(list, Options{}); len(diags) != 1 {
		t.Errorf("findings = %v, want the undeclared variable after NAMES", diags)
	}
}
