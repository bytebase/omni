package parser

import (
	"errors"
	"strconv"
	"strings"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// parseSampledRef parses one SELECT and returns its first FROM item as a table
// reference. Parse reads statements with no separator, so a statement count
// other than one means the clause was not read as part of the reference.
func parseSampledRef(t *testing.T, sql string) *ast.TableRef {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	sel, ok := list.Items[0].(*ast.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.SelectStmt", sql, list.Items[0])
	}
	ref, ok := sel.From[0].(*ast.TableRef)
	if !ok {
		t.Fatalf("Parse(%q): FROM item is %T, want *ast.TableRef", sql, sel.From[0])
	}
	return ref
}

// exprText reads a size or seed back as text, which is enough to tell one from
// the other. A nil expression, which the clause allows for both, reads as "".
func exprText(t *testing.T, e ast.ExprNode) string {
	t.Helper()
	switch v := e.(type) {
	case nil:
		return ""
	case *ast.IntLit:
		return strconv.FormatInt(v.Value, 10)
	case *ast.StringLit:
		return "'" + v.Value + "'"
	case *ast.FloatLit:
		return v.Value
	case *ast.ColumnRef:
		return v.Column
	default:
		t.Fatalf("unexpected expression %T", e)
		return ""
	}
}

// TestTableSampleFields holds every field of the clause to the text it came
// from, including the two the grammar lets a statement leave out. Its Loc runs
// from TABLESAMPLE to the last closing parenthesis.
func TestTableSampleFields(t *testing.T) {
	for _, tc := range []struct {
		sql    string
		alias  string
		parts  []string
		hints  int
		method ast.TableSampleMethod
		size   string
		unit   ast.TableSampleUnit
		seed   string
		text   string // what the clause's Loc covers
	}{
		{sql: "select * from t tablesample (50)", method: ast.SampleMethodDefault, size: "50", text: "tablesample (50)"},
		{sql: "select * from t tablesample system (50)", method: ast.SampleMethodSystem, size: "50", text: "tablesample system (50)"},
		{sql: "select * from t TABLESAMPLE BERNOULLI(50)", method: ast.SampleMethodBernoulli, size: "50", text: "TABLESAMPLE BERNOULLI(50)"},
		{sql: "select * from t tablesample regions ()", method: ast.SampleMethodRegions, text: "tablesample regions ()"},
		{sql: "select * from t tablesample ()", method: ast.SampleMethodDefault, text: "tablesample ()"},
		{sql: "select * from t tablesample (50 rows)", size: "50", unit: ast.SampleUnitRows, text: "tablesample (50 rows)"},
		{sql: "select * from t tablesample (49.9 PERCENT)", size: "49.9", unit: ast.SampleUnitPercent, text: "tablesample (49.9 PERCENT)"},
		{sql: "select * from t tablesample ('34')", size: "'34'", text: "tablesample ('34')"},
		// A column named percent is the size; only the word after it is the unit.
		{sql: "select * from t tablesample (percent)", size: "percent", text: "tablesample (percent)"},
		{sql: "select * from t tablesample (percent percent)", size: "percent", unit: ast.SampleUnitPercent, text: "tablesample (percent percent)"},
		{sql: "select * from t tablesample (50) repeatable (123456789)", size: "50", seed: "123456789", text: "tablesample (50) repeatable (123456789)"},
		{sql: "select * from t tablesample () repeatable (7)", seed: "7", text: "tablesample () repeatable (7)"},
		{sql: "select * from t as a tablesample (50) where a.b", alias: "a", size: "50", text: "tablesample (50)"},
		{sql: "select * from t `tablesample` tablesample (50)", alias: "tablesample", size: "50", text: "tablesample (50)"},
		{sql: "select * from t partition (p0) tablesample (50)", parts: []string{"p0"}, size: "50", text: "tablesample (50)"},
		{sql: "select * from t use index (i) tablesample (50) limit 1", hints: 1, size: "50", text: "tablesample (50)"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			ref := parseSampledRef(t, tc.sql)
			s := ref.Sample
			if s == nil {
				t.Fatal("no TableSample on the reference")
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
			if s.Method != tc.method {
				t.Errorf("Method = %d, want %d", s.Method, tc.method)
			}
			if got := exprText(t, s.Size); got != tc.size {
				t.Errorf("Size = %q, want %q", got, tc.size)
			}
			if s.Unit != tc.unit {
				t.Errorf("Unit = %d, want %d", s.Unit, tc.unit)
			}
			if got := exprText(t, s.Repeatable); got != tc.seed {
				t.Errorf("Repeatable = %q, want %q", got, tc.seed)
			}
			if got := tc.sql[s.Loc.Start:s.Loc.End]; got != tc.text {
				t.Errorf("Loc %+v covers %q, want %q", s.Loc, got, tc.text)
			}
			// The reference ends where the clause does.
			if ref.Loc.End != s.Loc.End {
				t.Errorf("TableRef Loc.End = %d, want the clause's %d", ref.Loc.End, s.Loc.End)
			}
		})
	}
}

// TestTableSampleEachSideOfAJoin checks each reference keeps its own clause,
// through a comma, a join and parentheses.
func TestTableSampleEachSideOfAJoin(t *testing.T) {
	list, err := Parse("select * from t1 a tablesample (10) join t2 b tablesample (20 rows) on a.id = b.id, (t3 tablesample (30))")
	if err != nil {
		t.Fatal(err)
	}
	var sizes []string
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		if ref, ok := n.(*ast.TableRef); ok {
			if ref.Sample == nil {
				t.Errorf("%s has no sample", ref.Name)
				return true
			}
			sizes = append(sizes, ref.Name+"="+exprText(t, ref.Sample.Size))
		}
		return true
	})
	if got, want := strings.Join(sizes, " "), "t1=10 t2=20 t3=30"; got != want {
		t.Errorf("samples = %q, want %q", got, want)
	}
}

// TestTableSampleWalk checks a consumer that walks the tree reaches the size and
// the seed, and that the node dump shows them.
func TestTableSampleWalk(t *testing.T) {
	list, err := Parse("select * from t tablesample system ((select 1) rows) repeatable ((select 2))")
	if err != nil {
		t.Fatal(err)
	}
	var subqueries int
	var sample *ast.TableSample
	ast.Inspect(list.Items[0], func(n ast.Node) bool {
		switch v := n.(type) {
		case *ast.SubqueryExpr:
			subqueries++
		case *ast.TableSample:
			sample = v
		}
		return true
	})
	if subqueries != 2 {
		t.Errorf("walk reached %d subqueries, want the size and the seed", subqueries)
	}
	if sample == nil {
		t.Fatal("walk did not reach the TableSample")
	}
	dump := ast.NodeToString(list.Items[0])
	for _, want := range []string{
		":sample {TABLESAMPLE",
		":method " + strconv.Itoa(int(ast.SampleMethodSystem)),
		":unit " + strconv.Itoa(int(ast.SampleUnitRows)),
		":size ",
		":repeatable ",
	} {
		if !strings.Contains(dump, want) {
			t.Errorf("NodeToString has no %q: %s", want, dump)
		}
	}
}

// TestTableSampleOnlyOnATableFactor pins where the clause is not read. These all
// parse some other way if a reader is too broad: DELETE and LOCK TABLES name a
// table without going through TableFactor, and a derived table or DUAL is not a
// table name.
func TestTableSampleOnlyOnATableFactor(t *testing.T) {
	for _, sql := range []string{
		"delete from t tablesample (50)",
		"delete from t a tablesample (50) where 1 = 1",
		"lock tables t tablesample (50) read",
		"table t tablesample (50)",
		"insert into t tablesample (50) values (1)",
		"analyze table t tablesample (50)",
		"select * from (select 1) a tablesample (50)",
		"select * from (t1 join t2) tablesample (50)",
		"select * from dual tablesample (50)",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Errorf("Parse(%q) accepted a clause TiDB's grammar has no place for", sql)
			}
		})
	}
}

// TestTableSampleTruncated cuts the longest form at every byte. Each prefix either
// parses or returns a ParseError holding an offset inside the text, and none panics.
func TestTableSampleTruncated(t *testing.T) {
	const sql = "select * from db.t partition (p0) as a use index (i) tablesample system (1 + 2 rows) repeatable (3) where b"
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

// TestTableSampleTreeHeight checks both expressions count toward the height the
// parser measures: a left-deep chain in either one builds a tree consumers recurse
// into, and the chain itself is read by a loop that never reaches enterNesting.
func TestTableSampleTreeHeight(t *testing.T) {
	chain := strings.Repeat("1+", maxTreeHeight+1) + "1"
	for _, sql := range []string{
		"select * from t tablesample (" + chain + ")",
		"select * from t tablesample (1) repeatable (" + chain + ")",
	} {
		_, err := Parse(sql)
		var pe *ParseError
		if !errors.As(err, &pe) || pe.Message != "statement is nested too deeply" {
			t.Errorf("%s...: err = %v, want the nesting error", sql[:40], briefErr(err))
		}
	}
	nested := strings.Repeat("(", maxNestingDepth+1) + "1" + strings.Repeat(")", maxNestingDepth+1)
	if _, err := Parse("select * from t tablesample (" + nested + ")"); err == nil {
		t.Error("a size nested past maxNestingDepth parsed")
	}
}

// TestLexerWordTouchingDot pins upstream's rule for a reserved word beside a dot,
// which the tables in tablesample_forms_test.go measure against the server. A
// reserved word is an identifier when the next byte is a dot. TABLESAMPLE is one too
// when the last byte before it that is not a space is a dot, where a tab, a newline
// or a comment ends the search. Other reserved words keep their token there because
// completion needs them to, and an unreserved keyword keeps its token throughout.
func TestLexerWordTouchingDot(t *testing.T) {
	for _, tc := range []struct {
		sql  string
		word string
		want int
	}{
		{"tablesample.a", "tablesample", tokIDENT},
		{"tablesample .a", "tablesample", kwTABLESAMPLE},
		{"tablesample\t.a", "tablesample", kwTABLESAMPLE},
		{"tablesample", "tablesample", kwTABLESAMPLE},
		{"t.tablesample", "tablesample", tokIDENT},
		{"t. tablesample", "tablesample", tokIDENT},
		{"t .  tablesample", "tablesample", tokIDENT},
		{"t.\ttablesample", "tablesample", kwTABLESAMPLE},
		{"t.\ntablesample", "tablesample", kwTABLESAMPLE},
		{"t./**/tablesample", "tablesample", kwTABLESAMPLE},
		{"a tablesample", "tablesample", kwTABLESAMPLE},
		{"select.a", "select", tokIDENT},
		{"o. from", "from", kwFROM},
		{"global.a", "global", kwGLOBAL},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			lex := NewLexer(tc.sql)
			for {
				tok := lex.NextToken()
				if tok.Type == tokEOF {
					t.Fatalf("no token for %q", tc.word)
				}
				if tok.Str == tc.word {
					if tok.Type != tc.want {
						t.Errorf("%q lexed as %d, want %d", tc.word, tok.Type, tc.want)
					}
					return
				}
			}
		})
	}
}
