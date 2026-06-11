package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Strict trailing-token rejection (parseSegment strictTrailing).
//
// Historically the parser consumed a valid statement PREFIX and silently
// ignored every token after it: `SELECT * FFROM users` parsed clean as
// `SELECT *`. The strict path (Parse / ParseStrict / parseSingle) now
// requires the token after a successfully-parsed statement to be `;` or
// end-of-segment and reports "syntax error at or near <token>" at the first
// stray token. ParseBestEffort keeps the historical tolerance for
// completion / partial-input callers.
// ---------------------------------------------------------------------------

func TestTrailing_GarbageShapesRejected(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantLoc int    // Loc.Start of the stray token
		wantTok string // token text in the error message
	}{
		// `SELECT *` parses, FFROM is not a star alias (stars take none).
		{"misspelled FROM", "SELECT * FFROM users", 9, "FFROM"},
		// `SELECT 1` parses; a number cannot be an implicit alias.
		{"numbers run", "SELECT 1 2 3", 9, "2"},
		// `garbage` IS a legal implicit table alias, `extra` is not.
		{"garbage after alias", "SELECT a FROM b garbage extra", 24, "extra"},
		// Stray token after a complete WHERE clause.
		{"garbage after where", "SELECT a FROM b WHERE c garbage", 24, "garbage"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			// Strict Parse: first error reported, prefix statement still in File.
			file, err := Parse(c.input)
			if err == nil {
				t.Fatalf("Parse(%q): expected error, got nil", c.input)
			}
			pe, ok := err.(*ParseError)
			if !ok {
				t.Fatalf("Parse(%q): error type = %T, want *ParseError", c.input, err)
			}
			if want := "syntax error at or near " + c.wantTok; pe.Msg != want {
				t.Errorf("Parse(%q): Msg = %q, want %q", c.input, pe.Msg, want)
			}
			if pe.Loc.Start != c.wantLoc {
				t.Errorf("Parse(%q): Loc.Start = %d, want %d", c.input, pe.Loc.Start, c.wantLoc)
			}
			if len(file.Stmts) != 1 {
				t.Errorf("Parse(%q): File.Stmts = %d, want 1 (prefix statement retained)", c.input, len(file.Stmts))
			}

			// ParseStrict: same error, all-errors shape.
			result := ParseStrict(c.input)
			if len(result.Errors) != 1 {
				t.Errorf("ParseStrict(%q): %d errors, want 1: %+v", c.input, len(result.Errors), result.Errors)
			}

			// ParseBestEffort keeps the historical tolerance.
			be := ParseBestEffort(c.input)
			if len(be.Errors) != 0 {
				t.Errorf("ParseBestEffort(%q): %d errors, want 0 (prefix tolerance): %+v", c.input, len(be.Errors), be.Errors)
			}
			if len(be.File.Stmts) != 1 {
				t.Errorf("ParseBestEffort(%q): File.Stmts = %d, want 1", c.input, len(be.File.Stmts))
			}
		})
	}
}

func TestTrailing_MultiStatementPerStatementReporting(t *testing.T) {
	// The first statement has trailing garbage; the second is valid. Strict
	// parsing reports the first statement's error AND still parses the second.
	input := "SELECT 1 2; SELECT 3"
	result := ParseStrict(input)
	if len(result.Errors) != 1 {
		t.Fatalf("errors = %d, want 1: %+v", len(result.Errors), result.Errors)
	}
	if result.Errors[0].Loc.Start != 9 {
		t.Errorf("error Loc.Start = %d, want 9 (the stray `2`)", result.Errors[0].Loc.Start)
	}
	if !strings.Contains(result.Errors[0].Msg, "at or near 2") {
		t.Errorf("error Msg = %q, want to mention `2`", result.Errors[0].Msg)
	}
	if len(result.File.Stmts) != 2 {
		t.Fatalf("stmts = %d, want 2 (both statements parsed)", len(result.File.Stmts))
	}
	// Both errors when both statements have garbage.
	result = ParseStrict("SELECT 1 2; SELECT 3 4")
	if len(result.Errors) != 2 {
		t.Errorf("two-garbage errors = %d, want 2: %+v", len(result.Errors), result.Errors)
	}
}

func TestTrailing_ValidStatementsUnchanged(t *testing.T) {
	valid := []string{
		"SELECT 1",
		"SELECT 1;",
		"SELECT 1 ;",
		"SELECT 1; ",
		"SELECT 1 -- trailing comment",
		"SELECT 1; -- trailing comment\n",
		"SELECT 1 /* block */;",
		"SELECT a FROM b garbage", // implicit table alias — legal SQL
		"SELECT * FROM t",
		"BEGIN SELECT 1; END;",                 // scripting block: inner ';' belongs to the block
		"SHOW TABLES ->> SELECT * FROM $1",     // result pipe consumed by parseTopStmt
		"SELECT 1;;",                           // empty statement segments are filtered
		"SELECT 1\n;\nSELECT 2;",               // newline-separated
		"EXECUTE IMMEDIATE 'SELECT 1'",         // EXECUTE dispatched by identifier text
		"SELECT column1 FROM VALUES (1), (2);", // bare VALUES row source
	}
	for _, sql := range valid {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) error = %v, want nil", sql, err)
		}
	}
}

func TestTrailing_SemicolonInsideSegmentAccepted(t *testing.T) {
	// parseSingle on un-split text: a trailing `;` is a legal terminator.
	node, errs := parseSingle("SELECT 1;", 0)
	if len(errs) != 0 {
		t.Errorf("parseSingle(\"SELECT 1;\") errors = %+v, want none", errs)
	}
	if node == nil {
		t.Errorf("parseSingle(\"SELECT 1;\") node = nil, want statement")
	}
	// ... but tokens AFTER the terminator within one segment are rejected
	// (one segment == one statement).
	_, errs = parseSingle("SELECT 1; SELECT 2", 0)
	if len(errs) != 1 {
		t.Errorf("parseSingle(\"SELECT 1; SELECT 2\") errors = %d, want 1: %+v", len(errs), errs)
	}
}

func TestTrailing_LexErrorInDroppedTailNowSurfaces(t *testing.T) {
	// The recovery scan consumes the garbage tail, so lex errors hiding in it
	// are promoted instead of vanishing with the silent drop.
	result := ParseStrict("SELECT 1 ⊕")
	if len(result.Errors) == 0 {
		t.Fatalf("expected errors for garbage tail with invalid byte, got none")
	}
}

// ---------------------------------------------------------------------------
// Construct regressions — every shape below previously "parsed clean" only
// because its tail was silently dropped (found empirically by running the
// 657-file corpus closure in strict mode). Each now parses COMPLETELY.
// ---------------------------------------------------------------------------

func TestTrailing_FixedConstructsParseCompletely(t *testing.T) {
	cases := []string{
		// Bare VALUES row source in FROM (order-by/pivot/values docs files).
		"SELECT column1 FROM VALUES (1), (null), (2) ORDER BY column1",
		"SELECT t.column1 FROM VALUES (1), (2) AS t",
		// Join variants with DIRECTED in documented position / NATURAL INNER.
		"SELECT t1.col1 FROM t1 INNER DIRECTED JOIN t2 ON t2.col1 = t1.col1",
		"SELECT * FROM d1 NATURAL INNER JOIN d2 ORDER BY id",
		"SELECT * FROM t1 LEFT OUTER DIRECTED JOIN t2 ON t1.a = t2.a",
		// SHOW scopes for failover / replication groups.
		"SHOW DATABASES IN FAILOVER GROUP grp",
		"SHOW SHARES IN REPLICATION GROUP grp",
		// DROP FUNCTION / PROCEDURE overload signatures.
		"DROP FUNCTION obj()",
		"DROP FUNCTION obj(int)",
		"DROP PROCEDURE obj(VARCHAR, NUMBER(10,2))",
		// Postfix IF NOT EXISTS (legacy grammar placement).
		"CREATE TABLE t1 IF NOT EXISTS (v VARCHAR(16777216))",
		// CREATE DATABASE from an inbound share.
		"CREATE DATABASE snow_sales FROM SHARE ab67890.sales_s",
		// IDENTIFIER(...) object names.
		"USE WAREHOUSE IDENTIFIER($current_wh_name)",
		"UNDROP TABLE IDENTIFIER(408578)",
		// Class instance role grantee (budget instance roles).
		"SHOW GRANTS TO SNOWFLAKE.CORE.BUDGET ROLE cost.budgets.my_budget!ADMIN",
		// ALTER VIEW comma-separated column actions + drop-and-add policies.
		"ALTER VIEW v MODIFY COLUMN a SET MASKING POLICY p1, COLUMN b SET MASKING POLICY p2",
		"ALTER VIEW v MODIFY COLUMN a UNSET MASKING POLICY, COLUMN b UNSET MASKING POLICY",
		"ALTER VIEW v1 DROP ROW ACCESS POLICY rap_1, ADD ROW ACCESS POLICY rap_2 ON (empl_id)",
		// Multi-line single-quoted string (Snowflake strings may span lines).
		"SELECT PARSE_JSON(column1) FROM VALUES ('{\n  \"a\": 1\n}')",
	}
	for _, sql := range cases {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) error = %v, want nil", sql, err)
		}
	}
}

func TestTrailing_FromValuesAST(t *testing.T) {
	file, err := Parse("SELECT column1 FROM VALUES (1), (2), (3)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	sel, ok := file.Stmts[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("stmt = %T, want *ast.SelectStmt", file.Stmts[0])
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("from = %T, want *ast.TableRef", sel.From[0])
	}
	values, ok := ref.Subquery.(*ast.ValuesClause)
	if !ok {
		t.Fatalf("subquery = %T, want *ast.ValuesClause", ref.Subquery)
	}
	if len(values.Rows) != 3 {
		t.Errorf("rows = %d, want 3", len(values.Rows))
	}
}

func TestTrailing_DirectedJoinAST(t *testing.T) {
	file, err := Parse("SELECT * FROM t1 INNER DIRECTED JOIN t2 ON t1.a = t2.a")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	sel := file.Stmts[0].(*ast.SelectStmt)
	join, ok := sel.From[0].(*ast.JoinExpr)
	if !ok {
		t.Fatalf("from = %T, want *ast.JoinExpr", sel.From[0])
	}
	if !join.Directed || join.Type != ast.JoinInner || join.Natural {
		t.Errorf("join = {Type:%v Natural:%v Directed:%v}, want inner directed non-natural",
			join.Type, join.Natural, join.Directed)
	}

	file, err = Parse("SELECT * FROM d1 NATURAL INNER JOIN d2")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	sel = file.Stmts[0].(*ast.SelectStmt)
	join = sel.From[0].(*ast.JoinExpr)
	if !join.Natural || join.Type != ast.JoinInner {
		t.Errorf("NATURAL INNER JOIN: got {Type:%v Natural:%v}", join.Type, join.Natural)
	}
}

func TestTrailing_DropFunctionSignatureAST(t *testing.T) {
	file, err := Parse("DROP FUNCTION obj(int, VARCHAR)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	drop := file.Stmts[0].(*ast.DropStmt)
	if !drop.HasArgs || len(drop.ArgTypes) != 2 {
		t.Errorf("HasArgs=%v ArgTypes=%d, want true/2", drop.HasArgs, len(drop.ArgTypes))
	}
	// Empty signature: HasArgs true, zero types.
	file, err = Parse("DROP PROCEDURE obj()")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	drop = file.Stmts[0].(*ast.DropStmt)
	if !drop.HasArgs || len(drop.ArgTypes) != 0 {
		t.Errorf("empty signature: HasArgs=%v ArgTypes=%d, want true/0", drop.HasArgs, len(drop.ArgTypes))
	}
	// No signature at all: HasArgs false.
	file, err = Parse("DROP FUNCTION obj")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	drop = file.Stmts[0].(*ast.DropStmt)
	if drop.HasArgs {
		t.Errorf("no signature: HasArgs = true, want false")
	}
}

func TestTrailing_ClassRoleGranteeAST(t *testing.T) {
	file, err := Parse("SHOW GRANTS TO SNOWFLAKE.CORE.BUDGET ROLE cost.budgets.my_budget!ADMIN")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	show := file.Stmts[0].(*ast.ShowStmt)
	g := show.GrantsTo
	if g == nil {
		t.Fatalf("GrantsTo = nil")
	}
	if g.Kind != ast.GranteeClassRole {
		t.Errorf("Kind = %v, want GranteeClassRole", g.Kind)
	}
	if g.Class == nil || g.Class.Normalize() != "SNOWFLAKE.CORE.BUDGET" {
		t.Errorf("Class = %v, want SNOWFLAKE.CORE.BUDGET", g.Class)
	}
	if g.Name == nil || g.Name.Normalize() != "COST.BUDGETS.MY_BUDGET" {
		t.Errorf("Name = %v, want COST.BUDGETS.MY_BUDGET", g.Name)
	}
	if g.InstanceRole.Normalize() != "ADMIN" {
		t.Errorf("InstanceRole = %q, want ADMIN", g.InstanceRole.Name)
	}
}

func TestTrailing_AlterViewColumnActionsAST(t *testing.T) {
	file, err := Parse("ALTER VIEW v MODIFY COLUMN a SET MASKING POLICY p1, COLUMN b UNSET MASKING POLICY")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	av := file.Stmts[0].(*ast.AlterViewStmt)
	if len(av.ColumnActions) != 2 {
		t.Fatalf("ColumnActions = %d, want 2", len(av.ColumnActions))
	}
	// First action mirrored into legacy fields.
	if av.Action != ast.AlterViewColumnSetMaskingPolicy || av.Column.Normalize() != "A" {
		t.Errorf("mirror: Action=%v Column=%q", av.Action, av.Column.Name)
	}
	if av.ColumnActions[1].Action != ast.AlterViewColumnUnsetMaskingPolicy ||
		av.ColumnActions[1].Column.Normalize() != "B" {
		t.Errorf("second action = {%v %q}", av.ColumnActions[1].Action, av.ColumnActions[1].Column.Name)
	}

	file, err = Parse("ALTER VIEW v1 DROP ROW ACCESS POLICY r1, ADD ROW ACCESS POLICY r2 ON (empl_id)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	av = file.Stmts[0].(*ast.AlterViewStmt)
	if av.Action != ast.AlterViewDropRowAccessPolicy {
		t.Errorf("Action = %v, want AlterViewDropRowAccessPolicy", av.Action)
	}
	if av.AddPolicyName == nil || av.AddPolicyName.Normalize() != "R2" || len(av.AddPolicyCols) != 1 {
		t.Errorf("AddPolicyName = %v, AddPolicyCols = %v", av.AddPolicyName, av.AddPolicyCols)
	}
}

func TestTrailing_PutFileURLDoesNotEatNextStatement(t *testing.T) {
	// `file://` used to open a `//` line comment in Split's lexer, hiding the
	// real `;` so the following statement fused into the PUT's segment and
	// was silently dropped. The `://` guard keeps the `;` visible.
	input := "put file:// @stage/;\nremove @stage/"
	segs := Split(input)
	if len(segs) != 2 {
		t.Fatalf("Split produced %d segments, want 2: %+v", len(segs), segs)
	}
	result := ParseStrict(input)
	if len(result.Errors) != 0 {
		t.Errorf("errors = %+v, want none", result.Errors)
	}
	if len(result.File.Stmts) != 2 {
		t.Fatalf("stmts = %d, want 2 (PUT and REMOVE)", len(result.File.Stmts))
	}
	if _, ok := result.File.Stmts[0].(*ast.PutStmt); !ok {
		t.Errorf("stmt[0] = %T, want *ast.PutStmt", result.File.Stmts[0])
	}
	if _, ok := result.File.Stmts[1].(*ast.RemoveStmt); !ok {
		t.Errorf("stmt[1] = %T, want *ast.RemoveStmt", result.File.Stmts[1])
	}
	// A real `//` line comment (not preceded by `:`) still comments.
	if _, err := Parse("SELECT 1 // trailing comment"); err != nil {
		t.Errorf("// comment broken: %v", err)
	}
}

func TestTrailing_IdentifierFuncName(t *testing.T) {
	file, err := Parse("USE WAREHOUSE IDENTIFIER($current_wh_name)")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	use := file.Stmts[0].(*ast.UseStmt)
	if use.Name == nil || use.Name.Name.Name != "IDENTIFIER($current_wh_name)" {
		t.Errorf("Name = %+v, want verbatim IDENTIFIER($current_wh_name)", use.Name)
	}
	// A table genuinely named identifier (no parens) is untouched.
	file, err = Parse("SELECT * FROM identifier")
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
}
