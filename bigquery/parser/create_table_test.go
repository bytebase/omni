package parser

import (
	"strings"
	"testing"

	"github.com/bytebase/omni/bigquery/ast"
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
// emulator in the *_oracle_test.go files.

// parseDDL parses a single DDL statement and fails the test on any parse error.
func parseDDL(t *testing.T, sql string) ast.Node {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) != 0 {
		t.Fatalf("parseForTest(%q): unexpected errors: %v", sql, errs)
	}
	if len(file.Stmts) != 1 {
		t.Fatalf("parseForTest(%q): got %d statements, want 1", sql, len(file.Stmts))
	}
	return file.Stmts[0]
}

// createTableOf parses sql and asserts the single statement is a *CreateTableStmt.
func createTableOf(t *testing.T, sql string) *ast.CreateTableStmt {
	t.Helper()
	n := parseDDL(t, sql)
	ct, ok := n.(*ast.CreateTableStmt)
	if !ok {
		t.Fatalf("parseForTest(%q): statement is %T, want *ast.CreateTableStmt", sql, n)
	}
	return ct
}

// assertReject parses sql and asserts at least one parse error is reported and
// no statement node is produced (the reject path returns a nil node).
func assertReject(t *testing.T, sql string) {
	t.Helper()
	file, errs := parseForTest(sql)
	if len(errs) == 0 {
		t.Errorf("parseForTest(%q): want a parse error, got none (stmts=%d)", sql, len(file.Stmts))
		return
	}
	if len(file.Stmts) != 0 {
		t.Errorf("parseForTest(%q): want 0 statements on reject, got %d", sql, len(file.Stmts))
	}
}

func TestCreateTable_GeneratedIdentity_Rejects(t *testing.T) {
	// The identity_column_info body is parsed structurally (NOT skipped as balanced
	// parens), so a malformed body rejects rather than silently parsing.
	cases := []string{
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY) PRIMARY KEY (id)",                               // missing parenthesized body
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (FOO BAR)) PRIMARY KEY (id)",                     // unknown clause tokens
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (START WITH)) PRIMARY KEY (id)",                  // START WITH missing numeric literal
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (INCREMENT 1)) PRIMARY KEY (id)",                 // INCREMENT missing BY
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (INCREMENT BY 1 START WITH 1)) PRIMARY KEY (id)", // out-of-order clauses
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (START WITH 1 START WITH 2)) PRIMARY KEY (id)",   // repeated clause
		"CREATE TABLE T (id INT64 GENERATED ALWAYS AS IDENTITY (NO)) PRIMARY KEY (id)",                          // NO without CYCLE
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}

func TestCreateTable_DefaultAndGeneratedRejected(t *testing.T) {
	// The grammar rejects a column with both DEFAULT and a generated clause.
	assertReject(t, "CREATE TABLE T (x INT64 DEFAULT (1) AS (2) STORED) PRIMARY KEY (x)")
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

func TestCreateTable_BareNullsRejected(t *testing.T) {
	// Regression (review): matchNullOrder must NOT consume a bare NULLS that is
	// not followed by FIRST/LAST; `PRIMARY KEY (x NULLS)` must reject (NULLS was
	// previously swallowed, wrongly accepting). oracle-confirmed reject.
	assertReject(t, "CREATE TABLE T (x INT64) PRIMARY KEY (x NULLS)")
}

func TestCreateTable_DuplicateInterleaveRejected(t *testing.T) {
	// Regression (review): the legacy grammar threads at most ONE
	// opt_spanner_interleave_in_parent_clause; a second INTERLEAVE must reject
	// (was previously over-accepted, last-wins). oracle-confirmed reject.
	assertReject(t, "CREATE TABLE T (x INT64) PRIMARY KEY (x), INTERLEAVE IN PARENT A ON DELETE CASCADE, INTERLEAVE IN PARENT B ON DELETE CASCADE")
}

func TestCreateTable_RowDeletionPolicyRequiresCommaAfterPK(t *testing.T) {
	// Regression (review/oracle): after a Spanner trailing PRIMARY KEY, ROW
	// DELETION POLICY is comma-led; the non-comma form must reject (was
	// over-accepted via the standalone opt_ttl_clause arm).
	assertReject(t, "CREATE TABLE T (x INT64, c TIMESTAMP) PRIMARY KEY (x) ROW DELETION POLICY (OLDER_THAN(c, INTERVAL 1 DAY))")
}

func TestCreateTable_Rejects(t *testing.T) {
	cases := []string{
		"CREATE TABLE",                                     // missing name
		"CREATE TABLE T (x INT64) PRIMARY KEY x",           // missing parens on PK
		"CREATE TABLE T (x INT64) PRIMARY KEY (x) garbage", // trailing junk
		"CREATE TABLE T (CONSTRAINT FOREIGN KEY (x))",      // CONSTRAINT missing name
		"CREATE TABLE T (x INT64,,)",                       // double comma in element list
	}
	for _, sql := range cases {
		assertReject(t, sql)
	}
}
