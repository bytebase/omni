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
// enumerated from the vendored pingcap/tidb v8.5.5 parser source (the corpus
// checkout). TiDB rejects statements at parse time through three paths, and
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
// row, never a silently wrong class. Re-derive this set when the pinned
// engine version changes.
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

// unsafeToAdjudicate reports whether sql's first keyword is unsafe to execute
// against the shared oracle. First-keyword match only — "SELECT shutdown_col
// FROM t" is safe — and identifier characters extend the token, so KILLER is
// not KILL. Leading comments are stripped the same way classifyFamily does.
func unsafeToAdjudicate(sql string) bool {
	s := strings.TrimSpace(leadingComment.ReplaceAllString(sql, ""))
	end := 0
	for end < len(s) && isIdentByte(s[end]) {
		end++
	}
	return unsafeKeywords[strings.ToUpper(s[:end])]
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

// adjudicateTiDB probes every non-agreeing row (GAP/OVER/INDETERMINATE)
// against the live container and reclassifies with the container as ground
// truth. Agreeing rows are left alone — label and omni concur; adjudicating
// them would only re-derive the label. Returns the container image digest
// (TIDB_CONTAINER_DIGEST, may be empty) for run meta.
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

	var candidates []int
	for i := range rows {
		switch rows[i].Class {
		case ClassGap, ClassOver, ClassIndeterminate:
			candidates = append(candidates, i)
		}
	}
	log.Printf("adjudicating %d rows against the container", len(candidates))
	start := time.Now()
	prevSQL := "(none)"
	for n, i := range candidates {
		r := &rows[i]
		if prepareAdjudication(r) {
			_, execErr := db.Exec(r.SQL)
			if infra := applyContainerVerdict(r, execErr); infra {
				if pingErr := pingRetry(db); pingErr != nil {
					return "", fmt.Errorf(
						"container died after executing %q (%s:%d; previous statement %q) — extend the unsafe-statement list: %w",
						r.SQL, r.SourcePath, r.Line, prevSQL, pingErr)
				}
			}
			prevSQL = r.SQL
		}
		if (n+1)%200 == 0 {
			log.Printf("adjudicated %d/%d rows", n+1, len(candidates))
		}
	}
	log.Printf("adjudication complete: %d rows in %s", len(candidates), time.Since(start).Round(time.Second))
	return os.Getenv("TIDB_CONTAINER_DIGEST"), nil
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
