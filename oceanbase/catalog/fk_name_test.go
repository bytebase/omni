package catalog

import (
	"slices"
	"testing"
)

// TestForeignKeyAndIndexNames pins how a foreign key is named when the
// CONSTRAINT symbol, the name after FOREIGN KEY, or both are written, and
// that no index is made for it, as OceanBase makes none (measured on
// 4.3.5.6: each table below has the indexes it declares). A UNIQUE key
// with a symbol and a name takes the name. OceanBase names a foreign key
// that has no symbol `<table>_OBFK_<time>`, and so does the loader, with
// a count for the time (madeUpConstraintName).
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
		table      string
		foreignKey string
		indexes    []string
	}{
		{"x1", "c", nil},
		{"x2", "x2_OBFK_<n>", nil},
		{"x3", "x3_OBFK_<n>", nil},
		{"x4", "c4", nil},
		{"x5", "c5", []string{"k5"}},
		{"x6", "", []string{"u6"}},
		{"y1", "c7", nil},
		{"y2", "y2_OBFK_<n>", nil},
		{"y3", "", []string{"u9"}},
	} {
		tbl := c.GetDatabase("testdb").GetTable(w.table)
		var fks, indexes []string
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey {
				fks = append(fks, withoutCount(con.Name))
			}
		}
		for _, idx := range tbl.Indexes {
			indexes = append(indexes, idx.Name)
		}
		if w.foreignKey != "" && !slices.Equal(fks, []string{w.foreignKey}) {
			t.Errorf("%s: foreign keys %v, want [%s]", w.table, fks, w.foreignKey)
		}
		if !slices.Equal(indexes, w.indexes) {
			t.Errorf("%s: indexes %v, want %v", w.table, indexes, w.indexes)
		}
	}
}
