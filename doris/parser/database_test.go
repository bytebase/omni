package parser

import (
	"fmt"
	"testing"

	"github.com/bytebase/omni/doris/ast"
)

// ---------------------------------------------------------------------------
// CREATE DATABASE
// ---------------------------------------------------------------------------

func TestCreateDatabase_Basic(t *testing.T) {
	file, errs := parseForTest("CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name == nil {
		t.Fatal("Name is nil")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.IfNotExists {
		t.Error("IfNotExists should be false")
	}
	if len(stmt.Properties) != 0 {
		t.Errorf("expected 0 properties, got %d", len(stmt.Properties))
	}
}

func TestCreateDatabase_IfNotExists(t *testing.T) {
	file, errs := parseForTest("CREATE DATABASE IF NOT EXISTS mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestCreateDatabase_WithProperties(t *testing.T) {
	file, errs := parseForTest(`CREATE DATABASE mydb PROPERTIES("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "key")
	}
	if stmt.Properties[0].Value != "value" {
		t.Errorf("Property value = %q, want %q", stmt.Properties[0].Value, "value")
	}
}

func TestCreateDatabase_WithMultipleProperties(t *testing.T) {
	file, errs := parseForTest(`CREATE DATABASE mydb PROPERTIES("replication_num"="3", "storage_medium"="SSD")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if len(stmt.Properties) != 2 {
		t.Fatalf("expected 2 properties, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "replication_num" {
		t.Errorf("Properties[0].Key = %q, want %q", stmt.Properties[0].Key, "replication_num")
	}
	if stmt.Properties[0].Value != "3" {
		t.Errorf("Properties[0].Value = %q, want %q", stmt.Properties[0].Value, "3")
	}
	if stmt.Properties[1].Key != "storage_medium" {
		t.Errorf("Properties[1].Key = %q, want %q", stmt.Properties[1].Key, "storage_medium")
	}
}

func TestCreateDatabase_IfNotExistsWithProperties(t *testing.T) {
	file, errs := parseForTest(`CREATE DATABASE IF NOT EXISTS mydb PROPERTIES("key"="val")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
}

// ---------------------------------------------------------------------------
// CREATE SCHEMA (alias for CREATE DATABASE)
// ---------------------------------------------------------------------------

func TestCreateSchema_Basic(t *testing.T) {
	file, errs := parseForTest("CREATE SCHEMA myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "myschema" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "myschema")
	}
}

func TestCreateSchema_IfNotExists(t *testing.T) {
	file, errs := parseForTest("CREATE SCHEMA IF NOT EXISTS myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.CreateDatabaseStmt)
	if !stmt.IfNotExists {
		t.Error("IfNotExists should be true")
	}
}

// ---------------------------------------------------------------------------
// DROP DATABASE
// ---------------------------------------------------------------------------

func TestDropDatabase_Basic(t *testing.T) {
	file, errs := parseForTest("DROP DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.IfExists {
		t.Error("IfExists should be false")
	}
	if stmt.Force {
		t.Error("Force should be false")
	}
}

func TestDropDatabase_IfExists(t *testing.T) {
	file, errs := parseForTest("DROP DATABASE IF EXISTS mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestDropDatabase_Force(t *testing.T) {
	file, errs := parseForTest("DROP DATABASE mydb FORCE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.Force {
		t.Error("Force should be true")
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

func TestDropDatabase_IfExistsForce(t *testing.T) {
	file, errs := parseForTest("DROP DATABASE IF EXISTS mydb FORCE")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
	if !stmt.Force {
		t.Error("Force should be true")
	}
}

// ---------------------------------------------------------------------------
// DROP SCHEMA (alias for DROP DATABASE)
// ---------------------------------------------------------------------------

func TestDropSchema_Basic(t *testing.T) {
	file, errs := parseForTest("DROP SCHEMA myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.DropDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "myschema" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "myschema")
	}
}

func TestDropSchema_IfExists(t *testing.T) {
	file, errs := parseForTest("DROP SCHEMA IF EXISTS myschema")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.DropDatabaseStmt)
	if !stmt.IfExists {
		t.Error("IfExists should be true")
	}
}

// ---------------------------------------------------------------------------
// ALTER DATABASE
// ---------------------------------------------------------------------------

func TestAlterDatabase_Rename(t *testing.T) {
	file, errs := parseForTest("ALTER DATABASE mydb RENAME newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if stmt.NewName == nil {
		t.Fatal("NewName is nil")
	}
	if stmt.NewName.String() != "newdb" {
		t.Errorf("NewName = %q, want %q", stmt.NewName.String(), "newdb")
	}
}

func TestAlterDatabase_RenameWithTo(t *testing.T) {
	file, errs := parseForTest("ALTER DATABASE mydb RENAME TO newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if stmt.NewName.String() != "newdb" {
		t.Errorf("NewName = %q, want %q", stmt.NewName.String(), "newdb")
	}
}

func TestAlterDatabase_SetProperties(t *testing.T) {
	file, errs := parseForTest(`ALTER DATABASE mydb SET PROPERTIES("key"="value")`)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	stmt := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
	if len(stmt.Properties) != 1 {
		t.Fatalf("expected 1 property, got %d", len(stmt.Properties))
	}
	if stmt.Properties[0].Key != "key" {
		t.Errorf("Property key = %q, want %q", stmt.Properties[0].Key, "key")
	}
	if stmt.Properties[0].Value != "value" {
		t.Errorf("Property value = %q, want %q", stmt.Properties[0].Value, "value")
	}
}

func TestAlterDatabase_SetQuota(t *testing.T) {
	// SET QUOTA is parsed in a best-effort fashion.
	file, errs := parseForTest("ALTER DATABASE mydb SET QUOTA 10737418240")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	stmt, ok := file.Stmts[0].(*ast.AlterDatabaseStmt)
	if !ok {
		t.Fatalf("expected *ast.AlterDatabaseStmt, got %T", file.Stmts[0])
	}
	if stmt.Name.String() != "mydb" {
		t.Errorf("Name = %q, want %q", stmt.Name.String(), "mydb")
	}
}

// ---------------------------------------------------------------------------
// NodeTag
// ---------------------------------------------------------------------------

func TestCreateDatabaseStmt_Tag(t *testing.T) {
	file, errs := parseForTest("CREATE DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_CreateDatabaseStmt {
		t.Errorf("Tag() = %v, want T_CreateDatabaseStmt", file.Stmts[0].Tag())
	}
}

func TestDropDatabaseStmt_Tag(t *testing.T) {
	file, errs := parseForTest("DROP DATABASE mydb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_DropDatabaseStmt {
		t.Errorf("Tag() = %v, want T_DropDatabaseStmt", file.Stmts[0].Tag())
	}
}

func TestAlterDatabaseStmt_Tag(t *testing.T) {
	file, errs := parseForTest("ALTER DATABASE mydb RENAME newdb")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if file.Stmts[0].Tag() != ast.T_AlterDatabaseStmt {
		t.Errorf("Tag() = %v, want T_AlterDatabaseStmt", file.Stmts[0].Tag())
	}
}

// ---------------------------------------------------------------------------
// Multiple statements
// ---------------------------------------------------------------------------

func TestCreateDropDatabase_MultiStatement(t *testing.T) {
	input := "CREATE DATABASE db1; DROP DATABASE db1"
	file, errs := parseForTest(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 2 {
		t.Fatalf("expected 2 stmts, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.CreateDatabaseStmt); !ok {
		t.Errorf("Stmts[0]: expected *ast.CreateDatabaseStmt, got %T", file.Stmts[0])
	}
	if _, ok := file.Stmts[1].(*ast.DropDatabaseStmt); !ok {
		t.Errorf("Stmts[1]: expected *ast.DropDatabaseStmt, got %T", file.Stmts[1])
	}
}

// ---------------------------------------------------------------------------
// Other CREATE/ALTER/DROP still unsupported
// ---------------------------------------------------------------------------

func TestCreateTable_NowSupported(t *testing.T) {
	file, errs := parseForTest("CREATE TABLE t (id INT)")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	_, ok := file.Stmts[0].(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("expected *ast.CreateTableStmt, got %T", file.Stmts[0])
	}
}

func TestAlterTable_NowSupported(t *testing.T) {
	// ALTER TABLE is now implemented (T2.2). Verify it parses without errors.
	file, errs := parseForTest("ALTER TABLE t ADD COLUMN c INT")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("expected 1 stmt, got %d", len(file.Stmts))
	}
	if _, ok := file.Stmts[0].(*ast.AlterTableStmt); !ok {
		t.Fatalf("expected *ast.AlterTableStmt, got %T", file.Stmts[0])
	}
}

func TestDropTable_NowSupported(t *testing.T) {
	// DROP TABLE is now implemented; verify it produces a DropTableStmt.
	file, errs := parseForTest("DROP TABLE t")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if _, ok := file.Stmts[0].(*ast.DropTableStmt); !ok {
		t.Errorf("expected *ast.DropTableStmt, got %T", file.Stmts[0])
	}
}

// TestPropertyListSeparators pins the shape the grammar gives every property
// list — propertyItemList : propertyItem (COMMA propertyItem)* — across the
// clauses that reach each of the parsers for one. Every clause below used to
// carry its own copy of the list loop with the comma optional, so each
// accepted an empty list, a missing comma and a trailing comma, none of which
// the engine takes (see the props_* rows in TestDorisSyntaxConformance for the
// container-verified verdicts).
func TestPropertyListSeparators(t *testing.T) {
	// %s is the property-list body.
	clauses := []struct{ name, tmpl string }{
		{"create_table", `CREATE TABLE zzp (a INT) DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES(%s)`},
		{"create_database", `CREATE DATABASE zzpdb PROPERTIES(%s)`},
		{"alter_catalog_set_properties", `ALTER CATALOG zzc SET PROPERTIES (%s)`},
		// The engine rejects this statement form outright (see
		// TestDorisSyntaxConformance), so these rows pin the list shape for
		// parsePropertyClause only -- they are not engine-backed.
		{"alter_catalog_property", `ALTER CATALOG zzc PROPERTY (%s)`},
		{"alter_table_set", `ALTER TABLE t SET (%s)`},
		{"admin_set_frontend_config", `ADMIN SET FRONTEND CONFIG (%s)`},
		{"create_mv", `CREATE MATERIALIZED VIEW zzpmv BUILD IMMEDIATE REFRESH COMPLETE DISTRIBUTED BY HASH(a) BUCKETS 1 PROPERTIES(%s) AS SELECT 1 AS a`},
		{"routine_load_job", `CREATE ROUTINE LOAD zzrl ON t PROPERTIES(%s) FROM KAFKA("kafka_broker_list"="localhost:9092","kafka_topic"="tp")`},
		{"routine_load_source", `CREATE ROUTINE LOAD zzrl ON t FROM KAFKA(%s)`},
		{"outfile", `SELECT 1 INTO OUTFILE '/tmp/x' PROPERTIES(%s)`},
		// The lists below were skipped unread or had their errors dropped.
		{"partition", `CREATE TABLE zzp (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (%s)) DISTRIBUTED BY HASH(a) BUCKETS 1`},
		{"partition_without_values", `CREATE TABLE zzp (a INT) PARTITION BY LIST(a) (PARTITION p1 ("1") (%s)) DISTRIBUTED BY HASH(a) BUCKETS 1`},
		{"load_broker", `LOAD LABEL zzl (DATA INFILE("hdfs://h/p") INTO TABLE t) WITH BROKER b (%s)`},
		{"export_broker", `EXPORT TABLE t TO "hdfs://h/p" WITH BROKER b (%s)`},
		{"admin_check_tablet", `ADMIN CHECK TABLET (1) PROPERTIES(%s)`},
		{"admin_copy_tablet", `ADMIN COPY TABLET 1 PROPERTIES(%s)`},
		{"admin_set_replica_status", `ADMIN SET REPLICA STATUS PROPERTIES(%s)`},
		{"admin_set_table_status", `ADMIN SET TABLE t STATUS PROPERTIES(%s)`},
		{"alter_system_modify_backend", `ALTER SYSTEM MODIFY BACKEND 'h:1' SET (%s)`},
		{"alter_system_add_backend", `ALTER SYSTEM ADD BACKEND 'h:1' PROPERTIES(%s)`},
	}
	bodies := []struct {
		name       string
		body       string
		wantAccept bool
	}{
		{"one item", `"k"="v"`, true},
		{"two items", `"k"="v","k2"="v2"`, true},
		{"three items", `"k"="v","k2"="v2","k3"="v3"`, true},
		{"empty", ``, false},
		{"missing comma", `"k"="v" "k2"="v2"`, false},
		{"missing commas", `"k"="v" "k2"="v2" "k3"="v3"`, false},
		{"trailing comma", `"k"="v",`, false},
		{"leading comma", `,"k"="v"`, false},
		{"doubled comma", `"k"="v",,"k2"="v2"`, false},
		{"comma only", `,`, false},
	}
	for _, cl := range clauses {
		for _, b := range bodies {
			t.Run(cl.name+"/"+b.name, func(t *testing.T) {
				sql := fmt.Sprintf(cl.tmpl, b.body)
				_, errs := parseForTest(sql)
				if got := len(errs) == 0; got != b.wantAccept {
					t.Errorf("accept = %v, want %v for %q (errs=%v)", got, b.wantAccept, sql, errs)
				}
			})
		}
	}
}

// TestPropertyListAtEndOfInput: a property list cut off by the end of the
// input is an error, reported there, in every clause. The ADMIN and ALTER
// SYSTEM parsers used to drop the list's error, the partition and broker lists
// were skipped to a ')' that never came, and nothing after the statement was
// left to notice, so each of these parsed (container-verified rejects).
func TestPropertyListAtEndOfInput(t *testing.T) {
	prefixes := []string{
		`ADMIN SET FRONTEND CONFIG `,
		`ALTER SYSTEM MODIFY BACKEND 'h:1' SET `,
		`ADMIN SET REPLICA STATUS PROPERTIES`,
		`ADMIN SET TABLE t STATUS PROPERTIES`,
		`ADMIN CHECK TABLET (1) PROPERTIES`,
		`ADMIN COPY TABLET 1 PROPERTIES`,
		`ALTER SYSTEM ADD BACKEND 'h:1' PROPERTIES`,
		`ALTER SYSTEM SET LOAD ERRORS HUB PROPERTIES`,
		`CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") `,
		`LOAD LABEL l (DATA INFILE("hdfs://h/p") INTO TABLE t) WITH BROKER b `,
		`EXPORT TABLE t TO "hdfs://h/p" WITH BROKER b `,
		`CREATE ROUTINE LOAD j ON t FROM KAFKA `,
	}
	cuts := []string{`(`, `("k"`, `("k"=`, `("k"="v"`, `("k"="v",`}
	for _, prefix := range prefixes {
		for _, cut := range cuts {
			sql := prefix + cut
			t.Run(sql, func(t *testing.T) {
				_, errs := parseForTest(sql)
				if len(errs) == 0 {
					t.Fatalf("%q parsed; want an error at the end of the input", sql)
				}
				if errs[0].Position != len(sql) {
					t.Errorf("%q: error at %d (%s), want %d", sql, errs[0].Position, errs[0].Message, len(sql))
				}
			})
		}
	}
}

// TestPropertyListEndsStatement: a statement that ends on a property list ends
// at its closing ')', and the list reaches the tree. These statements used to
// end at the last item's value, and the broker lists were not kept at all.
func TestPropertyListEndsStatement(t *testing.T) {
	systemProps := func(n ast.Node) []*ast.Property { return n.(*ast.SystemAlterStmt).Properties }
	cases := []struct {
		sql   string
		props func(ast.Node) []*ast.Property
	}{
		{`ADMIN SET FRONTEND CONFIG ("k"="v")`, adminProps},
		{`ADMIN SET REPLICA STATUS PROPERTIES("k"="v")`, adminProps},
		{`ADMIN SET TABLE t STATUS PROPERTIES("k"="v")`, adminProps},
		{`ADMIN CHECK TABLET (1) PROPERTIES("k"="v")`, adminProps},
		{`ADMIN COPY TABLET 1 PROPERTIES("k"="v")`, adminProps},
		{`ALTER SYSTEM ADD BACKEND 'h:1' PROPERTIES("k"="v")`, systemProps},
		{`ALTER SYSTEM SET LOAD ERRORS HUB PROPERTIES("k"="v")`, systemProps},
		{`ALTER SYSTEM MODIFY BACKEND 'h:1' SET ("k"="v")`, func(n ast.Node) []*ast.Property { return n.(*ast.SystemAlterStmt).SetClause }},
		{`LOAD LABEL l (DATA INFILE("hdfs://h/p") INTO TABLE t) WITH BROKER b ("k"="v")`, func(n ast.Node) []*ast.Property { return n.(*ast.LoadDataStmt).BrokerProperties }},
		{`EXPORT TABLE t TO "hdfs://h/p" WITH BROKER b ("k"="v")`, func(n ast.Node) []*ast.Property { return n.(*ast.ExportStmt).BrokerProperties }},
	}
	for _, tc := range cases {
		t.Run(tc.sql, func(t *testing.T) {
			node := parseOne(t, tc.sql)
			if loc := ast.NodeLoc(node); loc.Start != 0 || loc.End != len(tc.sql) {
				t.Errorf("Loc = %+v, want [0, %d)", loc, len(tc.sql))
			}
			props := tc.props(node)
			if len(props) != 1 || props[0].Key != "k" || props[0].Value != "v" {
				t.Fatalf("properties = %+v, want [k=v]", props)
			}
			if got := tc.sql[props[0].Loc.Start:props[0].Loc.End]; got != `"k"="v"` {
				t.Errorf("property text = %q, want %q", got, `"k"="v"`)
			}
		})
	}
}

func adminProps(n ast.Node) []*ast.Property { return n.(*ast.AdminStmt).Properties }

// TestAdminTabletPropertiesParsed: ADMIN CHECK TABLET and COPY TABLET collected
// everything after the object as raw Args, PROPERTIES list included, so the
// list was never held to its grammar and its branch never ran.
func TestAdminTabletPropertiesParsed(t *testing.T) {
	stmt := parseOne(t, `ADMIN CHECK TABLET (10000, 10001) PROPERTIES("type"="consistency")`).(*ast.AdminStmt)
	if stmt.Args != "( 10000 , 10001 )" {
		t.Errorf("Args = %q, want the tablet list only", stmt.Args)
	}
	if len(stmt.Properties) != 1 || stmt.Properties[0].Key != "type" || stmt.Properties[0].Value != "consistency" {
		t.Errorf("Properties = %+v, want [type=consistency]", stmt.Properties)
	}
	stmt = parseOne(t, `ADMIN COPY TABLET 10007 PROPERTIES("backend_id"="10001","version"="1")`).(*ast.AdminStmt)
	if stmt.Args != "10007" {
		t.Errorf("Args = %q, want the tablet id only", stmt.Args)
	}
	if len(stmt.Properties) != 2 {
		t.Errorf("Properties = %+v, want 2", stmt.Properties)
	}
}

// TestAlterSystemBrokerName: ALTER SYSTEM ADD / DROP [ALL] BROKER needs the
// broker's name. A missing one was an error that got dropped.
func TestAlterSystemBrokerName(t *testing.T) {
	for _, sql := range []string{`ALTER SYSTEM ADD BROKER`, `ALTER SYSTEM DROP BROKER`, `ALTER SYSTEM DROP ALL BROKER`} {
		_, errs := parseForTest(sql)
		if len(errs) == 0 || errs[0].Position != len(sql) {
			t.Errorf("%q: errs = %v, want one at the end of the input", sql, errs)
		}
	}
	stmt := parseOne(t, `ALTER SYSTEM DROP ALL BROKER b`).(*ast.SystemAlterStmt)
	if stmt.BrokerName != "b" || !stmt.DropAll {
		t.Errorf("BrokerName = %q DropAll = %v, want b true", stmt.BrokerName, stmt.DropAll)
	}
}

// TestShowAlterTableFromNeedsDatabase: the database after FROM is required, and
// a token that is not one is an error rather than raw text. Best-effort parsing
// keeps the raw recovery.
func TestShowAlterTableFromNeedsDatabase(t *testing.T) {
	for _, sql := range []string{`SHOW ALTER TABLE COLUMN FROM`, `SHOW ALTER TABLE COLUMN FROM 'x'`} {
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
	}
	if r := ParseBestEffort(`SHOW ALTER TABLE COLUMN FROM 'x'`); len(r.File.Stmts) != 1 {
		t.Errorf("ParseBestEffort: stmts = %d, want the raw recovery kept", len(r.File.Stmts))
	}
}

// TestPartitionProperties: a partition's inline property list is parsed into
// the tree, in the grammar's full item form. It used to be skipped to the next
// ')' unread.
func TestPartitionProperties(t *testing.T) {
	const sql = `CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("replication_num" = "1", storage_medium = HDD), PARTITION p2 VALUES LESS THAN ("20")) DISTRIBUTED BY HASH(a) BUCKETS 1`
	stmt := parseOne(t, sql).(*ast.CreateTableStmt)
	p1 := stmt.PartitionBy.Partitions[0]
	if len(p1.Properties) != 2 ||
		p1.Properties[0].Key != "replication_num" || p1.Properties[0].Value != "1" ||
		p1.Properties[1].Key != "storage_medium" || p1.Properties[1].Value != "HDD" {
		t.Fatalf("Properties = %+v", p1.Properties)
	}
	if got := sql[p1.Loc.Start:p1.Loc.End]; got != `PARTITION p1 VALUES LESS THAN ("10") ("replication_num" = "1", storage_medium = HDD)` {
		t.Errorf("partition text = %q", got)
	}
	if p2 := stmt.PartitionBy.Partitions[1]; p2.Properties != nil {
		t.Errorf("p2 Properties = %+v, want none", p2.Properties)
	}

	// propertyKey and propertyValue are identifier | constant.
	wrap := func(item string) string {
		return `CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") (` + item + `)) DISTRIBUTED BY HASH(a) BUCKETS 1`
	}
	accepted := []struct{ item, value string }{
		{`replication_num = "1"`, "1"},
		{`"k" = 1`, "1"},
		{`"k" = -1`, "-1"},
		{`"k" = 1.5`, "1.5"},
		{`"k" = 1e3`, "1e3"},
		{`"k" = NULL`, "NULL"},
		{`"k" = true`, "true"},
		{`"k" = HDD`, "HDD"},
		{"\"k\" = `v`", "v"},
		{`"k" = [1]`, "[1]"},
		{`"k" = DATE "2020-01-01"`, `DATE "2020-01-01"`},
		{`1 = "v"`, "v"},
		{`comment = "v"`, "v"},
		{`'k' = 'v'`, "v"},
	}
	for _, tc := range accepted {
		file, errs := parseForTest(wrap(tc.item))
		if len(errs) != 0 {
			t.Errorf("(%s): errs = %v, want accept", tc.item, errs)
			continue
		}
		props := file.Stmts[0].(*ast.CreateTableStmt).PartitionBy.Partitions[0].Properties
		if len(props) != 1 || props[0].Value != tc.value {
			t.Errorf("(%s): Properties = %+v, want value %q", tc.item, props, tc.value)
		}
	}
	for _, item := range []string{
		``, `"k"="v",`, `"k"="v" "k2"="v2"`, `"k" "v" SELECT DROP *`,
		`"k" = SELECT`, `"k" = a.b`, `"k" = 1+1`, `"k" = +1`, `"k" =`, `"k"`,
	} {
		if _, errs := parseForTest(wrap(item)); len(errs) == 0 {
			t.Errorf("(%s) parsed; the engine rejects it", item)
		}
	}
}

// TestPartitionWithoutValues: inPartitionDef leaves VALUES IN optional, so a
// list partition may carry its value list straight after the name, then its
// property list.
func TestPartitionWithoutValues(t *testing.T) {
	wrap := func(def string) string {
		return `CREATE TABLE t (a INT) PARTITION BY LIST(a) (` + def + `) DISTRIBUTED BY HASH(a) BUCKETS 1`
	}
	cases := []struct {
		def      string
		inValues [][]string
		props    int
	}{
		{`PARTITION p1`, nil, 0},
		{`PARTITION p1 ("1")`, [][]string{{"1"}}, 0},
		{`PARTITION p1 (("1"), ("2"))`, [][]string{{"1"}, {"2"}}, 0},
		{`PARTITION p1 ("replication_num"="1")`, nil, 1},
		{`PARTITION p1 ("1") ("replication_num"="1")`, [][]string{{"1"}}, 1},
	}
	for _, tc := range cases {
		file, errs := parseForTest(wrap(tc.def))
		if len(errs) != 0 {
			t.Errorf("%s: errs = %v", tc.def, errs)
			continue
		}
		item := file.Stmts[0].(*ast.CreateTableStmt).PartitionBy.Partitions[0]
		if fmt.Sprint(item.InValues) != fmt.Sprint(tc.inValues) || len(item.Properties) != tc.props {
			t.Errorf("%s: InValues = %v Properties = %+v, want %v and %d", tc.def, item.InValues, item.Properties, tc.inValues, tc.props)
		}
	}
	for _, def := range []string{`PARTITION p1 ()`, `PARTITION p1 (SELECT)`, `PARTITION p1 ("1") ()`} {
		if _, errs := parseForTest(wrap(def)); len(errs) == 0 {
			t.Errorf("%s parsed; the engine rejects it", def)
		}
	}
}

// TestBrokerPropertyItems: a broker load's list takes a string on both sides
// of each item, while EXPORT's takes the grammar's full item form; the engine
// differs between the two (container-verified).
func TestBrokerPropertyItems(t *testing.T) {
	const load = `LOAD LABEL l (DATA INFILE("hdfs://h/p") INTO TABLE t) WITH BROKER b (%s)`
	const export = `EXPORT TABLE t TO "hdfs://h/p" WITH BROKER b (%s)`
	for _, tc := range []struct {
		item                 string
		wantLoad, wantExport bool
	}{
		{`"username"="u"`, true, true},
		{`'k'='v'`, true, true},
		{`username="u"`, false, true},
		{`"k"=1`, false, true},
		{`"k"=v`, false, true},
		{`"k"=SELECT`, false, false},
		{`"k" "v" SELECT DROP *`, false, false},
	} {
		if _, errs := parseForTest(fmt.Sprintf(load, tc.item)); (len(errs) == 0) != tc.wantLoad {
			t.Errorf("LOAD (%s): accept = %v, want %v (errs=%v)", tc.item, len(errs) == 0, tc.wantLoad, errs)
		}
		if _, errs := parseForTest(fmt.Sprintf(export, tc.item)); (len(errs) == 0) != tc.wantExport {
			t.Errorf("EXPORT (%s): accept = %v, want %v (errs=%v)", tc.item, len(errs) == 0, tc.wantExport, errs)
		}
	}
}

// TestRoutineLoadSourcePropertyItems: a routine load's data-source list takes
// a string on both sides of each item. Its value used to be any single token,
// so ("a" = ,) and ("a" = )) parsed.
func TestRoutineLoadSourcePropertyItems(t *testing.T) {
	const tmpl = `CREATE ROUTINE LOAD j ON t FROM KAFKA (%s)`
	stmt := parseOne(t, fmt.Sprintf(tmpl, `"kafka_broker_list"="h:9092", 'kafka_topic'='tp'`)).(*ast.CreateRoutineLoadStmt)
	if len(stmt.DataSourceProperties) != 2 || stmt.DataSourceProperties[1].Value != "tp" {
		t.Errorf("DataSourceProperties = %+v", stmt.DataSourceProperties)
	}
	for _, item := range []string{
		`"a" = ,`, `"a" = )`, `"a" = (`, `"a" =`,
		`kafka_broker_list="h:9092"`, `comment="v"`, `"k"=1`, `"k"=1.5`, `"k"=tp`,
		`"k"=NULL`, `"k"=true`, `"k"=SELECT`, `"k"=*`, "\"k\"=`v`",
	} {
		sql := fmt.Sprintf(tmpl, item)
		if _, errs := parseForTest(sql); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", sql)
		}
		alter := fmt.Sprintf(`ALTER ROUTINE LOAD FOR j FROM KAFKA (%s)`, item)
		if _, errs := parseForTest(alter); len(errs) == 0 {
			t.Errorf("%q parsed; the engine rejects it", alter)
		}
	}
}

// TestWalkReachesPropertyLists: the walker visits the property lists that now
// reach the tree -- a partition's, and a broker load's or export's WITH BROKER
// list -- as it does every other property list.
func TestWalkReachesPropertyLists(t *testing.T) {
	for _, sql := range []string{
		`CREATE TABLE t (a INT) PARTITION BY RANGE(a) (PARTITION p1 VALUES LESS THAN ("10") ("k"="v")) DISTRIBUTED BY HASH(a) BUCKETS 1`,
		`LOAD LABEL l (DATA INFILE("hdfs://h/p") INTO TABLE t) WITH BROKER b ("k"="v")`,
		`EXPORT TABLE t TO "hdfs://h/p" WITH BROKER b ("k"="v")`,
	} {
		var keys []string
		ast.Inspect(parseOne(t, sql), func(n ast.Node) bool {
			if p, ok := n.(*ast.Property); ok {
				keys = append(keys, p.Key)
			}
			return true
		})
		if len(keys) != 1 || keys[0] != "k" {
			t.Errorf("%q: walked properties %v, want [k]", sql, keys)
		}
	}
}
