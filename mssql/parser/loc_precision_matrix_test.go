package parser

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/bytebase/omni/mssql/ast"
)

type locMatrixCase struct {
	name     string
	sql      string
	typeName string
	index    int
	want     string
}

var locPrecisionMatrixExactSliceCases = []locMatrixCase{
	// Universal statement/source invariants.
	{name: "leading whitespace select stmt", sql: "  \n\tSELECT 1", typeName: "SelectStmt", want: "SELECT 1"},
	{name: "statement semicolon excluded from ast loc", sql: "SELECT 1;", typeName: "SelectStmt", want: "SELECT 1"},
	{name: "multiline select stmt", sql: "SELECT\n  1", typeName: "SelectStmt", want: "SELECT\n  1"},
	{name: "unicode byte range select stmt", sql: "SELECT N'你好'", typeName: "SelectStmt", want: "SELECT N'你好'"},

	// Comparison and predicates.
	{name: "not in", sql: "SELECT c NOT IN (1, 2)", typeName: "InExpr", want: "c NOT IN (1, 2)"},
	{name: "not between", sql: "SELECT c NOT BETWEEN 1 AND 2", typeName: "BetweenExpr", want: "c NOT BETWEEN 1 AND 2"},
	{name: "not like escape", sql: "SELECT c NOT LIKE 'x%' ESCAPE '\\'", typeName: "LikeExpr", want: "c NOT LIKE 'x%' ESCAPE '\\'"},
	{name: "equals any subquery", sql: "SELECT * FROM t WHERE c = ANY (SELECT x FROM s)", typeName: "SubqueryComparisonExpr", want: "c = ANY (SELECT x FROM s)"},
	{name: "in subquery", sql: "SELECT c IN (SELECT x FROM s)", typeName: "InExpr", want: "c IN (SELECT x FROM s)"},

	// Boolean and arithmetic.
	{name: "or and precedence", sql: "SELECT a OR b AND c", typeName: "BinaryExpr", want: "a OR b AND c"},
	{name: "add subtract chain", sql: "SELECT a - b + c", typeName: "BinaryExpr", want: "a - b + c"},
	{name: "divide modulo chain", sql: "SELECT a / b % c", typeName: "BinaryExpr", want: "a / b % c"},
	{name: "unary minus", sql: "SELECT -c", typeName: "UnaryExpr", want: "-c"},
	{name: "unary not", sql: "SELECT NOT (c > 0)", typeName: "UnaryExpr", want: "NOT (c > 0)"},
	{name: "paren expr", sql: "SELECT (a + b)", typeName: "ParenExpr", want: "(a + b)"},

	// Postfix expressions.
	{name: "collate database_default", sql: "SELECT name COLLATE database_default", typeName: "CollateExpr", want: "name COLLATE database_default"},
	{name: "at time zone chain", sql: "SELECT dt AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'", typeName: "AtTimeZoneExpr", want: "dt AT TIME ZONE 'UTC' AT TIME ZONE 'Pacific Standard Time'"},

	// Functions and built-ins.
	{name: "count distinct", sql: "SELECT COUNT(DISTINCT c)", typeName: "FuncCallExpr", want: "COUNT(DISTINCT c)"},
	{name: "cast", sql: "SELECT CAST(c AS int)", typeName: "CastExpr", want: "CAST(c AS int)"},
	{name: "convert", sql: "SELECT CONVERT(varchar(20), c, 126)", typeName: "ConvertExpr", want: "CONVERT(varchar(20), c, 126)"},
	{name: "try cast", sql: "SELECT TRY_CAST(c AS int)", typeName: "TryCastExpr", want: "TRY_CAST(c AS int)"},
	{name: "try convert", sql: "SELECT TRY_CONVERT(int, c)", typeName: "TryConvertExpr", want: "TRY_CONVERT(int, c)"},
	{name: "coalesce", sql: "SELECT COALESCE(a, b, c)", typeName: "CoalesceExpr", want: "COALESCE(a, b, c)"},
	{name: "nullif", sql: "SELECT NULLIF(a, b)", typeName: "NullifExpr", want: "NULLIF(a, b)"},
	{name: "iif", sql: "SELECT IIF(c > 0, 1, 0)", typeName: "IifExpr", want: "IIF(c > 0, 1, 0)"},

	// CASE, subquery, and specialty expressions.
	{name: "case expr", sql: "SELECT CASE WHEN c > 0 THEN 1 ELSE 0 END", typeName: "CaseExpr", want: "CASE WHEN c > 0 THEN 1 ELSE 0 END"},
	{name: "case when", sql: "SELECT CASE WHEN c > 0 THEN 1 END", typeName: "CaseWhen", want: "WHEN c > 0 THEN 1"},
	{name: "exists", sql: "SELECT * FROM t WHERE EXISTS (SELECT 1 FROM s)", typeName: "ExistsExpr", want: "EXISTS (SELECT 1 FROM s)"},
	{name: "scalar subquery", sql: "SELECT (SELECT 1)", typeName: "SubqueryExpr", want: "(SELECT 1)"},
	{name: "current of", sql: "DELETE FROM t WHERE CURRENT OF cur", typeName: "CurrentOfExpr", want: "CURRENT OF cur"},
	{name: "method call", sql: "SELECT geography::Point(1, 2, 4326)", typeName: "MethodCallExpr", want: "geography::Point(1, 2, 4326)"},
	{name: "grouping sets", sql: "SELECT SUM(c) FROM t GROUP BY GROUPING SETS ((a), (b))", typeName: "GroupingSetsExpr", want: "GROUPING SETS ((a), (b))"},
	{name: "rollup", sql: "SELECT SUM(c) FROM t GROUP BY ROLLUP (a, b)", typeName: "RollupExpr", want: "ROLLUP (a, b)"},
	{name: "cube", sql: "SELECT SUM(c) FROM t GROUP BY CUBE (a, b)", typeName: "CubeExpr", want: "CUBE (a, b)"},

	// Column and variable references.
	{name: "four part column", sql: "SELECT db.dbo.t.c", typeName: "ColumnRef", want: "db.dbo.t.c"},
	{name: "five part column", sql: "SELECT server.db.dbo.t.c", typeName: "ColumnRef", want: "server.db.dbo.t.c"},
	{name: "variable ref", sql: "SELECT @v", typeName: "VariableRef", want: "@v"},
	{name: "system variable ref", sql: "SELECT @@ROWCOUNT", typeName: "VariableRef", want: "@@ROWCOUNT"},

	// Table references.
	{name: "table ref", sql: "SELECT * FROM t", typeName: "TableRef", want: "t"},
	{name: "schema table ref", sql: "SELECT * FROM dbo.t", typeName: "TableRef", want: "dbo.t"},
	{name: "database schema table ref", sql: "SELECT * FROM db.dbo.t", typeName: "TableRef", want: "db.dbo.t"},
	{name: "table variable ref", sql: "SELECT * FROM @t", typeName: "TableVarRef", want: "@t"},
	{name: "table variable method ref", sql: "SELECT * FROM @x.nodes('/r') n(c)", typeName: "TableVarMethodCallRef", want: "@x.nodes('/r') n(c)"},
	{name: "aliased table ref", sql: "SELECT * FROM t AS x", typeName: "TableRef", want: "t AS x"},
	{name: "openrowset table func", sql: "SELECT * FROM OPENROWSET('SQLNCLI', 'server=srv', 'SELECT 1') AS r", typeName: "FuncCallExpr", want: "OPENROWSET('SQLNCLI', 'server=srv', 'SELECT 1')"},

	// Data types.
	{name: "datatype int", sql: "CREATE TABLE t (c int)", typeName: "DataType", want: "int"},
	{name: "datatype varchar length", sql: "CREATE TABLE t (c varchar(50))", typeName: "DataType", want: "varchar(50)"},
	{name: "datatype decimal", sql: "CREATE TABLE t (c decimal(10, 2))", typeName: "DataType", want: "decimal(10, 2)"},
	{name: "datatype nvarchar max", sql: "CREATE TABLE t (c nvarchar(max))", typeName: "DataType", want: "nvarchar(max)"},
	{name: "datatype schema type", sql: "CREATE TABLE t (c dbo.MyType)", typeName: "DataType", want: "dbo.MyType"},

	// CREATE TABLE metadata.
	{name: "column def", sql: "CREATE TABLE t (c int)", typeName: "ColumnDef", want: "c int"},
	{name: "nullable null", sql: "CREATE TABLE t (c int NULL)", typeName: "NullableSpec", want: "NULL"},
	{name: "nullable not null", sql: "CREATE TABLE t (c int NOT NULL)", typeName: "NullableSpec", want: "NOT NULL"},
	{name: "identity", sql: "CREATE TABLE t (id int IDENTITY(1,1))", typeName: "IdentitySpec", want: "IDENTITY(1,1)"},
	{name: "computed expr", sql: "CREATE TABLE t (c AS (a + b))", typeName: "BinaryExpr", want: "a + b"},
	{name: "computed def persisted", sql: "CREATE TABLE t (c AS (a + b) PERSISTED)", typeName: "ComputedColumnDef", want: "AS (a + b) PERSISTED"},
	{name: "default constraint expr", sql: "CREATE TABLE t (c int CONSTRAINT df DEFAULT (0))", typeName: "Literal", want: "0"},
	{name: "check constraint between expr", sql: "CREATE TABLE t (c int CONSTRAINT ck CHECK (c BETWEEN 1 AND 10))", typeName: "BetweenExpr", want: "c BETWEEN 1 AND 10"},
	{name: "primary key constraint", sql: "CREATE TABLE t (c int, CONSTRAINT pk PRIMARY KEY (c))", typeName: "ConstraintDef", want: "CONSTRAINT pk PRIMARY KEY (c)"},
	{name: "unique constraint", sql: "CREATE TABLE t (c int, CONSTRAINT uq UNIQUE (c))", typeName: "ConstraintDef", want: "CONSTRAINT uq UNIQUE (c)"},
	{name: "foreign key constraint", sql: "CREATE TABLE t (c int, FOREIGN KEY (c) REFERENCES p(id))", typeName: "ConstraintDef", want: "FOREIGN KEY (c) REFERENCES p(id)"},
	{name: "edge constraint", sql: "CREATE TABLE edge_t AS EDGE (CONSTRAINT ec CONNECTION (dbo.a TO dbo.b))", typeName: "ConstraintDef", want: "CONSTRAINT ec CONNECTION (dbo.a TO dbo.b)"},
	{name: "memory optimized option", sql: "CREATE TABLE t (c int) WITH (MEMORY_OPTIMIZED = ON)", typeName: "TableOption", want: "MEMORY_OPTIMIZED = ON"},
	{name: "system versioning option", sql: "CREATE TABLE t (id int, valid_from datetime2 GENERATED ALWAYS AS ROW START, valid_to datetime2 GENERATED ALWAYS AS ROW END, PERIOD FOR SYSTEM_TIME (valid_from, valid_to)) WITH (SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.t_history))", typeName: "TableOption", want: "SYSTEM_VERSIONING = ON (HISTORY_TABLE = dbo.t_history)"},
	{name: "inline index", sql: "CREATE TABLE t (c int, INDEX ix (c))", typeName: "InlineIndexDef", want: "INDEX ix (c)"},
	{name: "inline index column", sql: "CREATE TABLE t (c int, INDEX ix (c DESC))", typeName: "IndexColumn", want: "c DESC"},

	// Index, spatial, fulltext, and search metadata.
	{name: "create index stmt", sql: "CREATE INDEX ix ON t (c)", typeName: "CreateIndexStmt", want: "CREATE INDEX ix ON t (c)"},
	{name: "create index column desc", sql: "CREATE INDEX ix ON t (c DESC)", typeName: "IndexColumn", want: "c DESC"},
	{name: "create index include column", sql: "CREATE INDEX ix ON t (c) INCLUDE (d)", typeName: "IndexColumn", index: 1, want: "d"},
	{name: "create index filter", sql: "CREATE INDEX ix ON t (c) WHERE c > 0", typeName: "BinaryExpr", want: "c > 0"},
	{name: "spatial index column", sql: "CREATE SPATIAL INDEX six ON dbo.t (geom) USING GEOMETRY_GRID", typeName: "ColumnRef", want: "geom"},
	{name: "spatial geometry grid stmt", sql: "CREATE SPATIAL INDEX six ON dbo.t (geom) USING GEOMETRY_GRID", typeName: "CreateSpatialIndexStmt", want: "CREATE SPATIAL INDEX six ON dbo.t (geom) USING GEOMETRY_GRID"},
	{name: "spatial geography auto grid stmt", sql: "CREATE SPATIAL INDEX six ON dbo.t (geom) USING GEOGRAPHY_AUTO_GRID", typeName: "CreateSpatialIndexStmt", want: "CREATE SPATIAL INDEX six ON dbo.t (geom) USING GEOGRAPHY_AUTO_GRID"},
	{name: "create fulltext stmt", sql: "CREATE FULLTEXT INDEX ON t(c) KEY INDEX pk", typeName: "CreateFulltextIndexStmt", want: "CREATE FULLTEXT INDEX ON t(c) KEY INDEX pk"},
	{name: "fulltext indexed column", sql: "CREATE FULLTEXT INDEX ON t(c) KEY INDEX pk", typeName: "ColumnRef", want: "c"},
	{name: "contains predicate", sql: "SELECT * FROM t WHERE CONTAINS(c, 'term')", typeName: "FullTextPredicate", want: "CONTAINS(c, 'term')"},
	{name: "freetext predicate", sql: "SELECT * FROM t WHERE FREETEXT(c, 'term')", typeName: "FullTextPredicate", want: "FREETEXT(c, 'term')"},
	{name: "containstable ref", sql: "SELECT * FROM CONTAINSTABLE(t, c, 'term') ct", typeName: "FullTextTableRef", want: "CONTAINSTABLE(t, c, 'term') ct"},
	{name: "semantic similarity ref", sql: "SELECT * FROM SEMANTICSIMILARITYTABLE(t, c, source_key) s", typeName: "SemanticTableRef", want: "SEMANTICSIMILARITYTABLE(t, c, source_key) s"},

	// SELECT clauses.
	{name: "res target alias", sql: "SELECT a AS x", typeName: "ResTarget", want: "a AS x"},
	{name: "select assign", sql: "SELECT @v = a", typeName: "SelectAssign", want: "@v = a"},
	{name: "top clause", sql: "SELECT TOP (10) a FROM t", typeName: "TopClause", want: "TOP (10)"},
	{name: "where expression", sql: "SELECT a FROM t WHERE a > 0", typeName: "BinaryExpr", want: "a > 0"},
	{name: "order by item", sql: "SELECT a FROM t ORDER BY a DESC", typeName: "OrderByItem", want: "a DESC"},
	{name: "fetch clause", sql: "SELECT a FROM t ORDER BY a OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY", typeName: "FetchClause", want: "FETCH NEXT 5 ROWS ONLY"},
	{name: "for xml clause", sql: "SELECT a FROM t FOR XML PATH", typeName: "ForClause", want: "FOR XML PATH"},
	{name: "common table expr", sql: "WITH cte AS (SELECT 1) SELECT * FROM cte", typeName: "CommonTableExpr", want: "cte AS (SELECT 1)"},
	{name: "xml namespace decl", sql: "WITH XMLNAMESPACES ('http://example.com' AS ns) SELECT 1", typeName: "XmlNamespaceDecl", want: "'http://example.com' AS ns"},

	// JOIN, WINDOW, GROUPING.
	{name: "join clause", sql: "SELECT * FROM a JOIN b ON a.id = b.id", typeName: "JoinClause", want: "a JOIN b ON a.id = b.id"},
	{name: "cross apply clause", sql: "SELECT * FROM a CROSS APPLY fn(a.id) f", typeName: "JoinClause", want: "a CROSS APPLY fn(a.id) f"},
	{name: "over full clause", sql: "SELECT SUM(x) OVER (PARTITION BY y ORDER BY z ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) FROM t", typeName: "OverClause", want: "OVER (PARTITION BY y ORDER BY z ROWS BETWEEN 1 PRECEDING AND CURRENT ROW)"},
	{name: "window frame", sql: "SELECT SUM(x) OVER (ORDER BY z ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) FROM t", typeName: "WindowFrame", want: "ROWS BETWEEN 1 PRECEDING AND CURRENT ROW"},
	{name: "window bound", sql: "SELECT SUM(x) OVER (ORDER BY z ROWS BETWEEN 1 PRECEDING AND CURRENT ROW) FROM t", typeName: "WindowBound", want: "1 PRECEDING"},
	{name: "window def", sql: "SELECT SUM(x) OVER w FROM t WINDOW w AS (PARTITION BY y)", typeName: "WindowDef", want: "w AS (PARTITION BY y)"},

	// DML.
	{name: "insert stmt", sql: "INSERT INTO t (a) VALUES (1)", typeName: "InsertStmt", want: "INSERT INTO t (a) VALUES (1)"},
	{name: "values clause", sql: "INSERT INTO t (a) VALUES (1), (2)", typeName: "ValuesClause", want: "VALUES (1), (2)"},
	{name: "update stmt", sql: "UPDATE t SET a = 1 WHERE b = 2", typeName: "UpdateStmt", want: "UPDATE t SET a = 1 WHERE b = 2"},
	{name: "delete stmt", sql: "DELETE FROM t WHERE a = 1", typeName: "DeleteStmt", want: "DELETE FROM t WHERE a = 1"},
	{name: "output clause", sql: "UPDATE t SET a = 1 OUTPUT inserted.a INTO audit WHERE b = 2", typeName: "OutputClause", want: "OUTPUT inserted.a INTO audit"},
	{name: "merge when matched", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a;", typeName: "MergeWhenClause", want: "WHEN MATCHED THEN UPDATE SET a = s.a"},
	{name: "merge update action", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN UPDATE SET a = s.a;", typeName: "MergeUpdateAction", want: "UPDATE SET a = s.a"},
	{name: "merge insert action", sql: "MERGE INTO t USING s ON t.id = s.id WHEN NOT MATCHED THEN INSERT (a) VALUES (s.a);", typeName: "MergeInsertAction", want: "INSERT (a) VALUES (s.a)"},
	{name: "merge delete action", sql: "MERGE INTO t USING s ON t.id = s.id WHEN MATCHED THEN DELETE;", typeName: "MergeDeleteAction", want: "DELETE"},

	// Top-level statement families.
	{name: "ddl statement", sql: "CREATE VIEW v AS SELECT 1", typeName: "CreateViewStmt", want: "CREATE VIEW v AS SELECT 1"},
	{name: "security statement", sql: "CREATE LOGIN testLogin WITH PASSWORD = 'P@ssw0rd'", typeName: "SecurityStmt", want: "CREATE LOGIN testLogin WITH PASSWORD = 'P@ssw0rd'"},
	{name: "admin statement", sql: "DBCC CHECKDB", typeName: "DbccStmt", want: "DBCC CHECKDB"},
}

func TestMSSQLLocPrecisionMatrix_ExactSlices(t *testing.T) {
	for _, tt := range locPrecisionMatrixExactSliceCases {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse failure: %v", err)
			}
			node := nthNodeByType(list, tt.typeName, tt.index)
			if node == nil {
				t.Fatalf("node %s[%d] not found; available types: %s", tt.typeName, tt.index, strings.Join(nodeTypeNames(list), ", "))
			}
			loc := ast.NodeLoc(node)
			if !loc.IsValid() {
				t.Fatalf("%s loc invalid: %+v", tt.typeName, loc)
			}
			if loc.Start > len(tt.sql) || loc.End > len(tt.sql) || loc.Start > loc.End {
				t.Fatalf("%s loc out of range: %+v for SQL length %d", tt.typeName, loc, len(tt.sql))
			}
			if got := tt.sql[loc.Start:loc.End]; got != tt.want {
				t.Fatalf("%s slice = %q, want %q (loc=%+v)", tt.typeName, got, tt.want, loc)
			}
		})
	}
}

func TestMSSQLLocPrecisionMatrix_OptionText(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{
			name: "create index fillfactor",
			sql:  "CREATE INDEX ix ON t (c) WITH (FILLFACTOR = 90)",
			want: []string{"FILLFACTOR=90"},
		},
		{
			name: "alter index online",
			sql:  "ALTER INDEX ix ON t REBUILD WITH (ONLINE = ON)",
			want: []string{"ONLINE=ON"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse failure: %v", err)
			}
			var got []string
			ast.Inspect(list, func(n ast.Node) bool {
				if s, ok := n.(*ast.String); ok {
					got = append(got, s.Str)
				}
				return true
			})
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("option strings = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestMSSQLLocPrecisionMatrix_ChildLocsWithinParent(t *testing.T) {
	for _, tt := range locPrecisionMatrixExactSliceCases {
		t.Run(tt.name, func(t *testing.T) {
			list, err := Parse(tt.sql)
			if err != nil {
				t.Fatalf("parse failure: %v", err)
			}
			assertChildLocsWithinParent(t, tt.sql, list)
		})
	}
}

func nthNodeByType(root ast.Node, typeName string, index int) ast.Node {
	var matches []ast.Node
	ast.Inspect(root, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		t := reflect.TypeOf(n)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Name() == typeName {
			matches = append(matches, n)
		}
		return true
	})
	if len(matches) == 0 {
		return nil
	}
	if index < 0 {
		index = len(matches) + index
	}
	if index < 0 || index >= len(matches) {
		return nil
	}
	return matches[index]
}

func nodeTypeNames(root ast.Node) []string {
	seen := make(map[string]bool)
	ast.Inspect(root, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		t := reflect.TypeOf(n)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		seen[t.Name()] = true
		return true
	})
	var names []string
	for name := range seen {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func assertChildLocsWithinParent(t *testing.T, sql string, root ast.Node) {
	t.Helper()

	type frame struct {
		node ast.Node
		loc  ast.Loc
	}
	var stack []frame
	ast.Walk(locContainmentVisitor{
		enter: func(n ast.Node) {
			loc := ast.NodeLoc(n)
			if !loc.IsValid() {
				stack = append(stack, frame{node: n, loc: loc})
				return
			}
			if loc.Start > len(sql) || loc.End > len(sql) || loc.Start > loc.End {
				t.Errorf("%T loc out of range: %+v", n, loc)
			}
			for i := len(stack) - 1; i >= 0; i-- {
				parent := stack[i]
				if !parent.loc.IsValid() {
					continue
				}
				if loc.Start < parent.loc.Start || loc.End > parent.loc.End {
					t.Errorf("%T loc %+v (%q) outside parent %T loc %+v (%q)",
						n, loc, locSliceForMessage(sql, loc), parent.node, parent.loc, locSliceForMessage(sql, parent.loc))
				}
				break
			}
			stack = append(stack, frame{node: n, loc: loc})
		},
		leave: func() {
			stack = stack[:len(stack)-1]
		},
	}, root)
}

func locSliceForMessage(sql string, loc ast.Loc) string {
	if !loc.IsValid() || loc.Start > len(sql) || loc.End > len(sql) || loc.Start > loc.End {
		return fmt.Sprintf("<bad %+v>", loc)
	}
	return sql[loc.Start:loc.End]
}

type locContainmentVisitor struct {
	enter func(ast.Node)
	leave func()
}

func (v locContainmentVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		v.leave()
		return nil
	}
	v.enter(n)
	return v
}
