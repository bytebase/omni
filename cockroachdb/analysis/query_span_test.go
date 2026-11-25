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
					// User tables named as system relations.
					{Name: "pg_class", Columns: []*metadatapb.ColumnMetadata{column("relname")}},
					{Name: "geometry_columns", Columns: []*metadatapb.ColumnMetadata{column("f_table_name")}},
				},
				Functions: []*metadatapb.FunctionMetadata{
					sqlFunction("get_ssn()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT ssn FROM db.public.t LIMIT 1 $$"),
					sqlFunction("note_of(i INT8)", "RETURNS STRING LANGUAGE SQL AS $$ SELECT note FROM db.public.u WHERE t_id = i $$"),
					sqlFunction("rows_of()", "RETURNS SETOF db.public.u LANGUAGE SQL AS $$ SELECT * FROM db.public.u $$"),
					sqlFunction("two_steps()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT 1 FROM db.public.u; SELECT name FROM db.public.t LIMIT 1 $$"),
					sqlFunction("atomic_salary()", "RETURNS DECIMAL LANGUAGE SQL BEGIN ATOMIC SELECT salary FROM db.hr.t LIMIT 1; END"),
					sqlFunction("pl_name()", "RETURNS STRING LANGUAGE PLpgSQL AS $$ DECLARE x STRING; BEGIN SELECT name INTO x FROM db.public.t LIMIT 1; RETURN x; END $$"),
					sqlFunction("dynamic()", "RETURNS STRING LANGUAGE PLpgSQL AS $$ BEGIN EXECUTE 'SELECT ssn FROM db.public.t'; RETURN ''; END $$"),
					sqlFunction("recurse()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT db.public.recurse() || ssn FROM db.public.t LIMIT 1 $$"),
					sqlFunction("relnames()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT relname FROM pg_catalog.pg_class LIMIT 1 $$"),
					sqlFunction("writer()", "RETURNS INT8 LANGUAGE SQL AS $$ INSERT INTO db.public.u (id) VALUES (1) RETURNING id $$"),
					{Name: "broken", Signature: "broken()", Definition: "not a function"},
					// Each PL/pgSQL form that returns what a query reads.
					plpgsqlFunction("star_into()", "DECLARE x RECORD; BEGIN SELECT * INTO x FROM db.public.t LIMIT 1; RETURN x; END"),
					plpgsqlFunction("return_query()", "BEGIN RETURN QUERY SELECT * FROM db.public.u; END"),
					plpgsqlFunction("bound_cursor()", "DECLARE c CURSOR FOR SELECT * FROM db.public.u; x RECORD; BEGIN OPEN c; FETCH c INTO x; RETURN x; END"),
					plpgsqlFunction("open_for()", "DECLARE c REFCURSOR; x RECORD; BEGIN OPEN c FOR SELECT * FROM db.public.u; FETCH c INTO x; RETURN x; END"),
					plpgsqlFunction("assigned()", "DECLARE x STRING; BEGIN x := (SELECT ssn FROM db.public.t LIMIT 1); RETURN x; END"),
					plpgsqlFunction("return_next()", "BEGIN RETURN NEXT (SELECT name FROM db.public.t LIMIT 1); END"),
					plpgsqlFunction("returning_into()", "DECLARE x RECORD; BEGIN INSERT INTO db.public.u (id) VALUES (1) RETURNING * INTO x; RETURN x; END"),
					plpgsqlFunction("diagnostics()", "DECLARE n INT8; BEGIN GET DIAGNOSTICS n = ROW_COUNT; RETURN n; END"),
					plpgsqlFunction("query_execute()", "BEGIN RETURN QUERY EXECUTE 'SELECT ssn FROM db.public.t'; END"),
					plpgsqlFunction("query_loop()", "DECLARE r RECORD; BEGIN FOR r IN SELECT * FROM db.public.t LOOP RETURN NEXT r; END LOOP; END"),
					plpgsqlFunction("performs()", "BEGIN PERFORM ssn FROM db.public.t; END"),
					sqlFunction("copied()", "RETURNS INT8 LANGUAGE SQL AS $$ INSERT INTO db.public.u (id) SELECT * FROM db.public.t RETURNING 1 $$"),
					sqlFunction("deleted()", "RETURNS SETOF db.public.u LANGUAGE SQL AS $$ DELETE FROM db.public.u RETURNING * $$"),
					// An overload whose arguments a later schema's overload matches.
					// Mutually recursive functions: f reads t.ssn at its base
					// case, g only through f; and a three-function cycle.
					sqlFunction("f(n INT8)", "RETURNS STRING LANGUAGE PLpgSQL AS $$ BEGIN IF n = 0 THEN RETURN (SELECT ssn FROM db.public.t LIMIT 1); END IF; RETURN db.public.g(n - 1); END $$"),
					sqlFunction("g(n INT8)", "RETURNS STRING LANGUAGE PLpgSQL AS $$ BEGIN RETURN db.public.f(n - 1); END $$"),
					sqlFunction("cycle_a()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT db.public.cycle_b() $$"),
					sqlFunction("cycle_b()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT db.public.cycle_c() $$"),
					sqlFunction("cycle_c()", "RETURNS STRING LANGUAGE SQL AS $$ SELECT db.public.cycle_a() || name FROM db.public.t LIMIT 1 $$"),
					// Synced functions named as built-ins, whose columns are
					// the ones they declare.
					sqlFunction("unnest(x STRING)", "RETURNS TABLE(a STRING, b STRING) LANGUAGE SQL AS $$ SELECT x, x $$"),
					sqlFunction("table_rows(x STRING)", "RETURNS SETOF db.public.t LANGUAGE SQL AS $$ SELECT * FROM db.public.t $$"),
					sqlFunction("one_column()", "RETURNS TABLE(a STRING) LANGUAGE SQL AS $$ SELECT name FROM db.public.t $$"),
					sqlFunction("pairs(x STRING, OUT k STRING, OUT v STRING)", "RETURNS RECORD LANGUAGE SQL AS $$ SELECT x, x $$"),
					sqlFunction("anything(x STRING)", "RETURNS SETOF RECORD LANGUAGE SQL AS $$ SELECT x, x $$"),
					sqlFunction("pick(i INT8)", "RETURNS STRING LANGUAGE SQL AS $$ SELECT ssn FROM db.public.t LIMIT 1 $$"),
				},
				Views: []*metadatapb.ViewMetadata{
					// A view whose sync stored no columns.
					{Name: "v", Definition: "SELECT ssn FROM public.t"},
					{Name: "v2", Definition: "SELECT ssn AS s, note FROM db.public.v JOIN db.public.u ON true", Columns: []*metadatapb.ColumnMetadata{column("s"), column("note")}},
					// Synced names differ from the definition's, and in
					// renamed_short also in number.
					{Name: "renamed", Definition: "SELECT ssn, name FROM db.public.t", Columns: []*metadatapb.ColumnMetadata{column("a"), column("b")}},
					{Name: "renamed_short", Definition: "SELECT ssn, name FROM db.public.t", Columns: []*metadatapb.ColumnMetadata{column("joined")}},
				},
				MaterializedViews: []*metadatapb.MaterializedViewMetadata{
					// The sync stores no columns of a materialized view.
					{Name: "region_counts", Definition: "SELECT name AS region, count(*) AS n FROM db.public.t GROUP BY name"},
					{Name: "synced_mv", Definition: "SELECT x FROM db.public.gone", Columns: []*metadatapb.ColumnMetadata{column("r"), column("n"), {Name: "rowid", IsInvisible: true}}},
				},
			},
			{
				Name:   "hr",
				Tables: []*metadatapb.TableMetadata{{Name: "t", Columns: []*metadatapb.ColumnMetadata{column("id"), column("salary")}}},
				Functions: []*metadatapb.FunctionMetadata{
					sqlFunction("get_ssn()", "RETURNS DECIMAL LANGUAGE SQL AS $$ SELECT salary FROM db.hr.t LIMIT 1 $$"),
					sqlFunction("pick(s STRING)", "RETURNS DECIMAL LANGUAGE SQL AS $$ SELECT salary FROM db.hr.t LIMIT 1 $$"),
				},
			},
		}},
		"other": {Name: "other", Schemas: []*metadatapb.SchemaMetadata{
			{Name: "public", Tables: []*metadatapb.TableMetadata{{Name: "x", Columns: []*metadatapb.ColumnMetadata{column("a")}}}},
		}},
	}
}

// sqlFunction is a function as the sync stores it: its signature, and
// pg_get_functiondef's CREATE FUNCTION in its schema.
func sqlFunction(signature, rest string) *metadatapb.FunctionMetadata {
	name, _, _ := strings.Cut(signature, "(")
	return &metadatapb.FunctionMetadata{Name: name, Signature: signature, Definition: "CREATE FUNCTION public." + signature + " " + rest}
}

func plpgsqlFunction(signature, body string) *metadatapb.FunctionMetadata {
	return sqlFunction(signature, "RETURNS RECORD LANGUAGE PLpgSQL AS $$ "+body+" $$")
}

func resolveSpan(statement string) (*analysis.ResolvedSpan, error) {
	return resolveSpanIn(statement, "public")
}

// resolveSpanIn resolves statement under searchPath, a search_path value.
func resolveSpanIn(statement, searchPath string) (*analysis.ResolvedSpan, error) {
	databases := spanTestMetadata()
	return ResolveQuerySpan(context.Background(), statement, ResolveOptions{
		DefaultDatabase: "db",
		SearchPath:      strings.Split(searchPath, ","),
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
		{"WITH RECURSIVE c AS (WITH x AS (SELECT t_id FROM u) SELECT id FROM t UNION ALL SELECT t_id FROM x) SELECT * FROM c", []string{"id: db.public.t.id db.public.u.t_id"}},
		{"SELECT id FROM t UNION SELECT t_id FROM u", []string{"id: db.public.t.id db.public.u.t_id"}},
		{"SELECT name FROM t EXCEPT SELECT note FROM u", []string{"name: db.public.t.name"}},
		{"SELECT name FROM t INTERSECT SELECT note FROM u", []string{"name: db.public.t.name db.public.u.note"}},
		{"SELECT (SELECT max(ssn) FROM t), (t).name FROM t", []string{"max: db.public.t.ssn", "name: db.public.t.name"}},
		{"SELECT row_to_json(t) FROM t", []string{"row_to_json: db.public.t.id db.public.t.name db.public.t.ssn"}},
		{"SELECT x.n FROM generate_series(1, 2) AS x(n)", []string{"n:"}},
		// A name a function's unknown columns may hold reads the function's
		// arguments too.
		{"SELECT (SELECT ssn FROM unknown_fn(t.name)) FROM t", []string{"ssn: db.public.t.name db.public.t.ssn"}},
		{"SELECT o.ordinality, o.name FROM t WITH ORDINALITY AS o", []string{"ordinality:", "name: db.public.t.name"}},
		// A tuple's fields are columns of their own; a value of unknown
		// fields stays one column.
		{"SELECT (ROW(t.id, t.ssn)).* FROM t", []string{"?column?: db.public.t.id", "?column?: db.public.t.ssn"}},
		{"SELECT ((t.id, upper(t.ssn))).* FROM t", []string{"?column?: db.public.t.id", "?column?: db.public.t.ssn"}},
		{"SELECT (((id, ssn) AS k, v)).* FROM t", []string{"k: db.public.t.id", "v: db.public.t.ssn"}},
		{"SELECT b FROM (SELECT (ROW(t.id, t.ssn)).* FROM t) AS q(a, b)", []string{"b: db.public.t.ssn"}},
		{"SELECT (unknown_fn(t.name)).* FROM t", []string{"?column?: db.public.t.name"}},
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
		// WITH RECURSIVE over a set operation that is not a UNION walks it
		// as one.
		{"WITH RECURSIVE c AS (SELECT name FROM t EXCEPT SELECT note FROM u) SELECT * FROM c", analysis.Select, "db.public.t db.public.u", "db.public.t.name db.public.u.note"},
		{"WITH RECURSIVE c AS (SELECT id FROM t UNION ALL SELECT id FROM t LIMIT (SELECT max(t_id) FROM u)) SELECT * FROM c", analysis.Select, "db.public.t db.public.u", "db.public.u.t_id"},
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
		{"SELECT * FROM t AS x, pg_catalog.unnest(ARRAY[x.ssn], ARRAY[x.name]) AS y(p, q)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "p: db.public.t.ssn", "q: db.public.t.name"}},
		{"SELECT (SELECT ssn FROM json_each(t.name::JSONB)) FROM t", []string{"ssn: db.public.t.ssn"}},
		// Past a function of unknown columns, every column reads every
		// argument of the item.
		{"SELECT r.b FROM t, ROWS FROM (unknown_fn(t.ssn), generate_series(1, t.id)) AS r(a, b)", []string{"b: db.public.t.id db.public.t.ssn"}},
		{"SELECT r.generate_series FROM t, ROWS FROM (generate_series(1, t.id), unknown_fn(t.ssn)) AS r", []string{"generate_series: db.public.t.id db.public.t.ssn"}},
		// A name an aggregate, a set operation or a values list makes reads
		// every input.
		{"SELECT x FROM (SELECT ssn FROM t UNION ALL SELECT note FROM u) AS q(x)", []string{"x: db.public.t.ssn db.public.u.note"}},
		{"SELECT column2 FROM (VALUES (1, (SELECT ssn FROM t LIMIT 1)))", []string{"column2: db.public.t.ssn"}},
		// A correlated name in a subquery of an unknown relation keeps the
		// outer column and the relation's reads.
		{"SELECT (SELECT ssn FROM unknown_fn(t.name)) FROM t", []string{"ssn: db.public.t.name db.public.t.ssn"}},
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

// TestResolveQuerySpanStoredColumns holds a relation's columns to the synced
// metadata, and a materialized view's to its own stored columns.
func TestResolveQuerySpanStoredColumns(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		// A materialized view the sync stored no columns for has its
		// definition's names and CockroachDB's hidden rowid.
		{"SELECT rowid FROM public.region_counts", []string{"rowid: db.public.region_counts.rowid"}},
		{"SELECT * FROM region_counts", []string{"region: db.public.region_counts.region", "n: db.public.region_counts.n"}},
		// A synced materialized view needs no definition.
		{"SELECT * FROM synced_mv", []string{"r: db.public.synced_mv.r", "n: db.public.synced_mv.n"}},
		{"SELECT rowid, r FROM synced_mv", []string{"rowid: db.public.synced_mv.rowid", "r: db.public.synced_mv.r"}},
		// A view's synced names stand; where their number differs from the
		// definition's, each reads everything the definition reads.
		{"SELECT a, b FROM renamed", []string{"a: db.public.renamed.a db.public.t.ssn", "b: db.public.renamed.b db.public.t.name"}},
		{"SELECT joined FROM renamed_short", []string{"joined: db.public.renamed_short.joined db.public.t.name db.public.t.ssn"}},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolveSpan(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if span.NotFoundError != nil {
				t.Fatalf("NotFoundError %v", span.NotFoundError)
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
		// A recursive CTE's body is walked whole: its WITH, ORDER BY and
		// LIMIT too.
		{"WITH RECURSIVE c AS (SELECT id FROM t UNION ALL SELECT id FROM t LIMIT (SELECT count(*) FROM u)) SELECT * FROM c", analysis.Select, "db.public.t db.public.u"},
		{"WITH RECURSIVE c AS (SELECT id FROM t UNION ALL SELECT c.id + 1 FROM c ORDER BY (SELECT max(t_id) FROM u) LIMIT 5) SELECT * FROM c", analysis.Select, "db.public.t db.public.u"},
		{"WITH RECURSIVE c AS (WITH x AS (SELECT t_id FROM u) SELECT id FROM t UNION ALL SELECT t_id FROM x) SELECT * FROM c", analysis.Select, "db.public.t db.public.u"},
		// A write anywhere a query runs one makes it DML.
		{"SELECT rank() OVER (PARTITION BY (SELECT 1 FROM [INSERT INTO u (id) VALUES (1) RETURNING id])) FROM t", analysis.DML, "db.public.t db.public.u"},
		{"SELECT name FROM t ORDER BY (SELECT max(id) FROM [DELETE FROM u RETURNING id])", analysis.DML, "db.public.t db.public.u"},
		{"SELECT * FROM [EXPLAIN DELETE FROM u]", analysis.Select, "db.public.u"},
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
		"WITH RECURSIVE c AS (SELECT id FROM t UNION ALL SELECT id FROM t LIMIT (SELECT count(*) FROM pg_class)) SELECT * FROM c",
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

// TestResolveQuerySpanFunctions holds a call of a user-defined function to
// what its synced body reads.
func TestResolveQuerySpanFunctions(t *testing.T) {
	for _, tc := range []struct {
		statement string
		queryType analysis.QueryType
		want      []string
		tables    string
	}{
		{"SELECT public.get_ssn()", analysis.Select, []string{"get_ssn: db.public.t.ssn"}, "db.public.t"},
		{"SELECT get_ssn(), get_ssn() AS again", analysis.Select, []string{"get_ssn: db.public.t.ssn", "again: db.public.t.ssn"}, "db.public.t"},
		{"SELECT hr.get_ssn()", analysis.Select, []string{"get_ssn: db.hr.t.salary"}, "db.hr.t"},
		// The body's names count, its parameters' arguments too.
		{"SELECT note_of(id) FROM t", analysis.Select, []string{"note_of: db.public.t.id db.public.u.note db.public.u.t_id"}, "db.public.t db.public.u"},
		{"SELECT upper(get_ssn())", analysis.Select, []string{"upper: db.public.t.ssn"}, "db.public.t"},
		{"SELECT rows_of()", analysis.Select, []string{"rows_of: db.public.u.id db.public.u.note db.public.u.t_id"}, "db.public.u"},
		{"SELECT x.y FROM rows_of() AS x(y, z)", analysis.Select, []string{"y: db.public.u.id db.public.u.note db.public.u.t_id"}, "db.public.u"},
		{"SELECT two_steps()", analysis.Select, []string{"two_steps: db.public.t.name"}, "db.public.t db.public.u"},
		{"SELECT atomic_salary()", analysis.Select, []string{"atomic_salary: db.hr.t.salary"}, "db.hr.t"},
		{"SELECT pl_name()", analysis.Select, []string{"pl_name: db.public.t.name"}, "db.public.t"},
		{"SELECT recurse()", analysis.Select, []string{"recurse: db.public.t.ssn"}, "db.public.t"},
		{"SELECT relnames()", analysis.SelectInfoSchema, nil, "db.pg_catalog.pg_class"},
		{"SELECT writer()", analysis.DML, nil, "db.public.u"},
		// A plain EXPLAIN plans the call without running its writes.
		{"SELECT * FROM [EXPLAIN SELECT writer()]", analysis.Select, []string{"*:"}, "db.public.u"},
		// Each query of a PL/pgSQL body returns what it reads, * too.
		{"SELECT star_into()", analysis.Select, []string{"star_into: db.public.t.id db.public.t.name db.public.t.ssn"}, "db.public.t"},
		{"SELECT return_query()", analysis.Select, []string{"return_query: db.public.u.id db.public.u.note db.public.u.t_id"}, "db.public.u"},
		{"SELECT bound_cursor()", analysis.Select, []string{"bound_cursor: db.public.u.id db.public.u.note db.public.u.t_id"}, "db.public.u"},
		{"SELECT open_for()", analysis.Select, []string{"open_for: db.public.u.id db.public.u.note db.public.u.t_id"}, "db.public.u"},
		{"SELECT assigned()", analysis.Select, []string{"assigned: db.public.t.ssn"}, "db.public.t"},
		{"SELECT return_next()", analysis.Select, []string{"return_next: db.public.t.name"}, "db.public.t"},
		{"SELECT returning_into()", analysis.DML, nil, "db.public.u"},
		{"SELECT diagnostics()", analysis.Select, []string{"diagnostics:"}, ""},
		{"SELECT x.id FROM deleted() AS x(id)", analysis.DML, nil, "db.public.u"},
		// Each function of a recursive cycle reads what any of them reads,
		// whichever is called first.
		{"SELECT f(0), g(1)", analysis.Select, []string{"f: db.public.t.ssn", "g: db.public.t.ssn"}, "db.public.t"},
		{"SELECT g(1), f(0)", analysis.Select, []string{"g: db.public.t.ssn", "f: db.public.t.ssn"}, "db.public.t"},
		{"SELECT cycle_a(), cycle_b(), cycle_c()", analysis.Select, []string{"cycle_a: db.public.t.name", "cycle_b: db.public.t.name", "cycle_c: db.public.t.name"}, "db.public.t"},
		{"SELECT cycle_b(), cycle_c(), cycle_a()", analysis.Select, []string{"cycle_b: db.public.t.name", "cycle_c: db.public.t.name", "cycle_a: db.public.t.name"}, "db.public.t"},
		// A builtin reads only its arguments.
		{"SELECT lower(name) FROM t", analysis.Select, []string{"lower: db.public.t.name"}, "db.public.t"},
	} {
		t.Run(tc.statement, func(t *testing.T) {
			span, err := resolveSpan(tc.statement)
			if err != nil {
				t.Fatal(err)
			}
			if span.Type != tc.queryType {
				t.Errorf("type %s, want %s", span.Type, tc.queryType)
			}
			if diff := cmp.Diff(tc.want, lineage(span)); diff != "" {
				t.Errorf("lineage (-want +got):\n%s", diff)
			}
			if got := columnNames(span.SourceColumns); got != tc.tables {
				t.Errorf("tables %q, want %q", got, tc.tables)
			}
		})
	}

	// What a body reads counts against a mix of user and system tables, and
	// its conditions are the statement's.
	if _, err := resolveSpan("SELECT relnames(), name FROM t"); !errors.Is(err, ErrMixUserSystemTables) {
		t.Errorf("got %v, want ErrMixUserSystemTables", err)
	}
	if span, err := resolveSpan("SELECT note_of(1)"); err != nil || columnNames(span.PredicateColumns) != "db.public.u.t_id" {
		t.Errorf("got predicates %v, error %v", span.PredicateColumns, err)
	}
	// A body that cannot be analyzed fails any statement, as does one with a
	// form CockroachDB's PL/pgSQL grammar does not take.
	for _, statement := range []string{"SELECT dynamic()", "SELECT broken()", "DELETE FROM t WHERE ssn = dynamic()", "SELECT query_execute()", "SELECT query_loop()", "SELECT performs()"} {
		var unsupported *analysis.FunctionNotSupportedError
		if _, err := resolveSpan(statement); !errors.As(err, &unsupported) {
			t.Errorf("%s: got %v, want a FunctionNotSupportedError", statement, err)
		}
	}
}

// TestResolveQuerySpanSearchPath holds unqualified names to CockroachDB's
// search path: pg_catalog first and pg_extension before public, each unless
// the path names it.
func TestResolveQuerySpanSearchPath(t *testing.T) {
	for _, tc := range []struct {
		searchPath, statement string
		queryType             analysis.QueryType
		tables                string
	}{
		{"public", "SELECT * FROM pg_class", analysis.SelectInfoSchema, "db.pg_catalog.pg_class"},
		{"public,pg_catalog", "SELECT * FROM pg_class", analysis.Select, "db.public.pg_class"},
		{"pg_catalog,public", "SELECT * FROM pg_class", analysis.SelectInfoSchema, "db.pg_catalog.pg_class"},
		{"public", "SELECT * FROM geometry_columns", analysis.SelectInfoSchema, "db.pg_extension.geometry_columns"},
		{"public,pg_extension", "SELECT * FROM geometry_columns", analysis.Select, "db.public.geometry_columns"},
		{"hr", "SELECT * FROM geometry_columns", analysis.QueryTypeUnknown, ""},
		// Every virtual schema the path names reaches its relations, after a
		// user table earlier in the path.
		{"public,information_schema", "SELECT table_name FROM tables", analysis.SelectInfoSchema, "db.information_schema.tables"},
		{"public,crdb_internal", "SELECT * FROM jobs", analysis.SelectInfoSchema, "db.crdb_internal.jobs"},
		{"information_schema,public", "SELECT * FROM columns", analysis.SelectInfoSchema, "db.information_schema.columns"},
		{"public,information_schema", "SELECT * FROM pg_class", analysis.SelectInfoSchema, "db.pg_catalog.pg_class"},
		{"hr,information_schema", "SELECT * FROM t", analysis.Select, "db.hr.t"},
	} {
		span, err := resolveSpanIn(tc.statement, tc.searchPath)
		if tc.queryType == analysis.QueryTypeUnknown {
			// With no public in the path, pg_extension is not searched.
			if err != nil || span.NotFoundError == nil {
				t.Errorf("%s under %s: got %v, error %v, want NotFoundError", tc.statement, tc.searchPath, span, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s under %s: %v", tc.statement, tc.searchPath, err)
			continue
		}
		if span.Type != tc.queryType || columnNames(span.SourceColumns) != tc.tables {
			t.Errorf("%s under %s: got %s %q, want %s %q", tc.statement, tc.searchPath, span.Type, columnNames(span.SourceColumns), tc.queryType, tc.tables)
		}
	}
	if _, err := resolveSpanIn("SELECT table_name, name FROM tables, t", "public,information_schema"); !errors.Is(err, ErrMixUserSystemTables) {
		t.Errorf("got %v, want ErrMixUserSystemTables", err)
	}
	span, err := resolveSpanIn("SELECT relname FROM pg_class", "public,pg_catalog")
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff([]string{"relname: db.public.pg_class.relname"}, lineage(span)); diff != "" {
		t.Errorf("lineage (-want +got):\n%s", diff)
	}
}

// TestResolveQuerySpanEmbeddedQueries holds every statement that holds a
// query or an expression to the relations they read.
func TestResolveQuerySpanEmbeddedQueries(t *testing.T) {
	for _, tc := range []struct {
		statement string
		queryType analysis.QueryType
		tables    string
	}{
		{"EXPORT INTO CSV 'nodelocal://1/export' FROM SELECT ssn FROM public.t", analysis.DML, "db.public.t"},
		{"CREATE CHANGEFEED FOR TABLE t INTO 'null://'", analysis.DDL, "db.public.t"},
		{"CREATE CHANGEFEED INTO 'null://' AS SELECT ssn FROM t", analysis.DDL, "db.public.t"},
		{"CREATE SCHEDULE s FOR CHANGEFEED TABLE t INTO 'null://' RECURRING '@daily'", analysis.DDL, "db.public.t"},
		{"ALTER TABLE t SPLIT AT SELECT id FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"ALTER TABLE t UNSPLIT AT SELECT id FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"ALTER TABLE t RELOCATE SELECT ARRAY[1], id FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"ALTER RANGE RELOCATE FROM 1 TO 2 FOR SELECT id FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"CANCEL JOBS SELECT id FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"DECLARE c CURSOR FOR SELECT ssn FROM t", analysis.QueryTypeUnknown, "db.public.t"},
		{"PREPARE p AS DELETE FROM u", analysis.QueryTypeUnknown, "db.public.u"},
		{"EXECUTE p((SELECT max(id) FROM u))", analysis.QueryTypeUnknown, "db.public.u"},
		{"REFRESH MATERIALIZED VIEW region_counts", analysis.DDL, "db.public.region_counts"},
		{"CREATE MATERIALIZED VIEW m AS SELECT ssn FROM t", analysis.DDL, "db.public.t"},
		{"SHOW JOBS SELECT id FROM u", analysis.SelectInfoSchema, "db.public.u"},
		// A definition of code that runs later reads nothing now.
		{"CREATE FUNCTION f() RETURNS STRING LANGUAGE SQL AS $$ SELECT ssn FROM t $$", analysis.DDL, ""},
		{"SET search_path = hr, public", analysis.Select, ""},
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
	for _, statement := range []string{
		"EXPORT INTO CSV 'nodelocal://1/export' FROM SELECT ssn FROM public.t, pg_catalog.pg_class",
		"ALTER TABLE t SPLIT AT SELECT u.id FROM u, pg_class",
	} {
		if _, err := resolveSpan(statement); !errors.Is(err, ErrMixUserSystemTables) {
			t.Errorf("%s: got %v, want ErrMixUserSystemTables", statement, err)
		}
	}
}

// TestResolveQuerySpanFunctionOverloads holds a call to every overload of its
// name the search path reaches, since CockroachDB picks one by its arguments.
func TestResolveQuerySpanFunctionOverloads(t *testing.T) {
	for _, tc := range []struct {
		searchPath, statement string
		want                  []string
	}{
		{"public,hr", "SELECT pick('x')", []string{"pick: db.hr.t.salary db.public.t.ssn"}},
		{"hr,public", "SELECT pick(1)", []string{"pick: db.hr.t.salary db.public.t.ssn"}},
		{"public", "SELECT pick('x')", []string{"pick: db.public.t.ssn"}},
		{"public,hr", "SELECT public.pick(1)", []string{"pick: db.public.t.ssn"}},
	} {
		span, err := resolveSpanIn(tc.statement, tc.searchPath)
		if err != nil {
			t.Errorf("%s under %s: %v", tc.statement, tc.searchPath, err)
			continue
		}
		if diff := cmp.Diff(tc.want, lineage(span)); diff != "" {
			t.Errorf("%s under %s: lineage (-want +got):\n%s", tc.statement, tc.searchPath, diff)
		}
	}
}

// TestResolveQuerySpanFunctionBody holds what a body returns through a
// RETURNING list, which a writing function's DML span does not show.
func TestResolveQuerySpanFunctionBody(t *testing.T) {
	databases := spanTestMetadata()
	for name, want := range map[string]string{
		"returning_into()": "db.public.u.id db.public.u.note db.public.u.t_id",
		"deleted()":        "db.public.u.id db.public.u.note db.public.u.t_id",
		// Every query of a body counts, INSERT ... SELECT's rows too.
		"copied()": "db.public.t.id db.public.t.name db.public.t.ssn",
	} {
		r := newSpanResolver(ResolveOptions{DefaultDatabase: "db", Catalog: Catalog{GetDatabase: func(name string) (Database, error) {
			return protoDatabase{databases[name]}, nil
		}}})
		var definition *metadatapb.FunctionMetadata
		for _, f := range databases["db"].GetSchemas()[0].GetFunctions() {
			if f.GetSignature() == name {
				definition = f
			}
		}
		body := r.analyzeFunction(name, definition)
		if got := columnNames(body.sources); got != want || !body.writes {
			t.Errorf("%s: reads %q, writes %v; want %q and a write", name, got, body.writes, want)
		}
	}
}

// TestResolveQuerySpanFunctionWidths holds a function in FROM to the columns
// it returns: a synced function's declared ones, and a built-in's only when
// the call may reach the built-in.
func TestResolveQuerySpanFunctionWidths(t *testing.T) {
	for _, tc := range []struct {
		statement string
		want      []string
	}{
		// A synced unnest whose body reads only its parameter is still the
		// synced function, two columns wide.
		{"SELECT r.b FROM t, ROWS FROM (public.unnest(t.ssn), pg_catalog.generate_series(1, 2)) AS r(a, b)", []string{"b: db.public.t.ssn"}},
		{"SELECT r.c FROM t, ROWS FROM (public.unnest(t.ssn), pg_catalog.generate_series(1, 2)) AS r(a, b, c)", []string{"c:"}},
		{"SELECT * FROM t AS x, public.unnest(x.ssn) AS u", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "a: db.public.t.ssn", "b: db.public.t.ssn"}},
		{"SELECT * FROM table_rows('x')", []string{"id: db.public.t.id db.public.t.name db.public.t.ssn", "ssn: db.public.t.id db.public.t.name db.public.t.ssn", "name: db.public.t.id db.public.t.name db.public.t.ssn"}},
		{"SELECT * FROM t AS x, pairs(x.ssn)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "k: db.public.t.ssn", "v: db.public.t.ssn"}},
		// A record of unknown columns, or a name the built-in may answer, is
		// of unknown width: every column reads everything.
		{"SELECT r.b FROM t, ROWS FROM (anything(t.ssn), generate_series(1, 2)) AS r(a, b)", []string{"b: db.public.t.ssn"}},
		{"SELECT r.b FROM t, ROWS FROM (unnest(t.ssn), generate_series(1, t.id)) AS r(a, b)", []string{"b: db.public.t.id db.public.t.ssn"}},
		{"SELECT r.c FROM t, ROWS FROM (unnest(t.ssn), generate_series(1, t.id)) AS r(a, b, c)", []string{"c: db.public.t.id db.public.t.ssn"}},
		// pg_catalog's is always the built-in.
		{"SELECT * FROM t AS x, pg_catalog.json_each(x.ssn::JSONB)", []string{"id: db.public.t.id", "ssn: db.public.t.ssn", "name: db.public.t.name", "key: db.public.t.ssn", "value: db.public.t.ssn"}},
		{"SELECT u FROM t, pg_catalog.unnest(ARRAY[t.ssn]) AS u", []string{"u: db.public.t.ssn"}},
		// A synced scalar function's column takes the alias's name; a table
		// function's does not.
		{"SELECT s FROM get_ssn() AS s", []string{"s: db.public.t.ssn"}},
		{"SELECT * FROM one_column() AS u", []string{"a: db.public.t.name"}},
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
