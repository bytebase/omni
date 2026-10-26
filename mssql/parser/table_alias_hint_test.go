package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// oldStyleTableHints are the hints SQL Server 2022 takes in parentheses
// without WITH after a table alias, one at a time.
var oldStyleTableHints = []string{
	"NOLOCK", "HOLDLOCK", "NOWAIT", "FORCESCAN", "FORCESEEK", "NOEXPAND",
	"PAGLOCK", "READCOMMITTED", "READCOMMITTEDLOCK", "READPAST",
	"READUNCOMMITTED", "REPEATABLEREAD", "ROWLOCK", "SERIALIZABLE", "SNAPSHOT",
	"TABLOCK", "TABLOCKX", "UPDLOCK", "XLOCK", "KEEPIDENTITY", "KEEPDEFAULTS",
	"IGNORE_CONSTRAINTS", "IGNORE_TRIGGERS",
}

// tableAliasHintCases covers WINDOW as a bare table alias and table hints in
// parentheses without WITH after an alias. WINDOW starts a WINDOW clause only
// as WINDOW name AS; anywhere else after a table source it is the alias. An
// alias may be followed by one hint in parentheses; more than one, INDEX, or
// a hint with arguments needs WITH. shape is the alias and hints of the first
// table source of an accepted statement, printed by tableAliasHintShape, or ""
// where that is a DML target.
var tableAliasHintCases = func() []struct {
	name    string
	sql     string
	accepts bool
	shape   string
} {
	cases := []struct {
		name    string
		sql     string
		accepts bool
		shape   string
	}{
		// WINDOW as an alias.
		{"window_alias_end", "SELECT * FROM t window", true, "t AS window"},
		{"window_alias_where", "SELECT * FROM t window WHERE 1 = 1", true, "t AS window"},
		{"window_alias_join", "SELECT * FROM t window JOIN t1 ON window.a = t1.a", true, "t AS window"},
		{"window_alias_cross_join", "SELECT * FROM t window CROSS JOIN t1", true, "t AS window"},
		{"window_alias_comma", "SELECT * FROM t window, t1", true, "t AS window"},
		{"window_alias_group_by", "SELECT * FROM t window GROUP BY a", true, "t AS window"},
		{"window_alias_order_by", "SELECT * FROM t window ORDER BY 1", true, "t AS window"},
		{"window_alias_union", "SELECT * FROM t window UNION SELECT * FROM t1", true, "t AS window"},
		{"window_alias_option", "SELECT * FROM t window OPTION (RECOMPILE)", true, "t AS window"},
		{"window_alias_for_xml", "SELECT * FROM t window FOR XML AUTO", true, "t AS window"},
		{"window_alias_pivot", "SELECT * FROM t window PIVOT (MAX(col) FOR a IN ([1])) p", true, "t AS window"},
		{"window_alias_tablesample", "SELECT * FROM t window TABLESAMPLE (10 PERCENT)", true, "t AS window"},
		{"window_alias_hint", "SELECT * FROM t window WITH (NOLOCK)", true, "t AS window WITH NOLOCK"},
		{"window_alias_then_clause", "SELECT a FROM t window WINDOW w AS (PARTITION BY a)", true, "t AS window"},
		{"window_alias_then_clause_named_window", "SELECT a FROM t window window w AS (PARTITION BY a)", true, "t AS window"},
		{"window_alias_join_right_then_clause", "SELECT * FROM t1 JOIN t window ON 1 = 1 WINDOW w AS (PARTITION BY a)", true, "t1"},
		{"window_alias_derived_table", "SELECT * FROM (SELECT 1 AS a) window WHERE 1 = 1", true, "derived AS window"},
		{"window_alias_table_variable", "DECLARE @tv TABLE (a INT); SELECT * FROM @tv window WHERE 1 = 1", true, "@tv AS window"},
		{"window_alias_rowset", "SELECT * FROM OPENJSON(N'[]') window WHERE 1 = 1", true, "derived AS window"},
		{"window_alias_join_right", "SELECT * FROM t1 JOIN t window ON window.a = t1.a", true, "t1"},
		{"window_alias_update_from", "UPDATE window SET a = 1 FROM t window", true, ""},
		{"window_clause", "SELECT a FROM t WINDOW w AS (PARTITION BY a)", true, "t"},
		{"window_clause_named_window", "SELECT a FROM t window window AS (PARTITION BY a)", true, "t"},
		{"window_clause_after_where", "SELECT a FROM t WHERE 1 = 1 WINDOW w AS (PARTITION BY a)", true, "t"},
		{"window_as_alias", "SELECT a FROM t AS window WINDOW w AS (PARTITION BY a)", true, "t AS window"},
		{"window_alias_then_name", "SELECT a FROM t window w", false, ""},
		{"window_alias_then_as", "SELECT a FROM t window AS (PARTITION BY a)", false, ""},

		// A table hint in parentheses after an alias.
		{"hint_after_as_alias", "SELECT * FROM t AS s (NOLOCK)", true, "t AS s WITH NOLOCK"},
		{"hint_lower_case", "SELECT * FROM t s (nolock)", true, "t AS s WITH NOLOCK"},
		{"hint_no_space", "SELECT * FROM t s(NOLOCK)", true, "t AS s WITH NOLOCK"},
		{"hint_schema_table", "SELECT * FROM dbo.t s (NOLOCK)", true, "t AS s WITH NOLOCK"},
		{"hint_window_alias", "SELECT * FROM t window (NOLOCK)", true, "t AS window WITH NOLOCK"},
		{"hint_then_join", "SELECT * FROM t s (NOLOCK) JOIN t1 ON s.a = t1.a", true, "t AS s WITH NOLOCK"},
		{"hint_then_where", "SELECT * FROM t s (NOLOCK) WHERE 1 = 1", true, "t AS s WITH NOLOCK"},
		{"hint_then_pivot", "SELECT * FROM t s (READPAST) PIVOT (MAX(col) FOR a IN ([1])) p", true, "t AS s WITH READPAST"},
		{"hint_after_tablesample", "SELECT * FROM t s TABLESAMPLE (10 PERCENT) (NOLOCK)", true, "t AS s WITH NOLOCK"},
		{"hint_after_for_system_time", "SELECT * FROM t FOR SYSTEM_TIME ALL s (NOLOCK)", true, "t AS s WITH NOLOCK"},
		{"hint_update_from", "UPDATE s SET a = 1 FROM t s (NOLOCK)", true, ""},
		{"hint_delete_from", "DELETE s FROM t s (NOLOCK)", true, ""},
		{"hint_merge_using", "MERGE INTO dbo.t USING t1 AS s (NOLOCK) ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, ""},
		{"hint_without_alias", "SELECT * FROM t (NOLOCK)", true, ""},
		{"query_after_alias", "SELECT * FROM t s (SELECT 1)", false, ""},
		{"hint_two", "SELECT * FROM t s (NOLOCK, READPAST)", false, ""},
		{"hint_two_without_comma", "SELECT * FROM t s (NOLOCK READPAST)", false, ""},
		{"hint_twice", "SELECT * FROM t s (NOLOCK) (READPAST)", false, ""},
		{"hint_index", "SELECT * FROM t s (INDEX(ix))", false, ""},
		{"hint_index_equals", "SELECT * FROM t s (INDEX = ix)", false, ""},
		{"hint_index_number", "SELECT * FROM t s (1)", false, ""},
		{"hint_forceseek_arguments", "SELECT * FROM t s (FORCESEEK(ix(a)))", false, ""},
		{"hint_spatial", "SELECT * FROM t s (SPATIAL_WINDOW_MAX_CELLS = 10)", false, ""},
		{"hint_unknown", "SELECT * FROM t s (FOO)", false, ""},
		{"hint_bracketed", "SELECT * FROM t s ([NOLOCK])", false, ""},
		{"hint_empty", "SELECT * FROM t s ()", false, ""},
		{"hint_unclosed", "SELECT * FROM t s (NOLOCK", false, ""},
		{"hint_then_with", "SELECT * FROM t s (NOLOCK) WITH (READPAST)", false, ""},
		{"with_then_hint", "SELECT * FROM t s WITH (NOLOCK) (READPAST)", false, ""},
		{"hint_then_alias", "SELECT * FROM t s (NOLOCK) AS x", false, ""},
		{"hint_then_tablesample", "SELECT * FROM t s (NOLOCK) TABLESAMPLE (10 PERCENT)", false, ""},
		{"hint_unparenthesized", "SELECT * FROM t s HOLDLOCK", false, ""},
		{"hint_bracketed_unparenthesized", "SELECT * FROM t s [NOLOCK]", false, ""},
		{"hint_update_target", "UPDATE t (NOLOCK) SET a = 1", false, ""},
	}
	for _, hint := range oldStyleTableHints {
		cases = append(cases, struct {
			name    string
			sql     string
			accepts bool
			shape   string
		}{
			"hint_" + strings.ToLower(hint),
			"SELECT * FROM t s (" + hint + ")",
			true,
			"t AS s WITH " + hint,
		})
	}
	return cases
}()

// tableAliasHintShape prints the name, alias, and hints of the first table
// source in pre-order.
func tableAliasHintShape(root ast.Node) string {
	var shape string
	ast.Inspect(root, func(n ast.Node) bool {
		if shape != "" {
			return false
		}
		var name, alias string
		var hints *ast.List
		switch v := n.(type) {
		case *ast.AliasedTableRef:
			name, alias, hints = "derived", v.Alias, v.Hints
			if ref, ok := v.Table.(*ast.TableRef); ok {
				name = ref.Object
				if hints == nil {
					hints = ref.Hints
				}
			}
		case *ast.TableRef:
			name, alias, hints = v.Object, v.Alias, v.Hints
		case *ast.TableVarRef:
			name, alias = v.Name, v.Alias
		default:
			return true
		}
		shape = name
		if alias != "" {
			shape += " AS " + alias
		}
		if hints != nil {
			var names []string
			for _, h := range hints.Items {
				names = append(names, h.(*ast.TableHint).Name)
			}
			shape += " WITH " + strings.Join(names, ", ")
		}
		return false
	})
	return shape
}

// TestTableAliasHint checks that omni accepts the aliases and old-style table
// hints SQL Server accepts, rejects the others, and keeps every Loc inside its
// parent.
func TestTableAliasHint(t *testing.T) {
	for _, tc := range tableAliasHintCases {
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
			if tc.shape != "" {
				if got := tableAliasHintShape(list); got != tc.shape {
					t.Errorf("shape = %q, want %q", got, tc.shape)
				}
			}
			assertChildLocsWithinParent(t, tc.sql, list)
		})
	}
}

// TestTableAliasHintOracle checks accept/reject parity with SQL Server for
// every case in tableAliasHintCases.
func TestTableAliasHintOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range tableAliasHintCases {
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

// TestTableAliasHintLoc checks that a table spans its alias and a hint in
// parentheses after it.
func TestTableAliasHintLoc(t *testing.T) {
	tests := []struct {
		sql  string
		node string
		want string
	}{
		{"SELECT * FROM dbo.t AS s (NOLOCK) WHERE 1 = 1", "TableRef", "dbo.t AS s (NOLOCK)"},
		{"SELECT * FROM dbo.t AS s (NOLOCK) WHERE 1 = 1", "TableHint", "NOLOCK"},
		{"SELECT * FROM t s TABLESAMPLE (10 PERCENT) (READPAST) WHERE 1 = 1", "AliasedTableRef", "t s TABLESAMPLE (10 PERCENT) (READPAST)"},
		{"SELECT * FROM t window WHERE 1 = 1", "TableRef", "t window"},
	}
	for _, tc := range tests {
		t.Run(tc.node+"/"+tc.want, func(t *testing.T) {
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
