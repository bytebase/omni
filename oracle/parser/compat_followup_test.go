package parser

import (
	"errors"
	"testing"
)

// TestParseFollowupForms pins the forms #1450 left open and this change
// reads, each with the near-miss the same compat.tsv row records Oracle
// refusing. The pages are the ones the parser cites: ALTER-SESSION,
// SELECT, DROP-MATERIALIZED-VIEW, PURGE, DROP-RESTORE-POINT, CREATE-TABLE,
// CREATE-MATERIALIZED-VIEW and its log, create-vector-index.
func TestParseFollowupForms(t *testing.T) {
	accepted := []string{
		`ALTER SESSION SET USE_PRIVATE_OUTLINES =`,
		`ALTER SESSION SET USE_STORED_OUTLINES =`,
		`SELECT * FROM t NATURAL`,
		`DELETE a FROM t NATURAL`,
		`UPDATE t SET a = 1 FROM t NATURAL`,
		`DROP MATERIALIZED VIEW a PRESERVE TABLE`,
		`PURGE TABLESPACE a USER b`,
		`DROP RESTORE POINT a FOR PLUGGABLE DATABASE p`,
		`CREATE TABLE a PARTITION BY REFERENCE(a)`,
		`ALTER TABLE a ADD (a)`,
		`ALTER TABLE a ADD CONSTRAINT`,
		`ALTER TABLE a ADD CONSTRAINTS c UNIQUE(a)`,
		`ALTER TABLE a ADD CONSTRAINTS c FOREIGN KEY(a) REFERENCES a`,
		`CREATE TABLE a (a NUMBER) PARTITION BY RANGE(a) INTERVAL(1) STORE IN(a)(PARTITION VALUES LESS THAN(1))`,
		`CREATE TABLE a (a NUMBER) PARTITION BY RANGE(a) INTERVAL(1) STORE IN(a,)(PARTITION VALUES LESS THAN(1))`,
		`CREATE MATERIALIZED VIEW a PARTITION BY RANGE(a) INTERVAL(1) STORE IN(a)(PARTITION VALUES LESS THAN(1)) AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW a PARTITION BY RANGE(a) INTERVAL(1) STORE IN(a) SUBPARTITION BY RANGE(a)(PARTITION VALUES LESS THAN(1)) AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW a PARTITION BY HASH(a) PARTITIONS 1 STORE IN(a,) AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW LOG ON a PARTITION BY HASH(a) PARTITIONS 1 STORE IN(a,)`,
		`CREATE MATERIALIZED VIEW a JSON(a) STORE AS AS SELECT 1 FROM dual`,
		`CREATE MATERIALIZED VIEW a JSON(a) AS AS SELECT 1 FROM dual`,
		`CREATE VECTOR INDEX a ON a(a) ORGANIZATION PARTITIONS`,
	}
	for _, sql := range accepted {
		if _, err := Parse(sql); err != nil {
			t.Errorf("%q: %v, want accepted", sql, err)
		}
	}

	// The bytes are compat.tsv's: Oracle reports a reference partitioning's
	// byte at the constraint, PURGE's TABLESPACE SET at the name after SET,
	// and a list's repeated bare column at the name.
	refusedAt := []struct {
		sql  string
		want int
	}{
		{`CREATE TABLE a PARTITION BY REFERENCE(a)(PARTITION p)`, 38},
		{`CREATE TABLE a PARTITION BY REFERENCE(a) zzz`, 38},
		{`PURGE TABLESPACE SET a`, 21},
		{`ALTER TABLE a ADD(a NUMBER, a)`, 28},
	}
	for _, c := range refusedAt {
		_, err := Parse(c.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("%q: %v, want a *ParseError", c.sql, err)
			continue
		}
		if pe.Position != c.want {
			t.Errorf("%q: refused at %d, want %d", c.sql, pe.Position, c.want)
		}
	}

	// Oracle refuses these too, at its own code: a plain ALTER SESSION
	// parameter without a value (ORA-00921 at the end), NATURAL before a
	// word that is no join (ORA-03049), STORE IN after a range without
	// INTERVAL (ORA-00922) or on a log (ORA-12044), an empty STORE IN
	// (ORA-02216), and one AS alone, which the JSON storage clause takes and
	// leaves the view without its own (ORA-00906).
	refused := []string{
		`ALTER SESSION SET optimizer_mode =`,
		`SELECT * FROM t NATURAL a`,
		`CREATE TABLE a PARTITION BY RANGE(a) STORE IN(a)(PARTITION VALUES LESS THAN(1))`,
		`CREATE MATERIALIZED VIEW LOG ON a PARTITION BY RANGE(a) INTERVAL(1) STORE IN(a)(PARTITION VALUES LESS THAN(1))`,
		`CREATE TABLE a (a NUMBER) PARTITION BY RANGE(a) INTERVAL(1) STORE IN()`,
		`CREATE MATERIALIZED VIEW a JSON(a) AS SELECT 1 FROM dual`,
	}
	for _, sql := range refused {
		if _, err := Parse(sql); err == nil {
			t.Errorf("%q: accepted, want refused", sql)
		}
	}
}
