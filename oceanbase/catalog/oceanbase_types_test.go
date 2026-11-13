package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// TestOracle_OceanBaseColumnTypes holds to the engine the type and character set the
// catalog gives a column, where OceanBase's is not MySQL 8.0's: it keeps an integer's
// display width and YEAR's, reads utf8 and the national types as utf8mb4, and takes a
// YEAR of any width. Each column is made alone in a table, on OceanBase and on a catalog,
// and both must answer the type, the character set and the error number written here.
func TestOracle_OceanBaseColumnTypes(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	const db = "obtypes"
	for _, s := range []string{"DROP DATABASE IF EXISTS " + db, "CREATE DATABASE " + db, "USE " + db} {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	for i, tc := range []struct {
		column    string // the column's definition
		typ       string // COLUMN_TYPE
		charset   string // CHARACTER_SET_NAME, "" for none
		errno     uint16
		collation string // COLLATION_NAME, where the column names one
	}{
		{"TINYINT", "tinyint(4)", "", taken, ""},
		{"SMALLINT", "smallint(6)", "", taken, ""},
		{"MEDIUMINT", "mediumint(9)", "", taken, ""},
		{"INT", "int(11)", "", taken, ""},
		{"INTEGER", "int(11)", "", taken, ""},
		{"BIGINT", "bigint(20)", "", taken, ""},
		{"TINYINT UNSIGNED", "tinyint(3) unsigned", "", taken, ""},
		{"SMALLINT UNSIGNED", "smallint(5) unsigned", "", taken, ""},
		{"MEDIUMINT UNSIGNED", "mediumint(8) unsigned", "", taken, ""},
		{"INT UNSIGNED", "int(10) unsigned", "", taken, ""},
		{"BIGINT UNSIGNED", "bigint(20) unsigned", "", taken, ""},
		{"INT(5)", "int(5)", "", taken, ""},
		{"INT(5) UNSIGNED", "int(5) unsigned", "", taken, ""},
		{"INT(5) ZEROFILL", "int(5) unsigned zerofill", "", taken, ""},
		{"INT ZEROFILL", "int(10) unsigned zerofill", "", taken, ""},
		{"BIGINT(3)", "bigint(3)", "", taken, ""},
		{"TINYINT(1)", "tinyint(1)", "", taken, ""},
		{"BOOLEAN", "tinyint(1)", "", taken, ""},
		{"BOOL", "tinyint(1)", "", taken, ""},
		{"SERIAL", "bigint(20) unsigned", "", taken, ""},
		{"YEAR", "year(4)", "", taken, ""},
		{"YEAR(4)", "year(4)", "", taken, ""},
		{"YEAR(2)", "year(4)", "", taken, ""},
		{"YEAR(3)", "year(4)", "", taken, ""},
		{"YEAR(5)", "year(4)", "", taken, ""},
		{"BIT", "bit(1)", "", taken, ""},
		{"DECIMAL", "decimal(10,0)", "", taken, ""},
		{"FLOAT", "float", "", taken, ""},
		{"DOUBLE", "double", "", taken, ""},
		{"VARCHAR(10) CHARACTER SET utf8", "varchar(10)", "utf8mb4", taken, ""},
		{"VARCHAR(10) CHARACTER SET utf8mb3", "varchar(10)", "utf8mb4", taken, ""},
		{"VARCHAR(10) CHARACTER SET utf8mb4", "varchar(10)", "utf8mb4", taken, ""},
		{"VARCHAR(10) COLLATE utf8_general_ci", "varchar(10)", "utf8mb4", taken, "utf8mb4_general_ci"},
		{"VARCHAR(10) COLLATE utf8mb3_bin", "varchar(10)", "utf8mb4", taken, "utf8mb4_bin"},
		{"VARCHAR(10) CHARACTER SET utf8 COLLATE utf8_bin", "varchar(10)", "utf8mb4", taken, "utf8mb4_bin"},
		{"VARCHAR(10) CHARACTER SET utf8mb4 COLLATE utf8_bin", "varchar(10)", "utf8mb4", taken, "utf8mb4_bin"},
		{"VARCHAR(10) CHARACTER SET latin1", "varchar(10)", "latin1", taken, ""},
		{"NATIONAL CHAR(10)", "char(10)", "utf8mb4", taken, ""},
		{"NCHAR(10)", "char(10)", "utf8mb4", taken, ""},
		{"NVARCHAR(10)", "varchar(10)", "utf8mb4", taken, ""},
		{"NATIONAL VARCHAR(10)", "varchar(10)", "utf8mb4", taken, ""},
		{"NCHAR VARCHAR(10)", "varchar(10)", "utf8mb4", taken, ""},
		{"NATIONAL CHAR VARYING(10)", "varchar(10)", "utf8mb4", taken, ""},
		{"VARCHAR(65535)", "varchar(65535)", "utf8mb4", taken, ""},
		{"VARCHAR(65536)", "varchar(65536)", "utf8mb4", taken, ""},
		{"VARCHAR(262144)", "varchar(262144)", "utf8mb4", taken, ""},
		{"VARCHAR(262145)", "", "", 1074, ""},
		{"CHAR(255)", "char(255)", "utf8mb4", taken, ""},
		{"CHAR(256)", "char(256)", "utf8mb4", taken, ""},
		{"CHAR(257)", "", "", 1074, ""},
		{"TEXT(10)", "tinytext", "utf8mb4", taken, ""},
		{"TEXT(1000)", "text", "utf8mb4", taken, ""},
		{"BLOB(10)", "tinyblob", "", taken, ""},
	} {
		t.Run(tc.column, func(t *testing.T) {
			table := fmt.Sprintf("t%02d", i)
			create := fmt.Sprintf("CREATE TABLE %s (c %s)", table, tc.column)
			c := New()
			if code := omniCode(c, "CREATE DATABASE "+db+"; USE "+db); code != 0 {
				t.Fatalf("omni: CREATE DATABASE: error %d", code)
			}
			_, err := o.db.ExecContext(ctx, create)
			var got uint16
			var myErr *gomysql.MySQLError
			switch {
			case errors.As(err, &myErr):
				got = myErr.Number
			case err != nil:
				t.Fatalf("%s: %v", create, err)
			}
			if got != tc.errno {
				t.Errorf("OceanBase answered %d (%v), recorded %d", got, err, tc.errno)
			}
			if code := omniCode(c, create); code != int(tc.errno) {
				t.Errorf("the catalog answered %d, OceanBase %d", code, tc.errno)
			}
			if got != taken {
				return
			}
			var typ string
			var charset, collation sql.NullString
			if err := o.db.QueryRowContext(ctx, `SELECT COLUMN_TYPE, CHARACTER_SET_NAME, COLLATION_NAME FROM information_schema.COLUMNS
				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND COLUMN_NAME = 'c'`, db, table).Scan(&typ, &charset, &collation); err != nil {
				t.Fatalf("COLUMNS: %v", err)
			}
			if typ != tc.typ || charset.String != tc.charset {
				t.Errorf("OceanBase stores %s, character set %q; recorded %s, %q", typ, charset.String, tc.typ, tc.charset)
			}
			if tc.collation != "" && collation.String != tc.collation {
				t.Errorf("OceanBase stores collation %q; recorded %q", collation.String, tc.collation)
			}
			tbl := c.GetDatabase(db).GetTable(table)
			if tbl == nil || tbl.GetColumn("c") == nil {
				return
			}
			col := tbl.GetColumn("c")
			if strings.ToLower(col.ColumnType) != tc.typ || strings.ToLower(col.Charset) != tc.charset {
				t.Errorf("the catalog stores %s, character set %q; OceanBase %s, %q", col.ColumnType, col.Charset, tc.typ, tc.charset)
			}
			if tc.collation != "" && strings.ToLower(col.Collation) != tc.collation {
				t.Errorf("the catalog stores collation %q; OceanBase %q", col.Collation, tc.collation)
			}
		})
	}
}

// TestOracle_OceanBaseIndexTypes holds to the engine the index type a USING clause leaves:
// OceanBase keeps HASH, which InnoDB makes BTREE, and takes USING BTREE on a SPATIAL key,
// which MySQL refuses. SHOW CREATE TABLE prints none of it (canonicalIndexType), so the type is
// read from information_schema.STATISTICS.
func TestOracle_OceanBaseIndexTypes(t *testing.T) {
	o := oceanBaseOracle(t)
	ctx := context.Background()
	const db = "obindextypes"
	c := New()
	for _, s := range []string{"DROP DATABASE IF EXISTS " + db, "CREATE DATABASE " + db, "USE " + db} {
		if _, err := o.db.ExecContext(ctx, s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
		if code := omniCode(c, s); code != 0 {
			t.Fatalf("omni: %s: error %d", s, code)
		}
	}
	for _, tc := range []struct {
		table, index string
		ddl          []string
		typ          string // INDEX_TYPE
		errno        uint16 // of the last statement
	}{
		{"t1", "k", []string{"CREATE TABLE t1 (a INT, KEY k (a))"}, "BTREE", taken},
		{"t2", "k", []string{"CREATE TABLE t2 (a INT, KEY k (a) USING BTREE)"}, "BTREE", taken},
		{"t3", "k", []string{"CREATE TABLE t3 (a INT, KEY k (a) USING HASH)"}, "HASH", taken},
		{"t4", "k", []string{"CREATE TABLE t4 (a INT, KEY k USING HASH (a))"}, "HASH", taken},
		{"t5", "k", []string{"CREATE TABLE t5 (a INT, UNIQUE KEY k (a) USING HASH) ENGINE=InnoDB"}, "HASH", taken},
		{"t6", "k", []string{"CREATE TABLE t6 (a INT)", "CREATE INDEX k ON t6 (a) USING HASH"}, "HASH", taken},
		{"t7", "k", []string{"CREATE TABLE t7 (a INT)", "ALTER TABLE t7 ADD KEY k (a) USING HASH"}, "HASH", taken},
		{"t8", "k", []string{"CREATE TABLE t8 (g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g))"}, "BTREE", taken},
		{"t9", "k", []string{"CREATE TABLE t9 (g GEOMETRY NOT NULL SRID 0, SPATIAL KEY k (g) USING BTREE)"}, "BTREE", taken},
		{"t11", "k", []string{"CREATE TABLE t11 (g GEOMETRY NOT NULL SRID 0)", "ALTER TABLE t11 ADD SPATIAL KEY k (g) USING BTREE"}, "BTREE", taken},
	} {
		t.Run(tc.table, func(t *testing.T) {
			var got uint16
			for i, stmt := range tc.ddl {
				_, err := o.db.ExecContext(ctx, stmt)
				var myErr *gomysql.MySQLError
				switch {
				case errors.As(err, &myErr):
					got = myErr.Number
				case err != nil:
					t.Fatalf("%s: %v", stmt, err)
				}
				want := uint16(taken)
				if i == len(tc.ddl)-1 {
					want = tc.errno
				}
				if got != want {
					t.Errorf("%s\n  OceanBase answered %d (%v), recorded %d", stmt, got, err, want)
				}
				if code := omniCode(c, stmt); code != int(want) {
					t.Errorf("%s\n  the catalog answered %d, OceanBase %d", stmt, code, want)
				}
			}
			if got != taken {
				return
			}
			var typ string
			if err := o.db.QueryRowContext(ctx, `SELECT DISTINCT INDEX_TYPE FROM information_schema.STATISTICS
				WHERE TABLE_SCHEMA = ? AND TABLE_NAME = ? AND INDEX_NAME = ?`, db, tc.table, tc.index).Scan(&typ); err != nil {
				t.Fatalf("STATISTICS: %v", err)
			}
			if typ != tc.typ {
				t.Errorf("OceanBase stores %s; recorded %s", typ, tc.typ)
			}
			tbl := c.GetDatabase(db).GetTable(tc.table)
			if tbl == nil {
				return
			}
			for _, idx := range tbl.Indexes {
				if !strings.EqualFold(idx.Name, tc.index) {
					continue
				}
				// The catalog writes "" for an index without USING and SPATIAL for a
				// spatial one, both of which the engine reports as BTREE.
				own := strings.ToUpper(idx.IndexType)
				if own == "" || own == "SPATIAL" {
					own = "BTREE"
				}
				if own != tc.typ {
					t.Errorf("the catalog stores %s (IndexType %q); OceanBase %s", own, idx.IndexType, tc.typ)
				}
			}
		})
	}
}
