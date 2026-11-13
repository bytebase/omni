package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/pg/ast"
	"github.com/bytebase/omni/pg/parser"
	"github.com/bytebase/omni/v1review"
)

// systemFinding is a finding as the system rule tests compare it: its code and
// name, the statement it is in, and the text its range anchors.
type systemFinding struct {
	Code      v1review.Code
	Name      string
	Statement int
	Text      string
}

// judgeSystem parses sql and returns what rule finds in its statements.
func judgeSystem(t *testing.T, sql string, rule func([]ast.Node) []v1review.Finding) []systemFinding {
	t.Helper()
	list, err := parser.Parse(sql)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	var got []systemFinding
	for _, f := range rule(list.Items) {
		got = append(got, systemFinding{f.Code, f.Name, f.Statement, sql[f.Start:f.End]})
	}
	return got
}

func TestSystemCharsetAllowlist(t *testing.T) {
	bad := func(name, text string) systemFinding {
		return systemFinding{Code: v1review.CharsetNotAllowed, Name: name, Text: text}
	}
	allowed := []string{"UTF8", "sql_ascii"}
	for _, tt := range []struct {
		name    string
		sql     string
		allowed []string
		want    []systemFinding
	}{
		{name: "allowed in any case", sql: "CREATE DATABASE d WITH ENCODING 'utf8'; CREATE DATABASE e ENCODING = 'SQL_ASCII';", allowed: allowed},
		{name: "no encoding", sql: "CREATE DATABASE d;", allowed: allowed},
		{
			name:    "disallowed is lower-cased",
			sql:     "CREATE DATABASE d WITH ENCODING 'LATIN1';",
			allowed: allowed,
			want:    []systemFinding{bad("latin1", "CREATE DATABASE d WITH ENCODING 'LATIN1'")},
		},
		{
			name:    "empty list allows nothing",
			sql:     "CREATE DATABASE d ENCODING 'UTF8';",
			allowed: nil,
			want:    []systemFinding{bad("utf8", "CREATE DATABASE d ENCODING 'UTF8'")},
		},
		{
			name:    "only the first encoding is judged",
			sql:     "CREATE DATABASE d ENCODING 'LATIN1' ENCODING 'UTF8';",
			allowed: allowed,
			want:    []systemFinding{bad("latin1", "CREATE DATABASE d ENCODING 'LATIN1' ENCODING 'UTF8'")},
		},
		{name: "a numeric encoding is skipped", sql: "CREATE DATABASE d ENCODING 8;", allowed: allowed},
		{name: "a default encoding is skipped", sql: "CREATE DATABASE d ENCODING DEFAULT;", allowed: allowed},
		{name: "an empty encoding is skipped", sql: "CREATE DATABASE d ENCODING '';", allowed: allowed},
		{
			name:    "other options are not read",
			sql:     "CREATE DATABASE d LOCALE 'C' LC_COLLATE 'latin1' TEMPLATE template0 ENCODING 'win1252';",
			allowed: allowed,
			want:    []systemFinding{bad("win1252", "CREATE DATABASE d LOCALE 'C' LC_COLLATE 'latin1' TEMPLATE template0 ENCODING 'win1252'")},
		},
		{name: "alter database is not judged", sql: "ALTER DATABASE d SET client_encoding = 'LATIN1';", allowed: allowed},
		{name: "set client encoding is not judged", sql: "SET client_encoding = 'LATIN1'; SET NAMES 'LATIN1';", allowed: allowed},
		{name: "a routine body is not judged", sql: "CREATE FUNCTION f() RETURNS void LANGUAGE sql AS $$ SELECT 1 $$;", allowed: allowed},
		{
			name:    "each statement in order",
			sql:     "CREATE DATABASE a ENCODING 'EUC_JP';\nCREATE DATABASE b ENCODING 'UTF8';\nCREATE DATABASE c ENCODING 'Big5';",
			allowed: allowed,
			want: []systemFinding{
				bad("euc_jp", "CREATE DATABASE a ENCODING 'EUC_JP'"),
				bad("big5", "CREATE DATABASE c ENCODING 'Big5'"),
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			got := judgeSystem(t, tt.sql, func(stmts []ast.Node) []v1review.Finding {
				return SystemCharsetAllowlist(stmts, tt.allowed)
			})
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("findings mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
