package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestAlterTableRenameSyntaxMatchesPG(t *testing.T) {
	assertSyntaxMatchesPG(t, []syntaxCase{
		{"ALTER TABLE t RENAME TO u", true},
		{"ALTER TABLE IF EXISTS t RENAME TO u", true},
		{"ALTER TABLE t RENAME COLUMN a TO b", true},
		{"ALTER TABLE ONLY t RENAME a TO b", true},
		{"ALTER TABLE t RENAME CONSTRAINT k TO l", true},
		{"ALTER TABLE t SET SCHEMA s", true},
		{"ALTER FOREIGN TABLE ft RENAME TO u", true},
		{"ALTER FOREIGN TABLE IF EXISTS ft RENAME COLUMN a TO b", true},
		{"ALTER FOREIGN TABLE ft RENAME a TO b", true},
		{"ALTER FOREIGN TABLE ft SET SCHEMA s", true},

		{"ALTER TABLE t RENAME CONSTRAINT k l", false},
		{"ALTER TABLE t RENAME a b", false},
		{"ALTER TABLE t RENAME TO", false},
		{"ALTER TABLE t RENAME COLUMN a TO", false},
		{"ALTER TABLE t RENAME CONSTRAINT k TO", false},
		{"ALTER TABLE t RENAME a TO", false},
		{"ALTER TABLE t RENAME", false},
		{"ALTER TABLE t SET SCHEMA", false},
		{"ALTER TABLE RENAME TO u", false},
		{"ALTER FOREIGN TABLE ft RENAME CONSTRAINT k TO l", false},
		{"ALTER FOREIGN TABLE ft RENAME TO", false},
		{"ALTER FOREIGN TABLE ft RENAME COLUMN a TO", false},
		{"ALTER FOREIGN TABLE ft SET SCHEMA", false},
		{"ALTER FOREIGN TABLE RENAME TO u", false},
	})
}

// TestAlterForeignTableRenameObjectType pins the object types gram.y sets
// for the foreign-table RENAME productions (gram.y:9590, gram.y:9706).
func TestAlterForeignTableRenameObjectType(t *testing.T) {
	tests := []struct {
		sql          string
		renameType   nodes.ObjectType
		relationType nodes.ObjectType
		subname      string
		newname      string
		missingOk    bool
	}{
		{"ALTER FOREIGN TABLE ft RENAME TO u", nodes.OBJECT_FOREIGN_TABLE, 0, "", "u", false},
		{"ALTER FOREIGN TABLE IF EXISTS ft RENAME TO u", nodes.OBJECT_FOREIGN_TABLE, 0, "", "u", true},
		{"ALTER FOREIGN TABLE ft RENAME COLUMN a TO b", nodes.OBJECT_COLUMN, nodes.OBJECT_FOREIGN_TABLE, "a", "b", false},
		{"ALTER FOREIGN TABLE ft RENAME a TO b", nodes.OBJECT_COLUMN, nodes.OBJECT_FOREIGN_TABLE, "a", "b", false},
		{"ALTER TABLE t RENAME TO u", nodes.OBJECT_TABLE, 0, "", "u", false},
		{"ALTER TABLE t RENAME a TO b", nodes.OBJECT_COLUMN, nodes.OBJECT_TABLE, "a", "b", false},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			rn, ok := singleStmt(t, tt.sql).(*nodes.RenameStmt)
			if !ok {
				t.Fatalf("expected RenameStmt")
			}
			if rn.RenameType != tt.renameType || rn.RelationType != tt.relationType {
				t.Fatalf("RenameType, RelationType = %v, %v; want %v, %v", rn.RenameType, rn.RelationType, tt.renameType, tt.relationType)
			}
			if rn.Subname != tt.subname || rn.Newname != tt.newname || rn.MissingOk != tt.missingOk {
				t.Fatalf("Subname, Newname, MissingOk = %q, %q, %v; want %q, %q, %v", rn.Subname, rn.Newname, rn.MissingOk, tt.subname, tt.newname, tt.missingOk)
			}
			if rn.Loc.Start != 0 || rn.Loc.End != len(tt.sql) {
				t.Fatalf("Loc = %+v, want {0 %d}", rn.Loc, len(tt.sql))
			}
		})
	}
}
