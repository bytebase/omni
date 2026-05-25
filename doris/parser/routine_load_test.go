package parser

import (
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestCreateRoutineLoad_MinimalKafka(t *testing.T) {
	sql := `CREATE ROUTINE LOAD mydb.job1 ON tbl1
FROM KAFKA (
    "kafka_broker_list" = "broker1:9092",
    "kafka_topic" = "my_topic"
)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
	if stmt.OnTable == nil || stmt.OnTable.String() != "tbl1" {
		t.Errorf("OnTable = %v, want tbl1", stmt.OnTable)
	}
	if stmt.DataSourceType != "KAFKA" {
		t.Errorf("DataSourceType = %q, want KAFKA", stmt.DataSourceType)
	}
	if len(stmt.DataSourceProperties) != 2 {
		t.Errorf("DataSourceProperties count = %d, want 2", len(stmt.DataSourceProperties))
	}
}

func TestCreateRoutineLoad_WithProperties(t *testing.T) {
	sql := `CREATE ROUTINE LOAD job2 ON tbl2
PROPERTIES (
    "desired_concurrent_number" = "3",
    "max_error_number" = "1000"
)
FROM KAFKA (
    "kafka_broker_list" = "broker1:9092",
    "kafka_topic" = "topic2"
)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if len(stmt.JobProperties) != 2 {
		t.Errorf("JobProperties count = %d, want 2", len(stmt.JobProperties))
	}
	if stmt.JobProperties[0].Key != "desired_concurrent_number" {
		t.Errorf("JobProperties[0].Key = %q, want desired_concurrent_number", stmt.JobProperties[0].Key)
	}
	if stmt.DataSourceType != "KAFKA" {
		t.Errorf("DataSourceType = %q, want KAFKA", stmt.DataSourceType)
	}
}

func TestCreateRoutineLoad_WithComment(t *testing.T) {
	sql := `CREATE ROUTINE LOAD job3 ON tbl3
FROM KAFKA (
    "kafka_broker_list" = "broker1:9092",
    "kafka_topic" = "t"
)
COMMENT "this is a test job"`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Comment != "this is a test job" {
		t.Errorf("Comment = %q, want %q", stmt.Comment, "this is a test job")
	}
}

func TestCreateRoutineLoad_Tag(t *testing.T) {
	sql := `CREATE ROUTINE LOAD job4 ON tbl4 FROM KAFKA ("kafka_broker_list"="b:9092", "kafka_topic"="t")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_CreateRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

func TestCreateRoutineLoad_S3Source(t *testing.T) {
	sql := `CREATE ROUTINE LOAD job5 ON tbl5
FROM S3 (
    "s3_bucket" = "mybucket",
    "s3_region" = "us-east-1"
)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.CreateRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.DataSourceType != "S3" {
		t.Errorf("DataSourceType = %q, want S3", stmt.DataSourceType)
	}
}

// ---------------------------------------------------------------------------
// ALTER ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestAlterRoutineLoad_PropertiesOnly(t *testing.T) {
	sql := `ALTER ROUTINE LOAD FOR mydb.job1
PROPERTIES (
    "desired_concurrent_number" = "5"
)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties count = %d, want 1", len(stmt.Properties))
	}
}

func TestAlterRoutineLoad_WithKafka(t *testing.T) {
	sql := `ALTER ROUTINE LOAD FOR job1
FROM KAFKA (
    "kafka_partitions" = "0,1,2",
    "kafka_offsets" = "100,200,300"
)`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.DataSourceType != "KAFKA" {
		t.Errorf("DataSourceType = %q, want KAFKA", stmt.DataSourceType)
	}
	if len(stmt.DataSourceProperties) != 2 {
		t.Errorf("DataSourceProperties count = %d, want 2", len(stmt.DataSourceProperties))
	}
}

func TestAlterRoutineLoad_Tag(t *testing.T) {
	sql := `ALTER ROUTINE LOAD FOR job1 PROPERTIES ("max_error_number"="100")`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_AlterRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// PAUSE ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestPauseRoutineLoad_Basic(t *testing.T) {
	sql := `PAUSE ROUTINE LOAD FOR mydb.job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.PauseRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.PauseRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
	if stmt.All {
		t.Error("All should be false")
	}
}

func TestPauseRoutineLoad_All(t *testing.T) {
	sql := `PAUSE ALL ROUTINE LOAD FOR mydb`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.PauseRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.PauseRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if !stmt.All {
		t.Error("All should be true")
	}
	if stmt.For != "mydb" {
		t.Errorf("For = %q, want mydb", stmt.For)
	}
}

func TestPauseRoutineLoad_AllNoFor(t *testing.T) {
	sql := `PAUSE ALL ROUTINE LOAD`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.PauseRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.PauseRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if !stmt.All {
		t.Error("All should be true")
	}
	if stmt.For != "" {
		t.Errorf("For = %q, want empty", stmt.For)
	}
}

func TestPauseRoutineLoad_Tag(t *testing.T) {
	sql := `PAUSE ROUTINE LOAD FOR job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_PauseRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_PauseRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// RESUME ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestResumeRoutineLoad_Basic(t *testing.T) {
	sql := `RESUME ROUTINE LOAD FOR mydb.job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ResumeRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ResumeRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
	if stmt.All {
		t.Error("All should be false")
	}
}

func TestResumeRoutineLoad_All(t *testing.T) {
	sql := `RESUME ALL ROUTINE LOAD FOR mydb`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ResumeRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ResumeRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if !stmt.All {
		t.Error("All should be true")
	}
	if stmt.For != "mydb" {
		t.Errorf("For = %q, want mydb", stmt.For)
	}
}

func TestResumeRoutineLoad_Tag(t *testing.T) {
	sql := `RESUME ROUTINE LOAD FOR job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_ResumeRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_ResumeRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// STOP ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestStopRoutineLoad_Basic(t *testing.T) {
	sql := `STOP ROUTINE LOAD FOR mydb.job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.StopRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.StopRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
}

func TestStopRoutineLoad_UnqualifiedName(t *testing.T) {
	sql := `STOP ROUTINE LOAD FOR job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.StopRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.StopRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "job1" {
		t.Errorf("Name = %v, want job1", stmt.Name)
	}
}

func TestStopRoutineLoad_Tag(t *testing.T) {
	sql := `STOP ROUTINE LOAD FOR job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_StopRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_StopRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SHOW ROUTINE LOAD
// ---------------------------------------------------------------------------

func TestShowRoutineLoad_Bare(t *testing.T) {
	sql := `SHOW ROUTINE LOAD`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.All {
		t.Error("All should be false")
	}
	if stmt.Name != nil {
		t.Errorf("Name should be nil, got %v", stmt.Name)
	}
}

func TestShowRoutineLoad_ForName(t *testing.T) {
	sql := `SHOW ROUTINE LOAD FOR mydb.job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "mydb.job1" {
		t.Errorf("Name = %v, want mydb.job1", stmt.Name)
	}
}

func TestShowRoutineLoad_All(t *testing.T) {
	sql := `SHOW ALL ROUTINE LOAD`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if !stmt.All {
		t.Error("All should be true")
	}
}

func TestShowRoutineLoad_Like(t *testing.T) {
	sql := `SHOW ROUTINE LOAD LIKE 'job%'`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Like != "job%" {
		t.Errorf("Like = %q, want job%%", stmt.Like)
	}
}

func TestShowRoutineLoad_FromDB(t *testing.T) {
	sql := `SHOW ROUTINE LOAD FROM mydb`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.From != "mydb" {
		t.Errorf("From = %q, want mydb", stmt.From)
	}
}

func TestShowRoutineLoad_Tag(t *testing.T) {
	sql := `SHOW ROUTINE LOAD FOR job1`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_ShowRoutineLoadStmt {
		t.Errorf("Tag() = %v, want T_ShowRoutineLoadStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SHOW ROUTINE LOAD TASK
// ---------------------------------------------------------------------------

func TestShowRoutineLoadTask_FromDB(t *testing.T) {
	sql := `SHOW ROUTINE LOAD TASK FROM mydb WHERE JobName = "job1"`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadTaskStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadTaskStmt, got %T", file.Stmts[0])
	}
	if stmt.From != "mydb" {
		t.Errorf("From = %q, want mydb", stmt.From)
	}
	if stmt.Where == nil {
		t.Error("Where should not be nil")
	}
}

func TestShowRoutineLoadTask_Tag(t *testing.T) {
	sql := `SHOW ROUTINE LOAD TASK FROM mydb WHERE JobName = "job1"`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_ShowRoutineLoadTaskStmt {
		t.Errorf("Tag() = %v, want T_ShowRoutineLoadTaskStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// SYNC
// ---------------------------------------------------------------------------

func TestSync_Basic(t *testing.T) {
	sql := `SYNC`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	_, ok := file.Stmts[0].(*ast.SyncStmt)
	if !ok {
		t.Fatalf("expected *ast.SyncStmt, got %T", file.Stmts[0])
	}
}

func TestSync_Tag(t *testing.T) {
	sql := `SYNC`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_SyncStmt {
		t.Errorf("Tag() = %v, want T_SyncStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus: the single test case from routine_load.sql
// ---------------------------------------------------------------------------

func TestShowCreateRoutineLoad_LegacyCorpus(t *testing.T) {
	// From doris/parser/testdata/legacy/regression/routine_load.sql
	sql := `show create routine load for uDA6TB9nkmLWHYWLFCdCP6XrykhxxNa4gXA9yxZJU0`
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse errors: %v", errs)
	}
	// This parses as SHOW ROUTINE LOAD FOR <name> — job name is the long identifier
	stmt, ok := file.Stmts[0].(*ast.ShowRoutineLoadStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowRoutineLoadStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil {
		t.Fatal("Name should not be nil")
	}
}
