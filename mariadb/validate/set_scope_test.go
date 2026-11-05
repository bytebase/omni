package validate

import (
	"testing"

	"github.com/bytebase/omni/mariadb/parser"
)

// TestSetAssignmentScope checks which SET targets count as a routine's local variables.
// The most recent scope reaches an assignment that writes none.
func TestSetAssignmentScope(t *testing.T) {
	for _, tc := range []struct {
		body string
		want int
	}{
		{"SET nope = 1", 1},
		{"SET @a = 1, nope = 2", 1},
		{"SET GLOBAL nope = 1", 0},
		{"SET SESSION nope = 1", 0},
		{"SET @a = 1, GLOBAL nope = 2", 0},
		{"SET @a = 1, SESSION nope = 2", 0},
		{"SET @a = 1, GLOBAL nope = 2, also_nope = 3", 0},
		{"SET SESSION sql_mode = '', also_nope = 3", 0},
		{"DECLARE x INT; SET x = 1, GLOBAL nope = 2", 0},
	} {
		t.Run(tc.body, func(t *testing.T) {
			list, err := parser.Parse("CREATE PROCEDURE p() BEGIN " + tc.body + "; END")
			if err != nil {
				t.Fatal(err)
			}
			if got := len(Validate(list, Options{})); got != tc.want {
				t.Errorf("%d findings, want %d", got, tc.want)
			}
		})
	}
}
