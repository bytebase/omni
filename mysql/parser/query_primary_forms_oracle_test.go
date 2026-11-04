package parser

import "testing"

// The TABLE and VALUES query primaries, the locking-read clause, and the INTO
// that follows them, held to a live MySQL 8.0. One case per form this parser
// reads and per malformed variant around it: a clause the parser swallows
// where MySQL rejects it leaves no trace in a suite that only checks what
// parses. These rows are where the clause order and the positions MySQL takes
// INTO in are pinned.
//
// The list is close to the same text in the sibling engine trees. What differs
// is this column, which is where the forks part: MySQL has no WAIT n modifier
// (TiDB, MariaDB and OceanBase each do), and it takes INTO both before and
// after the locking clause (TiDB and OceanBase take it only last).
//
// Measured on 2026-11-04 against mysql:8.0 via PREPARE.

// queryPrimaryCase is a statement and whether the engine's parser accepts it.
type queryPrimaryCase struct {
	sql    string
	accept bool
}

var queryPrimaryCases = []queryPrimaryCase{
	// TABLE and VALUES as statements, and the clauses after them.
	{"TABLE t", true},
	{"TABLE t ORDER BY a LIMIT 2,3", true},
	{"TABLE t ORDER BY a LIMIT 1 FOR UPDATE", true},
	{"TABLE t FOR SHARE", true},
	{"TABLE t LOCK IN SHARE MODE", true},
	{"TABLE t FOR UPDATE NOWAIT", true},
	{"TABLE t FOR UPDATE SKIP LOCKED", true},
	{"TABLE t FOR UPDATE OF t", true},
	{"TABLE t INTO OUTFILE '/tmp/o1.txt'", true},
	{"TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'", true},
	{"TABLE t INTO OUTFILE '/tmp/a2.txt' FOR UPDATE", true},
	{"TABLE t INTO OUTFILE '/tmp/o3.txt' ORDER BY a", false},
	{"TABLE t INTO OUTFILE '/tmp/pj.txt' LIMIT 1", false},
	{"TABLE t ORDER BY a ORDER BY b", false},
	{"TABLE t LIMIT 1 LIMIT 2", false},
	{"TABLE t1, t2", false},
	{"TABLE t (a)", false},
	{"VALUES ROW(1,2)", true},
	{"VALUES ROW(1,2) ORDER BY a LIMIT 1", true},
	{"VALUES ROW(1) FOR UPDATE", true},
	{"VALUES ROW(1) FOR UPDATE NOWAIT", true},
	{"VALUES ROW(1) LOCK IN SHARE MODE", true},
	{"VALUES ROW(1) FOR SHARE", true},
	{"VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'", true},
	{"VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'", true},
	{"VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'", true},
	{"VALUES ROW(1,2) INTO OUTFILE '/tmp/v3.txt' ORDER BY a", false},
	{"VALUES ROW(1) ORDER BY a ORDER BY b", false},
	{"VALUES (1,2)", false},
	{"VALUES ROW(1), (2)", false},

	// TABLE and VALUES as set-operation operands.
	{"TABLE t1 UNION TABLE t2", true},
	{"TABLE t1 UNION ALL TABLE t2", true},
	{"TABLE t1 UNION (TABLE t2)", true},
	{"TABLE t1 UNION SELECT * FROM t2", true},
	{"TABLE t1 EXCEPT TABLE t2", true},
	{"TABLE t1 INTERSECT TABLE t2", true},
	{"TABLE t1 UNION TABLE t2 ORDER BY c1", true},
	{"TABLE t1 UNION TABLE t2 FOR UPDATE", true},
	{"VALUES ROW(1,2) UNION VALUES ROW(3,4)", true},
	{"VALUES ROW(1,2) UNION SELECT 1,2", true},
	{"VALUES ROW(1,2) UNION TABLE t1", true},
	{"TABLE t1 UNION VALUES ROW(1,2)", true},
	{"VALUES ROW(1,2) EXCEPT VALUES ROW(3,4)", true},
	{"VALUES ROW(1,2) INTERSECT VALUES ROW(3,4)", true},
	{"WITH cte AS (SELECT 1) TABLE t", true},
	{"WITH cte AS (SELECT 1) VALUES ROW(1)", true},
	{"WITH cte AS (TABLE t) TABLE cte", true},
	{"WITH RECURSIVE cte AS (TABLE t) TABLE cte", true},
	{"EXPLAIN TABLE t1 UNION TABLE t2", true},
	{"INSERT INTO t1 TABLE t2 ORDER BY a", true},
	{"INSERT INTO t1 TABLE t2 UNION TABLE t3", true},
	{"CREATE TABLE tt3 AS TABLE t UNION TABLE t", true},

	// The locking-read clause. MySQL has no WAIT n.
	{"SELECT * FROM t FOR UPDATE", true},
	{"SELECT * FROM t FOR SHARE", true},
	{"SELECT * FROM t LOCK IN SHARE MODE", true},
	{"SELECT * FROM t FOR UPDATE NOWAIT", true},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED", true},
	{"SELECT * FROM t FOR UPDATE OF t", true},
	{"SELECT * FROM t FOR UPDATE OF t, t1", true},
	{"SELECT * FROM t FOR UPDATE OF t NOWAIT", true},
	{"SELECT * FROM t FOR SHARE OF t SKIP LOCKED", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE OF t WAIT 5", false},
	{"SELECT * FROM t LOCK IN SHARE MODE NOWAIT", false},
	{"SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED", false},
	{"SELECT * FROM t FOR UPDATE SKIP", false},
	{"SELECT * FROM t FOR UPDATE OF", false},
	{"SELECT * FROM t FOR UPDATE NOWAIT NOWAIT", false},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED SKIP LOCKED", false},

	// The INTO clause, which the locking clause may follow but no other clause may.
	{"SELECT 1 INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'", true},
	{"SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'", true},
	{"SELECT a INTO OUTFILE '/tmp/p1.txt' FROM t", true},
	{"SELECT a INTO @v FROM t", true},
	{"SELECT a FROM t INTO @v", true},
	{"SELECT a FROM t INTO DUMPFILE '/tmp/pd.txt'", true},
	{"SELECT * FROM t INTO OUTFILE '/tmp/f5.txt' ORDER BY a", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/p7.txt' LIMIT 1", false},
	{"SELECT a FROM t ORDER BY a INTO OUTFILE '/tmp/pb.txt' LIMIT 1", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/pe.txt' INTO OUTFILE '/tmp/pf.txt'", false},

	// The locking-clause keywords as names.
	{"SELECT * FROM t wait", true},
	{"SELECT * FROM t AS wait", true},
	{"DELETE FROM t wait WHERE wait.a = 1", true},
	{"DELETE FROM t locked WHERE locked.a = 1", true},
	{"DELETE FROM t nowait WHERE nowait.a = 1", true},
	{"UPDATE t wait SET wait.a = 1", true},
	{"CREATE TABLE nowait (locked int, wait int, `row` int)", true},
}

// TestQueryPrimaryFormsOracle asserts omni agrees with a live MySQL 8.0 on
// queryPrimaryCases (container test).
func TestQueryPrimaryFormsOracle(t *testing.T) {
	o := startParserOracle(t)
	for _, setup := range []string{
		"CREATE TABLE t (a INT, b INT)",
		"CREATE TABLE t1 (a INT, b INT, c1 INT, c2 INT)",
		"CREATE TABLE t2 (a INT, b INT, c1 INT, c2 INT)",
		"CREATE TABLE t3 (a INT, b INT, c1 INT, c2 INT)",
	} {
		if _, err := o.db.ExecContext(o.ctx, setup); err != nil {
			t.Fatalf("%s: %v", setup, err)
		}
	}
	for _, qpc := range queryPrimaryCases {
		t.Run(qpc.sql, func(t *testing.T) {
			if accepts := o.prepareAccepts(t, qpc.sql); accepts != qpc.accept {
				t.Fatalf("MySQL 8.0 accepts=%v, want %v for %q", accepts, qpc.accept, qpc.sql)
			}
			_, perr := Parse(qpc.sql)
			if omniAccepts := perr == nil; omniAccepts != qpc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MySQL 8.0 accepts=%v for %q: %v",
					omniAccepts, qpc.accept, qpc.sql, perr)
			}
		})
	}
}

// TestQueryPrimaryForms holds the parser to the recorded verdicts and needs no
// server, so it is the gate that runs without Docker.
func TestQueryPrimaryForms(t *testing.T) {
	for _, qpc := range queryPrimaryCases {
		t.Run(qpc.sql, func(t *testing.T) {
			_, err := Parse(qpc.sql)
			if got := err == nil; got != qpc.accept {
				t.Errorf("parser accepts=%v, want %v for %q: %v", got, qpc.accept, qpc.sql, err)
			}
		})
	}
}
