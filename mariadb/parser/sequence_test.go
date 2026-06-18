package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestSequenceAccept covers the CREATE/ALTER/DROP SEQUENCE accept surface
// (BYT-9135), sourced from mdbcheck/corpus/core-sequences-returning.sql and the
// container option-matrix recon (all verified parse-accepted by mariadb:11.8.8).
func TestSequenceAccept(t *testing.T) {
	accept := []string{
		// CREATE — flags + name forms
		"CREATE SEQUENCE s",
		"CREATE SEQUENCE IF NOT EXISTS s",
		"CREATE OR REPLACE SEQUENCE s",
		"CREATE SEQUENCE db.s",
		"CREATE TEMPORARY SEQUENCE s",
		// CREATE — START: WITH / = / bare
		"CREATE OR REPLACE SEQUENCE s START WITH 1000",
		"CREATE OR REPLACE SEQUENCE s START = 5",
		"CREATE OR REPLACE SEQUENCE s START 5",
		// CREATE — INCREMENT: BY / = / bare / negative
		"CREATE OR REPLACE SEQUENCE s INCREMENT BY 5",
		"CREATE OR REPLACE SEQUENCE s INCREMENT = 5",
		"CREATE OR REPLACE SEQUENCE s INCREMENT 5",
		"CREATE OR REPLACE SEQUENCE s INCREMENT BY -1 MAXVALUE -1 MINVALUE -100",
		"CREATE OR REPLACE SEQUENCE s INCREMENT BY +5",
		"CREATE OR REPLACE SEQUENCE s START WITH +5",
		// CREATE — MIN/MAX: value / = / NO-forms (one-word + spaced)
		"CREATE OR REPLACE SEQUENCE s MINVALUE 10 MAXVALUE 9999",
		"CREATE OR REPLACE SEQUENCE s MINVALUE = 1 MAXVALUE = 99",
		"CREATE OR REPLACE SEQUENCE s NOMINVALUE NOMAXVALUE",
		"CREATE OR REPLACE SEQUENCE s NO MINVALUE NO MAXVALUE",
		// CREATE — CACHE / CYCLE
		"CREATE OR REPLACE SEQUENCE s CACHE 50",
		"CREATE OR REPLACE SEQUENCE s CACHE = 50",
		"CREATE OR REPLACE SEQUENCE s NOCACHE",
		"CREATE OR REPLACE SEQUENCE s MINVALUE 1 MAXVALUE 5 CYCLE",
		"CREATE OR REPLACE SEQUENCE s MINVALUE 1 MAXVALUE 5 NOCYCLE",
		// CREATE — all options combined (corpus)
		"CREATE OR REPLACE SEQUENCE s START WITH 100 MINVALUE 10 MAXVALUE 100000 INCREMENT BY 5 CACHE 20 CYCLE",
		"CREATE OR REPLACE SEQUENCE s START WITH 1 NOMINVALUE NOMAXVALUE NOCACHE NOCYCLE INCREMENT BY 1",
		// CREATE — AS int_type (any-order with options)
		"CREATE OR REPLACE SEQUENCE s AS BIGINT",
		"CREATE OR REPLACE SEQUENCE s AS INT UNSIGNED",
		"CREATE OR REPLACE SEQUENCE s AS TINYINT",
		"CREATE OR REPLACE SEQUENCE s START WITH 1 AS BIGINT",
		"CREATE OR REPLACE SEQUENCE s AS BIGINT START WITH 1",
		"CREATE OR REPLACE SEQUENCE s INCREMENT BY 1 AS INT MINVALUE 1",
		// ALTER — RESTART forms + value options + IF EXISTS
		"ALTER SEQUENCE s RESTART",
		"ALTER SEQUENCE s RESTART WITH 500",
		"ALTER SEQUENCE s RESTART = 5",
		"ALTER SEQUENCE s RESTART 5",
		"ALTER SEQUENCE s INCREMENT BY 10",
		"ALTER SEQUENCE s MAXVALUE 1000000",
		"ALTER SEQUENCE s MINVALUE 1 NOCACHE",
		"ALTER SEQUENCE IF EXISTS s RESTART WITH 1",
		// ALTER — MariaDB also accepts START and AS (container-verified)
		"ALTER SEQUENCE s START WITH 5",
		"ALTER SEQUENCE s AS BIGINT",
		"ALTER SEQUENCE s RESTART INCREMENT BY 5",
		"ALTER SEQUENCE s INCREMENT BY 5 RESTART",
		// DROP — single / IF EXISTS / multi / TEMPORARY
		"DROP SEQUENCE s",
		"DROP SEQUENCE IF EXISTS s",
		"DROP SEQUENCE IF EXISTS a, b",
		"DROP SEQUENCE a, b, c",
		"DROP TEMPORARY SEQUENCE s",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSequenceReject covers the container-verified 1064 rejects: the donor
// SQL-Server/Oracle supersets MariaDB does not accept, and the value/grammar
// restrictions. Guards against over-acceptance.
func TestSequenceReject(t *testing.T) {
	reject := []string{
		// RESTART is ALTER-only
		"CREATE SEQUENCE s RESTART",
		"CREATE SEQUENCE s RESTART WITH 5",
		// option values are mandatory (bare keyword rejects)
		"CREATE SEQUENCE s CACHE",
		"CREATE SEQUENCE s INCREMENT",
		"CREATE SEQUENCE s START",
		"CREATE SEQUENCE s MINVALUE",
		"CREATE SEQUENCE s MAXVALUE",
		"CREATE SEQUENCE s INCREMENT BY",
		// spaced NO CACHE / NO CYCLE (one-word NOCACHE/NOCYCLE only)
		"CREATE SEQUENCE s NO CACHE",
		"CREATE SEQUENCE s NO CYCLE",
		// AS restricted to integer types without display width
		"CREATE SEQUENCE s AS VARCHAR(10)",
		"CREATE SEQUENCE s AS DECIMAL(10)",
		"CREATE SEQUENCE s AS INT(11)",
		"CREATE SEQUENCE s AS DATE",
		// option values are integer literals, not expressions
		"CREATE SEQUENCE s START WITH 1+1",
		"CREATE SEQUENCE s START WITH (5)",
		"CREATE SEQUENCE s CACHE CYCLE",
		"CREATE SEQUENCE s MINVALUE MAXVALUE 9",
		// ALTER requires at least one option (CREATE does not)
		"ALTER SEQUENCE s",
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestSequenceAST verifies field population for representative statements
// (robust against byte-offset churn, unlike a full NodeToString match).
func TestSequenceAST(t *testing.T) {
	t.Run("create or replace", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE OR REPLACE SEQUENCE s")
		if !stmt.OrReplace || stmt.IfNotExists {
			t.Errorf("OrReplace=%v IfNotExists=%v, want true/false", stmt.OrReplace, stmt.IfNotExists)
		}
		if stmt.Name == nil || stmt.Name.Name != "s" {
			t.Errorf("Name = %+v, want s", stmt.Name)
		}
	})
	t.Run("if not exists", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE IF NOT EXISTS s")
		if !stmt.IfNotExists {
			t.Errorf("IfNotExists = false, want true")
		}
	})
	t.Run("nocache", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE s NOCACHE")
		if !stmt.NoCache {
			t.Errorf("NoCache = false, want true")
		}
	})
	t.Run("cycle true", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE s CYCLE")
		if stmt.Cycle == nil || !*stmt.Cycle {
			t.Errorf("Cycle = %v, want &true", stmt.Cycle)
		}
	})
	t.Run("nocycle false", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE s NOCYCLE")
		if stmt.Cycle == nil || *stmt.Cycle {
			t.Errorf("Cycle = %v, want &false", stmt.Cycle)
		}
	})
	t.Run("as bigint", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE s AS BIGINT")
		if stmt.DataType == nil || stmt.DataType.Name != "BIGINT" {
			t.Errorf("DataType = %+v, want BIGINT", stmt.DataType)
		}
	})
	t.Run("start/increment populated", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE SEQUENCE s START WITH 1 INCREMENT BY 5")
		if stmt.Start == nil || stmt.Increment == nil {
			t.Errorf("Start=%v Increment=%v, want both set", stmt.Start, stmt.Increment)
		}
	})
	t.Run("alter if exists restart with", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterSequenceStmt](t, "ALTER SEQUENCE IF EXISTS s RESTART WITH 500")
		if !stmt.IfExists || !stmt.Restart || stmt.RestartWith == nil {
			t.Errorf("IfExists=%v Restart=%v RestartWith=%v", stmt.IfExists, stmt.Restart, stmt.RestartWith)
		}
	})
	t.Run("alter start", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterSequenceStmt](t, "ALTER SEQUENCE s START WITH 5")
		if stmt.Start == nil {
			t.Errorf("Start = nil, want set")
		}
	})
	t.Run("drop multi if exists", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.DropSequenceStmt](t, "DROP SEQUENCE IF EXISTS a, b")
		if !stmt.IfExists || len(stmt.Sequences) != 2 {
			t.Errorf("IfExists=%v len(Sequences)=%d, want true/2", stmt.IfExists, len(stmt.Sequences))
		}
	})
	t.Run("create temporary", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateSequenceStmt](t, "CREATE TEMPORARY SEQUENCE s")
		if !stmt.Temporary {
			t.Errorf("Temporary = false, want true")
		}
	})
	t.Run("drop temporary", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.DropSequenceStmt](t, "DROP TEMPORARY SEQUENCE s")
		if !stmt.Temporary {
			t.Errorf("Temporary = false, want true")
		}
	})
}

// parseSeqStmt parses sql and returns its single statement as type T.
func parseSeqStmt[T ast.Node](t *testing.T, sql string) T {
	t.Helper()
	result := ParseAndCheck(t, sql)
	if result.Len() != 1 {
		t.Fatalf("Parse(%q) returned %d statements, want 1", sql, result.Len())
	}
	stmt, ok := result.Items[0].(T)
	if !ok {
		t.Fatalf("Parse(%q) = %T, want %T", sql, result.Items[0], *new(T))
	}
	return stmt
}
