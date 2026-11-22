package analysis

import (
	"errors"
	"testing"

	"github.com/bytebase/omni/analysis"
	"github.com/bytebase/omni/clickhouse/parser"
)

func TestValidateQueryExplainAndViews(t *testing.T) {
	cases := []struct {
		name      string
		statement string
		readOnly  bool
		readsRows bool
	}{
		{name: "explain ast select", statement: "EXPLAIN AST SELECT 1", readOnly: true, readsRows: true},
		{name: "explain ast show", statement: "EXPLAIN AST SHOW TABLES", readOnly: true, readsRows: true},
		{name: "explain table override", statement: "EXPLAIN TABLE OVERRIDE db.t", readOnly: true, readsRows: true},
		{name: "explain current transaction", statement: "EXPLAIN CURRENT TRANSACTION", readOnly: true, readsRows: true},
		{name: "explain ast drop does not execute", statement: "EXPLAIN AST DROP TABLE t", readOnly: true, readsRows: true},
		{name: "explain analyze rejected as executing", statement: "EXPLAIN ANALYZE SELECT 1", readOnly: false, readsRows: true},
		{name: "explain into outfile rejected", statement: "EXPLAIN SELECT 1 INTO OUTFILE '/tmp/plan.txt'", readOnly: false, readsRows: false},
		{name: "describe into outfile rejected", statement: "DESCRIBE TABLE t INTO OUTFILE '/tmp/desc.txt'", readOnly: false, readsRows: false},
		{name: "create materialized view engine order by", statement: "CREATE MATERIALIZED VIEW mv_own ENGINE = MergeTree ORDER BY id AS SELECT id FROM src;", readOnly: false, readsRows: false},
		{name: "create view with comment", statement: "CREATE VIEW v AS SELECT id FROM src COMMENT 'view note';", readOnly: false, readsRows: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			ro, rd, err := ValidateQuery(tc.statement)
			if err != nil {
				t.Fatalf("ValidateQuery(%q) unexpected error: %v", tc.statement, err)
			}
			if ro != tc.readOnly || rd != tc.readOnly {
				t.Errorf("ValidateQuery(%q) = (%v, %v), want (%v, %v)", tc.statement, ro, rd, tc.readOnly, tc.readOnly)
			}
			rows, err := StatementReadsRows(tc.statement)
			if err != nil {
				t.Fatalf("StatementReadsRows(%q) unexpected error: %v", tc.statement, err)
			}
			if rows != tc.readsRows {
				t.Errorf("StatementReadsRows(%q) = %v, want %v", tc.statement, rows, tc.readsRows)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		statement string
		want      QueryType
		wantErr   bool
	}{
		{"", analysis.QueryTypeUnknown, false},
		{"-- comment only\n/* block */", analysis.QueryTypeUnknown, false},
		{"SELECT 1", analysis.Select, false},
		{"WITH x AS (SELECT 1) SELECT * FROM x", analysis.Select, false},
		{"FROM events WHERE user_id = 42", analysis.Select, false},
		{"SELECT 1 INTO OUTFILE '/tmp/x'", analysis.DML, false},
		{"SHOW TABLES", analysis.SelectInfoSchema, false},
		{"DESCRIBE TABLE t", analysis.SelectInfoSchema, false},
		{"EXISTS TABLE t", analysis.SelectInfoSchema, false},
		{"CHECK TABLE t", analysis.SelectInfoSchema, false},
		{"SHOW TABLES INTO OUTFILE '/tmp/x'", analysis.DML, false},
		{"EXPLAIN SELECT 1", analysis.Explain, false},
		{"EXPLAIN ANALYZE SELECT 1", analysis.Explain, false},
		{"EXPLAIN CURRENT TRANSACTION", analysis.Explain, false},
		{"INSERT INTO t VALUES (1)", analysis.DML, false},
		{"UPDATE t SET a = 1 WHERE id = 1", analysis.DML, false},
		{"DELETE FROM t WHERE id = 1", analysis.DML, false},
		{"CREATE TABLE t (a Int32) ENGINE = Memory", analysis.DDL, false},
		{"CREATE MATERIALIZED VIEW mv ENGINE = MergeTree ORDER BY id AS SELECT id FROM src", analysis.DDL, false},
		{"CREATE VIEW v AS SELECT 1 COMMENT 'note'", analysis.DDL, false},
		{"ALTER TABLE t ADD COLUMN b String", analysis.DDL, false},
		{"DROP TABLE t", analysis.DDL, false},
		{"TRUNCATE TABLE t", analysis.DDL, false},
		{"SET max_threads = 4", analysis.QueryTypeUnknown, false},
		{"USE analytics", analysis.QueryTypeUnknown, false},
		{"GRANT SELECT ON db.* TO analyst", analysis.QueryTypeUnknown, false},
		{"SELECT 1; SELECT 2", analysis.QueryTypeUnknown, true},
		{"SELECT 1 +", analysis.QueryTypeUnknown, true},
	}

	for _, tc := range cases {
		got, err := Classify(tc.statement)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Classify(%q) = (%v, nil), want error", tc.statement, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("Classify(%q) unexpected error: %v", tc.statement, err)
			continue
		}
		if got != tc.want {
			t.Errorf("Classify(%q) = %v, want %v", tc.statement, got, tc.want)
		}
	}

	_, err := Classify("  SELECT 1 +")
	var pe *parser.ParseError
	if !errors.As(err, &pe) || pe.Position != 12 {
		t.Errorf("Classify syntax error = %v, want *parser.ParseError at byte 12", err)
	}
}

func TestPlanFormatAndTrimStatement(t *testing.T) {
	for _, tc := range []struct {
		statement string
		format    string
		ok        bool
	}{
		{"SELECT 1", "", false},
		{"EXPLAIN SELECT 1", "TEXT", true},
		{"EXPLAIN json = 1 SELECT 1", "JSON", true},
		{"EXPLAIN CURRENT TRANSACTION", "TEXT", true},
	} {
		format, ok := PlanFormat(tc.statement)
		if ok != tc.ok || format != tc.format {
			t.Errorf("PlanFormat(%q) = (%q, %v), want (%q, %v)", tc.statement, format, ok, tc.format, tc.ok)
		}
	}

	for _, tc := range []struct{ in, want string }{
		{"  SELECT 1;  ", "SELECT 1"},
		{"\nSELECT 1;;;\n", "SELECT 1"},
		{"  ", ""},
		{";;", ""},
	} {
		if got := TrimStatement(tc.in); got != tc.want {
			t.Errorf("TrimStatement(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
