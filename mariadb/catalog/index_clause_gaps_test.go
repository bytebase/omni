package catalog

import (
	"context"
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"
)

// These hold what ADD ... IF NOT EXISTS does to the catalog, each read from mariadb:11.8.8
// on 2026-11-05: what it finds makes the clause a note and leaves the table as it was.

const addIfNotExistsSeed = `CREATE TABLE r (id INT PRIMARY KEY);
CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0));
CREATE TABLE np (id INT NOT NULL, a INT);
CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id));`

var addIfNotExistsCases = []struct {
	name    string
	table   string
	sql     string
	has     []string
	hasNot  []string
	errCode int
}{
	{name: "index found", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS k (a)", hasNot: []string{"KEY `k` (`a`)"}, has: []string{"KEY `k` (`b`)"}},
	{name: "index new", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a)", has: []string{"KEY `n` (`a`)"}},
	{name: "unnamed takes its column's name", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS (b)", has: []string{"KEY `b` (`b`)"}},
	{name: "unnamed found by its column's name", table: "t", sql: "ALTER TABLE t ADD UNIQUE IF NOT EXISTS (a)", hasNot: []string{"UNIQUE KEY"}},
	{name: "unique found", table: "t", sql: "ALTER TABLE t ADD UNIQUE IF NOT EXISTS k (a)", hasNot: []string{"UNIQUE KEY"}},
	{name: "the name after UNIQUE decides", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT nc UNIQUE IF NOT EXISTS k (a)", hasNot: []string{"UNIQUE KEY", "nc"}},
	{name: "the symbol names it otherwise", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT nc UNIQUE IF NOT EXISTS (a)", has: []string{"UNIQUE KEY `nc` (`a`)"}},
	{name: "found comes before the columns", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS k (nosuch)", hasNot: []string{"nosuch"}},
	{name: "twice in one statement", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX IF NOT EXISTS n (b)", has: []string{"KEY `n` (`a`)"}, hasNot: []string{"KEY `n` (`b`)"}},
	{name: "twice without it", table: "t", sql: "ALTER TABLE t ADD INDEX n (a), ADD INDEX n (b)", errCode: 1061},
	{name: "twice in the list", table: "t", sql: "ALTER TABLE t ADD (INDEX IF NOT EXISTS n (a), INDEX IF NOT EXISTS n (b))", has: []string{"KEY `n` (`a`)"}, hasNot: []string{"KEY `n` (`b`)"}},
	{name: "primary key twice on a table with none", table: "np", sql: "ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id), ADD PRIMARY KEY IF NOT EXISTS (id)", errCode: 1068},
	{name: "primary key twice on a table with one", table: "t", sql: "ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY IF NOT EXISTS (b)", has: []string{"PRIMARY KEY (`id`)"}},
	{name: "plain after IF NOT EXISTS", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD INDEX n (b)", errCode: 1061},
	{name: "IF NOT EXISTS after plain", table: "t", sql: "ALTER TABLE t ADD INDEX n (a), ADD INDEX IF NOT EXISTS n (b)", has: []string{"KEY `n` (`a`)"}, hasNot: []string{"KEY `n` (`b`)"}},
	{name: "check twice", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 1)", errCode: 1826},
	{name: "check found twice", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 1), ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 2)", has: []string{"`a` > 0"}},
	{name: "foreign key twice", table: "f", sql: "ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fq (a) REFERENCES r (id), ADD FOREIGN KEY IF NOT EXISTS fq (a) REFERENCES r (id)", has: []string{"CONSTRAINT `fq` FOREIGN KEY"}},
	{name: "primary key found", table: "t", sql: "ALTER TABLE t ADD PRIMARY KEY IF NOT EXISTS (a)", has: []string{"PRIMARY KEY (`id`)"}},
	{name: "primary key new", table: "np", sql: "ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id)", has: []string{"PRIMARY KEY (`id`)"}},
	{name: "primary key without it", table: "t", sql: "ALTER TABLE t ADD PRIMARY KEY (a)", errCode: 1068},
	{name: "check found", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck CHECK (a > 1)", has: []string{"`a` > 0"}, hasNot: []string{"`a` > 1"}},
	{name: "check new", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS ck2 CHECK (a > 1)", has: []string{"CONSTRAINT `ck2` CHECK"}},
	{name: "foreign key found by its name", table: "f", sql: "ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fkx (a) REFERENCES r (id)", hasNot: []string{"f_ibfk"}},
	{name: "foreign key found by its symbol", table: "f", sql: "ALTER TABLE f ADD CONSTRAINT fkx FOREIGN KEY IF NOT EXISTS (a) REFERENCES r (id)", hasNot: []string{"f_ibfk"}},
	{name: "foreign key with no name is always added", table: "f", sql: "ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS (a) REFERENCES r (id)", has: []string{"f_ibfk_1"}},
	{name: "foreign key new", table: "f", sql: "ALTER TABLE f ADD FOREIGN KEY IF NOT EXISTS fkq (a) REFERENCES r (id)", has: []string{"CONSTRAINT `fkq` FOREIGN KEY"}},
	// A primary key or CHECK the statement also adds is not there to find, in either order.
	{name: "primary key plain then IF NOT EXISTS", table: "np", sql: "ALTER TABLE np ADD PRIMARY KEY (id), ADD PRIMARY KEY IF NOT EXISTS (a)", errCode: 1068},
	{name: "primary key IF NOT EXISTS then plain", table: "np", sql: "ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (a), ADD PRIMARY KEY (id)", errCode: 1068},
	{name: "check plain then IF NOT EXISTS", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT c1 CHECK (a > 0), ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 1)", errCode: 1826},
	{name: "check IF NOT EXISTS then plain", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD CONSTRAINT c1 CHECK (a > 1)", errCode: 1826},
	// A column's own PRIMARY KEY counts as well.
	{name: "inline primary key then IF NOT EXISTS", table: "np", sql: "ALTER TABLE np ADD b INT PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)", errCode: 1068},
	{name: "IF NOT EXISTS then inline primary key", table: "np", sql: "ALTER TABLE np ADD PRIMARY KEY IF NOT EXISTS (id), ADD b INT PRIMARY KEY", errCode: 1068},
	{name: "listed column's primary key", table: "np", sql: "ALTER TABLE np ADD COLUMN (b INT PRIMARY KEY), ADD PRIMARY KEY IF NOT EXISTS (id)", errCode: 1068},
	{name: "MODIFY's primary key", table: "np", sql: "ALTER TABLE np MODIFY id INT NOT NULL PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)", errCode: 1068},
	{name: "CHANGE's primary key", table: "np", sql: "ALTER TABLE np CHANGE id id INT NOT NULL PRIMARY KEY, ADD PRIMARY KEY IF NOT EXISTS (id)", errCode: 1068},
	// PRIMARY is no index's name unless the clause finds the primary key and stands down.
	{name: "PRIMARY on a table with none", table: "np", sql: "ALTER TABLE np ADD INDEX IF NOT EXISTS `PRIMARY` (a)", errCode: 1280},
	{name: "primary, unique, on a table with none", table: "np", sql: "ALTER TABLE np ADD UNIQUE IF NOT EXISTS `primary` (a)", errCode: 1280},
	{name: "PRIMARY on a table with one", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS `PRIMARY` (a)", has: []string{"PRIMARY KEY (`id`)"}},
	// A CHECK and an index have a name each.
	{name: "check and index share a name", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT IF NOT EXISTS c1 CHECK (a > 0), ADD INDEX c1 (a)", has: []string{"KEY `c1` (`a`)", "CONSTRAINT `c1` CHECK"}},
	{name: "index and check share a name", table: "t", sql: "ALTER TABLE t ADD INDEX c2 (a), ADD CONSTRAINT IF NOT EXISTS c2 CHECK (a > 0)", has: []string{"KEY `c2` (`a`)", "CONSTRAINT `c2` CHECK"}},
	// A second index of a name goes as a note only of its own kind.
	{name: "plain then unique IF NOT EXISTS", table: "t", sql: "ALTER TABLE t ADD INDEX n (a), ADD UNIQUE IF NOT EXISTS n (a)", errCode: 1061},
	{name: "unique then plain IF NOT EXISTS", table: "t", sql: "ALTER TABLE t ADD UNIQUE n (a), ADD INDEX IF NOT EXISTS n (a)", errCode: 1061},
	{name: "both IF NOT EXISTS, two kinds", table: "t", sql: "ALTER TABLE t ADD INDEX IF NOT EXISTS n (a), ADD UNIQUE IF NOT EXISTS n (a)", errCode: 1061},
	{name: "both IF NOT EXISTS, one kind", table: "t", sql: "ALTER TABLE t ADD UNIQUE IF NOT EXISTS n (a), ADD UNIQUE IF NOT EXISTS n (b)", has: []string{"UNIQUE KEY `n` (`a`)"}, hasNot: []string{"UNIQUE KEY `n` (`b`)"}},
	{name: "KEY and INDEX are one kind", table: "t", sql: "ALTER TABLE t ADD INDEX n (a), ADD KEY IF NOT EXISTS n (b)", has: []string{"KEY `n` (`a`)"}, hasNot: []string{"KEY `n` (`b`)"}},
	{name: "two kinds on a key the table has", table: "t", sql: "ALTER TABLE t ADD UNIQUE IF NOT EXISTS k (a), ADD FULLTEXT IF NOT EXISTS k (a)", has: []string{"KEY `k` (`b`)"}, hasNot: []string{"FULLTEXT", "UNIQUE"}},
	// A foreign key looks for a foreign key of the name, on this table, and an index is not one.
	{name: "foreign key named as an index", table: "t", sql: "ALTER TABLE t ADD FOREIGN KEY IF NOT EXISTS k (b) REFERENCES r (id)", has: []string{"CONSTRAINT `k` FOREIGN KEY (`b`)"}},
	{name: "foreign key with a symbol named as an index", table: "t", sql: "ALTER TABLE t ADD CONSTRAINT k FOREIGN KEY IF NOT EXISTS (b) REFERENCES r (id)", has: []string{"CONSTRAINT `k` FOREIGN KEY (`b`)"}},
}

func TestAddIfNotExists(t *testing.T) {
	for _, tc := range addIfNotExistsCases {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			execAlterClauses(t, c, "CREATE DATABASE d; USE d;\n"+addIfNotExistsSeed)
			err := execAlterClausesExpectError(t, c, tc.sql+";")
			if tc.errCode != 0 {
				var cerr *Error
				if !errors.As(err, &cerr) || cerr.Code != tc.errCode {
					t.Fatalf("error %v, want %d", err, tc.errCode)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			ddl := c.ShowCreateTable("d", tc.table)
			for _, want := range tc.has {
				if !strings.Contains(ddl, want) {
					t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
				}
			}
			for _, not := range tc.hasNot {
				if strings.Contains(ddl, not) {
					t.Errorf("SHOW CREATE TABLE has %s:\n%s", not, ddl)
				}
			}
		})
	}
}

// TestAddIfNotExistsOracle holds addIfNotExistsCases to the server.
func TestAddIfNotExistsOracle(t *testing.T) {
	ctx := context.Background()
	conn := startDropPeriodMariaDB(t)
	reset := func() {
		for _, table := range []string{"f", "np", "t", "r"} {
			if _, err := conn.ExecContext(ctx, "DROP TABLE IF EXISTS "+table); err != nil {
				t.Fatal(err)
			}
		}
	}
	t.Cleanup(reset)
	for _, tc := range addIfNotExistsCases {
		t.Run(tc.name, func(t *testing.T) {
			reset()
			for _, stmt := range strings.Split(addIfNotExistsSeed, ";\n") {
				if _, err := conn.ExecContext(ctx, strings.TrimSuffix(stmt, ";")); err != nil {
					t.Fatalf("seed %q: %v", stmt, err)
				}
			}
			var code int
			if _, err := conn.ExecContext(ctx, tc.sql); err != nil {
				var myErr *gomysql.MySQLError
				if !errors.As(err, &myErr) {
					t.Fatal(err)
				}
				code = int(myErr.Number)
			}
			if code != tc.errCode {
				t.Fatalf("the server answered %d, the record says %d", code, tc.errCode)
			}
			if code != 0 {
				return
			}
			var name, ddl string
			if err := conn.QueryRowContext(ctx, "SHOW CREATE TABLE "+tc.table).Scan(&name, &ddl); err != nil {
				t.Fatal(err)
			}
			for _, want := range tc.has {
				if !strings.Contains(ddl, want) {
					t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
				}
			}
			for _, not := range tc.hasNot {
				if strings.Contains(ddl, not) {
					t.Errorf("SHOW CREATE TABLE has %s:\n%s", not, ddl)
				}
			}
		})
	}
}

// TestNamedScalarTypeSizeDropped: inet6(16) reads back as inet6, as it does on 11.8.8.
func TestNamedScalarTypeSizeDropped(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a inet6(16), b uuid(36), c inet4(1.5));")
	ddl := c.ShowCreateTable("d", "t")
	for _, want := range []string{"`a` inet6 DEFAULT NULL", "`b` uuid DEFAULT NULL", "`c` inet4 DEFAULT NULL"} {
		if !strings.Contains(ddl, want) {
			t.Errorf("SHOW CREATE TABLE lacks %s:\n%s", want, ddl)
		}
	}
}

// TestCreateAndDropIndexWait: WAIT changes nothing about the index.
func TestCreateAndDropIndexWait(t *testing.T) {
	c := New()
	execAlterClauses(t, c, "CREATE DATABASE d; USE d; CREATE TABLE t (a INT); CREATE INDEX i ON t (a) WAIT 3 COMMENT 'c';")
	if !strings.Contains(c.ShowCreateTable("d", "t"), "KEY `i` (`a`) COMMENT 'c'") {
		t.Errorf("the index with WAIT was lost:\n%s", c.ShowCreateTable("d", "t"))
	}
	execAlterClauses(t, c, "DROP INDEX i ON t NOWAIT;")
	if strings.Contains(c.ShowCreateTable("d", "t"), "KEY `i`") {
		t.Errorf("DROP INDEX ... NOWAIT left the index")
	}
}
