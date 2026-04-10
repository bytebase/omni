package catalog

// container_viewdef_parse_test.go — Oracle-driven parser gap discovery.
//
// Strategy: Create views in a real PostgreSQL container using every syntax
// feature we can think of, then retrieve the canonical view definition via
// pg_get_viewdef() and attempt to parse it with the omni parser. Any parse
// failure reveals a gap that would cause the Bytebase "failed to load schema
// DDL into catalog" error in production.
//
// Run: go test -v -run TestViewDefParse -tags containers -count=1 ./pg/catalog/

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"

	pg "github.com/bytebase/omni/pg"
)

// viewDefCase describes a PG feature to exercise. SetupSQL creates the
// prerequisite objects (tables, types, etc.) and the view. The test then
// fetches pg_get_viewdef() for the view and tries to parse it.
type viewDefCase struct {
	Name     string // short descriptive name
	SetupSQL string // DDL to run in the test schema (must CREATE VIEW named "v_<name>")
	ViewName string // view name to query (if empty, defaults to "v_" + lowercase name)
}

// viewDefCases is the full catalog of features to probe. Add new entries here
// whenever a new parser gap is suspected.
var viewDefCases = []viewDefCase{
	// ===================== XML =====================
	{
		Name: "xmltable",
		SetupSQL: `
			CREATE VIEW v_xmltable AS
			SELECT x.* FROM XMLTABLE('/root/row' PASSING '<root><row><a>1</a></row></root>'::xml
				COLUMNS a text PATH 'a') AS x;`,
	},
	{
		Name: "xmlexists",
		SetupSQL: `
			CREATE TABLE t_xmlexists (id serial, doc xml);
			CREATE VIEW v_xmlexists AS
			SELECT id FROM t_xmlexists WHERE XMLEXISTS('//a' PASSING BY VALUE doc);`,
	},
	{
		Name: "xmlparse",
		SetupSQL: `
			CREATE VIEW v_xmlparse AS SELECT XMLPARSE(DOCUMENT '<root/>');`,
	},
	{
		Name: "xmlserialize",
		SetupSQL: `
			CREATE VIEW v_xmlserialize AS SELECT XMLSERIALIZE(DOCUMENT '<r/>'::xml AS text);`,
	},
	{
		Name: "xmlelement",
		SetupSQL: `
			CREATE TABLE t_xmlelement (name text);
			CREATE VIEW v_xmlelement AS SELECT XMLELEMENT(NAME item, name) FROM t_xmlelement;`,
	},
	{
		Name: "xmlforest",
		SetupSQL: `
			CREATE TABLE t_xmlforest (a text, b int);
			CREATE VIEW v_xmlforest AS SELECT XMLFOREST(a, b) FROM t_xmlforest;`,
	},
	{
		Name: "xmlagg",
		SetupSQL: `
			CREATE TABLE t_xmlagg (doc xml);
			CREATE VIEW v_xmlagg AS SELECT XMLAGG(doc) FROM t_xmlagg;`,
	},
	{
		Name: "xmlconcat",
		SetupSQL: `
			CREATE VIEW v_xmlconcat AS SELECT XMLCONCAT('<a/>'::xml, '<b/>'::xml);`,
	},
	{
		Name: "xmlroot",
		SetupSQL: `
			CREATE VIEW v_xmlroot AS SELECT XMLROOT('<r/>'::xml, VERSION '1.0');`,
	},
	{
		Name: "xpath",
		SetupSQL: `
			CREATE TABLE t_xpath (doc xml);
			CREATE VIEW v_xpath AS SELECT XPATH('/a', doc) FROM t_xpath;`,
	},
	{
		Name: "xpath_exists",
		SetupSQL: `
			CREATE TABLE t_xpath_exists (doc xml);
			CREATE VIEW v_xpath_exists AS SELECT XPATH_EXISTS('/a', doc) FROM t_xpath_exists;`,
	},

	// ===================== JSON / JSONB =====================
	{
		Name: "jsonb_path_query",
		SetupSQL: `
			CREATE TABLE t_jpq (data jsonb);
			CREATE VIEW v_jsonb_path_query AS SELECT jsonb_path_query(data, '$.a') FROM t_jpq;`,
	},
	{
		Name: "jsonb_to_record",
		SetupSQL: `
			CREATE TABLE t_jtr (data jsonb);
			CREATE VIEW v_jsonb_to_record AS
			SELECT x.a, x.b FROM t_jtr, jsonb_to_record(data) AS x(a int, b text);`,
	},
	{
		Name: "jsonb_to_recordset",
		SetupSQL: `
			CREATE TABLE t_jtrs (data jsonb);
			CREATE VIEW v_jsonb_to_recordset AS
			SELECT x.a FROM t_jtrs, jsonb_to_recordset(data) AS x(a int);`,
	},
	{
		Name: "jsonb_populate_record",
		SetupSQL: `
			CREATE TYPE jpr_type AS (a int, b text);
			CREATE TABLE t_jpr (data jsonb);
			CREATE VIEW v_jsonb_populate_record AS
			SELECT (jsonb_populate_record(NULL::jpr_type, data)).* FROM t_jpr;`,
	},
	{
		Name: "jsonb_array_elements",
		SetupSQL: `
			CREATE TABLE t_jae (data jsonb);
			CREATE VIEW v_jsonb_array_elements AS
			SELECT jsonb_array_elements(data) FROM t_jae;`,
	},
	{
		Name: "json_each",
		SetupSQL: `
			CREATE TABLE t_je (data json);
			CREATE VIEW v_json_each AS
			SELECT (json_each(data)).* FROM t_je;`,
	},
	{
		Name: "jsonb_arrow",
		SetupSQL: `
			CREATE TABLE t_jarrow (data jsonb);
			CREATE VIEW v_jsonb_arrow AS
			SELECT data -> 'a' AS a, data ->> 'b' AS b FROM t_jarrow;`,
	},
	{
		Name: "jsonb_hash_arrow",
		SetupSQL: `
			CREATE TABLE t_jha (data jsonb);
			CREATE VIEW v_jsonb_hash_arrow AS
			SELECT data #> '{a,b}' AS nested, data #>> '{a,b}' AS nested_text FROM t_jha;`,
	},
	{
		Name: "jsonb_question",
		SetupSQL: `
			CREATE TABLE t_jq (data jsonb);
			CREATE VIEW v_jsonb_question AS
			SELECT data ? 'key' AS has_key FROM t_jq;`,
	},
	{
		Name: "jsonb_contains",
		SetupSQL: `
			CREATE TABLE t_jc (data jsonb);
			CREATE VIEW v_jsonb_contains AS
			SELECT data @> '{"a":1}'::jsonb AS contains_a FROM t_jc;`,
	},
	{
		Name: "jsonb_set",
		SetupSQL: `
			CREATE TABLE t_js (data jsonb);
			CREATE VIEW v_jsonb_set AS
			SELECT jsonb_set(data, '{a}', '"new"'::jsonb) FROM t_js;`,
	},
	{
		Name: "jsonb_strip_nulls",
		SetupSQL: `
			CREATE TABLE t_jsn (data jsonb);
			CREATE VIEW v_jsonb_strip_nulls AS
			SELECT jsonb_strip_nulls(data) FROM t_jsn;`,
	},
	{
		Name: "jsonb_minus",
		SetupSQL: `
			CREATE TABLE t_jm (data jsonb);
			CREATE VIEW v_jsonb_minus AS
			SELECT data - 'key' AS without_key, data #- '{a,b}' AS without_nested FROM t_jm;`,
	},
	{
		Name: "jsonb_typeof",
		SetupSQL: `
			CREATE TABLE t_jt (data jsonb);
			CREATE VIEW v_jsonb_typeof AS
			SELECT jsonb_typeof(data) FROM t_jt;`,
	},

	// ===================== OPERATOR() qualified =====================
	{
		Name: "operator_qualified_eq",
		SetupSQL: `
			CREATE TABLE t_opq (a int, b int);
			CREATE VIEW v_operator_qualified_eq AS
			SELECT * FROM t_opq WHERE a OPERATOR(pg_catalog.=) b;`,
	},
	{
		Name: "operator_qualified_plus",
		SetupSQL: `
			CREATE TABLE t_oqp (a int, b int);
			CREATE VIEW v_operator_qualified_plus AS
			SELECT a OPERATOR(pg_catalog.+) b AS total FROM t_oqp;`,
	},

	// ===================== OVERLAPS =====================
	{
		Name: "overlaps_dates",
		SetupSQL: `
			CREATE TABLE t_ovl (s date, e date);
			CREATE VIEW v_overlaps_dates AS
			SELECT * FROM t_ovl WHERE (s, e) OVERLAPS ('2020-01-01'::date, '2020-12-31'::date);`,
	},

	// ===================== Range / Multirange =====================
	{
		Name: "range_cast",
		SetupSQL: `
			CREATE VIEW v_range_cast AS SELECT '[1,10)'::int4range;`,
	},
	{
		Name: "range_operators",
		SetupSQL: `
			CREATE TABLE t_rng (r int4range);
			CREATE VIEW v_range_operators AS
			SELECT r, r @> 5 AS contains_5, r && '[3,7)'::int4range AS overlaps_37 FROM t_rng;`,
	},
	{
		Name: "numrange_func",
		SetupSQL: `
			CREATE VIEW v_numrange_func AS SELECT numrange(1.0, 10.0, '[)');`,
	},

	// ===================== Array operations =====================
	{
		Name: "array_agg_order",
		SetupSQL: `
			CREATE TABLE t_aao (a text, b int);
			CREATE VIEW v_array_agg_order AS
			SELECT array_agg(a ORDER BY b DESC) FROM t_aao;`,
	},
	{
		Name: "array_slice",
		SetupSQL: `
			CREATE TABLE t_as (arr int[]);
			CREATE VIEW v_array_slice AS
			SELECT arr[2:3] FROM t_as;`,
	},
	{
		Name: "array_any_all",
		SetupSQL: `
			CREATE TABLE t_aaa (val int, arr int[]);
			CREATE VIEW v_array_any_all AS
			SELECT val = ANY(arr) AS in_arr, val = ALL(arr) AS all_match FROM t_aaa;`,
	},
	{
		Name: "array_constructor_subquery",
		SetupSQL: `
			CREATE TABLE t_acs (a int);
			CREATE VIEW v_array_constructor_subquery AS
			SELECT ARRAY(SELECT a FROM t_acs ORDER BY a);`,
	},
	{
		Name: "unnest_multi",
		SetupSQL: `
			CREATE VIEW v_unnest_multi AS
			SELECT * FROM unnest(ARRAY[1,2], ARRAY['a','b']) AS u(num, letter);`,
	},

	// ===================== String aggregation =====================
	{
		Name: "string_agg_order",
		SetupSQL: `
			CREATE TABLE t_sao (a text, b int);
			CREATE VIEW v_string_agg_order AS
			SELECT string_agg(a, ',' ORDER BY b) FROM t_sao;`,
	},
	{
		Name: "string_agg_distinct",
		SetupSQL: `
			CREATE TABLE t_sad (a text);
			CREATE VIEW v_string_agg_distinct AS
			SELECT string_agg(DISTINCT a, ',') FROM t_sad;`,
	},

	// ===================== Regexp =====================
	{
		Name: "regexp_matches",
		SetupSQL: `
			CREATE TABLE t_rm (s text);
			CREATE VIEW v_regexp_matches AS
			SELECT regexp_matches(s, '(\w+)') FROM t_rm;`,
	},
	{
		Name: "regexp_split_to_table",
		SetupSQL: `
			CREATE TABLE t_rst (s text);
			CREATE VIEW v_regexp_split_to_table AS
			SELECT regexp_split_to_table(s, ',') FROM t_rst;`,
	},

	// ===================== Date/time =====================
	{
		Name: "extract_epoch",
		SetupSQL: `
			CREATE TABLE t_ee (ts timestamp);
			CREATE VIEW v_extract_epoch AS
			SELECT EXTRACT(EPOCH FROM ts) FROM t_ee;`,
	},
	{
		Name: "date_trunc",
		SetupSQL: `
			CREATE TABLE t_dt (ts timestamptz);
			CREATE VIEW v_date_trunc AS
			SELECT date_trunc('day', ts) FROM t_dt;`,
	},
	{
		Name: "at_time_zone",
		SetupSQL: `
			CREATE TABLE t_atz (ts timestamptz);
			CREATE VIEW v_at_time_zone AS
			SELECT ts AT TIME ZONE 'UTC' FROM t_atz;`,
	},
	{
		Name: "make_interval",
		SetupSQL: `
			CREATE VIEW v_make_interval AS
			SELECT make_interval(days => 1, hours => 2);`,
	},
	{
		Name: "age_func",
		SetupSQL: `
			CREATE TABLE t_age (born date);
			CREATE VIEW v_age_func AS
			SELECT age(now(), born::timestamp) FROM t_age;`,
	},
	{
		Name: "interval_arithmetic",
		SetupSQL: `
			CREATE TABLE t_ia (ts timestamp);
			CREATE VIEW v_interval_arithmetic AS
			SELECT ts + interval '1 day' AS tomorrow, ts - interval '1 hour' AS earlier FROM t_ia;`,
	},

	// ===================== Window functions =====================
	{
		Name: "window_rows_between",
		SetupSQL: `
			CREATE TABLE t_wrb (a int, b int);
			CREATE VIEW v_window_rows_between AS
			SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t_wrb;`,
	},
	{
		Name: "window_range_between",
		SetupSQL: `
			CREATE TABLE t_wranb (a int, b int);
			CREATE VIEW v_window_range_between AS
			SELECT sum(a) OVER (ORDER BY b RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t_wranb;`,
	},
	{
		Name: "window_groups_between",
		SetupSQL: `
			CREATE TABLE t_wgb (a int, b int);
			CREATE VIEW v_window_groups_between AS
			SELECT sum(a) OVER (ORDER BY b GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t_wgb;`,
	},
	{
		Name: "window_exclude",
		SetupSQL: `
			CREATE TABLE t_we (a int, b int);
			CREATE VIEW v_window_exclude AS
			SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE TIES) FROM t_we;`,
	},
	{
		Name: "window_named",
		SetupSQL: `
			CREATE TABLE t_wn (a int, b int);
			CREATE VIEW v_window_named AS
			SELECT sum(a) OVER w, avg(a) OVER w FROM t_wn WINDOW w AS (ORDER BY b);`,
	},
	{
		Name: "nth_value",
		SetupSQL: `
			CREATE TABLE t_nv (a text, b int);
			CREATE VIEW v_nth_value AS
			SELECT nth_value(a, 2) OVER (ORDER BY b) FROM t_nv;`,
	},

	// ===================== Aggregate FILTER =====================
	{
		Name: "count_filter",
		SetupSQL: `
			CREATE TABLE t_cf (a text, active bool);
			CREATE VIEW v_count_filter AS
			SELECT count(*) FILTER (WHERE active) AS active_count FROM t_cf;`,
	},
	{
		Name: "sum_filter",
		SetupSQL: `
			CREATE TABLE t_sf (amount int, category text);
			CREATE VIEW v_sum_filter AS
			SELECT sum(amount) FILTER (WHERE category = 'a') AS a_total FROM t_sf;`,
	},
	{
		Name: "agg_filter_with_cast",
		SetupSQL: `
			CREATE TABLE t_afc (a int, b bool);
			CREATE VIEW v_agg_filter_with_cast AS
			SELECT count(*)::int FILTER (WHERE b) FROM t_afc;`,
	},

	// ===================== WITHIN GROUP (ordered-set aggregates) =====================
	{
		Name: "percentile_cont",
		SetupSQL: `
			CREATE TABLE t_pc (val numeric);
			CREATE VIEW v_percentile_cont AS
			SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY val) FROM t_pc;`,
	},
	{
		Name: "percentile_disc",
		SetupSQL: `
			CREATE TABLE t_pd (val numeric);
			CREATE VIEW v_percentile_disc AS
			SELECT percentile_disc(0.5) WITHIN GROUP (ORDER BY val) FROM t_pd;`,
	},
	{
		Name: "mode_agg",
		SetupSQL: `
			CREATE TABLE t_mode (val text);
			CREATE VIEW v_mode_agg AS
			SELECT mode() WITHIN GROUP (ORDER BY val) FROM t_mode;`,
	},

	// ===================== GROUPING SETS / CUBE / ROLLUP =====================
	{
		Name: "grouping_sets",
		SetupSQL: `
			CREATE TABLE t_gs (a text, b text, c int);
			CREATE VIEW v_grouping_sets AS
			SELECT a, b, sum(c) FROM t_gs GROUP BY GROUPING SETS ((a), (b), ());`,
	},
	{
		Name: "cube",
		SetupSQL: `
			CREATE TABLE t_cube (a text, b text, c int);
			CREATE VIEW v_cube AS
			SELECT a, b, sum(c) FROM t_cube GROUP BY CUBE (a, b);`,
	},
	{
		Name: "rollup",
		SetupSQL: `
			CREATE TABLE t_rollup (a text, b text, c int);
			CREATE VIEW v_rollup AS
			SELECT a, b, sum(c) FROM t_rollup GROUP BY ROLLUP (a, b);`,
	},
	{
		Name: "grouping_func",
		SetupSQL: `
			CREATE TABLE t_gf (a text, c int);
			CREATE VIEW v_grouping_func AS
			SELECT a, sum(c), GROUPING(a) FROM t_gf GROUP BY ROLLUP (a);`,
	},

	// ===================== CTE variants =====================
	{
		Name: "cte_recursive",
		SetupSQL: `
			CREATE VIEW v_cte_recursive AS
			WITH RECURSIVE cte AS (
				SELECT 1 AS n
				UNION ALL
				SELECT n + 1 FROM cte WHERE n < 10
			) SELECT n FROM cte;`,
	},
	{
		Name: "cte_materialized",
		SetupSQL: `
			CREATE TABLE t_cm (a int);
			CREATE VIEW v_cte_materialized AS
			WITH cte AS MATERIALIZED (SELECT a FROM t_cm) SELECT * FROM cte;`,
	},
	{
		Name: "cte_not_materialized",
		SetupSQL: `
			CREATE TABLE t_cnm (a int);
			CREATE VIEW v_cte_not_materialized AS
			WITH cte AS NOT MATERIALIZED (SELECT a FROM t_cnm) SELECT * FROM cte;`,
	},

	// ===================== LATERAL =====================
	{
		Name: "lateral_subquery",
		SetupSQL: `
			CREATE TABLE t_lat (id int, n int);
			CREATE VIEW v_lateral_subquery AS
			SELECT t.id, s.x FROM t_lat t, LATERAL (SELECT generate_series(1, t.n) AS x) s;`,
	},
	{
		Name: "lateral_func",
		SetupSQL: `
			CREATE TABLE t_latf (id int, arr int[]);
			CREATE VIEW v_lateral_func AS
			SELECT t.id, u.val FROM t_latf t, LATERAL unnest(t.arr) AS u(val);`,
	},

	// ===================== TABLESAMPLE =====================
	{
		Name: "tablesample_system",
		SetupSQL: `
			CREATE TABLE t_ts (a int);
			CREATE VIEW v_tablesample_system AS
			SELECT * FROM t_ts TABLESAMPLE SYSTEM(10);`,
	},
	{
		Name: "tablesample_bernoulli",
		SetupSQL: `
			CREATE TABLE t_tsb (a int);
			CREATE VIEW v_tablesample_bernoulli AS
			SELECT * FROM t_tsb TABLESAMPLE BERNOULLI(10);`,
	},

	// ===================== VALUES in view =====================
	{
		Name: "values_view",
		SetupSQL: `
			CREATE VIEW v_values_view AS VALUES (1, 'a'), (2, 'b');`,
	},

	// ===================== DISTINCT ON =====================
	{
		Name: "distinct_on",
		SetupSQL: `
			CREATE TABLE t_do (a int, b text, c int);
			CREATE VIEW v_distinct_on AS
			SELECT DISTINCT ON (a) a, b, c FROM t_do ORDER BY a, c DESC;`,
	},

	// ===================== Full-text search =====================
	{
		Name: "to_tsvector_query",
		SetupSQL: `
			CREATE TABLE t_fts (doc text);
			CREATE VIEW v_to_tsvector_query AS
			SELECT doc FROM t_fts WHERE to_tsvector('english', doc) @@ to_tsquery('english', 'hello');`,
	},
	{
		Name: "ts_rank",
		SetupSQL: `
			CREATE TABLE t_tsr (doc text);
			CREATE VIEW v_ts_rank AS
			SELECT doc, ts_rank(to_tsvector('english', doc), to_tsquery('english', 'hello')) AS rank FROM t_tsr;`,
	},
	{
		Name: "ts_headline",
		SetupSQL: `
			CREATE TABLE t_tsh (doc text);
			CREATE VIEW v_ts_headline AS
			SELECT ts_headline('english', doc, to_tsquery('english', 'hello')) FROM t_tsh;`,
	},

	// ===================== Geometric types =====================
	{
		Name: "point_distance",
		SetupSQL: `
			CREATE TABLE t_geo (p point);
			CREATE VIEW v_point_distance AS
			SELECT p, p <-> point '(0,0)' AS dist FROM t_geo;`,
	},
	{
		Name: "box_intersect",
		SetupSQL: `
			CREATE TABLE t_box (b box);
			CREATE VIEW v_box_intersect AS
			SELECT b, b && box '((0,0),(1,1))' AS intersects FROM t_box;`,
	},

	// ===================== Network types =====================
	{
		Name: "inet_contains",
		SetupSQL: `
			CREATE TABLE t_inet (addr inet);
			CREATE VIEW v_inet_contains AS
			SELECT addr, addr << '192.168.0.0/16'::inet AS in_private FROM t_inet;`,
	},

	// ===================== Type casts — pg_get_viewdef patterns =====================
	{
		Name: "heavy_parens_cast",
		SetupSQL: `
			CREATE TABLE t_hpc (x numeric, y numeric);
			CREATE VIEW v_heavy_parens_cast AS
			SELECT ((x * 100.0) / y)::numeric(10,2) AS pct FROM t_hpc;`,
	},
	{
		Name: "pg_catalog_qualified_type",
		SetupSQL: `
			CREATE TABLE t_pgt (a int);
			CREATE VIEW v_pg_catalog_qualified_type AS
			SELECT a::pg_catalog.int8 FROM t_pgt;`,
	},
	{
		Name: "double_cast",
		SetupSQL: `
			CREATE VIEW v_double_cast AS SELECT ('now'::text)::date;`,
	},
	{
		Name: "cast_in_case",
		SetupSQL: `
			CREATE TABLE t_cic (a text);
			CREATE VIEW v_cast_in_case AS
			SELECT CASE WHEN a IS NULL THEN 'none'::text ELSE a END FROM t_cic;`,
	},
	{
		Name: "cast_in_coalesce",
		SetupSQL: `
			CREATE TABLE t_coal (a int, b int);
			CREATE VIEW v_cast_in_coalesce AS
			SELECT COALESCE(a, b, 0)::bigint FROM t_coal;`,
	},

	// ===================== Generated columns / Identity =====================
	{
		Name: "generated_stored",
		SetupSQL: `
			CREATE TABLE t_gen (a int, b int GENERATED ALWAYS AS (a * 2) STORED);
			CREATE VIEW v_generated_stored AS SELECT * FROM t_gen;`,
	},
	{
		Name: "identity_column",
		SetupSQL: `
			CREATE TABLE t_ident (id int GENERATED ALWAYS AS IDENTITY, val text);
			CREATE VIEW v_identity_column AS SELECT * FROM t_ident;`,
	},

	// ===================== Composite types =====================
	{
		Name: "composite_field_access",
		SetupSQL: `
			CREATE TYPE ct_type AS (x int, y text);
			CREATE TABLE t_ct (data ct_type);
			CREATE VIEW v_composite_field_access AS
			SELECT (data).x, (data).y FROM t_ct;`,
	},
	{
		Name: "row_constructor",
		SetupSQL: `
			CREATE TABLE t_rc (a int, b text);
			CREATE VIEW v_row_constructor AS
			SELECT ROW(a, b) FROM t_rc;`,
	},

	// ===================== Domain types =====================
	{
		Name: "domain_cast",
		SetupSQL: `
			CREATE DOMAIN posint AS int CHECK (VALUE > 0);
			CREATE TABLE t_dom (a posint);
			CREATE VIEW v_domain_cast AS SELECT a FROM t_dom;`,
	},

	// ===================== Enum types =====================
	{
		Name: "enum_cast",
		SetupSQL: `
			CREATE TYPE color AS ENUM ('red', 'green', 'blue');
			CREATE TABLE t_enum (c color);
			CREATE VIEW v_enum_cast AS SELECT c FROM t_enum WHERE c = 'red';`,
	},

	// ===================== Partitioned table =====================
	{
		Name: "partitioned_view",
		SetupSQL: `
			CREATE TABLE t_part (id int, created date) PARTITION BY RANGE (created);
			CREATE TABLE t_part_2020 PARTITION OF t_part FOR VALUES FROM ('2020-01-01') TO ('2021-01-01');
			CREATE VIEW v_partitioned_view AS SELECT * FROM t_part;`,
	},

	// ===================== RETURNING in CTE (writable CTE) =====================
	{
		Name: "writable_cte",
		SetupSQL: `
			CREATE TABLE t_wcte (id serial, val text);
			CREATE VIEW v_writable_cte AS
			WITH inserted AS (
				INSERT INTO t_wcte (val) VALUES ('x') RETURNING *
			) SELECT * FROM inserted;`,
	},

	// ===================== FETCH FIRST / LIMIT variants =====================
	{
		Name: "fetch_first",
		SetupSQL: `
			CREATE TABLE t_ff (a int);
			CREATE VIEW v_fetch_first AS
			SELECT a FROM t_ff ORDER BY a FETCH FIRST 5 ROWS ONLY;`,
	},
	{
		Name: "fetch_first_with_ties",
		SetupSQL: `
			CREATE TABLE t_fft (a int);
			CREATE VIEW v_fetch_first_with_ties AS
			SELECT a FROM t_fft ORDER BY a FETCH FIRST 5 ROWS WITH TIES;`,
	},
	{
		Name: "offset_fetch",
		SetupSQL: `
			CREATE TABLE t_of (a int);
			CREATE VIEW v_offset_fetch AS
			SELECT a FROM t_of ORDER BY a OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY;`,
	},

	// ===================== UNION / INTERSECT / EXCEPT =====================
	{
		Name: "union_all",
		SetupSQL: `
			CREATE TABLE t_ua1 (a int);
			CREATE TABLE t_ua2 (a int);
			CREATE VIEW v_union_all AS
			SELECT a FROM t_ua1 UNION ALL SELECT a FROM t_ua2;`,
	},
	{
		Name: "intersect",
		SetupSQL: `
			CREATE TABLE t_is1 (a int);
			CREATE TABLE t_is2 (a int);
			CREATE VIEW v_intersect AS
			SELECT a FROM t_is1 INTERSECT SELECT a FROM t_is2;`,
	},
	{
		Name: "except",
		SetupSQL: `
			CREATE TABLE t_ex1 (a int);
			CREATE TABLE t_ex2 (a int);
			CREATE VIEW v_except AS
			SELECT a FROM t_ex1 EXCEPT SELECT a FROM t_ex2;`,
	},

	// ===================== Subquery expressions =====================
	{
		Name: "exists_subquery",
		SetupSQL: `
			CREATE TABLE t_esq (id int);
			CREATE VIEW v_exists_subquery AS
			SELECT * FROM t_esq WHERE EXISTS (SELECT 1 FROM t_esq AS t2 WHERE t2.id = t_esq.id);`,
	},
	{
		Name: "in_subquery",
		SetupSQL: `
			CREATE TABLE t_insq (id int, parent_id int);
			CREATE VIEW v_in_subquery AS
			SELECT * FROM t_insq WHERE id IN (SELECT parent_id FROM t_insq);`,
	},
	{
		Name: "scalar_subquery",
		SetupSQL: `
			CREATE TABLE t_ssq (id int, val int);
			CREATE VIEW v_scalar_subquery AS
			SELECT id, (SELECT max(val) FROM t_ssq) AS max_val FROM t_ssq;`,
	},

	// ===================== NATURAL JOIN / USING =====================
	{
		Name: "natural_join",
		SetupSQL: `
			CREATE TABLE t_nj1 (id int, a text);
			CREATE TABLE t_nj2 (id int, b text);
			CREATE VIEW v_natural_join AS
			SELECT * FROM t_nj1 NATURAL JOIN t_nj2;`,
	},
	{
		Name: "join_using",
		SetupSQL: `
			CREATE TABLE t_ju1 (id int, a text);
			CREATE TABLE t_ju2 (id int, b text);
			CREATE VIEW v_join_using AS
			SELECT * FROM t_ju1 JOIN t_ju2 USING (id);`,
	},
	{
		Name: "full_outer_join",
		SetupSQL: `
			CREATE TABLE t_foj1 (id int, a text);
			CREATE TABLE t_foj2 (id int, b text);
			CREATE VIEW v_full_outer_join AS
			SELECT * FROM t_foj1 FULL OUTER JOIN t_foj2 ON t_foj1.id = t_foj2.id;`,
	},
	{
		Name: "cross_join",
		SetupSQL: `
			CREATE TABLE t_cj1 (a int);
			CREATE TABLE t_cj2 (b int);
			CREATE VIEW v_cross_join AS
			SELECT * FROM t_cj1 CROSS JOIN t_cj2;`,
	},

	// ===================== Special predicates =====================
	{
		Name: "is_distinct_from",
		SetupSQL: `
			CREATE TABLE t_idf (a int, b int);
			CREATE VIEW v_is_distinct_from AS
			SELECT * FROM t_idf WHERE a IS DISTINCT FROM b;`,
	},
	{
		Name: "is_not_distinct_from",
		SetupSQL: `
			CREATE TABLE t_indf (a int, b int);
			CREATE VIEW v_is_not_distinct_from AS
			SELECT * FROM t_indf WHERE a IS NOT DISTINCT FROM b;`,
	},
	{
		Name: "between_symmetric",
		SetupSQL: `
			CREATE TABLE t_bs (a int);
			CREATE VIEW v_between_symmetric AS
			SELECT * FROM t_bs WHERE a BETWEEN SYMMETRIC 10 AND 1;`,
	},
	{
		Name: "similar_to",
		SetupSQL: `
			CREATE TABLE t_st (a text);
			CREATE VIEW v_similar_to AS
			SELECT * FROM t_st WHERE a SIMILAR TO '%pattern%';`,
	},

	// ===================== Table functions =====================
	{
		Name: "generate_series_int",
		SetupSQL: `
			CREATE VIEW v_generate_series_int AS
			SELECT * FROM generate_series(1, 10) AS gs;`,
	},
	{
		Name: "generate_series_ts",
		SetupSQL: `
			CREATE VIEW v_generate_series_ts AS
			SELECT * FROM generate_series('2020-01-01'::timestamp, '2020-12-31'::timestamp, '1 month'::interval) AS gs;`,
	},
	{
		Name: "unnest_simple",
		SetupSQL: `
			CREATE VIEW v_unnest_simple AS
			SELECT * FROM unnest(ARRAY[1,2,3]) AS u;`,
	},

	// ===================== Security barrier view =====================
	{
		Name: "security_barrier",
		SetupSQL: `
			CREATE TABLE t_sb (id int, secret text, public_col text);
			CREATE VIEW v_security_barrier WITH (security_barrier=true) AS
			SELECT id, public_col FROM t_sb;`,
	},

	// ===================== Materialized view =====================
	{
		Name:     "materialized_view",
		ViewName: "mv_test",
		SetupSQL: `
			CREATE TABLE t_mv (a int, b text);
			INSERT INTO t_mv VALUES (1, 'x');
			CREATE MATERIALIZED VIEW mv_test AS SELECT a, b FROM t_mv;`,
	},

	// ===================== Recursive view =====================
	{
		Name: "recursive_view",
		SetupSQL: `
			CREATE RECURSIVE VIEW v_recursive_view (n) AS
			VALUES (1) UNION ALL SELECT n + 1 FROM v_recursive_view WHERE n < 10;`,
	},

	// ===================== Multiple FROM sources =====================
	{
		Name: "implicit_cross_join",
		SetupSQL: `
			CREATE TABLE t_icj1 (a int);
			CREATE TABLE t_icj2 (b int);
			CREATE VIEW v_implicit_cross_join AS
			SELECT a, b FROM t_icj1, t_icj2;`,
	},

	// ===================== GREATEST / LEAST / NULLIF =====================
	{
		Name: "greatest_least",
		SetupSQL: `
			CREATE TABLE t_gl (a int, b int, c int);
			CREATE VIEW v_greatest_least AS
			SELECT GREATEST(a, b, c) AS mx, LEAST(a, b, c) AS mn FROM t_gl;`,
	},
	{
		Name: "nullif_view",
		SetupSQL: `
			CREATE TABLE t_nif (a text);
			CREATE VIEW v_nullif_view AS
			SELECT NULLIF(a, '') FROM t_nif;`,
	},

	// ===================== Conditional / CASE =====================
	{
		Name: "case_searched",
		SetupSQL: `
			CREATE TABLE t_cs (a int);
			CREATE VIEW v_case_searched AS
			SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t_cs;`,
	},
	{
		Name: "case_simple",
		SetupSQL: `
			CREATE TABLE t_csm (status int);
			CREATE VIEW v_case_simple AS
			SELECT CASE status WHEN 1 THEN 'active' WHEN 2 THEN 'inactive' ELSE 'unknown' END FROM t_csm;`,
	},

	// ===================== Aggregate with ORDER BY =====================
	{
		Name: "json_agg_order",
		SetupSQL: `
			CREATE TABLE t_jao (a text, b int);
			CREATE VIEW v_json_agg_order AS
			SELECT json_agg(a ORDER BY b) FROM t_jao;`,
	},
	{
		Name: "jsonb_object_agg",
		SetupSQL: `
			CREATE TABLE t_joa (k text, v text);
			CREATE VIEW v_jsonb_object_agg AS
			SELECT jsonb_object_agg(k, v) FROM t_joa;`,
	},

	// ===================== System columns =====================
	{
		Name: "system_columns",
		SetupSQL: `
			CREATE TABLE t_sys (a int);
			CREATE VIEW v_system_columns AS
			SELECT tableoid, ctid, a FROM t_sys;`,
	},

	// ===================== FOR UPDATE / FOR SHARE (not in views, but subqueries) =====================
	// pg won't allow FOR UPDATE in views directly, skip

	// ===================== Complex real-world patterns =====================
	{
		Name: "complex_analytics",
		SetupSQL: `
			CREATE TABLE orders (id serial, customer_id int, amount numeric, created_at timestamptz, status text);
			CREATE VIEW v_complex_analytics AS
			SELECT
				customer_id,
				count(*) AS order_count,
				sum(amount) AS total_amount,
				avg(amount)::numeric(10,2) AS avg_amount,
				min(created_at) AS first_order,
				max(created_at) AS last_order,
				count(*) FILTER (WHERE status = 'completed') AS completed_count,
				sum(amount) FILTER (WHERE created_at > now() - interval '30 days') AS recent_total,
				array_agg(DISTINCT status ORDER BY status) AS statuses,
				row_number() OVER (ORDER BY sum(amount) DESC) AS rank
			FROM orders
			GROUP BY customer_id;`,
	},
	{
		Name: "complex_reporting",
		SetupSQL: `
			CREATE TABLE sales (id serial, region text, product text, amount numeric, sale_date date);
			CREATE VIEW v_complex_reporting AS
			WITH daily AS (
				SELECT region, product, sale_date, sum(amount) AS daily_total
				FROM sales GROUP BY region, product, sale_date
			),
			ranked AS (
				SELECT *,
					rank() OVER (PARTITION BY region ORDER BY daily_total DESC) AS rnk,
					sum(daily_total) OVER (PARTITION BY region ORDER BY sale_date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS rolling_7d
				FROM daily
			)
			SELECT region, product, sale_date, daily_total, rnk, rolling_7d,
				daily_total / NULLIF(sum(daily_total) OVER (PARTITION BY region, sale_date), 0) * 100 AS pct_of_day
			FROM ranked;`,
	},

	// ===================== Batch 2: deeper edge cases =====================

	// --- Subquery in FROM with column aliases ---
	{
		Name: "subquery_from_alias",
		SetupSQL: `
			CREATE TABLE t_sfa (a int, b text);
			CREATE VIEW v_subquery_from_alias AS
			SELECT sub.x, sub.y FROM (SELECT a AS x, b AS y FROM t_sfa) sub;`,
	},

	// --- Multiple LATERAL joins ---
	{
		Name: "multi_lateral",
		SetupSQL: `
			CREATE TABLE t_ml (id int, arr1 int[], arr2 text[]);
			CREATE VIEW v_multi_lateral AS
			SELECT t.id, u1.v1, u2.v2
			FROM t_ml t,
				LATERAL unnest(t.arr1) AS u1(v1),
				LATERAL unnest(t.arr2) AS u2(v2);`,
	},

	// --- CROSS JOIN LATERAL ---
	{
		Name: "cross_join_lateral",
		SetupSQL: `
			CREATE TABLE t_cjl (n int);
			CREATE VIEW v_cross_join_lateral AS
			SELECT t.n, gs FROM t_cjl t CROSS JOIN LATERAL generate_series(1, t.n) AS gs;`,
	},

	// --- Nested CTE ---
	{
		Name: "nested_cte",
		SetupSQL: `
			CREATE TABLE t_ncte (id int, parent_id int, name text);
			CREATE VIEW v_nested_cte AS
			WITH RECURSIVE tree AS (
				SELECT id, parent_id, name, 1 AS depth FROM t_ncte WHERE parent_id IS NULL
				UNION ALL
				SELECT c.id, c.parent_id, c.name, p.depth + 1
				FROM t_ncte c JOIN tree p ON c.parent_id = p.id
			)
			SELECT * FROM tree;`,
	},

	// --- USING with multiple columns ---
	{
		Name: "join_using_multi",
		SetupSQL: `
			CREATE TABLE t_jum1 (a int, b int, c text);
			CREATE TABLE t_jum2 (a int, b int, d text);
			CREATE VIEW v_join_using_multi AS
			SELECT * FROM t_jum1 JOIN t_jum2 USING (a, b);`,
	},

	// --- Self-join ---
	{
		Name: "self_join",
		SetupSQL: `
			CREATE TABLE t_sj (id int, parent_id int, name text);
			CREATE VIEW v_self_join AS
			SELECT c.name AS child, p.name AS parent
			FROM t_sj c LEFT JOIN t_sj p ON c.parent_id = p.id;`,
	},

	// --- Correlated subquery in SELECT ---
	{
		Name: "correlated_subquery",
		SetupSQL: `
			CREATE TABLE t_csq1 (id int, name text);
			CREATE TABLE t_csq2 (id int, parent_id int, val int);
			CREATE VIEW v_correlated_subquery AS
			SELECT t1.name,
				(SELECT sum(t2.val) FROM t_csq2 t2 WHERE t2.parent_id = t1.id) AS total
			FROM t_csq1 t1;`,
	},

	// --- ARRAY with subquery and cast ---
	{
		Name: "array_subquery_cast",
		SetupSQL: `
			CREATE TABLE t_asqc (id int, tag text);
			CREATE VIEW v_array_subquery_cast AS
			SELECT id, ARRAY(SELECT tag FROM t_asqc t2 WHERE t2.id = t.id)::text[] AS tags
			FROM (SELECT DISTINCT id FROM t_asqc) t;`,
	},

	// --- Multiple set-returning functions ---
	{
		Name: "multi_srf",
		SetupSQL: `
			CREATE VIEW v_multi_srf AS
			SELECT * FROM generate_series(1, 5) AS a, generate_series(1, 3) AS b;`,
	},

	// --- Aggregate with multiple ORDER BY ---
	{
		Name: "agg_multi_order",
		SetupSQL: `
			CREATE TABLE t_amo (grp text, a text, b int);
			CREATE VIEW v_agg_multi_order AS
			SELECT grp, string_agg(a, ',' ORDER BY b DESC, a ASC) FROM t_amo GROUP BY grp;`,
	},

	// --- COALESCE chain with different types ---
	{
		Name: "coalesce_chain_types",
		SetupSQL: `
			CREATE TABLE t_cct (a int, b bigint, c numeric);
			CREATE VIEW v_coalesce_chain_types AS
			SELECT COALESCE(a::numeric, b::numeric, c, 0) AS val FROM t_cct;`,
	},

	// --- Complex CASE with subqueries ---
	{
		Name: "case_with_subquery",
		SetupSQL: `
			CREATE TABLE t_cws (id int, status text);
			CREATE VIEW v_case_with_subquery AS
			SELECT id,
				CASE
					WHEN EXISTS (SELECT 1 FROM t_cws t2 WHERE t2.id = t.id AND t2.status = 'active') THEN 'has_active'
					ELSE 'no_active'
				END AS label
			FROM (SELECT DISTINCT id FROM t_cws) t;`,
	},

	// --- boolean expression patterns ---
	{
		Name: "boolean_is",
		SetupSQL: `
			CREATE TABLE t_bi (a bool, b bool);
			CREATE VIEW v_boolean_is AS
			SELECT * FROM t_bi WHERE a IS TRUE AND b IS NOT FALSE AND (a IS UNKNOWN) = false;`,
	},

	// --- IN with VALUES ---
	{
		Name: "in_values",
		SetupSQL: `
			CREATE TABLE t_iv (id int);
			CREATE VIEW v_in_values AS
			SELECT * FROM t_iv WHERE id IN (1, 2, 3, 4, 5);`,
	},

	// --- Nested function calls with casts ---
	{
		Name: "nested_func_cast",
		SetupSQL: `
			CREATE TABLE t_nfc (ts timestamptz);
			CREATE VIEW v_nested_func_cast AS
			SELECT
				date_trunc('day', ts AT TIME ZONE 'UTC')::date AS day,
				extract(epoch FROM ts)::bigint AS epoch_s,
				to_char(ts, 'YYYY-MM-DD"T"HH24:MI:SS') AS iso
			FROM t_nfc;`,
	},

	// --- String operators ---
	{
		Name: "string_operators",
		SetupSQL: `
			CREATE TABLE t_so (a text, b text);
			CREATE VIEW v_string_operators AS
			SELECT
				a || ' ' || b AS full_name,
				a LIKE 'test%' AS starts_with_test,
				a ILIKE '%TEST%' AS contains_test,
				a ~ '^[A-Z]' AS starts_upper,
				a ~* 'pattern' AS imatches
			FROM t_so;`,
	},

	// --- NULLS FIRST / LAST ---
	{
		Name: "nulls_ordering",
		SetupSQL: `
			CREATE TABLE t_no (a int);
			CREATE VIEW v_nulls_ordering AS
			SELECT a FROM t_no ORDER BY a NULLS FIRST;`,
	},

	// --- Subquery in WHERE with ALL/ANY ---
	{
		Name: "all_any_subquery",
		SetupSQL: `
			CREATE TABLE t_aas (val int);
			CREATE VIEW v_all_any_subquery AS
			SELECT * FROM t_aas WHERE val > ALL(SELECT val FROM t_aas WHERE val < 10);`,
	},

	// --- Complex window with PARTITION BY and frame ---
	{
		Name: "complex_window",
		SetupSQL: `
			CREATE TABLE t_cw (grp text, ts timestamp, val numeric);
			CREATE VIEW v_complex_window AS
			SELECT grp, ts, val,
				lag(val) OVER w AS prev_val,
				lead(val) OVER w AS next_val,
				sum(val) OVER (PARTITION BY grp ORDER BY ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
				ntile(4) OVER (PARTITION BY grp ORDER BY val) AS quartile,
				percent_rank() OVER (PARTITION BY grp ORDER BY val) AS pct_rank,
				cume_dist() OVER (PARTITION BY grp ORDER BY val) AS cume
			FROM t_cw
			WINDOW w AS (PARTITION BY grp ORDER BY ts);`,
	},

	// --- Aggregate FILTER with complex condition ---
	{
		Name: "agg_filter_complex",
		SetupSQL: `
			CREATE TABLE t_afc2 (category text, status text, amount numeric);
			CREATE VIEW v_agg_filter_complex AS
			SELECT category,
				count(*) FILTER (WHERE status IN ('active', 'pending')) AS active_count,
				sum(amount) FILTER (WHERE status = 'completed' AND amount > 100) AS big_completed,
				avg(amount) FILTER (WHERE amount BETWEEN 10 AND 1000) AS mid_avg
			FROM t_afc2
			GROUP BY category;`,
	},

	// --- WITH ORDINALITY ---
	{
		Name: "with_ordinality",
		SetupSQL: `
			CREATE VIEW v_with_ordinality AS
			SELECT val, ord FROM unnest(ARRAY['a','b','c']) WITH ORDINALITY AS t(val, ord);`,
	},

	// --- EXCEPT ALL / INTERSECT ALL ---
	{
		Name: "except_all",
		SetupSQL: `
			CREATE TABLE t_ea1 (a int);
			CREATE TABLE t_ea2 (a int);
			CREATE VIEW v_except_all AS
			SELECT a FROM t_ea1 EXCEPT ALL SELECT a FROM t_ea2;`,
	},

	// --- Multiple UNION with different casts ---
	{
		Name: "multi_union_cast",
		SetupSQL: `
			CREATE TABLE t_muc1 (a int);
			CREATE TABLE t_muc2 (a bigint);
			CREATE TABLE t_muc3 (a numeric);
			CREATE VIEW v_multi_union_cast AS
			SELECT a FROM t_muc1 UNION ALL SELECT a FROM t_muc2 UNION ALL SELECT a FROM t_muc3;`,
	},

	// --- HAVING ---
	{
		Name: "having_clause",
		SetupSQL: `
			CREATE TABLE t_hc (grp text, val int);
			CREATE VIEW v_having_clause AS
			SELECT grp, count(*), sum(val)
			FROM t_hc
			GROUP BY grp
			HAVING count(*) > 1 AND sum(val) > 100;`,
	},

	// --- Complex default expressions in table (test DDL parsing) ---
	{
		Name: "complex_defaults",
		SetupSQL: `
			CREATE TABLE t_cd (
				id bigint GENERATED ALWAYS AS IDENTITY,
				uid uuid DEFAULT gen_random_uuid(),
				created_at timestamptz DEFAULT now(),
				updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
				data jsonb DEFAULT '{}'::jsonb,
				tags text[] DEFAULT '{}'::text[],
				status text DEFAULT 'pending'::text,
				counter int DEFAULT 0
			);
			CREATE VIEW v_complex_defaults AS SELECT * FROM t_cd;`,
	},

	// --- Row comparison ---
	{
		Name: "row_comparison",
		SetupSQL: `
			CREATE TABLE t_rcmp (a int, b int);
			CREATE VIEW v_row_comparison AS
			SELECT * FROM t_rcmp WHERE ROW(a, b) > ROW(1, 2);`,
	},

	// --- TREAT/cast in composite type access ---
	{
		Name: "composite_nested",
		SetupSQL: `
			CREATE TYPE inner_type AS (x int, y int);
			CREATE TYPE outer_type AS (name text, coords inner_type);
			CREATE TABLE t_cn (data outer_type);
			CREATE VIEW v_composite_nested AS
			SELECT (data).name, ((data).coords).x, ((data).coords).y FROM t_cn;`,
	},

	// --- Enum comparison and ordering ---
	{
		Name: "enum_ordering",
		SetupSQL: `
			CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');
			CREATE TABLE t_eo (id int, p priority);
			CREATE VIEW v_enum_ordering AS
			SELECT * FROM t_eo WHERE p >= 'high' ORDER BY p;`,
	},

	// --- Inheritance ---
	{
		Name: "inheritance_view",
		SetupSQL: `
			CREATE TABLE t_parent (id int, name text);
			CREATE TABLE t_child (extra text) INHERITS (t_parent);
			CREATE VIEW v_inheritance_view AS
			SELECT * FROM ONLY t_parent;`,
	},

	// --- VARIADIC functions ---
	{
		Name: "variadic_func",
		SetupSQL: `
			CREATE FUNCTION debug_vd_concat(VARIADIC text[]) RETURNS text AS $$
				SELECT array_to_string($1, ',')
			$$ LANGUAGE SQL;
			CREATE TABLE t_vf (a text, b text, c text);
			CREATE VIEW v_variadic_func AS
			SELECT debug_vd_concat(a, b, c) FROM t_vf;`,
	},

	// --- ROWS FROM ---
	{
		Name: "rows_from",
		SetupSQL: `
			CREATE VIEW v_rows_from AS
			SELECT * FROM ROWS FROM (
				generate_series(1, 3),
				generate_series(1, 5)
			) AS t(a, b);`,
	},

	// --- hstore (extension) ---
	{
		Name:     "hstore_ext",
		ViewName: "v_hstore_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS hstore;
			CREATE TABLE t_hs (data hstore);
			CREATE VIEW v_hstore_ext AS
			SELECT data -> 'key' AS val, data ? 'key' AS has_key FROM t_hs;`,
	},

	// --- ltree (extension) ---
	{
		Name:     "ltree_ext",
		ViewName: "v_ltree_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS ltree;
			CREATE TABLE t_lt (path ltree);
			CREATE VIEW v_ltree_ext AS
			SELECT path, path <@ 'root.a'::ltree AS is_child FROM t_lt;`,
	},

	// --- citext (extension) ---
	{
		Name:     "citext_ext",
		ViewName: "v_citext_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS citext;
			CREATE TABLE t_ci (name citext);
			CREATE VIEW v_citext_ext AS
			SELECT name FROM t_ci WHERE name = 'Test';`,
	},

	// --- intarray (extension) ---
	{
		Name:     "intarray_ext",
		ViewName: "v_intarray_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS intarray;
			CREATE TABLE t_ia2 (arr int[]);
			CREATE VIEW v_intarray_ext AS
			SELECT arr, sort(arr) AS sorted FROM t_ia2;`,
	},

	// --- pg_trgm (extension) ---
	{
		Name:     "trgm_ext",
		ViewName: "v_trgm_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS pg_trgm;
			CREATE TABLE t_trgm (name text);
			CREATE VIEW v_trgm_ext AS
			SELECT name, similarity(name, 'test') AS sim FROM t_trgm;`,
	},

	// --- uuid-ossp (extension) ---
	{
		Name:     "uuid_ossp_ext",
		ViewName: "v_uuid_ossp_ext",
		SetupSQL: `
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
			CREATE VIEW v_uuid_ossp_ext AS
			SELECT uuid_generate_v4() AS id;`,
	},

	// --- Subquery in HAVING ---
	{
		Name: "having_subquery",
		SetupSQL: `
			CREATE TABLE t_hsq (grp int, val int);
			CREATE VIEW v_having_subquery AS
			SELECT grp, avg(val) FROM t_hsq
			GROUP BY grp
			HAVING avg(val) > (SELECT avg(val) FROM t_hsq);`,
	},

	// --- Multiple table expressions with same alias ---
	{
		Name: "duplicate_alias_subquery",
		SetupSQL: `
			CREATE TABLE t_das (id int, val int);
			CREATE VIEW v_duplicate_alias_subquery AS
			SELECT a.id, b.total
			FROM t_das a
			JOIN (SELECT id, sum(val) AS total FROM t_das GROUP BY id) b ON a.id = b.id;`,
	},

	// --- LIKE / NOT LIKE with ESCAPE ---
	{
		Name: "like_escape",
		SetupSQL: `
			CREATE TABLE t_le (pattern text);
			CREATE VIEW v_like_escape AS
			SELECT * FROM t_le WHERE pattern LIKE '%!%%' ESCAPE '!';`,
	},

	// --- Array comparison operators ---
	{
		Name: "array_compare",
		SetupSQL: `
			CREATE TABLE t_ac (arr int[]);
			CREATE VIEW v_array_compare AS
			SELECT arr,
				arr @> ARRAY[1,2] AS contains,
				arr <@ ARRAY[1,2,3,4,5] AS contained_by,
				arr && ARRAY[3,4] AS overlaps_arr
			FROM t_ac;`,
	},

	// --- ISNULL / NOTNULL (alternate syntax) ---
	{
		Name: "isnull_notnull",
		SetupSQL: `
			CREATE TABLE t_inn (a int);
			CREATE VIEW v_isnull_notnull AS
			SELECT * FROM t_inn WHERE a NOTNULL;`,
	},

	// --- OVERLAY ---
	{
		Name: "overlay_func",
		SetupSQL: `
			CREATE TABLE t_ov (s text);
			CREATE VIEW v_overlay_func AS
			SELECT OVERLAY(s PLACING 'XX' FROM 3 FOR 2) FROM t_ov;`,
	},

	// --- POSITION ---
	{
		Name: "position_func",
		SetupSQL: `
			CREATE TABLE t_pos (s text);
			CREATE VIEW v_position_func AS
			SELECT POSITION('needle' IN s) FROM t_pos;`,
	},

	// --- SUBSTRING with pattern ---
	{
		Name: "substring_pattern",
		SetupSQL: `
			CREATE TABLE t_sp (s text);
			CREATE VIEW v_substring_pattern AS
			SELECT SUBSTRING(s FROM 1 FOR 5) FROM t_sp;`,
	},

	// --- TRIM variants ---
	{
		Name: "trim_variants",
		SetupSQL: `
			CREATE TABLE t_tv (s text);
			CREATE VIEW v_trim_variants AS
			SELECT
				TRIM(BOTH ' ' FROM s) AS trimmed,
				TRIM(LEADING '0' FROM s) AS ltrimmed,
				TRIM(TRAILING FROM s) AS rtrimmed
			FROM t_tv;`,
	},

	// --- BIT operations ---
	{
		Name: "bit_operations",
		SetupSQL: `
			CREATE TABLE t_bo (a int, b int);
			CREATE VIEW v_bit_operations AS
			SELECT a & b AS band, a | b AS bor, a # b AS bxor, ~a AS bnot, a << 2 AS lshift, a >> 1 AS rshift FROM t_bo;`,
	},

	// --- CAST vs :: in different contexts ---
	{
		Name: "cast_vs_double_colon",
		SetupSQL: `
			CREATE TABLE t_cvdc (a text, b int);
			CREATE VIEW v_cast_vs_double_colon AS
			SELECT CAST(a AS int), b::text, CAST(b AS numeric(10,2)), a::varchar(100) FROM t_cvdc;`,
	},

	// --- NORMALIZE (PG13+) ---
	{
		Name: "normalize_func",
		SetupSQL: `
			CREATE TABLE t_norm (s text);
			CREATE VIEW v_normalize_func AS
			SELECT NORMALIZE(s, NFC) FROM t_norm;`,
	},

	// --- IS NORMALIZED (PG13+) ---
	{
		Name: "is_normalized",
		SetupSQL: `
			CREATE TABLE t_isn (s text);
			CREATE VIEW v_is_normalized AS
			SELECT s IS NFC NORMALIZED AS is_nfc FROM t_isn;`,
	},

	// --- Collation in aggregate ---
	{
		Name: "collation_in_agg",
		SetupSQL: `
			CREATE TABLE t_cia (name text);
			CREATE VIEW v_collation_in_agg AS
			SELECT min(name COLLATE "C") AS first_c FROM t_cia;`,
	},
}

func TestViewDefParse(t *testing.T) {
	ctr := startPGContainer(t)

	var passed, failed, skipped int
	var failures []string

	for _, tc := range viewDefCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			schema := ctr.freshSchema(t)

			// Run setup SQL in the isolated schema.
			_, err := ctr.db.ExecContext(ctr.ctx, fmt.Sprintf("SET search_path TO %q, public;\n%s", schema, tc.SetupSQL))
			if err != nil {
				t.Logf("SKIP: PG rejected setup SQL: %v", err)
				skipped++
				return
			}

			// Determine view name.
			viewName := tc.ViewName
			if viewName == "" {
				viewName = "v_" + tc.Name
			}

			// Get the canonical view definition from PG.
			var viewDef sql.NullString
			query := fmt.Sprintf("SELECT pg_get_viewdef('%q.%s'::regclass, true)", schema, viewName)
			err = ctr.db.QueryRowContext(ctr.ctx, query).Scan(&viewDef)
			if err != nil {
				// Try materialized view.
				query = fmt.Sprintf(`
					SELECT definition FROM pg_matviews
					WHERE schemaname = '%s' AND matviewname = '%s'`, schema, viewName)
				err = ctr.db.QueryRowContext(ctr.ctx, query).Scan(&viewDef)
				if err != nil {
					t.Logf("SKIP: could not get view definition: %v", err)
					skipped++
					return
				}
			}

			if !viewDef.Valid || viewDef.String == "" {
				t.Logf("SKIP: empty view definition")
				skipped++
				return
			}

			def := strings.TrimSpace(viewDef.String)
			def = strings.TrimSuffix(def, ";")

			// Wrap in CREATE VIEW to match the buildMinimalDDL pattern.
			fullSQL := fmt.Sprintf("CREATE VIEW test_parse AS %s;", def)

			// Attempt to parse with omni.
			_, parseErr := pg.Parse(fullSQL)
			if parseErr != nil {
				failed++
				failures = append(failures, fmt.Sprintf("[%s] %v\n  viewdef: %s", tc.Name, parseErr, def))
				t.Errorf("omni parser failed:\n  err: %v\n  viewdef: %s\n  full SQL: %s", parseErr, def, fullSQL)
			} else {
				passed++
			}
		})
	}

	// Summary.
	t.Logf("\n=== ViewDef Parse Summary ===")
	t.Logf("Passed:  %d", passed)
	t.Logf("Failed:  %d", failed)
	t.Logf("Skipped: %d", skipped)
	if len(failures) > 0 {
		t.Logf("\n--- Failures ---")
		for _, f := range failures {
			t.Logf("  %s", f)
		}
	}
}
