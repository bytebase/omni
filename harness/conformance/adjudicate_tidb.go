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

// tidbConnectionScopeCodes are MySQL-coded errors that describe the
// connection, not the statement — they are infra (VerdictNone), never a
// parse verdict. The sweep connects as root with the container's real
// credentials, so 1045 ER_ACCESS_DENIED_ERROR can never occur statement-
// level: it appears only when a probed batch mutated the credentials
// mid-sweep (e.g. `SELECT 1; SET PASSWORD = 'x'` — the container stays up,
// so the ping-abort never fires on its own) and every later fresh
// connection fails its HANDSHAKE with 1045. Mapping that to "parsed" would
// silently poison every remaining verdict; VerdictNone routes it to
// probeRow's ping-abort, which names the culprit statement.
//
// Deliberately NOT in this set:
//   - 1049 ER_BAD_DB_ERROR: a statement-level `USE nonexistent` on a healthy
//     connection legitimately returns 1049 for a statement that parsed.
//     Connection-level 1049 (handshake against a dropped default schema) is
//     made impossible instead: normalizeTiDBDSN forces an empty default
//     schema, so the handshake names nothing droppable.
//   - 1044 (db access denied) / 1046 (no database selected): statement-level
//     semantic codes — they prove the statement parsed.
var tidbConnectionScopeCodes = map[uint16]bool{1045: true}

// classifyTiDBExecError maps a driver error to an engine parse verdict:
// tidbParseRejectCodes = the parser rejected; anything else MySQL-coded =
// parsed (8108 "Unsupported type" = parsed but unexecutable; semantic and
// runtime errors = parsed) — except tidbConnectionScopeCodes, which are
// about the connection, not the statement. Non-MySQL errors are infra —
// VerdictNone, never accept/reject (fail-closed).
func classifyTiDBExecError(err error) (Verdict, int, string) {
	if err == nil {
		return VerdictAccept, 0, ""
	}
	var me *mysql.MySQLError
	if !errors.As(err, &me) {
		return VerdictNone, 0, err.Error()
	}
	if tidbConnectionScopeCodes[me.Number] {
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

// unsafeToAdjudicate reports whether any statement in sql is unsafe to
// execute against the shared oracle. One normalization pipeline, mirroring
// how the server reads the text: (1) executable-comment markers are
// neutralized — their content is EXECUTED, so it must survive as scannable
// SQL; (2) all remaining ordinary comments are stripped — the server treats
// them as whitespace ANYWHERE, so `SET /*c*/ PASSWORD` must not blind the
// keyword scan, and stripping before the split means a `;` inside an
// ordinary comment cannot fabricate a phantom statement (its content is
// never executed); (3) the text splits on `;` and every statement in the
// batch is scanned, not just the first: a mid-batch SHUTDOWN kills the
// oracle just as dead as a leading one, and a mid-batch SET PASSWORD doesn't
// kill it at all — it poisons the handshake of every later fresh connection,
// so the ping-abort backstop never fires. The `;` split is naive: a
// semicolon inside a string literal splits too, which can only over-match (a
// safe row skipped into INDETERMINATE unsafe_to_adjudicate), never miss an
// unsafe statement — the deny-list deliberately errs conservative.
// Best-effort deny-list, not a safety proof: the ping-abort in probeRow and
// the disposable container remain the backstop.
func unsafeToAdjudicate(sql string) bool {
	s := stripComments(neutralizeExecutableComments(sql))
	for _, stmt := range strings.Split(s, ";") {
		if unsafeStatement(stmt) {
			return true
		}
	}
	return false
}

// unsafeStatement checks a single comment-free statement (unsafeToAdjudicate
// has already neutralized executable comments and stripped ordinary ones): a
// first keyword in unsafeKeywords, a SET whose target outlives the probe's
// session (see unsafeSetTarget), or an account-mutation DCL statement (see
// accountMutationKeywords). First-keyword match only — "SELECT shutdown_col
// FROM t" is safe — and identifier characters extend the token, so KILLER is
// not KILL.
func unsafeStatement(stmt string) bool {
	first, rest := nextKeyword(stmt)
	if unsafeKeywords[first] {
		return true
	}
	if first == "SET" {
		return unsafeSetTarget(rest)
	}
	if accountMutationKeywords[first] {
		second, _ := nextKeyword(rest)
		return second == "USER"
	}
	return false
}

// accountMutationKeywords lead the account-mutation DCL forms — <keyword>
// USER — that can change the identity every later probe connection
// handshakes with: the corpus tests these on 'root' itself, and a probed
// `ALTER USER root IDENTIFIED BY 'x'` (container-verified on the pinned
// oracle) makes every later fresh handshake fail 1045. Same handshake-poison
// channel as SET PASSWORD: the container stays up, so the ping-abort
// backstop never fires. Deliberately NOT blocked: GRANT/REVOKE —
// container-verified unable to change our connection identity (no user
// auto-create: 1410 "You are not allowed to create a user with GRANT", and a
// 5.7-style `GRANT ... IDENTIFIED BY` on an existing user executes but
// leaves its credentials untouched) — and blocking them would move
// legitimately divergent DCL rows out of adjudication.
var accountMutationKeywords = map[string]bool{"CREATE": true, "ALTER": true, "DROP": true, "RENAME": true}

// neutralizeExecutableComments makes executable comments visible to the
// deny-list scan. `/*! SET PASSWORD = 'x' */`, `/*!40101 SET GLOBAL ... */`
// and TiDB's `/*T! ... */` look like comments but are EXECUTED by the server,
// so comment-stripping them would wave an unsafe statement through. The
// marker forms mirror omni's own splitter (tidb/parser/split.go,
// Segment.Empty): `/*!` and `/*T!` (uppercase T), matched by prefix. Each
// marker plus its optional version digits is blanked, as is its matching
// `*/` closer, so the content survives as plain scannable SQL; ordinary
// `/* ... */` comments keep their delimiters for stripComments. A marker
// inside a string literal over-matches (the row skips to INDETERMINATE) —
// acceptable, the deny-list errs conservative.
func neutralizeExecutableComments(stmt string) string {
	if !strings.Contains(stmt, "/*!") && !strings.Contains(stmt, "/*T!") {
		return stmt
	}
	b := []byte(stmt)
	i := 0
	for i+1 < len(b) {
		if b[i] != '/' || b[i+1] != '*' {
			i++
			continue
		}
		markerLen := 0
		switch {
		case i+2 < len(b) && b[i+2] == '!':
			markerLen = 3
		case i+3 < len(b) && b[i+2] == 'T' && b[i+3] == '!':
			markerLen = 4
		}
		if markerLen == 0 {
			// Ordinary comment: leave it intact for stripComments.
			end := strings.Index(stmt[i+2:], "*/")
			if end < 0 {
				break
			}
			i += 2 + end + 2
			continue
		}
		j := i + markerLen
		for j < len(b) && b[j] >= '0' && b[j] <= '9' {
			j++
		}
		for k := i; k < j; k++ {
			b[k] = ' '
		}
		end := strings.Index(stmt[j:], "*/")
		if end < 0 {
			break
		}
		b[j+end], b[j+end+1] = ' ', ' '
		i = j + end + 2
	}
	return string(b)
}

// stripComments replaces every ordinary comment — `/* ... */` anywhere,
// `-- ` and `#` line comments to end-of-line — with a single space, the way
// the server's lexer treats them. Executable-comment markers must already be
// neutralized, or their content (which the server EXECUTES) would be
// stripped here as if it were ordinary. String literals and backtick
// identifiers are skipped: a comment opener inside a literal is literal
// text, and stripping from it would swallow the real statements that follow
// (`SELECT 'a -- b'; KILL 5`) — the one direction the deny-list must never
// err. The literal lexing is still best-effort (default escape rules; a
// session sql_mode could diverge): divergence mangles comment-like literal
// text, which can only over-match into a conservative unsafe skip.
func stripComments(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	i := 0
	for i < len(s) {
		c := s[i]
		switch {
		case c == '\'' || c == '"' || c == '`':
			end := skipQuoted(s, i)
			b.WriteString(s[i:end])
			i = end
		case c == '/' && i+1 < len(s) && s[i+1] == '*':
			b.WriteByte(' ')
			end := strings.Index(s[i+2:], "*/")
			if end < 0 {
				return b.String() // unterminated: comment runs to end of input
			}
			i += 2 + end + 2
		case c == '#', isDashCommentStart(s, i):
			b.WriteByte(' ')
			for i < len(s) && s[i] != '\n' {
				i++
			}
		default:
			b.WriteByte(c)
			i++
		}
	}
	return b.String()
}

// isDashCommentStart reports whether s[i:] starts a MySQL `--` line comment:
// the dashes must be followed by whitespace or end of input (`--x` is a
// double negation, not a comment) — same rule as tidb/parser/split.go's
// isDashComment.
func isDashCommentStart(s string, i int) bool {
	if s[i] != '-' || i+1 >= len(s) || s[i+1] != '-' {
		return false
	}
	if i+2 >= len(s) {
		return true
	}
	c := s[i+2]
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// skipQuoted returns the index just past the quoted region opening at s[i]
// (', ", or `), using MySQL lexing: doubled-quote escapes for all three,
// backslash escapes for ' and " — mirroring tidb/parser/split.go's skip
// helpers. An unterminated region runs to end of input.
func skipQuoted(s string, i int) int {
	q := s[i]
	i++
	for i < len(s) {
		c := s[i]
		if c == '\\' && q != '`' {
			i += 2
			continue
		}
		if c != q {
			i++
			continue
		}
		i++
		if i < len(s) && s[i] == q {
			i++ // doubled-quote escape
			continue
		}
		return i
	}
	return i
}

// unsafeSetTarget reports whether a SET statement's tail mutates state that
// outlives the probe's session. Fresh-connection-per-row guards session
// state only: SET PASSWORD invalidates the credentials every later handshake
// uses (parser_test.go:1386-1387), and GLOBAL/PERSIST mutations (e.g. a
// global sql_mode of ANSI_QUOTES) change how every later probe *parses*.
// The @@GLOBAL/@@PERSIST forms are matched by prefix, not by token, because
// no space follows the scope in `@@global.sql_mode`; SET NAMES /
// SET sql_mode / SET @@session.x / SET @v stay adjudicable.
func unsafeSetTarget(rest string) bool {
	second, _ := nextKeyword(rest)
	if second == "PASSWORD" || second == "GLOBAL" || second == "PERSIST" {
		return true
	}
	r := strings.ToUpper(strings.TrimSpace(rest))
	return strings.HasPrefix(r, "@@GLOBAL") || strings.HasPrefix(r, "@@PERSIST")
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
// SQL; without it the server 1064s the whole batch — false parse-rejects),
// no default schema (a default schema is droppable by adjudicated DDL —
// `DROP DATABASE test` — after which every later fresh connection fails its
// handshake with 1049, silently poisoning the sweep; without one,
// unqualified-name statements fail statement-level 1046, which classifies
// identically as "parsed"), and dial/read/write timeouts (a hanging
// statement must not stall the sweep; a driver timeout is a non-MySQL
// error, so the row lands in INDETERMINATE infra_error). Explicit timeouts
// in the DSN are respected.
func normalizeTiDBDSN(dsn string) (string, error) {
	cfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return "", fmt.Errorf("invalid TIDB_DSN: %w", err)
	}
	cfg.MultiStatements = true
	cfg.DBName = ""
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
