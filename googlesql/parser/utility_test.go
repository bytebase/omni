package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Tests for the `parser-utility` node's metadata statements (§2.10 + §2.3):
// ASSERT / ANALYZE / DESCRIBE / RENAME / CALL.
//
// CORRECTNESS BASIS (correctness-protocol.md):
//   - the live Cloud Spanner emulator (the suite:utility rows of testdata/compat.tsv differential):
//     CALL is parsed in FULL (authoritative for the arg list); RENAME and bare
//     ANALYZE go through the real DDL parser (authoritative); ASSERT / DESCRIBE
//     are accepted by the shallow recognizer (leading-form authoritative,
//     trailing tokens swallowed → non-authoritative, follow the .g4).
//   - the canonical ZetaSQL corpus (assert.sql / call.sql / describe.sql /
//     analyze.sql) — the breadth oracle for the precise grammar.
//
// Spanner narrowings vs the union grammar are flagged divergences and covered
// by the hand-written cases, NOT diffed against the emulator:
//   - ANALYZE with targets (`ANALYZE t`) — Spanner rejects (its ANALYZE is the
//     bare whole-database form); the .g4 + BigQuery union accept targets.
//   - RENAME with a non-TABLE object kind — Spanner only allows `RENAME TABLE`.

func parseAssert(t *testing.T, sql string) *ast.AssertStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	a, ok := n.(*ast.AssertStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): got %T, want *ast.AssertStmt", sql, n)
	}
	return a
}

func parseAnalyze(t *testing.T, sql string) *ast.AnalyzeStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	a, ok := n.(*ast.AnalyzeStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): got %T, want *ast.AnalyzeStmt", sql, n)
	}
	return a
}

func parseDescribe(t *testing.T, sql string) *ast.DescribeStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	d, ok := n.(*ast.DescribeStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): got %T, want *ast.DescribeStmt", sql, n)
	}
	return d
}

func parseRename(t *testing.T, sql string) *ast.RenameStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	r, ok := n.(*ast.RenameStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): got %T, want *ast.RenameStmt", sql, n)
	}
	return r
}

func parseCall(t *testing.T, sql string) *ast.CallStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	c, ok := n.(*ast.CallStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): got %T, want *ast.CallStmt", sql, n)
	}
	return c
}

// ---------------------------------------------------------------------------
// ASSERT
// ---------------------------------------------------------------------------

func TestAssert(t *testing.T) {
	t.Run("bare boolean", func(t *testing.T) {
		a := parseAssert(t, "ASSERT TRUE")
		if a.Expr == nil {
			t.Fatal("Expr = nil")
		}
		if a.Description != nil {
			t.Errorf("Description = %v, want nil", a.Description)
		}
	})

	t.Run("with description", func(t *testing.T) {
		a := parseAssert(t, `ASSERT 1 = 1 AS 'must hold'`)
		if a.Expr == nil {
			t.Fatal("Expr = nil")
		}
		lit, ok := a.Description.(*ast.Literal)
		if !ok {
			t.Fatalf("Description = %T, want *ast.Literal", a.Description)
		}
		if lit.Value != "must hold" {
			t.Errorf("Description = %q, want %q", lit.Value, "must hold")
		}
	})

	t.Run("comparison expr is the whole condition", func(t *testing.T) {
		// `ASSERT(...) = 1 AS "..."` — the parenthesized subquery is the LHS of a
		// comparison that is the full assert expression (corpus assert.sql).
		a := parseAssert(t, `ASSERT (SELECT 1) = 1 AS "simple test"`)
		if _, ok := a.Expr.(*ast.CompareExpr); !ok {
			t.Errorf("Expr = %T, want *ast.CompareExpr", a.Expr)
		}
		if a.Description == nil {
			t.Error("Description = nil, want the AS string")
		}
	})

	t.Run("EXISTS subquery", func(t *testing.T) {
		parseAssert(t, "ASSERT NOT EXISTS(SELECT 1)")
	})

	t.Run("embedded subquery is filled for the lineage walker", func(t *testing.T) {
		// ASSERT carries a full expression; an embedded subquery must be re-parsed
		// (fillSubqueries) so the downstream query-span/lineage extractor can descend
		// into its *QueryStmt — not left as RawText with Query==nil.
		a := parseAssert(t, "ASSERT EXISTS(SELECT 1 FROM t)")
		assertSubqueriesFilled(t, a)
	})

	t.Run("param and sysvar conditions", func(t *testing.T) {
		parseAssert(t, `ASSERT @param = 1 AS "param test"`)
		parseAssert(t, `ASSERT @@sysvar = 1 AS "sysvar test"`)
	})

	t.Run("IN predicate", func(t *testing.T) {
		parseAssert(t, `ASSERT "123" IN ("123", "456")`)
	})

	t.Run("no space before paren", func(t *testing.T) {
		// `ASSERT(5 + a)` — the paren is part of the expression, not a call wrapper.
		parseAssert(t, "ASSERT(5 + a)")
	})
}

func TestAssert_Rejects(t *testing.T) {
	rejects := []string{
		"ASSERT",                  // missing expression
		"ASSERT AS 'x'",           // no expression before AS
		`ASSERT 1 AS notstring`,   // opt_description requires a string literal
		`ASSERT 1 AS 2`,           // ditto: integer is not a description string
		`ASSERT true AS 'd' junk`, // trailing junk after the description
		"ASSERT true extra",       // trailing token (no AS) — .g4 ends after the expr
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// ANALYZE
// ---------------------------------------------------------------------------

func TestAnalyze(t *testing.T) {
	t.Run("bare", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE")
		if a.Options != nil {
			t.Errorf("Options = %v, want nil", a.Options)
		}
		if a.Targets != nil {
			t.Errorf("Targets = %v, want nil", a.Targets)
		}
	})

	t.Run("single target", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE T")
		if len(a.Targets) != 1 {
			t.Fatalf("Targets = %d, want 1", len(a.Targets))
		}
		if pathString(a.Targets[0].Path) != "T" {
			t.Errorf("Targets[0].Path = %q, want T", pathString(a.Targets[0].Path))
		}
		if a.Targets[0].Columns != nil {
			t.Errorf("Columns = %v, want nil", a.Targets[0].Columns)
		}
	})

	t.Run("target with column list", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE T(a, b, c)")
		if len(a.Targets) != 1 {
			t.Fatalf("Targets = %d, want 1", len(a.Targets))
		}
		cols := a.Targets[0].Columns
		if len(cols) != 3 || cols[0] != "a" || cols[2] != "c" {
			t.Errorf("Columns = %v, want [a b c]", cols)
		}
	})

	t.Run("multiple targets", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE T, T2(a)")
		if len(a.Targets) != 2 {
			t.Fatalf("Targets = %d, want 2", len(a.Targets))
		}
		if pathString(a.Targets[1].Path) != "T2" || len(a.Targets[1].Columns) != 1 {
			t.Errorf("Targets[1] = %+v", a.Targets[1])
		}
	})

	t.Run("options then targets", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE OPTIONS(p1 = a1, p2 = a2) T(a, b, c)")
		if len(a.Options) != 2 {
			t.Fatalf("Options = %d, want 2", len(a.Options))
		}
		if len(a.Targets) != 1 || pathString(a.Targets[0].Path) != "T" {
			t.Errorf("Targets = %+v, want one T", a.Targets)
		}
	})

	t.Run("empty options no targets", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE OPTIONS()")
		if a.Options == nil {
			// distinct from "no OPTIONS clause": present-but-empty. We treat an
			// empty OPTIONS() as an empty (non-nil) slice so the clause's presence
			// is recorded.
			t.Error("Options = nil, want present-but-empty")
		}
		if len(a.Options) != 0 {
			t.Errorf("Options = %d entries, want 0", len(a.Options))
		}
		if a.Targets != nil {
			t.Errorf("Targets = %v, want nil", a.Targets)
		}
	})

	t.Run("options only", func(t *testing.T) {
		a := parseAnalyze(t, "ANALYZE OPTIONS(p1 = a1, p2 = a2)")
		if len(a.Options) != 2 {
			t.Fatalf("Options = %d, want 2", len(a.Options))
		}
		if a.Targets != nil {
			t.Errorf("Targets = %v, want nil", a.Targets)
		}
	})

	// Both BigQuery and Spanner shift `ANALYZE OPTIONS(...)` into opt_options_list
	// unconditionally, rejecting non-assignment entries like `ANALYZE OPTIONS(a, b, c)`
	// and `ANALYZE OPTIONS(zzz)`.
	t.Run("OPTIONS(...) requires assignment entries", func(t *testing.T) {
		for _, sql := range []string{"ANALYZE OPTIONS(a, b, c)", "ANALYZE OPTIONS(zzz)"} {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q): want error, got nil", sql)
			}
		}
	})

	t.Run("options-list then OPTIONS-as-table", func(t *testing.T) {
		// `ANALYZE OPTIONS(a = b) Options(a, b, c)` — first OPTIONS(a=b) is the
		// options-list; the second `Options(a,b,c)` is a table target.
		a := parseAnalyze(t, "ANALYZE OPTIONS(a = b) Options(a, b, c)")
		if len(a.Options) != 1 {
			t.Fatalf("Options = %d, want 1", len(a.Options))
		}
		if len(a.Targets) != 1 || pathString(a.Targets[0].Path) != "Options" {
			t.Errorf("Targets = %+v, want one Options", a.Targets)
		}
	})

	t.Run("target then OPTIONS-as-table", func(t *testing.T) {
		// `ANALYZE T, OPTIONS(a, b, c)` — no leading options-list (T is target 1),
		// OPTIONS is target 2 with a column list.
		a := parseAnalyze(t, "ANALYZE T, OPTIONS(a, b, c)")
		if a.Options != nil {
			t.Errorf("Options = %v, want nil", a.Options)
		}
		if len(a.Targets) != 2 {
			t.Fatalf("Targets = %d, want 2", len(a.Targets))
		}
		if pathString(a.Targets[1].Path) != "OPTIONS" || len(a.Targets[1].Columns) != 3 {
			t.Errorf("Targets[1] = %+v, want OPTIONS(a,b,c)", a.Targets[1])
		}
	})

	t.Run("target then bare OPTIONS table", func(t *testing.T) {
		// `ANALYZE T, OPTIONS` — OPTIONS as a bare table-name target (no column list).
		a := parseAnalyze(t, "ANALYZE T, OPTIONS")
		if len(a.Targets) != 2 || pathString(a.Targets[1].Path) != "OPTIONS" {
			t.Errorf("Targets = %+v, want T, OPTIONS", a.Targets)
		}
	})

	t.Run("dashed target path (BigQuery union)", func(t *testing.T) {
		// table_and_column_info: maybe_dashed_path_expression. A BigQuery dashed path
		// is accepted by the union grammar (Spanner syntax-rejects '-' in a table
		// name — divergence #85 class; non-authoritative on Spanner, omni follows
		// the union .g4).
		a := parseAnalyze(t, "ANALYZE my-project.ds.tbl")
		if len(a.Targets) != 1 || pathString(a.Targets[0].Path) != "my-project.ds.tbl" {
			t.Errorf("Targets = %+v, want one my-project.ds.tbl", a.Targets)
		}
	})
}

func TestAnalyze_Rejects(t *testing.T) {
	rejects := []string{
		"ANALYZE T,",          // trailing comma in target list
		"ANALYZE T(",          // unterminated column list
		"ANALYZE T(a,)",       // trailing comma in column list
		"ANALYZE T()",         // empty column list (column_list needs >= 1)
		"ANALYZE OPTIONS(a=)", // options entry missing value
		"ANALYZE , T",         // leading comma
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

func TestDescribe(t *testing.T) {
	t.Run("bare path", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE foo")
		if d.IsDesc {
			t.Error("IsDesc = true, want false (DESCRIBE)")
		}
		if d.ObjectType != "" {
			t.Errorf("ObjectType = %q, want empty", d.ObjectType)
		}
		if pathString(d.Path) != "foo" {
			t.Errorf("Path = %q, want foo", pathString(d.Path))
		}
		if d.FromPath != nil {
			t.Errorf("FromPath = %v, want nil", d.FromPath)
		}
	})

	t.Run("DESC abbreviation", func(t *testing.T) {
		d := parseDescribe(t, "DESC foo")
		if !d.IsDesc {
			t.Error("IsDesc = false, want true (DESC)")
		}
	})

	t.Run("dotted path", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE namespace.foo")
		if pathString(d.Path) != "namespace.foo" {
			t.Errorf("Path = %q, want namespace.foo", pathString(d.Path))
		}
	})

	t.Run("object type INDEX", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE INDEX myindex")
		if d.ObjectType != "INDEX" {
			t.Errorf("ObjectType = %q, want INDEX", d.ObjectType)
		}
		if pathString(d.Path) != "myindex" {
			t.Errorf("Path = %q, want myindex", pathString(d.Path))
		}
	})

	t.Run("object type TYPE with dotted path", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE TYPE mynamespace.mytype")
		if d.ObjectType != "TYPE" || pathString(d.Path) != "mynamespace.mytype" {
			t.Errorf("got ObjectType=%q Path=%q", d.ObjectType, pathString(d.Path))
		}
	})

	t.Run("backtick object type", func(t *testing.T) {
		// DESCRIBE `FUNCTION` myfunction — FUNCTION is a reserved-ish word given as
		// a backtick identifier object-type.
		d := parseDescribe(t, "DESCRIBE `FUNCTION` myfunction")
		if d.ObjectType != "FUNCTION" || pathString(d.Path) != "myfunction" {
			t.Errorf("got ObjectType=%q Path=%q", d.ObjectType, pathString(d.Path))
		}
	})

	t.Run("FROM path", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE foo FROM namespace")
		if pathString(d.Path) != "foo" {
			t.Errorf("Path = %q, want foo", pathString(d.Path))
		}
		if pathString(d.FromPath) != "namespace" {
			t.Errorf("FromPath = %q, want namespace", pathString(d.FromPath))
		}
	})

	t.Run("object type COLUMN with FROM dotted", func(t *testing.T) {
		d := parseDescribe(t, "DESCRIBE COLUMN foo FROM T.suffix")
		if d.ObjectType != "COLUMN" || pathString(d.Path) != "foo" || pathString(d.FromPath) != "T.suffix" {
			t.Errorf("got ObjectType=%q Path=%q FromPath=%q", d.ObjectType, pathString(d.Path), pathString(d.FromPath))
		}
	})

	t.Run("dashed path (BigQuery union)", func(t *testing.T) {
		// describe_info: maybe_slashed_or_dashed_path_expression — the dashed form is
		// part of the union grammar (Spanner's shallow recognizer accepts it too).
		d := parseDescribe(t, "DESCRIBE my-project.ds.tbl")
		if pathString(d.Path) != "my-project.ds.tbl" {
			t.Errorf("Path = %q, want my-project.ds.tbl", pathString(d.Path))
		}
	})
}

func TestDescribe_Rejects(t *testing.T) {
	rejects := []string{
		"DESCRIBE",             // describe_info requires a path
		"DESC",                 // ditto
		"DESCRIBE FROM s",      // no object path before FROM
		"DESCRIBE foo FROM",    // FROM with no path
		"DESCRIBE foo bar baz", // object-type + path consumes 2 words; a 3rd is junk
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// RENAME
// ---------------------------------------------------------------------------

func TestRename(t *testing.T) {
	t.Run("RENAME TABLE", func(t *testing.T) {
		r := parseRename(t, "RENAME TABLE a TO b")
		if r.ObjectType != "TABLE" {
			t.Errorf("ObjectType = %q, want TABLE", r.ObjectType)
		}
		if pathString(r.From) != "a" || pathString(r.To) != "b" {
			t.Errorf("got From=%q To=%q, want a -> b", pathString(r.From), pathString(r.To))
		}
	})

	t.Run("dotted paths", func(t *testing.T) {
		r := parseRename(t, "RENAME TABLE a.b TO c.d")
		if pathString(r.From) != "a.b" || pathString(r.To) != "c.d" {
			t.Errorf("got From=%q To=%q", pathString(r.From), pathString(r.To))
		}
	})

	t.Run("non-TABLE object kind", func(t *testing.T) {
		// .g4 rename_statement allows any object-kind identifier (Spanner narrows
		// to TABLE; the union grammar does not — flagged divergence).
		r := parseRename(t, "RENAME VIEW v1 TO v2")
		if r.ObjectType != "VIEW" {
			t.Errorf("ObjectType = %q, want VIEW", r.ObjectType)
		}
	})
}

func TestRename_Rejects(t *testing.T) {
	rejects := []string{
		"RENAME TABLE a",                    // missing TO + target
		"RENAME TABLE a TO",                 // missing target
		"RENAME a TO b",                     // missing object-kind: identifier then path needs 2 names
		"RENAME TABLE TO b",                 // missing source path
		"RENAME",                            // bare
		"RENAME TABLE a TO b c",             // trailing junk
		"RENAME TABLE proj-1.a TO proj-1.b", // rename_statement uses path_expression (NOT dashed) — '-' rejects (matches .g4 + Spanner)
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// CALL
// ---------------------------------------------------------------------------

func TestCall(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure()")
		if pathString(c.Proc) != "myprocedure" {
			t.Errorf("Proc = %q, want myprocedure", pathString(c.Proc))
		}
		if len(c.Args) != 0 {
			t.Errorf("Args = %d, want 0", len(c.Args))
		}
	})

	t.Run("dotted name", func(t *testing.T) {
		c := parseCall(t, "CALL schema.myprocedure()")
		if pathString(c.Proc) != "schema.myprocedure" {
			t.Errorf("Proc = %q, want schema.myprocedure", pathString(c.Proc))
		}
	})

	t.Run("expression args", func(t *testing.T) {
		c := parseCall(t, `CALL myprocedure(1 + 2, "a", CAST(NULL AS string))`)
		if len(c.Args) != 3 {
			t.Fatalf("Args = %d, want 3", len(c.Args))
		}
	})

	t.Run("MODEL arg", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure(MODEL my.model)")
		arg, ok := c.Args[0].(*ast.CallArg)
		if !ok {
			t.Fatalf("Args[0] = %T, want *ast.CallArg", c.Args[0])
		}
		if arg.Kind != ast.CallArgModel || pathString(arg.Path) != "my.model" {
			t.Errorf("got Kind=%v Path=%q, want MODEL my.model", arg.Kind, pathString(arg.Path))
		}
	})

	t.Run("CONNECTION DEFAULT", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure(CONNECTION DEFAULT)")
		arg := c.Args[0].(*ast.CallArg)
		if arg.Kind != ast.CallArgConnection || !arg.Default {
			t.Errorf("got Kind=%v Default=%v, want CONNECTION DEFAULT", arg.Kind, arg.Default)
		}
	})

	t.Run("CONNECTION path", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure(CONNECTION my.connection)")
		arg := c.Args[0].(*ast.CallArg)
		if arg.Kind != ast.CallArgConnection || pathString(arg.Path) != "my.connection" {
			t.Errorf("got Kind=%v Path=%q", arg.Kind, pathString(arg.Path))
		}
	})

	t.Run("TABLE + subquery + nested func args", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure(TABLE my.table, (SELECT * FROM my.another_table), mytvf(1, 2))")
		if len(c.Args) != 3 {
			t.Fatalf("Args = %d, want 3", len(c.Args))
		}
		ta, ok := c.Args[0].(*ast.CallArg)
		if !ok || ta.Kind != ast.CallArgTable || pathString(ta.Path) != "my.table" {
			t.Errorf("Args[0] = %+v, want TABLE my.table", c.Args[0])
		}
		// Args[1] is a parenthesized subquery expression; Args[2] is a function call.
		if _, ok := c.Args[2].(*ast.FuncCall); !ok {
			t.Errorf("Args[2] = %T, want *ast.FuncCall", c.Args[2])
		}
		// The embedded subquery arg must be re-parsed (fillSubqueries) so the
		// query-span/lineage extractor can descend into it.
		assertSubqueriesFilled(t, c)
	})

	t.Run("named argument", func(t *testing.T) {
		c := parseCall(t, "CALL ns.proc(a => 1, b => 2)")
		if len(c.Args) != 2 {
			t.Fatalf("Args = %d, want 2", len(c.Args))
		}
		na, ok := c.Args[0].(*ast.NamedArg)
		if !ok || na.Name != "a" {
			t.Errorf("Args[0] = %+v, want NamedArg a", c.Args[0])
		}
	})

	t.Run("named argument with lambda value", func(t *testing.T) {
		// named_argument: identifier '=>' (expression | lambda_argument). The value
		// may be a single-id or parenthesized lambda (oracle: both accept).
		for _, sql := range []string{
			"CALL p(f => x -> x)",
			"CALL p(f => (x) -> x)",
			"CALL p(f => (x, y) -> x + y)",
		} {
			c := parseCall(t, sql)
			na, ok := c.Args[0].(*ast.NamedArg)
			if !ok {
				t.Fatalf("%q: Args[0] = %T, want *ast.NamedArg", sql, c.Args[0])
			}
			if _, ok := na.Value.(*ast.LambdaExpr); !ok {
				t.Errorf("%q: NamedArg.Value = %T, want *ast.LambdaExpr", sql, na.Value)
			}
		}
	})

	t.Run("param and sysvar args", func(t *testing.T) {
		parseCall(t, "CALL myprocedure(@test_param_bool)")
		parseCall(t, "CALL myprocedure(@@sysvar)")
	})

	t.Run("DESCRIPTOR arg", func(t *testing.T) {
		c := parseCall(t, "CALL myprocedure(DESCRIPTOR(a, b))")
		arg, ok := c.Args[0].(*ast.CallArg)
		if !ok || arg.Kind != ast.CallArgDescriptor {
			t.Fatalf("Args[0] = %+v, want DESCRIPTOR", c.Args[0])
		}
		if len(arg.Columns) != 2 || arg.Columns[0] != "a" || arg.Columns[1] != "b" {
			t.Errorf("Columns = %v, want [a b]", arg.Columns)
		}
	})

	// The clause keywords (TABLE/MODEL/CONNECTION/DESCRIPTOR) are non-reserved, so
	// a bare keyword or a field access on it is an identifier EXPRESSION, not a
	// clause (oracle: the Spanner emulator — authoritative for CALL — accepts
	// `CALL p(model)`, `CALL p(table.x)`, `CALL p(descriptor)`). The keyword opens
	// its clause only when the next token begins the clause body.
	t.Run("keyword as bare identifier expression", func(t *testing.T) {
		for _, sql := range []string{
			"CALL p(model)",
			"CALL p(table)",
			"CALL p(connection)",
			"CALL p(descriptor)",
		} {
			c := parseCall(t, sql)
			if len(c.Args) != 1 {
				t.Fatalf("%q: Args = %d, want 1", sql, len(c.Args))
			}
			if _, isClause := c.Args[0].(*ast.CallArg); isClause {
				t.Errorf("%q: Args[0] is a *CallArg clause, want an identifier expression", sql)
			}
		}
	})

	// named_argument: `identifier => input_table_argument` and
	// named_non_expr_argument: `identifier => table_clause` — the two non-expression
	// value shapes a named tvf_argument has besides an expression / lambda.
	t.Run("positional input table", func(t *testing.T) {
		c := parseCall(t, "CALL p(INPUT TABLE, 1)")
		if len(c.Args) != 2 {
			t.Fatalf("Args = %d, want 2", len(c.Args))
		}
		ca, ok := c.Args[0].(*ast.CallArg)
		if !ok || ca.Kind != ast.CallArgInputTable {
			t.Errorf("Args[0] = %+v, want an INPUT TABLE CallArg", c.Args[0])
		}
		// A bare INPUT is the non-reserved keyword as an identifier expression.
		if c := parseCall(t, "CALL p(INPUT)"); len(c.Args) != 1 {
			t.Errorf("Args = %d, want 1", len(c.Args))
		} else if _, isClause := c.Args[0].(*ast.CallArg); isClause {
			t.Error("CALL p(INPUT): Args[0] is a clause, want an identifier expression")
		}
	})

	t.Run("named input table and table clause", func(t *testing.T) {
		c := parseCall(t, "CALL p(a => INPUT TABLE, b => TABLE t, c => 1)")
		if len(c.Args) != 3 {
			t.Fatalf("Args = %d, want 3", len(c.Args))
		}
		na, ok := c.Args[0].(*ast.NamedArg)
		if !ok {
			t.Fatalf("Args[0] = %T, want *ast.NamedArg", c.Args[0])
		}
		if ca, ok := na.Value.(*ast.CallArg); !ok || ca.Kind != ast.CallArgInputTable {
			t.Errorf("Args[0] value = %+v, want an INPUT TABLE CallArg", na.Value)
		}
		// The table_clause value is modelled as the clause itself, as a TVF's
		// positional `TABLE t` argument is (named_non_expr_argument wraps it in a
		// subquery expression).
		na, ok = c.Args[1].(*ast.NamedArg)
		if !ok {
			t.Fatalf("Args[1] = %T, want *ast.NamedArg", c.Args[1])
		}
		if _, ok := na.Value.(*ast.TableClause); !ok {
			t.Errorf("Args[1] value = %T, want *ast.TableClause", na.Value)
		}
	})

	t.Run("keyword field-access expression", func(t *testing.T) {
		for _, sql := range []string{
			"CALL p(model.col)",
			"CALL p(table.x)",
			"CALL p(connection.x)",
		} {
			c := parseCall(t, sql)
			if _, isClause := c.Args[0].(*ast.CallArg); isClause {
				t.Errorf("%q: Args[0] is a *CallArg clause, want a field-access expression", sql)
			}
		}
	})
}

func TestCall_Rejects(t *testing.T) {
	rejects := []string{
		"CALL",                // missing name
		"CALL proc",           // missing '(' arg list
		"CALL proc(",          // unterminated arg list
		"CALL proc(,)",        // leading comma
		"CALL proc(1,)",       // trailing comma
		"CALL proc(1 => 2)",   // named-arg name must be an identifier, not an expr
		"CALL proc() extra",   // trailing junk after the call
		"CALL my-proc.foo()",  // call_statement uses path_expression (NOT dashed) — '-' rejects (matches .g4 + Spanner)
		"CALL p(x -> x)",      // a bare POSITIONAL lambda is not a tvf_argument (oracle rejects at '->')
		"CALL p(arr, e -> e)", // ditto, positional lambda in a later arg
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// ---------------------------------------------------------------------------
// Completeness gate — the canonical ZetaSQL corpus
// ---------------------------------------------------------------------------

// TestUtility_CorpusAccepts parses every statement from the canonical ZetaSQL
// parser testdata files (assert.sql / analyze.sql / describe.sql / call.sql,
// lifted inline) and asserts each parses to a single node of the right kind with
// no errors. The legacy ANTLR grammar bytebase consumes is a hand-port of that
// ZetaSQL reference, so this corpus is the BREADTH oracle for the node — the
// completeness gate (correctness-protocol.md), authoritative for the
// BigQuery-union forms (ANALYZE targets, DESCRIBE object-types, CALL clauses)
// that the Spanner emulator narrows / cannot adjudicate. Each line is a `;`-
// terminated statement in the corpus; the terminator is stripped here.
func TestUtility_CorpusAccepts(t *testing.T) {
	type want int
	const (
		wAssert want = iota
		wAnalyze
		wDescribe
		wCall
	)
	cases := []struct {
		sql  string
		kind want
	}{
		// assert.sql
		{"ASSERT TRUE", wAssert},
		{"ASSERT(5 + a)", wAssert},
		{`ASSERT(SELECT 1) = 1 AS "simple test"`, wAssert},
		{"ASSERT NOT EXISTS(SELECT 1)", wAssert},
		{`ASSERT @param = 1 AS "param test"`, wAssert},
		{`ASSERT @@sysvar = 1 AS "sysvar test"`, wAssert},
		{`ASSERT "123" IN ("123", "456")`, wAssert},
		{`ASSERT IS_NAN(NULL) OR ENDS_WITH("suffix", "fix") AS "abc"`, wAssert},
		{`ASSERT "123" IS NOT NULL`, wAssert},
		{"ASSERT CASE TRUE WHEN TRUE THEN FALSE END", wAssert},
		{"ASSERT 123 BETWEEN 1 AND 456", wAssert},
		{"ASSERT(SELECT IS_NAN(NAN))", wAssert},
		{"ASSERT(ASSERT((SELECT TRUE)))", wAssert},
		// analyze.sql
		{"ANALYZE OPTIONS(p1 = a1, p2 = a2) T(a, b, c)", wAnalyze},
		{"ANALYZE OPTIONS(p1 = a1, p2 = a2) T(a, b, c), T2(a, b, c)", wAnalyze},
		{"ANALYZE OPTIONS(p1 = a1, p2 = a2) T, T2(a, b, c)", wAnalyze},
		{"ANALYZE T(a, b, c)", wAnalyze},
		{"ANALYZE T(a, b, c), T2(a)", wAnalyze},
		{"ANALYZE T", wAnalyze},
		{"ANALYZE T, T2", wAnalyze},
		{"ANALYZE OPTIONS()", wAnalyze},
		{"ANALYZE T, OPTIONS", wAnalyze},
		{"ANALYZE OPTIONS(a = b) Options(a, b, c)", wAnalyze},
		{"ANALYZE T, OPTIONS(a, b, c)", wAnalyze},
		{"ANALYZE OPTIONS(p1 = a1, p2 = a2)", wAnalyze},
		// describe.sql
		{"DESCRIBE foo", wDescribe},
		{"DESCRIBE namespace.foo", wDescribe},
		{"DESCRIBE INDEX myindex", wDescribe},
		{"DESCRIBE INDEX mynamespace.myindex", wDescribe},
		{"DESCRIBE `FUNCTION` myfunction", wDescribe},
		{"DESCRIBE `FUNCTION` mynamespace.myfunction", wDescribe},
		{"DESCRIBE TVF mytvf", wDescribe},
		{"DESCRIBE TVF mynamespace.mytvf", wDescribe},
		{"DESCRIBE TYPE mytype", wDescribe},
		{"DESCRIBE TYPE mynamespace.mytype", wDescribe},
		{"DESCRIBE foo FROM namespace", wDescribe},
		{"DESCRIBE COLUMN foo FROM T.suffix", wDescribe},
		{"DESCRIBE TYPE prefixed.name FROM Catalog.`With`.Dots", wDescribe},
		// call.sql
		{"CALL myprocedure()", wCall},
		{"CALL schema.myprocedure()", wCall},
		{`CALL myprocedure(1 + 2, "a", CAST(NULL AS string))`, wCall},
		{"CALL myprocedure(MODEL my.model)", wCall},
		{"CALL myprocedure(CONNECTION DEFAULT)", wCall},
		{"CALL myprocedure(CONNECTION my.connection)", wCall},
		{"CALL myprocedure(TABLE my.table, (SELECT * FROM my.another_table), mytvf(1, 2))", wCall},
		{"CALL myprocedure(@test_param_bool)", wCall},
		{"CALL myprocedure(@@sysvar)", wCall},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOneStmt(t, tc.sql)
			var ok bool
			switch tc.kind {
			case wAssert:
				_, ok = node.(*ast.AssertStmt)
			case wAnalyze:
				_, ok = node.(*ast.AnalyzeStmt)
			case wDescribe:
				_, ok = node.(*ast.DescribeStmt)
			case wCall:
				_, ok = node.(*ast.CallStmt)
			}
			if !ok {
				t.Errorf("parseForTest(%q): got %T, want kind %d", tc.sql, node, tc.kind)
			}
		})
	}
}

// assertSubqueriesFilled walks node and fails if any embedded SubqueryExpr /
// ExistsExpr / ArraySubqueryExpr still has a nil Query (i.e. was left as
// unresolved RawText). It also requires at least one subquery to be present, so
// the test is meaningful. The query-span / lineage extractor walks the AST and
// needs each subquery's inner *QueryStmt to resolve its tables/columns.
func assertSubqueriesFilled(t *testing.T, node ast.Node) {
	t.Helper()
	var total, unfilled int
	ast.Inspect(node, func(n ast.Node) bool {
		switch sq := n.(type) {
		case *ast.SubqueryExpr:
			total++
			if sq.Query == nil {
				unfilled++
			}
		case *ast.ExistsExpr:
			total++
			if sq.Query == nil {
				unfilled++
			}
		case *ast.ArraySubqueryExpr:
			total++
			if sq.Query == nil {
				unfilled++
			}
		}
		return true
	})
	if total == 0 {
		t.Error("no embedded subquery found; the fixture must contain one for this check to be meaningful")
	}
	if unfilled != 0 {
		t.Errorf("%d of %d embedded subqueries left unfilled (Query==nil); the lineage walker cannot descend into them", unfilled, total)
	}
}

// ---------------------------------------------------------------------------
// SHOW
// ---------------------------------------------------------------------------

func parseShow(t *testing.T, sql string) *ast.ShowStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	s, ok := n.(*ast.ShowStmt)
	if !ok {
		t.Fatalf("parseOneStmt(%q): got %T, want *ast.ShowStmt", sql, n)
	}
	return s
}

// TestParseShow covers show_statement's alternatives. The whole family is
// unjudged by the Spanner emulator (it answers "Statement not supported:
// ShowStmt" from dispatch and accepts `SHOW a zzz`), so the grammar and the
// canonical corpus (testdata/legacy/zetasql/parser/testdata/show.sql) govern.
func TestParseShow(t *testing.T) {
	s := parseShow(t, "SHOW TABLES")
	if s.Target != "TABLES" {
		t.Errorf("Target = %q, want TABLES", s.Target)
	}
	if s.From != nil || s.Like != nil {
		t.Errorf("From / Like = (%+v, %+v), want nil", s.From, s.Like)
	}
	if s.Loc.Start != 0 || s.Loc.End != len("SHOW TABLES") {
		t.Errorf("Loc = %+v, want (0, %d)", s.Loc, len("SHOW TABLES"))
	}

	// FROM, LIKE, and both together.
	s = parseShow(t, "SHOW TABLES FROM catalog.adwords LIKE 'KitchenSync%'")
	if s.From == nil || len(s.From.Parts) != 2 {
		t.Errorf("From = %+v, want catalog.adwords", s.From)
	}
	lit, ok := s.Like.(*ast.Literal)
	if !ok || lit.Value != "KitchenSync%" {
		t.Errorf("Like = %+v, want the KitchenSync%% literal", s.Like)
	}

	// show_target's two-word alternative, and the backtick-quoted spelling of the
	// same name (the `identifier` alternative).
	if got := parseShow(t, "SHOW MATERIALIZED VIEWS").Target; got != "MATERIALIZED VIEWS" {
		t.Errorf("Target = %q, want MATERIALIZED VIEWS", got)
	}
	if got := parseShow(t, "SHOW `MATERIALIZED VIEWS`").Target; got != "MATERIALIZED VIEWS" {
		t.Errorf("Target = %q, want MATERIALIZED VIEWS", got)
	}

	// from_path_expression takes maybe_slashed_or_dashed_path_expression, so a
	// dashed BigQuery project is one path component.
	if got := parseShow(t, "SHOW TABLES FROM my-project.ds"); len(got.From.Parts) != 2 || got.From.Parts[0] != "my-project" {
		t.Errorf("From = %+v, want [my-project ds]", got.From.Parts)
	}

	// A SHOW may carry a pipe suffix (sql_statement_body_maybe_pipe_suffix).
	n := parseOneStmt(t, "SHOW TABLES |> WHERE a > 1")
	if _, ok := n.(*ast.PipeStatement); !ok {
		t.Errorf("got %T, want *ast.PipeStatement for a SHOW with a pipe suffix", n)
	}
}

func TestParseShow_Rejects(t *testing.T) {
	for _, sql := range []string{
		"SHOW",
		"SHOW a zzz",
		"SHOW a FROM",
		"SHOW a LIKE",
		"SHOW a LIKE 1", // LIKE takes a string literal
	} {
		assertReject(t, sql)
	}
}

// ---------------------------------------------------------------------------
// IMPORT / MODULE
// ---------------------------------------------------------------------------

func parseImport(t *testing.T, sql string) *ast.ImportStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	s, ok := n.(*ast.ImportStmt)
	if !ok {
		t.Fatalf("parseOneStmt(%q): got %T, want *ast.ImportStmt", sql, n)
	}
	return s
}

func parseModule(t *testing.T, sql string) *ast.ModuleStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	s, ok := n.(*ast.ModuleStmt)
	if !ok {
		t.Fatalf("parseOneStmt(%q): got %T, want *ast.ModuleStmt", sql, n)
	}
	return s
}

// TestParseImport covers import_statement's alternatives. Like SHOW, the family
// is unjudged by the emulator ("Statement not supported: ImportStmt"), so the
// grammar and the canonical corpus (modules.sql) govern.
func TestParseImport(t *testing.T) {
	s := parseImport(t, "IMPORT MODULE foo.bar.baz")
	if s.Kind != ast.ImportModule {
		t.Errorf("Kind = %v, want ImportModule", s.Kind)
	}
	if path, ok := s.Name.(*ast.PathExpr); !ok || len(path.Parts) != 3 {
		t.Errorf("Name = %+v, want the foo.bar.baz path", s.Name)
	}
	if s.Alias != "" || s.AliasKeyword != "" || s.Options != nil {
		t.Errorf("Alias / keyword / Options = (%q, %q, %+v), want empty", s.Alias, s.AliasKeyword, s.Options)
	}

	// import_type PROTO, and a string literal module name.
	if s := parseImport(t, "IMPORT PROTO a.b"); s.Kind != ast.ImportProto {
		t.Errorf("Kind = %v, want ImportProto", s.Kind)
	}
	if s := parseImport(t, `IMPORT MODULE "file.proto"`); s.Name == nil {
		t.Error("Name = nil, want the string literal")
	}

	// opt_as_or_into_alias: both spellings. The alias is an identifier, and the
	// corpus imports under keywords too (`IMPORT MODULE foo AS TABLE`).
	for _, c := range []struct{ sql, kw, alias string }{
		{"IMPORT MODULE a.b AS alias", "AS", "alias"},
		{"IMPORT MODULE name.path INTO a", "INTO", "a"},
		{"IMPORT MODULE foo AS TABLE", "AS", "TABLE"},
	} {
		s := parseImport(t, c.sql)
		if s.Alias != c.alias || s.AliasKeyword != c.kw {
			t.Errorf("%q: alias = (%q, %q), want (%q, %q)", c.sql, s.Alias, s.AliasKeyword, c.alias, c.kw)
		}
	}

	// opt_options_list. An option value may embed a subquery, which the dispatch
	// fills (fillSubqueries) and a malformed one fails the strict parse.
	s = parseImport(t, "IMPORT MODULE a.foo AS foo OPTIONS(a = 1)")
	if len(s.Options) != 1 {
		t.Errorf("Options = %+v, want 1 entry", s.Options)
	}
	assertSubqueriesFilled(t, parseImport(t, "IMPORT MODULE a OPTIONS(a = (SELECT 1))"))
	assertReject(t, "IMPORT MODULE a OPTIONS(a = (SELECT 1 FROM t x y))")
	if s.Loc.End != len("IMPORT MODULE a.foo AS foo OPTIONS(a = 1)") {
		t.Errorf("Loc.End = %d, want the end of the statement", s.Loc.End)
	}
}

func TestParseImport_Rejects(t *testing.T) {
	for _, sql := range []string{
		"IMPORT a",          // import_type is MODULE | PROTO
		"IMPORT MODULE",     // the module name is required
		"IMPORT MODULE a b", // an alias needs AS or INTO
		"IMPORT MODULE a AS",
	} {
		assertReject(t, sql)
	}
}

func TestParseModule(t *testing.T) {
	s := parseModule(t, "MODULE foo.bar")
	if s.Name == nil || len(s.Name.Parts) != 2 {
		t.Errorf("Name = %+v, want foo.bar", s.Name)
	}
	if s.Options != nil {
		t.Errorf("Options = %+v, want nil", s.Options)
	}

	// A path component may be a non-reserved keyword (`MODULE import`), and an
	// empty OPTIONS() is accepted.
	if s := parseModule(t, "MODULE import.module.options OPTIONS(a = 1)"); len(s.Name.Parts) != 3 || len(s.Options) != 1 {
		t.Errorf("Name / Options = (%+v, %+v), want 3 parts and 1 entry", s.Name.Parts, s.Options)
	}
	parseModule(t, "MODULE a OPTIONS()")
	assertSubqueriesFilled(t, parseModule(t, "MODULE a OPTIONS(x = (SELECT 1))"))
	assertReject(t, "MODULE a OPTIONS(x = (SELECT 1 FROM t x y))")

	for _, sql := range []string{"MODULE", "MODULE a()", "MODULE a OPTIONS("} {
		assertReject(t, sql)
	}
}
