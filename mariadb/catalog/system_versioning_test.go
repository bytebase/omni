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

// TestSystemVersioningAlterCatalog verifies ALTER ... ADD/DROP SYSTEM VERSIONING
// mutates the catalog (not a silent no-op), and that the whole-table consistency
// check still applies after the ALTER. Container-verified vs mariadb:11.8.8.
func TestSystemVersioningAlterCatalog(t *testing.T) {
	t.Run("add system versioning reflects in show create", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE av (x INT)")
		if r, _ := c.Exec("ALTER TABLE av ADD SYSTEM VERSIONING", &ExecOptions{ContinueOnError: true}); r[0].Error != nil {
			t.Fatalf("ADD SYSTEM VERSIONING errored: %v", r[0].Error)
		}
		if ddl := c.ShowCreateTable("test", "av"); !strings.Contains(ddl, "WITH SYSTEM VERSIONING") {
			t.Errorf("ShowCreateTable missing WITH SYSTEM VERSIONING after ADD:\n%s", ddl)
		}
	})
	t.Run("drop system versioning reverts show create", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE av (x INT) WITH SYSTEM VERSIONING")
		c.Exec("ALTER TABLE av DROP SYSTEM VERSIONING", &ExecOptions{ContinueOnError: true})
		if ddl := c.ShowCreateTable("test", "av"); strings.Contains(ddl, "WITH SYSTEM VERSIONING") {
			t.Errorf("ShowCreateTable still versioned after DROP:\n%s", ddl)
		}
	})
	t.Run("multi-command convert plain table to versioned", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE conv (x INT)")
		r, _ := c.Exec("ALTER TABLE conv ADD COLUMN rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, ADD COLUMN re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, ADD PERIOD FOR SYSTEM_TIME(rs, re), ADD SYSTEM VERSIONING", &ExecOptions{ContinueOnError: true})
		if r[0].Error != nil {
			t.Fatalf("convert ALTER errored: %v", r[0].Error)
		}
		ddl := c.ShowCreateTable("test", "conv")
		if strings.Contains(ddl, "AS ()") {
			t.Errorf("renders broken AS () on the ALTER path:\n%s", ddl)
		}
		for _, want := range []string{"GENERATED ALWAYS AS ROW START", "PERIOD FOR SYSTEM_TIME (`rs`, `re`)", "WITH SYSTEM VERSIONING"} {
			if !strings.Contains(ddl, want) {
				t.Errorf("converted table missing %q:\n%s", want, ddl)
			}
		}
	})
	t.Run("drop system versioning with explicit period cols rejected (4125)", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE sv (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING")
		r, _ := c.Exec("ALTER TABLE sv DROP SYSTEM VERSIONING", &ExecOptions{ContinueOnError: true})
		catErr, ok := r[0].Error.(*Error)
		if !ok || catErr.Code != 4125 {
			t.Fatalf("expected error 4125, got %v", r[0].Error)
		}
		// The failed ALTER must roll back: the table is still versioned.
		if ddl := c.ShowCreateTable("test", "sv"); !strings.Contains(ddl, "WITH SYSTEM VERSIONING") {
			t.Errorf("rollback failed — table no longer versioned:\n%s", ddl)
		}
	})
}

// TestSystemVersioningConsistency covers the structural rules MariaDB enforces
// on system-versioned tables: ROW START/END columns and PERIOD FOR SYSTEM_TIME
// must be present together and reference each other. All reject vs mariadb:11.8.8.
func TestSystemVersioningConsistency(t *testing.T) {
	reject := []string{
		// row columns but no PERIOD
		"CREATE TABLE c1 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END) WITH SYSTEM VERSIONING",
		// PERIOD on ordinary columns (no ROW START/END)
		"CREATE TABLE c3 (a INT, b INT, PERIOD FOR SYSTEM_TIME(a, b)) WITH SYSTEM VERSIONING",
		// PERIOD references columns that are not the ROW START/END columns
		"CREATE TABLE c5 (a TIMESTAMP(6), b TIMESTAMP(6), rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(a, b)) WITH SYSTEM VERSIONING",
	}
	for _, sql := range reject {
		t.Run("create", func(t *testing.T) {
			if _, ok := execErr(t, sql).(*Error); !ok {
				t.Errorf("expected rejection for %q, got accepted", sql)
			}
		})
	}

	// ALTER ADD COLUMN carrying versioning metadata onto a plain table is rejected
	// (it would persist orphaned versioning metadata), while the multi-command
	// convert (which ends versioned) stays valid.
	t.Run("alter add row column to plain table rejected", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE pl (x INT)")
		r, _ := c.Exec("ALTER TABLE pl ADD COLUMN rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START", &ExecOptions{ContinueOnError: true})
		if _, ok := r[0].Error.(*Error); !ok {
			t.Errorf("expected rejection, got accepted; ShowCreate:\n%s", c.ShowCreateTable("test", "pl"))
		}
	})
	t.Run("alter add WITH SYSTEM VERSIONING column to plain table rejected", func(t *testing.T) {
		c := versionedCatalog(t, "CREATE TABLE pl (x INT)")
		r, _ := c.Exec("ALTER TABLE pl ADD COLUMN z INT WITH SYSTEM VERSIONING", &ExecOptions{ContinueOnError: true})
		if _, ok := r[0].Error.(*Error); !ok {
			t.Errorf("expected rejection, got accepted; ShowCreate:\n%s", c.ShowCreateTable("test", "pl"))
		}
	})
}

// TestSystemVersioningCaseInsensitivePeriod: MariaDB matches the PERIOD FOR
// SYSTEM_TIME columns against the ROW START/END column names case-insensitively.
func TestSystemVersioningCaseInsensitivePeriod(t *testing.T) {
	err := execErr(t, "CREATE TABLE ci (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(RS, RE)) WITH SYSTEM VERSIONING")
	if err != nil {
		t.Errorf("PERIOD(RS, RE) with columns rs/re should be accepted, got: %v", err)
	}
}

// TestSystemVersioningAlterColumnOps: a single-command DROP/MODIFY/CHANGE of a
// ROW START/END column on a versioned table must not silently leave PERIOD FOR
// SYSTEM_TIME pointing at a missing/non-row column. MariaDB rejects these.
func TestSystemVersioningAlterColumnOps(t *testing.T) {
	const create = "CREATE TABLE sv (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING"
	for _, alter := range []string{
		"ALTER TABLE sv DROP COLUMN rs",
		"ALTER TABLE sv MODIFY rs TIMESTAMP(6)",
		"ALTER TABLE sv CHANGE rs rs2 TIMESTAMP(6) GENERATED ALWAYS AS ROW START",
	} {
		t.Run(alter, func(t *testing.T) {
			c := versionedCatalog(t, create)
			r, _ := c.Exec(alter, &ExecOptions{ContinueOnError: true})
			if _, ok := r[0].Error.(*Error); !ok {
				t.Errorf("expected rejection (breaks PERIOD consistency), got accepted; ShowCreate:\n%s", c.ShowCreateTable("test", "sv"))
			}
		})
	}
}

// TestSystemVersioningRowColumnRules: ROW START/END columns must be unique and
// of type TIMESTAMP(6) (timestamp-precise) or BIGINT UNSIGNED (transaction-
// precise). Container-verified vs mariadb:11.8.8 (duplicate -> 4134, wrong type
// -> 4110).
func TestSystemVersioningRowColumnRules(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		// transaction-precise versioning uses BIGINT UNSIGNED row columns
		"CREATE TABLE t (x INT, rs BIGINT UNSIGNED GENERATED ALWAYS AS ROW START, re BIGINT UNSIGNED GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
	} {
		t.Run("accept", func(t *testing.T) {
			if err := execErr(t, sql); err != nil {
				t.Errorf("expected accepted for %q, got: %v", sql, err)
			}
		})
	}

	for _, sql := range []string{
		// duplicate ROW START
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, r2 TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		// duplicate ROW END
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, r2 TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		// row columns that are neither TIMESTAMP(6) nor BIGINT UNSIGNED
		"CREATE TABLE t (x INT, rs INT GENERATED ALWAYS AS ROW START, re INT GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		"CREATE TABLE t (x INT, rs TIMESTAMP GENERATED ALWAYS AS ROW START, re TIMESTAMP GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		"CREATE TABLE t (x INT, rs TIMESTAMP(3) GENERATED ALWAYS AS ROW START, re TIMESTAMP(3) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
	} {
		t.Run("reject", func(t *testing.T) {
			if _, ok := execErr(t, sql).(*Error); !ok {
				t.Errorf("expected rejection for %q, got accepted", sql)
			}
		})
	}
}

func versionedCatalog(t *testing.T, createSQL string) *Catalog {
	t.Helper()
	c := New()
	c.Exec("CREATE DATABASE test", nil)
	c.SetCurrentDatabase("test")
	c.Exec(createSQL, nil)
	return c
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
