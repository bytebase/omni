package parser

import (
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// TestFederationStatementShapes holds the parts CREATE, ALTER and DROP
// FEDERATION require. SQL Server 2022 does not know the statements — it
// answers Msg 343 "Unknown object type 'FEDERATION'" for every one of them,
// measured on the pinned server — so the rows of testdata/compat.tsv cannot
// show what their grammar is; omni reads it from TSql160.g for the Azure SQL
// Database these statements come from.
//
// Grammar: TSql160.g:7321 createFederationStatement_FEDERATION,
// TSql160.g:7334 alterFederationStatement_FEDERATION,
// TSql160.g:9010 dropFederationStatement_FEDERATION.
func TestFederationStatementShapes(t *testing.T) {
	t.Run("create", func(t *testing.T) {
		list := ParseAndCheck(t, "CREATE FEDERATION MyFederation (dist_key INT RANGE)")
		stmt, ok := list.Items[0].(*ast.CreateFederationStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.CreateFederationStmt", list.Items[0])
		}
		if stmt.Name != "MyFederation" || stmt.DistributionName != "dist_key" {
			t.Errorf("Name, DistributionName = %q, %q", stmt.Name, stmt.DistributionName)
		}
		if stmt.DataType == nil || stmt.DataType.Name != "INT" {
			t.Errorf("DataType = %v, want INT", stmt.DataType)
		}
		checkLocation(t, "CREATE FEDERATION MyFederation (dist_key INT RANGE)", "CreateFederationStmt", stmt.Loc)
	})

	// SPLIT AT and DROP AT each take one option: a name, an equals sign and a
	// value, in parentheses.
	t.Run("alter", func(t *testing.T) {
		for _, tt := range []struct{ sql, kind, distribution string }{
			{"ALTER FEDERATION a SPLIT AT (k = 1)", "SPLIT", "k"},
			{"ALTER FEDERATION a DROP AT (LOW k = 1)", "DROP LOW", "k"},
			{"ALTER FEDERATION a DROP AT (HIGH k = 1)", "DROP HIGH", "k"},
		} {
			list := ParseAndCheck(t, tt.sql)
			stmt, ok := list.Items[0].(*ast.AlterFederationStmt)
			if !ok {
				t.Fatalf("Parse(%q): %T, want *ast.AlterFederationStmt", tt.sql, list.Items[0])
			}
			if stmt.Kind != tt.kind || stmt.DistributionName != tt.distribution {
				t.Errorf("Parse(%q): Kind, DistributionName = %q, %q, want %q, %q", tt.sql, stmt.Kind, stmt.DistributionName, tt.kind, tt.distribution)
			}
			if stmt.Boundary == nil {
				t.Errorf("Parse(%q): no boundary", tt.sql)
			}
			checkLocation(t, tt.sql, "AlterFederationStmt", stmt.Loc)
		}
	})

	t.Run("drop", func(t *testing.T) {
		list := ParseAndCheck(t, "DROP FEDERATION MyFederation")
		stmt, ok := list.Items[0].(*ast.DropFederationStmt)
		if !ok {
			t.Fatalf("got %T, want *ast.DropFederationStmt", list.Items[0])
		}
		if stmt.Name != "MyFederation" {
			t.Errorf("Name = %q, want MyFederation", stmt.Name)
		}
	})

	// Each row is a near-miss of one of those statements in
	// testdata/compat.tsv: a part the grammar names is missing or another word
	// stands in for it.
	for _, sql := range []string{
		"CREATE FEDERATION(a a RANGE)",
		"CREATE FEDERATION a a a RANGE)",
		"CREATE FEDERATION a(a RANGE)",
		"CREATE FEDERATION a(a a)",
		"CREATE FEDERATION a(a a RANGE",
		"CREATE FEDERATION a()",
		"ALTER FEDERATION",
		"ALTER FEDERATION SPLIT AT(a = 1)",
		"ALTER FEDERATION a AT(a = 1)",
		"ALTER FEDERATION a SPLIT(a = 1)",
		"ALTER FEDERATION a SPLIT AT a = 1)",
		"ALTER FEDERATION a SPLIT AT(= 1)",
		"ALTER FEDERATION a SPLIT AT(a 1)",
		"ALTER FEDERATION a SPLIT AT(a =)",
		"ALTER FEDERATION a SPLIT AT(a = 1",
		"ALTER FEDERATION a SPLIT AT",
		"ALTER FEDERATION a SPLIT AT(zzz a = 1)",
		"ALTER FEDERATION a DROP(LOW a = 1)",
		"ALTER FEDERATION a DROP AT LOW a = 1)",
		"ALTER FEDERATION a DROP AT(a = 1)",
		"ALTER FEDERATION a DROP AT(",
		"ALTER FEDERATION a DROP AT(LOW zzz a = 1)",
		"ALTER FEDERATION a DROP AT(MAINTAIN_INDEX a = 1)",
		"ALTER FEDERATION a DROP AT(HINT a = 1)",
		"DROP FEDERATION",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q): accepted, want rejected", sql)
		} else if _, ok := err.(*ParseError); !ok {
			t.Errorf("Parse(%q): %T, want *ParseError", sql, err)
		}
	}
}
