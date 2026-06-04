package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/googlesql/ast"
)

// Unit tests for the parser-ddl node (core DDL: CREATE/ALTER/DROP over TABLE,
// VIEW, INDEX, SCHEMA, DATABASE for the BigQuery + Spanner union). These assert
// the AST STRUCTURE (the structural gate — accept/reject alone does not catch
// wrong nesting/attribute placement) and cover BigQuery-only forms the Spanner
// emulator feature/syntax-rejects (OR REPLACE, PARTITION/CLUSTER/OPTIONS, CTAS,
// inline PK, NOT ENFORCED, LIKE/COPY, GENERATED ALWAYS, SQL SECURITY DEFINER,
// VIRTUAL) which are NON-AUTHORITATIVE against that oracle and triangulated from
// the legacy GoogleSQLParser.g4 + the BigQuery truth1 corpus. The accept/reject
// behavior of the Spanner-authoritative forms is also proven against the live
// emulator in the *_oracle_test.go files (build tag googlesql_oracle).

// parseDDL parses a single DDL statement and fails the test on any parse error.
func parseDDL(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) != 0 {
		t.Fatalf("Parse(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("Parse(%q): got %d statements, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

// createTableOf parses sql and asserts the single statement is a *CreateTableStmt.
func createTableOf(t *testing.T, sql string) *ast.CreateTableStmt {
	t.Helper()
	n := parseDDL(t, sql)
	ct, ok := n.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("Parse(%q): statement is %T, want *ast.CreateTableStmt", sql, n)
	}
	return ct
}

// assertReject parses sql and asserts at least one parse error is reported and
// no statement node is produced (the reject path returns a nil node).
func assertReject(t *testing.T, sql string) {
	t.Helper()
	file, errs := Parse(sql)
	if len(errs) == 0 {
		t.Errorf("Parse(%q): want a parse error, got none (stmts=%d)", sql, len(file.Stmts))
		return
	}
	if len(file.Stmts) != 0 {
		t.Errorf("Parse(%q): want 0 statements on reject, got %d", sql, len(file.Stmts))
	}
}

func TestCreateTable_Basic(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64) PRIMARY KEY (x)")
	if ct.Name.String() != "T" {
		t.Errorf("Name = %q, want T", ct.Name.String())
	}
	if len(ct.Columns) != 1 || ct.Columns[0].Name != "x" {
		t.Fatalf("Columns = %+v, want one column x", ct.Columns)
	}
	if ct.Columns[0].Type == nil || ct.Columns[0].Type.Text != "INT64" {
		t.Errorf("column type = %+v, want INT64", ct.Columns[0].Type)
	}
	if !ct.HasPrimaryKey || len(ct.PrimaryKey) != 1 || ct.PrimaryKey[0].Name != "x" {
		t.Errorf("PrimaryKey = %+v (has=%v), want [x]", ct.PrimaryKey, ct.HasPrimaryKey)
	}
}

func TestCreateTable_ColumnsAndAttributes(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (a INT64 NOT NULL, b STRING(MAX), c STRING(10) HIDDEN) PRIMARY KEY (a)")
	if len(ct.Columns) != 3 {
		t.Fatalf("got %d columns, want 3", len(ct.Columns))
	}
	if !ct.Columns[0].NotNull {
		t.Errorf("column a: NotNull = false, want true")
	}
	if ct.Columns[1].Type.Text != "STRING(MAX)" {
		t.Errorf("column b type = %q, want STRING(MAX)", ct.Columns[1].Type.Text)
	}
	if !ct.Columns[2].Hidden {
		t.Errorf("column c: Hidden = false, want true")
	}
}

func TestCreateTable_IfNotExistsAndSchemaName(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE IF NOT EXISTS myschema.Emp (id INT64) PRIMARY KEY (id)")
	if !ct.IfNotExists {
		t.Error("IfNotExists = false, want true")
	}
	if ct.Name.String() != "myschema.Emp" {
		t.Errorf("Name = %q, want myschema.Emp", ct.Name.String())
	}
}

func TestCreateTable_PrimaryKeyDirection(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, y INT64) PRIMARY KEY (x DESC, y ASC)")
	if len(ct.PrimaryKey) != 2 {
		t.Fatalf("got %d PK parts, want 2", len(ct.PrimaryKey))
	}
	if ct.PrimaryKey[0].Direction != "DESC" {
		t.Errorf("PK[0].Direction = %q, want DESC", ct.PrimaryKey[0].Direction)
	}
	if ct.PrimaryKey[1].Direction != "ASC" {
		t.Errorf("PK[1].Direction = %q, want ASC", ct.PrimaryKey[1].Direction)
	}
}

func TestCreateTable_Interleave(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE Albums (s INT64, a INT64) PRIMARY KEY (s, a), INTERLEAVE IN PARENT Singers ON DELETE CASCADE")
	if ct.Interleave == nil {
		t.Fatal("Interleave = nil, want a clause")
	}
	if ct.Interleave.Parent.String() != "Singers" {
		t.Errorf("Interleave parent = %q, want Singers", ct.Interleave.Parent.String())
	}
	if ct.Interleave.OnDelete != ast.FKActionCascade {
		t.Errorf("Interleave OnDelete = %v, want CASCADE", ct.Interleave.OnDelete)
	}
}

func TestCreateTable_InterleaveNoAction(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE D (p INT64, d INT64) PRIMARY KEY (p, d), INTERLEAVE IN PARENT P ON DELETE NO ACTION")
	if ct.Interleave == nil || ct.Interleave.OnDelete != ast.FKActionNoAction {
		t.Errorf("Interleave OnDelete = %v, want NO ACTION", ct.Interleave)
	}
}

func TestCreateTable_GeneratedStored(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, f STRING(20) AS (CONCAT(x)) STORED) PRIMARY KEY (x)")
	col := ct.Columns[1]
	if col.Generated == nil {
		t.Fatal("column f: Generated = nil, want the AS expression")
	}
	if col.GenMode != ast.GenModeAs {
		t.Errorf("column f: GenMode = %v, want GenModeAs", col.GenMode)
	}
	if col.Stored != "STORED" {
		t.Errorf("column f: Stored = %q, want STORED", col.Stored)
	}
}

func TestCreateTable_GeneratedVirtual(t *testing.T) {
	// VIRTUAL is documented in the Spanner truth1 (DDL-005) but the Spanner
	// emulator rejects it (divergence — see create_table.go); the union parser
	// accepts it.
	ct := createTableOf(t, "CREATE TABLE T (x INT64, n INT64 AS (x) VIRTUAL) PRIMARY KEY (x)")
	if ct.Columns[1].Stored != "VIRTUAL" {
		t.Errorf("column n: Stored = %q, want VIRTUAL", ct.Columns[1].Stored)
	}
}

func TestCreateTable_GeneratedAlways(t *testing.T) {
	// BigQuery GENERATED ALWAYS AS … STORED (rejected by the Spanner emulator;
	// triangulated from the .g4 + BigQuery truth1).
	ct := createTableOf(t, "CREATE TABLE T (x INT64, n INT64 GENERATED ALWAYS AS (x + 1) STORED) PRIMARY KEY (x)")
	col := ct.Columns[1]
	if col.GenMode != ast.GenModeGeneratedAlways {
		t.Errorf("GenMode = %v, want GenModeGeneratedAlways", col.GenMode)
	}
	if col.Generated == nil {
		t.Error("Generated = nil, want the expression")
	}
}

func TestCreateTable_Default(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, s STRING(20) DEFAULT ('pending')) PRIMARY KEY (x)")
	if ct.Columns[1].Default == nil {
		t.Error("column s: Default = nil, want the DEFAULT expression")
	}
}

func TestCreateTable_DefaultAndGeneratedRejected(t *testing.T) {
	// The grammar rejects a column with both DEFAULT and a generated clause.
	assertReject(t, "CREATE TABLE T (x INT64 DEFAULT (1) AS (2) STORED) PRIMARY KEY (x)")
}

func TestCreateTable_CheckConstraint(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, CONSTRAINT chk CHECK (x > 0)) PRIMARY KEY (x)")
	if len(ct.Constraints) != 1 {
		t.Fatalf("got %d constraints, want 1", len(ct.Constraints))
	}
	c := ct.Constraints[0]
	if c.Kind != ast.ConstraintCheck {
		t.Errorf("constraint kind = %v, want CHECK", c.Kind)
	}
	if c.Name != "chk" {
		t.Errorf("constraint name = %q, want chk", c.Name)
	}
	if c.Check == nil {
		t.Error("Check expression = nil")
	}
}

func TestCreateTable_ForeignKeyConstraint(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, CONSTRAINT fk FOREIGN KEY (x) REFERENCES P (id) ON DELETE NO ACTION) PRIMARY KEY (x)")
	c := ct.Constraints[0]
	if c.Kind != ast.ConstraintForeignKey {
		t.Fatalf("constraint kind = %v, want FOREIGN KEY", c.Kind)
	}
	if len(c.Columns) != 1 || c.Columns[0] != "x" {
		t.Errorf("FK columns = %v, want [x]", c.Columns)
	}
	if c.ForeignKey == nil || c.ForeignKey.Table.String() != "P" {
		t.Errorf("FK ref table = %+v, want P", c.ForeignKey)
	}
	if c.ForeignKey.OnDelete != ast.FKActionNoAction {
		t.Errorf("FK OnDelete = %v, want NO ACTION", c.ForeignKey.OnDelete)
	}
}

func TestCreateTable_InlineForeignKey(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64 REFERENCES P (id)) PRIMARY KEY (x)")
	if ct.Columns[0].ForeignKey == nil {
		t.Fatal("column x: ForeignKey = nil, want inline REFERENCES")
	}
	if ct.Columns[0].ForeignKey.Table.String() != "P" {
		t.Errorf("inline FK table = %q, want P", ct.Columns[0].ForeignKey.Table.String())
	}
}

func TestCreateTable_RowDeletionPolicy(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, c TIMESTAMP) PRIMARY KEY (x), ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 30 DAY))")
	if ct.RowDeletion == nil {
		t.Error("RowDeletion = nil, want the policy expression")
	}
}

func TestCreateTable_ColumnOptions(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE T (x INT64, ts TIMESTAMP OPTIONS (allow_commit_timestamp = true)) PRIMARY KEY (x)")
	col := ct.Columns[1]
	if len(col.Options) != 1 {
		t.Fatalf("got %d column options, want 1", len(col.Options))
	}
	if col.Options[0].Name != "allow_commit_timestamp" || col.Options[0].Op != "=" {
		t.Errorf("option = %+v, want allow_commit_timestamp =", col.Options[0])
	}
}

func TestCreateTable_EmptyElementsAndEmptyPK(t *testing.T) {
	// oracle-confirmed accept: CREATE TABLE T () PRIMARY KEY ().
	ct := createTableOf(t, "CREATE TABLE T () PRIMARY KEY ()")
	if len(ct.Columns) != 0 {
		t.Errorf("got %d columns, want 0", len(ct.Columns))
	}
	if !ct.HasPrimaryKey || len(ct.PrimaryKey) != 0 {
		t.Errorf("PrimaryKey = %+v (has=%v), want present-and-empty", ct.PrimaryKey, ct.HasPrimaryKey)
	}
}

func TestCreateTable_TrailingCommaInElementList(t *testing.T) {
	// oracle-confirmed accept: a trailing comma before ')' (the grammar's
	// COMMA_SYMBOL? after the last table_element).
	ct := createTableOf(t, "CREATE TABLE T (x INT64,) PRIMARY KEY (x)")
	if len(ct.Columns) != 1 {
		t.Errorf("got %d columns, want 1", len(ct.Columns))
	}
}

func TestCreateTable_NoPrimaryKey(t *testing.T) {
	// A table without any PRIMARY KEY clause parses (the union grammar permits it
	// — BigQuery tables have no PK; Spanner's PK-required rule is semantic, not
	// syntactic).
	ct := createTableOf(t, "CREATE TABLE T (x INT64)")
	if ct.HasPrimaryKey {
		t.Error("HasPrimaryKey = true, want false (no PK clause)")
	}
}

// --- BigQuery-only CREATE TABLE forms (Spanner-emulator non-authoritative) ---

func TestCreateTable_OrReplaceTemp(t *testing.T) {
	ct := createTableOf(t, "CREATE OR REPLACE TEMP TABLE t (n INT64)")
	if !ct.OrReplace {
		t.Error("OrReplace = false, want true")
	}
	if ct.Scope != "TEMP" {
		t.Errorf("Scope = %q, want TEMP", ct.Scope)
	}
}

func TestCreateTable_PartitionClusterOptions(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE t (ts TIMESTAMP, v INT64) PARTITION BY TIMESTAMP_TRUNC(ts, DAY) CLUSTER BY v OPTIONS (description = 'd')")
	if len(ct.PartitionBy) != 1 {
		t.Errorf("PartitionBy = %d exprs, want 1", len(ct.PartitionBy))
	}
	if len(ct.ClusterBy) != 1 {
		t.Errorf("ClusterBy = %d exprs, want 1", len(ct.ClusterBy))
	}
	if len(ct.Options) != 1 || ct.Options[0].Name != "description" {
		t.Errorf("Options = %+v, want [description]", ct.Options)
	}
}

func TestCreateTable_CTAS(t *testing.T) {
	ct := createTableOf(t, "CREATE OR REPLACE TABLE t AS SELECT a, b FROM src")
	if ct.AsQuery == nil {
		t.Fatal("AsQuery = nil, want the SELECT body")
	}
	if _, ok := ct.AsQuery.(*ast.QueryStmt); !ok {
		t.Errorf("AsQuery is %T, want *ast.QueryStmt", ct.AsQuery)
	}
}

func TestCreateTable_LikeCopy(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE t LIKE src.tbl")
	if ct.Like == nil || ct.Like.String() != "src.tbl" {
		t.Errorf("Like = %+v, want src.tbl", ct.Like)
	}
	ct2 := createTableOf(t, "CREATE TABLE t COPY src.tbl")
	if ct2.Copy == nil || ct2.Copy.String() != "src.tbl" {
		t.Errorf("Copy = %+v, want src.tbl", ct2.Copy)
	}
}

func TestCreateTable_InlinePrimaryKeyNotEnforced(t *testing.T) {
	// BigQuery inline PRIMARY KEY NOT ENFORCED (Spanner rejects inline PK).
	ct := createTableOf(t, "CREATE TABLE t (id INT64 PRIMARY KEY NOT ENFORCED)")
	if !ct.Columns[0].PrimaryKey {
		t.Error("column id: PrimaryKey = false, want true")
	}
	if ct.Columns[0].Enforced != "NOT ENFORCED" {
		t.Errorf("column id: Enforced = %q, want NOT ENFORCED", ct.Columns[0].Enforced)
	}
}

func TestCreateTable_DefaultCollate(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE t (s STRING) DEFAULT COLLATE 'und:ci'")
	if !strings.Contains(ct.DefaultCollate, "und:ci") {
		t.Errorf("DefaultCollate = %q, want it to contain und:ci", ct.DefaultCollate)
	}
}

func TestCreateTable_StructArrayColumns(t *testing.T) {
	ct := createTableOf(t, "CREATE TABLE t (info STRUCT<a INT64, b STRING>, tags ARRAY<STRING>)")
	if ct.Columns[0].Type.Text != "STRUCT<a INT64, b STRING>" {
		t.Errorf("STRUCT column type = %q", ct.Columns[0].Type.Text)
	}
	if ct.Columns[1].Type.Text != "ARRAY<STRING>" {
		t.Errorf("ARRAY column type = %q", ct.Columns[1].Type.Text)
	}
}

// --- CREATE TABLE reject cases ---

func TestCreateTable_Rejects(t *testing.T) {
	cases := []string{
		"CREATE TABLE",                                    // missing name
		"CREATE TABLE T (x INT64) PRIMARY KEY x",          // missing parens on PK
		"CREATE TABLE T (x)",                              // column missing type
		"CREATE TABLE T (x INT64) PRIMARY KEY (x) garbage", // trailing junk
		"CREATE TABLE T (CONSTRAINT FOREIGN KEY (x))",      // CONSTRAINT missing name
		"CREATE TABLE T (x INT64,,)",                       // double comma in element list
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}
