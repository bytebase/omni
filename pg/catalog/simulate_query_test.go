package catalog

import (
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/bytebase/omni/enginetest/simulate"
	"github.com/bytebase/omni/pg/internal/testcontainer"
)

// TestQuerySimulation holds the catalog to testdata/queries/simulate.tsv:
// every generated SELECT statement PostgreSQL accepts in the parser's
// compat.tsv, wrapped as CREATE VIEW _v AS <statement> on one_table@column,
// with what PostgreSQL did and what Exec does. Wrapping a query in a view
// checks both error SQLSTATEs and the engine's stored reading of the query
// (derived columns and pg_get_viewdef text).
//
// PostgreSQL's side is on the test's tape, so the test starts no server.
// OMNI_RECORD=1 runs every row on the pinned PostgreSQL again; -update-compat
// commits the rows a catalog fix moved; and -simulate-show=<text> prints both
// schemas of the rows that hold the text.
func TestQuerySimulation(t *testing.T) {
	rows, err := simulate.FromCompat(filepath.Join("..", "parser", "testdata", "compat.tsv"), querySimulationFixtures, simulatedQuery)
	if err != nil {
		t.Fatal(err)
	}
	// Only the grammar's rows: compat.tsv's hand rows (suite:, issue:) are
	// corpora that grow on their own, and this file should not move with them.
	// FromCompat keeps a statement's first rule, and hand rows follow the
	// generated ones, so a generated statement is never lost here.
	rows = slices.DeleteFunc(rows, func(g simulate.Generated) bool { return handLabel(g.Rule) })
	seen := make(map[string]bool, len(rows))
	wrapped := rows[:0]
	for _, r := range rows {
		stmt := "CREATE VIEW _v AS " + strings.TrimSpace(strings.TrimRight(strings.TrimSpace(r.Statement), ";"))
		key := r.Fixture + "\x00" + stmt
		if seen[key] {
			continue
		}
		seen[key] = true
		r.Statement = stmt
		wrapped = append(wrapped, r)
	}
	simulate.Check(t, simulate.Suite{
		Path:     filepath.Join("testdata", "queries", "simulate.tsv"),
		Pins:     testcontainer.Pins,
		Fixtures: filepath.Join("testdata", "fixtures"),
		Rows:     wrapped,
		MadeUp:   nil,
		Omni:     simulationOmni,
		Engine:   simulationEngine,
	})
}

var querySimulationFixtures = []string{"one_table@column"}

func simulatedQuery(statement string) bool { return strings.HasPrefix(statement, "SELECT ") }

// handLabel reports whether a compat.tsv rule is a label a person wrote
// (issue:BYT-1, suite:cited) rather than a grammar rule: a word and a colon
// before any slash.
func handLabel(rule string) bool {
	i := strings.IndexAny(rule, ":/")
	return i >= 0 && rule[i] == ':'
}

func TestQuerySimulationFixes(t *testing.T) {
	setup := func(t *testing.T) *Catalog {
		t.Helper()
		c, err := LoadSQL(`
			CREATE TABLE a (b int PRIMARY KEY, a int, d int);
			CREATE TABLE b (b int PRIMARY KEY, c int);
			CREATE SCHEMA s1; CREATE TABLE s1.t (id int);
			CREATE SCHEMA s2; CREATE TABLE s2.t (id int);
		`)
		if err != nil {
			t.Fatal(err)
		}
		return c
	}

	t.Run("errors", func(t *testing.T) {
		cases := []struct {
			name string
			sql  string
			code string
		}{
			{"select_star_no_from", "CREATE VIEW _v AS SELECT *", CodeSyntaxError},
			{"for_update_of_missing", "CREATE VIEW _v AS SELECT FOR UPDATE OF a", CodeUndefinedTable},
			{"for_update_of_hidden_by_join_alias", "CREATE VIEW _v AS SELECT FROM (a CROSS JOIN b) AS j FOR UPDATE OF a", CodeUndefinedTable},
			{"dup_from_alias", "CREATE VIEW _v AS SELECT FROM a, a", CodeDuplicateAlias},
			{"dup_join_alias", "CREATE VIEW _v AS SELECT FROM a CROSS JOIN a", CodeDuplicateAlias},
			{"dup_join_on_before_ambig_col", "CREATE VIEW _v AS SELECT FROM a FULL JOIN a ON a", CodeDuplicateAlias},
			{"dup_join_using_alias", "CREATE VIEW _v AS SELECT FROM a JOIN b USING (b) AS a", CodeDuplicateAlias},
			{"func_coldeflist_non_record", "CREATE VIEW _v AS SELECT FROM CURRENT_DATE AS (x int)", CodeSyntaxError},
			{"select_default", "CREATE VIEW _v AS SELECT DEFAULT", CodeSyntaxError},
			{"view_param_ref", "CREATE VIEW _v AS SELECT $1", CodeUndefinedParameter},
			{"empty_array_untyped", "CREATE VIEW _v AS SELECT ARRAY[]", CodeIndeterminateDatatype},
			{"nested_empty_array_untyped", "CREATE VIEW _v AS SELECT ARRAY[[]]", CodeIndeterminateDatatype},
			{"row_pseudotype_view_col", "CREATE VIEW _v AS SELECT ROW()", CodeInvalidTableDefinition},
			{"tablesample_missing_method", "CREATE VIEW _v AS SELECT FROM a TABLESAMPLE nonexistent(10)", CodeUndefinedObject},
			{"window_duplicate_name", "CREATE VIEW _v AS SELECT WINDOW w AS (), w AS ()", CodeWindowingError},
			{"window_groups_needs_order_by", "CREATE VIEW _v AS SELECT WINDOW w AS (GROUPS UNBOUNDED PRECEDING)", CodeWindowingError},
			{"window_range_offset_needs_order_by", "CREATE VIEW _v AS SELECT WINDOW w AS (RANGE a PRECEDING)", CodeWindowingError},
			{"sublink_in_testexpr_first", "CREATE VIEW _v AS SELECT nonexistent IN (SELECT)", CodeUndefinedColumn},
			{"sublink_array_zero_cols", "CREATE VIEW _v AS SELECT ARRAY(SELECT)", CodeSyntaxError},
			{"json_array_subquery_zero_cols", "CREATE VIEW _v AS SELECT JSON_ARRAY(SELECT)", CodeSyntaxError},
			{"json_array_subquery_in_default", "CREATE TABLE t_def_jsq (j json DEFAULT JSON_ARRAY(SELECT 1))", CodeFeatureNotSupported},
			{"json_arrayagg_in_default", "CREATE TABLE t_def_jaa (j json DEFAULT JSON_ARRAYAGG(1))", CodeGroupingError},
			{"json_arrayagg_over_in_default", "CREATE TABLE t_def_jaaw (j json DEFAULT JSON_ARRAYAGG(1) OVER ())", CodeWindowingError},
			{"json_objectagg_in_default", "CREATE TABLE t_def_joa (j json DEFAULT JSON_OBJECTAGG('k' : 1))", CodeGroupingError},
			{"json_arrayagg_missing_col", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(nonexistent)", CodeUndefinedColumn},
			{"json_arrayagg_nonbool_filter", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(b) FILTER (WHERE 1) AS j FROM a", CodeDatatypeMismatch},
			{"json_objectagg_nonbool_filter", "CREATE VIEW _v AS SELECT JSON_OBJECTAGG('k' : b) FILTER (WHERE 1) AS j FROM a", CodeDatatypeMismatch},
			{"json_objectagg_missing_col", "CREATE VIEW _v AS SELECT JSON_OBJECTAGG(nonexistent VALUE 1)", CodeUndefinedColumn},
			{"json_object_bad_returning", "CREATE VIEW _v AS SELECT JSON_OBJECT(RETURNING a)", CodeCannotCoerce},
			{"json_query_bool_behavior", "CREATE VIEW _v AS SELECT JSON_QUERY(a, a TRUE ON EMPTY)", CodeSyntaxError},
			{"xmltable_missing_col", "CREATE VIEW _v AS SELECT FROM XMLTABLE(nonexistent PASSING nonexistent COLUMNS x int)", CodeUndefinedColumn},
			{"merge_action_outside_merge", "CREATE VIEW _v AS SELECT MERGE_ACTION()", CodeSyntaxError},
			{"typmod_disallowed_on_composite", "CREATE VIEW _v AS SELECT a(a) 's'", CodeSyntaxError},
			{"typmod_non_integer_on_bit", "CREATE VIEW _v AS SELECT BIT(a) 's'", CodeInvalidTextRepresentation},
			{"typed_literal_bad_int", "CREATE VIEW _v AS SELECT INT 's'", CodeInvalidTextRepresentation},
			{"typed_literal_bad_bit", "CREATE VIEW _v AS SELECT BIT 's'", CodeInvalidTextRepresentation},
			{"typed_literal_bad_json", "CREATE VIEW _v AS SELECT JSON 's'", CodeInvalidTextRepresentation},
			{"typed_literal_bad_composite", "CREATE VIEW _v AS SELECT a 's'", CodeInvalidTextRepresentation},
			{"typed_literal_bad_timestamp", "CREATE VIEW _v AS SELECT TIMESTAMP 's'", CodeInvalidDatetimeFormat},
			{"typed_literal_bad_interval", "CREATE VIEW _v AS SELECT INTERVAL 's'", CodeInvalidDatetimeFormat},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				c := setup(t)
				results, err := c.Exec(tc.sql, nil)
				if err == nil {
					for _, r := range results {
						if r.Error != nil {
							err = r.Error
							break
						}
					}
				}
				if err == nil {
					t.Fatalf("expected error %s, got nil", tc.code)
				}
				assertErrorCode(t, err, tc.code)
			})
		}
	})

	t.Run("views", func(t *testing.T) {
		cases := []struct {
			name    string
			sql     string
			wantSub string
		}{
			{"for_update_empty_from", "CREATE VIEW _v AS SELECT FOR UPDATE NOWAIT", "SELECT;"},
			{"unnamed_subquery", "CREATE VIEW _v AS SELECT FROM (SELECT)", "unnamed_subquery"},
			{"from_current_date", "CREATE VIEW _v AS SELECT FROM CURRENT_DATE", `"current_date"("current_date")`},
			{"rows_from_dup_col", "CREATE VIEW _v AS SELECT FROM ROWS FROM(CURRENT_DATE, CURRENT_DATE)", `("current_date", current_date_1)`},
			{"table_col_alias_dedup", "CREATE VIEW _v AS SELECT FROM a AS a(a)", "a(a, a_1, d)"},
			{"exists_col_name", "CREATE VIEW _v AS SELECT EXISTS(SELECT)", `AS "exists"`},
			{"user_deparse_and_name", "CREATE VIEW _v AS SELECT USER", `USER AS "user"`},
			{"bit_literal", "CREATE VIEW _v AS SELECT B'1'", `'1'::"bit"`},
			{"hex_bit_literal", "CREATE VIEW _v AS SELECT X'1F'", `'00011111'::"bit"`},
			{"typed_empty_array", "CREATE VIEW _v AS SELECT ARRAY[]::int[]", "ARRAY[]::integer[]"},
			{"nested_typed_empty_array", "CREATE VIEW _v AS SELECT ARRAY[[]]::int[]", "ARRAY[ARRAY[]::integer[]]"},
			{"valid_bit_hex_literal", "CREATE VIEW _v AS SELECT BIT(8) 'x1F'", `'x1F'::bit(8)`},
			{"array_sublink_deparse", "CREATE VIEW _v AS SELECT ARRAY(SELECT 1)", `ARRAY( SELECT 1) AS "array"`},
			{"array_type_with_typmod", "CREATE VIEW _v AS SELECT NULL::varchar(10)[]", "character varying[]"},
			{"named_window_refs", "CREATE VIEW _v AS SELECT count(*) OVER w AS c, sum(a) OVER w AS s FROM a WINDOW w AS (ORDER BY b)", "OVER w"},
			{"dedup_col_alias_var_refs", "CREATE VIEW _v(first, second, third) AS SELECT x.* FROM a AS x(a)", "a_1 AS second"},
			{"json_array_query_constructor", "CREATE VIEW _v AS SELECT JSON_ARRAY(SELECT b FROM a) AS j", "JSON_ARRAY( SELECT"},
			{"json_array_query_format_json", `CREATE VIEW _v AS SELECT JSON_ARRAY(SELECT '{"x":1}'::text AS j FORMAT JSON) AS a`, "FORMAT JSON"},
			{"json_arrayagg_preserve", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(b ORDER BY b ABSENT ON NULL RETURNING jsonb) FILTER (WHERE b > 0) AS j FROM a", "JSON_ARRAYAGG(b ORDER BY b ABSENT ON NULL RETURNING jsonb) FILTER (WHERE (b > 0))"},
			{"json_arrayagg_null_on_null", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(b NULL ON NULL) AS j FROM a", "JSON_ARRAYAGG(b NULL ON NULL RETURNING json)"},
			{"json_array_null_on_null", "CREATE VIEW _v AS SELECT JSON_ARRAY(b NULL ON NULL) AS j FROM a", "JSON_ARRAY(b NULL ON NULL RETURNING json)"},
			{"json_arrayagg_format_json", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(b::text FORMAT JSON) AS j FROM a", "FORMAT JSON"},
			{"json_objectagg_format_json", "CREATE VIEW _v AS SELECT JSON_OBJECTAGG('k' : b::text FORMAT JSON) AS j FROM a", "FORMAT JSON"},
			{"json_objectagg_over_preserve", "CREATE VIEW _v AS SELECT JSON_OBJECTAGG('k' : b WITH UNIQUE) OVER (PARTITION BY d) AS j FROM a", "JSON_OBJECTAGG('k' : b WITH UNIQUE RETURNING json) OVER (PARTITION BY d)"},
			{"xmltable_preserve", "CREATE VIEW _v AS SELECT x.id, x.ord FROM XMLTABLE(XMLNAMESPACES(DEFAULT 'http://ex'), '/r' PASSING '<r/>'::xml COLUMNS id int PATH '@id' NOT NULL, ord FOR ORDINALITY) x", "XMLTABLE(XMLNAMESPACES(DEFAULT 'http://ex'), '/r' PASSING '<r/>'::xml COLUMNS id integer PATH '@id' NOT NULL, ord FOR ORDINALITY) x"},
			{"cross_schema_unaliased_same_name", "CREATE VIEW _v AS SELECT 1 AS x FROM s1.t, s2.t", "s2.t"},
			{"inherited_window_order_for_groups", "CREATE VIEW _v AS SELECT count(*) OVER w2 AS n FROM a WINDOW w AS (ORDER BY b), w2 AS (w GROUPS UNBOUNDED PRECEDING)", "OVER w2"},
			{"multiple_unnamed_subqueries", "CREATE VIEW _v AS SELECT 1 AS x FROM (SELECT 1), (SELECT 2)", "unnamed_subquery_1"},
			{"aliased_join_deparse", "CREATE VIEW _v AS SELECT b FROM (a FULL JOIN b USING (b)) j", ") j"},
			{"aliased_cross_join_deparse", "CREATE VIEW _v AS SELECT d, c FROM (a CROSS JOIN b) j", "j.d"},
			{"aliased_join_dedup_col_refs", "CREATE VIEW _v(first, second, third, fourth, fifth) AS SELECT j.* FROM (a CROSS JOIN b) AS j(c)", "j.c_1 AS fifth"},
			{"json_returning_bytea", "CREATE VIEW _v AS SELECT JSON_ARRAY(1 RETURNING bytea) AS j, JSON_OBJECT('k' : 1 RETURNING bytea) AS o, JSON_ARRAYAGG(b RETURNING bytea) AS ja FROM a", "RETURNING bytea"},
			{"json_returning_typmod", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(b RETURNING varchar(5)) AS ja, JSON_ARRAY(1 RETURNING varchar(10)) AS j FROM a", "RETURNING character varying(5)"},
			{"unnamed_subquery_no_shadow_outer", "CREATE VIEW _v AS SELECT (SELECT unnamed_subquery.b FROM (SELECT 2 AS b)) AS z FROM a AS unnamed_subquery", "unnamed_subquery_1"},
			{"pg17_integer_literals", "CREATE VIEW _v AS SELECT '0xFF'::integer AS h, '1_000'::bigint AS u, '0o77'::smallint AS o, '0b1010'::int AS b", "'0xFF'::integer"},
			{"pg17_numeric_literals", "CREATE VIEW _v AS SELECT '1_000.2_5e1_0'::numeric AS u, '0xFF'::numeric AS h, '0o77'::numeric AS o, '0b1010'::numeric AS b", "'0xFF'::numeric"},
			{"signed_oid_literal", "CREATE VIEW _v AS SELECT '-1'::oid AS x, '-2147483648'::oid AS y", "'-1'::oid"},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				c := setup(t)
				results, err := c.Exec(tc.sql, nil)
				if err != nil {
					t.Fatalf("exec error: %v", err)
				}
				for _, r := range results {
					if r.Error != nil {
						t.Fatalf("stmt error: %v", r.Error)
					}
				}
				def, err := c.GetViewDefinition("public", "_v")
				if err != nil {
					t.Fatalf("GetViewDefinition: %v", err)
				}
				if !strings.Contains(def, tc.wantSub) {
					t.Errorf("view def %q does not contain %q", def, tc.wantSub)
				}
				// Verify the deparsed view definition can be recreated cleanly.
				switch tc.name {
				case "typed_empty_array", "dedup_col_alias_var_refs",
					"json_array_query_constructor", "json_array_query_format_json",
					"json_arrayagg_preserve", "json_arrayagg_null_on_null",
					"json_array_null_on_null", "json_arrayagg_format_json",
					"json_objectagg_format_json", "json_objectagg_over_preserve",
					"xmltable_preserve", "cross_schema_unaliased_same_name",
					"multiple_unnamed_subqueries", "aliased_join_deparse",
					"aliased_cross_join_deparse", "aliased_join_dedup_col_refs",
					"json_returning_bytea", "json_returning_typmod",
					"unnamed_subquery_no_shadow_outer", "pg17_integer_literals",
					"pg17_numeric_literals", "signed_oid_literal":
					c2 := setup(t)
					res2, err := c2.Exec("CREATE VIEW _v2 AS "+def, nil)
					if err != nil {
						t.Fatalf("recreate exec error for %q: %v", def, err)
					}
					for _, r := range res2 {
						if r.Error != nil {
							t.Fatalf("recreate stmt error for %q: %v", def, r.Error)
						}
					}
				}
			})
		}
	})

	t.Run("quoted_merge_action_func", func(t *testing.T) {
		c := setup(t)
		res, err := c.Exec(`
			CREATE FUNCTION public."merge_action"() RETURNS text LANGUAGE sql AS 'SELECT ''ok''::text';
			CREATE VIEW _v AS SELECT "merge_action"() AS x;
		`, nil)
		if err != nil {
			t.Fatalf("exec error: %v", err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("stmt error: %v", r.Error)
			}
		}
		def, err := c.GetViewDefinition("public", "_v")
		if err != nil {
			t.Fatalf("GetViewDefinition: %v", err)
		}
		if !strings.Contains(def, "merge_action()") {
			t.Errorf("unexpected view def %q", def)
		}
	})

	t.Run("array_subquery_typmod_and_collation", func(t *testing.T) {
		c := setup(t)
		res, err := c.Exec(`CREATE VIEW _v AS SELECT ARRAY(SELECT 'x'::varchar(5) COLLATE "C") AS arr`, nil)
		if err != nil {
			t.Fatalf("exec error: %v", err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("stmt error: %v", r.Error)
			}
		}
		rel := c.GetRelation("public", "_v")
		if rel == nil || len(rel.Columns) != 1 {
			t.Fatalf("expected 1 column on view _v")
		}
		col := rel.Columns[0]
		if col.TypeMod != 9 { // varchar(5) typmod is 5 + VARHDRSZ(4) = 9
			t.Errorf("got typmod %d, want 9", col.TypeMod)
		}
		if col.Collation != C_COLLATION_OID {
			t.Errorf("got collation %d, want %d", col.Collation, C_COLLATION_OID)
		}
	})

	t.Run("quoted_tablesample_handler", func(t *testing.T) {
		c := setup(t)
		res, err := c.Exec(`
			CREATE FUNCTION public."MySample"(internal) RETURNS tsm_handler LANGUAGE internal AS 'tsm_bernoulli_handler';
			CREATE VIEW _v AS SELECT b FROM a TABLESAMPLE "MySample"(10);
		`, nil)
		if err != nil {
			t.Fatalf("exec error: %v", err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("stmt error: %v", r.Error)
			}
		}
	})

	t.Run("xmltable_column_type_dep", func(t *testing.T) {
		c := setup(t)
		res, err := c.Exec(`
			CREATE TYPE mood AS ENUM ('happy');
			CREATE VIEW _v AS SELECT 1 AS n FROM XMLTABLE('/r' PASSING '<r/>'::xml COLUMNS x mood) t;
		`, nil)
		if err != nil {
			t.Fatalf("setup error: %v", err)
		}
		for _, r := range res {
			if r.Error != nil {
				t.Fatalf("setup stmt error: %v", r.Error)
			}
		}
		res, err = c.Exec("DROP TYPE mood", nil)
		if err != nil {
			t.Fatalf("exec error: %v", err)
		}
		assertErrorCode(t, res[0].Error, CodeDependentObjects)

		res, err = c.Exec("DROP TYPE mood CASCADE", nil)
		if err != nil {
			t.Fatalf("cascade exec error: %v", err)
		}
		if res[0].Error != nil {
			t.Fatalf("cascade stmt error: %v", res[0].Error)
		}
		if _, err := c.GetViewDefinition("public", "_v"); err == nil {
			t.Fatalf("expected view _v to be dropped by DROP TYPE mood CASCADE")
		}
	})
}
