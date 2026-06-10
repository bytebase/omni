package catalog

import (
	"fmt"
	"strings"
	"testing"

	nodes "github.com/bytebase/omni/redshift/ast"
	pgparser "github.com/bytebase/omni/redshift/parser"
)

type testRelationResolver map[string]*RelationSpec

func (r testRelationResolver) ResolveRelation(schemaName, relationName string, searchPath []string) (*RelationSpec, error) {
	if schemaName != "" {
		return r[schemaName+"."+relationName], nil
	}
	for _, schema := range searchPath {
		if spec := r[schema+"."+relationName]; spec != nil {
			return spec, nil
		}
	}
	return nil, nil
}

func TestRelationResolverLazyTable(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns: []RelationColumnSpec{
				{Name: "id", Type: "int4"},
				{Name: "name", Type: "text"},
			},
		},
	})

	q := analyzeSelectSQL(t, c, "SELECT id, name FROM accounts")

	if len(q.TargetList) != 2 {
		t.Fatalf("target list length = %d, want 2", len(q.TargetList))
	}
	_, rel, err := c.findRelation("public", "accounts")
	if err != nil {
		t.Fatalf("find materialized relation: %v", err)
	}
	if rel.RelKind != 'r' {
		t.Fatalf("relkind = %q, want table", rel.RelKind)
	}
	if got := columnNames(rel); strings.Join(got, ",") != "id,name" {
		t.Fatalf("columns = %v, want [id name]", got)
	}
}

func TestRelationResolverViewToTable(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.active_accounts": {
			SchemaName: "public",
			Name:       "active_accounts",
			Kind:       'v',
			Definition: "SELECT id FROM accounts WHERE active",
		},
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns: []RelationColumnSpec{
				{Name: "id", Type: "int4"},
				{Name: "active", Type: "bool"},
			},
		},
	})

	q := analyzeSelectSQL(t, c, "SELECT id FROM active_accounts")

	if len(q.TargetList) != 1 || q.TargetList[0].ResName != "id" {
		t.Fatalf("target list = %+v, want id", q.TargetList)
	}
	_, rel, err := c.findRelation("public", "active_accounts")
	if err != nil {
		t.Fatalf("find materialized view: %v", err)
	}
	if rel.RelKind != 'v' || rel.AnalyzedQuery == nil {
		t.Fatalf("view relkind/analyzed = %q/%v, want analyzed view", rel.RelKind, rel.AnalyzedQuery)
	}
}

func TestRelationResolverFullCreateViewDefinition(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.active_accounts": {
			SchemaName: "public",
			Name:       "active_accounts",
			Kind:       'v',
			Definition: "CREATE VIEW public.active_accounts(account_id) AS SELECT id FROM accounts WHERE active",
		},
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns: []RelationColumnSpec{
				{Name: "id", Type: "int4"},
				{Name: "active", Type: "bool"},
			},
		},
	})

	q := analyzeSelectSQL(t, c, "SELECT account_id FROM active_accounts")

	if len(q.TargetList) != 1 || q.TargetList[0].ResName != "account_id" {
		t.Fatalf("target list = %+v, want account_id", q.TargetList)
	}
	_, rel, err := c.findRelation("public", "active_accounts")
	if err != nil {
		t.Fatalf("find materialized view: %v", err)
	}
	if got := columnNames(rel); strings.Join(got, ",") != "account_id" {
		t.Fatalf("columns = %v, want [account_id]", got)
	}
}

func TestRelationResolverViewToLaterView(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.outer_v": {
			SchemaName: "public",
			Name:       "outer_v",
			Kind:       'v',
			Definition: "SELECT id FROM inner_v",
		},
		"public.inner_v": {
			SchemaName: "public",
			Name:       "inner_v",
			Kind:       'v',
			Definition: "SELECT id FROM base_t",
		},
		"public.base_t": {
			SchemaName: "public",
			Name:       "base_t",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "id", Type: "int8"}},
		},
	})

	analyzeSelectSQL(t, c, "SELECT id FROM outer_v")

	for _, name := range []string{"outer_v", "inner_v", "base_t"} {
		if _, _, err := c.findRelation("public", name); err != nil {
			t.Fatalf("relation %s was not materialized: %v", name, err)
		}
	}
}

func TestRelationResolverUsesSearchPath(t *testing.T) {
	c := New()
	c.SetSearchPath([]string{"tenant", "public"})
	c.SetRelationResolver(testRelationResolver{
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "public_id", Type: "int4"}},
		},
		"tenant.accounts": {
			SchemaName: "tenant",
			Name:       "accounts",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "tenant_id", Type: "int4"}},
		},
	})

	q := analyzeSelectSQL(t, c, "SELECT tenant_id FROM accounts")

	if got := q.TargetList[0].ResName; got != "tenant_id" {
		t.Fatalf("resolved column = %q, want tenant_id", got)
	}
	if _, _, err := c.findRelation("tenant", "accounts"); err != nil {
		t.Fatalf("tenant relation was not materialized: %v", err)
	}
}

func TestRelationResolverCyclicView(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.a_view": {
			SchemaName: "public",
			Name:       "a_view",
			Kind:       'v',
			Definition: "SELECT id FROM b_view",
		},
		"public.b_view": {
			SchemaName: "public",
			Name:       "b_view",
			Kind:       'v',
			Definition: "SELECT id FROM a_view",
		},
	})

	_, err := analyzeSelectSQLError(c, "SELECT id FROM a_view")
	if err == nil || !strings.Contains(err.Error(), "cyclic relation resolution") {
		t.Fatalf("error = %v, want cyclic relation resolution", err)
	}
}

func TestRelationResolverMissingDependency(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.v": {
			SchemaName: "public",
			Name:       "v",
			Kind:       'v',
			Definition: "SELECT id FROM missing_t",
		},
	})

	_, err := analyzeSelectSQLError(c, "SELECT id FROM v")
	if err == nil || !strings.Contains(err.Error(), "missing_t") {
		t.Fatalf("error = %v, want missing dependency error", err)
	}
	if _, _, findErr := c.findRelation("public", "v"); findErr == nil {
		t.Fatalf("view was materialized despite missing dependency")
	}
}

func TestRelationResolverUnsupportedDefinition(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.bad_v": {
			SchemaName: "public",
			Name:       "bad_v",
			Kind:       'v',
			Definition: "DELETE FROM accounts",
		},
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "id", Type: "int4"}},
		},
	})

	_, err := analyzeSelectSQLError(c, "SELECT id FROM bad_v")
	if err == nil || !strings.Contains(err.Error(), "requires a SELECT definition") {
		t.Fatalf("error = %v, want unsupported definition error", err)
	}
}

func TestRelationResolverLazyMaterializedView(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.recent_accounts": {
			SchemaName: "public",
			Name:       "recent_accounts",
			Kind:       'm',
			Definition: "SELECT id FROM accounts",
		},
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "id", Type: "int4"}},
		},
	})

	analyzeSelectSQL(t, c, "SELECT id FROM recent_accounts")

	_, rel, err := c.findRelation("public", "recent_accounts")
	if err != nil {
		t.Fatalf("find materialized view: %v", err)
	}
	if rel.RelKind != 'm' || rel.AnalyzedQuery == nil {
		t.Fatalf("relkind/analyzed = %q/%v, want analyzed matview", rel.RelKind, rel.AnalyzedQuery)
	}
}

func TestRelationResolverFullCreateMaterializedViewDefinition(t *testing.T) {
	c := New()
	c.SetRelationResolver(testRelationResolver{
		"public.recent_accounts": {
			SchemaName: "public",
			Name:       "recent_accounts",
			Kind:       'm',
			Definition: "CREATE MATERIALIZED VIEW public.recent_accounts(account_id) AS SELECT id FROM accounts",
		},
		"public.accounts": {
			SchemaName: "public",
			Name:       "accounts",
			Kind:       'r',
			Columns:    []RelationColumnSpec{{Name: "id", Type: "int4"}},
		},
	})

	analyzeSelectSQL(t, c, "SELECT account_id FROM recent_accounts")

	_, rel, err := c.findRelation("public", "recent_accounts")
	if err != nil {
		t.Fatalf("find materialized view: %v", err)
	}
	if rel.RelKind != 'm' || rel.AnalyzedQuery == nil {
		t.Fatalf("relkind/analyzed = %q/%v, want analyzed matview", rel.RelKind, rel.AnalyzedQuery)
	}
	if got := columnNames(rel); strings.Join(got, ",") != "account_id" {
		t.Fatalf("columns = %v, want [account_id]", got)
	}
}

func analyzeSelectSQL(t *testing.T, c *Catalog, sql string) *Query {
	t.Helper()
	q, err := analyzeSelectSQLError(c, sql)
	if err != nil {
		t.Fatalf("AnalyzeSelectStmt(%q): %v", sql, err)
	}
	return q
}

func analyzeSelectSQLError(c *Catalog, sql string) (*Query, error) {
	stmt, err := parseSingleSelect(sql)
	if err != nil {
		return nil, err
	}
	return c.AnalyzeSelectStmt(stmt)
}

func parseSingleSelect(sql string) (*nodes.SelectStmt, error) {
	list, err := pgparser.Parse(sql)
	if err != nil {
		return nil, err
	}
	if list == nil || len(list.Items) != 1 {
		return nil, fmt.Errorf("expected one statement")
	}
	raw, ok := list.Items[0].(*nodes.RawStmt)
	if !ok {
		return nil, fmt.Errorf("expected RawStmt, got %T", list.Items[0])
	}
	stmt, ok := raw.Stmt.(*nodes.SelectStmt)
	if !ok {
		return nil, fmt.Errorf("expected SelectStmt, got %T", raw.Stmt)
	}
	return stmt, nil
}

func columnNames(rel *Relation) []string {
	names := make([]string, len(rel.Columns))
	for i, col := range rel.Columns {
		names[i] = col.Name
	}
	return names
}
