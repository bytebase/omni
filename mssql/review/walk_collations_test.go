package review

import (
	"testing"

	"github.com/bytebase/omni/metadata"
)

// TestWalkThroughCollations pins where WALK_THROUGH reports an unknown
// collation (448) and where it leaves one to an error SQL Server reports
// first, or to a statement SQL Server does not bind.
func TestWalkThroughCollations(t *testing.T) {
	const missing = "Collation `nope` doesn't exist."
	// An operand of another type fails before the collation (447).
	const notChar = "SQL Server doesn't accept `COLLATE` on an expression of type `int`."
	tests := []struct {
		name      string
		sql       string
		statement int // -1 for no finding
		message   string
	}{
		// Declared collations, checked as the batch is parsed.
		{"column", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE nope);", 0, missing},
		{"column of a missing schema", "CREATE TABLE nosuch.x (a nvarchar(10) COLLATE nope);", 0, missing},
		{"column of a missing type", "CREATE TABLE dbo.x (a nosuch COLLATE nope);", 0, missing},
		{"listed collation in another case", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE latin1_general_100_ci_as_sc_utf8);", -1, ""},
		{"SQL Server 2025 collation", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE Chinese_Simplified_Pinyin_160_CI_AS_SC_UTF8);", -1, ""},
		{"compatibility collation", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE Compatibility_52_409_30003);", -1, ""},
		{"collation-shaped name no version lists", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE Klingon_100_CI_AS);", -1, ""},
		{"database_default column", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE database_default);", -1, ""},
		{"sensitivities out of order", "CREATE TABLE dbo.x (a nvarchar(10) COLLATE Latin1_General_CI_AS_WS_KS);", 0, "Collation `Latin1_General_CI_AS_WS_KS` doesn't exist."},
		{"alter column of a missing table", "ALTER TABLE dbo.nosuch ALTER COLUMN a nvarchar(10) COLLATE nope;", 0, missing},
		{"add column of a missing table", "ALTER TABLE dbo.nosuch ADD c nvarchar(10) COLLATE nope;", 0, missing},
		{"table variable", "DECLARE @t TABLE (c nvarchar(10) COLLATE nope);", 0, missing},
		{"in a branch that does not run", "IF 1 = 0 ALTER TABLE dbo.ct ADD c nvarchar(10) COLLATE nope;", 0, missing},
		{"database", "CREATE DATABASE zz COLLATE nope;", 0, missing},
		{"database_default database", "CREATE DATABASE zz COLLATE database_default;", 0, "Collation `database_default` doesn't exist."},
		{"alter database", "ALTER DATABASE CURRENT COLLATE nope;", 0, missing},
		{"partition function", "CREATE PARTITION FUNCTION pf (nvarchar(10) COLLATE database_default) AS RANGE FOR VALUES ();", 0, "Collation `database_default` doesn't exist."},
		{"procedure body", "CREATE PROCEDURE dbo.p AS ALTER TABLE dbo.nosuch ADD c nvarchar(10) COLLATE nope;", 0, missing},
		{"later statement of a procedure's body", "CREATE PROCEDURE dbo.p AS SELECT 1;\nDECLARE @t TABLE (c nvarchar(10) COLLATE nope);", 1, missing},
		{"after a valid statement", "CREATE TABLE dbo.x (a int);\nCREATE TABLE dbo.y (a nvarchar(10) COLLATE nope);", 1, missing},
		// Expressions, bound once the names they read are found.
		{"view column", "CREATE VIEW dbo.w AS SELECT n COLLATE nope AS n FROM dbo.ct;", 0, missing},
		{"view order by", "CREATE VIEW dbo.w AS SELECT i FROM dbo.ct ORDER BY n COLLATE nope OFFSET 0 ROWS;", 0, missing},
		{"view column of another type", "CREATE VIEW dbo.w AS SELECT i COLLATE nope AS n FROM dbo.ct;", 0, notChar},
		{"view of a number", "CREATE VIEW dbo.w AS SELECT 1 COLLATE nope AS n;", 0, notChar},
		{"view of a missing table", "CREATE VIEW dbo.w AS SELECT 'x' COLLATE nope AS c FROM dbo.nosuch;", 0, "`dbo.nosuch` doesn't exist."},
		{"view of a valid collation", "CREATE VIEW dbo.w AS SELECT n COLLATE Latin1_General_BIN2 AS n FROM dbo.ct;", -1, ""},
		{"procedure query", "CREATE PROCEDURE dbo.p AS SELECT n COLLATE nope FROM dbo.ct;", 0, missing},
		{"procedure query of a missing table", "CREATE PROCEDURE dbo.p AS SELECT 'x' COLLATE nope FROM dbo.nosuch;", -1, ""},
		{"procedure parameter", "CREATE PROCEDURE dbo.p @x nvarchar(10) AS SELECT @x COLLATE nope;", 0, missing},
		{"procedure parameter of another type", "CREATE PROCEDURE dbo.p @x int AS SELECT @x COLLATE nope;", -1, ""},
		{"scalar function", "CREATE FUNCTION dbo.f() RETURNS nvarchar(10) AS BEGIN RETURN N'x' COLLATE nope END;", 0, missing},
		{"trigger", "CREATE TRIGGER dbo.tr ON dbo.ct AFTER INSERT AS SELECT n COLLATE nope FROM inserted;", 0, missing},
		{"computed column", "CREATE TABLE dbo.x (c nvarchar(10), d AS (c COLLATE nope));", 0, missing},
		{"computed column of another type", "CREATE TABLE dbo.x (c int, d AS (c COLLATE nope));", 0, notChar},
		{"check of a missing schema", "CREATE TABLE nosuch.x (c nvarchar(10) DEFAULT (N'x' COLLATE nope));", 0, missing},
		{"added default", "ALTER TABLE dbo.ct ADD c nvarchar(10) DEFAULT (N'x' COLLATE nope);", 0, missing},
		{"added computed column of a missing table", "ALTER TABLE dbo.nosuch ADD c AS (N'x' COLLATE nope);", 0, "Table `dbo.nosuch` doesn't exist."},
		{"added check", "ALTER TABLE dbo.ct ADD CHECK (n = N'x' COLLATE nope);", 0, missing},
		{"default for a column", "ALTER TABLE dbo.ct ADD CONSTRAINT k DEFAULT (N'x' COLLATE nope) FOR n;", 0, missing},
		{"default for a missing column", "ALTER TABLE dbo.ct ADD CONSTRAINT k DEFAULT (N'x' COLLATE nope) FOR nosuch;", 0, "Column `nosuch` doesn't exist in `dbo.ct`."},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := walkTarget()
			target.Schema.Schemas[0].Tables = append(target.Schema.Schemas[0].Tables, &metadata.TableMetadata{
				Name:    "ct",
				Columns: []*metadata.ColumnMetadata{{Name: "n", Type: "nvarchar(10)"}, {Name: "i", Type: "int"}},
			})
			got := walkFindings(t, tt.sql, target)
			if tt.statement < 0 {
				if len(got) != 0 {
					t.Errorf("findings = %+v, want none", got)
				}
				return
			}
			if len(got) != 1 || got[0].Statement != tt.statement || got[0].Message != tt.message {
				t.Errorf("findings = %+v, want statement %d: %q", got, tt.statement, tt.message)
			}
			if (&walkEngine{}).Codes(tt.message) == nil {
				t.Errorf("walkCodes maps no error to %q", tt.message)
			}
		})
	}
}
