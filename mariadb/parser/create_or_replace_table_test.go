package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestCreateOrReplaceTable covers CREATE OR REPLACE TABLE: the AST records
// OrReplace and the statement's Loc starts at CREATE. Container-verified by
// TestMariaDBOrReplaceTableOracle.
func TestCreateOrReplaceTable(t *testing.T) {
	cases := []struct {
		sql       string
		orReplace bool
		temporary bool
	}{
		{"CREATE OR REPLACE TABLE t (id INT)", true, false},
		{"CREATE OR REPLACE TEMPORARY TABLE t (id INT)", true, true},
		{"CREATE OR REPLACE TABLE t LIKE t2", true, false},
		{"CREATE OR REPLACE TABLE t AS SELECT 1 AS a", true, false},
		{"create or replace table t (id INT)", true, false},
		{"CREATE TABLE t (id INT)", false, false},
		{"CREATE TEMPORARY TABLE t (id INT)", false, true},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSeqStmt[*ast.CreateTableStmt](t, tc.sql)
			if stmt.OrReplace != tc.orReplace || stmt.Temporary != tc.temporary {
				t.Errorf("OrReplace=%v Temporary=%v, want %v %v", stmt.OrReplace, stmt.Temporary, tc.orReplace, tc.temporary)
			}
			if stmt.IfNotExists {
				t.Error("IfNotExists set without IF NOT EXISTS")
			}
			if want := (ast.Loc{Start: 0, End: len(tc.sql)}); stmt.Loc != want {
				t.Errorf("Loc = %+v, want %+v", stmt.Loc, want)
			}
		})
	}

	t.Run("loc after a preceding statement", func(t *testing.T) {
		const sql = "SELECT 1; CREATE OR REPLACE TABLE t (id INT)"
		list := ParseAndCheck(t, sql)
		stmt, ok := list.Items[1].(*ast.CreateTableStmt)
		if !ok {
			t.Fatalf("Items[1] = %T, want *ast.CreateTableStmt", list.Items[1])
		}
		if want := (ast.Loc{Start: strings.Index(sql, "CREATE"), End: len(sql)}); stmt.Loc != want {
			t.Errorf("Loc = %+v, want %+v", stmt.Loc, want)
		}
	})

	t.Run("outfuncs", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateTableStmt](t, "CREATE OR REPLACE TABLE t (id INT)")
		if got := ast.NodeToString(stmt); !strings.HasPrefix(got, "{CREATE_TABLE :loc 0 :or_replace true") {
			t.Errorf("NodeToString = %s", got)
		}
	})
}

// TestCreateOrReplaceTableReject pins the forms MariaDB 11.8 rejects while
// parsing. OR REPLACE with IF NOT EXISTS fails with ER_WRONG_USAGE (1221),
// raised by the grammar action through LEX::check_create_options; the others
// are 1064.
func TestCreateOrReplaceTableReject(t *testing.T) {
	for _, sql := range []string{
		"CREATE OR REPLACE TABLE IF NOT EXISTS t (id INT)",
		"CREATE OR REPLACE TEMPORARY TABLE IF NOT EXISTS t (id INT)",
		"CREATE OR REPLACE TABLE IF NOT EXISTS t LIKE t2",
	} {
		t.Run(sql, func(t *testing.T) {
			_, err := Parse(sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want *ParseError", sql, err)
			}
			if pe.Message != "Incorrect usage of OR REPLACE and IF NOT EXISTS" {
				t.Errorf("Message = %q", pe.Message)
			}
			if want := strings.Index(sql, "IF NOT EXISTS"); pe.Position != want {
				t.Errorf("Position = %d, want %d", pe.Position, want)
			}
		})
	}
	for _, sql := range []string{
		"CREATE OR TABLE t (id INT)",
		"CREATE OR REPLACE TABLE IF EXISTS t (id INT)",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}
