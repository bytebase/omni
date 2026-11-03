package catalog

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/metadata"
	gomysql "github.com/go-sql-driver/mysql"
)

// invisiblePrimaryKeyCase is a CREATE TABLE of t whose primary key, explicit or
// implicit, may be invisible. mysql is MySQL 8.0's answer, 3522 or 0 for a
// created table; catalog is the catalog's, which differs only where noted.
type invisiblePrimaryKeyCase struct {
	name    string
	sql     string
	mysql   int
	catalog int
}

// invisiblePrimaryKeyCases are read off MySQL 8.0, which TestContainer_InvisiblePrimaryKey
// asks again. Without a primary key, MySQL takes the first UNIQUE key on NOT NULL
// parts, in the order the statement declares the keys, for the primary key, and
// refuses the table when that key is invisible.
var invisiblePrimaryKeyCases = []invisiblePrimaryKeyCase{
	{"explicit", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE)", 3522, 3522},
	{"explicit_named", "CREATE TABLE t (a INT, CONSTRAINT pk PRIMARY KEY (a) COMMENT 'x' INVISIBLE)", 3522, 3522},
	{"explicit_last_option_wins", "CREATE TABLE t (a INT, PRIMARY KEY (a) INVISIBLE VISIBLE)", 0, 0},

	{"sole_unique", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)", 3522, 3522},
	{"constraint_unique", "CREATE TABLE t (a INT NOT NULL, CONSTRAINT u UNIQUE (a) INVISIBLE)", 3522, 3522},
	{"unique_using_hash", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) USING HASH INVISIBLE)", 3522, 3522},
	{"composite", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a, b) INVISIBLE)", 3522, 3522},
	{"last_option_wins_invisible", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) VISIBLE INVISIBLE)", 3522, 3522},
	{"last_option_wins_visible", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE VISIBLE)", 0, 0},
	{"first_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, UNIQUE KEY v (b))", 3522, 3522},
	{"second_of_two", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", 0, 0},
	{"plain_key_first", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (b), UNIQUE KEY u (a) INVISIBLE)", 3522, 3522},
	{"invisible_plain_key", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, KEY k (a) INVISIBLE, UNIQUE KEY v (b))", 0, 0},

	// A key with a nullable part is never the primary key.
	{"nullable", "CREATE TABLE t (a INT, UNIQUE KEY u (a) INVISIBLE)", 0, 0},
	{"partly_nullable", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY u (a, b) INVISIBLE)", 0, 0},
	{"nullable_key_first", "CREATE TABLE t (a INT NOT NULL, b INT, UNIQUE KEY v (b), UNIQUE KEY u (a) INVISIBLE)", 3522, 3522},

	// An explicit primary key anywhere in the statement is the primary key.
	{"explicit_key", "CREATE TABLE t (a INT NOT NULL, b INT, PRIMARY KEY (b), UNIQUE KEY u (a) INVISIBLE)", 0, 0},
	{"explicit_key_last", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, PRIMARY KEY (b))", 0, 0},
	{"column_primary_key", "CREATE TABLE t (a INT PRIMARY KEY, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 0, 0},

	// A column's UNIQUE or SERIAL declares its key where the column stands.
	{"column_unique_first", "CREATE TABLE t (a INT NOT NULL UNIQUE, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 0, 0},
	{"serial_first", "CREATE TABLE t (a SERIAL, b INT NOT NULL, UNIQUE KEY u (b) INVISIBLE)", 0, 0},
	{"column_unique_after", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE, b INT NOT NULL UNIQUE)", 3522, 3522},
	{"executable_comment", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, /*!80000 UNIQUE KEY u (a) INVISIBLE, */ c INT NOT NULL UNIQUE)", 3522, 3522},

	// A prefix part counts, and statement order holds even though SHOW CREATE TABLE
	// lists a prefix key after the others.
	{"prefix", "CREATE TABLE t (a VARCHAR(10) NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE)", 3522, 3522},
	{"prefix_first_visible", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)), UNIQUE KEY v (b) INVISIBLE)", 0, 0},
	{"prefix_first_invisible", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY u (a(3)) INVISIBLE, UNIQUE KEY v (b))", 3522, 3522},
	{"blob_prefix", "CREATE TABLE t (a INT NOT NULL, b BLOB NOT NULL, UNIQUE KEY v (b(4)) INVISIBLE)", 3522, 3522},
	// MySQL's own SHOW CREATE TABLE of prefix_first_visible, which it refuses.
	{"prefix_readback", "CREATE TABLE t (a VARCHAR(10) NOT NULL, b INT NOT NULL, UNIQUE KEY v (b) INVISIBLE, UNIQUE KEY u (a(3)))", 3522, 3522},

	// NOT NULL however it comes.
	{"auto_increment", "CREATE TABLE t (a INT AUTO_INCREMENT, UNIQUE KEY u (a) INVISIBLE)", 3522, 3522},
	{"stored_generated", "CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) STORED NOT NULL, UNIQUE KEY u (g) INVISIBLE)", 3522, 3522},
	{"virtual_generated", "CREATE TABLE t (a INT NOT NULL, g INT AS (a + 1) VIRTUAL NOT NULL, UNIQUE KEY u (g) INVISIBLE)", 3522, 3522},
	{"create_select", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE) SELECT 1 AS a", 3522, 3522},

	// An expression part is NOT NULL when the expression is, which the catalog
	// does not work out, so it stops at such a key without refusing.
	{"expression_first_visible", "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)", 0, 0},
	{"expression_nullable", "CREATE TABLE t (a INT, UNIQUE KEY u ((a + 1)) INVISIBLE)", 0, 0},
	{"expression_not_null", "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u ((a + 1)) INVISIBLE)", 3522, 0},
	{"expression_nullable_first", "CREATE TABLE t (a INT, b INT NOT NULL, UNIQUE KEY f ((a + 1)), UNIQUE KEY v (b) INVISIBLE)", 3522, 0},
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
			if got := execCode(t, c, tc.sql); got != tc.catalog {
				t.Fatalf("got %d, want %d", got, tc.catalog)
			}
			if created := c.GetDatabase("testdb").GetTable("t") != nil; created != (tc.catalog == 0) {
				t.Errorf("table created: %v", created)
			}
		})
	}
}

// TestCreateTableGeneratedPrimaryKeyTakesInvisibleUnique pins that a generated
// invisible primary key is the primary key: MySQL 8.0 creates the table, with
// my_row_id for the key and u still invisible.
func TestCreateTableGeneratedPrimaryKeyTakesInvisibleUnique(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "SET SESSION sql_generate_invisible_primary_key = ON")
	if got := execCode(t, c, "CREATE TABLE t (a INT NOT NULL, UNIQUE KEY u (a) INVISIBLE)"); got != 0 {
		t.Fatalf("got %d, want the table", got)
	}
	if got := execCode(t, c, "CREATE TABLE t2 (a INT NOT NULL, PRIMARY KEY (a) INVISIBLE)"); got != 3522 {
		t.Fatalf("got %d, want 3522", got)
	}
}

// TestLoadMetadataKeepsTableWithInvisibleUniqueKey pins that the loader installs
// the table a server holds whatever order the snapshot lists its keys in. MySQL
// created it with z_visible declared first, which made that key the implicit
// primary key; a snapshot listing the keys by name puts a_hidden first, and
// checking the implicit primary key in that order would refuse a table that
// exists.
func TestLoadMetadataKeepsTableWithInvisibleUniqueKey(t *testing.T) {
	c, report := loadMySQLSnapshot(t, snapshotWithTables(&metadata.TableMetadata{
		Name:    "t",
		Columns: []*metadata.ColumnMetadata{{Name: "a", Type: "int"}, {Name: "b", Type: "int"}},
		Indexes: []*metadata.IndexMetadata{
			{Name: "a_hidden", Type: "BTREE", Unique: true, Expressions: []string{"a"}},
			{Name: "z_visible", Type: "BTREE", Unique: true, Visible: true, Expressions: []string{"b"}},
		},
	}))
	requireReport(t, report, nil, nil)
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

// TestContainer_InvisiblePrimaryKey asks MySQL 8.0 about every case, and for a
// table it creates, requires the same indexes with the same visibility as the
// catalog's, read from information_schema.STATISTICS. SHOW CREATE TABLE would not
// do: it lists a prefix or expression key after the others, and MySQL refuses
// that text for the cases ending in _first_visible ("prefix_readback")
// (container test).
func TestContainer_InvisiblePrimaryKey(t *testing.T) {
	ctr, cleanup := startContainer(t)
	defer cleanup()

	for _, tc := range invisiblePrimaryKeyCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := ctr.execSQLDirect("DROP TABLE IF EXISTS t"); err != nil {
				t.Fatalf("drop: %v", err)
			}
			got := 0
			if err := ctr.execSQLDirect(tc.sql); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatalf("MySQL gave no verdict: %v", err)
				}
				got = int(myErr.Number)
			}
			if got != tc.mysql {
				t.Fatalf("MySQL answered %d, want %d", got, tc.mysql)
			}
			if got != 0 || tc.catalog != 0 {
				return
			}
			rows, err := ctr.db.QueryContext(ctr.ctx, `SELECT DISTINCT INDEX_NAME, IS_VISIBLE FROM information_schema.STATISTICS
				WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 't'`)
			if err != nil {
				t.Fatalf("read STATISTICS: %v", err)
			}
			defer rows.Close()
			want := map[string]bool{}
			for rows.Next() {
				var name, visible string
				if err := rows.Scan(&name, &visible); err != nil {
					t.Fatalf("scan STATISTICS: %v", err)
				}
				want[name] = visible == "YES"
			}
			if err := rows.Err(); err != nil {
				t.Fatalf("read STATISTICS: %v", err)
			}
			c := setupWithDB(t)
			mustExec(t, c, tc.sql)
			have := indexVisibility(c, "testdb")
			if len(want) != len(have) {
				t.Fatalf("catalog indexes %v, MySQL has %v", have, want)
			}
			for name, visible := range want {
				if v, ok := have[name]; !ok || v != visible {
					t.Errorf("index %s: catalog visible=%v (present %v), MySQL visible=%v", name, v, ok, visible)
				}
			}
		})
	}
}
