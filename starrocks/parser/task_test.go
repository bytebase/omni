package parser

import (
	"testing"

	"github.com/bytebase/omni/starrocks/ast"
)

// TestSubmitTask covers submitTaskStatement (StarRocks.g4:625): the optional
// name, both task clauses, and the three statements it can submit.
func TestSubmitTask(t *testing.T) {
	tests := []struct {
		sql      string
		name     string
		props    int
		schedule *ast.TaskSchedule
		stmt     ast.NodeTag
	}{
		{"SUBMIT TASK AS INSERT INTO a SELECT CURRENT_USER", "", 0, nil, ast.T_InsertStmt},
		{"SUBMIT TASK zzz AS INSERT INTO a SELECT CURRENT_USER", "zzz", 0, nil, ast.T_InsertStmt},
		{"SUBMIT TASK db.job AS INSERT INTO a SELECT CURRENT_USER", "db.job", 0, nil, ast.T_InsertStmt},
		{"SUBMIT TASK properties AS INSERT INTO a SELECT CURRENT_USER", "properties", 0, nil, ast.T_InsertStmt},
		{"SUBMIT TASK schedule AS INSERT INTO a SELECT CURRENT_USER", "schedule", 0, nil, ast.T_InsertStmt},
		{"SUBMIT TASK PROPERTIES('a' = 'a') AS INSERT INTO a SELECT 1", "", 1, nil, ast.T_InsertStmt},
		{"SUBMIT TASK SCHEDULE START('20260101') EVERY(INTERVAL 10 SECOND) AS INSERT INTO a SELECT 1", "", 0,
			&ast.TaskSchedule{Start: "20260101", Interval: "10", Unit: "SECOND"}, ast.T_InsertStmt},
		{"SUBMIT TASK SCHEDULE EVERY(INTERVAL ((1)) DAY) AS INSERT INTO a SELECT 1", "", 0,
			&ast.TaskSchedule{Interval: "1", Unit: "DAY"}, ast.T_InsertStmt},
		{"SUBMIT TASK AS CREATE TABLE a AS SELECT CURRENT_USER", "", 0, nil, ast.T_CreateTableStmt},
		{"SUBMIT TASK AS CACHE SELECT CURRENT_USER FROM a", "", 0, nil, ast.T_DataCacheSelectStmt},
		{"SUBMIT TASK SCHEDULE START('2026-01-01') EVERY(INTERVAL 1 DAY) AS INSERT INTO a SELECT 1", "", 0,
			&ast.TaskSchedule{Start: "2026-01-01", Interval: "1", Unit: "DAY"}, ast.T_InsertStmt},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := parseForTest(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt, ok := file.Stmts[0].(*ast.SubmitTaskStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.SubmitTaskStmt", file.Stmts[0])
			}
			got := ""
			if stmt.Name != nil {
				got = stmt.Name.String()
			}
			if got != tc.name {
				t.Errorf("Name = %q, want %q", got, tc.name)
			}
			if len(stmt.Properties) != tc.props {
				t.Errorf("Properties = %d, want %d", len(stmt.Properties), tc.props)
			}
			if tc.schedule == nil {
				if stmt.Schedule != nil {
					t.Errorf("Schedule = %+v, want nil", stmt.Schedule)
				}
			} else if stmt.Schedule == nil ||
				stmt.Schedule.Start != tc.schedule.Start ||
				stmt.Schedule.Interval != tc.schedule.Interval ||
				stmt.Schedule.Unit != tc.schedule.Unit {
				t.Errorf("Schedule = %+v, want %+v", stmt.Schedule, tc.schedule)
			}
			if stmt.Statement == nil || stmt.Statement.Tag() != tc.stmt {
				t.Errorf("Statement = %v, want %v", stmt.Statement, tc.stmt)
			}
		})
	}
}

// The submitted INSERT keeps insertStatement's required INTO or OVERWRITE; the
// bare form parseInsert admits at the top level is not a task body.
func TestSubmitTaskRejectsBareInsert(t *testing.T) {
	if _, errs := parseForTest("SUBMIT TASK AS INSERT a SELECT CURRENT_USER"); len(errs) == 0 {
		t.Fatal("SUBMIT TASK AS INSERT a parsed, want rejected")
	}
}

// The submitted CREATE must be createTableAsSelectStatement, not any CREATE
// TABLE the general parser reads.
func TestSubmitTaskRejectsNonCTAS(t *testing.T) {
	for _, sql := range []string{
		"SUBMIT TASK AS CREATE TABLE a (x INT)",
		"SUBMIT TASK AS CREATE TABLE a LIKE b",
		"SUBMIT TASK AS CREATE EXTERNAL TABLE a AS SELECT 1",
		"SUBMIT TASK AS CREATE TABLE a ENGINE = OLAP AS SELECT 1",
		"SUBMIT TASK AS CREATE TABLE a ROLLUP (r (x)) AS SELECT 1",
		"SUBMIT TASK AS CREATE TABLE a BROKER PROPERTIES('x' = 'y') AS SELECT 1",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

// insertStatement takes an explainDesc prefix, so a submitted task can plan
// its INSERT; the other submitted forms cannot.
func TestSubmitTaskExplainedInsert(t *testing.T) {
	file, errs := parseForTest("SUBMIT TASK AS EXPLAIN INSERT INTO t SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.SubmitTaskStmt)
	explain, ok := stmt.Statement.(*ast.ExplainStmt)
	if !ok {
		t.Fatalf("Statement = %T, want *ast.ExplainStmt", stmt.Statement)
	}
	if _, ok := explain.Query.(*ast.InsertStmt); !ok {
		t.Errorf("explained query = %T, want *ast.InsertStmt", explain.Query)
	}

	for _, sql := range []string{
		"SUBMIT TASK AS EXPLAIN SELECT 1",
		"SUBMIT TASK AS DESC UPDATE t SET a = 1",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

func TestSubmitTaskRejectsBadSchedule(t *testing.T) {
	for _, sql := range []string{
		"SUBMIT TASK SCHEDULE START('not-a-date') EVERY(INTERVAL 10 SECOND) AS INSERT INTO t SELECT 1",
		"SUBMIT TASK SCHEDULE START('2026-13-01') EVERY(INTERVAL 10 SECOND) AS INSERT INTO t SELECT 1",
		"SUBMIT TASK SCHEDULE EVERY(INTERVAL 0x10 SECOND) AS INSERT INTO t SELECT 1",
		"SUBMIT TASK SCHEDULE EVERY(INTERVAL 0b10 SECOND) AS INSERT INTO t SELECT 1",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

func TestTaskNameArity(t *testing.T) {
	for _, sql := range []string{
		"SUBMIT TASK catalog.db.job AS INSERT INTO t SELECT 1",
		"DROP TASK catalog.db.job",
	} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed, want rejected", sql)
		}
	}
}

func TestDropTask(t *testing.T) {
	for _, tc := range []struct {
		sql   string
		force bool
	}{
		{"DROP TASK a", false},
		{"DROP TASK db.a FORCE", true},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := parseForTest(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt, ok := file.Stmts[0].(*ast.DropTaskStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.DropTaskStmt", file.Stmts[0])
			}
			if stmt.Name == nil || stmt.Name.String() == "" {
				t.Error("Name is empty")
			}
			if stmt.Force != tc.force {
				t.Errorf("Force = %v, want %v", stmt.Force, tc.force)
			}
		})
	}
}

// TestAlterLoad covers alterLoadStatement (StarRocks.g4:1499).
func TestAlterLoad(t *testing.T) {
	file, errs := parseForTest("ALTER LOAD FOR db.label PROPERTIES('a' = 'a')")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*ast.AlterLoadStmt)
	if !ok {
		t.Fatalf("got %T, want *ast.AlterLoadStmt", file.Stmts[0])
	}
	if stmt.Name == nil || stmt.Name.String() != "db.label" {
		t.Errorf("Name = %v, want db.label", stmt.Name)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "a" || stmt.Properties[0].Value != "a" {
		t.Errorf("Properties = %+v, want one a=a pair", stmt.Properties)
	}
}
