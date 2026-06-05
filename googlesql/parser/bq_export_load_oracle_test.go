//go:build googlesql_oracle

// Differential / triangulation gate for the parser-dml-ext node against the live
// Cloud Spanner emulator. Run with:
//
//	SPANNER_EMULATOR_HOST=localhost:9010 \
//	  go test -tags googlesql_oracle ./googlesql/parser/ -run TestExportLoadClone
//
// EMPIRICAL FINDING (probed 2026-06-05, this node) — these forms are NOT a
// grammar gap on the emulator. The analysis corpus (analysis.md / oracle.md)
// PRESUMED EXPORT DATA / LOAD DATA were BigQuery-only "hard syntax rejects" on
// Spanner. The live emulator says otherwise: it PARSES all four families and
// rejects them only at the RESOLVER, returning InvalidArgument with the message
// `Statement not supported: ExportDataStatement` (resp. ExportModelStatement /
// AuxLoadDataStatement / CloneDataStatement). The harness classifies that — an
// InvalidArgument with NO "Syntax error:" prefix — as ACCEPT (reason=semantic),
// exactly as it does `QUALIFY is not supported`. So the emulator's ZetaSQL-derived
// grammar SHARES these statements with omni; the union parser's accept is
// corroborated by a true accept/accept differential, not merely triangulated.
//
// This makes the differential a PLAIN accept/accept match (oracle=accept,
// omni=accept), and a far stronger PROVE result than the assumed reject/accept
// override. The guard still earns its keep: if an emulator bump ever flips one of
// these to a real `Syntax error:` (the emulator dropping the grammar), the
// recorded oracle=accept would drift to reject and this test fails, forcing a
// fresh look. omni's accept/reject + AST correctness is proven by the unit tests
// in bq_export_load_test.go; this file ties those to the live oracle. (Reuses the
// newDDLHarness plumbing from ddl_oracle_test.go under the googlesql_oracle tag.)
package parser

import (
	"os"
	"testing"
)

// exportLoadCloneOracleFixtures: each is a documented BigQuery data-movement form
// the union parser must ACCEPT, paired with the live Spanner-emulator verdict
// observed 2026-06-05 — all ACCEPT (reason=semantic: the emulator parses them and
// rejects only at the resolver with "Statement not supported: <node>"). See the
// EMPIRICAL FINDING note above.
var exportLoadCloneOracleFixtures = []bqOracleExpectation{
	// ---- EXPORT DATA (OTHER-001) ----
	{"EXPORT DATA OPTIONS(uri='gs://bucket/folder/*.csv', format='CSV') AS SELECT field1 FROM mydataset.table1", "accept", true, "emulator parses EXPORT DATA; rejects only semantically (Statement not supported: ExportDataStatement)"},
	{"EXPORT DATA WITH CONNECTION `myproject.us.my-connection` OPTIONS(uri='s3://bucket/path/file_*.json', format='JSON') AS SELECT * FROM mydataset.mytable", "accept", true, "EXPORT DATA WITH CONNECTION parses; semantic reject only"},
	{"EXPORT DATA AS SELECT 1 AS n", "accept", true, "EXPORT DATA without OPTIONS parses; semantic reject only"},

	// ---- EXPORT MODEL ----
	{"EXPORT MODEL mydataset.mymodel OPTIONS(uri='gs://bucket/path')", "accept", true, "EXPORT MODEL parses; semantic reject only (Statement not supported: ExportModelStatement)"},

	// ---- LOAD DATA (OTHER-003) ----
	{"LOAD DATA INTO mydataset.mytable FROM FILES (format = 'CSV', uris = ['gs://mybucket/myfile.csv'])", "accept", true, "LOAD DATA parses; semantic reject only (Statement not supported: AuxLoadDataStatement)"},
	{"LOAD DATA OVERWRITE mydataset.mytable OPTIONS(description=\"Refreshed table\") FROM FILES (format = 'PARQUET', uris = ['gs://mybucket/data*.parquet'])", "accept", true, "LOAD DATA OVERWRITE parses; semantic reject only"},
	{"LOAD DATA INTO mydataset.mytable FROM FILES (format = 'CSV', uris = ['gs://mybucket/year=*/data.csv'], hive_partition_uri_prefix = 'gs://mybucket/') WITH PARTITION COLUMNS (year INT64, month INT64)", "accept", true, "LOAD DATA WITH PARTITION COLUMNS parses; semantic reject only"},
	{"LOAD DATA INTO TEMP TABLE ds.t FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])", "accept", true, "LOAD DATA TEMP TABLE parses; semantic reject only"},

	// ---- CLONE DATA ----
	{"CLONE DATA INTO ds.dest FROM ds.src", "accept", true, "CLONE DATA parses; semantic reject only (Statement not supported: CloneDataStatement)"},
	{"CLONE DATA INTO ds.dest FROM ds.a UNION ALL ds.b", "accept", true, "CLONE DATA UNION ALL parses; semantic reject only"},
	{"CLONE DATA INTO ds.dest FROM ds.src FOR SYSTEM_TIME AS OF CURRENT_TIMESTAMP() WHERE x > 0", "accept", true, "CLONE DATA FOR SYSTEM_TIME … WHERE parses; semantic reject only"},
}

// exportLoadCloneDivergentFixtures — the NEGATIVE-polarity differential. These
// forms drop a clause the legacy GoogleSQLParser.g4 (and the BigQuery docs) make
// REQUIRED. omni REJECTS them (correct parity: it matches the pinned .g4 + the
// truth1 docs that bytebase's legacy parser enforces). The LIVE emulator, on a
// NEWER/looser ZetaSQL, parses them and rejects only at the resolver → ACCEPT
// (reason=semantic). This is a DEFENDED divergence: omni is intentionally
// stricter than the emulator, anchored to the authoritative .g4 + docs, NOT to
// the non-authoritative emulator (oracle.md: the emulator's GoogleSQL is a
// SUBSET/own-version of the union — its verdict on BigQuery-only forms is not
// binding). Recorded so the boundary is explicit and any future flip is caught.
var exportLoadCloneDivergentFixtures = []bqOracleExpectation{
	// export_data_statement REQUIRES `as_query` (= export_data_no_query AS query).
	{"EXPORT DATA OPTIONS(uri='gs://b/*.csv', format='CSV')", "accept", false, "no AS query — .g4 requires as_query; emulator parses, omni follows .g4+docs and rejects"},
	// aux_load_data_statement REQUIRES the FROM FILES options list.
	{"LOAD DATA INTO ds.t", "accept", false, "no FROM FILES — .g4 requires aux_load_data_from_files_options_list; omni rejects"},
	{"LOAD DATA INTO ds.t OPTIONS(description='d')", "accept", false, "no FROM FILES after OPTIONS — omni rejects per .g4"},
	// aux_load_data_statement REQUIRES append_or_overwrite (INTO|OVERWRITE).
	{"LOAD DATA ds.t FROM FILES (format='CSV', uris=['gs://b/f.csv'])", "accept", false, "no INTO/OVERWRITE — .g4 requires append_or_overwrite; omni rejects"},
	// clone_data_statement REQUIRES INTO.
	{"CLONE DATA ds.dest FROM ds.src", "accept", false, "no INTO — .g4 requires CLONE DATA INTO; omni rejects"},
	// clone_data_source_list separator is UNION ALL, not bare UNION.
	{"CLONE DATA INTO ds.dest FROM ds.a UNION ds.b", "accept", false, "UNION without ALL — .g4 separator is UNION ALL; omni rejects"},
	// export_model_statement uses path_expression (NO dashes); the looser emulator
	// accepts a dashed model name, omni rejects per the pinned .g4.
	{"EXPORT MODEL my-project.ds.m OPTIONS(uri='gs://b')", "accept", false, "dashed model name — .g4 export_model uses path_expression (no dashes); emulator parses, omni rejects"},
}

// TestExportLoadCloneTriangulation is the positive polarity: every documented
// BigQuery data-movement form, asserted oracle=accept (semantic) AND omni=accept.
func TestExportLoadCloneTriangulation(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live triangulation guard")
	}
	h := newDDLHarness(t)
	defer h.close()
	for _, fx := range exportLoadCloneOracleFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) { assertExportLoadCloneFixture(t, h, fx) })
	}
}

// TestExportLoadCloneDivergence is the negative polarity: forms that drop a
// .g4-required clause. omni REJECTS (parity with the pinned .g4 + the BigQuery
// docs); the looser live emulator ACCEPTs (semantic). A DEFENDED divergence — the
// authoritative grammar + docs win over the non-authoritative emulator.
func TestExportLoadCloneDivergence(t *testing.T) {
	if os.Getenv("SPANNER_EMULATOR_HOST") == "" {
		t.Skip("SPANNER_EMULATOR_HOST not set; skipping live triangulation guard")
	}
	h := newDDLHarness(t)
	defer h.close()
	for _, fx := range exportLoadCloneDivergentFixtures {
		fx := fx
		t.Run(fx.sql, func(t *testing.T) { assertExportLoadCloneFixture(t, h, fx) })
	}
}

// assertExportLoadCloneFixture checks one fixture: the live emulator verdict must
// equal the recorded fx.oracle (drift => the emulator changed; re-review), and
// omni's accept/reject must equal fx.omniAccepts. The two are recorded
// independently — a defended divergence has oracle=accept with omniAccepts=false.
func assertExportLoadCloneFixture(t *testing.T, h *ddlHarness, fx bqOracleExpectation) {
	t.Helper()
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

	// The recorded oracle expectation MUST still hold. If it drifts, the emulator
	// changed (e.g. it tightened a clause back to a hard syntax error, or added a
	// form) and the triangulation for this fixture needs a fresh look — do NOT
	// silently update.
	if v.Verdict != fx.oracle {
		t.Fatalf("oracle verdict drift on %q: recorded=%q now=%q (%s: %s) [%s]\n"+
			"the Spanner emulator changed; re-review the BigQuery-vs-emulator triangulation for this form",
			fx.sql, fx.oracle, v.Verdict, v.Reason, v.Message, fx.note)
	}

	// omni's verdict — asserted against the recorded expectation (which is anchored
	// to the legacy .g4 + the BigQuery docs, NOT to the non-authoritative emulator).
	_, errs := Parse(fx.sql)
	omniAccepts := len(errs) == 0
	if omniAccepts != fx.omniAccepts {
		t.Errorf("omni verdict on %q: accepts=%v, want %v; errs=%v [%s]",
			fx.sql, omniAccepts, fx.omniAccepts, errs, fx.note)
	}
}
