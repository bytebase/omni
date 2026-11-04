package parser

import "testing"

// The locking-read clause and the INTO after it, held to a live MariaDB
// 11.8.8. One case per form this parser reads and per malformed variant
// around it: a clause the parser swallows where MariaDB rejects it leaves no
// trace in a suite that only checks what parses.
//
// The list is close to the same text in the sibling engine trees. What differs
// is this column, which is where the forks part. MariaDB's locking reads are
// its own shape: it has no FOR SHARE and no FOR UPDATE ... OF (both MySQL's),
// it takes a WAIT n modifier that MySQL has not, after LOCK IN SHARE MODE as
// well as after FOR UPDATE, and its WAIT takes any numeric literal where
// TiDB's takes an integer.
//
// Measured on 2026-11-04 against mariadb:11.8.8 via PREPARE.

// lockingReadCase is a statement and whether the engine's parser accepts it.
type lockingReadCase struct {
	sql    string
	accept bool
}

var lockingReadCases = []lockingReadCase{
	// The locking-read clause. MariaDB has FOR UPDATE and LOCK IN SHARE MODE, each with a wait modifier, and neither FOR SHARE nor OF.
	{"SELECT * FROM t FOR UPDATE", true},
	{"SELECT * FROM t LOCK IN SHARE MODE", true},
	{"SELECT * FROM t FOR UPDATE NOWAIT", true},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5", true},
	{"SELECT * FROM t FOR UPDATE WAIT 0", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5.5", true},
	{"SELECT * FROM t FOR UPDATE WAIT 0.5", true},
	{"SELECT * FROM t FOR UPDATE WAIT 1e3", true},
	{"SELECT * FROM t FOR UPDATE WAIT 010", true},
	{"SELECT * FROM t FOR UPDATE WAIT 0x10", true},
	{"SELECT * FROM t FOR UPDATE WAIT 0b10", false},
	{"SELECT * FROM t FOR UPDATE WAIT X'10'", false},
	{"SELECT * FROM t FOR UPDATE WAIT 18446744073709551615", true},
	{"SELECT * FROM t FOR UPDATE WAIT 99999999999999999999999999", true},
	{"SELECT * FROM t FOR UPDATE WAIT", false},
	{"SELECT * FROM t FOR UPDATE WAIT 'a'", false},
	{"SELECT * FROM t FOR UPDATE WAIT -5", false},
	{"SELECT * FROM t FOR UPDATE WAIT @v", false},
	{"SELECT * FROM t FOR UPDATE WAIT ?", false},
	{"SELECT * FROM t FOR UPDATE WAIT (5)", false},
	{"SELECT * FROM t FOR UPDATE WAIT 1 WAIT 2", false},
	{"SELECT * FROM t FOR UPDATE NOWAIT WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5 NOWAIT", false},
	{"SELECT * FROM t FOR UPDATE SKIP LOCKED WAIT 5", false},
	{"SELECT * FROM t FOR UPDATE WAIT 5 SKIP LOCKED", false},
	{"SELECT * FROM t FOR UPDATE SKIP", false},
	{"SELECT * FROM t FOR UPDATE NOWAIT NOWAIT", false},
	{"SELECT * FROM t FOR SHARE", false},
	{"SELECT * FROM t FOR SHARE NOWAIT", false},
	{"SELECT * FROM t FOR SHARE OF t", false},
	{"SELECT * FROM t FOR UPDATE OF t", false},
	{"SELECT * FROM t FOR UPDATE OF t NOWAIT", false},
	{"SELECT * FROM t FOR UPDATE OF", false},
	{"SELECT * FROM t LOCK IN SHARE MODE NOWAIT", true},
	{"SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED", true},
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT 5", true},
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT 5.5", true},
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT 0x10", true},
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT", false},
	{"SELECT * FROM t LOCK IN SHARE MODE NOWAIT NOWAIT", false},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE", true},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11", true},
	{"SELECT 1 UNION SELECT 2 FOR UPDATE", true},

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
	{"SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'", true},

	// The locking-clause keywords as names.
	{"SELECT * FROM t wait", true},
	{"SELECT * FROM t AS wait", true},
	{"SELECT * FROM t wait WHERE wait.a = 1", true},
	{"DELETE FROM t wait WHERE wait.a = 1", true},
	{"DELETE FROM t locked WHERE locked.a = 1", true},
	{"DELETE FROM t nowait WHERE nowait.a = 1", true},
	{"UPDATE t wait SET wait.a = 1", true},
	{"SELECT 1 AS wait", true},
	{"CREATE TABLE nowait (locked int, wait int, `row` int)", true},

	// MariaDB has no TABLE statement and spells its table value constructor without ROW.
	{"TABLE t", false},
	{"TABLE t ORDER BY a", false},
	{"VALUES ROW(1,2)", false},
	{"VALUES (1,2)", true},
	{"VALUES (1,2), (3,4)", true},
}

// TestLockingReadFormsOracle asserts omni agrees with a live MariaDB 11.8.8 on
// lockingReadCases (container test).
func TestLockingReadFormsOracle(t *testing.T) {
	o := startMariaDB(t)
	for _, setup := range []string{"CREATE TABLE IF NOT EXISTS t (a INT, b INT)"} {
		if _, err := o.db.ExecContext(o.ctx, setup); err != nil {
			t.Fatalf("%s: %v", setup, err)
		}
	}
	for _, lrc := range lockingReadCases {
		t.Run(lrc.sql, func(t *testing.T) {
			if accepts := o.prepareAccepts(t, lrc.sql); accepts != lrc.accept {
				t.Fatalf("MariaDB 11.8.8 accepts=%v, want %v for %q", accepts, lrc.accept, lrc.sql)
			}
			_, perr := Parse(lrc.sql)
			if omniAccepts := perr == nil; omniAccepts != lrc.accept {
				t.Errorf("DIVERGENCE omni accepts=%v, MariaDB 11.8.8 accepts=%v for %q: %v",
					omniAccepts, lrc.accept, lrc.sql, perr)
			}
		})
	}
}

// TestLockingReadForms holds the parser to the recorded verdicts and needs no
// server, so it is the gate that runs without Docker.
func TestLockingReadForms(t *testing.T) {
	for _, lrc := range lockingReadCases {
		t.Run(lrc.sql, func(t *testing.T) {
			_, err := Parse(lrc.sql)
			if got := err == nil; got != lrc.accept {
				t.Errorf("parser accepts=%v, want %v for %q: %v", got, lrc.accept, lrc.sql, err)
			}
		})
	}
}
