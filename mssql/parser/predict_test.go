package parser

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bytebase/omni/mssql/ast"
)

// TestPredictWithListMakesProgress: the loop over PREDICT's WITH list
// skipped commas and otherwise called parsePredictColumnDef, which returned
// nil without consuming anything at a token that starts no column
// definition. A literal or a nested '(' in the list (WITH (1)) spun Parse
// forever. Completion stops parsing a data type once the cursor is behind
// it, which left the '(' of NVARCHAR(10) for the loop: Collect spun at 76 of
// the 94 cursor offsets of the first statement below, from the opening
// parenthesis through NVARCHAR.
func TestPredictWithListMakesProgress(t *testing.T) {
	within := func(what string, fn func()) {
		t.Helper()
		done := make(chan struct{})
		go func() {
			defer close(done)
			fn()
		}()
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("%s: no result within 10s", what)
		}
	}
	for _, sql := range []string{
		"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score FLOAT, Label NVARCHAR(10)) AS p",
		"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Score DECIMAL(10, 2) NOT NULL)",
		"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (Label NVARCHAR(10) COLLATE Latin1_General_CI_AS NULL)",
		"SELECT * FROM PREDICT (MODEL = @model, DATA = dbo.mytable AS d) WITH (",
	} {
		for off := 0; off <= len(sql); off++ {
			within(sql[:off]+"|"+sql[off:], func() { Collect(sql, off) })
		}
	}
	for _, sql := range []string{
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d) WITH (1)",
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT, (x))",
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT, 'a')",
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d) WITH (Score)",
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d) WITH (Score FLOAT",
		"SELECT * FROM PREDICT (MODEL = @m, DATA = t AS d",
		"SELECT * FROM PREDICT MODEL = @m",
	} {
		var err error
		within(sql, func() { _, err = Parse(sql) })
		if err == nil {
			t.Errorf("Parse(%q) = nil error, want a syntax error", sql)
		}
	}
}

// predictCases pin where SQL Server 2022 takes PREDICT, as SET PARSEONLY ON
// reports it, with the variables the statements use declared before them.
// PREDICT is a table source, not a statement. omni used to parse it as a
// statement and reject it in FROM.
//
// omni departs from SQL Server where reason says why: it keeps the Azure
// Synapse Analytics forms (RUNTIME = ONNX, a scalar subquery as MODEL), as
// it keeps the other Synapse syntax, and takes only the documented table
// sources as DATA and the documented column options in WITH.
var predictCases = []struct {
	sql    string
	server bool   // SQL Server's verdict
	omni   bool   // omni's verdict, when reason is set
	reason string // why omni departs from SQL Server
}{
	// Documented forms.
	{sql: "SELECT d.*, p.Score FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "INSERT INTO dbo.t (a) SELECT p.Score FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float COLLATE Latin1_General_CI_AS NOT NULL, Label nvarchar(10) NULL) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score decimal(10, 2), L nvarchar(max), u dbo.mytype) AS p", server: true},

	// The alias of DATA and of PREDICT are optional, with or without AS.
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t) WITH (Score float)", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t d) WITH (Score float) p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS [d]) WITH ([Score] float) AS [p]", server: true},

	// MODEL is a scalar expression; the arguments come in any order.
	{sql: "SELECT * FROM PREDICT(MODEL = 0x0102, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = N'abc', DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = CAST(@m AS VARBINARY(MAX)) + @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(DATA = dbo.t AS d, MODEL = @m) WITH (Score float) AS p", server: true},
	{sql: "select * from predict(model = @m, data = dbo.t as d) with (score float) as p", server: true},
	{sql: "SELECT * FROM [PREDICT](MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p", server: true},

	// DATA is a table, view, CTE, table variable, schema-qualified table-valued
	// function, or derived table.
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = testdb.dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = #tmp AS d) WITH (s float) AS p", server: true},
	{sql: "WITH c AS (SELECT 1 AS a) SELECT * FROM PREDICT(MODEL = @m, DATA = c AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = @tv AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.fn(1) AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT * FROM dbo.t) AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT 1 AS a UNION SELECT 2)) WITH (s float) AS p", server: true},

	// PREDICT joins, applies, pivots, and is a MERGE source like any table.
	{sql: "SELECT * FROM dbo.t AS x JOIN PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p ON 1 = 1", server: true},
	{sql: "SELECT * FROM dbo.t AS x CROSS APPLY PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p, PREDICT(MODEL = @m, DATA = dbo.t AS e) WITH (s float) AS q", server: true},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p PIVOT (SUM(Score) FOR a IN ([1])) AS q", server: true},
	{sql: "UPDATE dbo.t SET a = 1 FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "DELETE dbo.t FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p", server: true},
	{sql: "MERGE dbo.t AS x USING PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p ON 1 = 1 WHEN MATCHED THEN DELETE;", server: true},

	// PREDICT is not reserved: without name = value arguments it names a
	// table or a function.
	{sql: "SELECT predict FROM dbo.t", server: true},
	{sql: "SELECT * FROM predict AS p", server: true},
	{sql: "SELECT * FROM predict(1)", server: true},
	{sql: "SELECT * FROM predict(model)", server: true},
	{sql: "SELECT * FROM dbo.t CROSS APPLY predict(a) AS f", server: true},

	// Not a statement, and not an expression.
	{sql: "PREDICT (MODEL = @m, DATA = dbo.t AS d) WITH (Score float)"},
	{sql: "PREDICT (MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p"},
	{sql: "SELECT PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float)"},

	// The WITH list is required, and holds column definitions.
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH () AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float,) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (, Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (select int) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float NOT NULL COLLATE Latin1_General_CI_AS) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float NULL NOT NULL) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float PRIMARY KEY) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float DEFAULT 1) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float '$.x') AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float AS JSON) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float COLLATE [Latin1_General_CI_AS]) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float SPARSE) AS p", server: true, reason: "SQL Server's parser takes table-column storage options in WITH; the documented column takes none"},

	// MODEL and DATA once each, as name = value, and nothing else.
	{sql: "SELECT * FROM PREDICT(DATA = dbo.t AS d) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, DATA = dbo.t AS e) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d,) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m,, DATA = dbo.t AS d) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL @m, DATA = dbo.t AS d) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m DATA = dbo.t AS d) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, FOO = 1) WITH (Score float) AS p"},
	{sql: "SELECT * FROM predict(model = 1)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = ONNX) WITH (Score float) AS p", omni: true, reason: "RUNTIME = ONNX is documented for Azure Synapse Analytics"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = FOO) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = 'ONNX') WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, 'RUNTIME' = ONNX) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = (SELECT a FROM dbo.t), DATA = dbo.t AS d) WITH (Score float) AS p", omni: true, reason: "a scalar subquery as MODEL is documented for Azure Synapse Analytics"},

	// DATA takes no column alias list, table hint, FOR SYSTEM_TIME,
	// TABLESAMPLE, join, VALUES, rowset function, or unqualified function.
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d WITH (NOLOCK)) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t d (NOLOCK)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.fn(1) AS d(x)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = (SELECT 1 AS a) AS d(x)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t FOR SYSTEM_TIME AS OF '2020-01-01' AS d) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d TABLESAMPLE (10 PERCENT)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS a JOIN dbo.t AS b ON 1 = 1) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = (dbo.t AS a JOIN dbo.t AS b ON 1 = 1)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = (VALUES (1)) AS d(a)) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = OPENJSON(@j) AS d) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = CONTAINSTABLE(dbo.t, a, 'x') AS d) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = fn(1) AS d) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = PREDICT(MODEL = @m, DATA = dbo.t AS e) WITH (s float) AS d) WITH (Score float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS select) WITH (s float) AS p"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = 1) WITH (s float) AS p", server: true, reason: "SQL Server's parser reads DATA as an expression; 1 is no table source"},

	// After the alias: no column alias list, table hint, FOR SYSTEM_TIME, or
	// TABLESAMPLE, and no parentheses around PREDICT alone.
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p (x)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p WITH (NOLOCK)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS p (NOLOCK)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p FOR SYSTEM_TIME AS OF '2020-01-01'"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p TABLESAMPLE (10 PERCENT)"},
	{sql: "SELECT * FROM (PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (Score float) AS p)"},
	{sql: "SELECT * FROM PREDICT(MODEL = @m, DATA = dbo.t AS d) WITH (s float) AS select"},
}

// predictDeclarations declares the variables predictCases use, so that SQL
// Server's parser, which rejects an undeclared variable (Msg 137), checks
// the syntax alone.
const predictDeclarations = "DECLARE @m VARBINARY(MAX); DECLARE @j NVARCHAR(MAX); DECLARE @tv TABLE (a INT);\n"

// TestPredictOracle: omni takes PREDICT where SQL Server does, apart from
// the departures predictCases explains.
func TestPredictOracle(t *testing.T) {
	oracle := startParserOracle(t)
	for _, tc := range predictCases {
		t.Run(tc.sql, func(t *testing.T) {
			sql := predictDeclarations + tc.sql
			serverErr, err := oracle.parseError(sql)
			if err != nil {
				t.Fatal(err)
			}
			if (serverErr == nil) != tc.server {
				t.Fatalf("SQL Server: %v, want accept = %v", serverErr, tc.server)
			}
			want := tc.server
			if tc.reason != "" {
				want = tc.omni
			}
			if _, omniErr := Parse(sql); (omniErr == nil) != want {
				t.Fatalf("omni: %v, want accept = %v", omniErr, want)
			}
		})
	}
}

// TestPredictTableRef checks the tree PREDICT builds and that each node's
// Loc spans its text.
func TestPredictTableRef(t *testing.T) {
	const sql = "SELECT d.*, p.Score FROM PREDICT(MODEL = @m, DATA = dbo.t AS d, RUNTIME = ONNX) WITH (Score FLOAT NOT NULL, Label NVARCHAR(10) COLLATE Latin1_General_CI_AS NULL) AS p"
	list, err := Parse(sql)
	if err != nil {
		t.Fatal(err)
	}
	ref, ok := list.Items[0].(*ast.SelectStmt).FromClause.Items[0].(*ast.PredictTableRef)
	if !ok {
		t.Fatalf("FROM item is %T, want *ast.PredictTableRef", list.Items[0].(*ast.SelectStmt).FromClause.Items[0])
	}
	text := func(n ast.Node) string {
		loc := ast.NodeLoc(n)
		return sql[loc.Start:loc.End]
	}
	if got, want := text(ref), sql[strings.Index(sql, "PREDICT"):]; got != want {
		t.Errorf("PredictTableRef spans %q, want %q", got, want)
	}
	if got := text(ref.Model); got != "@m" {
		t.Errorf("Model spans %q", got)
	}
	data, ok := ref.Data.(*ast.TableRef)
	if !ok || data.Schema != "dbo" || data.Object != "t" || data.Alias != "d" || text(data) != "dbo.t AS d" {
		t.Errorf("Data = %#v", ref.Data)
	}
	if ref.Runtime != "ONNX" || ref.Alias != "p" {
		t.Errorf("Runtime %q, Alias %q", ref.Runtime, ref.Alias)
	}
	if ref.Columns == nil || len(ref.Columns.Items) != 2 {
		t.Fatalf("Columns = %v", ref.Columns)
	}
	score := ref.Columns.Items[0].(*ast.RowsetColumnDef)
	label := ref.Columns.Items[1].(*ast.RowsetColumnDef)
	if score.Name != "Score" || score.DataType.Name != "FLOAT" || score.Nullable == nil || !score.Nullable.NotNull || text(score) != "Score FLOAT NOT NULL" || text(score.Nullable) != "NOT NULL" {
		t.Errorf("first column %#v spans %q", score, text(score))
	}
	if label.Name != "Label" || label.Collation != "Latin1_General_CI_AS" || label.Nullable == nil || label.Nullable.NotNull || text(label.Nullable) != "NULL" {
		t.Errorf("second column %#v", label)
	}
	assertChildLocsWithinParent(t, sql, list)

	for _, tc := range []struct {
		data string
		want string // type of Data, and the text it spans
	}{
		{"dbo.t", "*ast.TableRef dbo.t"},
		{"@tv d", "*ast.TableVarRef @tv d"},
		{"dbo.fn(1, 2) AS d", "*ast.AliasedTableRef dbo.fn(1, 2) AS d"},
		{"(SELECT 1 AS a)", "*ast.SubqueryExpr (SELECT 1 AS a)"},
		{"(SELECT 1 AS a) AS d", "*ast.AliasedTableRef (SELECT 1 AS a) AS d"},
	} {
		sql := "SELECT * FROM PREDICT(DATA = " + tc.data + ", MODEL = 0x01) WITH (s float)"
		list, err := Parse(sql)
		if err != nil {
			t.Errorf("%s: %v", sql, err)
			continue
		}
		data := list.Items[0].(*ast.SelectStmt).FromClause.Items[0].(*ast.PredictTableRef).Data
		loc := ast.NodeLoc(data)
		if got := fmt.Sprintf("%T %s", data, sql[loc.Start:loc.End]); got != tc.want {
			t.Errorf("%s: Data is %s, want %s", sql, got, tc.want)
		}
		assertChildLocsWithinParent(t, sql, list)
	}
}
