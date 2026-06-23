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

// TestSystemVersionedRequiresExplicit pins MariaDB's rule (error 4125): a table
// declaring PERIOD FOR SYSTEM_TIME or a ROW START/END column must be explicitly
// system-versioned — via a table-level WITH SYSTEM VERSIONING or a column-level
// WITH SYSTEM VERSIONING. The catalog must reject the bare form, not fabricate
// versioning. Container-verified vs mariadb:11.8.8.
func TestSystemVersionedRequiresExplicit(t *testing.T) {
	reject := []string{
		"CREATE TABLE nov1 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re))",
		"CREATE TABLE nov2 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START)",
	}
	for _, sql := range reject {
		t.Run("reject", func(t *testing.T) {
			err := execErr(t, sql)
			if err == nil {
				t.Fatalf("expected error 4125, got nil for %q", sql)
			}
			catErr, ok := err.(*Error)
			if !ok {
				t.Fatalf("expected *Error, got %T", err)
			}
			if catErr.Code != 4125 {
				t.Errorf("Code = %d, want 4125 for %q", catErr.Code, sql)
			}
		})
	}
	accept := []string{
		"CREATE TABLE okv1 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		"CREATE TABLE okv2 (x INT WITH SYSTEM VERSIONING)",
		"CREATE TABLE okv3 (x INT WITHOUT SYSTEM VERSIONING)",
	}
	for _, sql := range accept {
		t.Run("accept", func(t *testing.T) {
			if err := execErr(t, sql); err != nil {
				t.Errorf("unexpected error for %q: %v", sql, err)
			}
		})
	}
}

func execErr(t *testing.T, ddl string) error {
	t.Helper()
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	results, _ := c.Exec(ddl, &ExecOptions{ContinueOnError: true})
	return results[0].Error
}

func defineAndShow(t *testing.T, table, ddl string) string {
	t.Helper()
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec(ddl, nil)
	return c.ShowCreateTable("test", table)
}
