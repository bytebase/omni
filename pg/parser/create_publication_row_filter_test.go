package parser

import (
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/pg/ast"
)

func TestCreatePublicationTablesInSchemaWhere(t *testing.T) {
	// gram.y's TABLES IN_P SCHEMA ColId takes no OptWhereClause, so
	// PostgreSQL rejects a WHERE after it as a syntax error at WHERE. A
	// bare schema name continuing the list does take one, and
	// preprocess_pubobj_list rejects it at the name.
	for _, tc := range []struct {
		sql, at string
	}{
		{"CREATE PUBLICATION p FOR TABLES IN SCHEMA s1 WHERE (a = 123)", "WHERE"},
		{"CREATE PUBLICATION p FOR TABLE t, TABLES IN SCHEMA s1 WHERE (a > 10)", "WHERE"},
		{"CREATE PUBLICATION p FOR TABLES IN SCHEMA s1, s2 WHERE (a > 10)", "s2"},
	} {
		t.Run(tc.sql, func(t *testing.T) {
			parseFailsAt(t, tc.sql, strings.Index(tc.sql, tc.at))
		})
	}

	t.Run("TABLE t, TABLES IN SCHEMA s1 (baseline)", func(t *testing.T) {
		stmt := singleStmt(t, "CREATE PUBLICATION p FOR TABLE t, TABLES IN SCHEMA s1")
		cp, ok := stmt.(*nodes.CreatePublicationStmt)
		if !ok {
			t.Fatalf("expected CreatePublicationStmt, got %T", stmt)
		}
		if cp.Pubobjects == nil || len(cp.Pubobjects.Items) != 2 {
			t.Fatalf("expected 2 pubobjects, got %v", cp.Pubobjects)
		}
		spec1 := cp.Pubobjects.Items[1].(*nodes.PublicationObjSpec)
		if spec1.Pubobjtype != nodes.PUBLICATIONOBJ_TABLES_IN_SCHEMA || spec1.Name != "s1" {
			t.Fatalf("item[1] = %v %q, want TABLES_IN_SCHEMA s1", spec1.Pubobjtype, spec1.Name)
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
