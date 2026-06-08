-- CREATE TABLE test cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_CREATE_TABLE.html

-- ====================================
-- Basic table creation
-- ====================================

-- Simple table with basic data types
CREATE TABLE users (
    user_id INTEGER,
    username VARCHAR(50),
    email VARCHAR(255),
    created_date DATE
);

-- Table with IF NOT EXISTS
CREATE TABLE IF NOT EXISTS products (
    product_id INTEGER,
    product_name VARCHAR(100),
    price DECIMAL(10,2),
    category VARCHAR(50)
);

-- Table with schema qualification
CREATE TABLE public.orders (
    order_id BIGINT,
    user_id INTEGER,
    order_date TIMESTAMP,
    total_amount DECIMAL(12,2)
);

-- ====================================
-- Data types testing
-- ====================================

-- All numeric types
CREATE TABLE numeric_types (
    col_smallint SMALLINT,
    col_integer INTEGER,
    col_int INT,
    col_int2 INT2,
    col_int4 INT4,
    col_int8 INT8,
    col_bigint BIGINT,
    col_decimal DECIMAL(10,2),
    col_numeric NUMERIC(15,4),
    col_real REAL,
    col_float4 FLOAT4,
    col_double DOUBLE PRECISION,
    col_float8 FLOAT8,
    col_float FLOAT
);

-- Boolean and character types
CREATE TABLE char_types (
    col_boolean BOOLEAN,
    col_bool BOOL,
    col_char CHAR(10),
    col_character CHARACTER(20),
    col_nchar NCHAR(15),
    col_bpchar BPCHAR(25),
    col_varchar VARCHAR(255),
    col_character_varying CHARACTER VARYING(100),
    col_nvarchar NVARCHAR(200),
    col_text TEXT,
    col_string STRING
);

-- Binary and special types
CREATE TABLE special_types (
    col_varbyte VARBYTE(1024),
    col_varbinary VARBINARY(512),
    col_binary_varying BINARY VARYING(256),
    col_date DATE,
    col_timestamp TIMESTAMP,
    col_timestamp_tz TIMESTAMP WITH TIME ZONE,
    col_timestamptz TIMESTAMPTZ,
    col_time TIME,
    col_time_tz TIME WITH TIME ZONE,
    col_timetz TIMETZ,
    col_geometry GEOMETRY,
    col_geography GEOGRAPHY,
    col_hllsketch HLLSKETCH,
    col_super SUPER
);

-- ====================================
-- Column constraints
-- ====================================

-- NOT NULL and DEFAULT constraints
CREATE TABLE customer_info (
    customer_id INTEGER NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);

-- IDENTITY columns
CREATE TABLE products_with_identity (
    product_id INTEGER IDENTITY(1,1),
    product_name VARCHAR(100) NOT NULL,
    price DECIMAL(10,2) NOT NULL
);

-- GENERATED IDENTITY columns
CREATE TABLE orders_with_generated_identity (
    order_id BIGINT GENERATED ALWAYS AS IDENTITY,
    customer_id INTEGER NOT NULL,
    order_date DATE DEFAULT CURRENT_DATE
);

-- ====================================
-- Primary and Foreign Keys
-- ====================================

-- Primary key constraint
CREATE TABLE categories (
    category_id INTEGER PRIMARY KEY,
    category_name VARCHAR(100) NOT NULL,
    description TEXT
);

-- Foreign key constraint
CREATE TABLE product_catalog (
    product_id INTEGER PRIMARY KEY,
    product_name VARCHAR(100) NOT NULL,
    category_id INTEGER REFERENCES categories(category_id),
    price DECIMAL(10,2) NOT NULL
);

-- Multiple foreign key constraints with referential actions
CREATE TABLE order_items (
    order_item_id INTEGER PRIMARY KEY,
    order_id INTEGER REFERENCES orders(order_id),
    product_id INTEGER REFERENCES product_catalog(product_id),
    quantity INTEGER NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL
);

-- ====================================
-- Redshift-specific: DISTKEY and SORTKEY
-- ====================================

-- Single column DISTKEY
CREATE TABLE user_activity (
    user_id INTEGER DISTKEY,
    activity_date DATE SORTKEY,
    page_views INTEGER,
    session_duration INTERVAL
);

-- Table-level DISTKEY
CREATE TABLE sales_summary (
    sale_id INTEGER,
    customer_id INTEGER,
    sale_date DATE,
    amount DECIMAL(10,2)
)
DISTKEY(customer_id);

-- Compound SORTKEY
CREATE TABLE time_series_data (
    timestamp TIMESTAMP,
    sensor_id INTEGER,
    temperature DECIMAL(5,2),
    humidity DECIMAL(5,2)
)
SORTKEY(timestamp, sensor_id);

-- Interleaved SORTKEY
CREATE TABLE customer_events (
    customer_id INTEGER,
    event_date DATE,
    event_type VARCHAR(50),
    event_data TEXT
)
INTERLEAVED SORTKEY(customer_id, event_date);

-- ====================================
-- Redshift-specific: DISTSTYLE
-- ====================================

-- DISTSTYLE ALL
CREATE TABLE lookup_table (
    lookup_id INTEGER PRIMARY KEY,
    lookup_value VARCHAR(100)
)
DISTSTYLE ALL;

-- DISTSTYLE EVEN
CREATE TABLE log_data (
    log_id BIGINT,
    timestamp TIMESTAMP,
    message TEXT
)
DISTSTYLE EVEN;

-- DISTSTYLE KEY with explicit DISTKEY
CREATE TABLE customer_orders (
    order_id BIGINT,
    customer_id INTEGER,
    order_date DATE,
    total_amount DECIMAL(12,2)
)
DISTSTYLE KEY
DISTKEY(customer_id);

-- DISTSTYLE AUTO
CREATE TABLE smart_distribution (
    id INTEGER,
    data VARCHAR(1000),
    created_at TIMESTAMP
)
DISTSTYLE AUTO;

-- ====================================
-- Redshift-specific: ENCODE
-- ====================================

-- Column-level encoding
CREATE TABLE encoded_data (
    id INTEGER ENCODE RAW,
    name VARCHAR(100) ENCODE LZO,
    description TEXT ENCODE ZSTD,
    amount DECIMAL(10,2) ENCODE AZ64,
    category_id INTEGER ENCODE DELTA,
    created_date DATE ENCODE DELTA32K
);

-- Table-level auto encoding
CREATE TABLE auto_encoded (
    id INTEGER,
    data VARCHAR(1000),
    timestamp TIMESTAMP
)
ENCODE AUTO;

-- Various encoding types
CREATE TABLE encoding_examples (
    col_raw INTEGER ENCODE RAW,
    col_az64 BIGINT ENCODE AZ64,
    col_bytedict VARCHAR(50) ENCODE BYTEDICT,
    col_delta INTEGER ENCODE DELTA,
    col_delta32k INTEGER ENCODE DELTA32K,
    col_lzo TEXT ENCODE LZO,
    col_mostly8 INTEGER ENCODE MOSTLY8,
    col_mostly16 INTEGER ENCODE MOSTLY16,
    col_mostly32 INTEGER ENCODE MOSTLY32,
    col_runlength VARCHAR(10) ENCODE RUNLENGTH,
    col_text255 VARCHAR(255) ENCODE TEXT255,
    col_text32k VARCHAR(32000) ENCODE TEXT32K,
    col_zstd TEXT ENCODE ZSTD
);

-- ====================================
-- Redshift-specific: BACKUP
-- ====================================

-- Table with backup enabled
CREATE TABLE important_data (
    id INTEGER PRIMARY KEY,
    critical_info TEXT NOT NULL
)
BACKUP YES;

-- Table with backup disabled
CREATE TABLE temporary_data (
    id INTEGER,
    temp_value VARCHAR(100)
)
BACKUP NO;

-- ====================================
-- Comprehensive table with all features
-- ====================================

-- Complex table combining multiple Redshift features
CREATE TABLE comprehensive_table (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    user_id INTEGER NOT NULL DISTKEY,
    transaction_date DATE NOT NULL SORTKEY,
    transaction_type VARCHAR(20) NOT NULL ENCODE BYTEDICT,
    amount DECIMAL(15,2) NOT NULL ENCODE AZ64,
    description TEXT ENCODE LZO,
    region VARCHAR(10) DEFAULT 'US' ENCODE RUNLENGTH,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    CONSTRAINT pk_comprehensive PRIMARY KEY (id),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id),
    CONSTRAINT chk_transaction_type
)
BACKUP YES
DISTSTYLE KEY;

-- ====================================
-- Table constraints (table-level)
-- ====================================

-- Table with multiple table-level constraints
CREATE TABLE order_management (
    order_id INTEGER,
    customer_id INTEGER,
    order_date DATE,
    ship_date DATE,
    order_status VARCHAR(20),
    total_amount DECIMAL(12,2),
    PRIMARY KEY (order_id),
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE RESTRICT,
    UNIQUE (order_id, customer_id)
);

-- ====================================
-- Temporary tables
-- ====================================

-- Temporary table
CREATE TEMPORARY TABLE temp_calculations (
    calc_id INTEGER,
    input_value DECIMAL(10,2),
    result_value DECIMAL(15,4)
);

-- Local temporary table
CREATE LOCAL TEMPORARY TABLE local_temp_data (
    session_id VARCHAR(50),
    data_point DECIMAL(10,2),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Temp table (short form)
CREATE TEMP TABLE quick_temp (
    id INTEGER,
    value VARCHAR(100)
);

-- ====================================
-- CREATE TABLE AS (CTAS)
-- ====================================

-- Simple CTAS
CREATE TABLE customer_summary AS
SELECT customer_id, COUNT(*) as order_count, SUM(total_amount) as total_spent
FROM orders
GROUP BY customer_id;

-- CTAS with table options
CREATE TABLE sales_analysis
DISTSTYLE KEY
DISTKEY(region)
SORTKEY(sale_date)
AS
SELECT region, sale_date, SUM(amount) as daily_sales
FROM sales_data
GROUP BY region, sale_date;

CREATE TABLE monthly_aggregates AS
SELECT
    EXTRACT(YEAR FROM order_date) as order_year,
    EXTRACT(MONTH FROM order_date) as order_month,
    COUNT(*) as order_count,
    AVG(total_amount) as avg_order_value
FROM orders
GROUP BY EXTRACT(YEAR FROM order_date), EXTRACT(MONTH FROM order_date);

-- ====================================
-- Complex data type specifications
-- ====================================

-- Precision and scale specifications
CREATE TABLE precise_numbers (
    tiny_decimal DECIMAL(3,2),
    medium_decimal DECIMAL(10,4),
    large_decimal DECIMAL(38,10),
    small_numeric NUMERIC(5,2),
    large_numeric NUMERIC(38,18)
);

-- Length specifications for character types
CREATE TABLE string_lengths (
    short_char CHAR(5),
    medium_varchar VARCHAR(100),
    long_varchar VARCHAR(65535),
    tiny_nchar NCHAR(10),
    medium_nvarchar NVARCHAR(500)
);

-- ====================================
-- Edge cases and boundary testing
-- ====================================

-- Maximum length identifiers and minimal table
CREATE TABLE a (a INT);

-- Mixed case and quoted identifiers
CREATE TABLE "MixedCaseTable" (
    "ColumnWithSpaces And Special-Chars" VARCHAR(100),
    "UPPERCASE_COLUMN" INTEGER,
    "lowercase_column" DECIMAL(10,2)
);

-- Table with all constraint types
CREATE TABLE constraint_showcase (
    id INTEGER IDENTITY(1,1) PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    parent_id INTEGER REFERENCES constraint_showcase(id),
    value DECIMAL(10,2) DEFAULT 0,
    status CHAR(1) DEFAULT 'A',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Table with complex referential integrity
CREATE TABLE parent_table (
    parent_id INTEGER PRIMARY KEY,
    parent_name VARCHAR(100) NOT NULL
);

CREATE TABLE child_table (
    child_id INTEGER PRIMARY KEY,
    parent_id INTEGER NOT NULL,
    child_name VARCHAR(100) NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES parent_table(parent_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- ====================================
-- Real-world examples
-- ====================================

-- E-commerce product catalog
CREATE TABLE product_inventory (
    sku VARCHAR(50) PRIMARY KEY,
    product_name VARCHAR(200) NOT NULL,
    brand VARCHAR(100) ENCODE BYTEDICT,
    category VARCHAR(50) ENCODE BYTEDICT,
    price DECIMAL(10,2) NOT NULL,
    cost DECIMAL(10,2),
    quantity_on_hand INTEGER DEFAULT 0,
    reorder_level INTEGER DEFAULT 10,
    supplier_id INTEGER,
    weight_kg DECIMAL(8,3),
    dimensions_cm VARCHAR(50),
    is_active BOOLEAN DEFAULT TRUE,
    created_date DATE DEFAULT CURRENT_DATE,
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)
DISTSTYLE KEY
DISTKEY(category)
SORTKEY(brand, product_name);

-- Customer relationship management
CREATE TABLE customer_profiles (
    customer_id BIGINT GENERATED ALWAYS AS IDENTITY,
    customer_type VARCHAR(20) DEFAULT 'INDIVIDUAL' ENCODE BYTEDICT,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    company_name VARCHAR(200),
    email VARCHAR(255) UNIQUE,
    phone VARCHAR(20),
    address_line1 VARCHAR(200),
    address_line2 VARCHAR(200),
    city VARCHAR(100),
    state_province VARCHAR(100),
    postal_code VARCHAR(20),
    country VARCHAR(2) DEFAULT 'US',
    registration_date DATE DEFAULT CURRENT_DATE,
    last_login TIMESTAMP,
    total_orders INTEGER DEFAULT 0,
    lifetime_value DECIMAL(12,2) DEFAULT 0,
    credit_limit DECIMAL(10,2),
    is_vip BOOLEAN DEFAULT FALSE,
    notes TEXT,
    PRIMARY KEY (customer_id)
)
BACKUP YES
DISTSTYLE AUTO;


-- Financial transactions table
CREATE TABLE financial_transactions (
    transaction_id BIGINT IDENTITY(1000000000, 1),
    account_id BIGINT NOT NULL DISTKEY,
    transaction_date TIMESTAMP NOT NULL SORTKEY,
    transaction_type VARCHAR(20) NOT NULL ENCODE BYTEDICT,
    amount DECIMAL(15,2) NOT NULL,
    currency CHAR(3) DEFAULT 'USD' ENCODE BYTEDICT,
    description VARCHAR(500) ENCODE LZO,
    reference_number VARCHAR(50),
    balance_after DECIMAL(15,2),
    fee_amount DECIMAL(8,2) DEFAULT 0,
    exchange_rate DECIMAL(10,6),
    processed_by VARCHAR(100),
    processed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(20) DEFAULT 'COMPLETED' ENCODE BYTEDICT,
    PRIMARY KEY (transaction_id)
)
BACKUP YES
DISTSTYLE KEY;
