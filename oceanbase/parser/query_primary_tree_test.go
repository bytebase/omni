package parser

import (
	"errors"
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

	// The clauses after a set operation's last operand are the set operation's, in the
	// order they were written.
	for _, sql := range []string{
		"VALUES ROW(1) UNION VALUES ROW(2) ORDER BY 1 FOR UPDATE INTO @x",
		"SELECT 1 UNION SELECT 2 ORDER BY 1 FOR UPDATE INTO @x",
		"SELECT 1 FROM t UNION SELECT 2 FROM t ORDER BY 1 FOR UPDATE INTO @x",
	} {
		sel := queryPrimaryStatement(t, sql).(*ast.SelectStmt)
		if sel.SetOp != ast.SetOpUnion || len(sel.OrderBy) != 1 || sel.ForUpdate == nil || sel.Into == nil {
			t.Errorf("%s: %+v", sql, sel)
		}
		if r := sel.Right; len(r.OrderBy) != 0 || r.ForUpdate != nil || r.Into != nil {
			t.Errorf("%s: the last operand kept clauses: %+v", sql, r)
		}
		if sel.Loc.End != len(sql) {
			t.Errorf("%s: Loc.End = %d, want %d", sql, sel.Loc.End, len(sql))
		}
	}
	// Inside parentheses they stay with the operand they scope.
	sel := queryPrimaryStatement(t, "SELECT 1 UNION (SELECT 2 FOR UPDATE)").(*ast.SelectStmt)
	if sel.ForUpdate != nil || sel.Right.ParenSource == nil || sel.Right.ParenSource.ForUpdate == nil {
		t.Errorf("the parenthesized operand lost its clause: %+v", sel)
	}
}

// TestIntoHoldsWhatWasWritten: an empty option is an option, the last of a repeated one
// wins, adjacent strings read as one, and a name that is also a keyword is a variable.
func TestIntoHoldsWhatWasWritten(t *testing.T) {
	into := func(sql string) *ast.IntoClause {
		t.Helper()
		return queryPrimaryStatement(t, sql).(*ast.SelectStmt).Into
	}
	i := into("VALUES ROW(1) INTO OUTFILE 'f' FIELDS TERMINATED BY '' ESCAPED BY '' LINES STARTING BY ''")
	if !i.HasOutfile || !i.HasFieldsTerminatedBy || !i.HasFieldsEscapedBy || !i.HasLinesStartingBy || i.HasFieldsEnclosedBy || i.HasLinesTerminatedBy {
		t.Errorf("%+v", i)
	}
	i = into("VALUES ROW(1) INTO OUTFILE 'a' 'b' FIELDS TERMINATED BY ',' TERMINATED BY ';' ENCLOSED BY '1' OPTIONALLY ENCLOSED BY '2'")
	if i.Outfile != "ab" || i.FieldsTerminatedBy != ";" || i.FieldsEnclosedBy != "2" || !i.FieldsOptionalEncl {
		t.Errorf("%+v", i)
	}
	i = into("VALUES ROW(1) INTO OUTFILE 'f' FIELDS OPTIONALLY ENCLOSED BY '1' ENCLOSED BY '2'")
	if i.FieldsEnclosedBy != "2" || i.FieldsOptionalEncl {
		t.Errorf("the last ENCLOSED BY should win, with its own OPTIONALLY: %+v", i)
	}
	i = into("VALUES ROW(1) INTO OUTFILE ''")
	if !i.HasOutfile || i.Outfile != "" {
		t.Errorf("%+v", i)
	}
	i = into("VALUES ROW(1) INTO DUMPFILE ''")
	if !i.HasDumpfile || i.HasOutfile || i.Dumpfile != "" {
		t.Errorf("%+v", i)
	}
	i = into("VALUES ROW(1) INTO dumpfile")
	if i.HasDumpfile || len(i.Vars) != 1 || i.Vars[0].Name != "dumpfile" || !i.Vars[0].Bare {
		t.Errorf("%+v", i)
	}
}

// TestIntoRefusalsPoint: a refusal names the token at fault, or the end of the input.
func TestIntoRefusalsPoint(t *testing.T) {
	for _, c := range []struct {
		sql string
		at  int
	}{
		{"VALUES ROW(1) INTO", 18},
		{"VALUES ROW(1) INTO OUTFILE", 26},
		{"VALUES ROW(1) INTO @x,", 22},
		{"VALUES ROW(1) INTO @@autocommit", 19},
		{"VALUES ROW(1) INTO OUTFILE 'f' FIELDS", 37},
		{"VALUES ROW(1) INTO OUTFILE 'f' LINES", 36},
		{"VALUES ROW(1) INTO OUTFILE 'f' FIELDS TERMINATED ','", 49},
		{"VALUES ROW(1) INTO OUTFILE 'f' FIELDS TERMINATED BY", 51},
		{"VALUES ROW(1) INTO OUTFILE 'f' CHARACTER utf8mb4", 41},
	} {
		_, err := Parse(c.sql)
		var perr *ParseError
		if !errors.As(err, &perr) {
			t.Errorf("Parse(%q) = %v, want a ParseError", c.sql, err)
			continue
		}
		if perr.Position != c.at {
			t.Errorf("Parse(%q) points at %d, want %d (%v)", c.sql, perr.Position, c.at, err)
		}
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
