package cassandra

import (
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
)

// parseOne parses sql and returns its single statement's AST.
func parseOne(t *testing.T, sql string) ast.Node {
	t.Helper()
	stmts, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(stmts) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(stmts))
	}
	return stmts[0].AST
}

func TestParseDescribe(t *testing.T) {
	tests := []struct {
		sql       string
		kind      ast.DescribeKind
		name      string // the name's parts joined by '.', "" for none
		full      bool
		only      bool
		internals bool
	}{
		{sql: "DESCRIBE CLUSTER", kind: ast.DescribeCluster},
		{sql: "DESC CLUSTER", kind: ast.DescribeCluster},
		{sql: "DESCRIBE SCHEMA", kind: ast.DescribeSchema},
		{sql: "DESCRIBE FULL SCHEMA", kind: ast.DescribeSchema, full: true},
		{sql: "DESCRIBE KEYSPACES", kind: ast.DescribeKeyspaces},
		{sql: "DESCRIBE KEYSPACE", kind: ast.DescribeKeyspace},
		{sql: "DESCRIBE KEYSPACE ks", kind: ast.DescribeKeyspace, name: "ks"},
		{sql: "DESCRIBE ONLY KEYSPACE ks", kind: ast.DescribeKeyspace, name: "ks", only: true},
		{sql: "DESCRIBE TABLES", kind: ast.DescribeTables},
		{sql: "DESCRIBE TABLE ks.t", kind: ast.DescribeTable, name: "ks.t"},
		{sql: "DESCRIBE INDEX i", kind: ast.DescribeIndex, name: "i"},
		{sql: "DESCRIBE MATERIALIZED VIEW v", kind: ast.DescribeMaterializedView, name: "v"},
		{sql: "DESCRIBE TYPES", kind: ast.DescribeTypes},
		{sql: "DESCRIBE TYPE ks.addr", kind: ast.DescribeType, name: "ks.addr"},
		{sql: "DESCRIBE FUNCTIONS", kind: ast.DescribeFunctions},
		{sql: "DESCRIBE FUNCTION f", kind: ast.DescribeFunction, name: "f"},
		{sql: "DESCRIBE AGGREGATES", kind: ast.DescribeAggregates},
		{sql: "DESCRIBE AGGREGATE ks.a", kind: ast.DescribeAggregate, name: "ks.a"},
		// TYPE, FUNCTION, AGGREGATE and ONLY with no name after them are the
		// name described; so are other unreserved keywords.
		{sql: "DESCRIBE TYPE", kind: ast.DescribeName, name: "TYPE"},
		{sql: "DESCRIBE FUNCTION", kind: ast.DescribeName, name: "FUNCTION"},
		{sql: "DESCRIBE ONLY", kind: ast.DescribeName, name: "ONLY"},
		{sql: "DESCRIBE KEYS", kind: ast.DescribeName, name: "KEYS"},
		{sql: "DESCRIBE TEXT", kind: ast.DescribeName, name: "TEXT"},
		{sql: `DESCRIBE "MyTable"`, kind: ast.DescribeName, name: "MyTable"},
		{sql: "DESCRIBE ks.t WITH INTERNALS", kind: ast.DescribeName, name: "ks.t", internals: true},
		{sql: "DESCRIBE KEYSPACE WITH INTERNALS", kind: ast.DescribeKeyspace, internals: true},
	}
	for _, tt := range tests {
		t.Run(tt.sql, func(t *testing.T) {
			d, ok := parseOne(t, tt.sql).(*ast.DescribeStmt)
			if !ok {
				t.Fatalf("got %T, want *ast.DescribeStmt", parseOne(t, tt.sql))
			}
			name := ""
			if d.Name != nil {
				for i, part := range d.Name.Parts {
					if i > 0 {
						name += "."
					}
					name += part.Name
				}
			}
			if d.Kind != tt.kind || name != tt.name || d.Full != tt.full || d.Only != tt.only || d.WithInternals != tt.internals {
				t.Fatalf("got kind=%d name=%q full=%v only=%v internals=%v", d.Kind, name, d.Full, d.Only, d.WithInternals)
			}
			if d.Loc.Start != 0 || d.Loc.End != len(tt.sql) {
				t.Fatalf("Loc = %+v, want {0 %d}", d.Loc, len(tt.sql))
			}
		})
	}
}

func TestParseGrammarDDL(t *testing.T) {
	cases := []struct {
		name  string
		sql   string
		check func(t *testing.T, n ast.Node)
	}{
		{"keyspace with any property", "CREATE KEYSPACE ks WITH replication = {'class': 'SimpleStrategy'} AND durable_writes = false AND x = KEYS", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateKeyspaceStmt)
			if len(s.Options) != 3 || s.Replication == nil || s.DurableWrites == nil || s.DurableWrites.Val {
				t.Fatalf("options=%d replication=%v durable=%v", len(s.Options), s.Replication, s.DurableWrites)
			}
		}},
		{"alter keyspace other property", "ALTER KEYSPACE ks WITH a = 'a'", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterKeyspaceStmt)
			if len(s.Options) != 1 || s.Replication != nil {
				t.Fatal("wrong options")
			}
		}},
		{"index on no column", "CREATE INDEX ON t()", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateIndexStmt)
			if s.Column != nil || len(s.Columns) != 0 {
				t.Fatal("want no target")
			}
		}},
		{"index on two targets", "CREATE INDEX ON t(a, KEYS(m)) WITH OPTIONS = {'k': 'v'} AND x = 1", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateIndexStmt)
			if len(s.Columns) != 2 || s.Column != s.Columns[0] || len(s.Properties) != 2 || s.Options == nil {
				t.Fatal("wrong targets or properties")
			}
			if fc, ok := s.Columns[1].(*ast.FunctionCall); !ok || fc.Name.Name != "KEYS" || len(fc.Args) != 1 {
				t.Fatalf("got %T", s.Columns[1])
			}
		}},
		{"index on a column named keys", "CREATE INDEX ON t(keys)", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.CreateIndexStmt).Column.(*ast.Identifier); !ok {
				t.Fatal("want an identifier")
			}
		}},
		{"aggregate with no arguments", "CREATE AGGREGATE a() SFUNC f STYPE int", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateAggregateStmt)
			if s.ParamType != nil || len(s.ParamTypes) != 0 {
				t.Fatal("want no argument types")
			}
		}},
		{"aggregate with two arguments", "CREATE AGGREGATE a(int, text) SFUNC f STYPE int", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateAggregateStmt)
			if len(s.ParamTypes) != 2 || s.ParamType != s.ParamTypes[0] {
				t.Fatal("want two argument types")
			}
		}},
		{"table with masks and a trailing comma", "CREATE TABLE t(a int PRIMARY KEY, b text MASKED WITH DEFAULT, c text MASKED WITH mask_inner(1, null),)", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateTableStmt)
			if len(s.Columns) != 3 || !s.Columns[1].Mask.Default || s.Columns[2].Mask.Function.Parts[0].Name != "mask_inner" || len(s.Columns[2].Mask.Args) != 2 {
				t.Fatal("wrong masks")
			}
		}},
		{"primary key among the columns", "CREATE TABLE t(a int, PRIMARY KEY(a), b int) WITH CLUSTERING ORDER BY (b ASC) AND CLUSTERING ORDER BY (c DESC)", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateTableStmt)
			if len(s.Columns) != 2 || s.PrimaryKey == nil || len(s.ClusteringOrders) != 2 {
				t.Fatal("wrong table")
			}
		}},
		{"compact as a property name", "CREATE TABLE t(a int PRIMARY KEY) WITH compact = 1 AND clustering = 'x'", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateTableStmt)
			if s.CompactStorage || len(s.Options) != 2 {
				t.Fatal("want two properties")
			}
		}},
		{"type with a trailing comma", "CREATE TYPE ks.addr(street text,)", func(t *testing.T, n ast.Node) {
			if len(n.(*ast.CreateTypeStmt).Fields) != 1 {
				t.Fatal("want one field")
			}
		}},
		{"alter column type", "ALTER TABLE t ALTER a TYPE text", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterTableStmt)
			if s.Op != ast.AlterTableAlterType || s.AlterColumn.Name != "a" || s.AlterType.Name.Name != "text" {
				t.Fatal("wrong ALTER TYPE")
			}
		}},
		{"alter column mask", "ALTER TABLE t ALTER IF EXISTS a MASKED WITH mask_null()", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterTableStmt)
			if s.Op != ast.AlterTableMask || !s.MaskIfExists || s.Mask == nil || s.DropMask {
				t.Fatal("wrong MASKED")
			}
		}},
		{"drop column mask", "ALTER TABLE t ALTER a DROP MASKED", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterTableStmt)
			if s.Op != ast.AlterTableMask || !s.DropMask || s.Mask != nil {
				t.Fatal("wrong DROP MASKED")
			}
		}},
		{"add columns in parentheses", "ALTER TABLE t ADD (a int STATIC, b text MASKED WITH DEFAULT)", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterTableStmt)
			if len(s.AddColumns) != 2 || !s.AddColumns[0].Static || s.AddColumns[1].Mask == nil {
				t.Fatal("wrong ADD")
			}
		}},
		{"drop columns with a timestamp", "ALTER TABLE t DROP (a, b) USING TIMESTAMP 10", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AlterTableStmt)
			if len(s.DropColumns) != 2 || s.DropTimestamp == nil {
				t.Fatal("wrong DROP")
			}
		}},
		{"drop a column named compact", "ALTER TABLE t DROP compact", func(t *testing.T, n ast.Node) {
			if n.(*ast.AlterTableStmt).Op != ast.AlterTableDrop {
				t.Fatal("want a column drop")
			}
		}},
		{"view without WHERE", "CREATE MATERIALIZED VIEW v AS SELECT a FROM t PRIMARY KEY(a)", nil},
		{"view with a relation", "CREATE MATERIALIZED VIEW v AS SELECT a FROM t WHERE a = 'x' AND b IS NOT NULL PRIMARY KEY(a, b)", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateMVStmt)
			if len(s.WhereRelations) != 1 || len(s.WhereNotNull) != 1 {
				t.Fatal("wrong WHERE")
			}
		}},
		{"frozen collection and qualified type", "CREATE TABLE t(a frozen<list<int>> PRIMARY KEY, b ks.addr, c frozen<addr>)", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateTableStmt)
			if s.Columns[1].Type.Keyspace == nil || s.Columns[1].Type.Keyspace.Name != "ks" {
				t.Fatal("want ks.addr")
			}
		}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			n := parseOne(t, tt.sql)
			if tt.check != nil {
				tt.check(t, n)
			}
		})
	}
}

func TestParseGrammarAuth(t *testing.T) {
	cases := []struct {
		name  string
		sql   string
		check func(t *testing.T, n ast.Node)
	}{
		{"LIST USERS", "LIST USERS", func(t *testing.T, n ast.Node) {
			if _, ok := n.(*ast.ListUsersStmt); !ok {
				t.Fatalf("got %T", n)
			}
		}},
		{"LIST two permissions", "LIST CREATE, SELECT PERMISSION, UNMASK NORECURSIVE", func(t *testing.T, n ast.Node) {
			s := n.(*ast.ListPermissionsStmt)
			if s.Privilege != "CREATE, SELECT, UNMASK" || len(s.Privileges) != 3 || !s.NoRecursive {
				t.Fatalf("got %q %v", s.Privilege, s.NoRecursive)
			}
		}},
		{"ALL TABLES IN KEYSPACE", "LIST ALL ON ALL TABLES IN KEYSPACE ks", func(t *testing.T, n ast.Node) {
			r := n.(*ast.ListPermissionsStmt).Resource
			if r.Type != "ALL TABLES IN KEYSPACE" || r.Name.Parts[0].Name != "ks" {
				t.Fatalf("got %s", r.Type)
			}
		}},
		{"a table named role", "LIST ALL ON ROLE", func(t *testing.T, n ast.Node) {
			r := n.(*ast.ListPermissionsStmt).Resource
			if r.Type != "TABLE" || r.Name.Parts[0].Name != "ROLE" {
				t.Fatalf("got %s", r.Type)
			}
		}},
		{"a table named mbeans", "LIST ALL ON MBEANS", func(t *testing.T, n ast.Node) {
			if n.(*ast.ListPermissionsStmt).Resource.Type != "TABLE" {
				t.Fatal("want a table")
			}
		}},
		{"MBEANS pattern", "GRANT SELECT ON MBEANS 'org.apache.cassandra.db:type=*' TO r", func(t *testing.T, n ast.Node) {
			if n.(*ast.GrantStmt).Resource.Type != "MBEANS" {
				t.Fatal("want MBEANS")
			}
		}},
		{"role named by a string", "CREATE ROLE 'Admin' WITH ACCESS FROM CIDRS {'a', 'b'} AND ACCESS TO ALL DATACENTERS", func(t *testing.T, n ast.Node) {
			s := n.(*ast.CreateRoleStmt)
			if s.Name.Name != "Admin" || !s.Name.Quoted || len(s.Options) != 2 || s.Options[0].Key != "ACCESS FROM CIDRS" {
				t.Fatal("wrong role")
			}
		}},
		{"ACCESS FROM ALL CIDRS", "ALTER ROLE r WITH ACCESS FROM ALL CIDRS", func(t *testing.T, n ast.Node) {
			if n.(*ast.AlterRoleStmt).Options[0].Key != "ACCESS FROM ALL CIDRS" {
				t.Fatal("wrong key")
			}
		}},
		{"user named by a string", "CREATE USER 'bob' WITH PASSWORD 'x'", func(t *testing.T, n ast.Node) {
			if n.(*ast.CreateUserStmt).Name.Name != "bob" {
				t.Fatal("wrong name")
			}
		}},
		{"ADD IDENTITY", "ADD IDENTITY IF NOT EXISTS 'spiffe://a' TO ROLE r", func(t *testing.T, n ast.Node) {
			s := n.(*ast.AddIdentityStmt)
			if !s.IfNotExists || s.Identity.Name != "spiffe://a" || s.Role.Name != "r" {
				t.Fatal("wrong ADD IDENTITY")
			}
		}},
		{"DROP IDENTITY", "DROP IDENTITY IF EXISTS 'spiffe://a'", func(t *testing.T, n ast.Node) {
			s := n.(*ast.DropIdentityStmt)
			if !s.IfExists || s.Identity.Name != "spiffe://a" {
				t.Fatal("wrong DROP IDENTITY")
			}
		}},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(t, parseOne(t, tt.sql))
		})
	}
}

// TestParseGrammarRejects holds Parse to what Cassandra refuses while it
// parses: forms Cql.g4 does not derive, and those an action of Parser.g
// refuses (testdata/grammar-overrides.tsv cites each kind).
func TestParseGrammarRejects(t *testing.T) {
	for _, sql := range []string{
		// Not derived.
		"CREATE TABLE t(a int) WITH CLUSTERING ORDER BY (a)",
		"CREATE TABLE t(a int) WITH CLUSTERING ORDER BY (ASC)",
		"CREATE INDEX ON t(VALUES())",
		"CREATE INDEX ON t(a) USING sai",
		"CREATE INDEX ON t(a) WITH {'k': 'v'}",
		"CREATE ROLE r WITH ACCESS TO DATACENTERS { }",
		"GRANT ALL PERMISSION ON TABLE t TO r",
		"ALTER TABLE t ADD a int, b int",
		"ALTER TABLE t DROP a, b",
		"ALTER TYPE addr ADD a int, b int",
		"ALTER MATERIALIZED VIEW v",
		"CREATE OR REPLACE TABLE t(a int PRIMARY KEY)",
		"CREATE KEYSPACE ks WITH a = b",
		"CREATE KEYSPACE ks WITH a = NULL",
		"CREATE KEYSPACE ks WITH replication = {'class': NULL}",
		"CREATE KEYSPACE ks WITH replication = {class: 'x'}",
		"CREATE TABLE t(a MASKED<int, int>)",
		"CREATE USER \"bob\"",
		"DESCRIBE FULL",
		"DESCRIBE TABLES.t",
		"DESCRIBE SELECT",
		"DESCRIBE KEYSPACE ks.t",
		// Refused by an action while parsing.
		"CREATE TABLE t(a int, a text)",
		"CREATE TABLE t(a int, A text)",
		"CREATE TABLE t(a int PRIMARY KEY, PRIMARY KEY(a))",
		"CREATE TABLE t(a int) WITH CLUSTERING ORDER BY (a ASC, a DESC)",
		"CREATE TABLE t(a int) WITH CLUSTERING ORDER BY (a ASC) AND CLUSTERING ORDER BY (a DESC)",
		"CREATE TABLE t(a int) WITH x = 1 AND X = 2",
		"CREATE MATERIALIZED VIEW v AS SELECT a FROM t PRIMARY KEY(a) WITH x = 1 AND x = 1",
		"CREATE KEYSPACE ks WITH replication = {'a': 'b', 'a': 'c'}",
		"CREATE ROLE r WITH LOGIN = true AND LOGIN = false",
		"CREATE ROLE r WITH PASSWORD = 'a' AND HASHED PASSWORD = 'b'",
		"CREATE USER u WITH PASSWORD 'a' AND HASHED PASSWORD 'b'",
		"CREATE ROLE r WITH ACCESS TO ALL DATACENTERS AND ACCESS TO DATACENTERS {'dc1'}",
		"CREATE ROLE r WITH ACCESS FROM CIDRS {'a'} AND ACCESS FROM ALL CIDRS",
		"ALTER TABLE t ADD a frozen<int>",
		"CREATE TABLE t(a frozen<frozen<text>>)",
		"CREATE TABLE t(a byte)",
		"CREATE TYPE enum(a int)",
		"LIST ALL ON MBEAN 'a'",
		"LIST ALL ON MBEAN ''",
		"LIST ALL ON MBEAN 'd:'",
		"LIST ALL ON MBEAN 'd:a=1,a=2'",
		"LIST ALL ON MBEANS '*:*'",
		"LIST ALL ON MBEANS 'org.apache.cassandra.db:type=Tables'",
	} {
		t.Run(sql, func(t *testing.T) {
			if _, err := Parse(sql); err == nil {
				t.Fatalf("Parse(%q) succeeded; Cassandra refuses it", sql)
			}
		})
	}
}
