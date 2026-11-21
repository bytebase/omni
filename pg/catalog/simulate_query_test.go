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
		c, err := LoadSQL("CREATE TABLE a (b int PRIMARY KEY, a int, d int); CREATE TABLE b (b int PRIMARY KEY, c int);")
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
			{"json_arrayagg_missing_col", "CREATE VIEW _v AS SELECT JSON_ARRAYAGG(nonexistent)", CodeUndefinedColumn},
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
			{"bit_literal", "CREATE VIEW _v AS SELECT B'1'", "'1'::bit"},
			{"hex_bit_literal", "CREATE VIEW _v AS SELECT X'1F'", "'00011111'::bit"},
			{"typed_empty_array", "CREATE VIEW _v AS SELECT ARRAY[]::int[]", "ARRAY[]"},
			{"valid_bit_hex_literal", "CREATE VIEW _v AS SELECT BIT 'x1F'", "'x1F'::bit"},
			{"array_sublink_deparse", "CREATE VIEW _v AS SELECT ARRAY(SELECT 1)", `ARRAY( SELECT 1) AS "array"`},
			{"array_type_with_typmod", "CREATE VIEW _v AS SELECT NULL::varchar(10)[]", "character varying[]"},
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
			})
		}
	})
}
