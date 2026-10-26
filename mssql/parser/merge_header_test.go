package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// mergeHeaderCases covers the head of a MERGE statement: the target and its
// alias, USING, the source and its alias, ON, and the WHEN clauses. USING,
// ON, and at least one WHEN clause are required; the target alias is not an
// undelimited USING; and the source takes one alias, its own. Every
// statement ends with a semicolon, which SQL Server requires after MERGE.
// shape is the target, source, and WHEN clause count of an accepted
// statement, printed by mergeHeaderShape.
var mergeHeaderCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	{"into", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t USING t1 WHEN 1"},
	{"no_into", "MERGE dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t USING t1 WHEN 1"},
	{"top", "MERGE TOP (1) INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t USING t1 WHEN 1"},
	{"using_lower_case", "MERGE INTO dbo.t using t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t USING t1 WHEN 1"},
	{"target_as_alias", "MERGE INTO dbo.t AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", true, "t AS x USING t1 WHEN 1"},
	{"target_bare_alias", "MERGE INTO dbo.t x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", true, "t AS x USING t1 WHEN 1"},
	{"target_bracketed_using_alias", "MERGE INTO dbo.t [USING] USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t AS USING USING t1 WHEN 1"},
	{"target_window_alias", "MERGE INTO dbo.t window USING t1 ON window.a = t1.a WHEN MATCHED THEN DELETE;", true, "t AS window USING t1 WHEN 1"},
	{"target_named_window", "MERGE INTO window USING t1 ON window.a = t1.a WHEN MATCHED THEN DELETE;", true, "window USING t1 WHEN 1"},
	{"target_hint_then_alias", "MERGE INTO dbo.t WITH (HOLDLOCK) AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", true, "t AS x USING t1 WHEN 1"},
	{"target_variable_as_alias", "DECLARE @tv TABLE (a INT); MERGE INTO @tv AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", true, "@tv AS x USING t1 WHEN 1"},
	{"target_variable_bare_alias", "DECLARE @tv TABLE (a INT); MERGE INTO @tv x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", true, "@tv AS x USING t1 WHEN 1"},
	{"target_variable", "DECLARE @tv TABLE (a INT); MERGE INTO @tv USING t1 ON a = t1.a WHEN MATCHED THEN DELETE;", true, "@tv USING t1 WHEN 1"},
	{"source_as_alias", "MERGE INTO dbo.t USING t1 AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t USING t1 AS s WHEN 1"},
	{"source_bare_alias", "MERGE INTO dbo.t USING t1 s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t USING t1 AS s WHEN 1"},
	{"source_hint_after_alias", "MERGE INTO dbo.t USING t1 s (NOLOCK) ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t USING t1 AS s WHEN 1"},
	{"source_derived_table", "MERGE INTO dbo.t USING (SELECT 1 AS a) s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t USING derived AS s WHEN 1"},
	{"source_derived_table_as", "MERGE INTO dbo.t USING (SELECT 1 AS a) AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t USING derived AS s WHEN 1"},
	{"source_values", "MERGE INTO dbo.t USING (VALUES (1)) v(a) ON dbo.t.a = v.a WHEN MATCHED THEN DELETE;", true, "t USING derived AS v WHEN 1"},
	{"source_join", "MERGE INTO dbo.t USING t1 JOIN t2 ON t1.a = t2.a ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", true, "t USING join WHEN 1"},
	{"when_three", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE WHEN NOT MATCHED THEN INSERT (a) VALUES (t1.a) WHEN NOT MATCHED BY SOURCE THEN DELETE;", true, "t USING t1 WHEN 3"},
	{"output_option", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE OUTPUT $action OPTION (RECOMPILE);", true, "t USING t1 WHEN 1"},
	{"target_as_using", "MERGE INTO dbo.t AS USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
	{"target_as_using_using", "MERGE INTO dbo.t AS USING USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
	{"target_alias_after_hint_order", "MERGE INTO dbo.t AS x WITH (HOLDLOCK) USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
	{"missing_using", "MERGE INTO dbo.t t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
	{"missing_using_after_alias", "MERGE INTO dbo.t x t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
	{"using_without_source", "MERGE INTO dbo.t USING ON dbo.t.a = 1 WHEN MATCHED THEN DELETE;", false, ""},
	{"missing_on", "MERGE INTO dbo.t USING t1 WHEN MATCHED THEN DELETE;", false, ""},
	{"on_without_condition", "MERGE INTO dbo.t USING t1 ON WHEN MATCHED THEN DELETE;", false, ""},
	{"missing_when", "MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a;", false, ""},
	{"source_two_aliases", "MERGE INTO dbo.t USING t1 a b ON dbo.t.a = a.a WHEN MATCHED THEN DELETE;", false, ""},
	{"source_two_as_aliases", "MERGE INTO dbo.t USING t1 AS a AS b ON dbo.t.a = a.a WHEN MATCHED THEN DELETE;", false, ""},
	{"source_derived_table_two_aliases", "MERGE INTO dbo.t USING (SELECT 1 AS a) s x ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", false, ""},
}

// mergeTableText prints a MERGE target or source.
func mergeTableText(n ast.Node) string {
	switch v := n.(type) {
	case *ast.TableRef:
		if v.Alias != "" {
			return v.Object + " AS " + v.Alias
		}
		return v.Object
	case *ast.TableVarRef:
		if v.Alias != "" {
			return v.Name + " AS " + v.Alias
		}
		return v.Name
	case *ast.AliasedTableRef:
		return "derived AS " + v.Alias
	case *ast.JoinClause:
		return "join"
	}
	return fmt.Sprintf("%T", n)
}

// mergeHeaderShape prints the target, source, and WHEN clause count of the
// first MERGE statement.
func mergeHeaderShape(stmt *ast.MergeStmt) string {
	return fmt.Sprintf("%s USING %s WHEN %d", mergeTableText(stmt.Target), mergeTableText(stmt.Source), len(stmt.WhenClauses.Items))
}

// TestMergeHeader checks that omni accepts the MERGE heads SQL Server
// accepts, rejects the others, and keeps every Loc inside its parent.
func TestMergeHeader(t *testing.T) {
	for _, tc := range mergeHeaderCases {
		t.Run(tc.name, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error:\n%s", tc.sql, ast.NodeToString(list))
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			stmt, ok := nthNodeByType(list, "MergeStmt", 0).(*ast.MergeStmt)
			if !ok {
				t.Fatalf("no MergeStmt in %s", ast.NodeToString(list))
			}
			if got := mergeHeaderShape(stmt); got != tc.shape {
				t.Errorf("shape = %q, want %q", got, tc.shape)
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestMergeHeaderOracle checks accept/reject parity with SQL Server for every
// case in mergeHeaderCases.
func TestMergeHeaderOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range mergeHeaderCases {
		t.Run(tc.name, func(t *testing.T) {
			ssErr, err := oracle.parseError(tc.sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, tc.sql)
			}

			_, omniErr := Parse(tc.sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, tc.sql)
			}
		})
	}
}

// TestMergeHeaderErrorPosition checks where a missing part of the MERGE head
// is reported.
func TestMergeHeaderErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"MERGE INTO dbo.t AS USING t1 ON dbo.t.a = t1.a WHEN MATCHED THEN DELETE;", `syntax error at or near "USING"`},
		{"MERGE INTO dbo.t x t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", `syntax error at or near "t1"`},
		{"MERGE INTO dbo.t USING t1 WHEN MATCHED THEN DELETE;", `syntax error at or near "WHEN"`},
		{"MERGE INTO dbo.t USING t1 ON dbo.t.a = t1.a;", `syntax error at or near ";"`},
		{"MERGE INTO dbo.t USING t1 a b ON dbo.t.a = a.a WHEN MATCHED THEN DELETE;", `syntax error at or near "b"`},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			if err == nil {
				t.Fatalf("Parse(%q) succeeded, want %s", tc.sql, tc.want)
			}
			if got := err.Error(); got != tc.want {
				t.Errorf("Parse(%q) error = %q, want %q", tc.sql, got, tc.want)
			}
		})
	}
}

// TestMergeHeaderCompletion checks that USING is offered after a MERGE target.
func TestMergeHeaderCompletion(t *testing.T) {
	for _, sql := range []string{"MERGE INTO dbo.t ", "MERGE INTO dbo.t AS x "} {
		t.Run(sql+"|", func(t *testing.T) {
			cs := Collect(sql, len(sql))
			if cs == nil || !cs.HasToken(kwUSING) {
				t.Errorf("missing token candidate USING")
			}
		})
	}
}

// TestMergeHeaderLoc checks that a MERGE target spans its hints and alias.
func TestMergeHeaderLoc(t *testing.T) {
	tests := []struct {
		sql  string
		node string
		want string
	}{
		{"MERGE INTO dbo.t WITH (HOLDLOCK) AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", "TableRef", "dbo.t WITH (HOLDLOCK) AS x"},
		{"DECLARE @tv TABLE (a INT); MERGE INTO @tv AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", "TableVarRef", "@tv AS x"},
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			n := nthNodeByType(list, tc.node, 0)
			if n == nil {
				t.Fatalf("no %s in %s", tc.node, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s text = %q, want %q", tc.node, got, tc.want)
			}
		})
	}
}
