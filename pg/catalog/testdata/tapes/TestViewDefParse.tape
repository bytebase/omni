-- pin: default postgres:17.12-alpine sha256:34261f594522cf41dfcb76aefd42b9d201098b0e6b23c9c68733f169e4218b6f REL_17_12

== TestViewDefParse default#1
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_1"
< ok 0 -
> exec
| SET search_path TO "test_1", public;
|
| 			CREATE VIEW v_xmltable AS
| 			SELECT x.* FROM XMLTABLE('/root/row' PASSING '<root><row><a>1</a></row></root>'::xml
| 				COLUMNS a text PATH 'a') AS x;
< error
| 42601 syntax error at or near "::"
| position: 150
> exec
| DROP SCHEMA IF EXISTS "test_1" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_2"
< ok 0 -
> exec
| SET search_path TO "test_2", public;
|
| 			CREATE TABLE t_xmlexists (id serial, doc xml);
| 			CREATE VIEW v_xmlexists AS
| 			SELECT id FROM t_xmlexists WHERE XMLEXISTS('//a' PASSING BY VALUE doc);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_2".v_xmlexists'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id
|    FROM t_xmlexists
|   WHERE XMLEXISTS(('//a'::text) PASSING (doc));
> exec
| DROP SCHEMA IF EXISTS "test_2" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_3"
< ok 0 -
> exec
| SET search_path TO "test_3", public;
|
| 			CREATE VIEW v_xmlparse AS SELECT XMLPARSE(DOCUMENT '<root/>');
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_3".v_xmlparse'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLPARSE(DOCUMENT '<root/>'::text STRIP WHITESPACE) AS "xmlparse";
> exec
| DROP SCHEMA IF EXISTS "test_3" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_4"
< ok 0 -
> exec
| SET search_path TO "test_4", public;
|
| 			CREATE VIEW v_xmlserialize AS SELECT XMLSERIALIZE(DOCUMENT '<r/>'::xml AS text);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_4".v_xmlserialize'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLSERIALIZE(DOCUMENT '<r/>'::xml AS text NO INDENT) AS "xmlserialize";
> exec
| DROP SCHEMA IF EXISTS "test_4" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_5"
< ok 0 -
> exec
| SET search_path TO "test_5", public;
|
| 			CREATE TABLE t_xmlelement (name text);
| 			CREATE VIEW v_xmlelement AS SELECT XMLELEMENT(NAME item, name) FROM t_xmlelement;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_5".v_xmlelement'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLELEMENT(NAME item, name) AS "xmlelement"
|    FROM t_xmlelement;
> exec
| DROP SCHEMA IF EXISTS "test_5" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_6"
< ok 0 -
> exec
| SET search_path TO "test_6", public;
|
| 			CREATE TABLE t_xmlforest (a text, b int);
| 			CREATE VIEW v_xmlforest AS SELECT XMLFOREST(a, b) FROM t_xmlforest;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_6".v_xmlforest'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLFOREST(a AS a, b AS b) AS "xmlforest"
|    FROM t_xmlforest;
> exec
| DROP SCHEMA IF EXISTS "test_6" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_7"
< ok 0 -
> exec
| SET search_path TO "test_7", public;
|
| 			CREATE TABLE t_xmlagg (doc xml);
| 			CREATE VIEW v_xmlagg AS SELECT XMLAGG(doc) FROM t_xmlagg;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_7".v_xmlagg'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT xmlagg(doc) AS xmlagg
|    FROM t_xmlagg;
> exec
| DROP SCHEMA IF EXISTS "test_7" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_8"
< ok 0 -
> exec
| SET search_path TO "test_8", public;
|
| 			CREATE VIEW v_xmlconcat AS SELECT XMLCONCAT('<a/>'::xml, '<b/>'::xml);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_8".v_xmlconcat'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLCONCAT('<a/>'::xml, '<b/>'::xml) AS "xmlconcat";
> exec
| DROP SCHEMA IF EXISTS "test_8" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_9"
< ok 0 -
> exec
| SET search_path TO "test_9", public;
|
| 			CREATE VIEW v_xmlroot AS SELECT XMLROOT('<r/>'::xml, VERSION '1.0');
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_9".v_xmlroot'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT XMLROOT('<r/>'::xml, VERSION '1.0'::text) AS "xmlroot";
> exec
| DROP SCHEMA IF EXISTS "test_9" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_10"
< ok 0 -
> exec
| SET search_path TO "test_10", public;
|
| 			CREATE TABLE t_xpath (doc xml);
| 			CREATE VIEW v_xpath AS SELECT XPATH('/a', doc) FROM t_xpath;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_10".v_xpath'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT xpath('/a'::text, doc) AS xpath
|    FROM t_xpath;
> exec
| DROP SCHEMA IF EXISTS "test_10" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_11"
< ok 0 -
> exec
| SET search_path TO "test_11", public;
|
| 			CREATE TABLE t_xpath_exists (doc xml);
| 			CREATE VIEW v_xpath_exists AS SELECT XPATH_EXISTS('/a', doc) FROM t_xpath_exists;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_11".v_xpath_exists'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT xpath_exists('/a'::text, doc) AS xpath_exists
|    FROM t_xpath_exists;
> exec
| DROP SCHEMA IF EXISTS "test_11" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_12"
< ok 0 -
> exec
| SET search_path TO "test_12", public;
|
| 			CREATE TABLE t_jpq (data jsonb);
| 			CREATE VIEW v_jsonb_path_query AS SELECT jsonb_path_query(data, '$.a') FROM t_jpq;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_12".v_jsonb_path_query'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_path_query(data, '$."a"'::jsonpath) AS jsonb_path_query
|    FROM t_jpq;
> exec
| DROP SCHEMA IF EXISTS "test_12" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_13"
< ok 0 -
> exec
| SET search_path TO "test_13", public;
|
| 			CREATE TABLE t_jtr (data jsonb);
| 			CREATE VIEW v_jsonb_to_record AS
| 			SELECT x.a, x.b FROM t_jtr, jsonb_to_record(data) AS x(a int, b text);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_13".v_jsonb_to_record'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT x.a,
|     x.b
|    FROM t_jtr,
|     LATERAL jsonb_to_record(t_jtr.data) x(a integer, b text);
> exec
| DROP SCHEMA IF EXISTS "test_13" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_14"
< ok 0 -
> exec
| SET search_path TO "test_14", public;
|
| 			CREATE TABLE t_jtrs (data jsonb);
| 			CREATE VIEW v_jsonb_to_recordset AS
| 			SELECT x.a FROM t_jtrs, jsonb_to_recordset(data) AS x(a int);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_14".v_jsonb_to_recordset'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT x.a
|    FROM t_jtrs,
|     LATERAL jsonb_to_recordset(t_jtrs.data) x(a integer);
> exec
| DROP SCHEMA IF EXISTS "test_14" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_15"
< ok 0 -
> exec
| SET search_path TO "test_15", public;
|
| 			CREATE TYPE jpr_type AS (a int, b text);
| 			CREATE TABLE t_jpr (data jsonb);
| 			CREATE VIEW v_jsonb_populate_record AS
| 			SELECT (jsonb_populate_record(NULL::jpr_type, data)).* FROM t_jpr;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_15".v_jsonb_populate_record'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (jsonb_populate_record(NULL::jpr_type, data)).a AS a,
|     (jsonb_populate_record(NULL::jpr_type, data)).b AS b
|    FROM t_jpr;
> exec
| DROP SCHEMA IF EXISTS "test_15" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_16"
< ok 0 -
> exec
| SET search_path TO "test_16", public;
|
| 			CREATE TABLE t_jae (data jsonb);
| 			CREATE VIEW v_jsonb_array_elements AS
| 			SELECT jsonb_array_elements(data) FROM t_jae;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_16".v_jsonb_array_elements'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_array_elements(data) AS jsonb_array_elements
|    FROM t_jae;
> exec
| DROP SCHEMA IF EXISTS "test_16" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_17"
< ok 0 -
> exec
| SET search_path TO "test_17", public;
|
| 			CREATE TABLE t_je (data json);
| 			CREATE VIEW v_json_each AS
| 			SELECT (json_each(data)).* FROM t_je;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_17".v_json_each'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (json_each(data)).key AS key,
|     (json_each(data)).value AS value
|    FROM t_je;
> exec
| DROP SCHEMA IF EXISTS "test_17" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_18"
< ok 0 -
> exec
| SET search_path TO "test_18", public;
|
| 			CREATE TABLE t_jarrow (data jsonb);
| 			CREATE VIEW v_jsonb_arrow AS
| 			SELECT data -> 'a' AS a, data ->> 'b' AS b FROM t_jarrow;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_18".v_jsonb_arrow'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data -> 'a'::text AS a,
|     data ->> 'b'::text AS b
|    FROM t_jarrow;
> exec
| DROP SCHEMA IF EXISTS "test_18" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_19"
< ok 0 -
> exec
| SET search_path TO "test_19", public;
|
| 			CREATE TABLE t_jha (data jsonb);
| 			CREATE VIEW v_jsonb_hash_arrow AS
| 			SELECT data #> '{a,b}' AS nested, data #>> '{a,b}' AS nested_text FROM t_jha;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_19".v_jsonb_hash_arrow'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data #> '{a,b}'::text[] AS nested,
|     data #>> '{a,b}'::text[] AS nested_text
|    FROM t_jha;
> exec
| DROP SCHEMA IF EXISTS "test_19" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_20"
< ok 0 -
> exec
| SET search_path TO "test_20", public;
|
| 			CREATE TABLE t_jq (data jsonb);
| 			CREATE VIEW v_jsonb_question AS
| 			SELECT data ? 'key' AS has_key FROM t_jq;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_20".v_jsonb_question'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data ? 'key'::text AS has_key
|    FROM t_jq;
> exec
| DROP SCHEMA IF EXISTS "test_20" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_21" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_21"
< ok 0 -
> exec
| SET search_path TO "test_21", public;
|
| 			CREATE TABLE t_jc (data jsonb);
| 			CREATE VIEW v_jsonb_contains AS
| 			SELECT data @> '{"a":1}'::jsonb AS contains_a FROM t_jc;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_21".v_jsonb_contains'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data @> '{"a": 1}'::jsonb AS contains_a
|    FROM t_jc;
> exec
| DROP SCHEMA IF EXISTS "test_21" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_22" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_22"
< ok 0 -
> exec
| SET search_path TO "test_22", public;
|
| 			CREATE TABLE t_js (data jsonb);
| 			CREATE VIEW v_jsonb_set AS
| 			SELECT jsonb_set(data, '{a}', '"new"'::jsonb) FROM t_js;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_22".v_jsonb_set'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_set(data, '{a}'::text[], '"new"'::jsonb) AS jsonb_set
|    FROM t_js;
> exec
| DROP SCHEMA IF EXISTS "test_22" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_23" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_23"
< ok 0 -
> exec
| SET search_path TO "test_23", public;
|
| 			CREATE TABLE t_jsn (data jsonb);
| 			CREATE VIEW v_jsonb_strip_nulls AS
| 			SELECT jsonb_strip_nulls(data) FROM t_jsn;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_23".v_jsonb_strip_nulls'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_strip_nulls(data) AS jsonb_strip_nulls
|    FROM t_jsn;
> exec
| DROP SCHEMA IF EXISTS "test_23" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_24" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_24"
< ok 0 -
> exec
| SET search_path TO "test_24", public;
|
| 			CREATE TABLE t_jm (data jsonb);
| 			CREATE VIEW v_jsonb_minus AS
| 			SELECT data - 'key' AS without_key, data #- '{a,b}' AS without_nested FROM t_jm;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_24".v_jsonb_minus'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data - 'key'::text AS without_key,
|     data #- '{a,b}'::text[] AS without_nested
|    FROM t_jm;
> exec
| DROP SCHEMA IF EXISTS "test_24" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_25" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_25"
< ok 0 -
> exec
| SET search_path TO "test_25", public;
|
| 			CREATE TABLE t_jt (data jsonb);
| 			CREATE VIEW v_jsonb_typeof AS
| 			SELECT jsonb_typeof(data) FROM t_jt;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_25".v_jsonb_typeof'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_typeof(data) AS jsonb_typeof
|    FROM t_jt;
> exec
| DROP SCHEMA IF EXISTS "test_25" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_26" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_26"
< ok 0 -
> exec
| SET search_path TO "test_26", public;
|
| 			CREATE TABLE t_opq (a int, b int);
| 			CREATE VIEW v_operator_qualified_eq AS
| 			SELECT * FROM t_opq WHERE a OPERATOR(pg_catalog.=) b;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_26".v_operator_qualified_eq'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_opq
|   WHERE a = b;
> exec
| DROP SCHEMA IF EXISTS "test_26" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_27" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_27"
< ok 0 -
> exec
| SET search_path TO "test_27", public;
|
| 			CREATE TABLE t_oqp (a int, b int);
| 			CREATE VIEW v_operator_qualified_plus AS
| 			SELECT a OPERATOR(pg_catalog.+) b AS total FROM t_oqp;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_27".v_operator_qualified_plus'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a + b AS total
|    FROM t_oqp;
> exec
| DROP SCHEMA IF EXISTS "test_27" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_28" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_28"
< ok 0 -
> exec
| SET search_path TO "test_28", public;
|
| 			CREATE TABLE t_ovl (s date, e date);
| 			CREATE VIEW v_overlaps_dates AS
| 			SELECT * FROM t_ovl WHERE (s, e) OVERLAPS ('2020-01-01'::date, '2020-12-31'::date);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_28".v_overlaps_dates'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT s,
|     e
|    FROM t_ovl
|   WHERE ((s, e) OVERLAPS ('2020-01-01'::date, '2020-12-31'::date));
> exec
| DROP SCHEMA IF EXISTS "test_28" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_29" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_29"
< ok 0 -
> exec
| SET search_path TO "test_29", public;
|
| 			CREATE VIEW v_range_cast AS SELECT '[1,10)'::int4range;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_29".v_range_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT '[1,10)'::int4range AS int4range;
> exec
| DROP SCHEMA IF EXISTS "test_29" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_30" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_30"
< ok 0 -
> exec
| SET search_path TO "test_30", public;
|
| 			CREATE TABLE t_rng (r int4range);
| 			CREATE VIEW v_range_operators AS
| 			SELECT r, r @> 5 AS contains_5, r && '[3,7)'::int4range AS overlaps_37 FROM t_rng;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_30".v_range_operators'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT r,
|     r @> 5 AS contains_5,
|     r && '[3,7)'::int4range AS overlaps_37
|    FROM t_rng;
> exec
| DROP SCHEMA IF EXISTS "test_30" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_31" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_31"
< ok 0 -
> exec
| SET search_path TO "test_31", public;
|
| 			CREATE VIEW v_numrange_func AS SELECT numrange(1.0, 10.0, '[)');
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_31".v_numrange_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT numrange(1.0, 10.0, '[)'::text) AS numrange;
> exec
| DROP SCHEMA IF EXISTS "test_31" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_32" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_32"
< ok 0 -
> exec
| SET search_path TO "test_32", public;
|
| 			CREATE TABLE t_aao (a text, b int);
| 			CREATE VIEW v_array_agg_order AS
| 			SELECT array_agg(a ORDER BY b DESC) FROM t_aao;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_32".v_array_agg_order'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT array_agg(a ORDER BY b DESC) AS array_agg
|    FROM t_aao;
> exec
| DROP SCHEMA IF EXISTS "test_32" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_33" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_33"
< ok 0 -
> exec
| SET search_path TO "test_33", public;
|
| 			CREATE TABLE t_as (arr int[]);
| 			CREATE VIEW v_array_slice AS
| 			SELECT arr[2:3] FROM t_as;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_33".v_array_slice'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT arr[2:3] AS arr
|    FROM t_as;
> exec
| DROP SCHEMA IF EXISTS "test_33" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_34" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_34"
< ok 0 -
> exec
| SET search_path TO "test_34", public;
|
| 			CREATE TABLE t_aaa (val int, arr int[]);
| 			CREATE VIEW v_array_any_all AS
| 			SELECT val = ANY(arr) AS in_arr, val = ALL(arr) AS all_match FROM t_aaa;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_34".v_array_any_all'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT val = ANY (arr) AS in_arr,
|     val = ALL (arr) AS all_match
|    FROM t_aaa;
> exec
| DROP SCHEMA IF EXISTS "test_34" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_35" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_35"
< ok 0 -
> exec
| SET search_path TO "test_35", public;
|
| 			CREATE TABLE t_acs (a int);
| 			CREATE VIEW v_array_constructor_subquery AS
| 			SELECT ARRAY(SELECT a FROM t_acs ORDER BY a);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_35".v_array_constructor_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT ARRAY( SELECT t_acs.a
|            FROM t_acs
|           ORDER BY t_acs.a) AS "array";
> exec
| DROP SCHEMA IF EXISTS "test_35" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_36" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_36"
< ok 0 -
> exec
| SET search_path TO "test_36", public;
|
| 			CREATE VIEW v_unnest_multi AS
| 			SELECT * FROM unnest(ARRAY[1,2], ARRAY['a','b']) AS u(num, letter);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_36".v_unnest_multi'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT num,
|     letter
|    FROM UNNEST(ARRAY[1, 2], ARRAY['a'::text, 'b'::text]) u(num, letter);
> exec
| DROP SCHEMA IF EXISTS "test_36" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_37" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_37"
< ok 0 -
> exec
| SET search_path TO "test_37", public;
|
| 			CREATE TABLE t_sao (a text, b int);
| 			CREATE VIEW v_string_agg_order AS
| 			SELECT string_agg(a, ',' ORDER BY b) FROM t_sao;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_37".v_string_agg_order'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(a, ','::text ORDER BY b) AS string_agg
|    FROM t_sao;
> exec
| DROP SCHEMA IF EXISTS "test_37" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_38" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_38"
< ok 0 -
> exec
| SET search_path TO "test_38", public;
|
| 			CREATE TABLE t_sad (a text);
| 			CREATE VIEW v_string_agg_distinct AS
| 			SELECT string_agg(DISTINCT a, ',') FROM t_sad;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_38".v_string_agg_distinct'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT string_agg(DISTINCT a, ','::text) AS string_agg
|    FROM t_sad;
> exec
| DROP SCHEMA IF EXISTS "test_38" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_39" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_39"
< ok 0 -
> exec
| SET search_path TO "test_39", public;
|
| 			CREATE TABLE t_rm (s text);
| 			CREATE VIEW v_regexp_matches AS
| 			SELECT regexp_matches(s, '(\w+)') FROM t_rm;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_39".v_regexp_matches'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT regexp_matches(s, '(\w+)'::text) AS regexp_matches
|    FROM t_rm;
> exec
| DROP SCHEMA IF EXISTS "test_39" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_40" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_40"
< ok 0 -
> exec
| SET search_path TO "test_40", public;
|
| 			CREATE TABLE t_rst (s text);
| 			CREATE VIEW v_regexp_split_to_table AS
| 			SELECT regexp_split_to_table(s, ',') FROM t_rst;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_40".v_regexp_split_to_table'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT regexp_split_to_table(s, ','::text) AS regexp_split_to_table
|    FROM t_rst;
> exec
| DROP SCHEMA IF EXISTS "test_40" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_41" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_41"
< ok 0 -
> exec
| SET search_path TO "test_41", public;
|
| 			CREATE TABLE t_ee (ts timestamp);
| 			CREATE VIEW v_extract_epoch AS
| 			SELECT EXTRACT(EPOCH FROM ts) FROM t_ee;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_41".v_extract_epoch'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT EXTRACT(epoch FROM ts) AS "extract"
|    FROM t_ee;
> exec
| DROP SCHEMA IF EXISTS "test_41" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_42" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_42"
< ok 0 -
> exec
| SET search_path TO "test_42", public;
|
| 			CREATE TABLE t_dt (ts timestamptz);
| 			CREATE VIEW v_date_trunc AS
| 			SELECT date_trunc('day', ts) FROM t_dt;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_42".v_date_trunc'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT date_trunc('day'::text, ts) AS date_trunc
|    FROM t_dt;
> exec
| DROP SCHEMA IF EXISTS "test_42" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_43" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_43"
< ok 0 -
> exec
| SET search_path TO "test_43", public;
|
| 			CREATE TABLE t_atz (ts timestamptz);
| 			CREATE VIEW v_at_time_zone AS
| 			SELECT ts AT TIME ZONE 'UTC' FROM t_atz;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_43".v_at_time_zone'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (ts AT TIME ZONE 'UTC'::text) AS timezone
|    FROM t_atz;
> exec
| DROP SCHEMA IF EXISTS "test_43" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_44" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_44"
< ok 0 -
> exec
| SET search_path TO "test_44", public;
|
| 			CREATE VIEW v_make_interval AS
| 			SELECT make_interval(days => 1, hours => 2);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_44".v_make_interval'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT make_interval(days => 1, hours => 2) AS make_interval;
> exec
| DROP SCHEMA IF EXISTS "test_44" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_45" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_45"
< ok 0 -
> exec
| SET search_path TO "test_45", public;
|
| 			CREATE TABLE t_age (born date);
| 			CREATE VIEW v_age_func AS
| 			SELECT age(now(), born::timestamp) FROM t_age;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_45".v_age_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT age(now(), born::timestamp without time zone::timestamp with time zone) AS age
|    FROM t_age;
> exec
| DROP SCHEMA IF EXISTS "test_45" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_46" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_46"
< ok 0 -
> exec
| SET search_path TO "test_46", public;
|
| 			CREATE TABLE t_ia (ts timestamp);
| 			CREATE VIEW v_interval_arithmetic AS
| 			SELECT ts + interval '1 day' AS tomorrow, ts - interval '1 hour' AS earlier FROM t_ia;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_46".v_interval_arithmetic'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT ts + '1 day'::interval AS tomorrow,
|     ts - '01:00:00'::interval AS earlier
|    FROM t_ia;
> exec
| DROP SCHEMA IF EXISTS "test_46" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_47" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_47"
< ok 0 -
> exec
| SET search_path TO "test_47", public;
|
| 			CREATE TABLE t_wrb (a int, b int);
| 			CREATE VIEW v_window_rows_between AS
| 			SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t_wrb;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_47".v_window_rows_between'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS sum
|    FROM t_wrb;
> exec
| DROP SCHEMA IF EXISTS "test_47" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_48" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_48"
< ok 0 -
> exec
| SET search_path TO "test_48", public;
|
| 			CREATE TABLE t_wranb (a int, b int);
| 			CREATE VIEW v_window_range_between AS
| 			SELECT sum(a) OVER (ORDER BY b RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) FROM t_wranb;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_48".v_window_range_between'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(a) OVER (ORDER BY b RANGE BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS sum
|    FROM t_wranb;
> exec
| DROP SCHEMA IF EXISTS "test_48" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_49" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_49"
< ok 0 -
> exec
| SET search_path TO "test_49", public;
|
| 			CREATE TABLE t_wgb (a int, b int);
| 			CREATE VIEW v_window_groups_between AS
| 			SELECT sum(a) OVER (ORDER BY b GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING) FROM t_wgb;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_49".v_window_groups_between'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(a) OVER (ORDER BY b GROUPS BETWEEN 1 PRECEDING AND 1 FOLLOWING) AS sum
|    FROM t_wgb;
> exec
| DROP SCHEMA IF EXISTS "test_49" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_50" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_50"
< ok 0 -
> exec
| SET search_path TO "test_50", public;
|
| 			CREATE TABLE t_we (a int, b int);
| 			CREATE VIEW v_window_exclude AS
| 			SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE TIES) FROM t_we;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_50".v_window_exclude'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(a) OVER (ORDER BY b ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW EXCLUDE TIES) AS sum
|    FROM t_we;
> exec
| DROP SCHEMA IF EXISTS "test_50" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_51" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_51"
< ok 0 -
> exec
| SET search_path TO "test_51", public;
|
| 			CREATE TABLE t_wn (a int, b int);
| 			CREATE VIEW v_window_named AS
| 			SELECT sum(a) OVER w, avg(a) OVER w FROM t_wn WINDOW w AS (ORDER BY b);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_51".v_window_named'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(a) OVER w AS sum,
|     avg(a) OVER w AS avg
|    FROM t_wn
|   WINDOW w AS (ORDER BY b);
> exec
| DROP SCHEMA IF EXISTS "test_51" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_52" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_52"
< ok 0 -
> exec
| SET search_path TO "test_52", public;
|
| 			CREATE TABLE t_nv (a text, b int);
| 			CREATE VIEW v_nth_value AS
| 			SELECT nth_value(a, 2) OVER (ORDER BY b) FROM t_nv;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_52".v_nth_value'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT nth_value(a, 2) OVER (ORDER BY b) AS nth_value
|    FROM t_nv;
> exec
| DROP SCHEMA IF EXISTS "test_52" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_53" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_53"
< ok 0 -
> exec
| SET search_path TO "test_53", public;
|
| 			CREATE TABLE t_cf (a text, active bool);
| 			CREATE VIEW v_count_filter AS
| 			SELECT count(*) FILTER (WHERE active) AS active_count FROM t_cf;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_53".v_count_filter'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT count(*) FILTER (WHERE active) AS active_count
|    FROM t_cf;
> exec
| DROP SCHEMA IF EXISTS "test_53" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_54" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_54"
< ok 0 -
> exec
| SET search_path TO "test_54", public;
|
| 			CREATE TABLE t_sf (amount int, category text);
| 			CREATE VIEW v_sum_filter AS
| 			SELECT sum(amount) FILTER (WHERE category = 'a') AS a_total FROM t_sf;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_54".v_sum_filter'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT sum(amount) FILTER (WHERE category = 'a'::text) AS a_total
|    FROM t_sf;
> exec
| DROP SCHEMA IF EXISTS "test_54" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_55" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_55"
< ok 0 -
> exec
| SET search_path TO "test_55", public;
|
| 			CREATE TABLE t_afc (a int, b bool);
| 			CREATE VIEW v_agg_filter_with_cast AS
| 			SELECT count(*)::int FILTER (WHERE b) FROM t_afc;
< error
| 42601 syntax error at or near "FILTER"
| position: 144
> exec
| DROP SCHEMA IF EXISTS "test_55" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_56" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_56"
< ok 0 -
> exec
| SET search_path TO "test_56", public;
|
| 			CREATE TABLE t_pc (val numeric);
| 			CREATE VIEW v_percentile_cont AS
| 			SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY val) FROM t_pc;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_56".v_percentile_cont'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT percentile_cont(0.5::double precision) WITHIN GROUP (ORDER BY (val::double precision)) AS percentile_cont
|    FROM t_pc;
> exec
| DROP SCHEMA IF EXISTS "test_56" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_57" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_57"
< ok 0 -
> exec
| SET search_path TO "test_57", public;
|
| 			CREATE TABLE t_pd (val numeric);
| 			CREATE VIEW v_percentile_disc AS
| 			SELECT percentile_disc(0.5) WITHIN GROUP (ORDER BY val) FROM t_pd;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_57".v_percentile_disc'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT percentile_disc(0.5::double precision) WITHIN GROUP (ORDER BY val) AS percentile_disc
|    FROM t_pd;
> exec
| DROP SCHEMA IF EXISTS "test_57" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_58" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_58"
< ok 0 -
> exec
| SET search_path TO "test_58", public;
|
| 			CREATE TABLE t_mode (val text);
| 			CREATE VIEW v_mode_agg AS
| 			SELECT mode() WITHIN GROUP (ORDER BY val) FROM t_mode;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_58".v_mode_agg'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT mode() WITHIN GROUP (ORDER BY val) AS mode
|    FROM t_mode;
> exec
| DROP SCHEMA IF EXISTS "test_58" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_59" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_59"
< ok 0 -
> exec
| SET search_path TO "test_59", public;
|
| 			CREATE TABLE t_gs (a text, b text, c int);
| 			CREATE VIEW v_grouping_sets AS
| 			SELECT a, b, sum(c) FROM t_gs GROUP BY GROUPING SETS ((a), (b), ());
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_59".v_grouping_sets'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b,
|     sum(c) AS sum
|    FROM t_gs
|   GROUP BY GROUPING SETS ((a), (b), ());
> exec
| DROP SCHEMA IF EXISTS "test_59" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_60" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_60"
< ok 0 -
> exec
| SET search_path TO "test_60", public;
|
| 			CREATE TABLE t_cube (a text, b text, c int);
| 			CREATE VIEW v_cube AS
| 			SELECT a, b, sum(c) FROM t_cube GROUP BY CUBE (a, b);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_60".v_cube'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b,
|     sum(c) AS sum
|    FROM t_cube
|   GROUP BY CUBE(a, b);
> exec
| DROP SCHEMA IF EXISTS "test_60" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_61" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_61"
< ok 0 -
> exec
| SET search_path TO "test_61", public;
|
| 			CREATE TABLE t_rollup (a text, b text, c int);
| 			CREATE VIEW v_rollup AS
| 			SELECT a, b, sum(c) FROM t_rollup GROUP BY ROLLUP (a, b);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_61".v_rollup'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b,
|     sum(c) AS sum
|    FROM t_rollup
|   GROUP BY ROLLUP(a, b);
> exec
| DROP SCHEMA IF EXISTS "test_61" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_62" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_62"
< ok 0 -
> exec
| SET search_path TO "test_62", public;
|
| 			CREATE TABLE t_gf (a text, c int);
| 			CREATE VIEW v_grouping_func AS
| 			SELECT a, sum(c), GROUPING(a) FROM t_gf GROUP BY ROLLUP (a);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_62".v_grouping_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     sum(c) AS sum,
|     GROUPING(a) AS "grouping"
|    FROM t_gf
|   GROUP BY ROLLUP(a);
> exec
| DROP SCHEMA IF EXISTS "test_62" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_63" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_63"
< ok 0 -
> exec
| SET search_path TO "test_63", public;
|
| 			CREATE VIEW v_cte_recursive AS
| 			WITH RECURSIVE cte AS (
| 				SELECT 1 AS n
| 				UNION ALL
| 				SELECT n + 1 FROM cte WHERE n < 10
| 			) SELECT n FROM cte;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_63".v_cte_recursive'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH RECURSIVE cte AS (
|          SELECT 1 AS n
|         UNION ALL
|          SELECT cte_1.n + 1
|            FROM cte cte_1
|           WHERE cte_1.n < 10
|         )
|  SELECT n
|    FROM cte;
> exec
| DROP SCHEMA IF EXISTS "test_63" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_64" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_64"
< ok 0 -
> exec
| SET search_path TO "test_64", public;
|
| 			CREATE TABLE t_cm (a int);
| 			CREATE VIEW v_cte_materialized AS
| 			WITH cte AS MATERIALIZED (SELECT a FROM t_cm) SELECT * FROM cte;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_64".v_cte_materialized'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH cte AS MATERIALIZED (
|          SELECT t_cm.a
|            FROM t_cm
|         )
|  SELECT a
|    FROM cte;
> exec
| DROP SCHEMA IF EXISTS "test_64" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_65" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_65"
< ok 0 -
> exec
| SET search_path TO "test_65", public;
|
| 			CREATE TABLE t_cnm (a int);
| 			CREATE VIEW v_cte_not_materialized AS
| 			WITH cte AS NOT MATERIALIZED (SELECT a FROM t_cnm) SELECT * FROM cte;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_65".v_cte_not_materialized'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH cte AS NOT MATERIALIZED (
|          SELECT t_cnm.a
|            FROM t_cnm
|         )
|  SELECT a
|    FROM cte;
> exec
| DROP SCHEMA IF EXISTS "test_65" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_66" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_66"
< ok 0 -
> exec
| SET search_path TO "test_66", public;
|
| 			CREATE TABLE t_lat (id int, n int);
| 			CREATE VIEW v_lateral_subquery AS
| 			SELECT t.id, s.x FROM t_lat t, LATERAL (SELECT generate_series(1, t.n) AS x) s;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_66".v_lateral_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     s.x
|    FROM t_lat t,
|     LATERAL ( SELECT generate_series(1, t.n) AS x) s;
> exec
| DROP SCHEMA IF EXISTS "test_66" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_67" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_67"
< ok 0 -
> exec
| SET search_path TO "test_67", public;
|
| 			CREATE TABLE t_latf (id int, arr int[]);
| 			CREATE VIEW v_lateral_func AS
| 			SELECT t.id, u.val FROM t_latf t, LATERAL unnest(t.arr) AS u(val);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_67".v_lateral_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     u.val
|    FROM t_latf t,
|     LATERAL unnest(t.arr) u(val);
> exec
| DROP SCHEMA IF EXISTS "test_67" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_68" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_68"
< ok 0 -
> exec
| SET search_path TO "test_68", public;
|
| 			CREATE TABLE t_ts (a int);
| 			CREATE VIEW v_tablesample_system AS
| 			SELECT * FROM t_ts TABLESAMPLE SYSTEM(10);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_68".v_tablesample_system'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_ts TABLESAMPLE system (10);
> exec
| DROP SCHEMA IF EXISTS "test_68" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_69" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_69"
< ok 0 -
> exec
| SET search_path TO "test_69", public;
|
| 			CREATE TABLE t_tsb (a int);
| 			CREATE VIEW v_tablesample_bernoulli AS
| 			SELECT * FROM t_tsb TABLESAMPLE BERNOULLI(10);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_69".v_tablesample_bernoulli'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_tsb TABLESAMPLE bernoulli (10);
> exec
| DROP SCHEMA IF EXISTS "test_69" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_70" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_70"
< ok 0 -
> exec
| SET search_path TO "test_70", public;
|
| 			CREATE VIEW v_values_view AS VALUES (1, 'a'), (2, 'b');
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_70".v_values_view'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  VALUES (1,'a'::text), (2,'b'::text);
> exec
| DROP SCHEMA IF EXISTS "test_70" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_71" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_71"
< ok 0 -
> exec
| SET search_path TO "test_71", public;
|
| 			CREATE TABLE t_do (a int, b text, c int);
| 			CREATE VIEW v_distinct_on AS
| 			SELECT DISTINCT ON (a) a, b, c FROM t_do ORDER BY a, c DESC;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_71".v_distinct_on'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT DISTINCT ON (a) a,
|     b,
|     c
|    FROM t_do
|   ORDER BY a, c DESC;
> exec
| DROP SCHEMA IF EXISTS "test_71" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_72" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_72"
< ok 0 -
> exec
| SET search_path TO "test_72", public;
|
| 			CREATE TABLE t_fts (doc text);
| 			CREATE VIEW v_to_tsvector_query AS
| 			SELECT doc FROM t_fts WHERE to_tsvector('english', doc) @@ to_tsquery('english', 'hello');
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_72".v_to_tsvector_query'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT doc
|    FROM t_fts
|   WHERE to_tsvector('english'::regconfig, doc) @@ to_tsquery('english'::regconfig, 'hello'::text);
> exec
| DROP SCHEMA IF EXISTS "test_72" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_73" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_73"
< ok 0 -
> exec
| SET search_path TO "test_73", public;
|
| 			CREATE TABLE t_tsr (doc text);
| 			CREATE VIEW v_ts_rank AS
| 			SELECT doc, ts_rank(to_tsvector('english', doc), to_tsquery('english', 'hello')) AS rank FROM t_tsr;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_73".v_ts_rank'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT doc,
|     ts_rank(to_tsvector('english'::regconfig, doc), to_tsquery('english'::regconfig, 'hello'::text)) AS rank
|    FROM t_tsr;
> exec
| DROP SCHEMA IF EXISTS "test_73" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_74" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_74"
< ok 0 -
> exec
| SET search_path TO "test_74", public;
|
| 			CREATE TABLE t_tsh (doc text);
| 			CREATE VIEW v_ts_headline AS
| 			SELECT ts_headline('english', doc, to_tsquery('english', 'hello')) FROM t_tsh;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_74".v_ts_headline'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT ts_headline('english'::regconfig, doc, to_tsquery('english'::regconfig, 'hello'::text)) AS ts_headline
|    FROM t_tsh;
> exec
| DROP SCHEMA IF EXISTS "test_74" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_75" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_75"
< ok 0 -
> exec
| SET search_path TO "test_75", public;
|
| 			CREATE TABLE t_geo (p point);
| 			CREATE VIEW v_point_distance AS
| 			SELECT p, p <-> point '(0,0)' AS dist FROM t_geo;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_75".v_point_distance'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT p,
|     p <-> '(0,0)'::point AS dist
|    FROM t_geo;
> exec
| DROP SCHEMA IF EXISTS "test_75" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_76" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_76"
< ok 0 -
> exec
| SET search_path TO "test_76", public;
|
| 			CREATE TABLE t_box (b box);
| 			CREATE VIEW v_box_intersect AS
| 			SELECT b, b && box '((0,0),(1,1))' AS intersects FROM t_box;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_76".v_box_intersect'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT b,
|     b && '(1,1),(0,0)'::box AS intersects
|    FROM t_box;
> exec
| DROP SCHEMA IF EXISTS "test_76" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_77" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_77"
< ok 0 -
> exec
| SET search_path TO "test_77", public;
|
| 			CREATE TABLE t_inet (addr inet);
| 			CREATE VIEW v_inet_contains AS
| 			SELECT addr, addr << '192.168.0.0/16'::inet AS in_private FROM t_inet;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_77".v_inet_contains'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT addr,
|     addr << '192.168.0.0/16'::inet AS in_private
|    FROM t_inet;
> exec
| DROP SCHEMA IF EXISTS "test_77" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_78" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_78"
< ok 0 -
> exec
| SET search_path TO "test_78", public;
|
| 			CREATE TABLE t_hpc (x numeric, y numeric);
| 			CREATE VIEW v_heavy_parens_cast AS
| 			SELECT ((x * 100.0) / y)::numeric(10,2) AS pct FROM t_hpc;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_78".v_heavy_parens_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (x * 100.0 / y)::numeric(10,2) AS pct
|    FROM t_hpc;
> exec
| DROP SCHEMA IF EXISTS "test_78" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_79" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_79"
< ok 0 -
> exec
| SET search_path TO "test_79", public;
|
| 			CREATE TABLE t_pgt (a int);
| 			CREATE VIEW v_pg_catalog_qualified_type AS
| 			SELECT a::pg_catalog.int8 FROM t_pgt;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_79".v_pg_catalog_qualified_type'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a::bigint AS a
|    FROM t_pgt;
> exec
| DROP SCHEMA IF EXISTS "test_79" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_80" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_80"
< ok 0 -
> exec
| SET search_path TO "test_80", public;
|
| 			CREATE VIEW v_double_cast AS SELECT ('now'::text)::date;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_80".v_double_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT 'now'::text::date AS date;
> exec
| DROP SCHEMA IF EXISTS "test_80" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_81" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_81"
< ok 0 -
> exec
| SET search_path TO "test_81", public;
|
| 			CREATE TABLE t_cic (a text);
| 			CREATE VIEW v_cast_in_case AS
| 			SELECT CASE WHEN a IS NULL THEN 'none'::text ELSE a END FROM t_cic;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_81".v_cast_in_case'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT
|         CASE
|             WHEN a IS NULL THEN 'none'::text
|             ELSE a
|         END AS a
|    FROM t_cic;
> exec
| DROP SCHEMA IF EXISTS "test_81" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_82" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_82"
< ok 0 -
> exec
| SET search_path TO "test_82", public;
|
| 			CREATE TABLE t_coal (a int, b int);
| 			CREATE VIEW v_cast_in_coalesce AS
| 			SELECT COALESCE(a, b, 0)::bigint FROM t_coal;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_82".v_cast_in_coalesce'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT COALESCE(a, b, 0)::bigint AS "coalesce"
|    FROM t_coal;
> exec
| DROP SCHEMA IF EXISTS "test_82" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_83" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_83"
< ok 0 -
> exec
| SET search_path TO "test_83", public;
|
| 			CREATE TABLE t_gen (a int, b int GENERATED ALWAYS AS (a * 2) STORED);
| 			CREATE VIEW v_generated_stored AS SELECT * FROM t_gen;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_83".v_generated_stored'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_gen;
> exec
| DROP SCHEMA IF EXISTS "test_83" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_84" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_84"
< ok 0 -
> exec
| SET search_path TO "test_84", public;
|
| 			CREATE TABLE t_ident (id int GENERATED ALWAYS AS IDENTITY, val text);
| 			CREATE VIEW v_identity_column AS SELECT * FROM t_ident;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_84".v_identity_column'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     val
|    FROM t_ident;
> exec
| DROP SCHEMA IF EXISTS "test_84" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_85" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_85"
< ok 0 -
> exec
| SET search_path TO "test_85", public;
|
| 			CREATE TYPE ct_type AS (x int, y text);
| 			CREATE TABLE t_ct (data ct_type);
| 			CREATE VIEW v_composite_field_access AS
| 			SELECT (data).x, (data).y FROM t_ct;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_85".v_composite_field_access'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (data).x AS x,
|     (data).y AS y
|    FROM t_ct;
> exec
| DROP SCHEMA IF EXISTS "test_85" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_86" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_86"
< ok 0 -
> exec
| SET search_path TO "test_86", public;
|
| 			CREATE TABLE t_rc (a int, b text);
| 			CREATE VIEW v_row_constructor AS
| 			SELECT ROW(a, b) FROM t_rc;
< error
| 42P16 column "row" has pseudo-type record
> exec
| DROP SCHEMA IF EXISTS "test_86" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_87" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_87"
< ok 0 -
> exec
| SET search_path TO "test_87", public;
|
| 			CREATE DOMAIN posint AS int CHECK (VALUE > 0);
| 			CREATE TABLE t_dom (a posint);
| 			CREATE VIEW v_domain_cast AS SELECT a FROM t_dom;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_87".v_domain_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_dom;
> exec
| DROP SCHEMA IF EXISTS "test_87" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_88" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_88"
< ok 0 -
> exec
| SET search_path TO "test_88", public;
|
| 			CREATE TYPE color AS ENUM ('red', 'green', 'blue');
| 			CREATE TABLE t_enum (c color);
| 			CREATE VIEW v_enum_cast AS SELECT c FROM t_enum WHERE c = 'red';
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_88".v_enum_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT c
|    FROM t_enum
|   WHERE c = 'red'::color;
> exec
| DROP SCHEMA IF EXISTS "test_88" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_89" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_89"
< ok 0 -
> exec
| SET search_path TO "test_89", public;
|
| 			CREATE TABLE t_part (id int, created date) PARTITION BY RANGE (created);
| 			CREATE TABLE t_part_2020 PARTITION OF t_part FOR VALUES FROM ('2020-01-01') TO ('2021-01-01');
| 			CREATE VIEW v_partitioned_view AS SELECT * FROM t_part;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_89".v_partitioned_view'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     created
|    FROM t_part;
> exec
| DROP SCHEMA IF EXISTS "test_89" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_90" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_90"
< ok 0 -
> exec
| SET search_path TO "test_90", public;
|
| 			CREATE TABLE t_wcte (id serial, val text);
| 			CREATE VIEW v_writable_cte AS
| 			WITH inserted AS (
| 				INSERT INTO t_wcte (val) VALUES ('x') RETURNING *
| 			) SELECT * FROM inserted;
< error
| 0A000 views must not contain data-modifying statements in WITH
> exec
| DROP SCHEMA IF EXISTS "test_90" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_91" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_91"
< ok 0 -
> exec
| SET search_path TO "test_91", public;
|
| 			CREATE TABLE t_ff (a int);
| 			CREATE VIEW v_fetch_first AS
| 			SELECT a FROM t_ff ORDER BY a FETCH FIRST 5 ROWS ONLY;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_91".v_fetch_first'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_ff
|   ORDER BY a
|  LIMIT 5;
> exec
| DROP SCHEMA IF EXISTS "test_91" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_92" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_92"
< ok 0 -
> exec
| SET search_path TO "test_92", public;
|
| 			CREATE TABLE t_fft (a int);
| 			CREATE VIEW v_fetch_first_with_ties AS
| 			SELECT a FROM t_fft ORDER BY a FETCH FIRST 5 ROWS WITH TIES;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_92".v_fetch_first_with_ties'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_fft
|   ORDER BY a
|  FETCH FIRST (5) ROWS WITH TIES;
> exec
| DROP SCHEMA IF EXISTS "test_92" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_93" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_93"
< ok 0 -
> exec
| SET search_path TO "test_93", public;
|
| 			CREATE TABLE t_of (a int);
| 			CREATE VIEW v_offset_fetch AS
| 			SELECT a FROM t_of ORDER BY a OFFSET 10 ROWS FETCH NEXT 5 ROWS ONLY;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_93".v_offset_fetch'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_of
|   ORDER BY a
|  OFFSET 10
|  LIMIT 5;
> exec
| DROP SCHEMA IF EXISTS "test_93" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_94" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_94"
< ok 0 -
> exec
| SET search_path TO "test_94", public;
|
| 			CREATE TABLE t_ua1 (a int);
| 			CREATE TABLE t_ua2 (a int);
| 			CREATE VIEW v_union_all AS
| 			SELECT a FROM t_ua1 UNION ALL SELECT a FROM t_ua2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_94".v_union_all'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_ua1.a
|    FROM t_ua1
| UNION ALL
|  SELECT t_ua2.a
|    FROM t_ua2;
> exec
| DROP SCHEMA IF EXISTS "test_94" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_95" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_95"
< ok 0 -
> exec
| SET search_path TO "test_95", public;
|
| 			CREATE TABLE t_is1 (a int);
| 			CREATE TABLE t_is2 (a int);
| 			CREATE VIEW v_intersect AS
| 			SELECT a FROM t_is1 INTERSECT SELECT a FROM t_is2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_95".v_intersect'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_is1.a
|    FROM t_is1
| INTERSECT
|  SELECT t_is2.a
|    FROM t_is2;
> exec
| DROP SCHEMA IF EXISTS "test_95" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_96" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_96"
< ok 0 -
> exec
| SET search_path TO "test_96", public;
|
| 			CREATE TABLE t_ex1 (a int);
| 			CREATE TABLE t_ex2 (a int);
| 			CREATE VIEW v_except AS
| 			SELECT a FROM t_ex1 EXCEPT SELECT a FROM t_ex2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_96".v_except'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_ex1.a
|    FROM t_ex1
| EXCEPT
|  SELECT t_ex2.a
|    FROM t_ex2;
> exec
| DROP SCHEMA IF EXISTS "test_96" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_97" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_97"
< ok 0 -
> exec
| SET search_path TO "test_97", public;
|
| 			CREATE TABLE t_esq (id int);
| 			CREATE VIEW v_exists_subquery AS
| 			SELECT * FROM t_esq WHERE EXISTS (SELECT 1 FROM t_esq AS t2 WHERE t2.id = t_esq.id);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_97".v_exists_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id
|    FROM t_esq
|   WHERE (EXISTS ( SELECT 1
|            FROM t_esq t2
|           WHERE t2.id = t_esq.id));
> exec
| DROP SCHEMA IF EXISTS "test_97" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_98" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_98"
< ok 0 -
> exec
| SET search_path TO "test_98", public;
|
| 			CREATE TABLE t_insq (id int, parent_id int);
| 			CREATE VIEW v_in_subquery AS
| 			SELECT * FROM t_insq WHERE id IN (SELECT parent_id FROM t_insq);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_98".v_in_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     parent_id
|    FROM t_insq
|   WHERE (id IN ( SELECT t_insq_1.parent_id
|            FROM t_insq t_insq_1));
> exec
| DROP SCHEMA IF EXISTS "test_98" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_99" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_99"
< ok 0 -
> exec
| SET search_path TO "test_99", public;
|
| 			CREATE TABLE t_ssq (id int, val int);
| 			CREATE VIEW v_scalar_subquery AS
| 			SELECT id, (SELECT max(val) FROM t_ssq) AS max_val FROM t_ssq;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_99".v_scalar_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     ( SELECT max(t_ssq_1.val) AS max
|            FROM t_ssq t_ssq_1) AS max_val
|    FROM t_ssq;
> exec
| DROP SCHEMA IF EXISTS "test_99" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_100" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_100"
< ok 0 -
> exec
| SET search_path TO "test_100", public;
|
| 			CREATE TABLE t_nj1 (id int, a text);
| 			CREATE TABLE t_nj2 (id int, b text);
| 			CREATE VIEW v_natural_join AS
| 			SELECT * FROM t_nj1 NATURAL JOIN t_nj2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_100".v_natural_join'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_nj1.id,
|     t_nj1.a,
|     t_nj2.b
|    FROM t_nj1
|      JOIN t_nj2 USING (id);
> exec
| DROP SCHEMA IF EXISTS "test_100" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_101" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_101"
< ok 0 -
> exec
| SET search_path TO "test_101", public;
|
| 			CREATE TABLE t_ju1 (id int, a text);
| 			CREATE TABLE t_ju2 (id int, b text);
| 			CREATE VIEW v_join_using AS
| 			SELECT * FROM t_ju1 JOIN t_ju2 USING (id);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_101".v_join_using'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_ju1.id,
|     t_ju1.a,
|     t_ju2.b
|    FROM t_ju1
|      JOIN t_ju2 USING (id);
> exec
| DROP SCHEMA IF EXISTS "test_101" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_102" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_102"
< ok 0 -
> exec
| SET search_path TO "test_102", public;
|
| 			CREATE TABLE t_foj1 (id int, a text);
| 			CREATE TABLE t_foj2 (id int, b text);
| 			CREATE VIEW v_full_outer_join AS
| 			SELECT * FROM t_foj1 FULL OUTER JOIN t_foj2 ON t_foj1.id = t_foj2.id;
< error
| 42701 column "id" specified more than once
> exec
| DROP SCHEMA IF EXISTS "test_102" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_103" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_103"
< ok 0 -
> exec
| SET search_path TO "test_103", public;
|
| 			CREATE TABLE t_cj1 (a int);
| 			CREATE TABLE t_cj2 (b int);
| 			CREATE VIEW v_cross_join AS
| 			SELECT * FROM t_cj1 CROSS JOIN t_cj2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_103".v_cross_join'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_cj1.a,
|     t_cj2.b
|    FROM t_cj1
|      CROSS JOIN t_cj2;
> exec
| DROP SCHEMA IF EXISTS "test_103" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_104" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_104"
< ok 0 -
> exec
| SET search_path TO "test_104", public;
|
| 			CREATE TABLE t_idf (a int, b int);
| 			CREATE VIEW v_is_distinct_from AS
| 			SELECT * FROM t_idf WHERE a IS DISTINCT FROM b;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_104".v_is_distinct_from'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_idf
|   WHERE a IS DISTINCT FROM b;
> exec
| DROP SCHEMA IF EXISTS "test_104" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_105" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_105"
< ok 0 -
> exec
| SET search_path TO "test_105", public;
|
| 			CREATE TABLE t_indf (a int, b int);
| 			CREATE VIEW v_is_not_distinct_from AS
| 			SELECT * FROM t_indf WHERE a IS NOT DISTINCT FROM b;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_105".v_is_not_distinct_from'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_indf
|   WHERE NOT a IS DISTINCT FROM b;
> exec
| DROP SCHEMA IF EXISTS "test_105" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_106" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_106"
< ok 0 -
> exec
| SET search_path TO "test_106", public;
|
| 			CREATE TABLE t_bs (a int);
| 			CREATE VIEW v_between_symmetric AS
| 			SELECT * FROM t_bs WHERE a BETWEEN SYMMETRIC 10 AND 1;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_106".v_between_symmetric'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_bs
|   WHERE a >= 10 AND a <= 1 OR a >= 1 AND a <= 10;
> exec
| DROP SCHEMA IF EXISTS "test_106" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_107" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_107"
< ok 0 -
> exec
| SET search_path TO "test_107", public;
|
| 			CREATE TABLE t_st (a text);
| 			CREATE VIEW v_similar_to AS
| 			SELECT * FROM t_st WHERE a SIMILAR TO '%pattern%';
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_107".v_similar_to'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_st
|   WHERE a ~ similar_to_escape('%pattern%'::text);
> exec
| DROP SCHEMA IF EXISTS "test_107" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_108" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_108"
< ok 0 -
> exec
| SET search_path TO "test_108", public;
|
| 			CREATE VIEW v_generate_series_int AS
| 			SELECT * FROM generate_series(1, 10) AS gs;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_108".v_generate_series_int'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT gs
|    FROM generate_series(1, 10) gs(gs);
> exec
| DROP SCHEMA IF EXISTS "test_108" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_109" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_109"
< ok 0 -
> exec
| SET search_path TO "test_109", public;
|
| 			CREATE VIEW v_generate_series_ts AS
| 			SELECT * FROM generate_series('2020-01-01'::timestamp, '2020-12-31'::timestamp, '1 month'::interval) AS gs;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_109".v_generate_series_ts'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT gs
|    FROM generate_series('2020-01-01 00:00:00'::timestamp without time zone, '2020-12-31 00:00:00'::timestamp without time zone, '1 mon'::interval) gs(gs);
> exec
| DROP SCHEMA IF EXISTS "test_109" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_110" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_110"
< ok 0 -
> exec
| SET search_path TO "test_110", public;
|
| 			CREATE VIEW v_unnest_simple AS
| 			SELECT * FROM unnest(ARRAY[1,2,3]) AS u;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_110".v_unnest_simple'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT u
|    FROM unnest(ARRAY[1, 2, 3]) u(u);
> exec
| DROP SCHEMA IF EXISTS "test_110" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_111" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_111"
< ok 0 -
> exec
| SET search_path TO "test_111", public;
|
| 			CREATE TABLE t_sb (id int, secret text, public_col text);
| 			CREATE VIEW v_security_barrier WITH (security_barrier=true) AS
| 			SELECT id, public_col FROM t_sb;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_111".v_security_barrier'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     public_col
|    FROM t_sb;
> exec
| DROP SCHEMA IF EXISTS "test_111" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_112" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_112"
< ok 0 -
> exec
| SET search_path TO "test_112", public;
|
| 			CREATE TABLE t_mv (a int, b text);
| 			INSERT INTO t_mv VALUES (1, 'x');
| 			CREATE MATERIALIZED VIEW mv_test AS SELECT a, b FROM t_mv;
< ok 1 -
> query
| SELECT pg_get_viewdef('"test_112".mv_test'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_mv;
> exec
| DROP SCHEMA IF EXISTS "test_112" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_113" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_113"
< ok 0 -
> exec
| SET search_path TO "test_113", public;
|
| 			CREATE RECURSIVE VIEW v_recursive_view (n) AS
| 			VALUES (1) UNION ALL SELECT n + 1 FROM v_recursive_view WHERE n < 10;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_113".v_recursive_view'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH RECURSIVE v_recursive_view(n) AS (
|          VALUES (1)
|         UNION ALL
|          SELECT v_recursive_view_1.n + 1
|            FROM v_recursive_view v_recursive_view_1
|           WHERE v_recursive_view_1.n < 10
|         )
|  SELECT n
|    FROM v_recursive_view;
> exec
| DROP SCHEMA IF EXISTS "test_113" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_114" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_114"
< ok 0 -
> exec
| SET search_path TO "test_114", public;
|
| 			CREATE TABLE t_icj1 (a int);
| 			CREATE TABLE t_icj2 (b int);
| 			CREATE VIEW v_implicit_cross_join AS
| 			SELECT a, b FROM t_icj1, t_icj2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_114".v_implicit_cross_join'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_icj1.a,
|     t_icj2.b
|    FROM t_icj1,
|     t_icj2;
> exec
| DROP SCHEMA IF EXISTS "test_114" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_115" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_115"
< ok 0 -
> exec
| SET search_path TO "test_115", public;
|
| 			CREATE TABLE t_gl (a int, b int, c int);
| 			CREATE VIEW v_greatest_least AS
| 			SELECT GREATEST(a, b, c) AS mx, LEAST(a, b, c) AS mn FROM t_gl;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_115".v_greatest_least'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT GREATEST(a, b, c) AS mx,
|     LEAST(a, b, c) AS mn
|    FROM t_gl;
> exec
| DROP SCHEMA IF EXISTS "test_115" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_116" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_116"
< ok 0 -
> exec
| SET search_path TO "test_116", public;
|
| 			CREATE TABLE t_nif (a text);
| 			CREATE VIEW v_nullif_view AS
| 			SELECT NULLIF(a, '') FROM t_nif;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_116".v_nullif_view'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT NULLIF(a, ''::text) AS "nullif"
|    FROM t_nif;
> exec
| DROP SCHEMA IF EXISTS "test_116" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_117" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_117"
< ok 0 -
> exec
| SET search_path TO "test_117", public;
|
| 			CREATE TABLE t_cs (a int);
| 			CREATE VIEW v_case_searched AS
| 			SELECT CASE WHEN a > 0 THEN 'pos' WHEN a < 0 THEN 'neg' ELSE 'zero' END FROM t_cs;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_117".v_case_searched'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT
|         CASE
|             WHEN a > 0 THEN 'pos'::text
|             WHEN a < 0 THEN 'neg'::text
|             ELSE 'zero'::text
|         END AS "case"
|    FROM t_cs;
> exec
| DROP SCHEMA IF EXISTS "test_117" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_118" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_118"
< ok 0 -
> exec
| SET search_path TO "test_118", public;
|
| 			CREATE TABLE t_csm (status int);
| 			CREATE VIEW v_case_simple AS
| 			SELECT CASE status WHEN 1 THEN 'active' WHEN 2 THEN 'inactive' ELSE 'unknown' END FROM t_csm;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_118".v_case_simple'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT
|         CASE status
|             WHEN 1 THEN 'active'::text
|             WHEN 2 THEN 'inactive'::text
|             ELSE 'unknown'::text
|         END AS "case"
|    FROM t_csm;
> exec
| DROP SCHEMA IF EXISTS "test_118" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_119" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_119"
< ok 0 -
> exec
| SET search_path TO "test_119", public;
|
| 			CREATE TABLE t_jao (a text, b int);
| 			CREATE VIEW v_json_agg_order AS
| 			SELECT json_agg(a ORDER BY b) FROM t_jao;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_119".v_json_agg_order'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT json_agg(a ORDER BY b) AS json_agg
|    FROM t_jao;
> exec
| DROP SCHEMA IF EXISTS "test_119" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_120" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_120"
< ok 0 -
> exec
| SET search_path TO "test_120", public;
|
| 			CREATE TABLE t_joa (k text, v text);
| 			CREATE VIEW v_jsonb_object_agg AS
| 			SELECT jsonb_object_agg(k, v) FROM t_joa;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_120".v_jsonb_object_agg'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT jsonb_object_agg(k, v) AS jsonb_object_agg
|    FROM t_joa;
> exec
| DROP SCHEMA IF EXISTS "test_120" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_121" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_121"
< ok 0 -
> exec
| SET search_path TO "test_121", public;
|
| 			CREATE TABLE t_sys (a int);
| 			CREATE VIEW v_system_columns AS
| 			SELECT tableoid, ctid, a FROM t_sys;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_121".v_system_columns'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT tableoid,
|     ctid,
|     a
|    FROM t_sys;
> exec
| DROP SCHEMA IF EXISTS "test_121" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_122" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_122"
< ok 0 -
> exec
| SET search_path TO "test_122", public;
|
| 			CREATE TABLE orders (id serial, customer_id int, amount numeric, created_at timestamptz, status text);
| 			CREATE VIEW v_complex_analytics AS
| 			SELECT
| 				customer_id,
| 				count(*) AS order_count,
| 				sum(amount) AS total_amount,
| 				avg(amount)::numeric(10,2) AS avg_amount,
| 				min(created_at) AS first_order,
| 				max(created_at) AS last_order,
| 				count(*) FILTER (WHERE status = 'completed') AS completed_count,
| 				sum(amount) FILTER (WHERE created_at > now() - interval '30 days') AS recent_total,
| 				array_agg(DISTINCT status ORDER BY status) AS statuses,
| 				row_number() OVER (ORDER BY sum(amount) DESC) AS rank
| 			FROM orders
| 			GROUP BY customer_id;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_122".v_complex_analytics'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT customer_id,
|     count(*) AS order_count,
|     sum(amount) AS total_amount,
|     avg(amount)::numeric(10,2) AS avg_amount,
|     min(created_at) AS first_order,
|     max(created_at) AS last_order,
|     count(*) FILTER (WHERE status = 'completed'::text) AS completed_count,
|     sum(amount) FILTER (WHERE created_at > (now() - '30 days'::interval)) AS recent_total,
|     array_agg(DISTINCT status ORDER BY status) AS statuses,
|     row_number() OVER (ORDER BY (sum(amount)) DESC) AS rank
|    FROM orders
|   GROUP BY customer_id;
> exec
| DROP SCHEMA IF EXISTS "test_122" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_123" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_123"
< ok 0 -
> exec
| SET search_path TO "test_123", public;
|
| 			CREATE TABLE sales (id serial, region text, product text, amount numeric, sale_date date);
| 			CREATE VIEW v_complex_reporting AS
| 			WITH daily AS (
| 				SELECT region, product, sale_date, sum(amount) AS daily_total
| 				FROM sales GROUP BY region, product, sale_date
| 			),
| 			ranked AS (
| 				SELECT *,
| 					rank() OVER (PARTITION BY region ORDER BY daily_total DESC) AS rnk,
| 					sum(daily_total) OVER (PARTITION BY region ORDER BY sale_date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS rolling_7d
| 				FROM daily
| 			)
| 			SELECT region, product, sale_date, daily_total, rnk, rolling_7d,
| 				daily_total / NULLIF(sum(daily_total) OVER (PARTITION BY region, sale_date), 0) * 100 AS pct_of_day
| 			FROM ranked;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_123".v_complex_reporting'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH daily AS (
|          SELECT sales.region,
|             sales.product,
|             sales.sale_date,
|             sum(sales.amount) AS daily_total
|            FROM sales
|           GROUP BY sales.region, sales.product, sales.sale_date
|         ), ranked AS (
|          SELECT daily.region,
|             daily.product,
|             daily.sale_date,
|             daily.daily_total,
|             rank() OVER (PARTITION BY daily.region ORDER BY daily.daily_total DESC) AS rnk,
|             sum(daily.daily_total) OVER (PARTITION BY daily.region ORDER BY daily.sale_date ROWS BETWEEN 6 PRECEDING AND CURRENT ROW) AS rolling_7d
|            FROM daily
|         )
|  SELECT region,
|     product,
|     sale_date,
|     daily_total,
|     rnk,
|     rolling_7d,
|     daily_total / NULLIF(sum(daily_total) OVER (PARTITION BY region, sale_date), 0::numeric) * 100::numeric AS pct_of_day
|    FROM ranked;
> exec
| DROP SCHEMA IF EXISTS "test_123" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_124" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_124"
< ok 0 -
> exec
| SET search_path TO "test_124", public;
|
| 			CREATE TABLE t_sfa (a int, b text);
| 			CREATE VIEW v_subquery_from_alias AS
| 			SELECT sub.x, sub.y FROM (SELECT a AS x, b AS y FROM t_sfa) sub;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_124".v_subquery_from_alias'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT x,
|     y
|    FROM ( SELECT t_sfa.a AS x,
|             t_sfa.b AS y
|            FROM t_sfa) sub;
> exec
| DROP SCHEMA IF EXISTS "test_124" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_125" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_125"
< ok 0 -
> exec
| SET search_path TO "test_125", public;
|
| 			CREATE TABLE t_ml (id int, arr1 int[], arr2 text[]);
| 			CREATE VIEW v_multi_lateral AS
| 			SELECT t.id, u1.v1, u2.v2
| 			FROM t_ml t,
| 				LATERAL unnest(t.arr1) AS u1(v1),
| 				LATERAL unnest(t.arr2) AS u2(v2);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_125".v_multi_lateral'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.id,
|     u1.v1,
|     u2.v2
|    FROM t_ml t,
|     LATERAL unnest(t.arr1) u1(v1),
|     LATERAL unnest(t.arr2) u2(v2);
> exec
| DROP SCHEMA IF EXISTS "test_125" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_126" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_126"
< ok 0 -
> exec
| SET search_path TO "test_126", public;
|
| 			CREATE TABLE t_cjl (n int);
| 			CREATE VIEW v_cross_join_lateral AS
| 			SELECT t.n, gs FROM t_cjl t CROSS JOIN LATERAL generate_series(1, t.n) AS gs;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_126".v_cross_join_lateral'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t.n,
|     gs.gs
|    FROM t_cjl t
|      CROSS JOIN LATERAL generate_series(1, t.n) gs(gs);
> exec
| DROP SCHEMA IF EXISTS "test_126" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_127" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_127"
< ok 0 -
> exec
| SET search_path TO "test_127", public;
|
| 			CREATE TABLE t_ncte (id int, parent_id int, name text);
| 			CREATE VIEW v_nested_cte AS
| 			WITH RECURSIVE tree AS (
| 				SELECT id, parent_id, name, 1 AS depth FROM t_ncte WHERE parent_id IS NULL
| 				UNION ALL
| 				SELECT c.id, c.parent_id, c.name, p.depth + 1
| 				FROM t_ncte c JOIN tree p ON c.parent_id = p.id
| 			)
| 			SELECT * FROM tree;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_127".v_nested_cte'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  WITH RECURSIVE tree AS (
|          SELECT t_ncte.id,
|             t_ncte.parent_id,
|             t_ncte.name,
|             1 AS depth
|            FROM t_ncte
|           WHERE t_ncte.parent_id IS NULL
|         UNION ALL
|          SELECT c.id,
|             c.parent_id,
|             c.name,
|             p.depth + 1
|            FROM t_ncte c
|              JOIN tree p ON c.parent_id = p.id
|         )
|  SELECT id,
|     parent_id,
|     name,
|     depth
|    FROM tree;
> exec
| DROP SCHEMA IF EXISTS "test_127" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_128" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_128"
< ok 0 -
> exec
| SET search_path TO "test_128", public;
|
| 			CREATE TABLE t_jum1 (a int, b int, c text);
| 			CREATE TABLE t_jum2 (a int, b int, d text);
| 			CREATE VIEW v_join_using_multi AS
| 			SELECT * FROM t_jum1 JOIN t_jum2 USING (a, b);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_128".v_join_using_multi'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_jum1.a,
|     t_jum1.b,
|     t_jum1.c,
|     t_jum2.d
|    FROM t_jum1
|      JOIN t_jum2 USING (a, b);
> exec
| DROP SCHEMA IF EXISTS "test_128" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_129" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_129"
< ok 0 -
> exec
| SET search_path TO "test_129", public;
|
| 			CREATE TABLE t_sj (id int, parent_id int, name text);
| 			CREATE VIEW v_self_join AS
| 			SELECT c.name AS child, p.name AS parent
| 			FROM t_sj c LEFT JOIN t_sj p ON c.parent_id = p.id;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_129".v_self_join'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT c.name AS child,
|     p.name AS parent
|    FROM t_sj c
|      LEFT JOIN t_sj p ON c.parent_id = p.id;
> exec
| DROP SCHEMA IF EXISTS "test_129" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_130" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_130"
< ok 0 -
> exec
| SET search_path TO "test_130", public;
|
| 			CREATE TABLE t_csq1 (id int, name text);
| 			CREATE TABLE t_csq2 (id int, parent_id int, val int);
| 			CREATE VIEW v_correlated_subquery AS
| 			SELECT t1.name,
| 				(SELECT sum(t2.val) FROM t_csq2 t2 WHERE t2.parent_id = t1.id) AS total
| 			FROM t_csq1 t1;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_130".v_correlated_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT name,
|     ( SELECT sum(t2.val) AS sum
|            FROM t_csq2 t2
|           WHERE t2.parent_id = t1.id) AS total
|    FROM t_csq1 t1;
> exec
| DROP SCHEMA IF EXISTS "test_130" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_131" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_131"
< ok 0 -
> exec
| SET search_path TO "test_131", public;
|
| 			CREATE TABLE t_asqc (id int, tag text);
| 			CREATE VIEW v_array_subquery_cast AS
| 			SELECT id, ARRAY(SELECT tag FROM t_asqc t2 WHERE t2.id = t.id)::text[] AS tags
| 			FROM (SELECT DISTINCT id FROM t_asqc) t;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_131".v_array_subquery_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     ARRAY( SELECT t2.tag
|            FROM t_asqc t2
|           WHERE t2.id = t.id) AS tags
|    FROM ( SELECT DISTINCT t_asqc.id
|            FROM t_asqc) t;
> exec
| DROP SCHEMA IF EXISTS "test_131" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_132" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_132"
< ok 0 -
> exec
| SET search_path TO "test_132", public;
|
| 			CREATE VIEW v_multi_srf AS
| 			SELECT * FROM generate_series(1, 5) AS a, generate_series(1, 3) AS b;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_132".v_multi_srf'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a.a,
|     b.b
|    FROM generate_series(1, 5) a(a),
|     generate_series(1, 3) b(b);
> exec
| DROP SCHEMA IF EXISTS "test_132" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_133" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_133"
< ok 0 -
> exec
| SET search_path TO "test_133", public;
|
| 			CREATE TABLE t_amo (grp text, a text, b int);
| 			CREATE VIEW v_agg_multi_order AS
| 			SELECT grp, string_agg(a, ',' ORDER BY b DESC, a ASC) FROM t_amo GROUP BY grp;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_133".v_agg_multi_order'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT grp,
|     string_agg(a, ','::text ORDER BY b DESC, a) AS string_agg
|    FROM t_amo
|   GROUP BY grp;
> exec
| DROP SCHEMA IF EXISTS "test_133" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_134" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_134"
< ok 0 -
> exec
| SET search_path TO "test_134", public;
|
| 			CREATE TABLE t_cct (a int, b bigint, c numeric);
| 			CREATE VIEW v_coalesce_chain_types AS
| 			SELECT COALESCE(a::numeric, b::numeric, c, 0) AS val FROM t_cct;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_134".v_coalesce_chain_types'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT COALESCE(a::numeric, b::numeric, c, 0::numeric) AS val
|    FROM t_cct;
> exec
| DROP SCHEMA IF EXISTS "test_134" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_135" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_135"
< ok 0 -
> exec
| SET search_path TO "test_135", public;
|
| 			CREATE TABLE t_cws (id int, status text);
| 			CREATE VIEW v_case_with_subquery AS
| 			SELECT id,
| 				CASE
| 					WHEN EXISTS (SELECT 1 FROM t_cws t2 WHERE t2.id = t.id AND t2.status = 'active') THEN 'has_active'
| 					ELSE 'no_active'
| 				END AS label
| 			FROM (SELECT DISTINCT id FROM t_cws) t;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_135".v_case_with_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|         CASE
|             WHEN (EXISTS ( SELECT 1
|                FROM t_cws t2
|               WHERE t2.id = t.id AND t2.status = 'active'::text)) THEN 'has_active'::text
|             ELSE 'no_active'::text
|         END AS label
|    FROM ( SELECT DISTINCT t_cws.id
|            FROM t_cws) t;
> exec
| DROP SCHEMA IF EXISTS "test_135" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_136" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_136"
< ok 0 -
> exec
| SET search_path TO "test_136", public;
|
| 			CREATE TABLE t_bi (a bool, b bool);
| 			CREATE VIEW v_boolean_is AS
| 			SELECT * FROM t_bi WHERE a IS TRUE AND b IS NOT FALSE AND (a IS UNKNOWN) = false;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_136".v_boolean_is'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_bi
|   WHERE a IS TRUE AND b IS NOT FALSE AND (a IS UNKNOWN) = false;
> exec
| DROP SCHEMA IF EXISTS "test_136" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_137" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_137"
< ok 0 -
> exec
| SET search_path TO "test_137", public;
|
| 			CREATE TABLE t_iv (id int);
| 			CREATE VIEW v_in_values AS
| 			SELECT * FROM t_iv WHERE id IN (1, 2, 3, 4, 5);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_137".v_in_values'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id
|    FROM t_iv
|   WHERE id = ANY (ARRAY[1, 2, 3, 4, 5]);
> exec
| DROP SCHEMA IF EXISTS "test_137" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_138" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_138"
< ok 0 -
> exec
| SET search_path TO "test_138", public;
|
| 			CREATE TABLE t_nfc (ts timestamptz);
| 			CREATE VIEW v_nested_func_cast AS
| 			SELECT
| 				date_trunc('day', ts AT TIME ZONE 'UTC')::date AS day,
| 				extract(epoch FROM ts)::bigint AS epoch_s,
| 				to_char(ts, 'YYYY-MM-DD"T"HH24:MI:SS') AS iso
| 			FROM t_nfc;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_138".v_nested_func_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT date_trunc('day'::text, (ts AT TIME ZONE 'UTC'::text))::date AS day,
|     EXTRACT(epoch FROM ts)::bigint AS epoch_s,
|     to_char(ts, 'YYYY-MM-DD"T"HH24:MI:SS'::text) AS iso
|    FROM t_nfc;
> exec
| DROP SCHEMA IF EXISTS "test_138" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_139" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_139"
< ok 0 -
> exec
| SET search_path TO "test_139", public;
|
| 			CREATE TABLE t_so (a text, b text);
| 			CREATE VIEW v_string_operators AS
| 			SELECT
| 				a || ' ' || b AS full_name,
| 				a LIKE 'test%' AS starts_with_test,
| 				a ILIKE '%TEST%' AS contains_test,
| 				a ~ '^[A-Z]' AS starts_upper,
| 				a ~* 'pattern' AS imatches
| 			FROM t_so;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_139".v_string_operators'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (a || ' '::text) || b AS full_name,
|     a ~~ 'test%'::text AS starts_with_test,
|     a ~~* '%TEST%'::text AS contains_test,
|     a ~ '^[A-Z]'::text AS starts_upper,
|     a ~* 'pattern'::text AS imatches
|    FROM t_so;
> exec
| DROP SCHEMA IF EXISTS "test_139" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_140" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_140"
< ok 0 -
> exec
| SET search_path TO "test_140", public;
|
| 			CREATE TABLE t_no (a int);
| 			CREATE VIEW v_nulls_ordering AS
| 			SELECT a FROM t_no ORDER BY a NULLS FIRST;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_140".v_nulls_ordering'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_no
|   ORDER BY a NULLS FIRST;
> exec
| DROP SCHEMA IF EXISTS "test_140" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_141" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_141"
< ok 0 -
> exec
| SET search_path TO "test_141", public;
|
| 			CREATE TABLE t_aas (val int);
| 			CREATE VIEW v_all_any_subquery AS
| 			SELECT * FROM t_aas WHERE val > ALL(SELECT val FROM t_aas WHERE val < 10);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_141".v_all_any_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT val
|    FROM t_aas
|   WHERE (val > ALL ( SELECT t_aas_1.val
|            FROM t_aas t_aas_1
|           WHERE t_aas_1.val < 10));
> exec
| DROP SCHEMA IF EXISTS "test_141" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_142" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_142"
< ok 0 -
> exec
| SET search_path TO "test_142", public;
|
| 			CREATE TABLE t_cw (grp text, ts timestamp, val numeric);
| 			CREATE VIEW v_complex_window AS
| 			SELECT grp, ts, val,
| 				lag(val) OVER w AS prev_val,
| 				lead(val) OVER w AS next_val,
| 				sum(val) OVER (PARTITION BY grp ORDER BY ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
| 				ntile(4) OVER (PARTITION BY grp ORDER BY val) AS quartile,
| 				percent_rank() OVER (PARTITION BY grp ORDER BY val) AS pct_rank,
| 				cume_dist() OVER (PARTITION BY grp ORDER BY val) AS cume
| 			FROM t_cw
| 			WINDOW w AS (PARTITION BY grp ORDER BY ts);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_142".v_complex_window'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT grp,
|     ts,
|     val,
|     lag(val) OVER w AS prev_val,
|     lead(val) OVER w AS next_val,
|     sum(val) OVER (PARTITION BY grp ORDER BY ts ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_total,
|     ntile(4) OVER (PARTITION BY grp ORDER BY val) AS quartile,
|     percent_rank() OVER (PARTITION BY grp ORDER BY val) AS pct_rank,
|     cume_dist() OVER (PARTITION BY grp ORDER BY val) AS cume
|    FROM t_cw
|   WINDOW w AS (PARTITION BY grp ORDER BY ts);
> exec
| DROP SCHEMA IF EXISTS "test_142" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_143" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_143"
< ok 0 -
> exec
| SET search_path TO "test_143", public;
|
| 			CREATE TABLE t_afc2 (category text, status text, amount numeric);
| 			CREATE VIEW v_agg_filter_complex AS
| 			SELECT category,
| 				count(*) FILTER (WHERE status IN ('active', 'pending')) AS active_count,
| 				sum(amount) FILTER (WHERE status = 'completed' AND amount > 100) AS big_completed,
| 				avg(amount) FILTER (WHERE amount BETWEEN 10 AND 1000) AS mid_avg
| 			FROM t_afc2
| 			GROUP BY category;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_143".v_agg_filter_complex'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT category,
|     count(*) FILTER (WHERE status = ANY (ARRAY['active'::text, 'pending'::text])) AS active_count,
|     sum(amount) FILTER (WHERE status = 'completed'::text AND amount > 100::numeric) AS big_completed,
|     avg(amount) FILTER (WHERE amount >= 10::numeric AND amount <= 1000::numeric) AS mid_avg
|    FROM t_afc2
|   GROUP BY category;
> exec
| DROP SCHEMA IF EXISTS "test_143" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_144" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_144"
< ok 0 -
> exec
| SET search_path TO "test_144", public;
|
| 			CREATE VIEW v_with_ordinality AS
| 			SELECT val, ord FROM unnest(ARRAY['a','b','c']) WITH ORDINALITY AS t(val, ord);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_144".v_with_ordinality'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT val,
|     ord
|    FROM unnest(ARRAY['a'::text, 'b'::text, 'c'::text]) WITH ORDINALITY t(val, ord);
> exec
| DROP SCHEMA IF EXISTS "test_144" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_145" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_145"
< ok 0 -
> exec
| SET search_path TO "test_145", public;
|
| 			CREATE TABLE t_ea1 (a int);
| 			CREATE TABLE t_ea2 (a int);
| 			CREATE VIEW v_except_all AS
| 			SELECT a FROM t_ea1 EXCEPT ALL SELECT a FROM t_ea2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_145".v_except_all'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_ea1.a
|    FROM t_ea1
| EXCEPT ALL
|  SELECT t_ea2.a
|    FROM t_ea2;
> exec
| DROP SCHEMA IF EXISTS "test_145" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_146" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_146"
< ok 0 -
> exec
| SET search_path TO "test_146", public;
|
| 			CREATE TABLE t_muc1 (a int);
| 			CREATE TABLE t_muc2 (a bigint);
| 			CREATE TABLE t_muc3 (a numeric);
| 			CREATE VIEW v_multi_union_cast AS
| 			SELECT a FROM t_muc1 UNION ALL SELECT a FROM t_muc2 UNION ALL SELECT a FROM t_muc3;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_146".v_multi_union_cast'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT t_muc1.a
|    FROM t_muc1
| UNION ALL
|  SELECT t_muc2.a
|    FROM t_muc2
| UNION ALL
|  SELECT t_muc3.a
|    FROM t_muc3;
> exec
| DROP SCHEMA IF EXISTS "test_146" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_147" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_147"
< ok 0 -
> exec
| SET search_path TO "test_147", public;
|
| 			CREATE TABLE t_hc (grp text, val int);
| 			CREATE VIEW v_having_clause AS
| 			SELECT grp, count(*), sum(val)
| 			FROM t_hc
| 			GROUP BY grp
| 			HAVING count(*) > 1 AND sum(val) > 100;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_147".v_having_clause'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT grp,
|     count(*) AS count,
|     sum(val) AS sum
|    FROM t_hc
|   GROUP BY grp
|  HAVING count(*) > 1 AND sum(val) > 100;
> exec
| DROP SCHEMA IF EXISTS "test_147" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_148" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_148"
< ok 0 -
> exec
| SET search_path TO "test_148", public;
|
| 			CREATE TABLE t_cd (
| 				id bigint GENERATED ALWAYS AS IDENTITY,
| 				uid uuid DEFAULT gen_random_uuid(),
| 				created_at timestamptz DEFAULT now(),
| 				updated_at timestamptz DEFAULT CURRENT_TIMESTAMP,
| 				data jsonb DEFAULT '{}'::jsonb,
| 				tags text[] DEFAULT '{}'::text[],
| 				status text DEFAULT 'pending'::text,
| 				counter int DEFAULT 0
| 			);
| 			CREATE VIEW v_complex_defaults AS SELECT * FROM t_cd;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_148".v_complex_defaults'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     uid,
|     created_at,
|     updated_at,
|     data,
|     tags,
|     status,
|     counter
|    FROM t_cd;
> exec
| DROP SCHEMA IF EXISTS "test_148" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_149" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_149"
< ok 0 -
> exec
| SET search_path TO "test_149", public;
|
| 			CREATE TABLE t_rcmp (a int, b int);
| 			CREATE VIEW v_row_comparison AS
| 			SELECT * FROM t_rcmp WHERE ROW(a, b) > ROW(1, 2);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_149".v_row_comparison'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM t_rcmp
|   WHERE (ROW(a, b) > ROW(1, 2));
> exec
| DROP SCHEMA IF EXISTS "test_149" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_150" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_150"
< ok 0 -
> exec
| SET search_path TO "test_150", public;
|
| 			CREATE TYPE inner_type AS (x int, y int);
| 			CREATE TYPE outer_type AS (name text, coords inner_type);
| 			CREATE TABLE t_cn (data outer_type);
| 			CREATE VIEW v_composite_nested AS
| 			SELECT (data).name, ((data).coords).x, ((data).coords).y FROM t_cn;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_150".v_composite_nested'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (data).name AS name,
|     (data).coords.x AS x,
|     (data).coords.y AS y
|    FROM t_cn;
> exec
| DROP SCHEMA IF EXISTS "test_150" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_151" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_151"
< ok 0 -
> exec
| SET search_path TO "test_151", public;
|
| 			CREATE TYPE priority AS ENUM ('low', 'medium', 'high', 'critical');
| 			CREATE TABLE t_eo (id int, p priority);
| 			CREATE VIEW v_enum_ordering AS
| 			SELECT * FROM t_eo WHERE p >= 'high' ORDER BY p;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_151".v_enum_ordering'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     p
|    FROM t_eo
|   WHERE p >= 'high'::priority
|   ORDER BY p;
> exec
| DROP SCHEMA IF EXISTS "test_151" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_152" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_152"
< ok 0 -
> exec
| SET search_path TO "test_152", public;
|
| 			CREATE TABLE t_parent (id int, name text);
| 			CREATE TABLE t_child (extra text) INHERITS (t_parent);
| 			CREATE VIEW v_inheritance_view AS
| 			SELECT * FROM ONLY t_parent;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_152".v_inheritance_view'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT id,
|     name
|    FROM ONLY t_parent;
> exec
| DROP SCHEMA IF EXISTS "test_152" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_153" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_153"
< ok 0 -
> exec
| SET search_path TO "test_153", public;
|
| 			CREATE FUNCTION debug_vd_concat(VARIADIC text[]) RETURNS text AS $$
| 				SELECT array_to_string($1, ',')
| 			$$ LANGUAGE SQL;
| 			CREATE TABLE t_vf (a text, b text, c text);
| 			CREATE VIEW v_variadic_func AS
| 			SELECT debug_vd_concat(a, b, c) FROM t_vf;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_153".v_variadic_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT debug_vd_concat(VARIADIC ARRAY[a, b, c]) AS debug_vd_concat
|    FROM t_vf;
> exec
| DROP SCHEMA IF EXISTS "test_153" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_154" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_154"
< ok 0 -
> exec
| SET search_path TO "test_154", public;
|
| 			CREATE VIEW v_rows_from AS
| 			SELECT * FROM ROWS FROM (
| 				generate_series(1, 3),
| 				generate_series(1, 5)
| 			) AS t(a, b);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_154".v_rows_from'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a,
|     b
|    FROM ROWS FROM(generate_series(1, 3), generate_series(1, 5)) t(a, b);
> exec
| DROP SCHEMA IF EXISTS "test_154" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_155" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_155"
< ok 0 -
> exec
| SET search_path TO "test_155", public;
|
| 			CREATE EXTENSION IF NOT EXISTS hstore;
| 			CREATE TABLE t_hs (data hstore);
| 			CREATE VIEW v_hstore_ext AS
| 			SELECT data -> 'key' AS val, data ? 'key' AS has_key FROM t_hs;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_155".v_hstore_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT data -> 'key'::text AS val,
|     data ? 'key'::text AS has_key
|    FROM t_hs;
> exec
| DROP SCHEMA IF EXISTS "test_155" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_156" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_156"
< ok 0 -
> exec
| SET search_path TO "test_156", public;
|
| 			CREATE EXTENSION IF NOT EXISTS ltree;
| 			CREATE TABLE t_lt (path ltree);
| 			CREATE VIEW v_ltree_ext AS
| 			SELECT path, path <@ 'root.a'::ltree AS is_child FROM t_lt;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_156".v_ltree_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT path,
|     path <@ 'root.a'::ltree AS is_child
|    FROM t_lt;
> exec
| DROP SCHEMA IF EXISTS "test_156" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_157" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_157"
< ok 0 -
> exec
| SET search_path TO "test_157", public;
|
| 			CREATE EXTENSION IF NOT EXISTS citext;
| 			CREATE TABLE t_ci (name citext);
| 			CREATE VIEW v_citext_ext AS
| 			SELECT name FROM t_ci WHERE name = 'Test';
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_157".v_citext_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT name
|    FROM t_ci
|   WHERE name = 'Test'::citext;
> exec
| DROP SCHEMA IF EXISTS "test_157" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_158" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_158"
< ok 0 -
> exec
| SET search_path TO "test_158", public;
|
| 			CREATE EXTENSION IF NOT EXISTS intarray;
| 			CREATE TABLE t_ia2 (arr int[]);
| 			CREATE VIEW v_intarray_ext AS
| 			SELECT arr, sort(arr) AS sorted FROM t_ia2;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_158".v_intarray_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT arr,
|     sort(arr) AS sorted
|    FROM t_ia2;
> exec
| DROP SCHEMA IF EXISTS "test_158" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_159" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_159"
< ok 0 -
> exec
| SET search_path TO "test_159", public;
|
| 			CREATE EXTENSION IF NOT EXISTS pg_trgm;
| 			CREATE TABLE t_trgm (name text);
| 			CREATE VIEW v_trgm_ext AS
| 			SELECT name, similarity(name, 'test') AS sim FROM t_trgm;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_159".v_trgm_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT name,
|     similarity(name, 'test'::text) AS sim
|    FROM t_trgm;
> exec
| DROP SCHEMA IF EXISTS "test_159" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_160" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_160"
< ok 0 -
> exec
| SET search_path TO "test_160", public;
|
| 			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
| 			CREATE VIEW v_uuid_ossp_ext AS
| 			SELECT uuid_generate_v4() AS id;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_160".v_uuid_ossp_ext'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT uuid_generate_v4() AS id;
> exec
| DROP SCHEMA IF EXISTS "test_160" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_161" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_161"
< ok 0 -
> exec
| SET search_path TO "test_161", public;
|
| 			CREATE TABLE t_hsq (grp int, val int);
| 			CREATE VIEW v_having_subquery AS
| 			SELECT grp, avg(val) FROM t_hsq
| 			GROUP BY grp
| 			HAVING avg(val) > (SELECT avg(val) FROM t_hsq);
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_161".v_having_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT grp,
|     avg(val) AS avg
|    FROM t_hsq
|   GROUP BY grp
|  HAVING avg(val) > (( SELECT avg(t_hsq_1.val) AS avg
|            FROM t_hsq t_hsq_1));
> exec
| DROP SCHEMA IF EXISTS "test_161" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_162" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_162"
< ok 0 -
> exec
| SET search_path TO "test_162", public;
|
| 			CREATE TABLE t_das (id int, val int);
| 			CREATE VIEW v_duplicate_alias_subquery AS
| 			SELECT a.id, b.total
| 			FROM t_das a
| 			JOIN (SELECT id, sum(val) AS total FROM t_das GROUP BY id) b ON a.id = b.id;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_162".v_duplicate_alias_subquery'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a.id,
|     b.total
|    FROM t_das a
|      JOIN ( SELECT t_das.id,
|             sum(t_das.val) AS total
|            FROM t_das
|           GROUP BY t_das.id) b ON a.id = b.id;
> exec
| DROP SCHEMA IF EXISTS "test_162" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_163" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_163"
< ok 0 -
> exec
| SET search_path TO "test_163", public;
|
| 			CREATE TABLE t_le (pattern text);
| 			CREATE VIEW v_like_escape AS
| 			SELECT * FROM t_le WHERE pattern LIKE '%!%%' ESCAPE '!';
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_163".v_like_escape'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT pattern
|    FROM t_le
|   WHERE pattern ~~ like_escape('%!%%'::text, '!'::text);
> exec
| DROP SCHEMA IF EXISTS "test_163" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_164" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_164"
< ok 0 -
> exec
| SET search_path TO "test_164", public;
|
| 			CREATE TABLE t_ac (arr int[]);
| 			CREATE VIEW v_array_compare AS
| 			SELECT arr,
| 				arr @> ARRAY[1,2] AS contains,
| 				arr <@ ARRAY[1,2,3,4,5] AS contained_by,
| 				arr && ARRAY[3,4] AS overlaps_arr
| 			FROM t_ac;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_164".v_array_compare'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT arr,
|     arr @> ARRAY[1, 2] AS contains,
|     arr <@ ARRAY[1, 2, 3, 4, 5] AS contained_by,
|     arr && ARRAY[3, 4] AS overlaps_arr
|    FROM t_ac;
> exec
| DROP SCHEMA IF EXISTS "test_164" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_165" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_165"
< ok 0 -
> exec
| SET search_path TO "test_165", public;
|
| 			CREATE TABLE t_inn (a int);
| 			CREATE VIEW v_isnull_notnull AS
| 			SELECT * FROM t_inn WHERE a NOTNULL;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_165".v_isnull_notnull'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a
|    FROM t_inn
|   WHERE a IS NOT NULL;
> exec
| DROP SCHEMA IF EXISTS "test_165" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_166" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_166"
< ok 0 -
> exec
| SET search_path TO "test_166", public;
|
| 			CREATE TABLE t_ov (s text);
| 			CREATE VIEW v_overlay_func AS
| 			SELECT OVERLAY(s PLACING 'XX' FROM 3 FOR 2) FROM t_ov;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_166".v_overlay_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT OVERLAY(s PLACING 'XX'::text FROM 3 FOR 2) AS "overlay"
|    FROM t_ov;
> exec
| DROP SCHEMA IF EXISTS "test_166" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_167" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_167"
< ok 0 -
> exec
| SET search_path TO "test_167", public;
|
| 			CREATE TABLE t_pos (s text);
| 			CREATE VIEW v_position_func AS
| 			SELECT POSITION('needle' IN s) FROM t_pos;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_167".v_position_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT POSITION(('needle'::text) IN (s)) AS "position"
|    FROM t_pos;
> exec
| DROP SCHEMA IF EXISTS "test_167" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_168" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_168"
< ok 0 -
> exec
| SET search_path TO "test_168", public;
|
| 			CREATE TABLE t_sp (s text);
| 			CREATE VIEW v_substring_pattern AS
| 			SELECT SUBSTRING(s FROM 1 FOR 5) FROM t_sp;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_168".v_substring_pattern'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT SUBSTRING(s FROM 1 FOR 5) AS "substring"
|    FROM t_sp;
> exec
| DROP SCHEMA IF EXISTS "test_168" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_169" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_169"
< ok 0 -
> exec
| SET search_path TO "test_169", public;
|
| 			CREATE TABLE t_tv (s text);
| 			CREATE VIEW v_trim_variants AS
| 			SELECT
| 				TRIM(BOTH ' ' FROM s) AS trimmed,
| 				TRIM(LEADING '0' FROM s) AS ltrimmed,
| 				TRIM(TRAILING FROM s) AS rtrimmed
| 			FROM t_tv;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_169".v_trim_variants'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT TRIM(BOTH ' '::text FROM s) AS trimmed,
|     TRIM(LEADING '0'::text FROM s) AS ltrimmed,
|     TRIM(TRAILING FROM s) AS rtrimmed
|    FROM t_tv;
> exec
| DROP SCHEMA IF EXISTS "test_169" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_170" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_170"
< ok 0 -
> exec
| SET search_path TO "test_170", public;
|
| 			CREATE TABLE t_bo (a int, b int);
| 			CREATE VIEW v_bit_operations AS
| 			SELECT a & b AS band, a | b AS bor, a # b AS bxor, ~a AS bnot, a << 2 AS lshift, a >> 1 AS rshift FROM t_bo;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_170".v_bit_operations'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT a & b AS band,
|     a | b AS bor,
|     a # b AS bxor,
|     ~ a AS bnot,
|     a << 2 AS lshift,
|     a >> 1 AS rshift
|    FROM t_bo;
> exec
| DROP SCHEMA IF EXISTS "test_170" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_171" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_171"
< ok 0 -
> exec
| SET search_path TO "test_171", public;
|
| 			CREATE TABLE t_cvdc (a text, b int);
| 			CREATE VIEW v_cast_vs_double_colon AS
| 			SELECT CAST(a AS int), b::text, CAST(b AS numeric(10,2)), a::varchar(100) FROM t_cvdc;
< error
| 42701 column "a" specified more than once
> exec
| DROP SCHEMA IF EXISTS "test_171" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_172" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_172"
< ok 0 -
> exec
| SET search_path TO "test_172", public;
|
| 			CREATE TABLE t_norm (s text);
| 			CREATE VIEW v_normalize_func AS
| 			SELECT NORMALIZE(s, NFC) FROM t_norm;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_172".v_normalize_func'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT NORMALIZE(s, NFC) AS "normalize"
|    FROM t_norm;
> exec
| DROP SCHEMA IF EXISTS "test_172" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_173" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_173"
< ok 0 -
> exec
| SET search_path TO "test_173", public;
|
| 			CREATE TABLE t_isn (s text);
| 			CREATE VIEW v_is_normalized AS
| 			SELECT s IS NFC NORMALIZED AS is_nfc FROM t_isn;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_173".v_is_normalized'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT (s IS NFC NORMALIZED) AS is_nfc
|    FROM t_isn;
> exec
| DROP SCHEMA IF EXISTS "test_173" CASCADE
< ok 0 -
> exec
| DROP SCHEMA IF EXISTS "test_174" CASCADE
< ok 0 -
> exec
| CREATE SCHEMA "test_174"
< ok 0 -
> exec
| SET search_path TO "test_174", public;
|
| 			CREATE TABLE t_cia (name text);
| 			CREATE VIEW v_collation_in_agg AS
| 			SELECT min(name COLLATE "C") AS first_c FROM t_cia;
< ok 0 -
> query
| SELECT pg_get_viewdef('"test_174".v_collation_in_agg'::regclass, true)
< result
< column "pg_get_viewdef" "TEXT"
< row
: string
|  SELECT min(name COLLATE "C") AS first_c
|    FROM t_cia;
> exec
| DROP SCHEMA IF EXISTS "test_174" CASCADE
< ok 0 -
