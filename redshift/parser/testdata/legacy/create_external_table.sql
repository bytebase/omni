-- CREATE EXTERNAL TABLE test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_EXTERNAL_TABLE.html

-- Basic external table with PARQUET format
CREATE EXTERNAL TABLE spectrum.sales (
    salesid INTEGER,
    listid INTEGER,
    sellerid INTEGER,
    buyerid INTEGER,
    eventid INTEGER,
    dateid SMALLINT,
    qtysold SMALLINT,
    pricepaid DECIMAL(8,2),
    commission DECIMAL(8,2),
    saletime TIMESTAMP
)
STORED AS PARQUET
LOCATION 's3://awssampledbuswest2/tickit/sales/';

-- External table with IF NOT EXISTS
CREATE EXTERNAL TABLE IF NOT EXISTS spectrum.events (
    eventid INTEGER,
    venueid SMALLINT,
    catid SMALLINT,
    dateid SMALLINT,
    eventname VARCHAR(200),
    starttime TIMESTAMP
)
STORED AS TEXTFILE
LOCATION 's3://mybucket/event_data/';

-- External table with partitioning
CREATE EXTERNAL TABLE spectrum.sales_part (
    salesid INTEGER,
    listid INTEGER,
    sellerid INTEGER,
    buyerid INTEGER,
    eventid INTEGER,
    qtysold SMALLINT,
    pricepaid DECIMAL(8,2),
    commission DECIMAL(8,2),
    saletime TIMESTAMP
)
PARTITIONED BY (
    saledate DATE,
    region VARCHAR(10)
)
STORED AS PARQUET
LOCATION 's3://mybucket/sales_partitioned/';

-- External table with ROW FORMAT DELIMITED
CREATE EXTERNAL TABLE spectrum.csv_table (
    id INT,
    name VARCHAR(100),
    age INT,
    city VARCHAR(50)
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY ','
LINES TERMINATED BY '\n'
NULL DEFINED AS 'NULL'
STORED AS TEXTFILE
LOCATION 's3://mybucket/csv_data/';

-- External table with custom SERDE
CREATE EXTERNAL TABLE spectrum.json_table (
    docid BIGINT,
    username VARCHAR(50),
    posted TIMESTAMP,
    message VARCHAR(500)
)
ROW FORMAT SERDE 'org.apache.hive.serde2.lazy.LazySimpleSerDe'
WITH SERDEPROPERTIES (
    'field.delim' = '\t',
    'serialization.null.format' = ''
)
STORED AS TEXTFILE
LOCATION 's3://mybucket/json_data/';

-- External table with ORC format and table properties
CREATE EXTERNAL TABLE spectrum.orc_table (
    col1 INTEGER,
    col2 VARCHAR(100),
    col3 BOOLEAN,
    col4 DOUBLE PRECISION
)
STORED AS ORC
LOCATION 's3://mybucket/orc_data/'
TABLE PROPERTIES (
    'compression_type' = 'snappy',
    'orc.compress' = 'SNAPPY'
);

-- External table with AVRO format
CREATE EXTERNAL TABLE spectrum.avro_table (
    id BIGINT,
    email VARCHAR(255),
    ts TIMESTAMP
)
STORED AS AVRO
LOCATION 's3://mybucket/avro_data/'
TABLE PROPERTIES (
    'avro.schema.url' = 's3://mybucket/schemas/user.avsc'
);

-- External table with RCFILE format
CREATE EXTERNAL TABLE spectrum.rcfile_table (
    key STRING,
    value STRING
)
STORED AS RCFILE
LOCATION 's3://mybucket/rcfile_data/';

-- External table with SEQUENCEFILE format
CREATE EXTERNAL TABLE spectrum.seqfile_table (
    user_id BIGINT,
    activity VARCHAR(100),
    timestamp TIMESTAMP
)
STORED AS SEQUENCEFILE
LOCATION 's3://mybucket/sequence_data/';

-- External table with custom input/output formats
CREATE EXTERNAL TABLE spectrum.custom_format_table (
    data VARCHAR(MAX)
)
STORED AS INPUTFORMAT 'com.example.CustomInputFormat'
OUTPUTFORMAT 'com.example.CustomOutputFormat'
LOCATION 's3://mybucket/custom_data/';

-- External table with comprehensive table properties
CREATE EXTERNAL TABLE spectrum.data_table (
    id INTEGER,
    data VARCHAR(500),
    created_date DATE
)
PARTITIONED BY (year INT, month INT)
STORED AS PARQUET
LOCATION 's3://mybucket/data/'
TABLE PROPERTIES (
    'compression_type' = 'gzip',
    'data_cleansing_enabled' = 'true',
    'invalid_char_handling' = 'REPLACE',
    'replacement_char' = '?',
    'numeric_overflow_handling' = 'SET_TO_NULL',
    'write.parallel' = 'on',
    'write.maxfilesize.mb' = '100'
);

-- External table using manifest file
CREATE EXTERNAL TABLE spectrum.sales_manifest (
    salesid INTEGER,
    listid INTEGER,
    sellerid INTEGER,
    buyerid INTEGER
)
STORED AS PARQUET
LOCATION 's3://mybucket/sales_manifest.json';

-- External table with complex data types
CREATE EXTERNAL TABLE spectrum.complex_types (
    id INTEGER,
    name VARCHAR(100),
    address SUPER,
    phone_numbers SUPER,
    metadata SUPER
)
STORED AS PARQUET
LOCATION 's3://mybucket/complex_data/';

-- External table with escaped delimiters
CREATE EXTERNAL TABLE spectrum.escaped_data (
    col1 VARCHAR(100),
    col2 VARCHAR(100),
    col3 VARCHAR(100)
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '|' ESCAPED BY '\\'
LINES TERMINATED BY '\n'
STORED AS TEXTFILE
LOCATION 's3://mybucket/escaped_data/';

-- External table with all ROW FORMAT DELIMITED options
CREATE EXTERNAL TABLE spectrum.full_delimited (
    id INT,
    data VARCHAR(1000),
    tags SUPER,
    properties SUPER
)
ROW FORMAT DELIMITED
FIELDS TERMINATED BY '\001'
COLLECTION ITEMS TERMINATED BY '\002'
MAP KEYS TERMINATED BY '\003'
LINES TERMINATED BY '\n'
NULL DEFINED AS '\\N'
STORED AS TEXTFILE
LOCATION 's3://mybucket/full_delimited_data/';