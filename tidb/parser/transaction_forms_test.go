package parser

import (
	"slices"
	"testing"

	"github.com/bytebase/omni/tidb/ast"
)

// The transaction statements as TiDB reads them, which is not as MySQL does.
//
// This parser began as a copy of the MySQL one and read MySQL's transaction
// grammar: BEGIN WORK, COMMIT WORK and ROLLBACK WORK, a comma-separated list
// of START TRANSACTION characteristics, and every combination of AND [NO]
// CHAIN and [NO] RELEASE. TiDB has no WORK at all, takes one characteristic
// at most, and refuses AND CHAIN RELEASE; it adds BEGIN PESSIMISTIC and BEGIN
// OPTIMISTIC, WITH CAUSAL CONSISTENCY ONLY, and READ ONLY AS OF TIMESTAMP. And
// SAVEPOINT is an unreserved keyword there, so it is also a savepoint name:
// "ROLLBACK TO SAVEPOINT" rolls back to a savepoint called savepoint.
//
// Not recorded here, because this parser still disagrees and the change is not
// about it: a stored program body takes COMMIT and ROLLBACK and no other
// transaction statement (ProcedureStatementStmt), so TiDB answers 1064 for
// START TRANSACTION, SAVEPOINT and RELEASE SAVEPOINT inside one; START REPLICA,
// START SLAVE and START GROUP_REPLICATION are 1064 on TiDB; and this parser
// reads no '?' parameter marker, so it rejects "... AS OF TIMESTAMP ?", which
// PREPARE takes.
//
// Every case below was measured against TiDB v8.5.5 (pingcap/tidb:v8.5.5)
// through PREPARE on 2026-11-05. TestTransactionFormsMatchTiDB keeps it a
// recording: an image upgrade that changes an answer fails there rather than
// leaving a stale expectation behind.
//
// Ref: BeginTransactionStmt, CommitStmt, RollbackStmt,
// CompletionTypeWithinTransaction, SavepointStmt and ReleaseSavepointStmt in
// TiDB v8.5.5 pkg/parser/parser.y.

// transactionFormCase is a statement and whether TiDB's parser accepts it.
type transactionFormCase struct {
	sql    string
	accept bool
}

var transactionFormCases = []transactionFormCase{
	// --- BEGIN [PESSIMISTIC | OPTIMISTIC], and no WORK ---
	{"begin", true},
	{"begin pessimistic", true},
	{"begin optimistic", true},
	{"BEGIN PESSIMISTIC", true},
	{"Begin Optimistic", true},
	{"begin /* c */ pessimistic", true},
	{"begin pessimistic;", true},
	{"begin work", false},
	{"begin work;", false},
	{"begin pessimistic work", false},
	{"begin work pessimistic", false},
	{"begin optimistic pessimistic", false},
	{"begin pessimistic, optimistic", false},
	{"begin transaction", false},
	{"begin read only", false},
	{"begin with consistent snapshot", false},
	// The modes are words, not names.
	{"begin `pessimistic`", false},
	{"begin 'pessimistic'", false},

	// --- START TRANSACTION and one characteristic ---
	{"start transaction", true},
	{"start transaction read only", true},
	{"start transaction read write", true},
	{"start transaction with consistent snapshot", true},
	{"start transaction with causal consistency only", true},
	{"START TRANSACTION WITH CAUSAL CONSISTENCY ONLY", true},
	{"start transaction /*!40108 with consistent snapshot */", true},
	{"start", false},
	{"start work", false},
	{"start transaction work", false},
	{"start transaction pessimistic", false},
	{"start transaction optimistic", false},
	{"start transaction isolation level read committed", false},
	// A characteristic is whole.
	{"start transaction read", false},
	{"start transaction with", false},
	{"start transaction with consistent", false},
	{"start transaction with snapshot", false},
	{"start transaction with consistent snapshot only", false},
	{"start transaction with causal", false},
	{"start transaction with causal consistency", false},
	{"start transaction with causal only", false},
	{"start transaction with `causal` consistency only", false},
	{"start transaction with causal `consistency` only", false},
	// One at most: MySQL's list is not TiDB's.
	{"start transaction read only, with consistent snapshot", false},
	{"start transaction with consistent snapshot, read only", false},
	{"start transaction read only, read write", false},
	{"start transaction read write, read only", false},
	{"start transaction with consistent snapshot, read write", false},
	{"start transaction with consistent snapshot, with consistent snapshot", false},
	{"start transaction with consistent snapshot read only", false},
	{"start transaction read only read write", false},
	{"start transaction read only,", false},
	{"start transaction, read only", false},
	// READ ONLY AS OF TIMESTAMP, with any expression, and READ ONLY alone.
	{"start transaction read only as of timestamp now()", true},
	{"start transaction read only as of timestamp '2015-09-21 00:07:01'", true},
	{"start transaction read only as of timestamp now() - interval 3 second", true},
	{"start transaction read only as of timestamp now() + 1", true},
	{"start transaction read only as of timestamp now() and 1", true},
	{"start transaction read only as of timestamp (select now())", true},
	{"start transaction read only as of timestamp a", true},
	{"start transaction read only as of timestamp @a", true},
	{"start transaction read only as of timestamp tidb_bounded_staleness(now(), now())", true},
	{"start transaction read only as", false},
	{"start transaction read only as of", false},
	{"start transaction read only as of timestamp", false},
	{"start transaction read only as of time now()", false},
	{"start transaction read only as of timestamp now() for update", false},
	{"start transaction read only as of timestamp now() read write", false},
	{"start transaction read only as of timestamp now(), read write", false},
	{"start transaction read write as of timestamp now()", false},
	{"start transaction as of timestamp now()", false},
	{"start transaction with consistent snapshot as of timestamp now()", false},

	// --- COMMIT [AND [NO] CHAIN] [[NO] RELEASE], and no WORK ---
	{"commit", true},
	{"commit;", true},
	{"commit and chain", true},
	{"commit and no chain", true},
	{"commit release", true},
	{"commit no release", true},
	{"commit and chain no release", true},
	{"commit and no chain release", true},
	{"commit and no chain no release", true},
	{"commit work", false},
	{"commit work and chain", false},
	{"commit work release", false},
	{"commit transaction", false},
	{"commit to savepoint s", false},
	// AND CHAIN RELEASE is the one combination there is not.
	{"commit and chain release", false},
	// Each part is whole, and they come in order.
	{"commit and", false},
	{"commit and no", false},
	{"commit no", false},
	{"commit chain", false},
	{"commit no chain", false},
	{"commit and no chain no", false},
	{"commit and no chain and", false},
	{"commit and chain no release no", false},
	{"commit release no", false},
	{"commit release release", false},
	{"commit no release release", false},
	{"commit and chain and chain", false},
	{"commit no release and chain", false},

	// --- ROLLBACK, as COMMIT, or TO [SAVEPOINT] name ---
	{"rollback", true},
	{"rollback and chain", true},
	{"rollback and no chain", true},
	{"rollback release", true},
	{"rollback no release", true},
	{"rollback and chain no release", true},
	{"rollback and no chain release", true},
	{"rollback and no chain no release", true},
	{"rollback work", false},
	{"rollback transaction", false},
	{"rollback and chain release", false},
	{"rollback and chain no", false},
	{"rollback and", false},
	{"rollback no", false},
	{"rollback no release and no chain", false},
	{"rollback to savepoint s", true},
	{"rollback to s", true},
	{"rollback to `s`", true},
	{"rollback to work", true},
	{"rollback to savepoint work", true},
	{"rollback to begin", true},
	{"rollback to chain", true},
	{"rollback work to savepoint s", false},
	{"rollback work to s", false},
	{"rollback to", false},
	{"rollback savepoint s", false},
	{"rollback to 's'", false},
	{"rollback to savepoint 's'", false},
	{"rollback to release", false},
	{"rollback to and", false},
	{"rollback to savepoint s s", false},
	// ROLLBACK TO takes no completion type.
	{"rollback to savepoint s and chain", false},
	{"rollback to s release", false},
	// SAVEPOINT is the name when nothing follows it.
	{"rollback to savepoint", true},
	{"rollback to savepoint;", true},
	{"rollback to savepoint savepoint", true},
	{"rollback to savepoint `savepoint`", true},
	{"rollback to `savepoint`", true},
	{"rollback to savepoint savepoint savepoint", false},
	{"rollback to savepoint and chain", false},

	// --- SAVEPOINT name and RELEASE SAVEPOINT name ---
	{"savepoint s", true},
	{"savepoint `s`", true},
	{"savepoint work", true},
	{"savepoint savepoint", true},
	{"savepoint begin", true},
	{"savepoint commit", true},
	{"savepoint chain", true},
	{"savepoint pessimistic", true},
	{"savepoint", false},
	{"savepoint s s", false},
	{"savepoint 's'", false},
	{"savepoint release", false},
	{"release savepoint s", true},
	{"release savepoint `s`", true},
	{"release savepoint savepoint", true},
	{"release s", false},
	{"release work", false},
	{"release", false},
	{"release savepoint", false},
	{"release savepoint 's'", false},
	{"release savepoint s s", false},
	{"release savepoint savepoint savepoint", false},

	// --- COMMIT and ROLLBACK in a stored program body, which reads them as
	// --- the statement does ---
	{"create procedure p() begin commit; end", true},
	{"create procedure p() begin rollback; end", true},
	{"create procedure p() begin rollback to s; end", true},
	{"create procedure p() begin rollback to savepoint; end", true},
	{"create procedure p() begin commit work; end", false},
	{"create procedure p() begin rollback work; end", false},
	{"create procedure p() begin commit and chain release; end", false},
	// A label right after BEGIN is not a transaction mode: WORK is no
	// keyword in TiDB, so it can name one.
	{"create procedure p() begin work: while 1 do select 1; end while work; end", true},
}

func TestTransactionForms(t *testing.T) {
	for _, c := range transactionFormCases {
		t.Run(c.sql, func(t *testing.T) {
			_, err := Parse(c.sql)
			if accepts := err == nil; accepts != c.accept {
				t.Errorf("omni accepts=%v, want %v for %q: %v", accepts, c.accept, c.sql, err)
			}
		})
	}
}

// TestTransactionTree holds what each accepted form records, and that every
// node ends at its last token: a BEGIN mode, a characteristic, the AS OF
// expression (which ast.Walk reaches), the completion type, and the
// savepoint name, SAVEPOINT included when it is the name.
func TestTransactionTree(t *testing.T) {
	for _, c := range []struct {
		sql   string
		want  string   // NodeToString of the statement
		texts []string // nodeTexts
	}{
		{"begin", "{BEGIN :loc 0}", []string{"BeginStmt: begin"}},
		{"begin pessimistic", "{BEGIN :loc 0 :mode PESSIMISTIC}", []string{"BeginStmt: begin pessimistic"}},
		{"BEGIN OPTIMISTIC", "{BEGIN :loc 0 :mode OPTIMISTIC}", []string{"BeginStmt: BEGIN OPTIMISTIC"}},
		{"start transaction", "{BEGIN :loc 0}", []string{"BeginStmt: start transaction"}},
		{"start transaction read write", "{BEGIN :loc 0 :read_write true}", nil},
		{"start transaction with consistent snapshot", "{BEGIN :loc 0 :consistent_snapshot true}", nil},
		{
			"start transaction with causal consistency only",
			"{BEGIN :loc 0 :causal_consistency_only true}",
			[]string{"BeginStmt: start transaction with causal consistency only"},
		},
		{"start transaction read only", "{BEGIN :loc 0 :read_only true}", nil},
		{
			"start transaction read only as of timestamp now() - interval 3 second",
			"",
			[]string{
				"BeginStmt: start transaction read only as of timestamp now() - interval 3 second",
				"BinaryExpr: now() - interval 3 second",
				"FuncCallExpr: now()",
				"IntervalExpr: interval 3 second",
				"IntLit: 3",
			},
		},
		{"commit", "{COMMIT :loc 0}", []string{"CommitStmt: commit"}},
		{"commit and chain", "{COMMIT :loc 0 :chain true}", []string{"CommitStmt: commit and chain"}},
		{"commit and chain no release", "{COMMIT :loc 0 :chain true}", []string{"CommitStmt: commit and chain no release"}},
		{"commit and no chain release", "{COMMIT :loc 0 :release true}", []string{"CommitStmt: commit and no chain release"}},
		{"rollback release", "{ROLLBACK :loc 0 :release true}", []string{"RollbackStmt: rollback release"}},
		{"rollback to s", "{ROLLBACK :loc 0 :savepoint s}", []string{"RollbackStmt: rollback to s"}},
		{"rollback to savepoint s", "{ROLLBACK :loc 0 :savepoint s}", []string{"RollbackStmt: rollback to savepoint s"}},
		{"rollback to savepoint", "{ROLLBACK :loc 0 :savepoint savepoint}", []string{"RollbackStmt: rollback to savepoint"}},
		{"rollback to savepoint savepoint", "{ROLLBACK :loc 0 :savepoint savepoint}", nil},
		{"release savepoint s", "{RELEASE_SAVEPOINT :loc 0 :name s}", []string{"SavepointStmt: release savepoint s"}},
	} {
		t.Run(c.sql, func(t *testing.T) {
			list, err := Parse(c.sql)
			if err != nil {
				t.Fatalf("Parse(%q): %v", c.sql, err)
			}
			if len(list.Items) != 1 {
				t.Fatalf("Parse(%q) = %d statements, want 1", c.sql, len(list.Items))
			}
			if c.want != "" {
				if got := ast.NodeToString(list.Items[0]); got != c.want {
					t.Errorf("Parse(%q) = %s, want %s", c.sql, got, c.want)
				}
			}
			if c.texts != nil {
				if got := nodeTexts(t, c.sql); !slices.Equal(got, c.texts) {
					t.Errorf("nodeTexts(%q) = %q, want %q", c.sql, got, c.texts)
				}
			}
		})
	}
}
