package parser

import "testing"

// TestConstraintNames pins the names a table constraint keeps apart: the
// CONSTRAINT symbol in Name, and the name after FOREIGN KEY in IndexName.
// TiDB names a foreign key after the symbol or else that name, and its index too, and a UNIQUE key after the symbol, or the name after UNIQUE without one
// (verified on v8.5.5).
func TestConstraintNames(t *testing.T) {
	cases := []struct {
		sql             string
		name, indexName string
	}{
		{"CREATE TABLE x (a int, CONSTRAINT c FOREIGN KEY f (a) REFERENCES r (id))", "c", "f"},
		{"CREATE TABLE x (a int, FOREIGN KEY f (a) REFERENCES r (id))", "", "f"},
		{"CREATE TABLE x (a int, CONSTRAINT FOREIGN KEY f (a) REFERENCES r (id))", "", "f"},
		{"CREATE TABLE x (a int, CONSTRAINT c FOREIGN KEY (a) REFERENCES r (id))", "c", ""},
		{"CREATE TABLE x (a int, CONSTRAINT c UNIQUE KEY u (a))", "c", ""},
		{"CREATE TABLE x (a int, CONSTRAINT c UNIQUE (a))", "c", ""},
		{"ALTER TABLE x ADD CONSTRAINT c FOREIGN KEY f (a) REFERENCES r (id)", "c", "f"},
		{"ALTER TABLE x ADD FOREIGN KEY f (a) REFERENCES r (id)", "", "f"},
		{"ALTER TABLE x ADD CONSTRAINT c UNIQUE INDEX u (a)", "c", ""},
	}
	for _, c := range cases {
		var name, indexName string
		if c.sql[0] == 'C' {
			con := parseCreateTable(t, c.sql).Constraints[0]
			name, indexName = con.Name, con.IndexName
		} else {
			con := parseAlterTable(t, c.sql).Commands[0].Constraint
			name, indexName = con.Name, con.IndexName
		}
		if name != c.name || indexName != c.indexName {
			t.Errorf("%s: Name %q, IndexName %q; want %q, %q", c.sql, name, indexName, c.name, c.indexName)
		}
	}
}
