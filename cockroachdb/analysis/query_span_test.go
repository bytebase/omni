package analysis

import (
	"context"
	"errors"
	"slices"
	"strings"
	"testing"

	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroachdb-parser/pkg/sql/pgwire/pgerror"
	"github.com/google/go-cmp/cmp"

	"github.com/bytebase/omni/analysis"
	metadatapb "github.com/bytebase/omni/metadata"
)

func spanTestMetadata() map[string]*metadatapb.DatabaseSchemaMetadata {
	column := func(name string) *metadatapb.ColumnMetadata { return &metadatapb.ColumnMetadata{Name: name} }
	return map[string]*metadatapb.DatabaseSchemaMetadata{
		"db": {Name: "db", Schemas: []*metadatapb.SchemaMetadata{
			{
				Name: "public",
				Tables: []*metadatapb.TableMetadata{
					{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("id"), column("ssn"), column("name")}},
					{Name: "u", Columns: []*metadatapb.ColumnMetadata{column("id"), column("t_id"), column("note"), {Name: "rowid", IsInvisible: true}}},
					{Name: "empty"},
				},
				Views: []*metadatapb.ViewMetadata{
					// A view whose sync stored no columns.
					{Name: "v", Definition: "SELECT ssn FROM public.t"},
					{Name: "v2", Definition: "SELECT ssn AS s, note FROM db.public.v JOIN db.public.u ON true", Columns: []*metadatapb.ColumnMetadata{column("s"), column("note")}},
				},
			},
			{Name: "hr", Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("id"), column("salary")}}}},
		}},
		"other": {Name: "other", Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "x", Columns: []*metadatapb.ColumnMetadata{column("a")}}}},
		}},
	}
}

func resolveSpan(statement string) (*analysis.ResolvedSpan, error) {
	databases := spanTestMetadata()
	return ResolveQuerySpan(context.Background(), statement, ResolveOptions{
		DefaultDatabase: "db",
		SearchPath:      []string{"public"},
		Catalog: Catalog{GetDatabase: func(name string) (Database, error) {
			if meta, ok := databases[name]; ok {
				return protoDatabase{meta}, nil
			}
			return nil, nil
		}},
	})
}

// lineage renders each result as "name: source source ...".
func lineage(span *analysis.ResolvedSpan) []string {
	var out []string
	for _, r := range span.Results {
		out = append(out, strings.TrimSpace(r.Name+": "+columnNames(r.SourceColumns)))
	}
	return out
}

func columnNames(set analysis.ColumnSet) string {
	var names []string
	for c := range set {
		names = append(names, c.String())
	}
	slices.Sort(names)
	return strings.Join(names, " ")
}

func TestResolveQuerySpanLineage(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		// Unquoted names fold to lower case.
		{"SELECT ID, Name FROM T", []string{"id: db.public.t.id", "name: db.public.t.name"}},
		{"SELECT count(*), upper(name), id + 1 FROM t", []string{"count:", "upper: db.public.t.name", "?column?: db.public.t.id"}},
		// * leaves out a hidden column, which a name still reads.
		{"SELECT * FROM u", []string{"id: db.public.u.id", "t_id: db.public.u.t_id", "note: db.public.u.note"}},
		{"SELECT rowid FROM u", []string{"rowid: db.public.u.rowid"}},
		{"SELECT * FROM t JOIN u USING (id)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "t_id: db.public.u.t_id", "note: db.public.u.note"}},
		{"SELECT id FROM t FULL JOIN u USING (id)", []string{"id: db.public.t.id db.public.u.id"}},
		// A view names its own column and the columns behind it.
		{"SELECT ssn FROM v", []string{"ssn: db.public.t.ssn db.public.v.ssn"}},
		{"SELECT * FROM v2", []string{"s: db.public.t.ssn db.public.v.ssn db.public.v2.s", "note: db.public.u.note db.public.v2.note"}},
		{"WITH c(k) AS (SELECT ssn FROM t) SELECT k FROM c", []string{"k: db.public.t.ssn"}},
		{"WITH RECURSIVE r AS (SELECT id, name FROM t UNION ALL SELECT u.id, r.name FROM u JOIN r ON r.id = u.t_id) SELECT * FROM r", []string{"id: db.public.t.id db.public.u.id", "name: db.public.t.name"}},
		{"SELECT id FROM t UNION SELECT t_id FROM u", []string{"id: db.public.t.id db.public.u.t_id"}},
		{"SELECT name FROM t EXCEPT SELECT note FROM u", []string{"name: db.public.t.name"}},
		{"SELECT name FROM t INTERSECT SELECT note FROM u", []string{"name: db.public.t.name"}},
		{"SELECT (SELECT max(ssn) FROM t), (t).name FROM t", []string{"max: db.public.t.ssn", "name: db.public.t.name"}},
		{"SELECT row_to_json(t) FROM t", []string{"row_to_json: db.public.t.id db.public.t.name db.public.t.ssn"}},
		{"SELECT x.n FROM generate_series(1, 2) AS x(n)", []string{"n:"}},
		// A name a function's unknown columns may hold reads the function's
		// arguments too.
		{"SELECT (SELECT ssn FROM json_each(t.name::JSONB)) FROM t", []string{"ssn: db.public.t.name db.public.t.ssn"}},
		{"SELECT o.ordinality, o.name FROM t WITH ORDINALITY AS o", []string{"ordinality:", "name: db.public.t.name"}},
		{"SELECT a FROM other.x", []string{"a: other.public.x.a"}},
		// CockroachDB's own syntax.
		{"SELECT name FROM t@t_pkey AS OF SYSTEM TIME '-1s'", []string{"name: db.public.t.name"}},
		{"SELECT s.ssn FROM t@{FORCE_INDEX=t_pkey} AS s INNER HASH JOIN u ON u.t_id = s.id", []string{"ssn: db.public.t.ssn"}},
		{"SELECT * FROM [SELECT ssn FROM t]", []string{"ssn: db.public.t.ssn"}},
		{"SELECT IF(id > 1, ssn, name) FROM t", []string{"if: db.public.t.id db.public.t.name db.public.t.ssn"}},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolveSpan(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if span.Type != analysis.Select {
				t.Errorf("type %s, want SELECT", span.Type)
			}
			if diff := cmp.Diff(tc.want, lineage(span)); diff != "" {
				t.Errorf("lineage (-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveQuerySpanTables(t *testing.T) {
	for _, tc := range []struct {
		statement  string
		queryType  analysis.QueryType
		tables     string
		predicates string
	}{
		// A view is read; what it reads is not.
		{"SELECT * FROM v2", analysis.Select, "db.public.v2", ""},
		{"WITH t AS (SELECT 1 AS id) SELECT id FROM t", analysis.Select, "", ""},
		{"SELECT id FROM t WHERE ssn = 'x' LIMIT (SELECT count(*) FROM hr.t)", analysis.Select, "db.hr.t db.public.t", "db.public.t.ssn"},
		{"SELECT t.id FROM t JOIN u ON u.t_id = t.id GROUP BY t.id HAVING max(u.note) > ''", analysis.Select, "db.public.t db.public.u", "db.public.t.id db.public.u.note db.public.u.t_id"},
		{"SELECT count(*) FILTER (WHERE name = 'a') FROM t", analysis.Select, "db.public.t", "db.public.t.name"},
		{"SELECT id FROM t WHERE EXISTS (SELECT 1 FROM u WHERE u.t_id = t.id)", analysis.Select, "db.public.t db.public.u", "db.public.t.id db.public.u.t_id"},
		{"SELECT id FROM t JOIN u USING (id)", analysis.Select, "db.public.t db.public.u", "db.public.t.id db.public.u.id"},
		{"SELECT name FROM t INTERSECT SELECT note FROM u", analysis.Select, "db.public.t db.public.u", "db.public.t.name db.public.u.note"},
		{"SELECT * FROM pg_class", analysis.SelectInfoSchema, "db.pg_catalog.pg_class", ""},
		{"SELECT * FROM information_schema.tables", analysis.SelectInfoSchema, "db.information_schema.tables", ""},
		{"SELECT * FROM crdb_internal.jobs", analysis.SelectInfoSchema, "db.crdb_internal.jobs", ""},
		{"SELECT * FROM system.users", analysis.SelectInfoSchema, "system.public.users", ""},
		{"SELECT * FROM [SHOW TABLES]", analysis.SelectInfoSchema, "", ""},
		{"SHOW COLUMNS FROM t", analysis.SelectInfoSchema, "", ""},
		{"UPSERT INTO t (id) SELECT t_id FROM u", analysis.DML, "db.public.t db.public.u", ""},
		// The target is the table, though a CTE has its name.
		{"WITH t AS (SELECT 1 AS id) INSERT INTO t (id) SELECT id FROM t", analysis.DML, "db.public.t", ""},
		{"UPDATE t SET name = u.note FROM u WHERE u.t_id = t.id", analysis.DML, "db.public.t db.public.u", ""},
		{"DELETE FROM missing WHERE id = 1", analysis.DML, "db.public.missing", ""},
		{"SELECT * FROM [INSERT INTO t (id) VALUES (1) RETURNING id]", analysis.DML, "db.public.t", ""},
		{"WITH d AS (DELETE FROM u RETURNING t_id) SELECT * FROM t", analysis.DML, "db.public.t db.public.u", ""},
		{"EXPLAIN (OPT) SELECT name FROM t", analysis.Explain, "db.public.t", ""},
		{"SELECT * FROM [EXPLAIN SELECT name FROM t]", analysis.Select, "db.public.t", ""},
		{"CREATE TABLE c AS SELECT ssn FROM t", analysis.DDL, "db.public.t", ""},
		{"SET search_path = hr", analysis.Select, "", ""},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolveSpan(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if span.Type != tc.queryType {
				t.Errorf("type %s, want %s", span.Type, tc.queryType)
			}
			if got := columnNames(span.SourceColumns); got != tc.tables {
				t.Errorf("tables %q, want %q", got, tc.tables)
			}
			if got := columnNames(span.PredicateColumns); got != tc.predicates {
				t.Errorf("predicates %q, want %q", got, tc.predicates)
			}
			if span.Type != analysis.Select && len(span.Results) != 0 {
				t.Errorf("results %v for a %s", lineage(span), span.Type)
			}
		})
	}
}

func TestResolveQuerySpanExplainAnalyze(t *testing.T) {
	span, err := resolveSpan("EXPLAIN ANALYZE SELECT name FROM t WHERE ssn = 'x'")
	if err != nil {
		t.Fatal(err)
	}
	if span.Type != analysis.Select || !span.ExplainAnalyze || len(span.Results) != 0 {
		t.Errorf("got type %s, ExplainAnalyze %v, results %v", span.Type, span.ExplainAnalyze, lineage(span))
	}
	if got := columnNames(span.PredicateColumns); got != "db.public.t.ssn" {
		t.Errorf("predicates %q", got)
	}

	span, err = resolveSpan("EXPLAIN ANALYZE UPSERT INTO t (id) VALUES (1)")
	if err != nil {
		t.Fatal(err)
	}
	if span.Type != analysis.DML || !span.ExplainAnalyze {
		t.Errorf("got type %s, ExplainAnalyze %v", span.Type, span.ExplainAnalyze)
	}
}

func TestResolveQuerySpanErrors(t *testing.T) {
	_, err := resolveSpan("SELECT FROM WHERE")
	var syntaxErr *SyntaxError
	if !errors.As(err, &syntaxErr) || syntaxErr.Position != 12 {
		t.Errorf("got %v, want a *SyntaxError at 12", err)
	}

	for _, statement := range []string{
		"SELECT * FROM t, pg_catalog.pg_class",
		"SELECT * FROM t, [SHOW TABLES]",
		"SELECT id FROM t WHERE id IN (SELECT oid FROM pg_class)",
	} {
		if _, err := resolveSpan(statement); !errors.Is(err, ErrMixUserSystemTables) {
			t.Errorf("%s: got %v, want ErrMixUserSystemTables", statement, err)
		}
	}

	for statement, code := range map[string]pgcode.Code{
		"SELECT nope FROM t":                     pgcode.UndefinedColumn,
		"SELECT id FROM t JOIN u ON true":        pgcode.AmbiguousColumn,
		"SELECT * FROM t JOIN u USING (rowid)":   pgcode.UndefinedColumn,
		"SELECT x.id FROM t":                     pgcode.UndefinedTable,
		"SELECT \"Name\" FROM t":                 pgcode.UndefinedColumn,
		"SELECT name FROM t ORDER BY nope":       pgcode.Uncategorized,
		"SELECT name FROM t GROUP BY name, nope": pgcode.Uncategorized,
	} {
		_, err := resolveSpan(statement)
		switch {
		case code == pgcode.Uncategorized && err != nil:
			t.Errorf("%s: %v; a name that decides neither values nor rows need not resolve", statement, err)
		case code != pgcode.Uncategorized && pgerror.GetPGCode(err) != code:
			t.Errorf("%s: got %v, want SQLSTATE %s", statement, err, code)
		}
	}
}

func TestResolveQuerySpanMissingMetadata(t *testing.T) {
	// A quoted name keeps its case, so "T" is not t.
	span, err := resolveSpan(`SELECT id FROM t JOIN "T" ON true`)
	if err != nil {
		t.Fatal(err)
	}
	nf := span.NotFoundError
	if nf == nil || *nf.Schema != "public" || *nf.Table != "T" || len(span.Results) != 0 {
		t.Fatalf("got NotFoundError %v and results %v", nf, lineage(span))
	}
	if got := columnNames(span.SourceColumns); got != "db.public.T db.public.t" {
		t.Errorf("tables %q", got)
	}

	span, err = resolveSpan("SELECT * FROM nowhere.x")
	if err != nil {
		t.Fatal(err)
	}
	if nf := span.NotFoundError; nf == nil || *nf.Database != "nowhere" {
		t.Errorf("got NotFoundError %v", nf)
	}

	// A table the sync stored no columns for reads any name.
	span, err = resolveSpan("SELECT a, b FROM empty")
	if err != nil {
		t.Fatal(err)
	}
	if span.UnresolvedColumnsError == nil || span.UnresolvedColumnsError.Error() != "the synced schema describes no columns for db.public.empty" {
		t.Errorf("got UnresolvedColumnsError %v", span.UnresolvedColumnsError)
	}
	if diff := cmp.Diff([]string{"a: db.public.empty.a", "b: db.public.empty.b"}, lineage(span)); diff != "" {
		t.Errorf("lineage (-want +got):\n%s", diff)
	}
}
