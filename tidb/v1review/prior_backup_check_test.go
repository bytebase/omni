package v1review

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/metadata"
	"github.com/bytebase/omni/tidb/ast"
	"github.com/bytebase/omni/tidb/parser"
	"github.com/bytebase/omni/v1review"
)

// priorBackupFinding is a prior-backup finding as the tests compare it.
type priorBackupFinding struct {
	Code      v1review.Code
	Name      string
	Table     string
	Value     string
	Statement int
}

// priorBackupFacts is a prior backup that can take 100 bytes and finds its
// database.
func priorBackupFacts() v1review.PriorBackupFacts {
	return v1review.PriorBackupFacts{TotalSize: 100, MaxSize: 100, Database: "bbdataarchive", DatabaseExists: true}
}

// priorBackupSchema is the pre-change snapshot of database test, holding the
// tables the statements the tests judge name, and the columns unqualified SET
// clauses are resolved against.
func priorBackupSchema() v1review.Schema {
	return v1review.Schema{Metadata: &metadata.DatabaseSchemaMetadata{
		Name: "test",
		Schemas: []*metadata.SchemaMetadata{{Tables: []*metadata.TableMetadata{
			{Name: "tech_book", Columns: []*metadata.ColumnMetadata{{Name: "id"}, {Name: "name"}}},
			{Name: "orders", Columns: []*metadata.ColumnMetadata{{Name: "order_id"}, {Name: "name"}}},
		}}},
	}}
}

// parsePriorBackup parses sql the way Bytebase does, one statement per segment.
func parsePriorBackup(t *testing.T, sql string) []ast.Node {
	t.Helper()
	var nodes []ast.Node
	for _, segment := range parser.Split(sql) {
		if segment.Empty() {
			continue
		}
		list, err := parser.Parse(segment.Text)
		if err != nil {
			t.Fatalf("Parse %q: %v", segment.Text, err)
		}
		nodes = append(nodes, list.Items...)
	}
	return nodes
}

// judgePriorBackup parses sql and runs PriorBackupCheck over the statements,
// an unqualified table reference resolved to database.
func judgePriorBackup(t *testing.T, sql, database string, snapshot v1review.Schema, facts v1review.PriorBackupFacts) []priorBackupFinding {
	t.Helper()
	var findings []priorBackupFinding
	for _, f := range PriorBackupCheck(parsePriorBackup(t, sql), database, snapshot, facts) {
		findings = append(findings, priorBackupFinding{Code: f.Code, Name: f.Name, Table: f.Table, Value: f.Value, Statement: f.Statement})
	}
	return findings
}

func TestPriorBackupCheckDDLForms(t *testing.T) {
	// Every form ast.IsDDL counts, so a type dropped from its list fails here.
	ddl := []struct{ name, sql string }{
		{"alter database", "ALTER DATABASE test CHARACTER SET utf8mb4"},
		{"alter placement policy", "ALTER PLACEMENT POLICY p FOLLOWERS = 2"},
		{"alter resource group", "ALTER RESOURCE GROUP g RU_PER_SEC = 100"},
		{"alter sequence", "ALTER SEQUENCE sq RESTART WITH 5"},
		{"alter table", "ALTER TABLE t ADD COLUMN b int"},
		{"create database", "CREATE DATABASE d"},
		{"create index", "CREATE INDEX i ON t (a)"},
		{"create placement policy", "CREATE PLACEMENT POLICY p FOLLOWERS = 2"},
		{"create resource group", "CREATE RESOURCE GROUP g RU_PER_SEC = 100"},
		{"create sequence", "CREATE SEQUENCE sq"},
		{"create table", "CREATE TABLE t (a int)"},
		{"create view", "CREATE VIEW v AS SELECT 1"},
		{"drop database", "DROP DATABASE d"},
		{"drop index", "DROP INDEX i ON t"},
		{"drop placement policy", "DROP PLACEMENT POLICY p"},
		{"drop resource group", "DROP RESOURCE GROUP g"},
		{"drop sequence", "DROP SEQUENCE sq"},
		{"drop table", "DROP TABLE t"},
		{"drop view", "DROP VIEW v"},
		{"flashback database", "FLASHBACK DATABASE d TO d2"},
		{"flashback table", "FLASHBACK TABLE t TO t2"},
		{"flashback to timestamp", "FLASHBACK TABLE t TO TIMESTAMP '2020-01-01 00:00:00'"},
		{"lock tables", "LOCK TABLES t READ"},
		{"unlock tables", "UNLOCK TABLES"},
		{"optimize table", "OPTIMIZE TABLE t"},
		{"recover table", "RECOVER TABLE t"},
		{"rename table", "RENAME TABLE t TO u"},
		{"repair table", "REPAIR TABLE t"},
		{"truncate", "TRUNCATE TABLE t"},
		{"admin cleanup table lock", "ADMIN CLEANUP TABLE LOCK t"},
		{"admin repair table", "ADMIN REPAIR TABLE t CREATE TABLE t (a int)"},
	}
	want := []priorBackupFinding{{Code: v1review.PriorBackupMixedDDLAndDML}}
	for _, tt := range ddl {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, "test", priorBackupSchema(), priorBackupFacts())
			if diff := cmp.Diff(want, got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}

	// An ADMIN variant that is not a schema change, and the DML and query
	// statements, are not DDL.
	notDDL := []struct{ name, sql string }{
		{"admin check", "ADMIN CHECK TABLE t"},
		{"admin show ddl", "ADMIN SHOW DDL"},
		{"insert", "INSERT INTO t VALUES (1)"},
		{"replace", "REPLACE INTO t VALUES (1)"},
		{"select", "SELECT 1"},
		{"set", "SET NAMES utf8mb4"},
		{"begin", "BEGIN"},
		{"commit", "COMMIT"},
		{"explain", "EXPLAIN UPDATE t SET a = 1"},
	}
	for _, tt := range notDDL {
		t.Run("not ddl "+tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, "test", priorBackupSchema(), priorBackupFacts())
			if diff := cmp.Diff([]priorBackupFinding(nil), got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckUpdateTargets(t *testing.T) {
	table := func(database, name string) priorBackupTable {
		return priorBackupTable{database: database, table: name}
	}
	tests := []struct {
		name     string
		sql      string
		database string
		schema   v1review.Schema
		want     []priorBackupTable
	}{
		{
			name: "an unqualified target resolves to the reviewed database",
			sql:  "UPDATE tech_book SET id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name:     "an unqualified target resolves to the given database",
			sql:      "UPDATE tech_book SET id = 1",
			database: "other",
			want:     []priorBackupTable{table("other", "tech_book")},
		},
		{
			name: "a qualified target keeps its database",
			sql:  "UPDATE other.tech_book SET id = 1",
			want: []priorBackupTable{table("other", "tech_book")},
		},
		{
			name: "an alias names its table",
			sql:  "UPDATE tech_book AS x SET x.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "an alias without AS names its table",
			sql:  "UPDATE tech_book x SET x.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a table-qualified SET column names its table",
			sql:  "UPDATE tech_book SET tech_book.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a schema-qualified SET column names its table",
			sql:  "UPDATE tech_book SET test.tech_book.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a joined table that is only read is not a target",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET tech_book.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a joined table the SET clause names is a target",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET orders.note = 'x'",
			want: []priorBackupTable{table("test", "orders")},
		},
		{
			name: "two targets in the order the SET clause names them",
			sql:  "UPDATE tech_book AS x JOIN orders AS y ON x.id = y.order_id SET y.note = 'x', x.id = 1",
			want: []priorBackupTable{table("test", "orders"), table("test", "tech_book")},
		},
		{
			name: "a schema-qualified SET column tells same-named tables apart",
			sql:  "UPDATE db1.tech_book JOIN db2.tech_book ON db1.tech_book.id = db2.tech_book.id SET db1.tech_book.id = 1",
			want: []priorBackupTable{table("db1", "tech_book")},
		},
		{
			name: "a bare name two tables share is skipped",
			sql:  "UPDATE db1.tech_book JOIN db2.tech_book ON db1.tech_book.id = db2.tech_book.id SET tech_book.id = 1",
			want: nil,
		},
		{
			name: "an unqualified column in a single-table update",
			sql:  "UPDATE tech_book SET name = 'x'",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "an unqualified column in an aliased single-table update",
			sql:  "UPDATE tech_book AS t SET id = 1 WHERE id = 2",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "an unqualified column unique to one joined table resolves to that table alone",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET order_id = 1",
			want: []priorBackupTable{table("test", "orders")},
		},
		{
			name: "an unqualified column shared by both joined tables resolves to both",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET name = 'x'",
			want: []priorBackupTable{table("test", "tech_book"), table("test", "orders")},
		},
		{
			name: "a column no table has falls back to every table",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET nonexistent = 1",
			want: []priorBackupTable{table("test", "tech_book"), table("test", "orders")},
		},
		{
			name: "a column of another database's table is not looked up",
			sql:  "UPDATE other.tech_book JOIN test.tech_book ON other.tech_book.id = test.tech_book.id SET name = 'x'",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name:   "without a snapshot every table keeps the change",
			sql:    "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET name = 'x'",
			schema: v1review.Schema{},
			want:   []priorBackupTable{table("test", "tech_book"), table("test", "orders")},
		},
		{
			name: "a derived table is not a target",
			sql:  "UPDATE tech_book JOIN (SELECT 1 AS id) d ON tech_book.id = d.id SET tech_book.id = 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a derived table alone names no target",
			sql:  "UPDATE (SELECT 1 AS id) d SET d.id = 1",
			want: nil,
		},
		{
			name: "a reference to a CTE is not a target",
			sql:  "WITH c AS (SELECT 1 AS id) UPDATE tech_book JOIN c ON tech_book.id = c.id SET c.id = 1",
			want: nil,
		},
		{
			name: "a reference to a CTE alongside a table",
			sql:  "WITH c AS (SELECT 1 AS id) UPDATE tech_book JOIN c ON tech_book.id = c.id SET tech_book.id = 1, c.id = 2",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a CTE shadows a table of its name",
			sql:  "WITH orders AS (SELECT 1 AS id) UPDATE tech_book JOIN orders ON tech_book.id = orders.id SET orders.id = 1",
			want: nil,
		},
		{
			name: "a target named twice is one target",
			sql:  "UPDATE tech_book SET id = 1, name = 'x'",
			want: []priorBackupTable{table("test", "tech_book")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database := tt.database
			if database == "" {
				database = "test"
			}
			snapshot := tt.schema
			if snapshot.Metadata == nil && tt.name != "without a snapshot every table keeps the change" {
				snapshot = priorBackupSchema()
			}
			nodes := parsePriorBackup(t, tt.sql)
			if len(nodes) != 1 {
				t.Fatalf("expected one statement, got %d", len(nodes))
			}
			update, ok := nodes[0].(*ast.UpdateStmt)
			if !ok {
				t.Fatalf("expected an UPDATE, got %T", nodes[0])
			}
			got := updateTargets(update.SetList, newUpdateTableAliases(update.Tables, database, cteNames(update.CTEs)), snapshot.Metadata)
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(priorBackupTable{})); diff != "" {
				t.Fatalf("target mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckDeleteTargets(t *testing.T) {
	table := func(database, name string) priorBackupTable {
		return priorBackupTable{database: database, table: name}
	}
	tests := []struct {
		name string
		sql  string
		want []priorBackupTable
	}{
		{
			name: "an unqualified table",
			sql:  "DELETE FROM tech_book",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "a qualified table",
			sql:  "DELETE FROM other.tech_book",
			want: []priorBackupTable{table("other", "tech_book")},
		},
		{
			name: "a USING clause adds no target",
			sql:  "DELETE FROM tech_book USING orders",
			want: []priorBackupTable{table("test", "tech_book")},
		},
		{
			name: "two targets",
			sql:  "DELETE tech_book, orders FROM tech_book JOIN orders ON 1",
			want: []priorBackupTable{table("test", "tech_book"), table("test", "orders")},
		},
		{
			name: "a derived table adds no target",
			sql:  "DELETE tech_book FROM tech_book JOIN (SELECT 1 AS id) d ON 1",
			want: []priorBackupTable{table("test", "tech_book")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nodes := parsePriorBackup(t, tt.sql)
			if len(nodes) != 1 {
				t.Fatalf("expected one statement, got %d", len(nodes))
			}
			del, ok := nodes[0].(*ast.DeleteStmt)
			if !ok {
				t.Fatalf("expected a DELETE, got %T", nodes[0])
			}
			got := deleteTables(del.Tables, "test", cteNames(del.CTEs))
			if diff := cmp.Diff(tt.want, got, cmp.AllowUnexported(priorBackupTable{})); diff != "" {
				t.Fatalf("target mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckMixedDML(t *testing.T) {
	mixed := func(table string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: table}
	}
	tests := []struct {
		name string
		sql  string
		want []priorBackupFinding
	}{
		{
			name: "update then delete",
			sql:  "UPDATE tech_book SET id = 1; DELETE FROM tech_book",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "a case difference is one table",
			sql:  "UPDATE Tech_Book SET id = 1; DELETE FROM tech_book",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "two tables in sorted key order",
			sql:  "UPDATE b SET a = 1; DELETE FROM b; UPDATE a SET a = 1; DELETE FROM a",
			want: []priorBackupFinding{mixed("test.a"), mixed("test.b")},
		},
		{
			name: "an insert is not grouped",
			sql:  "INSERT INTO tech_book VALUES (1); UPDATE tech_book SET id = 1",
			want: nil,
		},
		{
			name: "two databases are two tables",
			sql:  "UPDATE db1.t SET a = 1; DELETE FROM test.t",
			want: nil,
		},
		{
			name: "the same kind twice is not mixed",
			sql:  "UPDATE tech_book SET id = 1; UPDATE tech_book SET id = 2",
			want: nil,
		},
		{
			name: "a column two tables share changes both",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET name = 'x'; DELETE FROM orders",
			want: []priorBackupFinding{mixed("test.orders")},
		},
		{
			name: "a column no table has changes every table",
			sql:  "UPDATE tech_book JOIN orders ON tech_book.id = orders.order_id SET nonexistent = 1; DELETE FROM orders",
			want: []priorBackupFinding{mixed("test.orders")},
		},
		{
			name: "a schema-qualified column names one of two same-named tables",
			sql:  "UPDATE db1.tech_book JOIN db2.tech_book ON db1.tech_book.id = db2.tech_book.id SET db1.tech_book.id = 1; DELETE FROM db1.tech_book",
			want: []priorBackupFinding{mixed("db1.tech_book")},
		},
		{
			name: "a USING clause table is not a target",
			sql:  "UPDATE orders SET note = 'x'; DELETE FROM tech_book USING orders",
			want: nil,
		},
		{
			name: "a multi-table DELETE target resolves through its USING clause",
			sql:  "UPDATE b SET id = 1; DELETE b FROM tech_book AS b WHERE b.id = 1",
			want: nil,
		},
		{
			name: "a multi-table DELETE target resolves to the table it changes",
			sql:  "UPDATE tech_book SET id = 1; DELETE b FROM tech_book AS b WHERE b.id = 1",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "a multi-table DELETE target resolves by the table's name",
			sql:  "UPDATE tech_book SET id = 1; DELETE tech_book FROM tech_book AS b WHERE b.id = 1",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "a multi-table DELETE target the USING clause does not declare is kept",
			sql:  "UPDATE z SET id = 1; DELETE z FROM tech_book AS b WHERE b.id = 1",
			want: []priorBackupFinding{mixed("test.z")},
		},
		{
			name: "an alias of a physical table is a target beside a CTE of its name",
			sql:  "UPDATE test.c SET id = 2; WITH c AS (SELECT 1 AS id) DELETE c FROM test.c AS c WHERE c.id = 1",
			want: []priorBackupFinding{mixed("test.c")},
		},
		{
			name: "a qualified target finds its table among same-named ones",
			sql:  "UPDATE db2.t SET id = 1; DELETE db2.t FROM db1.t JOIN db2.t ON db1.t.id = db2.t.id",
			want: []priorBackupFinding{mixed("db2.t")},
		},
		{
			name: "a WITH reference is not a table",
			sql:  "WITH t AS (SELECT 1 AS id) UPDATE t SET id = 2; DELETE FROM t",
			want: nil,
		},
		{
			name: "a schema-qualified name is a table even when a CTE matches it",
			sql:  "WITH tech_book AS (SELECT 1 AS id) UPDATE test.tech_book SET id = 1; DELETE FROM test.tech_book",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "update and delete on different tables in the same database",
			sql:  "UPDATE tech_book SET id = 1 WHERE id = 2; DELETE FROM orders WHERE order_id = 3",
			want: nil,
		},
		{
			name: "update join does not mark the read-only joined table",
			sql:  "UPDATE tech_book INNER JOIN orders ON tech_book.id = orders.order_id SET tech_book.name = 'x'; DELETE FROM orders WHERE order_id = 5",
			want: nil,
		},
		{
			name: "aliased single-table update with unqualified SET mixes with delete",
			sql:  "UPDATE tech_book AS t SET id = 1 WHERE id = 2; DELETE FROM tech_book WHERE id = 3",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "qualified and unqualified references are one table",
			sql:  "UPDATE test.tech_book SET id = 1 WHERE id = 2; DELETE FROM tech_book WHERE id = 3",
			want: []priorBackupFinding{mixed("test.tech_book")},
		},
		{
			name: "unqualified SET matching a single joined table mixes only on that table",
			sql:  "UPDATE tech_book INNER JOIN orders ON tech_book.id = orders.order_id SET order_id = 1; DELETE FROM orders WHERE order_id = 5",
			want: []priorBackupFinding{mixed("test.orders")},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, "test", priorBackupSchema(), priorBackupFacts())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckSize(t *testing.T) {
	tooLarge := priorBackupFinding{Code: v1review.PriorBackupScriptTooLarge, Value: "100"}
	mixed := priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: "test.tech_book"}
	tests := []struct {
		name  string
		sql   string
		facts v1review.PriorBackupFacts
		want  []priorBackupFinding
	}{
		{
			name:  "at the limit",
			sql:   "UPDATE tech_book SET id = 1",
			facts: priorBackupFacts(),
			want:  nil,
		},
		{
			name: "one byte over",
			sql:  "UPDATE tech_book SET id = 1",
			facts: func() v1review.PriorBackupFacts {
				f := priorBackupFacts()
				f.TotalSize = 101
				return f
			}(),
			want: []priorBackupFinding{tooLarge},
		},
		{
			name: "over the limit before the other findings",
			sql:  "UPDATE tech_book SET id = 1; DELETE FROM tech_book",
			facts: func() v1review.PriorBackupFacts {
				f := priorBackupFacts()
				f.TotalSize = 101
				return f
			}(),
			want: []priorBackupFinding{tooLarge, mixed},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, "test", priorBackupSchema(), tt.facts)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckCount(t *testing.T) {
	tooMany := priorBackupFinding{Code: v1review.PriorBackupTooManyDMLStatements, Value: "5"}
	mixed := func(table string) priorBackupFinding {
		return priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: table}
	}
	tests := []struct {
		name string
		sql  string
		want []priorBackupFinding
	}{
		{
			name: "five statements on two tables",
			sql:  "UPDATE a SET x = 1; UPDATE a SET x = 2; UPDATE a SET x = 3; UPDATE b SET y = 1; UPDATE b SET y = 2",
			want: nil,
		},
		{
			name: "six statements on one table",
			sql:  "UPDATE a SET x = 1; UPDATE a SET x = 2; UPDATE a SET x = 3; UPDATE a SET x = 4; UPDATE a SET x = 5; UPDATE a SET x = 6",
			want: nil,
		},
		{
			name: "six statements on two tables",
			sql:  "UPDATE a SET x = 1; UPDATE a SET x = 2; UPDATE a SET x = 3; UPDATE b SET y = 1; UPDATE b SET y = 2; UPDATE b SET y = 3",
			want: []priorBackupFinding{tooMany},
		},
		{
			name: "a joined statement changes as many tables as it names",
			// The count is of changes, as the backup transformer counts them:
			// one per statement and table.
			sql:  "UPDATE a JOIN b JOIN c ON 1 SET a.x = 1, b.y = 1, c.z = 1; DELETE FROM a; DELETE FROM b; DELETE FROM c; DELETE FROM a",
			want: []priorBackupFinding{mixed("test.a"), mixed("test.b"), mixed("test.c"), tooMany},
		},
		{
			name: "a WITH reference does not count as a change",
			sql:  "WITH c AS (SELECT 1 AS id) UPDATE tech_book JOIN c ON tech_book.id = c.id SET tech_book.id = 1, c.id = 2; UPDATE a SET x = 1; UPDATE a SET x = 2; UPDATE a SET x = 3; UPDATE a SET x = 4",
			want: nil,
		},
		{
			name: "zero-match unqualified SET falls back to every joined table and trips the count cap",
			sql: "UPDATE tech_book INNER JOIN orders ON tech_book.id = orders.order_id SET nonexistent = 1 WHERE tech_book.id = 1; " +
				"UPDATE tech_book INNER JOIN orders ON tech_book.id = orders.order_id SET nonexistent = 2 WHERE tech_book.id = 2; " +
				"UPDATE tech_book INNER JOIN orders ON tech_book.id = orders.order_id SET nonexistent = 3 WHERE tech_book.id = 3",
			want: []priorBackupFinding{tooMany},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := judgePriorBackup(t, tt.sql, "test", priorBackupSchema(), priorBackupFacts())
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestPriorBackupCheckDatabase(t *testing.T) {
	missing := priorBackupFinding{Code: v1review.PriorBackupDatabaseMissing, Name: "bbdataarchive"}
	ddl := priorBackupFinding{Code: v1review.PriorBackupMixedDDLAndDML}
	mixed := priorBackupFinding{Code: v1review.PriorBackupMixedDMLOnTable, Table: "test.tech_book"}
	absent := func() v1review.PriorBackupFacts {
		f := priorBackupFacts()
		f.DatabaseExists = false
		return f
	}
	tests := []struct {
		name      string
		sql       string
		database  string
		snapshot  v1review.Schema
		facts     v1review.PriorBackupFacts
		want      []priorBackupFinding
		wantNoSet bool
	}{
		{
			name:     "the database exists",
			sql:      "UPDATE tech_book SET id = 1; DELETE FROM tech_book",
			snapshot: priorBackupSchema(),
			facts:    priorBackupFacts(),
			want:     []priorBackupFinding{mixed},
		},
		{
			name:     "the database is missing",
			sql:      "UPDATE tech_book SET id = 1",
			snapshot: priorBackupSchema(),
			facts:    absent(),
			want:     []priorBackupFinding{missing},
		},
		{
			name:     "a missing database does not end the check",
			sql:      "CREATE TABLE u (a int); UPDATE tech_book SET id = 1; DELETE FROM tech_book",
			snapshot: priorBackupSchema(),
			facts:    absent(),
			want:     []priorBackupFinding{ddl, missing, mixed},
		},
		{
			name:     "every finding in order",
			sql:      "CREATE TABLE u (a int);\nUPDATE tech_book SET id = 1;\nDELETE FROM tech_book",
			snapshot: priorBackupSchema(),
			facts:    absent(),
			want:     []priorBackupFinding{ddl, missing, mixed},
		},
		{
			name: "an empty script still reports the missing database",
			sql:  "",
			facts: func() v1review.PriorBackupFacts {
				f := priorBackupFacts()
				f.DatabaseExists = false
				return f
			}(),
			want: []priorBackupFinding{missing},
		},
		{
			name:     "without a snapshot unqualified targets use the given database",
			sql:      "UPDATE tech_book SET id = 1; DELETE FROM tech_book",
			database: "test_db",
			facts:    absent(),
			want: []priorBackupFinding{
				missing,
				{Code: v1review.PriorBackupMixedDMLOnTable, Table: "test_db.tech_book"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database := tt.database
			if database == "" {
				database = "test"
			}
			got := judgePriorBackup(t, tt.sql, database, tt.snapshot, tt.facts)
			if tt.want == nil && got == nil {
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("finding mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
