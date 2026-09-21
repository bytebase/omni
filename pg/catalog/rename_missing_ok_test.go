package catalog

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestRenameMissingOk checks that IF EXISTS suppresses the undefined-relation
// error for every relation-based RENAME form, matching PG's
// "relation does not exist, skipping" NOTICE. Verified against PG 17:
//
//	ALTER TABLE IF EXISTS missing RENAME TO x;                        -- skipped
//	ALTER VIEW IF EXISTS missing RENAME c TO d;                       -- skipped
//	ALTER MATERIALIZED VIEW IF EXISTS missing RENAME COLUMN c TO d;   -- skipped
//	ALTER VIEW IF EXISTS nosuch.missing RENAME c TO d;                -- skipped
//	ALTER VIEW missing RENAME c TO d;                                 -- ERROR
func TestRenameMissingOk(t *testing.T) {
	stmts := []struct {
		name string
		stmt *nodes.RenameStmt
	}{
		{"table rename to", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_TABLE,
			Relation:   &nodes.RangeVar{Relname: "missing"},
			Newname:    "x",
		}},
		{"table rename column", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLUMN, RelationType: nodes.OBJECT_TABLE,
			Relation: &nodes.RangeVar{Relname: "missing"},
			Subname:  "c", Newname: "d",
		}},
		{"view rename to", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_VIEW,
			Relation:   &nodes.RangeVar{Relname: "missing"},
			Newname:    "x",
		}},
		{"view rename column", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLUMN, RelationType: nodes.OBJECT_VIEW,
			Relation: &nodes.RangeVar{Relname: "missing"},
			Subname:  "c", Newname: "d",
		}},
		{"matview rename to", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_MATVIEW,
			Relation:   &nodes.RangeVar{Relname: "missing"},
			Newname:    "x",
		}},
		{"matview rename column", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLUMN, RelationType: nodes.OBJECT_MATVIEW,
			Relation: &nodes.RangeVar{Relname: "missing"},
			Subname:  "c", Newname: "d",
		}},
		{"view rename column in missing schema", &nodes.RenameStmt{
			RenameType: nodes.OBJECT_COLUMN, RelationType: nodes.OBJECT_VIEW,
			Relation: &nodes.RangeVar{Schemaname: "nosuch", Relname: "missing"},
			Subname:  "c", Newname: "d",
		}},
	}

	for _, tc := range stmts {
		t.Run(tc.name, func(t *testing.T) {
			c := New()

			// Without IF EXISTS the lookup error must surface.
			tc.stmt.MissingOk = false
			if err := c.ExecRenameStmt(tc.stmt); err == nil {
				t.Fatal("expected undefined relation error without IF EXISTS")
			}

			// With IF EXISTS the statement is a no-op.
			tc.stmt.MissingOk = true
			if err := c.ExecRenameStmt(tc.stmt); err != nil {
				t.Fatalf("expected IF EXISTS to suppress the error, got: %v", err)
			}
		})
	}
}
