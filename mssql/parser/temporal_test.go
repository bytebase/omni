package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// temporalDeclare declares the variables some temporal cases pass, which SQL
// Server requires even under SET PARSEONLY ON.
const temporalDeclare = "DECLARE @d DATETIME2 = SYSUTCDATETIME(), @e DATETIME2 = SYSUTCDATETIME(); "

// temporalCases covers FOR SYSTEM_TIME on a table source. It follows the name
// of a table, before the alias, TABLESAMPLE, and hints; it takes one of five
// forms; and each point in time is a literal or a variable. shape is the
// first table with FOR SYSTEM_TIME in an accepted statement, printed by
// temporalShape.
var temporalCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// The forms.
	{"as_of", "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01'", true, "t AS OF '2020-01-01'"},
	{"from_to", "SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' TO '2021-01-01'", true, "t FROM '2020-01-01' '2021-01-01'"},
	{"between", "SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' AND '2021-01-01'", true, "t BETWEEN '2020-01-01' '2021-01-01'"},
	{"contained_in", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('2020-01-01', '2021-01-01')", true, "t CONTAINED IN '2020-01-01' '2021-01-01'"},
	{"all", "SELECT * FROM t FOR SYSTEM_TIME ALL", true, "t ALL"},
	{"lower_case", "SELECT * FROM t for system_time all", true, "t ALL"},
	{"contained_lower_case", "SELECT * FROM t FOR system_time contained in ('a', 'b')", true, "t CONTAINED IN 'a' 'b'"},
	{"missing_form", "SELECT * FROM t FOR SYSTEM_TIME", false, ""},
	{"bracketed_system_time", "SELECT * FROM t FOR [SYSTEM_TIME] ALL", false, ""},
	{"bracketed_all", "SELECT * FROM t FOR SYSTEM_TIME [ALL]", false, ""},
	{"bracketed_contained", "SELECT * FROM t FOR SYSTEM_TIME [CONTAINED] IN ('a', 'b')", false, ""},
	{"as_without_of", "SELECT * FROM t FOR SYSTEM_TIME AS '2020-01-01'", false, ""},
	{"of_without_as", "SELECT * FROM t FOR SYSTEM_TIME OF '2020-01-01'", false, ""},
	{"as_of_without_time", "SELECT * FROM t FOR SYSTEM_TIME AS OF", false, ""},
	{"from_without_to", "SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01'", false, ""},
	{"from_and", "SELECT * FROM t FOR SYSTEM_TIME FROM '2020-01-01' AND '2021-01-01'", false, ""},
	{"between_to", "SELECT * FROM t FOR SYSTEM_TIME BETWEEN '2020-01-01' TO '2021-01-01'", false, ""},
	{"contained_without_in", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED ('2020-01-01', '2021-01-01')", false, ""},
	{"contained_one_time", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('2020-01-01')", false, ""},
	{"contained_three_times", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('a', 'b', 'c')", false, ""},
	{"contained_trailing_comma", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN ('a', 'b',)", false, ""},
	{"all_with_time", "SELECT * FROM t FOR SYSTEM_TIME ALL '2020'", false, ""},

	// A point in time.
	{"time_nchar", "SELECT * FROM t FOR SYSTEM_TIME AS OF N'2020-01-01'", true, "t AS OF N'2020-01-01'"},
	{"time_variable", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME AS OF @d", true, "t AS OF @d"},
	{"time_system_variable", "SELECT * FROM t FOR SYSTEM_TIME AS OF @@DBTS", true, "t AS OF @@DBTS"},
	{"time_integer", "SELECT * FROM t FOR SYSTEM_TIME AS OF 1", true, "t AS OF 1"},
	{"time_float", "SELECT * FROM t FOR SYSTEM_TIME AS OF 1.5", true, "t AS OF 1.5"},
	{"time_exponent", "SELECT * FROM t FOR SYSTEM_TIME AS OF 1e3", true, "t AS OF 1e3"},
	{"time_binary", "SELECT * FROM t FOR SYSTEM_TIME AS OF 0x01", true, "t AS OF 0x01"},
	{"time_money", "SELECT * FROM t FOR SYSTEM_TIME AS OF $5", true, "t AS OF $5"},
	{"time_null", "SELECT * FROM t FOR SYSTEM_TIME AS OF NULL", true, "t AS OF NULL"},
	{"from_to_variables", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME FROM @d TO @e", true, "t FROM @d @e"},
	{"from_to_integers", "SELECT * FROM t FOR SYSTEM_TIME FROM 1 TO 2", true, "t FROM 1 2"},
	{"between_null_variable", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME BETWEEN NULL AND @d", true, "t BETWEEN NULL @d"},
	{"contained_variables", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN (@d, @e)", true, "t CONTAINED IN @d @e"},
	{"contained_integer_null", "SELECT * FROM t FOR SYSTEM_TIME CONTAINED IN (1, NULL)", true, "t CONTAINED IN 1 NULL"},
	{"time_function", "SELECT * FROM t FOR SYSTEM_TIME AS OF GETDATE()", false, ""},
	{"time_builtin_function", "SELECT * FROM t FOR SYSTEM_TIME AS OF SYSUTCDATETIME()", false, ""},
	{"time_niladic", "SELECT * FROM t FOR SYSTEM_TIME AS OF CURRENT_TIMESTAMP", false, ""},
	{"time_cast", "SELECT * FROM t FOR SYSTEM_TIME AS OF CAST('2020-01-01' AS DATETIME2)", false, ""},
	{"time_expression", "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020' + '-01-01'", false, ""},
	{"time_negative", "SELECT * FROM t FOR SYSTEM_TIME AS OF -1", false, ""},
	{"time_positive_sign", "SELECT * FROM t FOR SYSTEM_TIME AS OF +1", false, ""},
	{"time_parenthesized", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME AS OF (@d)", false, ""},
	{"time_default", "SELECT * FROM t FOR SYSTEM_TIME AS OF DEFAULT", false, ""},
	{"time_collate", "SELECT * FROM t FOR SYSTEM_TIME AS OF 'a' COLLATE Latin1_General_CI_AS", false, ""},
	{"time_variable_member", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME AS OF @d.x", false, ""},

	// Where it goes.
	{"then_alias", "SELECT * FROM t FOR SYSTEM_TIME ALL AS s", true, "t ALL AS s"},
	{"then_bare_alias", "SELECT * FROM t FOR SYSTEM_TIME ALL s", true, "t ALL AS s"},
	{"from_to_then_alias", temporalDeclare + "SELECT * FROM t FOR SYSTEM_TIME FROM @d TO @e AS s", true, "t FROM @d @e AS s"},
	{"then_tablesample", "SELECT * FROM t FOR SYSTEM_TIME ALL TABLESAMPLE (10 PERCENT)", true, "t ALL"},
	{"then_alias_tablesample", "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' AS s TABLESAMPLE (10 PERCENT)", true, "t AS OF '2020-01-01' AS s"},
	{"then_hint", "SELECT * FROM t FOR SYSTEM_TIME ALL WITH (NOLOCK)", true, "t ALL"},
	{"then_alias_hint", "SELECT * FROM t FOR SYSTEM_TIME ALL AS s WITH (NOLOCK)", true, "t ALL AS s"},
	{"schema_qualified", "SELECT * FROM dbo.t FOR SYSTEM_TIME ALL", true, "t ALL"},
	{"four_part_name", "SELECT * FROM srv.db.dbo.t FOR SYSTEM_TIME ALL", true, "t ALL"},
	{"join_both", "SELECT * FROM t FOR SYSTEM_TIME ALL JOIN t1 FOR SYSTEM_TIME AS OF '2020-01-01' ON t.a = t1.a", true, "t ALL"},
	{"then_cross_apply", "SELECT * FROM t FOR SYSTEM_TIME ALL CROSS APPLY (SELECT 1 AS x) ca", true, "t ALL"},
	{"then_pivot", "SELECT * FROM t FOR SYSTEM_TIME ALL PIVOT (MAX(col) FOR a IN ([1])) p", true, "t ALL"},
	{"comma_both", "SELECT * FROM t FOR SYSTEM_TIME ALL, t1 FOR SYSTEM_TIME ALL", true, "t ALL"},
	{"then_where", "SELECT * FROM t FOR SYSTEM_TIME ALL WHERE 1 = 1", true, "t ALL"},
	{"then_for_xml", "SELECT * FROM t FOR SYSTEM_TIME ALL FOR XML AUTO", true, "t ALL"},
	{"as_of_then_for_xml", "SELECT * FROM t FOR SYSTEM_TIME AS OF '2020-01-01' FOR XML AUTO", true, "t AS OF '2020-01-01'"},
	{"then_order_by_for_json", "SELECT * FROM t FOR SYSTEM_TIME ALL ORDER BY 1 FOR JSON PATH", true, "t ALL"},
	{"union", "SELECT * FROM t FOR SYSTEM_TIME ALL UNION SELECT * FROM t1 FOR SYSTEM_TIME ALL", true, "t ALL"},
	{"cte", "WITH c AS (SELECT * FROM t FOR SYSTEM_TIME ALL) SELECT * FROM c", true, "t ALL"},
	{"scalar_subquery", "SELECT (SELECT COUNT(*) FROM t FOR SYSTEM_TIME ALL)", true, "t ALL"},
	{"in_subquery", "SELECT * FROM t WHERE a IN (SELECT a FROM t1 FOR SYSTEM_TIME ALL)", true, "t1 ALL"},
	{"update_from", "UPDATE s SET a = 1 FROM t FOR SYSTEM_TIME ALL AS s", true, "t ALL AS s"},
	{"delete_from", "DELETE s FROM t FOR SYSTEM_TIME ALL AS s", true, "t ALL AS s"},
	{"merge_using", "MERGE INTO dbo.t USING t1 FOR SYSTEM_TIME ALL AS s ON dbo.t.a = s.a WHEN MATCHED THEN DELETE;", true, "t1 ALL AS s"},
	{"for_xml_alone", "SELECT * FROM t FOR XML AUTO", true, ""},
	{"for_browse_alone", "SELECT * FROM t FOR BROWSE", true, ""},
	{"after_alias", "SELECT * FROM t AS s FOR SYSTEM_TIME ALL", false, ""},
	{"after_bare_alias", "SELECT * FROM t s FOR SYSTEM_TIME ALL", false, ""},
	{"after_hint", "SELECT * FROM t WITH (NOLOCK) FOR SYSTEM_TIME ALL", false, ""},
	{"twice", "SELECT * FROM t FOR SYSTEM_TIME ALL FOR SYSTEM_TIME ALL", false, ""},
	{"derived_table", "SELECT * FROM (SELECT 1 AS a) d FOR SYSTEM_TIME ALL", false, ""},
	{"derived_table_before_alias", "SELECT * FROM (SELECT 1 AS a) FOR SYSTEM_TIME ALL d", false, ""},
	{"table_variable", "DECLARE @tv TABLE (a INT); SELECT * FROM @tv FOR SYSTEM_TIME ALL", false, ""},
	{"rowset_function", "SELECT * FROM OPENJSON(N'[]') FOR SYSTEM_TIME ALL", false, ""},
	{"table_valued_function", "SELECT * FROM dbo.f() FOR SYSTEM_TIME ALL", false, ""},
	{"update_target", "UPDATE t FOR SYSTEM_TIME ALL SET a = 1", false, ""},
	{"delete_target", "DELETE FROM t FOR SYSTEM_TIME ALL WHERE a = 1", false, ""},
	{"insert_target", "INSERT INTO t FOR SYSTEM_TIME ALL (a) VALUES (1)", false, ""},
	{"merge_target", "MERGE INTO dbo.t FOR SYSTEM_TIME ALL AS x USING t1 ON x.a = t1.a WHEN MATCHED THEN DELETE;", false, ""},
}

// temporalTimeText prints a point in time of FOR SYSTEM_TIME.
func temporalTimeText(n ast.ExprNode) string {
	switch v := n.(type) {
	case *ast.Literal:
		switch v.Type {
		case ast.LitString:
			if v.IsNChar {
				return "N'" + v.Str + "'"
			}
			return "'" + v.Str + "'"
		case ast.LitNull:
			return "NULL"
		}
		return v.Str
	case *ast.VariableRef:
		return v.Name
	}
	return "?"
}

// temporalShape prints the table, form, times, and alias of the first table
// with FOR SYSTEM_TIME, or "" when there is none.
func temporalShape(root ast.Node) string {
	var ref *ast.TableRef
	ast.Inspect(root, func(n ast.Node) bool {
		if r, ok := n.(*ast.TableRef); ok && ref == nil && r.Temporal != nil {
			ref = r
		}
		return ref == nil
	})
	if ref == nil {
		return ""
	}
	parts := []string{ref.Object, ref.Temporal.Kind}
	if ref.Temporal.Start != nil {
		parts = append(parts, temporalTimeText(ref.Temporal.Start))
	}
	if ref.Temporal.End != nil {
		parts = append(parts, temporalTimeText(ref.Temporal.End))
	}
	if ref.Alias != "" {
		parts = append(parts, "AS", ref.Alias)
	}
	return strings.Join(parts, " ")
}

// TestTemporal checks that omni accepts the FOR SYSTEM_TIME forms SQL Server
// accepts, rejects the others, and keeps every Loc inside its parent.
func TestTemporal(t *testing.T) {
	for _, tc := range temporalCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := tc.sql
			list, err := Parse(sql)
			if !tc.accepts {
				if err == nil {
					t.Fatalf("Parse(%q) succeeded, want a syntax error:\n%s", sql, ast.NodeToString(list))
				}
				return
			}
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", sql, err)
			}
			if got := temporalShape(list); got != tc.shape {
				t.Errorf("shape = %q, want %q", got, tc.shape)
			}
			assertChildLocsWithinParent(t, sql, list)
		})
	}
}

// TestTemporalOracle checks accept/reject parity with SQL Server for every
// case in temporalCases.
func TestTemporalOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range temporalCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := tc.sql
			ssErr, err := oracle.parseError(sql)
			if err != nil {
				t.Fatalf("oracle error: %v", err)
			}
			ssAccepts := ssErr == nil
			if ssAccepts != tc.accepts {
				t.Errorf("table says SQL Server %s but it %s (%v):\n  SQL: %s",
					boolToAcceptReject(tc.accepts), boolToAcceptReject(ssAccepts), ssErr, sql)
			}

			_, omniErr := Parse(sql)
			omniAccepts := omniErr == nil
			if ssAccepts != omniAccepts {
				t.Errorf("MISMATCH: SQL Server %s (%v), omni %s (%v)\n  SQL: %s",
					boolToAcceptReject(ssAccepts), ssErr, boolToAcceptReject(omniAccepts), omniErr, sql)
			}
		})
	}
}

// TestTemporalLoc checks that FOR SYSTEM_TIME spans its clause and that the
// table spans its name, clause, and alias.
func TestTemporalLoc(t *testing.T) {
	const sql = "SELECT * FROM dbo.t FOR SYSTEM_TIME CONTAINED IN (@d, '2021-01-01') AS s WITH (NOLOCK) WHERE 1 = 1"
	tests := []struct {
		node  string
		index int // node of that type in pre-order
		want  string
	}{
		{"TableRef", 0, "dbo.t FOR SYSTEM_TIME CONTAINED IN (@d, '2021-01-01') AS s WITH (NOLOCK)"},
		{"TemporalClause", 0, "FOR SYSTEM_TIME CONTAINED IN (@d, '2021-01-01')"},
		{"VariableRef", 0, "@d"},
		{"Literal", 0, "'2021-01-01'"},
	}
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q) failed: %v", sql, err)
	}
	for _, tc := range tests {
		t.Run(tc.want, func(t *testing.T) {
			n := nthNodeByType(list, tc.node, tc.index)
			if n == nil {
				t.Fatalf("no %s #%d in %s", tc.node, tc.index, ast.NodeToString(list))
			}
			if got := locText(t, sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s #%d text = %q, want %q", tc.node, tc.index, got, tc.want)
			}
		})
	}
}

// TestTemporalCompletion checks the candidates offered after FOR SYSTEM_TIME.
func TestTemporalCompletion(t *testing.T) {
	sql := "SELECT * FROM t FOR SYSTEM_TIME "
	cs := Collect(sql, len(sql))
	if cs == nil {
		t.Fatal("Collect returned nil")
	}
	for _, tok := range []int{kwAS, kwFROM, kwBETWEEN, kwCONTAINED, kwALL} {
		if !cs.HasToken(tok) {
			t.Errorf("missing token candidate %s", TokenName(tok))
		}
	}
}
