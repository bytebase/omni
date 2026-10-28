package catalog

import (
	"slices"
	"testing"
)

// TestForeignKeyAndIndexNames pins how MariaDB names a foreign key and
// the index it adds when the CONSTRAINT symbol, the name after FOREIGN
// KEY, or both are written, and a UNIQUE key with a symbol and a name
// (verified on 11.8.9).
func TestForeignKeyAndIndexNames(t *testing.T) {
	c := wtSetup(t)
	wtExec(t, c, `CREATE TABLE r (id INT PRIMARY KEY);
CREATE TABLE x1 (a INT, CONSTRAINT c FOREIGN KEY f (a) REFERENCES r (id));
CREATE TABLE x2 (a INT, FOREIGN KEY f2 (a) REFERENCES r (id));
CREATE TABLE x3 (a INT, CONSTRAINT FOREIGN KEY f3 (a) REFERENCES r (id));
CREATE TABLE x4 (a INT, CONSTRAINT c4 FOREIGN KEY (a) REFERENCES r (id));
CREATE TABLE x5 (a INT, KEY k5 (a), CONSTRAINT c5 FOREIGN KEY f5 (a) REFERENCES r (id));
CREATE TABLE x6 (a INT, CONSTRAINT cu UNIQUE KEY u6 (a));
CREATE TABLE y1 (a INT);
ALTER TABLE y1 ADD CONSTRAINT c7 FOREIGN KEY f7 (a) REFERENCES r (id);
CREATE TABLE y2 (a INT);
ALTER TABLE y2 ADD FOREIGN KEY f8 (a) REFERENCES r (id);
CREATE TABLE y3 (a INT);
ALTER TABLE y3 ADD CONSTRAINT cu UNIQUE KEY u9 (a);`)
	for _, w := range []struct {
		table             string
		foreignKey, index string
	}{
		{"x1", "c", "c"},
		{"x2", "f2", "f2"},
		{"x3", "f3", "f3"},
		{"x4", "c4", "c4"},
		{"x5", "c5", "k5"},
		{"x6", "", "u6"},
		{"y1", "c7", "c7"},
		{"y2", "f8", "f8"},
		{"y3", "", "u9"},
	} {
		tbl := c.GetDatabase("testdb").GetTable(w.table)
		var fks, indexes []string
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey {
				fks = append(fks, con.Name)
			}
		}
		for _, idx := range tbl.Indexes {
			indexes = append(indexes, idx.Name)
		}
		if w.foreignKey != "" && !slices.Equal(fks, []string{w.foreignKey}) {
			t.Errorf("%s: foreign keys %v, want [%s]", w.table, fks, w.foreignKey)
		}
		if !slices.Equal(indexes, []string{w.index}) {
			t.Errorf("%s: indexes %v, want [%s]", w.table, indexes, w.index)
		}
	}
}
