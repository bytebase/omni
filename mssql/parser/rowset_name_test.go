package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// rowsetNameDeclare declares the XML document handle the OPENXML cases pass,
// which SQL Server requires even under SET PARSEONLY ON.
const rowsetNameDeclare = "DECLARE @h INT; "

// adHocSource is an OPENDATASOURCE call with a provider this SQL Server
// instance allows.
const adHocSource = "OPENDATASOURCE('MSOLEDBSQL', 'Data Source=x;Integrated Security=SSPI')"

// rowsetNameCases covers the table names that follow a rowset function:
// OPENXML ... WITH table_name, which takes the columns of a local table, and
// OPENDATASOURCE ( provider_name , init_string ) . database . schema . object.
// shape is the first such table source of an accepted statement, printed by
// rowsetNameShape. Every sql follows rowsetNameDeclare.
var rowsetNameCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// OPENXML WITH table_name.
	{"xml_table", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t", true, "OPENXML WITH t"},
	{"xml_schema_table", "SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.t", true, "OPENXML WITH dbo.t"},
	{"xml_database_table", "SELECT * FROM OPENXML(@h, '/r', 1) WITH db.dbo.t", true, "OPENXML WITH db.dbo.t"},
	{"xml_empty_schema", "SELECT * FROM OPENXML(@h, '/r', 1) WITH db..t", true, "OPENXML WITH db..t"},
	{"xml_bracketed", "SELECT * FROM OPENXML(@h, '/r', 1) WITH [dbo].[t]", true, "OPENXML WITH dbo.t"},
	{"xml_as_alias", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t AS y", true, "OPENXML WITH t AS y"},
	{"xml_bare_alias", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t y", true, "OPENXML WITH t AS y"},
	{"xml_then_comma", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t, t1", true, "OPENXML WITH t"},
	{"xml_then_join", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t JOIN t1 ON 1 = 1", true, "OPENXML WITH t"},
	{"xml_alias_then_join", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t y JOIN t1 ON 1 = 1", true, "OPENXML WITH t AS y"},
	{"xml_remote_table", "SELECT * FROM OPENXML(@h, '/r', 1) WITH srv.db.dbo.t", false, ""},
	{"xml_trailing_dot", "SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.", false, ""},
	{"xml_trailing_dots", "SELECT * FROM OPENXML(@h, '/r', 1) WITH db..", false, ""},
	{"xml_trailing_dot_alias", "SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo. AS y", false, ""},
	{"xml_three_parts_trailing_dot", "SELECT * FROM OPENXML(@h, '/r', 1) WITH db.dbo.", false, ""},
	{"xml_space_after_dot", "SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo. t", true, "OPENXML WITH dbo.t"},
	{"xml_table_variable", "DECLARE @tv TABLE (a INT); SELECT * FROM OPENXML(@h, '/r', 1) WITH @tv", false, ""},
	{"xml_column_list", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t y (c)", false, ""},
	{"xml_table_hint", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t WITH (NOLOCK)", false, ""},
	{"xml_alias_table_hint", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t y WITH (NOLOCK)", false, ""},
	{"xml_old_style_hint", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t (NOLOCK)", false, ""},
	{"xml_tablesample", "SELECT * FROM OPENXML(@h, '/r', 1) WITH t TABLESAMPLE (10 PERCENT)", false, ""},
	{"xml_derived_table", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (SELECT 1 AS a)", false, ""},
	{"xml_function", "SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.f()", false, ""},
	{"xml_reserved_name", "SELECT * FROM OPENXML(@h, '/r', 1) WITH select", false, ""},
	{"json_table", "SELECT * FROM OPENJSON(N'[]') WITH t", false, ""},
	{"json_table_alias", "SELECT * FROM OPENJSON(N'[]') WITH dbo.t AS y", false, ""},

	// OPENDATASOURCE ( provider_name , init_string ) . database . schema . object.
	{"ad_hoc_table", "SELECT * FROM " + adHocSource + ".db.dbo.t", true, "OPENDATASOURCE db.dbo.t"},
	{"ad_hoc_as_alias", "SELECT * FROM " + adHocSource + ".db.dbo.t AS x", true, "OPENDATASOURCE db.dbo.t AS x"},
	{"ad_hoc_bare_alias", "SELECT * FROM " + adHocSource + ".db.dbo.t x", true, "OPENDATASOURCE db.dbo.t AS x"},
	{"ad_hoc_empty_schema", "SELECT * FROM " + adHocSource + ".db..t", true, "OPENDATASOURCE db..t"},
	{"ad_hoc_empty_database_schema", "SELECT * FROM " + adHocSource + "...t", true, "OPENDATASOURCE t"},
	{"ad_hoc_bracketed", "SELECT * FROM " + adHocSource + ".[db].[dbo].[t]", true, "OPENDATASOURCE db.dbo.t"},
	{"ad_hoc_space_after_dot", "SELECT * FROM " + adHocSource + ". db.dbo.t", true, "OPENDATASOURCE db.dbo.t"},
	{"ad_hoc_then_join", "SELECT * FROM " + adHocSource + ".db.dbo.t JOIN t1 ON 1 = 1", true, "OPENDATASOURCE db.dbo.t"},
	{"ad_hoc_join_right", "SELECT * FROM t1 JOIN " + adHocSource + ".db.dbo.t x ON x.a = t1.a", true, "OPENDATASOURCE db.dbo.t AS x"},
	{"ad_hoc_call_alone", "SELECT * FROM " + adHocSource, false, ""},
	{"ad_hoc_call_alias", "SELECT * FROM " + adHocSource + " x", false, ""},
	{"ad_hoc_one_part", "SELECT * FROM " + adHocSource + ".t", false, ""},
	{"ad_hoc_two_parts", "SELECT * FROM " + adHocSource + ".dbo.t", false, ""},
	{"ad_hoc_four_parts", "SELECT * FROM " + adHocSource + ".x.db.dbo.t", false, ""},
	{"ad_hoc_column", "SELECT * FROM " + adHocSource + ".db.dbo.t.c", false, ""},
	{"ad_hoc_function", "SELECT * FROM " + adHocSource + ".db.dbo.f()", false, ""},
	{"ad_hoc_table_hint", "SELECT * FROM " + adHocSource + ".db.dbo.t WITH (NOLOCK)", false, ""},
	{"ad_hoc_alias_table_hint", "SELECT * FROM " + adHocSource + ".db.dbo.t x WITH (NOLOCK)", false, ""},
	{"ad_hoc_old_style_hint", "SELECT * FROM " + adHocSource + ".db.dbo.t x (NOLOCK)", false, ""},
	{"ad_hoc_tablesample", "SELECT * FROM " + adHocSource + ".db.dbo.t TABLESAMPLE (10 PERCENT)", false, ""},
	{"ad_hoc_for_system_time", "SELECT * FROM " + adHocSource + ".db.dbo.t FOR SYSTEM_TIME ALL", false, ""},
	{"ad_hoc_one_argument", "SELECT * FROM OPENDATASOURCE('MSOLEDBSQL').db.dbo.t", false, ""},
	{"ad_hoc_three_arguments", "SELECT * FROM OPENDATASOURCE('a', 'b', 'c').db.dbo.t", false, ""},
	{"ad_hoc_variable_argument", "SELECT * FROM OPENDATASOURCE(@h, 'b').db.dbo.t", false, ""},
	{"ad_hoc_identifier_argument", "SELECT * FROM OPENDATASOURCE(MSOLEDBSQL, 'b').db.dbo.t", false, ""},
	{"openquery_name", "SELECT * FROM OPENQUERY(srv, 'SELECT 1').db.dbo.t", false, ""},
	{"openrowset_name", "SELECT * FROM OPENROWSET('SQLNCLI', 'x', 'SELECT 1').db.dbo.t", false, ""},
}

// rowsetNameShape prints the first OPENXML WITH table or OPENDATASOURCE table
// in pre-order, or "" when there is none.
func rowsetNameShape(root ast.Node) string {
	name := func(ref *ast.TableRef) string {
		if ref.Database == "" && ref.Schema == "" {
			return ref.Object
		}
		if ref.Database == "" {
			return ref.Schema + "." + ref.Object
		}
		return strings.Join([]string{ref.Database, ref.Schema, ref.Object}, ".")
	}
	var shape string
	ast.Inspect(root, func(n ast.Node) bool {
		if shape != "" {
			return false
		}
		switch v := n.(type) {
		case *ast.AliasedTableRef:
			if v.WithTable != nil {
				shape = "OPENXML WITH " + name(v.WithTable)
				if v.Alias != "" {
					shape += " AS " + v.Alias
				}
			}
		case *ast.AdHocTableRef:
			shape = "OPENDATASOURCE " + name(v.Object)
			if v.Alias != "" {
				shape += " AS " + v.Alias
			}
		}
		return shape == ""
	})
	return shape
}

// TestRowsetName checks that omni accepts the rowset function table names SQL
// Server accepts, rejects the others, and keeps every Loc inside its parent.
func TestRowsetName(t *testing.T) {
	for _, tc := range rowsetNameCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := rowsetNameDeclare + tc.sql
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
			if got := rowsetNameShape(list); got != tc.shape {
				t.Errorf("shape = %q, want %q", got, tc.shape)
			}
			assertChildLocsWithinParent(t, sql, list)
		})
	}
}

// TestRowsetNameOracle checks accept/reject parity with SQL Server for every
// case in rowsetNameCases.
func TestRowsetNameOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range rowsetNameCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := rowsetNameDeclare + tc.sql
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

// TestRowsetNameLoc checks the spans of an OPENDATASOURCE table and an OPENXML
// WITH table.
func TestRowsetNameLoc(t *testing.T) {
	tests := []struct {
		sql   string
		node  string
		index int // node of that type in pre-order
		want  string
	}{
		{"SELECT * FROM " + adHocSource + ".db.dbo.t AS x WHERE 1 = 1", "AdHocTableRef", 0, adHocSource + ".db.dbo.t AS x"},
		{"SELECT * FROM " + adHocSource + ".db.dbo.t AS x WHERE 1 = 1", "TableRef", 1, "db.dbo.t"},
		{"SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.t AS y WHERE 1 = 1", "AliasedTableRef", 0, "OPENXML(@h, '/r', 1) WITH dbo.t AS y"},
		{"SELECT * FROM OPENXML(@h, '/r', 1) WITH dbo.t AS y WHERE 1 = 1", "TableRef", 1, "dbo.t"},
	}
	for _, tc := range tests {
		t.Run(tc.node+"/"+tc.want, func(t *testing.T) {
			list, err := Parse(tc.sql)
			if err != nil {
				t.Fatalf("Parse(%q) failed: %v", tc.sql, err)
			}
			n := nthNodeByType(list, tc.node, tc.index)
			if n == nil {
				t.Fatalf("no %s #%d in %s", tc.node, tc.index, ast.NodeToString(list))
			}
			if got := locText(t, tc.sql, ast.NodeLoc(n)); got != tc.want {
				t.Errorf("%s #%d text = %q, want %q", tc.node, tc.index, got, tc.want)
			}
		})
	}
}
