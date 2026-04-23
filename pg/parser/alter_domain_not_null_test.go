package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestParseAlterDomainAddNotNullConstraint(t *testing.T) {
	tests := []struct {
		name    string
		sql     string
		conname string
	}{
		{
			name: "unnamed",
			sql:  `ALTER DOMAIN d ADD NOT NULL`,
		},
		{
			name:    "named",
			sql:     `ALTER DOMAIN public.d ADD CONSTRAINT d_not_null NOT NULL`,
			conname: "d_not_null",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stmts, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}

			stmt := stmts.Items[0].(*nodes.RawStmt).Stmt.(*nodes.AlterDomainStmt)
			if stmt.Subtype != 'C' {
				t.Fatalf("expected ADD CONSTRAINT subtype 'C', got %q", stmt.Subtype)
			}

			constraint := stmt.Def.(*nodes.Constraint)
			if constraint.Contype != nodes.CONSTR_NOTNULL {
				t.Fatalf("expected CONSTR_NOTNULL, got %v", constraint.Contype)
			}
			if constraint.Conname != tt.conname {
				t.Fatalf("expected constraint name %q, got %q", tt.conname, constraint.Conname)
			}
		})
	}
}
