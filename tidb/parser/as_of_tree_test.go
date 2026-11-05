package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// TestAsOfFields holds the clause and the parts around it to the text they came from.
// Its Loc runs from AS to the last token of the expression.
func TestAsOfFields(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		alias  string
		parts  []string
		hints  int
		sample bool
		text   string // what the clause's Loc covers
	}{
		{sql: "select * from t as of timestamp now()", text: "as of timestamp now()"},
		{sql: "select * from t AS  OF\nTIMESTAMP '2020-01-01' where b", text: "AS  OF\nTIMESTAMP '2020-01-01'"},
		{sql: "select * from t a as of timestamp 1", alias: "a", text: "as of timestamp 1"},
		{sql: "select * from t as a as of timestamp 1", alias: "a", text: "as of timestamp 1"},
		{sql: "select * from t as `of` as of timestamp 1", alias: "of", text: "as of timestamp 1"},
		{sql: "select * from t timestamp as of timestamp 1", alias: "timestamp", text: "as of timestamp 1"},
		{sql: "select * from t partition (p0) as of timestamp 1", parts: []string{"p0"}, text: "as of timestamp 1"},
		{sql: "select * from t as of timestamp 1 + 2 limit 1", text: "as of timestamp 1 + 2"},
		{sql: "select * from t as of timestamp (select 1) where b", text: "as of timestamp (select 1)"},
		{sql: "select * from t as of timestamp a use index (i)", hints: 1, text: "as of timestamp a"},
		{sql: "select * from t as of timestamp a tablesample (50)", sample: true, text: "as of timestamp a"},
		{sql: "select * from t as of timestamp a use index (i) tablesample (50)", hints: 1, sample: true, text: "as of timestamp a"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			ref := parseSampledRef(t, tc.sql)
			if ref.AsOf == nil {
				t.Fatal("no AsOfClause on the reference")
			}
			if ref.Alias != tc.alias {
				t.Errorf("Alias = %q, want %q", ref.Alias, tc.alias)
			}
			if strings.Join(ref.Partitions, ",") != strings.Join(tc.parts, ",") {
				t.Errorf("Partitions = %v, want %v", ref.Partitions, tc.parts)
			}
			if len(ref.IndexHints) != tc.hints {
				t.Errorf("IndexHints = %d, want %d", len(ref.IndexHints), tc.hints)
			}
			if (ref.Sample != nil) != tc.sample {
				t.Errorf("Sample = %v, want present=%v", ref.Sample, tc.sample)
			}
			if got := tc.sql[ref.AsOf.Loc.Start:ref.AsOf.Loc.End]; got != tc.text {
				t.Errorf("Loc %+v covers %q, want %q", ref.AsOf.Loc, got, tc.text)
			}
			// The reference runs to the end of whatever came last in it.
			want := ref.AsOf.Loc.End
			if tc.hints > 0 {
				want = ref.IndexHints[tc.hints-1].Loc.End
			}
			if ref.Sample != nil {
				want = ref.Sample.Loc.End
			}
			if ref.Loc.End != want {
				t.Errorf("TableRef Loc.End = %d, want %d", ref.Loc.End, want)
			}
		})
	}
}

// TestAsOfOnlyWhereWritten checks a reference without the clause carries none, and
// that a backquoted of stays an alias.
func TestAsOfOnlyWhereWritten(t *testing.T) {
	for _, tc := range []struct{ sql, alias string }{
		{"select * from t", ""},
		{"select * from t a", "a"},
		{"select * from t as `of`", "of"},
		{"select * from t `of`", "of"},
	} {
		ref := parseSampledRef(t, tc.sql)
		if ref.AsOf != nil {
			t.Errorf("%q: AsOf = %+v, want nil", tc.sql, ref.AsOf)
		}
		if ref.Alias != tc.alias {
			t.Errorf("%q: Alias = %q, want %q", tc.sql, ref.Alias, tc.alias)
		}
	}
}

// TestAsOfEachSideOfAJoin checks each reference keeps its own clause, through a comma,
// a join and parentheses, on every statement that goes through a table factor.
func TestAsOfEachSideOfAJoin(t *testing.T) {
	for _, sql := range []string{
		"select * from t1 a as of timestamp 1 join t2 b as of timestamp 2 on a.id = b.id, (t3 as of timestamp 3)",
		"update t1 as of timestamp 1 join t2 as of timestamp 2 on 1, t3 as of timestamp 3 set a = 1",
		"delete t1 from t1 as of timestamp 1 join t2 as of timestamp 2 on 1, t3 as of timestamp 3",
		"delete from t1 using t1 as of timestamp 1 join t2 as of timestamp 2 on 1, t3 as of timestamp 3",
	} {
		t.Run(sql, func(t *testing.T) {
			list, err := Parse(sql)
			if err != nil {
				t.Fatal(err)
			}
			var got []string
			ast.Inspect(list.Items[0], func(n ast.Node) bool {
				if ref, ok := n.(*ast.TableRef); ok && ref.AsOf != nil {
					got = append(got, ref.Name+"="+exprText(t, ref.AsOf.Timestamp))
				}
				return true
			})
			if got, want := strings.Join(got, " "), "t1=1 t2=2 t3=3"; got != want {
				t.Errorf("clauses = %q, want %q", got, want)
			}
		})
	}
}

// TestAsOfWalk checks a consumer that walks the tree reaches the expression, and that
// the node dump shows it.
func TestAsOfWalk(t *testing.T) {
	list, err := Parse("select * from t as of timestamp (select 1)")
	if err != nil {
		t.Fatal(err)
	}
	var subqueries int
	var clause *ast.AsOfClause
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SubqueryExpr:
			subqueries++
		case *ast.AsOfClause:
			clause = v
		}
		return true
	})
	if subqueries != 1 {
		t.Errorf("walk reached %d subqueries, want the expression's one", subqueries)
	}
	if clause == nil {
		t.Fatal("walk did not reach the AsOfClause")
	}
	dump := ast.NodeToString(list.Items[0])
	for _, want := range []string{":asof {ASOF", ":timestamp {SUBQUERY"} {
		if !strings.Contains(dump, want) {
			t.Errorf("NodeToString has no %q: %s", want, dump)
		}
	}
}

// TestAsOfPlanReplayerKeepsItsExpression checks the other reader of the rule still
// stores only the expression.
func TestAsOfPlanReplayerKeepsItsExpression(t *testing.T) {
	list, err := Parse("plan replayer dump with stats as of timestamp 7 explain select 1")
	if err != nil {
		t.Fatal(err)
	}
	stmt, ok := list.Items[0].(*ast.PlanReplayerStmt)
	if !ok {
		t.Fatalf("%T, want *ast.PlanReplayerStmt", list.Items[0])
	}
	if got := exprText(t, stmt.HistoricalStats); got != "7" {
		t.Errorf("HistoricalStats = %q, want 7", got)
	}
}

// TestAsOfTruncated cuts the longest form at every byte. Each prefix either parses or
// returns a ParseError holding an offset inside the text, and none panics.
func TestAsOfTruncated(t *testing.T) {
	const sql = "select * from db.t partition (p0) as a as of timestamp now() - interval 1 second use index (i) tablesample (50) where b"
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

// TestAsOfErrorPosition checks a form the clause refuses reports the token it stopped
// at, the first one past what it could read.
func TestAsOfErrorPosition(t *testing.T) {
	for _, tc := range []struct{ sql, at string }{
		{"select * from t as of now()", "now()"},
		{"select * from t a of timestamp 1", "of timestamp 1"},
		{"select * from t as of timestamp", ""},
		{"select * from t as of timestamp now() a", "a"},
		{"select * from t as of timestamp now() partition (p0)", "partition (p0)"},
		{"select * from t use index (i) as of timestamp 1", "as of timestamp 1"},
		{"select * from t tablesample (50) as of timestamp 1", "as of timestamp 1"},
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

// TestAsOfTreeHeight checks the expression counts toward the height the parser
// measures: a left-deep chain builds a tree consumers recurse into, and the chain
// itself is read by a loop that never reaches enterNesting.
func TestAsOfTreeHeight(t *testing.T) {
	chain := strings.Repeat("1+", maxTreeHeight+1) + "1"
	_, err := Parse("select * from t as of timestamp " + chain)
	var pe *ParseError
	if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
		t.Errorf("err = %v, want the nesting error", briefErr(err))
	}
	nested := strings.Repeat("(", maxNestingDepth+1) + "1" + strings.Repeat(")", maxNestingDepth+1)
	if _, err := Parse("select * from t as of timestamp " + nested); err == nil {
		t.Error("an expression nested past maxNestingDepth parsed")
	}
}
