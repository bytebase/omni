package catalog

import (
	"strings"
	"testing"
)

// TestEnumDuplicateValueMessage holds 1291's message to the server's: the value
// named is the first spelling, which MariaDB compares case-insensitively.
// Measured on mariadb:11.8.8: `CREATE TABLE t (e ENUM('s','S'))` answers
// "Column 'e' has duplicated value 's' in ENUM", and a SET the same with SET.
func TestEnumDuplicateValueMessage(t *testing.T) {
	for _, tc := range []struct{ ddl, want string }{
		{"CREATE TABLE t (e ENUM('s', 'S'))", "Column 'e' has duplicated value 's' in ENUM"},
		{"CREATE TABLE t (e SET('a', 'b', 'A'))", "Column 'e' has duplicated value 'a' in SET"},
	} {
		err := execErr(t, tc.ddl)
		catErr, ok := err.(*Error)
		if !ok || catErr.Code != 1291 {
			t.Fatalf("%s: expected 1291, got %v", tc.ddl, err)
		}
		if !strings.Contains(catErr.Message, tc.want) {
			t.Errorf("%s: Message = %q, want %q", tc.ddl, catErr.Message, tc.want)
		}
	}
}
