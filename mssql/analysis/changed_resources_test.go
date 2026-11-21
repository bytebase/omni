package analysis

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/mssql"
)

func TestExtractChangedResources(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementChanges
	}{
		{
			statement: "CREATE TABLE t (a INT)",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}},
		},
		{
			statement: "DROP TABLE d2.s.t, u",
			want: StatementChanges{Resources: []ChangedResource{
				{Database: "d2", Schema: "s", Table: "t", AffectsData: true},
				{Database: "db", Schema: "dbo", Table: "u", AffectsData: true},
			}},
		},
		{
			statement: "INSERT INTO t VALUES (1), (2)",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}, InsertCount: 2},
		},
		{
			statement: "INSERT INTO t EXEC p",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "t"}}, DMLCount: 1},
		},
		{
			// An unqualified target binds to the FROM table it aliases.
			statement: "UPDATE x SET a = 1 FROM s.t AS x;  ",
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "db", Schema: "s", Table: "t"}},
				DMLStatements: []string{"UPDATE x SET a = 1 FROM s.t AS x;"},
				DMLCount:      1,
			},
		},
		{
			// A CTE target stays as written.
			statement: "WITH c AS (SELECT a FROM t) DELETE c FROM c",
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "db", Schema: "dbo", Table: "c"}},
				DMLStatements: []string{"WITH c AS (SELECT a FROM t) DELETE c FROM c;"},
				DMLCount:      1,
			},
		},
		{
			statement: "SELECT a INTO new_t FROM x UNION SELECT b FROM y",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "db", Schema: "dbo", Table: "new_t"}}},
		},
		{statement: "SELECT a INTO #t FROM x"},
		{
			statement: "CREATE VIEW other.dbo.v AS SELECT 1 AS a",
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other", DatabaseOnly: true}}},
		},
		{statement: "CREATE VIEW dbo.v AS SELECT 1 AS a"},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		got := ExtractChangedResources(stmts[0].AST, "db", "dbo", tc.statement)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q:\ngot  %+v\nwant %+v", tc.statement, got, tc.want)
		}
	}
}

func TestExtractChangedResourcesDMLWriteTarget(t *testing.T) {
	tests := []struct {
		name      string
		statement string
		database  string
		schema    string
		table     string
	}{
		{"alias", `UPDATE o SET o.status = 'x' FROM dbo.Orders AS o WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"alias without AS", `UPDATE o SET status = 'x' FROM dbo.Orders o WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"alias in DELETE", `DELETE o FROM dbo.Orders AS o WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"alias on the right of a JOIN", `UPDATE o SET o.status = 'x' FROM dbo.Customers AS c JOIN dbo.Orders AS o ON o.id = c.id WHERE c.id = 1;`, "DB", "dbo", "Orders"},
		{"alias beside an APPLY", `UPDATE o SET o.status = 'x' FROM dbo.Orders AS o CROSS APPLY (SELECT TOP 1 c.id FROM dbo.Customers AS c WHERE c.id = o.id) AS a WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"alias shadows a table of its name", `UPDATE Customers SET status = 'x' FROM dbo.Orders AS Customers WHERE Customers.id = 1;`, "DB", "dbo", "Orders"},
		{"alias of a table in another database", `UPDATE o SET o.status = 'x' FROM otherdb.dbo.Orders AS o WHERE o.id = 1;`, "otherdb", "dbo", "Orders"},
		{"unaliased table in another database", `UPDATE Orders SET status = 'x' FROM otherdb.dbo.Orders WHERE Orders.id = 1;`, "otherdb", "dbo", "Orders"},
		{"unaliased table in another schema", `UPDATE Orders SET status = 'x' FROM sales.Orders WHERE Orders.id = 1;`, "DB", "sales", "Orders"},
		{"unaliased table in DELETE FROM", `DELETE FROM Orders FROM otherdb.dbo.Orders WHERE Orders.id = 1;`, "otherdb", "dbo", "Orders"},
		{"unaliased join operand", `UPDATE Orders SET status = 'x' FROM dbo.Customers AS c LEFT JOIN otherdb.dbo.Orders ON c.id = Orders.id WHERE c.id = 1;`, "otherdb", "dbo", "Orders"},
		{"unaliased sampled table", `UPDATE Orders SET status = 'x' FROM otherdb.dbo.Orders TABLESAMPLE (100 PERCENT) WHERE Orders.id = 1;`, "otherdb", "dbo", "Orders"},
		{"aliased table hides its own name", `UPDATE Orders SET status = 'x' FROM otherdb.dbo.Orders AS x WHERE x.id = 1;`, "DB", "dbo", "Orders"},
		{"qualified target ignores an alias of its bare name", `UPDATE dbo.Orders SET status = 'x' FROM dbo.Orders AS o JOIN dbo.Customers AS Orders ON Orders.id = o.id WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"qualified DELETE target ignores an alias of its bare name", `DELETE dbo.Orders FROM dbo.Orders AS o JOIN dbo.Customers AS Orders ON Orders.id = o.id WHERE o.id = 1;`, "DB", "dbo", "Orders"},
		{"qualified target in another database ignores an alias of its bare name", `UPDATE otherdb.dbo.Orders SET status = 'x' FROM otherdb.dbo.Orders AS o JOIN dbo.Customers AS Orders ON Orders.id = o.id WHERE o.id = 1;`, "otherdb", "dbo", "Orders"},
		{"server-qualified target ignores an alias of its bare name", `UPDATE srv...Orders SET status = 'x' FROM srv...Orders AS t JOIN dbo.Customers AS Orders ON Orders.id = t.id WHERE t.id = 1;`, "DB", "dbo", "Orders"},
		{"differently cased alias stays as written", `UPDATE O SET status = 'x' FROM dbo.Orders AS o WHERE o.id = 1;`, "DB", "dbo", "O"},
		{"differently cased table name stays as written", `UPDATE ORDERS SET status = 'x' FROM otherdb.dbo.Orders WHERE Orders.id = 1;`, "DB", "dbo", "ORDERS"},
		{"alias of a CTE stays as written", `WITH c AS (SELECT id, status FROM dbo.Orders) UPDATE x SET status = 'v' FROM c AS x WHERE x.id = 1;`, "DB", "dbo", "x"},
		{"alias of a differently cased CTE name stays as written", `WITH c AS (SELECT id, status FROM dbo.Orders) UPDATE x SET status = 'v' FROM C AS x WHERE x.id = 1;`, "DB", "dbo", "x"},
		{"alias of a derived table stays as written", `UPDATE d SET status = 'x' FROM (SELECT * FROM dbo.Orders) AS d WHERE d.id = 1;`, "DB", "dbo", "d"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			stmts, err := mssql.Parse(tc.statement)
			if err != nil || len(stmts) != 1 {
				t.Fatalf("%q: %v", tc.statement, err)
			}
			got := ExtractChangedResources(stmts[0].AST, "DB", "dbo", tc.statement)
			want := []ChangedResource{{Database: tc.database, Schema: tc.schema, Table: tc.table}}
			if !reflect.DeepEqual(got.Resources, want) {
				t.Errorf("%q:\ngot  %+v\nwant %+v", tc.statement, got.Resources, want)
			}
		})
	}
}

func TestExtractChangedResourcesAdditionalDDLAndDML(t *testing.T) {
	tests := []struct {
		statement string
		want      StatementChanges
	}{
		{
			statement: `DROP INDEX idx1 ON t1;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "DB", Schema: "dbo", Table: "t1"}}},
		},
		{
			statement: `INSERT INTO t1 DEFAULT VALUES;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "DB", Schema: "dbo", Table: "t1"}}, InsertCount: 1},
		},
		{
			statement: `TRUNCATE TABLE dbo.t;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "DB", Schema: "dbo", Table: "t", AffectsData: true}}},
		},
		{
			statement: `MERGE INTO other_db.dbo.tgt AS t USING dbo.src AS s ON (t.id = s.id) WHEN MATCHED THEN UPDATE SET t.a = s.a;`,
			want: StatementChanges{
				Resources:     []ChangedResource{{Database: "other_db", Schema: "dbo", Table: "tgt"}},
				DMLStatements: []string{`MERGE INTO other_db.dbo.tgt AS t USING dbo.src AS s ON (t.id = s.id) WHEN MATCHED THEN UPDATE SET t.a = s.a;`},
				DMLCount:      1,
			},
		},
		{
			statement: `DROP VIEW other_db.dbo.v;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other_db", DatabaseOnly: true}}},
		},
		{
			statement: `CREATE SEQUENCE other_db.dbo.s;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other_db", DatabaseOnly: true}}},
		},
		{
			statement: `CREATE SYNONYM other_db.dbo.syn FOR t;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other_db", DatabaseOnly: true}}},
		},
		{
			statement: `DROP TYPE other_db.dbo.ty;`,
			want:      StatementChanges{Resources: []ChangedResource{{Database: "other_db", DatabaseOnly: true}}},
		},
	}
	for _, tc := range tests {
		stmts, err := mssql.Parse(tc.statement)
		if err != nil || len(stmts) != 1 {
			t.Fatalf("%q: %v", tc.statement, err)
		}
		got := ExtractChangedResources(stmts[0].AST, "DB", "dbo", tc.statement)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("%q:\ngot  %+v\nwant %+v", tc.statement, got, tc.want)
		}
	}
}
