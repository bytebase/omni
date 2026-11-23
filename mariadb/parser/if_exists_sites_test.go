package parser

import (
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestIfExistsSites covers the IF [NOT] EXISTS sites of sql_yacc.yy that
// omni read late: the flag reaches the AST. MariaDB's verdicts are the
// <alternative>+opt_if_[not_]exists rows of testdata/compat.tsv.
func TestIfExistsSites(t *testing.T) {
	t.Run("CREATE VIEW", func(t *testing.T) {
		for _, sql := range []string{
			"CREATE VIEW IF NOT EXISTS v AS SELECT 1",
			"CREATE ALGORITHM = MERGE DEFINER = CURRENT_USER SQL SECURITY INVOKER VIEW IF NOT EXISTS v AS SELECT 1",
		} {
			if v := parseSeqStmt[*ast.CreateViewStmt](t, sql); !v.IfNotExists || v.Name.Name != "v" {
				t.Errorf("%s: IfNotExists=%v Name=%q", sql, v.IfNotExists, v.Name.Name)
			}
		}
		if v := parseSeqStmt[*ast.CreateViewStmt](t, "CREATE VIEW v AS SELECT 1"); v.IfNotExists {
			t.Error("IfNotExists set without IF NOT EXISTS")
		}
		if _, err := Parse("CREATE OR REPLACE VIEW IF NOT EXISTS v AS SELECT 1"); err == nil {
			t.Error("OR REPLACE with IF NOT EXISTS: no error")
		}
	})

	t.Run("RENAME TABLE", func(t *testing.T) {
		r := parseSeqStmt[*ast.RenameTableStmt](t, "RENAME TABLE IF EXISTS a TO b, c TO d")
		if !r.IfExists || len(r.Pairs) != 2 {
			t.Errorf("IfExists=%v pairs=%d", r.IfExists, len(r.Pairs))
		}
	})

	for _, tc := range []struct {
		sql         string
		typ         ast.AlterTableCmdType
		ifExists    bool
		ifNotExists bool
	}{
		{"ALTER TABLE t ADD PARTITION IF NOT EXISTS (PARTITION p1 VALUES LESS THAN (10))", ast.ATAddPartition, false, true},
		{"ALTER TABLE t ADD PARTITION IF NOT EXISTS NO_WRITE_TO_BINLOG PARTITIONS 2", ast.ATAddPartition, false, true},
		{"ALTER TABLE t CHANGE COLUMN IF EXISTS a b INT", ast.ATChangeColumn, true, false},
		{"ALTER TABLE t MODIFY IF EXISTS a INT", ast.ATModifyColumn, true, false},
		{"ALTER TABLE t ALTER COLUMN IF EXISTS a SET DEFAULT 1", ast.ATAlterColumnDefault, true, false},
		{"ALTER TABLE t ALTER IF EXISTS a DROP DEFAULT", ast.ATAlterColumnDefault, true, false},
		{"ALTER TABLE t RENAME COLUMN IF EXISTS a TO b", ast.ATRenameColumn, true, false},
		{"ALTER TABLE t RENAME INDEX IF EXISTS a TO b", ast.ATRenameIndex, true, false},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSeqStmt[*ast.AlterTableStmt](t, tc.sql)
			if len(stmt.Commands) != 1 {
				t.Fatalf("commands = %d, want 1", len(stmt.Commands))
			}
			cmd := stmt.Commands[0]
			if cmd.Type != tc.typ || cmd.IfExists != tc.ifExists || cmd.IfNotExists != tc.ifNotExists {
				t.Errorf("Type=%v IfExists=%v IfNotExists=%v, want %v %v %v", cmd.Type, cmd.IfExists, cmd.IfNotExists, tc.typ, tc.ifExists, tc.ifNotExists)
			}
		})
	}

	// IF must be followed by EXISTS: the drop-4 near-misses of the
	// +opt_if_exists rows, which mariadb:11.8.8 answers 1064.
	for _, sql := range []string{
		"DROP FUNCTION IF a",
		"DROP PROCEDURE IF a",
		"DROP USER IF a",
		"DROP ROLE IF a",
		"DROP EVENT IF a",
		"DROP TRIGGER IF a",
		"DROP SERVER IF a",
		"ALTER USER IF a",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("%s: no error", sql)
		}
	}
}
