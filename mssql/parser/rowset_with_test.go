package parser

import (
	"errors"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

// rowsetWithDeclare declares the variables the rowset cases pass, which SQL
// Server requires even under SET PARSEONLY ON.
const rowsetWithDeclare = "DECLARE @j NVARCHAR(MAX) = N'[]', @h INT; "

// rowsetWithCases covers the WITH schema of OPENJSON and OPENXML. A column
// is a name and a data type, then an undelimited COLLATE, a string-literal
// path, and, in OPENJSON, AS JSON, in that order; no other rowset function
// takes a WITH schema. shape is the WITH schema of an accepted statement,
// printed by rowsetWithShape. Every sql follows rowsetWithDeclare.
var rowsetWithCases = []struct {
	name    string
	sql     string
	accepts bool
	shape   string
}{
	// OPENJSON.
	{"json_type", "SELECT * FROM OPENJSON(@j) WITH (a INT) AS x", true, "(a INT)"},
	{"json_two_columns", "SELECT * FROM OPENJSON(@j) WITH (a INT, b INT) AS x", true, "(a INT, b INT)"},
	{"json_path", "SELECT * FROM OPENJSON(@j) WITH (a INT '$.a') AS x", true, "(a INT '$.a')"},
	{"json_nchar_path", "SELECT * FROM OPENJSON(@j) WITH (a INT N'$.a') AS x", true, "(a INT N'$.a')"},
	{"json_empty_path", "SELECT * FROM OPENJSON(@j) WITH (a INT '') AS x", true, "(a INT '')"},
	{"json_as_json", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' AS JSON) AS x", true, "(a NVARCHAR '$.a' AS JSON)"},
	{"json_as_json_without_path", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON) AS x", true, "(a NVARCHAR AS JSON)"},
	{"json_as_json_lower_case", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' as json) AS x", true, "(a NVARCHAR '$.a' AS JSON)"},
	{"json_as_json_any_type", "SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' AS JSON) AS x", true, "(a INT '$.a' AS JSON)"},
	{"json_collate", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS) AS x", true, "(a NVARCHAR COLLATE Latin1_General_CI_AS)"},
	{"json_collate_path", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS '$.a') AS x", true, "(a NVARCHAR COLLATE Latin1_General_CI_AS '$.a')"},
	{"json_collate_path_as_json", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) COLLATE Latin1_General_CI_AS '$.a' AS JSON) AS x", true, "(a NVARCHAR COLLATE Latin1_General_CI_AS '$.a' AS JSON)"},
	{"json_bracketed_name", "SELECT * FROM OPENJSON(@j) WITH ([select] INT) AS x", true, "(select INT)"},
	{"json_window_name", "SELECT * FROM OPENJSON(@j) WITH (window INT) AS x", true, "(window INT)"},
	{"json_duplicate_names", "SELECT * FROM OPENJSON(@j) WITH (a INT, a INT) AS x", true, "(a INT, a INT)"},
	{"json_schema_type", "SELECT * FROM OPENJSON(@j) WITH (a dbo.mytype) AS x", true, "(a dbo.mytype)"},
	{"json_no_alias", "SELECT * FROM OPENJSON(@j) WITH (a INT)", true, "(a INT)"},
	{"json_path_argument", "SELECT * FROM OPENJSON(@j, '$.x') WITH (a INT) j", true, "(a INT)"},
	{"json_without_with", "SELECT * FROM OPENJSON(@j) AS x", true, "no WITH"},
	{"json_empty", "SELECT * FROM OPENJSON(@j) WITH () AS x", false, ""},
	{"json_trailing_comma", "SELECT * FROM OPENJSON(@j) WITH (a INT,) AS x", false, ""},
	{"json_no_type", "SELECT * FROM OPENJSON(@j) WITH (a) AS x", false, ""},
	{"json_null", "SELECT * FROM OPENJSON(@j) WITH (a INT NULL) AS x", false, ""},
	{"json_number_path", "SELECT * FROM OPENJSON(@j) WITH (a INT 5) AS x", false, ""},
	{"json_variable_path", "SELECT * FROM OPENJSON(@j) WITH (a INT @j) AS x", false, ""},
	{"json_expression_path", "SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' + '') AS x", false, ""},
	{"json_two_paths", "SELECT * FROM OPENJSON(@j) WITH (a INT '$.a' '$.b') AS x", false, ""},
	{"json_path_then_collate", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) '$.a' COLLATE Latin1_General_CI_AS) AS x", false, ""},
	{"json_as_json_then_path", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON '$.a') AS x", false, ""},
	{"json_as_json_then_collate", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) AS JSON COLLATE Latin1_General_CI_AS) AS x", false, ""},
	{"json_as_without_json", "SELECT * FROM OPENJSON(@j) WITH (a INT AS) AS x", false, ""},
	{"json_as_xml", "SELECT * FROM OPENJSON(@j) WITH (a INT AS XML) AS x", false, ""},
	{"json_as_bracketed_json", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' AS [JSON]) AS x", false, ""},
	{"json_json_without_as", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(MAX) '$.a' JSON) AS x", false, ""},
	{"json_collate_missing", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE) AS x", false, ""},
	{"json_collate_bracketed", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE [Latin1_General_CI_AS]) AS x", false, ""},
	{"json_collate_string", "SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE 'x') AS x", false, ""},
	{"json_reserved_name", "SELECT * FROM OPENJSON(@j) WITH (select INT) AS x", false, ""},
	{"json_qualified_name", "SELECT * FROM OPENJSON(@j) WITH (a.b INT) AS x", false, ""},
	{"json_column_aliases", "SELECT * FROM OPENJSON(@j) WITH (a INT) x (c)", false, ""},

	// OPENXML.
	{"xml_path", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a')", true, "(a INT '@a')"},
	{"xml_type", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT)", true, "(a INT)"},
	{"xml_two_columns", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a', b NVARCHAR(10) '@mp:id')", true, "(a INT '@a', b NVARCHAR '@mp:id')"},
	{"xml_collate_path", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a NVARCHAR(10) COLLATE Latin1_General_CI_AS '@a') x", true, "(a NVARCHAR COLLATE Latin1_General_CI_AS '@a')"},
	{"xml_two_arguments", "SELECT * FROM OPENXML(@h, '/r') WITH (a INT) y", true, "(a INT)"},
	{"xml_as_alias", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT) AS y", true, "(a INT)"},
	{"xml_without_with", "SELECT * FROM OPENXML(@h, '/r', 1) y", true, "no WITH"},
	{"xml_number_path", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT 1)", false, ""},
	{"xml_as_json", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a' AS JSON)", false, ""},
	{"xml_path_then_collate", "SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a' COLLATE Latin1_General_CI_AS) x", false, ""},

	// Other rowset functions take no WITH schema.
	{"openquery_with", "SELECT * FROM OPENQUERY(srv, 'SELECT 1') WITH (a INT) x", false, ""},
	{"openrowset_with", "SELECT * FROM OPENROWSET('SQLNCLI', 'x', 'SELECT 1') WITH (a INT) x", false, ""},
}

// rowsetWithShape prints the WITH schema of a rowset function.
func rowsetWithShape(ref *ast.AliasedTableRef) string {
	if ref.Columns == nil {
		return "no WITH"
	}
	var cols []string
	for _, item := range ref.Columns.Items {
		col := item.(*ast.RowsetColumnDef)
		s := col.Name + " " + col.DataType.Name
		if col.DataType.Schema != "" {
			s = col.Name + " " + col.DataType.Schema + "." + col.DataType.Name
		}
		if col.Collation != "" {
			s += " COLLATE " + col.Collation
		}
		if col.Path != nil {
			s += " "
			if col.Path.IsNChar {
				s += "N"
			}
			s += "'" + col.Path.Str + "'"
		}
		if col.AsJSON {
			s += " AS JSON"
		}
		cols = append(cols, s)
	}
	return "(" + strings.Join(cols, ", ") + ")"
}

// TestRowsetWith checks that omni accepts the OPENJSON and OPENXML WITH
// schemas SQL Server accepts, rejects the others, and keeps every Loc inside
// its parent.
func TestRowsetWith(t *testing.T) {
	for _, tc := range rowsetWithCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := rowsetWithDeclare + tc.sql
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
			ref, ok := nthNodeByType(list, "AliasedTableRef", 0).(*ast.AliasedTableRef)
			if !ok {
				t.Fatalf("no AliasedTableRef in %s", ast.NodeToString(list))
			}
			if got := rowsetWithShape(ref); got != tc.shape {
				t.Errorf("shape = %s, want %s", got, tc.shape)
			}
			assertChildLocsWithinParent(t, sql, list)
		})
	}
}

// TestRowsetWithOracle checks accept/reject parity with SQL Server for every
// case in rowsetWithCases.
func TestRowsetWithOracle(t *testing.T) {
	oracle := startParserOracle(t)

	for _, tc := range rowsetWithCases {
		t.Run(tc.name, func(t *testing.T) {
			sql := rowsetWithDeclare + tc.sql
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

// TestRowsetWithLoc checks the span of a WITH schema column, its path, and
// the rowset function it belongs to.
func TestRowsetWithLoc(t *testing.T) {
	const sql = "SELECT * FROM OPENJSON(@j) WITH (a INT, b NVARCHAR(MAX) COLLATE Latin1_General_CI_AS N'$.b' AS JSON) AS x WHERE 1 = 1"
	tests := []struct {
		node  string
		index int // node of that type in pre-order
		want  string
	}{
		{"AliasedTableRef", 0, "OPENJSON(@j) WITH (a INT, b NVARCHAR(MAX) COLLATE Latin1_General_CI_AS N'$.b' AS JSON) AS x"},
		{"RowsetColumnDef", 0, "a INT"},
		{"RowsetColumnDef", 1, "b NVARCHAR(MAX) COLLATE Latin1_General_CI_AS N'$.b' AS JSON"},
		{"Literal", 0, "N'$.b'"},
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

// TestRowsetWithErrorPosition checks where a malformed WITH schema column is
// reported.
func TestRowsetWithErrorPosition(t *testing.T) {
	tests := []struct {
		sql  string
		want string
	}{
		{"SELECT * FROM OPENJSON(@j) WITH (a) AS x", `syntax error at or near ")"`},
		{"SELECT * FROM OPENJSON(@j) WITH (a INT 5) AS x", `syntax error at or near "5"`},
		{"SELECT * FROM OPENJSON(@j) WITH (a INT NULL) AS x", `syntax error at or near "NULL"`},
		{"SELECT * FROM OPENJSON(@j) WITH (a INT AS XML) AS x", `syntax error at or near "AS"`},
		{"SELECT * FROM OPENJSON(@j) WITH (a NVARCHAR(10) COLLATE [x]) AS x", `syntax error at or near "["`},
		{"SELECT * FROM OPENXML(@h, '/r', 1) WITH (a INT '@a' AS JSON)", `syntax error at or near "AS"`},
	}
	for _, tc := range tests {
		t.Run(tc.sql, func(t *testing.T) {
			_, err := Parse(tc.sql)
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("Parse(%q) error = %v, want a *ParseError", tc.sql, err)
			}
			if got := pe.Error(); got != tc.want {
				t.Errorf("Error() = %q, want %q", got, tc.want)
			}
		})
	}
}
