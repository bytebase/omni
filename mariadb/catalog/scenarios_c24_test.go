package catalog

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestScenario_C24 records that MariaDB has no generated invisible primary key.
//
// MySQL 8.0.30+ has one: with `sql_generate_invisible_primary_key=ON` a CREATE
// TABLE that declares no PRIMARY KEY silently gains a hidden
// `my_row_id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT INVISIBLE` column and a
// PRIMARY KEY over it, shown only when
// `show_gipk_in_create_table_and_information_schema=ON`. The scenario this
// file was forked from covered that behavior.
//
// MariaDB has neither variable and generates no such column, so the catalog
// models neither. Measured on mariadb:11.8.8: both SETs answer 1193, and a
// table with no primary key has none.
func TestScenario_C24(t *testing.T) {
	mc, cleanup := scenarioContainer(t)
	defer cleanup()
	scenarioReset(t, mc)

	t.Run("24.1_gipk_variables_do_not_exist", func(t *testing.T) {
		for _, stmt := range []string{
			"SET SESSION sql_generate_invisible_primary_key = ON",
			"SET SESSION show_gipk_in_create_table_and_information_schema = ON",
		} {
			_, err := mc.db.ExecContext(mc.ctx, stmt)
			if err == nil {
				t.Fatalf("%q: expected MariaDB to refuse an unknown system variable", stmt)
			}
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("%q: MariaDB gave no verdict: %v", stmt, err)
			}
			if myErr.Number != 1193 {
				t.Errorf("%q: got error %d (%s), want 1193 ER_UNKNOWN_SYSTEM_VARIABLE",
					stmt, myErr.Number, myErr.Message)
			}
		}
	})

	t.Run("24.2_no_hidden_primary_key_is_generated", func(t *testing.T) {
		const ddl = "CREATE TABLE t_no_pk (a INT)"
		if err := mc.execSQL(ddl); err != nil {
			t.Fatalf("container exec %q: %v", ddl, err)
		}
		oracleDDL, err := mc.showCreateTable("t_no_pk")
		if err != nil {
			t.Fatalf("SHOW CREATE TABLE: %v", err)
		}
		if strings.Contains(oracleDDL, "my_row_id") || strings.Contains(oracleDDL, "PRIMARY KEY") {
			t.Errorf("MariaDB generated a primary key for a table that declares none:\n%s", oracleDDL)
		}

		c := scenarioNewCatalog(t)
		if _, err := c.Exec(ddl, nil); err != nil {
			t.Fatalf("catalog exec %q: %v", ddl, err)
		}
		omniDDL := c.ShowCreateTable("testdb", "t_no_pk")
		if strings.Contains(omniDDL, "my_row_id") || strings.Contains(omniDDL, "PRIMARY KEY") {
			t.Errorf("omni generated a primary key for a table that declares none:\n%s", omniDDL)
		}
	})
}
