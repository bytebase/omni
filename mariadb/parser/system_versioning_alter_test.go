package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestSystemVersioningAlterAccept covers ALTER TABLE system-versioning DDL,
// container-verified against mariadb:11.8.8 (GAP). The ADD COLUMN forms already
// work via the shared column-def path; they are pinned here as regression.
func TestSystemVersioningAlterAccept(t *testing.T) {
	accept := []string{
		"ALTER TABLE t ADD SYSTEM VERSIONING",
		"ALTER TABLE t DROP SYSTEM VERSIONING",
		"ALTER TABLE t ADD PERIOD FOR SYSTEM_TIME(rs, re)",
		"ALTER TABLE t ADD PERIOD FOR SYSTEM_TIME (rs, re)",
		"ALTER TABLE t DROP PERIOD FOR SYSTEM_TIME",
		// shared column-def path (regression)
		"ALTER TABLE t ADD COLUMN rs TIMESTAMP(6) GENERATED ALWAYS AS ROW START",
		"ALTER TABLE t ADD COLUMN z INT WITH SYSTEM VERSIONING",
	}
	for _, sql := range accept {
		t.Run(sql, func(t *testing.T) { ParseAndCheck(t, sql) })
	}
}

// TestSystemVersioningAlterReject covers the 1064 edges (AGREE_REJECT vs the
// container).
func TestSystemVersioningAlterReject(t *testing.T) {
	reject := []string{
		"ALTER TABLE t ADD SYSTEM",                 // no VERSIONING
		"ALTER TABLE t ADD PERIOD FOR SYSTEM_TIME", // no (start, end)
		"ALTER TABLE t DROP PERIOD FOR",            // no SYSTEM_TIME
	}
	for _, sql := range reject {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

// TestSystemVersioningAlterAST verifies the ALTER command type and period
// columns land on the AST (and reach outfuncs).
func TestSystemVersioningAlterAST(t *testing.T) {
	t.Run("add period", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE t ADD PERIOD FOR SYSTEM_TIME(rs, re)")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATAddPeriod {
			t.Errorf("Type = %v, want ATAddPeriod", cmd.Type)
		}
		if cmd.PeriodStartCol != "rs" || cmd.PeriodEndCol != "re" {
			t.Errorf("period cols = (%q, %q), want (rs, re)", cmd.PeriodStartCol, cmd.PeriodEndCol)
		}
		if got := ast.NodeToString(stmt); !strings.Contains(got, ":period_for_system_time (rs, re)") {
			t.Errorf("NodeToString missing period clause:\n%s", got)
		}
	})
	t.Run("add system versioning", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE t ADD SYSTEM VERSIONING")
		if stmt.Commands[0].Type != ast.ATAddSystemVersioning {
			t.Errorf("Type = %v, want ATAddSystemVersioning", stmt.Commands[0].Type)
		}
	})
	t.Run("drop system versioning", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE t DROP SYSTEM VERSIONING")
		if stmt.Commands[0].Type != ast.ATDropSystemVersioning {
			t.Errorf("Type = %v, want ATDropSystemVersioning", stmt.Commands[0].Type)
		}
	})
}
