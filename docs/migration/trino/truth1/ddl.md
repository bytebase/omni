# Trino DDL/Schema Syntax Reference — crawled from trino.io/docs/current/ (Trino 481, "current")

---

## create-schema

- **syntax:**
  ```
  CREATE SCHEMA [ IF NOT EXISTS ] schema_name
  [ AUTHORIZATION ( user | USER user | ROLE role ) ]
  [ WITH ( property_name = expression [, ...] ) ]
  ```
- **source_url:** https://trino.io/docs/current/sql/create-schema.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE SCHEMA web;

  CREATE SCHEMA hive.sales;

  CREATE SCHEMA IF NOT EXISTS traffic;

  CREATE SCHEMA web AUTHORIZATION alice;

  CREATE SCHEMA web AUTHORIZATION alice WITH ( LOCATION = '/hive/data/web' );

  CREATE SCHEMA web AUTHORIZATION ROLE PUBLIC;

  CREATE SCHEMA web AUTHORIZATION ROLE PUBLIC WITH ( LOCATION = '/hive/data/web' );
  ```
- **notes:** `IF NOT EXISTS` suppresses error if schema already exists. `AUTHORIZATION` accepts bare user name, `USER user`, or `ROLE role` (including `ROLE PUBLIC`). `WITH` passes connector-specific properties (e.g. `LOCATION`).

---

## drop-schema

- **syntax:**
  ```
  DROP SCHEMA [ IF EXISTS ] schema_name [ CASCADE | RESTRICT ]
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-schema.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP SCHEMA web;

  DROP SCHEMA IF EXISTS sales;

  DROP SCHEMA archive CASCADE;

  DROP SCHEMA archive RESTRICT;
  ```
- **notes:** Schema must be empty unless `CASCADE` is specified. `CASCADE` removes the schema together with all contained objects. `RESTRICT` (default behavior) only succeeds when schema is empty.

---

## alter-schema-rename

- **syntax:**
  ```
  ALTER SCHEMA name RENAME TO new_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-schema.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER SCHEMA web RENAME TO traffic;
  ```
- **notes:** Renames the schema to a new identifier.

---

## alter-schema-set-authorization

- **syntax:**
  ```
  ALTER SCHEMA name SET AUTHORIZATION ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-schema.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER SCHEMA web SET AUTHORIZATION alice;

  ALTER SCHEMA web SET AUTHORIZATION ROLE PUBLIC;
  ```
- **notes:** Transfers ownership to a user or role. Accepts bare user name, `USER user`, or `ROLE role`.

---

## create-table

- **syntax:**
  ```
  CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ]
  table_name (
    { column_name data_type [ DEFAULT default ] [ NOT NULL ]
        [ COMMENT comment ]
        [ WITH ( property_name = expression [, ...] ) ]
    | LIKE existing_table_name
        [ { INCLUDING | EXCLUDING } PROPERTIES ]
    }
    [, ...]
  )
  [ COMMENT table_comment ]
  [ WITH ( property_name = expression [, ...] ) ]
  ```
- **source_url:** https://trino.io/docs/current/sql/create-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE TABLE orders (
    orderkey bigint,
    orderstatus varchar,
    totalprice double,
    orderdate date
  )
  WITH (format = 'ORC');

  CREATE TABLE IF NOT EXISTS orders (
    orderkey bigint,
    orderstatus varchar,
    totalprice double COMMENT 'Price in cents.',
    orderdate date,
    status varchar DEFAULT 'created'
  )
  COMMENT 'A table to keep track of orders.';

  CREATE TABLE bigger_orders (
    another_orderkey bigint,
    LIKE orders,
    another_orderdate date
  );
  ```
- **notes:** `OR REPLACE` support is connector-dependent. Multiple `LIKE` clauses are allowed. `INCLUDING PROPERTIES` copies table properties from the source table; default is `EXCLUDING PROPERTIES`. Column-level `WITH` sets per-column connector properties. Available table properties can be discovered via `SELECT * FROM system.metadata.table_properties`.

---

## create-table-as

- **syntax:**
  ```
  CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ] table_name [ ( column_alias, ... ) ]
  [ COMMENT table_comment ]
  [ WITH ( property_name = expression [, ...] ) ]
  AS query
  [ WITH [ NO ] DATA ]
  ```
- **source_url:** https://trino.io/docs/current/sql/create-table-as.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE TABLE orders_column_aliased (order_date, total_price)
  AS
  SELECT orderdate, totalprice
  FROM orders;

  CREATE TABLE orders_by_date
  COMMENT 'Summary of orders by date'
  WITH (format = 'ORC')
  AS
  SELECT orderdate, sum(totalprice) AS price
  FROM orders
  GROUP BY orderdate;

  CREATE TABLE IF NOT EXISTS orders_by_date AS
  SELECT orderdate, sum(totalprice) AS price
  FROM orders
  GROUP BY orderdate;

  CREATE TABLE empty_nation AS
  SELECT *
  FROM nation
  WITH NO DATA;
  ```
- **notes:** `OR REPLACE` and `IF NOT EXISTS` cannot be used together. `WITH NO DATA` creates the table structure without inserting rows. Column aliases in the parenthesized list rename the SELECT output columns.

---

## drop-table

- **syntax:**
  ```
  DROP TABLE [ IF EXISTS ] table_name
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP TABLE orders_by_date;

  DROP TABLE IF EXISTS orders_by_date;
  ```
- **notes:** `IF EXISTS` suppresses errors when the table does not exist. Does not suppress errors if a view with the same name exists.

---

## alter-table-rename

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name RENAME TO new_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users RENAME TO people;

  ALTER TABLE IF EXISTS users RENAME TO people;
  ```
- **notes:** `IF EXISTS` suppresses errors if the table does not exist.

---

## alter-table-add-column

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name ADD COLUMN [ IF NOT EXISTS ] column_name data_type
    [ DEFAULT default ] [ NOT NULL ] [ COMMENT comment ]
    [ WITH ( property_name = expression [, ...] ) ]
    [ FIRST | LAST | AFTER after_column_name ]
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users ADD COLUMN zip varchar;

  ALTER TABLE users ADD COLUMN zip varchar DEFAULT '90210';

  ALTER TABLE IF EXISTS users ADD COLUMN IF NOT EXISTS zip varchar;

  ALTER TABLE users ADD COLUMN id varchar FIRST;

  ALTER TABLE users ADD COLUMN zip varchar AFTER country;
  ```
- **notes:** `IF NOT EXISTS` on the column name suppresses errors if the column already exists. Column position can be controlled with `FIRST`, `LAST`, or `AFTER column_name`. Column-level `WITH` sets connector-specific column properties.

---

## alter-table-drop-column

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name DROP COLUMN [ IF EXISTS ] column_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users DROP COLUMN zip;

  ALTER TABLE IF EXISTS users DROP COLUMN IF EXISTS zip;
  ```
- **notes:** Both `IF EXISTS` modifiers are independent; each suppresses the respective "not found" error.

---

## alter-table-rename-column

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name RENAME COLUMN [ IF EXISTS ] old_name TO new_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users RENAME COLUMN id TO user_id;

  ALTER TABLE IF EXISTS users RENAME COLUMN IF EXISTS id TO user_id;
  ```
- **notes:** `IF EXISTS` on the column suppresses errors if the column does not exist.

---

## alter-table-alter-column-set-default

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name ALTER COLUMN column_name SET DEFAULT expression
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';
  ```
- **notes:** Sets a default value on an existing column.

---

## alter-table-alter-column-drop-default

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name ALTER COLUMN column_name DROP DEFAULT
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users ALTER COLUMN status DROP DEFAULT;
  ```
- **notes:** Removes the default value from an existing column.

---

## alter-table-alter-column-set-data-type

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name ALTER COLUMN column_name SET DATA TYPE new_type
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users ALTER COLUMN id SET DATA TYPE bigint;
  ```
- **notes:** Changes the data type of an existing column. Connector support varies.

---

## alter-table-alter-column-drop-not-null

- **syntax:**
  ```
  ALTER TABLE [ IF EXISTS ] name ALTER COLUMN column_name DROP NOT NULL
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE users ALTER COLUMN id DROP NOT NULL;
  ```
- **notes:** Removes the `NOT NULL` constraint from an existing column.

---

## alter-table-set-authorization

- **syntax:**
  ```
  ALTER TABLE name SET AUTHORIZATION ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE people SET AUTHORIZATION alice;

  ALTER TABLE people SET AUTHORIZATION ROLE PUBLIC;
  ```
- **notes:** `IF EXISTS` is not documented for this sub-form. Accepts bare user name, `USER user`, or `ROLE role`.

---

## alter-table-set-properties

- **syntax:**
  ```
  ALTER TABLE name SET PROPERTIES property_name = expression [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE people SET PROPERTIES x = 'y';

  ALTER TABLE people SET PROPERTIES foo = 123, "foo bar" = 456;

  ALTER TABLE people SET PROPERTIES x = DEFAULT;
  ```
- **notes:** `DEFAULT` keyword resets a property to its default value. Omitting a property leaves it unchanged. `IF EXISTS` is not documented for this sub-form.

---

## alter-table-execute

- **syntax:**
  ```
  ALTER TABLE name EXECUTE command [ ( parameter => expression [, ... ] ) ]
      [ WHERE expression ]
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-table.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER TABLE example.test.test_table EXECUTE optimize(file_size_threshold => '16MB');
  ```
- **notes:** Named parameters use the `=>` operator. The optional `WHERE` clause restricts which partitions or rows the procedure operates on. `IF EXISTS` is not documented for this sub-form.

---

## truncate-table

- **syntax:**
  ```
  TRUNCATE TABLE table_name
  ```
- **source_url:** https://trino.io/docs/current/sql/truncate.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  TRUNCATE TABLE orders;
  ```
- **notes:** Deletes all rows from a table. No `IF EXISTS` or `CASCADE` clauses are documented.

---

## comment-on-table

- **syntax:**
  ```
  COMMENT ON TABLE name IS 'comments'
  ```
- **source_url:** https://trino.io/docs/current/sql/comment.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  COMMENT ON TABLE users IS 'master table';

  COMMENT ON TABLE users IS NULL;
  ```
- **notes:** Setting the comment to `NULL` removes the existing comment.

---

## comment-on-view

- **syntax:**
  ```
  COMMENT ON VIEW name IS 'comments'
  ```
- **source_url:** https://trino.io/docs/current/sql/comment.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  COMMENT ON VIEW users IS 'master view';

  COMMENT ON VIEW users IS NULL;
  ```
- **notes:** Setting the comment to `NULL` removes the existing comment.

---

## comment-on-column

- **syntax:**
  ```
  COMMENT ON COLUMN name IS 'comments'
  ```
- **source_url:** https://trino.io/docs/current/sql/comment.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  COMMENT ON COLUMN users.name IS 'full name';

  COMMENT ON COLUMN users.name IS NULL;
  ```
- **notes:** `name` is a qualified column reference in the form `table.column` or `catalog.schema.table.column`. Setting the comment to `NULL` removes the existing comment.

---

## create-view

- **syntax:**
  ```
  CREATE [ OR REPLACE ] VIEW view_name
  [ COMMENT view_comment ]
  [ SECURITY { DEFINER | INVOKER } ]
  AS query
  ```
- **source_url:** https://trino.io/docs/current/sql/create-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE VIEW test AS
  SELECT orderkey, orderstatus, totalprice / 2 AS half
  FROM orders;

  CREATE VIEW test_with_comment
  COMMENT 'A view to keep track of orders.'
  AS
  SELECT orderkey, orderstatus, totalprice
  FROM orders;

  CREATE OR REPLACE VIEW test AS
  SELECT orderkey, orderstatus, totalprice / 4 AS quarter
  FROM orders;
  ```
- **notes:** `SECURITY DEFINER` (default) executes the view query with the view creator's permissions. `SECURITY INVOKER` uses the executing user's permissions. The `current_user` function always returns the actual query executor regardless of security mode.

---

## drop-view

- **syntax:**
  ```
  DROP VIEW [ IF EXISTS ] view_name
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP VIEW orders_by_date;

  DROP VIEW IF EXISTS orders_by_date;
  ```
- **notes:** `IF EXISTS` suppresses errors when the view does not exist.

---

## alter-view-rename

- **syntax:**
  ```
  ALTER VIEW name RENAME TO new_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER VIEW people RENAME TO users;
  ```
- **notes:** Renames an existing view.

---

## alter-view-refresh

- **syntax:**
  ```
  ALTER VIEW name REFRESH
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER VIEW people REFRESH;
  ```
- **notes:** Refreshes the view's definition or cached metadata. Connector support varies.

---

## alter-view-set-authorization

- **syntax:**
  ```
  ALTER VIEW name SET AUTHORIZATION ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER VIEW people SET AUTHORIZATION alice;
  ```
- **notes:** Transfers ownership of the view to a user or role.

---

## create-materialized-view

- **syntax:**
  ```
  CREATE [ OR REPLACE ] MATERIALIZED VIEW
  [ IF NOT EXISTS ] view_name
  [ GRACE PERIOD interval ]
  [ WHEN STALE ( INLINE | FAIL ) ]
  [ COMMENT string ]
  [ WITH ( property_name = expression [, ...] ) ]
  AS query
  ```
- **source_url:** https://trino.io/docs/current/sql/create-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE MATERIALIZED VIEW cancelled_orders
  AS
      SELECT orderkey, totalprice
      FROM orders
      WHERE orderstatus = 3;

  CREATE OR REPLACE MATERIALIZED VIEW order_totals_by_date
  AS
      SELECT orderdate, sum(totalprice) AS price
      FROM orders
      GROUP BY orderdate;

  CREATE MATERIALIZED VIEW orders_nation_mkgsegment
  COMMENT 'Orders with nation and market segment data'
  WITH ( partitioning = ARRAY['mktsegment', 'nationkey'] )
  AS
      SELECT o.*, c.nationkey, c.mktsegment
      FROM orders AS o
      JOIN customer AS c
      ON o.custkey = c.custkey;

  CREATE MATERIALIZED VIEW orders_summary
  GRACE PERIOD INTERVAL '1' HOUR
  WHEN STALE FAIL
  AS
      SELECT orderdate, sum(totalprice) AS price
      FROM orders
      GROUP BY orderdate;

  CREATE MATERIALIZED VIEW orders_summary
  GRACE PERIOD INTERVAL '1' HOUR
  WHEN STALE INLINE
  AS
      SELECT orderdate, sum(totalprice) AS price
      FROM orders
      GROUP BY orderdate;
  ```
- **notes:** `OR REPLACE` and `IF NOT EXISTS` are mutually exclusive. A `REFRESH MATERIALIZED VIEW` must be issued after creation to populate data; the view is empty until refreshed. `GRACE PERIOD` specifies how long materialized data is considered fresh before recomputation; defaults to infinity if omitted. `WHEN STALE INLINE` (default) falls back to live computation when data is stale; `WHEN STALE FAIL` raises an error instead. Both `GRACE PERIOD` and `WHEN STALE` require connector support.

---

## drop-materialized-view

- **syntax:**
  ```
  DROP MATERIALIZED VIEW [ IF EXISTS ] view_name
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP MATERIALIZED VIEW orders_by_date;

  DROP MATERIALIZED VIEW IF EXISTS orders_by_date;
  ```
- **notes:** `IF EXISTS` suppresses errors when the materialized view does not exist.

---

## alter-materialized-view-rename

- **syntax:**
  ```
  ALTER MATERIALIZED VIEW [ IF EXISTS ] name RENAME TO new_name
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER MATERIALIZED VIEW people RENAME TO users;

  ALTER MATERIALIZED VIEW IF EXISTS people RENAME TO users;
  ```
- **notes:** `IF EXISTS` suppresses errors when the materialized view does not exist.

---

## alter-materialized-view-set-properties

- **syntax:**
  ```
  ALTER MATERIALIZED VIEW name SET PROPERTIES property_name = expression [, ...]
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER MATERIALIZED VIEW people SET PROPERTIES x = 'y';

  ALTER MATERIALIZED VIEW people SET PROPERTIES foo = 123, "foo bar" = 456;

  ALTER MATERIALIZED VIEW people SET PROPERTIES x = DEFAULT;
  ```
- **notes:** `DEFAULT` keyword resets a property to its connector-default value. Omitting a property leaves it unchanged. Connector support varies.

---

## alter-materialized-view-set-authorization

- **syntax:**
  ```
  ALTER MATERIALIZED VIEW name SET AUTHORIZATION ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/alter-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ALTER MATERIALIZED VIEW people SET AUTHORIZATION alice;
  ```
- **notes:** Transfers ownership of the materialized view to a user or role.

---

## refresh-materialized-view

- **syntax:**
  ```
  REFRESH MATERIALIZED VIEW view_name
  ```
- **source_url:** https://trino.io/docs/current/sql/refresh-materialized-view.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  REFRESH MATERIALIZED VIEW orders_by_date;
  ```
- **notes:** Required after `CREATE MATERIALIZED VIEW` to populate data. Subsequent refreshes may be less expensive if the underlying data has not changed and the connector implements change-awareness. The entire view contents are replaced; no filtering is possible in this statement.

---

## create-catalog

- **syntax:**
  ```
  CREATE CATALOG catalog_name
  USING connector_name
  [ WITH ( property_name = expression [, ...] ) ]
  ```
- **source_url:** https://trino.io/docs/current/sql/create-catalog.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE CATALOG tpch USING tpch;

  CREATE CATALOG brain USING memory
  WITH ("memory.max-data-per-node" = '128MB');

  CREATE CATALOG example USING postgresql
  WITH (
    "connection-url" = 'jdbc:pg:localhost:5432',
    "connection-user" = '${ENV:POSTGRES_USER}',
    "connection-password" = '${ENV:POSTGRES_PASSWORD}',
    "case-insensitive-name-matching" = 'true'
  );
  ```
- **notes:** Requires catalog management type `dynamic`. Property names containing hyphens or other special characters must be double-quoted. All property values are varchars in single quotes (including numeric and boolean values). Environmental variable references use `${ENV:VARIABLE_NAME}` syntax. The full query including credentials is logged and visible in the Web UI.

---

## drop-catalog

- **syntax:**
  ```
  DROP CATALOG catalog_name
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-catalog.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP CATALOG example;
  ```
- **notes:** Requires catalog management type `dynamic`. Does not stop queries already running against the catalog. Several connectors (Hive, Iceberg, Delta Lake, Hudi) may fail to fully release HDFS/S3/GCS/Azure resources on drop. No `IF EXISTS` clause is documented.

---

## create-function

- **syntax:**
  ```
  CREATE [ OR REPLACE ] FUNCTION udf_name ( [ parameter_name data_type [, ...] ] )
    RETURNS return_type
    [ LANGUAGE language ]
    [ DETERMINISTIC | NOT DETERMINISTIC ]
    [ RETURNS NULL ON NULL INPUT | CALLED ON NULL INPUT ]
    [ SECURITY { DEFINER | INVOKER } ]
    [ COMMENT description ]
    [ WITH ( property_name = expression [, ...] ) ]
    { RETURN expression | BEGIN statements END }
  ```
- **source_url:** https://trino.io/docs/current/sql/create-function.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  CREATE FUNCTION example.default.meaning_of_life()
    RETURNS bigint
    BEGIN
      RETURN 42;
    END;

  CREATE FUNCTION meaning_of_life() RETURNS bigint RETURN 42;
  ```
- **notes:** Function name must be fully qualified (`catalog.schema.function`) unless default UDF storage catalog/schema is configured. `OR REPLACE` replaces an existing UDF without error. Supported languages include `SQL` and `PYTHON` by default. `RETURNS NULL ON NULL INPUT` skips execution when any parameter is NULL; `CALLED ON NULL INPUT` (default) executes normally with NULL values. `SECURITY DEFINER` uses the creator's permissions; `SECURITY INVOKER` uses the caller's permissions (catalog UDFs only). `WITH` passes language-specific properties (e.g. Python handler name). The connector backing the catalog must support UDF storage.

---

## drop-function

- **syntax:**
  ```
  DROP FUNCTION [ IF EXISTS ] udf_name ( [ [ parameter_name ] data_type [, ...] ] )
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-function.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  DROP FUNCTION example.default.meaning_of_life();

  DROP FUNCTION multiply_by_two(bigint);

  DROP FUNCTION meaning_of_life();
  ```
- **notes:** Function name must be fully qualified unless default UDF storage catalog/schema is configured. Parameter data types must be included when the UDF accepts parameters (to disambiguate overloads). `IF EXISTS` suppresses errors when the function does not exist.

---

## analyze

- **syntax:**
  ```
  ANALYZE table_name [ WITH ( property_name = expression [, ...] ) ]
  ```
- **source_url:** https://trino.io/docs/current/sql/analyze.html
- **version:** 481 (current)
- **example_sql:**
  ```sql
  ANALYZE web;

  ANALYZE hive.default.stores;

  ANALYZE hive.default.sales WITH (partitions = ARRAY[ARRAY['1992-01-01'], ARRAY['1992-01-02']]);

  ANALYZE hive.default.customers WITH (partitions = ARRAY[ARRAY['CA', 'San Francisco'], ARRAY['NY', 'NY']]);

  ANALYZE hive.default.sales WITH (
      partitions = ARRAY[ARRAY['1992-01-01'], ARRAY['1992-01-02']],
      columns = ARRAY['department', 'product_id']);
  ```
- **notes:** Collects table and column statistics for the query optimizer. The `WITH` clause is connector-specific; available properties can be discovered via `SELECT * FROM system.metadata.analyze_properties`. The `partitions` property accepts a nested array where each inner array is one partition key tuple. The `columns` property restricts statistics collection to specific columns.
