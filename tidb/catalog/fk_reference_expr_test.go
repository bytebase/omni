package catalog

import (
	"strings"
	"testing"
)

// TestForeignKeyReferenceExpressionRejected pins that a foreign key whose
// REFERENCES list uses an expression key part is refused by the catalog, the
// way TiDB v8.5.5's DDL refuses it (it raises 1105 on
// `REFERENCES parent((id+1))`). omni's parser accepts the syntax -- it matches
// TiDB's parser -- so the catalog is where the invalid key is caught. A prefix
// length, which TiDB's DDL does run, is accepted and resolves to the column.
func TestForeignKeyReferenceExpressionRejected(t *testing.T) {
	c := setupWithDB(t)
	mustExec(t, c, "CREATE TABLE parent (id INT NOT NULL, PRIMARY KEY (id))")

	t.Run("expression reference is rejected", func(t *testing.T) {
		results, err := c.Exec("CREATE TABLE child1 (id INT, FOREIGN KEY (id) REFERENCES parent((id+1)))", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Fatalf("expected an exec error, got %+v", results)
		}
		if !strings.Contains(results[0].Error.Error(), "foreign key") {
			t.Errorf("error = %v, want a foreign-key error", results[0].Error)
		}
		if db := c.GetDatabase("testdb"); db != nil && db.GetTable("child1") != nil {
			t.Errorf("child1 was created despite the invalid foreign key")
		}
	})

	t.Run("prefix length reference is accepted and resolves to the column", func(t *testing.T) {
		mustExec(t, c, "CREATE TABLE child2 (id INT, FOREIGN KEY (id) REFERENCES parent(id(8)))")
		tbl := c.GetDatabase("testdb").GetTable("child2")
		if tbl == nil {
			t.Fatal("child2 not created")
		}
		fk := fkConstraint(t, tbl)
		if len(fk.RefColumns) != 1 || fk.RefColumns[0] != "id" {
			t.Errorf("RefColumns = %v, want [id]", fk.RefColumns)
		}
	})

	t.Run("expression reference is rejected by ALTER TABLE too", func(t *testing.T) {
		mustExec(t, c, "CREATE TABLE child3 (id INT)")
		results, err := c.Exec("ALTER TABLE child3 ADD FOREIGN KEY (id) REFERENCES parent((id+1))", nil)
		if err != nil {
			t.Fatalf("parse error: %v", err)
		}
		if len(results) != 1 || results[0].Error == nil {
			t.Fatalf("expected an exec error, got %+v", results)
		}
		tbl := c.GetDatabase("testdb").GetTable("child3")
		for _, con := range tbl.Constraints {
			if con.Type == ConForeignKey {
				t.Errorf("an invalid foreign key was added to child3 by ALTER")
			}
		}
	})
}

// fkConstraint returns the table's one FOREIGN KEY constraint.
func fkConstraint(t *testing.T, tbl *Table) *Constraint {
	t.Helper()
	for _, con := range tbl.Constraints {
		if con.Type == ConForeignKey {
			return con
		}
	}
	t.Fatalf("no FOREIGN KEY constraint on %s", tbl.Name)
	return nil
}
