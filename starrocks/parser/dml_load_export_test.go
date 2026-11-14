package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

func mustParseTruncateTable(t *testing.T, input string) *ast.TruncateTableStmt {
	t.Helper()
	file, errs := parseForTest(input)
	if len(errs) > 0 {
		t.Fatalf("parseForTest(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("parseForTest(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.TruncateTableStmt)
	if !ok {
		t.Fatalf("parseForTest(%q) got %T, want *ast.TruncateTableStmt", input, file.Stmts[0])
	}
	return stmt
}

func mustParseLoadData(t *testing.T, input string) *ast.LoadDataStmt {
	t.Helper()
	file, errs := parseForTest(input)
	if len(errs) > 0 {
		t.Fatalf("parseForTest(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("parseForTest(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.LoadDataStmt)
	if !ok {
		t.Fatalf("parseForTest(%q) got %T, want *ast.LoadDataStmt", input, file.Stmts[0])
	}
	return stmt
}

func mustParseExport(t *testing.T, input string) *ast.ExportStmt {
	t.Helper()
	file, errs := parseForTest(input)
	if len(errs) > 0 {
		t.Fatalf("parseForTest(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("parseForTest(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.ExportStmt)
	if !ok {
		t.Fatalf("parseForTest(%q) got %T, want *ast.ExportStmt", input, file.Stmts[0])
	}
	return stmt
}

func mustParseCopyInto(t *testing.T, input string) *ast.CopyIntoStmt {
	t.Helper()
	file, errs := parseForTest(input)
	if len(errs) > 0 {
		t.Fatalf("parseForTest(%q) errors: %v", input, errs)
	}
	if len(file.Stmts) == 0 {
		t.Fatalf("parseForTest(%q) returned no statements", input)
	}
	stmt, ok := file.Stmts[0].(*ast.CopyIntoStmt)
	if !ok {
		t.Fatalf("parseForTest(%q) got %T, want *ast.CopyIntoStmt", input, file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// TRUNCATE TABLE tests
// ---------------------------------------------------------------------------

func TestTruncateTableBasic(t *testing.T) {
	stmt := mustParseTruncateTable(t, "TRUNCATE TABLE example_db.tbl")
	if len(stmt.Target.Parts) != 2 {
		t.Fatalf("target parts = %v, want 2", stmt.Target.Parts)
	}
	if stmt.Target.Parts[0] != "example_db" {
		t.Errorf("target[0] = %q, want example_db", stmt.Target.Parts[0])
	}
	if stmt.Target.Parts[1] != "tbl" {
		t.Errorf("target[1] = %q, want tbl", stmt.Target.Parts[1])
	}
	if len(stmt.Partition) != 0 {
		t.Errorf("partition = %v, want empty", stmt.Partition)
	}
	if stmt.Force {
		t.Errorf("force = true, want false")
	}
}

func TestTruncateTableSimple(t *testing.T) {
	stmt := mustParseTruncateTable(t, "TRUNCATE TABLE tbl")
	if len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "tbl" {
		t.Fatalf("target = %v, want [tbl]", stmt.Target.Parts)
	}
}

func TestTruncateTableWithPartition(t *testing.T) {
	stmt := mustParseTruncateTable(t, "TRUNCATE TABLE tbl PARTITION(p1, p2)")
	if len(stmt.Partition) != 2 {
		t.Fatalf("partition len = %d, want 2", len(stmt.Partition))
	}
	if stmt.Partition[0] != "p1" {
		t.Errorf("partition[0] = %q, want p1", stmt.Partition[0])
	}
	if stmt.Partition[1] != "p2" {
		t.Errorf("partition[1] = %q, want p2", stmt.Partition[1])
	}
}

func TestTruncateTableWithForce(t *testing.T) {
	stmt := mustParseTruncateTable(t, "TRUNCATE TABLE tbl FORCE")
	if !stmt.Force {
		t.Errorf("force = false, want true")
	}
}

func TestTruncateTableWithPartitionAndForce(t *testing.T) {
	stmt := mustParseTruncateTable(t, "TRUNCATE TABLE tbl PARTITION(p1) FORCE")
	if len(stmt.Partition) != 1 || stmt.Partition[0] != "p1" {
		t.Fatalf("partition = %v, want [p1]", stmt.Partition)
	}
	if !stmt.Force {
		t.Errorf("force = false, want true")
	}
}

// TestTruncateTableLegacyCorpus covers both forms from the legacy test corpus.
func TestTruncateTableLegacyCorpus(t *testing.T) {
	tests := []struct {
		sql       string
		wantParts []string
		wantPart  []string
	}{
		{
			sql:       "TRUNCATE TABLE example_db.tbl",
			wantParts: []string{"example_db", "tbl"},
			wantPart:  nil,
		},
		{
			sql:       "TRUNCATE TABLE tbl PARTITION(p1, p2)",
			wantParts: []string{"tbl"},
			wantPart:  []string{"p1", "p2"},
		},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := mustParseTruncateTable(t, tc.sql)
			if len(stmt.Target.Parts) != len(tc.wantParts) {
				t.Fatalf("target parts = %v, want %v", stmt.Target.Parts, tc.wantParts)
			}
			for i, want := range tc.wantParts {
				if stmt.Target.Parts[i] != want {
					t.Errorf("target[%d] = %q, want %q", i, stmt.Target.Parts[i], want)
				}
			}
			if len(stmt.Partition) != len(tc.wantPart) {
				t.Fatalf("partition = %v, want %v", stmt.Partition, tc.wantPart)
			}
			for i, want := range tc.wantPart {
				if stmt.Partition[i] != want {
					t.Errorf("partition[%d] = %q, want %q", i, stmt.Partition[i], want)
				}
			}
		})
	}
}

// ---------------------------------------------------------------------------
// COPY INTO tests
// ---------------------------------------------------------------------------

func TestCopyIntoBasic(t *testing.T) {
	stmt := mustParseCopyInto(t, "COPY INTO my_table FROM 'stage/path'")
	if len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "my_table" {
		t.Fatalf("target = %v, want [my_table]", stmt.Target.Parts)
	}
	if stmt.Source != "stage/path" {
		t.Errorf("source = %q, want stage/path", stmt.Source)
	}
}

func TestCopyIntoWithProperties(t *testing.T) {
	stmt := mustParseCopyInto(t,
		`COPY INTO my_table FROM 'stage/path' PROPERTIES ("key"="val")`)
	if len(stmt.Properties) != 1 {
		t.Fatalf("properties len = %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" || stmt.Properties[0].Value != "val" {
		t.Errorf("property = %q=%q, want key=val", stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestCopyIntoWithFiles(t *testing.T) {
	stmt := mustParseCopyInto(t,
		`COPY INTO my_table FROM 'stage/' FILES = ('f1.csv', 'f2.csv')`)
	if len(stmt.Files) != 2 {
		t.Fatalf("files len = %d, want 2", len(stmt.Files))
	}
	if stmt.Files[0] != "f1.csv" || stmt.Files[1] != "f2.csv" {
		t.Errorf("files = %v, want [f1.csv f2.csv]", stmt.Files)
	}
}

// ---------------------------------------------------------------------------
// LOAD tests
// ---------------------------------------------------------------------------

func TestLoadBasic(t *testing.T) {
	sql := `LOAD LABEL example_db.label1
	(DATA INFILE('hdfs://hdfs_host:hdfs_port/file1') INTO TABLE my_table)
	WITH BROKER hdfs_broker
	PROPERTIES ("timeout"="1200")`

	stmt := mustParseLoadData(t, sql)
	if stmt.Label != "example_db.label1" {
		t.Errorf("label = %q, want example_db.label1", stmt.Label)
	}
	if len(stmt.DataDescs) != 1 {
		t.Fatalf("data descs = %d, want 1", len(stmt.DataDescs))
	}
	desc := stmt.DataDescs[0]
	if len(desc.SourceFiles) != 1 {
		t.Fatalf("source files = %d, want 1", len(desc.SourceFiles))
	}
	if desc.Target == nil || len(desc.Target.Parts) < 1 {
		t.Fatal("desc target is nil")
	}
	if desc.Target.Parts[len(desc.Target.Parts)-1] != "my_table" {
		t.Errorf("desc target = %v, want my_table", desc.Target.Parts)
	}
	if stmt.BrokerName != "hdfs_broker" {
		t.Errorf("broker = %q, want hdfs_broker", stmt.BrokerName)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("properties = %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "timeout" {
		t.Errorf("properties[0].key = %q, want timeout", stmt.Properties[0].Key)
	}
	// loadStatement (StarRocks.g4:1427) has no COMMENT.
	if _, err := Parse(sql + " COMMENT 'test load'"); err == nil {
		t.Error("LOAD ... COMMENT parsed; StarRocks's loadStatement has none")
	}
}

func TestLoadMultipleFiles(t *testing.T) {
	sql := `LOAD LABEL load1
	(DATA INFILE('file1', 'file2') INTO TABLE t)`

	stmt := mustParseLoadData(t, sql)
	if len(stmt.DataDescs) != 1 {
		t.Fatalf("data descs = %d, want 1", len(stmt.DataDescs))
	}
	if len(stmt.DataDescs[0].SourceFiles) != 2 {
		t.Fatalf("source files = %d, want 2", len(stmt.DataDescs[0].SourceFiles))
	}
}

func TestLoadWithPartition(t *testing.T) {
	sql := `LOAD LABEL load2
	(DATA INFILE('file1') INTO TABLE t PARTITION(p1, p2))`

	stmt := mustParseLoadData(t, sql)
	desc := stmt.DataDescs[0]
	if len(desc.Partition) != 2 {
		t.Fatalf("partition = %v, want [p1 p2]", desc.Partition)
	}
}

func TestLoadNegative(t *testing.T) {
	sql := `LOAD LABEL load3
	(DATA INFILE('file1') NEGATIVE INTO TABLE t)`

	stmt := mustParseLoadData(t, sql)
	if !stmt.DataDescs[0].Negative {
		t.Errorf("negative = false, want true")
	}
}

func TestLoadWithColumns(t *testing.T) {
	sql := `LOAD LABEL load4
	(DATA INFILE('file1') INTO TABLE t (c1, c2, c3))`

	stmt := mustParseLoadData(t, sql)
	desc := stmt.DataDescs[0]
	if len(desc.ColumnList) != 3 {
		t.Fatalf("column list = %v, want [c1 c2 c3]", desc.ColumnList)
	}
}

func TestLoadClauses(t *testing.T) {
	stmt := mustParseLoadData(t, `LOAD LABEL db.l
	(DATA INFILE('f') INTO TABLE t TEMPORARY PARTITION (p1, p2) COLUMNS TERMINATED BY ','
	 ROWS TERMINATED BY '\n' FORMAT AS 'csv' (SKIP_HEADER = 1 ENCLOSE = '"') (a, b)
	 COLUMNS FROM PATH AS (c) SET (a = b + 1, b) WHERE a > 1,
	 DATA FROM TABLE s NEGATIVE INTO TABLE t PARTITION (k = 1))
	WITH BROKER ('k' = 'v') BY 's' PROPERTIES ('timeout' = '1')`)
	d := stmt.DataDescs[0]
	if !d.TemporaryPart || len(d.Partition) != 2 || d.ColumnSeparator != "," || d.Format != "csv" ||
		len(d.FormatProps) != 2 || len(d.ColumnList) != 2 || len(d.ColumnsFromPath) != 1 ||
		len(d.ColumnMappings) != 2 || d.Where == nil {
		t.Errorf("desc 0 = %+v", d)
	}
	d = stmt.DataDescs[1]
	if d.SourceTable == nil || !d.Negative || len(d.Partition) != 1 {
		t.Errorf("desc 1 = %+v", d)
	}
	if stmt.StorageType != "BROKER" || stmt.BrokerName != "" || len(stmt.BrokerProperties) != 1 || stmt.System != "s" {
		t.Errorf("stmt = %+v", stmt)
	}
}

// TestLoadGrammar holds LOAD to StarRocks.g4's loadStatement;
// testdata/compat.tsv has StarRocks's verdicts.
func TestLoadGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"LOAD LABEL a", true},
		{"LOAD LABEL a WITH BROKER", true},
		{"LOAD LABEL a WITH RESOURCE a ('a' = 'a')", true},
		{"LOAD LABEL a WITH RESOURCE a BY a", false},
		{"LOAD LABEL a(DATA INFILE() INTO TABLE a)", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a COLUMNS FROM PATH AS())", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a())", true},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a(ESCAPE = 'a' ENCLOSE = 'a'))", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a(SKIP_HEADER = 'a'))", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a(SKIP_HEADER))", true},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a LINES TERMINATED BY 'a')", false},
		{"LOAD LABEL a(DATA FROM TABLE a INTO TABLE a COLUMNS TERMINATED BY 'a')", false},
	}
	for _, c := range cases {
		_, err := Parse(c.sql)
		if (err == nil) != c.ok {
			t.Errorf("Parse(%q): err = %v, want accepted=%v", c.sql, err, c.ok)
		}
	}
}

func TestLoadSimpleLabel(t *testing.T) {
	sql := `LOAD LABEL my_label (DATA INFILE('file1') INTO TABLE t)`

	stmt := mustParseLoadData(t, sql)
	if stmt.Label != "my_label" {
		t.Errorf("label = %q, want my_label", stmt.Label)
	}
}

// ---------------------------------------------------------------------------
// EXPORT tests
// ---------------------------------------------------------------------------

func TestExportBasic(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE my_table TO '/path/to/export/'`)
	if len(stmt.Target.Parts) < 1 || stmt.Target.Parts[len(stmt.Target.Parts)-1] != "my_table" {
		t.Fatalf("target = %v, want my_table", stmt.Target.Parts)
	}
	if stmt.Path != "/path/to/export/" {
		t.Errorf("path = %q, want /path/to/export/", stmt.Path)
	}
}

func TestExportWithPartition(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE t PARTITION(p1, p2) TO '/path/'`)
	if len(stmt.Partition) != 2 {
		t.Fatalf("partition = %v, want [p1 p2]", stmt.Partition)
	}
}

func TestExportWithWhere(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE t WHERE id > 100 TO '/path/'`)
	if stmt.Where == nil {
		t.Error("WHERE should not be nil")
	}
}

func TestExportWithProperties(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE t TO '/path/' PROPERTIES ("format"="csv", "column_separator"="|")`)
	if len(stmt.Properties) != 2 {
		t.Fatalf("properties len = %d, want 2", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "format" {
		t.Errorf("properties[0].key = %q, want format", stmt.Properties[0].Key)
	}
}

func TestExportWithBroker(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE t TO 'hdfs://path/' WITH BROKER my_broker`)
	if stmt.BrokerName != "my_broker" {
		t.Errorf("broker = %q, want my_broker", stmt.BrokerName)
	}
}

func TestExportWithBrokerAndProperties(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE t TO 'hdfs://path/' PROPERTIES ("format"="csv") WITH BROKER my_broker ("username"="user")`)
	if stmt.BrokerName != "my_broker" {
		t.Errorf("broker = %q, want my_broker", stmt.BrokerName)
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("properties len = %d, want 1", len(stmt.Properties))
	}
}

func TestExportQualifiedTable(t *testing.T) {
	stmt := mustParseExport(t,
		`EXPORT TABLE db.my_table TO '/path/'`)
	if len(stmt.Target.Parts) != 2 {
		t.Fatalf("target parts = %v, want 2 parts", stmt.Target.Parts)
	}
	if stmt.Target.Parts[0] != "db" || stmt.Target.Parts[1] != "my_table" {
		t.Errorf("target = %v, want [db my_table]", stmt.Target.Parts)
	}
}

// ---------------------------------------------------------------------------
// NodeTag tests
// ---------------------------------------------------------------------------

func TestT6NodeTags(t *testing.T) {
	truncate := &ast.TruncateTableStmt{}
	if truncate.Tag() != ast.T_TruncateTableStmt {
		t.Errorf("TruncateTableStmt tag = %v", truncate.Tag())
	}

	copyInto := &ast.CopyIntoStmt{}
	if copyInto.Tag() != ast.T_CopyIntoStmt {
		t.Errorf("CopyIntoStmt tag = %v", copyInto.Tag())
	}

	loadDesc := &ast.LoadDataDesc{}
	if loadDesc.Tag() != ast.T_LoadDataDesc {
		t.Errorf("LoadDataDesc tag = %v", loadDesc.Tag())
	}

	load := &ast.LoadDataStmt{}
	if load.Tag() != ast.T_LoadDataStmt {
		t.Errorf("LoadDataStmt tag = %v", load.Tag())
	}

	export := &ast.ExportStmt{}
	if export.Tag() != ast.T_ExportStmt {
		t.Errorf("ExportStmt tag = %v", export.Tag())
	}
}
