package catalog

import (
	"context"
	"database/sql"
	"strings"
	"testing"
)

// TestOracle_OceanBaseViewAttributes holds to the engine what the catalog keeps of a view,
// where OceanBase's answer is not MySQL's: it keeps no ALGORITHM and no SQL SECURITY, reports
// every view updatable, widens no column of an outer join to nullable, and gives LPAD and
// RPAD the character set of both their strings. Each view is made on OceanBase and on a
// catalog, and both must answer what is written here.
func TestOracle_OceanBaseViewAttributes(t *testing.T) {
	ctx := context.Background()
	engine, c, db := engineDatabase(t, "oceanbase_catalog_view_attributes")
	for _, s := range []string{
		"CREATE TABLE t (id INT NOT NULL, a INT NOT NULL, n INT)",
		"CREATE TABLE u (id INT NOT NULL, b INT NOT NULL)",
		"CREATE TABLE s (l VARCHAR(10) CHARACTER SET latin1 COLLATE latin1_swedish_ci, m VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_bin)",
	} {
		if _, err := engine.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
		if code := omniCode(c, s); code != 0 {
			t.Fatalf("omni: %s: error %d", s, code)
		}
	}
	for _, tc := range []struct {
		name, view string
		updatable  string // VIEWS.IS_UPDATABLE
		security   string // VIEWS.SECURITY_TYPE
		columns    string // each column's IS_NULLABLE and CHARACTER_SET_NAME
	}{
		{"v_plain", "CREATE VIEW v_plain AS SELECT a, n FROM t", "YES", "NONE", "a NO -, n YES -"},
		{"v_distinct", "CREATE VIEW v_distinct AS SELECT DISTINCT a FROM t", "YES", "NONE", "a NO -"},
		{"v_temptable", "CREATE ALGORITHM=TEMPTABLE VIEW v_temptable AS SELECT a FROM t", "YES", "NONE", "a NO -"},
		{"v_merge", "CREATE ALGORITHM=MERGE VIEW v_merge AS SELECT a FROM t", "YES", "NONE", "a NO -"},
		{"v_group", "CREATE VIEW v_group AS SELECT a, COUNT(*) AS c FROM t GROUP BY a", "YES", "NONE", "a NO -, c YES -"},
		{"v_expr", "CREATE VIEW v_expr AS SELECT a + 1 AS e, 1 AS one, 'x' AS str, ABS(a) AS f, CAST(a AS SIGNED) AS g, id FROM t", "YES", "NONE", "e NO -, one NO -, str NO utf8mb4, f NO -, g NO -, id NO -"},
		{"v_inner", "CREATE VIEW v_inner AS SELECT t.a, u.b, t.n FROM t JOIN u ON t.id = u.id", "YES", "NONE", "a NO -, b NO -, n YES -"},
		// IF is stored as a CASE, and reads as one: nullable when a result is or there
		// is no ELSE, with the character set of its results.
		{"v_if", "CREATE VIEW v_if AS SELECT IF(a > 0, NULL, 'x') AS c, IF(a > 0, a, 1) AS e, CASE WHEN a > 0 THEN a END AS f FROM t", "YES", "NONE", "c YES utf8mb4, e NO -, f YES -"},
		{"v_if_charset", "CREATE VIEW v_if_charset AS SELECT IF(l IS NULL, 'y', l) AS d, CASE WHEN l IS NULL THEN m ELSE l END AS g FROM s", "YES", "NONE", "d YES latin1, g YES utf8mb4"},
		{"v_union", "CREATE VIEW v_union AS SELECT a FROM t UNION SELECT b FROM u", "YES", "NONE", "a NO -"},
		{"v_invoker", "CREATE SQL SECURITY INVOKER VIEW v_invoker AS SELECT a FROM t", "YES", "NONE", "a NO -"},
		{"v_definer", "CREATE SQL SECURITY DEFINER VIEW v_definer AS SELECT a FROM t", "YES", "NONE", "a NO -"},
		{"v_left", "CREATE VIEW v_left AS SELECT t.a, u.b FROM t LEFT JOIN u ON t.id = u.id", "YES", "NONE", "a NO -, b NO -"},
		{"v_right", "CREATE VIEW v_right AS SELECT t.a, u.b FROM t RIGHT JOIN u ON t.id = u.id", "YES", "NONE", "a NO -, b NO -"},
		{"v_repeat", "CREATE VIEW v_repeat AS SELECT REPEAT(l, 3) AS c FROM s", "YES", "NONE", "c YES latin1"},
		{"v_lpad_lm", "CREATE VIEW v_lpad_lm AS SELECT LPAD(l, 20, m) AS c FROM s", "YES", "NONE", "c YES utf8mb4"},
		{"v_lpad_ml", "CREATE VIEW v_lpad_ml AS SELECT LPAD(m, 20, l) AS c FROM s", "YES", "NONE", "c YES utf8mb4"},
		{"v_rpad_lm", "CREATE VIEW v_rpad_lm AS SELECT RPAD(l, 20, m) AS c FROM s", "YES", "NONE", "c YES utf8mb4"},
		{"v_rpad_ml", "CREATE VIEW v_rpad_ml AS SELECT RPAD(m, 20, l) AS c FROM s", "YES", "NONE", "c YES utf8mb4"},
		{"v_lpad_ll", "CREATE VIEW v_lpad_ll AS SELECT LPAD(l, 20, l) AS c FROM s", "YES", "NONE", "c YES latin1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := engine.ExecContext(ctx, tc.view); err != nil {
				t.Fatalf("%s: %v", tc.view, err)
			}
			if code := omniCode(c, tc.view); code != 0 {
				t.Fatalf("omni: %s: error %d", tc.view, code)
			}
			var updatable, security string
			if err := engine.QueryRowContext(ctx, `SELECT IS_UPDATABLE, SECURITY_TYPE FROM information_schema.VIEWS
				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`, tc.name).Scan(&updatable, &security); err != nil {
				t.Fatalf("VIEWS: %v", err)
			}
			rows, err := engine.QueryContext(ctx, `SELECT COLUMN_NAME, IS_NULLABLE, CHARACTER_SET_NAME FROM information_schema.COLUMNS
				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ? ORDER BY ORDINAL_POSITION`, tc.name)
			if err != nil {
				t.Fatalf("COLUMNS: %v", err)
			}
			defer rows.Close()
			var columns []string
			for rows.Next() {
				var name, nullable string
				var charset sql.NullString
				if err := rows.Scan(&name, &nullable, &charset); err != nil {
					t.Fatalf("COLUMNS: %v", err)
				}
				columns = append(columns, name+" "+nullable+" "+orDash(charset.String))
			}
			if err := rows.Err(); err != nil {
				t.Fatalf("COLUMNS: %v", err)
			}
			if got := strings.Join(columns, ", "); updatable != tc.updatable || security != tc.security || got != tc.columns {
				t.Errorf("OceanBase: updatable %s, security %s, columns %s\n      recorded: %s, %s, %s", updatable, security, got, tc.updatable, tc.security, tc.columns)
			}

			v := c.GetDatabase(db).Views[strings.ToLower(tc.name)]
			if v == nil {
				t.Fatal("the catalog has no such view")
			}
			own := make([]string, len(v.ColumnMetadata))
			for i, col := range v.ColumnMetadata {
				nullable := "NO"
				if col.Nullable {
					nullable = "YES"
				}
				own[i] = col.Name + " " + nullable + " " + orDash(strings.ToLower(col.Charset))
			}
			ownUpdatable := "NO"
			if v.IsUpdatable {
				ownUpdatable = "YES"
			}
			ownSecurity := "NONE"
			if v.SqlSecurity != "" {
				ownSecurity = strings.ToUpper(v.SqlSecurity)
			}
			if got := strings.Join(own, ", "); ownUpdatable != tc.updatable || ownSecurity != tc.security || got != tc.columns {
				t.Errorf("the catalog: updatable %s, security %s, columns %s\n     OceanBase: %s, %s, %s", ownUpdatable, ownSecurity, got, tc.updatable, tc.security, tc.columns)
			}
			if v.Algorithm != "" {
				t.Errorf("the catalog keeps ALGORITHM %s, which OceanBase does not keep", v.Algorithm)
			}
		})
	}
}

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

// TestOracle_OceanBaseViewColumnNames holds to the engine where an unaliased item's text
// stops being its column's name: OceanBase counts the 64 in bytes, where MySQL counts
// characters, and names the column Name_exp_N past it.
func TestOracle_OceanBaseViewColumnNames(t *testing.T) {
	ctx := context.Background()
	engine, c, db := engineDatabase(t, "oceanbase_catalog_view_column_names")
	bytes64 := "concat('" + strings.Repeat("é", 27) + "')" // 37 characters, 64 bytes
	bytes66 := "concat('" + strings.Repeat("é", 28) + "')" // 38 characters, 66 bytes
	chars64 := "concat('" + strings.Repeat("e", 54) + "')"
	chars65 := "concat('" + strings.Repeat("e", 55) + "')"
	for _, tc := range []struct{ name, item, column string }{
		{"v64b", bytes64, bytes64},
		{"v66b", bytes66, "Name_exp_1"},
		{"v64c", chars64, chars64},
		{"v65c", chars65, "Name_exp_1"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			view := "CREATE VIEW " + tc.name + " AS SELECT " + tc.item
			if _, err := engine.ExecContext(ctx, view); err != nil {
				t.Fatalf("%s: %v", view, err)
			}
			if code := omniCode(c, view); code != 0 {
				t.Fatalf("omni: %s: error %d", view, code)
			}
			var column string
			if err := engine.QueryRowContext(ctx, `SELECT COLUMN_NAME FROM information_schema.COLUMNS
				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?`, tc.name).Scan(&column); err != nil {
				t.Fatalf("COLUMNS: %v", err)
			}
			if column != tc.column {
				t.Errorf("OceanBase names the column %q; recorded %q", column, tc.column)
			}
			if own := c.GetDatabase(db).Views[tc.name].Columns; len(own) != 1 || own[0] != tc.column {
				t.Errorf("the catalog names the column %q; OceanBase %q", own, tc.column)
			}
		})
	}
}
