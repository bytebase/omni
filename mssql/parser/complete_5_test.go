package parser

import (
	"testing"
)

// --- Section 5.1: CREATE TABLE ---

func TestCollect_CreateTable(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "CREATE TABLE |",
			sql:       "CREATE TABLE ",
			wantRules: []string{"identifier"},
		},
		{
			name:      "CREATE TABLE t (a INT, |)",
			sql:       "CREATE TABLE t (a INT, ",
			wantRules: []string{"identifier"},
			wantToks:  []int{kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK, kwINDEX},
		},
		{
			name:     "CREATE TABLE t (a INT |)",
			sql:      "CREATE TABLE t (a INT ",
			wantToks: []int{kwNULL, kwNOT, kwDEFAULT, kwIDENTITY, kwPRIMARY, kwUNIQUE, kwCHECK, kwREFERENCES, kwCONSTRAINT, kwCOLLATE},
		},
		{
			name:      "CREATE TABLE t (a |)",
			sql:       "CREATE TABLE t (a ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "CREATE TABLE t (a INT IDENTITY(|)",
			sql:       "CREATE TABLE t (a INT IDENTITY(",
			wantRules: []string{"numeric"},
		},
		{
			name:      "CREATE TABLE t (a INT REFERENCES |)",
			sql:       "CREATE TABLE t (a INT REFERENCES ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE TABLE t (a INT DEFAULT |)",
			sql:       "CREATE TABLE t (a INT DEFAULT ",
			wantRules: []string{"expression"},
		},
		{
			name:     "CREATE TABLE t (a INT) |",
			sql:      "CREATE TABLE t (a INT) ",
			wantToks: []int{kwON, kwWITH, ';'},
		},
		{
			name:      "CREATE TABLE t (CONSTRAINT pk PRIMARY KEY (|)",
			sql:       "CREATE TABLE t (CONSTRAINT pk PRIMARY KEY (",
			wantRules: []string{"columnref"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}

// --- Section 5.2: ALTER TABLE ---

func TestCollect_AlterTable(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "ALTER TABLE |",
			sql:       "ALTER TABLE ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "ALTER TABLE t |",
			sql:       "ALTER TABLE t ",
			wantToks:  []int{kwADD, kwDROP, kwALTER, kwSET, kwWITH},
			wantRules: []string{"ENABLE", "DISABLE", "SWITCH", "REBUILD"},
		},
		{
			name:      "ALTER TABLE t ADD |",
			sql:       "ALTER TABLE t ADD ",
			wantToks:  []int{kwCOLUMN, kwCONSTRAINT, kwPRIMARY, kwUNIQUE, kwFOREIGN, kwCHECK, kwDEFAULT, kwINDEX},
			wantRules: []string{"identifier"},
		},
		{
			name:      "ALTER TABLE t ADD COLUMN |",
			sql:       "ALTER TABLE t ADD COLUMN ",
			wantRules: []string{"identifier"},
		},
		{
			name:     "ALTER TABLE t DROP |",
			sql:      "ALTER TABLE t DROP ",
			wantToks: []int{kwCOLUMN, kwCONSTRAINT, kwINDEX},
		},
		{
			name:      "ALTER TABLE t DROP COLUMN |",
			sql:       "ALTER TABLE t DROP COLUMN ",
			wantRules: []string{"columnref"},
		},
		{
			name:      "ALTER TABLE t DROP CONSTRAINT |",
			sql:       "ALTER TABLE t DROP CONSTRAINT ",
			wantRules: []string{"constraint_name"},
		},
		{
			name:      "ALTER TABLE t ALTER COLUMN |",
			sql:       "ALTER TABLE t ALTER COLUMN ",
			wantRules: []string{"columnref"},
		},
		{
			name:      "ALTER TABLE t ALTER COLUMN a |",
			sql:       "ALTER TABLE t ALTER COLUMN a ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "ALTER TABLE t ENABLE TRIGGER |",
			sql:       "ALTER TABLE t ENABLE TRIGGER ",
			wantRules: []string{"trigger_ref"},
			wantToks:  []int{kwALL},
		},
		{
			name:      "ALTER TABLE t DISABLE TRIGGER |",
			sql:       "ALTER TABLE t DISABLE TRIGGER ",
			wantRules: []string{"trigger_ref"},
			wantToks:  []int{kwALL},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}

// --- Section 5.3: CREATE/DROP Index, View, Database ---

func TestCollect_CreateDropIndexViewDatabase(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "CREATE INDEX idx ON |",
			sql:       "CREATE INDEX idx ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE INDEX idx ON t (|)",
			sql:       "CREATE INDEX idx ON t (",
			wantRules: []string{"columnref"},
		},
		{
			name:      "CREATE UNIQUE INDEX idx ON |",
			sql:       "CREATE UNIQUE INDEX idx ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE CLUSTERED INDEX idx ON |",
			sql:       "CREATE CLUSTERED INDEX idx ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE NONCLUSTERED INDEX idx ON |",
			sql:       "CREATE NONCLUSTERED INDEX idx ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "DROP INDEX |",
			sql:       "DROP INDEX ",
			wantRules: []string{"index_name"},
		},
		{
			name:      "CREATE VIEW |",
			sql:       "CREATE VIEW ",
			wantRules: []string{"identifier"},
		},
		{
			name:     "CREATE VIEW v AS |",
			sql:      "CREATE VIEW v AS ",
			wantToks: []int{kwSELECT},
		},
		{
			name:     "ALTER VIEW v AS |",
			sql:      "ALTER VIEW v AS ",
			wantToks: []int{kwSELECT},
		},
		{
			name:      "DROP VIEW |",
			sql:       "DROP VIEW ",
			wantRules: []string{"view_name"},
		},
		{
			name:      "CREATE DATABASE |",
			sql:       "CREATE DATABASE ",
			wantRules: []string{"identifier"},
		},
		{
			name:      "DROP DATABASE |",
			sql:       "DROP DATABASE ",
			wantRules: []string{"database_ref"},
		},
		{
			name:      "DROP TABLE |",
			sql:       "DROP TABLE ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "DROP TABLE IF EXISTS |",
			sql:       "DROP TABLE IF EXISTS ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "TRUNCATE TABLE |",
			sql:       "TRUNCATE TABLE ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE SCHEMA |",
			sql:       "CREATE SCHEMA ",
			wantRules: []string{"identifier"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}

// --- Section 5.4: CREATE/ALTER PROCEDURE & FUNCTION ---

func TestCollect_CreateAlterProcFunction(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "CREATE PROCEDURE |",
			sql:       "CREATE PROCEDURE ",
			wantRules: []string{"identifier"},
		},
		{
			name:      "CREATE PROC |",
			sql:       "CREATE PROC ",
			wantRules: []string{"identifier"},
		},
		{
			name:     "CREATE PROCEDURE p |",
			sql:      "CREATE PROCEDURE p ",
			wantToks: []int{kwAS, kwWITH},
		},
		{
			name:     "CREATE PROCEDURE p AS |",
			sql:      "CREATE PROCEDURE p AS ",
			wantToks: []int{kwSELECT, kwINSERT, kwUPDATE, kwDELETE, kwBEGIN},
		},
		{
			name:      "ALTER PROCEDURE |",
			sql:       "ALTER PROCEDURE ",
			wantRules: []string{"proc_name"},
		},
		{
			name:      "DROP PROCEDURE |",
			sql:       "DROP PROCEDURE ",
			wantRules: []string{"proc_name"},
		},
		{
			name:      "CREATE FUNCTION |",
			sql:       "CREATE FUNCTION ",
			wantRules: []string{"identifier"},
		},
		{
			name:      "CREATE FUNCTION f (|)",
			sql:       "CREATE FUNCTION f (",
			wantRules: []string{"variable"},
		},
		{
			name:      "CREATE FUNCTION f () RETURNS |",
			sql:       "CREATE FUNCTION f () RETURNS ",
			wantRules: []string{"type_name"},
		},
		{
			name:      "ALTER FUNCTION |",
			sql:       "ALTER FUNCTION ",
			wantRules: []string{"func_name"},
		},
		{
			name:      "DROP FUNCTION |",
			sql:       "DROP FUNCTION ",
			wantRules: []string{"func_name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}

// --- Section 5.5: CREATE/DROP TRIGGER ---

func TestCollect_CreateDropTrigger(t *testing.T) {
	tests := []struct {
		name      string
		sql       string
		wantRules []string
		wantToks  []int
	}{
		{
			name:      "CREATE TRIGGER |",
			sql:       "CREATE TRIGGER ",
			wantRules: []string{"identifier"},
		},
		{
			name:      "CREATE TRIGGER tr ON |",
			sql:       "CREATE TRIGGER tr ON ",
			wantRules: []string{"table_ref"},
		},
		{
			name:      "CREATE TRIGGER tr ON t |",
			sql:       "CREATE TRIGGER tr ON t ",
			wantToks:  []int{kwFOR, kwWITH},
			wantRules: []string{"AFTER", "INSTEAD OF"},
		},
		{
			name:     "CREATE TRIGGER tr ON t FOR |",
			sql:      "CREATE TRIGGER tr ON t FOR ",
			wantToks: []int{kwINSERT, kwUPDATE, kwDELETE},
		},
		{
			name:      "DROP TRIGGER |",
			sql:       "DROP TRIGGER ",
			wantRules: []string{"trigger_name"},
		},
		{
			name:      "ALTER TRIGGER |",
			sql:       "ALTER TRIGGER ",
			wantRules: []string{"trigger_name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cs := Collect(tt.sql, len(tt.sql))
			if cs == nil {
				t.Fatal("Collect returned nil")
			}
			for _, r := range tt.wantRules {
				if !cs.HasRule(r) {
					t.Errorf("missing rule candidate %q", r)
				}
			}
			for _, tok := range tt.wantToks {
				if !cs.HasToken(tok) {
					t.Errorf("missing token candidate %s (%d)", TokenName(tok), tok)
				}
			}
		})
	}
}
