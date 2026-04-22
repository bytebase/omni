package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// KB-2d bug 1: CREATE SCHEMA's schema_element CREATE arm did not recognize
// OptTemp modifiers (TEMP/TEMPORARY/LOCAL TEMP/UNLOGGED) or OR REPLACE combined
// with OptTemp before VIEW/TABLE/SEQUENCE. The result was that
//
//	CREATE SCHEMA s CREATE TEMP VIEW v AS SELECT 1
//
// silently split into two RawStmts (the schema_stmt list terminated at TEMP,
// and TEMP VIEW ... started a new top-level statement via
// parseCreateTempDispatch). These tests pin the parser to absorb the full
// OptTemp variant into CreateSchemaStmt.SchemaElts.
func TestCreateSchemaSchemaElementOptTemp(t *testing.T) {
	schemaElts := func(t *testing.T, stmt nodes.Node) []nodes.Node {
		t.Helper()
		css, ok := stmt.(*nodes.CreateSchemaStmt)
		if !ok {
			t.Fatalf("expected *CreateSchemaStmt, got %T", stmt)
		}
		if css.SchemaElts == nil {
			return nil
		}
		return css.SchemaElts.Items
	}

	t.Run("CREATE SCHEMA with inline CREATE TEMP VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE TEMP VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		vs, ok := elts[0].(*nodes.ViewStmt)
		if !ok {
			t.Fatalf("expected *ViewStmt, got %T", elts[0])
		}
		if vs.View == nil || vs.View.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %+v", vs.View)
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE TEMPORARY VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE TEMPORARY VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		vs := elts[0].(*nodes.ViewStmt)
		if vs.View.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %q", string(vs.View.Relpersistence))
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE LOCAL TEMP VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE LOCAL TEMP VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		vs := elts[0].(*nodes.ViewStmt)
		if vs.View.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %q", string(vs.View.Relpersistence))
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE OR REPLACE TEMP VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE OR REPLACE TEMP VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		vs := elts[0].(*nodes.ViewStmt)
		if !vs.Replace {
			t.Errorf("expected Replace=true")
		}
		if vs.View.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %q", string(vs.View.Relpersistence))
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE OR REPLACE TEMPORARY VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE OR REPLACE TEMPORARY VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		vs := elts[0].(*nodes.ViewStmt)
		if !vs.Replace {
			t.Errorf("expected Replace=true")
		}
		if vs.View.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %q", string(vs.View.Relpersistence))
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE TEMP SEQUENCE", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE TEMP SEQUENCE seq")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		cs, ok := elts[0].(*nodes.CreateSeqStmt)
		if !ok {
			t.Fatalf("expected *CreateSeqStmt, got %T", elts[0])
		}
		if cs.Sequence == nil || cs.Sequence.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %+v", cs.Sequence)
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE UNLOGGED TABLE", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE UNLOGGED TABLE t (id int)")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		cs, ok := elts[0].(*nodes.CreateStmt)
		if !ok {
			t.Fatalf("expected *CreateStmt, got %T", elts[0])
		}
		if cs.Relation == nil || cs.Relation.Relpersistence != 'u' {
			t.Errorf("expected Relpersistence='u', got %+v", cs.Relation)
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE TEMP TABLE", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE TEMP TABLE t (id int)")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		cs := elts[0].(*nodes.CreateStmt)
		if cs.Relation.Relpersistence != 't' {
			t.Errorf("expected Relpersistence='t', got %q", string(cs.Relation.Relpersistence))
		}
	})

	t.Run("CREATE SCHEMA with mixed TEMP VIEW + regular TABLE", func(t *testing.T) {
		sql := "CREATE SCHEMA s CREATE TEMP VIEW v AS SELECT 1 CREATE TABLE t (id int)"
		stmt := singleStmt(t, sql)
		elts := schemaElts(t, stmt)
		if len(elts) != 2 {
			t.Fatalf("expected 2 schema elements, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.ViewStmt); !ok {
			t.Errorf("elt[0]: expected *ViewStmt, got %T", elts[0])
		}
		if _, ok := elts[1].(*nodes.CreateStmt); !ok {
			t.Errorf("elt[1]: expected *CreateStmt, got %T", elts[1])
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE RECURSIVE VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE RECURSIVE VIEW v (c) AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.ViewStmt); !ok {
			t.Fatalf("expected *ViewStmt, got %T", elts[0])
		}
	})

	// Direct repro of the KB-2d bug-1 input (test_view_schema) — must now
	// parse as a single CreateSchemaStmt with one ViewStmt element, rather
	// than splitting into 2 RawStmts.
	t.Run("KB-2d bug 1 repro", func(t *testing.T) {
		list, err := Parse("CREATE SCHEMA test_view_schema CREATE TEMP VIEW testview AS SELECT 1")
		if err != nil {
			t.Fatalf("parse: %v", err)
		}
		if list == nil || len(list.Items) != 1 {
			t.Fatalf("expected 1 RawStmt, got %d", len(list.Items))
		}
	})
}
