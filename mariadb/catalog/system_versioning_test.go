package catalog

import (
	"strings"
	"testing"
)

// TestShowCreateTableSystemVersioned guards that a system-versioned table
// round-trips through DefineTable + ShowCreateTable without producing invalid
// DDL: row-start/end columns must render as GENERATED ALWAYS AS ROW START/END
// (never the empty-expression "AS ()"), and the PERIOD / WITH SYSTEM VERSIONING
// clauses must survive. Targets match mariadb:11.8.8 SHOW CREATE TABLE.
func TestShowCreateTableSystemVersioned(t *testing.T) {
	t.Run("canonical period columns", func(t *testing.T) {
		ddl := defineAndShow(t, "sv2",
			"CREATE TABLE sv2 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING")
		if strings.Contains(ddl, "AS ()") {
			t.Errorf("renders broken empty-expr generated column:\n%s", ddl)
		}
		for _, want := range []string{
			"GENERATED ALWAYS AS ROW START",
			"GENERATED ALWAYS AS ROW END",
			"PERIOD FOR SYSTEM_TIME (`rs`, `re`)",
			"WITH SYSTEM VERSIONING",
		} {
			if !strings.Contains(ddl, want) {
				t.Errorf("missing %q in:\n%s", want, ddl)
			}
		}
	})
	t.Run("column-level versioning", func(t *testing.T) {
		ddl := defineAndShow(t, "sv4",
			"CREATE TABLE sv4 (x INT WITH SYSTEM VERSIONING, y INT WITHOUT SYSTEM VERSIONING)")
		if !strings.Contains(ddl, "WITHOUT SYSTEM VERSIONING") {
			t.Errorf("missing column WITHOUT SYSTEM VERSIONING in:\n%s", ddl)
		}
		if !strings.Contains(ddl, "WITH SYSTEM VERSIONING") {
			t.Errorf("missing table WITH SYSTEM VERSIONING in:\n%s", ddl)
		}
	})
}

func defineAndShow(t *testing.T, table, ddl string) string {
	t.Helper()
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec(ddl, nil)
	return c.ShowCreateTable("test", table)
}
