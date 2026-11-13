package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestAuditStmt pins the node AUDIT and NOAUDIT parse into. Whether each form
// parses is held to OceanBase by its rows in testdata/compat.tsv (rules
// audit_stmt, audit_clause and audit_all_shortcut).
func TestAuditStmt(t *testing.T) {
	for _, tc := range []struct {
		sql     string
		noAudit bool
	}{
		{"AUDIT ALL", false},
		{"AUDIT ALL STATEMENTS BY ACCESS WHENEVER NOT SUCCESSFUL", false},
		{"AUDIT SELECT TABLE, UPDATE TABLE BY u1, 'u2'@'%' BY ACCESS", false},
		{"AUDIT INSERT, DELETE ON d.t WHENEVER SUCCESSFUL", false},
		{"audit alter system", false},
		{"NOAUDIT SESSION", true},
		{"NOAUDIT SELECT ON DEFAULT", true},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.AuditStmt)
		if len(list.Items) != 1 || !ok || stmt.NoAudit != tc.noAudit {
			t.Errorf("Parse(%q) = %T, want one *ast.AuditStmt with NoAudit = %v", tc.sql, list.Items[0], tc.noAudit)
		}
	}
	for _, sql := range []string{
		"AUDIT",
		"AUDIT SELECT ON .t",            // normal_relation_factor has no leading dot
		"AUDIT ALL BY ACCESS BY ACCESS", // ACCESS names no user
		"AUDIT ALL WHENEVER SUCCESSFUL BY ACCESS",
		"AUDIT DROP",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
