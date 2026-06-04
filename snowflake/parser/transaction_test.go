package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bytebase/omni/snowflake/ast"
)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func testParseBeginStmt(input string) (*ast.BeginStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.BeginStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a BeginStmt"})
	}
	return stmt, result.Errors
}

func testParseCommitStmt(input string) (*ast.CommitStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.CommitStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a CommitStmt"})
	}
	return stmt, result.Errors
}

func testParseRollbackStmt(input string) (*ast.RollbackStmt, []ParseError) {
	result := ParseBestEffort(input)
	if len(result.File.Stmts) == 0 {
		return nil, result.Errors
	}
	stmt, ok := result.File.Stmts[0].(*ast.RollbackStmt)
	if !ok {
		return nil, append(result.Errors, ParseError{Msg: "not a RollbackStmt"})
	}
	return stmt, result.Errors
}

// ---------------------------------------------------------------------------
// BEGIN tests
//   Syntax (docs + legacy .g4):
//     BEGIN [ { WORK | TRANSACTION } ] [ NAME <name> ]
//     START TRANSACTION [ NAME <name> ]
// ---------------------------------------------------------------------------

func TestBegin_Bare(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginBare {
		t.Errorf("kind = %v, want BeginBare", stmt.Kind)
	}
	if !stmt.Name.IsEmpty() {
		t.Errorf("expected no NAME, got %q", stmt.Name.Normalize())
	}
}

func TestBegin_Work(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN WORK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginWork {
		t.Errorf("kind = %v, want BeginWork", stmt.Kind)
	}
	if !stmt.Name.IsEmpty() {
		t.Errorf("expected no NAME, got %q", stmt.Name.Normalize())
	}
}

func TestBegin_Transaction(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN TRANSACTION")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginTransaction {
		t.Errorf("kind = %v, want BeginTransaction", stmt.Kind)
	}
}

func TestBegin_Name(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN NAME T1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginBare {
		t.Errorf("kind = %v, want BeginBare", stmt.Kind)
	}
	if stmt.Name.Normalize() != "T1" {
		t.Errorf("name = %q, want T1", stmt.Name.Normalize())
	}
}

func TestBegin_WorkName(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN WORK NAME txn1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginWork {
		t.Errorf("kind = %v, want BeginWork", stmt.Kind)
	}
	if stmt.Name.Normalize() != "TXN1" {
		t.Errorf("name = %q, want TXN1", stmt.Name.Normalize())
	}
}

func TestBegin_TransactionName(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN TRANSACTION NAME my_txn")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginTransaction {
		t.Errorf("kind = %v, want BeginTransaction", stmt.Kind)
	}
	if stmt.Name.Normalize() != "MY_TXN" {
		t.Errorf("name = %q, want MY_TXN", stmt.Name.Normalize())
	}
}

func TestBegin_QuotedName(t *testing.T) {
	stmt, errs := testParseBeginStmt(`BEGIN NAME "My Txn"`)
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Name.Quoted {
		t.Error("expected quoted name")
	}
	// Quoted identifiers preserve case.
	if stmt.Name.Normalize() != "My Txn" {
		t.Errorf("name = %q, want %q", stmt.Name.Normalize(), "My Txn")
	}
}

func TestBegin_StartTransaction(t *testing.T) {
	stmt, errs := testParseBeginStmt("START TRANSACTION")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginStartTransaction {
		t.Errorf("kind = %v, want BeginStartTransaction", stmt.Kind)
	}
	if !stmt.Name.IsEmpty() {
		t.Errorf("expected no NAME, got %q", stmt.Name.Normalize())
	}
}

func TestBegin_StartTransactionName(t *testing.T) {
	stmt, errs := testParseBeginStmt("START TRANSACTION NAME T2")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginStartTransaction {
		t.Errorf("kind = %v, want BeginStartTransaction", stmt.Kind)
	}
	if stmt.Name.Normalize() != "T2" {
		t.Errorf("name = %q, want T2", stmt.Name.Normalize())
	}
}

func TestBegin_Lowercase(t *testing.T) {
	stmt, errs := testParseBeginStmt("begin transaction name t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Kind != ast.BeginTransaction {
		t.Errorf("kind = %v, want BeginTransaction", stmt.Kind)
	}
	if stmt.Name.Normalize() != "T1" {
		t.Errorf("name = %q, want T1", stmt.Name.Normalize())
	}
}

func TestBegin_Loc(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN TRANSACTION NAME t1")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	// End is the exclusive byte offset just past the final token ("t1").
	if stmt.Loc.End != len("BEGIN TRANSACTION NAME t1") {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len("BEGIN TRANSACTION NAME t1"))
	}
}

// Trailing-junk leniency: like every other omni snowflake statement parser
// (verified against parseDropStmt / parseTruncateStmt / parseUseStmt, which all
// accept "<stmt> junk" with zero errors), parseBeginStmt consumes exactly the
// grammar it owns and leaves any trailing tokens unconsumed without raising an
// error. Real Snowflake rejects "BEGIN FOO" as a syntax error; omni's
// engine-wide convention is to trust the F3 splitter and not enforce
// trailing-token strictness per node. This is a pre-existing engine-wide
// divergence (recorded in the divergence ledger), not specific to TCL. The
// statement still parses to the correct BeginStmt with the right Kind.
func TestBegin_TrailingJunkLenient(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN FOO")
	if len(errs) != 0 {
		t.Fatalf("BEGIN FOO: expected lenient parse (engine-wide convention), got errors: %v", errs)
	}
	if stmt.Kind != ast.BeginBare {
		t.Errorf("kind = %v, want BeginBare", stmt.Kind)
	}
	if !stmt.Name.IsEmpty() {
		t.Errorf("expected no NAME (FOO is trailing junk, not consumed), got %q", stmt.Name.Normalize())
	}
}

// Trailing junk after a complete BEGIN WORK: TRANSACTION is left unconsumed
// (WORK and TRANSACTION are mutually exclusive; at most one modifier is part of
// the grammar). Lenient parse per engine-wide convention.
func TestBegin_WorkTransactionLenient(t *testing.T) {
	stmt, errs := testParseBeginStmt("BEGIN WORK TRANSACTION")
	if len(errs) != 0 {
		t.Fatalf("BEGIN WORK TRANSACTION: expected lenient parse, got errors: %v", errs)
	}
	if stmt.Kind != ast.BeginWork {
		t.Errorf("kind = %v, want BeginWork (TRANSACTION is trailing junk)", stmt.Kind)
	}
}

// Negative (real grammar error this node owns): NAME keyword with no following
// identifier. parseIdent fails, so the whole statement fails.
func TestBegin_NegativeNameMissingIdent(t *testing.T) {
	_, errs := testParseBeginStmt("BEGIN NAME")
	if len(errs) == 0 {
		t.Fatal("expected error for BEGIN NAME (missing identifier), got none")
	}
}

// Negative (real grammar error): START without TRANSACTION is not a valid
// transaction opener — expect(kwTRANSACTION) fails.
func TestBegin_NegativeStartWithoutTransaction(t *testing.T) {
	_, errs := testParseBeginStmt("START")
	if len(errs) == 0 {
		t.Fatal("expected error for bare START, got none")
	}
}

// Negative (real grammar error): START WORK is invalid — only START TRANSACTION
// is documented; expect(kwTRANSACTION) fails on WORK.
func TestBegin_NegativeStartWork(t *testing.T) {
	_, errs := testParseBeginStmt("START WORK")
	if len(errs) == 0 {
		t.Fatal("expected error for START WORK, got none")
	}
}

// Negative (real grammar error): START TRANSACTION NAME with no following
// identifier — the NAME clause inside the START-form requires an identifier.
func TestBegin_NegativeStartTransactionNameMissingIdent(t *testing.T) {
	_, errs := testParseBeginStmt("START TRANSACTION NAME")
	if len(errs) == 0 {
		t.Fatal("expected error for START TRANSACTION NAME (missing identifier), got none")
	}
}

// ---------------------------------------------------------------------------
// COMMIT tests
//   Syntax: COMMIT [ WORK ]
// ---------------------------------------------------------------------------

func TestCommit_Bare(t *testing.T) {
	stmt, errs := testParseCommitStmt("COMMIT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false")
	}
}

func TestCommit_Work(t *testing.T) {
	stmt, errs := testParseCommitStmt("COMMIT WORK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Work {
		t.Error("expected Work=true")
	}
}

func TestCommit_Lowercase(t *testing.T) {
	stmt, errs := testParseCommitStmt("commit work")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Work {
		t.Error("expected Work=true")
	}
}

func TestCommit_Loc(t *testing.T) {
	stmt, errs := testParseCommitStmt("COMMIT WORK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len("COMMIT WORK") {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len("COMMIT WORK"))
	}
}

func TestCommit_BareLoc(t *testing.T) {
	stmt, errs := testParseCommitStmt("COMMIT")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.End != len("COMMIT") {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len("COMMIT"))
	}
}

// Trailing-junk leniency (engine-wide convention): a junk token after
// COMMIT [WORK] is left unconsumed without error. Real Snowflake rejects
// "COMMIT FOO"; omni does not enforce trailing-token strictness per node.
func TestCommit_TrailingJunkLenient(t *testing.T) {
	stmt, errs := testParseCommitStmt("COMMIT FOO")
	if len(errs) != 0 {
		t.Fatalf("COMMIT FOO: expected lenient parse, got errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false (FOO is trailing junk, not WORK)")
	}
}

// ---------------------------------------------------------------------------
// ROLLBACK tests
//   Syntax: ROLLBACK [ WORK ]
// ---------------------------------------------------------------------------

func TestRollback_Bare(t *testing.T) {
	stmt, errs := testParseRollbackStmt("ROLLBACK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false")
	}
}

func TestRollback_Work(t *testing.T) {
	stmt, errs := testParseRollbackStmt("ROLLBACK WORK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !stmt.Work {
		t.Error("expected Work=true")
	}
}

func TestRollback_Lowercase(t *testing.T) {
	stmt, errs := testParseRollbackStmt("rollback")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false")
	}
}

func TestRollback_Loc(t *testing.T) {
	stmt, errs := testParseRollbackStmt("ROLLBACK WORK")
	if len(errs) > 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if stmt.Loc.Start != 0 {
		t.Errorf("Loc.Start = %d, want 0", stmt.Loc.Start)
	}
	if stmt.Loc.End != len("ROLLBACK WORK") {
		t.Errorf("Loc.End = %d, want %d", stmt.Loc.End, len("ROLLBACK WORK"))
	}
}

// ROLLBACK TO [SAVEPOINT] is NOT supported by Snowflake — SAVEPOINT appears in
// neither the docs nor the legacy .g4, and WORK is the only modifier ROLLBACK
// accepts. The grammar therefore does NOT consume "TO SAVEPOINT sp1"; per the
// engine-wide trailing-token convention those tokens are left unconsumed
// without error, and only the leading ROLLBACK is captured. The important
// spec property — Snowflake has no savepoints, so omni must not grow a
// ROLLBACK-TO-SAVEPOINT form — holds: Work stays false and no SAVEPOINT branch
// exists.
func TestRollback_ToSavepointNotConsumed(t *testing.T) {
	stmt, errs := testParseRollbackStmt("ROLLBACK TO SAVEPOINT sp1")
	if len(errs) != 0 {
		t.Fatalf("ROLLBACK TO SAVEPOINT sp1: expected lenient parse, got errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false (TO is not WORK; no savepoint support)")
	}
	// Loc must end right after ROLLBACK, proving TO SAVEPOINT sp1 was not part
	// of the statement.
	if stmt.Loc.End != len("ROLLBACK") {
		t.Errorf("Loc.End = %d, want %d (only ROLLBACK consumed)", stmt.Loc.End, len("ROLLBACK"))
	}
}

// Trailing-junk leniency (engine-wide convention): ROLLBACK FOO.
func TestRollback_TrailingJunkLenient(t *testing.T) {
	stmt, errs := testParseRollbackStmt("ROLLBACK FOO")
	if len(errs) != 0 {
		t.Fatalf("ROLLBACK FOO: expected lenient parse, got errors: %v", errs)
	}
	if stmt.Work {
		t.Error("expected Work=false (FOO is trailing junk, not WORK)")
	}
}

// ---------------------------------------------------------------------------
// Multi-statement: a transaction-bracketed block parses every statement.
// ---------------------------------------------------------------------------

func TestTransaction_MultiStatement(t *testing.T) {
	sql := "BEGIN; COMMIT;"
	result := ParseBestEffort(sql)
	if len(result.Errors) > 0 {
		t.Fatalf("unexpected errors: %v", result.Errors)
	}
	if len(result.File.Stmts) != 2 {
		t.Fatalf("got %d statements, want 2", len(result.File.Stmts))
	}
	if _, ok := result.File.Stmts[0].(*ast.BeginStmt); !ok {
		t.Errorf("stmt[0] = %T, want *ast.BeginStmt", result.File.Stmts[0])
	}
	if _, ok := result.File.Stmts[1].(*ast.CommitStmt); !ok {
		t.Errorf("stmt[1] = %T, want *ast.CommitStmt", result.File.Stmts[1])
	}
}

// ---------------------------------------------------------------------------
// Corpus — official docs (truth1, authoritative) + legacy (truth2, regression)
//
// Oracle agreement: docs.snowflake.com BEGIN/COMMIT/ROLLBACK pages and the
// legacy SnowflakeParser.g4 begin_txn/commit/rollback rules agree exactly on
// the grammar, so the corpus is a pure accept-test. Statements owned by other
// nodes (SHOW TRANSACTIONS, SELECT CURRENT_TRANSACTION(), etc.) are skipped.
// ---------------------------------------------------------------------------

// txnCorpusDirs are official-docs corpora whose transaction-control statements
// are owned by this node. Other statements in these files (SHOW TRANSACTIONS,
// SELECT ...) are skipped as context.
var txnCorpusDirs = []string{
	"testdata/official/begin",
}

func TestTransaction_OfficialCorpus(t *testing.T) {
	for _, dir := range txnCorpusDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatalf("read corpus dir %s: %v", dir, err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
				continue
			}
			path := filepath.Join(dir, entry.Name())
			t.Run(path, func(t *testing.T) {
				data, err := os.ReadFile(path)
				if err != nil {
					t.Fatalf("read %s: %v", path, err)
				}
				assertTxnStatementsParse(t, string(data))
			})
		}
	}
}

// TestTransaction_LegacyCorpus exercises the COMMIT / ROLLBACK examples in the
// legacy other.sql file (the only legacy corpus file with standalone TCL
// statements). Statements owned by other nodes are skipped.
func TestTransaction_LegacyCorpus(t *testing.T) {
	path := "testdata/legacy/other.sql"
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	assertTxnStatementsParse(t, string(data))
}

// assertTxnStatementsParse parses sql and asserts that every transaction-control
// statement (BEGIN / START TRANSACTION / COMMIT / ROLLBACK) parses with no
// errors and to the expected AST type. Statements owned by other DAG nodes are
// skipped.
func assertTxnStatementsParse(t *testing.T, sql string) {
	t.Helper()
	for _, seg := range Split(sql) {
		text := strings.TrimSpace(seg.Text)
		if text == "" {
			continue
		}
		upper := strings.ToUpper(text)

		// Identify the statements this node owns.
		owned := strings.HasPrefix(upper, "BEGIN") ||
			strings.HasPrefix(upper, "START TRANSACTION") ||
			strings.HasPrefix(upper, "COMMIT") ||
			strings.HasPrefix(upper, "ROLLBACK")
		if !owned {
			continue
		}

		result := ParseBestEffort(text)
		if len(result.Errors) > 0 {
			t.Errorf("statement %q: unexpected errors: %v", text, result.Errors)
			continue
		}
		if len(result.File.Stmts) == 0 {
			t.Errorf("statement %q: produced no AST node", text)
			continue
		}
		switch result.File.Stmts[0].(type) {
		case *ast.BeginStmt, *ast.CommitStmt, *ast.RollbackStmt:
			// expected
		default:
			t.Errorf("statement %q: got %T, want a TCL statement node", text, result.File.Stmts[0])
		}
	}
}
