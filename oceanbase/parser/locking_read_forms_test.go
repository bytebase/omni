package parser

import "testing"

// The locking-read clause and the INTO after it, held to the verdicts
// OceanBase CE 4.3.5 gave. One case per form this parser reads and per
// malformed variant around it: a clause the parser swallows where the server
// rejects it leaves no trace in a suite that only checks what parses.
//
// The list is close to the same text in the sibling engine trees. What differs
// is this column, which is where the forks part. OceanBase takes a WAIT n
// modifier MySQL has not — any numeric literal, as MariaDB's does, where
// TiDB's takes an integer — and it takes INTO only last, after the locking
// clause, where MySQL and MariaDB take it before one too.
//
// Measured by hand on 2026-11-04 against oceanbase/oceanbase-ce:4.3.5-lts
// (server 5.7.25-OceanBase_CE-v4.3.5.6) via PREPARE. This package's oracle
// runs against mysql:8.0, which answers differently for every WAIT row and
// for INTO before the locking clause, so these rows are recorded rather than
// re-measured in CI; rerun the probe by hand before changing one.
//
// OceanBase CE 4.3.5 has no FOR SHARE and no FOR UPDATE ... OF, and
// query_primary_forms_test.go holds those refusals to the server.

// lockingReadCase is a statement and whether the engine's parser accepts it.
type lockingReadCase struct {
	sql    string
	accept bool
}

var lockingReadCases = []lockingReadCase{
	// The locking-read clause. OceanBase takes a WAIT n modifier MySQL has not.
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
	{"SELECT * FROM t FOR UPDATE WAIT 0x10", false},
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
	{"SELECT * FROM t LOCK IN SHARE MODE WAIT 5", false},
	{"SELECT * FROM t LOCK IN SHARE MODE NOWAIT", false},
	{"SELECT * FROM t LOCK IN SHARE MODE SKIP LOCKED", false},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE", true},
	{"SELECT * FROM t LIMIT 1 FOR UPDATE WAIT 11", true},
	{"SELECT 1 UNION SELECT 2 FOR UPDATE", true},

	// The INTO clause, which OceanBase takes only last, after the locking clause.
	{"SELECT 1 INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT 1 FOR UPDATE INTO OUTFILE '/tmp/f3.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p8.txt'", true},
	{"SELECT a FROM t ORDER BY a LIMIT 1 FOR UPDATE INTO OUTFILE '/tmp/pa.txt'", true},
	{"SELECT a FROM t GROUP BY a INTO OUTFILE '/tmp/p4.txt'", true},
	{"SELECT * FROM t FOR UPDATE WAIT 5 INTO OUTFILE '/tmp/f1.txt'", true},
	{"SELECT * FROM t INTO OUTFILE '/tmp/f5.txt' ORDER BY a", false},
	{"SELECT a FROM t INTO OUTFILE '/tmp/p7.txt' LIMIT 1", false},
	{"SELECT a FROM t ORDER BY a LIMIT 1 INTO OUTFILE '/tmp/p9.txt' FOR UPDATE", false},
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

// TestLockingReadForms holds the parser to the recorded verdicts. It needs no
// server: the verdicts are OceanBase's, and this package's container oracle is
// MySQL.
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
