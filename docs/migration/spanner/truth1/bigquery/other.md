# BigQuery GoogleSQL truth1 — Other Statements

Sources:
- EXPORT DATA: https://cloud.google.com/bigquery/docs/reference/standard-sql/other-statements
- ASSERT: https://cloud.google.com/bigquery/docs/reference/standard-sql/debugging-statements
- LOAD DATA: https://cloud.google.com/bigquery/docs/reference/standard-sql/load-data (page redirects/not found at time of corpus creation; syntax captured from known documentation)

---

## OTHER-001: EXPORT DATA statement

```
EXPORT DATA
[WITH CONNECTION connection_name]
OPTIONS (export_option_list) AS
query_statement

-- export_option_list required fields:
-- format = 'CSV' | 'AVRO' | 'JSON' | 'PARQUET'
-- uri = 'gs://...' | 's3://...'
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/other-statements#export_data_statement
version: bigquery-current

```sql
-- Export to Cloud Storage in CSV format:
EXPORT DATA
OPTIONS(
  uri='gs://bucket/folder/*.csv',
  format='CSV',
  overwrite=true,
  header=true,
  field_delimiter=';'
) AS
SELECT field1, field2 FROM mydataset.table1 ORDER BY field1;

-- Export to Amazon S3:
EXPORT DATA
WITH CONNECTION `myproject.us.my-connection`
OPTIONS(
  uri='s3://bucket/path/file_*.json',
  format='JSON'
) AS
SELECT * FROM mydataset.mytable;

-- Export to Cloud Storage in Avro format:
EXPORT DATA
OPTIONS(
  uri='gs://bucket/folder/*.avro',
  format='AVRO',
  use_avro_logical_types=true
) AS
SELECT * FROM mydataset.mytable;
```

---

## OTHER-002: ASSERT statement

```
ASSERT expression [AS description]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/debugging-statements#assert
version: bigquery-current

```sql
-- Succeeds if condition is TRUE:
ASSERT (
  (SELECT COUNT(*) > 5 FROM UNNEST([1, 2, 3, 4, 5, 6]))
) AS 'Table must contain more than 5 rows.';

-- Produces error if condition is FALSE:
ASSERT
  EXISTS(
    (SELECT X FROM UNNEST([7877, 7879]) AS X WHERE X = 7907))
AS 'Column X must contain the value 7907.';
```

---

## OTHER-003: LOAD DATA statement

```
LOAD DATA { OVERWRITE | INTO } table_name
[PARTITION BY partition_expression]
[CLUSTER BY cluster_col_list]
[OPTIONS(load_option_list)]
FROM FILES (load_files_option_list)
[WITH PARTITION COLUMNS [(col_name col_type[, ...])] ]
[WITH CONNECTION connection_name]

load_files_option_list (required):
  format = 'CSV' | 'AVRO' | 'JSON' | 'ORC' | 'PARQUET'
  uris = ['uri1', 'uri2', ...]
```

source_url: https://cloud.google.com/bigquery/docs/reference/standard-sql/load-data
version: bigquery-current

Note: The dedicated LOAD DATA page was not accessible at time of corpus creation (404). This syntax is extracted from referenced documentation and known BigQuery SQL behavior.

```sql
LOAD DATA INTO mydataset.mytable
FROM FILES (
  format = 'CSV',
  uris = ['gs://mybucket/myfile.csv']
);

LOAD DATA OVERWRITE mydataset.mytable
OPTIONS(description="Refreshed table")
FROM FILES (
  format = 'PARQUET',
  uris = ['gs://mybucket/data*.parquet']
);

-- With explicit partition columns:
LOAD DATA INTO mydataset.mytable
FROM FILES (
  format = 'CSV',
  uris = ['gs://mybucket/year=*/month=*/data.csv'],
  hive_partition_uri_prefix = 'gs://mybucket/'
)
WITH PARTITION COLUMNS (year INT64, month INT64);
```
