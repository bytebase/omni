package catalog

import (
	"errors"
	"testing"
)

// TestProgramSystemVariables pins the catalog to mariadb:11.8.8 on a stored program's SET of a
// name it doesn't declare: 1193 for a variable the server lacks, raised while MariaDB parses the
// program and so before the missing database (1049) or the existing procedure (1304).
func TestProgramSystemVariables(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		code int
	}{
		{"CREATE PROCEDURE p2() SET zzz = 1", 1193},
		{"CREATE PROCEDURE p() SET zzz = 1", 1193},
		{"CREATE FUNCTION nosuch.f2() RETURNS INT BEGIN SET zzz = 1; RETURN 1; END", 1193},
		{"CREATE TRIGGER tr BEFORE INSERT ON t FOR EACH ROW SET zzz = 1", 1193},
		{"CREATE EVENT e2 ON SCHEDULE EVERY 1 DAY DO SET zzz = 1", 1193},
		{"CREATE PROCEDURE p2(zzz INT) BEGIN DECLARE yyy INT; SET zzz = 1, yyy = 2, sql_mode = ''; END", 0},
		{"SET zzz = 1", 0},
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT); CREATE PROCEDURE p() COMMIT", nil); err != nil {
			t.Fatal(err)
		}
		r, err := c.Exec(tc.sql, &ExecOptions{ContinueOnError: true})
		if err != nil {
			t.Fatalf("%s: %v", tc.sql, err)
		}
		code := 0
		var e *Error
		if errors.As(r[0].Error, &e) {
			code = e.Code
		} else if r[0].Error != nil {
			t.Fatalf("%s: %v", tc.sql, r[0].Error)
		}
		if code != tc.code {
			t.Errorf("%s: code %d, want %d (%v)", tc.sql, code, tc.code, r[0].Error)
		}
	}
}
