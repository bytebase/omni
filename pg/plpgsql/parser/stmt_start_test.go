package parser

import (
	"testing"

	"github.com/bytebase/omni/pg/plpgsql/ast"
)

// TestAssignmentAtStatementStart holds the parser to PostgreSQL's regression
// suite (src/test/regress/sql/plpgsql.sql): a word followed by := at the
// start of a statement is an assignment even when it is spelled like a
// statement keyword, and $n is a target.
func TestAssignmentAtStatementStart(t *testing.T) {
	for body, target := range map[string]string{
		`DECLARE comment int := 21; BEGIN comment := comment * 2; RETURN comment; END`: "comment",
		`DECLARE return int := 42; BEGIN return := return + 1; RETURN return; END`:     "return",
		`BEGIN $1 := -1; RETURN NEXT; END`:                                             "$1",
		`BEGIN a[1] := 2; END`:                                                         "a[1]",
	} {
		block := parseOK(t, body)
		a, ok := block.Body[0].(*ast.PLAssign)
		if !ok {
			t.Errorf("%s: first statement is %T, want an assignment", body, block.Body[0])
			continue
		}
		if a.Target != target {
			t.Errorf("%s: target %q, want %q", body, a.Target, target)
		}
	}
}

func TestForEachTargetList(t *testing.T) {
	block := parseOK(t, `BEGIN FOREACH x, y IN ARRAY $1 LOOP RAISE NOTICE 'x = %, y = %', x, y; END LOOP; END`)
	f, ok := block.Body[0].(*ast.PLForEachA)
	if !ok {
		t.Fatalf("expected PLForEachA, got %T", block.Body[0])
	}
	if f.Var != "x, y" {
		t.Errorf("var %q, want %q", f.Var, "x, y")
	}
}
