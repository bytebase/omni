# BigQuery GoogleSQL truth1 — DDL Syntax

All forms extracted from:
https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language

---

## DDL-001: CREATE SCHEMA

```
CREATE SCHEMA [ IF NOT EXISTS ]
[project_name.]dataset_name
[DEFAULT COLLATE collate_specification]
[OPTIONS(schema_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_schema_statement
version: bigquery-current

```sql
CREATE SCHEMA mydataset
OPTIONS(
  location="us",
  default_table_expiration_days=3.75,
  labels=[("label1","value1"),("label2","value2")]
);

CREATE SCHEMA mydataset DEFAULT COLLATE 'und:ci';
```

---

## DDL-002: CREATE TABLE

```
CREATE [ OR REPLACE ] [ TEMP | TEMPORARY ] TABLE [ IF NOT EXISTS ]
table_name
[(
  column | constraint_definition[, ...]
)]
[DEFAULT COLLATE collate_specification]
[PARTITION BY partition_expression]
[CLUSTER BY clustering_column_list]
[WITH CONNECTION connection_name]
[OPTIONS(table_option_list)]
[AS query_statement]

column := column_definition

constraint_definition :=
  [primary_key]
  | [[CONSTRAINT constraint_name] foreign_key, ...]

primary_key :=
  PRIMARY KEY (column_name[, ...]) NOT ENFORCED

foreign_key :=
  FOREIGN KEY (column_name[, ...]) foreign_reference

foreign_reference :=
  REFERENCES primary_key_table(column_name[, ...]) NOT ENFORCED
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_table_statement
version: bigquery-current

```sql
CREATE TABLE mydataset.newtable (
  x INT64,
  y STRING NOT NULL,
  z ARRAY<FLOAT64>
)
PARTITION BY _PARTITIONDATE
CLUSTER BY x
OPTIONS(expiration_timestamp=TIMESTAMP "2025-01-01 00:00:00 UTC");

CREATE OR REPLACE TABLE `myproject.mydataset.new_table` AS
  WITH RECURSIVE T1 AS (SELECT 1 AS n UNION ALL SELECT n + 1 FROM T1 WHERE n < 3)
  SELECT * FROM T1;

CREATE TEMP TABLE tmp_table (n INT64);
```

---

## DDL-003: column_schema (column definition)

```
column :=
  column_name column_schema

column_schema :=
   {
     simple_type
     | STRUCT<field_list>
     | ARRAY<array_element_schema>
   }
   [PRIMARY KEY NOT ENFORCED | REFERENCES table_name(column_name) NOT ENFORCED]
   [ DEFAULT default_expression |
     GENERATED ALWAYS AS (generation_expression) STORED OPTIONS(generation_option_list) ]
   [NOT NULL]
   [OPTIONS(column_option_list)]

simple_type :=
  { data_type | STRING COLLATE collate_specification }

field_list :=
  field_name column_schema [, ...]

array_element_schema :=
  { simple_type | STRUCT<field_list> }
  [NOT NULL]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#column
version: bigquery-current

```sql
CREATE TABLE t (
  id INT64 NOT NULL,
  name STRING COLLATE 'und:ci',
  tags ARRAY<STRING>,
  info STRUCT<a INT64, b STRING>,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP()
);
```

---

## DDL-004: PARTITION BY (partition expression)

```
PARTITION BY partition_expression

partition_expression options:
  _PARTITIONDATE
  DATE(_PARTITIONTIME)
  <date_column>
  DATE({ <timestamp_column> | <datetime_column> })
  DATETIME_TRUNC(<datetime_column>, { DAY | HOUR | MONTH | YEAR })
  TIMESTAMP_TRUNC(<timestamp_column>, { DAY | HOUR | MONTH | YEAR })
  TIMESTAMP_TRUNC(_PARTITIONTIME, { DAY | HOUR | MONTH | YEAR })
  DATE_TRUNC(<date_column>, { MONTH | YEAR })
  RANGE_BUCKET(<int64_column>, GENERATE_ARRAY(<start>, <end>[, <interval>]))
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#partition_expression
version: bigquery-current

```sql
CREATE TABLE t (ts TIMESTAMP, val INT64)
PARTITION BY TIMESTAMP_TRUNC(ts, DAY);

CREATE TABLE t (dt DATE, val INT64)
PARTITION BY DATE_TRUNC(dt, MONTH);

CREATE TABLE t (id INT64, val STRING)
PARTITION BY RANGE_BUCKET(id, GENERATE_ARRAY(0, 100, 10));
```

---

## DDL-005: CLUSTER BY

```
CLUSTER BY clustering_column_list

clustering_column_list: column_reference [, ...] (up to 4 columns)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#clustering_column_list
version: bigquery-current

```sql
CREATE TABLE t (a INT64, b STRING, c DATE)
PARTITION BY DATE_TRUNC(c, MONTH)
CLUSTER BY a, b;
```

---

## DDL-006: CREATE TABLE LIKE

```
CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ]
table_name
LIKE [[project_name.]dataset_name.]source_table_name
...
[OPTIONS(table_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_table_like_statement
version: bigquery-current

```sql
CREATE TABLE mydataset.newtable LIKE mydataset.sourcetable;

CREATE TABLE mydataset.newtable
LIKE mydataset.sourcetable
AS SELECT * FROM mydataset.myothertable;
```

---

## DDL-007: CREATE TABLE COPY

```
CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ] table_name
COPY source_table_name
...
[OPTIONS(table_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_table_copy_statement
version: bigquery-current

```sql
CREATE TABLE mydataset.newtable COPY mydataset.sourcetable;
```

---

## DDL-008: CREATE SNAPSHOT TABLE

```
CREATE SNAPSHOT TABLE [ IF NOT EXISTS ] table_snapshot_name
CLONE source_table_name
[FOR SYSTEM_TIME AS OF time_expression]
[OPTIONS(snapshot_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_snapshot_table_statement
version: bigquery-current

```sql
CREATE SNAPSHOT TABLE mydataset.my_snapshot
CLONE mydataset.mytable
FOR SYSTEM_TIME AS OF TIMESTAMP_SUB(CURRENT_TIMESTAMP(), INTERVAL 1 HOUR);
```

---

## DDL-009: CREATE TABLE CLONE

```
CREATE [ OR REPLACE ] TABLE [ IF NOT EXISTS ] table_name
CLONE source_table_name
[FOR SYSTEM_TIME AS OF time_expression]
[OPTIONS(table_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_table_clone_statement
version: bigquery-current

```sql
CREATE TABLE mydataset.myclone CLONE mydataset.mytable;
```

---

## DDL-010: CREATE VIEW

```
CREATE [ OR REPLACE ] VIEW [ IF NOT EXISTS ] view_name
[(view_column_name_list)]
[OPTIONS(view_option_list)]
AS query_expression

view_column_name_list :=
  view_column[, ...]

view_column :=
  column_name [OPTIONS(view_column_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_view_statement
version: bigquery-current

```sql
CREATE VIEW mydataset.myview AS SELECT a, b FROM mydataset.mytable WHERE c > 0;

CREATE OR REPLACE VIEW mydataset.age_groups(age, count) AS
  SELECT age, COUNT(*) FROM mydataset.people GROUP BY age;
```

---

## DDL-011: CREATE MATERIALIZED VIEW

```
CREATE [ OR REPLACE ] MATERIALIZED VIEW [ IF NOT EXISTS ] materialized_view_name
[PARTITION BY partition_expression]
[CLUSTER BY clustering_column_list]
[OPTIONS(materialized_view_option_list)]
AS query_expression
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_materialized_view_statement
version: bigquery-current

```sql
CREATE MATERIALIZED VIEW myProject.myDataset.myView
AS SELECT product, SUM(sales) AS total FROM mydataset.sales GROUP BY product;

CREATE MATERIALIZED VIEW myProject.myDataset.myView
PARTITION BY DATE(ts)
CLUSTER BY product
AS SELECT product, ts, SUM(sales) AS total FROM mydataset.sales GROUP BY product, ts;
```

---

## DDL-012: CREATE MATERIALIZED VIEW AS REPLICA OF

```
CREATE MATERIALIZED VIEW [ IF NOT EXISTS ] materialized_view_name
AS REPLICA OF source_materialized_view_name
[OPTIONS(materialized_view_replica_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_materialized_view_as_replica_of_statement
version: bigquery-current

```sql
CREATE MATERIALIZED VIEW mydataset.my_replica AS REPLICA OF other_project.other_dataset.my_mv;
```

---

## DDL-013: CREATE EXTERNAL SCHEMA

```
CREATE EXTERNAL SCHEMA [ IF NOT EXISTS ]
[project_name.]dataset_name
WITH CONNECTION connection_name
OPTIONS(external_schema_option_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_external_schema_statement
version: bigquery-current

```sql
CREATE EXTERNAL SCHEMA mydataset
WITH CONNECTION myproject.us.my_connection
OPTIONS(location="s3://mybucket/data/");
```

---

## DDL-014: CREATE EXTERNAL TABLE

```
CREATE [ OR REPLACE ] EXTERNAL TABLE [ IF NOT EXISTS ] table_name
[(
  column[, ...]
)]
[DEFAULT COLLATE collate_specification]
[WITH CONNECTION connection_name]
[WITH PARTITION COLUMNS [(column[, ...])]]
OPTIONS(external_table_option_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_external_table_statement
version: bigquery-current

```sql
CREATE OR REPLACE EXTERNAL TABLE mydataset.myexternal
OPTIONS(
  format = 'CSV',
  uris = ['gs://mybucket/myfile.csv']
);
```

---

## DDL-015: CREATE FUNCTION (SQL)

```
CREATE [ OR REPLACE ] [ TEMPORARY | TEMP ] FUNCTION [ IF NOT EXISTS ]
    [[project_name.]dataset_name.]function_name
    ([named_parameter[, ...]])
  [RETURNS data_type]
  AS (sql_expression)
  [OPTIONS (function_option_list)]

named_parameter:
  param_name param_type
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_function_statement
version: bigquery-current

```sql
CREATE FUNCTION mydataset.square(x FLOAT64) RETURNS FLOAT64 AS (x * x);

CREATE OR REPLACE TEMP FUNCTION add_n(x INT64, n INT64) AS (x + n);
```

---

## DDL-016: CREATE FUNCTION (JavaScript)

```
CREATE [OR REPLACE] [TEMPORARY | TEMP] FUNCTION [IF NOT EXISTS]
    [[project_name.]dataset_name.]function_name
    ([named_parameter[, ...]])
  RETURNS data_type
  [{ DETERMINISTIC | NOT DETERMINISTIC }]
  LANGUAGE js
  [OPTIONS (function_option_list)]
  AS javascript_code

named_parameter:
  param_name param_type
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_function_statement
version: bigquery-current

```sql
CREATE OR REPLACE FUNCTION mydataset.multiplyInputs(x FLOAT64, y FLOAT64)
RETURNS FLOAT64
LANGUAGE js
AS r"""
  return x*y;
""";
```

---

## DDL-017: CREATE FUNCTION (Python / remote)

```
-- Python UDF:
CREATE [OR REPLACE] FUNCTION [IF NOT EXISTS]
    [project_name.]dataset_name.function_name
    ([named_parameter[, ...]])
  RETURNS data_type
  LANGUAGE python
  [WITH CONNECTION connection_name]
  OPTIONS (function_option_list)
  AS python_code

-- Remote function:
CREATE [OR REPLACE] [TEMPORARY | TEMP] FUNCTION [IF NOT EXISTS]
    [[project_name.]dataset_name.]function_name
    ([named_parameter[, ...]])
  RETURNS data_type
  REMOTE WITH CONNECTION connection_name
  [OPTIONS (function_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_function_statement
version: bigquery-current

```sql
CREATE FUNCTION mydataset.py_square(x FLOAT64)
RETURNS FLOAT64
LANGUAGE python
WITH CONNECTION `my-project.us.my-connection`
OPTIONS (entry_point='square', runtime_version='python-3.11')
AS r"""
def square(x):
  return x * x
""";
```

---

## DDL-018: CREATE AGGREGATE FUNCTION (SQL)

```
CREATE [ OR REPLACE ] AGGREGATE FUNCTION [ IF NOT EXISTS ]
    [[project_name.]dataset_name.]function_name
    ([named_parameter[, ...]])
  [RETURNS data_type]
  AS (sql_expression)
  [OPTIONS (function_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_aggregate_function_statement_sql
version: bigquery-current

```sql
CREATE AGGREGATE FUNCTION mydataset.sum_of_squares(x FLOAT64)
RETURNS FLOAT64
AS (SUM(x * x));
```

---

## DDL-019: CREATE TABLE FUNCTION

```
CREATE [ OR REPLACE ] TABLE FUNCTION [ IF NOT EXISTS ]
  [[project_name.]dataset_name.]function_name
  ( [ function_parameter [, ...] ] )
  [RETURNS TABLE < column_declaration [, ...] > ]
  [OPTIONS (table_function_options_list) ]
  AS sql_query

function_parameter:
  parameter_name { data_type | ANY TYPE | TABLE < column_declaration [, ...] > }

column_declaration:
  column_name data_type
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_table_function_statement
version: bigquery-current

```sql
CREATE OR REPLACE TABLE FUNCTION mydataset.names_by_year(y INT64)
AS
  SELECT year, name, SUM(number) AS total
  FROM `bigquery-public-data.usa_names.usa_1910_current`
  WHERE year = y
  GROUP BY year, name;

CREATE OR REPLACE TABLE FUNCTION mydataset.names_by_year(y INT64)
RETURNS TABLE<name STRING, year INT64, total INT64>
AS
  SELECT year, name, SUM(number) AS total
  FROM `bigquery-public-data.usa_names.usa_1910_current`
  WHERE year = y
  GROUP BY year, name;
```

---

## DDL-020: CREATE PROCEDURE

```
CREATE [OR REPLACE] PROCEDURE [IF NOT EXISTS]
[[project_name.]dataset_name.]procedure_name (procedure_argument[, ...] )
[OPTIONS(procedure_option_list)]
BEGIN
multi_statement_query
END;

procedure_argument: [procedure_argument_mode] argument_name argument_type

procedure_argument_mode: IN | OUT | INOUT
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_procedure_statement
version: bigquery-current

```sql
CREATE OR REPLACE PROCEDURE mydataset.SelectFromTable(IN tbl STRING, OUT result INT64)
BEGIN
  EXECUTE IMMEDIATE 'SELECT COUNT(*) FROM ' || tbl INTO result;
END;
```

---

## DDL-021: CREATE ROW ACCESS POLICY

```
CREATE [ OR REPLACE ] ROW ACCESS POLICY [ IF NOT EXISTS ]
row_access_policy_name ON table_name
[GRANT TO (grantee_list)]
FILTER USING (filter_expression);
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_row_access_policy_statement
version: bigquery-current

```sql
CREATE ROW ACCESS POLICY my_policy
ON mydataset.mytable
GRANT TO ("user:alice@example.com")
FILTER USING (department = SESSION_USER());
```

---

## DDL-022: CREATE SEARCH INDEX

```
CREATE SEARCH INDEX [ IF NOT EXISTS ] index_name
ON table_name({ALL COLUMNS [WITH COLUMN OPTIONS(column [, ...])] | column [, ...]})
[OPTIONS(index_option_list)]

column := column_name [OPTIONS(index_column_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_search_index_statement
version: bigquery-current

```sql
CREATE SEARCH INDEX my_index ON mydataset.mytable(ALL COLUMNS);
CREATE SEARCH INDEX my_index ON mydataset.mytable(name, description);
```

---

## DDL-023: CREATE VECTOR INDEX

```
CREATE [ OR REPLACE ] VECTOR INDEX [ IF NOT EXISTS ] index_name
ON table_name(column_name)
[STORING(stored_column_name [, ...])]
[PARTITION BY partition_expression]
OPTIONS(index_option_list);
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_vector_index_statement
version: bigquery-current

```sql
CREATE VECTOR INDEX my_index ON mydataset.mytable(embedding_column)
OPTIONS(index_type='TREE_AH', distance_type='COSINE');
```

---

## DDL-024: CREATE CAPACITY

```
CREATE CAPACITY
`project_id.location_id.commitment_id`
OPTIONS (capacity_commitment_option_list);
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_capacity_statement
version: bigquery-current

```sql
CREATE CAPACITY `admin_project.region-us.my-commitment`
OPTIONS (slot_count = 100, plan = 'ANNUAL');
```

---

## DDL-025: CREATE RESERVATION

```
CREATE RESERVATION
`project_id.location_id.reservation_id`
OPTIONS (reservation_option_list);
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_reservation_statement
version: bigquery-current

```sql
CREATE RESERVATION `admin_project.region-us.prod`
OPTIONS (slot_capacity = 100);
```

---

## DDL-026: CREATE ASSIGNMENT

```
CREATE ASSIGNMENT
`project_id.location_id.reservation_id.assignment_id`
OPTIONS (assignment_option_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#create_assignment_statement
version: bigquery-current

```sql
CREATE ASSIGNMENT `admin_project.region-us.prod.my-assignment`
OPTIONS (assignee = 'projects/my-project', job_type = 'QUERY');
```

---

## DDL-027: ALTER TABLE SET OPTIONS

```
ALTER TABLE [IF EXISTS] table_name
SET OPTIONS(table_set_options_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_set_options_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable SET OPTIONS(expiration_timestamp=TIMESTAMP "2025-01-01 00:00:00 UTC");
```

---

## DDL-028: ALTER TABLE ADD COLUMN

```
ALTER TABLE [IF EXISTS] table_name
ADD COLUMN [IF NOT EXISTS] column_definition
[, ADD COLUMN [IF NOT EXISTS] column_definition ...]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_add_column_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable ADD COLUMN new_col INT64;
ALTER TABLE mydataset.mytable ADD COLUMN IF NOT EXISTS new_col STRING;
```

---

## DDL-029: ALTER TABLE ADD FOREIGN KEY

```
ALTER TABLE [IF EXISTS] table_name
ADD [CONSTRAINT constraint_name] FOREIGN KEY (column_name[, ...])
REFERENCES ref_table(ref_column[, ...]) NOT ENFORCED
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_add_foreign_key_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.child_table
ADD CONSTRAINT fk1 FOREIGN KEY (parent_id) REFERENCES mydataset.parent_table(id) NOT ENFORCED;
```

---

## DDL-030: ALTER TABLE ADD PRIMARY KEY

```
ALTER TABLE [IF EXISTS] table_name
ADD PRIMARY KEY (column_name[, ...]) NOT ENFORCED
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_add_primary_key_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable ADD PRIMARY KEY (id) NOT ENFORCED;
```

---

## DDL-031: ALTER TABLE RENAME TO

```
ALTER TABLE [IF EXISTS] table_name
RENAME TO new_table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_rename_to_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.old_name RENAME TO new_name;
```

---

## DDL-032: ALTER TABLE RENAME COLUMN

```
ALTER TABLE [IF EXISTS] table_name
RENAME COLUMN [IF EXISTS] column_name TO new_column_name
[, RENAME COLUMN [IF EXISTS] column_name TO new_column_name ...]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_rename_column_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable RENAME COLUMN old_col TO new_col;
```

---

## DDL-033: ALTER TABLE DROP COLUMN

```
ALTER TABLE [IF EXISTS] table_name
DROP COLUMN [IF EXISTS] column_name
[, DROP COLUMN [IF EXISTS] column_name ...]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_drop_column_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable DROP COLUMN IF EXISTS old_col;
```

---

## DDL-034: ALTER TABLE DROP CONSTRAINT / DROP PRIMARY KEY

```
-- Drop named constraint:
ALTER TABLE [IF EXISTS] table_name
DROP CONSTRAINT [IF EXISTS] constraint_name

-- Drop primary key:
ALTER TABLE [IF EXISTS] table_name
DROP PRIMARY KEY [IF EXISTS]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_drop_constraint_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable DROP CONSTRAINT IF EXISTS fk1;
ALTER TABLE mydataset.mytable DROP PRIMARY KEY;
```

---

## DDL-035: ALTER TABLE SET DEFAULT COLLATE

```
ALTER TABLE [IF EXISTS] table_name
SET DEFAULT COLLATE collate_specification
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_table_set_default_collate_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable SET DEFAULT COLLATE 'und:ci';
```

---

## DDL-036: ALTER COLUMN SET OPTIONS / DROP NOT NULL / SET DATA TYPE / SET DEFAULT / DROP DEFAULT

```
ALTER TABLE [IF EXISTS] table_name
ALTER COLUMN [IF EXISTS] column_name { SET OPTIONS(col_option_list)
                                      | DROP NOT NULL
                                      | SET DATA TYPE new_type
                                      | SET DEFAULT expr
                                      | DROP DEFAULT }
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_column_set_options_statement
version: bigquery-current

```sql
ALTER TABLE mydataset.mytable ALTER COLUMN mycol SET OPTIONS(description="new description");
ALTER TABLE mydataset.mytable ALTER COLUMN mycol DROP NOT NULL;
ALTER TABLE mydataset.mytable ALTER COLUMN mycol SET DATA TYPE STRING;
ALTER TABLE mydataset.mytable ALTER COLUMN mycol SET DEFAULT 'unknown';
ALTER TABLE mydataset.mytable ALTER COLUMN mycol DROP DEFAULT;
```

---

## DDL-037: ALTER SCHEMA SET DEFAULT COLLATE

```
ALTER SCHEMA [IF EXISTS] [project_name.]dataset_name
SET DEFAULT COLLATE collate_specification
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_schema_set_default_collate_statement
version: bigquery-current

```sql
ALTER SCHEMA mydataset SET DEFAULT COLLATE 'und:ci';
```

---

## DDL-038: ALTER SCHEMA SET OPTIONS

```
ALTER SCHEMA [IF EXISTS] [project_name.]dataset_name
SET OPTIONS(schema_set_options_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_schema_set_options_statement
version: bigquery-current

```sql
ALTER SCHEMA mydataset SET OPTIONS(default_table_expiration_days=7);
```

---

## DDL-039: ALTER VIEW SET OPTIONS

```
ALTER VIEW [IF EXISTS] view_name
SET OPTIONS(view_set_options_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_view_set_options_statement
version: bigquery-current

```sql
ALTER VIEW mydataset.myview SET OPTIONS(expiration_timestamp=TIMESTAMP "2025-01-01 00:00:00 UTC");
```

---

## DDL-040: ALTER MATERIALIZED VIEW SET OPTIONS

```
ALTER MATERIALIZED VIEW [IF EXISTS] materialized_view_name
SET OPTIONS(mv_set_options_list)
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_materialized_view_set_options_statement
version: bigquery-current

```sql
ALTER MATERIALIZED VIEW mydataset.myview SET OPTIONS(enable_refresh=false);
```

---

## DDL-041: ALTER VECTOR INDEX REBUILD

```
ALTER VECTOR INDEX [IF EXISTS] index_name ON table_name REBUILD
[OPTIONS(index_rebuild_option_list)]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#alter_vector_index_rebuild_statement
version: bigquery-current

```sql
ALTER VECTOR INDEX my_index ON mydataset.mytable REBUILD;
```

---

## DDL-042: DROP SCHEMA

```
DROP [EXTERNAL] SCHEMA [IF EXISTS]
[project_name.]dataset_name
[ CASCADE | RESTRICT ]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_schema_statement
version: bigquery-current

```sql
DROP SCHEMA mydataset;
DROP SCHEMA IF EXISTS mydataset CASCADE;
```

---

## DDL-043: UNDROP SCHEMA

```
UNDROP SCHEMA [IF NOT EXISTS]
[project_name.]dataset_name
[OPTIONS (location="us")]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#undrop_schema_statement
version: bigquery-current

```sql
UNDROP SCHEMA mydataset;
```

---

## DDL-044: DROP TABLE

```
DROP TABLE [IF EXISTS] table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_table_statement
version: bigquery-current

```sql
DROP TABLE mydataset.mytable;
DROP TABLE IF EXISTS mydataset.mytable;
```

---

## DDL-045: DROP SNAPSHOT TABLE

```
DROP SNAPSHOT TABLE [IF EXISTS] table_snapshot_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_snapshot_table_statement
version: bigquery-current

```sql
DROP SNAPSHOT TABLE mydataset.mytablesnapshot;
DROP SNAPSHOT TABLE IF EXISTS mydataset.mytablesnapshot;
```

---

## DDL-046: DROP EXTERNAL TABLE

```
DROP EXTERNAL TABLE [IF EXISTS] table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_external_table_statement
version: bigquery-current

```sql
DROP EXTERNAL TABLE IF EXISTS mydataset.myexternal;
```

---

## DDL-047: DROP VIEW

```
DROP VIEW [IF EXISTS] view_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_view_statement
version: bigquery-current

```sql
DROP VIEW mydataset.myview;
DROP VIEW IF EXISTS mydataset.myview;
```

---

## DDL-048: DROP MATERIALIZED VIEW

```
DROP MATERIALIZED VIEW [IF EXISTS] materialized_view_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_materialized_view_statement
version: bigquery-current

```sql
DROP MATERIALIZED VIEW IF EXISTS mydataset.my_mv;
```

---

## DDL-049: DROP FUNCTION

```
DROP FUNCTION [IF EXISTS] [[project_name.]dataset_name.]function_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_function_statement
version: bigquery-current

```sql
DROP FUNCTION IF EXISTS mydataset.my_function;
```

---

## DDL-050: DROP TABLE FUNCTION

```
DROP TABLE FUNCTION [IF EXISTS] [[project_name.]dataset_name.]function_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_table_function_statement
version: bigquery-current

```sql
DROP TABLE FUNCTION IF EXISTS mydataset.my_tvf;
```

---

## DDL-051: DROP PROCEDURE

```
DROP PROCEDURE [IF EXISTS] [[project_name.]dataset_name.]procedure_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_procedure_statement
version: bigquery-current

```sql
DROP PROCEDURE IF EXISTS mydataset.my_procedure;
```

---

## DDL-052: DROP ROW ACCESS POLICY

```
DROP ROW ACCESS POLICY [IF EXISTS] policy_name ON table_name
DROP ALL ROW ACCESS POLICIES ON table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_row_access_policy_statement
version: bigquery-current

```sql
DROP ROW ACCESS POLICY IF EXISTS my_policy ON mydataset.mytable;
DROP ALL ROW ACCESS POLICIES ON mydataset.mytable;
```

---

## DDL-053: DROP CAPACITY / RESERVATION / ASSIGNMENT

```
DROP CAPACITY [IF EXISTS] `project_id.location_id.commitment_id`
DROP RESERVATION [IF EXISTS] `project_id.location_id.reservation_id`
DROP ASSIGNMENT [IF EXISTS] `project_id.location_id.reservation_id.assignment_id`
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_capacity_statement
version: bigquery-current

```sql
DROP CAPACITY IF EXISTS `admin_project.region-us.my-commitment`;
DROP RESERVATION IF EXISTS `admin_project.region-us.prod`;
DROP ASSIGNMENT IF EXISTS `admin_project.region-us.prod.my-assignment`;
```

---

## DDL-054: DROP SEARCH INDEX / DROP VECTOR INDEX

```
DROP SEARCH INDEX [IF EXISTS] index_name ON table_name
DROP VECTOR INDEX [IF EXISTS] index_name ON table_name
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/data-definition-language#drop_search_index_statement
version: bigquery-current

```sql
DROP SEARCH INDEX IF EXISTS my_index ON mydataset.mytable;
DROP VECTOR INDEX IF EXISTS my_index ON mydataset.mytable;
```
