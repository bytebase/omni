package validate

import (
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
)

// TestSetStatementInRoutine checks a routine's wrapped statement is validated, and its options
// are not taken for its variables.
func TestSetStatementInRoutine(t *testing.T) {
	for _, tc := range []struct {
		body string
		want int
	}{
		{"SET STATEMENT nope = 1 FOR SELECT 1", 0},
		{"SET STATEMENT max_statement_time = 1 FOR SET nope = 1", 1},
		{"DECLARE x INT; SET STATEMENT max_statement_time = 1 FOR SET x = 1", 0},
		{"SET STATEMENT max_statement_time = 1 FOR IF 1 THEN SET nope = 1; END IF", 1},
		{"SET STATEMENT max_statement_time = 1 FOR SET STATEMENT sql_mode = '' FOR SET nope = 1", 1},
	} {
		t.Run(tc.body, func(t *testing.T) {
			list, err := parser.Parse("CREATE PROCEDURE p() BEGIN " + tc.body + "; END")
			if err != nil {
				t.Fatal(err)
			}
			if diags := Validate(list, Options{}); len(diags) != tc.want {
				t.Errorf("findings = %v, want %d", diags, tc.want)
			}
		})
	}
}
