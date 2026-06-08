-- ALTER TABLE Test Cases for Redshift
-- Reference: https://docs.aws.amazon.com/redshift/latest/dg/r_ALTER_TABLE.html

-- ===========================
-- 1. ADD COLUMN Operations
-- ===========================

-- Basic ADD COLUMN
ALTER TABLE sales ADD COLUMN region VARCHAR(50);

-- ADD COLUMN with constraints
ALTER TABLE employees ADD COLUMN email VARCHAR(255) NOT NULL;

-- ADD COLUMN with encoding
ALTER TABLE events ADD COLUMN event_data VARCHAR(MAX) ENCODE ZSTD;

-- ADD COLUMN with default value
ALTER TABLE orders ADD COLUMN status VARCHAR(20) DEFAULT 'pending';

-- ADD COLUMN without COLUMN keyword
ALTER TABLE products ADD price DECIMAL(10,2);

-- ADD COLUMN with multiple constraints
ALTER TABLE users ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL;

-- ===========================
-- 2. DROP COLUMN Operations
-- ===========================

-- Basic DROP COLUMN
ALTER TABLE sales DROP COLUMN temp_data;

-- DROP COLUMN with CASCADE
ALTER TABLE employees DROP COLUMN department_id CASCADE;

-- DROP COLUMN with RESTRICT
ALTER TABLE products DROP COLUMN old_price RESTRICT;

-- DROP COLUMN without COLUMN keyword
ALTER TABLE orders DROP shipping_notes;

-- ===========================
-- 3. ALTER COLUMN Operations
-- ===========================

-- ALTER COLUMN TYPE
ALTER TABLE customers ALTER COLUMN phone TYPE VARCHAR(20);

-- ALTER COLUMN SET ENCODE
ALTER TABLE events ALTER COLUMN event_data ENCODE LZO;

-- ===========================
-- 4. RENAME Operations
-- ===========================

-- RENAME TABLE
ALTER TABLE old_customers RENAME TO customers_archive;

-- RENAME COLUMN
ALTER TABLE products RENAME COLUMN product_name TO name;

-- RENAME TABLE with IF EXISTS
ALTER TABLE IF EXISTS temp_table RENAME TO permanent_table;

-- ===========================
-- 5. CONSTRAINT Operations
-- ===========================

-- ADD PRIMARY KEY constraint
ALTER TABLE employees ADD CONSTRAINT pk_employees PRIMARY KEY (employee_id);

-- ADD UNIQUE constraint
ALTER TABLE users ADD CONSTRAINT uk_users_email UNIQUE (email);

-- ADD FOREIGN KEY constraint
ALTER TABLE orders ADD CONSTRAINT fk_orders_customer
    FOREIGN KEY (customer_id) REFERENCES customers(customer_id);

-- ADD FOREIGN KEY with referential actions
ALTER TABLE order_items ADD CONSTRAINT fk_items_order
    FOREIGN KEY (order_id) REFERENCES orders(order_id);

-- ADD unnamed constraint
ALTER TABLE inventory ADD UNIQUE (product_id, warehouse_id);

-- DROP CONSTRAINT
ALTER TABLE orders DROP CONSTRAINT fk_orders_customer;

-- DROP CONSTRAINT with CASCADE
ALTER TABLE products DROP CONSTRAINT pk_products CASCADE;

-- ===========================
-- 6. Distribution Key Operations
-- ===========================

-- ALTER DISTKEY
ALTER TABLE sales ALTER DISTKEY customer_id;

-- ALTER DISTSTYLE KEY with DISTKEY
ALTER TABLE events ALTER DISTSTYLE KEY DISTKEY event_id;

-- ALTER DISTSTYLE to ALL
ALTER TABLE small_lookup ALTER DISTSTYLE ALL;

-- ALTER DISTSTYLE to EVEN
ALTER TABLE large_table ALTER DISTSTYLE EVEN;

-- ALTER DISTSTYLE to AUTO
ALTER TABLE transactions ALTER DISTSTYLE AUTO;

-- ===========================
-- 7. Sort Key Operations
-- ===========================

-- ALTER SORTKEY with single column
ALTER TABLE sales ALTER SORTKEY (sale_date);

-- ALTER COMPOUND SORTKEY with multiple columns
ALTER TABLE events ALTER COMPOUND SORTKEY (event_date, event_type);

-- ALTER INTERLEAVED SORTKEY
ALTER TABLE customer_data ALTER SORTKEY (customer_id, order_date);

-- ALTER SORTKEY AUTO
ALTER TABLE analytics_table ALTER SORTKEY AUTO;

-- ALTER SORTKEY NONE
ALTER TABLE temp_table ALTER SORTKEY NONE;

-- ===========================
-- 8. Encoding Operations
-- ===========================

-- ALTER ENCODE AUTO
ALTER TABLE large_table ALTER ENCODE AUTO;

-- ALTER ENCODE AUTO PRESERVE YES
ALTER TABLE historical_data ALTER ENCODE AUTO;

-- ALTER ENCODE AUTO PRESERVE NO
ALTER TABLE staging_table ALTER ENCODE AUTO;

-- ALTER specific column encoding
ALTER TABLE events ALTER COLUMN event_data ENCODE ZSTD;

-- ===========================
-- 9. Owner Operations
-- ===========================

-- Change table owner
ALTER TABLE sales OWNER TO analytics_user;

-- Change owner with IF EXISTS
ALTER TABLE temp_data OWNER TO etl_user;

-- ===========================
-- 10. Row Level Security Operations
-- ===========================

-- Enable row level security
ALTER TABLE sensitive_data ROW LEVEL SECURITY ON;

-- Disable row level security
ALTER TABLE public_data ROW LEVEL SECURITY OFF;

-- Enable masking for datashares
ALTER TABLE shared_table MASKING ON FOR DATASHARES;
-- Disable masking for datashares
ALTER TABLE unmasked_table MASKING OFF FOR DATASHARES;

-- ===========================
-- 11. External Table Operations
-- ===========================

-- SET TABLE PROPERTIES
ALTER TABLE spectrum.sales SET TABLE PROPERTIES ('numRows'='1000000', 'compressionType'='gzip');

-- APPEND FROM another table
ALTER TABLE sales_history APPEND FROM sales_staging;

-- SET LOCATION for external table
ALTER TABLE spectrum.events SET LOCATION 's3://mybucket/data/events/2024/';

-- SET FILE FORMAT
ALTER TABLE spectrum.logs SET FILE FORMAT PARQUET;

-- ADD PARTITION
ALTER TABLE spectrum.sales ADD PARTITION (year=2024, month=1)
    LOCATION 's3://mybucket/sales/2024/01/';

-- ADD PARTITION IF NOT EXISTS
ALTER TABLE spectrum.events ADD IF NOT EXISTS PARTITION (date='2024-01-01')
    LOCATION 's3://mybucket/events/2024-01-01/';

-- DROP PARTITION
ALTER TABLE spectrum.sales DROP PARTITION (year=2023, month=12);

-- ===========================
-- 12. Complex Multi-Action Operations
-- ===========================

-- Multiple column operations in single statement
ALTER TABLE products
    ADD COLUMN category VARCHAR(50),
    ADD COLUMN subcategory VARCHAR(50),
    ALTER COLUMN price TYPE DECIMAL(12,2),
    DROP COLUMN old_category;

-- Add multiple constraints
ALTER TABLE orders
    ADD CONSTRAINT pk_orders PRIMARY KEY (order_id),
    ADD CONSTRAINT fk_customer FOREIGN KEY (customer_id) REFERENCES customers(id);

-- Comprehensive table modification
ALTER TABLE sales_fact
    ALTER DISTKEY customer_id,
    ALTER SORTKEY (sale_date, product_id),
    ADD COLUMN region_id INTEGER ENCODE AZ64,
    DROP CONSTRAINT old_constraint CASCADE;

-- ===========================
-- 13. Edge Cases and Special Syntax
-- ===========================

-- Table with quoted identifier
ALTER TABLE "MixedCase_Table" ADD COLUMN "New Column" VARCHAR(100);

-- Column names requiring quotes
ALTER TABLE sales RENAME COLUMN "old-column" TO "new_column";

-- Using non-reserved keywords as identifiers
ALTER TABLE data ADD COLUMN data VARCHAR(100);

-- Complex data type modifications
ALTER TABLE measurements ALTER COLUMN reading TYPE NUMERIC(15,6);

-- Multiple encoding types
ALTER TABLE performance_data
    ALTER COLUMN metric1 ENCODE MOSTLY32,
    ALTER COLUMN metric2 ENCODE DELTA32K;