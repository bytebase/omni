// Package testcontainer is how redshift's driver tests reach a cluster: through
// a tape (enginetest/tape) that replays a pinned cluster's answers, or, in a
// recording (OMNI_RECORD=1), through the cluster REDSHIFT_DSN names.
//
// Redshift is proprietary and has no image omni can start, so nothing is
// started and a replay reaches no server. A recording reads REDSHIFT_DSN and
// refuses to run without it: a recording that skipped would write an empty tape
// and call it a pass. The endpoint and the login never enter a tape; only the
// statements and the answers do.
//
// The pin row names the cluster's Redshift release, and the digest is the
// SHA-256 of SELECT version()'s answer, which is what an answer's text depends
// on. A cluster of another release is refused rather than replayed.
package testcontainer

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/stdlib"

	"github.com/bytebase/omni/enginetest"
	"github.com/bytebase/omni/enginetest/tape"
)

// Pins is redshift/testdata/pins.tsv.
var Pins = func() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "testdata", "pins.tsv")
}()

// Driver is the pin row of the cluster the driver's tapes were recorded on.
const Driver = "driver"

// Database is the database the tapes were recorded against. A test constructs
// its driver with it instead of reading a DSN, so a replay needs no
// environment.
const Database = "dev"

// server is the tape server the driver's recordings are written to.
var server = &tape.Server{Pins: Pins, Live: live, Codec: Codec{}}

// Open returns t's connection: a replay of its tape, or, in a recording, the
// pinned cluster with every exchange written down. The pool holds one
// connection, so session state (a search path, a temporary table) holds as it
// would on one live connection.
func Open(t testing.TB) *sql.DB {
	t.Helper()
	return tape.Open(t, server, Driver)
}

// releasePattern reads the Redshift release out of the version string, which
// ends in ", Redshift 1.0.539507".
var releasePattern = regexp.MustCompile(`Redshift (\d+\.\d+\.\d+)`)

// live connects to the cluster REDSHIFT_DSN names and proves it is the one the
// pin row records. It fails rather than skips when the variable is unset: a
// recording without a cluster must not produce an empty tape.
func live(t testing.TB, pin enginetest.Pin) driver.Connector {
	t.Helper()
	dsn := strings.TrimSpace(os.Getenv("REDSHIFT_DSN"))
	if dsn == "" {
		t.Fatalf("testcontainer: recording needs REDSHIFT_DSN, and OMNI_RECORD=1 was set")
	}
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatalf("testcontainer: REDSHIFT_DSN does not parse: %v", err)
	}
	// The tape holds each statement with its arguments, so the client must not
	// prepare statements of its own.
	cfg.DefaultQueryExecMode = pgx.QueryExecModeExec
	cfg.StatementCacheCapacity = 0

	ctx := context.Background()
	conn, err := pgx.ConnectConfig(ctx, cfg)
	if err != nil {
		t.Fatalf("testcontainer: connect: %v", err)
	}
	defer func() { _ = conn.Close(ctx) }()
	var version string
	if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
		t.Fatalf("testcontainer: SELECT version(): %v", err)
	}
	matches := releasePattern.FindStringSubmatch(version)
	if len(matches) != 2 {
		t.Fatalf("testcontainer: version %q names no Redshift release", version)
	}
	want := strings.TrimPrefix(pin.Image, "aws:redshift:")
	if matches[1] != want {
		t.Fatalf("testcontainer: the cluster runs Redshift %s, but pins.tsv pins %q", matches[1], pin.Image)
	}
	sum := sha256.Sum256([]byte(version))
	if got := "sha256:" + hex.EncodeToString(sum[:]); got != pin.Digest {
		t.Fatalf("testcontainer: SELECT version() hashes to %s, but pins.tsv pins %s", got, pin.Digest)
	}
	return stdlib.GetConnector(*cfg)
}

// errorFields are the fields a tape holds after the first line, in order.
var errorFields = []struct {
	name string
	get  func(*pgconn.PgError) *string
}{
	{"severity", func(e *pgconn.PgError) *string { return &e.Severity }},
	{"detail", func(e *pgconn.PgError) *string { return &e.Detail }},
	{"hint", func(e *pgconn.PgError) *string { return &e.Hint }},
	{"where", func(e *pgconn.PgError) *string { return &e.Where }},
	{"schema", func(e *pgconn.PgError) *string { return &e.SchemaName }},
	{"table", func(e *pgconn.PgError) *string { return &e.TableName }},
	{"column", func(e *pgconn.PgError) *string { return &e.ColumnName }},
	{"data type", func(e *pgconn.PgError) *string { return &e.DataTypeName }},
	{"constraint", func(e *pgconn.PgError) *string { return &e.ConstraintName }},
}

// plainSeverity is the severity of nearly every error, which a tape leaves out.
const plainSeverity = "ERROR"

var (
	escaper   = strings.NewReplacer(`\`, `\\`, "\n", `\n`, "\r", `\r`)
	unescaper = strings.NewReplacer(`\\`, `\`, `\n`, "\n", `\r`, "\r")
)

// Codec renders pgx's *pgconn.PgError as its SQLSTATE and message on one line,
// then a line for each other field the server sent, so a replayed test's
// errors.As finds the same type and the same text.
type Codec struct{}

func (Codec) EncodeError(err error) (string, bool) {
	var pe *pgconn.PgError
	if !errors.As(err, &pe) {
		return "", false
	}
	var b strings.Builder
	b.WriteString(pe.Code + " " + escaper.Replace(pe.Message))
	if pe.Position != 0 {
		fmt.Fprintf(&b, "\nposition: %d", pe.Position)
	}
	for _, f := range errorFields {
		v := *f.get(pe)
		if v == "" || f.name == "severity" && v == plainSeverity {
			continue
		}
		b.WriteString("\n" + f.name + ": " + escaper.Replace(v))
	}
	return b.String(), true
}

func (Codec) DecodeError(text string) error {
	lines := strings.Split(text, "\n")
	code, message, ok := strings.Cut(lines[0], " ")
	if !ok || len(code) != 5 {
		return fmt.Errorf("tape: a Redshift error %q", text)
	}
	pe := &pgconn.PgError{Severity: plainSeverity, Code: code, Message: unescaper.Replace(message)}
	for _, line := range lines[1:] {
		name, v, ok := strings.Cut(line, ": ")
		if !ok {
			return fmt.Errorf("tape: a Redshift error field %q", line)
		}
		if name == "position" {
			n, err := strconv.ParseInt(v, 10, 32)
			if err != nil {
				return fmt.Errorf("tape: a Redshift error position %q", v)
			}
			pe.Position = int32(n)
			continue
		}
		found := false
		for _, f := range errorFields {
			if f.name == name {
				*f.get(pe) = unescaper.Replace(v)
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("tape: a Redshift error field %q", line)
		}
	}
	pe.SeverityUnlocalized = pe.Severity
	return pe
}
