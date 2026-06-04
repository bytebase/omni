package parser

import (
	"testing"
)

// This file is the parser-ddl node's differential-oracle gate
// (correctness-protocol.md): for every statement in the corpus below, omni's
// Parse accept/reject verdict must equal Trino 481's verdict as reported by the
// live oracle (trinooracle.CheckSyntax, which classifies a SYNTAX_ERROR as
// reject and every other outcome — success or a semantic error like
// TABLE_NOT_FOUND / SCHEMA_NOT_FOUND / NOT_SUPPORTED / TYPE_NOT_FOUND — as
// accept).
//
// The corpus is the oracle's to adjudicate: each case runs through both omni
// and Trino and the two verdicts are compared. It covers the legacy grammar
// examples (truth2, examples/{create_table,create_table_as_select,create_schema,
// rename_schema,create_view,rename_view,alter_view_set_authorization,
// create_materialized_view,rename_materialized_view,refresh_materialized_view,
// drop_*,drop_column,alter_table_alter_column_set_data_type,
// alter_table_set_authorization,comment_*,analyze}.sql with assertions
// stripped), every documented form (truth1 ddl.md), the oracle-pinned grammar
// facts (D-CT*/D-DR*/D-SC*/D-VW*/D-MV*/D-AT*/D-CAT*/D-CM*/D-AN*), and a set of
// negative forms the engine must reject.
//
// Out of this node's scope (dispatched/owned elsewhere, NOT in this corpus):
// CREATE/DROP FUNCTION (parser-routines), CREATE/DROP ROLE (parser-dcl-tcl),
// TRUNCATE TABLE (parser-dml). Those have their own gates.
//
// Skipped cleanly when no Trino oracle is reachable.

// ddlOracleCorpus is the full DDL statement corpus for the differential gate.
var ddlOracleCorpus = []string{
	// =====================================================================
	// CREATE SCHEMA  (truth1 create-schema + truth2 create_schema.sql)
	// =====================================================================
	"CREATE SCHEMA web",
	"CREATE SCHEMA hive.sales",
	"CREATE SCHEMA IF NOT EXISTS traffic",
	"CREATE SCHEMA web AUTHORIZATION alice",
	"CREATE SCHEMA web AUTHORIZATION USER alice",
	"CREATE SCHEMA web AUTHORIZATION ROLE PUBLIC",
	"CREATE SCHEMA web AUTHORIZATION alice WITH ( LOCATION = '/hive/data/web' )",
	"CREATE SCHEMA web AUTHORIZATION ROLE PUBLIC WITH ( LOCATION = '/hive/data/web' )",
	"CREATE SCHEMA test",
	"CREATE SCHEMA IF NOT EXISTS test",
	"CREATE SCHEMA test WITH (a = 'apple', b = 123)",
	`CREATE SCHEMA "some name that contains space"`,
	// negatives
	"CREATE SCHEMA",                                      // missing name
	"CREATE SCHEMA a.b.c",                                // D-SC1: 3-part schema name
	"CREATE SCHEMA web WITH ()",                          // empty property list
	"CREATE SCHEMA web AUTHORIZATION",                    // dangling AUTHORIZATION
	"CREATE SCHEMA s WITH (a = '1') AUTHORIZATION alice", // clause order: AUTHORIZATION before WITH

	// =====================================================================
	// DROP SCHEMA  (truth1 drop-schema + truth2 drop_schema.sql)
	// =====================================================================
	"DROP SCHEMA web",
	"DROP SCHEMA IF EXISTS sales",
	"DROP SCHEMA archive CASCADE",
	"DROP SCHEMA archive RESTRICT",
	"DROP SCHEMA test",
	"DROP SCHEMA test CASCADE",
	"DROP SCHEMA IF EXISTS test",
	"DROP SCHEMA IF EXISTS test RESTRICT",
	`DROP SCHEMA "some schema that contains space"`,
	"DROP SCHEMA hive.sales", // 2-part
	// negatives
	"DROP SCHEMA",          // missing name
	"DROP SCHEMA a.b.c",    // 3-part schema name
	"DROP SCHEMA a FOObar", // trailing garbage

	// =====================================================================
	// ALTER SCHEMA  (truth1 alter-schema + truth2 rename_schema.sql)
	// =====================================================================
	"ALTER SCHEMA web RENAME TO traffic",
	"ALTER SCHEMA web SET AUTHORIZATION alice",
	"ALTER SCHEMA web SET AUTHORIZATION USER alice",
	"ALTER SCHEMA web SET AUTHORIZATION ROLE PUBLIC",
	"ALTER SCHEMA foo RENAME TO bar",
	"ALTER SCHEMA foo.bar RENAME TO baz",
	`ALTER SCHEMA "awesome schema"."awesome table" RENAME TO "even more awesome table"`,
	// negatives
	"ALTER SCHEMA web RENAME",               // dangling RENAME
	"ALTER SCHEMA web RENAME TO a.b",        // rename target must be a single identifier
	"ALTER SCHEMA web SET",                  // dangling SET
	"ALTER SCHEMA web SET PROPERTIES x = 1", // schema has no SET PROPERTIES

	// =====================================================================
	// CREATE TABLE — column-definition form  (truth1 create-table + truth2)
	// =====================================================================
	"CREATE TABLE orders (orderkey bigint, orderstatus varchar, totalprice double, orderdate date) WITH (format = 'ORC')",
	"CREATE TABLE IF NOT EXISTS orders (orderkey bigint, orderstatus varchar, totalprice double COMMENT 'Price in cents.', orderdate date, status varchar DEFAULT 'created') COMMENT 'A table to keep track of orders.'",
	"CREATE TABLE bigger_orders (another_orderkey bigint, LIKE orders, another_orderdate date)",
	"CREATE TABLE IF NOT EXISTS bar (LIKE like_table)",
	"CREATE TABLE IF NOT EXISTS bar (LIKE like_table INCLUDING PROPERTIES)",
	"CREATE TABLE t (a bigint, b varchar(10), c decimal(10, 2))",
	"CREATE TABLE t (a bigint NOT NULL)",
	"CREATE TABLE t (a bigint DEFAULT 5 NOT NULL)",             // D-CT2: DEFAULT then NOT NULL
	"CREATE TABLE t (a bigint NOT NULL DEFAULT 5)",             // D-CT2: NOT NULL then DEFAULT
	"CREATE TABLE t (a bigint COMMENT 'x' WITH (foo = 'bar'))", // column COMMENT + WITH
	"CREATE TABLE t (a row(x bigint, y double))",
	"CREATE TABLE t (a array(bigint))",
	"CREATE TABLE cat.sch.t (a bigint)",
	"CREATE TABLE t (LIKE a EXCLUDING PROPERTIES)",
	"CREATE OR REPLACE TABLE t (a bigint)",
	// negatives
	"CREATE TABLE t ()",                                  // D-CT3: empty body
	"CREATE TABLE t",                                     // D-CT4: neither body nor AS
	"CREATE TABLE a.b.c.d (a bigint)",                    // D-CT1: 4-part name
	"CREATE TABLE t (a)",                                 // bare ident, no type, no AS
	"CREATE OR REPLACE TABLE IF NOT EXISTS t (a bigint)", // D-CT6: OR REPLACE + IF NOT EXISTS conflict (SYNTAX_ERROR)
	"CREATE TABLE t (a bigint NOT NULL NOT NULL)",        // duplicate NOT NULL
	"CREATE TABLE t (a bigint,)",                         // trailing comma
	"CREATE TABLE t (a bigint COMMENT 'x' NOT NULL)",     // D-CT2 order: COMMENT before NOT NULL
	"CREATE TABLE t (a bigint WITH (p = 1) NOT NULL)",    // D-CT2 order: WITH before NOT NULL

	// =====================================================================
	// CREATE TABLE AS SELECT  (truth1 create-table-as + truth2)
	// =====================================================================
	"CREATE TABLE orders_column_aliased (order_date, total_price) AS SELECT orderdate, totalprice FROM orders",
	"CREATE TABLE orders_by_date COMMENT 'Summary of orders by date' WITH (format = 'ORC') AS SELECT orderdate, sum(totalprice) AS price FROM orders GROUP BY orderdate",
	"CREATE TABLE IF NOT EXISTS orders_by_date AS SELECT orderdate, sum(totalprice) AS price FROM orders GROUP BY orderdate",
	"CREATE TABLE empty_nation AS SELECT * FROM nation WITH NO DATA",
	"CREATE TABLE foo AS SELECT * FROM t",
	"CREATE TABLE foo(x) AS SELECT a FROM t",
	"CREATE TABLE foo(x,y) AS SELECT a,b FROM t",
	"CREATE TABLE IF NOT EXISTS foo AS SELECT * FROM t",
	"CREATE TABLE foo AS SELECT * FROM t WITH NO DATA",
	"CREATE TABLE foo(x) AS SELECT a FROM t WITH NO DATA",
	"CREATE TABLE foo AS SELECT * FROM t WITH DATA",
	"CREATE TABLE foo WITH ( string = 'bar', long = 42, computed = 'ban' || 'ana', a = ARRAY[ 'v1', 'v2' ] ) AS SELECT * FROM t",
	"CREATE TABLE foo COMMENT 'test' WITH ( string = 'bar' ) AS SELECT * FROM t WITH NO DATA",
	`CREATE TABLE foo(x,y) COMMENT 'test' WITH ( "string" = 'bar', "long" = 42, computed = 'ban' || 'ana', a = ARRAY[ 'v1', 'v2' ] ) AS SELECT a,b FROM t WITH NO DATA`,
	"CREATE TABLE t AS (SELECT 1)",                           // parenthesized AS query
	"CREATE TABLE t AS TABLE other",                          // TABLE query source
	"CREATE TABLE t AS WITH c AS (SELECT 1) SELECT * FROM c", // WITH source
	"CREATE OR REPLACE TABLE t AS SELECT * FROM x",
	// negatives
	"CREATE TABLE t (a bigint) AS SELECT * FROM x", // column defs + AS (alias list must be bare)
	"CREATE TABLE t AS",                            // dangling AS
	"CREATE TABLE t AS SELECT * FROM x WITH",       // dangling WITH after query

	// =====================================================================
	// DROP TABLE  (truth1 drop-table + truth2 drop_table.sql)
	// =====================================================================
	"DROP TABLE orders_by_date",
	"DROP TABLE IF EXISTS orders_by_date",
	"DROP TABLE a",
	"DROP TABLE a.b",
	"DROP TABLE a.b.c",
	`DROP TABLE a."b/y".c`,
	"DROP TABLE IF EXISTS a.b.c",
	// negatives
	"DROP TABLE",           // missing name
	"DROP TABLE a.b.c.d",   // D-DR1: 4-part name
	"DROP TABLE a CASCADE", // D-DR2: TABLE has no CASCADE

	// =====================================================================
	// ALTER TABLE — rename  (truth1 alter-table-rename)
	// =====================================================================
	"ALTER TABLE users RENAME TO people",
	"ALTER TABLE IF EXISTS users RENAME TO people",
	"ALTER TABLE a.b.c RENAME TO d.e.f",
	"ALTER TABLE a RENAME TO b.c.d.e", // D-AT5: RENAME TO target is UNBOUNDED (4-part accepted; fails semantically)
	// negatives
	"ALTER TABLE users RENAME",        // dangling
	"ALTER TABLE a.b.c.d RENAME TO x", // 4-part SOURCE is SYNTAX_ERROR (source capped at 3)

	// =====================================================================
	// ALTER TABLE — add column  (truth1 alter-table-add-column, D-AT2)
	// =====================================================================
	"ALTER TABLE users ADD COLUMN zip varchar",
	"ALTER TABLE users ADD COLUMN zip varchar DEFAULT '90210'",
	"ALTER TABLE IF EXISTS users ADD COLUMN IF NOT EXISTS zip varchar",
	"ALTER TABLE users ADD COLUMN id varchar FIRST",
	"ALTER TABLE users ADD COLUMN zip varchar AFTER country",
	"ALTER TABLE users ADD COLUMN zip varchar LAST",
	"ALTER TABLE users ADD COLUMN c bigint NOT NULL COMMENT 'x' WITH (p = 1)",
	// negatives
	"ALTER TABLE users ADD COLUMN zip",               // missing type
	"ALTER TABLE users ADD zip varchar",              // missing COLUMN keyword
	"ALTER TABLE users ADD COLUMN zip varchar AFTER", // dangling AFTER

	// =====================================================================
	// ALTER TABLE — drop / rename column  (truth1 + truth2 drop_column.sql)
	// =====================================================================
	"ALTER TABLE users DROP COLUMN zip",
	"ALTER TABLE IF EXISTS users DROP COLUMN IF EXISTS zip",
	"ALTER TABLE foo.t DROP COLUMN c",
	`ALTER TABLE "t x" DROP COLUMN "c d"`,
	"ALTER TABLE foo.t DROP COLUMN IF EXISTS c",
	"ALTER TABLE users RENAME COLUMN id TO user_id",
	"ALTER TABLE IF EXISTS users RENAME COLUMN IF EXISTS id TO user_id",
	"ALTER TABLE t DROP COLUMN a.b", // nested-field path
	// negatives
	"ALTER TABLE users DROP COLUMN",             // missing column
	"ALTER TABLE users RENAME COLUMN id",        // missing TO
	"ALTER TABLE users RENAME COLUMN id TO a.b", // new name must be single ident

	// =====================================================================
	// ALTER TABLE — alter column  (truth1 + truth2, D-AT3)
	// =====================================================================
	"ALTER TABLE users ALTER COLUMN id SET DATA TYPE bigint",
	"ALTER TABLE foo.bar.baz ALTER COLUMN col1 SET DATA TYPE bigint",
	"ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active'",
	"ALTER TABLE users ALTER COLUMN status DROP DEFAULT",
	"ALTER TABLE users ALTER COLUMN id DROP NOT NULL",
	"ALTER TABLE IF EXISTS users ALTER COLUMN id SET DATA TYPE bigint",
	// negatives
	"ALTER TABLE users ALTER COLUMN id SET DATA TYPE",   // missing type
	"ALTER TABLE users ALTER COLUMN id SET",             // dangling SET
	"ALTER TABLE users ALTER COLUMN id DROP",            // dangling DROP
	"ALTER TABLE users ALTER COLUMN id SET TYPE bigint", // SET TYPE (missing DATA)

	// =====================================================================
	// ALTER TABLE — set authorization / properties  (truth1 + truth2)
	// =====================================================================
	"ALTER TABLE people SET AUTHORIZATION alice",
	"ALTER TABLE people SET AUTHORIZATION ROLE PUBLIC",
	"ALTER TABLE foo.bar.baz SET AUTHORIZATION qux",
	"ALTER TABLE foo.bar.baz SET AUTHORIZATION USER qux",
	"ALTER TABLE foo.bar.baz SET AUTHORIZATION ROLE qux",
	"ALTER TABLE people SET PROPERTIES x = 'y'",
	`ALTER TABLE people SET PROPERTIES foo = 123, "foo bar" = 456`,
	"ALTER TABLE people SET PROPERTIES x = DEFAULT",
	// negatives
	"ALTER TABLE people SET PROPERTIES",               // empty
	"ALTER TABLE people SET PROPERTIES x",             // missing = value
	"ALTER TABLE people SET AUTHORIZATION",            // dangling
	"ALTER TABLE IF EXISTS t SET PROPERTIES x = 1",    // D-AT6: IF EXISTS not allowed on SET PROPERTIES
	"ALTER TABLE IF EXISTS t SET AUTHORIZATION alice", // D-AT6: IF EXISTS not allowed on SET AUTHORIZATION

	// =====================================================================
	// ALTER TABLE — execute  (truth1 alter-table-execute, D-AT4)
	// =====================================================================
	"ALTER TABLE example.test.test_table EXECUTE optimize(file_size_threshold => '16MB')",
	"ALTER TABLE t EXECUTE optimize",
	"ALTER TABLE t EXECUTE optimize()",
	"ALTER TABLE t EXECUTE optimize(1, 2)",
	"ALTER TABLE t EXECUTE optimize WHERE x > 1",
	"ALTER TABLE t EXECUTE optimize(threshold => '1MB') WHERE part = 'a'",
	// negatives
	"ALTER TABLE t EXECUTE",                    // missing procedure
	"ALTER TABLE t EXECUTE optimize WHERE",     // dangling WHERE
	"ALTER TABLE IF EXISTS t EXECUTE optimize", // D-AT6: IF EXISTS not allowed on EXECUTE

	// =====================================================================
	// CREATE VIEW  (truth1 create-view + truth2 create_view.sql)
	// =====================================================================
	"CREATE VIEW test AS SELECT orderkey, orderstatus, totalprice / 2 AS half FROM orders",
	"CREATE VIEW test_with_comment COMMENT 'A view to keep track of orders.' AS SELECT orderkey, orderstatus, totalprice FROM orders",
	"CREATE OR REPLACE VIEW test AS SELECT orderkey, orderstatus, totalprice / 4 AS quarter FROM orders",
	"CREATE VIEW a AS SELECT * FROM t",
	"CREATE VIEW a SECURITY DEFINER AS SELECT * FROM t",
	"CREATE VIEW a SECURITY INVOKER AS SELECT * FROM t",
	"CREATE VIEW a COMMENT 'comment' SECURITY DEFINER AS SELECT * FROM t",
	"CREATE VIEW a COMMENT '' AS SELECT * FROM t",
	"CREATE VIEW bar.foo AS SELECT * FROM t",
	`CREATE VIEW "awesome schema"."awesome view" AS SELECT * FROM t`,
	// negatives
	"CREATE VIEW a", // missing AS query
	"CREATE VIEW IF NOT EXISTS a AS SELECT 1",                // D-VW1: no IF NOT EXISTS for plain view
	"CREATE VIEW a SECURITY AS SELECT 1",                     // SECURITY without DEFINER/INVOKER
	"CREATE VIEW a AS",                                       // dangling AS
	"CREATE VIEW v SECURITY DEFINER COMMENT 'c' AS SELECT 1", // clause order: COMMENT before SECURITY

	// =====================================================================
	// DROP VIEW  (truth1 drop-view + truth2 drop_view.sql)
	// =====================================================================
	"DROP VIEW orders_by_date",
	"DROP VIEW IF EXISTS orders_by_date",
	"DROP VIEW a",
	"DROP VIEW a.b.c",
	// negatives
	"DROP VIEW",           // missing name
	"DROP VIEW a CASCADE", // VIEW has no CASCADE

	// =====================================================================
	// ALTER VIEW  (truth1 alter-view + truth2 rename_view / set_authorization)
	// =====================================================================
	"ALTER VIEW people RENAME TO users",
	"ALTER VIEW a RENAME TO b.c.d.e", // D-AT5: RENAME TO target unbounded
	"ALTER VIEW people REFRESH",
	"ALTER VIEW people SET AUTHORIZATION alice",
	"ALTER VIEW a RENAME TO b",
	"ALTER VIEW foo.bar.baz SET AUTHORIZATION qux",
	"ALTER VIEW foo.bar.baz SET AUTHORIZATION USER qux",
	"ALTER VIEW foo.bar.baz SET AUTHORIZATION ROLE qux",
	// negatives
	"ALTER VIEW a RENAME", // dangling
	"ALTER VIEW a SET",    // dangling SET

	// =====================================================================
	// CREATE MATERIALIZED VIEW  (truth1 + truth2, D-MV1)
	// =====================================================================
	"CREATE MATERIALIZED VIEW cancelled_orders AS SELECT orderkey, totalprice FROM orders WHERE orderstatus = 3",
	"CREATE OR REPLACE MATERIALIZED VIEW order_totals_by_date AS SELECT orderdate, sum(totalprice) AS price FROM orders GROUP BY orderdate",
	"CREATE MATERIALIZED VIEW orders_nation_mkgsegment COMMENT 'Orders with nation and market segment data' WITH ( partitioning = ARRAY['mktsegment', 'nationkey'] ) AS SELECT o.*, c.nationkey, c.mktsegment FROM orders AS o JOIN customer AS c ON o.custkey = c.custkey",
	"CREATE MATERIALIZED VIEW orders_summary GRACE PERIOD INTERVAL '1' HOUR WHEN STALE FAIL AS SELECT orderdate, sum(totalprice) AS price FROM orders GROUP BY orderdate",
	"CREATE MATERIALIZED VIEW orders_summary GRACE PERIOD INTERVAL '1' HOUR WHEN STALE INLINE AS SELECT orderdate, sum(totalprice) AS price FROM orders GROUP BY orderdate",
	"CREATE MATERIALIZED VIEW a AS SELECT * FROM t",
	"CREATE OR REPLACE MATERIALIZED VIEW catalog.schema.matview COMMENT 'A simple materialized view' AS SELECT * FROM catalog2.schema2.tab",
	"CREATE OR REPLACE MATERIALIZED VIEW catalog.schema.matview COMMENT 'A simple materialized view' WITH (partitioned_by = ARRAY ['dateint']) AS SELECT * FROM catalog2.schema2.tab",
	"CREATE OR REPLACE MATERIALIZED VIEW catalog.schema.matview COMMENT 'A partitioned materialized view' WITH (partitioned_by = ARRAY ['dateint']) AS WITH a (t, u) AS (SELECT * FROM x), b AS (SELECT * FROM a) TABLE b",
	"CREATE MATERIALIZED VIEW a GRACE PERIOD INTERVAL '2' DAY AS SELECT * FROM t",
	"CREATE MATERIALIZED VIEW IF NOT EXISTS a AS SELECT * FROM t",
	"CREATE MATERIALIZED VIEW a WHEN STALE INLINE AS SELECT * FROM t", // D-MV1 WHEN STALE alone
	// negatives
	"CREATE MATERIALIZED VIEW a",                          // missing AS
	"CREATE MATERIALIZED VIEW a GRACE PERIOD AS SELECT 1", // GRACE PERIOD without interval
	"CREATE MATERIALIZED VIEW a WHEN STALE AS SELECT 1",   // WHEN STALE without mode
	"CREATE MATERIALIZED VIEW a WHEN AS SELECT 1",         // WHEN without STALE
	// fixed-clause-order negatives (D-MV1): clauses out of GRACE/WHEN/COMMENT/WITH order
	"CREATE MATERIALIZED VIEW m COMMENT 'c' GRACE PERIOD INTERVAL '1' HOUR AS SELECT 1",       // COMMENT before GRACE PERIOD
	"CREATE MATERIALIZED VIEW m WITH (p = 1) COMMENT 'c' AS SELECT 1",                         // WITH before COMMENT
	"CREATE MATERIALIZED VIEW m WHEN STALE INLINE GRACE PERIOD INTERVAL '1' HOUR AS SELECT 1", // WHEN STALE before GRACE PERIOD

	// =====================================================================
	// DROP MATERIALIZED VIEW  (truth1 + truth2)
	// =====================================================================
	"DROP MATERIALIZED VIEW orders_by_date",
	"DROP MATERIALIZED VIEW IF EXISTS orders_by_date",
	"DROP MATERIALIZED VIEW a.b.c",
	// negatives
	"DROP MATERIALIZED VIEW", // missing name
	"DROP MATERIALIZED a",    // missing VIEW

	// =====================================================================
	// ALTER MATERIALIZED VIEW  (truth1 + truth2 rename_materialized_view.sql)
	// =====================================================================
	"ALTER MATERIALIZED VIEW people RENAME TO users",
	"ALTER MATERIALIZED VIEW a RENAME TO b.c.d.e", // D-AT5: RENAME TO target unbounded
	"ALTER MATERIALIZED VIEW IF EXISTS people RENAME TO users",
	"ALTER MATERIALIZED VIEW people SET PROPERTIES x = 'y'",
	`ALTER MATERIALIZED VIEW people SET PROPERTIES foo = 123, "foo bar" = 456`,
	"ALTER MATERIALIZED VIEW people SET PROPERTIES x = DEFAULT",
	"ALTER MATERIALIZED VIEW a RENAME TO b",
	"ALTER MATERIALIZED VIEW people SET AUTHORIZATION alice", // D-MV3: SET AUTHORIZATION (docs-ahead)
	"ALTER MATERIALIZED VIEW people SET AUTHORIZATION ROLE PUBLIC",
	"ALTER MATERIALIZED VIEW people SET AUTHORIZATION USER alice",
	// negatives
	"ALTER MATERIALIZED VIEW a RENAME",                                 // dangling
	"ALTER MATERIALIZED VIEW IF EXISTS people SET PROPERTIES x = 'y'",  // D-MV3: IF EXISTS not allowed on SET
	"ALTER MATERIALIZED VIEW IF EXISTS people SET AUTHORIZATION alice", // D-MV3: IF EXISTS not allowed on SET
	"ALTER MATERIALIZED VIEW a SET PROPERTIES",                         // empty

	// =====================================================================
	// REFRESH MATERIALIZED VIEW  (truth1 + truth2)
	// =====================================================================
	"REFRESH MATERIALIZED VIEW orders_by_date",
	"REFRESH MATERIALIZED VIEW test",
	`REFRESH MATERIALIZED VIEW "some name that contains space"`,
	// negatives
	"REFRESH MATERIALIZED VIEW", // missing name
	"REFRESH VIEW a",            // missing MATERIALIZED

	// =====================================================================
	// CREATE CATALOG  (truth1 create-catalog, D-CAT*)
	// =====================================================================
	"CREATE CATALOG tpch USING tpch",
	`CREATE CATALOG brain USING memory WITH ("memory.max-data-per-node" = '128MB')`,
	`CREATE CATALOG example USING postgresql WITH ("connection-url" = 'jdbc:pg:localhost:5432', "connection-user" = 'u', "connection-password" = 'p', "case-insensitive-name-matching" = 'true')`,
	"CREATE CATALOG IF NOT EXISTS c USING tpch",
	"CREATE CATALOG c USING tpch COMMENT 'my catalog'",
	"CREATE CATALOG c USING tpch AUTHORIZATION alice",
	"CREATE CATALOG c USING tpch COMMENT 'x' AUTHORIZATION alice WITH (a = '1')",
	// negatives
	"CREATE CATALOG c",           // D-CAT2: USING mandatory
	"CREATE CATALOG a.b USING t", // D-CAT1: catalog name is single ident
	"CREATE CATALOG c USING",     // dangling USING
	"CREATE CATALOG",             // missing name
	"CREATE CATALOG c USING tpch AUTHORIZATION alice COMMENT 'x'", // clause order: COMMENT before AUTHORIZATION
	"CREATE CATALOG c USING tpch WITH (a = '1') COMMENT 'x'",      // clause order: COMMENT before WITH

	// =====================================================================
	// DROP CATALOG  (truth1 drop-catalog)
	// =====================================================================
	"DROP CATALOG example",
	"DROP CATALOG IF EXISTS example",
	"DROP CATALOG example CASCADE",
	"DROP CATALOG example RESTRICT",
	// negatives
	"DROP CATALOG",     // missing name
	"DROP CATALOG a.b", // catalog name is single ident

	// =====================================================================
	// COMMENT ON  (truth1 comment + truth2 comment_*.sql, D-CM*)
	// =====================================================================
	"COMMENT ON TABLE users IS 'master table'",
	"COMMENT ON TABLE users IS NULL",
	"COMMENT ON VIEW users IS 'master view'",
	"COMMENT ON VIEW users IS NULL",
	"COMMENT ON COLUMN users.name IS 'full name'",
	"COMMENT ON COLUMN users.name IS NULL",
	"COMMENT ON TABLE a IS ''",
	"COMMENT ON COLUMN a.b IS 'test'",
	"COMMENT ON COLUMN a IS 'test'",
	"COMMENT ON COLUMN a.b.c IS 'test'",
	"COMMENT ON COLUMN a.b.c.d IS 'test'", // D-CM2: 4-part column reference
	"COMMENT ON TABLE cat.sch.t IS 'x'",
	// negatives
	"COMMENT ON TABLE a IS 5",         // D-CM1: value must be string or NULL
	"COMMENT ON TABLE a",              // missing IS clause
	"COMMENT ON SCHEMA a IS 'x'",      // COMMENT ON SCHEMA not supported in 481
	"COMMENT TABLE a IS 'x'",          // missing ON
	"COMMENT ON TABLE a.b.c.d IS 'x'", // D-CM2: TABLE name max 3 parts

	// =====================================================================
	// ANALYZE  (truth1 analyze + truth2 analyze.sql, D-AN1)
	// =====================================================================
	"ANALYZE web",
	"ANALYZE hive.default.stores",
	"ANALYZE hive.default.sales WITH (partitions = ARRAY[ARRAY['1992-01-01'], ARRAY['1992-01-02']])",
	"ANALYZE hive.default.customers WITH (partitions = ARRAY[ARRAY['CA', 'San Francisco'], ARRAY['NY', 'NY']])",
	"ANALYZE hive.default.sales WITH (partitions = ARRAY[ARRAY['1992-01-01'], ARRAY['1992-01-02']], columns = ARRAY['department', 'product_id'])",
	"ANALYZE foo",
	`ANALYZE foo WITH ( "string" = 'bar', "long" = 42, computed = concat('ban', 'ana'), a = ARRAY[ 'v1', 'v2' ] )`,
	// negatives
	"ANALYZE",             // missing name
	"ANALYZE a.b.c.d",     // D-AN1: 4-part name
	"ANALYZE foo WITH ()", // empty property list
}

// TestDDL_OracleDifferential is the authoritative accept/reject gate: omni's
// Parse verdict must equal Trino 481's verdict for every corpus statement.
func TestDDL_OracleDifferential(t *testing.T) {
	if testing.Short() {
		t.Skip("trino oracle: skipped in -short mode")
	}
	o := connectOracle(t)
	for _, sql := range ddlOracleCorpus {
		sql := sql
		t.Run(truncateName(sql), func(t *testing.T) {
			_, errs := Parse(sql)
			omniAccepts := len(errs) == 0

			trinoAccepts, ok := oracleAccepts(t, o, sql)
			if !ok {
				t.Skip("oracle unreachable for this case")
			}
			if omniAccepts != trinoAccepts {
				t.Errorf("MISMATCH %q: omni accepts=%v (errs=%v), Trino accepts=%v",
					sql, omniAccepts, errs, trinoAccepts)
			}
		})
	}
}
