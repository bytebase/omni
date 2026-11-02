package parser

import (
	"strings"
	"testing"
)

// TestDorisJobStatementsRejected pins the job statements the Doris fork brought
// along as rejected, and where. StarRocks 3.4.10 has none of them and reports
// each as a syntax error at the word after the verb (positions below are the
// engine's columns). TestStarRocksNoJobStatements holds the same forms to the
// engine; this test also holds the position, and that the error is a syntax
// error, not the "not yet supported" stub that callers read as a parser gap.
func TestDorisJobStatementsRejected(t *testing.T) {
	cases := []struct {
		sql string
		pos int
	}{
		{"CREATE JOB j ON SCHEDULE EVERY 1 MINUTE DO INSERT INTO t SELECT 1", 7},
		{"CREATE JOB j ON SCHEDULE AT '2030-01-01 00:00:00' DO INSERT INTO t SELECT 1", 7},
		{"CREATE JOB IF NOT EXISTS j ON SCHEDULE EVERY 1 HOUR DO INSERT INTO t SELECT 1", 7},
		{"CREATE JOB j AS STREAMING DO INSERT INTO t SELECT 1", 7},
		{`ALTER JOB j PROPERTIES("enable" = "false")`, 6},
		{"ALTER JOB j DO INSERT INTO t SELECT 1", 6},
		{"ALTER JOB j", 6},
		{"DROP JOB j", 5},
		{"DROP JOB IF EXISTS j", 5},
		{"DROP JOB WHERE jobName = 'j'", 5},
		{"PAUSE JOB j", 6},
		{"PAUSE JOB WHERE jobName = 'j'", 6},
		{"RESUME JOB j", 7},
		{"RESUME JOB WHERE jobName = 'j'", 7},
		{"CANCEL TASK FOR j 1", 7},
		{"CANCEL TASK WHERE jobName = 'j' AND taskId = 1", 7},
		{"SHOW JOB", 5},
		{"SHOW JOB LIKE 'j%'", 5},
		{"SHOW JOB WHERE jobName = 'j'", 5},
		{"SHOW JOB TASK FOR j", 5},
		{"SHOW JOBS", 5},
		{"SHOW JOBS WHERE jobName = 'j'", 5},
		{"SHOW FULL JOB", 10},
		{"SHOW ALL JOB", 9},
		{"SHOW TEMPORARY JOB", 15},
		{"PAUSE MATERIALIZED VIEW JOB ON mv", 6},
		{"RESUME MATERIALIZED VIEW JOB ON mv", 7},
		{"CANCEL MATERIALIZED VIEW TASK 1 ON mv", 7},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := parseForTest(tc.sql)
			if len(errs) != 1 {
				t.Fatalf("got %d errors, want 1: %v", len(errs), errs)
			}
			if len(file.Stmts) != 0 {
				t.Errorf("got %d statements, want none", len(file.Stmts))
			}
			if errs[0].Position != tc.pos {
				t.Errorf("error at %d, want %d (%q)", errs[0].Position, tc.pos, errs[0].Message)
			}
			if !strings.HasPrefix(errs[0].Message, "syntax error") {
				t.Errorf("message %q, want a syntax error", errs[0].Message)
			}
		})
	}
}

// TestJobWordsStayIdentifiers pins what the rejection above must not touch: the
// job and task words keep their keyword class, and the StarRocks statements that
// name a job or task keep parsing. StarRocks 3.4.10 parses each of these.
func TestJobWordsStayIdentifiers(t *testing.T) {
	for _, sql := range []string{
		"SELECT job FROM t",
		"SELECT jobs, task, tasks FROM t",
		"SELECT schedule, streaming, every, starts, ends FROM t",
		"SELECT 1 AS job",
		"SELECT job.a FROM job",
		"CREATE TABLE job (a INT)",
		"CREATE TABLE zzt (job INT, task INT, schedule INT)",
		"CREATE DATABASE job",
		"DROP TABLE IF EXISTS job",
		"ALTER TABLE job ADD COLUMN b INT",
		"SHOW ANALYZE JOB",
		"SHOW ROUTINE LOAD TASK WHERE JobName = 'x'",
		"PAUSE ROUTINE LOAD FOR job",
		"RESUME ROUTINE LOAD FOR job",
		"CANCEL REFRESH MATERIALIZED VIEW mv",
		"CANCEL LOAD WHERE LABEL = 'x'",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, errs := parseForTest(sql); len(errs) != 0 {
				t.Errorf("errors: %v", errs)
			}
		})
	}
}
