package catalog

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/internal/sharedserver"
	nodes "github.com/bytebase/omni/mysql/ast"
	"github.com/bytebase/omni/mysql/parser"
	gomysql "github.com/go-sql-driver/mysql"
)

// quantifiedSchema is the database every quantified-comparison test reads: qt
// holds a row on each side of 1, qn holds 1 and a NULL, and qe is empty, so
// ANY and ALL answer differently over each of them for every operator.
const quantifiedSchema = `CREATE TABLE qt (a INT);
CREATE TABLE qn (a INT);
CREATE TABLE qe (a INT);
INSERT INTO qt VALUES (0), (2);
INSERT INTO qn VALUES (1), (NULL)`

// quantifiedOps are the comparison operators MySQL 8.0 takes before ANY, SOME
// and ALL (sql_yacc.yy comp_op), with the operator the analyzer records. <=> is not
// one of them: the server answers 1064, and so does the parser.
var quantifiedOps = []struct{ src, op string }{
	{"=", "="},
	{"<>", "!="},
	{"!=", "!="},
	{"<", "<"},
	{"<=", "<="},
	{">", ">"},
	{">=", ">="},
}

var quantifiers = []struct {
	src  string
	kind SubLinkKind
}{
	{"ANY", SubLinkAny},
	{"SOME", SubLinkAny},
	{"ALL", SubLinkAll},
}

// quantifiedCatalog returns a catalog holding quantifiedSchema's tables in db.
func quantifiedCatalog(t *testing.T, db string) *Catalog {
	t.Helper()
	c := New()
	results, err := c.Exec(fmt.Sprintf("CREATE DATABASE %s; USE %s; CREATE TABLE qt (a INT); CREATE TABLE qn (a INT); CREATE TABLE qe (a INT);", db, db), nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, r := range results {
		if r.Error != nil {
			t.Fatal(r.Error)
		}
	}
	return c
}

func analyzeQuantified(t *testing.T, c *Catalog, sql string) *Query {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if list.Len() != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, list.Len())
	}
	sel, ok := list.Items[0].(*nodes.SelectStmt)
	if !ok {
		t.Fatalf("Parse(%q): %T, want *ast.SelectStmt", sql, list.Items[0])
	}
	q, err := c.AnalyzeSelectStmt(sel)
	if err != nil {
		t.Fatalf("AnalyzeSelectStmt(%q): %v", sql, err)
	}
	return q
}

// checkQuantifiedLink fails unless e is the quantified comparison `1 op
// quantifier (subquery)`.
func checkQuantifiedLink(t *testing.T, e AnalyzedExpr, kind SubLinkKind, op string) {
	t.Helper()
	link, ok := e.(*SubLinkExprQ)
	if !ok {
		t.Fatalf("expression is %T, want *SubLinkExprQ", e)
	}
	if link.Kind != kind {
		t.Errorf("Kind = %d, want %d", link.Kind, kind)
	}
	if link.Op != op {
		t.Errorf("Op = %q, want %q", link.Op, op)
	}
	if c, ok := link.TestExpr.(*ConstExprQ); !ok || c.Value != "1" {
		t.Errorf("TestExpr = %#v, want the constant 1", link.TestExpr)
	}
	if link.Subquery == nil {
		t.Errorf("Subquery is nil")
	}
}

// TestQuantifiedComparisonAnalyze checks `x op {ANY|SOME|ALL} (subquery)` is
// analyzed as the quantified link it is, not as a comparison with a scalar
// subquery, and that DeparseQuery writes it back so: without the quantifier the
// text fails on the server as soon as the subquery returns a second row.
func TestQuantifiedComparisonAnalyze(t *testing.T) {
	c := quantifiedCatalog(t, "d")
	for _, op := range quantifiedOps {
		for _, qu := range quantifiers {
			sql := fmt.Sprintf("SELECT 1 %s %s (SELECT a FROM qt) AS r", op.src, qu.src)
			t.Run(sql, func(t *testing.T) {
				q1 := analyzeQuantified(t, c, sql)
				checkQuantifiedLink(t, q1.TargetList[0].Expr, qu.kind, op.op)

				sql2 := DeparseQuery(q1)
				word := " any ("
				if qu.kind == SubLinkAll {
					word = " all ("
				}
				if want := "(1 " + op.op + word; !strings.Contains(sql2, want) {
					t.Errorf("DeparseQuery = %s, want it to hold %q", sql2, want)
				}
				q2 := analyzeQuantified(t, c, sql2)
				checkQuantifiedLink(t, q2.TargetList[0].Expr, qu.kind, op.op)
				if sql3 := DeparseQuery(q2); sql3 != sql2 {
					t.Errorf("deparse is not stable:\n first: %s\nsecond: %s", sql2, sql3)
				}
			})
		}
	}

	t.Run("scalar subquery stays scalar", func(t *testing.T) {
		q := analyzeQuantified(t, c, "SELECT 1 > (SELECT 0) AS r")
		op, ok := q.TargetList[0].Expr.(*OpExprQ)
		if !ok {
			t.Fatalf("expression is %T, want *OpExprQ", q.TargetList[0].Expr)
		}
		if link, ok := op.Right.(*SubLinkExprQ); !ok || link.Kind != SubLinkScalar {
			t.Errorf("right operand = %#v, want a scalar SubLinkExprQ", op.Right)
		}
	})

	t.Run("view", func(t *testing.T) {
		vc := quantifiedCatalog(t, "d")
		results, err := vc.Exec("CREATE VIEW v AS SELECT a FROM qt WHERE a >= ALL (SELECT a FROM qt)", nil)
		if err != nil {
			t.Fatal(err)
		}
		for _, r := range results {
			if r.Error != nil {
				t.Fatal(r.Error)
			}
		}
		v := vc.GetDatabase("d").Views["v"]
		if v == nil || v.AnalyzedQuery == nil {
			t.Fatalf("view v has no analyzed query: %#v", v)
		}
		link, ok := v.AnalyzedQuery.JoinTree.Quals.(*SubLinkExprQ)
		if !ok || link.Kind != SubLinkAll || link.Op != ">=" {
			t.Fatalf("WHERE = %#v, want a >= ALL link", v.AnalyzedQuery.JoinTree.Quals)
		}
		if got := DeparseQuery(v.AnalyzedQuery); !strings.Contains(got, ">= all (select") {
			t.Errorf("DeparseQuery = %s, lost the quantifier", got)
		}
	})
}

// quantifiedOracleQueries are the queries TestQuantifiedComparisonMySQL runs as
// written and as deparsed: every operator and quantifier over a set on both
// sides of 1, a set holding a NULL, and an empty set, plus the shapes where
// the comparison sits inside another expression.
func quantifiedOracleQueries() []string {
	var out []string
	for _, op := range quantifiedOps {
		for _, qu := range quantifiers {
			out = append(out, fmt.Sprintf(
				"SELECT 1 %[1]s %[2]s (SELECT a FROM qt) AS r1, 1 %[1]s %[2]s (SELECT a FROM qn) AS r2, 1 %[1]s %[2]s (SELECT a FROM qe) AS r3, NULL %[1]s %[2]s (SELECT a FROM qt) AS r4",
				op.src, qu.src))
		}
	}
	return append(out,
		"SELECT 1 > ALL (SELECT 0 UNION ALL SELECT 2) AS r1, 1 > ANY (SELECT 0 UNION ALL SELECT 2) AS r2",
		"SELECT q.a FROM qt AS q WHERE q.a >= ALL (SELECT a FROM qt) ORDER BY q.a",
		"SELECT q.a FROM qt AS q WHERE q.a < SOME (SELECT a FROM qt) ORDER BY q.a",
		"SELECT 1 = (2 > ANY (SELECT a FROM qt)) AS r",
		"SELECT NOT (1 = ANY (SELECT a FROM qt)) AS r",
		"SELECT (1 + 1) <= ALL (SELECT a + 2 FROM qt) AS r",
	)
}

// TestQuantifiedComparisonMySQL runs each query and its DeparseQuery text on
// MySQL 8.0 and requires the same answer: the same rows, or the same error.
func TestQuantifiedComparisonMySQL(t *testing.T) {
	ctx := context.Background()
	d, err := sharedserver.MySQL80.Database(ctx, "mysql_catalog_quantified")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close(ctx) })
	db, name := d.DB, d.Name
	for _, stmt := range strings.Split(quantifiedSchema, ";\n") {
		if _, err := db.Exec(stmt); err != nil {
			t.Fatalf("%s: %v", stmt, err)
		}
	}

	t.Run("<=> takes no quantifier", func(t *testing.T) {
		const sql = "SELECT 1 <=> ANY (SELECT a FROM qt)"
		_, err := db.Exec(sql)
		var myErr *gomysql.MySQLError
		if !errors.As(err, &myErr) || myErr.Number != 1064 {
			t.Errorf("server answered %v, want 1064", err)
		}
		if _, err := parser.Parse(sql); err == nil {
			t.Errorf("the parser accepts %q", sql)
		}
	})

	c := quantifiedCatalog(t, name)
	for _, sql := range quantifiedOracleQueries() {
		t.Run(sql, func(t *testing.T) {
			deparsed := DeparseQuery(analyzeQuantified(t, c, sql))
			want, wantErr := queryRowsText(db, sql)
			got, gotErr := queryRowsText(db, deparsed)
			if wantErr != nil {
				t.Fatalf("original query failed: %v", wantErr)
			}
			if gotErr != nil {
				t.Fatalf("deparsed: %s\nfailed: %v", deparsed, gotErr)
			}
			if got != want {
				t.Errorf("deparsed: %s\noriginal rows: %s\ndeparsed rows: %s", deparsed, want, got)
			}
		})
	}
}

// queryRowsText runs query and returns its rows as text, one row per line and
// NULL for a SQL NULL.
func queryRowsText(db *sql.DB, query string) (string, error) {
	rows, err := db.Query(query)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		return "", err
	}
	var b strings.Builder
	for rows.Next() {
		vals := make([]sql.NullString, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return "", err
		}
		for i, v := range vals {
			if i > 0 {
				b.WriteByte('\t')
			}
			if v.Valid {
				b.WriteString(v.String)
			} else {
				b.WriteString("NULL")
			}
		}
		b.WriteByte('\n')
	}
	return b.String(), rows.Err()
}
