package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/oracle/ast"
)

// Oracle's table_reference grammar places PIVOT/UNPIVOT between the
// query_table_expression and the trailing t_alias:
//
//	table_reference ::= query_table_expression
//	                    [ pivot_clause | unpivot_clause ] [ t_alias ]
//
// so the pivoted result can carry a bare-identifier alias and participate in
// joins like any other from-item. t_alias takes NO "AS" keyword and no column
// alias list. All shapes below verified against Oracle 23ai Free.

func fromItem(t *testing.T, sql string) ast.Node {
	t.Helper()
	stmt, ok := rawStmt(t, sql).(*ast.SelectStmt)
	if !ok {
		t.Fatalf("expected SelectStmt")
	}
	if stmt.FromClause == nil || stmt.FromClause.Len() != 1 {
		t.Fatalf("expected 1 from-item, got %v", stmt.FromClause)
	}
	return stmt.FromClause.Items[0]
}

func TestPivotTableAlias(t *testing.T) {
	item := fromItem(t, `SELECT p.a FROM t1 PIVOT (SUM(b) FOR m IN (1, 2, 3)) p`)
	pc, ok := item.(*ast.PivotClause)
	if !ok {
		t.Fatalf("expected PivotClause from-item, got %T", item)
	}
	if pc.Alias == nil || pc.Alias.Name != "P" {
		t.Fatalf("pivot alias = %v, want P", pc.Alias)
	}
	if _, ok := pc.Source.(*ast.TableRef); !ok {
		t.Fatalf("pivot source = %T, want *TableRef", pc.Source)
	}
}

func TestPivotSubqueryAliasReferencedInOrderBy(t *testing.T) {
	// The originally reported shape: (subquery) PIVOT (...) alias, with the
	// alias qualifying select-list and ORDER BY references.
	sql := `SELECT p.a, p.c1 FROM (SELECT a, b, m FROM t1 WHERE b > 0) PIVOT (SUM(b) FOR m IN (1 AS c1, 2 AS c2)) p ORDER BY p.a`
	item := fromItem(t, sql)
	pc, ok := item.(*ast.PivotClause)
	if !ok {
		t.Fatalf("expected PivotClause from-item, got %T", item)
	}
	if pc.Alias == nil || pc.Alias.Name != "P" {
		t.Fatalf("pivot alias = %v, want P", pc.Alias)
	}
	if _, ok := pc.Source.(*ast.SubqueryRef); !ok {
		t.Fatalf("pivot source = %T, want *SubqueryRef", pc.Source)
	}
}

func TestUnpivotAlias(t *testing.T) {
	item := fromItem(t, `SELECT u.q, u.v FROM t1 UNPIVOT (v FOR q IN (c1, c2, c3)) u`)
	uc, ok := item.(*ast.UnpivotClause)
	if !ok {
		t.Fatalf("expected UnpivotClause from-item, got %T", item)
	}
	if uc.Alias == nil || uc.Alias.Name != "U" {
		t.Fatalf("unpivot alias = %v, want U", uc.Alias)
	}
	if _, ok := uc.Source.(*ast.TableRef); !ok {
		t.Fatalf("unpivot source = %T, want *TableRef", uc.Source)
	}
}

func TestPivotXMLAlias(t *testing.T) {
	item := fromItem(t, `SELECT x.a FROM t1 PIVOT XML (SUM(b) FOR m IN (ANY)) x`)
	pc, ok := item.(*ast.PivotClause)
	if !ok {
		t.Fatalf("expected PivotClause from-item, got %T", item)
	}
	if !pc.XML {
		t.Fatal("expected XML pivot")
	}
	if pc.Alias == nil || pc.Alias.Name != "X" {
		t.Fatalf("pivot xml alias = %v, want X", pc.Alias)
	}
}

func TestPivotJoinParticipation(t *testing.T) {
	item := fromItem(t, `SELECT p.a, d.name FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) p JOIN d ON p.a = d.a`)
	jc, ok := item.(*ast.JoinClause)
	if !ok {
		t.Fatalf("expected JoinClause from-item, got %T", item)
	}
	pc, ok := jc.Left.(*ast.PivotClause)
	if !ok {
		t.Fatalf("join left = %T, want *PivotClause", jc.Left)
	}
	if pc.Alias == nil || pc.Alias.Name != "P" {
		t.Fatalf("pivot alias = %v, want P", pc.Alias)
	}
}

func TestPivotAliasRejectsASAndColumnList(t *testing.T) {
	// t_alias is a bare identifier: no AS keyword, no column alias list.
	// Oracle raises ORA-00933 on both.
	for _, sql := range []string{
		`SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) AS p`,
		`SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) p(x, y)`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("expected parse error for %q", sql)
		}
	}
}

func TestPivotAliasInsidePLSQLCursor(t *testing.T) {
	// The reported failure surfaced inside a function body: OPEN ... FOR a
	// pivoted subquery with a trailing alias.
	sql := `CREATE OR REPLACE FUNCTION f_report(p_from DATE, p_to DATE) RETURN SYS_REFCURSOR IS
  v_cur SYS_REFCURSOR;
BEGIN
  OPEN v_cur FOR
    SELECT p.a, p.c1, p.c2
    FROM (SELECT a, b, m FROM t1 WHERE d BETWEEN p_from AND p_to) PIVOT (SUM(b) FOR m IN (1 AS c1, 2 AS c2)) p
    ORDER BY p.a;
  RETURN v_cur;
END;`
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("expected 1 statement, got %d", result.Len())
	}
}

func TestPivotAliasOutputShape(t *testing.T) {
	// The pivoted from-item serializes with its alias AND its source: the
	// base relation now lives only in PivotClause.Source, so dropping it
	// from outfuncs would lose the table entirely.
	stmt := rawStmt(t, `SELECT p.a FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) p`)
	out := ast.NodeToString(stmt)
	for _, want := range []string{":alias", ":source", "T1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in output, got %s", want, out)
		}
	}

	stmt = rawStmt(t, `SELECT u.q FROM t1 UNPIVOT (v FOR q IN (c1, c2)) u`)
	out = ast.NodeToString(stmt)
	for _, want := range []string{":alias", ":source", "T1"} {
		if !strings.Contains(out, want) {
			t.Fatalf("expected %s in unpivot output, got %s", want, out)
		}
	}
}

// Oracle accepts pivot_clause / unpivot_clause after every
// query_table_expression form below (verified on 23ai Free), so each source
// parser must feed the pivot suffix rather than returning early.
func TestPivotAfterEveryTableSource(t *testing.T) {
	sources := []struct {
		name string
		from string
	}{
		{"table", `t1`},
		{"subquery", `(SELECT a, b, m FROM t1)`},
		{"parenthesized table", `(t1)`},
		{"table collection", `TABLE(f())`},
		{"xmltable", `XMLTABLE('/r' PASSING x COLUMNS m NUMBER PATH '@m', b NUMBER PATH '@b')`},
		{"json_table", `JSON_TABLE(j, '$[*]' COLUMNS (m NUMBER PATH '$.m', b NUMBER PATH '$.b'))`},
		{"containers", `CONTAINERS(t1)`},
		{"shards", `SHARDS(t1)`},
		{"sample", `t1 SAMPLE (10)`},
		{"flashback", `t1 AS OF SCN 1`},
		{"match_recognize", `t1 MATCH_RECOGNIZE (MEASURES 1 AS m ALL ROWS PER MATCH PATTERN (a) DEFINE a AS 1 = 1)`},
	}
	suffixes := []struct {
		name   string
		suffix string
	}{
		{"pivot", `PIVOT (SUM(b) FOR m IN (1, 2)) p`},
		{"unpivot", `UNPIVOT (v FOR q IN (m, b)) u`},
	}
	for _, src := range sources {
		for _, sfx := range suffixes {
			t.Run(src.name+"/"+sfx.name, func(t *testing.T) {
				sql := `SELECT * FROM ` + src.from + ` ` + sfx.suffix
				if _, err := Parse(sql); err != nil {
					t.Fatalf("parse failed: %v; sql: %s", err, sql)
				}
			})
		}
	}
}

func TestPivotChaining(t *testing.T) {
	// Oracle chains pivot/unpivot clauses, optionally aliased per step.
	for _, sql := range []string{
		`SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) PIVOT (COUNT(*) FOR a IN (1))`,
		`SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) p PIVOT (COUNT(*) FOR a IN (1)) q`,
		`SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1, 2)) UNPIVOT (v FOR q IN (a))`,
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("parse failed: %v; sql: %s", err, sql)
		}
	}
	// The chain nests: outer clause's Source is the inner clause.
	item := fromItem(t, `SELECT * FROM t1 PIVOT (SUM(b) FOR m IN (1)) p PIVOT (COUNT(*) FOR a IN (1)) q`)
	outer, ok := item.(*ast.PivotClause)
	if !ok {
		t.Fatalf("expected outer PivotClause, got %T", item)
	}
	if outer.Alias == nil || outer.Alias.Name != "Q" {
		t.Fatalf("outer alias = %v, want Q", outer.Alias)
	}
	inner, ok := outer.Source.(*ast.PivotClause)
	if !ok {
		t.Fatalf("outer source = %T, want inner *PivotClause", outer.Source)
	}
	if inner.Alias == nil || inner.Alias.Name != "P" {
		t.Fatalf("inner alias = %v, want P", inner.Alias)
	}
}

func TestIncompletePivotRejected(t *testing.T) {
	// PIVOT/UNPIVOT not followed by the clause body must not silently become
	// an empty clause plus an alias (Oracle: ORA-03049 on `PIVOT p`).
	// LATERAL views cannot be pivoted at all (Oracle: ORA-56905).
	for _, sql := range []string{
		`SELECT * FROM t1 PIVOT p`,
		`SELECT * FROM t1 PIVOT`,
		`SELECT * FROM t1 UNPIVOT u`,
		`SELECT * FROM t1 PIVOT XML junk`,
		`SELECT * FROM LATERAL (SELECT a, m FROM t1) PIVOT (SUM(a) FOR m IN (1))`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("expected parse error for %q", sql)
		}
	}
}
