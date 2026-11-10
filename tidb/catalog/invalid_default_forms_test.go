package catalog

import (
	"errors"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// invalidDefaultForm is one DDL statement with the verdict TiDB v8.5.5 gives
// it: wantCode is the error number, or 0 when the engine accepts the statement.
// Every case was measured against a live pingcap/tidb:v8.5.5 container on
// 2026-11-10, and TestContainerInvalidDefaultForms asks the engine again.
//
// Each case owns table `t`: setup (when present) creates it, and both tests
// drop it before the case runs, so the cases do not see each other.
type invalidDefaultForm struct {
	name     string
	setup    string
	sql      string
	wantCode int
}

// invalidDefaultForms are the forms TiDB v8.5.5 refuses and this catalog used to
// accept silently, plus the neighbouring forms it accepts, so the rules cannot
// drift into over-rejection. The five clusters are: an empty-string default on a
// type that cannot hold one (1067), a NULL default on a column that is not
// nullable (1067, or 1171 when the column is part of a primary key), more than
// one AUTO_INCREMENT column (1075), and a key part naming a column the table
// does not have (1072).
var invalidDefaultForms = []invalidDefaultForm{
	// Cluster 1: empty-string default. TiDB refuses '' (and a whitespace-only
	// string) on the temporal and numeric types and on ENUM, and takes it on
	// the string, binary and BIT types and on SET.
	{name: "date default empty", sql: "CREATE TABLE t (a date DEFAULT '')", wantCode: 1067},
	{name: "datetime default empty", sql: "CREATE TABLE t (a datetime DEFAULT '')", wantCode: 1067},
	{name: "datetime fsp default empty", sql: "CREATE TABLE t (a datetime(3) DEFAULT '')", wantCode: 1067},
	{name: "timestamp default empty", sql: "CREATE TABLE t (a timestamp DEFAULT '')", wantCode: 1067},
	{name: "year default empty", sql: "CREATE TABLE t (a year DEFAULT '')", wantCode: 1067},
	{name: "int default empty", sql: "CREATE TABLE t (a int DEFAULT '')", wantCode: 1067},
	{name: "tinyint default empty", sql: "CREATE TABLE t (a tinyint DEFAULT '')", wantCode: 1067},
	{name: "bigint default empty", sql: "CREATE TABLE t (a bigint DEFAULT '')", wantCode: 1067},
	{name: "decimal default empty", sql: "CREATE TABLE t (a decimal(10,2) DEFAULT '')", wantCode: 1067},
	{name: "double default empty", sql: "CREATE TABLE t (a double DEFAULT '')", wantCode: 1067},
	{name: "bool default empty", sql: "CREATE TABLE t (a bool DEFAULT '')", wantCode: 1067},
	{name: "enum default empty", sql: "CREATE TABLE t (a enum('x','y') DEFAULT '')", wantCode: 1067},
	// An ENUM is judged on its declared members, not on emptiness: a blank
	// string that the ENUM declares is a valid default, and one it does not is
	// 1067. Measured on 2026-11-10.
	{name: "enum declaring empty member default empty accepted", sql: "CREATE TABLE t (a enum('','x') DEFAULT '')"},
	{name: "enum declaring empty member last default empty accepted", sql: "CREATE TABLE t (a enum('a','') DEFAULT '')"},
	{name: "enum declaring empty member not null default empty accepted", sql: "CREATE TABLE t (a enum('','x') NOT NULL DEFAULT '')"},
	{name: "enum declaring blank member default blank accepted", sql: "CREATE TABLE t (a enum(' ','x') DEFAULT ' ')"},
	{name: "enum default blank not a member", sql: "CREATE TABLE t (a enum('x') DEFAULT ' ')", wantCode: 1067},
	{name: "date default blank", sql: "CREATE TABLE t (a date DEFAULT ' ')", wantCode: 1067},
	{name: "date not null default empty", sql: "CREATE TABLE t (a date NOT NULL DEFAULT '')", wantCode: 1067},
	{name: "varchar default empty accepted", sql: "CREATE TABLE t (a varchar(10) DEFAULT '')"},
	{name: "char default empty accepted", sql: "CREATE TABLE t (a char(5) DEFAULT '')"},
	{name: "set default empty accepted", sql: "CREATE TABLE t (a set('x','y') DEFAULT '')"},
	{name: "binary default empty accepted", sql: "CREATE TABLE t (a binary(4) DEFAULT '')"},
	{name: "bit default empty accepted", sql: "CREATE TABLE t (a bit(4) DEFAULT '')"},
	{name: "date default literal accepted", sql: "CREATE TABLE t (a date DEFAULT '2020-01-01')"},
	{name: "int default numeric string accepted", sql: "CREATE TABLE t (a int DEFAULT '5')"},
	{name: "decimal default numeric string accepted", sql: "CREATE TABLE t (a decimal DEFAULT '1.5')"},
	{name: "enum default member accepted", sql: "CREATE TABLE t (a enum('x') DEFAULT 'x')"},

	// Cluster 2: ALTER TABLE ... ALTER COLUMN c SET DEFAULT. A primary-key
	// column answers 1171; any other non-nullable column answers 1067.
	{name: "set default null on not null", setup: "CREATE TABLE t (a INT NOT NULL)",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", wantCode: 1067},
	{name: "set default null on pk", setup: "CREATE TABLE t (a INT PRIMARY KEY)",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", wantCode: 1171},
	{name: "set default null on composite pk part", setup: "CREATE TABLE t (a INT NOT NULL, b INT NOT NULL, PRIMARY KEY (a,b))",
		sql: "ALTER TABLE t ALTER COLUMN b SET DEFAULT NULL", wantCode: 1171},
	{name: "set default null on not null unique", setup: "CREATE TABLE t (a INT NOT NULL UNIQUE)",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL", wantCode: 1067},
	{name: "set default empty on date", setup: "CREATE TABLE t (x date)",
		sql: "ALTER TABLE t ALTER COLUMN x SET DEFAULT ''", wantCode: 1067},
	{name: "set default empty on int", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ALTER COLUMN x SET DEFAULT ''", wantCode: 1067},
	{name: "set default null on nullable accepted", setup: "CREATE TABLE t (a INT)",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT NULL"},
	{name: "set default literal on not null accepted", setup: "CREATE TABLE t (a INT NOT NULL)",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT 5"},
	{name: "set default empty on varchar accepted", setup: "CREATE TABLE t (x varchar(5))",
		sql: "ALTER TABLE t ALTER COLUMN x SET DEFAULT ''"},
	{name: "set default empty on enum declaring it accepted", setup: "CREATE TABLE t (a enum('','x'))",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT ''"},
	{name: "set default empty on enum not declaring it", setup: "CREATE TABLE t (a enum('x','y'))",
		sql: "ALTER TABLE t ALTER COLUMN a SET DEFAULT ''", wantCode: 1067},

	// Cluster 3: NOT NULL with DEFAULT NULL, written in the column definition.
	// A table-level PRIMARY KEY over the column answers 1171 instead, which is
	// the order TiDB reads the two rules in.
	{name: "not null default null", sql: "CREATE TABLE t (a INT NOT NULL DEFAULT NULL)", wantCode: 1067},
	{name: "inline pk default null", sql: "CREATE TABLE t (a INT PRIMARY KEY DEFAULT NULL)", wantCode: 1067},
	{name: "default null then not null", sql: "CREATE TABLE t (a INT DEFAULT NULL NOT NULL)", wantCode: 1067},
	{name: "table pk default null", sql: "CREATE TABLE t (a INT DEFAULT NULL, PRIMARY KEY(a))", wantCode: 1171},
	{name: "table pk not null default null", sql: "CREATE TABLE t (a INT NOT NULL DEFAULT NULL, PRIMARY KEY(a))", wantCode: 1171},
	{name: "table pk auto increment default null", sql: "CREATE TABLE t (a INT auto_increment DEFAULT NULL, PRIMARY KEY(a))", wantCode: 1171},
	{name: "nullable default null accepted", sql: "CREATE TABLE t (a INT DEFAULT NULL)"},
	{name: "unique default null accepted", sql: "CREATE TABLE t (a INT DEFAULT NULL UNIQUE)"},
	{name: "not null then null default null accepted", sql: "CREATE TABLE t (a INT NOT NULL NULL DEFAULT NULL)"},
	{name: "add column not null default null", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ADD COLUMN b INT NOT NULL DEFAULT NULL", wantCode: 1067},
	{name: "modify column not null default null", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t MODIFY COLUMN x INT NOT NULL DEFAULT NULL", wantCode: 1067},
	{name: "change column not null default null", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t CHANGE COLUMN x y INT NOT NULL DEFAULT NULL", wantCode: 1067},
	{name: "add column date default empty", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ADD COLUMN c date DEFAULT ''", wantCode: 1067},
	{name: "add column nullable default null accepted", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ADD COLUMN b INT DEFAULT NULL"},
	{name: "add column enum declaring empty member default empty accepted", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ADD COLUMN a enum('','x') DEFAULT ''"},
	{name: "add column enum not declaring empty member default empty", setup: "CREATE TABLE t (x int)",
		sql: "ALTER TABLE t ADD COLUMN c enum('x','y') DEFAULT ''", wantCode: 1067},

	// Cluster 4: more than one AUTO_INCREMENT column. Unlike MySQL, v8.5.5
	// takes a lone auto column that backs no key, so only the count is a
	// refusal. AUTO_RANDOM is not an auto column for this rule.
	{name: "two auto increment columns", sql: "CREATE TABLE t (a int auto_increment, b int auto_increment)", wantCode: 1075},
	{name: "two auto increment columns with pk", sql: "CREATE TABLE t (a int auto_increment primary key, b int auto_increment)", wantCode: 1075},
	{name: "two auto increment columns table pk", sql: "CREATE TABLE t (a int auto_increment, b int auto_increment, primary key(a))", wantCode: 1075},
	{name: "two auto increment columns both unique", sql: "CREATE TABLE t (a int auto_increment unique, b int auto_increment unique)", wantCode: 1075},
	{name: "two serial columns", sql: "CREATE TABLE t (a SERIAL, b SERIAL)", wantCode: 1075},
	{name: "serial and auto increment", sql: "CREATE TABLE t (a SERIAL, b int auto_increment)", wantCode: 1075},
	{name: "one auto increment pk accepted", sql: "CREATE TABLE t (a int auto_increment primary key)"},
	{name: "one unkeyed auto increment accepted", sql: "CREATE TABLE t (a int auto_increment)"},
	{name: "one serial accepted", sql: "CREATE TABLE t (a SERIAL)"},
	{name: "auto random beside auto increment accepted", sql: "CREATE TABLE t (a bigint auto_random primary key, b int auto_increment)"},

	// Cluster 5: a key part naming a column the table does not have. A
	// table-level PRIMARY KEY in CREATE TABLE is worded differently from the
	// same key added by ALTER, and both are 1072. A FULLTEXT key over a missing
	// column is accepted: v8.5.5 reads one and keeps nothing.
	{name: "create index missing column", setup: "CREATE TABLE t (id int)",
		sql: "CREATE INDEX idx_c ON t (c)", wantCode: 1072},
	{name: "create unique index missing column", setup: "CREATE TABLE t (id int)",
		sql: "CREATE UNIQUE INDEX uidx ON t (c)", wantCode: 1072},
	{name: "create table index missing column", sql: "CREATE TABLE t (id int, INDEX idx_c (c))", wantCode: 1072},
	{name: "create table unique missing column", sql: "CREATE TABLE t (id int, UNIQUE KEY uk_c (c))", wantCode: 1072},
	{name: "create table pk missing column", sql: "CREATE TABLE t (id int, PRIMARY KEY (c))", wantCode: 1072},
	{name: "create table index missing second column", sql: "CREATE TABLE t (id int, UNIQUE INDEX ui (id, c))", wantCode: 1072},
	{name: "create table prefix key missing column", sql: "CREATE TABLE t (id varchar(10), INDEX idx_p (c(3)))", wantCode: 1072},
	{name: "alter add index missing column", setup: "CREATE TABLE t (id int)",
		sql: "ALTER TABLE t ADD INDEX idx_c (c)", wantCode: 1072},
	{name: "alter add unique missing column", setup: "CREATE TABLE t (id int)",
		sql: "ALTER TABLE t ADD UNIQUE KEY uk (c)", wantCode: 1072},
	{name: "alter add pk missing column", setup: "CREATE TABLE t (id int)",
		sql: "ALTER TABLE t ADD PRIMARY KEY (c)", wantCode: 1072},
	{name: "key over existing column accepted", sql: "CREATE TABLE t (id int, KEY idx_x (id))"},
	{name: "key over existing column other case accepted", sql: "CREATE TABLE t (a int, INDEX idx (A))"},
	{name: "key over generated column accepted", sql: "CREATE TABLE t (a int, b int as (a+1), INDEX idx (b))"},
	{name: "functional key part accepted", sql: "CREATE TABLE t (id int, INDEX idx_e ((id+1)))"},
	{name: "fulltext key over missing column accepted", sql: "CREATE TABLE t (id int, FULLTEXT KEY ft_c (c))"},
}

// execError runs sql through the catalog and returns the one error it answers
// with: a parse error from Exec itself, or the first statement's own error,
// which Exec reports in the result rather than as a return value.
func execError(c *Catalog, sql string) error {
	results, err := c.Exec(sql, nil)
	if err != nil {
		return err
	}
	for _, r := range results {
		if r.Error != nil {
			return r.Error
		}
	}
	return nil
}

// TestCatalogInvalidDefaultForms holds the catalog to TiDB's verdict without a
// server: the refused forms come back with TiDB's own error number, and the
// accepted ones still go through.
func TestCatalogInvalidDefaultForms(t *testing.T) {
	for _, f := range invalidDefaultForms {
		t.Run(f.name, func(t *testing.T) {
			c := setupWithDB(t)
			if f.setup != "" {
				mustExec(t, c, f.setup)
			}
			err := execError(c, f.sql)
			if f.wantCode == 0 {
				if err != nil {
					t.Fatalf("Exec refused %q, which TiDB accepts: %v", f.sql, err)
				}
				return
			}
			var catErr *Error
			if !errors.As(err, &catErr) {
				t.Fatalf("Exec answered %v for %q, want a catalog error %d", err, f.sql, f.wantCode)
			}
			if catErr.Code != f.wantCode {
				t.Fatalf("Exec answered %d for %q, want %d (%v)", catErr.Code, f.sql, f.wantCode, catErr)
			}
		})
	}
}

// TestCatalogInvalidDefaultFormsLeaveNothingBehind asks that a refused
// statement changes no schema: a half-applied CREATE TABLE or ALTER would let a
// caller that ignores the error read a table the engine never made.
func TestCatalogInvalidDefaultFormsLeaveNothingBehind(t *testing.T) {
	for _, f := range invalidDefaultForms {
		if f.wantCode == 0 {
			continue
		}
		t.Run(f.name, func(t *testing.T) {
			c := setupWithDB(t)
			before := ""
			if f.setup != "" {
				mustExec(t, c, f.setup)
				before = c.ShowCreateTable("testdb", "t")
			}
			if err := execError(c, f.sql); err == nil {
				t.Fatalf("Exec accepted %q, which TiDB refuses", f.sql)
			}
			db := c.GetDatabase("testdb")
			if f.setup == "" {
				if len(db.Tables) != 0 {
					t.Fatalf("the refused statement left %d tables, want 0", len(db.Tables))
				}
				return
			}
			if after := c.ShowCreateTable("testdb", "t"); after != before {
				t.Fatalf("the refused statement changed table t:\n%s\nwant\n%s", after, before)
			}
		})
	}
}

// TestContainerInvalidDefaultForms asks TiDB v8.5.5 that the table above is its
// answer, error number for error number. This is the differential half: it is
// what keeps the catalog's rules from drifting away from the engine, and what
// would catch a form the engine starts accepting in a later version.
func TestContainerInvalidDefaultForms(t *testing.T) {
	tc := startTiDBForCatalog(t)
	mustExecTiDB(t, tc, "CREATE DATABASE IF NOT EXISTS omni_invalid_default_test")
	mustExecTiDB(t, tc, "USE omni_invalid_default_test")
	t.Cleanup(func() {
		_, _ = tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_invalid_default_test")
	})

	for _, f := range invalidDefaultForms {
		t.Run(f.name, func(t *testing.T) {
			mustExecTiDB(t, tc, "DROP TABLE IF EXISTS t")
			if f.setup != "" {
				mustExecTiDB(t, tc, f.setup)
			}
			_, err := tc.db.ExecContext(tc.ctx, f.sql)
			if f.wantCode == 0 {
				if err != nil {
					t.Fatalf("TiDB refused %q, which the catalog accepts: %v", f.sql, err)
				}
				return
			}
			var myErr *gomysql.MySQLError
			if !errors.As(err, &myErr) {
				t.Fatalf("TiDB answered %v for %q, want error %d", err, f.sql, f.wantCode)
			}
			if int(myErr.Number) != f.wantCode {
				t.Fatalf("TiDB answered %d for %q, want %d: %v", myErr.Number, f.sql, f.wantCode, err)
			}
		})
	}
}

// TestEnumMembersFromColumnType covers the reader that recovers an ENUM's
// declared members from its rendered column type, which is where the catalog
// keeps them once a column is built. A member holding a quote is the case no
// DDL test above reaches, and getting it wrong would silently mis-judge an
// ENUM default.
func TestEnumMembersFromColumnType(t *testing.T) {
	cases := []struct {
		columnType string
		want       []string
	}{
		{columnType: "enum('x','y')", want: []string{"x", "y"}},
		{columnType: "enum('','x')", want: []string{"", "x"}},
		{columnType: "enum('a','')", want: []string{"a", ""}},
		{columnType: "enum(' ','x')", want: []string{" ", "x"}},
		{columnType: "enum('a','b''c')", want: []string{"a", "b'c"}},
		{columnType: "enum('''')", want: []string{"'"}},
		{columnType: "enum('a,b','c')", want: []string{"a,b", "c"}},
		{columnType: "set('x','y')", want: []string{"x", "y"}},
		{columnType: "enum('x') character set binary", want: []string{"x"}},
		// Types that declare no members, and a rendering this reader does not
		// understand, both come back nil: a caller must treat them as unknown.
		{columnType: "int", want: nil},
		{columnType: "varchar(10)", want: nil},
		{columnType: "enum('x'", want: nil},
	}
	for _, c := range cases {
		t.Run(c.columnType, func(t *testing.T) {
			got := enumMembersFromColumnType(c.columnType)
			if len(got) != len(c.want) {
				t.Fatalf("read %q as %q, want %q", c.columnType, got, c.want)
			}
			for i := range got {
				if got[i] != c.want[i] {
					t.Fatalf("read %q as %q, want %q", c.columnType, got, c.want)
				}
			}
		})
	}
}
