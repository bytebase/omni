package catalog

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/metadata"
	gomysql "github.com/go-sql-driver/mysql"
)

// invisiblePrimaryKeyCases are CREATE TABLE statements of t whose primary key,
// explicit or implicit, may be invisible, with TiDB v8.5.5's answer: 3522, or 0
// for a created table. TestContainerCreateInvisiblePrimaryKey asks TiDB again.
//
// Without a primary key TiDB takes the first UNIQUE index whose columns are all
// NOT NULL for the primary key, passing over an expression index, and it orders
// the indexes of the table-level definitions before those a column definition
// declares: "column_unique_first" fails here and succeeds on MySQL 8.0.
var invisiblePrimaryKeyCases = []struct {
	name string
	sql  string
	code int
}{
	{"explicit", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE)", 3522},
	{"explicit_named", "CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a) INVISIBLE)", 3522},
	{"explicit_comment", "CREATE TABLE t (a INT, PRIMARY KEY (a) COMMENT 'x' INVISIBLE)", 3522},
	{"explicit_using_btree", "CREATE TABLE t (a INT, PRIMARY KEY USING BTREE (a) INVISIBLE)", 3522},
	{"explicit_clustered", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE CLUSTERED)", 3522},
	{"explicit_nonclustered", "CREATE TABLE t (v VARCHAR(10), PRIMARY KEY (v) INVISIBLE COMMENT 'x' NONCLUSTERED)", 3522},
	{"explicit_composite", "CREATE TABLE t (a INT, b INT, PRIMARY KEY (a, b) INVISIBLE)", 3522},
	{"explicit_auto_increment", "CREATE TABLE t (a INT AUTO_INCREMENT, PRIMARY KEY (a) INVISIBLE)", 3522},
	{"explicit_with_unique", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE, UNIQUE KEY u (a))", 3522},
	{"explicit_last_option_invisible", "CREATE TABLE t (a INT, PRIMARY KEY (a) VISIBLE INVISIBLE)", 3522},
	{"explicit_last_option_visible", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE)", 0},
	{"explicit_visible", "CREATE TABLE t (a INT, PRIMARY KEY (a) VISIBLE)", 0},

	{"sole_unique", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)", 3522},
	{"constraint_unique", "CREATE TABLE t (a INT NOT NULL, CONSTRAINT u UNIQUE (a) INVISIBLE)", 3522},
	{"unique_using_hash", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH INVISIBLE)", 3522},
	{"composite", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a, b) INVISIBLE)", 3522},
	{"last_option_wins_visible", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE VISIBLE)", 0},
	{"first_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b))", 3522},
	{"second_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", 0},
	{"plain_key_first", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (b), UNIQUE KEY u (a) INVISIBLE)", 3522},
	{"invisible_plain_key", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (a) INVISIBLE, UNIQUE KEY v (b))", 0},

	// An index with a nullable column is never the primary key.
	{"nullable", "CREATE TABLE t (a INT, UNIQUE KEY u (a) INVISIBLE)", 0},
	{"partly_nullable", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) INVISIBLE)", 0},
	{"nullable_key_first", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", 3522},

	// An explicit primary key anywhere in the statement is the primary key,
	// clustered or not.
	{"explicit_key", "CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (b), UNIQUE KEY u (a) INVISIBLE)", 0},
	{"explicit_key_last", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, PRIMARY KEY (b))", 0},
	{"column_primary_key", "CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 0},
	{"composite_key_last", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b), PRIMARY KEY (a, b))", 0},

	// A column definition's UNIQUE or SERIAL comes after every table-level index.
	{"column_unique_first", "CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 3522},
	{"column_unique_after", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, b INT NOT NULL UNIQUE)", 3522},
	{"serial_first", "CREATE TABLE t (a SERIAL, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 3522},
	{"column_unique_nullable_table_key", "CREATE TABLE t (a INT, b INT NOT NULL UNIQUE KEY, UNIQUE KEY u (a) INVISIBLE)", 0},
	{"table_keys_in_order", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a), b INT NOT NULL UNIQUE, UNIQUE KEY v (b) INVISIBLE)", 0},

	// A prefix column counts, in statement order.
	{"prefix", "CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE)", 3522},
	{"prefix_first_visible", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)), UNIQUE KEY v (b) INVISIBLE)", 0},
	{"prefix_first_invisible", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE, UNIQUE KEY v (b))", 3522},
	{"text_prefix_first", "CREATE TABLE t (a TEXT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(10)), UNIQUE KEY v (b) INVISIBLE)", 0},

	// NOT NULL however it comes.
	{"auto_increment", "CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) INVISIBLE)", 3522},
	{"stored_generated", "CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) STORED NOT NULL, UNIQUE KEY u (g) INVISIBLE)", 3522},
	{"virtual_generated", "CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) VIRTUAL NOT NULL, UNIQUE KEY u (g) INVISIBLE)", 3522},

	// An expression index is never the primary key.
	{"expression", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u ((a + 1)) INVISIBLE)", 0},
	{"expression_first", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)", 3522},
}

func execCode(t *testing.T, c *Catalog, sql string) int {
	t.Helper()
	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("parse %q: %v", sql, err)
	}
	if len(results) != 1 {
		t.Fatalf("%q: %d results", sql, len(results))
	}
	if results[0].Error == nil {
		return 0
	}
	var catErr *Error
	if !errors.As(results[0].Error, &catErr) {
		t.Fatalf("%q: %v", sql, results[0].Error)
	}
	return catErr.Code
}

func TestCreateTableInvisiblePrimaryKey(t *testing.T) {
	for _, tc := range invisiblePrimaryKeyCases {
		t.Run(tc.name, func(t *testing.T) {
			c := setupWithDB(t)
			if got := execCode(t, c, tc.sql); got != tc.code {
				t.Fatalf("got %d, want %d", got, tc.code)
			}
			if created := c.GetDatabase("testdb").GetTable("t") != nil; created != (tc.code == 0) {
				t.Errorf("table created: %v", created)
			}
		})
	}
}

// TestLoadMetadataKeepsTableWithInvisibleUniqueKey pins that the loader installs
// the table a server holds whatever order the snapshot lists its indexes in.
// TiDB created it with z_visible first, which made that index the implicit
// primary key; a snapshot listing the indexes by name puts a_hidden first, and
// checking the implicit primary key in that order would refuse a table that
// exists.
func TestLoadMetadataKeepsTableWithInvisibleUniqueKey(t *testing.T) {
	c, report := loadTiDBSnapshot(t, snapshot(&metadata.SchemaMetadata{Tables: []*metadata.TableMetadata{{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}},
		Indexes: []*metadata.IndexMetadata{
			{Name: "a_hidden", Type: "BTREE", Unique: true, Expressions: []string{"a"}},
			{Name: "z_visible", Type: "BTREE", Unique: true, Visible: true, Expressions: []string{"b"}},
		},
	}}}))
	requireReport(t, report)
	visible := map[string]bool{}
	for _, idx := range requireTable(t, c, "t").Indexes {
		visible[idx.Name] = idx.Visible
	}
	if len(visible) != 2 || visible["a_hidden"] || !visible["z_visible"] {
		t.Errorf("indexes %v, want a_hidden invisible and z_visible visible", visible)
	}
}

// indexVisibility maps each index of table t to whether it is visible.
func indexVisibility(c *Catalog, db string) map[string]bool {
	got := map[string]bool{}
	for _, idx := range c.GetDatabase(db).GetTable("t").Indexes {
		got[idx.Name] = idx.Visible
	}
	return got
}

// TestContainerCreateInvisiblePrimaryKey asks TiDB v8.5.5 about every case, and
// for a table it creates, loads its SHOW CREATE TABLE into a fresh catalog and
// requires the same indexes with the same visibility as the catalog's
// (container test).
func TestContainerCreateInvisiblePrimaryKey(t *testing.T) {
	tc := startTiDBForCatalog(t)
	for _, stmt := range []string{
		"DROP DATABASE IF EXISTS omni_invisible_pk",
		"CREATE DATABASE omni_invisible_pk",
		"USE omni_invisible_pk",
	} {
		mustExecTiDB(t, tc, stmt)
	}
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_invisible_pk")

	for _, cs := range invisiblePrimaryKeyCases {
		t.Run(cs.name, func(t *testing.T) {
			mustExecTiDB(t, tc, "DROP TABLE IF EXISTS t")
			got := 0
			if _, err := tc.db.ExecContext(tc.ctx, cs.sql); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatalf("TiDB gave no verdict: %v", err)
				}
				got = int(myErr.Number)
			}
			if got != cs.code {
				t.Fatalf("TiDB answered %d, want %d", got, cs.code)
			}
			if got != 0 {
				return
			}
			var name, readback string
			if err := tc.db.QueryRowContext(tc.ctx, "SHOW CREATE TABLE t").Scan(&name, &readback); err != nil {
				t.Fatalf("SHOW CREATE TABLE: %v", err)
			}
			engine := setupWithDB(t)
			mustExec(t, engine, readback)
			c := setupWithDB(t)
			mustExec(t, c, cs.sql)
			want, have := indexVisibility(engine, "testdb"), indexVisibility(c, "testdb")
			if len(want) != len(have) {
				t.Fatalf("catalog indexes %v, TiDB reads back %v:\n%s", have, want, readback)
			}
			for idx, visible := range want {
				if v, ok := have[idx]; !ok || v != visible {
					t.Errorf("index %s: catalog visible=%v (present %v), TiDB visible=%v", idx, v, ok, visible)
				}
			}
		})
	}
}
