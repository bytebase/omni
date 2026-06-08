-- UNLOAD command test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_UNLOAD.html

-- Basic UNLOAD with minimal options
UNLOAD ('SELECT * FROM sales')
TO 's3://mybucket/myprefix'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole';

-- UNLOAD with DEFAULT IAM role
UNLOAD ('SELECT * FROM customers')
TO 's3://mybucket/customers/'
IAM_ROLE default;

-- UNLOAD with CSV format
UNLOAD ('SELECT * FROM products')
TO 's3://mybucket/products/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT CSV;

-- UNLOAD with PARQUET format
UNLOAD ('SELECT * FROM orders')
TO 's3://mybucket/orders/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT PARQUET;

-- UNLOAD with JSON format
UNLOAD ('SELECT * FROM events')
TO 's3://mybucket/events/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT JSON;

-- UNLOAD with PARTITION BY
UNLOAD ('SELECT * FROM sales_data')
TO 's3://mybucket/partitioned/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
PARTITION BY (year, month);

-- UNLOAD with PARTITION BY and INCLUDE
UNLOAD ('SELECT * FROM sales_data')
TO 's3://mybucket/partitioned_include/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
PARTITION BY (region, category) INCLUDE;

-- UNLOAD with MANIFEST
UNLOAD ('SELECT * FROM inventory')
TO 's3://mybucket/inventory/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MANIFEST;

-- UNLOAD with MANIFEST VERBOSE
UNLOAD ('SELECT * FROM inventory')
TO 's3://mybucket/inventory_verbose/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MANIFEST VERBOSE;

-- UNLOAD with HEADER (for CSV)
UNLOAD ('SELECT * FROM employees')
TO 's3://mybucket/employees/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT CSV
HEADER;

-- UNLOAD with custom DELIMITER
UNLOAD ('SELECT * FROM data')
TO 's3://mybucket/data/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER AS '|';

-- UNLOAD with DELIMITER without AS
UNLOAD ('SELECT * FROM data')
TO 's3://mybucket/data_tab/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '\t';

-- UNLOAD with FIXEDWIDTH
UNLOAD ('SELECT * FROM fixed_data')
TO 's3://mybucket/fixed/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FIXEDWIDTH 'name:30,age:3,city:20';

-- UNLOAD with ENCRYPTED
UNLOAD ('SELECT * FROM sensitive_data')
TO 's3://mybucket/encrypted/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ENCRYPTED;

-- UNLOAD with ENCRYPTED AUTO
UNLOAD ('SELECT * FROM sensitive_data')
TO 's3://mybucket/encrypted_auto/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ENCRYPTED AUTO;

-- UNLOAD with KMS_KEY_ID
UNLOAD ('SELECT * FROM confidential')
TO 's3://mybucket/kms_encrypted/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ENCRYPTED
KMS_KEY_ID '1234abcd-12ab-34cd-56ef-1234567890ab';

-- UNLOAD with BZIP2 compression
UNLOAD ('SELECT * FROM large_table')
TO 's3://mybucket/compressed/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
BZIP2;

-- UNLOAD with GZIP compression
UNLOAD ('SELECT * FROM large_table')
TO 's3://mybucket/gzipped/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
GZIP;

-- UNLOAD with ZSTD compression
UNLOAD ('SELECT * FROM large_table')
TO 's3://mybucket/zstd_compressed/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ZSTD;

-- UNLOAD with ADDQUOTES
UNLOAD ('SELECT * FROM text_data')
TO 's3://mybucket/quoted/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ADDQUOTES;

-- UNLOAD with NULL AS
UNLOAD ('SELECT * FROM nullable_data')
TO 's3://mybucket/nulls/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
NULL AS '\\N';

-- UNLOAD with ESCAPE
UNLOAD ('SELECT * FROM special_chars')
TO 's3://mybucket/escaped/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ESCAPE;

-- UNLOAD with ALLOWOVERWRITE
UNLOAD ('SELECT * FROM daily_data')
TO 's3://mybucket/daily/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ALLOWOVERWRITE;

-- UNLOAD with CLEANPATH
UNLOAD ('SELECT * FROM temp_data')
TO 's3://mybucket/clean/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
CLEANPATH;

-- UNLOAD with PARALLEL ON
UNLOAD ('SELECT * FROM large_dataset')
TO 's3://mybucket/parallel/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
PARALLEL ON;

-- UNLOAD with PARALLEL OFF
UNLOAD ('SELECT * FROM small_dataset')
TO 's3://mybucket/sequential/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
PARALLEL OFF;

-- UNLOAD with MAXFILESIZE in MB
UNLOAD ('SELECT * FROM huge_table')
TO 's3://mybucket/sized/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MAXFILESIZE 100 MB;

-- UNLOAD with MAXFILESIZE AS and GB
UNLOAD ('SELECT * FROM huge_table')
TO 's3://mybucket/sized_gb/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MAXFILESIZE AS 1 GB;

-- UNLOAD with MAXFILESIZE without unit (defaults to MB)
UNLOAD ('SELECT * FROM huge_table')
TO 's3://mybucket/sized_default/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MAXFILESIZE 500;

-- UNLOAD with ROWGROUPSIZE for Parquet
UNLOAD ('SELECT * FROM parquet_data')
TO 's3://mybucket/parquet_rg/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT PARQUET
ROWGROUPSIZE 128 MB;

-- UNLOAD with ROWGROUPSIZE AS
UNLOAD ('SELECT * FROM parquet_data')
TO 's3://mybucket/parquet_rg_as/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT PARQUET
ROWGROUPSIZE AS 256 MB;

-- UNLOAD with REGION
UNLOAD ('SELECT * FROM regional_data')
TO 's3://mybucket/regional/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
REGION 'us-west-2';

-- UNLOAD with REGION AS
UNLOAD ('SELECT * FROM regional_data')
TO 's3://mybucket/regional_as/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
REGION AS 'eu-west-1';

-- UNLOAD with EXTENSION
UNLOAD ('SELECT * FROM custom_data')
TO 's3://mybucket/custom/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
EXTENSION 'txt';

-- Complex UNLOAD with multiple options
UNLOAD ('SELECT customer_id, order_date, total_amount FROM orders WHERE order_date >= ''2023-01-01''')
TO 's3://mybucket/orders/2023/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT CSV
HEADER
DELIMITER AS ','
ADDQUOTES
NULL AS ''
ESCAPE
PARALLEL ON
MAXFILESIZE 100 MB
GZIP
MANIFEST
ALLOWOVERWRITE;

-- UNLOAD with Parquet format and multiple options
UNLOAD ('SELECT * FROM fact_sales')
TO 's3://mybucket/fact_sales/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT PARQUET
PARTITION BY (year, month, day) INCLUDE
ROWGROUPSIZE 512 MB
MAXFILESIZE 1 GB
ENCRYPTED AUTO
MANIFEST VERBOSE
CLEANPATH;

-- UNLOAD with JSON format and options
UNLOAD ('SELECT event_id, event_type, event_data, timestamp FROM events')
TO 's3://mybucket/events/json/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT JSON
MAXFILESIZE 250 MB
ZSTD
MANIFEST
REGION 'ap-southeast-1'
EXTENSION 'json.zst';

-- UNLOAD with complex SELECT statement
UNLOAD ('
    WITH monthly_sales AS (
        SELECT
            DATE_TRUNC(''month'', order_date) as month,
            SUM(total_amount) as total_sales,
            COUNT(*) as order_count
        FROM orders
        WHERE order_date >= ''2023-01-01''
        GROUP BY 1
    )
    SELECT * FROM monthly_sales
    ORDER BY month
')
TO 's3://mybucket/monthly_reports/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT CSV
HEADER
PARALLEL OFF;

-- UNLOAD with JOIN query
UNLOAD ('
    SELECT
        c.customer_id,
        c.customer_name,
        o.order_id,
        o.order_date,
        o.total_amount
    FROM customers c
    JOIN orders o ON c.customer_id = o.customer_id
    WHERE o.order_date >= DATEADD(month, -3, CURRENT_DATE)
')
TO 's3://mybucket/customer_orders/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT PARQUET
PARTITION BY (order_date)
MAXFILESIZE 500 MB
MANIFEST VERBOSE;