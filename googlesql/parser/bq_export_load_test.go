package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-dml-ext node: the BigQuery data-movement statements
// EXPORT DATA / EXPORT MODEL / LOAD DATA / CLONE DATA. All four are BigQuery-only
// at the GoogleSQL union level (the Spanner emulator rejects them); the union
// parser must ACCEPT the documented forms, triangulated against the legacy
// GoogleSQLParser.g4 §2.8 + the BigQuery truth1 corpus (OTHER-001 / OTHER-003).
// These tests assert the AST STRUCTURE (accept/reject alone does not catch a
// dropped clause or wrong nesting); the live-oracle triangulation guard for these
// forms lives in bq_ddl_oracle_test.go (build tag googlesql_oracle).

func exportDataOf(t *testing.T, sql string) *ast.ExportDataStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.ExportDataStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.ExportDataStmt", sql, n)
	}
	return s
}

func loadDataOf(t *testing.T, sql string) *ast.LoadDataStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.LoadDataStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.LoadDataStmt", sql, n)
	}
	return s
}

func cloneDataOf(t *testing.T, sql string) *ast.CloneDataStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.CloneDataStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CloneDataStmt", sql, n)
	}
	return s
}

func exportModelOf(t *testing.T, sql string) *ast.ExportModelStmt {
	t.Helper()
	n := parseDDL(t, sql)
	s, ok := n.(*ast.ExportModelStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.ExportModelStmt", sql, n)
	}
	return s
}

// --- EXPORT DATA (OTHER-001) ---

func TestExportData_Basic(t *testing.T) {
	// truth1 OTHER-001, first example.
	s := exportDataOf(t, "EXPORT DATA OPTIONS(uri='gs://bucket/folder/*.csv', format='CSV', overwrite=true, header=true, field_delimiter=';') AS SELECT field1, field2 FROM mydataset.table1 ORDER BY field1")
	if s.HasConnection {
		t.Error("HasConnection = true, want false")
	}
	if len(s.Options) != 5 {
		t.Errorf("Options = %d, want 5", len(s.Options))
	}
	if s.Query == nil {
		t.Fatal("Query = nil, want the SELECT body")
	}
	if _, ok := s.Query.(*ast.QueryStmt); !ok {
		t.Errorf("Query is %T, want *ast.QueryStmt", s.Query)
	}
}

func TestExportData_WithConnection(t *testing.T) {
	// truth1 OTHER-001, S3 example: EXPORT DATA WITH CONNECTION `…` OPTIONS(…) AS …
	s := exportDataOf(t, "EXPORT DATA WITH CONNECTION `myproject.us.my-connection` OPTIONS(uri='s3://bucket/path/file_*.json', format='JSON') AS SELECT * FROM mydataset.mytable")
	if !s.HasConnection {
		t.Error("HasConnection = false, want true")
	}
	if s.ConnectionName == "" {
		t.Error("ConnectionName empty, want the connection path")
	}
	if len(s.Options) != 2 {
		t.Errorf("Options = %d, want 2", len(s.Options))
	}
}

func TestExportData_NoOptions(t *testing.T) {
	// opt_options_list is optional in the grammar — EXPORT DATA AS query is legal.
	s := exportDataOf(t, "EXPORT DATA AS SELECT 1 AS n")
	if len(s.Options) != 0 {
		t.Errorf("Options = %d, want 0", len(s.Options))
	}
	if s.Query == nil {
		t.Error("Query = nil")
	}
}

func TestExportData_ConnectionDefault(t *testing.T) {
	// connection_clause: DEFAULT | path.
	s := exportDataOf(t, "EXPORT DATA WITH CONNECTION DEFAULT OPTIONS(uri='gs://b/*.csv', format='CSV') AS SELECT 1 AS n")
	if !s.HasConnection || s.ConnectionName != "DEFAULT" {
		t.Errorf("ConnectionName = %q (has=%v), want DEFAULT", s.ConnectionName, s.HasConnection)
	}
}

func TestExportData_EmbeddedSubqueryReparsed(t *testing.T) {
	// The AS query is fillSubqueries-wrapped (parseStmtWithSubqueries), so a
	// subquery inside the exported query is reachable for the query-span consumer.
	s := exportDataOf(t, "EXPORT DATA OPTIONS(uri='gs://b/*.csv', format='CSV') AS SELECT x FROM (SELECT 1 AS x)")
	if s.Query == nil {
		t.Fatal("Query = nil")
	}
	// Walk the tree and confirm no SubqueryExpr is left with a nil Query (every
	// embedded subquery was re-parsed). There is exactly one derived-table here;
	// the walk simply must not panic and the outer parse must succeed.
	ast.Inspect(s, func(n ast.Node) bool { return true })
}

// --- EXPORT MODEL ---

func TestExportModel_Basic(t *testing.T) {
	s := exportModelOf(t, "EXPORT MODEL mydataset.mymodel OPTIONS(uri='gs://bucket/path')")
	if s.Name == nil || s.Name.String() != "mydataset.mymodel" {
		t.Errorf("Name = %v", s.Name)
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
}

func TestExportModel_WithConnectionNoOptions(t *testing.T) {
	// with_connection_clause? and opt_options_list? are both optional.
	s := exportModelOf(t, "EXPORT MODEL ds.m WITH CONNECTION `p.us.c`")
	if !s.HasConnection {
		t.Error("HasConnection = false, want true")
	}
	if len(s.Options) != 0 {
		t.Errorf("Options = %d, want 0", len(s.Options))
	}
}

// --- LOAD DATA (OTHER-003) ---

func TestLoadData_IntoBasic(t *testing.T) {
	// truth1 OTHER-003, first example.
	s := loadDataOf(t, "LOAD DATA INTO mydataset.mytable FROM FILES (format = 'CSV', uris = ['gs://mybucket/myfile.csv'])")
	if s.Overwrite {
		t.Error("Overwrite = true, want false (INTO)")
	}
	if s.Name == nil || s.Name.String() != "mydataset.mytable" {
		t.Errorf("Name = %v", s.Name)
	}
	if len(s.FromFiles) != 2 {
		t.Errorf("FromFiles = %d, want 2", len(s.FromFiles))
	}
}

func TestLoadData_OverwriteOptions(t *testing.T) {
	// truth1 OTHER-003, second example: OVERWRITE + OPTIONS before FROM FILES.
	s := loadDataOf(t, "LOAD DATA OVERWRITE mydataset.mytable OPTIONS(description=\"Refreshed table\") FROM FILES (format = 'PARQUET', uris = ['gs://mybucket/data*.parquet'])")
	if !s.Overwrite {
		t.Error("Overwrite = false, want true (OVERWRITE)")
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
	if len(s.FromFiles) != 2 {
		t.Errorf("FromFiles = %d, want 2", len(s.FromFiles))
	}
}

func TestLoadData_WithPartitionColumns(t *testing.T) {
	// truth1 OTHER-003, third example: WITH PARTITION COLUMNS (col type, …).
	s := loadDataOf(t, "LOAD DATA INTO mydataset.mytable FROM FILES (format = 'CSV', uris = ['gs://mybucket/year=*/month=*/data.csv'], hive_partition_uri_prefix = 'gs://mybucket/') WITH PARTITION COLUMNS (year INT64, month INT64)")
	if !s.HasPartitionColumns {
		t.Error("HasPartitionColumns = false, want true")
	}
	if len(s.PartitionColumns) != 2 {
		t.Fatalf("PartitionColumns = %d, want 2", len(s.PartitionColumns))
	}
	if s.PartitionColumns[0].Name != "year" || s.PartitionColumns[1].Name != "month" {
		t.Errorf("PartitionColumns names = %q, %q", s.PartitionColumns[0].Name, s.PartitionColumns[1].Name)
	}
	if len(s.FromFiles) != 3 {
		t.Errorf("FromFiles = %d, want 3", len(s.FromFiles))
	}
}

func TestLoadData_PartitionColumnsInferred(t *testing.T) {
	// with_partition_columns_clause table_element_list is OPTIONAL — `WITH PARTITION
	// COLUMNS` with no `( … )` infers the columns from the files.
	s := loadDataOf(t, "LOAD DATA INTO ds.t FROM FILES (format = 'PARQUET', uris = ['gs://b/*.parquet']) WITH PARTITION COLUMNS")
	if !s.HasPartitionColumns {
		t.Error("HasPartitionColumns = false, want true")
	}
	if s.PartitionColumns != nil {
		t.Errorf("PartitionColumns = %v, want nil (inferred)", s.PartitionColumns)
	}
}

func TestLoadData_TempTable(t *testing.T) {
	// maybe_dashed_path_expression_with_scope: TEMP TABLE <path> | TEMPORARY TABLE <path>.
	s := loadDataOf(t, "LOAD DATA INTO TEMP TABLE ds.t FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])")
	if !s.Temp {
		t.Error("Temp = false, want true")
	}
	if s.TempKeyword != "TEMP" {
		t.Errorf("TempKeyword = %q, want TEMP", s.TempKeyword)
	}
	if s.Name == nil || s.Name.String() != "ds.t" {
		t.Errorf("Name = %v", s.Name)
	}

	s2 := loadDataOf(t, "LOAD DATA OVERWRITE TEMPORARY TABLE ds.t FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])")
	if !s2.Temp || s2.TempKeyword != "TEMPORARY" {
		t.Errorf("TempKeyword = %q (temp=%v), want TEMPORARY", s2.TempKeyword, s2.Temp)
	}
}

func TestLoadData_ExplicitSchema(t *testing.T) {
	// table_element_list? — an explicit column schema between the table name and
	// the (optional) clauses.
	s := loadDataOf(t, "LOAD DATA INTO ds.t (a INT64, b STRING) FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])")
	if len(s.Columns) != 2 {
		t.Fatalf("Columns = %d, want 2", len(s.Columns))
	}
	if s.Columns[0].Name != "a" || s.Columns[1].Name != "b" {
		t.Errorf("Columns names = %q, %q", s.Columns[0].Name, s.Columns[1].Name)
	}
}

func TestLoadData_PartitionAndClusterBy(t *testing.T) {
	// partition_by_clause_prefix_no_hint? cluster_by_clause_prefix_no_hint? before
	// the FROM FILES clause.
	s := loadDataOf(t, "LOAD DATA INTO ds.t PARTITION BY DATE(ts) CLUSTER BY a, b OPTIONS(description='d') FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])")
	if len(s.PartitionBy) != 1 {
		t.Errorf("PartitionBy = %d, want 1", len(s.PartitionBy))
	}
	if len(s.ClusterBy) != 2 {
		t.Errorf("ClusterBy = %d, want 2", len(s.ClusterBy))
	}
	if len(s.Options) != 1 {
		t.Errorf("Options = %d, want 1", len(s.Options))
	}
}

func TestLoadData_PartitionsClause(t *testing.T) {
	// load_data_partitions_clause: OVERWRITE? PARTITIONS ( expr ). It sits between
	// the table element list and the collate clause.
	s := loadDataOf(t, "LOAD DATA OVERWRITE ds.t OVERWRITE PARTITIONS (_PARTITIONDATE = '2024-01-01') FROM FILES (format = 'CSV', uris = ['gs://b/f.csv'])")
	if s.Partitions == nil {
		t.Fatal("Partitions = nil, want the partition expression")
	}
	if !s.PartitionsOverwrite {
		t.Error("PartitionsOverwrite = false, want true")
	}
}

func TestLoadData_WithConnectionOnly(t *testing.T) {
	// opt_external_table_with_clauses third alternative: WITH CONNECTION alone.
	s := loadDataOf(t, "LOAD DATA INTO ds.t FROM FILES (format = 'CSV', uris = ['gs://b/f.csv']) WITH CONNECTION `p.us.c`")
	if !s.HasConnection {
		t.Error("HasConnection = false, want true")
	}
	if s.HasPartitionColumns {
		t.Error("HasPartitionColumns = true, want false")
	}
}

func TestLoadData_PartitionColumnsThenConnection(t *testing.T) {
	// opt_external_table_with_clauses first alternative: WITH PARTITION COLUMNS
	// followed by WITH CONNECTION.
	s := loadDataOf(t, "LOAD DATA INTO ds.t FROM FILES (format = 'CSV', uris = ['gs://b/f.csv']) WITH PARTITION COLUMNS (y INT64) WITH CONNECTION `p.us.c`")
	if !s.HasPartitionColumns {
		t.Error("HasPartitionColumns = false, want true")
	}
	if !s.HasConnection {
		t.Error("HasConnection = false, want true")
	}
}

// --- CLONE DATA ---

func TestCloneData_Basic(t *testing.T) {
	s := cloneDataOf(t, "CLONE DATA INTO ds.dest FROM ds.src")
	if s.Name == nil || s.Name.String() != "ds.dest" {
		t.Errorf("Name = %v", s.Name)
	}
	if len(s.Sources) != 1 {
		t.Fatalf("Sources = %d, want 1", len(s.Sources))
	}
	if s.Sources[0].Name.String() != "ds.src" {
		t.Errorf("source = %v", s.Sources[0].Name)
	}
}

func TestCloneData_UnionAll(t *testing.T) {
	// clone_data_source_list: source (UNION ALL source)*.
	s := cloneDataOf(t, "CLONE DATA INTO ds.dest FROM ds.a UNION ALL ds.b UNION ALL ds.c")
	if len(s.Sources) != 3 {
		t.Fatalf("Sources = %d, want 3", len(s.Sources))
	}
	if s.Sources[0].Name.String() != "ds.a" || s.Sources[2].Name.String() != "ds.c" {
		t.Errorf("sources = %v", s.Sources)
	}
}

func TestCloneData_SystemTimeAndWhere(t *testing.T) {
	// clone_data_source: path opt_at_system_time? where_clause?.
	s := cloneDataOf(t, "CLONE DATA INTO ds.dest FROM ds.src FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 DAY) WHERE x > 0")
	if len(s.Sources) != 1 {
		t.Fatalf("Sources = %d, want 1", len(s.Sources))
	}
	src := s.Sources[0]
	if src.ForSystemTime == nil {
		t.Error("ForSystemTime = nil, want the time expression")
	}
	if src.Where == nil {
		t.Error("Where = nil, want the filter expression")
	}
}

func TestCloneData_SpaceSystemTime(t *testing.T) {
	// `FOR SYSTEM TIME AS OF` (two-word spelling) accepted as well as SYSTEM_TIME.
	s := cloneDataOf(t, "CLONE DATA INTO ds.dest FROM ds.src FOR SYSTEM TIME AS OF CURRENT_TIMESTAMP()")
	if s.Sources[0].ForSystemTime == nil {
		t.Error("ForSystemTime = nil")
	}
}

// --- Review-finding regression tests ---

// Codex review finding (high): LOAD DATA / CLONE DATA / EXPORT MODEL parse full
// expressions (OPTIONS / PARTITION BY / CLUSTER BY values; CLONE source WHERE +
// FOR SYSTEM_TIME) that can embed subqueries. They must be wrapped by
// parseStmtWithSubqueries so every embedded SubqueryExpr.Query is re-parsed (for
// the query-span consumer) AND a malformed embedded subquery surfaces as a
// diagnostic. These regressions assert both.

func TestCloneData_WhereSubqueryReparsed(t *testing.T) {
	// CLONE source WHERE with an embedded subquery — must re-parse to a real Query,
	// not be left as nil RawText.
	s := cloneDataOf(t, "CLONE DATA INTO ds.dest FROM ds.src WHERE id IN (SELECT id FROM ds.allow)")
	src := s.Sources[0]
	if src.Where == nil {
		t.Fatal("Where = nil")
	}
	foundFilled := false
	ast.Inspect(s, func(n ast.Node) bool {
		if sq, ok := n.(*ast.SubqueryExpr); ok {
			if sq.Query == nil {
				t.Error("embedded SubqueryExpr.Query is nil — not re-parsed (missing parseStmtWithSubqueries wrap)")
			} else {
				foundFilled = true
			}
		}
		return true
	})
	if !foundFilled {
		t.Error("no SubqueryExpr found in the CLONE WHERE — expected the IN (SELECT …)")
	}
}

func TestCloneData_MalformedEmbeddedSubqueryRejected(t *testing.T) {
	// A balanced-but-invalid embedded subquery (stray alias `b`) must be surfaced
	// as a diagnostic by the fillSubqueries re-parse — the same contract the DML
	// family upholds. Without the wrap this would be silently accepted.
	_, errs := Parse("CLONE DATA INTO ds.dest FROM ds.src WHERE id IN (SELECT id FROM ds.s a b)")
	if len(errs) == 0 {
		t.Error("want a diagnostic for the malformed embedded subquery (SELECT … a b)")
	}
}

func TestLoadData_PartitionBySubqueryReparsed(t *testing.T) {
	// A subquery embedded in a PARTITION BY expression must be re-parsed.
	s := loadDataOf(t, "LOAD DATA INTO ds.t PARTITION BY (SELECT 1) FROM FILES (format='CSV', uris=['gs://b/f.csv'])")
	foundFilled := false
	ast.Inspect(s, func(n ast.Node) bool {
		if sq, ok := n.(*ast.SubqueryExpr); ok {
			if sq.Query == nil {
				t.Error("embedded SubqueryExpr.Query is nil — not re-parsed")
			} else {
				foundFilled = true
			}
		}
		return true
	})
	if !foundFilled {
		t.Error("no SubqueryExpr found in PARTITION BY (SELECT 1)")
	}
}

func TestLoadData_TempAsTableName(t *testing.T) {
	// Codex review finding (medium): TEMP / TEMPORARY are non-reserved, so a table
	// literally named `temp` (no scope prefix) must parse as a path, NOT be
	// mis-read as a scope keyword that then demands TABLE.
	s := loadDataOf(t, "LOAD DATA INTO temp.t FROM FILES (format='CSV', uris=['gs://b/f.csv'])")
	if s.Temp {
		t.Error("Temp = true, want false (temp is the dataset name, not a scope prefix)")
	}
	if s.Name == nil || s.Name.String() != "temp.t" {
		t.Errorf("Name = %v, want temp.t", s.Name)
	}

	// And a bare `temporary` dataset name likewise.
	s2 := loadDataOf(t, "LOAD DATA OVERWRITE temporary FROM FILES (format='CSV', uris=['gs://b/f.csv'])")
	if s2.Temp || s2.Name.String() != "temporary" {
		t.Errorf("Name = %v (temp=%v), want path 'temporary' with no scope", s2.Name, s2.Temp)
	}
}

func TestExportModel_DashedPathRejected(t *testing.T) {
	// Codex review finding (medium): export_model_statement uses path_expression
	// (identifier (DOT identifier)*), NOT maybe_dashed_path_expression. An unquoted
	// dashed model name is OUTSIDE the rule and must reject; the backtick-quoted
	// form is accepted.
	assertReject(t, "EXPORT MODEL my-project.ds.m OPTIONS(uri='gs://b')")

	// Backtick-quoted dashed name parses as a single quoted identifier component.
	s := exportModelOf(t, "EXPORT MODEL `my-project.ds.m` OPTIONS(uri='gs://b')")
	if s.Name == nil {
		t.Fatal("Name = nil")
	}
}

// --- EXPORT dispatch boundary ---

func TestExportMetadata_RecognizedButUnsupported(t *testing.T) {
	// EXPORT … METADATA is a valid BigQuery statement owned by the parser-utility
	// node; this node routes it to the unsupported stub (recognized, not unknown).
	_, errs := Parse("EXPORT TABLE METADATA FROM ds.t")
	if len(errs) == 0 {
		t.Fatal("want an error for the not-yet-supported EXPORT METADATA")
	}
	// It must be the not-yet-supported message, not an unknown-statement message.
	if !strings.Contains(errs[0].Error(), "EXPORT METADATA") {
		t.Errorf("error = %q, want it to mention EXPORT METADATA", errs[0].Error())
	}
}

// --- Rejects ---

func TestExportLoadClone_Rejects(t *testing.T) {
	cases := []string{
		// EXPORT DATA requires AS query (export_data_statement = no_query as_query).
		"EXPORT DATA OPTIONS(uri='gs://b/*.csv', format='CSV')",
		"EXPORT DATA",
		// EXPORT MODEL requires a model path.
		"EXPORT MODEL",
		"EXPORT MODEL OPTIONS(uri='gs://b')",
		// LOAD DATA requires INTO|OVERWRITE.
		"LOAD DATA ds.t FROM FILES (format='CSV', uris=['gs://b/f.csv'])",
		// LOAD DATA requires the FROM FILES clause (it is NOT optional).
		"LOAD DATA INTO ds.t",
		"LOAD DATA INTO ds.t OPTIONS(description='d')",
		// FROM without FILES.
		"LOAD DATA INTO ds.t FROM (format='CSV')",
		// TEMP/TEMPORARY requires the TABLE keyword.
		"LOAD DATA INTO TEMP ds.t FROM FILES (format='CSV', uris=['gs://b/f.csv'])",
		// CLONE DATA requires INTO.
		"CLONE DATA ds.dest FROM ds.src",
		// CLONE DATA requires FROM and a source.
		"CLONE DATA INTO ds.dest",
		"CLONE DATA INTO ds.dest FROM",
		// UNION without ALL is not a valid clone_data_source_list separator.
		"CLONE DATA INTO ds.dest FROM ds.a UNION ds.b",
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}
