package parser

import (
	"errors"
	"testing"
)

// strictCase is a statement and the error Parse refuses it with, at the
// byte offset pos; message is "" when Parse accepts it.
type strictCase struct {
	sql     string
	message string
	pos     int
}

func assertStrictCases(t *testing.T, cases []strictCase) {
	t.Helper()
	for _, tt := range cases {
		t.Run(tt.sql, func(t *testing.T) {
			_, err := Parse(tt.sql)
			if tt.message == "" {
				if err != nil {
					t.Fatalf("Parse: %v", err)
				}
				return
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse error = %v, want %q", err, tt.message)
			}
			if pe.Message != tt.message || pe.Position != tt.pos {
				t.Errorf("error = %q at %d, want %q at %d", pe.Message, pe.Position, tt.message, tt.pos)
			}
		})
	}
}

// TestListElements pins that a comma-separated list takes no empty
// element, as pg/parser's TestListElementsMatchPG does against PostgreSQL
// 17, from whose grammar this fork comes. Redshift's synopses write every
// such list as `item [, ...]` (r_SELECT_synopsis, r_SELECT_list,
// r_INSERT_30, r_CREATE_TABLE_NEW, r_GRANT), so a comma promises an item.
// No Redshift server is reachable here to say more.
func TestListElements(t *testing.T) {
	assertStrictCases(t, []strictCase{
		{"SELECT a, b, FROM t", `syntax error at or near "FROM"`, 13},
		{"SELECT a, FROM t", `syntax error at or near "FROM"`, 10},
		{"INSERT INTO t VALUES (1, 2,)", `syntax error at or near ")"`, 27},
		{"INSERT INTO t VALUES ()", `syntax error at or near ")"`, 22},
		{"SELECT 1 FROM t ORDER BY a,", "syntax error at end of input", 27},
		{"SELECT f(1, 2,)", `syntax error at or near ")"`, 14},
		{"SELECT 1 WHERE 1 IN (1, 2,)", `syntax error at or near ")"`, 26},
		{"SELECT (1, 2,)", `syntax error at or near ")"`, 13},
		{"SELECT * FROM t, WHERE true", `syntax error at or near "WHERE"`, 17},
		{"CREATE TABLE t (a int, UNIQUE (a, ))", `syntax error at or near ")"`, 34},
		{"CREATE TABLE t (a int) SORTKEY (a, )", `syntax error at or near ")"`, 35},
		{"CREATE INDEX i ON t (a, )", `syntax error at or near ")"`, 24},
		{"TRUNCATE a,", "syntax error at end of input", 11},
		{"GRANT SELECT ON t, TO a", `syntax error at or near "TO"`, 19},
		{"SELECT a, b FROM t ORDER BY a, b", "", 0},
		{"SELECT FROM t", "", 0},
		{"INSERT INTO t VALUES (1, 2), (3, 4)", "", 0},
		{"CREATE TABLE t (a int, b int) COMPOUND SORTKEY (a, b)", "", 0},
		{"GRANT SELECT, INSERT ON t, u TO a", "", 0},
	})
}

// TestTableConstraintSyntax pins the constraint productions' own tokens,
// as pg/parser's TestTableConstraintSyntaxMatchesPG does against
// PostgreSQL 17. CHECK, CONSTRAINT, FOREIGN, PRIMARY, and UNIQUE are
// Redshift reserved words (r_pg_keywords), usable as column names only
// when delimited, and CHECK, which Redshift's CREATE TABLE does not list
// among its constraints (r_CREATE_TABLE_NEW), takes its expression in
// parentheses in the grammar this fork comes from. EXCLUDE is not
// reserved, so it names a column unless USING or '(' follows it.
func TestTableConstraintSyntax(t *testing.T) {
	assertStrictCases(t, []strictCase{
		{"CREATE TABLE t (id int, check int)", `syntax error at or near "int"`, 30},
		{"CREATE TABLE t (id int, primary int)", `syntax error at or near "int"`, 32},
		{"CREATE TABLE t (id int, constraint int)", `syntax error at or near ")"`, 38},
		{"ALTER TABLE t ADD check int", `syntax error at or near "int"`, 24},
		{"CREATE TABLE t (id int CHECK id > 0)", `syntax error at or near "id"`, 29},
		{"CREATE TABLE t (id int, CHECK id > 0)", `syntax error at or near "id"`, 30},
		{"CREATE TABLE t (id int, \"check\" int)", "", 0},
		{"CREATE TABLE t (id int CHECK (id > 0))", "", 0},
		{"CREATE TABLE t (exclude int)", "", 0},
		{"CREATE TABLE t (id int, PRIMARY KEY (id), UNIQUE (id), FOREIGN KEY (id) REFERENCES u (id))", "", 0},
	})
}

// TestSelectLimit pins that a SELECT takes one LIMIT and one OFFSET, as
// Redshift's synopses give them (r_SELECT_synopsis, r_ORDER_BY_clause:
// [ LIMIT { number | ALL } ] [ OFFSET start ]) and as select_limit in the
// PostgreSQL grammar this fork comes from does; pg/parser's
// TestSelectLimitMatchesPG checks the same shape against PostgreSQL 17.
// gram.y's checks on a parenthesized select's own clauses, its LIMIT #,#
// refusal, and the FETCH forms' required words are not applied here: the
// Redshift reference does not settle them.
func TestSelectLimit(t *testing.T) {
	assertStrictCases(t, []strictCase{
		{"SELECT 1 LIMIT 5 LIMIT 10", `syntax error at or near "LIMIT"`, 17},
		{"SELECT 1 OFFSET 1 OFFSET 2", `syntax error at or near "OFFSET"`, 18},
		{"SELECT 1 LIMIT 1 OFFSET 2 LIMIT 3", `syntax error at or near "LIMIT"`, 26},
		{"SELECT 1 OFFSET 1 LIMIT 1 OFFSET 2", `syntax error at or near "OFFSET"`, 26},
		{"SELECT * FROM (SELECT 1 LIMIT 1 LIMIT 2) s", `syntax error at or near "LIMIT"`, 32},
		{"WITH a AS (SELECT 1) SELECT 1 LIMIT 1 LIMIT 2", `syntax error at or near "LIMIT"`, 38},
		{"SELECT 1 LIMIT", "syntax error at end of input", 14},
		{"SELECT 1 OFFSET", "syntax error at end of input", 15},
		{"SELECT 1 LIMIT ALL", "", 0},
		{"SELECT 1 LIMIT 1 OFFSET 2", "", 0},
		{"SELECT 1 OFFSET 2 LIMIT 1", "", 0},
		{"SELECT * FROM (SELECT 1 LIMIT 1) s LIMIT 2", "", 0},
		{"SELECT 1 UNION SELECT 2 ORDER BY 1 LIMIT 1 OFFSET 2", "", 0},
		{"WITH a AS (SELECT 1) SELECT * FROM a ORDER BY 1 LIMIT 1 OFFSET 1", "", 0},
	})
}

// TestSelectLimitCompletion pins the candidates after a LIMIT or OFFSET
// value: the operators that continue it as an a_expr, and only the clauses
// a SELECT still takes, so no second LIMIT or OFFSET.
func TestSelectLimitCompletion(t *testing.T) {
	operators := []int{'+', '*', AND, IS}
	tests := []struct {
		sql    string
		want   []int
		refuse []int
	}{
		{"SELECT 1 OFFSET 1 ", append([]int{LIMIT, FOR, ';'}, operators...), []int{OFFSET}},
		{"SELECT 1 LIMIT 1 OFFSET 1 ", append([]int{FOR, ';'}, operators...), []int{LIMIT, OFFSET}},
		{"SELECT 1 LIMIT 1 ", append([]int{OFFSET, FOR, ';'}, operators...), []int{LIMIT}},
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
