package parser

import (
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl-bigquery node: CREATE [SEARCH|VECTOR] INDEX,
// ALTER VECTOR INDEX … REBUILD, DROP [SEARCH|VECTOR] INDEX. BigQuery-only at the
// union level (the SEARCH-index ALL COLUMNS form, the DROP … ON form, and the
// REBUILD form all reject on the Spanner emulator, probed 2026-06-05); verdicts
// triangulated against the legacy GoogleSQLParser.g4 + BigQuery truth1
// (DDL-022/023/041/054). NB: ALTER VECTOR INDEX … REBUILD is NOT in the legacy
// .g4 — it is implemented from the BigQuery docs as an additive union form
// (flagged divergence; see bq_search_vector_index.go header).

func sviOf(t *testing.T, sql string) *ast.SearchVectorIndexStmt {
	t.Helper()
	n := parseDDL(t, sql)
	svi, ok := n.(*ast.SearchVectorIndexStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.SearchVectorIndexStmt", sql, n)
	}
	return svi
}

func TestCreateSearchIndex_AllColumns(t *testing.T) {
	// DDL-022.
	svi := sviOf(t, "CREATE SEARCH INDEX my_index ON mydataset.mytable(ALL COLUMNS)")
	if svi.IsVector {
		t.Error("IsVector = true, want false (SEARCH)")
	}
	if !svi.AllColumns {
		t.Error("AllColumns = false, want true")
	}
	if svi.Name.String() != "my_index" || svi.Table.String() != "mydataset.mytable" {
		t.Errorf("Name=%q Table=%q", svi.Name.String(), svi.Table.String())
	}
}

func TestCreateSearchIndex_ColumnList(t *testing.T) {
	// DDL-022.
	svi := sviOf(t, "CREATE SEARCH INDEX my_index ON mydataset.mytable(name, description)")
	if svi.AllColumns {
		t.Error("AllColumns = true, want false")
	}
	if len(svi.Keys) != 2 {
		t.Errorf("Keys = %d, want 2", len(svi.Keys))
	}
}

func TestCreateSearchIndex_AllColumnsWithOptions(t *testing.T) {
	svi := sviOf(t, "CREATE SEARCH INDEX i ON t(ALL COLUMNS WITH COLUMN OPTIONS(a))")
	if !svi.AllColumns {
		t.Error("AllColumns = false, want true")
	}
}

func TestCreateVectorIndex(t *testing.T) {
	// DDL-023.
	svi := sviOf(t, "CREATE VECTOR INDEX my_index ON mydataset.mytable(embedding_column) OPTIONS(index_type='TREE_AH', distance_type='COSINE')")
	if !svi.IsVector {
		t.Error("IsVector = false, want true")
	}
	if len(svi.Keys) != 1 {
		t.Errorf("Keys = %d, want 1", len(svi.Keys))
	}
	if len(svi.Options) != 2 {
		t.Errorf("Options = %d, want 2", len(svi.Options))
	}
}

func TestCreateVectorIndex_StoringPartition(t *testing.T) {
	// DDL-023: STORING + PARTITION BY.
	svi := sviOf(t, "CREATE OR REPLACE VECTOR INDEX i ON t(emb) STORING(a, b) PARTITION BY DATE(d) OPTIONS(distance_type='COSINE')")
	if !svi.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if len(svi.Storing) != 2 {
		t.Errorf("Storing = %d, want 2", len(svi.Storing))
	}
	if len(svi.PartitionBy) != 1 {
		t.Errorf("PartitionBy = %d, want 1", len(svi.PartitionBy))
	}
}

func TestCreateVectorIndex_NoOptionsAcceptedPerGrammar(t *testing.T) {
	// DEFEND (review finding): the legacy create_index_statement makes the
	// option/suffix tail OPTIONAL (`opt_create_index_statement_suffix?`). BigQuery
	// docs show OPTIONS on a VECTOR index, but with no Spanner oracle for this exact
	// shape and the migration target being legacy-grammar parity, we FOLLOW THE
	// GRAMMAR and accept a VECTOR index with no OPTIONS.
	svi := sviOf(t, "CREATE VECTOR INDEX i ON ds.t(embedding)")
	if !svi.IsVector || len(svi.Options) != 0 {
		t.Errorf("IsVector=%v Options=%d, want vector index with no options", svi.IsVector, len(svi.Options))
	}
}

func TestDropSearchVectorIndex_NoOnAcceptedPerGrammar(t *testing.T) {
	// DEFEND (review finding): drop_statement's `DROP index_type INDEX … path
	// on_path_expression?` makes the ON clause OPTIONAL. BigQuery docs show ON, but
	// the grammar permits its absence; we follow the grammar (accept, OnTable nil).
	d := bqDropOf(t, "DROP SEARCH INDEX i")
	if d.Object != ast.BQDropSearchIndex || d.OnTable != nil {
		t.Errorf("got Object=%v OnTable=%v, want SEARCH INDEX with nil OnTable", d.Object, d.OnTable)
	}
}

func TestCreateSearchVectorIndex_Rejects(t *testing.T) {
	cases := []string{
		"CREATE SEARCH INDEX i",              // missing ON
		"CREATE SEARCH INDEX i ON t",         // missing column list
		"CREATE VECTOR i ON t(c)",            // missing INDEX keyword
		"CREATE TEMP SEARCH INDEX i ON t(c)", // index takes no scope
		// Regression (review finding): index_all_columns requires COLUMNS after ALL;
		// a bare `( ALL )` must reject (not be skipped as balanced parens).
		"CREATE SEARCH INDEX i ON t(ALL)",
		"CREATE SEARCH INDEX i ON t(ALL COLUMNS WITH COLUMN)", // WITH COLUMN OPTIONS malformed (missing OPTIONS + list)
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- ALTER VECTOR INDEX REBUILD ---

func TestAlterVectorIndex_Rebuild(t *testing.T) {
	// DDL-041 (NOT in the legacy .g4; additive union form).
	a := bqAlterOf(t, "ALTER VECTOR INDEX my_index ON mydataset.mytable REBUILD")
	if a.Object != ast.BQAlterVectorIndex {
		t.Errorf("Object = %v, want VECTOR INDEX", a.Object)
	}
	if !a.Rebuild {
		t.Error("Rebuild = false, want true")
	}
	if a.OnTable == nil || a.OnTable.String() != "mydataset.mytable" {
		t.Errorf("OnTable = %v", a.OnTable)
	}
}

func TestAlterVectorIndex_RebuildOptions(t *testing.T) {
	a := bqAlterOf(t, "ALTER VECTOR INDEX IF EXISTS i ON t REBUILD OPTIONS(x=1)")
	if !a.IfExists || !a.Rebuild {
		t.Errorf("IfExists=%v Rebuild=%v", a.IfExists, a.Rebuild)
	}
	if len(a.SetOptions) != 1 {
		t.Errorf("SetOptions = %d, want 1", len(a.SetOptions))
	}
}

func TestAlterVectorIndex_Rejects(t *testing.T) {
	cases := []string{
		"ALTER VECTOR INDEX i ON t",    // missing REBUILD
		"ALTER VECTOR INDEX i REBUILD", // missing ON table
		"ALTER VECTOR i ON t REBUILD",  // missing INDEX
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

// --- DROP SEARCH/VECTOR INDEX ---

func TestDropSearchIndex(t *testing.T) {
	// DDL-054.
	d := bqDropOf(t, "DROP SEARCH INDEX IF EXISTS my_index ON mydataset.mytable")
	if d.Object != ast.BQDropSearchIndex {
		t.Errorf("Object = %v, want SEARCH INDEX", d.Object)
	}
	if d.OnTable == nil || d.OnTable.String() != "mydataset.mytable" {
		t.Errorf("OnTable = %v", d.OnTable)
	}
}

func TestDropVectorIndex(t *testing.T) {
	// DDL-054.
	d := bqDropOf(t, "DROP VECTOR INDEX IF EXISTS my_index ON mydataset.mytable")
	if d.Object != ast.BQDropVectorIndex {
		t.Errorf("Object = %v, want VECTOR INDEX", d.Object)
	}
}
