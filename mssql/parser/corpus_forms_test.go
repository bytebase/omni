package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestCorpusFormsTree pins the tree omni builds for the T-SQL forms of
// BYT-10538, which SQL Server accepts and omni used to refuse; whether each
// parses is SQL Server's answer in testdata/compat.tsv (issue:BYT-10538).
// Each case finds its node and checks its fields and the text its Loc covers.
func TestCorpusFormsTree(t *testing.T) {
	first := func(t *testing.T, sql string, match func(ast.Node) bool) ast.Node {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", sql, err)
		}
		var found ast.Node
		ast.Inspect(list, func(n ast.Node) bool {
			if found == nil && n != nil && match(n) {
				found = n
			}
			return found == nil
		})
		if found == nil {
			t.Fatalf("Parse(%q): no node matched", sql)
		}
		return found
	}
	covers := func(t *testing.T, sql string, n ast.Node, want string) {
		t.Helper()
		loc := ast.NodeLoc(n)
		if loc.Start < 0 || loc.End > len(sql) || loc.Start > loc.End {
			t.Fatalf("%T Loc %v out of %q", n, loc, sql)
		}
		if got := sql[loc.Start:loc.End]; got != want {
			t.Errorf("%T covers %q, want %q", n, got, want)
		}
	}
	is := func(n ast.Node) func(ast.Node) bool {
		return func(m ast.Node) bool {
			switch n.(type) {
			case *ast.MethodCallExpr:
				_, ok := m.(*ast.MethodCallExpr)
				return ok
			case *ast.SelectAssign:
				_, ok := m.(*ast.SelectAssign)
				return ok
			case *ast.ColumnDef:
				_, ok := m.(*ast.ColumnDef)
				return ok
			case *ast.ConstraintDef:
				_, ok := m.(*ast.ConstraintDef)
				return ok
			case *ast.SetExpr:
				_, ok := m.(*ast.SetExpr)
				return ok
			case *ast.ForClause:
				_, ok := m.(*ast.ForClause)
				return ok
			case *ast.BeginEndStmt:
				_, ok := m.(*ast.BeginEndStmt)
				return ok
			}
			return false
		}
	}

	t.Run("method on a subquery", func(t *testing.T) {
		sql := "SELECT (SELECT 'x' FOR XML PATH(''), TYPE).value('.', 'varchar(10)')"
		mc := first(t, sql, is(&ast.MethodCallExpr{})).(*ast.MethodCallExpr)
		if _, ok := mc.Expr.(*ast.SubqueryExpr); !ok || mc.Method != "value" || mc.Args == nil || len(mc.Args.Items) != 2 {
			t.Fatalf("got %s", ast.NodeToString(mc))
		}
		covers(t, sql, mc, "(SELECT 'x' FOR XML PATH(''), TYPE).value('.', 'varchar(10)')")
	})
	t.Run("methods and a property on a variable", func(t *testing.T) {
		sql := "SELECT @x.query('.').value('.', 'int'), @x.a"
		mc := first(t, sql, is(&ast.MethodCallExpr{})).(*ast.MethodCallExpr)
		inner, ok := mc.Expr.(*ast.MethodCallExpr)
		if !ok || mc.Method != "value" || inner.Method != "query" {
			t.Fatalf("got %s", ast.NodeToString(mc))
		}
		if v, ok := inner.Expr.(*ast.VariableRef); !ok || v.Name != "@x" {
			t.Fatalf("got %s", ast.NodeToString(inner))
		}
		covers(t, sql, mc, "@x.query('.').value('.', 'int')")
		covers(t, sql, inner, "@x.query('.')")
		var prop *ast.MethodCallExpr
		list, _ := Parse(sql)
		ast.Inspect(list, func(n ast.Node) bool {
			if m, ok := n.(*ast.MethodCallExpr); ok && m.Method == "a" {
				prop = m
			}
			return true
		})
		if prop == nil || prop.Args != nil {
			t.Fatalf("property @x.a: got %v", prop)
		}
		covers(t, sql, prop, "@x.a")
	})
	t.Run("FOR XML PATH with a Unicode string", func(t *testing.T) {
		sql := "SELECT 1 FOR XML PATH(N'r'), ROOT(N'x')"
		fc := first(t, sql, is(&ast.ForClause{})).(*ast.ForClause)
		if fc.ElementName != "r" || !fc.Root || fc.RootName != "x" {
			t.Fatalf("got %s", ast.NodeToString(fc))
		}
	})
	t.Run("a column's INDEX", func(t *testing.T) {
		sql := "DECLARE @t TABLE (a INT INDEX ix_a CLUSTERED, b INT)"
		col := first(t, sql, is(&ast.ColumnDef{})).(*ast.ColumnDef)
		if col.Name != "a" || col.Indexes == nil || len(col.Indexes.Items) != 1 {
			t.Fatalf("got %s", ast.NodeToString(col))
		}
		idx := col.Indexes.Items[0].(*ast.InlineIndexDef)
		if idx.Name != "ix_a" || idx.Clustered == nil || !*idx.Clustered {
			t.Fatalf("got %s", ast.NodeToString(col))
		}
		covers(t, sql, col, "a INT INDEX ix_a CLUSTERED")
		covers(t, sql, idx, "INDEX ix_a CLUSTERED")
		two := "DECLARE @t TABLE (a INT INDEX ix_a INDEX ix_b)"
		col = first(t, two, is(&ast.ColumnDef{})).(*ast.ColumnDef)
		if col.Indexes == nil || len(col.Indexes.Items) != 2 || col.Indexes.Items[1].(*ast.InlineIndexDef).Name != "ix_b" {
			t.Fatalf("got %s", ast.NodeToString(col))
		}
	})
	t.Run("a HASH index", func(t *testing.T) {
		sql := "CREATE TABLE a (a INT INDEX ix NONCLUSTERED HASH WITH (BUCKET_COUNT = 8))"
		col := first(t, sql, is(&ast.ColumnDef{})).(*ast.ColumnDef)
		if col.Indexes == nil || len(col.Indexes.Items) != 1 {
			t.Fatalf("got %s", ast.NodeToString(col))
		}
		idx := col.Indexes.Items[0].(*ast.InlineIndexDef)
		if !idx.Hash || idx.Options == nil {
			t.Fatalf("got %s", ast.NodeToString(col))
		}
		if s, ok := idx.Options.Items[0].(*ast.String); !ok || s.Str != "BUCKET_COUNT=8" {
			t.Fatalf("options: %s", ast.NodeToString(idx.Options))
		}
	})
	t.Run("a column's FOREIGN KEY", func(t *testing.T) {
		sql := "CREATE TABLE c (id INT NOT NULL CONSTRAINT fk FOREIGN KEY (id) REFERENCES dbo.p (id))"
		cd := first(t, sql, is(&ast.ConstraintDef{})).(*ast.ConstraintDef)
		if cd.Type != ast.ConstraintForeignKey || cd.Name != "fk" || cd.Columns == nil || len(cd.Columns.Items) != 1 ||
			cd.RefTable == nil || cd.RefTable.Schema != "dbo" || cd.RefTable.Object != "p" || cd.RefColumns == nil {
			t.Fatalf("got %s", ast.NodeToString(cd))
		}
		covers(t, sql, cd, "FOREIGN KEY (id) REFERENCES dbo.p (id)")
	})
	t.Run("a qualified SET target", func(t *testing.T) {
		sql := "UPDATE dbo.t SET srv.db.dbo.t.c += 1"
		se := first(t, sql, is(&ast.SetExpr{})).(*ast.SetExpr)
		c := se.Column
		if c == nil || c.Server != "srv" || c.Database != "db" || c.Schema != "dbo" || c.Table != "t" || c.Column != "c" || se.Operator != "+=" {
			t.Fatalf("got %s", ast.NodeToString(se))
		}
		covers(t, sql, c, "srv.db.dbo.t.c")
	})
	t.Run("a mutator", func(t *testing.T) {
		sql := "UPDATE t SET c.m(1, 2)"
		se := first(t, sql, is(&ast.SetExpr{})).(*ast.SetExpr)
		fc, ok := se.Value.(*ast.FuncCallExpr)
		if !se.Mutator || se.WriteMethod || se.Column == nil || se.Column.Column != "c" || se.Column.Table != "" ||
			!ok || fc.Name.Object != "m" || fc.Args == nil || len(fc.Args.Items) != 2 {
			t.Fatalf("got %s", ast.NodeToString(se))
		}
		covers(t, sql, fc, "m(1, 2)")
		spaced := "UPDATE t SET [c] . m(1)"
		se = first(t, spaced, is(&ast.SetExpr{})).(*ast.SetExpr)
		covers(t, spaced, se.Column, "[c]")
	})
	t.Run(".WRITE", func(t *testing.T) {
		sql := "UPDATE t SET c.WRITE(N'x', 0, 1)"
		se := first(t, sql, is(&ast.SetExpr{})).(*ast.SetExpr)
		if !se.WriteMethod || se.Mutator || se.Column.Column != "c" {
			t.Fatalf("got %s", ast.NodeToString(se))
		}
	})
	t.Run("a variable and a column", func(t *testing.T) {
		sql := "UPDATE t SET @v = dbo.t.c += 1"
		se := first(t, sql, is(&ast.SetExpr{})).(*ast.SetExpr)
		if se.Variable != "@v" || se.VarColumn == nil || se.VarColumn.Schema != "dbo" || se.VarColumn.Table != "t" ||
			se.VarColumn.Column != "c" || se.Operator != "+=" {
			t.Fatalf("got %s", ast.NodeToString(se))
		}
		covers(t, sql, se.VarColumn, "dbo.t.c")
	})
	t.Run("SELECT @c +=", func(t *testing.T) {
		sql := "SELECT @c += N'x' FROM sys.objects"
		sa := first(t, sql, is(&ast.SelectAssign{})).(*ast.SelectAssign)
		if sa.Variable != "@c" || sa.Operator != "+=" {
			t.Fatalf("got %s", ast.NodeToString(sa))
		}
		covers(t, sql, sa, "@c += N'x'")
	})
	t.Run("END CONVERSATION in a block", func(t *testing.T) {
		sql := "BEGIN SELECT 1; END CONVERSATION @h; END"
		be := first(t, sql, is(&ast.BeginEndStmt{})).(*ast.BeginEndStmt)
		if be.Stmts == nil || len(be.Stmts.Items) != 2 {
			t.Fatalf("got %s", ast.NodeToString(be))
		}
		end, ok := be.Stmts.Items[1].(*ast.ServiceBrokerStmt)
		if !ok || end.Action != "END" || end.ObjectType != "CONVERSATION" {
			t.Fatalf("second statement: %s", ast.NodeToString(be.Stmts.Items[1]))
		}
		covers(t, sql, end, "END CONVERSATION @h")
		covers(t, sql, be, sql)
	})
	t.Run("NBSP between tokens", func(t *testing.T) {
		sql := "SELECT a FROM t"
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		sel := list.Items[0].(*ast.SelectStmt)
		ref := sel.TargetList.Items[0].(*ast.ResTarget).Val.(*ast.ColumnRef)
		if ref.Column != "a" {
			t.Fatalf("column %q", ref.Column)
		}
		covers(t, sql, ref, "a")
	})
}
