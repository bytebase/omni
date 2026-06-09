package parser

// Unit tests for the parser-scripting node: GoogleSQL's procedural language
// (DECLARE / SET / IF / CASE / WHILE / LOOP / REPEAT / FOR-IN / BEGIN…END with
// EXCEPTION / RAISE / RETURN / BREAK / CONTINUE / EXECUTE IMMEDIATE / labels).
//
// These assert the AST shape and accept/reject polarity on hand-authored
// fixtures. The differential against the live Spanner emulator (the PROVE gate
// per correctness-protocol.md) lives in scripting_oracle_test.go (build tag
// googlesql_oracle); scripting is BigQuery-only except the BEGIN…END envelope /
// SET / EXECUTE IMMEDIATE, so the control-flow forms are triangulated against the
// legacy .g4 + the BigQuery truth1 corpus, not the Spanner oracle.

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// parseOneScript parses src (one top-level statement / block) and returns the
// single resulting node, failing on any parse error or a missing/extra node.
func parseOneScript(t *testing.T, src string) ast.Node {
	t.Helper()
	file, errs := Parse(src)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q) returned errors: %v", src, errs)
	}
	if file == nil || len(file.Stmts) != 1 {
		got := 0
		if file != nil {
			got = len(file.Stmts)
		}
		t.Fatalf("Parse(%q): expected exactly 1 statement, got %d", src, got)
	}
	return file.Stmts[0]
}

// ----------------------------------------------------------------------------
// DECLARE
// ----------------------------------------------------------------------------

func TestParseDeclareTyped(t *testing.T) {
	n := parseOneScript(t, "DECLARE x INT64")
	d, ok := n.(*ast.DeclareStmt)
	if !ok {
		t.Fatalf("expected *ast.DeclareStmt, got %T", n)
	}
	if len(d.Names) != 1 || d.Names[0].Name != "x" {
		t.Fatalf("Names = %+v, want [x]", d.Names)
	}
	if d.Type == nil || d.Type.Text != "INT64" {
		t.Fatalf("Type = %+v, want INT64", d.Type)
	}
	if d.Default != nil {
		t.Fatalf("Default = %+v, want nil", d.Default)
	}
}

func TestParseDeclareTypedWithDefault(t *testing.T) {
	n := parseOneScript(t, "DECLARE d DATE DEFAULT CURRENT_DATE()")
	d := n.(*ast.DeclareStmt)
	if d.Type == nil || d.Type.Text != "DATE" {
		t.Fatalf("Type = %+v, want DATE", d.Type)
	}
	if d.Default == nil {
		t.Fatalf("Default = nil, want CURRENT_DATE() expr")
	}
}

func TestParseDeclareMultipleNames(t *testing.T) {
	n := parseOneScript(t, "DECLARE x, y, z INT64 DEFAULT 0")
	d := n.(*ast.DeclareStmt)
	if len(d.Names) != 3 {
		t.Fatalf("Names = %+v, want 3 names", d.Names)
	}
	if d.Names[0].Name != "x" || d.Names[1].Name != "y" || d.Names[2].Name != "z" {
		t.Fatalf("Names spelling = %+v", d.Names)
	}
	if d.Type == nil || d.Type.Text != "INT64" {
		t.Fatalf("Type = %+v, want INT64", d.Type)
	}
}

func TestParseDeclareDefaultOnly(t *testing.T) {
	n := parseOneScript(t, "DECLARE item DEFAULT (SELECT 1)")
	d := n.(*ast.DeclareStmt)
	if d.Type != nil {
		t.Fatalf("Type = %+v, want nil (DEFAULT-only form)", d.Type)
	}
	if d.Default == nil {
		t.Fatalf("Default = nil, want subquery expr")
	}
}

func TestParseDeclareDefaultSubqueryFilled(t *testing.T) {
	// The DEFAULT initializer's subquery must be re-parsed (fillSubqueries) so the
	// query-span walker can reach it.
	n := parseOneScript(t, "DECLARE item DEFAULT (SELECT col FROM s.products LIMIT 1)")
	d := n.(*ast.DeclareStmt)
	sq, ok := d.Default.(*ast.SubqueryExpr)
	if !ok {
		t.Fatalf("Default = %T, want *ast.SubqueryExpr", d.Default)
	}
	if sq.Query == nil {
		t.Fatalf("subquery Query is nil — fillSubqueries did not reach the DECLARE DEFAULT")
	}
}

// ----------------------------------------------------------------------------
// SET
// ----------------------------------------------------------------------------

func TestParseSetSingleVariable(t *testing.T) {
	n := parseOneScript(t, "SET x = 5")
	s, ok := n.(*ast.SetStmt)
	if !ok {
		t.Fatalf("expected *ast.SetStmt, got %T", n)
	}
	if s.Kind != ast.SetVariable {
		t.Fatalf("Kind = %v, want SetVariable", s.Kind)
	}
	if len(s.Targets) != 1 {
		t.Fatalf("Targets = %+v, want 1", s.Targets)
	}
	id, ok := s.Targets[0].(*ast.Identifier)
	if !ok || id.Name != "x" {
		t.Fatalf("Targets[0] = %+v, want Identifier x", s.Targets[0])
	}
	if s.Value == nil {
		t.Fatalf("Value = nil, want 5")
	}
}

func TestParseSetTuple(t *testing.T) {
	n := parseOneScript(t, "SET (a, b, c) = (1 + 3, 'foo', false)")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTuple {
		t.Fatalf("Kind = %v, want SetTuple", s.Kind)
	}
	if len(s.Targets) != 3 {
		t.Fatalf("Targets = %+v, want 3", s.Targets)
	}
	if s.Value == nil {
		t.Fatalf("Value = nil, want RHS tuple")
	}
}

func TestParseSetParameter(t *testing.T) {
	n := parseOneScript(t, "SET @p = 5")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetVariable {
		t.Fatalf("Kind = %v, want SetVariable", s.Kind)
	}
	if _, ok := s.Targets[0].(*ast.Parameter); !ok {
		t.Fatalf("Targets[0] = %T, want *ast.Parameter", s.Targets[0])
	}
}

func TestParseSetSystemVariable(t *testing.T) {
	n := parseOneScript(t, "SET @@x.y = 5")
	s := n.(*ast.SetStmt)
	if _, ok := s.Targets[0].(*ast.SystemVariable); !ok {
		t.Fatalf("Targets[0] = %T, want *ast.SystemVariable", s.Targets[0])
	}
}

func TestParseSetTransaction(t *testing.T) {
	n := parseOneScript(t, "SET TRANSACTION READ ONLY")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTransaction {
		t.Fatalf("Kind = %v, want SetTransaction", s.Kind)
	}
	if len(s.Modes) != 1 || s.Modes[0].Kind != ast.TransactionModeReadOnly {
		t.Fatalf("Modes = %+v, want [READ ONLY]", s.Modes)
	}
}

func TestParseSetTransactionIsolation(t *testing.T) {
	n := parseOneScript(t, "SET TRANSACTION ISOLATION LEVEL SERIALIZABLE")
	s := n.(*ast.SetStmt)
	if s.Kind != ast.SetTransaction {
		t.Fatalf("Kind = %v, want SetTransaction", s.Kind)
	}
	if len(s.Modes) != 1 || s.Modes[0].Kind != ast.TransactionModeIsolationLevel {
		t.Fatalf("Modes = %+v, want [ISOLATION LEVEL]", s.Modes)
	}
}

// ----------------------------------------------------------------------------
// EXECUTE IMMEDIATE
// ----------------------------------------------------------------------------

func TestParseExecuteImmediateBare(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT 1"`)
	e, ok := n.(*ast.ExecuteImmediateStmt)
	if !ok {
		t.Fatalf("expected *ast.ExecuteImmediateStmt, got %T", n)
	}
	if e.SQL == nil {
		t.Fatalf("SQL = nil")
	}
	if e.Into != nil || e.Using != nil {
		t.Fatalf("Into/Using should be nil, got %+v / %+v", e.Into, e.Using)
	}
}

func TestParseExecuteImmediateIntoUsingPositional(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT ? * (? + 2)" INTO y USING 1, 3`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Into) != 1 || e.Into[0].Name != "y" {
		t.Fatalf("Into = %+v, want [y]", e.Into)
	}
	if len(e.Using) != 2 {
		t.Fatalf("Using = %+v, want 2 args", e.Using)
	}
	if e.Using[0].Alias != "" || e.Using[1].Alias != "" {
		t.Fatalf("positional USING args should have no alias, got %q / %q", e.Using[0].Alias, e.Using[1].Alias)
	}
}

func TestParseExecuteImmediateUsingNamed(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT @a * (@b + 2)" INTO y USING 1 AS a, 3 AS b`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Using) != 2 {
		t.Fatalf("Using = %+v, want 2 args", e.Using)
	}
	if e.Using[0].Alias != "a" || e.Using[1].Alias != "b" {
		t.Fatalf("USING aliases = %q / %q, want a / b", e.Using[0].Alias, e.Using[1].Alias)
	}
}

func TestParseExecuteImmediateIntoMultiple(t *testing.T) {
	n := parseOneScript(t, `EXECUTE IMMEDIATE "SELECT 1, 2" INTO a, b`)
	e := n.(*ast.ExecuteImmediateStmt)
	if len(e.Into) != 2 || e.Into[0].Name != "a" || e.Into[1].Name != "b" {
		t.Fatalf("Into = %+v, want [a b]", e.Into)
	}
}

// ----------------------------------------------------------------------------
// IF
// ----------------------------------------------------------------------------

func TestParseIfSimple(t *testing.T) {
	n := parseOneScript(t, "IF x > 0 THEN SELECT 1; END IF")
	f, ok := n.(*ast.IfStmt)
	if !ok {
		t.Fatalf("expected *ast.IfStmt, got %T", n)
	}
	if f.Cond == nil {
		t.Fatalf("Cond = nil")
	}
	if len(f.Then) != 1 {
		t.Fatalf("Then = %+v, want 1 statement", f.Then)
	}
	if f.HasElse {
		t.Fatalf("HasElse = true, want false")
	}
}

func TestParseIfElseifElse(t *testing.T) {
	src := "IF c1 THEN SELECT 1; ELSEIF c2 THEN SELECT 2; ELSE SELECT 3; END IF"
	n := parseOneScript(t, src)
	f := n.(*ast.IfStmt)
	if len(f.ElseIf) != 1 {
		t.Fatalf("ElseIf = %+v, want 1", f.ElseIf)
	}
	if f.ElseIf[0].Cond == nil || len(f.ElseIf[0].Then) != 1 {
		t.Fatalf("ElseIf[0] = %+v", f.ElseIf[0])
	}
	if !f.HasElse || len(f.Else) != 1 {
		t.Fatalf("Else = %+v (HasElse=%v), want 1 statement", f.Else, f.HasElse)
	}
}

func TestParseIfEmptyThen(t *testing.T) {
	// statement_list is optional after THEN.
	n := parseOneScript(t, "IF x THEN END IF")
	f := n.(*ast.IfStmt)
	if len(f.Then) != 0 {
		t.Fatalf("Then = %+v, want empty", f.Then)
	}
}

func TestParseIfMultipleStatements(t *testing.T) {
	n := parseOneScript(t, "IF x THEN SELECT 1; SELECT 2; SET y = 3; END IF")
	f := n.(*ast.IfStmt)
	if len(f.Then) != 3 {
		t.Fatalf("Then = %+v, want 3 statements", f.Then)
	}
}

// ----------------------------------------------------------------------------
// CASE (statement)
// ----------------------------------------------------------------------------

func TestParseCaseSearched(t *testing.T) {
	src := "CASE WHEN c1 THEN SELECT 1; WHEN c2 THEN SELECT 2; ELSE SELECT 3; END CASE"
	n := parseOneScript(t, src)
	c, ok := n.(*ast.CaseStmt)
	if !ok {
		t.Fatalf("expected *ast.CaseStmt, got %T", n)
	}
	if c.Operand != nil {
		t.Fatalf("Operand = %+v, want nil (searched form)", c.Operand)
	}
	if len(c.Whens) != 2 {
		t.Fatalf("Whens = %+v, want 2", c.Whens)
	}
	if !c.HasElse {
		t.Fatalf("HasElse = false, want true")
	}
}

func TestParseCaseSimple(t *testing.T) {
	src := "CASE product_id WHEN 1 THEN SELECT 'one'; WHEN 2 THEN SELECT 'two'; ELSE SELECT 'other'; END CASE"
	n := parseOneScript(t, src)
	c := n.(*ast.CaseStmt)
	if c.Operand == nil {
		t.Fatalf("Operand = nil, want product_id (simple form)")
	}
	if len(c.Whens) != 2 {
		t.Fatalf("Whens = %+v, want 2", c.Whens)
	}
}

func TestParseCaseNoElse(t *testing.T) {
	n := parseOneScript(t, "CASE WHEN c THEN SELECT 1; END CASE")
	c := n.(*ast.CaseStmt)
	if c.HasElse {
		t.Fatalf("HasElse = true, want false")
	}
}

// ----------------------------------------------------------------------------
// WHILE / LOOP / REPEAT
// ----------------------------------------------------------------------------

func TestParseWhile(t *testing.T) {
	n := parseOneScript(t, "WHILE x < 10 DO SET x = x + 1; END WHILE")
	w, ok := n.(*ast.WhileStmt)
	if !ok {
		t.Fatalf("expected *ast.WhileStmt, got %T", n)
	}
	if w.Cond == nil {
		t.Fatalf("Cond = nil")
	}
	if len(w.Body) != 1 {
		t.Fatalf("Body = %+v, want 1 statement", w.Body)
	}
}

func TestParseLoop(t *testing.T) {
	n := parseOneScript(t, "LOOP SET x = x + 1; IF x >= 10 THEN LEAVE; END IF; END LOOP")
	l, ok := n.(*ast.LoopStmt)
	if !ok {
		t.Fatalf("expected *ast.LoopStmt, got %T", n)
	}
	if len(l.Body) != 2 {
		t.Fatalf("Body = %+v, want 2 statements", l.Body)
	}
}

func TestParseRepeat(t *testing.T) {
	n := parseOneScript(t, "REPEAT SET x = x + 1; SELECT x; UNTIL x >= 3 END REPEAT")
	r, ok := n.(*ast.RepeatStmt)
	if !ok {
		t.Fatalf("expected *ast.RepeatStmt, got %T", n)
	}
	if len(r.Body) != 2 {
		t.Fatalf("Body = %+v, want 2 statements", r.Body)
	}
	if r.Until == nil {
		t.Fatalf("Until = nil, want condition")
	}
}

// ----------------------------------------------------------------------------
// FOR ... IN
// ----------------------------------------------------------------------------

func TestParseForIn(t *testing.T) {
	src := "FOR record IN (SELECT word, cnt FROM s.shakespeare LIMIT 5) DO SELECT record.word; END FOR"
	n := parseOneScript(t, src)
	f, ok := n.(*ast.ForInStmt)
	if !ok {
		t.Fatalf("expected *ast.ForInStmt, got %T", n)
	}
	if f.Var == nil || f.Var.Name != "record" {
		t.Fatalf("Var = %+v, want record", f.Var)
	}
	if f.Query == nil {
		t.Fatalf("Query = nil, want the source query")
	}
	if len(f.Body) != 1 {
		t.Fatalf("Body = %+v, want 1 statement", f.Body)
	}
}

// ----------------------------------------------------------------------------
// BEGIN ... END
// ----------------------------------------------------------------------------

func TestParseBeginEndSimple(t *testing.T) {
	n := parseOneScript(t, "BEGIN SELECT 1; END")
	b, ok := n.(*ast.BeginEndBlock)
	if !ok {
		t.Fatalf("expected *ast.BeginEndBlock, got %T", n)
	}
	if len(b.Body) != 1 {
		t.Fatalf("Body = %+v, want 1 statement", b.Body)
	}
	if b.HasException {
		t.Fatalf("HasException = true, want false")
	}
}

func TestParseBeginEndEmpty(t *testing.T) {
	n := parseOneScript(t, "BEGIN END")
	b := n.(*ast.BeginEndBlock)
	if len(b.Body) != 0 {
		t.Fatalf("Body = %+v, want empty", b.Body)
	}
}

func TestParseBeginEndMultiStatement(t *testing.T) {
	n := parseOneScript(t, "BEGIN DECLARE y INT64; SET y = 1; SELECT y; END")
	b := n.(*ast.BeginEndBlock)
	if len(b.Body) != 3 {
		t.Fatalf("Body = %+v, want 3 statements", b.Body)
	}
	if _, ok := b.Body[0].(*ast.DeclareStmt); !ok {
		t.Fatalf("Body[0] = %T, want DeclareStmt", b.Body[0])
	}
}

func TestParseBeginExceptionEnd(t *testing.T) {
	src := "BEGIN CALL s.proc2(); EXCEPTION WHEN ERROR THEN SELECT @@error.message; END"
	n := parseOneScript(t, src)
	b := n.(*ast.BeginEndBlock)
	if !b.HasException {
		t.Fatalf("HasException = false, want true")
	}
	if len(b.Exception) != 1 {
		t.Fatalf("Exception = %+v, want 1 statement", b.Exception)
	}
}

// TestParseBeginExceptionEmptyRejected pins the fix for the empty-handler
// over-accept (Codex review): opt_exception_handler's statement_list is non-empty
// in the .g4 (unterminated_non_empty_statement_list ';'), so `EXCEPTION WHEN
// ERROR THEN` with no handler statement is a syntax error. The Spanner emulator's
// BeginStmt recognizer is shallow here (divergence #161), so the grammar governs.
func TestParseBeginExceptionEmptyRejected(t *testing.T) {
	assertReject(t, "BEGIN EXCEPTION WHEN ERROR THEN END")
	// A non-empty handler is still accepted.
	parseOneScript(t, "BEGIN EXCEPTION WHEN ERROR THEN SELECT 1; END")
}

// TestParseLabeledEndLabelNotMatched defends a Codex finding: a labeled
// statement's optional END label is grammatically just `identifier?`
// (.g4: `label ':' unterminated_unlabeled_script_statement identifier?`), with NO
// syntactic requirement that it match the opening label — matching is a SEMANTIC
// check, not the parser's job. So omni accepts a mismatched end label (parity with
// the legacy ANTLR parser); this guards against a future change that wrongly
// rejects it at parse time. (parseOneScript fatals on any parse error.)
func TestParseLabeledEndLabelNotMatched(t *testing.T) {
	parseOneScript(t, "label_1: LOOP SELECT 1; END LOOP different")
}

func TestParseBeginEndNested(t *testing.T) {
	n := parseOneScript(t, "BEGIN BEGIN SELECT 1; END; SELECT 2; END")
	b := n.(*ast.BeginEndBlock)
	if len(b.Body) != 2 {
		t.Fatalf("Body = %+v, want 2 statements (nested block + select)", b.Body)
	}
	if _, ok := b.Body[0].(*ast.BeginEndBlock); !ok {
		t.Fatalf("Body[0] = %T, want nested BeginEndBlock", b.Body[0])
	}
}

func TestParseBeginEndSubqueryFilled(t *testing.T) {
	// An expression subquery inside a block statement (here a SET RHS) must be
	// re-parsed (fillSubqueries) so the query-span walker can reach its inner
	// query. (A FROM-clause derived table is parsed directly into a *QueryStmt by
	// parser-select and is not a SubqueryExpr, so it is not what this asserts.)
	n := parseOneScript(t, "BEGIN SET x = (SELECT a FROM s); END")
	b := n.(*ast.BeginEndBlock)
	found := false
	ast.Inspect(b, func(node ast.Node) bool {
		if sq, ok := node.(*ast.SubqueryExpr); ok && sq.Query != nil {
			found = true
		}
		return true
	})
	if !found {
		t.Fatalf("no filled subquery found inside the block — fillSubqueries did not descend")
	}
}

// ----------------------------------------------------------------------------
// BREAK / LEAVE / CONTINUE / ITERATE / RETURN
// ----------------------------------------------------------------------------

func TestParseBreak(t *testing.T) {
	n := parseOneScript(t, "BREAK")
	b, ok := n.(*ast.BreakStmt)
	if !ok {
		t.Fatalf("expected *ast.BreakStmt, got %T", n)
	}
	if b.IsLeave || b.Label != "" {
		t.Fatalf("got IsLeave=%v Label=%q, want false/empty", b.IsLeave, b.Label)
	}
}

func TestParseLeaveWithLabel(t *testing.T) {
	n := parseOneScript(t, "LEAVE my_label")
	b := n.(*ast.BreakStmt)
	if !b.IsLeave {
		t.Fatalf("IsLeave = false, want true")
	}
	if b.Label != "my_label" {
		t.Fatalf("Label = %q, want my_label", b.Label)
	}
}

func TestParseContinue(t *testing.T) {
	n := parseOneScript(t, "CONTINUE")
	c, ok := n.(*ast.ContinueStmt)
	if !ok {
		t.Fatalf("expected *ast.ContinueStmt, got %T", n)
	}
	if c.IsIterate || c.Label != "" {
		t.Fatalf("got IsIterate=%v Label=%q", c.IsIterate, c.Label)
	}
}

func TestParseIterateWithLabel(t *testing.T) {
	n := parseOneScript(t, "ITERATE outer_loop")
	c := n.(*ast.ContinueStmt)
	if !c.IsIterate {
		t.Fatalf("IsIterate = false, want true")
	}
	if c.Label != "outer_loop" {
		t.Fatalf("Label = %q, want outer_loop", c.Label)
	}
}

func TestParseReturn(t *testing.T) {
	n := parseOneScript(t, "RETURN")
	if _, ok := n.(*ast.ReturnStmt); !ok {
		t.Fatalf("expected *ast.ReturnStmt, got %T", n)
	}
}

// ----------------------------------------------------------------------------
// RAISE
// ----------------------------------------------------------------------------

func TestParseRaiseBare(t *testing.T) {
	n := parseOneScript(t, "RAISE")
	r, ok := n.(*ast.RaiseStmt)
	if !ok {
		t.Fatalf("expected *ast.RaiseStmt, got %T", n)
	}
	if r.Message != nil {
		t.Fatalf("Message = %+v, want nil", r.Message)
	}
}

func TestParseRaiseUsingMessage(t *testing.T) {
	n := parseOneScript(t, "RAISE USING MESSAGE = 'something went wrong'")
	r := n.(*ast.RaiseStmt)
	if r.Message == nil {
		t.Fatalf("Message = nil, want the message expr")
	}
}

// ----------------------------------------------------------------------------
// Labels
// ----------------------------------------------------------------------------

func TestParseLabeledLoop(t *testing.T) {
	n := parseOneScript(t, "label_1: LOOP SELECT 1; END LOOP label_1")
	l, ok := n.(*ast.LabeledStmt)
	if !ok {
		t.Fatalf("expected *ast.LabeledStmt, got %T", n)
	}
	if l.Label != "label_1" {
		t.Fatalf("Label = %q, want label_1", l.Label)
	}
	if l.EndLabel != "label_1" {
		t.Fatalf("EndLabel = %q, want label_1", l.EndLabel)
	}
	if _, ok := l.Stmt.(*ast.LoopStmt); !ok {
		t.Fatalf("Stmt = %T, want LoopStmt", l.Stmt)
	}
}

func TestParseLabeledBlockNoEndLabel(t *testing.T) {
	n := parseOneScript(t, "label_1: BEGIN SELECT 1; BREAK label_1; END")
	l := n.(*ast.LabeledStmt)
	if l.Label != "label_1" {
		t.Fatalf("Label = %q, want label_1", l.Label)
	}
	if l.EndLabel != "" {
		t.Fatalf("EndLabel = %q, want empty", l.EndLabel)
	}
	if _, ok := l.Stmt.(*ast.BeginEndBlock); !ok {
		t.Fatalf("Stmt = %T, want BeginEndBlock", l.Stmt)
	}
}

func TestParseLabeledWhile(t *testing.T) {
	n := parseOneScript(t, "lbl: WHILE x DO SELECT 1; END WHILE")
	l := n.(*ast.LabeledStmt)
	if _, ok := l.Stmt.(*ast.WhileStmt); !ok {
		t.Fatalf("Stmt = %T, want WhileStmt", l.Stmt)
	}
}

// ----------------------------------------------------------------------------
// Deep nesting (truth1 doc examples)
// ----------------------------------------------------------------------------

// TestParseNestedLabeledLoop is SCRIPT-009's nested labeled LOOP / WHILE / IF
// with CONTINUE / BREAK to the outer label — the deepest control-flow nesting in
// the BigQuery procedural docs.
func TestParseNestedLabeledLoop(t *testing.T) {
	src := "label_1: LOOP WHILE x < 1 DO IF y < 1 THEN CONTINUE label_1; ELSE BREAK label_1; END IF; END WHILE; END LOOP label_1"
	n := parseOneScript(t, src)
	lbl, ok := n.(*ast.LabeledStmt)
	if !ok {
		t.Fatalf("expected *ast.LabeledStmt, got %T", n)
	}
	loop, ok := lbl.Stmt.(*ast.LoopStmt)
	if !ok || len(loop.Body) != 1 {
		t.Fatalf("loop body = %+v", lbl.Stmt)
	}
	while, ok := loop.Body[0].(*ast.WhileStmt)
	if !ok || len(while.Body) != 1 {
		t.Fatalf("while body = %+v", loop.Body[0])
	}
	ifStmt, ok := while.Body[0].(*ast.IfStmt)
	if !ok {
		t.Fatalf("expected nested IfStmt, got %T", while.Body[0])
	}
	cont, ok := ifStmt.Then[0].(*ast.ContinueStmt)
	if !ok || cont.Label != "label_1" {
		t.Fatalf("THEN[0] = %+v, want CONTINUE label_1", ifStmt.Then[0])
	}
	brk, ok := ifStmt.Else[0].(*ast.BreakStmt)
	if !ok || brk.Label != "label_1" {
		t.Fatalf("ELSE[0] = %+v, want BREAK label_1", ifStmt.Else[0])
	}
}

// TestParseBlockTransactionWithException is SCRIPT-016's nested BEGIN
// TRANSACTION / COMMIT / ROLLBACK inside a block with an EXCEPTION handler.
func TestParseBlockTransactionWithException(t *testing.T) {
	src := "BEGIN BEGIN TRANSACTION; SELECT 1; COMMIT TRANSACTION; EXCEPTION WHEN ERROR THEN SELECT @@error.message; ROLLBACK TRANSACTION; END"
	n := parseOneScript(t, src)
	b := n.(*ast.BeginEndBlock)
	if len(b.Body) != 3 {
		t.Fatalf("Body = %+v, want 3 (BEGIN TRANSACTION / SELECT / COMMIT)", b.Body)
	}
	if !b.HasException || len(b.Exception) != 2 {
		t.Fatalf("Exception = %+v (HasException=%v), want 2 statements", b.Exception, b.HasException)
	}
	if _, ok := b.Body[0].(*ast.TransactionStmt); !ok {
		t.Fatalf("Body[0] = %T, want TransactionStmt (BEGIN TRANSACTION)", b.Body[0])
	}
}

// ----------------------------------------------------------------------------
// Reject cases (negative tests)
// ----------------------------------------------------------------------------

func TestScriptingRejects(t *testing.T) {
	rejects := []string{
		// SET shapes the .g4 rejects.
		"SET",           // missing target
		"SET x",         // missing = expr
		"SET x =",       // missing RHS
		"SET a, b = 1",  // multiple vars without parens (grammar error alt)
		"SET (a, b) = ", // missing RHS
		"SET () = (1)",  // empty identifier list
		// DECLARE shapes.
		"DECLARE",         // missing name
		"DECLARE x",       // missing type and DEFAULT
		"DECLARE 1 INT64", // name is not an identifier
		// Control flow missing required keywords.
		"IF x SELECT 1; END IF",                  // missing THEN
		"IF x THEN SELECT 1;",                    // missing END IF
		"WHILE x SELECT 1; END WHILE",            // missing DO
		"WHILE x DO SELECT 1; END",               // END not followed by WHILE
		"REPEAT SELECT 1; END REPEAT",            // missing UNTIL clause
		"LOOP SELECT 1; END WHILE",               // mismatched closer
		"FOR x (SELECT 1) DO SELECT 1; END FOR",  // missing IN
		"FOR x IN SELECT 1 DO SELECT 1; END FOR", // IN source not parenthesized
		"CASE WHEN c THEN SELECT 1; END",         // END not followed by CASE
		"CASE END CASE",                          // no WHEN clause
		// EXECUTE IMMEDIATE shapes.
		"EXECUTE",                      // missing IMMEDIATE
		"EXECUTE IMMEDIATE",            // missing sql expr
		`EXECUTE IMMEDIATE "x" INTO`,   // INTO with no variable
		`EXECUTE IMMEDIATE "x" USING`,  // USING with no arg
		`EXECUTE IMMEDIATE "x" INTO 1`, // INTO target not an identifier
		// BEGIN ... END.
		"BEGIN SELECT 1;", // unterminated block (no END)
		"BEGIN EXCEPTION WHEN ERROR THEN SELECT 1; END", // EXCEPTION needs a body before? actually valid; remove
	}
	for _, src := range rejects {
		if src == "BEGIN EXCEPTION WHEN ERROR THEN SELECT 1; END" {
			continue // handled below as accept
		}
		t.Run(src, func(t *testing.T) {
			_, errs := Parse(src)
			if len(errs) == 0 {
				t.Errorf("Parse(%q) accepted, want a syntax error", src)
			}
		})
	}
}
