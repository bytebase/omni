package main

import (
	"context"
	"errors"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClassify(t *testing.T) {
	cases := []struct {
		name        string
		kind        string
		err         error
		wantVerdict string
		wantReason  string
	}{
		{"nil", "query", nil, "accept", "ok"},
		// query / DML: InvalidArgument + "Syntax error:" => reject; other InvalidArgument => semantic accept.
		{"query syntax", "query", status.Error(codes.InvalidArgument, `Syntax error: Unexpected identifier "SELEC" [at 1:1]`), "reject", "syntax"},
		{"query unknown table", "query", status.Error(codes.InvalidArgument, "Table not found: no_such_table [at 1:15]"), "accept", "semantic"},
		{"query unrecognized name", "query", status.Error(codes.InvalidArgument, "Unrecognized name: bogus [at 1:8]"), "accept", "semantic"},
		{"query unsupported feature", "query", status.Error(codes.InvalidArgument, "QUALIFY is not supported [at 1:18]"), "accept", "semantic"},
		{"dml syntax", "dml", status.Error(codes.InvalidArgument, "Syntax error: Unexpected keyword VALUE [at 1:25]"), "reject", "syntax"},
		{"query FailedPrecondition is not a grammar verdict", "query", status.Error(codes.FailedPrecondition, "some precondition"), "error", "infra"},
		{"query generic Internal", "query", status.Error(codes.Internal, "internal emulator crash"), "error", "infra"},
		// DDL: keyed on CODE (robust to message drift). InvalidArgument => parse reject.
		{"ddl syntax (canonical message)", "ddl", status.Error(codes.InvalidArgument, "Error parsing Spanner DDL statement: CREATE TABL ... : Syntax error on line 1"), "reject", "syntax"},
		{"ddl syntax (drifted message, code-based)", "ddl", status.Error(codes.InvalidArgument, "Error parsing Spanner DDL statement [at 1:8]: bogus"), "reject", "syntax"},
		{"ddl dup name", "ddl", status.Error(codes.FailedPrecondition, "Duplicate name in schema: T."), "accept", "semantic"},
		{"ddl missing interleave parent", "ddl", status.Error(codes.NotFound, "Table not found: p"), "accept", "semantic"},
		{"ddl bad type (ret_check quirk)", "ddl", status.Error(codes.Internal, "GOOGLESQL_RET_CHECK failure (ddl_type_conversion.cc)"), "accept", "semantic"},
		{"ddl generic Internal", "ddl", status.Error(codes.Internal, "internal emulator crash"), "error", "infra"},
		// Resource-level lifecycle misses (scratch DB/session gone) are infra, NOT a grammar verdict.
		{"database gone (notfound)", "ddl", status.Error(codes.NotFound, "Database not found: projects/p/instances/i/databases/testdb"), "error", "infra"},
		{"database already exists", "ddl", status.Error(codes.AlreadyExists, "Database already exists: testdb"), "error", "infra"},
		{"session gone (notfound)", "query", status.Error(codes.NotFound, "Session not found: ..."), "error", "infra"},
		// Transport / availability — fail closed regardless of kind.
		{"emulator down (unavailable)", "query", status.Error(codes.Unavailable, "connection refused"), "error", "infra"},
		{"deadline exceeded (grpc)", "ddl", status.Error(codes.DeadlineExceeded, "context deadline exceeded"), "error", "infra"},
		{"canceled (grpc)", "query", status.Error(codes.Canceled, "context canceled"), "error", "infra"},
		{"aborted txn retry-exhausted", "dml", status.Error(codes.Aborted, "transaction was aborted"), "error", "infra"},
		{"resource exhausted", "query", status.Error(codes.ResourceExhausted, "quota exceeded"), "error", "infra"},
		{"context.DeadlineExceeded (non-status)", "query", context.DeadlineExceeded, "error", "infra"},
		{"context.Canceled (non-status)", "dml", context.Canceled, "error", "infra"},
		{"plain dial error (non-status)", "ddl", errors.New("dial tcp 127.0.0.1:9010: connection refused"), "error", "infra"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := classify(c.kind, c.err)
			if got.Verdict != c.wantVerdict || got.Reason != c.wantReason {
				t.Fatalf("classify(%q, %q) = {%s,%s}, want {%s,%s}", c.kind, c.name, got.Verdict, got.Reason, c.wantVerdict, c.wantReason)
			}
		})
	}
}

func TestKindOf(t *testing.T) {
	cases := []struct {
		sql  string
		want string
	}{
		{"SELECT 1", "query"},
		{"  with a as (select 1) select * from a", "query"},
		{"VALUES (1,2)", "query"},
		{"INSERT INTO t VALUES (1)", "dml"},
		{"update t set x=1", "dml"},
		{"DELETE FROM t", "dml"},
		{"MERGE INTO t USING s ON t.id=s.id WHEN MATCHED THEN DELETE", "dml"},
		{"CREATE TABLE t (x INT64) PRIMARY KEY (x)", "ddl"},
		{"  alter table t add column y INT64", "ddl"},
		{"DROP TABLE t", "ddl"},
		{"GRANT SELECT ON TABLE t TO ROLE r", "ddl"},
		{"-- leading comment\nSELECT 1", "query"},
		{"# hash comment\nINSERT INTO t VALUES (1)", "dml"},
		{"/* block */ CREATE TABLE t (x INT64) PRIMARY KEY (x)", "ddl"},
		{"@{OPTIMIZER_VERSION=1} SELECT 1", "query"},
		// Documented limitation: WITH-led DML routes to the query path (verdict
		// is still preserved). Pinned so a future change is intentional.
		{"WITH x AS (SELECT 1 n) INSERT INTO t (id) SELECT n FROM x", "query"},
	}
	for _, c := range cases {
		if got := kindOf(c.sql); got != c.want {
			t.Errorf("kindOf(%q) = %q, want %q", c.sql, got, c.want)
		}
	}
}

// TestOracleLive exercises the full path against a running emulator. It skips
// when SPANNER_EMULATOR_HOST is unset or the emulator is unreachable.
func TestOracleLive(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST unset; skipping live oracle test")
	}
	ctx := context.Background()
	o, err := newOracle(ctx)
	if err != nil {
		t.Skipf("emulator bootstrap failed (is it running?): %v", err)
	}
	defer o.close()

	cases := []struct {
		sql         string
		wantVerdict string
		wantKind    string
	}{
		{"SELECT 1", "accept", "query"},
		{"SELEC 1", "reject", "query"},
		{"@@@ garbage !!", "reject", "query"},
		{"SELECT * FROM no_such_table", "accept", "query"},
		{"INSERT INTO no_such (id) VALUES (1)", "accept", "dml"},
		{"INSERT INTO t (id) VALUE (1)", "reject", "dml"},
		{"CREATE TABLE live_t (x INT64) PRIMARY KEY (x)", "accept", "ddl"},
		{"CREATE TABL live_bad (x INT64) PRIMARY KEY (x)", "reject", "ddl"},
	}
	for _, c := range cases {
		v := o.evaluate(ctx, c.sql)
		if v.Verdict != c.wantVerdict || v.Kind != c.wantKind {
			t.Errorf("evaluate(%q) = {%s,%s,%s}, want verdict=%s kind=%s (msg=%q)",
				c.sql, v.Verdict, v.Kind, v.Reason, c.wantVerdict, c.wantKind, v.Message)
		}
	}
}
