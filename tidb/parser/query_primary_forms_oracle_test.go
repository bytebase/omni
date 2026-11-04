package parser

import "testing"

// The TABLE and VALUES query primaries, the locking-read clause, and the INTO
// that follows them, held to the pinned TiDB v8.5.5. One case per form this
// parser reads and per malformed variant around it, because the corpus in this
// package records only whether a statement parses and holds far more accept
// rows than reject rows: a clause the parser swallows where TiDB rejects it is
// invisible there. These rows are where the wait modifiers, the clause order,
// and the one position TiDB takes INTO in are pinned.
//
// Measured on 2026-11-04 against pingcap/tidb:v8.5.5 via PREPARE. The
// expectation is what the server answered, so an image upgrade that changes a
// verdict fails here instead of passing silently.

// queryPrimaryCase is a statement and whether the engine's parser accepts it.
type queryPrimaryCase struct {
	sql    string
	accept bool
}

var queryPrimaryCases = []queryPrimaryCase{
	// TABLE as a statement and the clauses after it.
	{"TABLE t", true},
	{"TABLE t ORDER BY a", true},
	{"TABLE t LIMIT 3", true},
	{"TABLE t ORDER BY a LIMIT 2,3", true},
	{"TABLE t ORDER BY a LIMIT 1 FOR UPDATE", true},
	{"TABLE t FOR SHARE", true},
	{"TABLE t LOCK IN SHARE MODE", true},
	{"TABLE t FOR UPDATE NOWAIT", true},
	{"TABLE t FOR UPDATE SKIP LOCKED", true},
	{"TABLE t FOR UPDATE WAIT 5", true},
	{"TABLE t FOR UPDATE OF t", true},
	{"TABLE t INTO OUTFILE '/tmp/o1.txt'", true},
	{"TABLE t ORDER BY a INTO OUTFILE '/tmp/o2.txt'", true},
	{"TABLE t FOR UPDATE INTO OUTFILE '/tmp/a1.txt'", true},
	{"TABLE t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/a3.txt'", true},
	{"TABLE t INTO OUTFILE '/tmp/o3.txt' ORDER BY a", false},
	{"TABLE t INTO OUTFILE '/tmp/a2.txt' FOR UPDATE", false},
	{"TABLE t INTO OUTFILE '/tmp/pj.txt' LIMIT 1", false},
	{"TABLE t ORDER BY a ORDER BY b", false},
	{"TABLE t LIMIT 1 LIMIT 2", false},
	{"TABLE t1, t2", false},
	{"TABLE t (a)", false},
	{"TABLE t WHERE a = 1", false},
	{"SELECT * FROM TABLE t", false},
	// VALUES as a statement and the clauses after it.
	{"VALUES ROW(1,2)", true},
	{"VALUES ROW()", true},
	{"VALUES ROW(1, default)", true},
	{"VALUES ROW(1,2), ROW(3)", true},
	{"VALUES ROW(1,2) ORDER BY a", true},
	{"VALUES ROW(1,2) LIMIT 1", true},
	{"VALUES ROW(1,2) ORDER BY a LIMIT 1", true},
	{"VALUES ROW(1) FOR UPDATE", true},
	{"VALUES ROW(1) FOR UPDATE NOWAIT", true},
	{"VALUES ROW(1) FOR UPDATE WAIT 5", true},
	{"VALUES ROW(1) LOCK IN SHARE MODE", true},
	{"VALUES ROW(1) FOR SHARE", true},
	{"VALUES ROW(1,-2,3), ROW(5,7,9) INTO OUTFILE '/tmp/v1.txt'", true},
	{"VALUES ROW(1,-2,3), ROW(5,7,9) ORDER BY a INTO OUTFILE '/tmp/v2.txt'", true},
	{"VALUES ROW(1) FOR UPDATE INTO OUTFILE '/tmp/a4.txt'", true},
	{"VALUES ROW(1,2) INTO OUTFILE '/tmp/v3.txt' ORDER BY a", false},
	{"VALUES ROW(1) INTO OUTFILE '/tmp/a5.txt' FOR UPDATE", false},
	{"VALUES ROW(1) ORDER BY a ORDER BY b", false},
	{"VALUES ROW(1) LIMIT 1 LIMIT 2", false},
	{"VALUES (1,2)", false},
	{"VALUES (1,2), (3,4)", false},
	{"VALUES ROW(1), (2)", false},
	// TABLE and VALUES as set-operation operands.
	{"TABLE t1 UNION TABLE t2", true},
	{"TABLE t1 UNION ALL TABLE t2", true},
	{"TABLE t1 UNION (TABLE t2)", true},
	{"TABLE t1 UNION SELECT * FROM t2", true},
	{"TABLE t1 UNION (SELECT c2 FROM t2) ORDER BY c1 LIMIT 1", true},
	{"TABLE t1 EXCEPT TABLE t2", true},
	{"TABLE t1 INTERSECT TABLE t2", true},
	{"TABLE t1 UNION TABLE t2 UNION TABLE t3", true},
	{"TABLE t1 UNION TABLE t2 INTERSECT TABLE t3", true},
	{"TABLE t1 UNION TABLE t2 ORDER BY c1", true},
	{"TABLE t1 UNION TABLE t2 FOR UPDATE", true},
	{"TABLE t UNION TABLE t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/a8.txt'", true},
	{"(TABLE t1) UNION SELECT c2 FROM t2 UNION (TABLE t3) ORDER BY c1 LIMIT 1", true},
	{"SELECT c1 FROM t1 UNION (TABLE t2) ORDER BY c1 LIMIT 1", true},
	{"VALUES ROW(1,2) UNION VALUES ROW(3,4)", true},
	{"VALUES ROW(1,2) UNION ALL VALUES ROW(3,4)", true},
	{"VALUES ROW(1,2) UNION (VALUES ROW(3,4))", true},
	{"VALUES ROW(1,2) UNION SELECT 1,2", true},
	{"VALUES ROW(1,2) UNION TABLE t1", true},
	{"TABLE t1 UNION VALUES ROW(1,2)", true},
	{"VALUES ROW(1,2) EXCEPT VALUES ROW(3,4)", true},
	{"VALUES ROW(1,2) INTERSECT VALUES ROW(3,4)", true},
	{"VALUES ROW(1) UNION VALUES ROW(2) ORDER BY 1 LIMIT 1", true},
	// TABLE and VALUES in every other query position.
	{"(TABLE t)", true},
	{"((TABLE t))", true},
	{"(TABLE t) ORDER BY a", true},
	{"(TABLE t) LIMIT 1", true},
	{"(TABLE t) UNION (TABLE t)", true},
	{"(TABLE t) FOR UPDATE", false},
	{"(VALUES ROW(1,2))", true},
	{"SELECT (TABLE t)", true},
	{"SELECT EXISTS (TABLE t)", true},
	{"SELECT * FROM (TABLE t) x", true},
	{"SELECT * FROM (VALUES ROW(1,2)) x", true},
	{"SELECT * FROM ((TABLE t)) x", true},
	{"SELECT * FROM ((VALUES ROW(1))) x", true},
	{"SELECT * FROM t1 WHERE c1 IN (TABLE t2)", true},
	{"SELECT * FROM t1 WHERE c1 IN (VALUES ROW(1))", true},
	{"SELECT * FROM t1 WHERE a = (TABLE t2)", true},
	{"SELECT 1 WHERE 1 IN ((TABLE t) FOR UPDATE)", false},
	{"UPDATE t SET a=1 WHERE b IN (TABLE t1)", true},
	{"DELETE FROM t WHERE b IN (TABLE t1)", true},
	{"WITH cte AS (TABLE t) TABLE cte", true},
	{"WITH RECURSIVE cte AS (TABLE t) TABLE cte", true},
	{"WITH cte AS (TABLE t) SELECT * FROM cte", true},
	{"WITH cte AS (SELECT 1) TABLE t", true},
	{"WITH cte AS (SELECT 1) VALUES ROW(1)", true},
	{"WITH cte AS (TABLE t) VALUES ROW(1)", true},
	{"CREATE VIEW v1 AS (TABLE t)", true},
	{"CREATE VIEW v2 AS TABLE t", true},
	{"CREATE VIEW v3 AS VALUES ROW(1,2)", true},
	{"CREATE VIEW v9 AS TABLE t1 UNION TABLE t2", true},
	{"CREATE TABLE tt AS TABLE t", true},
	{"CREATE TABLE tt2 AS VALUES ROW(1)", true},
	{"CREATE TABLE tt3 AS TABLE t UNION TABLE t", true},
	{"INSERT INTO t1 TABLE t2", true},
	{"REPLACE INTO t1 TABLE t2", true},
	{"INSERT INTO t1 TABLE t2 ORDER BY a", true},
	{"INSERT INTO t1 TABLE t2 LIMIT 1", true},
	{"INSERT INTO t1 TABLE t2 UNION TABLE t3", true},
	{"EXPLAIN TABLE t", true},
	{"EXPLAIN VALUES ROW(1)", true},
	{"EXPLAIN TABLE t1 UNION TABLE t2", true},
	// the locking-read clause.
	{"SELECT * FROM t FOR UPDATE", true},
	{"SELECT * FROM t FOR SHARE", true},
	{"SELECT * FROM t LOCK IN SHARE MODE", true},
	{"SELECT * FROM t FOR UPDATE NOWAIT", true},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5", true},
	{"SELECT * FROM t FOR UPDATE WAIT 0", true},
	{"SELECT * FROM t FOR UPDATE OF t", true},
	{"SELECT * FROM t FOR UPDATE OF t, t1", true},
	{"SELECT * FROM t FOR UPDATE OF t.t", true},
	{"SELECT * FROM t FOR UPDATE OF db.t", true},
	{"SELECT * FROM t FOR UPDATE OF `t`", true},
	{"SELECT * FROM t t1 FOR UPDATE OF t1", true},
	{"SELECT * FROM t FOR UPDATE OF t NOWAIT", true},
	{"SELECT * FROM t FOR UPDATE OF t SKIP LOCKED", true},
	{"SELECT * FROM t FOR UPDATE OF t WAIT 5", true},
	{"SELECT * FROM t FOR UPDATE OF t, t1 WAIT 5", true},
	{"SELECT * FROM t FOR SHARE NOWAIT", true},
	{"SELECT * FROM t FOR SHARE SKIP LOCKED", true},
	{"SELECT * FROM t FOR SHARE OF t", true},
	{"SELECT * FROM t FOR SHARE OF t, t1", true},
	{"SELECT * FROM t FOR SHARE OF t NOWAIT", true},
	{"SELECT * FROM t FOR SHARE OF t SKIP LOCKED", true},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE", true},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11", true},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE OF t WAIT 11", true},
	{"SELECT 1 UNION SELECT 2 FOR UPDATE", true},
	{"(SELECT 1) FOR UPDATE", false},
	{"SELECT * FROM t FOR SHARE WAIT 5", false},
	{"SELECT * FROM t FOR SHARE OF t WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE WAIT", false},
	{"SELECT * FROM t FOR UPDATE WAIT 'a'", false},
	{"SELECT * FROM t FOR UPDATE WAIT -1", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5.5", false},
	{"SELECT * FROM t FOR UPDATE WAIT 1e3", false},
	{"SELECT * FROM t FOR UPDATE WAIT 0x10", false},
	{"SELECT * FROM t FOR UPDATE WAIT 010", true},
	{"SELECT * FROM t FOR UPDATE WAIT 18446744073709551615", true},
	{"SELECT * FROM t FOR UPDATE WAIT 18446744073709551616", false},
	{"SELECT * FROM t FOR UPDATE WAIT 9223372036854775808", true},
	{"SELECT * FROM t FOR UPDATE WAIT 99999999999999999999999999", false},
	{"SELECT * FROM t FOR UPDATE WAIT @v", false},
	{"SELECT * FROM t FOR UPDATE WAIT ?", false},
	{"SELECT * FROM t FOR UPDATE WAIT (5)", false},
	{"SELECT * FROM t FOR UPDATE WAIT OF t", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5 OF t", false},
	{"SELECT * FROM t FOR UPDATE NOWAIT WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5 NOWAIT", false},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5 SKIP LOCKED", false},
	{"SELECT * FROM t FOR UPDATE WAIT 1 WAIT 2", false},
	{"SELECT * FROM t FOR UPDATE SKIP", false},
	{"SELECT * FROM t FOR UPDATE OF", false},
	{"SELECT * FROM t FOR UPDATE NOWAIT NOWAIT", false},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED SKIP LOCKED", false},
	{"SELECT * FROM t FOR UPDATE FOR UPDATE", false},
	{"SELECT * FROM t FOR UPDATE FOR SHARE", false},
	{"SELECT * FROM t FOR UPDATE LOCK IN SHARE MODE", false},
	{"SELECT * FROM t LOCK IN SHARE MODE FOR UPDATE", false},
	{"SELECT * FROM t LOCK IN SHARE MODE NOWAIT", false},
	{"SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED", false},
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT 5", false},
	// the INTO clause, which the locking clause comes before.
	{"SELECT 1 INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT * FROM t ORDER BY a INTO OUTFILE '/tmp/f4.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'", true},
	{"SELECT a FROM t WHERE a=1 INTO OUTFILE '/tmp/p2.txt'", true},
	{"SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'", true},
	{"SELECT a FROM t HAVING a>1 INTO OUTFILE '/tmp/p6.txt'", true},
	{"SELECT a FROM t LIMIT 1 INTO OUTFILE '/tmp/pc.txt'", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'", true},
	{"SELECT * FROM t FOR UPDATE OF t WAIT 5 INTO OUTFILE '/tmp/f6.txt'", true},
	{"SELECT 1 UNION SELECT 2 INTO OUTFILE '/tmp/pg.txt'", true},
	{"SELECT 1 UNION SELECT 2 ORDER BY 1 INTO OUTFILE '/tmp/ph.txt'", true},
	{"SELECT * FROM t INTO OUTFILE '/tmp/f5.txt' ORDER BY a", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/p7.txt' LIMIT 1", false},
	{"SELECT * FROM t INTO OUTFILE '/tmp/f2.txt' FOR UPDATE", false},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE", false},
	{"SELECT a FROM t ORDER BY a INTO OUTFILE '/tmp/pb.txt' LIMIT 1", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/pe.txt' INTO OUTFILE '/tmp/pf.txt'", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/p3.txt' WHERE a=1", false},
	{"SELECT a FROM t INTO DUMPFILE '/tmp/pd.txt'", false},
	{"SELECT a FROM t INTO @v", false},
	{"(SELECT 1) INTO OUTFILE '/tmp/pi.txt'", false},
	// the locking-clause keywords as names.
	{"SELECT * FROM t wait", true},
	{"SELECT wait FROM t", true},
	{"SELECT * FROM t AS wait", true},
	{"SELECT * FROM t wait WHERE wait.a = 1", true},
	{"DELETE FROM t wait WHERE wait.a = 1", true},
	{"UPDATE t wait SET wait.a = 1", true},
	{"SELECT 1 AS wait", true},
	{"SELECT * FROM t locked", true},
	{"SELECT * FROM t nowait", true},
	{"SELECT * FROM t AS locked WHERE locked.a=1", true},
	{"DELETE FROM t locked WHERE locked.a = 1", true},
	{"DELETE FROM t nowait WHERE nowait.a = 1", true},
	{"CREATE TABLE nowait (locked int, wait int, `row` int)", true},
}

// TestQueryPrimaryFormsOracle asserts omni agrees with the pinned TiDB on
// queryPrimaryCases (container test). The verdict comes from PREPARE, which
// parses the statement without running it: 1064 is the parser's rejection,
// and any other server error means the text parsed and failed later.
func TestQueryPrimaryFormsOracle(t *testing.T) {
	tc := startTiDB(t)
	assertTiDBVersion(t, tc)
	for i, qpc := range queryPrimaryCases {
		t.Run(qpc.sql, func(t *testing.T) {
			previous := "<start of cases>"
			if i > 0 {
				previous = queryPrimaryCases[i-1].sql
			}
			accepts, _, err := tc.syntaxVerdict(t, qpc.sql, previous)
			if err != nil {
				t.Fatalf("no verdict: %v", err)
			}
			if accepts != qpc.accept {
				t.Fatalf("TiDB %s accepts=%v, want %v for %q", tidbCorpusVersion, accepts, qpc.accept, qpc.sql)
			}
			_, perr := Parse(qpc.sql)
			if omniAccepts := perr == nil; omniAccepts != qpc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, TiDB %s accepts=%v for %q: %v",
					omniAccepts, tidbCorpusVersion, qpc.accept, qpc.sql, perr)
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
