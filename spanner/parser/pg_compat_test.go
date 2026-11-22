package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
	"github.com/bytebase/omni/spanner/internal/testcontainer"
)

// TestPostgresCompatibility holds spanner/parser.ParsePostgres (delegating to
// pg/parser.Parse) to testdata/pg/compat.tsv: a statement for each alternative
// of PostgreSQL 17's grammar (../../pg/testdata/gram.y) and its near-misses,
// with omni's verdict and the pinned Cloud Spanner emulator's verdict in
// DatabaseDialect_POSTGRESQL mode (../testdata/pins.tsv). The test reads the
// committed file and starts no server in offline replay. OMNI_RECORD=1 judges
// rows on the pinned emulator; OMNI_RECORD=verify re-judges every row read-only.
func TestPostgresCompatibility(t *testing.T) {
	pgtd := filepath.Join("..", "..", "pg", "testdata")
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "pg", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: filepath.Join(pgtd, "gram.y"),
		Tokens:  filepath.Join(pgtd, "tokens.tsv"),
		Omni:    omniPostgresVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatPostgresJudge,
	})
}

// omniPostgresVerdict reads the row as one statement via ParsePostgres. Because
// the Spanner emulator does not expose PostgreSQL byte-offset error positions
// in gRPC errors, verdicts are coarse accept/reject (matching TestCompatibility).
func omniPostgresVerdict(sql string) string {
	if _, err := ParsePostgres(sql); err != nil {
		return compat.Reject
	}
	return compat.Accept
}

// compatPostgresJudge asks the pinned Spanner emulator (in
// DatabaseDialect_POSTGRESQL mode) whether it accepts each statement. It runs
// only in a recording (OMNI_RECORD=1 or OMNI_RECORD=verify), outside any tape.
func compatPostgresJudge(t testing.TB, _ enginetest.Pin) func(string) (string, error) {
	t.Helper()
	o := testcontainer.JudgePostgres(t)
	return func(sql string) (string, error) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		kind, err := o.Eval(ctx, sql)
		return spannerPGVerdict(kind, err)
	}
}

// spannerPGVerdict classifies an error from a DatabaseDialect_POSTGRESQL
// Spanner emulator database into a compat verdict (accept or reject).
//
// FAIL-CLOSED: transport, availability, lifecycle, resource-level, non-gRPC,
// generic Internal, and unrecognized gRPC errors always return an error ("no
// verdict") rather than guessing accept or reject.
//
// For DDL (UpdateDatabaseDdl):
//   - Raw-parser and lexer rejections carry the suffix "; failed to parse the
//     DDL statements." (with InvalidArgument or Unimplemented), or match empty
//     DDL / Spanner DDL parse error prefixes -> reject.
//   - Post-parse semantic / unsupported DDL errors (FailedPrecondition,
//     non-resource NotFound, non-parse InvalidArgument, non-parse
//     Unimplemented, and Internal GOOGLESQL_RET_CHECK quirks after parsing) ->
//     accept.
//
// For query and DML (Single().Query / ReadWriteTransaction.Update):
//   - Raw-parser, lexer, single-statement cardinality ("SQL queries are limited
//     to single statements"), and PostgreSQL gram.y parser-action rejections
//     (InvalidArgument or Unimplemented) -> reject.
//   - Post-parse semantic / catalog / unsupported-statement errors
//     (InvalidArgument, Unimplemented, OutOfRange) -> accept.
func spannerPGVerdict(kind string, err error) (string, error) {
	if err == nil {
		return compat.Accept, nil
	}
	st, ok := status.FromError(err)
	if !ok {
		return "", fmt.Errorf("no verdict (%s, non-gRPC): %w", kind, err)
	}
	code := st.Code()
	msg := st.Message()

	// Transport, availability, and lifecycle errors are never grammar verdicts.
	switch code {
	case codes.Unavailable, codes.DeadlineExceeded, codes.Canceled,
		codes.Aborted, codes.ResourceExhausted, codes.Unknown,
		codes.Unauthenticated, codes.PermissionDenied:
		return "", fmt.Errorf("no verdict (%s, %s): %s", kind, code, msg)
	}
	// Resource-level lifecycle errors (scratch database/instance/session gone)
	// are infra failures, never per-statement semantic verdicts.
	if (code == codes.NotFound || code == codes.AlreadyExists) && isPGResourceLevel(msg) {
		return "", fmt.Errorf("no verdict (%s, resource %s): %s", kind, code, msg)
	}

	switch kind {
	case "ddl":
		switch code {
		case codes.InvalidArgument, codes.Unimplemented:
			if isPGDDLParseError(msg) {
				return compat.Reject, nil
			}
			return compat.Accept, nil
		case codes.FailedPrecondition, codes.NotFound:
			return compat.Accept, nil
		case codes.Internal:
			// Emulator post-parse quirk in schema_updater / pg_parse_tree_validator /
			// pg_to_spanner_ddl_translator after raw_parser succeeded.
			if strings.Contains(msg, "GOOGLESQL_RET_CHECK") {
				return compat.Accept, nil
			}
			return "", fmt.Errorf("no verdict (ddl, Internal): %s", msg)
		default:
			return "", fmt.Errorf("no verdict (ddl, %s): %s", code, msg)
		}

	case "query", "dml":
		switch code {
		case codes.InvalidArgument, codes.Unimplemented:
			if isPGQueryParseError(msg) {
				return compat.Reject, nil
			}
			return compat.Accept, nil
		case codes.OutOfRange:
			// Runtime error (e.g. division by zero) after successful parse.
			return compat.Accept, nil
		default:
			return "", fmt.Errorf("no verdict (%s, %s): %s", kind, code, msg)
		}

	default:
		return "", fmt.Errorf("no verdict (unknown kind %q, %s): %s", kind, code, msg)
	}
}

func isPGResourceLevel(msg string) bool {
	for _, p := range []string{
		"Database not found:", "Instance not found:", "Session not found:",
		"Database already exists:", "Instance already exists:",
	} {
		if strings.Contains(msg, p) {
			return true
		}
	}
	return false
}

func isPGDDLParseError(msg string) bool {
	if strings.HasSuffix(msg, "; failed to parse the DDL statements.") ||
		strings.HasPrefix(msg, "Error parsing Spanner DDL statement:") ||
		msg == "DDL statement is empty." ||
		msg == "No statement found." {
		return true
	}
	return isPGQueryParseError(msg)
}

// pgParseErrorSubstrings enumerates the exact PostgreSQL raw-parser (scan.l and
// gram.y) and Spanner PG parser-action rejection messages observed on the
// pinned emulator (1.5.54). Each of these fires during parsing (before AST
// analysis or execution).
var pgParseErrorSubstrings = []string{
	// Core Bison / GoogleSQL syntax errors
	"[ERROR] syntax error at ",
	"Syntax error:",

	// PostgreSQL lexer (scan.l) errors
	"[ERROR] unterminated quoted string",
	"[ERROR] unterminated quoted identifier",
	"[ERROR] unterminated dollar-quoted string",
	"[ERROR] unterminated /* comment",
	"[ERROR] unterminated bit string literal",
	"[ERROR] unterminated hexadecimal string literal",
	"[ERROR] invalid Unicode escape",
	"[ERROR] invalid Unicode surrogate pair",
	"[ERROR] invalid hexadecimal digit",
	"[ERROR] zero-length delimited identifier",
	"[ERROR] trailing junk after numeric literal",
	"[ERROR] trailing junk after parameter",

	// PostgreSQL gram.y & Spanner PG parser-action rejections
	"[ERROR] missing argument",
	"[ERROR] WITH TIES cannot be specified without ORDER BY clause",
	"[ERROR] wrong number of parameters on left side of OVERLAPS expression",
	"[ERROR] wrong number of parameters on right side of OVERLAPS expression",
	"[ERROR] frame starting from following row cannot end with current row",
	"[ERROR] frame end cannot be UNBOUNDED PRECEDING",
	"[ERROR] frame start cannot be UNBOUNDED FOLLOWING",
	"[ERROR] LIMIT #,# syntax is not supported",
	"[ERROR] unrecognized JSON encoding:",
	"[ERROR] conflicting or redundant NULL / NOT NULL declarations for column ",
	"[ERROR] unrecognized column option ",
	"[ERROR] UNIQUE predicate is not yet implemented",
	"[ERROR] current database cannot be changed",
	"[ERROR] XML expressions are not supported",
	"[ERROR] unrecognized role option ",
	"[ERROR] invalid publication object list",
	"[ERROR] unrecognized partitioning strategy ",
	"[ERROR] unrecognized hash partition bound specification ",
	"[ERROR] duplicate trigger events specified ",
	"[ERROR] unrecognized row security option ",
	"[ERROR] CREATE ASSERTION is not yet implemented",
	"[ERROR] CREATE EXTENSION ... FROM is no longer supported",
	"[ERROR] MATCH PARTIAL not yet implemented",
	"[ERROR] UNENCRYPTED PASSWORD is no longer supported",
}

func isPGQueryParseError(msg string) bool {
	if msg == "SQL queries are limited to single statements" {
		return true
	}
	for _, sub := range pgParseErrorSubstrings {
		if strings.Contains(msg, sub) {
			return true
		}
	}
	return false
}

func TestSpannerPGVerdict(t *testing.T) {
	cases := []struct {
		name    string
		kind    string
		err     error
		want    string
		wantErr bool
	}{
		// Success
		{name: "query ok", kind: "query", err: nil, want: compat.Accept},
		{name: "dml ok", kind: "dml", err: nil, want: compat.Accept},
		{name: "ddl ok", kind: "ddl", err: nil, want: compat.Accept},

		// Query & DML parser / lexer / cardinality rejections
		{name: "query syntax at or near", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] syntax error at or near "zzz"`), want: compat.Reject},
		{name: "query syntax eof", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] syntax error at end of input`), want: compat.Reject},
		{name: "query unterminated string", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] unterminated quoted string at or near "'abc"`), want: compat.Reject},
		{name: "query empty single statement", kind: "query", err: status.Error(codes.InvalidArgument, `SQL queries are limited to single statements`), want: compat.Reject},
		{name: "query overlaps action error", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] wrong number of parameters on left side of OVERLAPS expression`), want: compat.Reject},
		{name: "query with ties action error", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] WITH TIES cannot be specified without ORDER BY clause`), want: compat.Reject},
		{name: "query xml parser rejection", kind: "query", err: status.Error(codes.Unimplemented, `[ERROR] XML expressions are not supported`), want: compat.Reject},
		{name: "query unique predicate rejection", kind: "query", err: status.Error(codes.Unimplemented, `[ERROR] UNIQUE predicate is not yet implemented`), want: compat.Reject},
		{name: "query set catalog rejection", kind: "query", err: status.Error(codes.Unimplemented, `[ERROR] current database cannot be changed`), want: compat.Reject},
		{name: "dml syntax error", kind: "dml", err: status.Error(codes.InvalidArgument, `[ERROR] syntax error at or near "ON"`), want: compat.Reject},

		// Query & DML post-parse semantic / unsupported acceptances
		{name: "query missing table", kind: "query", err: status.Error(codes.InvalidArgument, `relation "nonexistent" does not exist`), want: compat.Accept},
		{name: "query missing column", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] column "a" does not exist`), want: compat.Accept},
		{name: "query invalid input syntax for type", kind: "query", err: status.Error(codes.InvalidArgument, `[ERROR] invalid input syntax for type bigint: "s"`), want: compat.Accept},
		{name: "query no output column", kind: "query", err: status.Error(codes.InvalidArgument, `Statement must produce at least one output column`), want: compat.Accept},
		{name: "query utility unsupported", kind: "query", err: status.Error(codes.Unimplemented, `Utility SQL statements are not supported`), want: compat.Accept},
		{name: "query window unsupported", kind: "query", err: status.Error(codes.Unimplemented, `Statements with WINDOW clauses are not supported`), want: compat.Accept},
		{name: "dml merge unsupported", kind: "dml", err: status.Error(codes.Unimplemented, `[ERROR] statement type "MergeStmt" is not supported`), want: compat.Accept},
		{name: "query out of range", kind: "query", err: status.Error(codes.OutOfRange, `division by zero: 1 / 0`), want: compat.Accept},

		// DDL parser / lexer / action rejections
		{name: "ddl syntax at or near", kind: "ddl", err: status.Error(codes.InvalidArgument, `[ERROR] syntax error at or near "zzz"; failed to parse the DDL statements.`), want: compat.Reject},
		{name: "ddl syntax eof", kind: "ddl", err: status.Error(codes.InvalidArgument, `[ERROR] syntax error at end of input; failed to parse the DDL statements.`), want: compat.Reject},
		{name: "ddl role option error", kind: "ddl", err: status.Error(codes.InvalidArgument, `[ERROR] unrecognized role option "zzz"; failed to parse the DDL statements.`), want: compat.Reject},
		{name: "ddl create assertion unimplemented", kind: "ddl", err: status.Error(codes.Unimplemented, `[ERROR] CREATE ASSERTION is not yet implemented; failed to parse the DDL statements.`), want: compat.Reject},
		{name: "ddl xml unimplemented", kind: "ddl", err: status.Error(codes.Unimplemented, `[ERROR] XML expressions are not supported; failed to parse the DDL statements.`), want: compat.Reject},
		{name: "ddl empty", kind: "ddl", err: status.Error(codes.InvalidArgument, `DDL statement is empty.`), want: compat.Reject},
		{name: "ddl spanner parse prefix", kind: "ddl", err: status.Error(codes.InvalidArgument, `Error parsing Spanner DDL statement: CREATE TABLE`), want: compat.Reject},

		// DDL post-parse semantic / unsupported acceptances
		{name: "ddl unsupported statement", kind: "ddl", err: status.Error(codes.FailedPrecondition, `Statement is not supported.`), want: compat.Accept},
		{name: "ddl missing primary key", kind: "ddl", err: status.Error(codes.FailedPrecondition, `Primary key must be defined for table "t".`), want: compat.Accept},
		{name: "ddl table not found", kind: "ddl", err: status.Error(codes.NotFound, `Table not found: t`), want: compat.Accept},
		{name: "ddl sequence kind missing", kind: "ddl", err: status.Error(codes.InvalidArgument, `The sequence does not have a valid sequence kind.`), want: compat.Accept},
		{name: "ddl view analysis error", kind: "ddl", err: status.Error(codes.InvalidArgument, "Error analyzing the definition of view `a`: [ERROR] column \"zzz\" does not exist"), want: compat.Accept},
		{name: "ddl ret_check quirk", kind: "ddl", err: status.Error(codes.Internal, `GOOGLESQL_RET_CHECK failure (backend/schema/updater/schema_updater.cc:1182) false Unsupported ddl statement: 1`), want: compat.Accept},

		// Fail-closed errors (must NEVER return accept or reject)
		{name: "non-grpc error", kind: "query", err: errors.New("connection reset by peer"), wantErr: true},
		{name: "unavailable", kind: "query", err: status.Error(codes.Unavailable, `transport is closing`), wantErr: true},
		{name: "deadline exceeded", kind: "ddl", err: status.Error(codes.DeadlineExceeded, `context deadline exceeded`), wantErr: true},
		{name: "canceled", kind: "dml", err: status.Error(codes.Canceled, `context canceled`), wantErr: true},
		{name: "aborted", kind: "dml", err: status.Error(codes.Aborted, `transaction aborted`), wantErr: true},
		{name: "database not found", kind: "ddl", err: status.Error(codes.NotFound, `Database not found: projects/omni/instances/pg-compat/databases/pgd_1`), wantErr: true},
		{name: "session not found", kind: "query", err: status.Error(codes.NotFound, `Session not found: projects/omni/...`), wantErr: true},
		{name: "generic internal ddl", kind: "ddl", err: status.Error(codes.Internal, `unexpected server crash`), wantErr: true},
		{name: "generic internal query", kind: "query", err: status.Error(codes.Internal, `GOOGLESQL_RET_CHECK in query`), wantErr: true},
		{name: "failed precondition on query", kind: "query", err: status.Error(codes.FailedPrecondition, `unexpected state`), wantErr: true},
		{name: "unknown kind", kind: "other", err: status.Error(codes.InvalidArgument, `boom`), wantErr: true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := spannerPGVerdict(tc.kind, tc.err)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("spannerPGVerdict(%q, %v) = %q, nil; want error", tc.kind, tc.err, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("spannerPGVerdict(%q, %v) unexpected error: %v", tc.kind, tc.err, err)
			}
			if got != tc.want {
				t.Fatalf("spannerPGVerdict(%q, %v) = %q, want %q", tc.kind, tc.err, got, tc.want)
			}
		})
	}
}
