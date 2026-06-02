# Trino 481 — Administrative / Session / Transaction / Prepared-Statement / Utility Statements

---

## use

- **syntax:**
  ```
  USE catalog.schema
  USE schema
  ```
- **source_url:** https://trino.io/docs/current/sql/use.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  USE hive.finance;
  USE information_schema;
  ```
- **notes:** Two forms: explicit `catalog.schema` or schema-only (resolves against the current catalog). Updates the session scope for subsequent queries.

---

## set-session

- **syntax:**
  ```
  SET SESSION name = expression
  SET SESSION catalog.name = expression
  ```
- **source_url:** https://trino.io/docs/current/sql/set-session.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SET SESSION query_max_run_time = '10m';
  SET SESSION example.incremental_refresh_enabled = false;
  ```
- **notes:** System session properties apply cluster-wide; catalog session properties are connector-defined and must be prefixed with the catalog name. Changes are session-scoped and are lost when the session ends.

---

## reset-session

- **syntax:**
  ```
  RESET SESSION name
  RESET SESSION catalog.name
  ```
- **source_url:** https://trino.io/docs/current/sql/reset-session.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  RESET SESSION query_max_run_time;
  RESET SESSION hive.optimized_reader_enabled;
  ```
- **notes:** Restores the named property to its default value. Catalog-scoped properties use dot notation `catalog.name`.

---

## set-session-authorization

- **syntax:**
  ```
  SET SESSION AUTHORIZATION username
  ```
- **source_url:** https://trino.io/docs/current/sql/set-session-authorization.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SET SESSION AUTHORIZATION 'John';
  SET SESSION AUTHORIZATION "John";
  SET SESSION AUTHORIZATION John;
  ```
- **notes:** Changes the current user of the session. The original connecting user must have impersonation privileges for the target user (configured in system access control). All three quoting styles (single quotes, double quotes, or unquoted) are accepted.

---

## reset-session-authorization

- **syntax:**
  ```
  RESET SESSION AUTHORIZATION
  ```
- **source_url:** https://trino.io/docs/current/sql/reset-session-authorization.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  RESET SESSION AUTHORIZATION;
  ```
- **notes:** Resets the current authorization user back to the original user (the authenticated principal or the session user provided by the client). No arguments. Counterpart to `SET SESSION AUTHORIZATION`.

---

## set-role

- **syntax:**
  ```
  SET ROLE ( role | ALL | NONE )
  [ IN catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/set-role.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SET ROLE admin;
  SET ROLE ALL;
  SET ROLE NONE;
  SET ROLE analyst IN hive;
  ```
- **notes:** `SET ROLE role` activates a specific role (user must have a grant for it); `ALL` activates every granted role; `NONE` deactivates all roles. Optional `IN catalog` restricts the effect to a specific catalog. Some connectors do not support role management.

---

## set-path

- **syntax:**
  ```
  SET PATH path-element[, ...]
  ```
  where each path-element is `catalog.schema` or `schema`.
- **source_url:** https://trino.io/docs/current/sql/set-path.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SET PATH example.system;
  ```
- **notes:** Defines a collection of catalog/schema paths for function lookup in the current session. Allows unqualified function references without full `catalog.schema.function` paths. Both `catalog.schema` and bare `schema` (resolved against current catalog) are valid path elements. Session-scoped.

---

## set-time-zone

- **syntax:**
  ```
  SET TIME ZONE LOCAL
  SET TIME ZONE expression
  ```
- **source_url:** https://trino.io/docs/current/sql/set-time-zone.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SET TIME ZONE LOCAL;
  SET TIME ZONE '-08:00';
  SET TIME ZONE INTERVAL '10' HOUR;
  SET TIME ZONE INTERVAL -'08:00' HOUR TO MINUTE;
  SET TIME ZONE 'America/Los_Angeles';
  SET TIME ZONE concat_ws('/', 'America', 'Los_Angeles');
  ```
- **notes:** `LOCAL` resets to the initial session time zone. The expression can be a region-based ID string (e.g., `'America/Los_Angeles'`), a UTC-offset string (e.g., `'-08:00'`), or an interval expression (UTC offset range −14 to +14 hours). Has no effect if the `sql.forced-session-time-zone` configuration property is set.

---

## create-role

- **syntax:**
  ```
  CREATE ROLE role_name
  [ WITH ADMIN ( user | USER user | ROLE role | CURRENT_USER | CURRENT_ROLE ) ]
  [ IN catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/create-role.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CREATE ROLE admin;
  CREATE ROLE moderator WITH ADMIN USER bob;
  ```
- **notes:** Optional `WITH ADMIN` clause designates a role administrator; without it the current user becomes admin. Optional `IN catalog` scopes the role to a specific catalog. Some connectors do not support role management.

---

## drop-role

- **syntax:**
  ```
  DROP ROLE [ IF EXISTS ] role_name
  [ IN catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/drop-role.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DROP ROLE admin;
  DROP ROLE IF EXISTS analyst IN hive;
  ```
- **notes:** Requires admin privileges on the role. `IF EXISTS` suppresses errors when the role does not exist. `IN catalog` scopes to a catalog-level role. Some connectors do not support role management.

---

## grant-privilege

- **syntax:**
  ```
  GRANT ( privilege [, ...] | ( ALL PRIVILEGES ) )
  ON [ BRANCH branch_name IN ] ( table_name | TABLE table_name | SCHEMA schema_name )
  TO ( user | USER user | ROLE role )
  [ WITH GRANT OPTION ]
  ```
  where `privilege` is one of: `SELECT`, `INSERT`, `UPDATE`, `DELETE`.
- **source_url:** https://trino.io/docs/current/sql/grant.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  GRANT INSERT, SELECT ON orders TO alice;
  GRANT DELETE ON SCHEMA finance TO bob;
  GRANT SELECT ON nation TO alice WITH GRANT OPTION;
  GRANT SELECT ON orders TO ROLE PUBLIC;
  GRANT INSERT ON BRANCH audit IN orders TO alice;
  ```
- **notes:** `ALL PRIVILEGES` covers DELETE, INSERT, UPDATE, and SELECT. Granting on a table applies to all current and future columns; granting on a schema applies to all columns of all current and future tables in the schema. `WITH GRANT OPTION` allows the grantee to further grant privileges. Granting to `ROLE PUBLIC` affects all users. The optional `BRANCH branch_name IN` targets a specific branch of a versioned table. Some connectors have no support for `GRANT`.

---

## grant-role

- **syntax:**
  ```
  GRANT role_name [, ...]
  TO ( user | USER user_name | ROLE role_name ) [, ...]
  [ GRANTED BY ( user | USER user | ROLE role | CURRENT_USER | CURRENT_ROLE ) ]
  [ WITH ADMIN OPTION ]
  [ IN catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/grant-roles.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  GRANT bar TO USER foo;
  GRANT bar, foo TO USER baz, ROLE qux WITH ADMIN OPTION;
  ```
- **notes:** Grants one or more roles to one or more principals. `WITH ADMIN OPTION` allows recipients to grant the role to others. `GRANTED BY` designates which principal acts as grantor; defaults to current user. `IN catalog` scopes to catalog-specific roles. The executing user must be a role admin or possess the grant option. Some connectors do not support role management.

---

## revoke-privilege

- **syntax:**
  ```
  REVOKE [ GRANT OPTION FOR ]
  ( privilege [, ...] | ALL PRIVILEGES )
  ON [ BRANCH branch_name IN ] ( table_name | TABLE table_name | SCHEMA schema_name )
  FROM ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/revoke.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  REVOKE INSERT, SELECT ON orders FROM alice;
  REVOKE DELETE ON SCHEMA finance FROM bob;
  REVOKE GRANT OPTION FOR SELECT ON nation FROM ROLE PUBLIC;
  REVOKE ALL PRIVILEGES ON test FROM alice;
  REVOKE INSERT ON BRANCH audit IN orders FROM alice;
  ```
- **notes:** `GRANT OPTION FOR` removes only the right to grant privileges, leaving the privilege itself intact. Revoking from `ROLE PUBLIC` affects only that role assignment; users may still retain directly-granted or other role-based permissions. Revoke on a table affects all columns; revoke on a schema affects all columns of all tables. The executor must hold the privileges being revoked AND the grant option. Some connectors have no support for `REVOKE`.

---

## revoke-role

- **syntax:**
  ```
  REVOKE
  [ ADMIN OPTION FOR ]
  role_name [, ...]
  FROM ( user | USER user | ROLE role ) [, ...]
  [ GRANTED BY ( user | USER user | ROLE role | CURRENT_USER | CURRENT_ROLE ) ]
  [ IN catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/revoke-roles.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  REVOKE bar FROM USER foo;
  REVOKE ADMIN OPTION FOR bar, foo FROM USER baz, ROLE qux;
  ```
- **notes:** `ADMIN OPTION FOR` revokes only the grant permission without removing the role assignment itself. `GRANTED BY` designates the revoker principal; defaults to current user. `IN catalog` scopes to catalog-level roles. Executor must be role admin or hold the grant option. Some connectors do not support role management.

---

## deny

- **syntax:**
  ```
  DENY ( privilege [, ...] | ( ALL PRIVILEGES ) )
  ON [ BRANCH branch_name IN ] ( table_name | TABLE table_name | SCHEMA schema_name )
  TO ( user | USER user | ROLE role )
  ```
- **source_url:** https://trino.io/docs/current/sql/deny.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DENY INSERT, SELECT ON orders TO alice;
  DENY DELETE ON SCHEMA finance TO bob;
  DENY SELECT ON orders TO ROLE PUBLIC;
  DENY INSERT ON BRANCH audit IN orders TO alice;
  ```
- **notes:** Explicitly rejects specified privileges. On a table, denies on all current and future columns; on a schema, denies on all columns of all current and future tables. The default Trino system access controls and connectors have no support for `DENY`; custom connector/access-control implementation is required. Supports branch-specific denials via the `BRANCH branch_name IN` clause.

---

## show-grants

- **syntax:**
  ```
  SHOW GRANTS [ ON [ TABLE ] table_name ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-grants.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW GRANTS ON TABLE orders;
  SHOW GRANTS;
  ```
- **notes:** Without `ON`, shows grants for the current user across all tables in all schemas of the current catalog. Requires the current catalog to be set. Authentication must be enabled. Not all connectors support this statement.

---

## show-roles

- **syntax:**
  ```
  SHOW [ CURRENT ] ROLES [ FROM catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-roles.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW ROLES;
  SHOW ROLES FROM hive;
  SHOW CURRENT ROLES;
  SHOW CURRENT ROLES FROM hive;
  ```
- **notes:** Without `CURRENT`, lists all system roles or all roles in the specified catalog. With `CURRENT`, lists only the enabled/active roles for the session. `FROM catalog` scopes to a specific catalog.

---

## show-role-grants

- **syntax:**
  ```
  SHOW ROLE GRANTS [ FROM catalog ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-role-grants.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW ROLE GRANTS;
  SHOW ROLE GRANTS FROM hive;
  ```
- **notes:** Lists non-recursively the system roles (or catalog roles) that have been directly granted to the current session user. Does not enumerate transitively inherited roles.

---

## show-catalogs

- **syntax:**
  ```
  SHOW CATALOGS [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-catalogs.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CATALOGS;
  SHOW CATALOGS LIKE 't%';
  ```
- **notes:** Without `LIKE`, returns all available catalogs. The `LIKE` clause uses standard SQL pattern syntax for filtering.

---

## show-schemas

- **syntax:**
  ```
  SHOW SCHEMAS [ FROM catalog ] [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-schemas.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW SCHEMAS;
  SHOW SCHEMAS FROM tpch;
  SHOW SCHEMAS FROM tpch LIKE '__3%';
  ```
- **notes:** Without `FROM`, uses the current catalog. The `LIKE` clause supports `_` (single-character wildcard) and `%` (multi-character wildcard).

---

## show-tables

- **syntax:**
  ```
  SHOW TABLES [ FROM schema ] [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-tables.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW TABLES;
  SHOW TABLES FROM tpch.tiny;
  SHOW TABLES FROM tpch.tiny LIKE 'p%';
  ```
- **notes:** Returns both tables and views together. `FROM schema` accepts a fully qualified `catalog.schema` form, enabling cross-catalog queries. Without `FROM`, uses the current schema. The `LIKE` clause filters by name pattern.

---

## show-columns

- **syntax:**
  ```
  SHOW COLUMNS FROM table [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-columns.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW COLUMNS FROM nation;
  SHOW COLUMNS FROM nation LIKE '%key';
  ```
- **notes:** Returns columns: `Column`, `Type`, `Extra`, `Comment`. The optional `LIKE` clause filters the result to a subset of columns.

---

## show-create-table

- **syntax:**
  ```
  SHOW CREATE TABLE table_name
  ```
- **source_url:** https://trino.io/docs/current/sql/show-create-table.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CREATE TABLE sf1.orders;
  ```
- **notes:** Returns the `CREATE TABLE` statement that would recreate the table, including column definitions, data types, and connector-specific `WITH` properties (e.g., `format`, `partitioned_by`).

---

## show-create-view

- **syntax:**
  ```
  SHOW CREATE VIEW view_name
  ```
- **source_url:** https://trino.io/docs/current/sql/show-create-view.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CREATE VIEW my_view;
  ```
- **notes:** Returns the `CREATE VIEW` statement that defines the specified view. No additional options.

---

## show-create-schema

- **syntax:**
  ```
  SHOW CREATE SCHEMA schema_name
  ```
- **source_url:** https://trino.io/docs/current/sql/show-create-schema.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CREATE SCHEMA hive.web;
  ```
- **notes:** Returns the `CREATE SCHEMA` statement that would recreate the schema. No additional options.

---

## show-create-materialized-view

- **syntax:**
  ```
  SHOW CREATE MATERIALIZED VIEW view_name
  ```
- **source_url:** https://trino.io/docs/current/sql/show-create-materialized-view.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CREATE MATERIALIZED VIEW my_mv;
  ```
- **notes:** Returns the `CREATE MATERIALIZED VIEW` statement that defines the specified materialized view.

---

## show-create-function

- **syntax:**
  ```
  SHOW CREATE FUNCTION function_name
  ```
- **source_url:** https://trino.io/docs/current/sql/show-create-function.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW CREATE FUNCTION example.default.meaning_of_life;
  ```
- **notes:** Displays the `CREATE FUNCTION` statement for a user-defined function. The function name should be fully qualified (`catalog.schema.function`).

---

## show-functions

- **syntax:**
  ```
  SHOW FUNCTIONS [ FROM schema ] [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-functions.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW FUNCTIONS FROM example.default;
  SHOW FUNCTIONS LIKE 'array%';
  SHOW FUNCTIONS LIKE 'cf%';
  ```
- **notes:** Lists built-in, plugin, and user-defined functions. Result columns: `Function`, `Return Type`, `Argument Types`, `Function Type`, `Deterministic`, `Description`. `FROM schema` must use `catalog.schema` fully qualified form. `LIKE` filters by function name.

---

## show-session

- **syntax:**
  ```
  SHOW SESSION [ LIKE pattern ]
  ```
- **source_url:** https://trino.io/docs/current/sql/show-session.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW SESSION;
  SHOW SESSION LIKE 'query%';
  ```
- **notes:** Displays current session properties and their values. Without `LIKE`, shows all properties. Read-only diagnostic statement.

---

## show-stats

- **syntax:**
  ```
  SHOW STATS FOR table
  SHOW STATS FOR ( query )
  ```
- **source_url:** https://trino.io/docs/current/sql/show-stats.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  SHOW STATS FOR nation;
  SHOW STATS FOR (SELECT * FROM nation WHERE regionkey = 1);
  ```
- **notes:** Returns approximate statistics with one row per column plus a summary row (identified by `NULL` in `column_name`). Result columns: `column_name`, `data_size` (string types only), `distinct_values_count`, `nulls_fractions`, `row_count` (summary row only), `low_value`, `high_value` (DATE/numeric types). Returns `NULL` for any statistics that are unavailable or not populated on the data source.

---

## describe

- **syntax:**
  ```
  DESCRIBE table_name
  ```
- **source_url:** https://trino.io/docs/current/sql/describe.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DESCRIBE nation;
  ```
- **notes:** `DESCRIBE` is an alias for `SHOW COLUMNS`. Returns the same result set: `Column`, `Type`, `Extra`, `Comment`. `DESC` is also accepted as a short form.

---

## describe-input

- **syntax:**
  ```
  DESCRIBE INPUT statement_name
  ```
- **source_url:** https://trino.io/docs/current/sql/describe-input.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  PREPARE my_select1 FROM
  SELECT ? FROM nation WHERE regionkey = ? AND name < ?;
  DESCRIBE INPUT my_select1;

  PREPARE my_select2 FROM
  SELECT * FROM nation;
  DESCRIBE INPUT my_select2;
  ```
- **notes:** Lists input parameters (`?` placeholders) of a prepared statement with zero-based `Position` and `Type`. Parameters with indeterminate types are shown as `unknown`. Returns empty result set if the prepared statement has no parameters.

---

## describe-output

- **syntax:**
  ```
  DESCRIBE OUTPUT statement_name
  ```
- **source_url:** https://trino.io/docs/current/sql/describe-output.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  PREPARE my_select1 FROM
  SELECT * FROM nation;
  DESCRIBE OUTPUT my_select1;

  PREPARE my_select2 FROM
  SELECT count(*) as my_count, 1+2 FROM nation;
  DESCRIBE OUTPUT my_select2;

  PREPARE my_create FROM
  CREATE TABLE foo AS SELECT * FROM nation;
  DESCRIBE OUTPUT my_create;

  DESCRIBE OUTPUT (SELECT *, n_name AS "name" FROM nation);
  ```
- **notes:** Lists output columns of a prepared statement (or an inline query). Result columns: column name/alias, catalog, schema, table, type, type size in bytes, aliased (boolean). For DDL-like statements (e.g., `CREATE TABLE AS`), returns a single `rows` column. The fourth form (`DESCRIBE OUTPUT (query)`) works without a prior `PREPARE`.

---

## explain

- **syntax:**
  ```
  EXPLAIN [ ( option [, ...] ) ] statement

  where option can be one of:
      FORMAT { TEXT | GRAPHVIZ | JSON }
      TYPE { LOGICAL | DISTRIBUTED | VALIDATE | IO }
  ```
- **source_url:** https://trino.io/docs/current/sql/explain.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  EXPLAIN (TYPE LOGICAL) SELECT regionkey, count(*) FROM nation GROUP BY 1;
  EXPLAIN (TYPE LOGICAL, FORMAT JSON) SELECT regionkey, count(*) FROM nation GROUP BY 1;
  EXPLAIN (TYPE DISTRIBUTED) SELECT regionkey, count(*) FROM nation GROUP BY 1;
  EXPLAIN (TYPE DISTRIBUTED, FORMAT JSON) SELECT regionkey, count(*) FROM nation GROUP BY 1;
  EXPLAIN (TYPE VALIDATE) SELECT regionkey, count(*) FROM nation GROUP BY 1;
  EXPLAIN (TYPE IO, FORMAT JSON) INSERT INTO test_lineitem SELECT * FROM lineitem WHERE shipdate = '2020-02-01' AND quantity > 10;
  ```
- **notes:** Default TYPE is `DISTRIBUTED`; default FORMAT is `TEXT`. `TYPE LOGICAL` is deprecated and will be removed in a future release (use `DISTRIBUTED` instead). `TYPE VALIDATE` returns a boolean (`true`/`false`) indicating whether the statement is valid. `TYPE IO` (JSON only) provides input/output object details. JSON output format is not guaranteed to be backward compatible across Trino versions. Fragment types in distributed plans: `SINGLE`, `HASH`, `ROUND_ROBIN`, `BROADCAST`, `SOURCE`.

---

## explain-analyze

- **syntax:**
  ```
  EXPLAIN ANALYZE [ VERBOSE ] statement
  ```
- **source_url:** https://trino.io/docs/current/sql/explain-analyze.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  EXPLAIN ANALYZE SELECT count(*), clerk FROM orders
  WHERE orderdate > date '1995-01-01' GROUP BY clerk;

  EXPLAIN ANALYZE VERBOSE SELECT count(clerk) OVER() FROM orders
  WHERE orderdate > date '1995-01-01';
  ```
- **notes:** Executes the statement and displays the distributed execution plan with actual runtime statistics (CPU time, wall time, input/output rows/bytes per plan node). `VERBOSE` adds detailed low-level statistics (e.g., percentile distributions of CPU time, input rows, scheduled time, active drivers) that may require knowledge of Trino internals to interpret. Stats may not be entirely accurate for queries that complete quickly.

---

## prepare

- **syntax:**
  ```
  PREPARE statement_name FROM statement
  ```
- **source_url:** https://trino.io/docs/current/sql/prepare.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  PREPARE my_select1 FROM SELECT * FROM nation;

  PREPARE my_select2 FROM
  SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?;

  PREPARE my_insert FROM
  INSERT INTO cities VALUES (1, 'San Francisco');
  ```
- **notes:** Saves a query in the session under `statement_name` for later execution. Parameters use `?` placeholders. Session-scoped; disappears when the session ends or is removed with `DEALLOCATE PREPARE`.

---

## execute

- **syntax:**
  ```
  EXECUTE statement_name [ USING parameter1 [ , parameter2, ... ] ]
  ```
- **source_url:** https://trino.io/docs/current/sql/execute.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  PREPARE my_select1 FROM SELECT name FROM nation;
  EXECUTE my_select1;

  PREPARE my_select2 FROM
  SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?;
  EXECUTE my_select2 USING 1, 3;
  ```
- **notes:** Runs a previously prepared statement. Parameters are supplied positionally via `USING`. The second example is equivalent to `SELECT name FROM nation WHERE regionkey = 1 AND nationkey < 3`.

---

## execute-immediate

- **syntax:**
  ```
  EXECUTE IMMEDIATE 'statement' [ USING parameter1 [ , parameter2, ... ] ]
  ```
- **source_url:** https://trino.io/docs/current/sql/execute-immediate.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  EXECUTE IMMEDIATE 'SELECT name FROM nation';

  EXECUTE IMMEDIATE
  'SELECT name FROM nation WHERE regionkey = ? AND nationkey < ?'
  USING 1, 3;
  ```
- **notes:** Executes a SQL string directly without the need for a prior `PREPARE` or subsequent `DEALLOCATE PREPARE`. Combines the full prepare-execute-deallocate lifecycle into one statement. Parameters use `?` placeholders supplied positionally via `USING`.

---

## deallocate-prepare

- **syntax:**
  ```
  DEALLOCATE PREPARE statement_name
  ```
- **source_url:** https://trino.io/docs/current/sql/deallocate-prepare.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  DEALLOCATE PREPARE my_query;
  ```
- **notes:** Removes a named prepared statement from the current session's list. After deallocation, the statement name is no longer available for `EXECUTE` or `DESCRIBE INPUT/OUTPUT`.

---

## start-transaction

- **syntax:**
  ```
  START TRANSACTION [ mode [, ...] ]

  where mode is one of:
      ISOLATION LEVEL { READ UNCOMMITTED | READ COMMITTED | REPEATABLE READ | SERIALIZABLE }
      READ { ONLY | WRITE }
  ```
- **source_url:** https://trino.io/docs/current/sql/start-transaction.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  START TRANSACTION;
  START TRANSACTION ISOLATION LEVEL REPEATABLE READ;
  START TRANSACTION READ WRITE;
  START TRANSACTION ISOLATION LEVEL READ COMMITTED, READ ONLY;
  START TRANSACTION READ WRITE, ISOLATION LEVEL SERIALIZABLE;
  ```
- **notes:** Starts a new transaction for the current session. Multiple modes can be combined in one statement (comma-separated). All four isolation levels are syntactically supported; connector support varies.

---

## commit

- **syntax:**
  ```
  COMMIT [ WORK ]
  ```
- **source_url:** https://trino.io/docs/current/sql/commit.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  COMMIT;
  COMMIT WORK;
  ```
- **notes:** Commits the current transaction, making all changes permanent. The `WORK` keyword is optional and has no functional effect.

---

## rollback

- **syntax:**
  ```
  ROLLBACK [ WORK ]
  ```
- **source_url:** https://trino.io/docs/current/sql/rollback.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  ROLLBACK;
  ROLLBACK WORK;
  ```
- **notes:** Rolls back the current transaction, discarding all changes. The `WORK` keyword is optional and has no functional effect.

---

## call

- **syntax:**
  ```
  CALL procedure_name ( [ name => ] expression [, ...] )
  ```
- **source_url:** https://trino.io/docs/current/sql/call.html
- **version:** Trino 481
- **example_sql:**
  ```sql
  CALL test(123, 'apple');
  CALL test(name => 'apple', id => 123);
  CALL catalog.schema.test();
  ```
- **notes:** Invokes connector-provided procedures for data manipulation or administrative tasks. Both positional and named (`name => value`) argument styles are supported. Procedure name can be fully qualified as `catalog.schema.procedure`. Note: stored procedures from external databases (e.g., PostgreSQL) are not callable via `CALL`; only Trino connector-defined procedures work.
