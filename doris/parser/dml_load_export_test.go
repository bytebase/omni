package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
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
	stmt := mustParseCopyInto(t, "COPY INTO my_table FROM @stage")
	if len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "my_table" {
		t.Fatalf("target = %v, want [my_table]", stmt.Target.Parts)
	}
	if stmt.Stage != "@stage" {
		t.Errorf("stage = %q, want @stage", stmt.Stage)
	}
}

func TestCopyIntoWithProperties(t *testing.T) {
	stmt := mustParseCopyInto(t,
		`COPY INTO my_table FROM @stage PROPERTIES ("key"="val")`)
	if len(stmt.Properties) != 1 {
		t.Fatalf("properties len = %d, want 1", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" || stmt.Properties[0].Value != "val" {
		t.Errorf("property = %q=%q, want key=val", stmt.Properties[0].Key, stmt.Properties[0].Value)
	}
}

func TestCopyIntoWithFiles(t *testing.T) {
	// FILES and PATTERN are properties a load names inside PROPERTIES; a load
	// that states them beside the stage is not a statement (DorisParser.g4:936
	// copyInto has nowhere else to put them).
	if _, err := Parse(`COPY INTO my_table FROM @stage FILES = ('f1.csv')`); err == nil {
		t.Error("Parse(COPY INTO ... FILES = (...)) succeeded, want a syntax error")
	}
}

func TestCopyIntoColumns(t *testing.T) {
	stmt := mustParseCopyInto(t, "COPY INTO my_table(a, b) FROM @stage")
	if len(stmt.Columns) != 2 || stmt.Columns[0] != "a" || stmt.Columns[1] != "b" {
		t.Errorf("columns = %v, want [a b]", stmt.Columns)
	}
}

// COPY INTO a FROM a names no stage: the source starts with '@'
// (DorisParser.g4:942 stageAndPattern), and the engine refuses the bare word
// (compat.tsv's ..._FROM_(s...!drop-5 rows).
func TestCopyIntoNeedsStage(t *testing.T) {
	for _, sql := range []string{
		"COPY INTO a FROM a",
		"COPY INTO a FROM 'stage/path'",
		"COPY INTO a FROM @",
		"COPY INTO a FROM @a()",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
	for _, sql := range []string{
		"COPY INTO a FROM @a",
		"COPY INTO a FROM @~",
		"COPY INTO a FROM @a('glob*')",
		"COPY INTO a FROM (SELECT b FROM @a)",
		"COPY INTO a FROM (SELECT b FROM @a WHERE b > 1)",
	} {
		if _, err := Parse(sql); err != nil {
			t.Errorf("Parse(%q) = %v, want it accepted", sql, err)
		}
	}
}

func TestCopyIntoStageAndSelect(t *testing.T) {
	stmt := mustParseCopyInto(t, "COPY INTO t(c) FROM (SELECT b FROM @s('p*') WHERE b > 1)")
	if stmt.Stage != "@s" || stmt.Pattern != "p*" {
		t.Errorf("stage = %q pattern = %q, want @s and p*", stmt.Stage, stmt.Pattern)
	}
	if len(stmt.SelectItems) != 1 {
		t.Fatalf("select items = %d, want 1", len(stmt.SelectItems))
	}
	if stmt.Where == nil {
		t.Error("Where is nil, want the query's WHERE")
	}
}

// Stage keeps the source's spelling, so the tilde and a quoted identifier that
// spells one are told apart, and a quoted name keeps its backticks. HasPattern
// tells @stage(”) — an empty glob — from @stage.
func TestCopyIntoStageSpelling(t *testing.T) {
	tests := []struct {
		sql, stage, pattern string
		hasPattern          bool
	}{
		{"COPY INTO t FROM @~", "@~", "", false},
		{"COPY INTO t FROM @`~`", "@`~`", "", false},
		{"COPY INTO t FROM @stage", "@stage", "", false},
		{"COPY INTO t FROM @`a.b`", "@`a.b`", "", false},
		{"COPY INTO t FROM @stage('')", "@stage", "", true},
		{"COPY INTO t FROM @stage('p*')", "@stage", "p*", true},
	}
	for _, tc := range tests {
		stmt := mustParseCopyInto(t, tc.sql)
		if stmt.Stage != tc.stage || stmt.Pattern != tc.pattern || stmt.HasPattern != tc.hasPattern {
			t.Errorf("Parse(%q): Stage = %q Pattern = %q HasPattern = %v, want %q %q %v",
				tc.sql, stmt.Stage, stmt.Pattern, stmt.HasPattern, tc.stage, tc.pattern, tc.hasPattern)
		}
	}
}

// ---------------------------------------------------------------------------
// LOAD tests
// ---------------------------------------------------------------------------

func TestLoadBasic(t *testing.T) {
	sql := `LOAD LABEL example_db.label1
	(DATA INFILE('hdfs://hdfs_host:hdfs_port/file1') INTO TABLE my_table)
	WITH BROKER hdfs_broker
	PROPERTIES ("timeout"="1200")
	COMMENT 'test load'`

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
	if stmt.Comment != "test load" {
		t.Errorf("comment = %q, want 'test load'", stmt.Comment)
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

// TestLoadNegative: DorisParser.g4's dataDesc (:983) has no NEGATIVE.
func TestLoadNegative(t *testing.T) {
	for _, sql := range []string{
		"LOAD LABEL load3 (NEGATIVE DATA INFILE('file1') INTO TABLE t)",
		"LOAD LABEL load3 (DATA INFILE('file1') NEGATIVE INTO TABLE t)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded; dataDesc has no NEGATIVE", sql)
		}
	}
}

func TestLoadWithColumns(t *testing.T) {
	sql := `LOAD LABEL load4
	(DATA INFILE('file1') INTO TABLE t (c1, c2, c3) COLUMNS FROM PATH AS (p)
	 SET (c1 = c2 + 1) PRECEDING FILTER c1 > 0 WHERE c2 = 1 DELETE ON c3 = 0 ORDER BY c1)`

	stmt := mustParseLoadData(t, sql)
	desc := stmt.DataDescs[0]
	if len(desc.ColumnList) != 3 || len(desc.ColumnsFromPath) != 1 || len(desc.ColumnMappings) != 1 ||
		desc.PrecedingFilter == nil || desc.Where == nil || desc.DeleteOn == nil || desc.SequenceColumn != "c1" {
		t.Errorf("desc = %+v", desc)
	}
}

func TestLoadFromTableAndMysql(t *testing.T) {
	stmt := mustParseLoadData(t, "LOAD LABEL l (WITH MERGE DATA FROM TABLE s INTO TABLE t PARTITION (p1) PROPERTIES (k = v)) WITH RESOURCE 'r' ('k' = 'v')")
	desc := stmt.DataDescs[0]
	if desc.MergeType != "MERGE" || desc.SourceTable == nil || len(desc.Partition) != 1 ||
		stmt.StorageType != "RESOURCE" || stmt.BrokerName != "r" {
		t.Errorf("stmt = %+v, desc = %+v", stmt, desc)
	}
	stmt = mustParseLoadData(t, "LOAD DATA LOCAL INFILE 'f' INTO TABLE db.t PARTITION (p) COLUMNS TERMINATED BY ',' IGNORE 1 LINES (a, b) SET (a = b) PROPERTIES ('k' = 'v') COMMENT 'c'")
	desc = stmt.DataDescs[0]
	if !desc.Local || desc.SkipLines != 1 || desc.ColumnSeparator != "," || len(desc.ColumnMappings) != 1 ||
		len(desc.Properties) != 1 || stmt.Comment != "c" || stmt.Label != "" {
		t.Errorf("stmt = %+v, desc = %+v", stmt, desc)
	}
}

// TestLoadGrammar holds LOAD to DorisParser.g4 (and, for LOAD DATA, the CUP
// parser's load_stmt); testdata/compat.tsv has Doris's verdicts.
func TestLoadGrammar(t *testing.T) {
	cases := []struct {
		sql string
		ok  bool
	}{
		{"LOAD LABEL a(DATA INFILE('a', 'a') INTO TABLE a)", true},
		{"LOAD LABEL a(DATA INFILE('a' 'a') INTO TABLE a)", false},
		{"LOAD LABEL a(DATA INFILE('a',) INTO TABLE a)", false},
		{"LOAD LABEL a(DATA INFILE() INTO TABLE a)", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a PARTITIONS a)", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a PARTITION a)", true},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a COLUMNS FROM PATH AS())", false},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a FORMAT AS a COMPRESS_TYPE AS a)", true},
		{"LOAD LABEL a(DATA INFILE('a') INTO TABLE a COMPRESS_TYPE AS a FORMAT AS a)", false},
		{"LOAD LABEL a(DATA FROM TABLE a INTO TABLE a) WITH S3 ('k' = 'v')", true},
		{"LOAD LABEL a(DATA FROM TABLE a INTO TABLE a) WITH S3", false},
		{"LOAD LABEL a(DATA FROM TABLE a INTO TABLE a ORDER BY a)", false},
		{"LOAD LABEL a(WITH DATA FROM TABLE a INTO TABLE a)", false},
		{"LOAD DATA INFILE 'a' INTO TABLE a IGNORE 1 ROWS", true},
		{"LOAD DATA INFILE 'a' INTO TABLE a PROPERTIES (a = a)", false},
		{"LOAD DATA INFILE ('a') INTO TABLE a", false},
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

func TestLoadLabelParts(t *testing.T) {
	// The parser keeps the label's parts as written. A quoted part may contain
	// a dot; the pinned engines reject such a label when they validate it, but
	// the parse still has to distinguish the part from a qualifier.
	for _, tc := range []struct {
		name      string
		sql       string
		wantLabel string
		wantParts []string
	}{
		{
			name:      "qualified",
			sql:       "LOAD LABEL db.job (DATA INFILE('file1') INTO TABLE t)",
			wantLabel: "db.job",
			wantParts: []string{"db", "job"},
		},
		{
			name:      "unqualified",
			sql:       "LOAD LABEL job (DATA INFILE('file1') INTO TABLE t)",
			wantLabel: "job",
			wantParts: []string{"job"},
		},
		{
			name:      "a quoted label part keeps its dot",
			sql:       "LOAD LABEL db.`job.with.dot` (DATA INFILE('file1') INTO TABLE t)",
			wantLabel: "db.job.with.dot",
			wantParts: []string{"db", "job.with.dot"},
		},
		{
			name:      "a quoted unqualified label keeps its dot",
			sql:       "LOAD LABEL `job.with.dot` (DATA INFILE('file1') INTO TABLE t)",
			wantLabel: "job.with.dot",
			wantParts: []string{"job.with.dot"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			stmt := mustParseLoadData(t, tc.sql)
			if stmt.Label != tc.wantLabel {
				t.Errorf("label = %q, want %q", stmt.Label, tc.wantLabel)
			}
			if len(stmt.LabelParts) != len(tc.wantParts) {
				t.Fatalf("label parts = %v, want %v", stmt.LabelParts, tc.wantParts)
			}
			for i := range tc.wantParts {
				if stmt.LabelParts[i] != tc.wantParts[i] {
					t.Errorf("label parts = %v, want %v", stmt.LabelParts, tc.wantParts)
					break
				}
			}
		})
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

// TestExportWithRemoteStorage pins withRemoteStorageSystem
// (DorisParser.g4:1121) as EXPORT takes it: WITH RESOURCE and WITH BROKER name
// a resource or broker and may leave the list out, and WITH S3, WITH HDFS and
// WITH LOCAL need one.
func TestExportWithRemoteStorage(t *testing.T) {
	tests := []struct {
		sql         string
		storageType string
		name        string
		props       int
	}{
		{`EXPORT TABLE t TO 'a' WITH RESOURCE r`, "RESOURCE", "r", 0},
		{`EXPORT TABLE t TO 'a' WITH RESOURCE r(a = a)`, "RESOURCE", "r", 1},
		{`EXPORT TABLE t TO 'a' WITH BROKER b`, "BROKER", "b", 0},
		{`EXPORT TABLE t TO 'a' WITH S3(a = a)`, "S3", "", 1},
		{`EXPORT TABLE t TO 'a' WITH S3('a' = 'a')`, "S3", "", 1},
		{`EXPORT TABLE t TO 'a' WITH HDFS(a = a)`, "HDFS", "", 1},
		{`EXPORT TABLE t TO 'a' WITH LOCAL('a' = 'a')`, "LOCAL", "", 1},
	}
	for _, c := range tests {
		stmt := mustParseExport(t, c.sql)
		if stmt.StorageType != c.storageType || stmt.BrokerName != c.name {
			t.Errorf("%s: storage = %q %q, want %q %q", c.sql, stmt.StorageType, stmt.BrokerName, c.storageType, c.name)
		}
		if len(stmt.BrokerProperties) != c.props {
			t.Errorf("%s: storage properties = %d, want %d", c.sql, len(stmt.BrokerProperties), c.props)
		}
	}

	for _, sql := range []string{
		`EXPORT TABLE t TO 'a' WITH S3`,
		`EXPORT TABLE t TO 'a' WITH HDFS`,
		`EXPORT TABLE t TO 'a' WITH LOCAL`,
		`EXPORT TABLE t TO 'a' WITH RESOURCE`,
		`EXPORT TABLE t TO 'a' WITH BROKER b(a = a`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want a rejection", sql)
		}
	}
}

// TestExportPropertyLists pins EXPORT's two lists to the grammar's full
// propertyItem (DorisParser.g4:1385), which takes a bare identifier as a key
// or a value, where most clauses take only a "key"="value" pair.
func TestExportPropertyLists(t *testing.T) {
	stmt := mustParseExport(t, `EXPORT TABLE t TO 'a' PROPERTIES(a = a)`)
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "a" || stmt.Properties[0].Value != "a" {
		t.Fatalf("properties = %v, want one a = a item", stmt.Properties)
	}

	for _, sql := range []string{
		`EXPORT TABLE t TO 'a' PROPERTIES()`,
		`EXPORT TABLE t TO 'a' PROPERTIES(a = a,)`,
		`EXPORT TABLE t TO 'a' PROPERTIES(a =)`,
		`EXPORT TABLE t TO 'a' PROPERTIES(a a)`,
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) = nil error, want a rejection", sql)
		}
	}
}

// TestExportPartitionKeyword pins EXPORT's keyword to PARTITION: both grammars
// Doris reads for it take PARTITION, so PARTITIONS is refused, unlike
// TRUNCATE's clause which takes both.
func TestExportPartitionKeyword(t *testing.T) {
	stmt := mustParseExport(t, `EXPORT TABLE t PARTITION(p1) TO 'a'`)
	if len(stmt.Partition) != 1 || stmt.Partition[0] != "p1" {
		t.Fatalf("partition = %v, want [p1]", stmt.Partition)
	}
	if _, err := Parse(`EXPORT TABLE t PARTITIONS(p1) TO 'a'`); err == nil {
		t.Error("Parse(EXPORT ... PARTITIONS(...)) = nil error, want a rejection")
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

func TestTruncateCommonHint(t *testing.T) {
	for _, sql := range []string{
		"TRUNCATE TABLE `/*+ 1 */`",
		"TRUNCATE TABLE a /*+ x, y */",
	} {
		if _, errs := parseForTest(sql); len(errs) != 0 {
			t.Errorf("parseForTest(%q) errors: %v", sql, errs)
		}
	}
	for _, sql := range []string{
		"TRUNCATE TABLE a /*+ 1 */",
		"TRUNCATE TABLE a /*+ x y */",
		"TRUNCATE TABLE a /*+ x",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("parseForTest(%q) accepted, want a hint syntax error", sql)
		}
	}
}
