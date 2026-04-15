package parser

import (
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestCreatePublicationTablesInSchemaWhere(t *testing.T) {
	// TABLES IN SCHEMA with trailing WHERE is accepted at parse time in
	// PG (errored out at validation time, because row filters don't
	// apply to schema-level specs). omni should match parse-time
	// permissiveness.
	t.Run("TABLES IN SCHEMA s1 WHERE (...)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE PUBLICATION p FOR TABLES IN SCHEMA s1 WHERE (a = 123)")
		cp, ok := stmt.(*nodes.CreatePublicationStmt)
		if !ok {
			t.Fatalf("expected CreatePublicationStmt, got %T", stmt)
		}
		if cp.Pubobjects == nil || len(cp.Pubobjects.Items) != 1 {
			t.Fatalf("expected 1 pubobject, got %v", cp.Pubobjects)
		}
		spec := cp.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
		if spec.Pubobjtype != nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA {
			t.Fatalf("expected TABLES_IN_SCHEMA, got %v", spec.Pubobjtype)
		}
		if spec.Name != "s1" {
			t.Fatalf("expected schema s1, got %q", spec.Name)
		}
	})

	t.Run("mixed TABLE t, TABLES IN SCHEMA s1 WHERE (...)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE PUBLICATION p FOR TABLE t, TABLES IN SCHEMA s1 WHERE (a > 10)")
		cp, ok := stmt.(*nodes.CreatePublicationStmt)
		if !ok {
			t.Fatalf("expected CreatePublicationStmt, got %T", stmt)
		}
		if cp.Pubobjects == nil || len(cp.Pubobjects.Items) != 2 {
			t.Fatalf("expected 2 pubobjects, got %v", cp.Pubobjects)
		}
		// Ensure the permissive WHERE parse did not reclassify item[0]
		// as TABLES_IN_SCHEMA or collapse both into one entry.
		spec0 := cp.Pubobjects.Items[0].(*nodes.PublicationObjSpec)
		if spec0.Pubobjtype != nodes.PUBLICATIONOBJ_TABLE {
			t.Fatalf("item[0] expected TABLE, got %v", spec0.Pubobjtype)
		}
		spec1 := cp.Pubobjects.Items[1].(*nodes.PublicationObjSpec)
		if spec1.Pubobjtype != nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA {
			t.Fatalf("item[1] expected TABLES_IN_SCHEMA, got %v", spec1.Pubobjtype)
		}
		if spec1.Name != "s1" {
			t.Fatalf("item[1] expected schema s1, got %q", spec1.Name)
		}
	})

	// Regression-sanity: existing forms must still work.
	t.Run("FOR TABLE t (baseline)", func(t *testing.T) {
		parseOK(t, "CREATE PUBLICATION p FOR TABLE t")
	})

	t.Run("FOR TABLE t WHERE (a > 10) (baseline)", func(t *testing.T) {
		parseOK(t, "CREATE PUBLICATION p FOR TABLE t WHERE (a > 10)")
	})
}
