-- SHOW EXTERNAL TABLE command test cases for Amazon Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_SHOW_EXTERNAL_TABLE.html

-- Basic SHOW EXTERNAL TABLE command
SHOW EXTERNAL TABLE spectrum.sales;

-- With schema and table name
SHOW EXTERNAL TABLE my_schema.alldatatypes_parquet_test_partitioned;

-- With database, schema, and table name
SHOW EXTERNAL TABLE my_database.my_schema.alldatatypes_parquet_test_partitioned;

-- With PARTITION option
SHOW EXTERNAL TABLE my_schema.alldatatypes_parquet_test_partitioned PARTITION;

-- With database and PARTITION option
SHOW EXTERNAL TABLE my_database.my_schema.sales_data PARTITION;

-- Various external table names
SHOW EXTERNAL TABLE spectrum.customer_data;
SHOW EXTERNAL TABLE external_schema.transactions;
SHOW EXTERNAL TABLE analytics_db.external_schema.user_events;

-- With quoted identifiers
SHOW EXTERNAL TABLE "My_Schema"."My_Table";
SHOW EXTERNAL TABLE my_db."External_Schema"."Table_Name";

-- Mixed case without quotes
SHOW EXTERNAL TABLE MyDatabase.MySchema.MyTable;

-- With underscores and numbers in names
SHOW EXTERNAL TABLE db_2024.schema_01.table_v2;
SHOW EXTERNAL TABLE external_db.ext_schema.sales_2024_q1;

-- Test with PARTITION keyword
SHOW EXTERNAL TABLE spectrum.sales_partitioned PARTITION;
SHOW EXTERNAL TABLE analytics.external_data.events_partitioned PARTITION;

-- Real-world examples from documentation
SHOW EXTERNAL TABLE spectrum.alldatatypes_parquet_test_partitioned;
SHOW EXTERNAL TABLE spectrum.alldatatypes_parquet_test_partitioned PARTITION;

-- Edge cases with special naming
SHOW EXTERNAL TABLE "database-name"."schema-name"."table-name";
SHOW EXTERNAL TABLE db123.schema456.table789;

-- Multiple levels of qualification
SHOW EXTERNAL TABLE prod_db.spectrum_schema.fact_sales;
SHOW EXTERNAL TABLE test_db.external_s3.dim_customer PARTITION;

-- Various partition table examples
SHOW EXTERNAL TABLE sales_db.spectrum.sales_by_region PARTITION;
SHOW EXTERNAL TABLE marketing.external_s3.campaigns_by_date PARTITION;
SHOW EXTERNAL TABLE finance.spectrum_tables.transactions_by_month PARTITION;

-- Comments to test parser's comment handling
-- This is a comment
SHOW EXTERNAL TABLE my_schema.my_table; -- inline comment

/* Block comment test */
SHOW EXTERNAL TABLE db.schema.table_name;

/* Multi-line
   block comment */
SHOW EXTERNAL TABLE external_database.external_schema.external_table PARTITION;