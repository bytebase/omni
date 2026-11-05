package parser

import (
	"testing"

	ast "github.com/bytebase/omni/tidb/ast"
)

func oneNode(t *testing.T, sql string) ast.Node {
	t.Helper()
	list, err := Parse(sql)
	if err != nil {
		t.Fatalf("Parse(%q): %v", sql, err)
	}
	if len(list.Items) != 1 {
		t.Fatalf("Parse(%q): %d statements, want 1", sql, len(list.Items))
	}
	return list.Items[0]
}

// TestVectorColumnTypeAST: the type comes out as VECTOR with its dimension, whichever
// spelling wrote it, and its Loc ends at the type's last token.
func TestVectorColumnTypeAST(t *testing.T) {
	sql := "CREATE TABLE t (a VECTOR, b VECTOR(3), c VECTOR<FLOAT>, d VECTOR<FLOAT>(4), e INT)"
	stmt := oneNode(t, sql).(*ast.CreateTableStmt)
	want := []struct {
		text   string
		length int
	}{{"VECTOR", 0}, {"VECTOR(3)", 3}, {"VECTOR<FLOAT>", 0}, {"VECTOR<FLOAT>(4)", 4}}
	for i, w := range want {
		dt := stmt.Columns[i].TypeName
		if dt.Name != "VECTOR" || dt.Length != w.length {
			t.Errorf("column %d: Name=%q Length=%d, want VECTOR and %d", i, dt.Name, dt.Length, w.length)
		}
		if got := sql[dt.Loc.Start:dt.Loc.End]; got != w.text {
			t.Errorf("column %d: Loc spans %q, want %q", i, got, w.text)
		}
	}
	if len(stmt.Columns) != 5 || stmt.Columns[4].TypeName.Name != "INT" {
		t.Errorf("the column after the vector types was lost: %+v", stmt.Columns)
	}
}

// TestCreateVectorIndexAST: the statement says it is a vector index, and every clause
// around the key part lands where a consumer reads it.
func TestCreateVectorIndexAST(t *testing.T) {
	sql := "CREATE VECTOR INDEX IF NOT EXISTS idx USING HNSW ON d.t ((VEC_COSINE_DISTANCE(a))) COMMENT 'c' ALGORITHM = INPLACE LOCK = NONE"
	stmt := oneNode(t, sql).(*ast.CreateIndexStmt)
	if !stmt.Vector || stmt.Unique || stmt.Fulltext || stmt.Spatial {
		t.Errorf("Vector=%v Unique=%v Fulltext=%v Spatial=%v, want only Vector", stmt.Vector, stmt.Unique, stmt.Fulltext, stmt.Spatial)
	}
	if !stmt.IfNotExists || stmt.IndexName != "idx" || stmt.IndexType != "HNSW" {
		t.Errorf("IfNotExists=%v IndexName=%q IndexType=%q", stmt.IfNotExists, stmt.IndexName, stmt.IndexType)
	}
	if stmt.Table.Schema != "d" || stmt.Table.Name != "t" {
		t.Errorf("Table=%s.%s, want d.t", stmt.Table.Schema, stmt.Table.Name)
	}
	if len(stmt.Columns) != 1 || !stmt.Columns[0].Functional {
		t.Fatalf("Columns=%+v, want one functional key part", stmt.Columns)
	}
	if _, ok := stmt.Columns[0].Expr.(*ast.FuncCallExpr); !ok {
		t.Errorf("the key part is %T, want a function call", stmt.Columns[0].Expr)
	}
	if len(stmt.Options) != 1 || stmt.Options[0].Name != "COMMENT" {
		t.Errorf("Options=%+v, want COMMENT", stmt.Options)
	}
	if stmt.Algorithm != "INPLACE" || stmt.Lock != "NONE" {
		t.Errorf("Algorithm=%q Lock=%q", stmt.Algorithm, stmt.Lock)
	}
	if stmt.Loc.Start != 0 || stmt.Loc.End != len(sql) {
		t.Errorf("Loc=%+v, want the whole statement of %d bytes", stmt.Loc, len(sql))
	}

	plain := oneNode(t, "CREATE INDEX i ON t (a)").(*ast.CreateIndexStmt)
	if plain.Vector || plain.Columns[0].Functional {
		t.Errorf("a plain index reads as Vector=%v Functional=%v", plain.Vector, plain.Columns[0].Functional)
	}
}

// TestVectorIndexTableElementAST: VECTOR INDEX in CREATE TABLE beside a column of the
// type, and a column that is only named vector.
func TestVectorIndexTableElementAST(t *testing.T) {
	stmt := oneNode(t, "CREATE TABLE vector (vector VECTOR(3), VECTOR INDEX vi USING HNSW ((VEC_L2_DISTANCE(vector))))").(*ast.CreateTableStmt)
	if len(stmt.Columns) != 1 || stmt.Columns[0].Name != "vector" || stmt.Columns[0].TypeName.Name != "VECTOR" {
		t.Fatalf("Columns=%+v", stmt.Columns)
	}
	if len(stmt.Constraints) != 1 {
		t.Fatalf("Constraints=%+v", stmt.Constraints)
	}
	c := stmt.Constraints[0]
	if c.Type != ast.ConstrVector || c.Name != "vi" || c.IndexType != "HNSW" {
		t.Errorf("Type=%v Name=%q IndexType=%q", c.Type, c.Name, c.IndexType)
	}
	if len(c.IndexColumns) != 1 || !c.IndexColumns[0].Functional {
		t.Errorf("IndexColumns=%+v, want one functional key part", c.IndexColumns)
	}

	// CONSTRAINT here takes vector as its symbol, and the key is an ordinary one.
	named := oneNode(t, "CREATE TABLE t (a VECTOR(3), CONSTRAINT VECTOR INDEX (a))").(*ast.CreateTableStmt)
	if got := named.Constraints[0]; got.Type != ast.ConstrIndex || got.Name != "VECTOR" {
		t.Errorf("Type=%v Name=%q, want an ordinary index with the symbol VECTOR", got.Type, got.Name)
	}
}

// TestVectorSyntaxRefused holds the shape check to CreateIndexStmt's and
// ConstraintVectorIndex's actions. Each case answered 1149 on v8.5.5 when refused is
// true, and any other error or none when it is false (read on 2026-11-05).
func TestVectorSyntaxRefused(t *testing.T) {
	for _, c := range []struct {
		sql     string
		refused bool
	}{
		{"CREATE VECTOR INDEX i ON t ((VEC_COSINE_DISTANCE(a)))", false},
		{"CREATE VECTOR INDEX i ON t ((a))", false},
		{"CREATE VECTOR INDEX i ON t ((a) DESC)", false},
		{"CREATE VECTOR INDEX i USING HNSW ON t ((a))", false},
		{"CREATE VECTOR INDEX i ON t ((a)) TYPE HNSW", false},
		{"CREATE VECTOR INDEX i USING BTREE ON t ((a)) USING HNSW", false},
		{"CREATE VECTOR INDEX i ON t (a)", true},
		{"CREATE VECTOR INDEX i ON t (a(3))", true},
		{"CREATE VECTOR INDEX i ON t (a, b)", true},
		{"CREATE VECTOR INDEX i ON t ((a), (b))", true},
		{"CREATE VECTOR INDEX i ON t ((a)) USING HYPO", true},
		{"CREATE VECTOR INDEX i USING HASH ON t ((a))", true},
		{"CREATE VECTOR INDEX i USING HNSW ON t ((a)) USING BTREE", true},
		{"CREATE INDEX i ON t (a)", false},
		{"CREATE INDEX i ON t (a) USING BTREE", false},
		{"CREATE INDEX i USING HNSW ON t (a)", true},
		{"CREATE UNIQUE INDEX i ON t (a) TYPE HNSW", true},
		{"CREATE FULLTEXT INDEX i ON t (a) USING HNSW", true},
	} {
		if got := oneNode(t, c.sql).(*ast.CreateIndexStmt).SyntaxRefused(); got != c.refused {
			t.Errorf("%q: SyntaxRefused=%v, want %v", c.sql, got, c.refused)
		}
	}
	for _, c := range []struct {
		sql     string
		refused bool
	}{
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX ((VEC_COSINE_DISTANCE(b))))", false},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX i USING HNSW ((b)))", false},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX i USING HNSW ((b)) USING BTREE)", false},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX (b))", true},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX ((b)) USING HASH)", true},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX i USING BTREE ((b)) USING HNSW)", true},
		{"CREATE TABLE t (b VECTOR(3), VECTOR INDEX ((b), a))", true},
		{"CREATE TABLE t (a INT, INDEX (a) USING HNSW)", false},
	} {
		if got := oneNode(t, c.sql).(*ast.CreateTableStmt).Constraints[0].SyntaxRefused(); got != c.refused {
			t.Errorf("%q: SyntaxRefused=%v, want %v", c.sql, got, c.refused)
		}
	}
}
