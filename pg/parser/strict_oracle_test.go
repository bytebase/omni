package parser

import (
	"context"
	"errors"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5/pgconn"
)

// syntaxCase is one statement and whether PostgreSQL's grammar accepts it.
type syntaxCase struct {
	sql    string
	accept bool
}

// assertSyntaxMatchesPG checks every case against PostgreSQL 17's raw
// parser and then against Parse. A statement the raw parser refuses must
// fail Parse with the same SQLSTATE and message, and at the same position
// when the server reports one: a check a grammar rule makes, such as
// processCASbits, reports none. Checking PG first keeps the expectation
// tied to the engine rather than to memory.
func assertSyntaxMatchesPG(t *testing.T, cases []syntaxCase) {
	t.Helper()
	o := startFirstSetOracle(t)
	for _, c := range cases {
		t.Run(c.sql, func(t *testing.T) {
			pgErr := o.rawParse(t, c.sql)
			if pgAccept := pgErr == nil; pgAccept != c.accept {
				t.Fatalf("postgres accept = %v, want %v (err = %v)", pgAccept, c.accept, pgErr)
			}
			_, err := Parse(c.sql)
			if omniAccept := err == nil; omniAccept != c.accept {
				t.Fatalf("omni accept = %v, want %v (err = %v)", omniAccept, c.accept, err)
			}
			if c.accept {
				return
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("omni error = %T, want *ParseError", err)
			}
			code := pe.Code
			if code == "" {
				code = "42601"
			}
			if code != pgErr.Code {
				t.Errorf("omni SQLSTATE = %s, postgres = %s (%s)", code, pgErr.Code, pgErr.Message)
			}
			if pe.Message != pgErr.Message {
				t.Errorf("omni message = %q, postgres = %q", pe.Message, pgErr.Message)
			}
			if pgErr.Position > 0 {
				if want := byteOffset(rawParsePrefix+c.sql, int(pgErr.Position)-1) - len(rawParsePrefix); pe.Position != want {
					t.Errorf("omni position = %d, postgres = %d", pe.Position, want)
				}
			}
		})
	}
}

// rawParsePrefix runs before the statement a raw-parse probe sends.
const rawParsePrefix = "SELECT 1/0;\n"

// rawParse returns the error PostgreSQL's raw parser refuses a statement
// with, or nil when it accepts it. The statement follows SELECT 1/0 in one
// simple query: the server parses the whole string before it runs any of
// it, so a division by zero means the statement parsed, and any other
// error is the parser's, whatever its SQLSTATE.
func (o *firstSetOracle) rawParse(t *testing.T, sqlStr string) *pgconn.PgError {
	t.Helper()
	ctx, cancel := context.WithTimeout(o.ctx, 5*time.Second)
	defer cancel()
	_, err := o.db.ExecContext(ctx, rawParsePrefix+sqlStr)
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		t.Fatalf("raw parse probe: %v", err)
	}
	if pgErr.Code == "22012" {
		return nil
	}
	return pgErr
}

// byteOffset converts a character offset into s, as the server reports
// positions, to a byte offset.
func byteOffset(s string, chars int) int {
	off := 0
	for i := 0; i < chars && off < len(s); i++ {
		_, size := utf8.DecodeRuneInString(s[off:])
		off += size
	}
	return off
}
