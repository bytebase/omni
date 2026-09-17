package parser

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// withClauseOf returns the WithClause attached to a statement-level DML node.
func withClauseOf(t *testing.T, n ast.Node) *ast.WithClause {
	t.Helper()
	switch s := n.(type) {
	case *ast.SelectStmt:
		return s.WithClause
	case *ast.InsertStmt:
		return s.WithClause
	case *ast.UpdateStmt:
		return s.WithClause
	case *ast.DeleteStmt:
		return s.WithClause
	case *ast.MergeStmt:
		return s.WithClause
	}
	t.Fatalf("unexpected statement type %T", n)
	return nil
}

func cteNames(wc *ast.WithClause) []string {
	var names []string
	for _, item := range wc.CTEs.Items {
		names = append(names, item.(*ast.CommonTableExpr).Name)
	}
	return names
}

// TestWithClauseDML covers WITH (CTE) followed by each statement T-SQL allows
// it on. Before the fix, WITH was always routed to the SELECT parser, which
// consumed the CTE list and returned a typed-nil *SelectStmt when the next
// token was not SELECT; the caller appended that nil as a statement and the
// DML that followed was parsed a second time without its CTEs (BYT-10223).
func TestWithClauseDML(t *testing.T) {
	tests := []struct {
		name     string
		sql      string
		wantType ast.Node
		wantCTEs []string
	}{
		{
			name:     "select",
			sql:      "WITH t AS (SELECT 1 AS a) SELECT a FROM t",
			wantType: &ast.SelectStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "insert select",
			sql:      "WITH t AS (SELECT 1 AS a) INSERT INTO x (a) SELECT a FROM t",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "insert without INTO",
			sql:      "WITH t AS (SELECT 1 AS a) INSERT x (a) SELECT a FROM t",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "insert top",
			sql:      "WITH t AS (SELECT 1 AS a) INSERT TOP (10) INTO x (a) SELECT a FROM t",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "update",
			sql:      "WITH t AS (SELECT 1 AS a) UPDATE x SET a = t.a FROM x JOIN t ON x.a = t.a",
			wantType: &ast.UpdateStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "update the cte itself",
			sql:      "WITH t AS (SELECT a FROM x) UPDATE t SET a = 1",
			wantType: &ast.UpdateStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "delete",
			sql:      "WITH t AS (SELECT 1 AS a) DELETE FROM x WHERE a IN (SELECT a FROM t)",
			wantType: &ast.DeleteStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "delete the cte itself",
			sql:      "WITH t AS (SELECT a FROM x) DELETE FROM t",
			wantType: &ast.DeleteStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "merge",
			sql:      "WITH t AS (SELECT 1 AS a) MERGE INTO x USING t ON x.a = t.a WHEN NOT MATCHED THEN INSERT (a) VALUES (t.a);",
			wantType: &ast.MergeStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "multiple ctes before insert",
			sql:      "WITH a AS (SELECT 1 AS v), b (v) AS (SELECT v FROM a) INSERT INTO x (v) SELECT v FROM b",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"a", "b"},
		},
		{
			name:     "xmlnamespaces only before select",
			sql:      "WITH XMLNAMESPACES ('http://x' AS ns) SELECT a FROM x",
			wantType: &ast.SelectStmt{},
			wantCTEs: nil,
		},
		{
			name:     "xmlnamespaces only before insert",
			sql:      "WITH XMLNAMESPACES ('http://x' AS ns) INSERT INTO x (a) SELECT a FROM y",
			wantType: &ast.InsertStmt{},
			wantCTEs: nil,
		},
		{
			name:     "xmlnamespaces then cte before insert",
			sql:      "WITH XMLNAMESPACES ('http://x' AS ns), t AS (SELECT 1 AS a) INSERT INTO x (a) SELECT a FROM t",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name:     "leading semicolon",
			sql:      ";WITH t AS (SELECT 1 AS a) INSERT INTO x (a) SELECT a FROM t;",
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"t"},
		},
		{
			name: "customer statement",
			sql: `-- txn-mode = off

;WITH TargetVenues AS (
    SELECT v.VenueId
    FROM Venue v
    WHERE v.IsActive = 1 AND (v.Name LIKE 'XP%')
)
INSERT INTO VenueMeta (VenueId, Attribute, Value, IsActive)
SELECT tv.VenueId, 'Roller.Venue.Cap', N'{"split": true}', 1
FROM TargetVenues tv
WHERE NOT EXISTS (
    SELECT 1
    FROM VenueMeta vm
    WHERE vm.VenueId = tv.VenueId
      AND vm.Attribute = 'Roller.Venue.Cap'
      AND vm.IsActive = 1
);`,
			wantType: &ast.InsertStmt{},
			wantCTEs: []string{"TargetVenues"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := ParseAndCheck(t, tt.sql)
			if list.Len() != 1 {
				t.Fatalf("got %d statements, want 1: %s", list.Len(), ast.NodeToString(list))
			}
			stmt := list.Items[0]
			if got, want := typeName(stmt), typeName(tt.wantType); got != want {
				t.Fatalf("statement type = %s, want %s", got, want)
			}

			wc := withClauseOf(t, stmt)
			if wc == nil {
				t.Fatal("WithClause is nil: the CTE list was dropped")
			}
			if got := cteNames(wc); strings.Join(got, ",") != strings.Join(tt.wantCTEs, ",") {
				t.Errorf("CTE names = %v, want %v", got, tt.wantCTEs)
			}

			// The statement starts at the WITH keyword, the CTE list ends
			// before the DML keyword, and the statement covers the DML.
			withPos := strings.Index(tt.sql, "WITH")
			loc := ast.NodeLoc(stmt)
			if loc.Start != withPos {
				t.Errorf("Loc.Start = %d, want %d (the WITH keyword)", loc.Start, withPos)
			}
			if wc.Loc.Start != withPos {
				t.Errorf("WithClause.Loc.Start = %d, want %d", wc.Loc.Start, withPos)
			}
			if wc.Loc.End <= wc.Loc.Start || wc.Loc.End > loc.End {
				t.Errorf("WithClause.Loc = %+v not inside statement Loc %+v", wc.Loc, loc)
			}
			trimmed := strings.TrimRight(tt.sql, "; \n")
			if loc.End != len(trimmed) {
				t.Errorf("Loc.End = %d, want %d (end of statement text)", loc.End, len(trimmed))
			}

			// The dump must carry the CTE so deparse/compare tooling sees it.
			if dump := ast.NodeToString(stmt); !strings.Contains(dump, ":with ") {
				t.Errorf("NodeToString lacks :with clause: %s", dump)
			}

			// Walking must reach the CTE body and must not panic.
			sawCTE := false
			ast.Inspect(stmt, func(n ast.Node) bool {
				if _, ok := n.(*ast.CommonTableExpr); ok {
					sawCTE = true
				}
				return true
			})
			if !sawCTE && len(tt.wantCTEs) > 0 {
				t.Error("Inspect did not visit the CommonTableExpr")
			}
		})
	}
}

// TestWithClauseDMLNested checks WITH ... DML inside compound statements,
// which dispatch through parseStmt just like the top level.
func TestWithClauseDMLNested(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{
			name: "begin end block",
			sql:  "BEGIN WITH t AS (SELECT 1 AS a) INSERT INTO x (a) SELECT a FROM t; SELECT 1; END",
		},
		{
			name: "if body",
			sql:  "IF 1 = 1 WITH t AS (SELECT 1 AS a) DELETE FROM x WHERE a IN (SELECT a FROM t)",
		},
		{
			name: "procedure body",
			sql:  "CREATE PROCEDURE p AS WITH t AS (SELECT 1 AS a) UPDATE x SET a = 1 FROM x JOIN t ON x.a = t.a",
		},
		{
			name: "trigger body",
			sql:  "CREATE TRIGGER tr ON x AFTER INSERT AS WITH t AS (SELECT a FROM inserted) INSERT INTO y (a) SELECT a FROM t",
		},
		{
			name: "two statements back to back",
			sql:  "WITH t AS (SELECT 1 AS a) INSERT INTO x (a) SELECT a FROM t\nWITH u AS (SELECT 2 AS a) INSERT INTO x (a) SELECT a FROM u",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := ParseAndCheck(t, tt.sql)
			nWith := 0
			for _, item := range list.Items {
				ast.Inspect(item, func(n ast.Node) bool {
					switch s := n.(type) {
					case *ast.InsertStmt:
						if s.WithClause != nil {
							nWith++
						}
					case *ast.UpdateStmt:
						if s.WithClause != nil {
							nWith++
						}
					case *ast.DeleteStmt:
						if s.WithClause != nil {
							nWith++
						}
					case *ast.SelectStmt:
						if s.WithClause != nil {
							nWith++
						}
					}
					return true
				})
			}
			wantWith := strings.Count(tt.sql, "WITH")
			if nWith != wantWith {
				t.Errorf("found %d DML nodes with a WithClause, want %d: %s", nWith, wantWith, ast.NodeToString(list))
			}
		})
	}
}

// TestWithClauseDMLErrors checks that a CTE list not followed by a statement
// that can own it is a syntax error at the offending token, instead of a
// silently dropped CTE.
func TestWithClauseDMLErrors(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		wantErr string
	}{
		{"end of input", "WITH t AS (SELECT 1 AS a)", `syntax error at end of input`},
		{"semicolon", "WITH t AS (SELECT 1 AS a);", `syntax error at or near ";"`},
		{"ddl", "WITH t AS (SELECT 1 AS a) CREATE TABLE x (a INT)", `syntax error at or near "CREATE"`},
		{"exec", "WITH t AS (SELECT 1 AS a) EXEC p", `syntax error at or near "EXEC"`},
		{"trailing comma before insert", "WITH t AS (SELECT 1 AS a), INSERT INTO x (a) SELECT a FROM t", `syntax error at or near "INSERT"`},
		{"trailing comma before select", "WITH t AS (SELECT 1 AS a), SELECT a FROM t", `syntax error at or near "SELECT"`},
		{"cte in subquery position", "SELECT * FROM (WITH t AS (SELECT 1 AS a) INSERT INTO x SELECT a FROM t) d", `syntax error at or near "INSERT"`},
		{"cte in view body", "CREATE VIEW v AS WITH t AS (SELECT 1 AS a) INSERT INTO x SELECT a FROM t", `syntax error at or near "INSERT"`},
		{"empty cte list before insert", "WITH INSERT INTO x (a) VALUES (1)", `syntax error at or near "INSERT"`},
		{"empty cte list before select", "WITH SELECT 1", `syntax error at or near "SELECT"`},
		{"dml inside EXISTS", "SELECT 1 WHERE EXISTS (WITH c AS (SELECT 1 AS a) DELETE FROM x)", `syntax error at or near "DELETE"`},
		{"dml inside EXISTS without cte", "SELECT 1 WHERE EXISTS (DELETE FROM x)", `syntax error at or near "DELETE"`},
		{"dangling comma after xmlnamespaces before insert", "WITH XMLNAMESPACES ('http://x' AS ns), INSERT INTO x (a) SELECT a FROM y", `syntax error at or near "INSERT"`},
		{"dangling comma after xmlnamespaces before select", "WITH XMLNAMESPACES ('http://x' AS ns), SELECT 1", `syntax error at or near "SELECT"`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded with %d statements, want error %q", tt.sql, list.Len(), tt.wantErr)
			}
			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Parse(%q) error = %q, want it to contain %q", tt.sql, err.Error(), tt.wantErr)
			}
		})
	}
}

// TestWithClauseDMLOracle checks accept/reject parity with SQL Server for
// WITH followed by each statement kind.
func TestWithClauseDMLOracle(t *testing.T) {
	oracle := startParserOracle(t)

	tests := []struct {
		name string
		sql  string
	}{
		{"with_select", "WITH c AS (SELECT a FROM dbo.t) SELECT a FROM c"},
		{"with_insert_select", "WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c"},
		{"with_insert_no_into", "WITH c AS (SELECT a FROM dbo.t) INSERT dbo.t (a) SELECT a FROM c"},
		{"with_insert_top", "WITH c AS (SELECT a FROM dbo.t) INSERT TOP (1) INTO dbo.t (a) SELECT a FROM c"},
		{"with_update_from", "WITH c AS (SELECT a FROM dbo.t) UPDATE dbo.t SET a = c.a FROM dbo.t JOIN c ON dbo.t.a = c.a"},
		{"with_update_cte", "WITH c AS (SELECT a FROM dbo.t) UPDATE c SET a = 1"},
		{"with_delete_where", "WITH c AS (SELECT a FROM dbo.t) DELETE FROM dbo.t WHERE a IN (SELECT a FROM c)"},
		{"with_delete_cte", "WITH c AS (SELECT a FROM dbo.t) DELETE FROM c"},
		{"with_merge", "WITH c AS (SELECT a FROM dbo.t) MERGE INTO dbo.t USING c ON dbo.t.a = c.a WHEN NOT MATCHED THEN INSERT (a) VALUES (c.a);"},
		{"with_two_ctes_insert", "WITH c1 AS (SELECT a FROM dbo.t), c2 (a) AS (SELECT a FROM c1) INSERT INTO dbo.t (a) SELECT a FROM c2"},
		{"with_leading_semicolon", ";WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c;"},
		{"with_in_begin_end", "BEGIN WITH c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c; END"},
		{"with_alone_invalid", "WITH c AS (SELECT a FROM dbo.t)"},
		{"with_semicolon_invalid", "WITH c AS (SELECT a FROM dbo.t);"},
		{"with_create_invalid", "WITH c AS (SELECT a FROM dbo.t) CREATE TABLE dbo.u (a INT)"},
		{"with_exec_invalid", "WITH c AS (SELECT a FROM dbo.t) EXEC sp_who"},
		{"with_trailing_comma_invalid", "WITH c AS (SELECT a FROM dbo.t), INSERT INTO dbo.t (a) SELECT a FROM c"},
		{"with_empty_before_insert_invalid", "WITH INSERT INTO dbo.t (a) VALUES (1)"},
		{"with_empty_before_select_invalid", "WITH SELECT 1"},
		{"with_xmlnamespaces_only_select", "WITH XMLNAMESPACES ('http://x' AS ns) SELECT a FROM dbo.t"},
		{"with_xmlnamespaces_only_insert", "WITH XMLNAMESPACES ('http://x' AS ns) INSERT INTO dbo.t (a) SELECT a FROM dbo.t"},
		{"with_dml_in_exists_invalid", "SELECT 1 WHERE EXISTS (WITH c AS (SELECT a FROM dbo.t) DELETE FROM dbo.t)"},
		{"with_xmlnamespaces_comma_cte_insert", "WITH XMLNAMESPACES ('http://x' AS ns), c AS (SELECT a FROM dbo.t) INSERT INTO dbo.t (a) SELECT a FROM c"},
		{"with_xmlnamespaces_dangling_comma_insert_invalid", "WITH XMLNAMESPACES ('http://x' AS ns), INSERT INTO dbo.t (a) SELECT a FROM dbo.t"},
		{"with_xmlnamespaces_dangling_comma_select_invalid", "WITH XMLNAMESPACES ('http://x' AS ns), SELECT a FROM dbo.t"},
		// Not covered here: SQL Server rejects a CTE inside a derived table
		// (SELECT * FROM (WITH c AS (...) SELECT ...) d) while omni accepts
		// it. That leniency predates the WITH ... DML fix and needs its own
		// oracle pass over every nested SELECT position.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ssAccepts, err := oracle.canParse(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil

			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s, omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

func typeName(n ast.Node) string {
	return fmt.Sprintf("%T", n)
}
