//go:build oracle

package splittest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	tcpg "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github.com/bytebase/omni/pg"
	"github.com/bytebase/omni/pg/internal/metacmd"
)

// S3 server-side differential (v1.2 S line, framework G4 proposition).
//
// Local runs need the build tag or the tests are silently absent:
//
//	S3DIFF_N=300 go test -tags=oracle ./pg/splittest/ -run TestS3
//
// The server is the ultimate splitting authority: a multi-statement
// script sent as ONE simple-query message is split by PostgreSQL
// itself, and any error it reports carries a 1-based byte Position
// into the script we sent. That position must land inside the segment
// our splitter produced for that statement — if our boundaries drift,
// the server's own error positions expose it. Scripts that succeed
// whole must also succeed segment-by-segment (equivalence is defined
// modulo the implicit transaction of the simple-query batch, so the
// segment side runs inside an explicit transaction; contract v1.2 #5).

var (
	s3Once    sync.Once
	s3DB      *sql.DB
	s3InitErr error
)

func s3Conn(t *testing.T) *sql.DB {
	t.Helper()
	if testing.Short() {
		t.Skip("skipping S3 differential in short mode")
	}
	s3Once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
		defer cancel()
		container, err := tcpg.Run(ctx, "postgres:17-alpine",
			tcpg.WithDatabase("s3diff"),
			tcpg.WithUsername("postgres"),
			tcpg.WithPassword("test"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2)),
		)
		if err != nil {
			s3InitErr = fmt.Errorf("start PG container: %w", err)
			return
		}
		connStr, err := container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			s3InitErr = fmt.Errorf("connection string: %w", err)
			return
		}
		// Force the simple query protocol explicitly. The whole-script
		// oracle depends on PostgreSQL splitting a multi-command string
		// itself; pgx's default (extended, cache_statement) prepares
		// single statements and would not exercise the server splitter.
		// Do not rely on pgx's arg-less heuristic — pin it.
		if strings.Contains(connStr, "?") {
			connStr += "&default_query_exec_mode=simple_protocol"
		} else {
			connStr += "?default_query_exec_mode=simple_protocol"
		}
		db, err := sql.Open("pgx", connStr)
		if err != nil {
			s3InitErr = fmt.Errorf("open: %w", err)
			return
		}
		if err := db.PingContext(ctx); err != nil {
			s3InitErr = fmt.Errorf("ping: %w", err)
			return
		}
		s3DB = db
	})
	if s3InitErr != nil {
		// In CI a startup failure must fail the job, not silently turn
		// the differential into a green no-op.
		if os.Getenv("CI") != "" {
			t.Fatalf("PG container required in CI: %v", s3InitErr)
		}
		t.Skipf("PG container unavailable: %v", s3InitErr)
	}
	return s3DB
}

// wholeScriptError sends the script as one simple-query batch and
// returns the server's first error (nil on full success) with its
// 0-based byte position (-1 when the server reports none).
func wholeScriptError(ctx context.Context, db *sql.DB, script string) (*pgconn.PgError, int, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, -1, err
	}
	defer conn.Close()
	_, execErr := conn.ExecContext(ctx, script)
	if execErr == nil {
		return nil, -1, nil
	}
	var pgErr *pgconn.PgError
	if !errors.As(execErr, &pgErr) {
		return nil, -1, fmt.Errorf("non-PG error from whole-script exec: %w", execErr)
	}
	pos := -1
	if pgErr.Position > 0 {
		pos = int(pgErr.Position) - 1 // server position is 1-based
	}
	return pgErr, pos, nil
}

// sqlForServer returns the SQL a segment would actually send to the
// database: metacommand lines are trivia in the real pipeline (parsed
// out before execution), so they must be stripped before handing the
// text to PostgreSQL, which has no notion of psql backslash commands.
func sqlForServer(text string) string {
	var b strings.Builder
	i := 0
	for i < len(text) {
		if metacmd.IsLineStart(text, i, i == 0 || text[i-1] == '\n') {
			i = metacmd.SkipLine(text, i)
			continue
		}
		b.WriteByte(text[i])
		i++
	}
	return b.String()
}

// segmentedRun executes the non-empty segments one by one inside an
// explicit transaction (normalizing the simple-query batch's implicit
// transaction). It returns the index (into segs) of the first failing
// segment and the error, or -1 on full success. Each segment must
// contain exactly one statement — if a splitter regression leaves a
// top-level semicolon inside a segment, executing it as a batch would
// hide the very drift this differential exists to catch, so that is a
// hard error.
func segmentedRun(t *testing.T, ctx context.Context, db *sql.DB, segs []pg.Segment) (int, *pgconn.PgError, error) {
	conn, err := db.Conn(ctx)
	if err != nil {
		return -1, nil, err
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, "BEGIN"); err != nil {
		return -1, nil, err
	}
	defer func() { _, _ = conn.ExecContext(ctx, "ROLLBACK") }()
	for i, s := range segs {
		if s.Empty() {
			continue
		}
		if sub := pg.Split(s.Text); nonEmptyCount(sub) > 1 {
			t.Errorf("segment %d contains %d statements (splitter drift): %q", i, nonEmptyCount(sub), s.Text)
		}
		if _, execErr := conn.ExecContext(ctx, sqlForServer(s.Text)); execErr != nil {
			var pgErr *pgconn.PgError
			if !errors.As(execErr, &pgErr) {
				return i, nil, fmt.Errorf("non-PG error at segment %d: %w", i, execErr)
			}
			return i, pgErr, nil
		}
	}
	return -1, nil, nil
}

// nonEmptyCount counts the non-Empty segments in a slice.
func nonEmptyCount(segs []pg.Segment) int {
	n := 0
	for _, s := range segs {
		if !s.Empty() {
			n++
		}
	}
	return n
}

// segmentContaining maps a 0-based byte position in the script to the
// index of the segment whose [ByteStart,ByteEnd) range contains it.
func segmentContaining(segs []pg.Segment, pos int) int {
	for i, s := range segs {
		if pos >= s.ByteStart && pos < s.ByteEnd {
			return i
		}
	}
	return -1
}

func s3Scale() int {
	if v := os.Getenv("S3DIFF_N"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			return p
		}
	}
	return 200
}

// TestS3Differential is the G4 gate: generated scripts dual-run against
// the server, modeling both the server's phase separation and the real
// pipeline's trivia filtering:
//
//   - The whole-script side sends the concatenation of NON-EMPTY
//     segments (the pipeline never sends trivia segments — meta-command
//     lines would be 42601 to the server but are dropped by
//     FilterEmptyStatements downstream).
//
//   - A simple-query batch is parsed IN FULL before execution, so a
//     late syntax error wins over an early semantic error. Alignment
//     rules are therefore phase-split:
//
//     A1 whole succeeds            ⇒ segmented succeeds
//     A2 whole fails 42601 at P    ⇒ the segment containing P, executed
//     alone, also fails 42601 (boundary signal from the server's
//     own parser position)
//     A3 whole fails non-42601     ⇒ the batch parsed clean, so
//     execution order matches: segmented first failure has the same
//     SQLSTATE and contains P (when the server reports a position)
func TestS3Differential(t *testing.T) {
	db := s3Conn(t)
	ctx := context.Background()
	n := s3Scale()
	r := rand.New(rand.NewSource(20260716))
	atoms := EnabledAtoms("D2", "D3", "D5")
	t.Logf("S3 differential: n=%d", n)

	for i := 0; i < n; i++ {
		script := Compose(r, atoms, 1+r.Intn(4))
		segs := pg.Split(script.SQL)

		// Pipeline-faithful batch: non-empty segments only, with each
		// segment's offset in the batch for position mapping.
		var batch strings.Builder
		type span struct{ segIdx, start, end int }
		var spans []span
		var live []pg.Segment
		for si, s := range segs {
			if s.Empty() {
				continue
			}
			// Send what the real pipeline sends: the segment with its
			// metacommand trivia lines stripped. Writing the raw
			// backslash lines would make the server 42601 at the
			// metacommand, so A2 could pass for the wrong reason.
			start := batch.Len()
			batch.WriteString(sqlForServer(s.Text))
			spans = append(spans, span{si, start, batch.Len()})
			live = append(live, s)
		}
		if len(live) == 0 {
			continue
		}
		batchInSeg := func(pos int) int {
			for _, sp := range spans {
				if pos >= sp.start && pos < sp.end {
					return sp.segIdx
				}
			}
			return -1
		}

		wholeErr, wholePos, infra := wholeScriptError(ctx, db, batch.String())
		if infra != nil {
			t.Fatalf("case %d: infra: %v", i, infra)
		}
		segIdx, segErr, infra := segmentedRun(t, ctx, db, segs)
		if infra != nil {
			t.Fatalf("case %d: infra: %v", i, infra)
		}

		switch {
		case wholeErr == nil && segErr == nil:
			// A1: agreement on success.
		case wholeErr == nil && segErr != nil:
			t.Errorf("case %d: batch succeeded but segment %d failed (%s %s)\nscript: %q",
				i, segIdx, segErr.Code, segErr.Message, script.SQL)
		case wholeErr.Code == "42601":
			// A2: syntax error — parse-phase, order-independent. The
			// segment the server points into must itself be a syntax
			// error in isolation.
			if wholePos < 0 {
				break // no position: no boundary signal to check
			}
			wantIdx := batchInSeg(wholePos)
			if wantIdx < 0 {
				t.Errorf("case %d: server syntax position %d falls outside all segments\nbatch: %q",
					i, wholePos, batch.String())
				break
			}
			// The pointed-at segment must itself be a syntax error alone.
			soloErr, _, infra := wholeScriptError(ctx, db, sqlForServer(segs[wantIdx].Text))
			if infra != nil {
				t.Fatalf("case %d: infra: %v", i, infra)
			}
			if soloErr == nil || soloErr.Code != "42601" {
				got := "success"
				if soloErr != nil {
					got = soloErr.Code
				}
				t.Errorf("case %d: server syntax error points into segment %d, but that segment alone yields %s — boundary drift\nsegment: %q\nbatch: %q",
					i, wantIdx, got, segs[wantIdx].Text, batch.String())
			}
			// Cross-check the segmented run: since a 42601 in the batch
			// means SOME segment is invalid, the segmented run must also
			// fail with 42601, and no EARLIER segment may fail first (an
			// extra boundary before wantIdx would surface here).
			if segErr == nil {
				t.Errorf("case %d: batch syntax-errored (pos %d, segment %d) but segmented run succeeded\nscript: %q",
					i, wholePos, wantIdx, script.SQL)
			} else if segErr.Code == "42601" && segIdx > wantIdx {
				t.Errorf("case %d: segmented run failed at segment %d, before the server's error segment %d — extra boundary\nscript: %q",
					i, segIdx, wantIdx, script.SQL)
			}
		default:
			// A3: batch parsed clean; execution order is preserved on
			// both sides, so first failures must align.
			if segErr == nil {
				t.Errorf("case %d: batch failed (%s at %d) but segmented run succeeded\nscript: %q",
					i, wholeErr.Code, wholePos, script.SQL)
				break
			}
			if wholeErr.Code != segErr.Code {
				t.Errorf("case %d: SQLSTATE mismatch: batch=%s segmented=%s (segment %d)\nscript: %q",
					i, wholeErr.Code, segErr.Code, segIdx, script.SQL)
			}
			if wholePos >= 0 {
				if wantIdx := batchInSeg(wholePos); wantIdx >= 0 && wantIdx != segIdx {
					t.Errorf("case %d: batch error at pos %d (segment %d) but segmented failed at segment %d\nscript: %q",
						i, wholePos, wantIdx, segIdx, script.SQL)
				}
			}
		}
	}
}

// TestS3KnownBetterWhitelist pins that every known-better input (where
// this splitter intentionally diverges from psql toward server truth)
// gets past the server's PARSER when sent whole — the claim is about
// syntax/splitting, so per the calibration bucketing only SQLSTATE
// 42601 (syntax_error) disqualifies an entry; semantic errors such as
// undefined_table mean the grammar accepted the statement.
func TestS3KnownBetterWhitelist(t *testing.T) {
	db := s3Conn(t)
	ctx := context.Background()
	for i, input := range KnownBetterThanPsql {
		wholeErr, _, infra := wholeScriptError(ctx, db, sqlForServer(input))
		if infra != nil {
			t.Fatalf("whitelist[%d]: infra: %v", i, infra)
		}
		if wholeErr != nil && wholeErr.Code == "42601" {
			t.Errorf("whitelist[%d]: server parser rejected (%s) — entry no longer qualifies as known-better: %q",
				i, wholeErr.Message, input)
		}
	}
}
