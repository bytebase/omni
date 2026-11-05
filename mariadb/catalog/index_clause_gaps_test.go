package catalog

import (
	"errors"
	"strings"
	"testing"
)

// These hold what ADD ... IF NOT EXISTS does to the catalog, each read from mariadb:11.8.8
// on 2026-11-05: what it finds makes the clause a note and leaves the table as it was.

func TestAddIfNotExists(t *testing.T) {
	const seed = `CREATE DATABASE d; USE d;
CREATE TABLE r (id INT PRIMARY KEY);
CREATE TABLE t (id INT PRIMARY KEY, a INT, b INT, KEY a (a), KEY k (b), CONSTRAINT ck CHECK (a > 0));
CREATE TABLE np (id INT NOT NULL, a INT);
CREATE TABLE f (id INT PRIMARY KEY, a INT, CONSTRAINT fkx FOREIGN KEY (a) REFERENCES r (id));`
	for _, tc := range []struct {
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
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := New()
			execAlterClauses(t, c, seed)
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
