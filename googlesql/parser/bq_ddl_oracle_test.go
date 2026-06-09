//go:build googlesql_oracle

// Differential / triangulation gate for the parser-ddl-bigquery node against the
// live Cloud Spanner emulator. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestBQDDL
//
// ⚠ Every form this node owns is BigQuery-ONLY at the GoogleSQL union level
// (oracle.md): the Spanner emulator is NOT an authoritative oracle for them. So
// this is NOT a plain accept/reject differential — it is a TRIANGULATION GUARD.
// For each fixture we record:
//   - omni's verdict (the union parser MUST accept these documented BigQuery
//     forms — triangulated against the legacy GoogleSQLParser.g4 + truth1), and
//   - the live emulator's verdict, asserted to equal the EMPIRICALLY-OBSERVED
//     value (probed 2026-06-05). Most BigQuery-only forms hard-syntax-REJECT on
//     Spanner (a non-authoritative reject the union parser correctly overrides);
//     a NARROW few (bare CREATE FUNCTION … AS (expr), DROP FUNCTION, CREATE
//     VECTOR INDEX) ACCEPT because Spanner happens to share that shape.
//
// The guard's value: if an emulator bump ever changes a recorded oracle verdict
// (e.g. Spanner starts parsing a form it used to reject), this test fails and
// forces a re-review of the triangulation for that form, rather than silently
// trusting a now-stale assumption. omni's own accept/reject correctness is
// proven by the hand-written unit tests (bq_*_test.go); this file ties those to
// the live oracle so the BigQuery-vs-Spanner boundary stays honest.
package parser

import (
	"os"
	"testing"
)

// bqOracleExpectation is the empirically-observed live-emulator verdict for a
// BigQuery-only form, plus omni's required verdict (always accept — the union
// parser must parse the documented form).
type bqOracleExpectation struct {
	sql         string
	oracle      string // "accept" | "reject" — the OBSERVED Spanner-emulator verdict (2026-06-05)
	omniAccepts bool   // omni's required verdict (the union parser)
	note        string
}

var bqOracleFixtures = []bqOracleExpectation{
	// ---- hard BigQuery-only (Spanner has NO such grammar) → oracle REJECTs, omni accepts ----
	{"CREATE OR REPLACE TABLE FUNCTION ds.f(y INT64) AS SELECT 1 AS n", "reject", true, "TABLE FUNCTION not in Spanner"},
	{"CREATE PROCEDURE ds.p(IN tbl STRING) BEGIN SELECT 1; END", "reject", true, "PROCEDURE not in Spanner"},
	{"CREATE MATERIALIZED VIEW ds.mv AS SELECT 1 AS n", "reject", true, "MATERIALIZED VIEW not in Spanner"},
	{"CREATE APPROX VIEW ds.v AS SELECT 1 AS n", "reject", true, "APPROX VIEW not in Spanner"},
	{"CREATE SEARCH INDEX my_index ON ds.t(ALL COLUMNS)", "reject", true, "SEARCH INDEX ALL COLUMNS not in Spanner"},
	{"CREATE SNAPSHOT TABLE ds.s CLONE ds.t", "reject", true, "SNAPSHOT not in Spanner"},
	{"CREATE ROW ACCESS POLICY p ON ds.t GRANT TO ('user:a@b.com') FILTER USING (TRUE)", "reject", true, "ROW ACCESS POLICY not in Spanner"},
	{"CREATE CAPACITY `p.r.c` OPTIONS (slot_count = 100)", "reject", true, "generic-entity CAPACITY not in Spanner"},
	{"CREATE RESERVATION `p.r` OPTIONS (slot_capacity = 100)", "reject", true, "generic-entity RESERVATION not in Spanner"},
	{"CREATE TEMP FUNCTION f(x INT64) AS (x)", "reject", true, "TEMP function not in Spanner"},
	{"CREATE AGGREGATE FUNCTION ds.f(x FLOAT64) AS (SUM(x))", "reject", true, "AGGREGATE function not in Spanner"},
	{"CREATE FUNCTION ds.f(x FLOAT64) RETURNS FLOAT64 LANGUAGE js AS 'return x'", "reject", true, "LANGUAGE js not in Spanner"},
	{"CREATE FUNCTION ds.f(x INT64) RETURNS INT64 DETERMINISTIC AS (x)", "reject", true, "DETERMINISTIC not in Spanner"},
	{"DROP TABLE FUNCTION ds.f", "reject", true, "DROP TABLE FUNCTION not in Spanner"},
	{"DROP VECTOR INDEX idx ON ds.t", "reject", true, "DROP VECTOR INDEX … ON not in Spanner"},
	{"DROP SEARCH INDEX idx ON ds.t", "reject", true, "DROP SEARCH INDEX … ON not in Spanner"},
	{"ALTER VECTOR INDEX idx ON ds.t REBUILD", "reject", true, "ALTER VECTOR INDEX REBUILD not in Spanner (nor in legacy .g4 — flagged divergence)"},
	// ALTER SEARCH INDEX … {ADD|DROP} STORED COLUMN is a DOCUMENTED Spanner form
	// (DDL-049, Spanner DDL reference) the union parser must accept, but the live
	// emulator (sha256:caf1bd24) REJECTs it at the syntax layer ("Encountered
	// 'SEARCH' while parsing: alter_statement") — its grammar is a subset that lags
	// the docs. Non-authoritative reject the union parser correctly overrides
	// (triangulated against the Spanner DDL reference + legacy .g4, which has
	// INDEX but not SEARCH INDEX in schema_object_kind). Divergence #203.
	{"ALTER SEARCH INDEX idx ADD STORED COLUMN c", "reject", true, "ALTER SEARCH INDEX ADD STORED COLUMN: documented Spanner form (DDL-049) the emulator rejects — non-authoritative (#203)"},
	{"ALTER SEARCH INDEX idx DROP STORED COLUMN c", "reject", true, "ALTER SEARCH INDEX DROP STORED COLUMN: documented Spanner form (DDL-049) the emulator rejects — non-authoritative (#203)"},

	// ---- narrow shared shapes → oracle ACCEPTs (semantic), omni accepts ----
	{"CREATE FUNCTION ds.f(x INT64) AS (x + 1)", "accept", true, "bare CREATE FUNCTION … AS (expr) is shared with Spanner"},
	{"CREATE OR REPLACE FUNCTION ds.f(x INT64) RETURNS INT64 AS (x)", "accept", true, "shared scalar SQL function"},
	{"CREATE VECTOR INDEX my_index ON ds.t(embedding) OPTIONS(distance_type='COSINE')", "accept", true, "Spanner has a vector-index grammar for this shape"},
	{"DROP FUNCTION ds.f", "accept", true, "DROP FUNCTION is shared with Spanner"},
	{"DROP FUNCTION IF EXISTS ds.f", "accept", true, "DROP FUNCTION IF EXISTS shared"},
}

func TestBQDDLTriangulation(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live triangulation guard")
	}
	h := newDDLHarness(t)
	defer h.close()

	for _, fx := range bqOracleFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) {
			// 1. Live oracle verdict.
			v := h.verdict(t, fx.sql)
			switch v.Verdict {
			case "accept", "reject":
				// proceed
			case "error":
				t.Fatalf("oracle returned ERROR (no verdict) for %q: reason=%s code=%s msg=%s",
					fx.sql, v.Reason, v.Code, v.Message)
			default:
				t.Fatalf("unexpected harness verdict %q for %q", v.Verdict, fx.sql)
			}

			// The recorded oracle expectation MUST still hold. If it drifts, the
			// emulator changed (e.g. Spanner added the form) and the triangulation
			// for this fixture needs a fresh look — do NOT silently update.
			if v.Verdict != fx.oracle {
				t.Fatalf("oracle verdict drift on %q: recorded=%q now=%q (%s: %s) [%s]\n"+
					"the Spanner emulator changed; re-review the BigQuery-vs-Spanner triangulation for this form",
					fx.sql, fx.oracle, v.Verdict, v.Reason, v.Message, fx.note)
			}

			// 2. omni's verdict — the union parser must match omniAccepts. For the
			// BigQuery-only forms that is ALWAYS accept (the documented form is valid
			// GoogleSQL); the Spanner reject above is non-authoritative.
			_, errs := Parse(fx.sql)
			omniAccepts := len(errs) == 0
			if omniAccepts != fx.omniAccepts {
				t.Errorf("omni verdict on %q: accepts=%v, want %v; errs=%v [%s]",
					fx.sql, omniAccepts, fx.omniAccepts, errs, fx.note)
			}
		})
	}
}
