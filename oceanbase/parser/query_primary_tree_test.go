package parser

import (
	"testing"

	"github.com/bytebase/omni/oceanbase/ast"
)

// The tree of the forms query_primary_forms_test.go holds to the server.

func queryPrimaryStatement(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q) read %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

// TestValuesStatementIsAQuery: a VALUES statement reads as a SelectStmt over a
// ValuesStmt, with its clauses on the SelectStmt, and a trailing INTO is held by the
// query expression and covered by its Loc.
func TestValuesStatementIsAQuery(t *testing.T) {
	for _, c := range []struct {
		sql                      string
		order, limit, lock, into bool
	}{
		{"VALUES ROW(1)", false, false, false, false},
		{"VALUES ROW(1) ORDER BY 1", true, false, false, false},
		{"VALUES ROW(1) LIMIT 1", false, true, false, false},
		{"VALUES ROW(1) FOR UPDATE", false, false, true, false},
		{"VALUES ROW(1) ORDER BY 1 LIMIT 1 LOCK IN SHARE MODE", true, true, true, false},
		{"VALUES ROW(1) INTO @x", false, false, false, true},
		{"VALUES ROW(1) ORDER BY 1 FOR UPDATE INTO @x", true, false, true, true},
		{"(SELECT 1) LIMIT 1 FOR UPDATE INTO @x", false, true, true, true},
		{"(VALUES ROW(1)) ORDER BY 1 INTO @x", true, false, false, true},
		{"(SELECT 1) UNION (SELECT 2) INTO @x", false, false, false, true},
		{"(SELECT 1) UNION (SELECT 2) FOR UPDATE", false, false, true, false},
	} {
		t.Run(c.sql, func(t *testing.T) {
			sel, ok := queryPrimaryStatement(t, c.sql).(*ast.SelectStmt)
			if !ok {
				t.Fatalf("not a SelectStmt")
			}
			if got := len(sel.OrderBy) > 0; got != c.order {
				t.Errorf("ORDER BY held = %v, want %v", got, c.order)
			}
			if got := sel.Limit != nil; got != c.limit {
				t.Errorf("LIMIT held = %v, want %v", got, c.limit)
			}
			if got := sel.ForUpdate != nil; got != c.lock {
				t.Errorf("locking clause held = %v, want %v", got, c.lock)
			}
			if got := sel.Into != nil; got != c.into {
				t.Errorf("INTO held = %v, want %v", got, c.into)
			}
			if sel.Loc.Start != 0 || sel.Loc.End != len(c.sql) {
				t.Errorf("Loc = %+v, want the whole statement of %d bytes", sel.Loc, len(c.sql))
			}
		})
	}

	// A locking clause after a set operation's last operand stays on it, as it does
	// for a SELECT.
	sel := queryPrimaryStatement(t, "VALUES ROW(1) UNION VALUES ROW(2) FOR UPDATE").(*ast.SelectStmt)
	if sel.SetOp != ast.SetOpUnion || sel.Left.ValuesSource == nil || sel.Right.ValuesSource == nil || sel.Right.ForUpdate == nil {
		t.Errorf("%+v", sel)
	}
	if sel.Left.ForUpdate != nil || sel.ForUpdate != nil {
		t.Errorf("the locking clause is on the wrong node: %+v", sel)
	}
}

// TestValuesSourcesHoldTheirQuery: CREATE TABLE and EXPLAIN keep the query inside the
// statement rather than leaving it to read as a second one.
func TestValuesSourcesHoldTheirQuery(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE n VALUES ROW(1)",
		"CREATE TABLE n AS VALUES ROW(1)",
		"CREATE TABLE n (VALUES ROW(1))",
		"CREATE TABLE n (a INT) VALUES ROW(1)",
		"CREATE TABLE n (a INT) (VALUES ROW(1))",
		"CREATE TABLE n VALUES ROW(1) UNION SELECT 1",
		"CREATE TABLE n ENGINE = InnoDB VALUES ROW(1)",
	} {
		stmt, ok := queryPrimaryStatement(t, sql).(*ast.CreateTableStmt)
		if !ok || stmt.Select == nil {
			t.Errorf("%s: %+v", sql, stmt)
			continue
		}
		if stmt.Loc.End != len(sql) {
			t.Errorf("%s: Loc.End = %d, want %d", sql, stmt.Loc.End, len(sql))
		}
	}
	explain := queryPrimaryStatement(t, "EXPLAIN VALUES ROW(1) UNION VALUES ROW(2)").(*ast.ExplainStmt)
	if sel, ok := explain.Stmt.(*ast.SelectStmt); !ok || sel.SetOp != ast.SetOpUnion {
		t.Errorf("EXPLAIN target = %+v", explain.Stmt)
	}
	insert := queryPrimaryStatement(t, "INSERT INTO t VALUES ROW(1) FOR UPDATE").(*ast.InsertStmt)
	if insert.Values != nil || insert.Select == nil || insert.Select.ForUpdate == nil {
		t.Errorf("INSERT source = Values %v, Select %+v", insert.Values, insert.Select)
	}
}
