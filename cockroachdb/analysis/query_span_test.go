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
		{"SELECT name FROM t INTERSECT SELECT note FROM u", []string{"name: db.public.t.name db.public.u.note"}},
		{"SELECT (SELECT max(ssn) FROM t), (t).name FROM t", []string{"max: db.public.t.ssn", "name: db.public.t.name"}},
		{"SELECT row_to_json(t) FROM t", []string{"row_to_json: db.public.t.id db.public.t.name db.public.t.ssn"}},
		{"SELECT x.n FROM generate_series(1, 2) AS x(n)", []string{"n:"}},
		// A name a function's unknown columns may hold reads the function's
		// arguments too.
		{"SELECT (SELECT ssn FROM f(t.name)) FROM t", []string{"ssn: db.public.t.name db.public.t.ssn"}},
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

	// The conditions of a mutation decide which rows it changes, which its
	// plan reports.
	for statement, predicates := range map[string]string{
		"EXPLAIN ANALYZE UPDATE t SET name = 'x' WHERE ssn = 'guess'":                                                                         "db.public.t.ssn",
		"EXPLAIN ANALYZE DELETE FROM t WHERE ssn = 'guess'":                                                                                   "db.public.t.ssn",
		"EXPLAIN ANALYZE UPDATE t SET name = u.note FROM u WHERE u.t_id = t.id":                                                               "db.public.t.id db.public.u.t_id",
		"EXPLAIN ANALYZE DELETE FROM t USING u JOIN hr.t AS h ON h.id = u.id WHERE u.t_id = t.id":                                             "db.hr.t.id db.public.t.id db.public.u.id db.public.u.t_id",
		"EXPLAIN ANALYZE DELETE FROM t WHERE id > 1 LIMIT (SELECT count(*) FROM u WHERE note = 'n')":                                          "db.public.t.id db.public.u.note",
		"EXPLAIN ANALYZE INSERT INTO t (id) SELECT t_id FROM u WHERE note = 'n'":                                                              "db.public.u.note",
		"EXPLAIN ANALYZE INSERT INTO t (id) VALUES (1) ON CONFLICT (id) WHERE ssn > '' DO UPDATE SET name = excluded.name WHERE t.name < 'm'": "db.public.t.name db.public.t.ssn",
	} {
		span, err := resolveSpan(statement)
		if err != nil {
			t.Errorf("%s: %v", statement, err)
			continue
		}
		if got := columnNames(span.PredicateColumns); got != predicates {
			t.Errorf("%s: predicates %q, want %q", statement, got, predicates)
		}
	}
	// Without ANALYZE a mutation's span holds no conditions.
	if span, err := resolveSpan("UPDATE t SET name = 'x' WHERE ssn = 'guess'"); err != nil || len(span.PredicateColumns) != 0 {
		t.Errorf("got predicates %v, error %v", span.PredicateColumns, err)
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
		"SELECT name FROM t ORDER BY nope":       pgcode.UndefinedColumn,
		"SELECT name FROM t GROUP BY name, nope": pgcode.UndefinedColumn,
		"SELECT * FROM t AS x(a, b, c, d)":       pgcode.InvalidColumnReference,
		// A non-LATERAL subquery in FROM does not see the items before it.
		"SELECT * FROM t, (SELECT t.id) AS s": pgcode.UndefinedTable,
	} {
		if _, err := resolveSpan(statement); pgerror.GetPGCode(err) != code {
			t.Errorf("%s: got %v, want SQLSTATE %s", statement, err, code)
		}
	}
}

// TestResolveQuerySpanDerivedColumns holds every column a relation or a query
// makes to the sources it derives from.
func TestResolveQuerySpanDerivedColumns(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		// Each column of a function in FROM reads the function's arguments,
		// those its alias names beyond the known ones too.
		{"SELECT j.value FROM t, json_each(t.name::JSONB) AS j(key, value)", []string{"value: db.public.t.name"}},
		{"SELECT * FROM t AS x, LATERAL json_each(x.ssn::JSONB) WITH ORDINALITY AS j(k, v)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "k: db.public.t.ssn", "v: db.public.t.ssn", "ordinality:"}},
		{"SELECT r.a, r.c FROM t, ROWS FROM (generate_series(1, 2), json_each(t.name::JSONB)) AS r(a, b, c)", []string{"a:", "c: db.public.t.name"}},
		// A function of known columns places each one.
		{"SELECT r.b FROM t, ROWS FROM (json_each(t.ssn::JSONB), generate_series(1, 2)) AS r(a, b, c)", []string{"b: db.public.t.ssn"}},
		{"SELECT r.c FROM t, ROWS FROM (json_each(t.ssn::JSONB), generate_series(1, 2)) AS r(a, b, c)", []string{"c:"}},
		{"SELECT * FROM t AS x, json_each(x.ssn::JSONB)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "key: db.public.t.ssn", "value: db.public.t.ssn"}},
		{"SELECT * FROM t AS x, unnest(ARRAY[x.ssn], ARRAY[x.name]) AS y(p, q)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "p: db.public.t.ssn", "q: db.public.t.name"}},
		{"SELECT (SELECT ssn FROM json_each(t.name::JSONB)) FROM t", []string{"ssn: db.public.t.ssn"}},
		// Past a function of unknown columns, every column reads every
		// argument of the item.
		{"SELECT r.b FROM t, ROWS FROM (f(t.ssn), generate_series(1, t.id)) AS r(a, b)", []string{"b: db.public.t.id db.public.t.ssn"}},
		{"SELECT r.generate_series FROM t, ROWS FROM (generate_series(1, t.id), f(t.ssn)) AS r", []string{"generate_series: db.public.t.id db.public.t.ssn"}},
		// A name an aggregate, a set operation or a values list makes reads
		// every input.
		{"SELECT x FROM (SELECT ssn FROM t UNION ALL SELECT note FROM u) AS q(x)", []string{"x: db.public.t.ssn db.public.u.note"}},
		{"SELECT column2 FROM (VALUES (1, (SELECT ssn FROM t LIMIT 1)))", []string{"column2: db.public.t.ssn"}},
		// A correlated name in a subquery of an unknown relation keeps the
		// outer column and the relation's reads.
		{"SELECT (SELECT ssn FROM f(t.name)) FROM t", []string{"ssn: db.public.t.name db.public.t.ssn"}},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolveSpan(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.want, lineage(span)); diff != "" {
				t.Errorf("lineage (-want +got):\n%s", diff)
			}
		})
	}
}

// TestResolveQuerySpanEveryClause holds that every clause's relations are
// read, in the scope the clause sees, and that a name a query cannot resolve
// fails it only after every relation it names is recorded.
func TestResolveQuerySpanEveryClause(t *testing.T) {
	for _, tc := range []struct {
		statement string
		queryType analysis.QueryType
		tables    string
	}{
		// ORDER BY sees the output columns, then the FROM items.
		{"SELECT t.name FROM t ORDER BY (SELECT t.id + (SELECT count(*) FROM u))", analysis.Select, "db.public.t db.public.u"},
		{"SELECT name AS n FROM t ORDER BY n, (SELECT count(*) FROM hr.t WHERE hr.t.id = t.id)", analysis.Select, "db.hr.t db.public.t"},
		{"SELECT name AS n FROM t GROUP BY n HAVING count(*) > (SELECT count(*) FROM u)", analysis.Select, "db.public.t db.public.u"},
		{"SELECT DISTINCT ON (n) name AS n FROM t ORDER BY n", analysis.Select, "db.public.t"},
		{"SELECT name FROM t UNION SELECT note FROM u ORDER BY 1 LIMIT (SELECT count(*) FROM hr.t)", analysis.Select, "db.hr.t db.public.t db.public.u"},
		{"SELECT rank() OVER (PARTITION BY (SELECT max(note) FROM u) ORDER BY id) FROM t", analysis.Select, "db.public.t db.public.u"},
		{"SELECT name FROM t WINDOW w AS (ORDER BY (SELECT max(note) FROM u))", analysis.Select, "db.public.t db.public.u"},
		{"SELECT id FROM t AS OF SYSTEM TIME (SELECT '-1s' FROM u LIMIT 1)", analysis.Select, "db.public.t db.public.u"},
		{"SELECT id FROM t, LATERAL (SELECT note FROM u WHERE u.t_id = t.id) AS l", analysis.Select, "db.public.t db.public.u"},
		// UPDATE and DELETE read their ORDER BY, LIMIT and BATCH.
		{"DELETE FROM t LIMIT (SELECT count(*) FROM u)", analysis.DML, "db.public.t db.public.u"},
		{"DELETE FROM t ORDER BY (SELECT max(note) FROM u) LIMIT 1", analysis.DML, "db.public.t db.public.u"},
		{"DELETE BATCH (SIZE (SELECT count(*) FROM u)) FROM t", analysis.DML, "db.public.t db.public.u"},
		{"UPDATE t SET name = 'x' ORDER BY (SELECT max(note) FROM u) LIMIT (SELECT count(*) FROM hr.t)", analysis.DML, "db.hr.t db.public.t db.public.u"},
		// INSERT and UPSERT read their rows, ON CONFLICT and RETURNING.
		{"INSERT INTO t (id) VALUES (1) ON CONFLICT (id) WHERE (SELECT true FROM u LIMIT 1) DO UPDATE SET name = excluded.name WHERE t.ssn = (SELECT max(salary) FROM hr.t)::STRING", analysis.DML, "db.hr.t db.public.t db.public.u"},
		{"UPSERT INTO t (id) SELECT t_id FROM u RETURNING (SELECT max(salary) FROM hr.t)", analysis.DML, "db.hr.t db.public.t db.public.u"},
		{"COPY t FROM STDIN", analysis.DML, "db.public.t"},
		{"COPY (SELECT ssn FROM t) TO STDOUT", analysis.DML, "db.public.t"},
		// A name a mutation cannot resolve leaves its relations whole.
		{"DELETE FROM t WHERE nope = (SELECT count(*) FROM u)", analysis.DML, "db.public.t db.public.u"},
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
		})
	}

	// A nested statement sees the CTEs around it; a mutation's target is a
	// table whatever CTE has its name.
	for statement, tables := range map[string]string{
		"WITH c AS (SELECT id FROM u), d AS (DELETE FROM t WHERE id IN (SELECT id FROM c) RETURNING id) SELECT * FROM d":  "db.public.t db.public.u",
		"WITH c AS (SELECT id FROM u) SELECT * FROM [DELETE FROM t WHERE id IN (SELECT id FROM c) RETURNING id]":          "db.public.t db.public.u",
		"WITH c AS (SELECT id FROM u) SELECT * FROM [EXPLAIN DELETE FROM t WHERE id IN (SELECT id FROM c)]":               "db.public.t db.public.u",
		"WITH c AS (SELECT id FROM u) SELECT * FROM [UPSERT INTO t (id) SELECT id FROM c RETURNING id]":                   "db.public.t db.public.u",
		"WITH c AS (SELECT id FROM u), d AS (UPDATE t SET name = 'x' FROM c WHERE c.id = t.id RETURNING t.id) SELECT 1":   "db.public.t db.public.u",
		"WITH t AS (SELECT t_id AS id FROM u), d AS (DELETE FROM t WHERE id IN (SELECT id FROM t) RETURNING id) SELECT 1": "db.public.t db.public.u",
		"CREATE TABLE x AS WITH c AS (SELECT ssn FROM t) SELECT * FROM c":                                                 "db.public.t",
		"EXPLAIN WITH c AS (SELECT id FROM u) DELETE FROM t WHERE id IN (SELECT id FROM c)":                               "db.public.t db.public.u",
	} {
		span, err := resolveSpan(statement)
		if err != nil {
			t.Errorf("%s: %v", statement, err)
			continue
		}
		if got := columnNames(span.SourceColumns); got != tables {
			t.Errorf("%s: tables %q, want %q", statement, got, tables)
		}
	}

	// A system relation in any clause makes the statement mixed.
	for _, statement := range []string{
		"SELECT name FROM t ORDER BY (SELECT count(*) FROM pg_class)",
		"DELETE FROM t LIMIT (SELECT count(*) FROM pg_class)",
		"UPDATE t SET name = 'x' ORDER BY (SELECT max(relname) FROM pg_class)",
		// A name that fails to resolve before the system read does not hide it.
		"SELECT nope, (SELECT count(*) FROM crdb_internal.jobs) FROM t",
	} {
		if _, err := resolveSpan(statement); !errors.Is(err, ErrMixUserSystemTables) {
			t.Errorf("%s: got %v, want ErrMixUserSystemTables", statement, err)
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
