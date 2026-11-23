package analysis

import (
	"reflect"
	"testing"
)

func TestExtractViewDependencies(t *testing.T) {
	cases := []struct {
		name           string
		sql            string
		defaultCatalog string
		defaultSchema  string
		want           []ViewDependency
		wantErr        bool
	}{
		{
			name:           "bare select view definition",
			sql:            "SELECT a.id, b.name FROM orders a JOIN main.dim.customers b ON a.cid = b.id",
			defaultCatalog: "workspace",
			defaultSchema:  "sales",
			want: []ViewDependency{
				{Catalog: "workspace", Schema: "sales", Table: "orders"},
				{Catalog: "main", Schema: "dim", Table: "customers"},
			},
		},
		{
			name:           "create view with schema compensation",
			sql:            "CREATE VIEW omni_record.recent (\n  id,\n  lookup_id)\nWITH SCHEMA COMPENSATION\nAS SELECT id, lookup_id FROM workspace.omni_record.facts WHERE id > 0",
			defaultCatalog: "workspace",
			defaultSchema:  "omni_record",
			want: []ViewDependency{
				{Catalog: "workspace", Schema: "omni_record", Table: "facts"},
			},
		},
		{
			name:           "create materialized view",
			sql:            "CREATE MATERIALIZED VIEW `workspace`.`omni_record`.`totals` (\n  lookup_id INT,\n  total BIGINT)\nAS SELECT lookup_id, SUM(id) AS total FROM workspace.omni_record.facts GROUP BY lookup_id",
			defaultCatalog: "workspace",
			defaultSchema:  "omni_record",
			want: []ViewDependency{
				{Catalog: "workspace", Schema: "omni_record", Table: "facts"},
			},
		},
		{
			name:           "alter view as query",
			sql:            "ALTER VIEW analytics.v_daily AS SELECT * FROM raw.events JOIN dim_users ON raw.events.uid = dim_users.id",
			defaultCatalog: "main",
			defaultSchema:  "analytics",
			want: []ViewDependency{
				{Catalog: "main", Schema: "raw", Table: "events"},
				{Catalog: "main", Schema: "analytics", Table: "dim_users"},
			},
		},
		{
			name:           "ctes and recursive ctes are excluded while external tables inside them are kept",
			sql:            "WITH RECURSIVE r(n) AS (SELECT id FROM base_t UNION ALL SELECT n + 1 FROM r WHERE n < 5), c AS (SELECT * FROM r JOIN ext_sch.ext_t ON r.n = ext_t.id) SELECT * FROM c",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "sch", Table: "base_t"},
				{Catalog: "cat", Schema: "ext_sch", Table: "ext_t"},
			},
		},
		{
			name:           "non-recursive cte referencing table with same name extracts external table",
			sql:            "WITH t AS (SELECT * FROM t) SELECT * FROM t",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "sch", Table: "t"},
			},
		},
		{
			name:           "qualified name matching cte name is still an external dependency",
			sql:            "WITH c AS (SELECT 1 AS id) SELECT * FROM c JOIN other_schema.c ON c.id = other_schema.c.id",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "other_schema", Table: "c"},
			},
		},
		{
			name:           "subqueries in where exists and union and table shorthand",
			sql:            "SELECT * FROM a WHERE EXISTS (SELECT 1 FROM b WHERE b.id = a.id) UNION ALL TABLE sch2.c",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "sch", Table: "a"},
				{Catalog: "cat", Schema: "sch", Table: "b"},
				{Catalog: "cat", Schema: "sch2", Table: "c"},
			},
		},
		{
			name:           "from-first pipe query and backtick identifiers",
			sql:            "FROM `my cat`.`my sch`.`my``tbl` |> JOIN ref_tbl ON ref_tbl.id = 1 |> SELECT *",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "my cat", Schema: "my sch", Table: "my`tbl"},
				{Catalog: "cat", Schema: "sch", Table: "ref_tbl"},
			},
		},
		{
			name:           "case-insensitive deduplication",
			sql:            "SELECT * FROM Orders JOIN orders ON Orders.id = orders.id JOIN cat.sch.ORDERS ON Orders.id = ORDERS.id",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "sch", Table: "Orders"},
			},
		},
		{
			name:           "dynamic identifier skipped",
			sql:            "SELECT * FROM IDENTIFIER('main.sch.dyn') JOIN real_t ON true",
			defaultCatalog: "cat",
			defaultSchema:  "sch",
			want: []ViewDependency{
				{Catalog: "cat", Schema: "sch", Table: "real_t"},
			},
		},
		{
			name:    "invalid syntax returns error",
			sql:     "SELECT 1 +",
			wantErr: true,
		},
		{
			name:    "unclosed comment returns error",
			sql:     "/* unclosed SELECT * FROM t",
			wantErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ExtractViewDependencies(tc.sql, tc.defaultCatalog, tc.defaultSchema)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("ExtractViewDependencies(%q) = %+v, nil; want error", tc.sql, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("ExtractViewDependencies(%q) unexpected error: %v", tc.sql, err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("ExtractViewDependencies(%q) =\n  %+v\nwant:\n  %+v", tc.sql, got, tc.want)
			}
		})
	}
}

func TestQualifyViewDDL(t *testing.T) {
	const target = "`workspace`.`omni_record`.`recent`"
	for _, tc := range []struct {
		name string
		ddl  string
		want string
	}{
		{
			name: "schema-qualified view gains catalog",
			ddl:  "CREATE VIEW omni_record.recent (\n  id INT)\nAS SELECT 1",
			want: "CREATE VIEW " + target + " (\n  id INT)\nAS SELECT 1",
		},
		{
			name: "create or replace view if not exists style",
			ddl:  "/* header /* nested */ */\nCREATE OR REPLACE VIEW sch.v AS SELECT 1",
			want: "/* header /* nested */ */\nCREATE OR REPLACE VIEW " + target + " AS SELECT 1",
		},
		{
			name: "create view if not exists",
			ddl:  "CREATE VIEW IF NOT EXISTS sch.v AS SELECT 1",
			want: "CREATE VIEW IF NOT EXISTS " + target + " AS SELECT 1",
		},
		{
			name: "backtick identifier containing AS and parens and doubled backticks",
			ddl:  "CREATE VIEW sch.`sales ``AS`` archived (old)` (\n  id INT)\nAS SELECT 1",
			want: "CREATE VIEW " + target + " (\n  id INT)\nAS SELECT 1",
		},
		{
			name: "create materialized view is untouched",
			ddl:  "CREATE MATERIALIZED VIEW `workspace`.`omni_record`.`totals` AS SELECT 1",
			want: "CREATE MATERIALIZED VIEW `workspace`.`omni_record`.`totals` AS SELECT 1",
		},
		{
			name: "create table is untouched",
			ddl:  "CREATE TABLE sch.t (id INT)",
			want: "CREATE TABLE sch.t (id INT)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := QualifyViewDDL(tc.ddl, target); got != tc.want {
				t.Errorf("QualifyViewDDL =\n  %q\nwant:\n  %q", got, tc.want)
			}
		})
	}
}
