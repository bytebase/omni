package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/v1review"
)

func TestStatementDisallowTruncate(t *testing.T) {
	truncate := func(table, text string) stmtFinding {
		return stmtFinding{Code: v1review.Truncate, Table: table, Text: text}
	}
	testStatementRule(t, StatementDisallowTruncate, []statementCase{
		{name: "one table", sql: "TRUNCATE TABLE t;", want: []stmtFinding{truncate("t", "t")}},
		{name: "each table anchors itself", sql: "TRUNCATE s.a,\n  b;", want: []stmtFinding{truncate("s.a", "s.a"), truncate("b", "b")}},
		{name: "ONLY is part of the reference", sql: "TRUNCATE ONLY t;", want: []stmtFinding{truncate("t", "ONLY t")}},
		{name: "delete not judged", sql: "DELETE FROM t;"},
		{name: "function body not judged", sql: "CREATE FUNCTION f() RETURNS void AS $$ TRUNCATE t; $$ LANGUAGE sql;"},
	})
}

func TestStatementDisallowTruncateOptions(t *testing.T) {
	for sql, want := range map[string]string{
		"TRUNCATE t;":                            "",
		"TRUNCATE t CASCADE;":                    "CASCADE",
		"TRUNCATE t RESTART IDENTITY;":           "RESTART IDENTITY",
		"TRUNCATE t RESTART IDENTITY CASCADE;":   "CASCADE,RESTART IDENTITY",
		"TRUNCATE t CONTINUE IDENTITY RESTRICT;": "",
	} {
		findings := StatementDisallowTruncate(parseStatements(t, sql))
		var got []string
		for _, f := range findings {
			got = append(got, f.Value)
		}
		if diff := cmp.Diff([]string{want}, got); diff != "" {
			t.Errorf("%s: Detail (-want +got):\n%s", sql, diff)
		}
	}
}
