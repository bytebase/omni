package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

func parseCreateJobStmt(t *testing.T, sql string) *ast.CreateJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateJobStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseAlterJobStmt(t *testing.T, sql string) *ast.AlterJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterJobStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseDropJobStmt(t *testing.T, sql string) *ast.DropJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropJobStmt)
	if !ok {
		t.Fatalf("expected *ast.DropJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parsePauseJobStmt(t *testing.T, sql string) *ast.PauseJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.PauseJobStmt)
	if !ok {
		t.Fatalf("expected *ast.PauseJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseResumeJobStmt(t *testing.T, sql string) *ast.ResumeJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ResumeJobStmt)
	if !ok {
		t.Fatalf("expected *ast.ResumeJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseCancelTaskStmt(t *testing.T, sql string) *ast.CancelTaskStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CancelTaskStmt)
	if !ok {
		t.Fatalf("expected *ast.CancelTaskStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseShowJobStmt(t *testing.T, sql string) *ast.ShowJobStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowJobStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowJobStmt, got %T", file.Stmts[0])
	}
	return stmt
}

func parseShowJobTaskStmt(t *testing.T, sql string) *ast.ShowJobTaskStmt {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.ShowJobTaskStmt)
	if !ok {
		t.Fatalf("expected *ast.ShowJobTaskStmt, got %T", file.Stmts[0])
	}
	return stmt
}

// ---------------------------------------------------------------------------
// CREATE JOB
// ---------------------------------------------------------------------------

func TestCreateJobEveryMinute(t *testing.T) {
	// From legacy corpus: CREATE JOB my_job ON SCHEDULE EVERY 1 MINUTE DO INSERT INTO ...
	stmt := parseCreateJobStmt(t, "CREATE JOB my_job ON SCHEDULE EVERY 1 MINUTE DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2")
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.JobType != "SCHEDULE" {
		t.Errorf("JobType: got %q, want SCHEDULE", stmt.JobType)
	}
	if stmt.IfNotExists {
		t.Errorf("IfNotExists: got true, want false")
	}
	if stmt.Schedule == nil {
		t.Fatalf("Schedule: got nil, want non-nil")
	}
	if stmt.Schedule.Every != "1 MINUTE" {
		t.Errorf("Schedule.Every: got %q, want %q", stmt.Schedule.Every, "1 MINUTE")
	}
	if stmt.Schedule.At != "" {
		t.Errorf("Schedule.At: got %q, want empty", stmt.Schedule.At)
	}
	if stmt.DoStmt == nil {
		t.Errorf("DoStmt: got nil, want non-nil")
	}
}

func TestCreateJobAt(t *testing.T) {
	// From legacy corpus: CREATE JOB my_job ON SCHEDULE AT '2020-01-01 00:00:00' DO INSERT INTO ...
	stmt := parseCreateJobStmt(t, "CREATE JOB my_job ON SCHEDULE AT '2020-01-01 00:00:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2")
	if stmt.JobType != "SCHEDULE" {
		t.Errorf("JobType: got %q, want SCHEDULE", stmt.JobType)
	}
	if stmt.Schedule == nil {
		t.Fatalf("Schedule: got nil, want non-nil")
	}
	if stmt.Schedule.At != "2020-01-01 00:00:00" {
		t.Errorf("Schedule.At: got %q, want %q", stmt.Schedule.At, "2020-01-01 00:00:00")
	}
	if stmt.Schedule.Every != "" {
		t.Errorf("Schedule.Every: got %q, want empty", stmt.Schedule.Every)
	}
	if stmt.DoStmt == nil {
		t.Errorf("DoStmt: got nil, want non-nil")
	}
}

func TestCreateJobEveryWithStarts(t *testing.T) {
	// From legacy corpus: CREATE JOB my_job ON SCHEDULE EVERY 1 DAY STARTS '...' DO INSERT INTO ...
	stmt := parseCreateJobStmt(t,
		"CREATE JOB my_job ON SCHEDULE EVERY 1 DAY STARTS '2020-01-01 00:00:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2 WHERE create_time >= days_add(now(),-1)")
	if stmt.Schedule == nil {
		t.Fatalf("Schedule: got nil, want non-nil")
	}
	if stmt.Schedule.Every != "1 DAY" {
		t.Errorf("Schedule.Every: got %q, want %q", stmt.Schedule.Every, "1 DAY")
	}
	if stmt.Schedule.Starts != "2020-01-01 00:00:00" {
		t.Errorf("Schedule.Starts: got %q, want %q", stmt.Schedule.Starts, "2020-01-01 00:00:00")
	}
	if stmt.Schedule.Ends != "" {
		t.Errorf("Schedule.Ends: got %q, want empty", stmt.Schedule.Ends)
	}
}

func TestCreateJobEveryWithStartsAndEnds(t *testing.T) {
	// From legacy corpus: CREATE JOB ... ON SCHEDULE EVERY 1 DAY STARTS '...' ENDS '...' DO INSERT INTO ...
	stmt := parseCreateJobStmt(t,
		"CREATE JOB my_job ON SCHEDULE EVERY 1 DAY STARTS '2020-01-01 00:00:00' ENDS '2020-01-01 00:10:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2 WHERE create_time >= days_add(now(),-1)")
	if stmt.Schedule == nil {
		t.Fatalf("Schedule: got nil, want non-nil")
	}
	if stmt.Schedule.Every != "1 DAY" {
		t.Errorf("Schedule.Every: got %q, want %q", stmt.Schedule.Every, "1 DAY")
	}
	if stmt.Schedule.Starts != "2020-01-01 00:00:00" {
		t.Errorf("Schedule.Starts: got %q, want %q", stmt.Schedule.Starts, "2020-01-01 00:00:00")
	}
	if stmt.Schedule.Ends != "2020-01-01 00:10:00" {
		t.Errorf("Schedule.Ends: got %q, want %q", stmt.Schedule.Ends, "2020-01-01 00:10:00")
	}
}

func TestCreateJobIfNotExists(t *testing.T) {
	stmt := parseCreateJobStmt(t, "CREATE JOB IF NOT EXISTS my_job ON SCHEDULE EVERY 1 HOUR DO INSERT INTO t1 SELECT * FROM t2")
	if !stmt.IfNotExists {
		t.Errorf("IfNotExists: got false, want true")
	}
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.Schedule == nil {
		t.Fatalf("Schedule is nil")
	}
	if stmt.Schedule.Every != "1 HOUR" {
		t.Errorf("Schedule.Every: got %q, want %q", stmt.Schedule.Every, "1 HOUR")
	}
}

func TestCreateJobWithComment(t *testing.T) {
	stmt := parseCreateJobStmt(t, "CREATE JOB my_job ON SCHEDULE EVERY 1 DAY COMMENT 'daily load' DO INSERT INTO t1 SELECT * FROM t2")
	if stmt.Comment != "daily load" {
		t.Errorf("Comment: got %q, want %q", stmt.Comment, "daily load")
	}
}

func TestCreateJobStreaming(t *testing.T) {
	stmt := parseCreateJobStmt(t, "CREATE JOB my_job AS STREAMING DO INSERT INTO t1 SELECT * FROM t2")
	if stmt.JobType != "STREAMING" {
		t.Errorf("JobType: got %q, want STREAMING", stmt.JobType)
	}
	if stmt.Schedule != nil {
		t.Errorf("Schedule: got non-nil, want nil for STREAMING job")
	}
	if stmt.DoStmt == nil {
		t.Errorf("DoStmt: got nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// ALTER JOB
// ---------------------------------------------------------------------------

func TestAlterJobProperties(t *testing.T) {
	stmt := parseAlterJobStmt(t, `ALTER JOB my_job PROPERTIES("enable" = "false")`)
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if len(stmt.Properties) == 0 {
		t.Errorf("Properties: got empty, want non-empty")
	}
}

func TestAlterJobDoStmt(t *testing.T) {
	stmt := parseAlterJobStmt(t, "ALTER JOB my_job DO INSERT INTO t1 SELECT * FROM t2")
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.NewStmt == nil {
		t.Errorf("NewStmt: got nil, want non-nil")
	}
}

// ---------------------------------------------------------------------------
// DROP JOB
// ---------------------------------------------------------------------------

func TestDropJobByName(t *testing.T) {
	stmt := parseDropJobStmt(t, "DROP JOB my_job")
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.Where != nil {
		t.Errorf("Where: got non-nil, want nil")
	}
}

func TestDropJobIfExists(t *testing.T) {
	stmt := parseDropJobStmt(t, "DROP JOB IF EXISTS my_job")
	if !stmt.IfExists {
		t.Errorf("IfExists: got false, want true")
	}
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
}

func TestDropJobWhere(t *testing.T) {
	// From legacy corpus: DROP JOB WHERE jobName='example'
	stmt := parseDropJobStmt(t, "DROP JOB WHERE jobName='example'")
	if stmt.Where == nil {
		t.Errorf("Where: got nil, want non-nil")
	}
	if stmt.Name != nil {
		t.Errorf("Name: got non-nil, want nil for WHERE form")
	}
}

// ---------------------------------------------------------------------------
// PAUSE JOB
// ---------------------------------------------------------------------------

func TestPauseJobByName(t *testing.T) {
	stmt := parsePauseJobStmt(t, "PAUSE JOB my_job")
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.Where != nil {
		t.Errorf("Where: got non-nil, want nil")
	}
}

func TestPauseJobWhere(t *testing.T) {
	// From legacy corpus: PAUSE JOB WHERE jobname='example'
	stmt := parsePauseJobStmt(t, "PAUSE JOB WHERE jobname='example'")
	if stmt.Where == nil {
		t.Errorf("Where: got nil, want non-nil")
	}
	if stmt.Name != nil {
		t.Errorf("Name: got non-nil, want nil for WHERE form")
	}
}

// ---------------------------------------------------------------------------
// RESUME JOB
// ---------------------------------------------------------------------------

func TestResumeJobByName(t *testing.T) {
	stmt := parseResumeJobStmt(t, "RESUME JOB my_job")
	if stmt.Name == nil || stmt.Name.Parts[0] != "my_job" {
		t.Errorf("Name: got %v, want my_job", stmt.Name.Parts)
	}
	if stmt.Where != nil {
		t.Errorf("Where: got non-nil, want nil")
	}
}

func TestResumeJobWhere(t *testing.T) {
	// From legacy corpus: RESUME JOB WHERE jobName='example'
	stmt := parseResumeJobStmt(t, "RESUME JOB WHERE jobName='example'")
	if stmt.Where == nil {
		t.Errorf("Where: got nil, want non-nil")
	}
	if stmt.Name != nil {
		t.Errorf("Name: got non-nil, want nil for WHERE form")
	}
}

// ---------------------------------------------------------------------------
// CANCEL TASK
// ---------------------------------------------------------------------------

func TestCancelTask(t *testing.T) {
	stmt := parseCancelTaskStmt(t, "CANCEL TASK FOR my_job 12345")
	if stmt.For == nil || stmt.For.Parts[0] != "my_job" {
		t.Errorf("For: got %v, want my_job", stmt.For.Parts)
	}
	if stmt.TaskID != 12345 {
		t.Errorf("TaskID: got %d, want 12345", stmt.TaskID)
	}
}

// ---------------------------------------------------------------------------
// SHOW JOB / SHOW JOB TASK
// ---------------------------------------------------------------------------

func TestShowJob(t *testing.T) {
	stmt := parseShowJobStmt(t, "SHOW JOB")
	if stmt.Like != "" {
		t.Errorf("Like: got %q, want empty", stmt.Like)
	}
	if stmt.Where != nil {
		t.Errorf("Where: got non-nil, want nil")
	}
}

func TestShowJobLike(t *testing.T) {
	stmt := parseShowJobStmt(t, "SHOW JOB LIKE 'my_%'")
	if stmt.Like != "my_%" {
		t.Errorf("Like: got %q, want %q", stmt.Like, "my_%")
	}
}

func TestShowJobWhere(t *testing.T) {
	stmt := parseShowJobStmt(t, "SHOW JOB WHERE name='my_job'")
	if stmt.Where == nil {
		t.Errorf("Where: got nil, want non-nil")
	}
}

func TestShowJobTaskFor(t *testing.T) {
	stmt := parseShowJobTaskStmt(t, "SHOW JOB TASK FOR my_job")
	if stmt.For == nil || stmt.For.Parts[0] != "my_job" {
		t.Errorf("For: got %v, want my_job", stmt.For.Parts)
	}
}

// ---------------------------------------------------------------------------
// Legacy corpus round-trip
// ---------------------------------------------------------------------------

func TestLegacyJobCorpus(t *testing.T) {
	// All statements from doris/parser/testdata/legacy/job.sql
	cases := []string{
		"CREATE JOB my_job ON SCHEDULE EVERY 1 MINUTE DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2",
		"CREATE JOB my_job ON SCHEDULE AT '2020-01-01 00:00:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2",
		"CREATE JOB my_job ON SCHEDULE EVERY 1 DAY STARTS '2020-01-01 00:00:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2 WHERE create_time >= days_add(now(),-1)",
		"CREATE JOB my_job ON SCHEDULE EVERY 1 DAY STARTS '2020-01-01 00:00:00' ENDS '2020-01-01 00:10:00' DO INSERT INTO db1.tbl1 SELECT * FROM db2.tbl2 WHERE create_time >= days_add(now(),-1)",
		"PAUSE JOB WHERE jobname='example'",
		"RESUME JOB WHERE jobName='example'",
		"DROP JOB WHERE jobName='example'",
	}

	for _, sql := range cases {
		t.Run(sql[:min(len(sql), 60)], func(t *testing.T) {
			file, errs := parseForTest(sql)
			if len(errs) != 0 {
				t.Fatalf("parseForTest(%q) errors: %v", sql, errs)
			}
			if len(file.Stmts) != 1 {
				t.Fatalf("parseForTest(%q): got %d stmts, want 1", sql, len(file.Stmts))
			}
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// TestJobDoBodyIsDML: a job's DO takes the grammar's supportedDmlStatement --
// INSERT, UPDATE or DELETE, optionally under EXPLAIN, or LOAD / EXPORT -- and
// nothing else. The body used to be any statement, so DO SELECT 1 INTO OUTFILE
// 'p' parsed (container-verified engine rejects). The error is at the body's
// first token, as the engine reports it, or at the query an EXPLAIN prefix
// took.
func TestJobDoBodyIsDML(t *testing.T) {
	const create = "CREATE JOB j ON SCHEDULE AT '2030-01-01 00:00:00' DO "
	const alter = "ALTER JOB j DO "
	accepted := []struct {
		body string
		want ast.NodeTag
	}{
		{"INSERT INTO t SELECT 1", ast.T_InsertStmt},
		{"TRUNCATE TABLE t", ast.T_TruncateTableStmt},
		{"MERGE INTO t USING s ON t.a = s.a WHEN MATCHED THEN DELETE", ast.T_MergeStmt},
		{"INSERT INTO t VALUES (1)", ast.T_InsertStmt},
		{"INSERT OVERWRITE TABLE t SELECT 1", ast.T_InsertStmt},
		{"INSERT INTO t SELECT 1 UNION SELECT 2", ast.T_InsertStmt},
		{"UPDATE t SET a = 1 WHERE a = 2", ast.T_UpdateStmt},
		{"DELETE FROM t WHERE a = 1", ast.T_DeleteStmt},
		{"EXPLAIN INSERT INTO t SELECT 1", ast.T_ExplainStmt},
		{"EXPLAIN VERBOSE INSERT INTO t SELECT 1", ast.T_ExplainStmt},
		{"EXPLAIN UPDATE t SET a = 1", ast.T_ExplainStmt},
		{`LOAD LABEL l (DATA INFILE("hdfs://h/p") INTO TABLE t)`, ast.T_LoadDataStmt},
		{`EXPORT TABLE t TO "hdfs://h/p"`, ast.T_ExportStmt},
	}
	for _, tc := range accepted {
		if got := parseCreateJobStmt(t, create+tc.body).DoStmt.Tag(); got != tc.want {
			t.Errorf("CREATE JOB ... DO %s: body %v, want %v", tc.body, got, tc.want)
		}
		if got := parseAlterJobStmt(t, alter+tc.body).NewStmt.Tag(); got != tc.want {
			t.Errorf("ALTER JOB ... DO %s: body %v, want %v", tc.body, got, tc.want)
		}
	}

	rejected := []struct{ body, at string }{
		{"SELECT 1", "SELECT 1"},
		{"SELECT 1 INTO OUTFILE 'p'", "SELECT 1"},
		// Refused at its first token, before the query is read at all.
		{"SELECT 1 FROM", "SELECT 1"},
		{"SELECT 1 UNION SELECT 2", "SELECT 1"},
		{"(SELECT 1)", "(SELECT 1)"},
		{"WITH c AS (SELECT 1) SELECT * FROM c", "WITH"},
		{"EXPLAIN SELECT 1", "SELECT 1"},
		{"CREATE TABLE x (a INT)", "CREATE TABLE"},
		{"SHOW TABLES", "SHOW"},
		{"SET x = 1", "SET"},
		{"CREATE JOB j2 ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT INTO t SELECT 1", "CREATE JOB j2"},
	}
	for _, tc := range rejected {
		for _, prefix := range []string{create, alter} {
			sql := prefix + tc.body
			_, errs := parseForTest(sql)
			if len(errs) == 0 {
				t.Errorf("%q parsed; the engine rejects it", sql)
				continue
			}
			if want := len(prefix) + strings.Index(tc.body, tc.at); errs[0].Position != want {
				t.Errorf("%q: error at %d (%s), want %d", sql, errs[0].Position, errs[0].Message, want)
			}
		}
	}
}
