package analysis

import (
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractChangedResources(t *testing.T) {
	type table struct {
		database string
		schema   string
		name     string
		affected bool
	}
	for _, tc := range []struct {
		name          string
		statement     string
		tables        []table
		dmlStatements []string
		dmlCount      int
		insertCount   int
	}{
		{
			name: "ddl_affects_a_table_that_dml_also_changes",
			statement: `CREATE TABLE t1 (c1 INT);
ALTER TABLE t1 ADD COLUMN c2 INT;
INSERT INTO t1 VALUES (1);
UPDATE t1 SET c2 = 2;`,
			tables:        []table{{database: "db", schema: "public", name: "t1", affected: true}},
			dmlStatements: []string{"UPDATE t1 SET c2 = 2;"},
			dmlCount:      1,
			insertCount:   1,
		},
		{
			name: "dml_targets_are_not_affected",
			statement: `UPDATE t SET c = 1 WHERE id = 1;
DELETE FROM t WHERE id = 2;`,
			tables: []table{{database: "db", schema: "public", name: "t"}},
			dmlStatements: []string{
				"UPDATE t SET c = 1 WHERE id = 1;",
				"DELETE FROM t WHERE id = 2;",
			},
			dmlCount: 2,
		},
		{
			name: "insert_values_counts_rows_without_a_sample",
			statement: `INSERT INTO t VALUES (1), (2), (3);
INSERT INTO t (c) VALUES (4);
INSERT INTO t DEFAULT VALUES;`,
			tables:      []table{{database: "db", schema: "public", name: "t"}},
			insertCount: 5,
		},
		{
			name:          "insert_select_is_sampled",
			statement:     `INSERT INTO t SELECT * FROM s WHERE c > 1;`,
			tables:        []table{{database: "db", schema: "public", name: "t"}},
			dmlStatements: []string{"INSERT INTO t SELECT * FROM s WHERE c > 1;"},
			dmlCount:      1,
		},
		{
			name: "merge_is_counted_without_a_sample",
			statement: `MERGE INTO m USING s ON m.id = s.id WHEN MATCHED THEN UPDATE SET c = s.c;
UPDATE t SET c = 1 WHERE id = 1;`,
			tables:        []table{{database: "db", schema: "public", name: "m"}, {database: "db", schema: "public", name: "t"}},
			dmlStatements: []string{"UPDATE t SET c = 1 WHERE id = 1;"},
			dmlCount:      2,
		},
		{
			name: "truncate_affects_tables_resolved_against_the_search_path",
			statement: `TRUNCATE a;
SET search_path TO analytics;
TRUNCATE TABLE b, public.c, other.s.d;
RESET search_path;
TRUNCATE e;`,
			tables: []table{
				{database: "db", schema: "public", name: "a", affected: true},
				{database: "db", schema: "analytics", name: "b", affected: true},
				{database: "db", schema: "public", name: "c", affected: true},
				{database: "other", schema: "s", name: "d", affected: true},
				{database: "db", schema: "public", name: "e", affected: true},
			},
		},
		{
			name: "search_path_drop_and_select_into",
			statement: `SET search_path TO analytics;
CREATE TABLE unqualified(id INT);
SELECT * INTO copied_rows FROM public.rows;
DROP TABLE old_rows;
DELETE FROM copied_rows WHERE id = 1;`,
			tables: []table{
				{database: "db", schema: "analytics", name: "unqualified"},
				{database: "db", schema: "analytics", name: "copied_rows"},
				{database: "db", schema: "analytics", name: "old_rows", affected: true},
			},
			dmlStatements: []string{"SET LOCAL search_path TO \"analytics\";\nDELETE FROM copied_rows WHERE id = 1;"},
			dmlCount:      1,
		},
		{
			name: "sampled_dml_replays_the_search_path",
			statement: `SET search_path TO '$user', analytics, public;
UPDATE t SET c = 1;
RESET search_path;
UPDATE t SET c = 2;`,
			tables: []table{
				{database: "db", schema: "analytics", name: "t"},
				{database: "db", schema: "public", name: "t"},
			},
			dmlStatements: []string{
				"SET LOCAL search_path TO \"analytics\", \"public\";\nUPDATE t SET c = 1;",
				"UPDATE t SET c = 2;",
			},
			dmlCount: 2,
		},
		{
			name: "comment_before_semicolon_between_ddl_and_dml",
			statement: `ALTER TABLE a ADD COLUMN c INT /* note */;
INSERT INTO b VALUES (1);
DROP TABLE c;`,
			tables: []table{
				{database: "db", schema: "public", name: "a", affected: true},
				{database: "db", schema: "public", name: "b"},
				{database: "db", schema: "public", name: "c", affected: true},
			},
			insertCount: 1,
		},
		{
			name:          "every_dml_statement_is_listed",
			statement:     strings.Repeat("UPDATE t SET c = 1;\n", 12),
			tables:        []table{{database: "db", schema: "public", name: "t"}},
			dmlStatements: slices.Repeat([]string{"UPDATE t SET c = 1;"}, 12),
			dmlCount:      12,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractChangedResources(tc.statement, "db", "public")
			require.NoError(t, err)

			// The caller deduplicates the resources it is handed, keeping a flag
			// any occurrence set.
			affected := make(map[table]bool)
			var order []table
			for _, resource := range got.Resources {
				key := table{database: resource.Database, schema: resource.Schema, name: resource.Table}
				if _, ok := affected[key]; !ok {
					order = append(order, key)
				}
				affected[key] = affected[key] || resource.AffectsData
			}
			var gotTables []table
			for _, key := range order {
				key.affected = affected[key]
				gotTables = append(gotTables, key)
			}
			require.ElementsMatch(t, tc.tables, gotTables)
			require.Equal(t, tc.dmlStatements, got.DMLStatements)
			require.Equal(t, tc.dmlCount, got.DMLCount)
			require.Equal(t, tc.insertCount, got.InsertCount)
		})
	}
}
