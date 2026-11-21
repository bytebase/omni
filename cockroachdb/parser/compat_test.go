package parser

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/bytebase/omni/cockroachdb/internal/testcontainer"
	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/compat"
)

// TestCompatibility holds omni to testdata/compat.tsv: a statement for each
// alternative of CockroachDB's sql.y and its near-misses, with omni's verdict
// and CockroachDB's. It reads the file and starts no server. The grammar is
// subject to the CockroachDB Software License (BSL), so the repo does not hold
// it: the test downloads it once (../testdata/sources.tsv). OMNI_RECORD=1
// judges the rows that have no verdict yet on the pinned CockroachDB.
func TestCompatibility(t *testing.T) {
	compat.Check(t, compat.Suite{
		Path:    filepath.Join("testdata", "compat.tsv"),
		Pins:    filepath.Join("..", "testdata", "pins.tsv"),
		Grammar: enginetest.Fetched(t, filepath.Join("..", "testdata", "sources.tsv"), "sql.y"),
		Tokens:  filepath.Join("..", "testdata", "tokens.tsv"),
		Omni:    omniVerdict,
		Source:  compat.SourceEngine,
		Judge:   compatJudge,
	})
}

// omniVerdict is Parse's answer as a compat.tsv cell: accept or reject.
// CockroachDB's wire-protocol errors leave PgError.Position at 0, so compat.tsv
// verdicts are plain accept/reject while Parse itself still surfaces *ParseError
// with a byte-offset Position extracted from CockroachDB's error detail.
func omniVerdict(sql string) string {
	_, err := Parse(sql)
	if err == nil {
		return compat.Accept
	}
	var pe *ParseError
	if !errors.As(err, &pe) {
		panic(fmt.Sprintf("cockroachdb/parser.Parse(%q) returned %T, want *ParseError", sql, err))
	}
	return compat.Reject
}

// compatJudge asks the pinned CockroachDB whether its parser accepts each
// statement, with the probe compat.RawParseVerdict reads: the statement after
// SELECT 1/0 in one simple query, which CockroachDB parses whole before it
// runs any of it. Nothing of the statement runs, so the judge needs no
// transaction and leaves nothing behind. It runs only in a recording or
// verification, on a connection outside any tape.
func compatJudge(t testing.TB, pin enginetest.Pin) func(string) (string, error) {
	t.Helper()
	db := testcontainer.Judge(t, pin)
	// Sanity-check the probe before judging rows so a release that changes
	// simple-query parsing fails immediately rather than misclassifying rows.
	probe := func(sql string) (string, error) {
		// The protocol sends a query as a C string: the server would read a
		// statement only up to its first zero byte.
		if strings.ContainsRune(sql, 0) {
			return "", fmt.Errorf("the statement holds a zero byte, which the protocol cannot carry")
		}
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, err := db.ExecContext(ctx, compat.RawParseHead+sql)
		return compat.RawParseVerdict(sql, err, crdbState)
	}
	if v, err := probe("SELECT 1"); err != nil || v != compat.Accept {
		t.Fatalf("probe sanity check failed on valid SQL: verdict=%q err=%v", v, err)
	}
	if v, err := probe("SELECT FROM"); err != nil || v != compat.Reject {
		t.Fatalf("probe sanity check failed on invalid SQL: verdict=%q err=%v", v, err)
	}
	return probe
}

// crdbState reads a CockroachDB server error for compat.RawParseVerdict.
// Operational / transport SQLSTATEs (connection, insufficient resources,
// operator intervention / timeout, internal errors) return ok=false so they
// fail the recording rather than silently recording reject.
func crdbState(err error) (sqlstate string, position int, ok bool) {
	var pe *pgconn.PgError
	if !errors.As(err, &pe) {
		return "", 0, false
	}
	if isOperationalSQLState(pe.Code) {
		return "", 0, false
	}
	// CockroachDB reports error locations as a caret diagram in Detail rather
	// than in PgError.Position (which is 0), so verdicts are plain accept/reject.
	return pe.Code, 0, true
}

func isOperationalSQLState(code string) bool {
	if len(code) < 2 {
		return true
	}
	switch code[:2] {
	case "08", // connection exception
		"53", // insufficient resources
		"57", // operator intervention (statement_timeout, shutdown)
		"58", // system error (I/O, OS error)
		"XX": // internal error
		return true
	}
	return false
}

func TestCRDBRawParseVerdict(t *testing.T) {
	cases := []struct {
		name    string
		stmt    string
		err     error
		want    string
		wantErr bool
	}{
		{
			name: "division by zero means parsed",
			stmt: "SELECT 1",
			err:  &pgconn.PgError{Code: "22012", Message: "division by zero"},
			want: compat.Accept,
		},
		{
			name: "syntax error rejected",
			stmt: "SELECT FROM",
			err:  &pgconn.PgError{Code: "42601", Message: `at or near "EOF": syntax error`},
			want: compat.Reject,
		},
		{
			name: "lexical error rejected",
			stmt: "SELECT 'unterminated",
			err:  &pgconn.PgError{Code: "42601", Message: "lexical error: unterminated string"},
			want: compat.Reject,
		},
		{
			name: "grammar action unimplemented rejected",
			stmt: "CREATE TABLE a (b INT, CONSTRAINT c CHECK (b > 0) DEFERRABLE)",
			err:  &pgconn.PgError{Code: "0A000", Message: `at or near ")": syntax error: unimplemented: this syntax`},
			want: compat.Reject,
		},
		{
			name: "grammar action invalid parameter value rejected",
			stmt: "SELECT 1::FLOAT(99)",
			err:  &pgconn.PgError{Code: "22023", Message: "precision for type float must be less than 54 bits"},
			want: compat.Reject,
		},
		{
			name: "help token rejected",
			stmt: "??",
			err:  &pgconn.PgError{Code: "42601", Message: "help token in input: help upon syntax error"},
			want: compat.Reject,
		},
		{
			name: "non-zero position ignored for plain reject",
			stmt: "SELECT FROM",
			err:  &pgconn.PgError{Code: "42601", Position: 18, Message: "syntax error"},
			want: compat.Reject,
		},
		{
			name:    "nil error is not a verdict",
			stmt:    "SELECT 1",
			err:     nil,
			wantErr: true,
		},
		{
			name:    "non-pg error is not a verdict",
			stmt:    "SELECT 1",
			err:     errors.New("connection reset by peer"),
			wantErr: true,
		},
		{
			name:    "statement timeout is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "57014", Message: "query execution canceled"},
			wantErr: true,
		},
		{
			name:    "internal error is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "XX000", Message: "internal error"},
			wantErr: true,
		},
		{
			name:    "insufficient resources is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "53200", Message: "out of memory"},
			wantErr: true,
		},
		{
			name:    "connection exception is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "08006", Message: "connection failure"},
			wantErr: true,
		},
		{
			name:    "system error is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "58030", Message: "io error"},
			wantErr: true,
		},
		{
			name:    "malformed sqlstate is not a verdict",
			stmt:    "SELECT 1",
			err:     &pgconn.PgError{Code: "", Message: "unknown"},
			wantErr: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := compat.RawParseVerdict(tc.stmt, tc.err, crdbState)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got verdict %q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("verdict = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestParseBasicAndEdgeCases(t *testing.T) {
	accepts := []string{
		"SELECT 1",
		"UPSERT INTO a VALUES (1)",
		"CREATE TABLE a (b INT) LOCALITY REGIONAL BY ROW",
		"IMPORT INTO a (b) CSV DATA ('s')",
		"BACKUP INTO 's'",
		"SHOW TABLES",
		"USE defaultdb",
		"BEGIN",
		"COMMIT",
		"SET application_name = 'x'",
		"DO $$ BEGIN RETURN; END $$",
	}
	for _, sql := range accepts {
		if got := omniVerdict(sql); got != compat.Accept {
			t.Errorf("omniVerdict(%q) = %q, want %q", sql, got, compat.Accept)
		}
	}
	rejects := []string{
		"SELECT FROM",
		"DO $$ BEGIN INVALID SYNTAX; END $$",
		"RESTORE FROM 's'",
		"??",
		"SELECT 'unterminated",
	}
	for _, sql := range rejects {
		if got := omniVerdict(sql); got != compat.Reject {
			t.Errorf("omniVerdict(%q) = %q, want %q", sql, got, compat.Reject)
		}
	}
}
