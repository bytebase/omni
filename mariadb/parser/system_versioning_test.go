package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestSystemVersioningCreateAccept covers system-versioned table DDL in
// CREATE TABLE, container-verified against mariadb:11.8.8 (GAP: MariaDB parses,
// omni rejected). DDL donor = omni/mssql (PERIOD FOR SYSTEM_TIME shape).
func TestSystemVersioningCreateAccept(t *testing.T) {
	accept := []string{
		// table-level WITH SYSTEM VERSIONING (implicit invisible period columns)
		"CREATE TABLE sv1 (x INT) WITH SYSTEM VERSIONING",
		// canonical: explicit period columns + PERIOD FOR SYSTEM_TIME
		"CREATE TABLE sv2 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		// bare "AS ROW START/END" (no GENERATED ALWAYS) + INVISIBLE
		"CREATE TABLE sv3 (x INT, rs TIMESTAMP(6) AS ROW START INVISIBLE, re TIMESTAMP(6) AS ROW END INVISIBLE, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		// column-level WITH / WITHOUT SYSTEM VERSIONING
		"CREATE TABLE sv4 (x INT WITH SYSTEM VERSIONING, y INT WITHOUT SYSTEM VERSIONING)",
		// PERIOD FOR SYSTEM_TIME with a space before the paren
		"CREATE TABLE sv5 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME (rs, re)) WITH SYSTEM VERSIONING",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSystemVersioningCreateReject covers the 1064 edges (AGREE_REJECT vs the
// container): each clause needs its full spec.
func TestSystemVersioningCreateReject(t *testing.T) {
	reject := []string{
		"CREATE TABLE bad1 (x INT) WITH SYSTEM",                           // no VERSIONING
		"CREATE TABLE bad2 (x INT, PERIOD FOR SYSTEM_TIME)",               // no (start, end)
		"CREATE TABLE bad3 (x INT, rs TIMESTAMP GENERATED ALWAYS AS ROW)", // no START/END
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestSystemVersioningAST verifies the period columns, table option, generated
// row bound, and column-level versioning land on the AST (and reach outfuncs).
func TestSystemVersioningAST(t *testing.T) {
	t.Run("period for system_time", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateTableStmt](t,
			"CREATE TABLE sv2 (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING")
		if stmt.PeriodStartCol != "rs" || stmt.PeriodEndCol != "re" {
			t.Errorf("period cols = (%q, %q), want (rs, re)", stmt.PeriodStartCol, stmt.PeriodEndCol)
		}
		if !hasSystemVersioningOption(stmt.Options) {
			t.Errorf("missing SYSTEM VERSIONING table option; got %v", stmt.Options)
		}
		if stmt.Columns[1].Generated == nil || stmt.Columns[1].Generated.RowBound != ast.RowBoundStart {
			t.Errorf("column rs should be GENERATED ... ROW START")
		}
		if stmt.Columns[2].Generated == nil || stmt.Columns[2].Generated.RowBound != ast.RowBoundEnd {
			t.Errorf("column re should be GENERATED ... ROW END")
		}
		if got := ast.NodeToString(stmt); !strings.Contains(got, ":period_for_system_time (rs, re)") {
			t.Errorf("NodeToString missing period clause:\n%s", got)
		}
	})
	t.Run("column versioning", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.CreateTableStmt](t,
			"CREATE TABLE sv4 (x INT WITH SYSTEM VERSIONING, y INT WITHOUT SYSTEM VERSIONING)")
		if stmt.Columns[0].SystemVersioning != ast.ColVersioningWith {
			t.Errorf("column x = %v, want ColVersioningWith", stmt.Columns[0].SystemVersioning)
		}
		if stmt.Columns[1].SystemVersioning != ast.ColVersioningWithout {
			t.Errorf("column y = %v, want ColVersioningWithout", stmt.Columns[1].SystemVersioning)
		}
	})
}

// TestSystemVersioningGeneratedColumnReject: column-level WITH/WITHOUT SYSTEM
// VERSIONING is not allowed on an expression-generated column (MariaDB 1064),
// in either attribute order.
func TestSystemVersioningGeneratedColumnReject(t *testing.T) {
	for _, sql := range []string{
		"CREATE TABLE t (x INT, g INT GENERATED ALWAYS AS (x + 1) WITH SYSTEM VERSIONING) WITH SYSTEM VERSIONING",
		"CREATE TABLE t (x INT, g INT GENERATED ALWAYS AS (x + 1) WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING",
		"CREATE TABLE t (x INT, g INT WITH SYSTEM VERSIONING GENERATED ALWAYS AS (x + 1)) WITH SYSTEM VERSIONING",
		// ROW START/END period columns also reject column-level versioning
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START WITH SYSTEM VERSIONING, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
		"CREATE TABLE t (x INT, rs TIMESTAMP(6) AS ROW START WITHOUT SYSTEM VERSIONING, re TIMESTAMP(6) GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME(rs, re)) WITH SYSTEM VERSIONING",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
	// A generated column without versioning, and a plain WITHOUT column, parse.
	ParseAndCheck(t, "CREATE TABLE t (x INT, g INT GENERATED ALWAYS AS (x + 1)) WITH SYSTEM VERSIONING")
	ParseAndCheck(t, "CREATE TABLE t (x INT, y INT WITHOUT SYSTEM VERSIONING) WITH SYSTEM VERSIONING")
}

func hasSystemVersioningOption(opts []*ast.TableOption) bool {
	for _, o := range opts {
		if o.Name == "SYSTEM VERSIONING" {
			return true
		}
	}
	return false
}
