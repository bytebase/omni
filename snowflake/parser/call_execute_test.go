package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func mustCall(t *testing.T, input string) *ast.CallStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.CallStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.CallStmt", input, node)
	}
	return stmt
}

func mustExecImm(t *testing.T, input string) *ast.ExecuteImmediateStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ExecuteImmediateStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ExecuteImmediateStmt", input, node)
	}
	return stmt
}

func mustExecTask(t *testing.T, input string) *ast.ExecuteTaskStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ExecuteTaskStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ExecuteTaskStmt", input, node)
	}
	return stmt
}

func mustExplain(t *testing.T, input string) *ast.ExplainStmt {
	t.Helper()
	node := mustParseOne(t, input)
	stmt, ok := node.(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("parse %q: got %T, want *ast.ExplainStmt", input, node)
	}
	return stmt
}

// assertValidLoc fails if loc has a negative or inverted span.
func assertValidLoc(t *testing.T, what string, loc ast.Loc) {
	t.Helper()
	if loc.Start < 0 || loc.End < 0 {
		t.Errorf("%s: Loc = %+v, want Start/End >= 0", what, loc)
	}
	if loc.End < loc.Start {
		t.Errorf("%s: Loc = %+v, want End >= Start", what, loc)
	}
}

// callArg returns the i-th argument of a CALL statement as a *ast.CallArg.
func callArg(t *testing.T, stmt *ast.CallStmt, i int) *ast.CallArg {
	t.Helper()
	if i >= len(stmt.Args) {
		t.Fatalf("arg index %d out of range (len=%d)", i, len(stmt.Args))
	}
	arg, ok := stmt.Args[i].(*ast.CallArg)
	if !ok {
		t.Fatalf("arg %d: got %T, want *ast.CallArg", i, stmt.Args[i])
	}
	return arg
}

// ---------------------------------------------------------------------------
// CALL
//   Syntax (legacy .g4 + docs): CALL <proc_name> ( [ <arg> [ , ... ] ] )
//   Args are expressions; named args (name => value) are docs-driven.
// ---------------------------------------------------------------------------

func TestCall_PositionalArgs(t *testing.T) {
	// Corpus call/example_01: CALL stproc1(5.14::FLOAT);
	stmt := mustCall(t, "CALL stproc1(5.14::FLOAT)")
	if stmt.Name.String() != "stproc1" {
		t.Errorf("Name = %q, want stproc1", stmt.Name.String())
	}
	if len(stmt.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(stmt.Args))
	}
	arg := callArg(t, stmt, 0)
	if !arg.Name.IsEmpty() {
		t.Errorf("arg 0 Name = %q, want empty (positional)", arg.Name.Normalize())
	}
	if _, ok := arg.Value.(*ast.CastExpr); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.CastExpr", arg.Value)
	}
	assertValidLoc(t, "CallStmt", stmt.Loc)
	assertValidLoc(t, "CallArg", arg.Loc)
}

func TestCall_ArithmeticArg(t *testing.T) {
	// Corpus call/example_02: CALL stproc1(2 * 5.14::FLOAT);
	stmt := mustCall(t, "CALL stproc1(2 * 5.14::FLOAT)")
	if len(stmt.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(stmt.Args))
	}
	if _, ok := callArg(t, stmt, 0).Value.(*ast.BinaryExpr); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.BinaryExpr", callArg(t, stmt, 0).Value)
	}
}

func TestCall_BareSubqueryArg(t *testing.T) {
	// Corpus call/example_03: CALL stproc1(SELECT COUNT(*) FROM stproc_test_table1);
	// DOCS DIVERGENCE: the docs permit a bare SELECT as a CALL argument; the
	// legacy expr_list rule does not. Docs win.
	stmt := mustCall(t, "CALL stproc1(SELECT COUNT(*) FROM stproc_test_table1)")
	if len(stmt.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(stmt.Args))
	}
	if _, ok := callArg(t, stmt, 0).Value.(*ast.SelectStmt); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.SelectStmt", callArg(t, stmt, 0).Value)
	}
}

func TestCall_StringAndNumberArgs(t *testing.T) {
	// Corpus call/example_04: CALL sv_proc1('Manitoba', 127.4);
	stmt := mustCall(t, "CALL sv_proc1('Manitoba', 127.4)")
	if len(stmt.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(stmt.Args))
	}
	if _, ok := callArg(t, stmt, 0).Value.(*ast.Literal); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.Literal", callArg(t, stmt, 0).Value)
	}
	if _, ok := callArg(t, stmt, 1).Value.(*ast.Literal); !ok {
		t.Errorf("arg 1 Value = %T, want *ast.Literal", callArg(t, stmt, 1).Value)
	}
}

func TestCall_NamedArgs(t *testing.T) {
	// Corpus call/example_05: CALL sv_proc1(province => 'Manitoba', amount => 127.4);
	stmt := mustCall(t, "CALL sv_proc1(province => 'Manitoba', amount => 127.4)")
	if len(stmt.Args) != 2 {
		t.Fatalf("len(Args) = %d, want 2", len(stmt.Args))
	}
	a0 := callArg(t, stmt, 0)
	if a0.Name.Normalize() != "PROVINCE" {
		t.Errorf("arg 0 Name = %q, want PROVINCE", a0.Name.Normalize())
	}
	if _, ok := a0.Value.(*ast.Literal); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.Literal", a0.Value)
	}
	a1 := callArg(t, stmt, 1)
	if a1.Name.Normalize() != "AMOUNT" {
		t.Errorf("arg 1 Name = %q, want AMOUNT", a1.Name.Normalize())
	}
	assertValidLoc(t, "named CallArg", a0.Loc)
}

func TestCall_SystemFunctionArg(t *testing.T) {
	// SYSTEM$... contains a mid-identifier '$' that lexes as one identifier token.
	stmt := mustCall(t, "CALL my_sp(SYSTEM$WAIT(1))")
	if len(stmt.Args) != 1 {
		t.Fatalf("len(Args) = %d, want 1", len(stmt.Args))
	}
	if _, ok := callArg(t, stmt, 0).Value.(*ast.FuncCallExpr); !ok {
		t.Errorf("arg 0 Value = %T, want *ast.FuncCallExpr", callArg(t, stmt, 0).Value)
	}
}

func TestCall_EmptyArgs(t *testing.T) {
	stmt := mustCall(t, "CALL my_sp()")
	if len(stmt.Args) != 0 {
		t.Errorf("len(Args) = %d, want 0", len(stmt.Args))
	}
	assertValidLoc(t, "CallStmt", stmt.Loc)
}

func TestCall_QualifiedName(t *testing.T) {
	stmt := mustCall(t, "CALL db.sch.my_sp(1, 2)")
	if stmt.Name.Normalize() != "DB.SCH.MY_SP" {
		t.Errorf("Name = %q, want DB.SCH.MY_SP", stmt.Name.Normalize())
	}
}

func TestCall_LocCoversWholeStatement(t *testing.T) {
	input := "CALL my_sp(1)"
	stmt := mustCall(t, input)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

func TestCall_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"missing parens", "CALL my_sp"},
		{"missing close paren", "CALL my_sp(1, 2"},
		{"missing name", "CALL ()"},
		{"trailing comma then close", "CALL my_sp(1,)"},
		{"named arg missing value", "CALL my_sp(x =>)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mustNotParse(t, tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// EXECUTE IMMEDIATE
//   Syntax (docs): EXECUTE IMMEDIATE { '<string>' | $$<body>$$ | <variable>
//     | $<session_variable> } [ USING ( <bind_var> [ , ... ] ) ]
// ---------------------------------------------------------------------------

func TestExecuteImmediate_StringBody(t *testing.T) {
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE 'SELECT 1'")
	if stmt.Source != ast.ExecImmString {
		t.Errorf("Source = %v, want ExecImmString", stmt.Source)
	}
	if stmt.Body != "'SELECT 1'" {
		t.Errorf("Body = %q, want \"'SELECT 1'\" (verbatim, quotes included)", stmt.Body)
	}
	assertValidLoc(t, "ExecuteImmediateStmt", stmt.Loc)
}

func TestExecuteImmediate_DollarBody(t *testing.T) {
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE $$ SELECT PI(); $$")
	if stmt.Source != ast.ExecImmDollar {
		t.Errorf("Source = %v, want ExecImmDollar", stmt.Source)
	}
	if !strings.HasPrefix(stmt.Body, "$$") || !strings.HasSuffix(stmt.Body, "$$") {
		t.Errorf("Body = %q, want it to retain the $$ delimiters", stmt.Body)
	}
	if !strings.Contains(stmt.Body, "SELECT PI()") {
		t.Errorf("Body = %q, want it to contain SELECT PI()", stmt.Body)
	}
}

func TestExecuteImmediate_MultilineDollarBody(t *testing.T) {
	// Corpus execute-immediate/example_08: a multi-line $$...$$ scripting block.
	input := "EXECUTE IMMEDIATE $$\nDECLARE\n  x FLOAT;\nBEGIN\n  RETURN x;\nEND;\n$$"
	stmt := mustExecImm(t, input)
	if stmt.Source != ast.ExecImmDollar {
		t.Errorf("Source = %v, want ExecImmDollar", stmt.Source)
	}
	if !strings.Contains(stmt.Body, "DECLARE") || !strings.Contains(stmt.Body, "END;") {
		t.Errorf("Body = %q, want the full scripting block verbatim", stmt.Body)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d (end of $$ body)", stmt.Loc.End, len(input))
	}
}

func TestExecuteImmediate_MultilineStringBody(t *testing.T) {
	// A single-quoted body that spans newlines: the lexer alone rejects it
	// (unterminated string on the newline); the raw-scan captures it verbatim
	// with no spurious lex error leaking to the result.
	input := "EXECUTE IMMEDIATE 'CREATE TABLE t (\n  i INTEGER\n)'"
	stmt := mustExecImm(t, input)
	if stmt.Source != ast.ExecImmString {
		t.Errorf("Source = %v, want ExecImmString", stmt.Source)
	}
	if !strings.Contains(stmt.Body, "CREATE TABLE t") {
		t.Errorf("Body = %q, want the verbatim multi-line string", stmt.Body)
	}
	if !strings.HasPrefix(stmt.Body, "'") || !strings.HasSuffix(stmt.Body, "'") {
		t.Errorf("Body = %q, want surrounding quotes retained", stmt.Body)
	}
}

func TestExecuteImmediate_Variable(t *testing.T) {
	// Corpus execute-immediate/example_01: EXECUTE IMMEDIATE v1;
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE v1")
	if stmt.Source != ast.ExecImmVariable {
		t.Errorf("Source = %v, want ExecImmVariable", stmt.Source)
	}
	if stmt.Var.Normalize() != "V1" {
		t.Errorf("Var = %q, want V1", stmt.Var.Normalize())
	}
}

func TestExecuteImmediate_SessionVariable(t *testing.T) {
	// Corpus execute-immediate/example_07: EXECUTE IMMEDIATE $stmt;
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE $stmt")
	if stmt.Source != ast.ExecImmSessionVar {
		t.Errorf("Source = %v, want ExecImmSessionVar", stmt.Source)
	}
	if stmt.Var.Name != "stmt" {
		t.Errorf("Var.Name = %q, want stmt (no leading $)", stmt.Var.Name)
	}
	assertValidLoc(t, "ExecuteImmediateStmt", stmt.Loc)
}

func TestExecuteImmediate_UsingClause(t *testing.T) {
	// Corpus execute-immediate/example_04 (inner): EXECUTE IMMEDIATE :query
	// USING (minimum_price, maximum_price). The :query bind is a scripting-var
	// reference; here we exercise the USING list against a variable form.
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE query USING (minimum_price, maximum_price)")
	if stmt.Source != ast.ExecImmVariable {
		t.Errorf("Source = %v, want ExecImmVariable", stmt.Source)
	}
	if len(stmt.Using) != 2 {
		t.Fatalf("len(Using) = %d, want 2", len(stmt.Using))
	}
	if stmt.Using[0].Normalize() != "MINIMUM_PRICE" || stmt.Using[1].Normalize() != "MAXIMUM_PRICE" {
		t.Errorf("Using = %v, want [MINIMUM_PRICE MAXIMUM_PRICE]", stmt.Using)
	}
	if stmt.Loc.End <= stmt.Loc.Start {
		t.Errorf("Loc = %+v, want End past the USING clause", stmt.Loc)
	}
}

func TestExecuteImmediate_StringWithUsing(t *testing.T) {
	stmt := mustExecImm(t, "EXECUTE IMMEDIATE 'SELECT ?' USING (a)")
	if stmt.Source != ast.ExecImmString {
		t.Errorf("Source = %v, want ExecImmString", stmt.Source)
	}
	if len(stmt.Using) != 1 || stmt.Using[0].Normalize() != "A" {
		t.Errorf("Using = %v, want [A]", stmt.Using)
	}
}

func TestExecuteImmediate_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"no body", "EXECUTE IMMEDIATE"},
		// Unterminated bodies MUST return fast with an error, never loop
		// (LOOP-GUARD coverage).
		{"unterminated dollar body", "EXECUTE IMMEDIATE $$ SELECT 1"},
		{"unterminated string body", "EXECUTE IMMEDIATE 'SELECT 1"},
		{"empty using list", "EXECUTE IMMEDIATE 'x' USING ()"},
		{"using missing close", "EXECUTE IMMEDIATE 'x' USING (a"},
		{"using missing open", "EXECUTE IMMEDIATE 'x' USING a)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mustNotParse(t, tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// EXECUTE TASK
//   Syntax (docs + legacy .g4):
//     EXECUTE TASK <name> [ USING CONFIG = <config_string> ]
//     EXECUTE TASK <name> RETRY LAST
// ---------------------------------------------------------------------------

func TestExecuteTask_Bare(t *testing.T) {
	stmt := mustExecTask(t, "EXECUTE TASK mytask")
	if stmt.Name.Normalize() != "MYTASK" {
		t.Errorf("Name = %q, want MYTASK", stmt.Name.Normalize())
	}
	if stmt.RetryLast {
		t.Error("RetryLast = true, want false")
	}
	if stmt.UsingConfig != "" {
		t.Errorf("UsingConfig = %q, want empty", stmt.UsingConfig)
	}
	assertValidLoc(t, "ExecuteTaskStmt", stmt.Loc)
}

func TestExecuteTask_QualifiedName(t *testing.T) {
	stmt := mustExecTask(t, "EXECUTE TASK db.sch.mytask")
	if stmt.Name.Normalize() != "DB.SCH.MYTASK" {
		t.Errorf("Name = %q, want DB.SCH.MYTASK", stmt.Name.Normalize())
	}
}

func TestExecuteTask_RetryLast(t *testing.T) {
	stmt := mustExecTask(t, "EXECUTE TASK mytask RETRY LAST")
	if !stmt.RetryLast {
		t.Error("RetryLast = false, want true")
	}
	assertValidLoc(t, "ExecuteTaskStmt", stmt.Loc)
}

func TestExecuteTask_UsingConfig(t *testing.T) {
	stmt := mustExecTask(t, "EXECUTE TASK mytask USING CONFIG = '{\"k\": 1}'")
	if stmt.UsingConfig == "" || !strings.Contains(stmt.UsingConfig, "\"k\"") {
		t.Errorf("UsingConfig = %q, want the config string literal", stmt.UsingConfig)
	}
}

func TestExecuteTask_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"missing name", "EXECUTE TASK"},
		{"retry without last", "EXECUTE TASK t RETRY"},
		{"config missing equals", "EXECUTE TASK t USING CONFIG '{}'"},
		{"config missing value", "EXECUTE TASK t USING CONFIG ="},
		{"using missing config keyword", "EXECUTE TASK t USING '{}'"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mustNotParse(t, tc.input)
		})
	}
}

func TestExecute_Reject(t *testing.T) {
	// EXECUTE followed by neither IMMEDIATE nor TASK is a syntax error.
	mustNotParse(t, "EXECUTE FOO")
	mustNotParse(t, "EXECUTE")
}

// ---------------------------------------------------------------------------
// EXPLAIN
//   Syntax (docs + legacy .g4): EXPLAIN [ USING { TABULAR | JSON | TEXT } ] <stmt>
// ---------------------------------------------------------------------------

func TestExplain_Default(t *testing.T) {
	stmt := mustExplain(t, "EXPLAIN SELECT 1")
	if stmt.Format != ast.ExplainDefault {
		t.Errorf("Format = %v, want ExplainDefault", stmt.Format)
	}
	if _, ok := stmt.Stmt.(*ast.SelectStmt); !ok {
		t.Errorf("Stmt = %T, want *ast.SelectStmt", stmt.Stmt)
	}
	assertValidLoc(t, "ExplainStmt", stmt.Loc)
}

func TestExplain_UsingFormats(t *testing.T) {
	cases := []struct {
		input string
		want  ast.ExplainFormat
	}{
		// Corpus explain/example_02..04.
		{"EXPLAIN USING TABULAR SELECT Z1.ID FROM Z1", ast.ExplainTabular},
		{"EXPLAIN USING JSON SELECT Z1.ID FROM Z1", ast.ExplainJSON},
		{"EXPLAIN USING TEXT SELECT Z1.ID FROM Z1", ast.ExplainText},
	}
	for _, tc := range cases {
		t.Run(tc.want.String(), func(t *testing.T) {
			stmt := mustExplain(t, tc.input)
			if stmt.Format != tc.want {
				t.Errorf("Format = %v, want %v", stmt.Format, tc.want)
			}
			if _, ok := stmt.Stmt.(*ast.SelectStmt); !ok {
				t.Errorf("Stmt = %T, want *ast.SelectStmt", stmt.Stmt)
			}
			assertValidLoc(t, "ExplainStmt", stmt.Loc)
		})
	}
}

func TestExplain_InnerInsert(t *testing.T) {
	// EXPLAIN over a non-SELECT statement (the inner uses parseStmt dispatch).
	stmt := mustExplain(t, "EXPLAIN INSERT INTO t VALUES (1)")
	if _, ok := stmt.Stmt.(*ast.InsertStmt); !ok {
		t.Errorf("Stmt = %T, want *ast.InsertStmt", stmt.Stmt)
	}
	assertValidLoc(t, "ExplainStmt", stmt.Loc)
}

func TestExplain_LocCoversInner(t *testing.T) {
	input := "EXPLAIN SELECT 1"
	stmt := mustExplain(t, input)
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len(input) {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len(input))
	}
}

func TestExplain_Reject(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"using without format", "EXPLAIN USING SELECT 1"},
		{"bad format", "EXPLAIN USING XML SELECT 1"},
		{"no inner statement", "EXPLAIN"},
		{"using format no inner", "EXPLAIN USING JSON"},
		{"inner garbage", "EXPLAIN FROBNICATE"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			mustNotParse(t, tc.input)
		})
	}
}

// ---------------------------------------------------------------------------
// Walker — confirm CALL argument expressions (and EXPLAIN's inner statement)
// are reachable via ast.Inspect so analysis / query-span can see them.
// ---------------------------------------------------------------------------

func TestCall_ArgsAreWalkable(t *testing.T) {
	stmt := mustCall(t, "CALL my_sp(SELECT id FROM users)")
	var sawSelect bool
	ast.Inspect(stmt, func(n ast.Node) bool {
		if _, ok := n.(*ast.SelectStmt); ok {
			sawSelect = true
		}
		return true
	})
	if !sawSelect {
		t.Error("walker did not reach the subquery argument's SelectStmt")
	}
}

func TestExplain_InnerIsWalkable(t *testing.T) {
	stmt := mustExplain(t, "EXPLAIN SELECT id FROM users")
	var sawSelect bool
	ast.Inspect(stmt, func(n ast.Node) bool {
		if _, ok := n.(*ast.SelectStmt); ok {
			sawSelect = true
		}
		return true
	})
	if !sawSelect {
		t.Error("walker did not reach EXPLAIN's inner SelectStmt")
	}
}

// ---------------------------------------------------------------------------
// Official docs corpus — every CALL / EXECUTE IMMEDIATE / EXPLAIN statement in
// the call, execute-immediate, and explain corpora must parse with zero errors
// and to its expected AST type. The official docs are the authoritative oracle
// (truth1). Context / setup statements owned by other DAG nodes (CREATE / SET /
// INSERT / the DECLARE…BEGIN scripting block) are skipped, and the $-variable
// argument forms blocked by the shared expression parser's missing tokVariable
// support (CALL p($v)) are filtered as a flagged divergence (see
// execCallDollarLimited).
// ---------------------------------------------------------------------------

var execCallCorpusDirs = []string{
	"testdata/official/call",
	"testdata/official/execute-immediate",
	"testdata/official/explain",
}

// execCallDollarLimited reports whether a CALL statement uses a $-prefixed
// session-variable ARGUMENT (CALL p($v)), which depends on the shared
// expression parser's tokVariable support owned by a separate node
// (expr-dollar-refs). EXECUTE IMMEDIATE's own $<session_var> form is handled by
// this node directly (it reads the $-token, not an expression) and is therefore
// NOT filtered here. Flagged in the divergence ledger.
func execCallDollarLimited(upper string) bool {
	return strings.HasPrefix(upper, "CALL") && strings.Contains(upper, "($")
}

func TestExecCall_OfficialCorpus(t *testing.T) {
	for _, dir := range execCallCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertExecCallStatementsParse(t, string(data))
			})
		}
	}
}

// assertExecCallStatementsParse parses sql and asserts every CALL / EXECUTE
// IMMEDIATE / EXECUTE TASK / EXPLAIN statement parses with no errors and to the
// expected AST type. Statements owned by other DAG nodes are skipped; the
// $-argument-limited CALL statements are filtered (flagged divergence).
func assertExecCallStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)
		kind, want := execCallStmtKind(upper)
		if kind == "" {
			continue // context statement owned by another DAG node
		}
		if execCallDollarLimited(upper) {
			// Known dependency limitation: must currently fail to parse. If it
			// starts parsing, the dependency lifted the limitation — surface it.
			if _, errs := parseSingle(seg.Text, seg.ByteStart); len(errs) == 0 {
				t.Logf("note: $-arg CALL now parses, drop it from execCallDollarLimited: %q", text)
			}
			continue
		}
		node, errs := parseSingle(seg.Text, seg.ByteStart)
		if len(errs) > 0 {
			t.Errorf("statement %q produced %d error(s): %v", text, len(errs), errs)
			continue
		}
		if !want(node) {
			t.Errorf("statement %q (%s) parsed to unexpected type %T", text, kind, node)
		}
	}
}

// execCallStmtKind classifies an uppercased statement by its leading keyword(s),
// returning the kind label and a predicate checking the parsed node type.
// Returns ("", nil) for statements this node does not own (CREATE / SET /
// INSERT / CALL setup lines and DECLARE…BEGIN scripting blocks).
func execCallStmtKind(upper string) (string, func(ast.Node) bool) {
	switch {
	case hasWordPrefix(upper, "EXECUTE") && strings.HasPrefix(strings.TrimSpace(strings.TrimPrefix(upper, "EXECUTE")), "IMMEDIATE"):
		return "EXECUTE-IMMEDIATE", func(n ast.Node) bool { _, ok := n.(*ast.ExecuteImmediateStmt); return ok }
	case hasWordPrefix(upper, "EXECUTE") && strings.HasPrefix(strings.TrimSpace(strings.TrimPrefix(upper, "EXECUTE")), "TASK"):
		return "EXECUTE-TASK", func(n ast.Node) bool { _, ok := n.(*ast.ExecuteTaskStmt); return ok }
	case hasWordPrefix(upper, "EXPLAIN"):
		return "EXPLAIN", func(n ast.Node) bool { _, ok := n.(*ast.ExplainStmt); return ok }
	case hasWordPrefix(upper, "CALL"):
		return "CALL", func(n ast.Node) bool { _, ok := n.(*ast.CallStmt); return ok }
	}
	return "", nil
}
