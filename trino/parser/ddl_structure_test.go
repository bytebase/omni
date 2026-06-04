package parser

import (
	"testing"
)

// This file is the parser-ddl node's STRUCTURAL gate (correctness-protocol.md):
// the differential oracle (oracle_ddl_test.go) proves accept/reject parity with
// Trino 481, but cannot see the AST omni builds. These tests assert the parsed
// node type and its field values for one representative of each DDL form, so a
// regression that still parses but mis-assigns a clause (e.g. swapping RENAME
// source/target, dropping a property, losing IF EXISTS) is caught.

func TestDDLStructure_CreateTable(t *testing.T) {
	file, errs := Parse("CREATE OR REPLACE TABLE cat.sch.orders (orderkey bigint NOT NULL, name varchar COMMENT 'n', LIKE base INCLUDING PROPERTIES) COMMENT 'tbl' WITH (format = 'ORC')")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*CreateTableStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateTableStmt", file.Stmts[0])
	}
	if !stmt.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists = true, want false")
	}
	if got := stmt.Name.String(); got != "cat.sch.orders" {
		t.Errorf("Name = %q, want cat.sch.orders", got)
	}
	if len(stmt.Elements) != 3 {
		t.Fatalf("len(Elements) = %d, want 3", len(stmt.Elements))
	}
	// Element 0: orderkey bigint NOT NULL
	c0 := stmt.Elements[0].Column
	if c0 == nil || c0.Name.Value != "orderkey" || !c0.NotNull {
		t.Errorf("element0 = %+v, want column orderkey NOT NULL", c0)
	}
	// Element 1: name varchar COMMENT 'n'
	c1 := stmt.Elements[1].Column
	if c1 == nil || c1.Comment == nil || *c1.Comment != "n" {
		t.Errorf("element1 comment = %v, want 'n'", c1)
	}
	// Element 2: LIKE base INCLUDING PROPERTIES
	l2 := stmt.Elements[2].Like
	if l2 == nil || l2.Source.String() != "base" || !l2.HasOption || !l2.Including {
		t.Errorf("element2 = %+v, want LIKE base INCLUDING PROPERTIES", l2)
	}
	if stmt.Comment == nil || *stmt.Comment != "tbl" {
		t.Errorf("table comment = %v, want 'tbl'", stmt.Comment)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Name.Value != "format" {
		t.Errorf("properties = %+v, want [format=...]", stmt.Properties)
	}
}

func TestDDLStructure_CreateTableAs(t *testing.T) {
	file, errs := Parse("CREATE TABLE t (a, b) WITH (x = 1) AS SELECT 1, 2 WITH NO DATA")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*CreateTableAsStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateTableAsStmt", file.Stmts[0])
	}
	if len(stmt.ColumnAliases) != 2 || stmt.ColumnAliases[0].Value != "a" || stmt.ColumnAliases[1].Value != "b" {
		t.Errorf("ColumnAliases = %+v, want [a b]", stmt.ColumnAliases)
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
	if !stmt.HasDataClause || !stmt.NoData {
		t.Errorf("data clause: HasDataClause=%v NoData=%v, want true/true", stmt.HasDataClause, stmt.NoData)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("properties = %+v, want 1", stmt.Properties)
	}
}

func TestDDLStructure_ColumnDefault(t *testing.T) {
	// DEFAULT then NOT NULL (the only legal order, D-CT2).
	file, errs := Parse("CREATE TABLE t (status varchar DEFAULT 'created' NOT NULL)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*CreateTableStmt)
	col := stmt.Elements[0].Column
	if col.Default == nil {
		t.Error("Default is nil, want 'created' expression")
	}
	if !col.NotNull {
		t.Error("NotNull = false, want true")
	}
}

func TestDDLStructure_DropVariants(t *testing.T) {
	cases := []struct {
		sql      string
		object   DropObjectKind
		ifExists bool
		name     string
		behavior DropBehavior
	}{
		{"DROP TABLE IF EXISTS a.b.c", DropTable, true, "a.b.c", DropBehaviorDefault},
		{"DROP SCHEMA s CASCADE", DropSchema, false, "s", DropBehaviorCascade},
		{"DROP SCHEMA IF EXISTS s RESTRICT", DropSchema, true, "s", DropBehaviorRestrict},
		{"DROP VIEW v", DropView, false, "v", DropBehaviorDefault},
		{"DROP MATERIALIZED VIEW mv", DropMaterializedView, false, "mv", DropBehaviorDefault},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("unexpected errors: %v", errs)
			}
			stmt, ok := file.Stmts[0].(*DropStmt)
			if !ok {
				t.Fatalf("got %T, want *DropStmt", file.Stmts[0])
			}
			if stmt.Object != tc.object {
				t.Errorf("Object = %v, want %v", stmt.Object, tc.object)
			}
			if stmt.IfExists != tc.ifExists {
				t.Errorf("IfExists = %v, want %v", stmt.IfExists, tc.ifExists)
			}
			if stmt.Name.String() != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name.String(), tc.name)
			}
			if stmt.Behavior != tc.behavior {
				t.Errorf("Behavior = %v, want %v", stmt.Behavior, tc.behavior)
			}
		})
	}
}

func TestDDLStructure_AlterTableRename(t *testing.T) {
	file, _ := Parse("ALTER TABLE IF EXISTS users RENAME TO people")
	stmt := file.Stmts[0].(*AlterTableStmt)
	if stmt.Kind != AlterTableRename {
		t.Fatalf("Kind = %v, want AlterTableRename", stmt.Kind)
	}
	if !stmt.IfExists {
		t.Error("IfExists = false, want true")
	}
	if stmt.Name.String() != "users" || stmt.NewName.String() != "people" {
		t.Errorf("rename %q -> %q, want users -> people", stmt.Name, stmt.NewName)
	}
}

func TestDDLStructure_AlterTableAddColumn(t *testing.T) {
	file, _ := Parse("ALTER TABLE users ADD COLUMN IF NOT EXISTS zip varchar AFTER country")
	stmt := file.Stmts[0].(*AlterTableStmt)
	if stmt.Kind != AlterTableAddColumn {
		t.Fatalf("Kind = %v, want AlterTableAddColumn", stmt.Kind)
	}
	if !stmt.ColumnIfNotExists {
		t.Error("ColumnIfNotExists = false, want true")
	}
	if stmt.NewColumn == nil || stmt.NewColumn.Name.Value != "zip" {
		t.Errorf("NewColumn = %+v, want zip", stmt.NewColumn)
	}
	if stmt.Position != ColumnPositionAfter || stmt.PositionAfter.Value != "country" {
		t.Errorf("position = %v after %v, want AFTER country", stmt.Position, stmt.PositionAfter)
	}
}

func TestDDLStructure_AlterTableAlterColumn(t *testing.T) {
	cases := []struct {
		sql    string
		action AlterColumnAction
	}{
		{"ALTER TABLE t ALTER COLUMN c SET DATA TYPE bigint", AlterColumnSetDataType},
		{"ALTER TABLE t ALTER COLUMN c SET DEFAULT 1", AlterColumnSetDefault},
		{"ALTER TABLE t ALTER COLUMN c DROP DEFAULT", AlterColumnDropDefault},
		{"ALTER TABLE t ALTER COLUMN c DROP NOT NULL", AlterColumnDropNotNull},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt := file.Stmts[0].(*AlterTableStmt)
			if stmt.Kind != AlterTableAlterColumn {
				t.Fatalf("Kind = %v, want AlterTableAlterColumn", stmt.Kind)
			}
			if stmt.ColumnAction != tc.action {
				t.Errorf("ColumnAction = %v, want %v", stmt.ColumnAction, tc.action)
			}
			if stmt.AlterColumn.String() != "c" {
				t.Errorf("AlterColumn = %q, want c", stmt.AlterColumn)
			}
		})
	}
}

func TestDDLStructure_AlterTableExecute(t *testing.T) {
	file, errs := Parse("ALTER TABLE t EXECUTE optimize(file_size_threshold => '16MB') WHERE p = 1")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt := file.Stmts[0].(*AlterTableStmt)
	if stmt.Kind != AlterTableExecute {
		t.Fatalf("Kind = %v, want AlterTableExecute", stmt.Kind)
	}
	if stmt.Procedure.Value != "optimize" {
		t.Errorf("Procedure = %q, want optimize", stmt.Procedure.Value)
	}
	if !stmt.HasArgList || len(stmt.ExecuteArgs) != 1 {
		t.Errorf("args: HasArgList=%v n=%d, want true/1", stmt.HasArgList, len(stmt.ExecuteArgs))
	}
	if stmt.ExecuteArgs[0].Name == nil || stmt.ExecuteArgs[0].Name.Value != "file_size_threshold" {
		t.Errorf("arg name = %+v, want file_size_threshold", stmt.ExecuteArgs[0].Name)
	}
	if stmt.ExecuteWhere == nil {
		t.Error("ExecuteWhere is nil")
	}
}

func TestDDLStructure_AlterTableSetProperties(t *testing.T) {
	file, _ := Parse("ALTER TABLE t SET PROPERTIES a = 1, b = DEFAULT")
	stmt := file.Stmts[0].(*AlterTableStmt)
	if stmt.Kind != AlterTableSetProperties {
		t.Fatalf("Kind = %v, want AlterTableSetProperties", stmt.Kind)
	}
	if len(stmt.Properties) != 2 {
		t.Fatalf("len(Properties) = %d, want 2", len(stmt.Properties))
	}
	if stmt.Properties[0].IsDefault {
		t.Error("property[0] IsDefault = true, want false")
	}
	if !stmt.Properties[1].IsDefault {
		t.Error("property[1] IsDefault = false, want true (= DEFAULT)")
	}
}

func TestDDLStructure_CreateSchema(t *testing.T) {
	file, _ := Parse("CREATE SCHEMA IF NOT EXISTS hive.web AUTHORIZATION ROLE PUBLIC WITH (location = '/x')")
	stmt, ok := file.Stmts[0].(*CreateSchemaStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateSchemaStmt", file.Stmts[0])
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if stmt.Name.String() != "hive.web" {
		t.Errorf("Name = %q, want hive.web", stmt.Name)
	}
	if stmt.Authorization == nil || stmt.Authorization.Kind != PrincipalRole || stmt.Authorization.Name.Value != "PUBLIC" {
		t.Errorf("Authorization = %+v, want ROLE PUBLIC", stmt.Authorization)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %+v, want 1", stmt.Properties)
	}
}

func TestDDLStructure_AlterSchema(t *testing.T) {
	rename, _ := Parse("ALTER SCHEMA foo.bar RENAME TO baz")
	rs := rename.Stmts[0].(*AlterSchemaStmt)
	if rs.Kind != AlterSchemaRename || rs.Name.String() != "foo.bar" || rs.NewName.Value != "baz" {
		t.Errorf("rename = %+v, want foo.bar RENAME TO baz", rs)
	}
	auth, _ := Parse("ALTER SCHEMA web SET AUTHORIZATION USER alice")
	as := auth.Stmts[0].(*AlterSchemaStmt)
	if as.Kind != AlterSchemaSetAuthorization || as.Authorization.Kind != PrincipalUser {
		t.Errorf("set auth = %+v, want SET AUTHORIZATION USER alice", as)
	}
}

func TestDDLStructure_CreateView(t *testing.T) {
	file, _ := Parse("CREATE OR REPLACE VIEW v COMMENT 'c' SECURITY INVOKER AS SELECT 1")
	stmt, ok := file.Stmts[0].(*CreateViewStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateViewStmt", file.Stmts[0])
	}
	if !stmt.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if stmt.Comment == nil || *stmt.Comment != "c" {
		t.Errorf("Comment = %v, want 'c'", stmt.Comment)
	}
	if stmt.Security != ViewSecurityInvoker {
		t.Errorf("Security = %v, want ViewSecurityInvoker", stmt.Security)
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
}

func TestDDLStructure_AlterView(t *testing.T) {
	cases := []struct {
		sql  string
		kind AlterViewKind
	}{
		{"ALTER VIEW v RENAME TO w", AlterViewRename},
		{"ALTER VIEW v SET AUTHORIZATION alice", AlterViewSetAuthorization},
		{"ALTER VIEW v REFRESH", AlterViewRefresh},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt := file.Stmts[0].(*AlterViewStmt)
			if stmt.Kind != tc.kind {
				t.Errorf("Kind = %v, want %v", stmt.Kind, tc.kind)
			}
		})
	}
}

func TestDDLStructure_CreateMaterializedView(t *testing.T) {
	file, errs := Parse("CREATE MATERIALIZED VIEW mv GRACE PERIOD INTERVAL '1' HOUR WHEN STALE FAIL COMMENT 'c' WITH (p = 1) AS SELECT 1")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*CreateMaterializedViewStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateMaterializedViewStmt", file.Stmts[0])
	}
	if stmt.GracePeriod == nil {
		t.Error("GracePeriod is nil, want INTERVAL '1' HOUR")
	}
	if stmt.StaleMode != MVStaleFail {
		t.Errorf("StaleMode = %v, want MVStaleFail", stmt.StaleMode)
	}
	if stmt.Comment == nil || *stmt.Comment != "c" {
		t.Errorf("Comment = %v, want 'c'", stmt.Comment)
	}
	if len(stmt.Properties) != 1 {
		t.Errorf("Properties = %+v, want 1", stmt.Properties)
	}
	if stmt.Query == nil {
		t.Error("Query is nil")
	}
}

func TestDDLStructure_CreateCatalog(t *testing.T) {
	file, errs := Parse("CREATE CATALOG IF NOT EXISTS c USING postgresql COMMENT 'x' AUTHORIZATION alice WITH (\"connection-url\" = 'u')")
	if len(errs) != 0 {
		t.Fatalf("errors: %v", errs)
	}
	stmt, ok := file.Stmts[0].(*CreateCatalogStmt)
	if !ok {
		t.Fatalf("got %T, want *CreateCatalogStmt", file.Stmts[0])
	}
	if !stmt.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if stmt.Name.Value != "c" {
		t.Errorf("Name = %q, want c", stmt.Name.Value)
	}
	if stmt.Connector.Value != "postgresql" {
		t.Errorf("Connector = %q, want postgresql", stmt.Connector.Value)
	}
	if stmt.Comment == nil || *stmt.Comment != "x" {
		t.Errorf("Comment = %v, want 'x'", stmt.Comment)
	}
	if stmt.Authorization == nil || stmt.Authorization.Name.Value != "alice" {
		t.Errorf("Authorization = %+v, want alice", stmt.Authorization)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Name.Value != "connection-url" {
		t.Errorf("Properties = %+v, want connection-url", stmt.Properties)
	}
}

func TestDDLStructure_DropCatalog(t *testing.T) {
	file, _ := Parse("DROP CATALOG IF EXISTS c CASCADE")
	stmt, ok := file.Stmts[0].(*DropCatalogStmt)
	if !ok {
		t.Fatalf("got %T, want *DropCatalogStmt", file.Stmts[0])
	}
	if !stmt.IfExists || stmt.Name.Value != "c" || stmt.Behavior != DropBehaviorCascade {
		t.Errorf("got %+v, want IF EXISTS c CASCADE", stmt)
	}
}

func TestDDLStructure_Comment(t *testing.T) {
	cases := []struct {
		sql     string
		object  CommentObjectKind
		isNull  bool
		comment string
		name    string
	}{
		{"COMMENT ON TABLE users IS 'master'", CommentOnTable, false, "master", "users"},
		{"COMMENT ON TABLE users IS NULL", CommentOnTable, true, "", "users"},
		{"COMMENT ON VIEW v IS 'x'", CommentOnView, false, "x", "v"},
		{"COMMENT ON COLUMN a.b.c.d IS 'col'", CommentOnColumn, false, "col", "a.b.c.d"},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			file, errs := Parse(tc.sql)
			if len(errs) != 0 {
				t.Fatalf("errors: %v", errs)
			}
			stmt := file.Stmts[0].(*CommentStmt)
			if stmt.Object != tc.object {
				t.Errorf("Object = %v, want %v", stmt.Object, tc.object)
			}
			if stmt.IsNull != tc.isNull {
				t.Errorf("IsNull = %v, want %v", stmt.IsNull, tc.isNull)
			}
			if !tc.isNull && stmt.Comment != tc.comment {
				t.Errorf("Comment = %q, want %q", stmt.Comment, tc.comment)
			}
			if stmt.Name.String() != tc.name {
				t.Errorf("Name = %q, want %q", stmt.Name.String(), tc.name)
			}
		})
	}
}

func TestDDLStructure_Analyze(t *testing.T) {
	bare, _ := Parse("ANALYZE hive.default.stores")
	b := bare.Stmts[0].(*AnalyzeStmt)
	if b.Name.String() != "hive.default.stores" || len(b.Properties) != 0 {
		t.Errorf("bare ANALYZE = %+v, want hive.default.stores no props", b)
	}
	withProps, _ := Parse("ANALYZE t WITH (columns = ARRAY['a', 'b'])")
	w := withProps.Stmts[0].(*AnalyzeStmt)
	if len(w.Properties) != 1 || w.Properties[0].Name.Value != "columns" {
		t.Errorf("ANALYZE WITH = %+v, want columns property", w.Properties)
	}
}

func TestDDLStructure_RefreshMaterializedView(t *testing.T) {
	file, _ := Parse("REFRESH MATERIALIZED VIEW cat.sch.mv")
	stmt, ok := file.Stmts[0].(*RefreshMaterializedViewStmt)
	if !ok {
		t.Fatalf("got %T, want *RefreshMaterializedViewStmt", file.Stmts[0])
	}
	if stmt.Name.String() != "cat.sch.mv" {
		t.Errorf("Name = %q, want cat.sch.mv", stmt.Name)
	}
}
