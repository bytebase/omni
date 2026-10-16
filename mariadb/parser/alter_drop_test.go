package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mariadb/ast"
)

// TestAlterTableDropIfExists covers IF EXISTS on the DROP clauses that take
// opt_if_exists_table_element in MariaDB's alter_list_item, and on DROP
// PARTITION. Container-verified by TestMariaDBDropClauseOracle. MySQL 8.0
// rejects the COLUMN, CONSTRAINT, INDEX, and FOREIGN KEY forms with IF EXISTS.
func TestAlterTableDropIfExists(t *testing.T) {
	cases := []struct {
		sql      string
		typ      ast.AlterTableCmdType
		name     string
		ifExists bool
	}{
		{"ALTER TABLE t DROP CONSTRAINT IF EXISTS c", ast.ATDropConstraint, "c", true},
		{"ALTER TABLE t DROP CONSTRAINT c", ast.ATDropConstraint, "c", false},
		{"ALTER TABLE t DROP INDEX IF EXISTS i", ast.ATDropIndex, "i", true},
		{"ALTER TABLE t DROP KEY IF EXISTS i", ast.ATDropIndex, "i", true},
		{"ALTER TABLE t DROP INDEX i", ast.ATDropIndex, "i", false},
		{"ALTER TABLE t DROP FOREIGN KEY IF EXISTS f", ast.ATDropConstraint, "f", true},
		{"ALTER TABLE t DROP FOREIGN KEY f", ast.ATDropConstraint, "f", false},
		{"ALTER TABLE t DROP COLUMN IF EXISTS c", ast.ATDropColumn, "c", true},
		{"ALTER TABLE t DROP IF EXISTS c", ast.ATDropColumn, "c", true},
		{"ALTER TABLE t DROP INDEX IF EXISTS `if`", ast.ATDropIndex, "if", true},
		// DROP [CONSTRAINT] PRIMARY KEY never takes IF EXISTS.
		{"ALTER TABLE t DROP CONSTRAINT PRIMARY KEY", ast.ATDropConstraint, "PRIMARY", false},
		{"ALTER TABLE t DROP PRIMARY KEY", ast.ATDropConstraint, "PRIMARY", false},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			stmt := parseSeqStmt[*ast.AlterTableStmt](t, tc.sql)
			cmd := stmt.Commands[0]
			if cmd.Type != tc.typ || cmd.Name != tc.name || cmd.IfExists != tc.ifExists {
				t.Errorf("got Type=%v Name=%q IfExists=%v, want Type=%v Name=%q IfExists=%v",
					cmd.Type, cmd.Name, cmd.IfExists, tc.typ, tc.name, tc.ifExists)
			}
			wantLoc := ast.Loc{Start: strings.Index(tc.sql, "DROP"), End: len(tc.sql)}
			if cmd.Loc != wantLoc {
				t.Errorf("cmd Loc = %+v, want %+v", cmd.Loc, wantLoc)
			}
		})
	}

	t.Run("partition", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t, "ALTER TABLE t DROP PARTITION IF EXISTS p0, p1")
		cmd := stmt.Commands[0]
		if cmd.Type != ast.ATDropPartition || !cmd.IfExists {
			t.Errorf("got Type=%v IfExists=%v, want ATDropPartition with IfExists", cmd.Type, cmd.IfExists)
		}
		if len(cmd.PartitionNames) != 2 || cmd.PartitionNames[0] != "p0" || cmd.PartitionNames[1] != "p1" {
			t.Errorf("PartitionNames = %v, want [p0 p1]", cmd.PartitionNames)
		}
	})

	t.Run("several commands", func(t *testing.T) {
		stmt := parseSeqStmt[*ast.AlterTableStmt](t,
			"ALTER TABLE t DROP CONSTRAINT IF EXISTS c1, DROP INDEX i1, DROP FOREIGN KEY IF EXISTS f1")
		var got []bool
		for _, cmd := range stmt.Commands {
			got = append(got, cmd.IfExists)
		}
		if len(got) != 3 || !got[0] || got[1] || !got[2] {
			t.Errorf("IfExists per command = %v, want [true false true]", got)
		}
		if s := ast.NodeToString(stmt); strings.Count(s, ":if_exists true") != 2 {
			t.Errorf("NodeToString should carry two :if_exists, got:\n%s", s)
		}
	})
}

// TestAlterTableDropReject pins the DROP forms MariaDB 11.8 rejects at parse
// (1064), container-verified by TestMariaDBDropClauseOracle.
func TestAlterTableDropReject(t *testing.T) {
	for _, sql := range []string{
		"ALTER TABLE t DROP CHECK ck",           // MySQL 8.0 syntax; MariaDB drops a CHECK with DROP CONSTRAINT
		"ALTER TABLE t DROP CHECK IF EXISTS ck", // same
		"ALTER TABLE t DROP CONSTRAINT IF EXISTS",
		"ALTER TABLE t DROP INDEX IF i",
		"ALTER TABLE t DROP INDEX EXISTS i",
		"ALTER TABLE t DROP CONSTRAINT PRIMARY",               // KEY required after CONSTRAINT PRIMARY
		"ALTER TABLE t DROP CONSTRAINT IF EXISTS PRIMARY KEY", // PRIMARY KEY takes no IF EXISTS
		"ALTER TABLE t DROP PRIMARY KEY IF EXISTS",
		"ALTER TABLE t DROP PARTITION IF EXISTS",
	} {
		t.Run(sql, func(t *testing.T) { ParseExpectError(t, sql) })
	}
}

func TestCollectAlterTableDropIfExists(t *testing.T) {
	for _, tc := range []struct{ sql, rule string }{
		{"ALTER TABLE t1 DROP INDEX IF EXISTS ", "index_ref"},
		{"ALTER TABLE t1 DROP KEY IF EXISTS ", "index_ref"},
		{"ALTER TABLE t1 DROP CONSTRAINT IF EXISTS ", "constraint_ref"},
		{"ALTER TABLE t1 DROP FOREIGN KEY IF EXISTS ", "constraint_ref"},
		{"ALTER TABLE t1 DROP COLUMN IF EXISTS ", "columnref"},
		{"ALTER TABLE t1 DROP CONSTRAINT ", "constraint_ref"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			cs := Collect(tc.sql, len(tc.sql))
			if cs == nil || !cs.HasRule(tc.rule) {
				t.Errorf("expected %s rule", tc.rule)
			}
		})
	}

	t.Run("no CHECK after DROP", func(t *testing.T) {
		cs := Collect("ALTER TABLE t1 DROP ", 20)
		if cs == nil {
			t.Fatal("expected non-nil candidates")
		}
		if cs.HasToken(kwCHECK) {
			t.Error("DROP CHECK is MySQL syntax; MariaDB rejects it, so CHECK must not be offered")
		}
		if !cs.HasToken(kwCONSTRAINT) || !cs.HasToken(kwPARTITION) {
			t.Error("expected CONSTRAINT and PARTITION after DROP")
		}
	})
}
