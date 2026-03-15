package parser

import (
	"testing"
)

// TestParseSessionStmts tests parsing of SET ROLE and SET CONSTRAINT(S) statements.
func TestParseSessionStmts(t *testing.T) {
	tests := []string{
		// ---- SET ROLE ----
		// Single role
		`SET ROLE dba`,
		// Multiple roles
		`SET ROLE dba, resource, connect`,
		// Role with IDENTIFIED BY
		`SET ROLE secure_role IDENTIFIED BY secret`,
		// Multiple roles with IDENTIFIED BY
		`SET ROLE role1 IDENTIFIED BY pass1, role2, role3 IDENTIFIED BY pass3`,
		// ALL
		`SET ROLE ALL`,
		// ALL EXCEPT
		`SET ROLE ALL EXCEPT dba`,
		`SET ROLE ALL EXCEPT dba, resource`,
		// NONE
		`SET ROLE NONE`,

		// ---- SET CONSTRAINT(S) ----
		// ALL IMMEDIATE
		`SET CONSTRAINTS ALL IMMEDIATE`,
		// ALL DEFERRED
		`SET CONSTRAINTS ALL DEFERRED`,
		// Specific constraint IMMEDIATE
		`SET CONSTRAINT pk_emp IMMEDIATE`,
		// Specific constraints DEFERRED
		`SET CONSTRAINTS fk_dept, fk_mgr DEFERRED`,
		// Single form
		`SET CONSTRAINT ALL IMMEDIATE`,
		// Multiple constraints IMMEDIATE
		`SET CONSTRAINTS pk_emp, fk_dept, uq_email IMMEDIATE`,
	}
	for _, sql := range tests {
		name := sql
		if len(name) > 60 {
			name = name[:60]
		}
		t.Run(name, func(t *testing.T) {
			result := ParseAndCheck(t, sql)
			if result.Len() < 1 {
				t.Fatalf("expected at least 1 statement, got %d", result.Len())
			}
		})
	}
}
