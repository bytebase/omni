package parser

import (
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
)

// TestForeignKeyReferenceIndexParts pins the tree a REFERENCES column list
// parses to now that it reads an IndexPartSpecificationList (a column with an
// optional prefix length, or an expression) rather than plain column names, so
// omni accepts the forms TiDB v8.5.5 does. These are the gap rows the §3.4
// rest corpus surfaced; each statement here is measured accept/accept there.
//
// The referenced-column names are what consumers read: a prefix length is
// dropped (the name is still the column), and an expression key part yields no
// name at all, as indexColumnsToNames does for the local column list.
func TestForeignKeyReferenceIndexParts(t *testing.T) {
	for _, c := range []struct {
		name     string
		sql      string
		wantRef  []string // referenced column names
		wantCols []string // foreign key's own column names
	}{
		{
			name:     "prefix length on one referenced column",
			sql:      "CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent(id,hello(12)) ON DELETE CASCADE)",
			wantRef:  []string{"id", "hello"},
			wantCols: []string{"parent_id"},
		},
		{
			name:     "prefix length on both referenced columns",
			sql:      "CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent(id(8),hello(12)) ON DELETE CASCADE)",
			wantRef:  []string{"id", "hello"},
			wantCols: []string{"parent_id"},
		},
		{
			name:     "expression referenced key part keeps no name",
			sql:      "CREATE TABLE child (id INT, parent_id INT, INDEX par_ind (parent_id), FOREIGN KEY (parent_id) REFERENCES parent((id+1)))",
			wantRef:  []string{},
			wantCols: []string{"parent_id"},
		},
		{
			name:     "prefix lengths in the FK column list, expression reference",
			sql:      "CREATE TABLE child (id INT, parent_id INT, CONSTRAINT FOREIGN KEY (parent_id(2),hello(4)) REFERENCES parent((id+1)) ON DELETE CASCADE)",
			wantRef:  []string{},
			wantCols: []string{"parent_id", "hello"},
		},
		{
			name:     "expression reference with ON DELETE and ON UPDATE",
			sql:      "CREATE TABLE child (id INT, parent_id INT, CONSTRAINT FOREIGN KEY (parent_id) REFERENCES parent((id+1)) ON DELETE CASCADE ON UPDATE RESTRICT)",
			wantRef:  []string{},
			wantCols: []string{"parent_id"},
		},
		{
			name:     "plain column reference is unchanged",
			sql:      "CREATE TABLE orders (id INT, user_id INT, FOREIGN KEY (user_id) REFERENCES users(id))",
			wantRef:  []string{"id"},
			wantCols: []string{"user_id"},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if _, err := Parse(c.sql); err != nil {
				t.Fatalf("Parse rejected %q: %v", c.sql, err)
			}
			stmt := parseCreateTable(t, c.sql)
			fk := foreignKeyConstraint(t, stmt)
			if !equalStrings(fk.RefColumns, c.wantRef) {
				t.Errorf("RefColumns = %v, want %v", fk.RefColumns, c.wantRef)
			}
			if !equalStrings(fk.Columns, c.wantCols) {
				t.Errorf("Columns = %v, want %v", fk.Columns, c.wantCols)
			}
			// RefKeyParts keeps the list as written, so the catalog can see an
			// expression part that RefColumns drops.
			if len(fk.RefKeyParts) == 0 {
				t.Errorf("RefKeyParts is empty; the key parts were not retained")
			}
		})
	}
}

// TestReferenceKeyPartRequiresLength pins that an empty prefix `()` is a syntax
// error, not a zero-length prefix: TiDB's FieldLen is `'(' LengthNum ')'`, and
// v8.5.5 answers 1064 to each of these. The fix is in the shared key-part
// reader, so it holds for a REFERENCES list and for an index alike.
func TestReferenceKeyPartRequiresLength(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE child (id INT, FOREIGN KEY (id) REFERENCES parent(id()))",
		"CREATE TABLE t (a INT, KEY (a()))",
		"CREATE INDEX i ON t (a())",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse accepted %q; want a syntax error for the empty prefix length", sql)
		}
	}
}

// foreignKeyConstraint returns the table's one FOREIGN KEY constraint, failing
// if there is not exactly one.
func foreignKeyConstraint(t *testing.T, stmt *ast.CreateTableStmt) *ast.Constraint {
	t.Helper()
	var fk *ast.Constraint
	for _, c := range stmt.Constraints {
		if c.Type == ast.ConstrForeignKey {
			if fk != nil {
				t.Fatalf("want one FOREIGN KEY constraint, got more than one")
			}
			fk = c
		}
	}
	if fk == nil {
		t.Fatalf("no FOREIGN KEY constraint in %+v", stmt.Constraints)
	}
	return fk
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
