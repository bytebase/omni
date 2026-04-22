package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

// TestCreateSchemaSchemaElementList exercises KB-2c: CREATE SCHEMA must
// absorb a whitespace-separated list of schema_elements (CreateStmt,
// CreateSeqStmt, CreateTrigStmt, GrantStmt, ViewStmt, CreateIndexStmt)
// into CreateSchemaStmt.SchemaElts rather than silently leaving them as
// follow-on top-level statements.
func TestCreateSchemaSchemaElementList(t *testing.T) {
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

	t.Run("plain CREATE SCHEMA foo (no elements)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA foo")
		elts := schemaElts(t, stmt)
		if len(elts) != 0 {
			t.Fatalf("expected 0 schema elements, got %d", len(elts))
		}
		css := stmt.(*nodes.CreateSchemaStmt)
		if css.Schemaname != "foo" {
			t.Errorf("expected schemaname=foo, got %q", css.Schemaname)
		}
		if css.Authrole != nil {
			t.Errorf("expected nil authrole, got %+v", css.Authrole)
		}
	})

	t.Run("CREATE SCHEMA foo AUTHORIZATION bob (no elements)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA foo AUTHORIZATION bob")
		elts := schemaElts(t, stmt)
		if len(elts) != 0 {
			t.Fatalf("expected 0 schema elements, got %d", len(elts))
		}
		css := stmt.(*nodes.CreateSchemaStmt)
		if css.Schemaname != "foo" {
			t.Errorf("expected schemaname=foo, got %q", css.Schemaname)
		}
		if css.Authrole == nil {
			t.Fatalf("expected authrole, got nil")
		}
	})

	t.Run("CREATE SCHEMA foo with inline CREATE TABLE", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA foo CREATE TABLE bar (id int)")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.CreateStmt); !ok {
			t.Fatalf("expected *CreateStmt, got %T", elts[0])
		}
	})

	t.Run("CREATE SCHEMA AUTHORIZATION bob with inline CREATE TABLE and CREATE VIEW", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA AUTHORIZATION bob CREATE TABLE x (id int) CREATE VIEW v AS SELECT 1")
		elts := schemaElts(t, stmt)
		if len(elts) != 2 {
			t.Fatalf("expected 2 schema elements, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.CreateStmt); !ok {
			t.Errorf("elt[0]: expected *CreateStmt, got %T", elts[0])
		}
		if _, ok := elts[1].(*nodes.ViewStmt); !ok {
			t.Errorf("elt[1]: expected *ViewStmt, got %T", elts[1])
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE TRIGGER (pgregress KB-2c)", func(t *testing.T) {
		sql := `CREATE SCHEMA AUTHORIZATION bob
			CREATE TRIGGER schema_trig BEFORE INSERT ON schema_not_existing.tab
			EXECUTE FUNCTION schema_trig.no_func()`
		stmt := singleStmt(t, sql)
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.CreateTrigStmt); !ok {
			t.Fatalf("expected *CreateTrigStmt, got %T", elts[0])
		}
	})

	t.Run("CREATE SCHEMA with inline CREATE INDEX and CREATE SEQUENCE", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA s CREATE SEQUENCE seq CREATE INDEX ON t (id)")
		elts := schemaElts(t, stmt)
		if len(elts) != 2 {
			t.Fatalf("expected 2 schema elements, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.CreateSeqStmt); !ok {
			t.Errorf("elt[0]: expected *CreateSeqStmt, got %T", elts[0])
		}
		if _, ok := elts[1].(*nodes.IndexStmt); !ok {
			t.Errorf("elt[1]: expected *IndexStmt, got %T", elts[1])
		}
	})

	t.Run("CREATE SCHEMA with inline GRANT", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE SCHEMA foo GRANT USAGE ON SCHEMA foo TO public")
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.GrantStmt); !ok {
			t.Fatalf("expected *GrantStmt, got %T", elts[0])
		}
	})

	t.Run("CREATE SCHEMA name AUTHORIZATION CURRENT_ROLE with inline CREATE TRIGGER (line 45)", func(t *testing.T) {
		sql := `CREATE SCHEMA regress_schema_1 AUTHORIZATION CURRENT_ROLE
			CREATE TRIGGER schema_trig BEFORE INSERT ON schema_not_existing.tab
			EXECUTE FUNCTION schema_trig.no_func()`
		stmt := singleStmt(t, sql)
		css := stmt.(*nodes.CreateSchemaStmt)
		if css.Schemaname != "regress_schema_1" {
			t.Errorf("expected schemaname=regress_schema_1, got %q", css.Schemaname)
		}
		if css.Authrole == nil {
			t.Fatal("expected authrole, got nil")
		}
		elts := schemaElts(t, stmt)
		if len(elts) != 1 {
			t.Fatalf("expected 1 schema element, got %d", len(elts))
		}
		if _, ok := elts[0].(*nodes.CreateTrigStmt); !ok {
			t.Fatalf("expected *CreateTrigStmt, got %T", elts[0])
		}
	})
}
