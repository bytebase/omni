package catalog

import (
	"errors"
	"testing"
)

// TestUDFArgAliasRefused holds the catalog to MySQL 8.0.46's answers for an
// argument alias: 1583 on a native function, 1584 on any other name, which a
// stored function's call is.
func TestUDFArgAliasRefused(t *testing.T) {
	for sql, code := range map[string]int{
		"CREATE VIEW v AS SELECT abs(a x) FROM t":            1583,
		"CREATE VIEW v AS SELECT abs(a AS '') FROM t":        1583,
		"CREATE VIEW v AS SELECT nosuch(a x) FROM t":         1584,
		"CREATE VIEW v AS SELECT concat_ws(',', a x) FROM t": 1583,
	} {
		c := New()
		if _, err := c.Exec("CREATE DATABASE d; USE d; CREATE TABLE t (a INT)", nil); err != nil {
			t.Fatal(err)
		}
		results, err := c.Exec(sql, nil)
		if err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
		var e *Error
		if !errors.As(results[0].Error, &e) || e.Code != code {
			t.Errorf("%s: got %v, want %d", sql, results[0].Error, code)
		}
	}
}
