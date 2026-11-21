package catalog

import (
	"context"
	"database/sql"
	"strings"
	"testing"

	"github.com/bytebase/omni/metadata"
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
		// IF and CASE: nullable when a result is or a CASE has no ELSE, with the
		// character set of the results. OceanBase stores the IF as a CASE.
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

func TestOracle_OceanBaseViewOperatorTypesAndMVParenIndex(t *testing.T) {
	ctx := context.Background()
	engine, c, db := engineDatabase(t, "oceanbase_catalog_view_ops_and_mv")
	setup := []string{
		"CREATE TABLE t_ops (id INT NOT NULL PRIMARY KEY, name VARCHAR(20) NOT NULL, ch CHAR(5), txt TEXT, doc JSON, `a(b` INT)",
		`CREATE VIEW v_ops AS SELECT
			id & 1 AS c_band,
			id | 1 AS c_bor,
			id ^ 1 AS c_bxor,
			id << 1 AS c_shl,
			id >> 1 AS c_shr,
			~id AS c_bnot,
			id XOR 1 AS c_lxor,
			BINARY name AS c_bin_vc,
			BINARY ch AS c_bin_c,
			BINARY id AS c_bin_int,
			BINARY txt AS c_bin_txt,
			BINARY 'abc' AS c_bin_lit,
			doc->'$.x' AS c_jextract,
			doc->>'$.x' AS c_junquote
		FROM t_ops`,
	}
	for _, s := range setup {
		if _, err := engine.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
		if code := omniCode(c, s); code != 0 {
			t.Fatalf("omni: %s: error %d", s, code)
		}
	}
	for _, s := range []string{
		"CREATE MATERIALIZED VIEW mv_paren (PRIMARY KEY (id)) AS SELECT id, `a(b` FROM t_ops",
		"CREATE INDEX idx_paren ON mv_paren (`a(b`)",
	} {
		if _, err := engine.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	var mvQuery, idxCol string
	if err := engine.QueryRowContext(ctx, `SELECT QUERY FROM oceanbase.DBA_MVIEWS WHERE OWNER = DATABASE() AND MVIEW_NAME = 'mv_paren'`).Scan(&mvQuery); err != nil {
		t.Fatalf("DBA_MVIEWS: %v", err)
	}
	if err := engine.QueryRowContext(ctx, `SELECT s.COLUMN_NAME FROM oceanbase.DBA_MVIEWS m
		JOIN information_schema.STATISTICS s ON s.TABLE_SCHEMA = m.OWNER AND s.TABLE_NAME = m.CONTAINER_NAME
		WHERE m.OWNER = DATABASE() AND m.MVIEW_NAME = 'mv_paren' AND s.INDEX_NAME = 'idx_paren'`).Scan(&idxCol); err != nil {
		t.Fatalf("STATISTICS: %v", err)
	}
	if idxCol != "a(b" {
		t.Fatalf("expected OceanBase MV index column a(b, got %q", idxCol)
	}

	rows, err := engine.QueryContext(ctx, `SELECT COLUMN_NAME, COLUMN_TYPE FROM information_schema.COLUMNS
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'v_ops' ORDER BY ORDINAL_POSITION`)
	if err != nil {
		t.Fatalf("COLUMNS: %v", err)
	}
	defer rows.Close()
	engineTypes := map[string]string{}
	for rows.Next() {
		var name, colType string
		if err := rows.Scan(&name, &colType); err != nil {
			t.Fatal(err)
		}
		engineTypes[name] = colType
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	meta := c.DatabaseMetadata(db)
	catTypes := map[string]string{}
	for _, v := range meta.GetSchemas()[0].GetViews() {
		if v.GetName() == "v_ops" {
			for _, col := range v.GetColumns() {
				catTypes[col.GetName()] = col.GetType()
			}
		}
	}
	for colName, wantCat := range map[string]string{
		"c_band":     "bigint unsigned",
		"c_bor":      "bigint unsigned",
		"c_bxor":     "bigint unsigned",
		"c_shl":      "bigint unsigned",
		"c_shr":      "bigint unsigned",
		"c_bnot":     "bigint unsigned",
		"c_lxor":     "int",
		"c_bin_vc":   "varbinary(20)",
		"c_bin_c":    "varbinary(5)",
		"c_bin_int":  "varbinary(11)",
		"c_bin_txt":  "varbinary(65535)",
		"c_bin_lit":  "varbinary(3)",
		"c_jextract": "json",
		"c_junquote": "longtext",
	} {
		if catTypes[colName] != wantCat {
			t.Errorf("catalog v_ops.%s type = %q, want %q (engine reported %q)", colName, catTypes[colName], wantCat, engineTypes[colName])
		}
	}
	if engineTypes["c_jextract"] != "json" || engineTypes["c_junquote"] != "longtext" || engineTypes["c_bin_vc"] != "varbinary(20)" {
		t.Errorf("unexpected engine types: %v", engineTypes)
	}

	// Verify live MV index on quoted column with `(` survives LoadMetadata + DatabaseDefinition.
	meta.Schemas[0].MaterializedViews = append(meta.Schemas[0].MaterializedViews, &metadata.MaterializedViewMetadata{
		Name:       "mv_paren",
		Definition: mvQuery,
		Indexes: []*metadata.IndexMetadata{
			{Name: "PRIMARY", Primary: true, Unique: true, Visible: true, Type: "BTREE", Expressions: []string{"id"}, KeyLength: []int64{-1}, Descending: []bool{false}},
			{Name: "idx_paren", Visible: true, Type: "BTREE", Expressions: []string{idxCol}, KeyLength: []int64{-1}, Descending: []bool{false}},
		},
	})
	reloaded := New()
	report, err := reloaded.LoadMetadata(ctx, meta)
	if err != nil {
		t.Fatalf("LoadMetadata: %v", err)
	}
	if len(report.Incomplete) != 0 {
		t.Fatalf("unexpected incomplete items: %v", report.Incomplete)
	}
	reMV := reloaded.GetDatabase(db).Views["mv_paren"]
	if reMV == nil || len(reMV.Indexes) != 2 || len(reMV.Indexes[1].Columns) != 1 || reMV.Indexes[1].Columns[0].Name != "a(b" {
		t.Fatalf("unexpected reloaded MV indexes: %+v", reMV)
	}
	reMeta := reloaded.DatabaseMetadata(db)
	def, err := DatabaseDefinition(reMeta, DefinitionOptions{})
	if err != nil {
		t.Fatalf("DatabaseDefinition: %v", err)
	}
	if !strings.Contains(def, "CREATE INDEX `idx_paren` ON `mv_paren` (`a(b`)") {
		t.Fatalf("expected CREATE INDEX `idx_paren` ON `mv_paren` (`a(b`) in definition:\n%s", def)
	}
}
