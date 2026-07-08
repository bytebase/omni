package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
)

// tidbParseRejectCodes is the parse-abort error space of the pinned oracle,
// enumerated from the parser source in the pinned pingcap/tidb v8.5.5 corpus
// checkout. TiDB rejects statements at parse time through three paths, and
// only the first uses the classic parse-error codes:
//
//   - yacc/scanner errors: ErrParse 1064 (plain scanner errors are wrapped to
//     1064 by the server) and ErrSyntax 1149 (pkg/parser/yy_parser.go:33-36).
//   - named grammar-action aborts (yylex.AppendError + return 1 in parser.go,
//     plus two AppendError accumulations without return 1, which still fail
//     ParseSQL): ErrWrongFieldTerminators 1083, ErrWrongDBName 1102,
//     ErrUnknownCharacterSet 1115, ErrWrongArguments 1210, ErrWrongUsage 1221,
//     ErrUnknownCollation 1273, ErrTooBigPrecision 1426, ErrWrongValue 1525,
//     ErrUnknownAlterAlgorithm 1800, ErrUnknownAlterLock 1801,
//     ErrInvalidYearColumnLength 1818.
//   - pkg/parser/ast validators called from grammar actions (ColumnDef and
//     PartitionOptions Validate, ast/ddl.go): ErrPartitionRequiresValues 1479,
//     ErrPartitionWrongValues 1480, ErrPartitionWrongNoPart 1484,
//     ErrPartitionWrongNoSubpart 1485, ErrPartitionsMustBeDefined 1492,
//     ErrSubpartition 1500, ErrNoParts 1504, ErrCoalescePartitionNoPartition
//     1515, ErrPartitionColumnList 1653, ErrTooManyValues 1657,
//     ErrRowSinglePartitionField 1658, ErrWrongPartitionTypeExpectedSystemTime
//     4113, ErrSystemVersioningWrongPartitions 4128.
//
// Some of these codes are also raised at runtime for statements that parse
// (observed once in the v8.5.5 sweep: 1102 on an accept-labeled row). Such a
// collision misreads runtime-reject as parse-reject, which classify() then
// fails closed into INDETERMINATE label_container_disagree — a manual-queue
// row, never a silently wrong class. That net requires an upstream label:
// a label-less row (H3-cleared duplicate_label_conflict) hitting such a
// collision classifies from the container verdict alone and would land in
// OVER/AGREE_REJECT silently. Re-derive this set when the pinned engine
// version changes.
var tidbParseRejectCodes = map[uint16]bool{
	1064: true, 1149: true,
	1083: true, 1102: true, 1115: true, 1210: true, 1221: true, 1273: true,
	1426: true, 1525: true, 1800: true, 1801: true, 1818: true,
	1479: true, 1480: true, 1484: true, 1485: true, 1492: true, 1500: true,
	1504: true, 1515: true, 1653: true, 1657: true, 1658: true,
	4113: true, 4128: true,
}

// classifyTiDBExecError maps a driver error to an engine parse verdict:
// tidbParseRejectCodes = the parser rejected; anything else MySQL-coded =
// parsed (8108 "Unsupported type" = parsed but unexecutable; semantic and
// runtime errors = parsed). Non-MySQL errors are infra — VerdictNone, never
// accept/reject (fail-closed).
func classifyTiDBExecError(err error) (Verdict, int, string) {
	if err == nil {
		return VerdictAccept, 0, ""
	}
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return VerdictNone, 0, err.Error()
	}
	if tidbParseRejectCodes[me.Number] {
		return VerdictReject, int(me.Number), me.Message
	}
	return VerdictAccept, int(me.Number), me.Message
}

// unsafeKeywords lead statements that can take down or destabilize the shared
// oracle: the corpus is a parser test suite, so it literally contains
// `shutdown`, `restart`, and KILL variants (parser_test.go:5958-5968).
var unsafeKeywords = map[string]bool{"SHUTDOWN": true, "KILL": true, "RESTART": true}

// unsafeToAdjudicate reports whether sql leads with a statement unsafe to
// execute against the shared oracle: a first keyword in unsafeKeywords, or
// SET PASSWORD (would change the oracle's credentials mid-sweep,
// parser_test.go:1386-1387). First-statement match only — "SELECT
// shutdown_col FROM t" is safe — and identifier characters extend the token,
// so KILLER is not KILL. Leading comments are stripped the same way
// classifyFamily does. Best-effort deny-list, not a safety proof: the
// ping-abort in probeRow and the disposable container remain the backstop.
func unsafeToAdjudicate(sql string) bool {
	s := strings.TrimSpace(leadingComment.ReplaceAllString(sql, ""))
	first, rest := nextKeyword(s)
	if unsafeKeywords[first] {
		return true
	}
	if first == "SET" {
		second, _ := nextKeyword(rest)
		return second == "PASSWORD"
	}
	return false
}

// nextKeyword returns s's leading identifier-shaped token upper-cased, plus
// the remainder after it. Empty token when s starts with a non-ident byte.
func nextKeyword(s string) (string, string) {
	s = strings.TrimSpace(s)
	end := 0
	for end < len(s) && isIdentByte(s[end]) {
		end++
	}
	return strings.ToUpper(s[:end]), s[end:]
}

func isIdentByte(c byte) bool {
	return c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_'
}

// normalizeTiDBDSN forces the settings the sweep is incorrect or unbounded
// without (H1): multiStatements=true (corpus rows contain multi-statement
// SQL; without it the server 1064s the whole batch — false parse-rejects)
// and dial/read/write timeouts (a hanging statement must not stall the
// sweep; a driver timeout is a non-MySQL error, so the row lands in
// INDETERMINATE infra_error). Explicit timeouts in the DSN are respected.
func normalizeTiDBDSN(dsn string) (string, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("invalid TIDB_DSN: %w", err)
	}
	cfg.MultiStatements = true
	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}
	if cfg.ReadTimeout == 0 {
		cfg.ReadTimeout = 10 * time.Second
	}
	if cfg.WriteTimeout == 0 {
		cfg.WriteTimeout = 10 * time.Second
	}
	return cfg.FormatDSN(), nil
}

// prepareAdjudication applies the pre-exec hazard checks to one candidate row
// and reports whether it may be sent to the container.
//
// H2: unsafe statements are marked INDETERMINATE/unsafe_to_adjudicate without
// ever touching the container (idempotent across runs — the predicate depends
// only on the SQL). H3: a duplicate_label_conflict row kept an arbitrary
// first-seen label; that label is not ground truth, so it is cleared and the
// container verdict becomes the sole truth — otherwise adjudication would
// coin-flip such rows into label_container_disagree.
func prepareAdjudication(r *Row) bool {
	if unsafeToAdjudicate(r.SQL) {
		r.Class = ClassIndeterminate
		r.ClassifierReason = "unsafe_to_adjudicate"
		r.EngineVerdict = VerdictNone
		r.DivergenceKey = "" // INDETERMINATE rows are not clustered
		return false
	}
	if r.ClassifierReason == "duplicate_label_conflict" {
		r.Expected = VerdictNone
	}
	return true
}

// applyContainerVerdict folds one Exec outcome into the row and reclassifies.
// Infra errors (verdict none) are fail-closed: the row becomes INDETERMINATE
// infra_error and the caller must check whether the container is still alive.
func applyContainerVerdict(r *Row, execErr error) (infra bool) {
	v, code, msg := classifyTiDBExecError(execErr)
	if v == VerdictNone {
		r.Class = ClassIndeterminate
		r.ClassifierReason = "infra_error"
		r.EngineVerdict = VerdictNone
		r.RawErrorCode = 0
		r.RawErrorMessage = msg
		r.DivergenceKey = "" // INDETERMINATE rows are not clustered
		return true
	}
	r.EngineVerdict = v
	r.RawErrorCode = code
	r.RawErrorMessage = msg
	classify(r)
	return false
}

// adjudicateTiDB probes every non-agreeing row against the live container
// and reclassifies with the container as ground truth. Returns the container
// image digest (TIDB_CONTAINER_DIGEST, may be empty) for run meta.
func adjudicateTiDB(rows []Row) (string, error) {
	dsn := os.Getenv("TIDB_DSN")
	if dsn == "" {
		return "", errors.New("TIDB_DSN is not set: run ./start_tidb.sh and export the DSN line it prints")
	}
	dsn, err := normalizeTiDBDSN(dsn)
	if err != nil {
		return "", err
	}
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	// No idle reuse: every row gets a fresh session, so session state (USE,
	// SET sql_mode — which changes how later statements *parse*) cannot leak
	// across rows. Localhost dials are cheap; verdict fidelity is not.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(0)
	if err := pingRetry(db); err != nil {
		return "", fmt.Errorf("TiDB not reachable (is ./start_tidb.sh running?): %w", err)
	}

	candidates := adjudicationCandidates(rows)
	log.Printf("adjudicating %d rows against the container", len(candidates))
	start := time.Now()
	prevSQL := "(none)"
	for n, i := range candidates {
		if n > 0 && n%200 == 0 {
			log.Printf("adjudicated %d/%d rows", n, len(candidates))
		}
		r := &rows[i]
		if !prepareAdjudication(r) {
			continue // unsafe — never touches the container
		}
		if err := probeRow(db, r, prevSQL); err != nil {
			return "", err
		}
		prevSQL = r.SQL
	}
	log.Printf("adjudication complete: %d rows in %s", len(candidates), time.Since(start).Round(time.Second))
	return os.Getenv("TIDB_CONTAINER_DIGEST"), nil
}

// adjudicationCandidates returns the indexes of the rows the container should
// arbitrate: the non-agreeing classes (GAP/OVER/INDETERMINATE). Agreeing rows
// are left alone — label and omni concur; adjudicating them would only
// re-derive the label.
func adjudicationCandidates(rows []Row) []int {
	var idx []int
	for i := range rows {
		switch rows[i].Class {
		case ClassGap, ClassOver, ClassIndeterminate:
			idx = append(idx, i)
		}
	}
	return idx
}

// probeRow sends one prepared row to the container and folds the outcome in.
// After an infra error it verifies the oracle is still alive: a dead
// container aborts the sweep, naming the statements that preceded the death
// (so the unsafe-statement list can be extended), instead of silently
// poisoning every remaining row with infra_error.
func probeRow(db *sql.DB, r *Row, prevSQL string) error {
	_, execErr := db.Exec(r.SQL)
	if infra := applyContainerVerdict(r, execErr); !infra {
		return nil
	}
	if pingErr := pingRetry(db); pingErr != nil {
		return fmt.Errorf(
			"container died after executing %q (%s:%d; previous statement %q) — extend the unsafe-statement list: %w",
			r.SQL, r.SourcePath, r.Line, prevSQL, pingErr)
	}
	return nil
}

// pingRetry gives the server a few chances: the start script only waits for
// the port to listen, which is not ready-to-serve. Also used as the liveness
// check after a connection-level Exec error, where the retries guard against
// declaring a transient blip a container death.
func pingRetry(db *sql.DB) error {
	var err error
	for range 20 {
		if err = db.Ping(); err == nil {
			return nil
		}
		time.Sleep(500 * time.Millisecond)
	}
	return err
}
