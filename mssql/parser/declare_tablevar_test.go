package parser

import "testing"

// TestDeclareTableVariableConstraints pins table-level constraints and
// indexes inside DECLARE @t TABLE bodies — the same element grammar as
// CREATE TABLE. A customer function declared table variables with
// PRIMARY KEY CLUSTERED entries and failed SQL review (BYT-9909 stmt 6).
func TestDeclareTableVariableConstraints(t *testing.T) {
	accepts := []struct {
		name string
		sql  string
	}{
		{"table-level primary key", "DECLARE @t TABLE (p_year INT, PRIMARY KEY (p_year));"},
		{"primary key clustered", "DECLARE @t TABLE (p_year INT, PRIMARY KEY CLUSTERED (p_year));"},
		{"primary key nonclustered", "DECLARE @t TABLE (role_code VARCHAR(10), PRIMARY KEY NONCLUSTERED (role_code));"},

		{"unique constraint", "DECLARE @t TABLE (a INT, UNIQUE (a));"},
		{"check constraint", "DECLARE @t TABLE (a INT, CHECK (a > 0));"},
		{"inline index", "DECLARE @t TABLE (a INT, INDEX ix_a (a));"},
		{"inline column pk control", "DECLARE @t TABLE (value NVARCHAR(100) PRIMARY KEY);"},
		{"customer shape", "DECLARE @nhom_ns TABLE (value NVARCHAR(100) PRIMARY KEY);\nDECLARE @y TABLE (p_year INT, PRIMARY KEY CLUSTERED (p_year));"},
	}
	t.Run("named constraint rejected like the engine", func(t *testing.T) {
		// Engine-verified (SQL Server 2019, Msg 156): the CONSTRAINT name
		// prefix is invalid inside a table variable body.
		if _, err := Parse("DECLARE @t TABLE (a INT, CONSTRAINT pk_a PRIMARY KEY (a));"); err == nil {
			t.Fatal("expected parse error for named constraint in table variable")
		}
	})

	// Documented over-accept (parser-unsure -> accept direction): a
	// trailing comma before ')' is accepted here while SQL Server rejects
	// it. Recorded by cross-validation; revisit only with engine-anchored
	// tightening.

	for _, tc := range accepts {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Parse(tc.sql); err != nil {
				t.Fatalf("parse failed: %v", err)
			}
		})
	}
}
