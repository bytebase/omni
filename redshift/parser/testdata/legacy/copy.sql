-- Redshift COPY command test cases
-- These examples focus on Redshift-specific features that differ from PostgreSQL COPY

-- Basic S3 COPY with IAM Role
COPY customer
FROM 's3://mybucket/customer/customer.tbl'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole';

-- S3 COPY with default IAM Role
COPY sales
FROM 's3://mybucket/data/sales/'
IAM_ROLE default;

-- S3 COPY with multiple options
COPY listing
FROM 's3://mybucket/data/listing/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '|'
IGNOREHEADER 1
DATEFORMAT 'auto'
TIMEFORMAT 'auto'
MAXERROR 10
ACCEPTINVCHARS
GZIP;

-- S3 COPY with manifest file
COPY customer
FROM 's3://mybucket/cust.manifest'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
MANIFEST
DELIMITER '|'
GZIP
DATEFORMAT 'auto'
MAXERROR 10;

-- S3 COPY with encryption
COPY venue
FROM 's3://mybucket/encrypted/venue/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ENCRYPTED
DELIMITER '|'
GZIP
REMOVEQUOTES
BLANKSASNULL
EMPTYASNULL
EXPLICIT_IDS;

-- S3 COPY with region specification
COPY sales
FROM 's3://mybucket-eu/sales/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
REGION 'eu-west-1'
DELIMITER '\t'
BZIP2;

-- DynamoDB COPY
COPY favoritemovies
FROM 'dynamodb://Movies'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
READRATIO 50;

-- DynamoDB COPY with column mapping
COPY myredshifttable (col1, col2, col3)
FROM 'dynamodb://ProductCatalog'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
READRATIO 25
DATEFORMAT 'auto'
TIMEFORMAT 'auto'
COMPUPDATE OFF
STATUPDATE ON;

-- EMR COPY
COPY sales
FROM 'emr://j-SAMPLE2B500FC/myoutput/part-*'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '\t'
LZOP;

-- CSV format COPY
COPY category
FROM 's3://mybucket/data/category_csv.txt'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT CSV
QUOTE '"'
DELIMITER ','
IGNOREHEADER 1;

-- CSV format with AS keyword
COPY venue
FROM 's3://mybucket/venue.csv'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT AS CSV
IGNOREHEADER AS 1
DATEFORMAT AS 'MM/DD/YYYY'
REGION AS 'us-west-2';

-- JSON format with auto
COPY category
FROM 's3://mybucket/category_object_auto.json'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
JSON 'auto'
TRUNCATECOLUMNS;

-- JSON format with JSONPaths file
COPY venue
FROM 's3://mybucket/venue_pipe.json'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
JSON 's3://mybucket/venue_jsonpath.json'
ACCEPTINVCHARS
DATEFORMAT 'auto'
MAXERROR 10
REGION 'us-west-2';

-- Avro format
COPY category
FROM 's3://mybucket/category_object.avro'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT AS AVRO 'auto';

-- Avro with schema file
COPY product
FROM 's3://mybucket/data/product.avro'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
AVRO 's3://mybucket/schema/product.avsc'
ACCEPTANYDATE
BLANKSASNULL;

-- Parquet format
COPY listing
FROM 's3://mybucket/data/listings/parquet/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT AS PARQUET;

-- ORC format
COPY sales
FROM 's3://mybucket/data/sales/orc/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT AS ORC
DATEFORMAT 'auto'
TIMEFORMAT 'auto';

-- Fixed width format
COPY venue_fw
FROM 's3://mybucket/data/venue_fw.txt'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FIXEDWIDTH 'venueid:3,venuename:25,venuecity:12,venuestate:2,venueseats:6'
ACCEPTINVCHARS ' '
BLANKSASNULL;

-- Multiple compression formats
COPY sales_gzip FROM 's3://mybucket/sales.gz'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' GZIP;

COPY sales_bzip2 FROM 's3://mybucket/sales.bz2'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' BZIP2;

COPY sales_lzop FROM 's3://mybucket/sales.lzo'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' LZOP;

COPY sales_zstd FROM 's3://mybucket/sales.zst'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole' ZSTD;

-- Data conversion options
COPY event
FROM 's3://mybucket/data/allevents_pipe.txt'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '|'
TIMEFORMAT 'YYYY-MM-DD HH:MI:SS'
IGNOREHEADER 1
ACCEPTINVCHARS '?'
NULL AS '\000'
ESCAPE
ROUNDEC
TRIMBLANKS
TRUNCATECOLUMNS
IGNOREBLANKLINES
ACCEPTANYDATE
EMPTYASNULL
BLANKSASNULL
FILLRECORD
REMOVEQUOTES;

-- Load operation options
COPY sales
FROM 's3://mybucket/sales/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '\t'
COMPUPDATE PRESET
STATUPDATE ON
MAXERROR 100
NOLOAD;

-- Error handling options
COPY venue
FROM 's3://mybucket/venue_errors/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER ','
IGNOREALLERRORS;

-- Time format options
COPY timestamp_table
FROM 's3://mybucket/timestamps/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
TIMEFORMAT 'auto';

COPY epoch_table
FROM 's3://mybucket/epochs/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
TIMEFORMAT 'epochsecs';

COPY epoch_milli_table
FROM 's3://mybucket/epochs_milli/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
TIMEFORMAT 'epochmillisecs';

-- Encoding options
COPY international_sales
FROM 's3://mybucket/international/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
ENCODING UTF16
DELIMITER '\t'
ACCEPTINVCHARS;

-- Access key authorization (not recommended but supported)
COPY sales
FROM 's3://mybucket/sales/'
ACCESS_KEY_ID 'AKIAIOSFODNN7EXAMPLE'
SECRET_ACCESS_KEY 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'
DELIMITER '|';

-- Session token for temporary credentials
COPY venue
FROM 's3://mybucket/venue/'
ACCESS_KEY_ID 'AKIAIOSFODNN7EXAMPLE'
SECRET_ACCESS_KEY 'wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY'
SESSION_TOKEN 'AQoDYXdzEPT//////////wEXAMPLEtc764bNrC9SAPBSM22wDOk4x4HIZ8j4FZTwdQWLWsKWHGBuFqwAeMicRXmxfpSPfIeoIYRqTflfKD8YUuwthAx7mSEI/qkPpKPi/kMcGd'
DELIMITER ',';

-- CREDENTIALS clause
COPY customer
FROM 's3://mybucket/customer/'
CREDENTIALS 'aws_iam_role=arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '|';

-- Complex multi-option example
COPY sales_fact (
    sale_id,
    product_id,
    customer_id,
    sale_date,
    quantity,
    amount
)
FROM 's3://mybucket/sales/2023/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
FORMAT AS CSV
QUOTE AS '"'
DELIMITER ','
IGNOREHEADER 1
DATEFORMAT 'YYYY-MM-DD'
TIMEFORMAT 'YYYY-MM-DD HH24:MI:SS'
REGION 'us-east-1'
GZIP
MAXERROR 50
ACCEPTINVCHARS '?'
BLANKSASNULL
EMPTYASNULL
ESCAPE
EXPLICIT_IDS
ROUNDEC
TRIMBLANKS
COMPUPDATE OFF
STATUPDATE ON;

-- COPY with COMPROWS option
COPY large_table
FROM 's3://mybucket/large_table/'
IAM_ROLE 'arn:aws:iam::123456789012:role/MyRedshiftRole'
DELIMITER '|'
COMPROWS 1000000
COMPUPDATE PRESET;