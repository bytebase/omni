package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestExprSelectFormsMatchPG holds expression and SELECT forms to the
// verdicts PostgreSQL 17.12's raw parser gives them, in compat.tsv's terms:
// accept, reject@<byte>, or reject where the server reports no position.
func TestExprSelectFormsMatchPG(t *testing.T) {
	for _, tc := range []struct {
		sql, want string
	}{
		// qual_Op as a prefix operator; OPERATOR alone is a column name.
		{`SELECT OPERATOR(pg_catalog.~) 1`, "accept"},
		{`SELECT operator ~ 1`, "accept"},
		{`SELECT 1 + OPERATOR(-) 2`, "accept"},
		// AexprConst: func_name '(' func_arg_list opt_sort_clause ')' Sconst.
		{`SELECT numeric(10, 2) '1.5'`, "accept"},
		{`SELECT a.b(1, 2) 's'`, "accept"},
		{`SELECT a(b => 1) 's'`, "reject@9"},
		{`SELECT a(1 ORDER BY 1) 's'`, "reject@11"},
		{`SELECT a(VARIADIC 1) 's'`, "reject@21"},
		// func_name takes no col_name_keyword.
		{`SELECT boolean(1) 's'`, "reject@14"},
		{`SELECT merge_action()`, "accept"},
		{`SELECT merge_action(1)`, "reject@20"},
		{`SELECT collation(a)`, "accept"},
		{`SELECT collation 1`, "reject@17"},
		{`SELECT is`, "reject@9"},
		{`SELECT NOT IN (1)`, "reject@11"},
		// Unreserved keywords that lead a clause only before what goes on with it.
		{`SELECT XMLROOT(a, VERSION no)`, "accept"},
		{`SELECT XMLROOT(a, VERSION NO VALUE)`, "accept"},
		{`SELECT FROM rows(1)`, "accept"},
		{`SELECT FROM rows`, "accept"},
		{`SELECT FROM ROWS FROM (f())`, "accept"},
		{`SELECT FETCH FIRST rows ROWS ONLY`, "accept"},
		{`SELECT FETCH FIRST ROW ONLY`, "accept"},
		// Window frames.
		{`SELECT WINDOW w AS (ROWS current PRECEDING)`, "accept"},
		{`SELECT WINDOW w AS (ROWS unbounded + 1 PRECEDING)`, "accept"},
		{`SELECT WINDOW w AS (ROWS 1)`, "reject@26"},
		{`SELECT WINDOW w AS (ROWS 1 FOLLOWING)`, "reject@25"},
		{`SELECT WINDOW w AS (ROWS BETWEEN CURRENT ROW AND 1 PRECEDING)`, "reject@49"},
		{`SELECT WINDOW w AS (ROWS BETWEEN 1 FOLLOWING AND CURRENT ROW)`, "reject@49"},
		{`SELECT WINDOW w AS (ROWS BETWEEN UNBOUNDED PRECEDING AND UNBOUNDED FOLLOWING EXCLUDE TIES)`, "accept"},
		{`SELECT WINDOW w AS (ROWS UNBOUNDED PRECEDING EXCLUDE)`, "reject@52"},
		{`SELECT WINDOW w AS )`, "reject@19"},
		// row OVERLAPS row.
		{`SELECT ROW(1, 2) OVERLAPS ROW(3, 4)`, "accept"},
		{`SELECT (1, 2) OVERLAPS (3, 4) AND true`, "accept"},
		{`SELECT ((1, 2)) OVERLAPS (3, 4)`, "reject@16"},
		{`SELECT (1, 2, 3) OVERLAPS (3, 4)`, "reject@7"},
		{`SELECT (1, 2) OVERLAPS (3, 4, 5)`, "reject@23"},
		{`SELECT f() OVERLAPS (3, 4)`, "reject@11"},
		// FROM items and their aliases.
		{`SELECT FROM t AS`, "reject@16"},
		{`SELECT FROM t a(b`, "reject@17"},
		{`SELECT FROM f() AS (a)`, "reject@21"},
		{`SELECT FROM f() AS (a int)`, "accept"},
		{`SELECT FROM t TABLESAMPLE system(1`, "reject@34"},
		{`SELECT FROM LATERAL`, "reject@19"},
		{`SELECT FROM left('a', 1)`, "accept"},
		{`SELECT FROM left`, "reject@16"},
		{`SELECT FROM xmltable`, "accept"},
		{`SELECT FROM char(1)`, "reject@16"},
		{`SELECT FROM ROWS FROM (f() AS a)`, "reject@30"},
		{`SELECT FROM a JOIN b ON`, "reject@23"},
		// Clauses that need what follows their keyword.
		{`SELECT FOR`, "reject@10"},
		{`SELECT FOR KEY`, "reject@14"},
		{`SELECT FOR NO KEY UPDATE OF t SKIP LOCKED`, "accept"},
		{`SELECT 1 UNION`, "reject@14"},
		{`SELECT DISTINCT`, "reject@15"},
		{`SELECT ALL`, "accept"},
		{`SELECT 1 AS`, "reject@11"},
		{`SELECT CASE END`, "reject@12"},
		{`SELECT CASE WHEN THEN 1 END`, "reject@17"},
		{`SELECT CASE 1 WHEN 1 THEN 2 ELSE END`, "reject@33"},
		{`SELECT +`, "reject@8"},
		{`SELECT count(*) FILTER (WHERE)`, "reject@29"},
		{`SELECT (SELECT 1).`, "reject@18"},
		{`SELECT a[]`, "reject@9"},
		{`SELECT a ~ ANY()`, "reject@15"},
		{`SELECT 1 FROM t HAVING`, "reject@22"},
		{`()`, "reject@1"},
		{`CREATE DOMAIN d int DEFAULT 1 IS`, "reject@32"},
		{`SELECT FETCH FIRST ROW WITH TIES`, "reject"},
		// WITH clauses.
		{`WITH a AS (SELECT) SEARCH FIRST BY a SET b SELECT`, "reject@26"},
		{`WITH a AS (SELECT) CYCLE a SET b TO DEFAULT 1 USING c SELECT`, "reject@36"},
		{`WITH a AS (SELECT) CYCLE a SET b TO 1 DEFAULT 0 USING c SELECT`, "accept"},
		{`WITH a AS NOT (SELECT) SELECT`, "reject@14"},
		{`WITH a AS NOT MATERIALIZED (SELECT) SELECT`, "accept"},
		// JSON and XML.
		{`SELECT a IS JSON WITH UNIQUE KEYS`, "accept"},
		{`CREATE TABLE t AS SELECT a IS JSON WITH DATA`, "reject@40"},
		{`SELECT JSON_ARRAYAGG(a FORMAT)`, "reject@23"},
		{`SELECT FROM JSON_TABLE(a, '$' COLUMNS (nested text))`, "accept"},
		{`SELECT FROM JSON_TABLE(a, '$' COLUMNS (NESTED '$' COLUMNS (b int)))`, "accept"},
		{`SELECT XMLSERIALIZE(DOCUMENT a AS text NO)`, "reject@41"},
		// The UNIQUE predicate is read, then refused as not implemented.
		{`SELECT UNIQUE (SELECT 1)`, "reject@7"},
		{`SELECT UNIQUE NULLS NOT DISTINCT (SELECT 1)`, "reject@7"},
		{`SELECT UNIQUE ()`, "reject@15"},
	} {
		if got := omniVerdict(tc.sql); got != tc.want {
			t.Errorf("%s: omni says %s, PostgreSQL %s", tc.sql, got, tc.want)
		}
	}
}

// TestStartsColNameKeywordColumnRef checks startsColNameKeywordColumnRef
// against the whole Keywords table: before any token but '(' every
// col_name_keyword where a c_expr starts is a column name, as PostgreSQL
// 17.12 reads SELECT k, SELECT k.a and OFFSET k for each of them.
func TestStartsColNameKeywordColumnRef(t *testing.T) {
	for _, kw := range Keywords {
		if kw.Category != ColNameKeyword {
			continue
		}
		for _, sql := range []string{"SELECT " + kw.Name, "SELECT " + kw.Name + ".a", "SELECT 1 OFFSET " + kw.Name} {
			if _, err := Parse(sql); err != nil {
				t.Errorf("Parse(%q): %v", sql, err)
			}
		}
	}
}

// TestTypmodConstShape: func_name '(' func_arg_list ')' Sconst is the
// string cast to func_name's type, with the arguments as its typmods.
func TestTypmodConstShape(t *testing.T) {
	list, err := Parse(`SELECT a.b(1, 2) 's'`)
	if err != nil {
		t.Fatal(err)
	}
	rt := list.Items[0].(*nodes.RawStmt).Stmt.(*nodes.SelectStmt).TargetList.Items[0].(*nodes.ResTarget)
	tc, ok := rt.Val.(*nodes.TypeCast)
	if !ok {
		t.Fatalf("got %T, want *nodes.TypeCast", rt.Val)
	}
	if n := len(tc.TypeName.Names.Items); n != 2 {
		t.Errorf("type name has %d parts, want 2", n)
	}
	if tc.TypeName.Typmods == nil || len(tc.TypeName.Typmods.Items) != 2 {
		t.Errorf("typmods = %v, want two", tc.TypeName.Typmods)
	}
	if s, ok := tc.Arg.(*nodes.A_Const).Val.(*nodes.String); !ok || s.Str != "s" {
		t.Errorf("arg = %v, want 's'", tc.Arg)
	}
}
