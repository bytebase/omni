package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Tests for the `parser-utility` node's transaction / batch statements (§2.9):
// BEGIN / START TRANSACTION / COMMIT / ROLLBACK and START / RUN / ABORT BATCH.
//
// CORRECTNESS BASIS (correctness-protocol.md). Two oracles:
//   - the live Cloud Spanner emulator (transaction_oracle differential in
//     utility_oracle_test.go) — it PARSES every one of these (accepts the
//     leading form, then feature-rejects "Statement not supported: …"), so the
//     leading-form ACCEPT is authoritative.
//   - the canonical ZetaSQL corpus (transaction.sql / batch.sql) — the breadth
//     oracle for the precise grammar.
//
// The emulator's recognizer for these keywords swallows arbitrary trailing
// tokens (it accepts `COMMIT WORK`, `START BATCH a b`, `COMMIT garbage`), so the
// PRECISE trailing grammar is NON-authoritative on Spanner and follows the
// ZetaSQL .g4 (the grammar bytebase consumes): `commit_statement: COMMIT
// TRANSACTION?` rejects `COMMIT WORK`; `start_batch_statement: START BATCH id?`
// rejects a second word. These rejects are proven against the .g4, recorded in
// the divergence ledger, and exercised in the *_Rejects tests below — NOT diffed
// against the emulator (a Spanner accept there would be a false divergence).

func parseTransaction(t *testing.T, sql string) *ast.TransactionStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	tx, ok := n.(*ast.TransactionStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.TransactionStmt", sql, n)
	}
	return tx
}

func parseBatch(t *testing.T, sql string) *ast.BatchStmt {
	t.Helper()
	n := parseOneStmt(t, sql)
	b, ok := n.(*ast.BatchStmt)
	if !ok {
		t.Fatalf("Parse(%q): got %T, want *ast.BatchStmt", sql, n)
	}
	return b
}

func TestTransaction_Begin(t *testing.T) {
	t.Run("bare BEGIN", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN")
		if tx.Kind != ast.TransactionBegin {
			t.Errorf("Kind = %v, want BEGIN", tx.Kind)
		}
		if tx.Transaction {
			t.Error("Transaction = true, want false (no TRANSACTION keyword)")
		}
		if tx.Modes != nil {
			t.Errorf("Modes = %v, want nil", tx.Modes)
		}
	})

	t.Run("BEGIN TRANSACTION", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION")
		if tx.Kind != ast.TransactionBegin {
			t.Errorf("Kind = %v, want BEGIN", tx.Kind)
		}
		if !tx.Transaction {
			t.Error("Transaction = false, want true")
		}
	})

	t.Run("START TRANSACTION", func(t *testing.T) {
		tx := parseTransaction(t, "START TRANSACTION")
		if tx.Kind != ast.TransactionStart {
			t.Errorf("Kind = %v, want START TRANSACTION", tx.Kind)
		}
		if !tx.Transaction {
			t.Error("Transaction = false, want true (START TRANSACTION implies it)")
		}
	})

	t.Run("READ ONLY", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION READ ONLY")
		if len(tx.Modes) != 1 || tx.Modes[0].Kind != ast.TransactionModeReadOnly {
			t.Fatalf("Modes = %+v, want one READ ONLY", tx.Modes)
		}
	})

	t.Run("READ WRITE", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION READ WRITE")
		if len(tx.Modes) != 1 || tx.Modes[0].Kind != ast.TransactionModeReadWrite {
			t.Fatalf("Modes = %+v, want one READ WRITE", tx.Modes)
		}
	})

	t.Run("ISOLATION LEVEL one word", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION ISOLATION LEVEL Serializable")
		if len(tx.Modes) != 1 {
			t.Fatalf("Modes = %+v, want one", tx.Modes)
		}
		m := tx.Modes[0]
		if m.Kind != ast.TransactionModeIsolationLevel {
			t.Errorf("Kind = %v, want ISOLATION LEVEL", m.Kind)
		}
		if len(m.Levels) != 1 || m.Levels[0] != "Serializable" {
			t.Errorf("Levels = %v, want [Serializable] (case preserved)", m.Levels)
		}
	})

	t.Run("ISOLATION LEVEL two words", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION ISOLATION LEVEL READ COMMITED")
		m := tx.Modes[0]
		if len(m.Levels) != 2 || m.Levels[0] != "READ" || m.Levels[1] != "COMMITED" {
			t.Errorf("Levels = %v, want [READ COMMITED]", m.Levels)
		}
	})

	t.Run("mode list", func(t *testing.T) {
		tx := parseTransaction(t, "BEGIN TRANSACTION READ WRITE, ISOLATION LEVEL READ COMMITED")
		if len(tx.Modes) != 2 {
			t.Fatalf("Modes = %d, want 2", len(tx.Modes))
		}
		if tx.Modes[0].Kind != ast.TransactionModeReadWrite {
			t.Errorf("Modes[0] = %v, want READ WRITE", tx.Modes[0].Kind)
		}
		if tx.Modes[1].Kind != ast.TransactionModeIsolationLevel {
			t.Errorf("Modes[1] = %v, want ISOLATION LEVEL", tx.Modes[1].Kind)
		}
	})

	t.Run("START TRANSACTION with modes", func(t *testing.T) {
		tx := parseTransaction(t, "START TRANSACTION READ ONLY, ISOLATION LEVEL SERIALIZABLE")
		if tx.Kind != ast.TransactionStart {
			t.Errorf("Kind = %v, want START TRANSACTION", tx.Kind)
		}
		if len(tx.Modes) != 2 {
			t.Fatalf("Modes = %d, want 2", len(tx.Modes))
		}
	})
}

func TestTransaction_CommitRollback(t *testing.T) {
	t.Run("COMMIT", func(t *testing.T) {
		tx := parseTransaction(t, "COMMIT")
		if tx.Kind != ast.TransactionCommit {
			t.Errorf("Kind = %v, want COMMIT", tx.Kind)
		}
		if tx.Transaction {
			t.Error("Transaction = true, want false")
		}
	})

	t.Run("COMMIT TRANSACTION", func(t *testing.T) {
		tx := parseTransaction(t, "COMMIT TRANSACTION")
		if tx.Kind != ast.TransactionCommit || !tx.Transaction {
			t.Errorf("got Kind=%v Transaction=%v, want COMMIT + TRANSACTION", tx.Kind, tx.Transaction)
		}
	})

	t.Run("ROLLBACK", func(t *testing.T) {
		tx := parseTransaction(t, "ROLLBACK")
		if tx.Kind != ast.TransactionRollback {
			t.Errorf("Kind = %v, want ROLLBACK", tx.Kind)
		}
	})

	t.Run("ROLLBACK TRANSACTION", func(t *testing.T) {
		tx := parseTransaction(t, "ROLLBACK TRANSACTION")
		if tx.Kind != ast.TransactionRollback || !tx.Transaction {
			t.Errorf("got Kind=%v Transaction=%v, want ROLLBACK + TRANSACTION", tx.Kind, tx.Transaction)
		}
	})
}

func TestBatch(t *testing.T) {
	t.Run("START BATCH", func(t *testing.T) {
		b := parseBatch(t, "START BATCH")
		if b.Kind != ast.BatchStart {
			t.Errorf("Kind = %v, want START BATCH", b.Kind)
		}
		if b.Name != "" {
			t.Errorf("Name = %q, want empty", b.Name)
		}
	})

	t.Run("START BATCH with type (case preserved)", func(t *testing.T) {
		b := parseBatch(t, "START BATCH dDL")
		if b.Kind != ast.BatchStart || b.Name != "dDL" {
			t.Errorf("got Kind=%v Name=%q, want START BATCH dDL", b.Kind, b.Name)
		}
	})

	t.Run("RUN BATCH", func(t *testing.T) {
		b := parseBatch(t, "RUN BATCH")
		if b.Kind != ast.BatchRun {
			t.Errorf("Kind = %v, want RUN BATCH", b.Kind)
		}
	})

	t.Run("ABORT BATCH", func(t *testing.T) {
		b := parseBatch(t, "ABORT BATCH")
		if b.Kind != ast.BatchAbort {
			t.Errorf("Kind = %v, want ABORT BATCH", b.Kind)
		}
	})
}

// TestTransaction_Rejects covers the trailing-grammar rejects that follow the
// ZetaSQL .g4 (NOT the Spanner emulator, which over-accepts them — see header).
func TestTransaction_Rejects(t *testing.T) {
	rejects := []string{
		// commit_statement: COMMIT TRANSACTION? — no WORK/CHAIN/other word.
		"COMMIT WORK",
		"COMMIT FOO",
		"COMMIT 1",
		"COMMIT TRANSACTION extra",
		// rollback_statement: ROLLBACK TRANSACTION?.
		"ROLLBACK WORK",
		"ROLLBACK TO sp", // SAVEPOINT not in this grammar
		// transaction_mode requires READ ONLY/WRITE or ISOLATION LEVEL. The
		// follower after BEGIN here is TRANSACTION (a TCL follower), so these are
		// genuinely begin_statement inputs whose mode list is malformed. (A
		// `BEGIN <non-TCL-word>` is instead a BEGIN…END block opener owned by
		// parser-scripting — NOT a transaction reject — so it is not tested here.)
		"BEGIN TRANSACTION READ", // READ without ONLY/WRITE
		"BEGIN TRANSACTION ISOLATION foo",
		"BEGIN TRANSACTION ISOLATION LEVEL", // missing level identifier
		"BEGIN TRANSACTION READ ONLY,",      // trailing comma
		"BEGIN READ",                        // READ follower (TCL) but incomplete mode
		// start_batch_statement: START BATCH identifier? — at most one word.
		"START BATCH ddl extra",
		// run_batch_statement: RUN BATCH (no trailer); abort_batch_statement likewise.
		"RUN BATCH foo",
		"ABORT BATCH foo",
		"RUN", // RUN must be followed by BATCH
		"START BATCH ,",
	}
	for _, sql := range rejects {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			assertReject(t, sql)
		})
	}
}

// TestTransaction_CorpusAccepts parses the transaction / batch statements of the
// canonical ZetaSQL parser testdata (transaction.sql / batch.sql, lifted inline,
// EXCLUDING the SET TRANSACTION lines — set_statement is a separate node) and
// asserts each parses to the right node kind. The breadth/completeness oracle
// (correctness-protocol.md): the legacy ANTLR grammar bytebase consumes is a
// hand-port of ZetaSQL, so its own testdata is authoritative for the precise
// grammar (the Spanner emulator over-accepts the trailing tokens here).
func TestTransaction_CorpusAccepts(t *testing.T) {
	tcl := []string{
		// transaction.sql (begin/commit/rollback forms only).
		"BEGIN TRANSACTION",
		"BEGIN TRANSACTION READ ONLY",
		"BEGIN TRANSACTION ISOLATION LEVEL read uncommited",
		"BEGIN TRANSACTION READ WRITE, ISOLATION LEVEL READ COMMITED",
		"BEGIN TRANSACTION ISOLATION LEVEL READ repeatable",
		"BEGIN TRANSACTION ISOLATION LEVEL Serializable",
		"BEGIN TRANSACTION ISOLATION LEVEL FOO bar",
		"COMMIT",
		"ROLLBACK",
	}
	for _, sql := range tcl {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			if _, ok := parseOneStmt(t, sql).(*ast.TransactionStmt); !ok {
				t.Errorf("Parse(%q): want *ast.TransactionStmt", sql)
			}
		})
	}
	batch := []string{
		// batch.sql
		"START BATCH",
		"START BATCH dDL",
		"RUN BATCH",
		"ABORT BATCH",
	}
	for _, sql := range batch {
		sql := sql
		t.Run(sql, func(t *testing.T) {
			if _, ok := parseOneStmt(t, sql).(*ast.BatchStmt); !ok {
				t.Errorf("Parse(%q): want *ast.BatchStmt", sql)
			}
		})
	}
}
