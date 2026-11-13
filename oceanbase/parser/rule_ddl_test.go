package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/oceanbase/ast"
)

// TestRuleDDLStmt pins the node outline, snapshot, limiting-rule and
// sensitive-rule statements parse into. Whether each form parses is held to
// OceanBase by its rows in testdata/compat.tsv.
func TestRuleDDLStmt(t *testing.T) {
	for _, tc := range []struct {
		sql, action, schema, name string
	}{
		{"CREATE OUTLINE o ON SELECT 1", "CREATE OUTLINE", "", "o"},
		{"CREATE OR REPLACE FORMAT OUTLINE o ON SELECT * FROM t TO SELECT * FROM t", "CREATE OR REPLACE FORMAT OUTLINE", "", "o"},
		{"CREATE /*+ hint */ OUTLINE o ON 'select 1' USING HINT /*+ full(t) */", "CREATE OUTLINE", "", "o"},
		{"ALTER OUTLINE o ADD UPDATE t SET a = 1", "ALTER OUTLINE", "", "o"},
		{"DROP FORMAT OUTLINE d.o", "DROP FORMAT OUTLINE", "d", "o"},
		{"CREATE SNAPSHOT", "CREATE SNAPSHOT", "", ""},
		{"CREATE SNAPSHOT s FOR TENANT t", "CREATE SNAPSHOT", "", "s"},
		{"DROP SNAPSHOT s", "DROP SNAPSHOT", "", "s"},
		{"CREATE CONCURRENT_LIMITING_RULE IF NOT EXISTS r ON d.t TO u@'%' FOR SELECT FILTER BY KEYWORD('a', 'b') WITH MAX_CONCURRENCY = 3 PER SQL", "CREATE CONCURRENT_LIMITING_RULE", "", "r"},
		{"CREATE CONCURRENT_LIMITING_RULE r ON * . * TO CURRENT_USER() FOR ALL WITH MAX_CONCURRENCY 1", "CREATE CONCURRENT_LIMITING_RULE", "", "r"},
		{"DROP CONCURRENT_LIMITING_RULE IF EXISTS d.r", "DROP CONCURRENT_LIMITING_RULE", "d", "r"},
		{"CREATE SENSITIVE RULE r ON t(a, b), d.u(c) USING ENCRYPTION = _utf8mb4 'k'", "CREATE SENSITIVE RULE", "", "r"},
		{"ALTER SENSITIVE RULE r DROP COLUMN t(a)", "ALTER SENSITIVE RULE", "", "r"},
		{"DROP SENSITIVE RULE r", "DROP SENSITIVE RULE", "", "r"},
	} {
		list, err := Parse(tc.sql)
		if err != nil {
			t.Errorf("Parse(%q): %v", tc.sql, err)
			continue
		}
		stmt, ok := list.Items[0].(*nodes.ObjectDDLStmt)
		if len(list.Items) != 1 || !ok {
			t.Errorf("Parse(%q) = %T, want one *ast.ObjectDDLStmt", tc.sql, list.Items[0])
			continue
		}
		if stmt.Action != tc.action || stmt.Schema != tc.schema || stmt.Name != tc.name {
			t.Errorf("Parse(%q) = %q %q.%q, want %q %q.%q", tc.sql, stmt.Action, stmt.Schema, stmt.Name, tc.action, tc.schema, tc.name)
		}
		// A hint is a comment, so the statement ends at HINT, its last token.
		if end := len(strings.TrimSuffix(tc.sql, " /*+ full(t) */")); stmt.Loc.Start != 0 || stmt.Loc.End != end {
			t.Errorf("Parse(%q).Loc = %+v, want 0..%d", tc.sql, stmt.Loc, end)
		}
	}
	for _, sql := range []string{
		"CREATE OUTLINE o ON 'select 1' USING HINT",         // no hint after HINT
		"CREATE OUTLINE o ON 'select 1' USING HINT /* x */", // a comment, not a hint
		"CREATE OR REPLACE SNAPSHOT s",
		"CREATE CONCURRENT_LIMITING_RULE r ON d . t TO u FOR ALL WITH MAX_CONCURRENCY 1", // spaced, d . t is no ID_DOT_ID
		"CREATE SENSITIVE RULE r ON t(a) USING ENCRYPTION = 1",
		"ALTER OUTLINE o ADD CREATE TABLE t (a int)",
	} {
		if _, err := Parse(sql); err == nil {
			t.Errorf("Parse(%q) was accepted", sql)
		}
	}
}
