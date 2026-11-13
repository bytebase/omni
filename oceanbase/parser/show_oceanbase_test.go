package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestOceanBaseShowForms pins the ShowStmt OceanBase's own SHOW forms parse
// into. Whether each form parses is held to OceanBase by the show_stmt rows of
// testdata/compat.tsv.
func TestOceanBaseShowForms(t *testing.T) {
	for _, tc := range []struct {
		sql, typ, from, like string
	}{
		{"SHOW PARAMETERS LIKE 'a%' TENANT = t1", "PARAMETERS", "t1", "a%"},
		{"SHOW PARAMETERS TENANT 'x'", "PARAMETERS", "x", ""},
		{"SHOW ALL PARAMETERS", "ALL PARAMETERS", "", ""},
		{"SHOW SEQUENCES LIKE 's' 'q' IN d", "SEQUENCES", "d", "sq"},
		{"SHOW TABLEGROUPS", "TABLEGROUPS", "", ""},
		{"SHOW TENANT STATUS", "TENANT STATUS", "", ""},
		{"SHOW CREATE TENANT t", "CREATE TENANT", "t", ""},
		{"SHOW CREATE /*+ x */ TABLEGROUP g", "CREATE TABLEGROUP", "g", ""},
		{"SHOW SENSITIVE RULES FROM t IN d", "SENSITIVE RULES", "t", ""},
		{"SHOW SENSITIVE RULES IN DATABASE d LIKE 'r'", "SENSITIVE RULES", "d", "r"},
		{"SHOW TRACE FORMAT = 'json'", "TRACE", "", ""},
		{"SHOW JOB STATUS WHERE JOB = 'j'", "JOB STATUS", "", ""},
		{"SHOW DATABASES STATUS LIKE 'd'", "DATABASES STATUS", "", "d"},
		{"SHOW LOCAL VARIABLES", "SESSION VARIABLES", "", ""},
		{"SHOW TABLES LIKE 'a!%' ESCAPE '!'", "TABLES", "", "a!%"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.ShowStmt)
		if len(list.Items) != 1 || !ok {
			t.Errorf("Parse(%q) = %T, want one *ast.ShowStmt", tc.sql, list.Items[0])
			continue
		}
		from, like := "", ""
		if stmt.From != nil {
			from = stmt.From.Name
		}
		if lit, ok := stmt.Like.(*nodes.StringLit); ok {
			like = lit.Value
		}
		if stmt.Type != tc.typ || from != tc.from || like != tc.like {
			t.Errorf("Parse(%q) = %q from %q like %q, want %q from %q like %q", tc.sql, stmt.Type, from, like, tc.typ, tc.from, tc.like)
		}
	}
	for _, sql := range []string{
		"SHOW ALL PARAMETERS TENANT t", // only SHOW PARAMETERS names a tenant
		"SHOW TABLES LIKE 'a' ESCAPE 'xy'",
		"SHOW JOB STATUS WHERE 1",
		"SHOW TRACE FORMAT 'json'",
		"SHOW GRANTS FOR CURRENT_USER(",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
