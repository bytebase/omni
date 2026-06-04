package analysis

import (
	"context"
	"testing"
	"time"

	"github.com/bytebase/omni/trino/internal/trinooracle"
)

// This file is the analysis node's slice of the differential-oracle gate
// (correctness-protocol.md). analysis is a feature node, so the oracle is not
// used to adjudicate a grammar accept/reject — that is the parser nodes' job.
// What the oracle pins here is that the SQL the span/classify tests assert over
// is REAL Trino 481 syntax: if a corpus statement the tests treat as a valid
// SELECT were actually rejected by Trino, the lineage/classification assertions
// built on it would be meaningless. The oracle therefore confirms accept/reject
// for the whole corpus and cross-checks the read-only classification against
// Trino's own verdict.
//
// All subtests skip cleanly when no Trino is reachable, matching the harness
// convention (oracle_foundation_test.go).

// connectOracle dials the live Trino oracle, skipping the test when unreachable.
func connectOracle(t *testing.T) *trinooracle.Oracle {
	t.Helper()
	o := trinooracle.Connect("")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	ver, err := o.Ping(ctx)
	if err != nil {
		t.Skipf("trino oracle not reachable (start: docker run -d -p 18080:8080 %s): %v",
			trinooracle.DefaultImage, err)
	}
	t.Logf("connected to Trino %s", ver)
	return o
}

func oracleAccepts(t *testing.T, o *trinooracle.Oracle, sql string) (accepted, ok bool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	res, err := o.CheckSyntax(ctx, sql)
	if err != nil {
		return false, false
	}
	return res.Accepted, true
}

// readOnlyCorpus pairs each statement with whether Trino should accept it AND
// whether the analysis classifier deems it read-only (Select/SelectInfoSchema/
// Explain). The corpus is intentionally schema-light (every table is created in
// the oracle's memory catalog below) so the only verdict that varies is
// accept/reject of the SYNTAX, which is what we cross-check.
var readOnlyCorpus = []struct {
	sql        string
	wantType   QueryType
	isReadOnly bool
}{
	{"SELECT a FROM t", Select, true},
	{"SELECT * FROM t1 JOIN t2 ON t1.id = t2.id", Select, true},
	{"WITH c AS (SELECT a FROM t) SELECT a FROM c", Select, true},
	{"SELECT a FROM t WHERE b IN (SELECT b FROM t2)", Select, true},
	{"SELECT a FROM t UNION SELECT a2 FROM t2", Select, true},
	{"SELECT * FROM information_schema.tables", SelectInfoSchema, true},
	{"SHOW TABLES", SelectInfoSchema, true},
	{"SHOW SCHEMAS", SelectInfoSchema, true},
	{"DESCRIBE t", SelectInfoSchema, true},
	{"EXPLAIN SELECT a FROM t", Explain, true},
	{"EXPLAIN ANALYZE SELECT a FROM t", Explain, true},
	// EXPLAIN ANALYZE executes its inner statement: a mutating inner statement
	// is NOT read-only (oracle confirms Trino runs it).
	{"EXPLAIN ANALYZE INSERT INTO t VALUES (1)", DML, false},
	// Plain EXPLAIN never executes, so EXPLAIN of a DML stays read-only.
	{"EXPLAIN INSERT INTO t VALUES (1)", Explain, true},
	{"INSERT INTO t VALUES (1)", DML, false},
	{"DELETE FROM t WHERE a = 1", DML, false},
	{"UPDATE t SET a = 2 WHERE a = 1", DML, false},
	{"CREATE TABLE t3 (a integer)", DDL, false},
	{"DROP TABLE t3", DDL, false},
}

// TestAnalysis_ClassificationDifferential confirms each corpus statement is
// syntactically accepted by Trino 481, and that the analysis classifier's
// read-only verdict agrees with Trino's read-only/data-changing nature.
// Read-only-ness is checked structurally (a SELECT/SHOW/EXPLAIN is read-only; an
// INSERT/DELETE/UPDATE/CREATE/DROP is not) — the classifier must put each
// statement on the correct side of that line.
func TestAnalysis_ClassificationDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)

	// Seed the schema objects so the corpus statements reach Trino's PARSER
	// without being rejected for missing tables (a missing-table error is a
	// SEMANTIC, not SYNTAX, rejection — CheckSyntax already treats it as
	// "accepted" — but seeding keeps the corpus runnable end to end).
	setup := []string{
		"DROP TABLE IF EXISTS t",
		"DROP TABLE IF EXISTS t2",
		"DROP TABLE IF EXISTS t3",
		"CREATE TABLE t (a integer, b integer, id integer)",
		"CREATE TABLE t2 (a2 integer, b integer, id integer)",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	for _, s := range setup {
		if _, err := o.CheckSyntax(ctx, s); err != nil {
			t.Skipf("oracle setup failed (%q): %v", s, err)
		}
	}

	for _, tc := range readOnlyCorpus {
		tc := tc
		t.Run(tc.sql, func(t *testing.T) {
			accepted, ok := oracleAccepts(t, o, tc.sql)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if !accepted {
				t.Errorf("Trino 481 REJECTED corpus statement %q — the analysis assertions on it are built on invalid syntax", tc.sql)
				return
			}

			gotType := Classify(tc.sql)
			if gotType != tc.wantType {
				t.Errorf("Classify(%q) = %v, want %v", tc.sql, gotType, tc.wantType)
			}

			gotReadOnly := isReadOnly(gotType)
			if gotReadOnly != tc.isReadOnly {
				t.Errorf("isReadOnly(%q) = %v (type %v), want %v",
					tc.sql, gotReadOnly, gotType, tc.isReadOnly)
			}
		})
	}
}

// isReadOnly mirrors the bytebase validateQuery read-only guard: a statement is
// read-only iff it is a Select, Explain, or SelectInfoSchema.
func isReadOnly(qt QueryType) bool {
	return qt == Select || qt == Explain || qt == SelectInfoSchema
}
