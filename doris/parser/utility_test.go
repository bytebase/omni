package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// DESCRIBE / DESC
// ---------------------------------------------------------------------------

func TestDescribeTable(t *testing.T) {
	file, errs := parseForTest("DESC test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.DescribeStmt)
	if !ok {
		t.Fatalf("expected *ast.DescribeStmt, got %T", file.Stmts[0])
	}
	if n.Target == nil || n.Target.Parts[0] != "test_table" {
		t.Errorf("Target = %v, want test_table", n.Target)
	}
	if n.Full {
		t.Error("Full should be false")
	}
}

func TestDescribeKeyword(t *testing.T) {
	file, errs := parseForTest("DESCRIBE test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if n.Target == nil || n.Target.Parts[0] != "test_table" {
		t.Errorf("Target = %v, want test_table", n.Target)
	}
}

func TestDescribeFullTable(t *testing.T) {
	file, errs := parseForTest("DESCRIBE FULL test_table")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.Full {
		t.Error("Full should be true")
	}
}

func TestDescribeAllVerbose(t *testing.T) {
	file, errs := parseForTest("DESCRIBE test_table ALL VERBOSE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.AllVerbose {
		t.Error("AllVerbose should be true")
	}
}

func TestDescribeLocValid(t *testing.T) {
	file, errs := parseForTest("DESC t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.DescribeStmt)
	if !n.Loc.IsValid() {
		t.Error("DescribeStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// EXPLAIN
// ---------------------------------------------------------------------------

func TestExplainSelect(t *testing.T) {
	file, errs := parseForTest("EXPLAIN SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("expected *ast.ExplainStmt, got %T", file.Stmts[0])
	}
	if n.Type != "" {
		t.Errorf("Type = %q, want empty", n.Type)
	}
	if n.Query == nil {
		t.Error("Query should not be nil")
	}
}

func TestExplainVerbose(t *testing.T) {
	file, errs := parseForTest("EXPLAIN VERBOSE SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "VERBOSE" {
		t.Errorf("Type = %q, want VERBOSE", n.Type)
	}
}

func TestExplainGraph(t *testing.T) {
	file, errs := parseForTest("EXPLAIN GRAPH SELECT 1 FROM t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "GRAPH" {
		t.Errorf("Type = %q, want GRAPH", n.Type)
	}
}

func TestExplainPlan(t *testing.T) {
	file, errs := parseForTest("EXPLAIN PLAN SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "PLAN" {
		t.Errorf("Type = %q, want PLAN", n.Type)
	}
}

func TestExplainShape(t *testing.T) {
	file, errs := parseForTest("EXPLAIN SHAPE SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "SHAPE" {
		t.Errorf("Type = %q, want SHAPE", n.Type)
	}
}

func TestExplainMemo(t *testing.T) {
	file, errs := parseForTest("EXPLAIN MEMO SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.ExplainStmt)
	if n.Type != "MEMO" {
		t.Errorf("Type = %q, want MEMO", n.Type)
	}
}

// ---------------------------------------------------------------------------
// USE
// ---------------------------------------------------------------------------

func TestUseDatabase(t *testing.T) {
	file, errs := parseForTest("USE demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.UseStmt)
	if !ok {
		t.Fatalf("expected *ast.UseStmt, got %T", file.Stmts[0])
	}
	if n.Database != "demo" {
		t.Errorf("Database = %q, want demo", n.Database)
	}
	if n.Catalog != "" {
		t.Errorf("Catalog = %q, want empty", n.Catalog)
	}
}

func TestUseCatalogDB(t *testing.T) {
	file, errs := parseForTest("USE hms_catalog.demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if n.Catalog != "hms_catalog" {
		t.Errorf("Catalog = %q, want hms_catalog", n.Catalog)
	}
	if n.Database != "demo" {
		t.Errorf("Database = %q, want demo", n.Database)
	}
}

func TestUseDBAtCluster(t *testing.T) {
	file, errs := parseForTest("USE mydb@mycluster")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if n.Database != "mydb" {
		t.Errorf("Database = %q, want mydb", n.Database)
	}
	if n.Cluster != "mycluster" {
		t.Errorf("Cluster = %q, want mycluster", n.Cluster)
	}
}

func TestUseLocValid(t *testing.T) {
	file, errs := parseForTest("USE demo")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UseStmt)
	if !n.Loc.IsValid() {
		t.Error("UseStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// SET (generic variable assignment)
// ---------------------------------------------------------------------------

func TestSetVariable(t *testing.T) {
	file, errs := parseForTest(`SET time_zone = "Asia/Shanghai"`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.SetStmt)
	if !ok {
		t.Fatalf("expected *ast.SetStmt, got %T", file.Stmts[0])
	}
	if n.Type != "VARIABLE" {
		t.Errorf("Type = %q, want VARIABLE", n.Type)
	}
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	if n.Items[0].Name != "time_zone" {
		t.Errorf("Name = %q, want time_zone", n.Items[0].Name)
	}
}

func TestSetGlobalVariable(t *testing.T) {
	file, errs := parseForTest("SET GLOBAL exec_mem_limit = 137438953472")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	item := n.Items[0]
	if item.Scope != "GLOBAL" {
		t.Errorf("Scope = %q, want GLOBAL", item.Scope)
	}
	if item.Name != "exec_mem_limit" {
		t.Errorf("Name = %q, want exec_mem_limit", item.Name)
	}
}

func TestSetDoubleAtVariable(t *testing.T) {
	file, errs := parseForTest("SET @@exec_mem_limit = 137438953472")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if len(n.Items) != 1 {
		t.Fatalf("Items count = %d, want 1", len(n.Items))
	}
	item := n.Items[0]
	if item.Name != "exec_mem_limit" {
		t.Errorf("Name = %q, want exec_mem_limit", item.Name)
	}
}

func TestSetNames(t *testing.T) {
	file, errs := parseForTest("SET NAMES 'utf8'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "NAMES" {
		t.Errorf("Type = %q, want NAMES", n.Type)
	}
	if len(n.Items) != 1 || n.Items[0].Raw != "utf8" {
		t.Errorf("Items[0].Raw = %q, want utf8", n.Items[0].Raw)
	}
}

func TestSetCharset(t *testing.T) {
	file, errs := parseForTest("SET CHARSET 'utf8mb4'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "CHARSET" {
		t.Errorf("Type = %q, want CHARSET", n.Type)
	}
}

func TestSetTransaction(t *testing.T) {
	file, errs := parseForTest("SET TRANSACTION READ ONLY")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if n.Type != "TRANSACTION" {
		t.Errorf("Type = %q, want TRANSACTION", n.Type)
	}
}

func TestSetLocValid(t *testing.T) {
	file, errs := parseForTest("SET x = 1")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.SetStmt)
	if !n.Loc.IsValid() {
		t.Error("SetStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// UNSET
// ---------------------------------------------------------------------------

func TestUnsetVariable(t *testing.T) {
	file, errs := parseForTest("UNSET VARIABLE myvar")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.UnsetStmt)
	if !ok {
		t.Fatalf("expected *ast.UnsetStmt, got %T", file.Stmts[0])
	}
	if n.Type != "VARIABLE" {
		t.Errorf("Type = %q, want VARIABLE", n.Type)
	}
	if len(n.Names) != 1 || n.Names[0] != "myvar" {
		t.Errorf("Names = %v, want [myvar]", n.Names)
	}
}

func TestUnsetVariableAll(t *testing.T) {
	file, errs := parseForTest("UNSET VARIABLE ALL")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if !n.All {
		t.Error("All should be true")
	}
}

func TestUnsetGlobalVariable(t *testing.T) {
	file, errs := parseForTest("UNSET GLOBAL VARIABLE myvar")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if n.Scope != "GLOBAL" {
		t.Errorf("Scope = %q, want GLOBAL", n.Scope)
	}
	if len(n.Names) != 1 || n.Names[0] != "myvar" {
		t.Errorf("Names = %v, want [myvar]", n.Names)
	}
}

func TestUnsetSessionVariable(t *testing.T) {
	file, errs := parseForTest("UNSET SESSION VARIABLE x")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if n.Scope != "SESSION" {
		t.Errorf("Scope = %q, want SESSION", n.Scope)
	}
}

func TestUnsetLocValid(t *testing.T) {
	file, errs := parseForTest("UNSET VARIABLE x")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.UnsetStmt)
	if !n.Loc.IsValid() {
		t.Error("UnsetStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// HELP
// ---------------------------------------------------------------------------

func TestHelp(t *testing.T) {
	file, errs := parseForTest("HELP 'SELECT'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n, ok := file.Stmts[0].(*ast.HelpStmt)
	if !ok {
		t.Fatalf("expected *ast.HelpStmt, got %T", file.Stmts[0])
	}
	if n.Mask != "SELECT" {
		t.Errorf("Mask = %q, want SELECT", n.Mask)
	}
}

func TestHelpIdentifier(t *testing.T) {
	file, errs := parseForTest("HELP contents")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.HelpStmt)
	if n.Mask != "contents" {
		t.Errorf("Mask = %q, want contents", n.Mask)
	}
}

func TestHelpLocValid(t *testing.T) {
	file, errs := parseForTest("HELP 'x'")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	n := file.Stmts[0].(*ast.HelpStmt)
	if !n.Loc.IsValid() {
		t.Error("HelpStmt.Loc should be valid")
	}
}

// ---------------------------------------------------------------------------
// NodeTag checks
// ---------------------------------------------------------------------------

func TestUtilityNodeTags(t *testing.T) {
	tests := []struct {
		input   string
		wantTag ast.NodeTag
	}{
		{"SHOW TABLES", ast.T_ShowStmt},
		{"DESC t", ast.T_DescribeStmt},
		{"DESCRIBE t", ast.T_DescribeStmt},
		{"EXPLAIN SELECT 1", ast.T_ExplainStmt},
		{"USE db", ast.T_UseStmt},
		{"SET x = 1", ast.T_SetStmt},
		{"UNSET VARIABLE x", ast.T_UnsetStmt},
		{"HELP 'x'", ast.T_HelpStmt},
	}
	for _, tt := range tests {
		file, errs := parseForTest(tt.input)
		if len(errs) != 0 {
			t.Errorf("parseForTest(%q) errors: %v", tt.input, errs)
			continue
		}
		if len(file.Stmts) == 0 {
			t.Errorf("parseForTest(%q): no stmts", tt.input)
			continue
		}
		if got := file.Stmts[0].Tag(); got != tt.wantTag {
			t.Errorf("parseForTest(%q): Tag() = %v, want %v", tt.input, got, tt.wantTag)
		}
	}
}

// TestDescribeFunction: describe_command FUNCTION table_valued_function_ref
// (DorisLegacyParser.g4:1631; DorisParser.g4:948 derives it too).
func TestDescribeFunction(t *testing.T) {
	stmt := parseOne(t, "DESCRIBE FUNCTION a()").(*ast.DescribeStmt)
	if !stmt.Function {
		t.Error("Function should be true")
	}
	if stmt.Target == nil || len(stmt.Target.Parts) != 1 || stmt.Target.Parts[0] != "a" {
		t.Errorf("Target = %+v, want a", stmt.Target)
	}
}

func TestDescribeFunctionArgsAndAlias(t *testing.T) {
	stmt := parseOne(t, "DESCRIBE FUNCTION a('k' = 'v') zzz").(*ast.DescribeStmt)
	if len(stmt.Args) != 1 || stmt.Args[0].Key != "k" || stmt.Args[0].Value != "v" {
		t.Errorf("Args = %+v, want k=v", stmt.Args)
	}
	if stmt.TableAlias != "zzz" {
		t.Errorf("TableAlias = %q, want zzz", stmt.TableAlias)
	}
}

// TestDescribePartitions: describe_stmt's opt_partition_names
// (DorisLegacyParser.g4:1629, DorisParser.g4:951 specifiedPartition).
func TestDescribePartitions(t *testing.T) {
	tests := []struct {
		sql       string
		partition []string
		temporary bool
	}{
		{"DESCRIBE a PARTITION(a)", []string{"a"}, false},
		{"DESCRIBE a PARTITION(a, b)", []string{"a", "b"}, false},
		{"DESCRIBE a PARTITION a", []string{"a"}, false},
		{"DESCRIBE a PARTITIONS(a)", []string{"a"}, false},
		{"DESCRIBE a TEMPORARY PARTITION(a)", []string{"a"}, true},
		{"DESCRIBE a TEMPORARY PARTITIONS(a)", []string{"a"}, true},
		{"DESCRIBE a TEMPORARY PARTITION a", []string{"a"}, true},
	}
	for _, tc := range tests {
		stmt := parseOne(t, tc.sql).(*ast.DescribeStmt)
		if len(stmt.Partition) != len(tc.partition) {
			t.Fatalf("%q: Partition = %q, want %q", tc.sql, stmt.Partition, tc.partition)
		}
		for i := range tc.partition {
			if stmt.Partition[i] != tc.partition[i] {
				t.Errorf("%q: Partition = %q, want %q", tc.sql, stmt.Partition, tc.partition)
			}
		}
		if stmt.TempPartition != tc.temporary {
			t.Errorf("%q: TempPartition = %v, want %v", tc.sql, stmt.TempPartition, tc.temporary)
		}
	}
}

// TestDescribePartitionStarRejected: specifiedPartition (DorisParser.g4:1627)
// has no star form, so the engine refuses DESCRIBE t PARTITION(*), though the
// CUP parser's partition_names (:2311) alone would take it.
func TestDescribePartitionStarRejected(t *testing.T) {
	for _, sql := range []string{"DESCRIBE a PARTITION(*)", "DESCRIBE a PARTITIONS(*)"} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) succeeded, want a syntax error", sql)
		}
	}
}

// TestDescribeQueryIsExplain: DESCRIBE names an explain command when a query
// or an INSERT follows (DorisParser.g4:52 statementBase is explain? query,
// :1067 explainCommand is EXPLAIN, DESC or DESCRIBE).
func TestDescribeQueryIsExplain(t *testing.T) {
	stmt := parseOne(t, "DESCRIBE VERBOSE SELECT 1").(*ast.ExplainStmt)
	if stmt.Type != "VERBOSE" {
		t.Errorf("Type = %q, want VERBOSE", stmt.Type)
	}
	if _, ok := stmt.Query.(*ast.SelectStmt); !ok {
		t.Errorf("Query = %T, want *ast.SelectStmt", stmt.Query)
	}

	stmt = parseOne(t, "DESCRIBE INSERT INTO a SELECT 1").(*ast.ExplainStmt)
	if _, ok := stmt.Query.(*ast.InsertStmt); !ok {
		t.Errorf("Query = %T, want *ast.InsertStmt", stmt.Query)
	}
}

// TestDescribeInsertNeedsInto: the explain command states INSERT (INTO |
// OVERWRITE TABLE) (DorisParser.g4:139, DorisLegacyParser.g4:1777), so the
// bare head parseInsert tolerates is refused here.
func TestDescribeInsertNeedsInto(t *testing.T) {
	if _, err := Parse("DESCRIBE INSERT a SELECT 1"); err == nil {
		t.Error("Parse(DESCRIBE INSERT a SELECT 1) succeeded, want a syntax error")
	}
}

// tableAlias (:1393) carries an optional identifier list, and the name after
// AS is required.
func TestDescribeFunctionTableAliasColumns(t *testing.T) {
	stmt := parseOne(t, "DESCRIBE FUNCTION a() AS v(c)").(*ast.DescribeStmt)
	if stmt.TableAlias != "v" || len(stmt.TableColumns) != 1 || stmt.TableColumns[0] != "c" {
		t.Errorf("alias = %q columns = %q, want v(c)", stmt.TableAlias, stmt.TableColumns)
	}
	if _, err := Parse("DESCRIBE FUNCTION a() AS"); err == nil {
		t.Error("Parse(DESCRIBE FUNCTION a() AS) succeeded, want a syntax error")
	}
}

// FULL scopes to describe_stmt's table-name alternative, so it does not prefix
// FUNCTION.
func TestDescribeFullFunctionRejected(t *testing.T) {
	if _, err := Parse("DESCRIBE FULL FUNCTION a()"); err == nil {
		t.Error("Parse(DESCRIBE FULL FUNCTION a()) succeeded, want a syntax error")
	}
}

// explain : explainCommand planType? level? PROCESS? (:1061), so DESCRIBE takes
// the same options EXPLAIN does.
func TestDescribeExplainOptions(t *testing.T) {
	tests := []struct {
		sql      string
		wantType string
	}{
		{"DESCRIBE TREE SELECT 1", "TREE"},
		{"DESCRIBE PROCESS SELECT 1", "PROCESS"},
		{"DESCRIBE OPTIMIZED VERBOSE SELECT 1", "OPTIMIZED VERBOSE"},
		{"EXPLAIN PLAN PROCESS SELECT 1", "PLAN PROCESS"},
	}
	for _, tc := range tests {
		stmt := parseOne(t, tc.sql).(*ast.ExplainStmt)
		if stmt.Type != tc.wantType {
			t.Errorf("Parse(%q): Type = %q, want %q", tc.sql, stmt.Type, tc.wantType)
		}
	}
}

// A partition clause needs the table name describe_stmt states first.
// (TEMPORARY is no test for this: it is non-reserved, so DESCRIBE TEMPORARY
// PARTITION(a) names a table TEMPORARY.)
func TestDescribePartitionNeedsTarget(t *testing.T) {
	if _, err := Parse("DESCRIBE PARTITION a"); err == nil {
		t.Error("Parse(DESCRIBE PARTITION a) succeeded, want a syntax error")
	}
	stmt := parseOne(t, "DESCRIBE TEMPORARY PARTITION(a)").(*ast.DescribeStmt)
	if stmt.Target == nil || stmt.Target.Parts[0] != "TEMPORARY" || len(stmt.Partition) != 1 {
		t.Errorf("Target = %+v Partition = %q, want table TEMPORARY partition a", stmt.Target, stmt.Partition)
	}
}

// explain gives PROCESS one slot after the level.
func TestExplainOneProcess(t *testing.T) {
	if _, err := Parse("EXPLAIN PLAN PROCESS PROCESS SELECT 1"); err == nil {
		t.Error("Parse(EXPLAIN PLAN PROCESS PROCESS SELECT 1) succeeded, want a syntax error")
	}
}
