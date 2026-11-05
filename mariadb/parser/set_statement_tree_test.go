package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mariadb/ast"
)

// parseSetStatement parses one SET STATEMENT.
func parseSetStatement(t *testing.T, sql string) *nodes.SetStatementStmt {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	stmt, ok := list.Items[0].(*nodes.SetStatementStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.SetStatementStmt", sql, list.Items[0])
	}
	return stmt
}

// TestSetStatementTrees pins the options and the statement each form holds.
func TestSetStatementTrees(t *testing.T) {
	for _, c := range []struct {
		sql, options, inner string
	}{
		{"set statement a = 1 for select 1", "a", "*ast.SelectStmt"},
		{"set statement a = 1, b = 2 for select 1", "a,b", "*ast.SelectStmt"},
		{"set statement a := 1 for select 1", "a", "*ast.SelectStmt"},
		{"set statement x.a = 1 for select 1", "x.a", "*ast.SelectStmt"},
		{"set statement default.a = 1 for select 1", "default.a", "*ast.SelectStmt"},
		{"set statement names = 1, password = 2 for select 1", "names,password", "*ast.SelectStmt"},
		{"SET STATEMENT a = 1 FOR DROP TABLE t", "a", "*ast.DropTableStmt"},
		{"set statement a = 1 for set @b = 2", "a", "*ast.SetStmt"},
		{"set statement a = 1 for set transaction read only", "a", "*ast.SetTransactionStmt"},
		{"set statement a = 1 for begin", "a", "*ast.BeginStmt"},
		{"set statement a = 1 for set statement b = 2 for select 1", "a", "*ast.SetStatementStmt"},
		{"set statement a = 1 for if 1 then select 1; end if", "a", "*ast.IfStmt"},
		{"set statement a = 1 for case 1 when 1 then select 1; end case", "a", "*ast.CaseStmtNode"},
		{"set statement a = 1 for while 1 do select 1; end while", "a", "*ast.WhileStmt"},
		{"set statement a = 1 for repeat select 1; until 1 end repeat", "a", "*ast.RepeatStmt"},
		{"set statement a = 1 for loop select 1; end loop", "a", "*ast.LoopStmt"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			stmt := parseSetStatement(t, c.sql)
			var options []string
			for _, o := range stmt.Options {
				name := o.Column.Column
				if o.Column.Table != "" {
					name = o.Column.Table + "." + name
				}
				options = append(options, name)
				if o.Scope != "" {
					t.Errorf("option %s has scope %q", name, o.Scope)
				}
			}
			if got := strings.Join(options, ","); got != c.options {
				t.Errorf("options = %q, want %q", got, c.options)
			}
			if got := fmt.Sprintf("%T", stmt.Stmt); got != c.inner {
				t.Errorf("statement = %s, want %s", got, c.inner)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("Loc = %+v, want {0 %d}", stmt.Loc, len(c.sql))
			}
		})
	}
}

// TestSetStatementNotAKeyword checks STATEMENT stays a variable name where `=`, `:=` or `.`
// follows it, and that a quoted one starts no form.
func TestSetStatementNotAKeyword(t *testing.T) {
	for _, c := range []struct{ sql, table, column string }{
		{"set statement = 1", "", "statement"},
		{"set statement := 1", "", "statement"},
		{"set statement.a = 1", "statement", "a"},
		{"set `statement` = 1", "", "statement"},
		{"set global statement = 1", "", "statement"},
	} {
		list, err := Parse(c.sql)
		if err != nil {
			t.Fatalf("Parse(%q): %v", c.sql, err)
		}
		set, ok := list.Items[0].(*nodes.SetStmt)
		if !ok {
			t.Fatalf("%q: got %T, want *ast.SetStmt", c.sql, list.Items[0])
		}
		col := set.Assignments[0].Column
		if col.Table != c.table || col.Column != c.column {
			t.Errorf("%q: table=%q column=%q", c.sql, col.Table, col.Column)
		}
	}
	if _, err := Parse("set `statement` a = 1 for select 1"); err == nil {
		t.Error("a quoted STATEMENT started SET STATEMENT")
	}
}

// TestSetStatementWalk checks a consumer that walks the tree reaches the wrapped statement and
// the options' values, and that the node dump shows both.
func TestSetStatementWalk(t *testing.T) {
	stmt := parseSetStatement(t, "set statement a = (select 1) for drop table t")
	var subqueries, drops int
	nodes.Inspect(stmt, func(n nodes.Node) bool {
		switch n.(type) {
		case *nodes.SubqueryExpr:
			subqueries++
		case *nodes.DropTableStmt:
			drops++
		}
		return true
	})
	if subqueries != 1 || drops != 1 {
		t.Errorf("walk reached %d subqueries and %d drops, want one of each", subqueries, drops)
	}
	dump := nodes.NodeToString(stmt)
	if !strings.Contains(dump, "{SET_STATEMENT") || !strings.Contains(dump, ":stmt {DROP_TABLE") {
		t.Errorf("NodeToString = %s", dump)
	}
}

// TestSetStatementErrorPosition checks each refused form reports the token it stopped at, an
// empty string meaning the end of the text.
func TestSetStatementErrorPosition(t *testing.T) {
	for _, tc := range []struct{ sql, at string }{
		{"set statement", ""},
		{"set statement for select 1", "for select 1"},
		{"set statement a = 1", ""},
		{"set statement a = 1 for", ""},
		{"set statement a = 1,", ""},
		{"set statement a = 1, for select 1", "for select 1"},
		{"set statement a = 1 select 1", "select 1"},
		{"set statement @a = 1 for select 1", "@a = 1 for select 1"},
		{"set statement @@a = 1 for select 1", "@@a = 1 for select 1"},
		{"set statement global a = 1 for select 1", "a = 1 for select 1"},
		{"set statement names utf8mb4 for select 1", "utf8mb4 for select 1"},
		{"set statement a = 1 for leave l", "leave l"},
		{"set statement a = 1 for declare b int", "declare b int"},
		{"set statement a = 1 for l: loop select 1; end loop", "l: loop select 1; end loop"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("err = %v, want a ParseError", err)
			}
			if got := tc.sql[pe.Position:]; got != tc.at {
				t.Errorf("stopped at %q, want %q", got, tc.at)
			}
		})
	}
}

// TestSetStatementTruncated cuts the longest form at every byte. Each prefix either parses or
// returns a ParseError holding an offset inside the text, and none panics.
func TestSetStatementTruncated(t *testing.T) {
	const sql = "set statement a = 1, b.c = default for set statement d = on for if 1 then select 1; else drop table t; end if"
	for i := 0; i <= len(sql); i++ {
		_, err := Parse(sql[:i])
		if err == nil {
			continue
		}
		var pe *ParseError
		if !errors.As(err, &pe) {
			t.Fatalf("Parse(%q): %T is not a ParseError", sql[:i], err)
		}
		if pe.Position < 0 || pe.Position > i {
			t.Errorf("Parse(%q): error position %d is outside the text", sql[:i], pe.Position)
		}
	}
}

// TestSetStatementNesting checks a deep chain of wrappers fails with the nesting error rather
// than overflowing the stack.
func TestSetStatementNesting(t *testing.T) {
	sql := strings.Repeat("set statement a = 1 for ", maxNestingDepth+1) + "select 1"
	_, err := Parse(sql)
	if !isNestingError(err) {
		t.Errorf("err = %v, want the nesting error", briefErr(err))
	}
	if _, err := Parse(strings.Repeat("set statement a = 1 for ", 100) + "select 1"); err != nil {
		t.Errorf("100 wrappers: %v", briefErr(err))
	}
}

// TestSetStatementSplit checks the statement stays one segment, with what it wraps.
func TestSetStatementSplit(t *testing.T) {
	for _, sql := range []string{
		"set statement a = 1 for select 1;",
		"set statement a = 1 for if 1 then select 1; end if;",
		"SET STATEMENT a = 1 FOR /* c */ IF 1 THEN SELECT 1; END IF;",
		"set statement a = 1 for while 1 do select 1; end while;",
		"set statement a = 1 for repeat select 1; until 1 end repeat;",
		"set statement a = 1 for loop select 1; end loop;",
		"set statement a = 1 for case 1 when 1 then select 1; end case;",
		"set statement a = 1 for create procedure p() begin select 1; select 2; end;",
	} {
		t.Run(sql, func(t *testing.T) {
			if got := countNonEmptySegments(Split(sql + "\nselect 2;")); got != 2 {
				t.Errorf("%d segments, want the statement and the one after it", got)
			}
		})
	}
	// A word after FOR elsewhere stays what it was.
	for _, sql := range []string{
		"select 1 for update;\nselect if(1, 2, 3);",
		"create trigger g before insert on t for each row set @a = if(1, 2, 3);\nselect 1;",
	} {
		if got := countNonEmptySegments(Split(sql)); got != 2 {
			t.Errorf("%q: %d segments, want 2", sql, got)
		}
	}
}

func countNonEmptySegments(segs []Segment) int {
	n := 0
	for _, s := range segs {
		if !s.Empty() {
			n++
		}
	}
	return n
}
