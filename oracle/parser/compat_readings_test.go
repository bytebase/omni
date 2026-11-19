package parser

import (
	"errors"
	"testing"
)

// TestConsistentHashAndExternalTableReadings pins the two readings this
// change adds, at the bytes compat.tsv records the pinned Oracle refusing
// at.
//
// consistent_hash_partitions and consistent_hash_with_subpartitions belong
// to sharded tables (CREATE-TABLE.bnf): Oracle reads CONSISTENT HASH and
// answers ORA-14151 at what follows HASH, before the columns, because this
// server has no sharding. No form of them parses here, so they have no
// control beside them. An external table's column list comes before
// ORGANIZATION EXTERNAL: Oracle reads the access driver type and then
// answers ORA-00906 at the table's name, and the form beside it, with its
// column list, parses.
func TestConsistentHashAndExternalTableReadings(t *testing.T) {
	refusedAt := []struct {
		sql  string
		want int
	}{
		// The token after HASH, which is the column list's '('.
		{`CREATE TABLE a PARTITION BY CONSISTENT HASH(a) TABLESPACE SET a`, 43},
		{`ALTER TABLE a MODIFY PARTITION BY CONSISTENT HASH(a) TABLESPACE SET a`, 49},
		{`CREATE MATERIALIZED VIEW a PARTITION BY CONSISTENT HASH(a) TABLESPACE SET a AS SELECT 1 FROM dual`, 55},
		{`CREATE MATERIALIZED VIEW LOG ON a PARTITION BY CONSISTENT HASH(a) TABLESPACE SET a`, 62},
		// The subpartitioned form stops at the same token.
		{`CREATE TABLE a PARTITION BY CONSISTENT HASH(a) SUBPARTITION BY RANGE(a) TABLESPACE SET a`, 43},
		// A near-miss that drops HASH stops at the '('.
		{`CREATE TABLE a PARTITION BY CONSISTENT(a) TABLESPACE SET a`, 38},
		{`ALTER TABLE a MODIFY PARTITION BY CONSISTENT(a) TABLESPACE SET a`, 44},
		// An external table with no column list: Oracle reads TYPE and its
		// name, then answers ORA-00906 at the table's name, before the rest
		// of the clause.
		{`CREATE TABLE a ORGANIZATION EXTERNAL(TYPE a)`, 13},
		{`CREATE TABLE a ORGANIZATION EXTERNAL(TYPE a zzz)`, 13},
		{`CREATE TABLE a ORGANIZATION EXTERNAL(TYPE a) AS SELECT 1 FROM dual`, 13},
		// A clause that does not start with TYPE is refused inside it, so
		// its own byte answers.
		{`CREATE TABLE a ORGANIZATION EXTERNAL(a)`, 37},
		{`CREATE TABLE a ORGANIZATION EXTERNAL(TYPE)`, 41},
	}
	for _, tc := range refusedAt {
		_, err := Parse(tc.sql)
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Errorf("Parse(%q): error = %v, want a ParseError", tc.sql, err)
			continue
		}
		if pe.Position != tc.want {
			t.Errorf("Parse(%q): refused at %d, want %d", tc.sql, pe.Position, tc.want)
		}
	}

	// The controls: an external table's clause after its column list, which
	// is the form the page derives and Oracle takes.
	for _, sql := range []string{
		`CREATE TABLE a (a NUMBER) ORGANIZATION EXTERNAL(TYPE a)`,
		`CREATE TABLE a (a NUMBER) ORGANIZATION EXTERNAL(TYPE a DEFAULT DIRECTORY d LOCATION('s'))`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q): %v, want accepted", sql, err)
		}
	}
}
