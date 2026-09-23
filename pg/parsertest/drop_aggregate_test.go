package parsertest

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// aggregateObjects returns the ObjectWithArgs entries of a parsed DROP AGGREGATE.
func aggregateObjects(t *testing.T, sql string) []*nodes.ObjectWithArgs {
	t.Helper()
	result, err := parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) error: %v", sql, err)
	}
	stmt, ok := result.Items[0].(*nodes.DropStmt)
	if !ok {
		t.Fatalf("expected *nodes.DropStmt, got %T", result.Items[0])
	}
	if stmt.RemoveType != int(nodes.OBJECT_AGGREGATE) {
		t.Fatalf("expected OBJECT_AGGREGATE, got %d", stmt.RemoveType)
	}
	var objects []*nodes.ObjectWithArgs
	for _, item := range stmt.Objects.Items {
		owa, ok := item.(*nodes.ObjectWithArgs)
		if !ok {
			t.Fatalf("expected *nodes.ObjectWithArgs, got %T", item)
		}
		objects = append(objects, owa)
	}
	return objects
}

// lastName returns the unqualified name of a TypeName, e.g. "int4" for int.
func lastName(t *testing.T, tn *nodes.TypeName) string {
	t.Helper()
	if tn == nil || tn.Names == nil || len(tn.Names.Items) == 0 {
		t.Fatalf("TypeName has no names: %#v", tn)
	}
	s, ok := tn.Names.Items[len(tn.Names.Items)-1].(*nodes.String)
	if !ok {
		t.Fatalf("expected *nodes.String, got %T", tn.Names.Items[len(tn.Names.Items)-1])
	}
	return s.Str
}

// TestDropAggregateEmptyArgs tests: DROP AGGREGATE b()
//
// gram.y aggr_args is '(' '*' ')' or a non-empty aggr_args_list, so "b()" is a
// syntax error rather than a no-argument aggregate.
func TestDropAggregateEmptyArgs(t *testing.T) {
	_, err := parse("DROP AGGREGATE b()")
	if err == nil {
		t.Fatal("expected a syntax error for DROP AGGREGATE b(), got nil")
	}
	if !strings.Contains(err.Error(), `syntax error at or near ")"`) {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestDropAggregateStar tests: DROP AGGREGATE a(*)
//
// "(*)" is the no-argument form and yields NIL objargs.
func TestDropAggregateStar(t *testing.T) {
	objects := aggregateObjects(t, "DROP AGGREGATE a(*)")
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}
	if objects[0].Objargs != nil {
		t.Errorf("expected nil Objargs for a(*), got %v", nodes.NodeToString(objects[0].Objargs))
	}
	if objects[0].ArgsUnspecified {
		t.Error("expected ArgsUnspecified to be false")
	}
}

// TestDropAggregateTypedArgs tests: DROP AGGREGATE c(int, text)
func TestDropAggregateTypedArgs(t *testing.T) {
	objects := aggregateObjects(t, "DROP AGGREGATE c(int, text)")
	if len(objects) != 1 {
		t.Fatalf("expected 1 object, got %d", len(objects))
	}
	args := objects[0].Objargs
	if args == nil || len(args.Items) != 2 {
		t.Fatalf("expected 2 Objargs, got %v", nodes.NodeToString(args))
	}
	want := []string{"int4", "text"}
	for i, item := range args.Items {
		tn, ok := item.(*nodes.TypeName)
		if !ok || tn == nil {
			t.Fatalf("Objargs[%d]: expected non-nil *nodes.TypeName, got %T (%v)", i, item, item)
		}
		if got := lastName(t, tn); got != want[i] {
			t.Errorf("Objargs[%d]: expected %q, got %q", i, want[i], got)
		}
	}
	// A well-formed tree must be walkable without a panic.
	nodes.Inspect(objects[0], func(n nodes.Node) bool { return true })
}

// TestAggregateEmptyArgsRejected covers every statement that takes
// aggregate_with_argtypes: each must reject the empty "()" argument list.
func TestAggregateEmptyArgsRejected(t *testing.T) {
	for _, sql := range []string{
		"DROP AGGREGATE IF EXISTS b()",
		"DROP AGGREGATE a(int), b()",
		"ALTER AGGREGATE b() RENAME TO c",
		"ALTER AGGREGATE b() OWNER TO u",
		"ALTER AGGREGATE b() SET SCHEMA s",
		"COMMENT ON AGGREGATE b() IS 'x'",
		"SECURITY LABEL ON AGGREGATE b() IS 'x'",
		"ALTER EXTENSION e ADD AGGREGATE b()",
		"ALTER EXTENSION e DROP AGGREGATE b()",
	} {
		_, err := parse(sql)
		if err == nil {
			t.Errorf("Parse(%q): expected a syntax error, got nil", sql)
			continue
		}
		if !strings.Contains(err.Error(), `syntax error at or near ")"`) {
			t.Errorf("Parse(%q): unexpected error: %v", sql, err)
		}
	}
}

// TestAggregateStarAndTypedArgsAccepted is the positive counterpart of
// TestAggregateEmptyArgsRejected for the same statements.
func TestAggregateStarAndTypedArgsAccepted(t *testing.T) {
	for _, sql := range []string{
		"ALTER AGGREGATE a(*) RENAME TO c",
		"ALTER AGGREGATE c(int, text) OWNER TO u",
		"COMMENT ON AGGREGATE a(*) IS 'x'",
		"COMMENT ON AGGREGATE c(int, text) IS 'x'",
		"ALTER EXTENSION e ADD AGGREGATE a(*)",
		"ALTER EXTENSION e ADD AGGREGATE c(int, text)",
		"ALTER EXTENSION e DROP AGGREGATE o(int ORDER BY text)",
	} {
		result, err := parse(sql)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", sql, err)
			continue
		}
		nodes.Inspect(result, func(n nodes.Node) bool { return true })
	}
}

// TestCreateAggregateEmptyArgs tests: CREATE AGGREGATE b() (SFUNC = f, STYPE = int)
//
// Neither aggr_args nor old_aggr_definition accepts "()", so gram.y reports the
// ")" token.
func TestCreateAggregateEmptyArgs(t *testing.T) {
	_, err := parse("CREATE AGGREGATE b() (SFUNC = f, STYPE = int)")
	if err == nil {
		t.Fatal("expected a syntax error for CREATE AGGREGATE b() (...), got nil")
	}
	if !strings.Contains(err.Error(), `syntax error at or near ")"`) {
		t.Errorf("unexpected error: %v", err)
	}
}
