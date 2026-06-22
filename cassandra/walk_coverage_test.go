package cassandra

import (
	"reflect"
	"testing"

	"github.com/bytebase/omni/cassandra/ast"
	"github.com/bytebase/omni/cassandra/parser"
)

func TestWalkCoversAllChildren(t *testing.T) {
	tests := []struct {
		name string
		sql  string
	}{
		{"CREATE TABLE complex", "CREATE TABLE t (id int, name text, age int, PRIMARY KEY ((id, name), age)) WITH CLUSTERING ORDER BY (age DESC) AND comment = 'test'"},
		{"CREATE MV", "CREATE MATERIALIZED VIEW mv AS SELECT col1, col2 FROM t WHERE col1 IS NOT NULL AND col2 IS NOT NULL PRIMARY KEY (col1, col2)"},
		{"CREATE ROLE WITH", "CREATE ROLE myrole WITH PASSWORD = 'secret' AND LOGIN = true AND SUPERUSER = false"},
		{"SELECT complex", "SELECT DISTINCT name AS n FROM ks.users WHERE id = 1 ORDER BY name ASC LIMIT 10 ALLOW FILTERING"},
		{"UPDATE with IF", "UPDATE users USING TTL 3600 SET name = 'Bob' WHERE id = 2 IF name = 'old'"},
		{"INSERT JSON", "INSERT INTO users JSON '{\"id\": 1}' DEFAULT UNSET IF NOT EXISTS USING TTL 86400"},
		{"DELETE complex", "DELETE name FROM ks.users WHERE id = 2 IF EXISTS"},
		{"BATCH", "BEGIN UNLOGGED BATCH USING TIMESTAMP 12345 INSERT INTO t (id) VALUES (1); DELETE FROM t WHERE id = 2; APPLY BATCH"},
		{"GRANT", "GRANT SELECT ON TABLE users TO reader"},
		{"CREATE KEYSPACE", "CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'} AND DURABLE_WRITES = true"},
		{"CREATE INDEX", "CREATE INDEX idx ON users (name)"},
		{"CREATE TYPE", "CREATE TYPE address (street text, city text, zip int)"},
		{"ALTER TYPE RENAME", "ALTER TYPE address RENAME street TO road AND city TO town"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, err := parser.Parse(tt.sql)
			if err != nil {
				t.Fatal(err)
			}

			reflectTypes := make(map[string]bool)
			for _, item := range list.Items {
				raw := item.(*ast.RawStmt)
				collectNodeTypes(reflect.ValueOf(raw.Stmt), reflectTypes)
			}

			walkTypes := make(map[string]bool)
			for _, item := range list.Items {
				raw := item.(*ast.RawStmt)
				ast.Inspect(raw.Stmt, func(n ast.Node) bool {
					if n == nil {
						return false
					}
					walkTypes[reflect.TypeOf(n).Elem().Name()] = true
					return true
				})
			}

			for typeName := range reflectTypes {
				if !walkTypes[typeName] {
					t.Errorf("reflection found Node type %s but ast.Walk did not visit it", typeName)
				}
			}
		})
	}
}

var nodeIface = reflect.TypeOf((*ast.Node)(nil)).Elem()

func collectNodeTypes(v reflect.Value, types map[string]bool) {
	switch v.Kind() {
	case reflect.Ptr:
		if v.IsNil() {
			return
		}
		if v.Type().Implements(nodeIface) {
			types[v.Type().Elem().Name()] = true
		}
		collectNodeTypes(v.Elem(), types)
	case reflect.Interface:
		if v.IsNil() {
			return
		}
		elem := v.Elem()
		if elem.Type().Implements(nodeIface) {
			typeName := elem.Type().Elem().Name()
			types[typeName] = true
		}
		collectNodeTypes(elem, types)
	case reflect.Struct:
		t := v.Type()
		for i := range t.NumField() {
			f := t.Field(i)
			if !f.IsExported() || f.Name == "Loc" {
				continue
			}
			collectNodeTypes(v.Field(i), types)
		}
	case reflect.Slice:
		for i := range v.Len() {
			collectNodeTypes(v.Index(i), types)
		}
	}
}

func TestWalkChildrenCoverage(t *testing.T) {
	walkSrc := "/Users/rebeliceyang/.slock/agents/3ddcbf07-d89e-4266-aa7b-12b9445bc400/omni-cassandra/cassandra/ast/walk_children.go"
	_ = walkSrc

	nodeTypes := collectAllNodeStructTypes()

	walkCases := make(map[string]bool)
	for _, sql := range []string{
		"SELECT * FROM users",
		"SELECT name AS n FROM ks.users WHERE id = 1 ORDER BY name ASC LIMIT 10 ALLOW FILTERING",
		"SELECT count(*) FROM users WHERE id IN (1, 2, 3)",
		"SELECT * FROM users WHERE tags CONTAINS 'admin'",
		"SELECT * FROM users WHERE tags CONTAINS KEY 'role'",
		"SELECT * FROM users WHERE (a, b) = (1, 2)",
		"SELECT * FROM users WHERE (a, b) IN ((1, 2), (3, 4))",
		"INSERT INTO t (id) VALUES (1)",
		"INSERT INTO t (id, name) VALUES (1, 'Alice') USING TTL 86400 AND TIMESTAMP 12345",
		"INSERT INTO t (id, f, b, n) VALUES (1, 3.14, true, null)",
		"INSERT INTO t (id, u) VALUES (1, 550e8400-e29b-41d4-a716-446655440000)",
		"INSERT INTO t (id, h) VALUES (1, 0xDEADBEEF)",
		"INSERT INTO t (id, m) VALUES (1, {'key': 'val'})",
		"INSERT INTO t (id, s) VALUES (1, {'a', 'b'})",
		"INSERT INTO t (id, l) VALUES (1, ['x', 'y'])",
		"INSERT INTO t (id, t2) VALUES (1, (1, 'a', true))",
		"UPDATE t SET x = 1 WHERE id = 1",
		"UPDATE t SET x = 1 WHERE id = 1 IF x = 0",
		"UPDATE t SET m['key'] = 'val' WHERE id = 1",
		"DELETE FROM t WHERE id = 1",
		"DELETE FROM t WHERE id = 1 IF EXISTS",
		"BEGIN BATCH INSERT INTO t (id) VALUES (1); APPLY BATCH",
		"TRUNCATE t",
		"USE ks",
		"SELECT token(id) FROM t WHERE token(id) > token(1)",
		"CREATE KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'} AND DURABLE_WRITES = true",
		"ALTER KEYSPACE ks WITH REPLICATION = {'class': 'SimpleStrategy', 'replication_factor': '1'}",
		"DROP KEYSPACE ks",
		"CREATE TABLE t (id int, name text, age int, PRIMARY KEY ((id, name), age)) WITH CLUSTERING ORDER BY (age DESC) AND comment = 'test'",
		"ALTER TABLE t ADD col text",
		"DROP TABLE t",
		"CREATE INDEX ON t (col)",
		"DROP INDEX idx",
		"CREATE TYPE mytype (f1 text)",
		"ALTER TYPE mytype ADD f2 int",
		"ALTER TYPE mytype RENAME f1 TO field1 AND f2 TO field2",
		"DROP TYPE mytype",
		"CREATE MATERIALIZED VIEW mv AS SELECT * FROM t WHERE id IS NOT NULL PRIMARY KEY (id)",
		"ALTER MATERIALIZED VIEW mv WITH comment = 'test'",
		"DROP MATERIALIZED VIEW mv",
		"CREATE FUNCTION ks.f(input text) CALLED ON NULL INPUT RETURNS text LANGUAGE java AS $$return input;$$",
		"DROP FUNCTION IF EXISTS ks.f",
		"CREATE AGGREGATE ks.agg(int) SFUNC plus STYPE int FINALFUNC fin INITCOND 0",
		"DROP AGGREGATE IF EXISTS ks.agg",
		"CREATE TRIGGER tr ON t USING 'org.example.Trigger'",
		"DROP TRIGGER tr ON t",
		"GRANT SELECT ON TABLE t TO r",
		"REVOKE ALL ON ALL KEYSPACES FROM r",
		"LIST ALL PERMISSIONS",
		"LIST ROLES",
		"CREATE ROLE r",
		"ALTER ROLE r WITH PASSWORD = 'x'",
		"DROP ROLE r",
		"CREATE USER u WITH PASSWORD 'x'",
		"ALTER USER u WITH PASSWORD 'y'",
		"DROP USER u",
		"SELECT ks.* FROM ks.t",
		"SELECT * FROM t ORDER BY v ANN OF [1.0, 2.0, 3.0] LIMIT 5",
	} {
		list, err := parser.Parse(sql)
		if err != nil {
			continue
		}
		for _, item := range list.Items {
			raw := item.(*ast.RawStmt)
			ast.Inspect(raw.Stmt, func(n ast.Node) bool {
				if n == nil {
					return false
				}
				walkCases[reflect.TypeOf(n).Elem().Name()] = true
				return true
			})
		}
	}

	var missing []string
	for _, nt := range nodeTypes {
		if nt == "List" || nt == "RawStmt" {
			continue
		}
		if !walkCases[nt] {
			missing = append(missing, nt)
		}
	}
	if len(missing) > 0 {
		t.Errorf("AST node types not reached by Walk in any test SQL: %v", missing)
		t.Log("Add test SQL that exercises these node types, or add cases to walkChildren")
	}
}

func collectAllNodeStructTypes() []string {
	nodeType := reflect.TypeOf((*ast.Node)(nil)).Elem()
	candidates := []reflect.Type{
		reflect.TypeOf(ast.Identifier{}),
		reflect.TypeOf(ast.QualifiedName{}),
		reflect.TypeOf(ast.StringLit{}),
		reflect.TypeOf(ast.IntegerLit{}),
		reflect.TypeOf(ast.FloatLit{}),
		reflect.TypeOf(ast.BoolLit{}),
		reflect.TypeOf(ast.NullLit{}),
		reflect.TypeOf(ast.UUIDLit{}),
		reflect.TypeOf(ast.HexLit{}),
		reflect.TypeOf(ast.CodeBlock{}),
		reflect.TypeOf(ast.StarExpr{}),
		reflect.TypeOf(ast.MapLit{}),
		reflect.TypeOf(ast.SetLit{}),
		reflect.TypeOf(ast.ListLit{}),
		reflect.TypeOf(ast.TupleLit{}),
		reflect.TypeOf(ast.VectorLit{}),
		reflect.TypeOf(ast.FunctionCall{}),
		reflect.TypeOf(ast.BinaryExpr{}),
		reflect.TypeOf(ast.InExpr{}),
		reflect.TypeOf(ast.ContainsExpr{}),
		reflect.TypeOf(ast.TupleCompareExpr{}),
		reflect.TypeOf(ast.TupleInExpr{}),
		reflect.TypeOf(ast.IndexAccess{}),
		reflect.TypeOf(ast.DotAccess{}),
		reflect.TypeOf(ast.DataType{}),
		reflect.TypeOf(ast.ColumnDef{}),
		reflect.TypeOf(ast.PrimaryKeyDef{}),
		reflect.TypeOf(ast.ClusteringOrder{}),
		reflect.TypeOf(ast.TableOption{}),
		reflect.TypeOf(ast.OptionHash{}),
		reflect.TypeOf(ast.OptionHashItem{}),
		reflect.TypeOf(ast.SelectElement{}),
		reflect.TypeOf(ast.AssignmentElement{}),
		reflect.TypeOf(ast.IfCondition{}),
		reflect.TypeOf(ast.UsingClause{}),
		reflect.TypeOf(ast.OrderByElement{}),
		reflect.TypeOf(ast.SelectStmt{}),
		reflect.TypeOf(ast.InsertStmt{}),
		reflect.TypeOf(ast.UpdateStmt{}),
		reflect.TypeOf(ast.DeleteStmt{}),
		reflect.TypeOf(ast.BatchStmt{}),
		reflect.TypeOf(ast.TruncateStmt{}),
		reflect.TypeOf(ast.UseStmt{}),
		reflect.TypeOf(ast.CreateKeyspaceStmt{}),
		reflect.TypeOf(ast.AlterKeyspaceStmt{}),
		reflect.TypeOf(ast.DropKeyspaceStmt{}),
		reflect.TypeOf(ast.CreateTableStmt{}),
		reflect.TypeOf(ast.AlterTableStmt{}),
		reflect.TypeOf(ast.DropTableStmt{}),
		reflect.TypeOf(ast.CreateIndexStmt{}),
		reflect.TypeOf(ast.DropIndexStmt{}),
		reflect.TypeOf(ast.CreateTypeStmt{}),
		reflect.TypeOf(ast.AlterTypeStmt{}),
		reflect.TypeOf(ast.AlterTypeRenameItem{}),
		reflect.TypeOf(ast.DropTypeStmt{}),
		reflect.TypeOf(ast.CreateMVStmt{}),
		reflect.TypeOf(ast.AlterMVStmt{}),
		reflect.TypeOf(ast.DropMVStmt{}),
		reflect.TypeOf(ast.CreateFunctionStmt{}),
		reflect.TypeOf(ast.FunctionParam{}),
		reflect.TypeOf(ast.DropFunctionStmt{}),
		reflect.TypeOf(ast.CreateAggregateStmt{}),
		reflect.TypeOf(ast.DropAggregateStmt{}),
		reflect.TypeOf(ast.CreateTriggerStmt{}),
		reflect.TypeOf(ast.DropTriggerStmt{}),
		reflect.TypeOf(ast.CreateRoleStmt{}),
		reflect.TypeOf(ast.RoleOption{}),
		reflect.TypeOf(ast.AlterRoleStmt{}),
		reflect.TypeOf(ast.DropRoleStmt{}),
		reflect.TypeOf(ast.CreateUserStmt{}),
		reflect.TypeOf(ast.AlterUserStmt{}),
		reflect.TypeOf(ast.DropUserStmt{}),
		reflect.TypeOf(ast.GrantStmt{}),
		reflect.TypeOf(ast.RevokeStmt{}),
		reflect.TypeOf(ast.Resource{}),
		reflect.TypeOf(ast.ListPermissionsStmt{}),
		reflect.TypeOf(ast.ListRolesStmt{}),
		reflect.TypeOf(ast.List{}),
		reflect.TypeOf(ast.RawStmt{}),
	}
	var result []string
	for _, c := range candidates {
		ptr := reflect.PointerTo(c)
		if ptr.Implements(nodeType) {
			result = append(result, c.Name())
		}
	}
	return result
}
