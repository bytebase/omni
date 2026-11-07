package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestListElementsMatchPG pins that a comma-separated list takes no empty
// element. gram.y writes every list as `x | list ',' x` (target_list,
// expr_list, func_arg_list, sortby_list, from_list, columnList,
// index_params, privilege_list, var_list, ...), so a comma promises an
// element, and only a list the grammar makes optional as a whole
// (opt_target_list) may be empty. The error is at the token after the
// comma.
func TestListElementsMatchPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		// target_list
		{"SELECT a, b, FROM t", false},
		{"SELECT a, FROM t", false},
		{"SELECT a, b,", false},
		{"SELECT a, b, WHERE true", false},
		{"INSERT INTO t VALUES (1) RETURNING a,", false},
		{"UPDATE t SET a = 1 RETURNING a,", false},
		{"DELETE FROM t RETURNING a,", false},
		{"SELECT a, b FROM t", true},
		{"SELECT FROM t", true},
		{"SELECT", true},

		// expr_list
		{"INSERT INTO t VALUES (1, 2,)", false},
		{"VALUES (1, 2,)", false},
		{"VALUES ()", false},
		{"INSERT INTO t VALUES ()", false},
		{"SELECT 1 WHERE 1 IN (1, 2,)", false},
		{"SELECT 1 WHERE (1, 2) IN ((1, 2), )", false},
		{"SELECT ARRAY[1, 2,]", false},
		{"SELECT ARRAY[[1], ]", false},
		{"SELECT ROW(1, 2,)", false},
		{"SELECT (1, 2,)", false},
		{"SELECT ()", false},
		{"SELECT COALESCE(1, )", false},
		{"SELECT GREATEST(1, )", false},
		{"SELECT xmlconcat(1, )", false},
		{"SELECT xmlelement(NAME x, )", false},
		{"SELECT DISTINCT ON (a, ) a FROM t", false},
		{"SELECT DISTINCT ON () a FROM t", false},
		{"SELECT 1 FROM t WINDOW w AS (PARTITION BY a, )", false},
		{"SELECT 1 FROM t GROUP BY ROLLUP (a, )", false},
		{"SELECT * FROM t TABLESAMPLE system (1, )", false},
		{"EXECUTE p (1, )", false},
		{"MERGE INTO t USING u ON true WHEN NOT MATCHED THEN INSERT VALUES (1, )", false},
		{"CREATE TABLE t (a numeric(10, ))", false},
		{"CREATE TABLE t (a numeric())", false},
		{"SELECT 1::numeric(10, )", false},
		{"SELECT 1::bit(1, )", false},
		{"INSERT INTO t VALUES (1, 2), (3, 4)", true},
		{"SELECT ARRAY[]", true},
		{"SELECT ROW()", true},
		{"SELECT 1 FROM t GROUP BY ()", true},
		{"SELECT 1 FROM t GROUP BY GROUPING SETS ((), (a))", true},
		{"EXECUTE p (1, 2)", true},
		{"MERGE INTO t USING u ON true WHEN NOT MATCHED THEN INSERT VALUES (1, 2)", true},

		// func_arg_list
		{"SELECT f(1, 2,)", false},
		{"SELECT f(a => 1, )", false},
		{"SELECT f(a => )", false},
		{"SELECT f(VARIADIC )", false},
		{"SELECT count(DISTINCT )", false},
		{"SELECT f()", true},
		{"SELECT f(1, VARIADIC ARRAY[2])", true},
		{"SELECT f(a => 1, b => 2)", true},

		// sortby_list
		{"SELECT 1 FROM t ORDER BY a,", false},
		{"SELECT 1 ORDER BY", false},
		{"SELECT 1 FROM t WINDOW w AS (ORDER BY a, )", false},
		{"SELECT 1 FROM t ORDER BY a, b DESC", true},

		// from_list
		{"SELECT * FROM t,", false},
		{"SELECT * FROM t, WHERE true", false},
		{"UPDATE t SET a = 1 FROM u,", false},
		{"DELETE FROM t USING u,", false},
		{"SELECT * FROM t, u", true},

		// columnList, index_params, ExclusionConstraintList
		{"CREATE VIEW v (a, ) AS SELECT 1", false},
		{"CREATE TABLE t (a int, UNIQUE (a, ))", false},
		{"CREATE TABLE t (a int, PRIMARY KEY (a) INCLUDE (a, ))", false},
		{"CREATE TABLE t (a int, FOREIGN KEY (a, ) REFERENCES u)", false},
		{"CREATE TABLE t (a int REFERENCES u (b, ))", false},
		{"CREATE TABLE t (a int, EXCLUDE USING gist (a WITH =, ))", false},
		{"CREATE INDEX i ON t (a, )", false},
		{"CREATE INDEX i ON t (a) INCLUDE (b, )", false},
		{"INSERT INTO t VALUES (1) ON CONFLICT (a, ) DO NOTHING", false},
		{"COPY t FROM stdin WITH CSV FORCE NOT NULL a,", false},
		{"CREATE TRIGGER tr AFTER UPDATE OF a, ON t FOR EACH ROW EXECUTE FUNCTION f()", false},
		{"WITH RECURSIVE a(x) AS (SELECT 1) SEARCH DEPTH FIRST BY x, SET o SELECT 1", false},
		{"WITH RECURSIVE a(x) AS (SELECT 1) CYCLE x, SET c USING p SELECT 1", false},
		{"CREATE VIEW v (a, b) AS SELECT 1, 2", true},
		{"CREATE INDEX i ON t (a, lower(b), (a + 1)) INCLUDE (c, d)", true},
		{"COPY t FROM stdin CSV FORCE QUOTE a, b", true},
		{"CREATE TRIGGER tr AFTER UPDATE OF a, b ON t FOR EACH ROW EXECUTE FUNCTION f()", true},

		// relation_expr_list
		{"TRUNCATE a,", false},
		{"LOCK TABLE a,", false},
		{"IMPORT FOREIGN SCHEMA s LIMIT TO (a, ) FROM SERVER x INTO y", false},
		{"TRUNCATE a, b", true},

		// GRANT lists
		{"GRANT SELECT, ON t TO a", false},
		{"GRANT SELECT ON t, TO a", false},
		{"GRANT SELECT ON TABLE a, TO r", false},
		{"GRANT USAGE ON SCHEMA a, TO r", false},
		{"GRANT USAGE ON TYPE a, TO r", false},
		{"GRANT CONNECT ON DATABASE a, TO r", false},
		{"GRANT SELECT ON ALL TABLES IN SCHEMA a, TO r", false},
		{"GRANT EXECUTE ON FUNCTION f(int, ) TO a", false},
		{"GRANT SELECT, INSERT ON t TO r", true},
		{"GRANT SELECT (a, b), UPDATE (a) ON t TO r", true},
		{"GRANT USAGE ON SCHEMA a, b TO r", true},
		{"GRANT EXECUTE ON FUNCTION f(int), g() TO r", true},

		// var_list, transaction_mode_list
		{"SET search_path = a,", false},
		{"SET search_path TO a, b,", false},
		{"BEGIN ISOLATION LEVEL SERIALIZABLE,", false},
		{"BEGIN ISOLATION LEVEL SERIALIZABLE READ ONLY,", false},
		{"SET search_path = a, b", true},
		{"BEGIN ISOLATION LEVEL SERIALIZABLE, READ ONLY", true},
		{"BEGIN ISOLATION LEVEL SERIALIZABLE READ ONLY", true},
	})
}

// TestTableConstraintSyntaxMatchesPG pins the constraint productions'
// own tokens. CHECK, PRIMARY, UNIQUE, FOREIGN, and CONSTRAINT are
// reserved keywords, so at the head of a table element they start a
// TableConstraint and never name a column (TableElement, ConstraintElem);
// CHECK takes its expression in parentheses that belong to the
// production (ColConstraintElem: CHECK '(' a_expr ')' opt_no_inherit,
// ConstraintElem: CHECK '(' a_expr ')' ConstraintAttributeSpec), as does
// a generated column ('(' a_expr ')' STORED). EXCLUDE is unreserved and
// starts a constraint only before USING or '('.
func TestTableConstraintSyntaxMatchesPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		// Reserved keywords as bare column names.
		{"CREATE TABLE t (id int, check int)", false},
		{"CREATE TABLE t (check int)", false},
		{"CREATE TABLE t (id int, primary int)", false},
		{"CREATE TABLE t (id int, unique int)", false},
		{"CREATE TABLE t (id int, foreign int)", false},
		{"CREATE TABLE t (id int, constraint int)", false},
		{"CREATE TABLE t (id int, select int)", false},
		{"CREATE FOREIGN TABLE t (id int, check int) SERVER s", false},
		{"ALTER TABLE t ADD check int", false},
		{"ALTER TABLE t ADD primary int", false},
		{"ALTER TABLE t ADD constraint int", false},
		{"ALTER TABLE t ADD COLUMN check int", false},
		{"CREATE TABLE t OF ty (check WITH OPTIONS NOT NULL)", false},
		{"CREATE TABLE t (id int, \"check\" int)", true},
		{"ALTER TABLE t ADD COLUMN \"check\" int", true},
		{"CREATE TABLE t (between int)", true},
		{"CREATE TABLE t (id int, like int)", true},
		{"CREATE TABLE t (exclude int)", true},
		{"CREATE TABLE t (id int, exclude int)", true},
		{"ALTER TABLE t ADD exclude int", true},
		{"CREATE TABLE t OF ty (exclude WITH OPTIONS NOT NULL)", true},
		{"CREATE TABLE t (exclude) AS SELECT 1", true},
		{"CREATE TABLE t (a int, exclude)", false},
		{"CREATE TABLE t (id int, exclude using gist (id with =))", true},
		{"CREATE TABLE t (id int, exclude (id with =))", true},
		{"ALTER TABLE t ADD exclude (a WITH =)", true},
		{"CREATE TABLE t (id int, primary key (id), unique (id), foreign key (id) references u)", true},
		{"CREATE TABLE t (id int, constraint c check (id > 0))", true},

		// CHECK takes parentheses.
		{"CREATE TABLE t (id int CHECK id > 0)", false},
		{"CREATE TABLE t (id int, CHECK id > 0)", false},
		{"CREATE TABLE t (id int CONSTRAINT c CHECK id > 0)", false},
		{"CREATE TABLE t (id int DEFAULT 1 CHECK id > 0)", false},
		{"CREATE TABLE t (id int, CHECK)", false},
		{"CREATE TABLE t (id int, CHECK ())", false},
		{"ALTER TABLE t ADD CHECK id > 0", false},
		{"ALTER TABLE t ADD CONSTRAINT c CHECK id > 0", false},
		{"CREATE DOMAIN d AS int CHECK VALUE > 0", false},
		{"CREATE TABLE t (id int CHECK (id > 0))", true},
		{"CREATE TABLE t (id int, CHECK (id > 0))", true},
		{"CREATE TABLE t (id int CHECK (id > 0) NO INHERIT)", true},
		{"ALTER TABLE t ADD CONSTRAINT c CHECK (id > 0)", true},
		{"CREATE DOMAIN d AS int CHECK (VALUE > 0)", true},

		// The rest of each constraint's tokens.
		{"CREATE TABLE t (a int GENERATED ALWAYS AS a + 1 STORED)", false},
		{"CREATE TABLE t (a int GENERATED ALWAYS AS (a + 1))", true},
		{"CREATE TABLE t (a int CONSTRAINT c)", false},
		{"CREATE TABLE t (a int CONSTRAINT c NOT 5)", false},
		{"CREATE TABLE t (a int PRIMARY)", false},
		{"CREATE TABLE t (a int, CONSTRAINT c)", false},
		{"CREATE TABLE t (a int, CONSTRAINT CHECK (a > 0))", false},
		{"ALTER TABLE t ADD CONSTRAINT c UNIQUE USING INDEX", false},
		{"CREATE TABLE t (a int, FOREIGN KEY (a) REFERENCES)", false},
		{"CREATE TABLE t (a int, FOREIGN KEY (a) u)", false},
		{"CREATE TABLE t (a int, EXCLUDE (a =))", false},
		{"CREATE TABLE t (a int, EXCLUDE (a WITH =) WHERE )", false},
		{"CREATE INDEX i ON t ((a + ))", false},
		{"CREATE TABLE t (a int GENERATED ALWAYS AS (a + 1) STORED)", true},
		{"CREATE TABLE t (a int GENERATED BY DEFAULT AS IDENTITY)", true},
		{"CREATE TABLE t (a int, UNIQUE USING INDEX i)", true},
		{"ALTER TABLE t ADD CONSTRAINT c UNIQUE USING INDEX i", true},
		{"CREATE TABLE t (a int, EXCLUDE (a WITH OPERATOR(pg_catalog.=)))", true},
		{"CREATE TABLE t (a text, EXCLUDE USING gist (a gist_trgm_ops WITH =) WHERE (a <> ''))", true},
	})
}

// TestSelectLimitMatchesPG pins select_limit and insertSelectOptions.
// select_no_parens takes, after its sort clause, a for_locking_clause and
// an opt_select_limit, or a select_limit and an opt_for_locking_clause;
// a select_limit is at most one limit_clause (LIMIT or FETCH) and one
// offset_clause, in either order, so a repeated clause is a syntax error
// at it. A parenthesized select keeps its own clauses, and
// insertSelectOptions refuses a second ORDER BY, OFFSET, LIMIT, or WITH
// on it at the second clause's expression, a second limit option, and
// WITH TIES without ORDER BY or with SKIP LOCKED. Each nested query has
// its own.
func TestSelectLimitMatchesPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		// One limit_clause and one offset_clause.
		{"SELECT 1 LIMIT 5 LIMIT 10", false},
		{"SELECT 1 OFFSET 1 OFFSET 2", false},
		{"SELECT 1 LIMIT ALL LIMIT 1", false},
		{"SELECT 1 LIMIT 1 FETCH FIRST 1 ROW ONLY", false},
		{"SELECT 1 FETCH FIRST 1 ROW ONLY LIMIT 1", false},
		{"SELECT 1 FETCH FIRST 1 ROW ONLY FETCH FIRST 2 ROWS ONLY", false},
		{"SELECT 1 OFFSET 1 LIMIT 1 OFFSET 2", false},
		{"SELECT 1 LIMIT 1 OFFSET 2 LIMIT 3", false},
		{"SELECT 1 OFFSET 1 ROWS OFFSET 2", false},
		{"SELECT 1 LIMIT", false},
		{"SELECT 1 OFFSET", false},
		{"SELECT * FROM (SELECT 1 LIMIT 1 LIMIT 2) s", false},
		{"SELECT 1 WHERE EXISTS (SELECT 1 LIMIT 1 LIMIT 2)", false},
		{"SELECT ARRAY(SELECT 1 LIMIT 1 LIMIT 2)", false},
		{"WITH a AS (SELECT 1) SELECT 1 LIMIT 1 LIMIT 2", false},
		{"INSERT INTO t SELECT 1 LIMIT 1 LIMIT 2", false},
		{"CREATE TABLE x AS SELECT 1 LIMIT 1 LIMIT 2", false},
		{"CREATE VIEW v AS SELECT 1 OFFSET 1 OFFSET 2", false},
		{"CREATE RECURSIVE VIEW v (a) AS SELECT 1 LIMIT 1 LIMIT 2", false},
		{"DECLARE c CURSOR FOR SELECT 1 LIMIT 1 LIMIT 2", false},
		{"TABLE t LIMIT 1 LIMIT 2", false},
		{"SELECT 1 LIMIT ALL", true},
		{"SELECT 1 LIMIT NULL", true},
		{"SELECT 1 LIMIT 1 OFFSET 2", true},
		{"SELECT 1 OFFSET 2 LIMIT 1", true},
		{"SELECT 1 OFFSET 1 FETCH FIRST 1 ROW ONLY", true},
		{"SELECT 1 FETCH FIRST 1 ROW ONLY OFFSET 1", true},
		{"VALUES (1) LIMIT 1 OFFSET 1", true},
		{"SELECT * FROM (SELECT 1 LIMIT 1) s LIMIT 2", true},
		{"SELECT 1 UNION SELECT 2 LIMIT 1", true},
		{"(SELECT 1 LIMIT 1) UNION (SELECT 2 LIMIT 1) LIMIT 1", true},
		{"(SELECT 1 LIMIT 1) UNION SELECT 2 LIMIT 1", true},
		{"SELECT 1 UNION SELECT 2 ORDER BY 1 LIMIT 1 OFFSET 2 FOR UPDATE", true},
		{"WITH a AS (SELECT 1) SELECT * FROM a ORDER BY 1 LIMIT 1 OFFSET 1", true},

		// Locking clauses on one side of the select_limit.
		{"SELECT 1 LIMIT 1 FOR UPDATE LIMIT 2", false},
		{"SELECT 1 FOR UPDATE LIMIT 2 FOR SHARE", false},
		{"SELECT 1 LIMIT 1 FOR UPDATE OFFSET 2", false},
		{"SELECT 1 OFFSET 1 FOR UPDATE LIMIT 2", false},
		{"SELECT 1 FOR READ ONLY FOR UPDATE", false},
		{"SELECT 1 FOR UPDATE LIMIT 2", true},
		{"SELECT 1 FOR UPDATE LIMIT 2 OFFSET 1", true},
		{"SELECT 1 LIMIT 2 FOR UPDATE FOR SHARE", true},
		{"SELECT 1 FOR UPDATE FOR SHARE LIMIT 2", true},
		{"SELECT 1 FOR READ ONLY LIMIT 1", true},
		{"SELECT 1 LIMIT 1 FOR READ ONLY", true},

		// LIMIT #,# is read only to be refused, at LIMIT.
		{"SELECT 1 LIMIT 1, 2", false},
		{"SELECT 1 LIMIT ALL, 2", false},
		{"SELECT 1 LIMIT 1,", false},

		// offset_clause: ROW or ROWS follows only a c_expr or a signed
		// numeric constant (select_fetch_first_value).
		{"SELECT 1 OFFSET 1 + 1 ROWS", false},
		{"SELECT 1 OFFSET 1::int ROWS", false},
		{"SELECT 1 OFFSET -a ROWS", false},
		{"SELECT 1 OFFSET -1::int ROWS", false},
		{"SELECT 1 OFFSET 1 ROW", true},
		{"SELECT 1 OFFSET 1 ROWS", true},
		{"SELECT 1 OFFSET (1 + 1) ROWS", true},
		{"SELECT 1 OFFSET -1 ROWS", true},
		{"SELECT 1 OFFSET +1 ROWS", true},
		{"SELECT 1 OFFSET - 1.5 ROWS", true},
		{"SELECT 1 OFFSET $1 ROWS", true},
		{"SELECT 1 OFFSET a[1] ROWS", true},
		{"SELECT 1 OFFSET (ARRAY[1])[1] ROWS", true},
		{"SELECT 1 OFFSET f(1) ROWS", true},
		{"SELECT 1 OFFSET 1 + 1", true},
		{"SELECT 1 OFFSET -1 * 2", true},
		{"SELECT 1 OFFSET NOT true", true},
		{"SELECT 1 OFFSET 1 IS NULL", true},

		// FETCH: FIRST or NEXT, then ROW or ROWS, are required.
		{"SELECT 1 FETCH FIRST 1 ONLY", false},
		{"SELECT 1 FETCH 1 ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST 1 ROWS", false},
		{"SELECT 1 FETCH FIRST +a ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST -a ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST 1+1 ROWS ONLY", false},
		{"SELECT 1 FETCH NEXT ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST ROW ONLY", true},
		{"SELECT 1 FETCH FIRST -1 ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST +1 ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST 1.5 ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST -1.5 ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST (1+1) ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST $1 ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST a[1] ROWS ONLY", true},

		// select_fetch_first_value's c_expr takes subscripts only as
		// columnref, PARAM, '(' a_expr ')', or select_with_parens
		// (gram.y c_expr); a constant, ARRAY, CASE, or function call
		// takes none.
		{"SELECT 1 FETCH FIRST ARRAY[1][1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST f(1)[1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST 1[1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST 'a'[1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST CASE WHEN true THEN ARRAY[1] END[1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST ARRAY(SELECT 1)[1] ROWS ONLY", false},
		{"SELECT 1 FETCH FIRST (ARRAY[1])[1] ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST $1[1] ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST a.b[1] ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST a[1][2] ROWS ONLY", true},
		{"SELECT 1 FETCH FIRST (SELECT ARRAY[1])[1] ROWS ONLY", true},
		{"SELECT 1 OFFSET $1[1] ROWS", true},
		{"SELECT 1 OFFSET (SELECT ARRAY[1])[1] ROWS", true},
		{"SELECT 1 ORDER BY 1 FETCH NEXT 2 ROW WITH TIES", true},

		// insertSelectOptions on a parenthesized select.
		{"(SELECT 1 LIMIT 1) LIMIT 2", false},
		{"(SELECT 1 LIMIT ALL) LIMIT 2", false},
		{"(SELECT 1 LIMIT 1) LIMIT ALL", false},
		{"(SELECT 1 FETCH FIRST 1 ROW ONLY) LIMIT 2", false},
		{"(SELECT 1 LIMIT 1) FETCH FIRST 2 ROWS ONLY", false},
		{"(SELECT 1 LIMIT 1) FETCH FIRST ROW ONLY", false},
		{"(SELECT 1 LIMIT 1) FETCH FIRST -2 ROWS ONLY", false},
		{"(SELECT 1 LIMIT 1) FETCH FIRST +2 ROWS ONLY", false},
		{"(SELECT 1 LIMIT 1) FETCH FIRST (2) ROWS ONLY", false},
		{"(SELECT 1 OFFSET 1) OFFSET 2", false},
		{"(SELECT 1 OFFSET 1) OFFSET 2 ROWS", false},
		{"(SELECT 1 OFFSET 1) OFFSET -2 ROWS", false},
		{"(SELECT 1 LIMIT 1 OFFSET 1) OFFSET 2", false},
		{"((SELECT 1 LIMIT 1)) LIMIT 2", false},
		{"(SELECT 1 ORDER BY 1) ORDER BY 1", false},
		{"(SELECT 1 ORDER BY 1) ORDER BY 2, 3", false},
		{"(SELECT 1 ORDER BY 1) ORDER BY a.b", false},
		{"(SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES) OFFSET 1", true},
		{"(SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES) LIMIT 1", false},
		{"SELECT 1 FETCH FIRST 1 ROW WITH TIES", false},
		{"WITH a AS (SELECT 1) SELECT * FROM a FETCH FIRST 1 ROW WITH TIES", false},
		{"SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES FOR UPDATE SKIP LOCKED", false},
		{"SELECT 1 ORDER BY 1 FOR UPDATE SKIP LOCKED FETCH FIRST 1 ROW WITH TIES", false},
		{"(SELECT 1 FOR UPDATE SKIP LOCKED) ORDER BY 1 FETCH FIRST 1 ROW WITH TIES", false},
		{"WITH a AS (SELECT 1) (WITH b AS (SELECT 2) SELECT 1)", false},
		{"WITH a AS (SELECT 1) (WITH b AS (SELECT 2) SELECT 1) ORDER BY 1", false},
		{"WITH a AS (SELECT 1) (SELECT 1 LIMIT 1) LIMIT 2", false},
		{"CREATE VIEW v AS (SELECT 1 LIMIT 1) LIMIT 2", false},
		{"CREATE MATERIALIZED VIEW v AS (SELECT 1 OFFSET 1) OFFSET 2", false},
		{"CREATE TABLE x AS (SELECT 1 LIMIT 1) LIMIT 2", false},
		{"(SELECT 1 LIMIT 1) OFFSET 2", true},
		{"(SELECT 1 OFFSET 1) LIMIT 2", true},
		{"(SELECT 1 ORDER BY 1) LIMIT 1", true},
		{"(SELECT 1) ORDER BY 1 LIMIT 1", true},
		{"(SELECT 1 FOR UPDATE) FOR SHARE", true},
		{"(SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW ONLY) OFFSET 1", true},
		{"(SELECT 1 ORDER BY 1) FETCH FIRST 1 ROW WITH TIES", true},
		{"SELECT 1 ORDER BY 1 FETCH FIRST 1 ROW WITH TIES", true},

		// Where the second clause's error points: the server's
		// exprLocation, which skips the parentheses around an
		// expression but not those of a subquery.
		{"(SELECT 1 LIMIT 1) LIMIT (2)", false},
		{"(SELECT 1 LIMIT 1) LIMIT ((2))", false},
		{"(SELECT 1 LIMIT 1) LIMIT ( /* c */ 2)", false},
		{"(SELECT 1 LIMIT 1) LIMIT (2)::int", false},
		{"(SELECT 1 LIMIT 1) LIMIT (SELECT 2)", false},
		{"(SELECT 1 LIMIT 1) LIMIT ((SELECT 2) + 1)", false},
		{"(SELECT 1 LIMIT 1) LIMIT (-2)", false},
		{"(SELECT 1 LIMIT 1) LIMIT 1 + 1", false},
		{"(SELECT 1 LIMIT 1) LIMIT -1", false},
		{"(SELECT 1 LIMIT 1) LIMIT x.y", false},
		{"(SELECT 1 LIMIT 1) LIMIT f(1)", false},
		{"(SELECT 1 LIMIT 1) LIMIT CAST(1 AS int)", false},
		{"(SELECT 1 OFFSET 1) OFFSET (2)", false},
		{"(SELECT 1 OFFSET 1) OFFSET (ARRAY[1])[1] ROWS", false},
		{"(SELECT 1 ORDER BY 1) ORDER BY (1) DESC", false},
		{"(SELECT 1 ORDER BY 1) ORDER BY (a).b", false},
	})
}

// TestSelectLimitCompletion pins the candidates after a LIMIT, OFFSET, or
// FETCH value. A value that may go on as an a_expr offers the operators
// that continue it, as after any a_expr; the clauses offered next are
// those select_limit still takes, so no second LIMIT or OFFSET; and FETCH
// offers the FIRST or NEXT and ROW or ROWS it requires.
func TestSelectLimitCompletion(t *testing.T) {
	operators := []int{'+', '*', AND, IS}
	tests := []struct {
		sql    string
		want   []int
		refuse []int
	}{
		{"SELECT 1 OFFSET 1 ", append([]int{LIMIT, FETCH, FOR, ';'}, operators...), []int{OFFSET}},
		{"SELECT 1 OFFSET a ", append([]int{LIMIT, FETCH, FOR, ';'}, operators...), []int{OFFSET}},
		{"SELECT 1 OFFSET (1) ", append([]int{LIMIT, FETCH, FOR, ';'}, operators...), []int{OFFSET}},
		{"SELECT 1 OFFSET -1 ", append([]int{LIMIT, FETCH, FOR, ';'}, operators...), []int{OFFSET}},
		{"SELECT 1 LIMIT 1 OFFSET 1 ", append([]int{FOR, ';'}, operators...), []int{LIMIT, OFFSET, FETCH}},
		{"SELECT 1 LIMIT 1 ", append([]int{OFFSET, FOR, ';'}, operators...), []int{LIMIT, FETCH}},
		{"SELECT 1 OFFSET 1 ROWS ", []int{LIMIT, FETCH, FOR, ';'}, append([]int{OFFSET}, operators...)},
		{"SELECT 1 FETCH ", []int{FIRST_P, NEXT}, nil},
		{"SELECT 1 FETCH FIRST 1 ", []int{ROW, ROWS}, []int{ONLY}},
		{"SELECT 1 FETCH FIRST 1 ROWS ", []int{ONLY}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			for _, tok := range tt.want {
				if !cs.HasToken(tok) {
					t.Errorf("missing candidate %s", TokenName(tok))
				}
			}
			for _, tok := range tt.refuse {
				if cs.HasToken(tok) {
					t.Errorf("unexpected candidate %s", TokenName(tok))
				}
			}
		})
	}
}

// TestIsSelectLimitStartMatchesSlice checks isSelectLimitStart against
// every keyword: it accepts exactly selectLimitLeadTokens.
func TestIsSelectLimitStartMatchesSlice(t *testing.T) {
	for _, kw := range Keywords {
		p := &Parser{cur: Token{Type: kw.Token, Str: kw.Name}}
		if got, want := p.isSelectLimitStart(), selectLimitLeadSet[kw.Token]; got != want {
			t.Errorf("isSelectLimitStart(%q) = %v, want %v", kw.Name, got, want)
		}
	}
	for _, tok := range selectLimitLeadTokens {
		p := &Parser{cur: Token{Type: tok}}
		if !p.isSelectLimitStart() {
			t.Errorf("isSelectLimitStart must accept token %d", tok)
		}
	}
}

// TestSelectOptionsShape pins the tree each accepted select_limit and
// locking combination builds.
func TestSelectOptionsShape(t *testing.T) {
	intVal := func(n nodes.Node) int64 {
		t.Helper()
		c, ok := n.(*nodes.A_Const)
		if !ok {
			t.Fatalf("got %T, want *A_Const", n)
		}
		i, ok := c.Val.(*nodes.Integer)
		if !ok {
			t.Fatalf("got %T, want *Integer", c.Val)
		}
		return i.Ival
	}
	sel := func(sql string) *nodes.SelectStmt {
		t.Helper()
		s, ok := singleStmt(t, sql).(*nodes.SelectStmt)
		if !ok {
			t.Fatalf("%s: not a SelectStmt", sql)
		}
		return s
	}

	for _, sql := range []string{"SELECT 1 LIMIT 3 OFFSET 4", "SELECT 1 OFFSET 4 LIMIT 3", "SELECT 1 OFFSET 4 ROWS FETCH NEXT 3 ROWS ONLY"} {
		s := sel(sql)
		if intVal(s.LimitCount) != 3 || intVal(s.LimitOffset) != 4 || s.LimitOption != nodes.LIMIT_OPTION_COUNT {
			t.Errorf("%s: count/offset/option = %v/%v/%v", sql, s.LimitCount, s.LimitOffset, s.LimitOption)
		}
	}

	s := sel("SELECT 1 FETCH FIRST ROW ONLY")
	if intVal(s.LimitCount) != 1 {
		t.Errorf("FETCH FIRST ROW ONLY: count = %v, want 1", s.LimitCount)
	}
	s = sel("SELECT 1 FETCH FIRST -2 ROWS ONLY")
	if intVal(s.LimitCount) != -2 {
		t.Errorf("FETCH FIRST -2 ROWS ONLY: count = %v, want -2", s.LimitCount)
	}
	s = sel("SELECT 1 LIMIT ALL")
	if c, ok := s.LimitCount.(*nodes.A_Const); !ok || !c.Isnull {
		t.Errorf("LIMIT ALL: count = %#v, want a NULL constant", s.LimitCount)
	}
	s = sel("SELECT 1 ORDER BY 1 FETCH FIRST 2 ROWS WITH TIES")
	if s.LimitOption != nodes.LIMIT_OPTION_WITH_TIES || intVal(s.LimitCount) != 2 {
		t.Errorf("WITH TIES: option/count = %v/%v", s.LimitOption, s.LimitCount)
	}

	// A parenthesized select keeps its clauses, and the outer ones that
	// do not repeat them are added to the same statement.
	s = sel("(SELECT 1 ORDER BY 1 LIMIT 2) OFFSET 3")
	if s.SortClause == nil || intVal(s.LimitCount) != 2 || intVal(s.LimitOffset) != 3 {
		t.Errorf("(... LIMIT 2) OFFSET 3: sort/count/offset = %v/%v/%v", s.SortClause, s.LimitCount, s.LimitOffset)
	}
	// Locking clauses accumulate.
	s = sel("(SELECT 1 FOR UPDATE) FOR SHARE")
	if s.LockingClause == nil || len(s.LockingClause.Items) != 2 {
		t.Fatalf("(... FOR UPDATE) FOR SHARE: locking = %v, want 2 clauses", s.LockingClause)
	}
	for i, want := range []nodes.LockClauseStrength{nodes.LCS_FORUPDATE, nodes.LCS_FORSHARE} {
		if lc := s.LockingClause.Items[i].(*nodes.LockingClause); lc.Strength != int(want) {
			t.Errorf("locking clause %d strength = %d, want %d", i, lc.Strength, want)
		}
	}
	// A top-level WITH reaches the statement with its sort and limit, and
	// the statement spans from WITH to its last token.
	sql := "WITH a AS (SELECT 1) SELECT * FROM a ORDER BY 1 LIMIT 2"
	s = sel(sql)
	if s.WithClause == nil || s.SortClause == nil || intVal(s.LimitCount) != 2 {
		t.Errorf("%s: with/sort/count = %v/%v/%v", sql, s.WithClause, s.SortClause, s.LimitCount)
	}
	if s.Loc.Start != 0 || s.Loc.End != len(sql) {
		t.Errorf("%s: Loc = %+v, want {0 %d}", sql, s.Loc, len(sql))
	}
}

// TestExcludeColumnShape pins that EXCLUDE names a column unless USING
// or '(' follows it.
func TestExcludeColumnShape(t *testing.T) {
	cs, ok := singleStmt(t, "CREATE TABLE t (exclude int, exclude2 int, EXCLUDE USING gist (exclude WITH =))").(*nodes.CreateStmt)
	if !ok {
		t.Fatal("not a CreateStmt")
	}
	if len(cs.TableElts.Items) != 3 {
		t.Fatalf("got %d table elements, want 3", len(cs.TableElts.Items))
	}
	if cd, ok := cs.TableElts.Items[0].(*nodes.ColumnDef); !ok || cd.Colname != "exclude" {
		t.Errorf("element 0 = %#v, want column exclude", cs.TableElts.Items[0])
	}
	if c, ok := cs.TableElts.Items[2].(*nodes.Constraint); !ok || c.Contype != nodes.CONSTR_EXCLUSION || c.AccessMethod != "gist" {
		t.Errorf("element 2 = %#v, want a gist exclusion constraint", cs.TableElts.Items[2])
	}
}
