package parser

import "testing"

// parenOracleSchema is the schema every parenthesized-expression oracle case
// runs against: the BYT-10316 report's tables a and b, plus t and u for the
// generic shapes.
var parenOracleSchema = []string{
	"CREATE TABLE a (id INT, k VARCHAR(20))",
	"CREATE TABLE b (id INT, k VARCHAR(20))",
	"CREATE TABLE t (a INT, b INT, c INT, d INT)",
	"CREATE TABLE u (x INT)",
}

// parenOracleCorpus pins how a leading run of '(' is read. The deciding token
// is the first one past the run: a query primary opens a subquery, anything
// else a parenthesized scalar expression or, in a FROM clause, a join group.
// Every case is checked against a live MariaDB 11.8.8, which settles want.
var parenOracleCorpus = []struct {
	sql  string
	want bool
}{
	// BYT-10316: the report's statements, rejected at columns 33 and 43 by a
	// lookahead that read any '((' as the start of a subquery. The last is the
	// view definition MySQL stores; MariaDB stores the same view as
	// `on(a and b)`, so here the shape arrives only as SQL a user writes.
	{"SELECT count(*) FROM a WHERE ((a.id = 1) AND (a.k = 'x'))", true},
	{"SELECT count(*) FROM a WHERE (((a.id = 1) AND (a.k = 'x')))", true},
	{"SELECT count(*) FROM a t1 JOIN b t2 ON (((t1.id = t2.id) AND (t1.k = t2.k)))", true},
	{"select `t1`.`id` AS `id`,`t2`.`k` AS `k` from (`a` `t1` join `b` `t2` on(((`t1`.`id` = `t2`.`id`) and (`t1`.`k` = `t2`.`k`) and (`t2`.`k` <> 'x') and (`t2`.`k` <> 'y'))))", true},

	// Parenthesized scalar expressions.
	{"SELECT ((1))", true},
	{"SELECT (((1 + 2)))", true},
	{"SELECT ((a = b)) FROM t", true},
	{"SELECT * FROM t WHERE ((((a = b))))", true},
	{"SELECT ((a)) + ((b)) FROM t", true},
	{"SELECT * FROM t WHERE a IN ((1), (2))", true},
	{"SELECT (SELECT ((1)))", true},

	// Parenthesized join groups nested deeper than one '('.
	{"SELECT * FROM (((a JOIN b ON a.id = b.id) JOIN t ON b.id = t.a) JOIN u ON t.a = u.x)", true},

	// Subqueries behind a run of '(' stay subqueries.
	{"SELECT ((SELECT 1))", true},
	{"SELECT (((SELECT 1)))", true},
	{"SELECT * FROM t WHERE a IN ((SELECT x FROM u))", true},
	{"SELECT * FROM t WHERE EXISTS ((SELECT 1))", true},
	{"SELECT * FROM t WHERE a > ALL ((SELECT x FROM u))", true},
	{"SELECT * FROM ((SELECT 1)) x", true},
	{"SELECT * FROM (((SELECT 1))) x", true},
	{"SELECT * FROM ((SELECT 1) UNION (SELECT 2)) x", true},

	// EXISTS and a quantified comparison take only a subquery, however many
	// parens wrap a scalar.
	{"SELECT * FROM t WHERE EXISTS ((1))", false},
	{"SELECT * FROM t WHERE a > ALL ((1))", false},
}

// TestParenthesizedExprOracle asserts omni and a live MariaDB 11.8.8 agree on
// every parenOracleCorpus case, and that both match its want.
func TestParenthesizedExprOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, ddl := range parenOracleSchema {
		if _, err := o.db.ExecContext(o.ctx, ddl); err != nil {
			t.Fatalf("seed %q: %v", ddl, err)
		}
	}

	for _, tc := range parenOracleCorpus {
		v, code, msg := o.classify(tc.sql)
		if v == vErrored {
			t.Errorf("ERRORED (infra, not a parse verdict) for %q: %s", tc.sql, msg)
			continue
		}
		if mdbAccepts := v == vAccepted; mdbAccepts != tc.want {
			t.Errorf("mariadb accepts=%v (code %d), want %v, for %q: %s", mdbAccepts, code, tc.want, tc.sql, msg)
		}
		_, err := Parse(tc.sql)
		if omniAccepts := err == nil; omniAccepts != tc.want {
			t.Errorf("omni accepts=%v, want %v, for %q: %v", omniAccepts, tc.want, tc.sql, err)
		}
	}
}
