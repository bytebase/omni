package parser

import (
	"errors"
	"strings"
	"testing"

	gomysql "github.com/go-sql-driver/mysql"

	"github.com/bytebase/omni/tidb/ast"
)

// TestBatchKeywordsStillIdentifiers verifies that introducing the BATCH/DRY/RUN
// keyword tokens does not regress their use as identifiers. Upstream classifies
// all three as non-reserved (BATCH/DRY/RUN are TiDBKeyword), so they must remain
// usable as table, column, and qualifier names.
func TestBatchKeywordsStillIdentifiers(t *testing.T) {
	cases := []string{
		"SELECT batch FROM t",
		"SELECT dry, run FROM t",
		"SELECT * FROM t WHERE run = 1",
		"SELECT * FROM t WHERE dry > 0 AND batch < 10",
		"CREATE TABLE batch (id INT)",
		"CREATE TABLE run (dry INT, batch VARCHAR(10))",
		"INSERT INTO run (dry) VALUES (1)",
		"SELECT batch.id FROM batch",
		"UPDATE batch SET run = 1 WHERE dry = 0",
	}
	for _, sql := range cases {
		t.Run(sql, func(t *testing.T) {
			ParseAndCheck(t, sql)
		})
	}
}

// batchCase is a BATCH grammar acceptance case, shared by the pure-parser test
// and the TiDB container lockstep so both assert against the same corpus.
type batchCase struct {
	name       string
	sql        string
	wantAccept bool
}

// batchCases enumerates BATCH grammar positives and sibling-arm negatives.
// Negatives are derived from the upstream production
// "BATCH" OptionalShardColumn "LIMIT" NUM DryRunOptions ShardableStmt:
// every decorator a sibling arm accepts that BATCH's arm rejects gets a case.
var batchCases = []batchCase{
	// Positives.
	{"shard col + delete", "BATCH ON id LIMIT 5000 DELETE FROM t WHERE 1=1", true},
	{"no shard col + delete", "BATCH LIMIT 5000 DELETE FROM t WHERE 1=1", true},
	{"dry run (split dml)", "BATCH ON id LIMIT 1000 DRY RUN DELETE FROM t WHERE 1=1", true},
	{"dry run query", "BATCH ON id LIMIT 1000 DRY RUN QUERY DELETE FROM t WHERE 1=1", true},
	{"update", "BATCH ON id LIMIT 100 UPDATE t SET x=1 WHERE 1=1", true},
	{"insert select", "BATCH ON id LIMIT 100 INSERT INTO t2 SELECT * FROM t WHERE 1=1", true},
	{"insert ignore", "BATCH ON id LIMIT 100 INSERT IGNORE INTO t2 SELECT * FROM t WHERE 1=1", true},
	{"replace select", "BATCH ON id LIMIT 100 REPLACE INTO t2 SELECT * FROM t WHERE 1=1", true},
	{"qualified shard col", "BATCH ON t.id LIMIT 100 DELETE FROM t WHERE 1=1", true},

	// Negatives — sibling-arm decorators BATCH rejects.
	{"missing LIMIT", "BATCH DELETE FROM t WHERE 1=1", false},
	{"ON without column", "BATCH ON LIMIT 100 DELETE FROM t WHERE 1=1", false},
	{"LIMIT without NUM", "BATCH ON id LIMIT DELETE FROM t WHERE 1=1", false},
	{"LIMIT expression not NUM", "BATCH ON id LIMIT (1+1) DELETE FROM t WHERE 1=1", false},
	{"LIMIT negative", "BATCH ON id LIMIT -100 DELETE FROM t WHERE 1=1", false},
	{"SELECT not shardable", "BATCH ON id LIMIT 100 SELECT * FROM t", false},
	{"DDL not shardable", "BATCH ON id LIMIT 100 CREATE TABLE x (id INT)", false},
	{"DRY without RUN", "BATCH ON id LIMIT 100 DRY DELETE FROM t WHERE 1=1", false},
	{"DRY RUN junk modifier", "BATCH ON id LIMIT 100 DRY RUN FOO DELETE FROM t WHERE 1=1", false},
	{"no DML statement", "BATCH ON id LIMIT 100", false},
	{"bare BATCH", "BATCH", false},
}

// TestBatchParse locks BATCH grammar acceptance in the pure parser (no Docker).
func TestBatchParse(t *testing.T) {
	for _, c := range batchCases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Parse(c.sql)
			accepts := err == nil
			if accepts != c.wantAccept {
				t.Errorf("Parse(%q): accepts=%v, want %v (err=%v)", c.sql, accepts, c.wantAccept, err)
			}
		})
	}
}

// TestBatchAST locks the semantic mapping of the parsed BatchStmt: dry-run mode,
// shard column shape, limit value, and DML node type (incl. REPLACE unification).
func TestBatchAST(t *testing.T) {
	mustBatch := func(t *testing.T, sql string) *ast.BatchStmt {
		t.Helper()
		l := ParseAndCheck(t, sql)
		b, ok := l.Items[0].(*ast.BatchStmt)
		if !ok {
			t.Fatalf("Parse(%q): got %T, want *ast.BatchStmt", sql, l.Items[0])
		}
		return b
	}

	t.Run("dry run query → Query", func(t *testing.T) {
		b := mustBatch(t, "BATCH ON id LIMIT 1000 DRY RUN QUERY DELETE FROM t WHERE 1=1")
		if b.DryRun != ast.BatchDryRunQuery {
			t.Errorf("DryRun=%d, want BatchDryRunQuery(%d)", b.DryRun, ast.BatchDryRunQuery)
		}
		if _, ok := b.DML.(*ast.DeleteStmt); !ok {
			t.Errorf("DML=%T, want *ast.DeleteStmt", b.DML)
		}
	})

	t.Run("dry run → SplitDML", func(t *testing.T) {
		b := mustBatch(t, "BATCH ON id LIMIT 1000 DRY RUN DELETE FROM t WHERE 1=1")
		if b.DryRun != ast.BatchDryRunSplitDML {
			t.Errorf("DryRun=%d, want BatchDryRunSplitDML(%d)", b.DryRun, ast.BatchDryRunSplitDML)
		}
	})

	t.Run("no dry run → None", func(t *testing.T) {
		b := mustBatch(t, "BATCH ON id LIMIT 5000 DELETE FROM t WHERE 1=1")
		if b.DryRun != ast.BatchDryRunNone {
			t.Errorf("DryRun=%d, want BatchDryRunNone(%d)", b.DryRun, ast.BatchDryRunNone)
		}
		if b.Limit != 5000 {
			t.Errorf("Limit=%d, want 5000", b.Limit)
		}
		if b.ShardColumn == nil || b.ShardColumn.Column != "id" {
			t.Errorf("ShardColumn=%+v, want {Column:id}", b.ShardColumn)
		}
	})

	t.Run("no shard column → nil", func(t *testing.T) {
		b := mustBatch(t, "BATCH LIMIT 5000 DELETE FROM t WHERE 1=1")
		if b.ShardColumn != nil {
			t.Errorf("ShardColumn=%+v, want nil", b.ShardColumn)
		}
	})

	t.Run("qualified shard column", func(t *testing.T) {
		b := mustBatch(t, "BATCH ON t.id LIMIT 100 DELETE FROM t WHERE 1=1")
		if b.ShardColumn == nil || b.ShardColumn.Table != "t" || b.ShardColumn.Column != "id" {
			t.Errorf("ShardColumn=%+v, want {Table:t Column:id}", b.ShardColumn)
		}
	})

	t.Run("replace → InsertStmt IsReplace", func(t *testing.T) {
		b := mustBatch(t, "BATCH ON id LIMIT 100 REPLACE INTO t2 SELECT * FROM t")
		ins, ok := b.DML.(*ast.InsertStmt)
		if !ok {
			t.Fatalf("DML=%T, want *ast.InsertStmt", b.DML)
		}
		if !ins.IsReplace {
			t.Errorf("InsertStmt.IsReplace=false, want true for REPLACE")
		}
	})
}

// TestBatchDryRunEnumValues locks the numeric enum values to mirror pingcap's
// ast (NoDryRun=0, DryRunQuery=1, DryRunSplitDml=2). The omni parser and any
// downstream consumer rely on these exact values.
func TestBatchDryRunEnumValues(t *testing.T) {
	if ast.BatchDryRunNone != 0 {
		t.Errorf("BatchDryRunNone=%d, want 0", ast.BatchDryRunNone)
	}
	if ast.BatchDryRunQuery != 1 {
		t.Errorf("BatchDryRunQuery=%d, want 1", ast.BatchDryRunQuery)
	}
	if ast.BatchDryRunSplitDML != 2 {
		t.Errorf("BatchDryRunSplitDML=%d, want 2", ast.BatchDryRunSplitDML)
	}
}

// TestBatchSerialize verifies the AST serialization (outfuncs writeBatchStmt)
// for both DRY RUN variants and the default mode. omni has no statement-level
// SQL deparse (the deparse package handles only expressions and SELECT), so
// parse→NodeToString→inspect is the round-trip analog for BatchStmt.
func TestBatchSerialize(t *testing.T) {
	cases := []struct {
		sql          string
		wantContains string
	}{
		{"BATCH ON id LIMIT 1000 DRY RUN DELETE FROM t WHERE 1=1", ":dry_run split_dml"},
		{"BATCH ON id LIMIT 1000 DRY RUN QUERY UPDATE t SET x=1 WHERE 1=1", ":dry_run query"},
		{"BATCH ON id LIMIT 5000 DELETE FROM t WHERE 1=1", ":limit 5000"},
		{"BATCH ON id LIMIT 100 REPLACE INTO t2 SELECT * FROM t", ":replace true"},
	}
	for _, c := range cases {
		got := ast.NodeToString(ParseAndCheck(t, c.sql).Items[0])
		if !strings.Contains(got, c.wantContains) {
			t.Errorf("NodeToString(%q) = %q, want substring %q", c.sql, got, c.wantContains)
		}
	}

	// Default (no DRY RUN) mode must omit the :dry_run field entirely.
	noDryRun := ast.NodeToString(ParseAndCheck(t, "BATCH ON id LIMIT 5000 DELETE FROM t WHERE 1=1").Items[0])
	if strings.Contains(noDryRun, ":dry_run") {
		t.Errorf("default mode should omit :dry_run, got %q", noDryRun)
	}
}

// tidbRejectedSyntax reports whether a TiDB execution error is a parse error
// (ER_PARSE_ERROR, 1064). Any other error (e.g. 1146 table-not-found) means
// TiDB parsed the statement, so the syntax is accepted.
func tidbRejectedSyntax(err error) bool {
	if err == nil {
		return false
	}
	var myErr *gomysql.MySQLError
	if errors.As(err, &myErr) {
		return myErr.Number == 1064
	}
	return false
}

// TestBatchTiDBOracle lockstep-verifies every batchCase against real TiDB
// v8.5.5: our parser's accept/reject must match TiDB's syntax acceptance.
// Skips under -short (CI) and when the container is unavailable.
func TestBatchTiDBOracle(t *testing.T) {
	tc := startTiDB(t)

	tc.db.ExecContext(tc.ctx, "CREATE DATABASE IF NOT EXISTS omni_batch_test")
	tc.db.ExecContext(tc.ctx, "USE omni_batch_test")
	defer tc.db.ExecContext(tc.ctx, "DROP DATABASE IF EXISTS omni_batch_test")

	for _, c := range batchCases {
		t.Run(c.name, func(t *testing.T) {
			_, omniErr := Parse(c.sql)
			omniAccepts := omniErr == nil

			tidbAccepts := !tidbRejectedSyntax(tc.canExecute(c.sql))

			if omniAccepts != tidbAccepts {
				t.Errorf("MISMATCH %q: omni accepts=%v (err=%v), TiDB accepts=%v",
					c.sql, omniAccepts, omniErr, tidbAccepts)
			}
			if omniAccepts != c.wantAccept {
				t.Errorf("omni %q: accepts=%v, want %v (err=%v)", c.sql, omniAccepts, c.wantAccept, omniErr)
			}
		})
	}
}
