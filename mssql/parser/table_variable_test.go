package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestTableVariableTableSource exercises T-SQL table variables (@t) used as
// table sources in DML. This mirrors SqlScriptDOM's VariableTableReference and
// variableDmlTarget productions.
func TestTableVariableTableSource(t *testing.T) {
	t.Run("select from table variable", func(t *testing.T) {
		list, err := Parse("SELECT * FROM @t")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		tv, ok := sel.FromClause.Items[0].(*ast.TableVarRef)
		if !ok {
			t.Fatalf("expected *TableVarRef in FROM, got %T", sel.FromClause.Items[0])
		}
		if tv.Name != "@t" {
			t.Errorf("Name = %q, want %q", tv.Name, "@t")
		}
	})

	t.Run("select from table variable with alias", func(t *testing.T) {
		list, err := Parse("SELECT * FROM @t AS x")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		tv := sel.FromClause.Items[0].(*ast.TableVarRef)
		if tv.Alias != "x" {
			t.Errorf("Alias = %q, want %q", tv.Alias, "x")
		}
	})

	t.Run("select from table variable with bare alias", func(t *testing.T) {
		list, err := Parse("SELECT * FROM @t x")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		tv := sel.FromClause.Items[0].(*ast.TableVarRef)
		if tv.Alias != "x" {
			t.Errorf("Alias = %q, want %q", tv.Alias, "x")
		}
	})

	t.Run("join two table variables", func(t *testing.T) {
		list, err := Parse("SELECT * FROM @t1 a JOIN @t2 b ON a.id = b.id")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		join := sel.FromClause.Items[0].(*ast.JoinClause)
		if _, ok := join.Left.(*ast.TableVarRef); !ok {
			t.Errorf("join.Left = %T, want *TableVarRef", join.Left)
		}
		if _, ok := join.Right.(*ast.TableVarRef); !ok {
			t.Errorf("join.Right = %T, want *TableVarRef", join.Right)
		}
	})

	t.Run("insert into table variable", func(t *testing.T) {
		list, err := Parse("INSERT INTO @t VALUES (1)")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		ins := list.Items[0].(*ast.InsertStmt)
		tv, ok := ins.Relation.(*ast.TableVarRef)
		if !ok {
			t.Fatalf("Relation = %T, want *TableVarRef", ins.Relation)
		}
		if tv.Name != "@t" {
			t.Errorf("Name = %q, want %q", tv.Name, "@t")
		}
	})

	t.Run("insert into table variable with columns", func(t *testing.T) {
		list, err := Parse("INSERT INTO @t (id, name) VALUES (1, 'a')")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		ins := list.Items[0].(*ast.InsertStmt)
		if _, ok := ins.Relation.(*ast.TableVarRef); !ok {
			t.Fatalf("Relation = %T, want *TableVarRef", ins.Relation)
		}
		if ins.Cols == nil || len(ins.Cols.Items) != 2 {
			t.Errorf("expected 2 columns, got %v", ins.Cols)
		}
	})

	t.Run("insert into table variable with select source", func(t *testing.T) {
		_, err := Parse("INSERT INTO @t SELECT id FROM src")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
	})

	t.Run("update table variable", func(t *testing.T) {
		list, err := Parse("UPDATE @t SET id = 1")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		upd := list.Items[0].(*ast.UpdateStmt)
		if _, ok := upd.Relation.(*ast.TableVarRef); !ok {
			t.Fatalf("Relation = %T, want *TableVarRef", upd.Relation)
		}
	})

	t.Run("delete from table variable", func(t *testing.T) {
		list, err := Parse("DELETE FROM @t WHERE id = 1")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		del := list.Items[0].(*ast.DeleteStmt)
		if _, ok := del.Relation.(*ast.TableVarRef); !ok {
			t.Fatalf("Relation = %T, want *TableVarRef", del.Relation)
		}
	})

	t.Run("delete bare table variable", func(t *testing.T) {
		_, err := Parse("DELETE @t")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
	})

	t.Run("merge into table variable", func(t *testing.T) {
		sql := `MERGE @t AS tgt USING src AS s ON tgt.id = s.id
			WHEN MATCHED THEN UPDATE SET tgt.id = s.id;`
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		m := list.Items[0].(*ast.MergeStmt)
		tv, ok := m.Target.(*ast.TableVarRef)
		if !ok {
			t.Fatalf("Target = %T, want *TableVarRef", m.Target)
		}
		if tv.Alias != "tgt" {
			t.Errorf("Alias = %q, want %q", tv.Alias, "tgt")
		}
	})

	t.Run("declare then use table variable in same batch", func(t *testing.T) {
		sql := `declare @t table(id int); insert into @t values(1); select * from @t;`
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		if len(list.Items) != 3 {
			t.Errorf("expected 3 statements, got %d", len(list.Items))
		}
	})

	t.Run("table variable method call in FROM", func(t *testing.T) {
		// @xmlcol.nodes('/x') AS T(c)
		list, err := Parse("SELECT c FROM @x.nodes('/a') AS T(c)")
		if err != nil {
			t.Fatalf("parse failed: %v", err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		mc, ok := sel.FromClause.Items[0].(*ast.TableVarMethodCallRef)
		if !ok {
			t.Fatalf("expected *TableVarMethodCallRef, got %T", sel.FromClause.Items[0])
		}
		if mc.Var != "@x" {
			t.Errorf("Var = %q, want %q", mc.Var, "@x")
		}
		if mc.Method != "nodes" {
			t.Errorf("Method = %q, want %q", mc.Method, "nodes")
		}
		if mc.Alias != "T" {
			t.Errorf("Alias = %q, want %q", mc.Alias, "T")
		}
		if len(mc.Columns) != 1 || mc.Columns[0] != "c" {
			t.Errorf("Columns = %v, want [c]", mc.Columns)
		}
	})
}

// TestTableVariableNegative verifies that @var is rejected where SqlScriptDOM
// rejects it — namely DDL contexts and as a hint target.
func TestTableVariableNegative(t *testing.T) {
	// Note: `CREATE VIEW v AS SELECT * FROM @t` parses successfully because the
	// body is a SELECT and SELECT FROM @t is a legal table source — it would
	// fail at bind time, not parse time. Only purely DDL-contextual uses of
	// @t as an object name should fail at parse.
	cases := []string{
		"CREATE INDEX ix ON @t (id)",
		"DROP TABLE @t",
	}
	for _, sql := range cases {
		_, err := Parse(sql)
		if err == nil {
			t.Errorf("expected parse error for %q, got nil", sql)
		}
	}
}

// TestTableVariableScalarStillWorks confirms we did not regress scalar @var
// usage in expressions (WHERE, VALUES, SET RHS, SELECT list) — these go
// through parseExpr, not parseTableRef.
func TestTableVariableScalarStillWorks(t *testing.T) {
	cases := []string{
		"SELECT @x",
		"SELECT @x + 1",
		"SET @x = 1",
		"DECLARE @x INT = 5",
		"UPDATE t SET c = @x WHERE id = @y",
		"INSERT INTO t VALUES (@x)",
		"SELECT * FROM t WHERE id = @x",
	}
	for _, sql := range cases {
		if _, err := Parse(sql); err != nil {
			t.Errorf("parse %q failed: %v", sql, err)
		}
	}
}

// Sanity: ensure the DECLARE @t TABLE(...) standalone case still parses,
// as it did before this change. Catches regressions in the declare_set path.
func TestDeclareTableVariableStillParses(t *testing.T) {
	_, err := Parse("DECLARE @t TABLE (id INT, name NVARCHAR(50))")
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}
	// Also multiple columns + constraints
	sql := `DECLARE @t TABLE (id INT PRIMARY KEY, name NVARCHAR(50) NOT NULL)`
	if _, err := Parse(sql); err != nil {
		t.Fatalf("parse failed: %v", err)
	}
}
