package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-pipe node (GoogleSQL's pipe syntax: the suffix
// statement, a query's suffix, a statement's suffix, and the pipe operators).
// The accept/reject behavior is also recorded in testdata/compat.tsv; those
// rows judge the pinned Spanner emulator, which answers "Statement not
// supported" for a pipe statement without deciding the grammar (see
// testdata/README.md), so these tests hold the grammar itself: the statement
// form `|> AGGREGATE` is a complete statement, each operator parses, a
// malformed operator rejects, and each node's Loc covers its text.

// firstStmt parses sql and returns its single statement.
func firstStmt(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d statements, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

// subpipelineOf returns the Pipeline of a statement form that carries one.
func subpipelineOf(t *testing.T, node ast.Node) *ast.Subpipeline {
	t.Helper()
	switch n := node.(type) {
	case *ast.QueryStmt:
		if n.Pipe == nil {
			t.Fatalf("query has no pipe suffix")
		}
		return n.Pipe
	case *ast.SubpipelineStmt:
		return n.Pipeline
	case *ast.PipeStatement:
		return n.Pipeline
	default:
		t.Fatalf("statement is %T, want a pipe statement", node)
		return nil
	}
}

// TestPipe_StatementForms holds the three ways the grammar reaches a pipe
// suffix. The first is the fragment question: `|> AGGREGATE` alone is a
// statement (sql_statement_body_no_pipe_suffix → subpipeline_statement →
// subpipeline_no_parens), not a fragment.
func TestPipe_StatementForms(t *testing.T) {
	t.Run("suffix alone is a statement", func(t *testing.T) {
		stmt, ok := firstStmt(t, "|> AGGREGATE").(*ast.SubpipelineStmt)
		if !ok {
			t.Fatalf("statement is not an *ast.SubpipelineStmt")
		}
		if len(stmt.Pipeline.Ops) != 1 || stmt.Pipeline.Ops[0].Tag() != ast.T_PipeAggregate {
			t.Fatalf("operators = %v, want one PipeAggregate", opTags(stmt.Pipeline))
		}
		if stmt.Loc.Start != 0 || stmt.Loc.End != len("|> AGGREGATE") {
			t.Errorf("statement Loc = %+v, want [0,%d)", stmt.Loc, len("|> AGGREGATE"))
		}
	})

	t.Run("query suffix belongs to the query", func(t *testing.T) {
		q, ok := firstStmt(t, "SELECT x FROM t |> WHERE x > 1 |> ORDER BY x").(*ast.QueryStmt)
		if !ok {
			t.Fatalf("statement is not an *ast.QueryStmt")
		}
		if q.Pipe == nil {
			t.Fatalf("query has no pipe suffix")
		}
		want := []ast.NodeTag{ast.T_PipeWhere, ast.T_PipeOrderBy}
		if got := opTags(q.Pipe); len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
			t.Errorf("operators = %v, want %v", got, want)
		}
	})

	t.Run("table-returning statement suffix", func(t *testing.T) {
		stmt, ok := firstStmt(t, "DESCRIBE t |> AGGREGATE COUNT(*)").(*ast.PipeStatement)
		if !ok {
			t.Fatalf("statement is not an *ast.PipeStatement")
		}
		if stmt.Head.Tag() != ast.T_DescribeStmt {
			t.Errorf("head is %v, want DescribeStmt", stmt.Head.Tag())
		}
		if stmt.Loc.Start != 0 || stmt.Loc.End != len("DESCRIBE t |> AGGREGATE COUNT(*)") {
			t.Errorf("statement Loc = %+v, want [0,%d)", stmt.Loc, len("DESCRIBE t |> AGGREGATE COUNT(*)"))
		}
	})
}

// opTags lists the tags of a subpipeline's operators.
func opTags(p *ast.Subpipeline) []ast.NodeTag {
	tags := make([]ast.NodeTag, len(p.Ops))
	for i, op := range p.Ops {
		tags[i] = op.Tag()
	}
	return tags
}

// TestPipe_Operators parses one statement per pipe operator. Each is expected
// to be a single operator of the named node type.
func TestPipe_Operators(t *testing.T) {
	cases := []struct {
		sql  string
		want ast.NodeTag
	}{
		{"|> WHERE x > 1", ast.T_PipeWhere},
		{"|> SELECT x, y", ast.T_PipeSelect},
		{"|> EXTEND x AS y", ast.T_PipeExtend},
		{"|> RENAME a AS b", ast.T_PipeRename},
		{"|> AGGREGATE COUNT(*) AS n GROUP BY x", ast.T_PipeAggregate},
		{"|> LIMIT 10 OFFSET 2", ast.T_PipeLimitOffset},
		{"|> ORDER BY x DESC", ast.T_PipeOrderBy},
		{"|> UNION ALL (SELECT 1)", ast.T_PipeSetOperation},
		{"|> RECURSIVE UNION ALL (|> AGGREGATE)", ast.T_PipeRecursiveUnion},
		{"|> JOIN a ON true", ast.T_PipeJoin},
		{"|> CALL f() AS g", ast.T_PipeCall},
		{"|> WINDOW x AS y", ast.T_PipeWindow},
		{"|> DISTINCT", ast.T_PipeDistinct},
		{"|> TABLESAMPLE SYSTEM (10 PERCENT)", ast.T_PipeTablesample},
		{"|> AS t", ast.T_PipeAs},
		{"|> DESCRIBE", ast.T_PipeDescribe},
		{"|> STATIC_DESCRIBE", ast.T_PipeStaticDescribe},
		{"|> ASSERT x > 0", ast.T_PipeAssert},
		{"|> LOG(|> AGGREGATE)", ast.T_PipeLog},
		{"|> DROP a, b", ast.T_PipeDrop},
		{"|> SET a = 1", ast.T_PipeSet},
		{"|> PIVOT(SUM(x) FOR y IN (1))", ast.T_PipePivot},
		{"|> UNPIVOT(a FOR b IN (a))", ast.T_PipeUnpivot},
		{"|> IF true THEN (|> AGGREGATE) ELSE (|> DISTINCT)", ast.T_PipeIf},
		{"|> FORK (|> AGGREGATE), (|> DISTINCT)", ast.T_PipeFork},
		{"|> TEE (|> AGGREGATE)", ast.T_PipeTee},
		{"|> WITH c AS (SELECT 1)", ast.T_PipeWith},
		{"|> EXPORT DATA OPTIONS(uri = 'x')", ast.T_PipeExportData},
		{"|> CREATE TABLE t (x INT64)", ast.T_PipeCreateTable},
		{"|> INSERT INTO t (x) ON CONFLICT DO NOTHING", ast.T_PipeInsert},
		{"|> MATCH_RECOGNIZE(ORDER BY x MEASURES 1 AS m PATTERN(a) DEFINE a AS true)", ast.T_PipeMatchRecognize},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			pipe := subpipelineOf(t, firstStmt(t, tc.sql))
			if len(pipe.Ops) != 1 {
				t.Fatalf("got %d operators, want 1", len(pipe.Ops))
			}
			if got := pipe.Ops[0].Tag(); got != tc.want {
				t.Fatalf("operator is %v, want %v", got, tc.want)
			}
		})
	}
}

// TestPipe_OperatorShapes checks the fields a few operators are parsed into,
// not just their types.
func TestPipe_OperatorShapes(t *testing.T) {
	t.Run("aggregate item order and group by", func(t *testing.T) {
		pipe := subpipelineOf(t, firstStmt(t, "|> AGGREGATE x AS y DESC NULLS LAST GROUP BY z"))
		agg := pipe.Ops[0].(*ast.PipeAggregate)
		if len(agg.Select.Items) != 1 {
			t.Fatalf("got %d aggregate items, want 1", len(agg.Select.Items))
		}
		item := agg.Select.Items[0]
		if item.Alias != "y" || item.Order == nil || !item.Order.HasDir || !item.Order.Desc {
			t.Errorf("item = %+v (order %+v), want alias y descending", item, item.Order)
		}
		if item.Order.NullsFirst == nil || *item.Order.NullsFirst {
			t.Errorf("NULLS LAST not recorded: %+v", item.Order)
		}
		if agg.Select.GroupBy == nil || len(agg.Select.GroupBy.Items) != 1 {
			t.Errorf("GROUP BY = %+v, want one item", agg.Select.GroupBy)
		}
	})

	t.Run("limit all", func(t *testing.T) {
		pipe := subpipelineOf(t, firstStmt(t, "|> LIMIT ALL OFFSET 3"))
		lim := pipe.Ops[0].(*ast.PipeLimitOffset)
		if !lim.All || lim.Limit != nil || lim.Offset == nil {
			t.Errorf("limit = %+v, want LIMIT ALL OFFSET 3", lim)
		}
	})

	t.Run("set operation metadata", func(t *testing.T) {
		pipe := subpipelineOf(t, firstStmt(t, "|> UNION ALL TABLE a, (SELECT 1)"))
		set := pipe.Ops[0].(*ast.PipeSetOperation)
		if set.Set.Op != ast.SetOpUnion || !set.Set.All {
			t.Errorf("set = %+v, want UNION ALL", set.Set)
		}
		if len(set.Right) != 2 {
			t.Fatalf("got %d operands, want 2", len(set.Right))
		}
		if set.Right[0].Tag() != ast.T_TableClause || set.Right[1].Tag() != ast.T_QueryStmt {
			t.Errorf("operands = %v, %v; want TableClause then QueryStmt", set.Right[0].Tag(), set.Right[1].Tag())
		}
	})

	t.Run("recursive union alias", func(t *testing.T) {
		pipe := subpipelineOf(t, firstStmt(t, "|> RECURSIVE UNION ALL WITH DEPTH MAX 2 (|> DISTINCT) AS r"))
		ru := pipe.Ops[0].(*ast.PipeRecursiveUnion)
		if ru.Alias != "r" || ru.Depth == nil || ru.Depth.Max == nil {
			t.Errorf("recursive union = %+v (depth %+v), want alias r and MAX 2", ru, ru.Depth)
		}
		if _, ok := ru.RHS.(*ast.Subpipeline); !ok {
			t.Errorf("rhs is %T, want *ast.Subpipeline", ru.RHS)
		}
	})

	t.Run("row pattern parentheses are in the Loc", func(t *testing.T) {
		sql := "|> MATCH_RECOGNIZE(ORDER BY x MEASURES 1 AS m PATTERN((a)) DEFINE a AS true)"
		pipe := subpipelineOf(t, firstStmt(t, sql))
		pat := pipe.Ops[0].(*ast.PipeMatchRecognize).Clause.Pattern
		want := strings.Index(sql, "PATTERN(") + len("PATTERN(")
		if pat.Loc.Start != want {
			t.Errorf("pattern Loc = %+v, want it to start at the outer '(' (%d)", pat.Loc, want)
		}
	})

	t.Run("insert has no source", func(t *testing.T) {
		pipe := subpipelineOf(t, firstStmt(t, "|> INSERT INTO t (x) ON CONFLICT DO NOTHING"))
		ins := pipe.Ops[0].(*ast.PipeInsert).Insert
		if ins.Rows != nil || ins.Query != nil || ins.TableClause != nil {
			t.Errorf("pipe INSERT carries a data source: %+v", ins)
		}
		if ins.OnConflict == nil || !ins.OnConflict.DoNothing {
			t.Errorf("ON CONFLICT = %+v, want DO NOTHING", ins.OnConflict)
		}
	})
}

// TestPipe_RejectsFragments holds inputs the grammar refuses. None is a
// statement: `|>` alone, an operator with its operand missing, an unknown
// operator, and the operator-level custom errors for a misplaced GROUP/ELSE/
// ELSEIF or a path in RENAME/SET/DROP.
func TestPipe_RejectsFragments(t *testing.T) {
	reject := []string{
		"|>",
		"|> ALL",
		"|> WHERE",
		"|> SELECT",
		"|> EXTEND",
		"|> RENAME a",
		"|> GROUP BY x",
		"|> LIMIT",
		"|> ORDER BY",
		"|> ORDER",
		"|> JOIN",
		"|> CALL()",
		"|> WINDOW",
		"|> TABLESAMPLE a(1 ROWS) zzz",
		"|> AS",
		"|> DESCRIBE zzz",
		"|> ASSERT",
		"|> DROP",
		"|> RENAME a.",
		"|> SET a.",
		"|> DROP a.",
		"|> ELSE",
		"|> ELSEIF",
		"|> TEE() zzz",
		"|> FORK() zzz",
		"|> WITH c AS (SELECT 1) zzz",
		"|> CREATE TABLE t (x INT64) AS (SELECT 1)",
		"|> EXPORT DATA AS (SELECT 1)",
		"|> RECURSIVE UNION ALL() a",
		"SELECT 1 |>",
		"SELECT 1 |> WHERE",
		"DESCRIBE t |>",
		"|> LOG(SELECT 1)",
		"|> LOG(WHERE true)",
		"|> IF true THEN () ELSEIF@1 false THEN ()",
		"|> IF true THEN () ELSE IF false THEN ()",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) {
			if _, errs := parseForTest(sql); len(errs) == 0 {
				t.Fatalf("parseForTest(%q) accepted, want a syntax error", sql)
			}
		})
	}
}

// TestPipe_AcceptsValidNearMisses holds inputs that look broken but the
// grammar accepts: an inserted token that is a valid continuation, an empty
// parenthesized subpipeline, and a trailing comma in an operator's list.
func TestPipe_AcceptsValidNearMisses(t *testing.T) {
	accept := []string{
		"|> AGGREGATE zzz",
		"|> AGGREGATE NULL ASC",
		"|> AGGREGATE NULL, GROUP BY x",
		"|> TEE()",
		"|> FORK()",
		"|> LOG()",
		"|> RECURSIVE UNION ALL()",
		"|> RECURSIVE UNION ALL() AS a",
		"|> UNION ALL TABLE a",
		"|> UNION ALL (SELECT 1), TABLE a",
		"|> DROP a, b,",
		"|> SET a = 1,",
		"|> ORDER BY x,",
		"|> WITH c AS (SELECT 1),",
		"|> WINDOW x AS y,",
		"|> EXTEND a WINDOW w AS ()",
		"|> AGGREGATE x AS y DESC NULLS LAST, GROUP AND ORDER BY z NULLS FIRST",
		"|> IF true THEN () ELSEIF false THEN (|> DISTINCT) ELSE ()",
		"SELECT x, |> WHERE true",
		"|> SELECT x, |> WHERE true",
		"|> MATCH_RECOGNIZE(ORDER BY x MEASURES 1 AS m PATTERN(a || b) DEFINE a AS true, b AS true)",
		"|> MATCH_RECOGNIZE(ORDER BY x MEASURES 1 AS m PATTERN((a)) DEFINE a AS true)",
		"|> JOIN a USING (x, y)",
		"|> FULL OUTER HASH JOIN a ON true",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) {
			if _, errs := parseForTest(sql); len(errs) != 0 {
				t.Fatalf("parseForTest(%q) rejected: %v", sql, errs)
			}
		})
	}
}

// TestPipe_Locs checks byte offsets: an operator's Loc starts at its `|>` and
// ends at its last token, and the statement covers the whole text.
func TestPipe_Locs(t *testing.T) {
	sql := "|> WHERE x > 1 |> DISTINCT"
	pipe := subpipelineOf(t, firstStmt(t, sql))
	if len(pipe.Ops) != 2 {
		t.Fatalf("got %d operators, want 2", len(pipe.Ops))
	}
	first := pipe.Ops[0].(*ast.PipeWhere).Loc
	if first.Start != 0 || first.End != len("|> WHERE x > 1") {
		t.Errorf("first operator Loc = %+v, want [0,%d)", first, len("|> WHERE x > 1"))
	}
	second := pipe.Ops[1].(*ast.PipeDistinct).Loc
	wantSecond := strings.Index(sql, "|> DISTINCT")
	if second.Start != wantSecond || second.End != len(sql) {
		t.Errorf("second operator Loc = %+v, want [%d,%d)", second, wantSecond, len(sql))
	}
	if pipe.Loc.Start != 0 || pipe.Loc.End != len(sql) {
		t.Errorf("subpipeline Loc = %+v, want [0,%d)", pipe.Loc, len(sql))
	}
}
