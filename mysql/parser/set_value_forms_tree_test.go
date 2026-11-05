package parser

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/mysql/ast"
)

// TestSetCharsetFormTrees pins the assignments each form stands for, wherever it sits in the
// list.
func TestSetCharsetFormTrees(t *testing.T) {
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"set names utf8mb4", []string{"NAMES=utf8mb4"}},
		{"set names 'utf8mb4'", []string{"NAMES=utf8mb4"}},
		{"set names binary", []string{"NAMES=binary"}},
		{"set names default", []string{"NAMES=DEFAULT"}},
		{"set names utf8mb4 collate utf8mb4_bin", []string{"NAMES=utf8mb4", "COLLATE=utf8mb4_bin"}},
		{"set names utf8mb4 collate 'utf8mb4_bin'", []string{"NAMES=utf8mb4", "COLLATE=utf8mb4_bin"}},
		{"set charset utf8mb4", []string{"CHARACTER SET=utf8mb4"}},
		{"set character set default", []string{"CHARACTER SET=DEFAULT"}},
		{"set char set binary", []string{"CHARACTER SET=binary"}},
		{"set char set 'utf8mb4'", []string{"CHARACTER SET=utf8mb4"}},
		{"set @a = 'x', names utf8mb4", []string{"@a=x", "NAMES=utf8mb4"}},
		{"set names utf8mb4, @a = 'x'", []string{"NAMES=utf8mb4", "@a=x"}},
		{"set @a = 'x', names utf8mb4 collate utf8mb4_bin, @b = 'y'", []string{"@a=x", "NAMES=utf8mb4", "COLLATE=utf8mb4_bin", "@b=y"}},
		{"set names utf8mb4, charset latin1", []string{"NAMES=utf8mb4", "CHARACTER SET=latin1"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			stmt := list.Items[0].(*nodes.SetStmt)
			var got []string
			for _, a := range stmt.Assignments {
				got = append(got, a.Column.Column+"="+a.Value.(*nodes.StringLit).Value)
			}
			if !slices.Equal(got, c.want) {
				t.Errorf("assignments = %q, want %q", got, c.want)
			}
			if stmt.Loc.Start != 0 || stmt.Loc.End != len(c.sql) {
				t.Errorf("statement Loc = %+v, want {0 %d}", stmt.Loc, len(c.sql))
			}
		})
	}
}

// TestSetValueTrees pins what a value's word turned into. ON stays a true BoolLit, as review
// and catalog read it.
func TestSetValueTrees(t *testing.T) {
	for _, c := range []struct{ sql, typ, want string }{
		{"set a = on", "*ast.BoolLit", "true"},
		{"set a = binary", "*ast.StringLit", "BINARY"},
		{"set @@a = binary", "*ast.StringLit", "BINARY"},
		{"set global a = binary", "*ast.StringLit", "BINARY"},
		{"set a = all", "*ast.StringLit", "ALL"},
		{"set a = row", "*ast.StringLit", "ROW"},
		{"set a = system", "*ast.StringLit", "SYSTEM"},
		{"set a = default", "*ast.DefaultExpr", ""},
		{"set a = default(b)", "*ast.DefaultExpr", "b"},
		{"set a = binary 'x'", "*ast.UnaryExpr", ""},
		{"set a = 1", "*ast.IntLit", "1"},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatal(err)
			}
			stmt := list.Items[0].(*nodes.SetStmt)
			val := stmt.Assignments[0].Value
			got := ""
			switch v := val.(type) {
			case *nodes.BoolLit:
				got = fmt.Sprint(v.Value)
			case *nodes.StringLit:
				got = v.Value
			case *nodes.DefaultExpr:
				got = v.Column
			case *nodes.IntLit:
				got = fmt.Sprint(v.Value)
			}
			if typ := fmt.Sprintf("%T", val); typ != c.typ || got != c.want {
				t.Errorf("value = %s %q, want %s %q", typ, got, c.typ, c.want)
			}
			if stmt.Assignments[0].Loc.End != len(c.sql) {
				t.Errorf("assignment Loc = %+v, want it to end at %d", stmt.Assignments[0].Loc, len(c.sql))
			}
		})
	}
}

// TestQuotedUserVariableTrees pins the name a quoted user variable reads as in each position,
// and that the user specs that share its text stay themselves.
func TestQuotedUserVariableTrees(t *testing.T) {
	refs := func(t *testing.T, sql string) (got []string) {
		t.Helper()
		list, err := Parse(sql)
		if err != nil {
			t.Fatal(err)
		}
		nodes.Inspect(list.Items[0], func(n nodes.Node) bool {
			if v, ok := n.(*nodes.VariableRef); ok {
				if text := sql[v.Loc.Start:v.Loc.End]; v.System || text[0] != '@' || !strings.ContainsAny(text[len(text)-1:], `'"`) {
					t.Errorf("%q: Loc covers %q", v.Name, text)
				}
				got = append(got, v.Name)
			}
			return true
		})
		return got
	}
	for _, c := range []struct {
		sql  string
		want []string
	}{
		{"select @'a'", []string{"a"}},
		{`select @"a"`, []string{"a"}},
		{"select @'a b'", []string{"a b"}},
		{"select @'a''b'", []string{"a'b"}},
		{`select @"a", @'b'`, []string{"a", "b"}},
		{"select @'a' := 1", []string{"a"}},
		{"select @'a' + 1", []string{"a"}},
		{"select 1 from t where a = @'b'", []string{"b"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			if got := refs(t, c.sql); !slices.Equal(got, c.want) {
				t.Errorf("variables = %q, want %q", got, c.want)
			}
		})
	}

	list, err := Parse("set @'a b' = 1, @\"c\" := 2")
	if err != nil {
		t.Fatal(err)
	}
	stmt := list.Items[0].(*nodes.SetStmt)
	if got := []string{stmt.Assignments[0].Column.Column, stmt.Assignments[1].Column.Column}; !slices.Equal(got, []string{"@a b", "@c"}) {
		t.Errorf("columns = %q", got)
	}

	list, err = Parse("select 1 into @'a', @\"b\"")
	if err != nil {
		t.Fatal(err)
	}
	into := list.Items[0].(*nodes.SelectStmt).Into
	if len(into.Vars) != 2 || into.Vars[0].Name != "a" || into.Vars[1].Name != "b" {
		t.Errorf("into = %+v", into.Vars)
	}

	list, err = Parse("create user 'a'@'b'")
	if err != nil {
		t.Fatal(err)
	}
	user := list.Items[0].(*nodes.CreateUserStmt).Users[0]
	if user.Name != "a" || user.Host != "b" {
		t.Errorf("user = %+v, want a at b", user)
	}
}

// TestSetValueFormsErrorPosition checks each refused form reports the token it stopped at, an
// empty string meaning the end of the text.
func TestSetValueFormsErrorPosition(t *testing.T) {
	for _, tc := range []struct{ sql, at string }{
		{"set names", ""},
		{"set names @a", "@a"},
		{"set names @'a'", "@'a'"},
		{"set charset @a", "@a"},
		{"set names utf8mb4 collate", ""},
		{"set names utf8mb4 collate @a", "@a"},
		{"set global names utf8mb4", "utf8mb4"},
		{"set @a = 1, names", ""},
		{"set @a = 1, global names utf8mb4", "utf8mb4"},
		{"set @a = on", "on"},
		{"set @a = default", "default"},
		{"set @a = binary", ""},
		{"set a = (default)", "default)"},
		{"set a = (on)", "on)"},
		{"set a = 1 + default", "default"},
		{"set a = default + 1", "+ 1"},
		{"set a = on + 1", "+ 1"},
		{"select on", "on"},
		{"set @'a", "'a"},
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

// TestSetValueFormsTruncated cuts the longest form at every byte. Each prefix either parses or
// returns a ParseError holding an offset inside the text, and none panics.
func TestSetValueFormsTruncated(t *testing.T) {
	const sql = "set @'a b' = 1, names 'utf8mb4' collate 'utf8mb4_bin', char set binary, sql_mode = binary, @@sql_mode = default, global sql_mode = on"
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
