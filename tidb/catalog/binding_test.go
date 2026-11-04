package catalog

import (
	"sort"
	"testing"
)

// TestBindingStatementsChangeNothing pins what a SQL plan binding does to
// the catalog: nothing. A binding records a plan for later statements, and
// the statements it carries are not run, so neither the binding nor its
// inner CREATE-shaped text may reach the schema. The DELETE inside one
// deletes no row and the table it names need not exist.
func TestBindingStatementsChangeNothing(t *testing.T) {
	c := New()
	if _, err := c.Exec("CREATE DATABASE test; USE test; CREATE TABLE t (a INT, b INT);", nil); err != nil {
		t.Fatalf("seed: %v", err)
	}
	before := schemaShape(t, c)

	sql := "CREATE GLOBAL BINDING FOR SELECT * FROM t USING SELECT * FROM t USE INDEX(a);\n" +
		"CREATE BINDING USING UPDATE t SET a = 1;\n" +
		"CREATE BINDING FROM HISTORY USING PLAN DIGEST 'abc';\n" +
		"DROP BINDING FOR DELETE FROM nosuchtable;\n" +
		"DROP SESSION BINDING FOR SQL DIGEST 'abc';\n" +
		"SET BINDING DISABLED FOR SELECT * FROM t;"

	results, err := c.Exec(sql, nil)
	if err != nil {
		t.Fatalf("Exec: %v", err)
	}
	if len(results) != 6 {
		t.Fatalf("got %d results, want 6", len(results))
	}
	// Each statement is applied, not skipped as DML is, and each applies
	// nothing. A statement the catalog does not know would also report no
	// error, so the line below is what tells the two apart: an unknown
	// statement has no case in stmtLocStart and reports line 1.
	for i, r := range results {
		if r.Error != nil {
			t.Errorf("result[%d]: %v", i, r.Error)
		}
		if r.Skipped {
			t.Errorf("result[%d] was skipped, want applied", i)
		}
		if want := i + 1; r.Line != want {
			t.Errorf("result[%d].Line = %d, want %d", i, r.Line, want)
		}
	}

	if after := schemaShape(t, c); after != before {
		t.Errorf("the schema changed:\n  before %s\n  after  %s", before, after)
	}
	// The binding named a table that does not exist, and that is no error
	// here: the catalog never resolved it.
	if db := c.GetDatabase("test"); db == nil || db.GetTable("nosuchtable") != nil {
		t.Errorf("the catalog made a table the binding named")
	}
}

// schemaShape writes every database, table, and column the catalog holds,
// so a test can compare the whole schema before and after a statement.
func schemaShape(t *testing.T, c *Catalog) string {
	t.Helper()
	var out []string
	for _, db := range c.Databases() {
		for name, table := range db.Tables {
			cols := make([]string, 0, len(table.Columns))
			for _, col := range table.Columns {
				cols = append(cols, col.Name)
			}
			sort.Strings(cols)
			out = append(out, db.Name+"."+name+"("+joinComma(cols)+")")
		}
		for name := range db.Views {
			out = append(out, db.Name+"."+name+" view")
		}
		if len(db.Tables) == 0 && len(db.Views) == 0 {
			out = append(out, db.Name)
		}
	}
	sort.Strings(out)
	return joinComma(out)
}

func joinComma(parts []string) string {
	out := ""
	for i, p := range parts {
		if i > 0 {
			out += ", "
		}
		out += p
	}
	return out
}
